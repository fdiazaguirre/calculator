# Implementation Plan: Full-Stack Calculator

**Branch**: `001-fullstack-calculator` | **Date**: 2026-07-14 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/001-fullstack-calculator/spec.md`

## Summary

Build a full-stack calculator: a Go REST API performing arithmetic (add, subtract, multiply, divide) and optional advanced operations (power, sqrt, percentage), with a React TypeScript frontend that calls the API and presents results or structured error messages. All calculation logic lives in the backend; the frontend is a pure presentation layer (constitution Architecture Constraint).

## Technical Context

**Language/Version**: Go 1.22+ (backend), TypeScript 5 + React 18 (frontend)
**Primary Dependencies**:
- Backend: `net/http` (stdlib), `encoding/json` (stdlib), `github.com/stretchr/testify` (test assertions only)
- Frontend: React 18, Vite 6, React Testing Library, Vitest, msw
**Storage**: None (stateless; no DB, no files)
**Testing**: Backend — stdlib `testing` + `testify/assert` + `httptest`; Frontend — Vitest + RTL + msw
**Target Platform**: Local server (Linux/macOS); Docker Compose for full-stack deployment
**Project Type**: Web service (backend) + single-page application (frontend)
**Performance Goals**: Backend p99 < 100 ms (pure computation; trivially achievable); frontend outcome visible within 2 s of submission (SC-004)
**Constraints**: Frontend renders without horizontal scroll from 320px to 1920px (FR-013); CORS allow `localhost:5173` in dev; backend `PORT` configurable via env var
**Scale/Scope**: Single-user demo; no concurrency requirements beyond standard Go net/http defaults

## Constitution Check

*Complexity score: 3 → Standard tier (Spec Builder → TDD gate)*

| Principle | Status | Evidence |
|-----------|--------|---------|
| I — Spec-Driven | ✅ | spec.md exists, complexity scored (3), decision table D1–D9 exhaustive, no vague words |
| II — Deep Modules | ✅ | `internal/calculator` exposes one function; HTTP handler is a thin adapter; CORS is a middleware wrapper |
| III — Test-First (NON-NEGOTIABLE) | ✅ | TDD gate required; every task includes RED step before implementation |
| IV — Minimal Abstractions | ✅ | stdlib `net/http` only (no framework); no repository pattern (no storage); no custom error hierarchy; plain `fmt.Errorf` |
| VI — Evidence-Based Verification | ✅ | SC-006 mandates ≥ 80% coverage with coverage report as task completion artifact |

**Gates applicable:**
- Discovery Gate: Passed (spec complete, no unresolved unknowns).
- TDD Gate: Required for `internal/calculator`, `internal/api`, `src/services/calculatorApi`, `src/components/Calculator`. Blocks task completion.
- Integration Gate: Required before merge to `main`. End-to-end: frontend calls live backend, all D1–D9 rows pass.

**Complexity violations**: None. No abstractions introduced beyond what one present use case requires.

## Project Structure

### Documentation (this feature)

```text
specs/001-fullstack-calculator/
├── plan.md              ← this file
├── research.md          ← Phase 0 output
├── data-model.md        ← Phase 1 output
├── quickstart.md        ← Phase 1 output
├── contracts/
│   └── api.md           ← Phase 1 output
└── tasks.md             ← Phase 2 output (/speckit.tasks — NOT created here)
```

### Source Code (repository root)

```text
backend/
├── main.go                    # wire up server; read PORT env var
├── go.mod                     # module: calculator/backend; Go 1.22
├── go.sum
├── internal/
│   ├── calculator/
│   │   ├── calculator.go      # Compute(op, a, b) → (float64, error); pure functions
│   │   └── calculator_test.go # table-driven; covers D1–D9 happy + negative cases
│   └── api/
│       ├── handler.go         # POST /api/v1/calculate; GET /health; CORS middleware
│       └── handler_test.go    # httptest-based; covers HTTP 200/400/405/500
└── Dockerfile                 # multi-stage: go builder → distroless/static

frontend/
├── package.json
├── vite.config.ts
├── tsconfig.json
├── index.html
├── src/
│   ├── main.tsx
│   ├── App.tsx
│   ├── components/
│   │   ├── Calculator/
│   │   │   ├── Calculator.tsx      # main container; owns all state
│   │   │   ├── Calculator.test.tsx # RTL: user events; msw for API mock
│   │   │   └── calculator.css
│   │   └── ui/
│   │       ├── Button.tsx          # reusable accessible button
│   │       └── Display.tsx         # result/error display area
│   ├── services/
│   │   ├── calculatorApi.ts        # fetch wrapper; timeout 5 s (FR-012)
│   │   └── calculatorApi.test.ts   # msw; tests all response shapes
│   └── styles/
│       ├── tokens.css              # CSS custom properties (palette, spacing, type)
│       └── global.css
└── Dockerfile                      # node build → nginx:alpine

docker-compose.yml                  # backend :8080; frontend nginx :80 (proxies /api/ → backend)
README.md
```

**Structure decision**: Web application layout (Option 2 from template). Backend and frontend are independent runtimes with no shared code, so a flat two-directory layout without monorepo tooling is sufficient (research Decision 8).

## Implementation Build Order

The tasks file (`/speckit.tasks`) will sequence these phases:

### Phase A — Backend Core (P1, unblocked)

1. **A1** Scaffold `backend/` Go module; `main.go` with configurable port.
2. **A2** TDD `internal/calculator/calculator.go`: write failing tests for D1–D4 (add/sub/mul/div/power), then implement `Compute`.
3. **A3** Extend TDD for D5–D9 (sqrt, percentage, edge cases, range, unknown op).
4. **A4** TDD `internal/api/handler.go`: write failing handler tests, then implement HTTP handler + CORS middleware.
5. **A5** Run `go test ./... -coverprofile=coverage.out`; verify ≥ 80%.

### Phase B — Frontend Core (P1, after A4 API contract locked)

1. **B1** Scaffold `frontend/` Vite + React TS; install RTL, Vitest, msw.
2. **B2** TDD `calculatorApi.ts`: write failing service tests (msw), then implement fetch + timeout.
3. **B3** TDD `Calculator.tsx`: write failing component tests (user-event), then implement UI.
4. **B4** Add CSS tokens + responsive layout (FR-013: 320–1920px).
5. **B5** Run `npm run test:ci`; verify ≥ 80% on service + component.

### Phase C — Integration & Docs (after A5 + B5)

1. **C1** Manual smoke test: start backend + `npm run dev`; exercise all D1–D9 rows in browser.
2. **C2** Write `backend/Dockerfile`, `frontend/Dockerfile`, `docker-compose.yml`.
3. **C3** Write `README.md` (setup instructions, API examples, design decisions).
4. **C4** Integration gate: `docker compose up --build`; repeat D1–D9 smoke test against composed stack.

## Complexity Tracking

> No violations to justify. All constraints remain within Standard-tier baseline.
