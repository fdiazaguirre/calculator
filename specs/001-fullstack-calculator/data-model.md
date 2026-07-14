# Data Model: Full-Stack Calculator

**Branch**: `001-fullstack-calculator` | **Date**: 2026-07-14

No persistent storage. All entities are in-memory request/response structures.

---

## Entities

### CalculationRequest

Sent by the frontend; received and validated by the backend.

| Field | Type | Required | Constraints | Notes |
|-------|------|----------|-------------|-------|
| `operation` | string (enum) | yes | one of: `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`, `percentage` | Case-sensitive. Unknown value → D9 error. |
| `a` | float64 | yes | ∈ [−1e15, 1e15] | For `sqrt`, the radicand. For all binary ops, the left operand. |
| `b` | float64 \| null | conditional | ∈ [−1e15, 1e15] when present | Required for all ops except `sqrt`. Must be absent (null) for `sqrt`. |

**Validation rules (applied before computation):**
1. `operation` must be one of the seven known values.
2. `a` must be a finite number within [−1e15, 1e15].
3. For binary operations: `b` must be present and finite within [−1e15, 1e15].
4. For `sqrt`: `b` must be absent or null.
5. All validation failures return structured errors per FR-005, FR-006, FR-014.

---

### CalculationResponse

Returned by the backend for every request (success and error alike, per FR-007).

| Field | Type | On success | On error |
|-------|------|-----------|---------|
| `success` | bool | `true` | `false` |
| `result` | float64 \| null | computed value, rounded to 12 sig figs (FR-008) | `null` |
| `error` | string \| null | `null` | human-readable message (FR-014) |

**Invariants:**
- `success = true` ⟺ `result ≠ null` AND `error = null`
- `success = false` ⟺ `result = null` AND `error ≠ null`
- Never both present, never both null.

---

### Operation (closed enum)

| Identifier | Arity | Description | Special cases |
|------------|-------|-------------|---------------|
| `add` | binary | a + b | — |
| `subtract` | binary | a − b | — |
| `multiply` | binary | a × b | — |
| `divide` | binary | a ÷ b | b = 0 → D2 error |
| `power` | binary | a ^ b | 0 ^ negative → D4 error; negative base + fractional exponent → D5 error |
| `sqrt` | unary | √a | a < 0 → D6 error |
| `percentage` | binary | a × b ÷ 100 | — |

---

## State Transitions

```
[User enters operands & operation]
         │
         ▼
[Frontend validates locally: non-empty, numeric]
         │ invalid → show field error (FR-005, FR-014), no request sent
         │ valid ↓
[POST /api/v1/calculate]
         │
         ▼
[Backend validates operands (range, type, arity)]
         │ invalid → HTTP 400, success=false, error=<message>
         │ valid ↓
[Compute result]
         │ calc error (div/0, sqrt(-x), etc.) → HTTP 200, success=false, error=<message>
         │ result out of range → HTTP 200, success=false, error="Result out of range"
         │ success ↓
[Round to 12 sig figs]
         │
         ▼
[HTTP 200, success=true, result=<number>]
         │
         ▼
[Frontend displays result, clears previous error]
```

**Network failure path**: Frontend `fetch` times out or throws → display "Service unavailable" within 5 s (FR-012). No backend state changed.
