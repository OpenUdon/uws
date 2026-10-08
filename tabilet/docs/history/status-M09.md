# Retired milestone M09 — Stage 11 parser and binding contract remediation

**Milestone.** M09
**Outcome.** completed
**Retired.** 2026-10-08
**Source status.** tabilet/memory-bank/status-M09.md
**Source specification.** tabilet/memory-bank/milestone.md#m09--stage-11-parser-and-binding-contract-remediation
**Evidence.** 51a74545b016b8ab75e30d454338c52f7945a836
**Worktree.** includes uncommitted changes
**Review.** passed
**Review iterations.** 7
**Verification.** Root/codec ordinary module-only full suites/races/vet/build, exact source/origin/version/file-content and archive/GoMod sum proof, 52/50/52 selected closures, public consumer races, strict docs, protected versions/corpora/pins, whitespace and retirement-envelope checks passed; whole review 7 passed.
**Consolidated into.** [product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), [technical stack](../../memory-bank/tech-stack.md), [lessons](../../memory-bank/lessons.md), [knowledge](knowledge.md), [publication](../../../docs/m09-publication.md), [qualification](../../../docs/m09-qualification.md); APItools:M83, OpenUdon:M99 and Kinet:M49 exact prerequisite reconciliation. Retirement metadata includes coordinator closure changes after the published evidence baseline. No new evolution version or live authority.

## Milestone specification

``````markdown
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

**Accepted result.** All five task units and whole review 7/10 pass. Root/codec
runtime b099f6803277ae94c7e9f1da0904a0140b278f20 independently resolves as
v0.0.0-20261008043726-b099f6803277; ordinary module/consumer proofs, full selected
closures and downloaded full suites/races/vet/builds pass. M09.5 evidence
51a74545b016b8ab75e30d454338c52f7945a836 is independently published. Current
knowledge and exact pending consumer assumptions are consolidated; normal
retirement closure publication precedes consumer execution. Frozen historical
contracts/pins and installed M44 remain unchanged; no live authority follows.
``````

## Status record

``````markdown
# M09 — Stage 11 parser and binding contract remediation

**Stage:** STG-11 post-acceptance remediation. **Owner:** UWS.
**State:** All five rows complete; whole review 7/10 passed; accepted and retired on 2026-10-08 after ordinary publication proof and exact consumer reconciliation.
**Authority:** The confirmed GOAL request names Udon:M49 → UWS:M09 → APItools:M83
→ OpenUdon:M99 → Kinet:M49, task commits and no external mutations. The separately
granted UWS-only normal main source/closure publication exception is recorded
below; it grants no deployment, installed-consumer change or live run.
**Review source:** stage11-siblings-review.md — Stage 11 code review — sibling packages; uws section.
**Review baseline/range:** `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` → `0a4597122a7baa7e79e46e79e4ec60dbfffc3720`.
**Revalidation HEAD:** `0a4597122a7baa7e79e46e79e4ec60dbfffc3720`; clean worktree, no relevant uncommitted code in the evidence. Approved planning changes are not implementation evidence.
**Lineage:** [C08](../docs/history/status-C08.md), [C09](../docs/history/status-C09.md) and [M08](../docs/history/status-M08.md). Existing acceptance, review counters, statuses and frozen evidence stay preserved.
**Coordinator:** [Stage 11](../../../kinet/docs/stage11.md#post-acceptance-remediation--2026-10-08); package-local specification and status own acceptance.

## Dependencies and handoff

Accepted/published C08/C09/M08 contracts. Serial scheduling follows Udon:M49; this is not a new runtime dependency. Root and nested codec qualification/publication must be independently resolved before APItools:M83, OpenUdon:M99 and Kinet:M49 adoption.

**Exact successor publication/build identities:** independently observed runtime
b099f6803277ae94c7e9f1da0904a0140b278f20, root/codec
v0.0.0-20261008043726-b099f6803277; initial public evidence head
af225f8b0b5cfbd5848e28e5ce49eaebc0f228bb. Actual ordinary artifact/closure
identities are recorded in [publication](../../docs/m09-publication.md).
All five rows and whole review 7 pass. Coordinator consolidation and exact
consumer reconciliation close acceptance; the full retired record preserves
source/publication proof without substituting local HEAD or replacements.

**Accepted root/codec source:** b099f6803277ae94c7e9f1da0904a0140b278f20,
v0.0.0-20261008043726-b099f6803277, whole review 7/10 passed. Local/bootstrap
qualification remains separate from independent ordinary module/consumer proof
in [publication](../../docs/m09-publication.md). Published M09.5 evidence head
51a74545b016b8ab75e30d454338c52f7945a836 was independently observed before
coordinator closure; final retirement publication is observed before adoption.

**Downstream:** APItools:M83 native shapes, OpenUdon:M99 package/source verification and Kinet:M49 author/private-exec workers. Retained browser profiles and frozen consumers remain separately pinned.

## Scope and acceptance

Bound untrusted HCL and shape JSON before recursive parsing, enforce deterministic verified-view bytes, make deprecated inert import structurally symmetric, repair supported binding containment/path/draft context and scoped flow/strict-portability diagnostics. Keep existing published grammar/wire/schema and frozen conformance bytes unchanged.

Inputs within byte limits but above depth 100 refuse before unbounded HCL/strict-JSON recursion; comments, strings and template/interpolation delimiters are handled correctly by preflight. Verify accepts only exact deterministic Render bytes plus independent value/numeric-lexeme proof, without recursive Render/Verify calls. Deprecated Import rejects invalid typed blocks; non-NFC refusal stays fail-closed and is documented. Pattern-property containment, false leaves, nullable/pattern paths and inherited draft/reference contexts never produce unsupported compatibility or false missing-field claims. Flow references respect node kind/workflow scope; opt-in portability traverses trigger routes and diagnoses absent iteration context without narrowing ordinary Parse/validation. Root and codec standalone builds, immutable version/corpus guards and owner review pass.

## Tasks

| Item | State | Notes |
|---|---|---|
| M09.1 — Bound HCL and shape-table parsing before recursion | `[+]` | Iterative lexical preflight counts delimiters and refuses unsupported active templates before HCL parsing; strict JSON uses an iterative depth-100 stack. Exact duplicate/Unicode/numeric/trailing checks remain. Fatal-depth Import, strict-JSON and ParseTable child processes passed with a 1 MiB stack. Focused root tests/races/vet and separate full codec tests/races/vet passed, GOWORK=off GOPROXY=off Go 1.26.6. Sources W1/W2. |
| M09.2 — Verify canonical HCL views and retain inert import symmetry | `[+]` | Shared nonrecursive writer path enforces exact deterministic bytes before independent complete-value/numeric-lexeme proof. Self-hashed comment/escape/formatting tests, typed-block empty/null/mismatch symmetry and JSON/YAML NFD string/key fail-closed tests pass. Deprecated is a separate Go doc paragraph; README states the limitation. Separate full codec races and vet passed. Sources W6 and P3.2/P3.3/P3.8. |
| M09.3 — Restore sound binding containment and schema context | `[+]` | Original compiled child nodes preserve dialect/local-reference context. Pattern-driven closure and nullable/pattern paths are indeterminate where unsupported; false leaves refuse. Draft-07 tuple/future-keyword/local-reference tests agree with whole-schema literal validation; binding full tests/races/vet passed. The original dependencies claim remains unconfirmed. Sources W3/W4/W5/P3.1. |
| M09.4 — Scope flow references and strict portability contexts | `[+]` | Explicit identities retain generic group barriers and step→workflow→operation dependency precedence; genuine duplicate step definitions remain ambiguous. Output references use explicit workflow/result owners and known operation invocation contexts, with no foreign-name fallback. Strict checker separates loop/item availability and pre-forEach controls, including transitive trigger/dependency contexts. Root full tests/races/vet pass; ordinary Parse remains compatible. Sources P3.4/P3.5. |
| M09.5 — Qualify and hand off exact root and nested codec revisions | `[+]` | Separately authorized normal initial af225f8 publication independently observed; runtime b099f6803277ae94c7e9f1da0904a0140b278f20 root/codec independently resolve with exact version/origin/sums/file bytes. Downloaded module-only full suites/races/vet/build and ordinary 52-module consumer proof pass. Whole review 7/10 passed; all five rows qualified. Coordinator owns downstream reconciliation and final acceptance/retirement. No deployment or live authority follows. |

## Active finding provenance

Only approved active findings are recorded here. Source priorities and local severity are separate; task references identify one owner for each required outcome. Unsupported and unscheduled findings remain in the conversational handoff; optional directions belong only in milestone.md.

| Source finding | Source priority | Local severity | Disposition | Current repository evidence | Task owner |
|---|---|---|---|---|---|
| W1 | P2 | P2 | confirmed | hcl/parse.go:28 calls hclsyntax.ParseConfig before depth-100 budget | M09.1 |
| W2 | P2 | P2 | confirmed | binding/table.go:180; internal/strictjson/strictjson.go:107 recursively consumes unbounded nesting | M09.1 |
| W3 | P2 | P2 | confirmed | binding/template.go objectContainment ignores patternProperties; false-compatible closed-object probe | M09.3 |
| W4 | P2 | P2 | confirmed | binding/validate.go schemaPath accepts false leaf at #/x; offline probe | M09.3 |
| W5 | P2 | P2 | confirmed | binding/validate.go schemaPath rejects reachable nullable/pattern fields; offline probes | M09.3 |
| W6 | P2 | P2 | confirmed | hcl/codec.go:101 semantic-only verification accepts a self-hashed added approval comment | M09.2 |
| P3.1 | P3 | P2 | partially confirmed | binding/template.go schemaFrom loses root draft; offline draft-07 tuple/future-keyword probes diverge from whole-schema validation; dependencies example did not reproduce | M09.3 |
| P3.2 | P3 | Lower | confirmed; canonical codec acceptance | hcl/codec.go:109 Deprecated text is not its own paragraph | M09.2 |
| P3.3 | P3 | Lower | confirmed; documented codec limitation acceptance | hcl/parse.go cty normalization; offline NFD Render refusal | M09.2 |
| P3.4 | P3 | P2 | confirmed | binding/flow.go:333–341,367 mixes bare namespaces and cross-workflow output fallback | M09.4 |
| P3.5 | P3 | P2 | partially confirmed | expressions/portability.go:340 visits main only; parse.go accepts item/index outside loop; ordinary parser remains compatible | M09.4 |
| P3.8 | P3 | Lower | confirmed; inert import acceptance | hcl/parse.go:64 imports info=[] while Render refuses typed block shape | M09.2 |

## Verification and compatibility

Root go test ./... and go vet ./...; affected binding/expressions/strictjson race tests; separate hcl-module tests/vet/races with GOWORK=off GOPROXY=off; adversarial-depth subprocess fixtures; canonical-comment/escape/numeric/NFD/import symmetry cases; whole-schema versus nested draft-07 array proof; flow namespace and trigger-loop fixtures; immutable published versions/schema/conformance digest checks; exact root/nested-codec archive and ordinary downstream-consumer qualification; git diff --check.

Use retained Go 1.26.6 and exact ordinary modules without ambient workspace substitution. Default verification is offline, credential-free and model-free. No live user ledger, host, provider, account, mail, Cloudflare, registration or consumer adoption operation is included.

Public schemas/wires, published grammar/version bytes, accepted historical qualification and independently retained browser/media/legacy/frozen-consumer pins stay preserved. Corrected derived metadata and new worker/package identities require fresh consumer assessment and explicit authority; historical approvals are never upgraded automatically. Source publication requires a new separately named request and independent resolution before downstream adoption. Planning rows may remain pending on this external prerequisite; none is started here.

## Closing review

**Review iterations:** 7/10.
**Review state:** iteration 7 passed, 2026-10-08, after independently observed
normal source publication and passing ordinary downloaded module/consumer proof.
Both read-only reviewers rechecked the whole implementation at
b099f6803277ae94c7e9f1da0904a0140b278f20 and current publication evidence;
no remaining P1/P2 or higher finding. Coordinator acceptance closes after
current-truth consolidation and exact downstream reconciliation.
**Findings/fixes:** iterations 1–5 findings are addressed and whole reviews 6/7
passed. Root full tests/vet, affected races, strict docs and protected-byte/pin
guards pass. Final corrected exact source passes root/codec module-only full
qualification and isolated owner-bootstrap consumer proof; earlier candidate
proofs remain preliminary at their recorded sources. Standalone codec retains
C09 root. Separately authorized initial publication and independent ordinary
resolution/full verification pass. Coordinator acceptance and retirement retain
this full evidence after completed review and downstream reconciliation.
**Execution owner:** sole serial UWS:M09 owner through the M09.5 publication
handoff; all five rows are complete and none is in progress. Coordinator resumes
cross-package reconciliation/acceptance/retirement after the final handoff.
**Commit policy:** Confirmed GOAL COMMIT_POLICY: task; verified task commits and
substantive review/closure commits, without amend/rewrite/tag authority. Fresh
UWS-only main publication authority is recorded below.
**Closure:** persist each started review iteration before reviewing; resume an interrupted pass at the same number. No open P1/P2 may remain at acceptance. Required verification, exact downstream reconciliation and owner-specific consolidation/retirement follow implementation; never reopen completed Stage 11 history.

### Iteration 1 findings (persisted before fixes)

- R1-F01 / P2 / M09.1: delimiter/unary preflight still admitted a shallow
  20,000-term conditional chain (`true ? 0 : ...`), which overflowed a 1 MiB
  child stack in hclsyntax.parseTernaryConditional. Read-only anonymous-overlay
  regression reproduced the failure on Go 1.26.6. Refuse already-unsupported
  conditional syntax before recursive parsing and add isolated child coverage.
- R1-F02 / P2 / M09.3: a declared property intersected by a false
  patternProperties schema was reported reachable/compatible. The read-only
  overlay reproduced whole-schema rejection versus schemaPath compatibility.
  Pattern intersections need conservative proof or indeterminate refusal.
- R1-F03 / P2 / M09.4: owner-kind-only dependency prefixes narrowed the frozen
  generic dependsOn contract, which permits cross-kind targets and defines
  step→workflow→operation precedence plus parallel groups. Preserve that
  precedence without global cross-kind ambiguity; output references remain
  scoped to their workflow.
- R1-F04 / P2 / M09.3: draft-2019 unevaluatedItems was incorrectly treated as
  an ignored future keyword; an expression template was compatible while the
  corresponding literal was rejected. Keep effective unsupported constraints
  indeterminate under their inherited draft.
- R1-F05 / P2 / M09.4: structural results with From=main.join retained an empty
  reference workflow after removing fallback, falsely classifying the named
  join output as unused. Resolve the result's explicit owner context.
- R1-F06 / P2 / M09.4: direct-trigger context was not propagated through a
  step's transitive generic dependsOn, so an ordinary-valid dependency operation
  using item/batchIndex was incorrectly portable outside a loop.
- R1-F07 / P2 / M09.4: the checker created forEach item/index availability before
  its when condition, but executeRunnable evaluates when before forEach binds an
  item. Preserve outer iteration context for pre-iteration controls.
- R1-F08 / P2 / M09.3: draft-07 $ref ignores sibling type, but raw source type
  was used to prove containment. Effective unsupported source references must
  remain indeterminate rather than proving the ignored sibling.
- R1-F09 / P2 / M09.4: dependency contexts also incorrectly inherited an
  operation's not-yet-created forEach iteration, while workflow/step generic
  dependencies were omitted. Traverse dependencies in incoming execution context,
  including group members, before creating body iteration state.
- R1-F10 / P2 / M09.4: source-bound operation request expressions retained an
  empty workflow after fallback removal, falsely marking a same-workflow producer
  output unused. Project operation expressions into their known invocation owners
  rather than selecting a foreign step by name.

**Iteration 1 fix verification:** Unsupported conditionals refuse before parsing;
the existing isolated low-stack child suite also proves codec source JSON is
bounded while its standalone dependency remains C09. Pattern intersections,
draft-2019 constraints and source references remain indeterminate without proof;
local-reference literal proofs retain original compiled children. Generic/group
dependencies, result/operation owners, transitive trigger contexts and incoming
pre-forEach controls have new deterministic regressions. All ten findings are
fixed for full re-review; review counter remains 1 until that pass starts.

### Iteration 2 findings (persisted before fixes)

- R2-F01 / P2 / M09.4: parallel-group members were stored as typed keys,
  bypassing the published member-name step→workflow→operation precedence. An
  ordinary-valid step sharing its operation's ID executes as a group dependency
  but was marked unreachable. Preserve bare member resolution with the same
  generic dependency rules; the explicit step/operation identity remains intact.
- R2-F02 / P2 / M09.4: cross-declaration step dependencies inherit the caller's
  invocation frame, but flow reset every step to its declaration workflow. A
  harmless pure-runtime probe resolved the caller's fetch output while flow
  marked it unused. Record/evaluate step and operation expressions in their known
  invocation owners; never choose a foreign step by bare-name fallback.
- R2-F03 / P2 / M09.3: terminal composed leaves (allOf with false, or not:{})
  were compatible although they admitted no value. Preserve indeterminate for
  unproved terminal applicators rather than dropping their constraints.
- R2-F04 / P2 / M09.4: goto workflow/entry-step targets were omitted from strict
  invocation contexts. A root goto to a loop-declared step falsely retained
  item/index availability. Visit exact terminal goto targets in fresh root context.

**Iteration 2 fix verification:** Bare group-member names retain published
dependency precedence. Cross-declaration steps retain known incoming caller
frames, and step/operation references use recorded same-invocation owners.
Unsupported applicators refuse proof at intermediate and terminal schema paths.
Success/failure goto targets contribute fresh exact root iteration contexts.
New ordinary-valid regression fixtures, full root tests/vet, affected races,
strict docs and protected-byte/module-pin guards passed before committing fixes.

### Iteration 3 findings (persisted before fixes)

- R3-F01 / P2 / M09.4: goto StepID lookup in binding/flow.go and
  expressions/portability.go searches only main/sole-entry declarations, while
  uws1/execution_actions.go:124 selects the globally indexed exact step. The
  returned ordinary/executable-valid pure fixture calls helper.target inside a
  main loop, then goes to target at root; flow reports a missing target and
  strict portability omits its absent batchIndex context. Resolve the exact
  indexed step without changing generic dependency/group precedence.
- R3-F02 / P2 / M09.4: flow stores goto-step edges as dependencies and carries
  the current helper invocation frame, although terminal goto executes the
  target from the original root context. The returned supported pure fixture
  proves the root target cannot resolve helper outputs that flow marks used.
  Keep root transfer context separate from dependency caller frames.
- R3-F03 / P2 / M09.3: binding/validate.go schemaPath returns Compatible for
  an empty fragment before its terminal applicator checks. Composed root
  schemas allOf:[false] and not:{} reject every literal but receive root-path
  compatibility. Apply conservative terminal proof to root and child paths,
  retaining the existing offline compiled resource/dialect context.

The interrupted pass resumes at iteration 3. These returned read-only findings
were revalidated against exact source 8b2f8dfb5da419284ee98e9404b5279bec463bae
before corrections. No acceptance or publication is inferred. Preliminary
docs/m09-root-module-closure.json and docs/m09-codec-module-closure.json remain
at their recorded 8b2f8dfb5da419284ee98e9404b5279bec463bae context.

**Iteration 3 fix verification:** Global goto step resolution and separate root
transfers pass bounded ordinary/executable-valid regressions. The pure reference
runtime rejects batchIndex without an iteration and helper-only outputs after
root transfer. Both success/failure targets are checked. Empty root pointers now
receive the same conservative terminal applicator proof as child pointers.
All three regressions fail against the exact 8b2 baseline through a disposable
read-only overlay and pass on corrected source. Root full tests/vet, affected
binding/expressions/strictjson races, strict MkDocs and protected-byte/module-pin
guards passed. No new large parser inputs or external actions were generated.

### Iteration 4 findings (persisted before fixes)

- R4-F01 / P2 / M09.4: generic workflow dependencies always reset flow to the
  declaration workflow, although ExecuteWorkflow invoked from a root dependency
  retains the unqualified root record frame (uws1/execution.go:145,155–160).
  An ordinary/executable-valid main.fetch → main.join dependsOn helper fixture
  resolves fetch in helper through the pure reference evaluator, but flow marks
  main.fetch's output unreferenced. Preserve the actual incoming root frame;
  explicit/nested workflow calls continue creating their separate frames.
- R4-F02 / P2 / M09.3: schemaPath ignores effective maxItems when traversing
  an array index. An array with maxItems:0 and string items reports #/0 compatible
  despite admitting only an empty array. Compare canonical indexes with the
  compiled node's effective length limit before proving reachability.
- R4-F03 / P2 / M09.4: a core step-output reference blindly marks a same-named
  operation output used. ExecuteStep copies only the operation Result, while
  finalizeOutputs stores the step's separate definitions. An ordinary/executable
  fixture produces different step/operation values and the pure consumer reads
  only the step value; flow suppresses the operation's unreferenced finding.
  Preserve separate output ownership rather than inferring an alias.
- R4-F04 / Lower / M09.5: active milestone summary/index still says all five
  rows are pending and review is 0/10. The authoritative status is executing
  M09.5 at review 4/10. Reconcile routine current-state metadata without changing
  the approved scope, publication boundary or frozen planning provenance.

The read-only correctness reviewer reproduced both findings against exact
adb0a5aa74179f335ca5c5dddab4bbd84089e39d using bounded disposable fixtures at
/tmp/uws-m09-review4-9eyy90_u. The independent contracts/tests reviewer reported
no additional finding. Exact adb0a5a root/codec archive full races/vet/build
passed, but remain preliminary because these P2s block the whole review.

**Iteration 4 fix verification:** Root entry/trigger/goto now use the unqualified
root frame. Root workflow dependencies retain that frame; explicit calls and
workflow dependencies inside a call create their child frame. Paired supported
fixtures prove the root reference succeeds and the nested caller reference
fails, matching native execution. Compiled maxItems excludes unreachable array
indexes. Step references no longer infer operation-output aliases; separate
ordinary/executable-valid native fixtures pass with distinct and equivalent
step/operation IDs. Two older non-frozen flow unit tests now assert the native
separate-output contract rather than the faulty alias. All five reviewer probe
assertions reproduce against exact adb through its preserved original overlay
and pass after fixes. Root full tests/vet, affected races, strict docs and
protected-byte/module-pin guards passed. The active index now records actual
execution. Preliminary adb closure files are retained separately as
docs/m09-review4-{root,codec}-module-closure.json; no source is accepted.

### Iteration 5 findings (persisted before fixes)

- R5-F01 / P2 / M09.3: schemaPath drops containing object constraints.
  Complete schemas declaring x:string with maxProperties:0, const:{}, enum:[{}]
  or propertyNames:false admit {} and reject {x:"fixture"}, but public
  ValidateBinding reports #/x compatible. Check effective supported bounds and
  preserve indeterminate for unproved parent constraints before walking.
  The same issue includes array predecessors: a false index-zero prefix/tuple
  followed by string items reports #/1 compatible although reaching index one
  requires the forbidden prior index. Preserve compiled dialect and all required
  predecessor constraints without discarding them when selecting a later item.
- R5-F02 / P2 / M09.4: flow assigns all workflow expressions to its body frame.
  Workflow controls retain executeOnce's incoming caller record snapshot;
  changing WorkflowScope metadata does not refresh Records. Ordinary/executable
  main.fetch → call helper fixtures whose forEach or when reads fetch outputs
  succeed through the pure reference runtime, but flow marks those outputs
  unreferenced. Record incoming control snapshots separately from refreshed
  body/output contexts. The initial when/body interpretation was corrected by
  exact native execution before code fixes; runtime semantics stay unchanged.
- R5-F03 / P2 / M09.3: templateSchema emits prefixItems:[] for an empty array,
  violating the effective 2020-12 metaschema. A mixed body with a known empty
  array and reviewed string expression cannot be proved compatible, while the
  corresponding literal passes. Represent the exact empty array without an
  invalid prefixItems value and retain inherited target dialect context.
- R5-F04 / Lower / M09.5: the active index's duplicated review number becomes
  stale as soon as the authoritative status starts another pass. Replace that
  duplicate with a status pointer so the persisted counter has one source.

Both read-only reviewers confirm the grouped parent-constraint finding at exact
e89060b456ebcedb6cf3e1427fe60af9d734c04e; bounded public/pure-runtime probes
remain at /tmp/uws-m09-review5-q5gkkijj. The contracts reviewer additionally
reported a convert test linker infrastructure failure (mapping output file:
disk quota exceeded) during a fresh full root run; other root packages and the
separate fresh codec suite passed. Earlier owner-required checks passed. Record
this as a verification gap until a successful fresh required run, not as a
product test failure or an excuse to accept. No capability refusal occurred.

**Iteration 5 fix verification:** Parent cardinality, required property and array
predecessor constraints remain in path proofs; unsupported parent constraints
stay indeterminate. Native literal witnesses prove viable leaves, while failed
samples never prove absence; false and exhausted finite const/enum schemas can
prove absence. Empty arrays use exact const proof, including inherited draft-07
local references. Workflow controls have separately recorded incoming contexts;
their native record snapshots are not conflated with refreshed body outputs.
All six control forms (forEach, duration wait, when, loop items, switch case when
and immediate await wait), plus the negative body-output counterpart, pass
ordinary/executable validation and pure reference checks. Public parent-path and
empty-template probes reproduce against preserved e890 overlay and pass after
fixes. Fresh serial full root tests (including convert), affected races, vet,
strict MkDocs and protected-byte/module-pin guards passed. The quota gap is
closed without a waiver; the coordinator removed only old regenerative build
cache entries, preserving source, module downloads and proofs. The active index
now points to the status counter rather than duplicating it. No new evolution
version or publication is claimed; prior candidate archives stay preliminary.

### Iteration 6 whole review and qualification

Both independent reviewers rechecked the complete milestone implementation,
including all preceding fixes, and found no remaining blocker. Source
b099f6803277ae94c7e9f1da0904a0140b278f20 is the reviewed local release candidate.
Root and nested codec exact module-only archives pass full races/vet/build with
GOWORK=off GOPROXY=off and no replacements. The standalone codec still selects
accepted C09 root; that declaration is not mislabeled as the new root.

The /tmp mount's separate user quota caused initial archive linker/unpack
failures despite the root filesystem's free space. Successful serial rechecks
use owned GOTMPDIR/TMPDIR under /home/peter/.cache/uws-m09-proof-b099f680 and
an isolated bootstrap module cache there. No source, module-download cache,
retained proof or pin was deleted; no quota waiver or skipped check applies.
The archive go.mod/sum bytes remain exact after verification. Whole fresh root
and codec suites, affected races, strict docs and protected guards passed.

Exact source/archive/closure identities, a reusable public consumer fixture and
the main-only publication proposal are in docs/m09-qualification.md and
docs/m09-publication-proposal.md. Local VCS archives and the isolated file-proxy
consumer are owner bootstrap evidence, not independently published Go modules.
M09.5 remains in progress until fresh separately named UWS publication authority,
normal source publication and independent ordinary module/consumer proof.
M09 is not accepted or retired; downstream adoption and external actions remain
outside this owner's present authority. Original Stage 11 history stays frozen.

## Accepted serial predecessor — Udon:M49, 2026-10-08

Udon:M49 is accepted/published after final whole review4/10. Accepted runtime
cd99ccdaa84bd7b0df16e8b9c81bd2e52abd1fd5 independently resolves as
v0.0.0-20261008031303-cd99ccdaa84b; normal publication/closure is independently
observed at 10abc1ffeb5f2b2d10322817b9be9e31c1ff43da.
[Owner handoff](../../../udon/docs/m49-release-handoff.md) records full ordinary
module/CLI proof and frozen dependency/browser boundaries. This satisfies only
the serial scheduling predecessor; UWS imports no new Udon dependency. M09
may start under the confirmed task-commit goal. UWS source publication still
requires its own fresh named authority; the Udon-only grant is consumed.

## Publication gate handoff — 2026-10-08

Local implementation/qualification and whole review6/10 pass at
b099f6803277ae94c7e9f1da0904a0140b278f20. Clean reviewed initial publication
head af225f8b0b5cfbd5848e28e5ce49eaebc0f228bb [skip ci] is the concrete
[proposal](../../docs/m09-publication-proposal.md). The coordinator independently
checked clean source/head, implementation ancestry, documentation-only later
changes, protected versions/corpora/module pins/history, closure bytes and
fresh remote main0a4597122a7baa7e79e46e79e4ec60dbfffc3720.
A fresh named UWS-only grant was requested; no answer/authority is inferred.
M09.5 is blocked on that grant and independent ordinary root/codec upstream
proof. No row is in progress; coordinator owns continuation. Acceptance and
retirement remain incomplete. The consumed Udon/original source grants are not
UWS authority. No push, module adoption, deployment or live action occurred.

## Active status marker correction — 2026-10-08

Coordinator verification found four completed task rows using noncanonical
lowercase [x]. They now use the governing GOAL completed marker [+], in
backticks, while M09.5 remains blocked. This is a current active-record factual
format correction, not a completed-history edit or acceptance/publication claim.
No runtime, source archive, module pin, review counter or frozen byte changes.
The four verified implementation outcomes and review6 remain preserved.

## Fresh UWS source-publication authority — 2026-10-08

The user explicitly grants: "I give UWS:M09 separate source-publication", in
response to the concrete proposal. This grants normal fast-forward publication
of reviewed initial head af225f8b0b5cfbd5848e28e5ce49eaebc0f228bb and reviewed
evidence/acceptance-retirement closure heads to the unchanged exact origin
git@github.com-tabilet:OpenUdon/uws.git refs/heads/main. All push heads use
[skip ci] to suppress the existing force docs deployment to another ref.
The grant includes independent ordinary root/codec resolution and consumer
proof. It grants no tag, force push, gh-pages/docs deployment, host deployment,
live run, installed-consumer migration or unrelated sibling change.

Sole serial ownership resumes in UWS; M09.5 is in progress before any push.
Fresh ls-remote observes main 0a4597122a7baa7e79e46e79e4ec60dbfffc3720, an
ancestor of the exact initial head. Runtime source b099f6803277ae94c7e9f1da0904a0140b278f20
is also its ancestor; later committed changes are evidence/docs only. Live
parent marker/gate bookkeeping and this grant are uncommitted at initial push
and are not claimed to be contained in outgoing af225f8. They will be retained
in the substantive M09.5 evidence commit. Parent owns later cross-package
reconciliation/retirement; this owner writes no sibling ledger or goal/audit
lifecycle. Historical grants and records remain frozen.

## Initial source publication and independent resolution — 2026-10-08

Normal push published exact reviewed af225f8b0b5cfbd5848e28e5ce49eaebc0f228bb
to origin/main from 0a4597122a7baa7e79e46e79e4ec60dbfffc3720. Fresh independent
ls-remote observed that exact head. Its [skip ci] subject preserves main-only
scope; no hosted-CI/docs-deployment result is claimed. The grant and parent
bookkeeping were uncommitted during this initial push, as recorded above.

A fresh separate ordinary module cache resolves root and codec through the
configured proxy.golang.org/direct and sum.golang.org source. Both queries bind
v0.0.0-20261008043726-b099f6803277 to Origin.Hash
b099f6803277ae94c7e9f1da0904a0140b278f20, with codec Subdir hcl. Both Go archive
and GoMod sums match the reviewed candidate. All 360 root and 21 codec archive
files match exact reviewed VCS archive names/content. Ordinary ZIP compression
differs, so raw ordinary byte/SHA identities are recorded separately in
docs/m09-ordinary-module-proof.json rather than relabeling local ZIP artifacts.
Standalone go.mod/sum bytes remain exact. The 52-root/50-codec selected closures
match local qualification; the ordinary combined 52-module consumer closure is
docs/m09-ordinary-consumer-module-closure.json. No replacement, workspace,
bootstrap file proxy or bootstrap module cache supplies this independent proof.

Fresh full tests, full races, vet and builds pass separately from both exact
ordinary downloaded module-only copies. The ordinary combined consumer passes
full races, selecting both exact new modules without replacements. The complete
selected artifacts downloaded successfully from the ordinary configured source;
no network failure remains. Whole review 7 starts after these checks, before
reviewing source/evidence. This proof paragraph alone does not mark the milestone
accepted or retired.

## Iteration 7 outcome and M09.5 handoff — 2026-10-08

The persisted seventh whole pass reviewed complete runtime source and the actual
ordinary publication/module/consumer evidence. Both independent read-only
reviewers found no remaining P1/P2 or higher and no evidence discrepancy.
All required tests, full races, vet/build, strict MkDocs, protected immutable
versions/corpora/module pins and staged/new-file whitespace checks pass.
Only evidence/current-state documentation changed after runtime source b099.

All five task rows are complete with canonical [+] markers; M09.5 is a focused
task commit under the confirmed policy and fresh source-only grant. The
coordinator's marker correction and prior blocked-gate bookkeeping are retained
in that same commit, with honest uncommitted provenance during initial af225f8
publication. The reviewed evidence head may be normally published under the
same grant using [skip ci], then independently observed before ownership returns.
No amendment/history rewrite, force/tag/docs deployment or native goal/audit
lifecycle change is included.

Parent still owns exact downstream specification reconciliation, shared-memory
knowledge consolidation/journaling, final acceptance and retirement, and later
reviewed normal closure publication. Retained source/proof materials and all
original Stage 11/history/browser/consumer pins remain preserved. Implementation
source is b099f6803277ae94c7e9f1da0904a0140b278f20 and ordinary root/codec version
v0.0.0-20261008043726-b099f6803277; no consumer adoption is inferred by publication.

## Coordinator acceptance and retirement — 2026-10-08

All five rows and whole review 7/10 pass at exact runtime
b099f6803277ae94c7e9f1da0904a0140b278f20, root/codec
v0.0.0-20261008043726-b099f6803277. Independent ordinary resolution, every
reviewed file, canonical sums, complete 52/50/52 selected closures, downloaded
full suites/races/vet/builds and ordinary consumer proof pass. Published evidence
51a74545b016b8ab75e30d454338c52f7945a836 is independently observed on main.

Product/architecture/technical truth and reusable lessons are consolidated;
materially superseded lesson paragraphs are preserved literally in knowledge.md.
APItools:M83, OpenUdon:M99 and Kinet:M49 retain their pending rows with the exact
accepted root/codec contract and ordinary proof. Standalone codec still declares
C09 root; combined consumers require both corrected modules. Canonical views
must be regenerated from exact source, and changed package/shape/worker identities
require fresh assessments/confirmations/grants. Ordinary validation/execution,
published versions/corpora and browser/media/Phase A/legacy/frozen-consumer pins
remain unchanged. No new evolution version is triggered.

The complete specification/status are retired after envelope/byte validation;
normal reviewed retirement closure publication under the same UWS-only grant is
independently observed before APItools implementation. Exact final head is
recorded by the coordinator and downstream ledgers after publication, avoiding
a self-referential commit identity. No deployment, installed migration,
registration, provider/API/model/mail action or live run follows. Sole execution
ownership returns to the coordinator; no row remains in progress here.
``````
