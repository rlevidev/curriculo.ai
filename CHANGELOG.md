# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Security

- **LaTeX injection:** escape all special characters (`{ } \ & % $ # _ ~ ^`) in user input before PDF compilation, preventing file content leakage and DoS via malformed LaTeX payloads
- **Rate limiter:** key by IP address instead of IP:port, so opening new connections no longer bypasses the limit; tokens replenish after a 1-hour window; periodic cleanup prevents memory leak from stale entries

### Fixed

- **PDF template completeness:** add Education, Projects, Skills (languages + technologies), LinkedIn, GitHub, and Location sections with visible separators (`·` and `\item`) so the generated resume matches the editor
- **Skills input:** preserve commas during typing by keeping local string state and committing to array on blur, preventing skills from merging unexpectedly
- **Draft migration:** normalize old localStorage drafts field-by-field against defaults, version the storage key (`resume-draft-v2`) to prevent crashes from changed data structures
- **Payload size enforcement:** use `MaxBytesReader` to reject chunked requests over 50KB regardless of `ContentLength` header
- **Export error handling:** check `content-type` before parsing response JSON, show a friendly error message when the server returns HTML (e.g. 502 Bad Gateway)
- **Test isolation:** add `resetVisitors()` helper to prevent global state leakage between rate limiter tests

### Changed

- **Dockerfile multi-stage:** separate Go compilation (`golang:1.24-alpine`) from runtime (`texlive`), reducing image size and cold start time; remove personal PATH entries for reproducible builds
- **Server timeouts:** add `ReadHeaderTimeout: 5s` to `http.Server` to mitigate Slowloris attacks
- **Go APIs:** replace deprecated `ioutil.TempDir`, `ioutil.ReadFile`, `ioutil.WriteFile` with `os` equivalents
- **Lint:** fix `exhaustive-deps` warning by isolating health check interval in its own effect
- **UI cleanup:** remove dead "Templates"/"Histórico" tabs, unused `previewRef`, and 10-second export timeout
