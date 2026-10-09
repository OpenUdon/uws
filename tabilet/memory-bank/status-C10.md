# C10 — Browser-profile shape activation

**Stage:** Kinet STG-12, Phase A. **Owner:** UWS.
**State:** Pending. No row has started. Review 0/10.
**Source baseline:** `989e3f2c88cac5c0f5a2911dfe04c36a61e43126` (clean at planning).
**Coordinator:** [Stage 12 contract](../../../kinet/docs/stage12.md). This
package-local milestone and status own acceptance. Planning was approved on
2026-10-09 (Kinet R65). It authorizes these planning files only.

**Parallel planning provenance.** Stage 12 parallel execution review,
2026-10-09; source priority and separate review baseline not supplied. Current
revalidation `989e3f2c88cac5c0f5a2911dfe04c36a61e43126`, including uncommitted planning files. F01 (confirmed local P2) adopts scoped workflow ownership; F02 (confirmed Lower) replaces strict ordering with readiness; F04 (confirmed local P2) requires frozen consumer checks.
Evidence: the old Kinet goal/launch rules, owning agent rules, existing pending
dependencies and sibling consumer checks; Udon `go.mod`/`pkg/execute` import no
OpenUdon code. This approved intake changes no row state or review counter.

## Dispatch and lease boundaries

**Depends on.** Kinet:M51. All required prerequisites must have
accepted closure at exact revisions; sibling producer adoption also needs
independently verified publication. A priority position never supplies authority.

**Downstream impacts.** Browsertools:M33, OpenUdon:P10, Udon:M53, Kinet:M56, Kinet:W20.

**Write set.** The owning `uws/` package's implementation, tests, ordinary
documentation, manifests and qualification outputs only as required by this
milestone's existing scope, plus `tabilet/memory-bank/status-C10.md` in its
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

**Parallel-safe.** yes. Eligible only under the explicit Stage 12 lease opt-in, with no dependency path or bidirectional read/write conflict against any running lease.
At most one live milestone per package. All tests use private lease ports,
disposable stores/caches/browser profiles and unique output directories.

## Dependencies and handoff

**Upstream.** [Kinet:M51](../../../kinet/tabilet/memory-bank/status-M51.md)
browser consumer contract. Build on accepted C09: `KindBrowser` is reserved in
`binding/types.go:15`, and browser templates are skipped in `binding/flow.go:311`
and `expressions/portability.go:43`.

**Downstream.**

- [Browsertools:M33](../../../browsertools/tabilet/memory-bank/status-M33.md)
- [OpenUdon:P10](../../../openudon/tabilet/memory-bank/status-P10.md)
- [Udon:M53](../../../udon/tabilet/memory-bank/status-M53.md)
- [Kinet:W20](../../../kinet/tabilet/memory-bank/status-W20.md)

Consumers adopt only the exact accepted and independently published revision.

## Tasks

| Item | State | Notes |
|---|---|---|
| C10.1 — Activate the browser-profile shape kind | `[ ]` | Define the `browser-profile` ShapeTable contract: profile version, action key as native selector, typed inputs and outputs, side-effect class, origins, authentication requirement and confirmation policy. Carry no HTTP method, path or server. Versioning is additive; existing `uws.shape-table.v1` bytes and APItools tables stay valid. |
| C10.2 — Browser binding and flow rules | `[ ]` | Binding validation, deterministic flow and strict portability for browser request templates, which are skipped today. Stable value-free diagnostic codes. Keep indeterminate results. Diagnostics are observation, not authorization. |
| C10.3 — Conformance vectors | `[ ]` | Positive and negative vectors for Browser 1.5–1.10 actions and for browser-authentication and browser-registration calls. Tampered, ambiguous and incomplete cases. |
| C10.4 — Qualify and publish | `[ ]` | Docs, published-version immutability, consumer builds (APItools, OpenUdon, Udon, Browsertools, Kinet workers) and an exact publication handoff. Publish only with named authority. |

## Acceptance and verification

**Acceptance.**

- Browser shapes validate and bind with stable codes and deterministic fixtures.
- Shape production stays outside UWS.
- No published profile, schema or `uws.shape-table.v1` byte changes.
- No diagnostic grants execution.

**Verification.**

- `go test ./...`, `go test -race ./...` and `go vet ./...`.
- `GOWORK=off go test ./...`, also run from `hcl/`.
- Schema, conformance and published-version immutability checks.
- `mkdocs build --strict` and `git diff --check`.
- Use only disposable roots and fixtures.

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
