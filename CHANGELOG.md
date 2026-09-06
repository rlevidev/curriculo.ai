# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- **Structured JSON logs:** slog stdlib with UTC timestamps and `service=curriculo-api` for production log aggregation
- **Request ID tracking:** generate unique ID via `crypto/rand`, propagate via `X-Request-ID` header and all log entries
- **Request metadata logging:** method, path, status code, and latency in milliseconds for every HTTP request
- **Error context logging:** validation failures (400), rate limits (429), and PDF generation errors (500) with operation tags
- **Startup log:** server port and startup event on boot

## [1.0.0] - 2026-09-04

### Security

- **LaTeX injection:** escape all special characters (`{ } \ & % $ # _ ~ ^`) in user input before PDF compilation, preventing file content leakage and DoS via malformed LaTeX payloads
- **Rate limiter:** key by IP address instead of IP:port, so opening new connections no longer bypasses the limit; tokens replenish after a 1-hour window; periodic cleanup prevents memory leak from stale entries
- **CORS origin:** configurable via environment variable instead of hardcoded
- **Payload size enforcement:** use `MaxBytesReader` to reject chunked requests over 50KB regardless of `ContentLength` header
- **Server timeouts:** add `ReadHeaderTimeout: 5s` to `http.Server` to mitigate Slowloris attacks

### Added

- **MVP implementation:** complete resume builder with editor/preview panes, PDF export, and ATS optimization
- **Main application layout:** editor and preview components with responsive design
- **Missing components and types:** additional UI components and TypeScript types for full functionality
- **Unit tests:** backend tests for texEscape function and rate limiter, frontend tests for EditorPane, Export, and normalize utility
- **CI/CD:** GitHub Actions workflow for continuous integration and deployment
- **GitHub templates:** issue and PR templates for consistent contribution workflow
- **Environment variables:** `VITE_API_URL` for frontend build, server port configuration via env
- **CORS middleware:** HTTP handlers with configurable CORS support

### Fixed

- **PDF template completeness:** add Education, Projects, Skills (languages + technologies), LinkedIn, GitHub, and Location sections with visible separators (`·` and `\item`) so the generated resume matches the editor
- **Skills input:** preserve commas during typing by keeping local string state and committing to array on blur, preventing skills from merging unexpectedly
- **Draft migration:** normalize old localStorage drafts field-by-field against defaults, version the storage key (`resume-draft-v2`) to prevent crashes from changed data structures
- **Export error handling:** check `content-type` before parsing response JSON, show a friendly error message when the server returns HTML (e.g. 502 Bad Gateway)
- **Test isolation:** add `resetVisitors()` helper to prevent global state leakage between rate limiter tests
- **Deploy workflow permissions:** add `pages: write` permission to fix 403 error on frontend deployment
- **Vite configuration:** set base path for correct asset loading
- **CORS origin via env, pdflatex timeout context, export gating, ATS inline**
- **Overflow warning:** remove dead state and effect, restore useCallback
- **Coderabbit suggestions:** fix resume skills state crash

### Changed

- **Dockerfile multi-stage:** separate Go compilation (`golang:1.24-alpine`) from runtime (`texlive`), reducing image size and cold start time; remove personal PATH entries for reproducible builds
- **Go APIs:** replace deprecated `ioutil.TempDir`, `ioutil.ReadFile`, `ioutil.WriteFile` with `os` equivalents
- **CI workflow:** improve deploy concurrency
- **UI cleanup:** remove dead "Templates"/"Histórico" tabs, unused `previewRef`, and 10-second export timeout
- **Lint:** fix `exhaustive-deps` warning by isolating health check interval in its own effect

### Removed

- **Dead code:** remove unused imports, dead overflowWarning state, and phantom go.mod file
- **Local-only files:** move local-only files out of repo, remove node_modules/.vite cache and AVALIACAO.md from index
