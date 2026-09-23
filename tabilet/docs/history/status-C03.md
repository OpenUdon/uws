# C03 History

Milestone: `C03`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-C03.md`
Source specification: `tabilet/memory-bank/milestone.md#c03---portable-execution-contract`
Evidence: `c43b4dd504258cb1f2d46afaddfd06b251cdb970`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `2`
Verification: `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`; immutable historical artifact checks
Consolidated into: [UWS 1.10.0 specification](../../../versions/1.10.0.md), [validation guide](../../../docs/09-Validation.md), [product truth](../../memory-bank/product.md), [release notes](../../../versions/CHANGELOG.md), [conformance vectors](../../../testdata/conformance/1.10.0.json); deferred candidate directions remain in [milestone planning](../../memory-bank/milestone.md)

## Final milestone specification

`````markdown
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
`````

## Final status

`````markdown
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
per-call scoping of nested steps, operations, dependencies, and merge records,
clear rejection of recursive workflow calls, and numeric ordering of nested
iteration results. C02 establishes exact declared-version schema selection,
no fallback for prerelease or unpublished versions, SemVer-aware feature gates,
and compatibility coverage for all published UWS 1.x schemas. M03's accepted
Go mapping escapes ordinary dynamic keys with `__uws_literal__` when they
collide with legacy dollar-key spellings, while preserving legacy dollar
decoding. C03 must preserve that behavior and decide only whether a separate
portable HCL mapping document is justified. C01 settles invocation and
ordering; C02 settles release gating; B01 settles the browser 1.8 profile
cross-reference. C03 selected UWS 1.10.0 under C02 after implementing and
version-gating the new semantics.

**Downstream impacts.** Reconsider the unnumbered interoperability, profile,
and UWS 2.0 candidates after acceptance, without promoting them automatically.

**Review provenance.** “UWS spec and architecture review: 1.9.2 and earlier
(2026-09-23)”; stated baseline
`8382d0f26b3b10870125760643078d1a1a3e31b6`; revalidated at full HEAD
`e765fed4ba0a241e9481f99fc324c6a8afb9a8a1` with untracked memory-bank
initialization content present. No relevant uncommitted core implementation
or versioned-document change formed the evidence. M01/M02 are completed
historical documentation and immutability work, not reopened work. The source
finding IDs below point to this one review. “Partial” notes mark a narrower
confirmed gap than the review claimed; no task assumes its unverified details.
C5 duplicates A3 and does not create another task.

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
| D2 | P3 | P3 | Confirmed duplicate await bullet in `versions/1.9.2.md` §4.5.6.3. |
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

**Review gate.** Ordinary intake after all task rows and acceptance checks
were complete. Review the full milestone diff, including schema/spec/code/
version compatibility, security boundaries, conformance evidence, generated
artifacts, and downstream docs. Iteration 1 findings were fixed and the
affected tests, strict docs build, hash checks, and diff check passed; iteration
2 of 10 is now started for another whole-milestone review.

**Iteration 1 findings (resolved).** P2: UWS 1.10 truthiness treated every
non-empty reflected map/slice as a valid JSON value without checking nested
members or map-key types, and `big.Rat.SetString` could allocate excessively
for very large `json.Number` exponents. The replacement lexical zero check
must inspect only the significand, not the exponent, so valid values such as
`0e10` stay false. The cross-type audit also found Go's `uintptr` numeric kind
missing from 1.10 truthiness and wait-duration conversion despite the accepted
all-numeric-kinds contract. Whole-diff review also found a present-tense
browser-profile sentence in the new 1.10 spec still calling Browser 1.7 the
latest contract even though Browser 1.8 is current. Fix and re-review these
paths. Composite values are now recursively checked, JSON numbers are checked
without exponent expansion and zero detection uses only the significand,
`uintptr` is handled consistently, and the stale profile reference is fixed.
Focused Go tests, affected package tests, strict MkDocs, artifact hash checks,
and `git diff --check` passed.

**Iteration 2 findings: none.** The complete milestone diff was re-reviewed
after iteration 1 fixes across compatibility, schema/spec/code synchronization,
security boundaries, conformance fixtures, generated artifacts, docs, and
downstream candidate directions. No P1, P2, or higher-severity issue remains;
the bounded review gate passes in iteration 2 of 10.

| Item | State | Notes |
|---|---|---|
| Specify and verify non-`await` wait | `[+]` | For UWS 1.10+, non-`await` `wait` evaluates once after dependencies and `when`, then delays before the operation/workflow/step body. It must resolve to a finite JSON number from 0 through 86,400 seconds inclusive. The orchestrator performs a cancellable delay with fractional seconds rounded to the nearest nanosecond; invalid values/evaluation failures fail before the body. Workflow/step delay is inside their serialized timeout; operation delay is before the leaf attempt and does not consume the attempt timeout. UWS 1.0–1.9.2 behavior is preserved; `await.wait` remains a truthy polling predicate at every version. Focused tests cover timing, scopes, invalid values, cancellation, and old-version no-op behavior. Owner: A3/C5. |
| Align criteria, truthiness, expressions, and runtime values | `[+]` | UWS 1.10 truthiness now handles `json.Number` and all Go numeric kinds consistently, rejects non-finite numeric values, and retains the existing empty string/array/object rules; older versions retain old behavior. UWS 1.10 JSON Pointer fragments percent-decode before token splitting, validate `~0`/`~1`, preserve explicit-null presence separately from missing, and retain C02's canonical array-index gate; older versions retain raw fragment behavior. XPath 1.0 number truthiness treats NaN as false only from 1.10. Regex contexts in 1.10 require text/UTF-8 bytes; prior JSON coercion remains versioned for older documents. UWS 1.10 rejects dotted names that the expression grammar cannot address in outputs, top-level/component variables, workflow input schemas, and step inputs; legacy dotted names remain accepted. Workflow targets are literal declared IDs, not expressions. Decisions for task 5 publication: regex uses RE2 search semantics; JSONPath uses RFC 9535 and succeeds when a selected value is truthy; XPath uses XPath 1.0 boolean conversion; top-level `variables` shadow same-name `components.variables`; `$steps` is scoped to one workflow invocation; `$trigger` exposes payload only; non-HTTP `$response` shapes remain profile-defined. `go test ./uws1 ./validation -count=1` passes. Owners: A7, B3, B4, C1, C6–C9. |
| Specify entry, invocation, and structural results | `[+]` | UWS 1.10 entry selection: a sole workflow is the entry; with multiple workflows, exactly one `main` is required. Direct operation/workflow/step APIs target the named object but executable validation still requires an unambiguous `main` when multiple workflows are present. Executable identifiers are globally unique across operation, workflow, step, and parallel-group namespaces, except parallel-group membership may recur and a step may equal its directly referenced operation ID. Workflow-call targets are literal IDs; recursive calls fail. Switch evaluates cases in declaration order; the first truthy case runs, and an unguarded case is an unconditional match at its position. Await evaluates immediately, polls after an executor-configured interval (200 ms default), executes its nested steps once on truthy, and obeys its serialized or executor-default timeout plus context cancellation. Parallel branches run concurrently; first branch error cancels siblings and is returned; goto/end signals are rejected inside a parallel branch. Loop iterations and batches are sequential and preserve source order; the loop result is an array of `{index,batchIndex,item}` records, with zero-based indexes. `forEach` is sequential and preserves order; its parent result contains per-item `{index,item,status,error,result,outputs}` records and each declared output is aggregated into an ordered array. Merge reads only declared dependencies, expands parallel groups by member declaration order, returns dependency records in declared order and nested iterations in numeric order; each record has `{id,kind,status,error,result,outputs}`. Sequence, parallel, switch, and await do not synthesize result values; workflow `results[]` declarations are metadata and do not add executor behavior. Added regression tests for switch declaration order, loop result shape/batch indexes, and ordered `forEach` output aggregation. Existing tests cover entry selection, identifier ambiguity, await polling/cancellation, parallel control-signal rejection, merge expansion/order, and C01 invocation scoping. `go test ./uws1 -count=1` and `git diff --check` pass. Owners: B1/B2/B5–B10/C12. |
| Specify failure, actions, triggers, and security | `[+]` | Action semantics for publication: runtime leaf error or unsatisfied `successCriteria` enters ordered `onFailure`; all criteria on an action must match, empty criteria match, first matching action wins, and absent match preserves the failure. `retryLimit` counts retries after the initial attempt; each retry gets a fresh operation timeout and optional cancellable `retryAfter`. Criteria-evaluation errors fail directly. After success, ordered `onSuccess` uses the same first-match/all-criteria rules; no matching action means normal continuation. `goto` and `end` are control signals, and parallel rejects them. Trigger dispatch validates a zero-based output index, routes targets in route/target declaration order, deduplicates repeated targets, resolves labels before numeric indexes, and errors on missing trigger, invalid index, or no matched targets. Route targets are top-level steps of the selected entry workflow or declared workflows. Added an output-range guard because direct `Orchestrator.ExecuteTrigger` could otherwise accept a fabricated numeric route/index combination despite document-level validation. Security publication will distinguish normative core guarantees from runtime advice: authenticate trigger ingress before dispatch; constrain source URL loading/egress; keep secrets in external stores and redact logs/outputs; bound expression cost, item counts, concurrency, payload and record growth, and timeouts; assess retry side effects/idempotency; treat content-trust as advisory provenance rather than sanitization or authorization. Core does not implement those runtime mitigations. No new trigger kinds or stable error taxonomy. Focused and full `uws1` tests and diff check pass. Owners: B11/B12/C10/C11/C19. |
| Publish portable core contract and conformance evidence | `[+]` | Published `versions/1.10.0.{json,md}` under C02's SemVer policy, retaining every earlier core artifact. Added strict schema patterns and Go enforcement for expression-addressable names, language-neutral vectors at `testdata/conformance/1.10.0.json`, the core conformance classes, release history, docs and cross-references; regenerated the embedded archive and added exact SHA-256 immutability entries. Preserved M03 HCL helper behavior without adding a normative HCL mapping. D1–D5/D11/D13, C14/C15/C17/C18, and E3/E7 are reconciled. Full tests, race, vet, strict MkDocs, diff check, immutable-history checks, and the bounded whole-milestone review pass in iteration 2 of 10. |
`````
