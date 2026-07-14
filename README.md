# Full-Stack Calculator

A calculator with a **Go** REST backend and a **React + TypeScript** frontend.
All arithmetic runs server-side; the frontend is a pure presentation layer.

- Backend: Go standard library (`net/http`), no web framework.
- Frontend: React 19 + Vite + TypeScript.
- Stateless: no database, no persistence, no authentication.

Built with [Spec-Driven Development](specs/001-fullstack-calculator/) — the full
spec, plan, contracts, and task list live under `specs/`.

---

## Features

| Operation | Identifier | Arity | Example |
|-----------|-----------|-------|---------|
| Addition | `add` | 2 | `7 + 5 = 12` |
| Subtraction | `subtract` | 2 | `10 − 4 = 6` |
| Multiplication | `multiply` | 2 | `6 × 7 = 42` |
| Division | `divide` | 2 | `9 ÷ 4 = 2.25` |
| Exponentiation | `power` | 2 | `2 ^ 10 = 1024` |
| Square root | `sqrt` | 1 | `√144 = 12` |
| Percentage | `percentage` | 2 | `15% of 200 = 30` |

- Input validation with specific, human-readable error messages.
- Division-by-zero, negative-square-root, and out-of-range guards.
- Exact decimal display up to 12 significant digits (`0.1 + 0.2` shows `0.3`, not `0.30000000000000004`).
- Responsive down to 320px; keyboard-accessible; screen-reader live region.

---

## Prerequisites

| Tool | Version |
|------|---------|
| Go | 1.22+ |
| Node.js | 20 LTS+ |
| Docker + Compose | recent (optional, for the containerised path) |

---

## Run locally (two terminals)

### Backend

```bash
cd backend
go run ./main.go
# calculator backend listening on :8080
```

### Frontend

```bash
cd frontend
npm install
npm run dev
# Vite dev server on http://localhost:5173
```

Open <http://localhost:5173>.

Environment variables:

| Service | Variable | Default | Purpose |
|---------|----------|---------|---------|
| backend | `PORT` | `8080` | HTTP listen port |
| backend | `ALLOWED_ORIGIN` | `http://localhost:5173` | CORS origin for dev |
| frontend | `VITE_API_URL` | `http://localhost:8080` | Backend base URL (empty = same-origin) |

---

## Run with Docker Compose (one command)

```bash
docker compose up --build
# Frontend: http://localhost
# API (proxied by nginx): http://localhost/api/v1/calculate
```

The frontend nginx container serves the built SPA and proxies `/api/` to the
backend container. Stop with `docker compose down`.

---

## API Reference

### `POST /api/v1/calculate`

Request body:

```json
{ "operation": "add", "a": 7, "b": 5 }
```

- `operation` — one of the seven identifiers in the table above.
- `a` — first operand (required, |a| ≤ 1e15).
- `b` — second operand (required for binary ops; omit or `null` for `sqrt`).

Every response uses the same envelope. Exactly one of `result` / `error` is non-null:

```json
{ "success": true, "result": 12, "error": null }
```

**Calculation-domain errors** (division by zero, negative root, out-of-range
result) return **HTTP 200** with `success: false` — the request was well-formed,
the math failed:

```json
{ "success": false, "result": null, "error": "Division by zero is not allowed" }
```

**Structural errors** (bad JSON, unknown operation, missing operand, operand out
of range) return **HTTP 400** with the same envelope shape.

#### Examples

```bash
# Addition
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":7,"b":5}'
# {"success":true,"result":12,"error":null}

# Division
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":9,"b":4}'
# {"success":true,"result":2.25,"error":null}

# Floating point stays exact
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":0.1,"b":0.2}'
# {"success":true,"result":0.3,"error":null}

# Exponentiation
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"power","a":2,"b":10}'
# {"success":true,"result":1024,"error":null}

# Square root (unary — omit b)
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":144}'
# {"success":true,"result":12,"error":null}

# Percentage: 15% of 200
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"percentage","a":15,"b":200}'
# {"success":true,"result":30,"error":null}

# Division by zero (HTTP 200, success:false)
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":5,"b":0}'
# {"success":false,"result":null,"error":"Division by zero is not allowed"}

# Negative square root
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":-4}'
# {"success":false,"result":null,"error":"Square root of a negative number is not allowed"}

# Unknown operation (HTTP 400)
curl -s -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"modulo","a":10,"b":3}'
# {"success":false,"result":null,"error":"Unsupported operation 'modulo'. Supported: add, subtract, multiply, divide, power, sqrt, percentage"}
```

### `GET /health`

```bash
curl -s localhost:8080/health
# {"status":"ok"}
```

---

## Testing

### Backend

```bash
cd backend
go test ./... -cover
```

Table-driven tests cover every operation and every error case in the spec's
decision table. Coverage: `internal/calculator` ~90%, `internal/api` ~96%.

### Frontend

```bash
cd frontend
npm run test:ci
```

Vitest + React Testing Library + MSW (mocked backend). Coverage ~100% statements
on the service and calculator component.

---

## Design Decisions

- **Backend does all math.** The frontend never computes a result. This keeps a
  single source of truth for correctness and matches the project constitution's
  language-boundary rule (Go = logic, React = UI).
- **No Go web framework.** With one calculation endpoint plus a health check,
  Go 1.22's method-qualified `net/http` routing is enough. A framework would add
  a dependency with no payoff (YAGNI).
- **Uniform response envelope** (`success` / `result` / `error`). The frontend
  has one code path for every outcome instead of branching on HTTP status codes.
- **Calculation errors are HTTP 200, not 4xx.** A division-by-zero request is
  syntactically valid — the HTTP transaction succeeded; the *math* failed. Only
  malformed requests get 4xx.
- **Float noise stripped server-side** by formatting to 12 significant digits and
  re-parsing (`strconv.FormatFloat(v, 'g', 12, 64)`), so `0.1 + 0.2` → `0.3`. No
  decimal library needed for the ±1e15 domain.
- **Timeout + service-unavailable handling.** The frontend aborts a request after
  5s and shows a clear "service is unavailable" message rather than hanging.

### Assumptions

- Percentage means `a% of b = a × b ÷ 100`.
- `0 ^ 0 = 1`; `0 ^ negative` is division by zero.
- Numeric domain is real decimals within ±1e15.
- No history, accounts, or persistence — single anonymous session.
- Rate limiting is out of scope (local/demo deployment).

---

## Project Structure

```text
backend/
├── main.go                      # server bootstrap (PORT, CORS, routes)
├── internal/calculator/         # pure arithmetic + validation (no HTTP)
└── internal/api/                # HTTP handlers, CORS middleware, routing
frontend/
├── src/services/calculatorApi.ts        # typed fetch client
├── src/components/Calculator/           # main UI container
└── src/components/ui/                   # Button, Display
docker-compose.yml               # backend + nginx-served frontend
specs/001-fullstack-calculator/  # spec, plan, contracts, tasks
```
