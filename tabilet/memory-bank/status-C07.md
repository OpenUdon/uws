# Status C07 — Effect Classification And Pending Steps

**State:** Planned; no implementation or acceptance verification performed.

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
| C07.1 Operation effect | `[ ]` | Add optional `effect: read\|write\|unknown`; omission means unknown. Cover source-bound and extension-owned operations, invalid values, version gates, tags/known fields, and JSON/YAML/HCL preservation. Specify classification by meaning, with no HTTP-method inference. Unknown is treated like write for hybrid routing; effect metadata alone grants no execution permission. |
| C07.2 Pending steps | `[ ]` | Add the mutually exclusive `pending` contract with purpose, input/output declarations, and effect, without an operation. Keep declarations distinct from executable inputs and output expressions. Cover valid pending-only drafts without dummy operations, malformed/mixed forms, and normal reference integrity. Schema and semantic validation accept valid drafts; executable validation rejects pending nodes throughout nested steps/cases/defaults and unselected branches. Public execution entry points must reject before invoking runtime hooks. Cover conversions and rejection on earlier versions. The purpose/input/output/effect declaration is the canonical shape reused by OpenUdon M87 step contracts; agree it with OpenUdon M87.1 before freezing, and treat later shape changes as a coordinated contract change. Aligned 2026-09-27 at the owner's request (Kinet `docs/kinet-order.md` X2). |
| C07.3 Publish UWS 1.12 | `[ ]` | Publish the combined `versions/1.12.0.json` and `.md` once with exact-version admission, schema/model parity, archive regeneration, protected digests, changelog, examples, and release surfaces. Preserve every earlier published artifact and the pinned 1.11 corpus. Update README, AGENTS, documentation/navigation, and current memory-bank facts as supported by the release; follow the existing knowledge-preservation policy for replaced facts. Verify both new features together, earlier-version compatibility, and downstream handoff. |

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

Not started; zero iterations consumed. After all rows and acceptance commands
pass, persist iteration 1 before reviewing the complete milestone. Follow the
[review gate](milestone.md#milestone-review-gate): fix every P1/P2 or higher
finding, rerun affected verification, and review the whole milestone again.
The limit is 10 persisted iterations; resume interrupted passes at their saved
number. No terminal row alone establishes milestone acceptance. Retire only
after the gate, verification, current-truth consolidation, and downstream
reconciliation pass under the existing retirement policy.
