# C09 — Binding contracts

**Stage:** Kinet STG-11, Phase A. **Owner:** UWS.
**State:** Approved planning on 2026-10-06; 4 pending rows, no implementation or acceptance.
**Source baseline:** `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` (clean at planning).
**Coordinator:** [Stage 11 contract](../../../kinet/docs/stage11.md); the package-local milestone/status owns acceptance.

## Dependencies and handoff

[UWS:C08](status-C08.md).
The serial predecessor is a scheduling gate; direct contract and regression impacts are also listed. Every prerequisite must pass its whole review, and required publication must be independently verified before adoption. Record exact accepted/published sources and fixture/build hashes; no Stage 11 acceptance or future pin is claimed yet.

**Downstream:** [APItools:M82](../../../apitools/tabilet/memory-bank/status-M82.md), [Udon:M47](../../../udon/tabilet/memory-bank/status-M47.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [OpenUdon:P09](../../../openudon/tabilet/memory-bank/status-P09.md). Reconcile every affected consumer against the accepted prerequisite revision before advancing.

## Tasks

| Item | State | Notes |
|---|---|---|
| C09.1 — Define operation shapes and resolvers | `[ ]` | Define source-neutral OperationShape, Resolver and versioned ShapeTable contracts with source identity, native selector, protocol, supported schemas and security alternatives. Reserve extensible browser and runtime-function kinds without claiming their implementations. |
| C09.2 — Validate bindings and security declarations | `[ ]` | Check operation resolution, required inputs, literal/schema compatibility, known expression types, response-field references and symbolic security requirements. Missing or unsupported evidence remains indeterminate, never silently compatible. |
| C09.3 — Analyze flow deterministically | `[ ]` | Add stable, value-free diagnostics for structural reachability, unused outputs, effect/pending ordering and unbounded loops. Shared analysis is observation; authorization and enforcement remain consumer policy. |
| C09.4 — Qualify and publish the shape contract | `[ ]` | Compare semantic findings against accepted OpenUdon fixtures, add tampered/ambiguous/incomplete conformance cases, verify no APItools dependency, and publish the exact accepted contract only with named authority. |

## Acceptance and verification

Portable binding and flow diagnostics have stable codes and deterministic fixtures. Shape production stays outside UWS; no diagnostic grants execution or silently tightens legacy validation.

go test ./...; go test -race ./...; go vet ./...; schema/conformance and published-version immutability checks; mkdocs build --strict; git diff --check. Run the separate codec module checks once it exists.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.

## Execution policy

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## Reconciled consumer contract — 2026-10-06

Udon:M47 consumes this contract for the runtime-function catalog. OpenUdon:P09 owns source/shape-to-authority verification; M98 no longer exposes synthesis-coupled construction.

## Persisted review

- Review iteration: **0/10**; not started.
- Closing-review findings: none; the whole-milestone review has not started. Approved intake requirements above remain pending.
- Accepted revision: not available.
- Published revision / artifact evidence: not available.
- Verification: pending implementation; no test result is claimed by this planning record.

After all tasks finish, perform the whole-milestone review with persisted iteration/finding state and fix every P1/P2 before acceptance. Resume an interrupted pass at the same counter. Consolidate current facts, reconcile downstream work and retire under this package’s normal procedure.
