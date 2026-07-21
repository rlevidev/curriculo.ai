# Spec: curriculo.ai MVP

> Source: [GitHub Issue #1](https://github.com/rlevidev/curriculo.ai/issues/1)

## Problem Statement

Pessoas que buscam emprego na área de tecnologia têm seus currículos filtrados automaticamente por sistemas ATS (Applicant Tracking System). Currículos mal formatados — criados em Word, Google Docs ou Canva — são frequentemente rejeitados antes de chegar a um humano. O usuário precisa de uma ferramenta que gere currículos em formato ATS-friendly, com preview em tempo real e exportação em PDF de alta qualidade compilado em LaTeX, sem precisar saber LaTeX.

## Solution

Um web app split em duas partes:

- **Frontend (React + Vite no GitHub Pages):** Formulário com seções dinâmicas (N educações, N experiências, etc.), preview client-side em HTML/CSS que imita o layout do PDF final, ATS score baseado em checklist determinístico, e persistência de rascunho via localStorage.
- **Backend (Go no Render):** Endpoint stateless que recebe JSON, renderiza um template LaTeX via `text/template`, compila com `pdflatex`, e devolve o PDF. Sem banco de dados, sem autenticação, sem conceito de usuário.

O resultado é um PDF de uma página em formato ATS otimizado, compilado em LaTeX com fontes Libertine/Biolinum.

## User Stories

1. As a job seeker, I want to fill in my resume sections in a structured form, so that I don't need to know LaTeX to produce a professionally formatted resume.
2. As a job seeker, I want to see a live preview of my resume as I type, so that I can iterate on the content without waiting for PDF compilation.
3. As a job seeker, I want to export my resume as a PDF, so that I can submit it to job applications.
4. As a job seeker, I want the PDF to be ATS-friendly (compiled in LaTeX, no images, selectable text), so that automated systems can parse my resume correctly.
5. As a job seeker, I want to add multiple education entries, so that I can list all my degrees.
6. As a job seeker, I want to add multiple professional experiences with bullet points, so that I can describe my achievements at each company.
7. As a job seeker, I want to add multiple projects with links, so that I can showcase my portfolio.
8. As a job seeker, I want to list my programming languages separately from my technologies, so that recruiters can quickly identify my technical skills.
9. As a job seeker, I want to list my spoken languages with proficiency levels, so that international roles can see my language capabilities.
10. As a job seeker, I want to list my certifications, so that I can prove my professional qualifications.
11. As a job seeker, I want sections I didn't fill in to automatically disappear from the PDF, so that my resume doesn't have empty sections wasting space.
12. As a job seeker, I want my draft to be saved automatically in my browser, so that I don't lose my work if I close the tab or refresh the page.
13. As a job seeker, I want to see an ATS score based on a checklist of best practices, so that I know how to improve my resume before submitting.
14. As a job seeker, I want to see the ATS score update in real time as I fill in the form, so that I get immediate feedback on completeness.
15. As a job seeker, I want to see which specific criteria I'm missing to reach a higher ATS score, so that I know exactly what to improve.
16. As a job seeker, I want to see a clear server status indicator, so that I know whether the backend is ready before I try to export.
17. As a job seeker, I want to receive a meaningful loading message when the server is waking up (cold start), so that I don't think the app is broken.
18. As a job seeker, I want to be told my position in the queue if the server is busy, so that I know how long to wait.
19. As a job seeker, I want the form to prevent me from entering excessively long text, so that my resume stays a reasonable length.
20. As a job seeker, I want to see a warning when my content might exceed one page, so that I can edit before exporting.
21. As a job seeker, I want the preview to display a disclaimer that it's an approximation of the final PDF, so that I'm not surprised by minor formatting differences.
22. As a job seeker, I want to use the app on my desktop browser without issues, so that I can fill in my resume comfortably.
23. As a job seeker, I want the app to be minimally functional on mobile (stacked layout), so that I can review my resume on my phone even if filling it in is better on desktop.
24. As a job seeker, I want to use special characters in my name and text (accents, &, %, etc.), so that my resume renders correctly regardless of my language or content.

## Implementation Decisions

### Architecture

- **Frontend:** React + Vite, deployed to GitHub Pages as static files. Single-page app with client-side preview. No server-side rendering.
- **Backend:** Go HTTP server, deployed to Render free tier inside a Docker container with TeXLive. Stateless — receives JSON, returns PDF. Two endpoints: `POST /generate-pdf` and `GET /health`.
- **Communication:** Frontend calls backend via CORS-enabled fetch. Backend URL configured via Vite environment variable (`VITE_API_URL`), not hardcoded.

### Monorepo Structure

```
curriculo-ai/
├── frontend/              # React + Vite (GitHub Pages)
│   ├── package.json
│   ├── vite.config.ts
│   ├── index.html
│   ├── public/
│   └── src/
│       ├── App.tsx
│       ├── main.tsx
│       ├── components/
│       └── ...
├── backend/               # Go API + LaTeX (Render)
│   ├── go.mod
│   ├── main.go
│   ├── template.tex
│   ├── Dockerfile
│   └── ...
├── .github/
│   └── workflows/
│       ├── deploy-frontend.yml
│       ├── deploy-backend.yml
│       └── keep-alive.yml
└── README.md
```

### Domain Model — JSON Schema (API Contract)

All code (variables, functions, keys) is in English. UI-facing text (labels, messages, placeholders) is in Portuguese.

```json
{
  "name": "string (max 100)",
  "title": "string (max 100)",
  "contact": {
    "email": "string (optional)",
    "phone": "string (optional)",
    "linkedin_url": "string URL (optional)",
    "github_url": "string URL (optional)",
    "location": "string (optional)"
  },
  "education": [
    {
      "institution": "string (max 100)",
      "degree": "string (max 100)",
      "period": "string (max 50)",
      "notes": ["string (max 500)"]
    }
  ],
  "skills": {
    "languages": ["string"],
    "technologies": ["string"]
  },
  "experiences": [
    {
      "company": "string (max 100)",
      "location": "string (max 100, optional)",
      "period": "string (max 50)",
      "role": "string (max 100)",
      "bullets": ["string (max 500)"]
    }
  ],
  "projects": [
    {
      "name": "string (max 100)",
      "link_url": "string URL (optional)",
      "link_label": "string (max 50, defaults to 'Link')",
      "bullets": ["string (max 500)"]
    }
  ],
  "languages": [
    {
      "language": "string (max 50)",
      "level": "string (max 50)"
    }
  ],
  "certifications": ["string (max 200)"]
}
```

Key naming: `skills.languages` (programming languages) is nested inside `skills`. `languages` at root level is for spoken languages. This resolves the naming collision.

### Template Strategy

- Single template for MVP. No template selection UI.
- Go backend uses `text/template` with custom delimiters `<[` and `]>` to avoid conflict with LaTeX braces.
- LaTeX special characters (`&`, `%`, `_`, `#`, `$`, `{`, `}`, `~`, `^`, `\`) are escaped via a registered template function (`texEscape`) applied at render time.
- Empty sections are conditionally omitted from the LaTeX output using `<[if .Education]>` blocks.
- The `-no-shell-escape` flag is passed explicitly to `pdflatex` for defense in depth.

### Section Order

Fixed in MVP: Header → Contact → Education → Skills → Experience → Projects → Languages → Certifications. Dictated by the backend template.

### Validation

- **Both sides (defense in depth):** Frontend validates for UX (inline errors, character counters). Backend validates as the authority (rejects with 400).
- **Limits:** Name/title/company/institution: 100 chars. Bullets: 500 chars. Max 8 entries per section. Total payload: 50KB.
- **URL fields:** Allowlist regex for `https?://` with safe characters.

### Rate Limiting

- Token bucket: 5 requests/minute per IP.
- Key source: `r.RemoteAddr` (TCP-level, not spoofable).
- Cleanup: background goroutine every 5 minutes removes buckets inactive for 10+ minutes.

### Concurrency Control

- Semaphore (`chan struct{}`) with capacity 3 limits concurrent `pdflatex` processes.
- Requests beyond capacity wait with a timeout. Timeout → `503` with queue position.
- Frontend shows "Servidor ocupado — posição N na fila" with automatic retry.

### Cold Start Handling

- `GET /health` pinged on React app mount via `useEffect`.
- Topbar indicator: green "servidor online" / yellow "servidor acordando...".
- GitHub Actions pings `/health` every 14 minutes.
- Frontend uses `AbortController` with 60s timeout on PDF export.

### ATS Score

- Deterministic checklist computed in the frontend (no AI).
- Updates in real time. Breakdown visible to the user.
- Criteria: name, email, experience with bullets, metrics in bullets, skills, education, LinkedIn, GitHub, projects, all sections populated.

### Persistence

- `localStorage` only. No backend state.
- Draft auto-saved on every form change (debounced).
- Restored on page load.

### Frontend Details

- Desktop-first. Responsive stacking at < 880px.
- Tabs "Templates" and "Histórico" removed from topbar in MVP.
- Page overflow warning when content exceeds A4 height (1123px at 96dpi).
- Disclaimer: "prévia aproximada — o PDF final pode variar levemente."

### Error Handling

- Backend errors always JSON `{ "message": "..." }`, never PDF bytes.
- Log sanitization: `pdflatex` stderr only logged in full when `DEBUG_LOGS=1`. Production uses generic messages.

## Testing Decisions

Tests verify **external behavior at seam boundaries**, not implementation details.

### Seam 1 — API Contract (Integration)

- Valid JSON → `200` with `Content-Type: application/pdf` and non-empty body.
- Missing required fields → `400` with error JSON.
- Exceeding size limits → `400`.
- LaTeX-breaking characters (`&`, `%`, `_`, `#`) → PDF generates successfully.
- Malicious URL fields → `400`.
- Rate limit exceeded → `429`.
- `GET /health` → `200`.

### Seam 2 — Resume State → Preview (Unit, Frontend)

- Complete state → all sections rendered.
- Empty certifications → section not rendered.
- HTML-dangerous characters → escaped (no XSS).
- ATS score matches expected criteria.

### Seam 3 — Go Handler → LaTeX Renderer (Unit, Backend)

- Valid struct → correct LaTeX output.
- Empty sections → `\section{}` blocks omitted.
- `texEscape` handles all 10 special characters.
- Invalid input → error before reaching template.

### Seam 4 — LaTeX → pdflatex (Integration, Backend)

- Known-good `.tex` → non-empty PDF.
- `-no-shell-escape` → `\write18` blocked.
- Compilation respects 20s timeout.

### Prior Art

Greenfield. Go `*_test.go` with `testing` package. Frontend with Vitest.

## Out of Scope

- **AI features:** No LLM. The ".ai" is branding only.
- **Multiple templates:** Single template. Selection UI deferred.
- **User accounts:** No login, no server-side persistence.
- **Section reordering:** Fixed order.
- **Per-experience technology tags:** Deferred. Users mention tech in bullets.
- **Mobile-optimized UX:** Desktop-first, mobile tolerável.
- **Internationalization:** Portuguese UI only.
- **Custom contact fields:** Fixed set (email, phone, LinkedIn, GitHub, location).
- **Multi-page resume:** Single-page target with overflow warning.
- **Tabs "Templates" and "Histórico":** Removed from MVP.

## Further Notes

- **Maior risco:** Abandono do projeto. Arquitetura desenhada para minimizar escopo.
- **Divergência preview/PDF:** Arial ≠ Libertine/Biolinum. Disclaimer explicita isso.
- **Docker image:** ~1GB com TeXLive. Todos os 4 pacotes são necessários.
- **Keep-alive:** Ping de 14 min não cobre 100% dos casos. Health ping no frontend é segunda defesa.
- **Backend URL:** `VITE_API_URL` baked at build time. `.env` files para dev/prod.
