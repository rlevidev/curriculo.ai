package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
)

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simple validation
	if len(data.Name) == 0 || len(data.Title) == 0 {
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
\documentclass{article}
\usepackage{hyperref}
\begin{document}
\section*{<[ .Name ]>}
\subsection*{<[ .Title ]>}
<[ .Email ]> | <[ .Phone ]><[ if .Location ]> | <[ .Location ]><[ end ]><[ if .LinkedIn ]> | \href{https://<[ .LinkedIn ]>}{<[ .LinkedIn ]>}<[ end ]><[ if .GitHub ]> | \href{https://<[ .GitHub ]>}{<[ .GitHub ]>}<[ end ]>

<[ if .Education ]>
\section*{Education}
<[ range .Education ]>
\textbf{<[ .Institution ]>} -- <[ .Period ]> \\
<[ .Degree ]>
<[ if .Notes ]>
\begin{itemize}
<[ range .Notes ]> \item <[ . ]> <[ end ]>
\end{itemize}
<[ end ]>
<[ end ]>
<[ end ]>

<[ if .Skills.Languages ]>
\section*{Skills}
\textbf{Languages:} <[ $lenLangs := len .Skills.Languages ]><[ range $i, $lang := .Skills.Languages ]><[ $lang ]><[ if  lt (add $i 1) $lenLangs ]> $\cdot$ <[ end ]><[ end ]>
<[ end ]>

<[ if .Skills.Technologies ]>
\textbf{Technologies:} <[ $lenTechs := len .Skills.Technologies ]><[ range $i, $tech := .Skills.Technologies ]><[ $tech ]><[ if lt (add $i 1) $lenTechs ]> $\cdot$ <[ end ]><[ end ]>
<[ end ]>

<[ if .Experiences ]>
\section*{Experience}
<[ range .Experiences ]>
\textbf{<[ .Role ]>} @ <[ .Company ]> (<[ .Period ]>)<[ if .Location ]> -- <[ .Location ]><[ end ]>
\begin{itemize}
<[ range .Bullets ]> \item <[ . ]> <[ end ]>
\end{itemize}
<[ end ]>
<[ end ]>

<[ if .Projects ]>
\section*{Projects}
<[ range .Projects ]>
\textbf{<[ .Name ]>} <[ if .LinkURL ]>-- \href{<[ .LinkURL ]>}{<[ if .LinkLabel ]><[ .LinkLabel ]><[ else ]><[ .LinkURL ]><[ end ]>}<[ end ]>
<[ if .Bullets ]>
\begin{itemize}
<[ range .Bullets ]> \item <[ . ]> <[ end ]>
\end{itemize}
<[ end ]>
<[ end ]>
<[ end ]>

<[ if .SpokenLanguages ]>
\section*{Languages}
<[ $lenSpoken := len .SpokenLanguages ]><[ range $i, $lang := .SpokenLanguages ]><[ .Language ]> (<[ .Level ]>)<[ if lt (add $i 1) $lenSpoken ]> $\cdot$ <[ end ]><[ end ]>
<[ end ]>

<[ if .Certifications ]>
\section*{Certifications}
<[ $lenCerts := len .Certifications ]><[ range $i, $cert := .Certifications ]><[ . ]><[ if lt (add $i 1) $lenCerts ]> $\cdot$ <[ end ]><[ end ]>
<[ end ]>
\end{document}
`)

	var buf bytes.Buffer
	tmpl.Execute(&buf, data)

	dir, err := os.MkdirTemp("", "resume")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	texPath := filepath.Join(dir, "resume.tex")
	if err := os.WriteFile(texPath, buf.Bytes(), 0644); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pdflatex", "-no-shell-escape", "-output-directory", dir, texPath)
	if err := cmd.Run(); err != nil {
		http.Error(w, "failed to generate PDF", http.StatusInternalServerError)
		return
	}

	pdfPath := filepath.Join(dir, "resume.pdf")
	pdf, err := os.ReadFile(pdfPath)
	if err != nil {
		http.Error(w, "failed to read generated PDF", http.StatusInternalServerError)
		return
	}

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
	mux.HandleFunc("/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
	}))
	mux.HandleFunc("/generate-pdf", corsMiddleware(rateLimiter(generatePdfHandler)))

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

	server := setupServer(port)
	server.ListenAndServe()
}
