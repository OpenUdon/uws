# M09 — Stage 11 parser and binding contract remediation

**Stage:** STG-11 post-acceptance remediation. **Owner:** UWS.
**State:** Approved planning, 2026-10-08; implementation not started. All 5 task rows are pending.
**Authority:** The user approved the complete reconciliation proposal with “Implement the plan.” This applies planning files only. A separate execution request is required; no code, commit, source-publication, deployment or live-operation authority follows.
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
| M09.1 — Bound HCL and shape-table parsing before recursion | `[ ]` | Preflight lexical HCL nesting including parentheses, brackets, braces and template/interpolation contexts while ignoring inert string/comment content; enforce depth 100 before hclsyntax.ParseConfig. Replace recursive strict-JSON prevalidation with bounded iterative scanning under existing byte/duplicate/Unicode/numeric rules. Prove fatal-depth inputs refuse in disposable child processes. Sources W1/W2. |
| M09.2 — Verify canonical HCL views and retain inert import symmetry | `[ ]` | Derive canonical bytes through a nonrecursive shared rendering path, then independently parse and compare complete values and numeric lexemes. Reject self-hashed comments/misleading escapes. Reject typed-block shape mismatches in Import, fix the Deprecated paragraph, and document/test non-NFC failure without promising normalization. Keep Import deprecated and nonauthoritative. Sources W6 and P3.2/P3.3/P3.8. |
| M09.3 — Restore sound binding containment and schema context | `[ ]` | Account for patternProperties and true closure evidence; reject false leaf schemas; avoid false missing-field claims for nullable objects/pattern keys, using indeterminate when proof is unavailable. Preserve inherited draft/reference context through nested containment. The review's dependencies example did not reproduce, but draft-07 array probes produced indeterminate/incompatible while whole-schema literal validation was compatible. Sources W3/W4/W5/P3.1. |
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
**Execution owner:** one serial owner across the five ledgers; no row is in progress.
**Commit policy:** The user separately authorized a planning commit on 2026-10-08 with “git commit and then report the index refresh issue in ~/skill-index.md”. This authorizes one commit of the approved planning changes in this owner repository; implementation, publication and deployment remain outside this request. Future task commits follow the separately invoked GOAL/request policy.
**Closure:** persist each started review iteration before reviewing; resume an interrupted pass at the same number. No open P1/P2 may remain at acceptance. Required verification, exact downstream reconciliation and owner-specific consolidation/retirement follow implementation; never reopen completed Stage 11 history.
