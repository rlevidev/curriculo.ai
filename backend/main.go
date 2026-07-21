package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func rateLimiter(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		visitor, ok := visitors[r.RemoteAddr]
		if !ok {
			visitor = &Visitor{lastSeen: time.Now(), tokens: 5}
			visitors[r.RemoteAddr] = visitor
		}

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
	return s
}

func generatePdfHandler(w http.ResponseWriter, r *http.Request) {
	// Validate size
	if r.ContentLength > 50*1024 {
		http.Error(w, "payload too large", http.StatusBadRequest)
		return
	}

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

	dir, err := ioutil.TempDir("", "resume")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	texPath := filepath.Join(dir, "resume.tex")
	if err := ioutil.WriteFile(texPath, buf.Bytes(), 0644); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("pdflatex", "-no-shell-escape", "-output-directory", dir, texPath)
	if err := cmd.Run(); err != nil {
		http.Error(w, "failed to generate PDF", http.StatusInternalServerError)
		return
	}

	pdfPath := filepath.Join(dir, "resume.pdf")
	pdf, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		http.Error(w, "failed to read generated PDF", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"curriculo.pdf\"")
	w.Write(pdf)
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
	})
	http.HandleFunc("/generate-pdf", rateLimiter(generatePdfHandler))
	http.ListenAndServe(":8080", nil)
}
