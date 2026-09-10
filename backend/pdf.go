package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

const resumeTemplateSrc = `
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
`

var (
	resumeTmpl     *template.Template
	resumeTmplOnce sync.Once
	resumeTmplErr  error
)

// getResumeTemplate compiles the LaTeX template once (sync.Once) and reuses it on subsequent calls.
func getResumeTemplate() (*template.Template, error) {
	resumeTmplOnce.Do(func() {
		t := template.New("resume").Delims("<[", "]>").Funcs(template.FuncMap{
			"add": func(a, b int) int { return a + b },
		})
		resumeTmpl, resumeTmplErr = t.Parse(resumeTemplateSrc)
	})
	return resumeTmpl, resumeTmplErr
}

// generatePDF renders the template into a temp dir, runs pdflatex with a 60s timeout, and returns the PDF bytes.
func generatePDF(ctx context.Context, data ResumeData) ([]byte, error) {
	tmpl, err := getResumeTemplate()
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("render template: %w", err)
	}
	dir, err := os.MkdirTemp("", "resume")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(dir)
	texPath := filepath.Join(dir, "resume.tex")
	if err := os.WriteFile(texPath, buf.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("write tex: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "pdflatex", "-no-shell-escape", "-output-directory", dir, texPath)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pdflatex: %w", err)
	}
	pdfPath := filepath.Join(dir, "resume.pdf")
	pdf, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("read pdf: %w", err)
	}
	return pdf, nil
}
