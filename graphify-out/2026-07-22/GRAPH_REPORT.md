# Graph Report - curriculo.ai  (2026-07-22)

## Corpus Check
- 34 files · ~18,614 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 235 nodes · 234 edges · 27 communities (16 shown, 11 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 6 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `5a91eeda`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- devDependencies
- compilerOptions
- compilerOptions
- main.go
- Implementation Decisions
- Possíveis Problemas
- Tickets: curriculo.ai MVP
- Spec: curriculo.ai MVP
- System Design
- App.tsx
- Inputs necessários
- Criador de Currículo
- .oxlintrc.json
- React + TypeScript + Vite
- tsconfig.json
- CLAUDE.md
- Exemplo de front-end 3a3e3b528f7380eea019e28d31ab28cc.md
- Template de Currículo 3a1e3b528f7380d28276e0241a56ee75.md
- README.md
- github.com/rlevidev/curriculo.ai
- patch-preview.js
- update-css.js
- update-editor.js
- Spec: curriculo.ai MVP
- SPEC.md
- Tickets: Missing Components Fix

## God Nodes (most connected - your core abstractions)
1. `compilerOptions` - 18 edges
2. `compilerOptions` - 15 edges
3. `Implementation Decisions` - 14 edges
4. `Possíveis Problemas` - 13 edges
5. `Tickets: curriculo.ai MVP` - 10 edges
6. `System Design` - 10 edges
7. `Inputs necessários` - 9 edges
8. `Spec: curriculo.ai MVP` - 8 edges
9. `rateLimiter()` - 7 edges
10. `ResumeData` - 6 edges

## Surprising Connections (you probably didn't know these)
- `TestGeneratePdfHandler_MissingFields()` --calls--> `rateLimiter()`  [INFERRED]
  backend/main_test.go → backend/main.go
- `TestGeneratePdfHandler_Success()` --calls--> `rateLimiter()`  [INFERRED]
  backend/main_test.go → backend/main.go
- `TestRateLimiter()` --calls--> `rateLimiter()`  [INFERRED]
  backend/main_test.go → backend/main.go
- `TestRateLimiter_SameIPAfterExhaustion()` --calls--> `rateLimiter()`  [INFERRED]
  backend/main_test.go → backend/main.go
- `TestTexEscape()` --calls--> `texEscape()`  [INFERRED]
  backend/main_test.go → backend/main.go

## Import Cycles
- None detected.

## Communities (27 total, 11 thin omitted)

### Community 0 - "devDependencies"
Cohesion: 0.10
Nodes (20): dependencies, react, react-dom, devDependencies, oxlint, @types/node, @types/react, @types/react-dom (+12 more)

### Community 1 - "compilerOptions"
Cohesion: 0.10
Nodes (19): compilerOptions, allowArbitraryExtensions, allowImportingTsExtensions, erasableSyntaxOnly, jsx, lib, module, moduleDetection (+11 more)

### Community 2 - "compilerOptions"
Cohesion: 0.12
Nodes (16): compilerOptions, allowImportingTsExtensions, erasableSyntaxOnly, lib, module, moduleDetection, noEmit, noFallthroughCasesInSwitch (+8 more)

### Community 3 - "main.go"
Cohesion: 0.14
Nodes (22): Education, Experience, Language, generatePdfHandler(), main(), rateLimiter(), TestGeneratePdfHandler_MissingFields(), TestGeneratePdfHandler_Success() (+14 more)

### Community 4 - "Implementation Decisions"
Cohesion: 0.14
Nodes (14): Architecture, ATS Score, Cold Start Handling, Concurrency Control, Domain Model — JSON Schema (API Contract), Error Handling, Frontend Details, Implementation Decisions (+6 more)

### Community 5 - "Possíveis Problemas"
Cohesion: 0.14
Nodes (13): Bundle do TeXLive pesado., Cold start do Render., CORS mal configurado., Educação/Experiência fixas no frontend., Maior risco de todos: abandono do projeto., Para evitar Bundle do TeXLive pesado esses são os pacotes do Dockerfile necessário para criar o currículo seguindo o template:, pdflatex trava sem erro claro., Possíveis Problemas (+5 more)

### Community 6 - "Tickets: curriculo.ai MVP"
Cohesion: 0.17
Nodes (10): 1. Monorepo scaffolding, 2. Header tracer bullet: form → preview → PDF, 3. Dynamic sections: Education + Skills, 4. Dynamic sections: Experience + Projects, 5. Simple sections: Languages + Certifications + auto-disappear, 6. Input validation + security hardening, 7. ATS score + page overflow warning, 8. localStorage draft + cold start UX (+2 more)

### Community 8 - "System Design"
Cohesion: 0.10
Nodes (19): Botão de Download, Cabeçalho:, Certificações (Opcional: Sim):, Componentes:, Educação (Opcional: Sim):, Exemplo de Request JSON:, Experiência (Opcional: Sim):, **Fluxo de dados**: (+11 more)

### Community 9 - "App.tsx"
Cohesion: 0.12
Nodes (15): App(), defaultResumeData, EditorPane(), EditorPaneProps, PreviewPaneProps, TopBarProps, ATSCriteria, ATSResult (+7 more)

### Community 11 - "Criador de Currículo"
Cohesion: 0.25
Nodes (7): A ideia é criar um site que crie um currículo para o usuário em formato ATS:, Como devem ser feitos os currículos:, Criador de Currículo, Exemplos:, Nome do projeto: Curriculo.ai, Problemática a ser resolvida:, Tópicos que devem ter:

### Community 12 - ".oxlintrc.json"
Cohesion: 0.33
Nodes (5): plugins, rules, react/only-export-components, react/rules-of-hooks, $schema

### Community 13 - "React + TypeScript + Vite"
Cohesion: 0.50
Nodes (3): Expanding the Oxlint configuration, React Compiler, React + TypeScript + Vite

### Community 24 - "Spec: curriculo.ai MVP"
Cohesion: 0.17
Nodes (12): Further Notes, Out of Scope, Prior Art, Problem Statement, Seam 1 — API Contract (Integration), Seam 2 — Resume State → Preview (Unit, Frontend), Seam 3 — Go Handler → LaTeX Renderer (Unit, Backend), Seam 4 — LaTeX → pdflatex (Integration, Backend) (+4 more)

### Community 25 - "SPEC.md"
Cohesion: 0.25
Nodes (7): Further Notes, Implementation Decisions, Out of Scope, Problem Statement, Solution, Testing Decisions, User Stories

### Community 26 - "Tickets: Missing Components Fix"
Cohesion: 0.50
Nodes (3): Create Data Contracts and Shared Logic, Create UI Component Stubs, Tickets: Missing Components Fix

## Knowledge Gaps
- **157 isolated node(s):** `github.com/rlevidev/curriculo.ai`, `fs`, `css`, `$schema`, `plugins` (+152 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **11 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Spec: curriculo.ai MVP` connect `Spec: curriculo.ai MVP` to `Implementation Decisions`, `Tickets: curriculo.ai MVP`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Why does `Implementation Decisions` connect `Implementation Decisions` to `Spec: curriculo.ai MVP`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **What connects `github.com/rlevidev/curriculo.ai`, `fs`, `css` to the rest of the system?**
  _157 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `devDependencies` be split into smaller, more focused modules?**
  _Cohesion score 0.09523809523809523 - nodes in this community are weakly interconnected._
- **Should `compilerOptions` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._
- **Should `compilerOptions` be split into smaller, more focused modules?**
  _Cohesion score 0.11764705882352941 - nodes in this community are weakly interconnected._
- **Should `main.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14130434782608695 - nodes in this community are weakly interconnected._