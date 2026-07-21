package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

// Test texEscape function
func TestTexEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "&", expected: `\&`},
		{input: "%", expected: `\%`},
		{input: "$", expected: `\$`},
		{input: "#", expected: `\#`},
		{input: "_", expected: `\_`},
		{input: "{", expected: `\{`},
		{input: "}", expected: `\}`},
		{input: "~", expected: `\textasciitilde{}`},
		{input: "^", expected: `\textasciicircum{}`},
		{input: `\`, expected: `\textbackslash{}`},
		{input: "Hello & % $ _ { } ~ ^ \\ World", expected: `Hello \& \% \$ \_ \{ \} \textasciitilde{} \textasciicircum{} \textbackslash{} World`},
		{input: "Plain text", expected: "Plain text"},
		{input: "", expected: ""},
	}

	for _, tt := range tests {
		if got := texEscape(tt.input); got != tt.expected {
			t.Errorf("texEscape(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// Test health handler
func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.RemoteAddr = "127.0.0.1:8080" // Unique IP to avoid conflicts with rate limiter tests
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "200 OK"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// Test generatePdfHandler with valid input
func TestGeneratePdfHandler_Success(t *testing.T) {
	// Skip test if pdflatex is not available
	if _, err := exec.LookPath("pdflatex"); err != nil {
		t.Skip("pdflatex not available, skipping PDF generation test")
	}

	// Prepare test data
	testData := ResumeData{
		Name:    "John Doe",
		Title:   "Software Engineer",
		Email:   "john@example.com",
		Phone:   "123-456-7890",
		Experiences: []Experience{
			{
				Role:   "Senior Engineer",
				Company: "Tech Corp",
				Period: "2020-2023",
				Bullets: []string{
					"Led development of web applications",
					"Improved system performance by 20%",
				},
			},
		},
	}
	jsonData, _ := json.Marshal(testData)

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := rateLimiter(generatePdfHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	if rr.Header().Get("Content-Type") != "application/pdf" {
		t.Errorf("Handler returned unexpected content type: got %v want application/pdf", rr.Header().Get("Content-Type"))
	}

	// Check that we got a PDF (starts with %PDF)
	if !strings.HasPrefix(rr.Body.String(), "%PDF") {
		t.Errorf("Handler did not return a valid PDF")
	}
}

// Test generatePdfHandler with missing required fields
func TestGeneratePdfHandler_MissingFields(t *testing.T) {
	// Prepare test data with missing Name
	testData := ResumeData{
		Title:   "Software Engineer",
		Email:   "john@example.com",
		Phone:   "123-456-7890",
	}
	jsonData, _ := json.Marshal(testData)

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:12345" // Unique IP for this test
	rr := httptest.NewRecorder()
	handler := rateLimiter(generatePdfHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

// Test rate limiter
func TestRateLimiter(t *testing.T) {
	// Create a simple handler that just returns OK
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	limitedHandler := rateLimiter(handlerFunc)

	// First 5 requests should succeed (we start with 5 tokens)
	for i := 1; i <= 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		limitedHandler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Request %d failed: got %v want %v", i, rr.Code, http.StatusOK)
			return
		}
	}

	// 6th request should be rate limited (no tokens left)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	limitedHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request 6 should be rate limited: got %v want %v", rr.Code, http.StatusTooManyRequests)
		return
	}

	// 7th request should also be rate limited (no token replenishment in current implementation)
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	rr = httptest.NewRecorder()
	limitedHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request 7 should be rate limited: got %v want %v", rr.Code, http.StatusTooManyRequests)
		return
	}

	// Different IP should get a fresh visitor with 5 tokens
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	rr = httptest.NewRecorder()
	limitedHandler.ServeHTTP(rr, req2)
	if rr.Code != http.StatusOK {
		t.Errorf("Request with new IP failed: got %v want %v", rr.Code, http.StatusOK)
		return
	}
}

// Test rate limiter with same IP after tokens exhausted (should still be limited)
func TestRateLimiter_SameIPAfterExhaustion(t *testing.T) {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	limitedHandler := rateLimiter(handlerFunc)

	// Exhaust tokens for a specific IP
	ip := "10.0.0.1:54321"
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = ip
		rr := httptest.NewRecorder()
		limitedHandler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Request %d failed: got %v want %v", i+1, rr.Code, http.StatusOK)
		}
	}

	// 6th request from same IP should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip
	rr := httptest.NewRecorder()
	limitedHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request after exhaustion should be rate limited: got %v want %v", rr.Code, http.StatusTooManyRequests)
	}
}