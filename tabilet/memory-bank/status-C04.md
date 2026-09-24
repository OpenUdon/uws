# Status C04 - UWS 1.11 Contract Repair

**Goal.** Publish a coherent UWS 1.11 contract that closes the confirmed
control-flow, grammar, conformance, and guide gaps in 1.10 and incorporates
[C05](../docs/history/status-C05.md)'s exact-version admission correction.
This status records future work; no UWS 1.11 artifact exists because of this
reconciliation.

**Scope and contract choices.** In UWS 1.11, `goto` from a sub-workflow, loop,
or `forEach` item unwinds the entire top-level run. Remaining siblings and
iterations are skipped. The globally declared step or workflow target, with
its dependencies, is then invoked in root orchestration scope; it does not
select a prior caller-scoped workflow invocation. After unwinding, a root-key
target record with terminal status `success`, `error`, or `skipped` counts as
already completed and fails explicitly; an in-flight root target also fails
rather than waiting or replaying. A same-named record in another workflow-call
scope does not count. After the target completes, the top-level run ends; it
does not resume after the target. The 1.11 grammar will admit
`$response.body.<path>`, `$batchIndex` in loop contexts, and complete
JSON-number literals for numeric `wait` and `batchSize`, with version-aware
advisory content-trust parsing. A merge over `forEach` includes iteration
records when present, otherwise its parent record for a skipped or zero-item
dependency; it never counts both. A failed dependency still aborts merge.
`await` re-evaluates its predicate but does not re-execute completed status
operations; portable operation reexecution remains a UWS 2.0 candidate.
Correct genuine omissions and contradictions in a new version; preserve
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
their versions and receive exact-path, documented allowlist entries for
intentional legacy or §5.5 implementation-dependent expressions; they are not
represented as conforming 1.11 examples. Separate 1.11 fixtures prove the
new grammar. Test prior versions and regression-check C05's direct/file entry
paths. Run focused tests, `go test ./...`, `go test -race ./...`, `go vet ./...`,
`mkdocs build --strict`, `git diff --check`, then the bounded whole-milestone
review gate in [milestone.md](milestone.md).

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
| N10 | P2 | P2 documentation / P3 new feature | Partially confirmed: `uws1/execution_structural.go` re-evaluates `await` but memoizes completed status steps; `docs/02-Six-Structural-Constructs.md` suggests portable polling that does not occur. Reexecution is deferred, guide correction active. |
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
| Align expression grammar, numeric literals, and lint | `[x]` | N4/N5/N20/F01/F03. Added UWS 1.11-gated `$response.body.<path>`, loop-only `$batchIndex`, complete JSON-number literals for `wait`/`batchSize`, and version/context-aware content-trust parsing. Grammar lint executes YAML/HCL examples in `docs/03` and the versioned grammar fixture. The UWS 1.0 sample retains only four documented exact-JSON-Pointer exceptions; `testdata/big/big.json` has an exact file-level test-runtime exception and is explicitly not conformance evidence. Added `testdata/grammar/1.11.0.json`; schema/document validation is exercised with C04's publication/artifact task after its 1.11 schema exists. Verification: `go test ./contenttrust -count=1`, race equivalent, `go vet ./contenttrust`, strict MkDocs build, and `git diff --check`. |
| Publish corrected UWS 1.11 contract and artifacts | `[ ]` | N9/N12/N13/N15/N18/F04/F05. In new schema/spec and release history, audit §3.5 step definition; §4.5.1.1 variables; §4.5.6.1 items and §4.5.6.3 blank line; §5.1 loop scope; §5.2 and §4.5.11 criteria; §4.6 dangling “as listed above” and patch-level gates; §5.6 bare `~` and `$batchIndex`; §1.1 exact-version policy and replacement of the unsupported `1.0.0-beta.1` example; §6.1 extension registry; §7.4 browser references; header-name case; and XPath 1.0, RE2, YAML 1.2, and JSON Schema 2020-12 references. Preserve historically accurate “UWS 1.9” wording, document C05 admission and other cross-version bug-fix exceptions, align content-trust output keys, pin references/corpus by digest, and regenerate archive/hashes without changing frozen versions. |
| Make conformance vectors executable and protected | `[ ]` | N6/N18. Define typed vector operations/inputs, run every vector through Go with scripted execution where needed, add the missing structural/version/wait cases, and pin the corpus by digest. |
| Synchronize guides and release surfaces | `[ ]` | N8/N10. Update docs 02 and 04–07 plus affected examples/nav/README/AGENTS/current memory-bank facts; correct result shapes, trigger order, regex context, retry/goto wording, and `await` polling claims. Run grammar lint and full acceptance checks, then the whole-milestone review. |
