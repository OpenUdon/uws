# Milestones

The active horizon now protects the accepted current-release Markdown surface
from later drift. Work follows the order below, one pending status row and one
task commit at a time.

## Active Horizon

| Milestone | Goal | Dependencies | Status |
|---|---|---|---|
| M02 | Add an automated immutability guard for published Markdown artifacts. | [M01](../docs/history/status-M01.md) completed and accepted | [status-M02.md](status-M02.md) |

Execution order: `M02`.

## M02 - Published Markdown Immutability Guard

**Goal.** Detect changes, removals, and unregistered additions among published
Markdown documents under `versions/`, while leaving the release changelog
intentionally mutable.

**Scope.** Extend `schemas/version_immutability_test.go` with SHA-256 and exact
membership coverage for every `versions/*.md` file except
`versions/CHANGELOG.md`, using the final M01 bytes. Update `AGENTS.md`,
`architecture.md`, and `tech-stack.md` in the same task to describe the new
current rule. Do not modify the protected Markdown documents merely to satisfy
their hashes.

**Acceptance.** `go test ./schemas -run
TestPublishedVersionDocumentsAreImmutable`, `go test ./...`, `go test -race
./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` pass.
The focused test proves exact non-changelog Markdown membership and SHA-256
bytes. Review confirms that no versioned JSON document or embedded archive
changed.

**Dependencies.** [M01](../docs/history/status-M01.md) completed and accepted.

**Downstream impacts.** none in the active horizon.

**Task breakdown.** Add the Markdown digest/membership guard and reconcile the
current instructions and memory-bank facts in one commit.

## Candidate Directions

Candidate directions have no lane, permanent ID, status file, or execution
order. A trigger causes a fresh proposal and approval; it does not schedule the
work automatically.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| MCP public supplement consideration | The OpenUdon experiment is unimplemented and has no interoperability evidence. | Stage 1 produces real workflow evidence and a second independent consumer requests interoperable exchange. |
| Concrete content-trust resolvers | No source/profile resolver has a named in-repository owner or representative acceptance corpus. | A runtime or profile owner supplies reviewed channel contracts and fixtures. |
| New browser or account-lifecycle profile version | No unmet portable semantic requiring another version is established. | Multi-runtime evidence demonstrates a portable contract gap that existing versions cannot express. |
| UWS 2.0 enforcement or wire changes | The current outcome is backward-compatible documentation and tooling maintenance. | Evidence shows required behavior cannot be delivered compatibly in UWS 1.x. |

## Review Finding Severity

P1 and P2 are engineering-review priorities, not product-domain terms,
milestone priority, or status markers. A linked project-specific review policy
or `AGENTS.md` definition overrides these defaults.

- **P1:** a likely defect with critical or broad impact, such as corrupting a
  published contract, bypassing a security boundary, breaking compatibility,
  or making the accepted outcome unsafe or unusable. Any confirmed P1 blocks
  milestone acceptance.
- **P2:** a material correctness, compatibility, security, or operability defect
  in normal supported use. Any confirmed P2 blocks milestone acceptance.

Classify findings by impact, likelihood, and affected scope rather than the
size or convenience of the fix. Lower-severity improvements do not block unless
the milestone's own acceptance contract requires them.

## New Review Intake

For an engineering review received after initialization:

1. Revalidate every finding against the current tree and active specifications.
2. Preserve the review's source severity and record the locally assessed
   severity separately in affected milestone or status notes.
3. Propose dispositions and every file action, and obtain approval before
   writing planning changes.
4. Put confirmed in-scope work into an open or pending owning milestone. Never
   reopen completed history; create a remediation milestone with lineage.
5. Add P1, P2, or higher confirmed work to the dependency-closed active
   horizon. Keep optional lower-severity work as an unnumbered candidate.
6. Record portable provenance in the affected notes without copying the review
   or creating a separate review ledger.

A review counts as a bounded-gate iteration only when it was explicitly
requested as the next pass of an already active persisted milestone gate.

## Milestone Review Gate

Every milestone receives a whole-milestone deep review after its task rows and
acceptance commands are complete. Before each pass, persist the iteration number
in its status notes. The first pass is iteration 1; session or reviewer changes
never reset the count, and an interrupted pass resumes at the same number.

After every P1, P2, or higher-severity fix, rerun affected verification and
review the whole milestone again. The gate passes only when a complete review
finds no P1, P2, or higher issue. The limit is 10 iterations. If iteration 10
still finds a blocker, keep the milestone active and add or update a `[!]`
status row with the finding, owner/source, impact, and unblock condition.

## Closure, Retirement, And Retrieval

After the review gate passes within 10 iterations, all acceptance evidence
passes, current facts and applicable lessons are consolidated, and downstream
work is reconciled, automatically retire the milestone. Terminal row markers
alone are insufficient, and unresolved or triggered conditional work keeps it
active.

Retirement creates a history record named `status-ID.md` under
`tabilet/docs/history/`, where `ID` is the permanent milestone ID. It contains
the complete final milestone specification and status document in separate
literal `markdown` fences. Before those sections, record single-line fields for
Milestone as the permanent ID; Outcome as `completed`, `cancelled`, or
`superseded`; Retired as a UTC `YYYY-MM-DD` date; Source status as its original
memory-bank path; Source specification as the original milestone path and
heading anchor; Evidence as the full Git commit or `unversioned`; Worktree as
`clean`, `includes uncommitted changes`, or `unversioned`; Review as `passed`;
Review iterations from 1 through 10; Verification; and Consolidated into with
current-document or lesson links, or explicit `no current-truth change`. Use a
fence longer than any fence inside the retained source. Obtain Git evidence
with `git rev-parse --verify HEAD`; never imply that a commit contains
uncommitted work. Validate the complete envelope and retained source bytes
before removing active records.

Cancelled or superseded outcomes also record Disposition with authority,
rationale, and dependency treatment; supersession records Successor. Neither
automatically satisfies a completion dependency. Every historical `[-]` row
names its accepted successor and is never retried.

Create `tabilet/docs/history/index.md` only on first retirement, with Milestone,
Outcome, Retired, Record, and Summary columns. Remove the retired status file,
its active index row, and its specification together, repair maintained links,
and preserve the history record and index entry permanently. Index IDs and
dates must match record metadata, and Record uses a relative link. Resolve IDs
across active and retired locations and never reuse them. Frozen archives and
evolution snapshots retain their original paths and are resolved through the
history record's provenance. Read history only for a relevant dependency or
question.

Before materially replacing or removing a current fact or lesson, append its
source heading and literal old wording, reason, evidence, and replacement link
under a unique dated heading in `tabilet/docs/history/knowledge.md`; link that
journal from the history index. Merge duplicate lessons and remove obsolete
ones only after preserving that evidence. Ordinary editorial changes need no
entry. A retirement needs no separate archive run. Follow the governing commit
policy, then refresh disposable goal input for remaining work or remove it when
the active horizon is empty.
