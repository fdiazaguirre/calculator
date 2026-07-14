# REST API Contract: Calculator Backend

**Version**: v1 | **Branch**: `001-fullstack-calculator` | **Date**: 2026-07-14

Base URL (local): `http://localhost:8080`

---

## Endpoints

### POST /api/v1/calculate

Perform a calculation. Returns a structured response for every outcome — including business-logic errors — so the frontend never needs to parse HTTP status codes for result display.

#### Request

```
POST /api/v1/calculate
Content-Type: application/json
```

```json
{
  "operation": "<string>",
  "a": <number>,
  "b": <number | null>
}
```

| Field | Type | Required | Valid values |
|-------|------|----------|-------------|
| `operation` | string | yes | `"add"`, `"subtract"`, `"multiply"`, `"divide"`, `"power"`, `"sqrt"`, `"percentage"` |
| `a` | number | yes | finite float, absolute value ≤ 1e15 |
| `b` | number \| null | conditional | finite float, absolute value ≤ 1e15; **required** for all operations except `sqrt`; **must be omitted or null** for `sqrt` |

#### Responses

**HTTP 200 — calculation succeeded**

```json
{
  "success": true,
  "result": 2.25,
  "error": null
}
```

**HTTP 200 — calculation-level error** (e.g., division by zero, sqrt of negative, result out of range)

```json
{
  "success": false,
  "result": null,
  "error": "Division by zero is not allowed"
}
```

Calculation errors are **not** HTTP 4xx — the HTTP transaction itself succeeded; the calculation is what failed. This keeps the frontend error-handling path uniform.

**HTTP 400 — malformed request** (invalid JSON, missing required field, unknown operation, operand out of range)

```json
{
  "success": false,
  "result": null,
  "error": "Operand 'a' is out of the supported range (±1e15)"
}
```

**HTTP 405 — wrong HTTP method**

Body: `Method Not Allowed` (plain text). Frontend must not call with GET/PUT/etc.

**HTTP 500 — unexpected server error** (should never occur during normal operation)

```json
{
  "success": false,
  "result": null,
  "error": "Internal server error"
}
```

#### Error message catalogue

| Trigger | HTTP status | `error` value |
|---------|------------|---------------|
| Division by zero; 0 ^ negative | 200 | `"Division by zero is not allowed"` |
| sqrt of negative | 200 | `"Square root of a negative number is not allowed"` |
| Negative base with fractional exponent | 200 | `"Result is not a real number"` |
| Result magnitude > 1e15 | 200 | `"Result out of range"` |
| Missing or null required field | 400 | `"Field '<name>' is required"` |
| Non-numeric operand | 400 | `"Operand '<a\|b>' must be a number"` |
| Operand magnitude > 1e15 | 400 | `"Operand '<a\|b>' is out of the supported range (±1e15)"` |
| Unknown operation | 400 | `"Unsupported operation '<value>'. Supported: add, subtract, multiply, divide, power, sqrt, percentage"` |
| `b` present for `sqrt` | 400 | `"Operation 'sqrt' takes one operand; 'b' must be omitted or null"` |
| `b` missing for binary op | 400 | `"Operation '<op>' requires operand 'b'"` |

#### Examples

**7 + 5 = 12**
```json
// Request
{ "operation": "add", "a": 7, "b": 5 }

// Response 200
{ "success": true, "result": 12, "error": null }
```

**0.1 + 0.2 = 0.3 (FR-008)**
```json
// Request
{ "operation": "add", "a": 0.1, "b": 0.2 }

// Response 200
{ "success": true, "result": 0.3, "error": null }
```

**5 ÷ 0**
```json
// Request
{ "operation": "divide", "a": 5, "b": 0 }

// Response 200
{ "success": false, "result": null, "error": "Division by zero is not allowed" }
```

**√144**
```json
// Request
{ "operation": "sqrt", "a": 144, "b": null }

// Response 200
{ "success": true, "result": 12, "error": null }
```

**Unknown operation**
```json
// Request
{ "operation": "modulo", "a": 10, "b": 3 }

// Response 400
{ "success": false, "result": null, "error": "Unsupported operation 'modulo'. Supported: add, subtract, multiply, divide, power, sqrt, percentage" }
```

---

### GET /health

Health check for Docker and load-balancer probes.

#### Response

```
HTTP 200
Content-Type: application/json
```

```json
{ "status": "ok" }
```

No error cases — if the server is running, this returns 200.

---

## CORS Policy

| Environment | Allowed origins |
|------------|----------------|
| Development | `http://localhost:5173` (Vite dev server) |
| Docker Compose | same-origin (nginx proxy; no CORS header needed) |

Allowed methods: `POST`, `GET`, `OPTIONS`
Allowed headers: `Content-Type`

---

## Content-Type Requirements

- Request `Content-Type` must be `application/json`.
- Requests missing or with wrong `Content-Type` return HTTP 400 with structured error.
- Response `Content-Type` is always `application/json; charset=utf-8` for JSON endpoints.
