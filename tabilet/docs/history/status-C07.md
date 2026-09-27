# Retired milestone C07 — Effect Classification And Pending Steps

Milestone: `C07`
Outcome: `completed`
Retired: `2026-09-27`
Source status: `tabilet/memory-bank/status-C07.md`
Source specification: `tabilet/memory-bank/milestone.md#c07--effect-classification-and-pending-steps`
Evidence: `9d092664a6062563e0414527f997a2475aeab003`
Worktree: `clean`
Review: `passed`
Review iterations: `1`
Verification: `go generate ./schemas`; focused schema/parity/effect/pending and conversion tests; `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`; post-review immutable-document and embedded-source checks.
Consolidated into: [product contract](../../memory-bank/product.md#current-contract-surface), [README](../../../README.md), [agent instructions](../../../AGENTS.md), [UWS 1.12 specification](../../../versions/1.12.0.md), and [release changelog](../../../versions/CHANGELOG.md).

## Final milestone specification

````markdown
### C07 — Effect Classification And Pending Steps

Add optional operation `effect: read|write|unknown`, with omission interpreted
as `unknown`. Meaning determines classification; HTTP methods do not. Then
add a mutually exclusive `pending` step contract declaring purpose, inputs,
outputs, and effect without an operation. That declaration is the canonical
shape for those four fields: OpenUdon M87's step contracts reuse it for Kinet's
step-based authoring, so agree it with OpenUdon M87.1 before freezing and keep
a single definition. Pending input/output declarations
remain distinct from executable values and output expressions. Structurally
and semantically valid pending-only drafts require no dummy operation, but
executable validation rejects pending steps anywhere, including unselected
branches. Public execution entry points reject them before runtime invocation.

Develop the two additions in order against draft artifacts outside the
published-document set, then publish their combined UWS 1.12.0 contract once.
Keep earlier contracts immutable, gate the new fields to 1.12, and preserve
exact-version admission until publication. Synchronize schema, spec, Go model,
validation, conversion, embedded archive, digests, examples, release surfaces,
and evidenced current facts in their owning implementation rows. Status and
acceptance: [C07](status-C07.md).
````

## Final status document

````markdown
# Status C07 — Effect Classification And Pending Steps

**State:** Complete. C07.1, C07.2, and C07.3 are complete; the whole-milestone review passed in iteration 1.

**Specification:** [C07](milestone.md#c07--effect-classification-and-pending-steps).

**Authority and provenance:** User approved C07 -> M06 on 2026-09-26.
[S1 provenance](milestone.md#s1-provenance-and-downstream-boundaries) identifies
the clean UWS baseline and Kinet's uncommitted source document and digest.

**Prerequisite:** None. **Downstream:** M06; effect/pending adoption in the
independently owned OpenUdon, Udon, APItools, and Kinet ledgers. OpenUdon
M87.1 reuses C07.2's declaration shape for its step contracts; the two agree
that shape before either freezes.

## Tasks

Execute in this order with one active row owner. Each row includes its tests
and any evidenced current-truth corrections. Rows C07.1 and C07.2 use draft
artifacts outside `versions/`; public exact-version admission continues
rejecting unpublished UWS 1.12 until C07.3. Test candidate contracts through
isolated draft-schema and validator coverage without bypassing public admission.

| Item | State | Notes |
|---|---|---|
| C07.1 Operation effect | `[+]` | Added typed `Operation.Effect` with `read`, `write`, and `unknown`; omission remains unknown and the core does not infer from HTTP methods. Semantic validation gates the field to UWS 1.12.0 and rejects other values. Added isolated `testdata/candidate/1.12.0.json`; candidate schema accepts the three values and omission, while published 1.11 rejects the field and remains byte-identical. Kept Go known-field/schema parity version-aware and tested candidate parity. JSON/YAML/HCL round trips pass. Verification: `go test ./...`, `go vet ./...`, `git diff --check`. No published contract or current release fact changed in this draft row. |
| C07.2 Pending steps | `[+]` | Added public `PendingStep` on `Step.pending` with required purpose, one recursive `ParamSchema` each for inputs/outputs, and required effect. Both schema roots must be `type: object`; nested properties/items, required names, and `x-*` extensions are preserved. This matches OpenUdon M87.1's draft schema and fixture, rechecked against its active ledger on 2026-09-27. Pending declarations cannot combine with executable values, references, execution controls, output expressions, or structural children; ordinary `dependsOn` references retain normal integrity checks. UWS 1.12 pending-only documents accept an empty `operations` array without dummy operations. The candidate schema and internal candidate semantic check accept a valid draft while public admission still rejects unpublished 1.12 and published 1.11 rejects the new field. Executable validation and direct document, workflow, step, orchestrator, and trigger entry points reject pending work before runtime hooks, including nested, unselected, and other workflow branches. JSON/YAML/HCL round trips pass; HCL restores the required empty operations array when its repeated operation blocks are absent. Verification: focused pending/schema/conversion tests; `go test ./...`; `go vet ./...`; `go test ./uws1 -run TestSchemaConformance`; `go test ./schemas -run TestPublishedVersionDocumentsAreImmutable`; `go test ./convert -run 'RoundTrip|RoundTrips'`; `git diff --check`. Published artifacts and the pinned 1.11 corpus remain unchanged. |
| C07.3 Publish UWS 1.12 | `[+]` | Published combined `versions/1.12.0.json` and `.md` with exact-version admission, effect/pending schema-model parity, semantic and executable validation, recursive field-set conversion, regenerated embedded archive, protected digests, changelog, pending-only example, and current release surfaces. Empty `operations` is accepted only with a valid pending step under UWS 1.12+; executable validation rejects pending steps before runtime calls, including unselected branches. Effect remains descriptive and defaults to unknown; no HTTP-method inference or execution authorization was added. Preserved every earlier published digest and the pinned 1.11 corpus. Updated README, AGENTS, MkDocs, feature guides, product facts, and knowledge history. Verification passed: `go generate ./schemas`; focused schema/parity/effect/pending tests; focused conversion, immutable-document/archive, and exact-schema validation tests; `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`. |

## Acceptance And Verification

- Effect and pending semantics agree across schema, Go validation, normative
  prose, serialization, and executable rejection.
- A pending-only workflow draft passes structural and semantic validation and
  fails executable validation. A spy runtime observes zero calls on rejection,
  including when the pending step is nested or not selected by control flow.
- New fields are admitted only for 1.12 and later supported contracts. Earlier
  documents retain their behavior; no existing published digest is replaced
  to introduce these semantics.
- Focused tests cover candidate schemas before publication, then public
  exact-version admission, schema conformance, conversion, and immutable bytes:

  ```bash
  go test ./uws1 -run TestSchemaConformance
  go test ./convert -run TestRoundtrip
  go test ./schemas -run TestPublishedVersionDocumentsAreImmutable
  ```

- Full acceptance commands after the combined release:

  ```bash
  go test ./...
  go test -race ./...
  go vet ./...
  mkdocs build --strict
  git diff --check
  ```

- Record actual commands/results in the completed rows. Reconcile M06 against
  the accepted core contract and provide downstream contract/version handoffs.
  OpenUdon approval integration remains downstream; UWS executable rejection
  does not by itself prove that downstream approval checks have been adopted.

## Whole-Milestone Review Gate

**Review state:** Passed.
**Review iterations started:** 1 of at most 10.
**Review result:** No P1/P2-or-higher findings. One P3 documentation inconsistency
in the new 1.12 specification still called Browser 1.9 the current opt-in
profile despite identifying Browser 1.10 as current elsewhere. Corrected the
new 1.12 wording to identify Browser 1.10 and clarify that the page/frame
topology began in Browser 1.9; older published profile documents remain
unchanged. Updated only the 1.12 Markdown digest. The immutability and embedded
source checks and `git diff --check` passed after the correction. The separate
Browser 1.9 profile-documentation candidate remains deferred.

The whole-milestone review gate, full verification, current-truth consolidation,
and downstream handoff reconciliation are complete. OpenUdon M87.1's shared
field contract is compatible with UWS 1.12; M87.2–M87.5 may proceed without
waiting for UWS publication, while M87.6 should pin a published UWS 1.12
revision for its wrapper-mapping fixture. The full downstream handoff and
read-only sibling observations are recorded in active
[`status-M06.md`](status-M06.md). C07 is retired in
[`tabilet/docs/history/status-C07.md`](../docs/history/status-C07.md) with
implementation evidence at `9d092664a6062563e0414527f997a2475aeab003`.
````
