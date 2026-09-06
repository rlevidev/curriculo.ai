package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

// Rate Limiter
type Visitor struct {
	lastSeen time.Time
	tokens   int
}

var (
	visitors = make(map[string]*Visitor)
	mu       sync.Mutex
	// stdout JSON para coletor do host/Docker, sem agregador dedicado por enquanto
	baseLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.UTC().Format(time.RFC3339))
				}
			}
			return a
		},
	})).With("service", "curriculo-api")
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func newRequestID() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}

func getRequestID(r *http.Request) string {
	if v := r.Context().Value(requestIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// sem PII do ResumeData nos logs, só metadados
func loggerFor(r *http.Request) *slog.Logger {
	return baseLogger.With("request_id", getRequestID(r))
}

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

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.status = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

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

func init() {
	go cleanupVisitors()
}

func getVisitorIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return ip
}

func cleanupVisitors() {
	for {
		time.Sleep(10 * time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 1*time.Hour {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func rateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getVisitorIP(r.RemoteAddr)

		mu.Lock()
		visitor, ok := visitors[ip]
		if !ok {
			visitor = &Visitor{lastSeen: time.Now(), tokens: 5}
			visitors[ip] = visitor
		}

		if visitor.tokens <= 0 {
			mu.Unlock()
			loggerFor(r).Warn("rate limited", "op", "rate_limit", "status", http.StatusTooManyRequests)
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		visitor.tokens--
		mu.Unlock()
		next(w, r)
	}
}

// Semaphore
var semaphore = make(chan struct{}, 3)

type Education struct {
	Institution string   `json:"institution"`
	Degree      string   `json:"degree"`
	Period      string   `json:"period"`
	Notes       []string `json:"notes"`
}

type Experience struct {
	Company string   `json:"company"`
	Role    string   `json:"role"`
	Period  string   `json:"period"`
	Bullets []string `json:"bullets"`
}

type Project struct {
	Name    string   `json:"name"`
	Link    string   `json:"link"`
	Bullets []string `json:"bullets"`
}

type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
}

type Skills struct {
	Languages    []string `json:"languages"`
	Technologies []string `json:"technologies"`
}

type ResumeData struct {
	Name            string       `json:"name"`
	Title           string       `json:"title"`
	Email           string       `json:"email"`
	Phone           string       `json:"phone"`
	LinkedIn        string       `json:"linkedin"`
	GitHub          string       `json:"github"`
	Location        string       `json:"location"`
	Education       []Education  `json:"education"`
	Experiences     []Experience `json:"experiences"`
	Projects        []Project    `json:"projects"`
	SpokenLanguages []Language   `json:"spokenLanguages"`
	Certifications  []string     `json:"certifications"`
	Skills          Skills       `json:"skills"`
}

func texEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "\\&",
		"%", "\\%",
		"$", "\\$",
		"#", "\\#",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"~", "\\textasciitilde{}",
		"^", "\\textasciicircum{}",
		"\\", "\\textbackslash{}",
	)
	return replacer.Replace(s)
}

func generatePdfHandler(w http.ResponseWriter, r *http.Request) {
	// Validate size
	if r.ContentLength > 50*1024 {
		loggerFor(r).Warn("payload too large", "op", "validation", "status", http.StatusBadRequest)
		http.Error(w, "payload too large", http.StatusBadRequest)
		return
	}

	var data ResumeData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		loggerFor(r).Warn("invalid payload", "op", "validation", "status", http.StatusBadRequest, "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simple validation
	if len(data.Name) == 0 || len(data.Title) == 0 {
		loggerFor(r).Warn("missing required fields", "op", "validation", "status", http.StatusBadRequest)
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	// Semaphore
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	tmpl := template.New("resume").Delims("<[", "]>")
	tmpl, _ = tmpl.Funcs(template.FuncMap{"texEscape": texEscape}).Parse(`
\documentclass{article}
\usepackage{hyperref}
\begin{document}
\section*{<[ .Name ]>}
\subsection*{<[ .Title ]>}
<[ .Email ]> | <[ .Phone ]>

<[ if .Experiences ]>
\section*{Experience}
<[ range .Experiences ]>
\textbf{<[ .Role ]>} @ <[ .Company ]> (<[ .Period ]>)
\begin{itemize}
<[ range .Bullets ]> \item <[ . ]> <[ end ]>
\end{itemize}
<[ end ]>
<[ end ]>

<[ if .SpokenLanguages ]>
\section*{Languages}
<[ range .SpokenLanguages ]><[ .Name ]> (<[ .Proficiency ]>)<[ end ]>
<[ end ]>

<[ if .Certifications ]>
\section*{Certifications}
<[ range .Certifications ]><[ . ]><[ end ]>
<[ end ]>
\end{document}
`)

	var buf bytes.Buffer
	tmpl.Execute(&buf, data)

	dir, err := os.MkdirTemp("", "resume")
	if err != nil {
		loggerFor(r).Error("tempdir failed", "op", "pdf_generate", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	texPath := filepath.Join(dir, "resume.tex")
	if err := os.WriteFile(texPath, buf.Bytes(), 0644); err != nil {
		loggerFor(r).Error("write tex failed", "op", "pdf_generate", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	pdfStart := time.Now()
	cmd := exec.CommandContext(ctx, "pdflatex", "-no-shell-escape", "-output-directory", dir, texPath)
	if err := cmd.Run(); err != nil {
		loggerFor(r).Error("pdflatex failed", "op", "pdf_generate", "status", http.StatusInternalServerError, "latency_ms", time.Since(pdfStart).Milliseconds(), "error", err.Error())
		http.Error(w, "failed to generate PDF", http.StatusInternalServerError)
		return
	}

	pdfPath := filepath.Join(dir, "resume.pdf")
	pdf, err := os.ReadFile(pdfPath)
	if err != nil {
		loggerFor(r).Error("read pdf failed", "op", "pdf_generate", "status", http.StatusInternalServerError, "error", err.Error())
		http.Error(w, "failed to read generated PDF", http.StatusInternalServerError)
		return
	}

	loggerFor(r).Info("pdf generated", "op", "pdf_generate", "status", http.StatusOK, "latency_ms", time.Since(pdfStart).Milliseconds(), "pdf_bytes", len(pdf))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"curriculo.pdf\"")
	w.Write(pdf)
}

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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseLogger.Info("starting", "op", "startup", "port", port)
	http.HandleFunc("/health", requestIDMiddleware(loggingMiddleware(corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
	}))))
	http.HandleFunc("/generate-pdf", requestIDMiddleware(loggingMiddleware(corsMiddleware(rateLimiter(generatePdfHandler)))))
	http.ListenAndServe(":"+port, nil)
}
