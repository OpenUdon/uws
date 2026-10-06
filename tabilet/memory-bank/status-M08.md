# M08 — Verified HCL presentation

**Stage:** Kinet STG-11, Phase A. **Owner:** UWS.
**State:** Approved planning on 2026-10-06; 4 pending rows, no implementation or acceptance.
**Source baseline:** `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` (clean at planning).
**Coordinator:** [Stage 11 contract](../../../kinet/docs/stage11.md); the package-local milestone/status owns acceptance.

## Dependencies and handoff

[APItools:M82](../../../apitools/tabilet/memory-bank/status-M82.md).
The serial predecessor is a scheduling gate; direct contract and regression impacts are also listed. Every prerequisite must pass its whole review, and required publication must be independently verified before adoption. Record exact accepted/published sources and fixture/build hashes; no Stage 11 acceptance or future pin is claimed yet.

**Downstream:** [Udon:M47](../../../udon/tabilet/memory-bank/status-M47.md), [Udon:M48](../../../udon/tabilet/memory-bank/status-M48.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [Kinet:W17](../../../kinet/tabilet/memory-bank/status-W17.md). Reconcile every affected consumer against the accepted prerequisite revision before advancing.

## Tasks

| Item | State | Notes |
|---|---|---|
| M08.1 — Specify the additive 1.13 transition | `[ ]` | Publish new versioned documents for the HCL presentation/deprecation contract without editing older published specifications. Keep legacy core APIs and input support during Stage 11; do not schedule removal here. |
| M08.2 — Add render verify and import APIs | `[ ]` | Add the separate uws/hcl module with deterministic Render, lossless Verify and a deprecated public Import. Bind provenance to exact source bytes and codec revision; preserve existing typed key/extension mappings. |
| M08.3 — Prove lossless presentation | `[ ]` | Test large integers, precise decimals, exponent notation, required lexeme preservation, strings, keys, extensions and malformed HCL. Never use float64 or JCS equality as the losslessness oracle; refuse a misleading view. |
| M08.4 — Qualify and publish both modules | `[ ]` | Verify old APIs and immutable version hashes, schema/code/docs/archive parity, nested-module builds and codec round trips. Publish exact root/codec sources with named authority and record consumer handoffs. |

## Acceptance and verification

A deterministic HCL view can be proved against the approved document. Existing HCL readers remain available; removal and core dependency extraction are deferred, not falsely claimed complete.

go test ./...; go test -race ./...; go vet ./...; schema/conformance and published-version immutability checks; mkdocs build --strict; git diff --check. Run the separate codec module checks once it exists.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.

## Execution policy

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## Reconciled consumer contract — 2026-10-06

Kinet:W17 consumes verified presentation via M46; legacy packaged HCL remains a distinct artifact.

## Persisted review

- Review iteration: **0/10**; not started.
- Closing-review findings: none; the whole-milestone review has not started. Approved intake requirements above remain pending.
- Accepted revision: not available.
- Published revision / artifact evidence: not available.
- Verification: pending implementation; no test result is claimed by this planning record.

After all tasks finish, perform the whole-milestone review with persisted iteration/finding state and fix every P1/P2 before acceptance. Resume an interrupted pass at the same counter. Consolidate current facts, reconcile downstream work and retire under this package’s normal procedure.
