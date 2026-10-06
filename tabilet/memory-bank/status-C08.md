# C08 — Portable expressions

**Stage:** Kinet STG-11, Phase A. **Owner:** UWS.
**State:** Confirmed Stage 11 goal execution; C08.1 complete, C08.2–C08.4 pending; review 0/10 not started.
**Source baseline:** `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` (clean at planning).
**Coordinator:** [Stage 11 contract](../../../kinet/docs/stage11.md); the package-local milestone/status owns acceptance.

## Dependencies and handoff

[Kinet:M45](../../../kinet/tabilet/docs/history/status-M45.md).
The serial predecessor is a scheduling gate; direct contract and regression impacts are also listed. Every prerequisite must pass its whole review, and required publication must be independently verified before adoption. Record exact accepted/published sources and fixture/build hashes; no Stage 11 acceptance or future pin is claimed yet.

**Downstream:** [UWS:C09](status-C09.md), [Udon:M48](../../../udon/tabilet/memory-bank/status-M48.md), [OpenUdon:M98](../../../openudon/tabilet/memory-bank/status-M98.md). Reconcile every affected consumer against the accepted prerequisite revision before advancing.

## Tasks

| Item | State | Notes |
|---|---|---|
| C08.1 — Parse the existing expression grammar | `[+]` | Implement the existing sources, JSON Pointer/dot traversal, comparison operators, expression contexts and document-version gates. Add no language extension and do not reinterpret function templates or browser-profile strings. Evidence: expressions/parse.go and parser vectors cover all core sources, operator longest-match/whitespace, scalar lexemes, pointer escapes and version/field/loop refusals. Parser race/vet and published-version immutability checks passed; no old schema/spec/archive, ordinary validator or mock behavior changed. docs/expression-reference.md records the public boundary. |
| C08.2 — Evaluate expressions with shared vectors | `[ ]` | Add a reference evaluator and parse/evaluation vectors for null propagation, comparisons, exact values and workflow/iteration visibility. Preserve current normative source/context distinctions. |
| C08.3 — Adopt without narrowing ordinary validation | `[ ]` | Use the reference evaluator in mockruntime with compatibility tests. Keep strict portability checking opt-in and preserve ordinary validation of implementation-specific expressions; diagnose legacy expr wrappers explicitly. Kinet:W18 explicitly enables strict portability for all newly authored packages; this does not narrow ordinary legacy validation. |
| C08.4 — Qualify and prepare the parser release | `[ ]` | Run schema/version, evaluator, mock-runtime and race checks; prepare exact accepted source and a separately pinned conformance supplement. Publish only under named authority before a consumer requires publication. |

## Acceptance and verification

The existing normative grammar has one tested reference implementation. Existing document acceptance and frozen conformance artifacts remain compatible; strict portability checks are explicit.

go test ./...; go test -race ./...; go vet ./...; schema/conformance and published-version immutability checks; mkdocs build --strict; git diff --check. Run the separate codec module checks once it exists.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.

## Execution policy

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## Reconciled consumer contract — 2026-10-06

Kinet:W18 requires strict portability for new packages; the opt-in UWS API preserves ordinary legacy validation. OpenUdon:M98 consumes affected mockruntime regression vectors.

## Accepted transition prerequisite — 2026-10-06

Kinet:M45 is accepted/retired at `76c5a7cc577cd1dc86d21e9c3a3bd372e3c807b7` after whole review 3/10. Consume the [behavior/corpus baseline](../../../kinet/docs/m45-transition-baseline.md), [bounded worker consumer requirements](../../../kinet/docs/stage11-worker-contract.md) and [qualified development baseline](../../../kinet/docs/m45-qualification.md). The final three-sample profile is docs/m45-profile-final.json (retained worker median 8.939s); matching semantics and recorded investigation thresholds govern W17/Phase B comparisons. M45 proves neither live-package completeness nor new library/worker acceptance. Kinet publication is not a prerequisite; the eight sibling publication gates remain explicitly authorized by the confirmed request.

## Persisted review

- Review iteration: **0/10**; not started.
- Closing-review findings: none; the whole-milestone review has not started. Approved intake requirements above remain pending.
- Accepted revision: not available.
- Published revision / artifact evidence: not available.
- Verification: pending implementation; no test result is claimed by this planning record.

After all tasks finish, perform the whole-milestone review with persisted iteration/finding state and fix every P1/P2 before acceptance. Resume an interrupted pass at the same counter. Consolidate current facts, reconcile downstream work and retire under this package’s normal procedure.
