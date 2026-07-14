# calculator Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-07-14

## Active Technologies

- Go 1.22+ (backend), TypeScript 5 + React 18 (frontend) (001-fullstack-calculator)

## Project Structure

```text
backend/
├── main.go
├── go.mod
├── internal/
│   ├── calculator/   # pure math logic + tests
│   └── api/          # HTTP handlers + tests
└── Dockerfile

frontend/
├── src/
│   ├── components/Calculator/
│   ├── services/
│   └── styles/
└── Dockerfile

docker-compose.yml
README.md
specs/001-fullstack-calculator/  # spec, plan, research, data-model, contracts
```

## Commands

```bash
# Backend
cd backend && go test ./... -v -count=1          # run tests
cd backend && go run ./main.go                   # start (PORT=8080)

# Frontend
cd frontend && npm install && npm run dev         # start (localhost:5173)
cd frontend && npm run test:ci                    # tests + coverage

# Full stack
docker compose up --build                        # both services
```

## Code Style

- Backend: idiomatic Go; no framework (stdlib net/http only); table-driven tests; `fmt.Errorf` for errors; no custom error hierarchy.
- Frontend: React functional components + TypeScript; Vitest + RTL for tests; CSS custom properties for tokens; no CSS-in-JS.
- Constitution IV YAGNI: no interface before second concrete use case; no abstraction for a single handler.

## Recent Changes

- 001-fullstack-calculator: Added Go 1.22+ (backend), TypeScript 5 + React 18 (frontend)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
