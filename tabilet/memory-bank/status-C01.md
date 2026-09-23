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

**Review gate.** Ordinary intake; 0 of 10 iterations. Start only after task
acceptance checks.

| Item | State | Notes |
|---|---|---|
| Distinguish workflow call invocations | `[+]` | Completed 2026-09-23. Workflow calls now use a deterministic caller-step identity; nested steps, operations, dependencies, and merge records are scoped to that invocation. Entry workflow keys remain unchanged, and dependency memoization remains within its caller scope. Recursive calls fail clearly rather than waiting on their own in-flight records. Regression tests verify distinct inputs, outputs, child execution, scoped merge results, and recursion failure. `go test ./uws1 -run 'TestOrchestratorExecuteStepWorkflowReference|TestWorkflowCallsFromDifferentStepsUseDistinctInputsAndRecords|TestRecursiveWorkflowCallFailsInsteadOfWaitingOnItself' -count=1`, `go test ./uws1 -count=1`, `go test ./uws1 -race`, `go test ./...`, `go vet ./...`, and `git diff --check` passed. Owner: A2. |
| Preserve numeric iteration order | `[ ]` | Sort iteration identity numerically where merge/dependency results are collected; cover 12 or more iterations and nested cases. Owner: A6. |
