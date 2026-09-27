# Retired milestone M07 — Mock runtime enclosing-step visibility

Milestone: `M07`
Outcome: `completed`
Retired: `2026-09-27`
Source status: `tabilet/memory-bank/status-M07.md`
Source specification: `tabilet/memory-bank/milestone.md#m07--mock-runtime-enclosing-step-visibility`
Evidence: `7f843af78e508fee140b3b43f28a6b73b37a66a8`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `1`
Verification: focused mock runtime tests; `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check` (all passed on 2026-09-27)
Consolidated into: [architecture](../../memory-bank/architecture.md) for mock `$steps` scope behavior; no product, tech-stack, or reusable-lesson change was needed. No evolution bump.

## Final milestone specification

````markdown
### M07 — Mock runtime enclosing-step visibility

**Goal.** Correct the confirmed P2 R5 finding from the 2026-09-27 UWS/APItools/OpenUdon stage 1 review. This is M06 remediation; its retired record remains unchanged.

**Acceptance.** A mock or hybrid runtime leaf inside `forEach` can resolve a completed visible outer step output, while same-iteration records remain preferred and other iterations, workflow invocations, and parallel branches stay isolated. Missing or ambiguous records continue to fail explicitly. Use the existing UWS 1.12 wire contract and public Go API; no new UWS version, schema field, or real transport is needed.

**Owner and checks.** [M07.1](status-M07.md) owns implementation, orchestrator-backed regressions and any current-truth correction. M06 acceptance is the prerequisite. Run focused tests, `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, `git diff --check`, then the persisted whole-milestone review gate. OpenUdon simulation and Kinet's later W04 are downstream consumers; M07 does not gate Kinet W03.
````

## Final status document

````markdown
# Status M07 — Mock runtime enclosing-step visibility

**State:** Complete. Whole-milestone review passed in iteration 1.

**Provenance:** 2026-09-27 UWS/APItools/OpenUdon stage 1 review R5; source P2, local P2, confirmed against clean UWS `7f843af78e508fee140b3b43f28a6b73b37a66a8`. Review baseline is the same commit; no uncommitted repository changes were part of the evidence. `mockruntime/expressions.go` filters an outer `step:config` record when a `forEach` item resolves `$steps.config.outputs.token`. Historical owner M06 is retired.

| Item | State | Notes |
| --- | --- | --- |
| M07.1 Enclosing-step lookup | `[+]` | A forEach leaf resolves completed outer outputs through the public mock runtime; nearest visible iteration wins, while sibling iterations, nested invocations, missing records, and duplicate nearest records fail or remain isolated. Orchestrator-backed execution plus direct scope regressions passed. `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed. Architecture current truth updated. |

## Acceptance and review

M07.1 and the verification in `milestone.md` must pass before a whole-milestone review. Persist the review iteration here before each pass (maximum ten).

**Review iteration 1:** Started after M07.1 and all specified checks passed. Reviewed the full milestone diff from `7f843af78e508fee140b3b43f28a6b73b37a66a8`, including iteration ancestry, invocation isolation, ambiguity, executable behavior, current truth, and regression coverage. No P1/P2 finding remains. Current truth is consolidated in `architecture.md`; there is no new wire or API contract, so no evolution bump or new lesson is needed. OpenUdon simulation and Kinet W04 remain their own downstream consumers.
````
