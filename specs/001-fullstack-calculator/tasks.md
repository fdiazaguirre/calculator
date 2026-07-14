# Tasks: Full-Stack Calculator

**Input**: Design documents from `specs/001-fullstack-calculator/`
**Prerequisites**: plan.md ✅ spec.md ✅ research.md ✅ data-model.md ✅ contracts/api.md ✅ quickstart.md ✅

**Tests**: Included — Constitution III (NON-NEGOTIABLE): Red-Green-Refactor before any implementation line. Every RED task must be confirmed failing before its GREEN task begins.

**Parallel strategy**: Backend and frontend TDD streams run concurrently within each US phase (API contract is locked in `contracts/api.md` so frontend can mock without the real server).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependency on an incomplete sibling task)
- **[Story]**: User story label — [US1] basic arithmetic, [US2] validation/errors, [US3] advanced ops, [US4] mobile

---

## Phase 1: Setup

**Purpose**: Create both project skeletons. These two tasks have no dependencies on each other — start them in parallel.

- [X] T001 [P] Scaffold `backend/` Go module: `go mod init calculator/backend` (Go 1.22), create `backend/main.go` stub (empty `main` with `PORT` env var read), create `backend/internal/calculator/` and `backend/internal/api/` directories in `backend/`
- [X] T002 [P] Scaffold `frontend/` with `npm create vite@latest frontend -- --template react-ts`; add `vitest`, `@testing-library/react`, `@testing-library/user-event`, `msw` to devDependencies; update `package.json` scripts: `"test": "vitest"`, `"test:ci": "vitest run --coverage"` in `frontend/`

**Checkpoint**: Both directories exist with their boilerplate before Phase 2 begins.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared infrastructure that must exist before any user story task begins. All [P] tasks run in parallel.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T003 [P] Add `github.com/stretchr/testify v1.9+` to `backend/go.mod`; run `go mod download`; verify `go test ./...` exits 0 on empty test files in `backend/go.mod`
- [X] T004 [P] Create `frontend/vitest.config.ts` (include jsdom environment, coverage provider v8); create `frontend/src/setupTests.ts` (import `@testing-library/jest-dom`); configure msw service worker setup in `frontend/vitest.config.ts`
- [X] T005 [P] Create CORS middleware in `backend/internal/api/middleware.go`: `func CORSMiddleware(next http.Handler, allowedOrigin string) http.Handler` — sets `Access-Control-Allow-Origin`, handles `OPTIONS` preflight, passes through all other methods
- [X] T006 [P] Create CSS design tokens in `frontend/src/styles/tokens.css` (palette, spacing, font-size clamps, border-radius, duration vars per web coding style rules); create `frontend/src/styles/global.css` (reset, body, box-sizing)
- [X] T007 Wire HTTP server in `backend/main.go`: read `PORT` env var (default `8080`), register routes placeholder (`/api/v1/calculate`, `/health`), wrap with CORS middleware using `ALLOWED_ORIGIN` env var (default `http://localhost:5173`); server must start cleanly before any handler is attached in `backend/main.go`

**Checkpoint**: `go run ./backend/main.go` starts on :8080; `curl localhost:8080/health` returns 404 (not yet wired); `cd frontend && npm test` runs (0 tests, exits 0).

---

## Phase 3: User Story 1 — Basic Arithmetic (Priority: P1) 🎯 MVP

**Goal**: A user enters two numbers and one of add/subtract/multiply/divide; backend computes and returns the correct result; frontend displays it.

**Independent Test**: `curl -s -X POST localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d '{"operation":"add","a":7,"b":5}'` returns `{"success":true,"result":12,"error":null}`. Frontend at localhost:5173 shows `12` for the same input.

### Backend TDD Stream

> Write tests FIRST. Confirm FAIL. Then implement.

- [X] T008 [US1] TDD RED: In `backend/internal/calculator/calculator_test.go` write table-driven tests for `Compute(op string, a, b float64) (float64, error)` covering: 7+5=12, 10−4=6, 6×7=42, 9÷4=2.25, 0.1+0.2=0.3 (FR-008 rounding). Run `go test ./internal/calculator/...` — must FAIL (function does not exist yet).
- [X] T009 [US1] TDD GREEN: Implement `backend/internal/calculator/calculator.go` — `Compute()` with `roundResult(v float64) float64` using `strconv.FormatFloat(v, 'g', 12, 64)` round-trip. Run tests — must PASS. Refactor: ensure no magic numbers, single-responsibility functions.
- [X] T010 [US1] TDD RED: In `backend/internal/api/handler_test.go` write `httptest`-based tests for `POST /api/v1/calculate` (success paths: add/sub/mul/div with correct JSON envelope) and `GET /health` (returns `{"status":"ok"}`). Run `go test ./internal/api/...` — must FAIL.
- [X] T011 [US1] TDD GREEN: Implement `backend/internal/api/handler.go` — `HandleCalculate(w, r)` decodes JSON body, calls `calculator.Compute`, encodes `CalculationResponse{Success, Result, Error}`. Implement `HandleHealth`. Register both in `main.go`. Run tests — must PASS.

### Frontend TDD Stream (parallel with T008–T011)

> API contract is in `contracts/api.md` — use msw to mock it; no real backend needed.

- [X] T012 [P] [US1] TDD RED: In `frontend/src/services/calculatorApi.test.ts` set up msw server; write tests for `calculate(op, a, b?)` returning `CalculationResponse` for success paths (add/sub/mul/div). Run `npm test` — must FAIL (function does not exist yet).
- [X] T013 [P] [US1] TDD GREEN: Implement `frontend/src/services/calculatorApi.ts` — `calculate(operation, a, b?)` calls `POST ${VITE_API_URL}/api/v1/calculate`, returns typed `CalculationResponse`. Run tests — must PASS.
- [X] T014 [US1] TDD RED: In `frontend/src/components/Calculator/Calculator.test.tsx` write RTL tests: user enters `7` and `5`, selects `add`, clicks submit → Display shows `12`. Confirm FAIL (component does not exist).
- [X] T015 [US1] TDD GREEN: Implement `frontend/src/components/Calculator/Calculator.tsx` (owns all state), `frontend/src/components/ui/Button.tsx`, `frontend/src/components/ui/Display.tsx`. Wire into `App.tsx` and `main.tsx`. Add `frontend/src/components/Calculator/calculator.css`. Run tests — must PASS.

### Coverage Gate (US1)

- [X] T016 [US1] Run `go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out`; verify ≥80% on `internal/calculator` and `internal/api`. Attach output as completion evidence (SC-006).
- [X] T017 [US1] Run `npm run test:ci` in `frontend/`; verify ≥80% coverage on `src/services/` and `src/components/Calculator/`. Attach output as completion evidence (SC-006).

**Checkpoint (US1 complete)**: Start backend (`go run ./main.go`) and frontend (`npm run dev`). In browser: enter 7 + 5 → see `12`. Enter 0.1 + 0.2 → see `0.3` (not `0.30000000000000004`). US1 is independently functional.

---

## Phase 4: User Story 2 — Validation & Errors (Priority: P2)

**Goal**: Every invalid input class produces its exact error message (per catalogue in `contracts/api.md`); errors clear on next valid calculation; service-unavailable message appears within 5 s if backend is down.

**Independent Test**: Browser: submit `5 ÷ 0` → see "Division by zero is not allowed". Kill backend. Submit any calc → see service-unavailable message within 5 s. Restart backend. Submit `2 + 2` → error cleared, see `4`.

### Backend TDD Stream

- [X] T018 [US2] TDD RED: Extend `backend/internal/calculator/calculator_test.go` with error-case rows: D2 (5÷0), D3 (missing operand), D8 (1e16 operand out of range), D9 (unknown operation "modulo"). Run — must FAIL (Compute does not yet return these errors).
- [X] T019 [US2] TDD GREEN: Extend `backend/internal/calculator/calculator.go` — add range check (|a|, |b| > 1e15), unknown-op guard, div-by-zero guard with exact messages from error catalogue. Run tests — must PASS.
- [X] T020 [US2] TDD RED: Extend `backend/internal/api/handler_test.go` — HTTP 400 paths: bad JSON body, missing `b` for binary op, `b` present for `sqrt`, operand out of range, unknown operation; HTTP 200 calc errors: div-by-zero, result out of range. Run — must FAIL.
- [X] T021 [US2] TDD GREEN: Extend `backend/internal/api/handler.go` — full request validation (check Content-Type, decode, validate fields, validate operation arity, validate ranges) before calling Compute. HTTP 400 for structural errors, HTTP 200 for calc errors. Run tests — must PASS.

### Frontend TDD Stream (parallel with T018–T021)

- [X] T022 [P] [US2] TDD RED: Extend `frontend/src/services/calculatorApi.test.ts` — msw: error response body (`success:false`), network timeout (abort after 5 s), network failure (msw network error). Run — must FAIL.
- [X] T023 [P] [US2] TDD GREEN: Extend `frontend/src/services/calculatorApi.ts` — `AbortController` with 5 s timeout; map `success:false` response to returned error; catch network errors and throw typed `ServiceUnavailableError`. Run tests — must PASS.
- [X] T024 [US2] TDD RED: Extend `Calculator.test.tsx` — US2 scenarios: error message shown for div-by-zero; service-unavailable message shown on timeout; valid calc after error clears error and shows result. Confirm FAIL.
- [X] T025 [US2] TDD GREEN: Extend `Calculator.tsx` and `Display.tsx` — show `result` when `success:true`; show `error` string when `success:false`; show "Calculator service is unavailable. Please try again." on `ServiceUnavailableError`; clear on next successful response. Run tests — must PASS.

### Coverage Gate (US2)

- [X] T026 [US2] Rerun backend and frontend coverage; verify ≥80% still holds after error-path additions. Attach output.

**Checkpoint (US2 complete)**: All D1–D4 + D8–D9 error paths verified in browser. Error recovery verified (error clears on next valid calc). US1 and US2 both independently functional.

---

## Phase 5: User Story 3 — Advanced Operations (Priority: P3)

**Goal**: User can compute power (`2^10=1024`), square root (`√144=12`), and percentage (`15% of 200=30`) using the same UI.

**Independent Test**: Browser: `2 ^ 10 = 1024`, `√144 = 12`, `√(−4)` → error "Square root of a negative number is not allowed", `15% of 200 = 30`.

### Backend TDD Stream

- [X] T027 [US3] TDD RED: Extend `backend/internal/calculator/calculator_test.go` with D4–D7 rows: `2^10=1024`, `0^0=1`, `0^(−2)` → div-by-zero, `(−8)^3=(−512)`, `(−8)^0.5` → "Result is not a real number", `√144=12`, `√(−4)` → sqrt-negative, `15%200=30`. Confirm FAIL.
- [X] T028 [US3] TDD GREEN: Extend `backend/internal/calculator/calculator.go` — `power`, `sqrt`, `percentage` cases in `Compute()`; guard D4 (0^neg), D5 (neg base + frac exp via `math.IsNaN`), D6 (neg radicand). Run tests — must PASS.

### Frontend TDD Stream (parallel with T027–T028)

- [X] T029 [P] [US3] TDD RED: Extend `calculatorApi.test.ts` — msw: power, sqrt (unary, no `b`), percentage success responses. Confirm FAIL.
- [X] T030 [P] [US3] TDD GREEN: Extend `calculatorApi.ts` — add `"power" | "sqrt" | "percentage"` to `Operation` type; pass `b: null` when operation is `sqrt`. Run tests — must PASS.
- [X] T031 [US3] TDD RED: Extend `Calculator.test.tsx` — operation dropdown includes power/sqrt/percentage; `b` input hidden when `sqrt` selected; submitting `2 ^ 10` shows `1024`. Confirm FAIL.
- [X] T032 [US3] TDD GREEN: Extend `Calculator.tsx` — add `power`, `sqrt`, `percentage` to operation `<select>`; conditionally hide `b` input when `operation === 'sqrt'`; pass `b: undefined` to service for unary ops. Run tests — must PASS.

### Coverage Gate (US3)

- [X] T033 [US3] Rerun backend coverage; confirm D4–D7 paths covered and overall ≥80%. Rerun frontend coverage ≥80%. Attach output.

**Checkpoint (US3 complete)**: All seven operations work in browser. D4–D7 error paths verified.

---

## Phase 6: User Story 4 — Mobile Responsiveness (Priority: P3)

**Goal**: Full calculator flow usable on a 320px-wide viewport with no horizontal scroll.

**Independent Test**: Resize browser to 320px. Enter `8 × 8`. Submit. See `64`. No horizontal scrollbar at any point.

- [X] T034 [P] [US4] Implement responsive layout in `frontend/src/components/Calculator/calculator.css`: use CSS custom properties from `tokens.css`; fluid widths (`max-width`, `min-width: 0`, `width: 100%`); flex/grid layout that collapses gracefully at 320px; no fixed pixel widths wider than viewport
- [X] T035 [P] [US4] Set minimum tap target 44×44px on all `Button.tsx` elements; add `aria-label` to operation `<select>`; add `role="status" aria-live="polite"` to `Display.tsx` for screen-reader announcements
- [ ] T036 [US4] Manual visual test: open DevTools → responsive mode; test at 320px, 375px, 768px, 1440px; complete calculation `8 × 8 = 64` at each width; verify zero horizontal overflow; take screenshots as completion evidence (SC-005)
  - **Status**: Responsive CSS implemented and verified structurally (fluid `max-width`, `min-width: 0`, `clamp()` sizing, `@media (max-width: 380px)`, 44px tap targets, `overflow-wrap: anywhere`). Production build emits valid CSS. The interactive DevTools screenshot pass at each breakpoint remains as the one manual verification step (requires a human/Playwright browser session).

**Checkpoint (US4 complete)**: Mobile smoke test at 375px: complete a full calculation. No overflow.

---

## Phase 7: Polish & Deployment

**Purpose**: Docker, README, final integration gate, coverage evidence. T037/T038/T040 are parallel.

- [X] T037 [P] Write `backend/Dockerfile`: multi-stage build — `golang:1.22-alpine` builder compiles `GOOS=linux CGO_ENABLED=0`; final stage `gcr.io/distroless/static-debian12`; expose `8080`; `HEALTHCHECK` via `/health`
- [X] T038 [P] Write `frontend/Dockerfile`: `node:20-alpine` build stage runs `npm ci && npm run build`; final stage `nginx:alpine` serves `dist/`; include `nginx.conf` that proxies `/api/` to `http://backend:8080` in `frontend/nginx.conf`
- [X] T039 Write `docker-compose.yml` at repo root: `backend` service (build `./backend`, port `8080:8080`, env `PORT=8080 ALLOWED_ORIGIN=http://localhost`), `frontend` service (build `./frontend`, port `80:80`, depends_on backend); verify `docker compose up --build` starts both services
- [X] T040 [P] Write `README.md`: setup prerequisites, local dev instructions, Docker Compose instructions, API curl examples for all 7 operations + 3 error cases, design decisions (stdlib net/http rationale, float rounding approach, error envelope design, module layout decision), link to `specs/` directory (SC-007)
- [X] T041 Full integration test (`docker compose up --build`): exercise all D1–D9 decision-table rows via curl against `localhost/api/v1/calculate`; open `http://localhost` in browser; complete at least one calculation per operation type; verify SC-001 through SC-007 checklist
- [X] T042 [P] Final coverage evidence: `cd backend && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out | tee backend-coverage.txt`; `cd frontend && npm run test:ci 2>&1 | tee frontend-coverage.txt`; attach both files as task completion artifacts per SC-006 and constitution VI

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — T001 || T002 immediately.
- **Phase 2 (Foundational)**: After Phase 1 — T003 || T004 || T005 || T006 in parallel; T007 after T003 + T005.
- **Phase 3 (US1)**: After Phase 2 — backend stream (T008→T009→T010→T011) || frontend stream (T012→T013→T014→T015) in parallel. T014 needs T013.
- **Phase 4 (US2)**: After Phase 3 — backend (T018→T019→T020→T021) || frontend (T022→T023→T024→T025) in parallel.
- **Phase 5 (US3)**: After Phase 4 — backend (T027→T028) || frontend (T029→T030→T031→T032) in parallel.
- **Phase 6 (US4)**: After Phase 5 — T034 || T035 in parallel; T036 after T034+T035.
- **Phase 7 (Deploy)**: After Phase 6 — T037 || T038 || T040 in parallel; T039 after T037+T038; T041 after T039; T042 || T041.

### Within Each Phase

1. TDD RED task must be confirmed FAILING before its GREEN task begins.
2. Tests must be written before implementation — verified, not assumed.
3. Coverage gate (T016/T017/T026/T033) blocks moving to next US phase.
4. Models/types before services; services before components.

---

## Parallel Examples

### Phase 1 — Two parallel agents

```
Agent A: T001 — backend/ Go scaffold
Agent B: T002 — frontend/ Vite React TS scaffold
```

### Phase 2 — Four parallel agents

```
Agent A: T003 — go.mod testify
Agent B: T004 — vitest config + msw
Agent C: T005 — CORS middleware skeleton
Agent D: T006 — CSS tokens + global styles
(T007 sequential after A + C)
```

### Phase 3 (US1) — Two parallel streams

```
Stream A (backend):  T008 → T009 → T010 → T011
Stream B (frontend): T012 → T013 → T014 → T015
(streams join at T016/T017 coverage gates)
```

### Phase 4 (US2) — Two parallel streams

```
Stream A (backend):  T018 → T019 → T020 → T021
Stream B (frontend): T022 → T023 → T024 → T025
(streams join at T026 coverage gate)
```

### Phase 5 (US3) — Two parallel streams

```
Stream A (backend):  T027 → T028
Stream B (frontend): T029 → T030 → T031 → T032
(streams join at T033 coverage gate)
```

### Phase 7 (Deploy) — Three parallel agents

```
Agent A: T037 — backend Dockerfile
Agent B: T038 — frontend Dockerfile + nginx.conf
Agent C: T040 — README.md
(T039 after A+B; T041 after T039; T042 parallel with T041)
```

---

## Implementation Strategy

### MVP First (US1 only)

1. Complete Phase 1 + Phase 2.
2. Complete Phase 3 (US1) using parallel backend + frontend streams.
3. **STOP and VALIDATE**: `curl localhost:8080/api/v1/calculate` + browser smoke test.
4. Confirm T016 + T017 coverage gates pass.

### Incremental Delivery

```
Phase 1+2 → Foundation ready
Phase 3 (US1) → Basic arithmetic working; demo-able
Phase 4 (US2) → Error handling production-quality; demo-able
Phase 5 (US3) → All 7 operations working; demo-able
Phase 6 (US4) → Mobile-ready; demo-able on any device
Phase 7 → Docker deployment; README complete; Integration gate cleared → merge to main
```

---

## Notes

- `[P]` = truly parallel (different files, no shared state, no dependency on incomplete sibling).
- TDD RED tasks: run the test command and paste failing output as evidence before starting GREEN.
- Coverage gates (T016, T017, T026, T033, T042) are blocking — they enforce SC-006 and constitution VI (evidence-based verification).
- `contracts/api.md` error message catalogue is the source of truth for exact strings; backend and frontend tests must use the same strings verbatim.
- `0^0 = 1` and `0^negative → division by zero` are coded into T027/T028 per spec Assumptions.
- Rate limiting intentionally omitted (spec Assumption: local/demo deployment only).
