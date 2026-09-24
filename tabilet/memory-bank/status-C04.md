# Status C04 - UWS 1.11 Contract Repair

**Goal.** Publish a coherent UWS 1.11 contract that closes the confirmed
control-flow, grammar, conformance, and guide gaps in 1.10 and incorporates
[C05](../docs/history/status-C05.md)'s exact-version admission correction.
The UWS 1.11 schema/specification and executable conformance corpus have been
published during C04 execution; the final guide/release-surface row and
milestone acceptance review remain active.

**Scope and contract choices.** In UWS 1.11, `goto` from a sub-workflow, loop,
or `forEach` item unwinds the entire top-level run. Remaining siblings and
iterations are skipped. The globally declared step or workflow target, with
its dependencies, is then invoked in root orchestration scope; it does not
select a prior caller-scoped workflow invocation. After unwinding, a root-key
target record with terminal status `success`, `error`, or `skipped` counts as
already completed and fails explicitly; an in-flight root target also fails
rather than waiting or replaying. A same-named record in another workflow-call
scope does not count. After the target completes, the top-level run ends; it
does not resume after the target. The 1.11 grammar admits
`$response.body.<path>`, `$batchIndex` in loop contexts, and complete
JSON-number literals for numeric `wait` and `batchSize`, with version-aware
advisory content-trust parsing. A merge over `forEach` includes iteration
records when present, otherwise its parent record for a skipped or zero-item
dependency; it never counts both. A failed dependency still aborts merge.
`await` re-evaluates its predicate but does not re-execute completed status
operations; portable operation reexecution remains a UWS 2.0 candidate.
Corrected genuine omissions and contradictions in UWS 1.11 while preserving
frozen 1.10.0 artifacts.

**Dependencies.** [C02](../docs/history/status-C02.md) established
declared-version gates and exact file-schema selection; [C03](../docs/history/status-C03.md)
published UWS 1.10.0; completed [C05](../docs/history/status-C05.md) enforces
exact-version admission; completed [B02](../docs/history/status-B02.md) publishes
the accepted, opt-in Browser 1.9 contract for UWS 1.11 release surfaces.
**Downstream impact.** No active product milestone depends on C04; M04 runs
later but has no product dependency.

**Compatibility and rollout.** New grammar, goto failure, merge shape, and
numeric-literal semantics are opt-in by declared 1.11 version. Previously
published schemas/specifications stay byte-for-byte unchanged and older-version
behavior is covered by compatibility tests. C05's exact-version admission
applies across declared versions, including older documents; the 1.11
specification and changelog state this deliberate bug-fix exception alongside
earlier cross-version corrections. Rollback is selecting a published earlier
version only where the document does not rely on 1.11 features.

**Acceptance and verification.** New schema/spec, Go validator and executor,
content-trust analyzer, executable conformance corpus, fixtures, docs,
changelog, immutable membership/digests, and embedded archive agree.
Conformance vectors use typed inputs/operations and every vector is executed
by Go, including forward/backward goto, merge over `forEach`, version
rejection, and numeric waits. Grammar lint checks docs and fixtures against
each example's declared version. Existing 1.0/1.1 conversion fixtures retain
their versions and use exact-path allowlist entries for intentional legacy or
§5.5 implementation-dependent expressions; they are not represented as
conforming 1.11 examples. Separate 1.11 fixtures prove the new grammar. Test
prior versions and regression-check C05's direct/file entry paths. The final
acceptance commands are `go test ./...`, `go test -race ./...`, `go vet ./...`,
`mkdocs build --strict`, and `git diff --check`, followed by the bounded
whole-milestone review gate in [milestone.md](milestone.md).

**Review provenance.** Source: “UWS review, pass 2: 1.10.0, browser 1.8 and
the remediation (2026-09-23)” (`uws-review-2.md`), stated baseline `5923cc0`;
revalidated at full HEAD `5923cc0d3b2c852e45aee3771201c642b9daeca9`.
No relevant uncommitted core code or versioned-document change formed the
evidence. The worktree contained only pre-existing untracked Tabilet GOAL,
archives, and evolution files. Lineage: completed C02/C03, with first-review
residuals B3, B10, C7, C9, C13–C15, C18, D3/D5/D13 where applicable. Neither
retired status is reopened. This is ordinary intake, not review iteration 1.

**Follow-up plan review.** The pasted “Plan review of second-review
remediation (2026-09-23)” supplied no source priority or baseline; revalidated
at full HEAD `5923cc0d3b2c852e45aee3771201c642b9daeca9` with the
previous reconciliation's uncommitted planning edits and pre-existing
untracked Tabilet evidence present. No new core implementation or published
artifact formed the evidence. F01/F03–F07 below refine existing pending C04
work; N3/N14 moved without implementation to C05. This is ordinary intake,
not a review-gate iteration.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| N2 | P2 | P2 | Confirmed: `uws1/execution.go` ends after a goto target; `versions/1.10.0.md` and `docs/06-Success-Criteria-and-Actions.md` do not clearly state terminal transfer. |
| N4 | P2 | P2 | Confirmed: `versions/1.10.0.md` §5.6 excludes `$response.body.<path>` used in `docs/03-Runtime-Expression-Grammar.md`, other guides, fixtures, and opaque-handled by `contenttrust/expression.go`. |
| N5 | P2 | P3 | Confirmed ergonomic gap: §5.6 lacks complete numeric literals for 1.10 numeric `wait`; numeric `batchSize` shares the gap. Literal `items` remains a 2.0 candidate. |
| N6 | P2 | P2 | Confirmed: `uws1/schema_artifacts_test.go` checks only shape of `testdata/conformance/1.10.0.json`, not executable behavior or digest. |
| N8 | P2 | P2 | Confirmed: `docs/02-Six-Structural-Constructs.md`, `docs/04-Triggers-and-Route-Dispatch.md`, `docs/05-Structural-Results.md`, `docs/06-Success-Criteria-and-Actions.md`, and `docs/07-Execution-Model.md` retain stale or contradictory 1.10 execution/result/action examples. |
| N9 | P2 | P2 | Partially confirmed: `versions/1.10.0.md` has genuine definition, reference, and criteria inconsistencies; some “UWS 1.9” references describe history correctly and must not be rewritten blindly. |
| N10 | P2 | P2 documentation / P3 new feature | Confirmed: `uws1/execution_structural.go` re-evaluates `await` but memoizes completed status steps. Portable operation reexecution remains a UWS 2.0 candidate; `docs/02-Six-Structural-Constructs.md` now states that `await` cannot portably poll a status operation by placing it in a preceding sequence step. |
| N11 | P3 | P2 | Confirmed: `uws1/execution_structural.go` merge expansion includes both a `forEach` parent aggregate and its iteration records, double-counting the data. |
| N12 | P3 | P3 | Confirmed: `uws1/execution_criteria.go` strips matching context prefixes and accepts bare `/` pointers, but `versions/1.10.0.md` does not describe both. |
| N13 | P3 | P3 | Confirmed: `uws1` workflow-call isolation, merge ordering, and trigger output-index guard apply as bug fixes to older documents; `versions/CHANGELOG.md` needs explicit compatibility exceptions. |
| N15 | P3 | P3 | Confirmed: `versions/1.10.0.json` content-trust output key pattern is wider than the operation-output name pattern. Correct in the new schema. |
| N18 | P3 | P3 | Confirmed: `versions/1.10.0.md` cites moving informative references and mutable `testdata/conformance/1.10.0.json`; pin new-release references and corpus by digest. |
| N20 | P3 | P3 | Confirmed: no shared grammar lint tests docs/fixtures against §5.6; `contenttrust/expression.go` is the existing parser but a new public Go API is not required. |
| F01 | not supplied | P2 | Confirmed: `testdata/sample.uws.json` declares 1.0.0 and `testdata/big/big.json` declares 1.1.0; both include expressions outside their declared grammar, including `$error.*`, `$context.*`, `$signals.*` and response dot-walks. Lineage: N4/N20. |
| F03 | not supplied | P2 | Confirmed: `versions/1.10.0.md` promises loop `$batchIndex` to the runtime but §5.6 omits its grammar production. Lineage: N9/C03. |
| F04 | not supplied | P3 | Confirmed: `versions/1.10.0.md` §1.1 presents `1.0.0-beta.1` as an example although C05's exact admission must reject it. Correct only in new 1.11 prose. Lineage: N3/N9. |
| F05 | not supplied | P2 | Confirmed: C04's pending publication row lacked a clause-by-clause checklist for `versions/1.10.0.md` N9 corrections. Lineage: N9/C03. |
| F06 | not supplied | P2 | Confirmed: `uws1/execution.go` catches goto after stack unwinding; `uws1/execution_actions.go` resolves globally and `uws1/execution_runnable.go` records terminal/skipped state. Prior C04 wording did not define nested/root-call identity. Lineage: N2/C03. |
| F07 | not supplied | P2 | Partially confirmed: `uws1/execution_structural.go` returns both parent and iteration records today; zero/skipped dependencies need a parent fallback, but a failed dependency aborts merge before collection. Lineage: N11/C03. |

| Item | State | Notes |
|---|---|---|
| Correct goto and merge structural behavior | `[x]` | N2/N11/F06/F07. Implemented root-scope goto dispatch with UWS 1.11 fail-closed handling for terminal and in-flight step/workflow targets; tests cover whole-run unwind from top-level, sub-workflow, loop, and `forEach`, distinct caller-scoped records, root dependencies, and completed targets. UWS 1.11 merge keeps deepest iteration records, falls back to the parent for zero/skipped dependencies, and preserves dependency-failure propagation. UWS 1.10 compatibility remains covered. Focused verification: `go test ./uws1 -run 'Test(Goto|MergeForEach)' -count=1`. |
| Align expression grammar, numeric literals, and lint | `[x]` | N4/N5/N20/F01/F03. Added UWS 1.11-gated `$response.body.<path>`, loop-only `$batchIndex`, complete JSON-number literals for `wait`/`batchSize`, and version/context-aware content-trust parsing. Grammar lint executes YAML/HCL examples in `docs/03` and the versioned grammar fixture. The UWS 1.0 sample retains only four documented exact-JSON-Pointer exceptions; `testdata/big/big.json` has an exact file-level test-runtime exception and is explicitly not conformance evidence. Added `testdata/grammar/1.11.0.json`; schema/document validation was added with the 1.11 schema in the publication/artifact row. Verification: `go test ./contenttrust -count=1`, race equivalent, `go vet ./contenttrust`, strict MkDocs build, and `git diff --check`. |
| Publish corrected UWS 1.11 contract and artifacts | `[x]` | N9/N12/N13/N15/F04/F05 and N18 reference portion. Published `versions/1.11.0.{json,md}` plus the changelog entry, corrected step/variable/loop/criteria/version-gate/grammar/extension/browser-reference wording, and pinned the normative RFC/XPath/RE2/YAML/JSON Schema references. The 1.11 schema aligns `contentTrust.operations.*.outputs` keys with expression-addressable operation output names; the Go published-version registry, latest-artifact tests, digest manifest, and embedded JSON archive are synchronized. Exact-version admission and cross-version workflow-call isolation, merge ordering, and trigger-index corrections are documented without enabling 1.11 features on earlier documents. Prior published schema/spec hashes remain unchanged. Added a schema-and-semantic-validation test for the 1.11 grammar fixture. The following row added the executable conformance corpus and pinned its digest. Verification: `go test ./uws1 ./schemas ./contenttrust -count=1`, matching `go vet`, `mkdocs build --strict`, and `git diff --check`. |
| Make conformance vectors executable and protected | `[x]` | N6/N18. Added a typed UWS 1.11 corpus with `expression.parse`, `document.validate`, and `document.execute` operations; Go executes every core vector and `contenttrust` executes every parser vector. The UWS 1.10 corpus is now dispatched through Go semantic/execution tests instead of shape-only validation. Added version-gating, exact-version rejection, numeric wait/batch-size, goto, and merge cases. Execution observations use logical selectors, not Go record keys. The corpus SHA-256 is pinned in the spec and regression-tested. The vectors exposed and fixed exact `json.Number` handling for positive integral batch sizes, including fractional/overflow rejection; documented this representation correction without retroactively admitting 1.11 numeric-literal syntax. Verification: `go test ./uws1 ./contenttrust ./schemas -count=1`, race equivalent, `go vet ./uws1 ./contenttrust ./schemas`, `mkdocs build --strict`, `git diff --check`, and corpus/spec digest checks. |
| Synchronize guides and release surfaces | `[x]` | N8/N10. Updated docs 02 and 04–07, the home and validation/content-trust guides, README, `AGENTS.md`, MkDocs references, and the current product/status facts. Corrected portable loop/forEach/merge result shapes, trigger target order/concurrency, regex context and matching, retry counts/fallback behavior, root-scoped terminal `goto`, and `await`'s inability to re-execute completed status operations. Current release pointers now identify UWS 1.11 while preserving accurately historical 1.10 semantics. Full verification passed: `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check`. Whole-milestone review gate remains to be run. |

## Whole-Milestone Review Gate

- Iteration 1 started 2026-09-24 UTC after all five task rows and full
  acceptance verification passed.
- Scope: the complete C04 implementation and documentation change set; C05 and
  B02 are reviewed as completed dependencies, not reopened work.
- Findings (iteration 1):
  - **P2 — await predicate is treated as a numeric delay in advisory grammar
    checks.** `versions/1.11.0.md` §4.5.4, §4.5.5, §4.5.6, and §5.6 permit bare
    JSON-number strings for non-`await` `wait` fields only; `await.wait` remains
    a truthiness predicate. `contenttrust/analyze.go` enables numeric parsing
    for every path ending in `.wait`; `contenttrust/grammar_lint_test.go` and
    `contenttrust/conformance_vectors_test.go` do the same based only on the
    field name. Consequently these required advisory/lint/conformance checks
    accept a duration such as `"30"` as an `await.wait` predicate. This is a
    contract-accuracy and required-verification defect, not an execution
    validator finding. Fixed by making numeric-literal admission depend on
    field and structural construct, adding a rejected numeric-await conformance
    vector and parser/analyzer regressions, and updating the pinned corpus
    digest and its specification. Focused verification passed:
    `go test ./contenttrust ./uws1 -count=1`.
- Iteration 2 started 2026-09-24 UTC after the iteration 1 fix's focused
  verification passed. Scope: repeat the complete C04 implementation and
  documentation review, including the fix and its conformance updates.
- Iteration 2 completed 2026-09-24 UTC. The whole-milestone review found no
  remaining P1/P2 or higher-severity issue; the bounded review-fix gate passes
  after two iterations.
- Final acceptance verification passed after the review fix:
  `go test ./... -count=1`, `go test -race ./...`, `go vet ./...`,
  `mkdocs build --strict`, and `git diff --check`.
