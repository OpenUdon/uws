# Milestones

M01, M02, M03, C01, C02, B01, B02, C03, C05, C04, and M04 completed
acceptance and bounded review; their IDs remain reserved in the history index.
C03 published UWS 1.10.0 and C04 published UWS 1.11.0.

## Active Horizon

The second 2026-09-23 engineering review opened the required remediation
horizon. C05's exact-version admission correction, B02's Browser 1.9 profile,
C04's UWS 1.11 contract, and M04's evidence tracking are implemented, reviewed,
and retired in history. Planning does not implement other product findings;
archive display labels were clarified editorially and the reproducible
evidence is now clone-visible. No active required or conditional milestone
remains.

**Active order.** None. **Reconciled downstream impacts.** B02 -> C04 was
reconciled: UWS 1.11 references the completed Browser 1.9 contract. C04 had no
active product consumers; M04 was independent in the dependency graph and ran
after C04 for the approved goal order. All required milestones are now
complete, and no conditional milestone was triggered. Published prior
versions and the existing GOAL protocol remain unchanged.

## Candidate Directions

Candidate directions have no lane, permanent ID, status file, or execution
order. A trigger causes a fresh proposal and approval; it does not schedule the
work automatically.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| MCP public supplement consideration | The OpenUdon experiment is unimplemented and has no interoperability evidence. | Stage 1 produces real workflow evidence and a second independent consumer requests interoperable exchange. |
| Concrete content-trust resolvers | No source/profile resolver has a named in-repository owner or representative acceptance corpus. | A runtime or profile owner supplies reviewed channel contracts and fixtures. |
| UWS 2.0 expression, trigger, polling, and enforcement redesign | C2–C4 and E1/E2/E6/E8/E11 require new wire or governance choices: an expression marker/escape and interpolation, richer operators and names, literal `items`, decoupled profile/core versions, extensible source types, trigger kinds, possible content-trust enforcement, and portable `await` operation reexecution. C04 added numeric `batchSize` literals and corrected the existing `await` guide; portable operation reexecution remains deferred. | A concrete multi-runtime need and compatibility analysis support a separately approved 2.0 proposal. |
| Profile documentation | D7's remaining runtime-supplement ambiguity and D10's registration 1.2 authoring detail concern separately versioned profiles; editing frozen published documents is not an automatic review fix. Browser BCP 14 language belongs to B02. | The next relevant profile version or an explicitly approved meaning-preserving editorial amendment. |
| Interoperability formats | D6 file extensions, C16/E9 content-trust wire reports/resolvers, E4 stable error codes, E10 normative HCL mapping, and C10 portable error taxonomy require independent consumer and compatibility evidence beyond C04's executable conformance corpus. | A named independent consumer or portable conformance requirement and separately approved contract. |
| `uws.*` profile-name namespace reservation | C18's proposed reservation is a governance change, not a correction to the present core list of `x-uws-*` fields. | An approved namespace/governance or 2.0 design with migration analysis. |
| Evaluation-cost hardening | N16 identifies repeated recursive truthiness scans and expression-pattern compilation, but supplies no measured cost or representative workload. | A reproducible benchmark shows material latency or resource impact and a compatible cache/validation design is approved. |

The second review established a concrete Browser 1.8 gap and promoted that
work to B02, together with the separate Browser 1.9 text-safety contract. It
also promoted core defects to C04 and evidence reproducibility to M04. C05's
urgent version-admission correction was completed first, followed by B02, C04,
and M04. MCP still lacks Stage 1 evidence and a second consumer; the remaining
candidate directions require their own promotion triggers and approval.

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
