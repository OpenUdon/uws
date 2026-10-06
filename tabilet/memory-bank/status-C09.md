# C09 — Binding contracts

**Stage:** Kinet STG-11, Phase A. **Owner:** UWS.
**State:** Confirmed Stage 11 execution; All four tasks complete; whole review 2/10 started, acceptance pending; review 0/10 not started.
**Source baseline:** `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` (clean at planning).
**Coordinator:** [Stage 11 contract](../../../kinet/docs/stage11.md); the package-local milestone/status owns acceptance.

## Dependencies and handoff

[UWS:C08](../docs/history/status-C08.md).
The serial predecessor is a scheduling gate; direct contract and regression impacts are also listed. Every prerequisite must pass its whole review, and required publication must be independently verified before adoption. Record exact accepted/published sources and fixture/build hashes; no Stage 11 acceptance or future pin is claimed yet.

**Downstream:** [APItools:M82](../../../apitools/tabilet/memory-bank/status-M82.md), [Udon:M47](../../../udon/tabilet/memory-bank/status-M47.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [OpenUdon:P09](../../../openudon/tabilet/memory-bank/status-P09.md). Reconcile every affected consumer against the accepted prerequisite revision before advancing.

## Tasks

| Item | State | Notes |
|---|---|---|
| C09.1 — Define operation shapes and resolvers | `[+]` | Define source-neutral OperationShape, Resolver and versioned ShapeTable contracts with source identity, native selector, protocol, supported schemas and security alternatives. Reserve extensible browser and runtime-function kinds without claiming their implementations. Evidence: binding/types.go/table.go and contract tests qualify exact identity/aliases, unknown completeness, security OR/AND alternatives, snapshot immutability, malformed/duplicate/unknown fields, forged-source mismatch and non-HTTP refusal. Race/vet pass; no source parser, provider I/O or dependency change. docs/binding-reference.md records the public contract. |
| C09.2 — Validate bindings and security declarations | `[+]` | Check operation resolution, required inputs, literal/schema compatibility, known expression types, response-field references and symbolic security requirements. Missing or unsupported evidence remains indeterminate, never silently compatible. Evidence: binding/validate.go and regressions prove required input/literal schema, known expression type containment, response references and symbolic security OR/AND checks. Unsupported refs use a refusing loader and remain indeterminate; custom resolver source/selector mismatches are rejected without value/error text exposure. Race/vet pass; no key/provider loading or ordinary validation change. |
| C09.3 — Analyze flow deterministically | `[+]` | Add stable, value-free diagnostics for structural reachability, unused outputs, effect/pending ordering and unbounded loops. Shared analysis is observation; authorization and enforcement remain consumer policy. Evidence: binding/flow.go reports sorted bounded code/path findings for possible reachability, core output references, unknown effects, sequence/branch pending-before-write order, cycles and static work bounds. Privacy/nonmutation/cancellation/merge exclusion regressions and binding races/vet pass. Conditions and opaque profiles are not executed or inferred; findings never supply permission. |
| C09.4 — Qualify and publish the shape contract | `[+]` | Compare semantic findings against accepted OpenUdon fixtures, add tampered/ambiguous/incomplete conformance cases, verify no APItools dependency, and publish the exact accepted contract only with named authority. Evidence: separately pinned OpenUdon source/step-check fixtures and neutral projection pass semantic/negative checks; full offline tests/races/vet/strict docs, old immutability/conformance and no-APItools/private-runtime dependency proofs pass. docs/c09-qualification.md records scope; qualified publication proceeds under the confirmed origin/main envelope before closing-review adoption. |

## Acceptance and verification

Portable binding and flow diagnostics have stable codes and deterministic fixtures. Shape production stays outside UWS; no diagnostic grants execution or silently tightens legacy validation.

go test ./...; go test -race ./...; go vet ./...; schema/conformance and published-version immutability checks; mkdocs build --strict; git diff --check. Run the separate codec module checks once it exists.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.

## Execution policy

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## Reconciled consumer contract — 2026-10-06

Udon:M47 consumes this contract for the runtime-function catalog. OpenUdon:P09 owns source/shape-to-authority verification; M98 no longer exposes synthesis-coupled construction.

## Accepted expression prerequisite — 2026-10-06

UWS:C08 is accepted and independently observed on origin/main at `0411eea6fc84fbd6aa97cef94f53f301260f4844`, whole review 2/10. [Qualification](../../docs/c08-qualification.md) and the [supplement manifest](../../docs/examples/expressions/v1/manifest.json) identify exact source/vector bytes. Ordinary validation and frozen published artifacts are unchanged; strict portability is opt-in, required for new Kinet:W18 packages. Legacy mock numeric/encoded-root adapters remain explicit. Values must use lossless json.Number projections before constructing snapshots because outer UseNumber does not override legacy custom model decoders. C08 supplies no source parsing, credential/provider I/O or execution authority. Retirement closure `5c0c74f48d84588e3ff4f994f0713f199cfcc67c` was independently observed on origin/main; that publication gate is satisfied. [Publication evidence](../../docs/c08-publication.md) records the accepted-source ancestry.

## Persisted review

- Review iteration: **2/10**; started on 2026-10-06 after R1-F01/F02 corrections; full rereview underway.
- Closing-review findings, iteration 1: **R1-F01 (P2)** — a body object with nested core expressions is validated as a literal string-bearing object, falsely rejecting a reviewed integer expression. **R1-F02 (P2)** — schemaCompatibility treats number-to-integer as disjoint although it has partial overlap. Regressions reproduce both supported scenarios. Add recursively projected template schemas and proven structural containment, keep unsupported constraints indeterminate, and correct numeric overlap before rereview. Both findings are fixed and their nested-object/array, type-overlap, unsupported-constraint and open-source-schema regressions pass. Full tests/races/vet passed after fixes.
- Closing-review findings, iteration 2: **R2-F01 (P2)** — template discovery handles only map[string]any/[]any, while literal validation already accepts equivalent typed JSON-compatible containers. A map[string]string containing a reviewed integer expression is falsely rejected. Normalize consumer values losslessly before template discovery and add the typed-container equivalence regression; preserve error/value redaction. R2-F01 is fixed: typed values are normalized losslessly before template discovery. The regression and full tests/races/vet pass. No other P1/P2 found on whole rereview.
- Accepted revision: not available.
- Published revision / artifact evidence: qualified source `df54c6644439114fa749ab962405b71c1bfcfff8` pushed fast-forward to authorized UWS origin/main, independently observed by git ls-remote. Accepted review/closure is required before adoption.
- Verification: full offline tests/races/vet/strict docs, immutable-artifact and no-APItools/private-runtime checks passed; pinned OpenUdon semantic fixture comparison and negative/tampered/ambiguous/incomplete tests pass.

After all tasks finish, perform the whole-milestone review with persisted iteration/finding state and fix every P1/P2 before acceptance. Resume an interrupted pass at the same counter. Consolidate current facts, reconcile downstream work and retire under this package’s normal procedure.
