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
| M09.4 — Scope flow references and strict portability contexts | `[ ]` | Resolve dependencies/output references within the appropriate workflow/node namespace and retain genuine ambiguity diagnostics. Traverse trigger-route entrypoints as well as main for effective loop context; diagnose item/index availability in the opt-in checker only. Preserve profile field ownership, ordinary grammar/validation and all published version semantics. Sources P3.4/P3.5. |
| M09.5 — Qualify and hand off exact root and nested codec revisions | `[ ]` | Run root/codec conformance, races, immutable-version/corpus guards and independent module-only builds. Add successor regression vectors without rewriting frozen evidence. Persist normal review counts and publish only under fresh named authority; record exact accepted root/codec sources, sums and APItools/OpenUdon/Kinet handoffs. |

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

**Review iterations:** 0/10.
**Review state:** not started; this is review intake, not a pass of an existing gate.
**Findings/fixes:** no implementation or fix verification claimed.
**Execution owner:** sole serial UWS:M09 owner; coordinator makes no writes during this handoff. No row is in progress after M09.1 verification.
**Commit policy:** Confirmed GOAL COMMIT_POLICY: task; verified task commits and substantive review/closure commits, without amend/rewrite/push/tag authority.
**Closure:** persist each started review iteration before reviewing; resume an interrupted pass at the same number. No open P1/P2 may remain at acceptance. Required verification, exact downstream reconciliation and owner-specific consolidation/retirement follow implementation; never reopen completed Stage 11 history.

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
