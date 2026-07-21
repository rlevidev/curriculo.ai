# Tickets: curriculo.ai MVP

Vertical slices for the curriculo.ai MVP. Source spec: [SPEC.md](SPEC.md) / [Issue #1](https://github.com/rlevidev/curriculo.ai/issues/1).

Work the **frontier**: any ticket whose blockers are all done.

```
#1 Scaffolding
├── #2 Header tracer bullet
│   ├── #3 Education + Skills ──────┐
│   ├── #4 Experience + Projects ───┤── #7 ATS Score
│   ├── #5 Languages + Certs ───────┘
│   ├── #6 Validation + Security
│   └── #8 localStorage + Cold start
└── #9 Deploy pipeline
```

---

## 1. Monorepo scaffolding

**What to build:** React + Vite app in `frontend/` and Go HTTP server in `backend/`, both running locally. Dockerfile that compiles and serves the Go binary with TeXLive installed. `GET /health` returns `200`. Frontend shows an empty shell with the topbar layout (brand, server status indicator placeholder, export button placeholder). No features yet — just the skeleton that everything else builds on.

**Blocked by:** None — can start immediately.

- [ ] `cd frontend && npm run dev` starts a React dev server with the app shell (topbar + two-pane layout)
- [ ] `cd backend && go run .` starts an HTTP server that responds `200` on `GET /health`
- [ ] `docker build` in `backend/` produces a working image with TeXLive (`pdflatex` available inside the container)
- [ ] CORS middleware configured — frontend dev server can reach backend without browser errors
- [ ] `VITE_API_URL` environment variable wired (defaults to `http://localhost:8080` in dev)
- [ ] Go module initialized (`go.mod`) with project name
- [ ] TypeScript configured in frontend with strict mode

---

## 2. Header tracer bullet: form → preview → PDF

**What to build:** The full end-to-end pipeline, scoped to just the header (name, title, contact fields). User fills in the form on the left, sees a live preview on the right that mimics the LaTeX layout, clicks "Exportar PDF", and downloads a real PDF compiled from LaTeX. This proves every layer works together: React state → preview rendering → JSON payload → Go handler → `text/template` with `<[` `]>` delimiters → `texEscape` function → `pdflatex -no-shell-escape` → PDF bytes → blob download in the browser.

**Blocked by:** #1 Monorepo scaffolding

- [ ] Form with fields: name, title, email, phone, LinkedIn URL, GitHub URL, location — all wired to React state
- [ ] Preview pane renders header in HTML/CSS matching the LaTeX template style (centered name, title, contact row with separators)
- [ ] Preview updates in real time with debounce (~200ms)
- [ ] `POST /generate-pdf` endpoint accepts the resume JSON and returns `application/pdf` with `Content-Disposition` header
- [ ] Go backend renders LaTeX via `text/template` with `<[` `]>` delimiters and a registered `texEscape` template function
- [ ] `pdflatex` runs with `-no-shell-escape` flag, 20s timeout, in an isolated temp directory per request
- [ ] Temp files cleaned up after each request (success or failure)
- [ ] Backend error responses are JSON `{ "message": "..." }`, not PDF bytes
- [ ] Frontend handles export: loading state on button, blob download with `Content-Disposition` filename, error display
- [ ] Export button disabled during request (double-click guard)
- [ ] Integration test: send valid header JSON to `/generate-pdf`, verify non-empty PDF response
- [ ] Unit test: `texEscape` correctly escapes all 10 LaTeX special characters (`&`, `%`, `_`, `#`, `$`, `{`, `}`, `~`, `^`, `\`)

---

## 3. Dynamic sections: Education + Skills

**What to build:** Education as the first dynamic list section — user adds/removes up to 8 education entries, each with institution, degree, period, and optional notes. Skills as two separate input lists (programming languages and technologies) nested under `skills{}` in the JSON. Both sections appear in the preview and in the exported PDF.

**Blocked by:** #2 Header tracer bullet

- [ ] Education form: add/remove entries (max 8), each with institution, degree, period, and optional notes (list of strings)
- [ ] Skills form: two separate inputs for `skills.languages` and `skills.technologies`
- [ ] Preview renders Education section with entry layout matching LaTeX template (`entryhead` for institution/period, `entrysub` for degree, bulleted notes)
- [ ] Preview renders Skills section with two labeled lines ("Languages:" and "Technologies:")
- [ ] LaTeX template updated with `<[range .Education]>` loop and `<[if .Skills.Languages]>` / `<[if .Skills.Technologies]>` conditionals
- [ ] Backend Go struct updated with `Education` and `Skills` fields matching the JSON schema
- [ ] Unit test: rendered `.tex` contains correct LaTeX for multiple education entries
- [ ] Integration test: full resume with education + skills produces valid PDF

---

## 4. Dynamic sections: Experience + Projects

**What to build:** Experience entries with nested bullet points — the most complex section. User adds/removes up to 8 experiences, each with company, location, period, role, and up to 8 bullets. Projects with name, optional link (URL + label), and bullets. Both render in preview and PDF.

**Blocked by:** #2 Header tracer bullet

- [ ] Experience form: add/remove entries (max 8), each with company, location (optional), period, role, and add/remove bullets (max 8 per entry, max 500 chars each)
- [ ] Projects form: add/remove entries (max 8), each with name, optional link URL, optional link label (defaults to "Link"), and add/remove bullets
- [ ] Preview renders Experience with split layout (company left, period right; role left, location right) and bulleted list
- [ ] Preview renders Projects with name + link and bulleted list
- [ ] LaTeX template updated with `<[range .Experiences]>` and `<[range .Projects]>` loops, including `\href` for project links
- [ ] Backend Go struct updated with `Experiences` and `Projects` fields
- [ ] Unit test: rendered `.tex` contains correct LaTeX for experiences with bullets and projects with links
- [ ] Integration test: resume with experiences + projects produces valid PDF

---

## 5. Simple sections: Languages + Certifications + auto-disappear

**What to build:** Spoken languages (language + proficiency level) and certifications (flat string list). Plus the auto-disappear behavior: any section with no data is omitted from both the preview and the PDF. This applies retroactively to all sections from tickets #3 and #4.

**Blocked by:** #2 Header tracer bullet

- [ ] Languages form: add/remove entries (language + level)
- [ ] Certifications form: add/remove string entries
- [ ] Preview renders Languages inline (e.g., "**Portuguese:** Native · **English:** Fluent")
- [ ] Preview renders Certifications inline with separators
- [ ] LaTeX template wraps every section in `<[if]>` conditionals — empty sections produce no LaTeX output
- [ ] Preview conditionally renders sections — empty ones are hidden
- [ ] Backend Go struct updated with `Languages` and `Certifications` fields
- [ ] Unit test: rendered `.tex` for a resume with empty education/certifications omits those `\section{}` blocks entirely
- [ ] Frontend test: preview with empty certifications does not render that section

---

## 6. Input validation + security hardening

**What to build:** Defense-in-depth validation on both frontend and backend, rate limiting, concurrency control, and LaTeX security. The frontend shows inline errors and character counters. The backend rejects invalid input with clear `400` responses. Rate limiter uses `RemoteAddr` with cleanup. Semaphore limits concurrent compilations.

**Blocked by:** #2 Header tracer bullet

- [ ] Frontend: character counters on text fields, max-length enforcement, inline error messages when limits exceeded
- [ ] Frontend: export button disabled when required fields (name, title) are empty
- [ ] Backend: validate all field lengths (100/500 chars), entry counts (max 8), payload size (max 50KB) — reject with `400` and `{ "message": "..." }`
- [ ] Backend: URL allowlist regex on `linkedin_url`, `github_url`, `link_url` — reject non-`https?://` URLs
- [ ] Backend: rate limiter with token bucket (5 req/min per `r.RemoteAddr`)
- [ ] Backend: goroutine cleans up inactive rate limit buckets every 5 minutes (buckets idle > 10 min)
- [ ] Backend: concurrency semaphore (`chan struct{}` capacity 3) — excess requests wait with timeout, then `503` with queue info
- [ ] Frontend: handles `429` (rate limit) and `503` (queue) responses with user-friendly messages and automatic retry
- [ ] Backend: log sanitization — `pdflatex` stderr only fully logged when `DEBUG_LOGS=1`, generic message in production
- [ ] Integration test: payload exceeding 50KB → `400`
- [ ] Integration test: field exceeding max length → `400`
- [ ] Integration test: malicious URL → `400`
- [ ] Integration test: 6th request in 1 minute → `429`
- [ ] Unit test: rate limiter cleanup removes stale buckets

---

## 7. ATS score + page overflow warning

**What to build:** A deterministic ATS score computed in the frontend based on resume completeness. Displays in the topbar area, updates in real time, and shows a breakdown of criteria so the user knows what to improve. Also: a page overflow warning when content might exceed one A4 page, and a disclaimer that the preview is an approximation.

**Blocked by:** #3 Education + Skills, #4 Experience + Projects, #5 Languages + Certifications

- [ ] ATS score engine: pure function that takes resume state and returns score (0–100) + list of criteria with pass/fail
- [ ] Criteria include: name filled (+5), email (+10), at least 1 experience (+15), bullets with numbers/metrics (+15), both skill lists filled (+10), education filled (+10), LinkedIn (+10), GitHub (+5), at least 1 project (+10), all visible sections populated (+10)
- [ ] Score displayed in topbar area, updates on every form change
- [ ] Score breakdown visible to user (expandable list or tooltip showing what's missing)
- [ ] Page overflow warning: when preview container `scrollHeight > 1123px`, show estimated page count and suggestion to reduce content
- [ ] Overflow check runs after preview DOM update (using `requestAnimationFrame`)
- [ ] Disclaimer text near preview: "prévia aproximada — o PDF final pode variar levemente"
- [ ] Frontend test: ATS score function returns correct score for various resume states (empty, partial, complete)
- [ ] Frontend test: overflow warning appears when content exceeds threshold

---

## 8. localStorage draft + cold start UX

**What to build:** Auto-save the resume draft to `localStorage` so users don't lose work on page refresh. Health check ping on app mount with server status indicator in the topbar. Contextual messaging during cold start and export timeout.

**Blocked by:** #2 Header tracer bullet

- [ ] `localStorage.setItem` on every form change (debounced, same timer as preview)
- [ ] On page load, restore draft from `localStorage` if present — populate form and preview
- [ ] `GET /health` fetched on React app mount (`useEffect`)
- [ ] Topbar indicator: green dot + "servidor online" when health responds, yellow dot + "servidor acordando..." when pending/failed
- [ ] Health check retries periodically (every 10s) while status is not OK
- [ ] Export fetch uses `AbortController` with 60s timeout
- [ ] If export takes > 10s, button text changes to "Servidor acordando — aguarde..."
- [ ] If export times out (60s), show error with retry button
- [ ] Frontend test: draft saved to localStorage contains expected resume state
- [ ] Frontend test: server status indicator reflects health check result

---

## 9. Deploy pipeline + keep-alive

**What to build:** GitHub Actions workflows for frontend deployment to GitHub Pages, keep-alive pings to prevent Render cold starts, and production environment configuration. The backend deploys via Render's auto-deploy on push (no custom workflow needed unless desired).

**Blocked by:** #1 Monorepo scaffolding

- [ ] `deploy-frontend.yml`: triggers on push to `main` (paths: `frontend/**`), runs `npm ci && npm run build` in `frontend/`, deploys `dist/` to GitHub Pages
- [ ] `keep-alive.yml`: runs on schedule (`*/14 * * * *`), pings the Render backend health endpoint
- [ ] `.env.production` in `frontend/` with `VITE_API_URL` pointing to the Render URL
- [ ] Backend CORS middleware reads allowed origin from `ALLOWED_ORIGIN` env var (set in Render dashboard)
- [ ] Backend Dockerfile optimized: multi-stage build, non-root user, only necessary TeXLive packages
- [ ] Render configured to auto-deploy from `backend/` directory on push to `main`
- [ ] Verify: frontend build succeeds in CI, pages deploy works, keep-alive ping returns 200
