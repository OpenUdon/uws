# C01 History

Milestone: `C01`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-C01.md`
Source specification: `tabilet/memory-bank/milestone.md#c01---workflow-execution-correctness`
Evidence: `135221019ab575e3bbcc5ca2d6c79a9407e4da17`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `2`
Verification: `go test ./...`; `go test -race ./...`; `go vet ./...`; `git diff --check`
Consolidated into: [workflow execution](../../../uws1/execution.go), [workflow-call regression tests](../../../uws1/execution_test.go), [numeric merge ordering](../../../uws1/execution_structural.go), [downstream contract](../../memory-bank/status-C03.md)

## Final milestone specification

````markdown
## C01 - Workflow Execution Correctness

**Goal.** Make independent workflow calls execute independently and preserve
numeric iteration order in merged results.

**Scope.** Separate workflow invocation identity from entry/dependency
memoization so two calling steps can pass distinct inputs and receive distinct
results. Order iteration keys numerically, not lexically, including after ten
iterations. Keep deliberate memoization and dependency behavior intact.

**Acceptance.** Focused tests cover two calls to the same workflow and merges
with at least twelve iterations; full tests, race tests, vet, and diff check
pass. No historical contract is silently rewritten.

**Dependencies.** None. **Downstream impact.** C03 documents the verified
invocation and result-order semantics.
````

## Final status

````markdown
# Status C01 - Workflow Execution Correctness

**Acceptance.** Two separate steps calling one workflow with different inputs
run independently and receive distinct results; dependency and entry
memoization remain intentional. Merge results after at least twelve loop
iterations use numeric order. Focused tests, full tests, race tests, vet, and
diff check pass.

**Dependencies.** None. **Downstream impacts.** C03 documents verified
invocation and ordered-result semantics.

**Review provenance.** “UWS spec and architecture review: 1.9.2 and earlier
(2026-09-23)”; stated baseline
`8382d0f26b3b10870125760643078d1a1a3e31b6`; revalidated at full HEAD
`e765fed4ba0a241e9481f99fc324c6a8afb9a8a1` with untracked
memory-bank initialization content present. No relevant uncommitted executor
change formed the evidence.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| A2 | P0 | P1 | Confirmed: workflow execution key uses workflow ID in `uws1/execution.go` and `uws1/execution_runnable.go`, collapsing distinct calling steps. |
| A6 | P1 | P1 | Confirmed: lexical key sorting in `uws1/execution_structural.go` misorders iteration 10 before 2. |

**Review gate.** Iteration 1 of 10 found G1 (P1): called-workflow dependencies
and `$steps` expression snapshots were not fully scoped. The remediation row
below is complete and its focused, full, race, vet, and diff checks pass.
Iteration 2 of 10 passed 2026-09-23. Whole-milestone review found no new P1,
P2, or higher-severity issue. The G1 fix was specifically re-reviewed along
with caller identity, dependency memoization, scoped records, recursion
handling, and nested numeric iteration ordering.

| Item | State | Notes |
|---|---|---|
| Distinguish workflow call invocations | `[+]` | Completed 2026-09-23. Workflow calls now use a deterministic caller-step identity; nested steps, operations, dependencies, and merge records are scoped to that invocation. Entry workflow keys remain unchanged, and dependency memoization remains within its caller scope. Recursive calls fail clearly rather than waiting on their own in-flight records. Regression tests verify distinct inputs, outputs, child execution, scoped merge results, and recursion failure. `go test ./uws1 -run 'TestOrchestratorExecuteStepWorkflowReference|TestWorkflowCallsFromDifferentStepsUseDistinctInputsAndRecords|TestRecursiveWorkflowCallFailsInsteadOfWaitingOnItself' -count=1`, `go test ./uws1 -count=1`, `go test ./uws1 -race`, `go test ./...`, `go vet ./...`, and `git diff --check` passed. Owner: A2. |
| Preserve numeric iteration order | `[+]` | Completed 2026-09-23. Dependency and operation invocation records now sort by base key and numeric iteration path, including nested paths. Tests verify step and operation merge results in order across twelve nested iteration keys. Focused merge test, `go test ./...`, `go test -race ./...`, `go vet ./...`, and `git diff --check` passed. Owner: A6. |
| Scope workflow-call dependencies and expression records | `[+]` | Completed 2026-09-23. Workflow-level dependencies and output evaluation now use the call's scope; runtime expression snapshots expose local child records under their ordinary step keys. Regression coverage proves both callers execute a workflow dependency with their own inputs and each workflow output reads its own child-step value. `go test ./uws1 -run 'TestOrchestratorExecuteStepWorkflowReference|TestWorkflowCallsFromDifferentStepsUseDistinctInputsAndRecords|TestRecursiveWorkflowCallFailsInsteadOfWaitingOnItself|TestMergeDependencyRecordsSortsNestedIterationIndexesNumerically' -count=1`, `go test ./...`, `go test -race ./...`, `go vet ./...`, and `git diff --check` passed. Review finding G1 is resolved. |
````
