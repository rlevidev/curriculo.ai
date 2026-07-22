# Tickets: Missing Components Fix

Fixes the frontend TypeScript build errors by creating the missing UI components and type definitions imported by `App.tsx`. Ref: `docs/spec-missing-components-2026-07-22/SPEC.md`.

Work the **frontier**: any ticket whose blockers are all done. For a purely linear chain that means top to bottom.

## Create Data Contracts and Shared Logic

**What to build:** The foundation for the resume application's data layer, enabling the application to strongly type the user's resume details and compute their ATS score in real-time.

**Blocked by:** None — can start immediately

- [ ] `frontend/src/types.ts` is created.
- [ ] Exports the `ResumeData` interface with all necessary personal, education, and experience fields.
- [ ] Exports the pure function `calculateATSScore(data: ResumeData)` which evaluates the resume and returns a score and criteria list.

## Create UI Component Stubs

**What to build:** The visual interface of the application, allowing the user to see the top bar, edit their resume details, and view the live preview without the application crashing on load.

**Blocked by:** Create Data Contracts and Shared Logic

- [ ] `frontend/src/components/TopBar.tsx` is created and renders the server status and ATS score.
- [ ] `frontend/src/components/EditorPane.tsx` is created and provides input fields bound to the `ResumeData` state.
- [ ] `frontend/src/components/PreviewPane.tsx` is created and renders a visual approximation of the PDF resume.
- [ ] Running `tsc --noEmit` succeeds without "Cannot find module" errors.