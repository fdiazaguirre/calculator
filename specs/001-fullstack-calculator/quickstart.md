# Quickstart: Full-Stack Calculator

**Branch**: `001-fullstack-calculator` | **Date**: 2026-07-14

Two ways to run: **local development** (two terminals) or **Docker Compose** (one command).

---

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.22+ | https://go.dev/dl/ |
| Node.js | 20 LTS+ | https://nodejs.org |
| Docker + Compose | any recent | https://docs.docker.com/get-docker/ |

---

## Local Development

### 1. Start the backend

```bash
cd backend
go mod download
go run ./main.go
# Listening on :8080
```

Backend URL: `http://localhost:8080`

Verify:
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### 2. Start the frontend

```bash
cd frontend
npm install
npm run dev
# Vite dev server on http://localhost:5173
```

Open `http://localhost:5173` in a browser.

### 3. Quick smoke test

```bash
curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","a":7,"b":5}' | jq .
# { "success": true, "result": 12, "error": null }

curl -s -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"divide","a":5,"b":0}' | jq .
# { "success": false, "result": null, "error": "Division by zero is not allowed" }
```

---

## Docker Compose (full stack)

```bash
docker compose up --build
# Frontend: http://localhost
# Backend: http://localhost/api/v1/calculate (proxied by nginx)
```

To stop:
```bash
docker compose down
```

---

## Running Tests

### Backend

```bash
cd backend
go test ./... -v -count=1
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out
```

Coverage must be ≥ 80% on `internal/calculator/` and `internal/api/` (SC-006).

### Frontend

```bash
cd frontend
npm test              # watch mode
npm run test:ci       # single run + coverage report
```

Coverage must be ≥ 80% on `src/services/` and `src/components/Calculator/` (SC-006).

---

## Environment Variables

### Backend

| Variable | Default | Purpose |
|----------|---------|---------|
| `PORT` | `8080` | HTTP listen port |
| `ALLOWED_ORIGIN` | `http://localhost:5173` | CORS allowed origin (dev) |

### Frontend

| Variable | Default | Purpose |
|----------|---------|---------|
| `VITE_API_URL` | `http://localhost:8080` | Backend base URL |

In Docker Compose both are set automatically via `docker-compose.yml`; no `.env` file needed.
