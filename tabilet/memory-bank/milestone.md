# Milestones

## Stage 11 post-acceptance remediation — 2026-10-08

The user approved the complete reconciliation proposal and its planning-file application. Five new package-local milestones / 28 task rows restore accepted supported contracts. Original Stage 11 acceptance and all completed records remain frozen; the planning intake started no closing review.

**Serial order:** Udon:M49 → UWS:M09 → APItools:M83 → OpenUdon:M99 → Kinet:M49.
**Direct dependency graph:** UWS:M09 → APItools:M83, OpenUdon:M99, Kinet:M49; APItools:M83 → OpenUdon:M99, Kinet:M49; OpenUdon:M99 → Kinet:M49; Udon:M49 → Kinet:M49. Serial order is a scheduling gate, not an additional library dependency.

**Local owner:** M09 — Stage 11 parser and binding contract remediation; four implementation rows complete, M09.5 in progress. The confirmed serial GOAL authorizes implementation/task commits; the linked status owns the single persisted closing-review counter and findings. One serial execution owner across the five ledgers. Planning itself authorized no execution. Fresh separately named source publication authority remains required; consumed Stage 11 publication/deployment envelopes are not reusable.

**Coordinator:** [Stage 11](../../../kinet/docs/stage11.md#post-acceptance-remediation--2026-10-08). Exact accepted/published successor revisions remain unset until independently observed; downstream adoption requires them. Installed M44 and independent browser/media/Phase A/legacy/frozen-consumer pins remain unchanged. Stage 12 remains provisional.

## Stage 11 active horizon

Approved 2026-10-06: both phases of [Kinet STG-11](../../../kinet/docs/stage11.md), with one serial execution owner across Kinet, UWS, APItools, Udon and OpenUdon. All UWS Stage 11 rows are accepted/retired; M08 accepted/published `c0b19385a3b034cd45de16726668b9150f0633f2` after review 3; C09 is accepted/retired at 6a267306032edc687a298cefc8bba7019d3ad059 after review 3. C08 is accepted/retired at 0411eea6fc84fbd6aa97cef94f53f301260f4844 after review 2. The [history index](../docs/history/index.md) retains the complete accepted status/specification; downstream adoption remains package-local. Earlier completed horizons and records remain historical; planning grants no implementation or external authority.

M01, M02, M03, C01, C02, B01, B02, C03, C05, C04, M04, C06, C07, M06, and M07
completed acceptance and bounded review; their IDs remain reserved in the
history index. C03 published UWS 1.10.0, C04 published UWS 1.11.0, and C07
committed UWS 1.12.0 at `9d092664a6062563e0414527f997a2475aeab003`.

## Active Horizon

The third engineering review's [C06](../docs/history/status-C06.md) amendment
to the published UWS 1.11 specification passed verification and its bounded
review gate, then retired. Earlier milestones also remain in the
[history index](../docs/history/index.md).

M05 published and verified Browser 1.10 at commit
`80ee9bfb24a688b5e875dadf9ecacdc65398f1ff`, passed bounded review iteration 1,
and was retired to the [milestone history](../docs/history/status-M05.md).
Browser 1.9 and older published profiles remain immutable; Browser 1.10 is the
current opt-in profile. UWS 1.13.0 is the current core contract, qualified/published by M08. Earlier
UWS 1.12 and browser contracts remain immutable. Browsertools,
Browserdriver, Udon, OpenUdon, and W8M own their downstream compatibility work
in their respective repositories and active ledgers. The GOAL protocol remains
unchanged.

M07 corrected mock-runtime enclosing-step lookup, passed its whole-milestone
review in iteration 1, and was retired to the [milestone history](../docs/history/status-M07.md). M06 completed versioned fixtures and the public
mock/hybrid runtime and remains in the [milestone history](../docs/history/status-M06.md). C07 passed
its whole-milestone review in iteration 1 and is archived at
[`status-C07.md`](../docs/history/status-C07.md). The user approved the C07 ->
M06 horizon on 2026-09-26. No sibling implementation is included in the UWS
milestone; sibling adoption remains in each package's own ledger.

### S1 Provenance And Downstream Boundaries

The approved source is Kinet `docs/icot.md` §7 UWS, with sequencing in §8.
It was inspected as working-tree evidence at Kinet HEAD
`9719a364904743c0d2cb8b3f47cb4e50ed6322e8`; the document contains uncommitted
edits and has SHA-256
`7d2903d0ef70f046e9b711933749791271e77c53622e258b47a9fb4a6c610fe6`.
UWS was clean at `1d5535ec75d5693a66bcced5bd98f5c4c824fb2a`.
The [evolution v6 direction](../evolution/prompt-v6.md) and
[planning result](../evolution/result-v6.md) record this material change.

C07 supplies the effect/pending contract, including the shared step-contract
field shape adopted by OpenUdon M87.1. M06 delivered versioned mock fixtures
and the public pure/hybrid runtime. OpenUdon owns simulation commands,
assessment, and approval refusal; Udon owns real-read enforcement and
credentials; APItools owns classification discovery; Browsertools owns snapshot
simulation. Kinet owns the authoring loop and later W04 adoption. M06 does not
block Kinet W03; W03 still requires released OpenUdon step-command schemas,
fixtures, and commands. Downstream adoption is recorded and executed in each
package's own ledger. Existing candidates remain deferred, including broader
interoperability formats and optional expression portability tooling. Ordinary
UWS validation continues accepting implementation-specific expressions.

## Candidate Directions

Candidate directions have no lane, permanent ID, status file, or execution
order. A trigger causes a fresh proposal and approval; it does not schedule the
work automatically.

| Direction | Why Deferred | Promotion Trigger |
|---|---|---|
| MCP public supplement consideration | The OpenUdon experiment is unimplemented and has no interoperability evidence. | Stage 1 produces real workflow evidence and a second independent consumer requests interoperable exchange. |
| Concrete content-trust resolvers | No source/profile resolver has a named in-repository owner or representative acceptance corpus. | A runtime or profile owner supplies reviewed channel contracts and fixtures. |
| UWS 2.0 expression, trigger, polling, and enforcement redesign | C2–C4 and E1/E2/E6/E8/E11 require new wire or governance choices: an expression marker/escape and interpolation, richer operators and names, literal `items`, decoupled profile/core versions, extensible source types, trigger kinds, possible content-trust enforcement, and portable `await` operation reexecution. C04 added numeric `batchSize` literals and corrected the existing `await` guide; portable operation reexecution remains deferred. | A concrete multi-runtime need and compatibility analysis support a separately approved 2.0 proposal. |
| Other profile documentation | D7's remaining runtime-supplement ambiguity and D10's registration 1.2 authoring detail remain separate from Browser 1.10. M05 published the Browser profile versioning procedure and count profile; earlier profile versions remain immutable. | A relevant later profile version or an explicitly approved meaning-preserving editorial amendment for the remaining topics. |
| Authoring diagnostics and examples | Third-review R3/R10: numeric `wait`/`batchSize` tokens must stay quoted under the existing string wire shape; raw-number diagnostics and the runtime-specific `$error.*` guide excerpt could be clearer. | A named authoring consumer and approved diagnostics or guide amendment with fixtures. Unquoted numbers require a separately versioned wire decision. |
| Remaining conformance corpus supplement | Stage 11 C08/C09 now own separately pinned expression/binding vectors. Other third-review R9 gaps remain deferred; changing the frozen 1.11 corpus would alter published evidence. | A named need for the remaining behaviors and a separately approved supplement or later release corpus. |
| Browser text-safety expansion | Third-review R8: the Browser 1.9 text rule omits some Unicode `Cf` characters; a multilingual-safe replacement rule is undecided. | A reviewed allowlist or profile-version proposal with Unicode safety and compatibility evidence. |
| Browser default migration | Third-review R12: empty profile lookup intentionally selects Browser 1.8 for compatibility, while Browser 1.10 remains opt-in. | Caller migration analysis and an approved default/deprecation policy. |
| Interoperability formats | D6 file extensions, C16/E9 content-trust wire reports/resolvers, E4 stable error codes, E10 normative HCL mapping, and C10 portable error taxonomy require independent consumer and compatibility evidence beyond C04's executable conformance corpus. | A named independent consumer or portable conformance requirement and separately approved contract. |
| `uws.*` profile-name namespace reservation | C18's proposed reservation is a governance change, not a correction to the present core list of `x-uws-*` fields. | An approved namespace/governance or 2.0 design with migration analysis. |
| Evaluation-cost hardening | N16 identifies repeated recursive truthiness scans and expression-pattern compilation, but supplies no measured cost or representative workload. | A reproducible benchmark shows material latency or resource impact and a compatible cache/validation design is approved. |
| URL participation in shape source identity (stage11-siblings-review.md uws P3.6; source P3/local Lower; UWS binding owner) | binding/table.go compares the full Source, including URL. Provenance-only wording forbids fetch authority but does not establish URL-independent identity. | A named consumer and approved identity contract decide whether URL differences represent the same source; preserve exact hashes/native selectors and no-fetch rules. |
| Literal versus self-referential input expressions (stage11-siblings-review.md uws P3.7; source P3/local Lower; UWS expression owner) | expressions/evaluate.go bounds self-reference but the supported literal/expression interpretation is not established by this review. | A named consumer needs an explicit literal/escape/self-reference policy with compatibility fixtures; reconcile before changing evaluation. |
| Non-NFC presentation support (stage11-siblings-review.md uws P3.3 normalization component; source P3/local Lower; UWS codec owner) | cty normalization causes fail-closed Render refusal; M09 documents/tests the limitation without promising arbitrary Unicode normalization. | A named consumer requires lossless NFD text and an approved codec representation can retain exact values/numeric lexemes without changing frozen view or authority contracts. |

The second-review remediation (C05, B02, C04, M04) is complete. The third
review's two normative UWS 1.11 contradictions were corrected in C06; the
other improvement directions remain unnumbered. MCP still lacks Stage 1
evidence and a second consumer; candidates require fresh promotion approval.

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

## Stage 11 cross-package refactoring

Approved review-intake amendment, 2026-10-06: 18 required milestones / 87 pending rows across the five owners. Scope and milestone IDs are unchanged; the coordinator records the approved publication proposal, which needs a separate execution-time grant. Full local source baseline remains `a7688f54c68f5a75c7cc95aa2b31cea98b31af41`, including the reviewed uncommitted planning state. This intake does not start a closing review or authorize implementation.

**Approved source.** User-approved complete proposal, 2026-10-06; source baseline `a7688f54c68f5a75c7cc95aa2b31cea98b31af41`. [Coordinated contract](../../../kinet/docs/stage11.md) defines both phases, cross-package order, compatibility and acceptance. The request to implement the proposal authorizes its planning files only.

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

### Stage 11 status index

| ID | Milestone | Status file | State |
|---|---|---|---|
| M09 | Stage 11 parser and binding contract remediation | [status-M09.md](status-M09.md) | M09.5 in progress; review in linked status |

Completed IDs and full accepted records remain reserved in history.

## Stage 11 candidate dispositions

C08 promotes opt-in expression portability. C08/C09 promote the expression/binding subset of a separately pinned conformance supplement; the frozen 1.11 corpus and broader remaining conformance work stay unchanged. M08 promotes only the additive HCL presentation/deprecation subset of interoperability. New grammar, trigger/polling redesign, general diagnostic wire standardization and HCL removal remain deferred.

## M09 — Stage 11 parser and binding contract remediation

**Stage/owner.** STG-11 post-acceptance remediation; uws. Approved planning 2026-10-08. **Placement.** Core/cross-cutting lane; new remediation, never a reopened historical gate.
**Lineage.** [C08](../docs/history/status-C08.md), [C09](../docs/history/status-C09.md) and [M08](../docs/history/status-M08.md).
**Dependencies.** Accepted/published C08/C09/M08 contracts. Serial scheduling follows Udon:M49; this is not a new runtime dependency. Root and nested codec qualification/publication must be independently resolved before APItools:M83, OpenUdon:M99 and Kinet:M49 adoption.
**Scope.** Bound untrusted HCL and shape JSON before recursive parsing, enforce deterministic verified-view bytes, make deprecated inert import structurally symmetric, repair supported binding containment/path/draft context and scoped flow/strict-portability diagnostics. Keep existing published grammar/wire/schema and frozen conformance bytes unchanged.
**Acceptance.** Inputs within byte limits but above depth 100 refuse before unbounded HCL/strict-JSON recursion; comments, strings and template/interpolation delimiters are handled correctly by preflight. Verify accepts only exact deterministic Render bytes plus independent value/numeric-lexeme proof, without recursive Render/Verify calls. Deprecated Import rejects invalid typed blocks; non-NFC refusal stays fail-closed and is documented. Pattern-property containment, false leaves, nullable/pattern paths and inherited draft/reference contexts never produce unsupported compatibility or false missing-field claims. Flow references respect node kind/workflow scope; opt-in portability traverses trigger routes and diagnoses absent iteration context without narrowing ordinary Parse/validation. Root and codec standalone builds, immutable version/corpus guards and owner review pass.
**Verification.** Root go test ./... and go vet ./...; affected binding/expressions/strictjson race tests; separate hcl-module tests/vet/races with GOWORK=off GOPROXY=off; adversarial-depth subprocess fixtures; canonical-comment/escape/numeric/NFD/import symmetry cases; whole-schema versus nested draft-07 array proof; flow namespace and trigger-loop fixtures; immutable published versions/schema/conformance digest checks; exact root/nested-codec archive and ordinary downstream-consumer qualification; git diff --check.
**Compatibility/recovery.** Preserve public wires/schemas, declared grammar versions and frozen evidence/pins. Corrected derived tables/packages/workers require fresh consumer assessment, confirmation and grants; never upgrade historical authority or replay unknown writes. Installed M44 is unchanged. Rollout ends at a new exact-source qualified handoff; real installation, migration, sends and registration need separate named authority.
**Downstream.** APItools:M83 native shapes, OpenUdon:M99 package/source verification and Kinet:M49 author/private-exec workers. Retained browser profiles and frozen consumers remain separately pinned.
**Tasks.** 5 task/commit units in status-M09.md: M09.1 Bound HCL and shape-table parsing before recursion; M09.2 Verify canonical HCL views and retain inert import symmetry; M09.3 Restore sound binding containment and schema context; M09.4 Scope flow references and strict portability contexts; M09.5 Qualify and hand off exact root and nested codec revisions.
**Review/authority.** Closing review under confirmed execution; the linked status owns the single persisted counter and findings/fixes. Publication is an external prerequisite requiring fresh named authority and independent resolution. Planning itself grants no code execution, commit, publication, deployment or goal launch.
