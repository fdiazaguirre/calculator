# Feature Specification: Full-Stack Calculator

**Feature Branch**: `001-fullstack-calculator`
**Created**: 2026-07-14
**Status**: Draft
**Input**: User description: "from @GOAL.md create a spec using the @.specify/memory/constitution.md"

## Complexity Routing *(constitution I)*

| Dimension | Score | Rationale |
|-----------|-------|-----------|
| Scope (A) | 2 | Two new surfaces (user interface + calculation service), multiple operations, but a well-understood domain |
| Risk (B) | 1 | No financial, auth, or personal data; worst failure is a wrong or missing calculation result |
| Reversibility (C) | 0 | Greenfield code; everything can be rewritten with no migration cost |

**Total: 3 → Standard tier → Spec Builder → TDD gate.** Discovery Gate applies before implementation starts; BDD and Debate gates are not required.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Perform Basic Arithmetic (Priority: P1)

A user opens the calculator, enters two numbers, chooses one of the four basic operations (addition, subtraction, multiplication, division), and sees the correct result displayed.

**Why this priority**: This is the core value of the product. Without it nothing else matters; with only this, the product is already a usable calculator (viable MVP).

**Independent Test**: Can be fully tested by entering two numbers, selecting each of the four operations, and comparing the displayed result against the mathematically correct value.

**Acceptance Scenarios**:

1. **Given** the calculator is loaded, **When** the user computes 7 + 5, **Then** the result 12 is displayed.
2. **Given** the calculator is loaded, **When** the user computes 10 − 4, **Then** the result 6 is displayed.
3. **Given** the calculator is loaded, **When** the user computes 6 × 7, **Then** the result 42 is displayed.
4. **Given** the calculator is loaded, **When** the user computes 9 ÷ 4, **Then** the result 2.25 is displayed.
5. **Given** the calculator is loaded, **When** the user computes 0.1 + 0.2, **Then** the displayed result is 0.3 (no floating-point noise such as 0.30000000000000004).

---

### User Story 2 - Receive Clear Feedback on Invalid Input (Priority: P2)

A user who enters invalid input (empty field, non-numeric text, division by zero) receives a specific, human-readable error message and the application remains usable.

**Why this priority**: Error handling is what separates a demo from a usable product. It protects the P1 flow from dead ends but is meaningless without P1 existing first.

**Independent Test**: Can be fully tested by submitting each invalid input class and verifying the exact error message, that no result is shown, and that a subsequent valid calculation still succeeds.

**Acceptance Scenarios**:

1. **Given** the calculator is loaded, **When** the user computes 5 ÷ 0, **Then** an error message "Division by zero is not allowed" is displayed and no result value is shown.
2. **Given** the calculator is loaded, **When** the user submits a non-numeric operand (e.g., "abc"), **Then** an error message identifying the invalid operand is displayed before any calculation is attempted.
3. **Given** an error message is currently displayed, **When** the user performs a valid calculation (e.g., 2 + 2), **Then** the error is cleared and the result 4 is displayed.
4. **Given** the calculation service is unreachable, **When** the user submits any calculation, **Then** the user sees a message indicating the service is unavailable, not a frozen or blank screen.

---

### User Story 3 - Perform Advanced Operations (Priority: P3)

A user can compute exponentiation, square root, and percentage in the same interface used for basic operations.

**Why this priority**: Listed as optional in the goal; adds convenience value on top of a complete P1+P2 product without changing its architecture.

**Independent Test**: Can be fully tested by computing each advanced operation with known inputs and comparing against mathematically correct values, including the negative-square-root rejection case.

**Acceptance Scenarios**:

1. **Given** the calculator is loaded, **When** the user computes 2 ^ 10, **Then** the result 1024 is displayed.
2. **Given** the calculator is loaded, **When** the user computes √144, **Then** the result 12 is displayed.
3. **Given** the calculator is loaded, **When** the user computes √(−4), **Then** an error message "Square root of a negative number is not allowed" is displayed.
4. **Given** the calculator is loaded, **When** the user computes 15% of 200, **Then** the result 30 is displayed.

---

### User Story 4 - Use the Calculator on a Phone (Priority: P3)

A user on a mobile phone (viewport width down to 320px) can read all controls, tap every button, and complete any calculation without horizontal scrolling.

**Why this priority**: Basic mobile support is required by the goal, but desktop is the primary evaluation surface; this story adapts existing behaviour rather than adding new behaviour.

**Independent Test**: Can be fully tested by rendering the interface at 320px, 375px, and 768px widths and completing a full calculation at each width with no horizontal overflow.

**Acceptance Scenarios**:

1. **Given** a 320px-wide viewport, **When** the calculator renders, **Then** all controls are visible with no horizontal scrolling.
2. **Given** a 320px-wide viewport, **When** the user completes the calculation 8 × 8, **Then** the result 64 is displayed and every tap target was reachable.

---

### Edge Cases

- Division by zero (covered by US2, decision table row D2).
- Square root of a negative number (US3, row D6).
- Non-numeric, empty, or missing operands — rejected before calculation with a field-specific message.
- Results exceeding the supported magnitude (see FR-009) — rejected with "Result out of range" rather than displaying infinity or garbage.
- Floating-point representation noise — results are presented per FR-008 so 0.1 + 0.2 displays as 0.3.
- 0 ^ 0 — returns 1 (documented convention, see Assumptions).
- 0 ^ negative exponent (e.g., 0 ^ −2) — rejected as division by zero.
- Calculation service unavailable — user-facing unavailability message (US2 scenario 4).
- Rapid repeated submissions — each submission produces exactly one result or one error; no interleaved or stale results are displayed.

### Operation Decision Table *(exhaustive — constitution I)*

| Rule | Operation | Condition | Outcome | Example (happy) | Example (negative) |
|------|-----------|-----------|---------|-----------------|--------------------|
| D1 | Add / Subtract / Multiply | Both operands valid numbers, result in range | Correct result | 7 + 5 = 12 | 1e300 × 1e300 → "Result out of range" |
| D2 | Divide | Divisor ≠ 0 | Correct quotient | 9 ÷ 4 = 2.25 | 5 ÷ 0 → "Division by zero is not allowed" |
| D3 | Any | Any operand non-numeric or missing | Validation error naming the operand | — | "abc" + 2 → "Operand 'a' must be a number" |
| D4 | Exponent | Base ≠ 0, or exponent ≥ 0 | Correct power (0^0 = 1) | 2 ^ 10 = 1024 | 0 ^ −2 → "Division by zero is not allowed" |
| D5 | Exponent | Result not representable as a real number (e.g., negative base with fractional exponent) | Validation error | (−8) ^ 3 = −512 | (−8) ^ 0.5 → "Result is not a real number" |
| D6 | Square root | Operand ≥ 0 | Correct root | √144 = 12 | √(−4) → "Square root of a negative number is not allowed" |
| D7 | Percentage | Both operands valid numbers | a% of b = a × b ÷ 100 | 15% of 200 = 30 | "x"% of 200 → validation error |
| D8 | Any | Result magnitude exceeds supported range (FR-009) | "Result out of range" error | 1e15 + 1 = 1000000000000001 | 1e308 × 10 → "Result out of range" |
| D9 | Any | Unknown operation identifier requested | Validation error listing supported operations | — | operation "modulo" → "Unsupported operation" |

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST compute addition, subtraction, multiplication, and division of two numeric operands and return the mathematically correct result.
- **FR-002**: System MUST compute exponentiation (base, exponent), square root (single operand), and percentage (a% of b = a × b ÷ 100).
- **FR-003**: System MUST reject division by zero (including 0 raised to a negative exponent) with the message "Division by zero is not allowed" and MUST NOT return a numeric result for it.
- **FR-004**: System MUST reject square root of a negative operand with the message "Square root of a negative number is not allowed".
- **FR-005**: System MUST validate every operand before calculation: present, numeric, and within the supported input range (FR-009). Validation failures MUST identify which operand is invalid and why.
- **FR-006**: System MUST validate the requested operation identifier and reject unknown operations with a message listing the supported operations.
- **FR-007**: The calculation service MUST return every response — success or error — as structured data containing: a success indicator, the result value (null on error), and an error message (null on success).
- **FR-008**: Displayed results MUST be exact for exact decimal inputs up to 12 significant digits (e.g., 0.1 + 0.2 displays as 0.3) and MUST NOT display binary floating-point noise.
- **FR-009**: System MUST accept operands with absolute value up to 1e15 (and zero) and MUST reject any operand or result outside ±1e15 with "Result out of range" / "Operand out of range" respectively.
- **FR-010**: The user interface MUST display exactly one outcome per submission: either one result or one error message, never both, and never a stale outcome from a previous submission.
- **FR-011**: The user interface MUST remain fully operable after any error (validation, calculation, or service-unavailable) with no reload required.
- **FR-012**: When the calculation service is unreachable or responds abnormally, the user interface MUST display a service-unavailable message within 5 seconds of submission.
- **FR-013**: The user interface MUST render all controls without horizontal scrolling at viewport widths from 320px to 1920px.
- **FR-014**: All error messages shown to the user MUST state what was wrong and which input caused it; raw internal error text (stack traces, technical codes) MUST never be shown.

### Key Entities

- **Calculation Request**: the operation identifier plus one or two numeric operands (square root takes one; all other operations take two).
- **Calculation Result**: the outcome envelope — success indicator, numeric result (null on error), error message (null on success).
- **Operation**: one of the seven supported operations: add, subtract, multiply, divide, power, square root, percentage. This set is closed for this feature (D9).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A first-time user can complete a basic calculation (two operands, one operation, result shown) in under 30 seconds without instructions.
- **SC-002**: 100% of the decision-table rows D1–D9 pass automated verification with at least one happy-path and one negative example each.
- **SC-003**: Every invalid input class (empty, non-numeric, division by zero, negative square root, out of range, unknown operation) produces its specific error message — zero occurrences of a generic "something went wrong" for these classes.
- **SC-004**: A user sees the calculation outcome (result or error) within 2 seconds of submission under normal conditions, and a service-unavailable message within 5 seconds when the service is down.
- **SC-005**: The full calculation flow is completable at 320px, 375px, 768px, and 1440px viewport widths with zero horizontal overflow.
- **SC-006**: Automated test coverage of calculation and validation logic is at least 80%, with evidence (coverage report) attached to task completion per constitution III/VI.
- **SC-007**: A new developer can set up and run the full application locally following only the written setup documentation, in under 15 minutes.

## Assumptions

- Advanced operations (exponentiation, square root, percentage) are in scope as P3 — the goal lists them as optional and they are included because they reuse the P1 architecture unchanged.
- Percentage semantics: `a% of b = a × b ÷ 100` (binary operation). A unary "x% = x/100" mode is out of scope.
- `0 ^ 0 = 1`, following the dominant programming-language convention; `0 ^ negative` is division by zero.
- Supported numeric domain is real decimal numbers within ±1e15 (FR-009); complex numbers, arbitrary precision, and scientific-notation entry are out of scope.
- No calculation history, persistence, user accounts, or authentication — single anonymous user per session, nothing stored between sessions.
- Expression parsing (e.g., "2 + 3 × 4" as one string) is out of scope; each calculation is one operation with explicit operands.
- Rate limiting and abuse protection are out of scope for this feature: the service is intended for local/demo deployment, not public exposure (documented exception to the global security checklist).
- Containerised deployment (Dockerfile) is a deliverable of the project but not part of this feature's user-facing behaviour; it is handled at the plan/tasks level.
- Constitution architecture constraint applies at plan time: calculation logic lives in the backend service, the frontend only presents (Architecture Constraints, constitution). This spec stays technology-agnostic; stack choices belong to the plan.
