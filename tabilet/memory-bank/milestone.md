# Milestones

M01, M02, M03, C01, C02, and B01 completed acceptance and bounded review; their IDs
remain reserved in the history index. C03 remains active from the
2026-09-23 review remediation.

## Active Horizon

Required remaining execution order: **C03**. C03 depends on completed M03,
C01, C02, and B01. B01's downstream profile cross-reference has been
reconciled. C03 is required, not conditional.

| Milestone | Goal | Status |
|---|---|---|
| [C03](status-C03.md) | Portable execution contract | Pending |

## C03 - Portable Execution Contract

**Goal.** Make the executable UWS contract reviewable and implementable beyond
the Go reference implementation without changing older published artifacts.

**Scope.** Resolve the confirmed specification/code mismatches and ambiguities
in its [status](status-C03.md), including non-`await` wait, criteria, expression
and result shapes, structural/control-flow semantics, actions, triggers,
security guidance, conformance classes, references, and examples. The
orchestrator owns non-`await` wait as bounded numeric seconds; `await` retains
its predicate meaning. Test compatibility before assigning this interpretation
to older documents. Define currently implemented behavior where intentional;
change implementation only where an approved portable contract requires it.
Publish the appropriate new versioned core specification/schema, changelog,
archive/parity updates, and documentation without modifying immutable history.

**Acceptance.** Spec, schema, Go validation/execution, changelog, conformance
fixtures, and embedded archive agree; prior versions are protected. Full and
race tests, vet, strict MkDocs build, and diff check pass. The bounded
whole-milestone review gate passes.

**Dependencies.** [M03](../docs/history/status-M03.md),
[C01](../docs/history/status-C01.md), [C02](../docs/history/status-C02.md),
and [B01](../docs/history/status-B01.md) completed and accepted. Browser 1.8
is a separate UWS 1.9 profile; C03 should update references only, not import
its substitution semantics into core. C01
establishes deterministic caller-step identity for workflow
invocations, per-call scoping of nested steps, operations, dependencies, and
merge records, clear rejection of recursive workflow calls, and numeric
ordering of nested iteration results. C02 establishes exact declared-version
schema selection, no fallback for prerelease or unpublished versions, and
SemVer-aware feature gates, including 1.5 step inputs and 1.9.2 semantic
changes. **Downstream impacts.** None currently active; candidate directions
below may be reconsidered after C03.

## Candidate Directions

Candidate directions have no lane, permanent ID, status file, or execution
order. A trigger causes a fresh proposal and approval; it does not schedule the
work automatically.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| MCP public supplement consideration | The OpenUdon experiment is unimplemented and has no interoperability evidence. | Stage 1 produces real workflow evidence and a second independent consumer requests interoperable exchange. |
| Concrete content-trust resolvers | No source/profile resolver has a named in-repository owner or representative acceptance corpus. | A runtime or profile owner supplies reviewed channel contracts and fixtures. |
| New browser or account-lifecycle profile direction beyond browser 1.8 | Browser 1.8 resolves the confirmed template-safety gap; broader profile or account-lifecycle changes lack a proved portable contract. | A new concrete gap beyond browser 1.8 is demonstrated; portable semantics additionally require multi-runtime demand. |
| UWS 2.0 expression, trigger, and enforcement redesign | C2–C4 and E1/E2/E6/E8/E11 require new wire or governance choices: explicit expression marker/escape and interpolation, richer operators and names, literal `items`/`batchSize`, decoupled profile/core versions, extensible source types, trigger kinds, and possible content-trust enforcement. The review does not establish a compatible 1.x design or an injection exploit. | A concrete multi-runtime need and compatibility analysis support a separately approved 2.0 proposal. |
| Profile documentation | D7's remaining runtime-supplement ambiguity, D9's BCP 14 declarations, and D10's registration 1.2 authoring detail concern separately versioned profiles; editing frozen published documents is not an automatic review fix. | The next relevant profile version or an explicitly approved meaning-preserving editorial amendment. |
| Interoperability formats | D6 file extensions, C16/E9 content-trust wire reports/resolvers, E4 stable error codes, E10 normative HCL mapping, E3 fixture expansion, and C10 portable error taxonomy require independent consumer and compatibility evidence beyond C03's core semantics. | A named independent consumer or portable conformance requirement and separately approved contract. |
| `uws.*` profile-name namespace reservation | C18's proposed reservation is a governance change, not a correction to the present core list of `x-uws-*` fields. | An approved namespace/governance or 2.0 design with migration analysis. |

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
