# M09 — Stage 11 parser and binding contract remediation

**Stage:** STG-11 post-acceptance remediation. **Owner:** UWS.
**State:** Executing under the confirmed serial GOAL request, 2026-10-08.
**Authority:** The confirmed GOAL request names Udon:M49 → UWS:M09 → APItools:M83 → OpenUdon:M99 → Kinet:M49, task commits and no external mutations. UWS source publication remains a separate exact gate.
**Review source:** stage11-siblings-review.md — Stage 11 code review — sibling packages; uws section.
**Review baseline/range:** `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` → `0a4597122a7baa7e79e46e79e4ec60dbfffc3720`.
**Revalidation HEAD:** `0a4597122a7baa7e79e46e79e4ec60dbfffc3720`; clean worktree, no relevant uncommitted code in the evidence. Approved planning changes are not implementation evidence.
**Lineage:** [C08](../docs/history/status-C08.md), [C09](../docs/history/status-C09.md) and [M08](../docs/history/status-M08.md). Existing acceptance, review counters, statuses and frozen evidence stay preserved.
**Coordinator:** [Stage 11](../../../kinet/docs/stage11.md#post-acceptance-remediation--2026-10-08); package-local specification and status own acceptance.

## Dependencies and handoff

Accepted/published C08/C09/M08 contracts. Serial scheduling follows Udon:M49; this is not a new runtime dependency. Root and nested codec qualification/publication must be independently resolved before APItools:M83, OpenUdon:M99 and Kinet:M49 adoption.

**Exact successor acceptance/publication/build identities:** unset; record full independently observed revisions and hashes during the later execution. Never substitute local HEAD, directory replacements or prior consumed publication authority.

**Downstream:** APItools:M83 native shapes, OpenUdon:M99 package/source verification and Kinet:M49 author/private-exec workers. Retained browser profiles and frozen consumers remain separately pinned.

## Scope and acceptance

Bound untrusted HCL and shape JSON before recursive parsing, enforce deterministic verified-view bytes, make deprecated inert import structurally symmetric, repair supported binding containment/path/draft context and scoped flow/strict-portability diagnostics. Keep existing published grammar/wire/schema and frozen conformance bytes unchanged.

Inputs within byte limits but above depth 100 refuse before unbounded HCL/strict-JSON recursion; comments, strings and template/interpolation delimiters are handled correctly by preflight. Verify accepts only exact deterministic Render bytes plus independent value/numeric-lexeme proof, without recursive Render/Verify calls. Deprecated Import rejects invalid typed blocks; non-NFC refusal stays fail-closed and is documented. Pattern-property containment, false leaves, nullable/pattern paths and inherited draft/reference contexts never produce unsupported compatibility or false missing-field claims. Flow references respect node kind/workflow scope; opt-in portability traverses trigger routes and diagnoses absent iteration context without narrowing ordinary Parse/validation. Root and codec standalone builds, immutable version/corpus guards and owner review pass.

## Tasks

| Item | State | Notes |
|---|---|---|
| M09.1 — Bound HCL and shape-table parsing before recursion | `[x]` | Iterative lexical preflight counts delimiters and refuses unsupported active templates before HCL parsing; strict JSON uses an iterative depth-100 stack. Exact duplicate/Unicode/numeric/trailing checks remain. Fatal-depth Import, strict-JSON and ParseTable child processes passed with a 1 MiB stack. Focused root tests/races/vet and separate full codec tests/races/vet passed, GOWORK=off GOPROXY=off Go 1.26.6. Sources W1/W2. |
| M09.2 — Verify canonical HCL views and retain inert import symmetry | `[x]` | Shared nonrecursive writer path enforces exact deterministic bytes before independent complete-value/numeric-lexeme proof. Self-hashed comment/escape/formatting tests, typed-block empty/null/mismatch symmetry and JSON/YAML NFD string/key fail-closed tests pass. Deprecated is a separate Go doc paragraph; README states the limitation. Separate full codec races and vet passed. Sources W6 and P3.2/P3.3/P3.8. |
| M09.3 — Restore sound binding containment and schema context | `[x]` | Original compiled child nodes preserve dialect/local-reference context. Pattern-driven closure and nullable/pattern paths are indeterminate where unsupported; false leaves refuse. Draft-07 tuple/future-keyword/local-reference tests agree with whole-schema literal validation; binding full tests/races/vet passed. The original dependencies claim remains unconfirmed. Sources W3/W4/W5/P3.1. |
| M09.4 — Scope flow references and strict portability contexts | `[x]` | Explicit identities retain generic group barriers and step→workflow→operation dependency precedence; genuine duplicate step definitions remain ambiguous. Output references use explicit workflow/result owners and known operation invocation contexts, with no foreign-name fallback. Strict checker separates loop/item availability and pre-forEach controls, including transitive trigger/dependency contexts. Root full tests/races/vet pass; ordinary Parse remains compatible. Sources P3.4/P3.5. |
| M09.5 — Qualify and hand off exact root and nested codec revisions | `[~]` | Running root/codec conformance, races, immutable-version/corpus guards and independent module-only builds. Successor regressions are additive. Publication and ordinary new-root downstream proof require fresh named authority; prepare reviewed exact artifacts before that gate. |

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

**Review iterations:** 1/10.
**Review state:** iteration 1 completed with findings, 2026-10-08; fixes verified before the next whole pass. No pass/acceptance claimed.
**Findings/fixes:** ten persisted P2 findings below are addressed. Root full tests/races/vet, separate full codec races/vet, strict docs build and protected version/corpus/module-pin guards pass. Standalone corrected-source artifacts are the next qualification step.
**Execution owner:** sole serial UWS:M09 owner; coordinator makes no writes during this handoff. Only M09.5 is in progress.
**Commit policy:** Confirmed GOAL COMMIT_POLICY: task; verified task commits and substantive review/closure commits, without amend/rewrite/push/tag authority.
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
