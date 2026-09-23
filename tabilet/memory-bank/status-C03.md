# Status C03 - Portable Execution Contract

**Acceptance.** A newly versioned core specification, schema, Go behavior,
release history, conformance fixtures, documentation, and generated embedded
archive agree on the accepted contract. Older versions remain protected and
compatible under C02's policy. Full tests, race tests, vet, strict MkDocs
build, diff check, and the bounded whole-milestone review pass. No release
number is preselected; choose it under C02's version policy.

**Dependencies.** [M03](../docs/history/status-M03.md),
[C01](../docs/history/status-C01.md), [C02](../docs/history/status-C02.md),
and [B01](../docs/history/status-B01.md) are completed and accepted. Browser
1.8 is a separate UWS 1.9 profile; C03 should cross-reference it without
importing its template or formatting semantics into the core contract. C01
establishes deterministic caller-step identity for workflow invocations,
per-call scoping of nested steps, operations,
dependencies, and merge records, clear rejection of recursive workflow calls,
and numeric ordering of nested iteration results. C02 establishes exact
declared-version schema selection, no fallback for prerelease or unpublished
versions, SemVer-aware feature gates, and compatibility coverage for all
published UWS 1.x schemas. M03's
accepted Go mapping escapes ordinary dynamic keys with `__uws_literal__` when
they collide with legacy dollar-key spellings, while preserving legacy dollar
decoding. C03 must preserve that behavior and decide only whether a separate
portable HCL mapping document is justified. C01 settles invocation and
ordering; C02 settles release gating; B01 settles the browser 1.8 profile
cross-reference, while UWS core remains at 1.9.2 until C03 chooses its release.

**Downstream impacts.** Reconsider the unnumbered interoperability, profile,
and UWS 2.0 candidates after acceptance, without promoting them automatically.

**Review provenance.** “UWS spec and architecture review: 1.9.2 and earlier
(2026-09-23)”; stated baseline
`8382d0f26b3b10870125760643078d1a1a3e31b6`; revalidated at full HEAD
`e765fed4ba0a241e9481f99fc324c6a8afb9a8a1` with untracked
memory-bank initialization content present. No relevant uncommitted core
implementation or versioned-document change formed the evidence. M01/M02 are
completed historical documentation and immutability work, not reopened work.
The source finding IDs below point to this one review. “Partial” notes mark a
narrower confirmed gap than the review claimed; no task assumes its unverified
details. C5 duplicates A3 and does not create another task.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| A3 / C5 | P1 / P2 | P1 | Confirmed/duplicate: `uws1/execution_structural.go` handles `wait` for `await`; `versions/1.9.2.md` §7.2 requires broader enforcement but lacks a duration. |
| A7 | P2 | P2 | Confirmed: `uws1/execution_criteria.go` truthiness omits `json.Number` and several numeric kinds. |
| B1 | P1 | P2 | Confirmed: entry selection in `uws1/executable_validation.go` and `uws1/document_index.go` lacks a matching portable core rule. |
| B2 | P1 | P2 | Confirmed: `uws1/executable_validation.go` applies cross-kind name uniqueness not stated in `versions/1.9.2.md` §7.1. |
| B3 | P1 | P1 | Confirmed: criterion implementation in `uws1/execution_criteria.go` differs from `versions/1.9.2.md` §5.2 context/condition prose; regex and XPath details need tests. |
| B4 | P1 | P1 | Confirmed: `uws1/validation_workflow.go` requires literal declared workflow ID while `versions/1.9.2.md` describes `workflow` as expression-capable. |
| B5 | P1 | P2 | Confirmed: `uws1/execution_structural.go` switch declaration-order and unguarded-case rules are not fully stated in `versions/1.9.2.md`. |
| B6 | P1 | P2 | Confirmed: `uws1/execution_structural.go` await nested-step and polling behavior is underspecified in `versions/1.9.2.md`. |
| B7 | P1 | P1 | Confirmed: `uws1/execution_structural.go` parallel cancellation/control-flow restrictions lack corresponding core prose. |
| B8 | P1 | P2 | Confirmed: `uws1/execution_structural.go` loop serial/batch execution and result shape are not fully portable in `versions/1.9.2.md`. |
| B9 | P1 | P2 | Confirmed: `uws1/execution_structural.go` merge result shape/expansion are not fully described; ordering defect belongs to C01/A6. |
| B10 | P1 | P2 | Confirmed: `uws1/execution_structural.go` forEach aggregation and loop `$item`/`$index` scope exceed `versions/1.9.2.md` §5.1 prose. |
| B11 | P1 | P1 | Confirmed: `uws1/execution.go` and action helpers define retry, goto, end, cached-node, and criteria details absent from `versions/1.9.2.md`. |
| B12 | P1 | P2 | Confirmed: `uws1/trigger_dispatch.go` route order, duplicate handling, no-route error, and label precedence are underspecified. |
| C1 | P2 | P2 | Confirmed: `versions/1.9.2.md` uses truthiness without a complete rule; `uws1/execution_criteria.go` supplies only implementation behavior. |
| C6 | P2 | P2 | Confirmed: variable/output name schema and expression grammar in `versions/1.9.2.md` allow declarations that cannot be referenced; root/component precedence is unstated. |
| C7 | P2 | P2 | Confirmed: `versions/1.9.2.md` gives HTTP-centric `$response`/`$trigger` meaning without portable extension/trigger shape and header rules. |
| C8 | P2 | P2 | Confirmed: `versions/1.9.2.md` and `uws1/execution.go` leave nested `$steps` visibility and aggregated outputs unclear. |
| C9 | P2 | P2 | Partially confirmed: pointer grammar in `versions/1.9.2.md` needs clarification; percent-decoding order proposed by review must be tested against `uws1` behavior before changing it. |
| C10 | P2 | P2 | Partially confirmed: default success is implemented in `uws1/execution.go` but underdocumented; portable error taxonomy/idempotency remain separate candidate work. |
| C11 | P2 | P2 | Partially confirmed: optional `routes` in `versions/1.9.2.json` conflicts with `uws1/trigger_dispatch.go` no-route failure; new trigger kinds remain a 2.0 candidate. |
| C12 | P2 | P2 | Partially confirmed: result value can intentionally be implementation-defined, but loop/merge shapes in `uws1/execution_structural.go` need portable contract; overlaps B8/B9. |
| C14 | P2 | P3 | Confirmed editorial gap: `versions/1.9.2.md` §4.6 is incomplete although other normative sections cover some listed rules. |
| C15 | P2 | P2 | Confirmed: `versions/1.9.2.md` normative Go API names lack language-neutral conformance classes. |
| C17 | P2 | P3 | Confirmed cross-reference gap: `versions/1.9.2.md` does not surface separate registration documents; do not embed those profiles into core. |
| C18 | P2 | P3 | Confirmed current `x-uws-*` cross-reference gap in `versions/1.9.2.md`; proposed `uws.*` namespace reservation is a separate governance candidate. |
| C19 | P2 | P2 | Confirmed: `versions/1.9.2.md` §8 has some security guidance but lacks specific source-URL, evaluation-cost, resource-limit, trigger-auth, and secret-handling coverage. |
| D1 | P3 | P3 | Confirmed stale release numeral in `versions/1.9.2.md` §4.5.3; fix in the new release, preserving immutable history. |
| D2 | P3 | P3 | Confirmed duplicate await bullet in `versions/1.9.2.md` §4.5.6.3; fix in new release. |
| D3 | P3 | P3 | Confirmed incomplete Step definition in `versions/1.9.2.md` §3.5. |
| D4 | P3 | P3 | Confirmed wrong entry-workflow cross-reference in `versions/1.9.2.md` §4.6. |
| D5 | P3 | P3 | Partially confirmed: `versions/1.9.2.md` references omit core authorities and retain a Medium article; browser-only authorities need not be imported. |
| D11 | P3 | P3 | Partially confirmed: registration entries were supplied by completed M01; `versions/CHANGELOG.md` still needs release wording for output-name order and version-gating decision. |
| D13 | P3 | P3 | Confirmed: `docs/02*` and `docs/03*` contain expression/await examples unsupported or ambiguous under `versions/1.9.2.md` grammar. |
| E3 | Recommendation | P2 | Confirmed language-neutral conformance gap in `versions/1.9.2.md` and current Go-centric tests; broader fixture distribution is a candidate. |
| E7 | Recommendation | P2 | Confirmed normative execution-semantics gap across `versions/1.9.2.md` and `uws1/execution*.go`; overlaps B findings, not a separate implementation. |

**Decisions and boundaries.** Non-`await` `wait` is orchestrator-owned,
measured in bounded numeric seconds; `await` retains a predicate. Test old
documents before assigning new behavior. MCP, content-trust wire reports,
stable error codes, normative HCL mapping, expression-marker redesign, new
trigger kinds, and namespace reservation are not part of this milestone.
Current profile semantics remain separately versioned. Security findings are
documented without claiming runtime mitigations that do not exist.

**Review gate.** Ordinary intake; 0 of 10 iterations. Start only after all
task rows and acceptance checks are complete.

| Item | State | Notes |
|---|---|---|
| Specify and verify non-`await` wait | `[+]` | For UWS 1.10+, non-`await` `wait` evaluates once after dependencies and `when`, then delays before the operation/workflow/step body. It must resolve to a finite JSON number from 0 through 86,400 seconds inclusive. The orchestrator performs a cancellable delay with fractional seconds rounded to the nearest nanosecond; invalid values/evaluation failures fail before the body. Workflow/step delay is inside their serialized timeout; operation delay is before the leaf attempt and does not consume the attempt timeout. UWS 1.0–1.9.2 behavior is preserved; `await.wait` remains a truthy polling predicate at every version. Focused tests cover timing, scopes, invalid values, cancellation, and old-version no-op behavior. Owner: A3/C5. |
| Align criteria, truthiness, expressions, and runtime values | `[+]` | UWS 1.10 truthiness now handles `json.Number` and all Go numeric kinds consistently, rejects non-finite numeric values, and retains the existing empty string/array/object rules; older versions retain old behavior. UWS 1.10 JSON Pointer fragments percent-decode before token splitting, validate `~0`/`~1`, preserve explicit-null presence separately from missing, and retain C02's canonical array-index gate; older versions retain raw fragment behavior. XPath 1.0 number truthiness treats NaN as false only from 1.10. Regex contexts in 1.10 require text/UTF-8 bytes; prior JSON coercion remains versioned for older documents. UWS 1.10 rejects dotted names that the expression grammar cannot address in outputs, top-level/component variables, workflow input schemas, and step inputs; legacy dotted names remain accepted. Workflow targets are literal declared IDs, not expressions. Decisions for task 5 publication: regex uses RE2 search semantics; JSONPath uses RFC 9535 and succeeds when a selected value is truthy; XPath uses XPath 1.0 boolean conversion; top-level `variables` shadow same-name `components.variables`; `$steps` is scoped to one workflow invocation; `$trigger` exposes payload only; non-HTTP `$response` shapes remain profile-defined. `go test ./uws1 ./validation -count=1` passes. Owners: A7, B3, B4, C1, C6–C9. |
| Specify entry, invocation, and structural results | `[+]` | UWS 1.10 entry selection: a sole workflow is the entry; with multiple workflows, exactly one `main` is required. Explicit operation/workflow/step APIs do not require a document entry. Executable identifiers are globally unique across operation, workflow, step, and parallel-group namespaces, except parallel-group membership may recur and a step may equal its directly referenced operation ID. Workflow-call targets are literal IDs; recursive calls fail. Switch evaluates cases in declaration order; the first truthy case runs, and an unguarded case is an unconditional match at its position. Await evaluates immediately, polls after an executor-configured interval (200 ms default), executes its nested steps once on truthy, and obeys its serialized or executor-default timeout plus context cancellation. Parallel branches run concurrently; first branch error cancels siblings and is returned; goto/end signals are rejected inside a parallel branch. Loop iterations and batches are sequential and preserve source order; the loop result is an array of `{index,batchIndex,item}` records, with zero-based indexes. `forEach` is sequential and preserves order; its parent result contains per-item `{index,item,status,error,result,outputs}` records and each declared output is aggregated into an ordered array. Merge reads only declared dependencies, expands parallel groups by member declaration order, returns dependency records in declared order and nested iterations in numeric order; each record has `{id,kind,status,error,result,outputs}`. Sequence, parallel, switch, and await do not synthesize result values; workflow `results[]` declarations are metadata and do not add executor behavior. Added regression tests for switch declaration order, loop result shape/batch indexes, and ordered `forEach` output aggregation. Existing tests cover entry selection, identifier ambiguity, await polling/cancellation, parallel control-signal rejection, merge expansion/order, and C01 invocation scoping. `go test ./uws1 -count=1` and `git diff --check` pass. Owners: B1/B2/B5–B10/C12. |
| Specify failure, actions, triggers, and security | `[ ]` | Resolve B11/B12/C10/C11/C19: default success, retry/goto/end, route dispatch/no-route behavior, failure boundaries, and concrete security guidance. Leave new trigger kinds and portable error taxonomy for candidates. |
| Publish portable core contract and conformance evidence | `[ ]` | Resolve C14/C15/C17/C18, D1–D5/D11/D13, E3/E7; choose release under C02, preserve M03's accepted HCL helper mapping, update schema/spec/changelog/docs/cross-references, add language-neutral conformance cases, regenerate embedded archive, and preserve immutable versions. |
