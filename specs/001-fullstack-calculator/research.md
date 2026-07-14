# Research: Full-Stack Calculator

**Branch**: `001-fullstack-calculator` | **Date**: 2026-07-14

All decisions resolved from GOAL.md, constitution constraints, and framework ecosystem state. No external web search required.

---

## Decision 1: Go HTTP Framework

**Decision**: `net/http` stdlib only (Go 1.22+).

**Rationale**: Go 1.22 added native path parameter routing (`{name}`), eliminating the main reason to reach for chi/gorilla/gin on small services. This feature has one endpoint plus a health check — zero basis for a framework abstraction (constitution IV YAGNI: no abstraction before second concrete use case exists).

**Alternatives considered**:
- `github.com/go-chi/chi/v5` — solid choice, but zero benefit for one handler; adds a dependency with no payoff.
- `github.com/gin-gonic/gin` — heavier; designed for APIs with many routes; overkill and adds reflection overhead for a pure-computation service.

---

## Decision 2: Go Testing

**Decision**: stdlib `testing` + `github.com/stretchr/testify/assert`.

**Rationale**: Table-driven tests with `testify/assert` are idiomatic Go. Calculator logic is pure functions — no mocks needed. `net/http/httptest` (stdlib) covers handler integration tests. `testify` is a single, stable dep for readable assertions; no need for a heavier suite.

**Alternatives considered**:
- stdlib `testing` only: verbose equality assertions; no clean diff output on failure.
- `github.com/onsi/ginkgo`: BDD-style; incompatible with the table-driven idiom that is most natural for decision-table D1–D9.

---

## Decision 3: Frontend Build Tool

**Decision**: Vite 6 + React 18 + TypeScript 5.

**Rationale**: CRA is deprecated (React team removed it from docs, 2023). Vite is the current community standard: instant HMR, native ESM, no Babel config, minimal boilerplate. Bundle budget target (< 150 kB gzipped) is achievable trivially with React + Vite with no tree-shaking effort.

**Alternatives considered**:
- Next.js: SSR/file-based routing adds complexity not warranted for a single-page calculator (YAGNI).
- Parcel: Less config, but less community momentum; Vite preferred.

---

## Decision 4: Frontend Testing

**Decision**: Vitest + React Testing Library (RTL) + `msw` for API-layer tests.

**Rationale**: Vitest shares the Vite config — no separate Babel transform, zero extra config file. RTL tests components as users do (DOM assertions, not implementation details). `msw` intercepts `fetch` at the service-worker level for realistic API integration tests without hitting the actual backend.

**Alternatives considered**:
- Jest: Requires separate transform config for ESM + TypeScript; slower cold start.
- Cypress component tests: Good for E2E but heavy for unit/integration; overkill for this scope.

---

## Decision 5: Floating-Point Precision (FR-008)

**Decision**: Backend rounds result to 12 significant digits via `strconv.FormatFloat` → `strconv.ParseFloat` before JSON serialization. No external decimal library.

**Rationale**: `0.1 + 0.2` in IEEE 754 float64 = `0.30000000000000004`. Formatting with `'g', 12` produces `"0.3"`; re-parsing gives `0.3`. This is deterministic, dependency-free, and covers the ±1e15 range (FR-009) without precision loss for values that have clean decimal representations within 12 significant digits.

```go
// Example (illustrative — not production code)
s := strconv.FormatFloat(v, 'g', 12, 64)  // "0.3"
r, _ := strconv.ParseFloat(s, 64)          // 0.3
```

**Alternatives considered**:
- `github.com/shopspring/decimal`: Correct arbitrary-precision decimal, but pulls in a dep and is heavier than needed for this numeric domain.
- `math/big.Float`: Arbitrary precision but significantly more code; YAGNI for ±1e15 with 12 sig figs.
- Frontend-side rounding: Moves spec-mandated calculation correctness out of the backend, violating constitution Architecture Constraint ("Go handles the logic").

---

## Decision 6: CORS

**Decision**: Thin CORS middleware as a Go function (~15 lines, stdlib only). No external CORS library.

**Rationale**: Requirements: allow `http://localhost:5173` (Vite dev server) + `*` in dev. Production: same-origin (Docker Compose serves both on one port). A simple `http.Handler` wrapper is sufficient; no library has enough value to justify the dependency.

---

## Decision 7: Containerisation

**Decision**: `backend/Dockerfile` + `frontend/Dockerfile` + root `docker-compose.yml`.

**Rationale**: GOAL.md requests optional Dockerfile for full-stack deployment. Docker Compose is the standard local full-stack runner; single YAML at root is the simplest layout.

Backend: multi-stage Go build (builder + `gcr.io/distroless/static`).
Frontend: `node:alpine` build stage + `nginx:alpine` serve stage.
Compose: backend on `:8080`, frontend nginx on `:80` (proxies `/api/` to backend).

---

## Decision 8: Module / Monorepo Layout

**Decision**: `backend/` (Go module) + `frontend/` (npm package), both in repo root. No monorepo tooling.

**Rationale**: Two independent runtimes with no shared code — a monorepo tool (Turborepo, Nx) would only add config. Each directory is self-contained. `docker-compose.yml` at root ties them together.
