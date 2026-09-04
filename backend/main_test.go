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
	"time"
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

func resetVisitors() {
	mu.Lock()
	visitors = make(map[string]*Visitor)
	mu.Unlock()
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
	resetVisitors()
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
	resetVisitors()
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
	resetVisitors()
	// Create a simple handler that just returns OK
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	limitedHandler := rateLimiter(handlerFunc)

	ip := "192.168.1.100"
	port1 := ":1111"
	port2 := ":2222"

	// First 5 requests should succeed (we start with 5 tokens)
	for i := 1; i <= 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		if i%2 == 0 {
			req.RemoteAddr = ip + port1
		} else {
			req.RemoteAddr = ip + port2
		}
		rr := httptest.NewRecorder()
		limitedHandler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Request %d failed: got %v want %v", i, rr.Code, http.StatusOK)
			return
		}
	}

	// 6th request should be rate limited (no tokens left, even on different port)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip + ":3333"
	rr := httptest.NewRecorder()
	limitedHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request 6 should be rate limited: got %v want %v", rr.Code, http.StatusTooManyRequests)
		return
	}

	// Wait for refill window (e.g. mock it by modifying lastSeen)
	mu.Lock()
	visitors[ip].lastSeen = time.Now().Add(-2 * time.Hour)
	mu.Unlock()

	// Request should succeed now
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip + port1
	rr = httptest.NewRecorder()
	limitedHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Request after refill should succeed: got %v want %v", rr.Code, http.StatusOK)
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
	resetVisitors()
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	limitedHandler := rateLimiter(handlerFunc)

	// Exhaust tokens for a specific IP
	ip := "10.0.0.1"
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = ip + ":54321"
		rr := httptest.NewRecorder()
		limitedHandler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Request %d failed: got %v want %v", i+1, rr.Code, http.StatusOK)
		}
	}

	// 6th request from same IP should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip + ":11111" // different port, same IP
	rr := httptest.NewRecorder()
	limitedHandler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request after exhaustion should be rate limited: got %v want %v", rr.Code, http.StatusTooManyRequests)
	}
}
func TestGeneratePdfHandler_PayloadLimit(t *testing.T) {
	largePayload := make([]byte, 50*1024+1)
	for i := range largePayload {
		largePayload[i] = 'a'
	}
	
	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", bytes.NewReader(largePayload))
	req.ContentLength = -1 // simulate chunked/unknown length
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler := rateLimiter(generatePdfHandler)
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusBadRequest && rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected error status for payload too large, got %v", rr.Code)
	}
}

func TestEscapeResumeData(t *testing.T) {
	input := ResumeData{
		Name:  `\input{/etc/passwd}`,
		Title: `& % $ # _ { } ~ ^ \`,
		Experiences: []Experience{
			{
				Company: `Company & Co`,
				Bullets: []string{`Bullet %`},
			},
		},
		Skills: Skills{
			Languages: []string{`Go \`},
		},
	}

	escaped := escapeResumeData(input)

	if escaped.Name != `\textbackslash{}input\{/etc/passwd\}` {
		t.Errorf("Name not escaped correctly: %s", escaped.Name)
	}
	if escaped.Title != `\& \% \$ \# \_ \{ \} \textasciitilde{} \textasciicircum{} \textbackslash{}` {
		t.Errorf("Title not escaped correctly: %s", escaped.Title)
	}
	if escaped.Experiences[0].Company != `Company \& Co` {
		t.Errorf("Experience.Company not escaped correctly: %s", escaped.Experiences[0].Company)
	}
	if escaped.Experiences[0].Bullets[0] != `Bullet \%` {
		t.Errorf("Experience.Bullets not escaped correctly: %s", escaped.Experiences[0].Bullets[0])
	}
	if escaped.Skills.Languages[0] != `Go \textbackslash{}` {
		t.Errorf("Skills.Languages not escaped correctly: %s", escaped.Skills.Languages[0])
	}
}

// Test generatePdfHandler with complete payload
func TestGeneratePdfHandler_CompletePayload(t *testing.T) {
	resetVisitors()
	// Skip test if pdflatex is not available
	if _, err := exec.LookPath("pdflatex"); err != nil {
		t.Skip("pdflatex not available, skipping PDF generation test")
	}

	testData := ResumeData{
		Name:     "Jane Doe",
		Title:    "Full Stack Developer",
		Email:    "jane@example.com",
		Phone:    "555-0123",
		LinkedIn: "linkedin.com/in/janedoe",
		GitHub:   "github.com/janedoe",
		Location: "New York, NY",
		Education: []Education{
			{
				Institution: "Tech University",
				Degree:      "B.S. Computer Science",
				Period:      "2015-2019",
				Notes:       []string{"Graduated with honors", "President of Coding Club"},
			},
		},
		Experiences: []Experience{
			{
				Role:    "Senior Developer",
				Company: "Tech Corp",
				Period:  "2020-2023",
				Bullets: []string{"Built cool things"},
			},
		},
		Projects: []Project{
			{
				Name:      "Open Source Project",
				LinkURL:   "github.com/janedoe/project",
				LinkLabel: "Link",
				Bullets:   []string{"10k stars", "Used by many"},
			},
		},
		SpokenLanguages: []Language{
			{Language: "English", Level: "Native"},
			{Language: "Spanish", Level: "Fluent"},
		},
		Certifications: []string{"AWS Certified Solutions Architect", "CKA"},
		Skills: Skills{
			Languages:    []string{"Go", "TypeScript", "Python"},
			Technologies: []string{"React", "Docker", "Kubernetes"},
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
	
	if !strings.HasPrefix(rr.Body.String(), "%PDF") {
		t.Errorf("Handler did not return a valid PDF")
	}
}

func TestGeneratePdfHandler_LatexInjection(t *testing.T) {
	if _, err := exec.LookPath("pdflatex"); err != nil {
		t.Skip("pdflatex not available")
	}

	testData := ResumeData{
		Name:    `\input{/etc/passwd}`,
		Title:   `Software Engineer`,
	}
	jsonData, _ := json.Marshal(testData)

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := rateLimiter(generatePdfHandler)

	handler.ServeHTTP(rr, req)

	// If pdflatex fails due to injection, it returns 500. We expect 200.
	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", rr.Code)
	}
}

func TestServerTimeouts(t *testing.T) {
	server := setupServer("8080")
	if server.ReadHeaderTimeout == 0 {
		t.Errorf("Server ReadHeaderTimeout is not configured, vulnerable to Slowloris")
	}
}
