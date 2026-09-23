package uws1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGotoFromNestedExecutionTransfersToRootAndEndsTopLevelRun(t *testing.T) {
	for _, source := range []string{"top-level", "sub-workflow", "loop", "forEach"} {
		t.Run(source, func(t *testing.T) {
			jump := &Operation{
				OperationID: "jump",
				OnSuccess:   []*SuccessAction{{Name: "transfer", Type: "goto", StepID: "target"}},
			}
			target := &Operation{OperationID: "target-op"}
			after := &Operation{OperationID: "after-op"}
			nestedAfter := &Operation{OperationID: "nested-after-op"}
			doc := testDocument(jump, target, after, nestedAfter)
			doc.UWS = "1.11.0"
			runtime := &mockRuntime{items: map[string][]any{"items": {"one", "two"}}}

			targetStep := &Step{StepID: "target", OperationRef: "target-op"}
			afterStep := &Step{StepID: "after", OperationRef: "after-op"}
			switch source {
			case "top-level":
				doc.Workflows = []*Workflow{{
					WorkflowID: "main", Type: WorkflowTypeSequence,
					Steps: []*Step{{StepID: "jump", OperationRef: "jump"}, afterStep, targetStep},
				}}
			case "sub-workflow":
				doc.Workflows = []*Workflow{
					{WorkflowID: "secondary", Type: WorkflowTypeSequence, Steps: []*Step{
						{StepID: "nested_jump", OperationRef: "jump"},
						{StepID: "nested_after", OperationRef: "nested-after-op"},
					}},
					{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{
						{StepID: "call", StepExecutionFields: StepExecutionFields{Workflow: "secondary"}},
						targetStep, afterStep,
					}},
				}
			case "loop":
				doc.Workflows = []*Workflow{{
					WorkflowID: "main", Type: WorkflowTypeLoop,
					StructuralFields: StructuralFields{Items: "items"},
					Steps:            []*Step{{StepID: "jump", OperationRef: "jump"}, afterStep, targetStep},
				}}
			case "forEach":
				doc.Workflows = []*Workflow{{
					WorkflowID: "main", Type: WorkflowTypeSequence,
					Steps: []*Step{
						{StepID: "each", OperationRef: "jump", StepExecutionFields: StepExecutionFields{ForEach: "items"}},
						afterStep, targetStep,
					},
				}}
			}

			err := NewOrchestrator(doc, runtime).Execute(context.Background())
			require.NoError(t, err)
			require.Equal(t, []string{"jump", "target-op"}, runtime.leafs())
			require.Equal(t, "success", doc.ExecutionRecords()["step:target"].Status)
			require.NotContains(t, doc.ExecutionRecords(), "step:after")
			if source == "sub-workflow" {
				require.NotContains(t, runtime.leafs(), "nested-after-op")
			}
			if source == "loop" {
				require.NotContains(t, doc.ExecutionRecords(), "step:after#iter:0")
			}
			if source == "forEach" {
				require.NotContains(t, doc.ExecutionRecords(), "step:each#iter:1")
			}
		})
	}
}

func TestGotoTargetUsesRootInvocationInsteadOfPriorWorkflowCallScope(t *testing.T) {
	jump := &Operation{
		OperationID: "jump",
		OnSuccess:   []*SuccessAction{{Name: "transfer", Type: "goto", StepID: "target"}},
	}
	doc := testDocument(jump, &Operation{OperationID: "dependency"}, &Operation{OperationID: "target-op"}, &Operation{OperationID: "after"})
	doc.UWS = "1.11.0"
	doc.Workflows = []*Workflow{
		{WorkflowID: "secondary", Type: WorkflowTypeSequence, Steps: []*Step{
			{StepID: "dependency", OperationRef: "dependency"},
			{StepID: "target", OperationRef: "target-op", StepExecutionFields: StepExecutionFields{DependsOn: []string{"dependency"}}},
		}},
		{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{
			{StepID: "call", StepExecutionFields: StepExecutionFields{Workflow: "secondary"}},
			{StepID: "jump", OperationRef: "jump"},
			{StepID: "after", OperationRef: "after"},
		}},
	}
	runtime := &mockRuntime{}

	require.NoError(t, NewOrchestrator(doc, runtime).Execute(context.Background()))
	require.Equal(t, []string{"dependency", "target-op", "jump", "dependency", "target-op"}, runtime.leafs())
	records := doc.ExecutionRecords()
	scopedTarget := workflowCallKey("secondary", "step:call") + "::step:target"
	require.Equal(t, "success", records[scopedTarget].Status)
	require.Equal(t, "success", records["step:target"].Status)
	require.NotContains(t, records, "step:after")
}

func TestGotoWorkflowTargetUsesRootInvocation(t *testing.T) {
	jump := &Operation{
		OperationID: "jump",
		OnSuccess:   []*SuccessAction{{Name: "transfer", Type: "goto", WorkflowID: "secondary"}},
	}
	doc := testDocument(jump, &Operation{OperationID: "target-op"}, &Operation{OperationID: "after"})
	doc.UWS = "1.11.0"
	doc.Workflows = []*Workflow{
		{WorkflowID: "secondary", Type: WorkflowTypeSequence, Steps: []*Step{{StepID: "target", OperationRef: "target-op"}}},
		{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{
			{StepID: "call", StepExecutionFields: StepExecutionFields{Workflow: "secondary"}},
			{StepID: "jump", OperationRef: "jump"},
			{StepID: "after", OperationRef: "after"},
		}},
	}
	runtime := &mockRuntime{}

	require.NoError(t, NewOrchestrator(doc, runtime).Execute(context.Background()))
	require.Equal(t, []string{"target-op", "jump", "target-op"}, runtime.leafs())
	records := doc.ExecutionRecords()
	require.Equal(t, "success", records[workflowCallKey("secondary", "step:call")].Status)
	require.Equal(t, "success", records["wf:secondary"].Status)
	require.NotContains(t, records, "step:after")
}

func TestGotoToCompletedRootTargetFailsOnlyFromUWS111(t *testing.T) {
	for _, tc := range []struct {
		version string
		wantErr bool
	}{
		{version: "1.10.0"},
		{version: "1.11.0", wantErr: true},
	} {
		t.Run(tc.version, func(t *testing.T) {
			doc := testDocument(&Operation{OperationID: "target-op"})
			doc.UWS = tc.version
			doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{
				{StepID: "target", OperationRef: "target-op"},
				{StepID: "jump", OperationRef: "jump"},
			}}}
			doc.Operations = append(doc.Operations, &Operation{
				OperationID: "jump",
				OnSuccess:   []*SuccessAction{{Name: "back", Type: "goto", StepID: "target"}},
			})
			runtime := &mockRuntime{}
			err := NewOrchestrator(doc, runtime).Execute(context.Background())
			if tc.wantErr {
				require.ErrorContains(t, err, `goto target step "target" is already completed`)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, []string{"target-op", "jump"}, runtime.leafs())
		})
	}
}

func TestGotoRejectsSkippedErrorAndInFlightRootTargetsInUWS111(t *testing.T) {
	for _, targetKind := range []string{"step", "workflow"} {
		key := "step:target"
		targetDescription := `step "target"`
		signal := &gotoSignal{stepID: "target"}
		if targetKind == "workflow" {
			key = "wf:target"
			targetDescription = `workflow "target"`
			signal = &gotoSignal{workflowID: "target"}
		}
		for _, status := range []string{"success", "error", "skipped"} {
			t.Run(targetKind+"/"+status, func(t *testing.T) {
				doc := testDocument(&Operation{OperationID: "target-op"})
				doc.UWS = "1.11.0"
				doc.Workflows = []*Workflow{
					{WorkflowID: "target", Type: WorkflowTypeSequence, Steps: []*Step{{StepID: "target", OperationRef: "target-op"}}},
					{WorkflowID: "main", Type: WorkflowTypeSequence},
				}
				orch := NewOrchestrator(doc, &mockRuntime{})
				kind := "workflow:sequence"
				if targetKind == "step" {
					kind = "step:operation"
				}
				orch.setRecord(key, ExecutionRecord{ID: "target", Kind: kind, Status: status, Error: "prior failure"})

				err := orch.executeWithSignals(context.Background(), func(context.Context) error { return signal })
				require.ErrorContains(t, err, "goto target "+targetDescription+" is already completed")
			})
		}

		t.Run(targetKind+"/in-flight", func(t *testing.T) {
			doc := testDocument(&Operation{OperationID: "target-op"})
			doc.UWS = "1.11.0"
			doc.Workflows = []*Workflow{
				{WorkflowID: "target", Type: WorkflowTypeSequence, Steps: []*Step{{StepID: "target", OperationRef: "target-op"}}},
				{WorkflowID: "main", Type: WorkflowTypeSequence},
			}
			orch := NewOrchestrator(doc, &mockRuntime{})
			inFlight := make(chan struct{})
			kind := "workflow:sequence"
			if targetKind == "step" {
				kind = "step:operation"
			}
			orch.mu.Lock()
			orch.inFlight[key] = inFlight
			orch.writeRecordLocked(key, ExecutionRecord{ID: "target", Kind: kind, Status: "running"})
			orch.mu.Unlock()

			done := make(chan error, 1)
			go func() {
				done <- orch.executeWithSignals(context.Background(), func(context.Context) error { return signal })
			}()
			select {
			case err := <-done:
				require.ErrorContains(t, err, "goto target "+targetDescription+" is in flight")
			case <-time.After(100 * time.Millisecond):
				close(inFlight)
				err := <-done
				t.Fatalf("goto waited for an in-flight target instead of failing; result %v", err)
			}
		})
	}
}

func TestGotoToCompletedRootWorkflowFailsFromUWS111(t *testing.T) {
	jump := &Operation{
		OperationID: "jump",
		OnSuccess:   []*SuccessAction{{Name: "repeat", Type: "goto", WorkflowID: "target"}},
	}
	doc := testDocument(jump, &Operation{OperationID: "target-op"})
	doc.UWS = "1.11.0"
	target := &Workflow{WorkflowID: "target", Type: WorkflowTypeSequence, Steps: []*Step{{StepID: "target-step", OperationRef: "target-op"}}}
	jumpStep := &Step{StepID: "jump", OperationRef: "jump"}
	doc.Workflows = []*Workflow{target, {WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{jumpStep}}}
	runtime := &mockRuntime{}
	orch := NewOrchestrator(doc, runtime)
	require.NoError(t, orch.ExecuteWorkflow(context.Background(), target))
	err := orch.executeWithSignals(context.Background(), func(ctx context.Context) error { return orch.ExecuteStep(ctx, jumpStep) })
	require.ErrorContains(t, err, `goto target workflow "target" is already completed`)
	require.Equal(t, []string{"target-op", "jump"}, runtime.leafs())
}

func TestMergeForEachUsesIterationRecordsOrParentFallbackByVersion(t *testing.T) {
	for _, tc := range []struct {
		name      string
		version   string
		items     []any
		when      string
		wantIDs   []string
		wantLeafs []string
	}{
		{name: "1.10 keeps parent and iterations", version: "1.10.0", items: []any{"a", "b"}, wantIDs: []string{"each", "each", "each"}, wantLeafs: []string{"each-op", "each-op"}},
		{name: "1.11 uses iterations", version: "1.11.0", items: []any{"a", "b"}, wantIDs: []string{"each", "each"}, wantLeafs: []string{"each-op", "each-op"}},
		{name: "1.11 zero items falls back to parent", version: "1.11.0", wantIDs: []string{"each"}},
		{name: "1.11 skipped dependency falls back to parent", version: "1.11.0", items: []any{"a"}, when: "run", wantIDs: []string{"each"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := testDocument(&Operation{OperationID: "each-op"})
			doc.UWS = tc.version
			doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{
				{StepID: "each", OperationRef: "each-op", StepExecutionFields: StepExecutionFields{ForEach: "items", When: tc.when}},
				{StepID: "join", Type: WorkflowTypeMerge, StepExecutionFields: StepExecutionFields{DependsOn: []string{"each"}}},
			}}}
			runtime := &mockRuntime{items: map[string][]any{"items": tc.items}}

			require.NoError(t, NewOrchestrator(doc, runtime).Execute(context.Background()))
			require.Equal(t, tc.wantLeafs, runtime.leafs())
			merged, ok := doc.ExecutionRecords()["step:join"].Result.([]map[string]any)
			require.True(t, ok)
			require.Len(t, merged, len(tc.wantIDs))
			for index, id := range tc.wantIDs {
				require.Equal(t, id, merged[index]["id"])
			}
			if tc.when != "" {
				require.Equal(t, "skipped", merged[0]["status"])
			} else if len(tc.items) == 0 {
				require.Equal(t, "success", merged[0]["status"])
				require.Empty(t, merged[0]["result"])
			} else if tc.version == "1.11.0" {
				for _, row := range merged {
					require.Nil(t, row["result"], "merge must use iteration records, not the parent aggregate")
					require.Equal(t, "success", row["status"])
				}
				parentRows, ok := doc.ExecutionRecords()["step:each"].Result.([]map[string]any)
				require.True(t, ok)
				require.Len(t, parentRows, len(tc.items))
				require.Len(t, merged, len(tc.items))
			}
		})
	}
}

func TestMergeForEachDropsOnlyIterationAncestors(t *testing.T) {
	keys := []string{
		"step:each",
		"step:each#iter:0",
		"step:each#iter:0.0",
		"step:each#iter:0.1",
		"step:each#iter:1",
	}
	require.Equal(t, []string{
		"step:each#iter:0.0",
		"step:each#iter:0.1",
		"step:each#iter:1",
	}, removeIterationAggregateKeys(keys))
	require.Equal(t, []string{"step:empty"}, removeIterationAggregateKeys([]string{"step:empty"}))
}

func TestMergeForEachFailureStillAbortsMerge(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "each-op"})
	doc.UWS = "1.11.0"
	doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{
		{StepID: "each", OperationRef: "each-op", StepExecutionFields: StepExecutionFields{ForEach: "items"}},
		{StepID: "join", Type: WorkflowTypeMerge, StepExecutionFields: StepExecutionFields{DependsOn: []string{"each"}}},
	}}}
	runtime := &mockRuntime{
		items: map[string][]any{"items": {"a"}},
		execute: func(context.Context, *Operation) error {
			return errors.New("dependency failed")
		},
	}

	err := NewOrchestrator(doc, runtime).Execute(context.Background())
	require.ErrorContains(t, err, "dependency failed")
	require.NotContains(t, doc.ExecutionRecords(), "step:join")
}
