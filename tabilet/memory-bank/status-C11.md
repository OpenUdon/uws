# C11 — HCL input removal version (conditional)

**Owner-approved Stage 12 closure and Phase C deferral — 2026-10-25.** The owner approved the recommendation to defer Phase C and close Stage 12 after required work was accepted. All five untouched conditional units (Kinet:M54/M55, Udon:M54, UWS:C11 and APItools:M86), nineteen rows total, are withdrawn together as `[X]`; none is implemented or accepted by this disposition. Their original IDs, task text, contracts and review counters are retained. Deferred HCL-removal intent belongs to Kinet STG-13 provisional context only, without executable successors. Future work requires a fresh approved plan/IDs and fresh concrete authority; old triggers cannot reactivate withdrawn rows. Original dependency edges below are historical, with no active dependents left in this family. No deployment or inventory is required to close Stage 12. HCL input removal is UNDELIVERED; existing behavior stays. W8M remains the only Stage 12 website qualification target. All accepted evidence, historical UNKNOWNs, counters, denials and disabled audit persist. No execution, publication, inventory or deployment is authorized.

**Stage:** Kinet STG-12, Phase C. **Owner:** UWS.
**State:** Withdrawn from Stage 12; all original rows `[X]`. Not implemented or accepted. Review 0/10 unchanged.
**Original trigger (withdrawn):** [Kinet:M54](../../../kinet/tabilet/memory-bank/status-M54.md) records
a clean, authorized, value-free live inventory, and
[Udon:M54](../../../udon/tabilet/memory-bank/status-M54.md) is accepted and
published. If the trigger is absent, the goal skips this milestone without
completing or cancelling it.
**Source baseline:** `989e3f2c88cac5c0f5a2911dfe04c36a61e43126` (clean at planning).
**Coordinator:** [Stage 12 contract](../../../kinet/docs/stage12.md#conditional-hcl-phase).
Boundary: UWS 1.13 §12.1 (`versions/1.13.0.md`) requires a separately versioned
compatibility decision plus consumer migration evidence.

**Planning reconciliation.** Stage 12 planning review, 2026-10-09; source priority not supplied. Review baseline and current revalidation: `989e3f2c88cac5c0f5a2911dfe04c36a61e43126`; includes the uncommitted Stage 12 planning files. F05, F06 and F07 (confirmed P2 each) are owned by C11.2, C11.1 and C11.3 respectively. Evidence: `../kinet/workers/author/workflow_view.go` uses `hcl.Import` for JSON/YAML comparison; `uws1/hcl.go` retains Horizon marshaling and `../udon/pkg/uwsprofile/document_io.go` exports through it; original C11.3 required unmodified consumers before their M86/M55 migrations. Udon:M54 migrates output first; Kinet:M55 migrates views/conversion and proves final adoption. This is approved review intake, not a closing-review iteration; all rows and review counters remain pending/0.

**Parallel planning provenance.** Stage 12 parallel execution review,
2026-10-09; source priority and separate review baseline not supplied. Current
revalidation `989e3f2c88cac5c0f5a2911dfe04c36a61e43126`, including uncommitted planning files. F01 (confirmed local P2) adopts scoped workflow ownership; F02 (confirmed Lower) replaces strict ordering with readiness; F04 (confirmed local P2) requires frozen consumer checks. F05 (confirmed conditional Lower) permits independent retirement/inventory branches while preserving operational triggers and the C11 join.
Evidence: the old Kinet goal/launch rules, owning agent rules, existing pending
dependencies and sibling consumer checks; Udon `go.mod`/`pkg/execute` import no
OpenUdon code. This approved intake changes no row state or review counter.

## Status rules

See [`milestone.md`](milestone.md) for milestone definitions and the status ID
pattern that names this file.

Status markers:

| Symbol | Suggested Status | Interpretation |
|---|---|---|
| `[ ]` | Pending | Item is actionable when its dependencies pass. |
| `[+]` | Completed | Item finished or done. |
| `[~]` | In Progress | Item is the one general row currently being worked. |
| `[!]` | Blocked | A current unresolved blocker still prevents progress. |
| `[X]` | Cancelled | The approved item was withdrawn; any later replacement needs fresh planning. Cancellation does not prove delivery. |
| `[-]` | Closed Historical | A consumed failed attempt or superseded row retained for audit; it is never retried and does not block its accepted successor. |

Write markers with backticks, exactly as in the table above: `` `[ ]` ``, not
`[ ]`.

Across the active ledger, zero or one general row may be `[~]`. Resume it before
selecting another row. An operational launcher additionally requires its exact
authorized operation row to be `[~]` before invocation; the marker records the
selected operation but does not grant missing external-mutation authority.

Read this milestone's optional `AUTHORIZATION_REQUIREMENTS` in
[`milestone.md`](milestone.md#authorization-requirements); do not duplicate its
declaration here. Check protected actions against the current human-approved
request and effective grants before acting. Notes may record a grant ID,
performed-action evidence, or a missing-authority blocker; they never prove
approval. Required unperformed actions prevent closure. No new task marker is
introduced, and absent authorization metadata grants no additional authority.

With human-approved `CORRECTIVE_PLANNING: auto`, the owner may prepare and
check necessary repair plans within the original delivery outcome, existing
grants and explicit limits. Preserve this row's identity and interrupted work;
move the interrupted `[~]` row to `[!]` before activating a repair row. Record
prerequisite IDs, standing approval source and resumption condition. Keep zero
or one `[~]` row. Once every prerequisite is verified and accepted and all
blockers clear, restore this row from `[!]` to `[~]` under the same approval.
Recover this chain on resume; do not select unrelated work. Preserve partial
implementation; suspension bookkeeping may join the scoped repair checkpoint,
but unfinished implementation cannot. Pause if overlapping changes cannot be
separated safely. Explicit project restrictions still govern.
Preserve review counts and frozen history. Children report needs
to the owner; a marker never grants planning or action authority.

Each table row is a commit unit: after changing its state, verify the change,
update the memory bank/docs, and make a scoped `git commit` before starting the
next row, except the scoped repair suspension bookkeeping above. A `[-]` row's
notes must record the consumed attempt or supersession
and identify its accepted successor. If two rows are inseparable, redefine them
as one row before starting rather than closing several rows in one commit.
During a goal run, the resolved `COMMIT_POLICY` governs instead; `none` means
no commits even when a requirement or grant mentions `git.commit`.

Keep this file a current-state record, not a log. Replace a superseded checkpoint
instead of appending dated narrative, and link per-attempt detail in evidence
files. A row whose acceptance depends on an expensive check records its attempt
count here under the [verification-attempt gate](milestone.md#verification-attempt-gate).

After the milestone's review, consolidation, and downstream reconciliation
pass, follow [the retirement procedure](milestone.md#long-term-memory-and-retirement).
Preserve this complete document and the full milestone specification in its
retired record; earlier rows, notes, and original path context must survive.
Retirement follows the governing commit policy and never changes the task
markers or permits reuse of this milestone ID.

## Dispatch and lease boundaries

**Depends on.** Kinet:M54, Udon:M54, OpenUdon:A32. All required prerequisites must have
accepted closure at exact revisions; sibling producer adoption also needs
independently verified publication. A priority position never supplies authority.

**Downstream impacts.** APItools:M86?, Kinet:M55?.

**Write set.** The owning `uws/` package's implementation, tests, ordinary
documentation, manifests and qualification outputs only as required by this
milestone's existing scope, plus `tabilet/memory-bank/status-C11.md` in its
assigned worktree. Excludes `AGENTS.md`, `tabilet/GOAL.md`, shared memory-bank
files, other statuses, evolution, stages, history/knowledge, the package audit
database/sidecars, coordination docs and launch input. The coordinator alone applies shared-memory and closure
changes serially; no child writes a sibling repository or user ledger.

**Contracts read.** Immutable exact prerequisite artifacts listed above, the
M51 native-owner-reviewed contract/fixtures when applicable, the assigned
package baseline and frozen shared-memory/consumer snapshots captured at
dispatch. Cross-package checks use read-only exact snapshots or approved
published module inputs, never changing sibling checkouts. Record full source,
artifact and fixture hashes in the later execution brief; contract drift pauses
affected leases for coordinator reconciliation. Existing no-workspace/no-directory
substitution requirements for ordinary published adoption remain in force.

**Parallel-safe.** no. Execute sequentially; this milestone is on a required integration/contract chain.
At most one live milestone per package. All tests use private lease ports,
disposable stores/caches/browser profiles and unique output directories.

## Dependencies and handoff

**Upstream.**

- Kinet:M54 (clean).
- Udon:M54.
- [OpenUdon:A32](../../../openudon/tabilet/docs/history/status-A32.md), which
  proves no OpenUdon HCL-input use remains.

**Downstream.**

- [APItools:M86](../../../apitools/tabilet/memory-bank/status-M86.md)
- [Kinet:M55](../../../kinet/tabilet/memory-bank/status-M55.md)

Ramen stays on its pinned UWS; older published versions stay immutable and
fetchable.

**A32 reconciliation.** Source `d339d6879eee21b82e31d39e88ed27dd89299a04` and normal
closure `c8ae90b0a802f4c8c57d7dcd226f45ca9a7739a5` are accepted/published and
independently verified. A32 removes the private/CLI and owned HCL-input closure while preserving public APIs
and history readers. C11.3 checks the retained public SDK through its existing
adapted disposable fixture. A32 qualification does not qualify C11: clean
Kinet:M54, accepted/published Udon:M54 and verified A32 normal closure
publication remain dispatch gates. Rows stay pending and review0; no trigger,
ordinary adoption or global HCL-free claim is added.

## Tasks

| Item | State | Notes |
|---|---|---|
| C11.1 — Versioned legacy core codec removal | `[X]` | Under UWS 1.13 §12.1 choose/document the compatibility version and remove every UWS core HCL input entrypoint: `convert.HCLToJSON`, `HCLToJSONIndent`, `HCLToYAML`, `UnmarshalHCL`, all `uws1` UnmarshalHCL hooks and deprecated `hcl.Import`. Remove the complete legacy core codec including `uws1.Document.MarshalHCL`, core `convert.MarshalHCL`/`JSONToHCL`/`YAMLToHCL` and Horizon-backed helpers. Preserve historical published versions and independent `uws/hcl` Render/Verify (`uws.hcl-view.v1`). **Withdrawn 2026-10-25:** owner-approved Stage 12 closure; avoid adding deployment/tooling work solely for conditional HCL cleanup. Deferred intent: Kinet STG-13, not scheduled. Original dependencies withdrawn together; no delivered acceptance. A replacement needs fresh approved IDs/authority. |
| C11.2 — HCL-free core and bounded JSON/YAML projection | `[X]` | Remove HashiCorp/Horizon HCL from UWS core selected/compiled module closures. Supply a bounded, context-cancellable `CanonicalJSON(context.Context, Source)` API accepting JSON/YAML only, preserving complete values, numeric lexemes (large integers, decimal/exponent spellings, signed zero) and extensions without lossy typed-model projection. Define limits/refusals and independent Render/Verify comparison tests. Keep the renderer in its separate module; this does not remove native Udon dependencies. **Withdrawn 2026-10-25:** owner-approved Stage 12 closure; avoid adding deployment/tooling work solely for conditional HCL cleanup. Deferred intent: Kinet STG-13, not scheduled. Original dependencies withdrawn together; no delivered acceptance. A replacement needs fresh approved IDs/authority. |
| C11.3 — Pre-publication migration fixtures | `[X]` | Build explicitly adapted disposable consumer fixtures at recorded source/patch identities for APItools, Udon after M54, OpenUdon’s retained public SDK after A32 (no retired CLI/internal fixture), Browsertools and Kinet workers. Exercise removal and replacement projection/output APIs without changing ordinary consumer worktrees or claiming their adoption. Record Kinet:M54 live inventory as separate readiness evidence. These pre-publication proofs precede C11 publication; APItools:M86 and Kinet:M55 own actual ordinary published adoption afterward. **Withdrawn 2026-10-25:** owner-approved Stage 12 closure; avoid adding deployment/tooling work solely for conditional HCL cleanup. Deferred intent: Kinet STG-13, not scheduled. Original dependencies withdrawn together; no delivered acceptance. A replacement needs fresh approved IDs/authority. |
| C11.4 — Qualify and publish | `[X]` | Docs, changelog, version records, immutability of earlier versions, and an exact publication handoff under named authority. **Withdrawn 2026-10-25:** owner-approved Stage 12 closure; avoid adding deployment/tooling work solely for conditional HCL cleanup. Deferred intent: Kinet STG-13, not scheduled. Original dependencies withdrawn together; no delivered acceptance. A replacement needs fresh approved IDs/authority. |

## Acceptance and verification

**Acceptance.** UWS core accepts no HCL input and has no legacy marshaling or
HCL selected/compiled dependency. The bounded JSON/YAML projection preserves
complete values and numeric lexemes. The independent renderer still renders and
verifies. Explicitly adapted disposable consumers build before publication;
M86/M55 later prove actual published adoption. Earlier published versions and
historical fixtures are unchanged; native Udon HCL remains outside this removal.

**Verification.** `go test ./...`, `go test -race ./...`, `go vet ./...`,
`GOWORK=off` module checks including `hcl/`, a module-graph proof,
`mkdocs build --strict` and `git diff --check`.

## Execution policy

One coordinator owns the integrated ledgers, shared memory and serialized
integration/closure. Serial execution remains the default. Concurrent leases
require this milestone's declared safety, frozen inputs, a complete explicit
Kinet goal request and the Stage 12 agent-rule opt-in. Each lease has one
in-progress row and one assigned milestone; its persisted review count survives
resume/rebase. Commit policy comes from that later request. Source publication,
deployment and live operations retain separate named authority. Planning and
status markers grant none; audit stays disabled.

## Review

Whole-milestone review: 0/10, not started.
