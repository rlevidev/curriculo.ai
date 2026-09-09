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
		} else {
			// replenish if enough time passed
			elapsed := time.Since(visitor.lastSeen)
			if elapsed > 1*time.Hour {
				visitor.tokens = 5
			}
		}

		visitor.lastSeen = time.Now()

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
	Company  string   `json:"company"`
	Location string   `json:"location"`
	Role     string   `json:"role"`
	Period   string   `json:"period"`
	Bullets  []string `json:"bullets"`
}

type Project struct {
	Name      string   `json:"name"`
	LinkURL   string   `json:"link_url"`
	LinkLabel string   `json:"link_label"`
	Bullets   []string `json:"bullets"`
}

type Language struct {
	Language string `json:"language"`
	Level    string `json:"level"`
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

func escapeResumeData(d ResumeData) ResumeData {
	d.Name = texEscape(d.Name)
	d.Title = texEscape(d.Title)
	d.Email = texEscape(d.Email)
	d.Phone = texEscape(d.Phone)
	d.LinkedIn = texEscape(d.LinkedIn)
	d.GitHub = texEscape(d.GitHub)
	d.Location = texEscape(d.Location)

	for i := range d.Education {
		d.Education[i].Institution = texEscape(d.Education[i].Institution)
		d.Education[i].Degree = texEscape(d.Education[i].Degree)
		d.Education[i].Period = texEscape(d.Education[i].Period)
		for j := range d.Education[i].Notes {
			d.Education[i].Notes[j] = texEscape(d.Education[i].Notes[j])
		}
	}

	for i := range d.Experiences {
		d.Experiences[i].Company = texEscape(d.Experiences[i].Company)
		d.Experiences[i].Location = texEscape(d.Experiences[i].Location)
		d.Experiences[i].Role = texEscape(d.Experiences[i].Role)
		d.Experiences[i].Period = texEscape(d.Experiences[i].Period)
		for j := range d.Experiences[i].Bullets {
			d.Experiences[i].Bullets[j] = texEscape(d.Experiences[i].Bullets[j])
		}
	}

	for i := range d.Projects {
		d.Projects[i].Name = texEscape(d.Projects[i].Name)
		d.Projects[i].LinkURL = texEscape(d.Projects[i].LinkURL)
		d.Projects[i].LinkLabel = texEscape(d.Projects[i].LinkLabel)
		for j := range d.Projects[i].Bullets {
			d.Projects[i].Bullets[j] = texEscape(d.Projects[i].Bullets[j])
		}
	}

	for i := range d.SpokenLanguages {
		d.SpokenLanguages[i].Language = texEscape(d.SpokenLanguages[i].Language)
		d.SpokenLanguages[i].Level = texEscape(d.SpokenLanguages[i].Level)
	}

	for i := range d.Certifications {
		d.Certifications[i] = texEscape(d.Certifications[i])
	}

	for i := range d.Skills.Languages {
		d.Skills.Languages[i] = texEscape(d.Skills.Languages[i])
	}

	for i := range d.Skills.Technologies {
		d.Skills.Technologies[i] = texEscape(d.Skills.Technologies[i])
	}

	return d
}

func generatePdfHandler(w http.ResponseWriter, r *http.Request) {
	// Validate size - limit body read to 50KB regardless of ContentLength
	r.Body = http.MaxBytesReader(w, r.Body, 50*1024)

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

	// Escape latex characters in the payload
	data = escapeResumeData(data)

	// Semaphore
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	tmpl := template.New("resume").Delims("<[", "]>").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	})
	tmpl, _ = tmpl.Parse(`
%!TEX program = pdflatex
\documentclass[11pt,a4paper]{article}

\usepackage[T1]{fontenc}
\usepackage[left=0.6in,right=0.6in,top=0.45in,bottom=0.45in]{geometry}
\usepackage{array}
\usepackage{xcolor}
\usepackage{hyperref}
\usepackage{enumitem}
\usepackage{titlesec}
\usepackage[utf8]{inputenc}
\usepackage{libertine}

\definecolor{ink}{HTML}{211F1A}
\definecolor{muted}{HTML}{7A715F}
\definecolor{hline}{HTML}{E1DCCD}
\definecolor{bodytext}{HTML}{3A362E}

\color{ink}

\hypersetup{
    colorlinks=true,
    linkcolor=ink,
    urlcolor=ink,
    citecolor=ink,
    pdfborder={0 0 0}
}

\titleformat{\section}
  {\normalfont\sffamily\bfseries\small\color{ink}}
  {}{0em}{}
  [\vspace{2pt}{\color{hline}\titlerule[0.6pt]}]
\titlespacing{\section}{0pt}{4pt}{2pt}

\pagestyle{empty}
\setlength{\parindent}{0pt}
\setlength{\parskip}{0.5pt}
\setlist[itemize]{itemsep=0.5pt, parsep=0.5pt, topsep=1pt, leftmargin=15pt}

\clubpenalty=10000
\widowpenalty=10000
\displaywidowpenalty=10000
\raggedbottom

\newcommand{\entryhead}[2]{%
  \noindent{\rmfamily\bfseries #1}\hfill{\sffamily\small\color{muted}#2}\\
}
\newcommand{\entrysub}[1]{%
  {\sffamily\itshape\small\color{muted}#1}\par
}

\begin{document}

\begin{center}
    {\rmfamily\Huge\bfseries <[ .Name ]>}\\[2pt]
    {\sffamily\large\color{muted}<[ .Title ]>}\\[2pt]
    {\sffamily\small\color{muted}
      <[ .Email ]><[ if .Phone ]> \textperiodcentered\ <[ .Phone ]><[ end ]><[ if .Location ]> \textperiodcentered\ <[ .Location ]><[ end ]><[ if .LinkedIn ]> \textperiodcentered\ \href{https://<[ .LinkedIn ]>}{LinkedIn}<[ end ]><[ if .GitHub ]> \textperiodcentered\ \href{https://<[ .GitHub ]>}{GitHub}<[ end ]>%
    }
\end{center}

\vspace{4pt}
\noindent{\color{ink}\rule{\linewidth}{1.1pt}}
\vspace{1pt}

<[ if .Education ]>
\section{EDUCATION}
<[ range .Education ]>
\entryhead{<[ .Institution ]>}{<[ .Period ]>}
\entrysub{<[ .Degree ]>}
<[ if .Notes ]>
{\sffamily\small\color{bodytext}
\begin{itemize}
<[ range .Notes ]>    \item <[ . ]>
<[ end ]>
\end{itemize}
}
<[ end ]>
<[ end ]>
<[ end ]>

<[ if or .Skills.Languages .Skills.Technologies ]>
\section{TECHNICAL SKILLS}
{\sffamily\small\color{bodytext}
<[ if .Skills.Languages ]>
\textbf{\color{ink}Languages:} <[ $lenLangs := len .Skills.Languages ]><[ range $i, $lang := .Skills.Languages ]><[ $lang ]><[ if lt (add $i 1) $lenLangs ]>, <[ end ]><[ end ]>
<[ end ]>
<[ if .Skills.Technologies ]>
<[ if .Skills.Languages ]> \par\vspace{4pt}
<[ end ]>
\textbf{\color{ink}Technologies:} <[ $lenTechs := len .Skills.Technologies ]><[ range $i, $tech := .Skills.Technologies ]><[ $tech ]><[ if lt (add $i 1) $lenTechs ]>, <[ end ]><[ end ]>
<[ end ]>
}
<[ end ]>

<[ if .Experiences ]>
\section{PROFESSIONAL EXPERIENCE}
<[ range .Experiences ]>
\entryhead{<[ .Company ]>}{<[ .Period ]>}
\entrysub{<[ .Role ]><[ if .Location ]> \hfill <[ .Location ]><[ end ]>}
{\sffamily\small\color{bodytext}
\begin{itemize}
<[ range .Bullets ]>    \item <[ . ]>
<[ end ]>
\end{itemize}
}
<[ end ]>
<[ end ]>

<[ if .Projects ]>
\section{PROJECTS}
<[ range .Projects ]>
\entryhead{<[ .Name ]>}{<[ if .LinkURL ]>\href{<[ .LinkURL ]>}{<[ if .LinkLabel ]><[ .LinkLabel ]><[ else ]>Link<[ end ]>}<[ end ]>}
<[ if .Bullets ]>
{\sffamily\small\color{bodytext}
\begin{itemize}
<[ range .Bullets ]>    \item <[ . ]>
<[ end ]>
\end{itemize}
}
<[ end ]>
<[ end ]>
<[ end ]>

<[ if .SpokenLanguages ]>
\section{LANGUAGES}
{\sffamily\small\color{bodytext}
<[ range $i, $lang := .SpokenLanguages ]><[ if $i ]> \textperiodcentered\ <[ end ]>\textbf{\color{ink}<[ .Language ]>:} <[ .Level ]><[ end ]>
}
<[ end ]>

<[ if .Certifications ]>
\section{CERTIFICATIONS}
{\sffamily\small\color{bodytext}
<[ $lenCerts := len .Certifications ]><[ range $i, $cert := .Certifications ]><[ . ]><[ if lt (add $i 1) $lenCerts ]> \textperiodcentered\ <[ end ]><[ end ]>
}
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
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		logOutput := output.String()
		if len(logOutput) > 8000 {
			logOutput = logOutput[len(logOutput)-8000:]
		}
		texContent := string(buf.Bytes())
		if len(texContent) > 4000 {
			texContent = texContent[len(texContent)-4000:]
		}
		loggerFor(r).Error("pdflatex failed", "op", "pdf_generate", "status", http.StatusInternalServerError, "latency_ms", time.Since(pdfStart).Milliseconds(), "error", err.Error(), "output", logOutput, "tex", texContent)
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseLogger.Info("starting", "op", "startup", "port", port)
	server := setupServer(port)
	server.ListenAndServe()
}
