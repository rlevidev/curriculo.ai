package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// newRequestID generates a random 8-byte hex ID; falls back to timestamp if rand fails.
func newRequestID() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}

// getRequestID reads request_id from context; returns "" if absent.
func getRequestID(r *http.Request) string {
	if v := r.Context().Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// loggerFor returns baseLogger with request_id attached for log correlation.
func loggerFor(r *http.Request) *slog.Logger {
	return baseLogger.With("request_id", getRequestID(r))
}

// requestIDMiddleware propagates X-Request-ID from the header or generates a new one and injects it into context.
func requestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = newRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		next(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader records the status before forwarding to the underlying ResponseWriter.
func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Write records an implicit 200 when called without a prior WriteHeader.
func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

// loggingMiddleware measures latency, captures status, and logs with level by status code (debug for /health, warn for 4xx, error for 5xx).
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 0}
		next(sr, r)
		if sr.status == 0 {
			sr.status = http.StatusOK
		}
		latency := time.Since(start)
		level := slog.LevelInfo
		if sr.status >= 500 {
			level = slog.LevelError
		} else if sr.status >= 400 {
			level = slog.LevelWarn
		} else if r.URL.Path == "/health" {
			level = slog.LevelDebug
		}
		loggerFor(r).Log(r.Context(), level, "request",
			"op", "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sr.status,
			"latency_ms", latency.Milliseconds(),
		)
	}
}

// corsMiddleware sets Allow-Origin/Methods/Headers and short-circuits preflight OPTIONS requests.
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	allowedOrigin := os.Getenv("FRONTEND_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "https://rlevidev.github.io"
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

var semaphore = make(chan struct{}, 3)

// generatePdfHandler validates payload (50KB limit, name+title required), escapes LaTeX, limits concurrency (semaphore cap 3), and writes the PDF.
func generatePdfHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 50*1024)
	var data ResumeData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		loggerFor(r).Warn("invalid payload", "op", "validation", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		return
	}
	if len(data.Name) == 0 || len(data.Title) == 0 {
		loggerFor(r).Warn("missing required fields", "op", "validation", "status", http.StatusBadRequest)
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}
	data = escapeResumeData(data)
	semaphore <- struct{}{}
	defer func() { <-semaphore }()
	pdfStart := time.Now()
	pdf, err := generatePDF(r.Context(), data)
	if err != nil {
		loggerFor(r).Error("pdf generation failed", "op", "pdf_generate", "status", http.StatusInternalServerError, "latency_ms", time.Since(pdfStart).Milliseconds(), "error", err.Error())
		http.Error(w, "failed to generate PDF", http.StatusInternalServerError)
		return
	}
	loggerFor(r).Info("pdf generated", "op", "pdf_generate", "status", http.StatusOK, "latency_ms", time.Since(pdfStart).Milliseconds(), "pdf_bytes", len(pdf))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="curriculo.pdf"`)
	if _, err := w.Write(pdf); err != nil {
		loggerFor(r).Error("write pdf failed", "op", "pdf_generate", "error", err.Error())
	}
}

// setupServer builds the mux, chains middlewares (requestID → logging → cors → rateLimit), and configures timeouts.
func setupServer(port string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", requestIDMiddleware(loggingMiddleware(corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
	}))))
	mux.HandleFunc("/generate-pdf", requestIDMiddleware(loggingMiddleware(corsMiddleware(rateLimiter(generatePdfHandler)))))

	return &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
