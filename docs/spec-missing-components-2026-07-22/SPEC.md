## Problem Statement

The frontend application build fails during the TypeScript compilation phase before reaching Vite. This occurs because `App.tsx` imports four modules (`./components/TopBar`, `./components/EditorPane`, `./components/PreviewPane`, and `./types`) that have not yet been created in the repository.

## Solution

Create the missing TypeScript and TSX files as functional stubs to satisfy the imports in `App.tsx`. This will resolve the module resolution errors, allow the TypeScript compiler to pass, and enable the Vite development server and build process to complete successfully.

## User Stories

1. As a developer, I want the TypeScript build to pass, so that I can run the local development server without fatal errors.
2. As a user, I want to see the TopBar component, so that I can view the current server status and my real-time ATS score.
3. As a user, I want to interact with the EditorPane component, so that I can input my personal details, experiences, and skills into the resume form.
4. As a user, I want to view the PreviewPane component, so that I can see a real-time approximation of how my final PDF resume will look.
5. As a user, I want to click the "Exportar PDF" button in the TopBar, so that I can download my finished resume.
6. As a developer, I want a shared `types.ts` file, so that the data contract for `ResumeData` is strongly typed across all components.

## Implementation Decisions

- **Modules built/modified:**
  - `frontend/src/types.ts`: Will contain the central `ResumeData` interface and the `calculateATSScore` business logic.
  - `frontend/src/components/TopBar.tsx`: Will render the application header, ATS score indicator, server status, and export controls.
  - `frontend/src/components/EditorPane.tsx`: Will render the input fields (Name, Email, LinkedIn, etc.) using controlled components bound to the `ResumeData` state.
  - `frontend/src/components/PreviewPane.tsx`: Will render the visual representation of the resume, including handling overflow warnings and mobile view states.
- **Architectural decisions:**
  - The ATS scoring logic (`calculateATSScore`) is pure and extracted to `types.ts` to easily share it between the UI and any potential future backend validations.
  - Components are stateless (dumb) and rely on props (like `resumeData`, `onChange`) passed down from the parent `App.tsx`, following a unidirectional data flow.

## Testing Decisions

- **Good tests for this feature:** External behavior testing here primarily revolves around successful compilation and component mounting.
- **Modules to be tested:** 
  - Compilation tests (implicitly handled by `tsc` during the build step).
  - Unit tests for the pure function `calculateATSScore` in `types.ts` to ensure criteria calculations add up correctly.
  - Smoke tests for `TopBar`, `EditorPane`, and `PreviewPane` to ensure they mount without crashing when provided with default/empty `ResumeData` props.
- **Seams:** The primary seam is the TypeScript compiler interface (running `tsc --noEmit`). The secondary seam is the React component props interface—components can be tested in isolation by passing mock `ResumeData`.

## Out of Scope

- Implementing the actual PDF generation backend logic (the UI will just show the "Gerando..." loading state stub).
- Deep integration with the Go/Java backend beyond displaying the simple `serverStatus` prop.
- Advanced styling or CSS animations (relying on existing CSS classes).

## Further Notes

Once these files are committed and pushed to the `main` branch, the `deploy-frontend` workflow should naturally recover and pass, as the missing dependencies are the sole cause of the current build breakage.
