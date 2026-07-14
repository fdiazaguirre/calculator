<!--
Sync Impact Report
==================
Version change: none → 1.0.0
Source documents:
  - agentic-ai-playbook.md (SDD methodology, quality gates, parallel execution, complexity gating)
  - global_rules.md (KISS, LoB, YAGNI, deep modules, information hiding, strategic programming)
Modified principles: N/A (initial version)
Added sections:
  - Core Principles (6 principles)
  - Architecture Constraints
  - Quality Gates
  - Governance
Removed sections: N/A
Templates:
  ✅ .specify/templates/plan-template.md — Constitution Check gate aligns with principles I–VI
  ✅ .specify/templates/spec-template.md — User stories + acceptance scenarios align with SDD pipeline
  ✅ .specify/templates/tasks-template.md — Phase/checkpoint structure aligns with SDD stages
Deferred TODOs: None
-->

# Calculator Constitution

## Core Principles

### I. Spec-Driven Development (NON-NEGOTIABLE)

Every non-trivial change MUST pass through the Zero-Gap Pipeline before any production code is written:

**Discovery → Spec → BDD → TDD → Integration → Verification**

Complexity is scored 0–6 across Scope (A), Risk (B), and Reversibility (C) before routing:

| Score | Tier | Pipeline |
|-------|------|----------|
| 0 | Minimal | Direct execution — no pipeline overhead |
| 1–2 | Light | TDD gate only |
| 3–4 | Standard | Spec Builder → TDD gate |
| 5–6 | Full | Spec Builder → Debate → BDD gate → TDD gate |

- All red-card unknowns MUST be resolved before implementation proceeds; no guessing.
- Vague words (`fast`, `secure`, `robust`, `handles errors`) are forbidden in specs without measurable replacements (e.g., "responds within 200ms at p99").
- Decision tables MUST be exhaustive — no "TBD" cells, no blank entries.
- Every acceptance rule MUST have ≥ 2 concrete examples (at least one happy path, at least one negative case).
- Adversarial quality checklist (10 points from the playbook) MUST pass before any spec is locked.

**Rationale:** Rework from unresolved ambiguity is the leading cost multiplier in agentic systems.
The pipeline forces full resolution upfront, eliminating downstream surprises in a multi-service,
multi-language codebase.

### II. Deep Modules & Information Hiding

Each service and module MUST expose a simple interface that hides significant internal complexity:

- Service interfaces MUST hide implementation details (algorithms, storage mechanics, retry logic) from callers.
- Cross-service business context MUST be expressed via the `TransactionContext` schema; ad-hoc business headers are forbidden.
- Modules MUST NOT push configuration or edge-case handling onto callers — provide sensible defaults internally.
- Prefer one deep module over several shallow ones that merely delegate.
- Information hiding MUST be applied at every boundary: service-to-service, package-to-package, function-to-caller.

**Rationale:** Calculator's decision accuracy is its value. Shallow modules expose internal complexity to every
caller — that complexity compounds across a distributed system and erodes maintainability.

### III. Test-First (NON-NEGOTIABLE)

All business logic MUST be developed using Red-Green-Refactor:

1. Write a failing test (RED) — confirmed failing before implementation begins.
2. Write the minimum implementation to pass it (GREEN).
3. Refactor without changing behaviour (IMPROVE).
4. Verify coverage ≥ 80% before marking a task complete.

- Tests MUST be written and confirmed failing before a single line of implementation is written.
- Integration tests MUST verify end-to-end reachability via live service endpoints, not just unit correctness.
- The Security Gate MUST include threat modelling for all auth, crypto, and agent-decision code paths.
- No task may be marked complete without test output as evidence. "I fixed it" is not evidence.

**Rationale:** Agent decision logic is high-stakes — approve/decline on financial transactions.
Silent regressions have direct business consequences. The test-first discipline makes regressions
visible immediately.


### IV. Minimal Abstractions & Strategic Simplicity

Code MUST use the fewest abstractions that solve the actual, present problem:

- Prefer `errors.New` / `fmt.Errorf` over custom exception hierarchies unless domain semantics require distinct handling.
- Composition over inheritance in all service and handler construction.
- YAGNI: no interface, pattern, or abstraction may be introduced before a second concrete use case exists.
- Locality of Behaviour (LoB): behaviour MUST be co-located with the code it affects; scattered logic is a violation.
- "Design it Twice": before committing to any module interface, sketch two radically different approaches
  and select the simpler one.
- Invest 10–20% of implementation time on root-cause fixes and design improvement (Strategic Programming);
  tactical shortcuts that defer complexity are not acceptable.

**Rationale:** Calculator is already architecturally complex (multi-service, multi-language).
Every unnecessary abstraction multiplies that complexity for every future contributor and every future agent.

### VI. Evidence-Based Verification

No task, feature, or claim of correctness is accepted without observable evidence:

- Automated test output (passing, coverage ≥ 80%) is mandatory for every completed task.
- Integration verification MUST confirm the feature is reachable via live service endpoints.
- State checkpoints MUST be saved when agent context pressure reaches 60% or before any risky multi-file operation.
- "It works" without test output is treated as unverified. Reviewers MUST request evidence before approving.

**Rationale:** In agentic systems where agents report on their own work, unverified claims propagate
silently through downstream stages. Evidence gates break that propagation.

## Architecture Constraints

These constraints govern all services in the Calculator stack and MUST NOT be violated without an ADR:
- **Language boundaries**: Go handles the logic React the UI.

## Quality Gates

All implementation MUST clear these gates at the indicated stage:

| Gate | Trigger | Blocks |
|------|---------|--------|
| Discovery Gate | Any Standard+ tier (score ≥ 3) | Implementation start |
| Debate Gate | Full tier (score 5–6) or architectural decision | Spec lock |
| BDD Gate | Full tier (score 5–6) | Coding start |
| TDD Gate | All business logic | Task completion |
| Integration Gate | All feature branches | Merge to main |


Gate violations in PRs MUST be documented as explicit exceptions with written justification.

## Governance

- This constitution supersedes all other practices, patterns, and preferences documented in this project.
- Amendments require: (1) written rationale, (2) impact analysis on existing services, (3) migration plan if behaviour changes, (4) version bump per the policy below.
- `CONSTITUTION_VERSION` follows semantic versioning:
    - MAJOR: principle removed, renamed, or incompatibly redefined.
    - MINOR: new principle or section added.
    - PATCH: clarification, wording fix, or example added.
- All PRs MUST include a Constitution Check section verifying compliance with applicable principles.
- Complexity violations (e.g., unnecessary abstraction layer, custom exception hierarchy) MUST be
  justified in the plan's Complexity Tracking table before implementation proceeds.
- Use `.specify/` artifacts (spec, plan, tasks) as the runtime development workflow for all Standard+ tier work.

**Version**: 1.0.0 | **Ratified**: 2026-05-31 | **Last Amended**: 2026-05-31