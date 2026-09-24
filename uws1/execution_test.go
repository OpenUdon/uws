package uws1

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRuntime struct {
	mu            sync.Mutex
	executedLeafs []string
	expressions   map[string]any
	items         map[string][]any
	execute       func(context.Context, *Operation) error
	eval          func(context.Context, string) (any, error)
}

func (m *mockRuntime) ExecuteLeaf(ctx context.Context, op *Operation) error {
	m.mu.Lock()
	m.executedLeafs = append(m.executedLeafs, op.OperationID)
	m.mu.Unlock()
	if m.execute != nil {
		return m.execute(ctx, op)
	}
	return nil
}

func (m *mockRuntime) leafs() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.executedLeafs...)
}

func (m *mockRuntime) ResolveItems(ctx context.Context, itemsExpr string) ([]any, error) {
	if m.items == nil {
		return nil, nil
	}
	return m.items[itemsExpr], nil
}

func (m *mockRuntime) EvaluateExpression(ctx context.Context, expr string) (any, error) {
	if m.eval != nil {
		return m.eval(ctx, expr)
	}
	if m.expressions == nil {
		return nil, nil
	}
	return m.expressions[expr], nil
}

func testDocument(ops ...*Operation) *Document {
	for _, op := range ops {
		if op == nil {
			continue
		}
		if op.Extensions == nil {
			op.Extensions = map[string]any{ExtensionOperationProfile: "test"}
		}
	}
	return &Document{
		UWS: "1.0.0",
		Info: &Info{
			Title:   "test",
			Version: "1.0.0",
		},
		Operations: ops,
	}
}

func TestOrchestratorExecuteSequenceWorkflow(t *testing.T) {
	doc := testDocument(
		&Operation{OperationID: "op1"},
		&Operation{OperationID: "op2"},
	)
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "step1", OperationRef: "op1"},
			{StepID: "step2", OperationRef: "op2"},
		},
	}}
	runtime := &mockRuntime{}
	doc.SetRuntime(runtime)

	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(runtime.executedLeafs) != 2 || runtime.executedLeafs[0] != "op1" || runtime.executedLeafs[1] != "op2" {
		t.Fatalf("unexpected execution order: %v", runtime.executedLeafs)
	}
}

func TestNonAwaitWaitDelaysRunnableOnlyFromUWS110(t *testing.T) {
	runtime := &mockRuntime{expressions: map[string]any{"delay": 0.02}}
	doc := waitTestDocument(&Operation{
		OperationID:              "delayed",
		OperationExecutionFields: OperationExecutionFields{Wait: "delay"},
	})
	doc.Runtime = runtime
	start := time.Now()
	require.NoError(t, doc.Execute(context.Background()))
	require.GreaterOrEqual(t, time.Since(start), 15*time.Millisecond)
	require.Equal(t, []string{"delayed"}, runtime.leafs())

	legacyRuntime := &mockRuntime{expressions: map[string]any{"delay": 1.0}}
	legacy := waitTestDocument(&Operation{
		OperationID:              "legacy",
		OperationExecutionFields: OperationExecutionFields{Wait: "delay"},
	})
	legacy.UWS = "1.9.2"
	legacy.Runtime = legacyRuntime
	start = time.Now()
	require.NoError(t, legacy.Execute(context.Background()))
	require.Less(t, time.Since(start), 100*time.Millisecond)
	require.Equal(t, []string{"legacy"}, legacyRuntime.leafs())
}

func TestNonAwaitWaitAcceptsUintptrNumber(t *testing.T) {
	runtime := &mockRuntime{expressions: map[string]any{"delay": uintptr(0)}}
	doc := waitTestDocument(&Operation{
		OperationID:              "zero_delay",
		OperationExecutionFields: OperationExecutionFields{Wait: "delay"},
	})
	doc.Runtime = runtime
	require.NoError(t, doc.Execute(context.Background()))
	require.Equal(t, []string{"zero_delay"}, runtime.leafs())
}

func TestNonAwaitWaitRejectsInvalidDurationsBeforeLeafExecution(t *testing.T) {
	for name, value := range map[string]any{
		"text": "1", "boolean": true, "negative": -1.0,
		"too large": float64(maxWaitSeconds) + 1, "nan": math.NaN(), "infinity": math.Inf(1),
	} {
		t.Run(name, func(t *testing.T) {
			runtime := &mockRuntime{expressions: map[string]any{"delay": value}}
			doc := waitTestDocument(&Operation{
				OperationID:              "invalid_wait",
				OperationExecutionFields: OperationExecutionFields{Wait: "delay"},
			})
			doc.Runtime = runtime
			err := doc.Execute(context.Background())
			require.ErrorContains(t, err, "wait duration")
			require.Empty(t, runtime.leafs())
		})
	}
}

func TestNonAwaitWaitHonorsContextCancellation(t *testing.T) {
	runtime := &mockRuntime{expressions: map[string]any{"delay": 1.0}}
	doc := waitTestDocument(&Operation{
		OperationID:              "cancelled_wait",
		OperationExecutionFields: OperationExecutionFields{Wait: "delay"},
	})
	doc.Runtime = runtime
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(10*time.Millisecond, cancel)
	err := doc.Execute(ctx)
	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, runtime.leafs())
}

func TestNonAwaitWaitAppliesToWorkflowStepAndOperation(t *testing.T) {
	runtime := &mockRuntime{expressions: map[string]any{"delay": 0.01}}
	doc := waitTestDocument(&Operation{
		OperationID:              "nested_wait",
		OperationExecutionFields: OperationExecutionFields{Wait: "delay"},
	})
	doc.Workflows[0].Wait = "delay"
	doc.Workflows[0].Steps[0].Wait = "delay"
	doc.Runtime = runtime
	start := time.Now()
	require.NoError(t, doc.Execute(context.Background()))
	require.GreaterOrEqual(t, time.Since(start), 25*time.Millisecond)
	require.Equal(t, []string{"nested_wait"}, runtime.leafs())
}

func TestAwaitWaitRemainsPredicateUnderUWS110(t *testing.T) {
	runtime := &mockRuntime{expressions: map[string]any{"ready": true}}
	doc := testDocument(&Operation{OperationID: "after_await"})
	doc.UWS = "1.10.0"
	doc.Workflows = []*Workflow{{
		WorkflowID:              "main",
		Type:                    WorkflowTypeAwait,
		WorkflowExecutionFields: WorkflowExecutionFields{Wait: "ready"},
		Steps:                   []*Step{{StepID: "run", OperationRef: "after_await"}},
	}}
	doc.Runtime = runtime
	start := time.Now()
	require.NoError(t, doc.Execute(context.Background()))
	require.Less(t, time.Since(start), 100*time.Millisecond)
	require.Equal(t, []string{"after_await"}, runtime.leafs())
}

func waitTestDocument(op *Operation) *Document {
	doc := testDocument(op)
	doc.UWS = "1.10.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{{
			StepID:       "run",
			OperationRef: op.OperationID,
		}},
	}}
	return doc
}

func TestStepInputsAreVisiblePerOperationInvocation(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "shared"})
	doc.UWS = "1.5.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "first", OperationRef: "shared", Inputs: map[string]any{"value": "one"}},
			{StepID: "second", OperationRef: "shared", Inputs: map[string]any{"value": "two"}},
		},
	}}
	var seen []map[string]any
	runtime := &mockRuntime{
		execute: func(ctx context.Context, op *Operation) error {
			state, ok := ExecutionContextFromContext(ctx)
			require.True(t, ok)
			seen = append(seen, cloneInputs(state.Inputs))
			return nil
		},
	}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.Execute(context.Background()))
	require.Equal(t, []string{"shared", "shared"}, runtime.executedLeafs)
	require.Equal(t, []map[string]any{
		{"value": "one"},
		{"value": "two"},
	}, seen)
	records := doc.ExecutionRecords()
	require.Contains(t, records, "stepop:first:shared")
	require.Contains(t, records, "stepop:second:shared")
}

func TestOrchestratorSkipsWhenFalse(t *testing.T) {
	doc := testDocument(
		&Operation{OperationID: "op1"},
		&Operation{OperationID: "op2"},
	)
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "step1", OperationRef: "op1", StepExecutionFields: RunnableExecutionFields{When: "false"}},
			{StepID: "step2", OperationRef: "op2", StepExecutionFields: RunnableExecutionFields{When: "true"}},
		},
	}}
	runtime := &mockRuntime{expressions: map[string]any{"false": false, "true": true}}
	doc.SetRuntime(runtime)

	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(runtime.executedLeafs) != 1 || runtime.executedLeafs[0] != "op2" {
		t.Fatalf("unexpected execution result: %v", runtime.executedLeafs)
	}
}

func TestOrchestratorParallelGroupDependencyBarrier(t *testing.T) {
	doc := testDocument(
		&Operation{OperationID: "op1", OperationExecutionFields: OperationExecutionFields{ParallelGroup: "grp"}},
		&Operation{OperationID: "op2", OperationExecutionFields: OperationExecutionFields{ParallelGroup: "grp"}},
		&Operation{OperationID: "op3"},
	)
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "step1", OperationRef: "op1"},
			{StepID: "step2", OperationRef: "op2"},
			{StepID: "step3", OperationRef: "op3", StepExecutionFields: RunnableExecutionFields{DependsOn: []string{"grp"}}},
		},
	}}
	runtime := &mockRuntime{}
	doc.SetRuntime(runtime)

	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(runtime.executedLeafs) != 3 {
		t.Fatalf("expected 3 executions, got %v", runtime.executedLeafs)
	}
}

func TestOrchestratorRejectsControlSignalsInsideParallel(t *testing.T) {
	for _, tc := range []struct {
		name   string
		action *SuccessAction
	}{
		{name: "end", action: &SuccessAction{Name: "stop", Type: "end"}},
		{name: "goto", action: &SuccessAction{Name: "jump", Type: "goto", StepID: "target"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := testDocument(
				&Operation{OperationID: "op1", OnSuccess: []*SuccessAction{tc.action}},
				&Operation{OperationID: "op2"},
			)
			doc.Workflows = []*Workflow{{
				WorkflowID: "main",
				Type:       WorkflowTypeParallel,
				Steps: []*Step{
					{StepID: "branch", OperationRef: "op1"},
					{StepID: "target", OperationRef: "op2"},
				},
			}}
			doc.SetRuntime(&mockRuntime{})

			err := doc.Execute(context.Background())
			require.Error(t, err)
			assert.Contains(t, err.Error(), "control signal")
			assert.Contains(t, err.Error(), "parallel workflow")
		})
	}
}

func TestOrchestratorExecuteSwitch(t *testing.T) {
	doc := testDocument(
		&Operation{OperationID: "op1"},
		&Operation{OperationID: "op2"},
	)
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSwitch,
		Cases: []*Case{
			{CaseFields: CaseFields{Name: "case1", When: "false"}, Steps: []*Step{{StepID: "s1", OperationRef: "op1"}}},
			{CaseFields: CaseFields{Name: "case2", When: "true"}, Steps: []*Step{{StepID: "s2", OperationRef: "op2"}}},
		},
	}}
	runtime := &mockRuntime{expressions: map[string]any{"false": false, "true": true}}
	doc.SetRuntime(runtime)

	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(runtime.executedLeafs) != 1 || runtime.executedLeafs[0] != "op2" {
		t.Fatalf("unexpected switch execution: %v", runtime.executedLeafs)
	}
}

func TestSwitchUsesDeclarationOrderAndFirstUnguardedCase(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "first"}, &Operation{OperationID: "later"})
	doc.UWS = "1.10.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSwitch,
		Cases: []*Case{
			{CaseFields: CaseFields{Name: "fallback"}, Steps: []*Step{{StepID: "first_step", OperationRef: "first"}}},
			{CaseFields: CaseFields{Name: "later", When: "true"}, Steps: []*Step{{StepID: "later_step", OperationRef: "later"}}},
		},
	}}
	runtime := &mockRuntime{expressions: map[string]any{"true": true}}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.Execute(context.Background()))
	assert.Equal(t, []string{"first"}, runtime.executedLeafs)
}

func TestOrchestratorExecuteLoop(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.Workflows = []*Workflow{{
		WorkflowID:       "main",
		Type:             WorkflowTypeLoop,
		StructuralFields: StructuralFields{Items: "items"},
		Steps: []*Step{
			{StepID: "step1", OperationRef: "op1"},
		},
	}}
	runtime := &mockRuntime{items: map[string][]any{"items": []any{1, 2, 3}}}
	doc.SetRuntime(runtime)

	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(runtime.executedLeafs) != 3 {
		t.Fatalf("expected 3 loop executions, got %v", runtime.executedLeafs)
	}
}

func TestLoopResultUsesOrderedItemRecordsAndBatchIndexes(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "noop"})
	doc.UWS = "1.10.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeLoop,
		StructuralFields: StructuralFields{
			Items:     "items",
			BatchSize: "batch",
		},
	}}
	doc.SetRuntime(&mockRuntime{
		items:       map[string][]any{"items": {"a", "b", "c"}},
		expressions: map[string]any{"batch": 2},
	})

	require.NoError(t, doc.Execute(context.Background()))
	result := doc.ExecutionRecords()["wf:main"].Result
	rows, ok := result.([]map[string]any)
	require.True(t, ok, "unexpected loop result type: %#v", result)
	require.Len(t, rows, 3)
	for i, want := range []struct {
		item       string
		batchIndex int
	}{{"a", 0}, {"b", 0}, {"c", 1}} {
		assert.Equal(t, i, rows[i]["index"])
		assert.Equal(t, want.batchIndex, rows[i]["batchIndex"])
		assert.Equal(t, want.item, rows[i]["item"])
	}
}

func TestResolveBatchSizeAcceptsIntegralJSONNumbers(t *testing.T) {
	for _, test := range []struct {
		value json.Number
		want  int
		valid bool
	}{
		{value: "2", want: 2, valid: true},
		{value: "2.0", want: 2, valid: true},
		{value: "2e0", want: 2, valid: true},
		{value: "2.5", valid: false},
		{value: "9007199254740993.5", valid: false},
		{value: "0", valid: false},
		{value: "-1", valid: false},
		{value: "1e1000", valid: false},
		{value: "9223372036854775808", valid: false},
	} {
		t.Run(string(test.value), func(t *testing.T) {
			runtime := &mockRuntime{expressions: map[string]any{"batch": test.value}}
			orch := NewOrchestrator(testDocument(), runtime)
			got, err := orch.resolveBatchSize(context.Background(), "batch")
			if !test.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestLoopWithNoItemsProducesEmptyResultArray(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "noop"})
	doc.UWS = "1.10.0"
	doc.Workflows = []*Workflow{{
		WorkflowID:       "main",
		Type:             WorkflowTypeLoop,
		StructuralFields: StructuralFields{Items: "items"},
	}}
	doc.SetRuntime(&mockRuntime{items: map[string][]any{"items": {}}})

	require.NoError(t, doc.Execute(context.Background()))
	rows, ok := doc.ExecutionRecords()["wf:main"].Result.([]map[string]any)
	require.True(t, ok, "unexpected empty loop result type: %#v", doc.ExecutionRecords()["wf:main"].Result)
	assert.Empty(t, rows)
}

func TestLoopWithNoItemsPreservesLegacyNilResult(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "noop"})
	doc.UWS = "1.9.2"
	doc.Workflows = []*Workflow{{
		WorkflowID:       "main",
		Type:             WorkflowTypeLoop,
		StructuralFields: StructuralFields{Items: "items"},
	}}
	doc.SetRuntime(&mockRuntime{items: map[string][]any{"items": {}}})

	require.NoError(t, doc.Execute(context.Background()))
	assert.Nil(t, doc.ExecutionRecords()["wf:main"].Result)
}

func TestForEachResultAndOutputsPreserveIterationOrder(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "op",
		Outputs:     map[string]string{"item": "$item"},
	})
	doc.UWS = "1.10.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main", Type: WorkflowTypeSequence,
		Steps: []*Step{{StepID: "each", OperationRef: "op", StepExecutionFields: StepExecutionFields{ForEach: "items"}, Outputs: map[string]string{"item": "$item"}}},
	}}
	doc.SetRuntime(&mockRuntime{
		items: map[string][]any{"items": {"a", "b", "c"}},
		eval: func(ctx context.Context, expr string) (any, error) {
			state, _ := ExecutionContextFromContext(ctx)
			if expr == "$item" && state != nil && state.Iteration != nil {
				return state.Iteration.Item, nil
			}
			return nil, nil
		},
	})

	require.NoError(t, doc.Execute(context.Background()))
	record := doc.ExecutionRecords()["step:each"]
	rows, ok := record.Result.([]map[string]any)
	require.True(t, ok, "unexpected forEach result type: %#v", record.Result)
	require.Len(t, rows, 3)
	assert.Equal(t, []any{"a", "b", "c"}, record.Outputs["item"])
	for i, row := range rows {
		assert.Equal(t, i, row["index"])
		assert.Equal(t, []any{"a", "b", "c"}[i], row["item"])
		assert.Equal(t, "success", row["status"])
	}
}

func TestOperationExecute_UsesOrchestratorSemantics(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "op1",
		Outputs: map[string]string{
			"status": "$response.body.status",
		},
	})
	doc.SetRuntime(&mockRuntime{
		expressions: map[string]any{
			"$response.body.status": "ok",
		},
	})

	require.NoError(t, doc.Operations[0].Execute(context.Background(), doc))

	records := doc.ExecutionRecords()
	require.Contains(t, records, "op:op1")
	assert.Equal(t, "success", records["op:op1"].Status)
	assert.Equal(t, "ok", records["op:op1"].Outputs["status"])
}

func TestWorkflowAndStepExecute_PersistExecutionRecords(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "leaf"})
	doc.SetRuntime(&mockRuntime{})
	wf := &Workflow{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{{
			StepID:       "step1",
			OperationRef: "leaf",
		}},
	}

	require.NoError(t, wf.Execute(context.Background(), doc))
	require.Contains(t, doc.ExecutionRecords(), "step:step1")

	doc.setExecutionRecords(nil)
	require.NoError(t, wf.Steps[0].Execute(context.Background(), doc))
	require.Contains(t, doc.ExecutionRecords(), "step:step1")
}

func TestOrchestratorExecuteAwait(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeAwait,
		WorkflowExecutionFields: WorkflowExecutionFields{
			Wait: "ready",
		},
		Steps: []*Step{{StepID: "step1", OperationRef: "op1"}},
	}}
	runtime := &mockRuntime{expressions: map[string]any{"ready": true}}
	doc.SetRuntime(runtime)

	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(runtime.executedLeafs) != 1 || runtime.executedLeafs[0] != "op1" {
		t.Fatalf("unexpected await execution: %v", runtime.executedLeafs)
	}
}

func TestDocumentExecuteRequiresEntryWorkflow(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.SetRuntime(&mockRuntime{})

	err := doc.Execute(context.Background())
	require.ErrorContains(t, err, "entry workflow")
}

func TestExplicitEntryPointsDoNotRequireDocumentEntryWorkflow(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.SetRuntime(&mockRuntime{})

	require.NoError(t, doc.Operations[0].Execute(context.Background(), doc))

	workflowDoc := testDocument(&Operation{OperationID: "op2"})
	workflowDoc.SetRuntime(&mockRuntime{})
	wf := &Workflow{
		WorkflowID: "secondary",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "step1", OperationRef: "op2"}},
	}
	require.NoError(t, wf.Execute(context.Background(), workflowDoc))

	stepDoc := testDocument(&Operation{OperationID: "op3"})
	stepDoc.SetRuntime(&mockRuntime{})
	require.NoError(t, (&Step{StepID: "step3", OperationRef: "op3"}).Execute(context.Background(), stepDoc))
}

func TestOrchestratorExecuteAwaitTimesOut(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.ExecutionOptions = ExecutionOptions{AwaitTimeout: 25 * time.Millisecond}
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeAwait,
		WorkflowExecutionFields: WorkflowExecutionFields{
			Wait: "ready",
		},
		Steps: []*Step{{StepID: "step1", OperationRef: "op1"}},
	}}
	doc.SetRuntime(&mockRuntime{expressions: map[string]any{"ready": false}})

	start := time.Now()
	err := doc.Execute(context.Background())
	require.Error(t, err)
	var timeoutErr *AwaitTimeoutError
	assert.True(t, errors.As(err, &timeoutErr))
	assert.GreaterOrEqual(t, time.Since(start), 20*time.Millisecond)
}

func TestOrchestratorExecuteAwaitHonorsContextCancellation(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.ExecutionOptions = ExecutionOptions{AwaitTimeout: time.Second}
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeAwait,
		WorkflowExecutionFields: WorkflowExecutionFields{
			Wait: "ready",
		},
		Steps: []*Step{{StepID: "step1", OperationRef: "op1"}},
	}}
	doc.SetRuntime(&mockRuntime{expressions: map[string]any{"ready": false}})

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	err := doc.Execute(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestOrchestratorDoesNotEvaluateAwaitAfterCancellation(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeAwait,
		WorkflowExecutionFields: WorkflowExecutionFields{
			Wait: "ready",
		},
	}}
	evaluations := 0
	doc.SetRuntime(&mockRuntime{eval: func(context.Context, string) (any, error) {
		evaluations++
		return false, nil
	}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.ErrorIs(t, doc.Execute(ctx), context.Canceled)
	assert.Zero(t, evaluations)
}

func TestSerializedAwaitTimeoutOverridesExecutorDefault(t *testing.T) {
	timeout := 0.05
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.UWS = "1.9.0"
	doc.ExecutionOptions = ExecutionOptions{
		AwaitTimeout:      5 * time.Millisecond,
		AwaitPollInterval: time.Millisecond,
	}
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeAwait,
		WorkflowExecutionFields: WorkflowExecutionFields{
			Wait:    "ready",
			Timeout: &timeout,
		},
	}}
	doc.SetRuntime(&mockRuntime{expressions: map[string]any{"ready": false}})

	start := time.Now()
	err := doc.Execute(context.Background())
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.GreaterOrEqual(t, time.Since(start), 35*time.Millisecond)
	var awaitErr *AwaitTimeoutError
	assert.False(t, errors.As(err, &awaitErr))
}

func TestRunnableTimeoutStartsAfterDependenciesAndWhen(t *testing.T) {
	timeout := 0.02
	doc := testDocument(&Operation{OperationID: "dependency"})
	doc.UWS = "1.9.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		WorkflowExecutionFields: WorkflowExecutionFields{
			DependsOn: []string{"dependency"},
			When:      "ready",
			Timeout:   &timeout,
		},
	}}
	doc.SetRuntime(&mockRuntime{
		execute: func(context.Context, *Operation) error {
			time.Sleep(35 * time.Millisecond)
			return nil
		},
		expressions: map[string]any{"ready": true},
	})

	require.NoError(t, doc.Execute(context.Background()))
}

func TestExecuteWithTimeoutSaturatesLargeDuration(t *testing.T) {
	timeout := 1e10
	called := false
	err := executeWithTimeout(context.Background(), &timeout, func(ctx context.Context) error {
		called = true
		return ctx.Err()
	})

	require.NoError(t, err)
	assert.True(t, called)
}

func TestWorkflowTimeoutCancelsDescendantOperation(t *testing.T) {
	timeout := 0.02
	doc := testDocument(&Operation{OperationID: "leaf"})
	doc.UWS = "1.9.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		WorkflowExecutionFields: WorkflowExecutionFields{
			Timeout: &timeout,
		},
		Steps: []*Step{{StepID: "leaf_step", OperationRef: "leaf"}},
	}}
	doc.SetRuntime(&mockRuntime{execute: func(ctx context.Context, _ *Operation) error {
		<-ctx.Done()
		return ctx.Err()
	}})

	require.ErrorIs(t, doc.Execute(context.Background()), context.DeadlineExceeded)
}

func TestOperationRetriesReceiveFreshTimeoutBudgets(t *testing.T) {
	timeout := 0.05
	doc := testDocument(&Operation{
		OperationID: "leaf",
		OperationExecutionFields: OperationExecutionFields{
			Timeout: &timeout,
		},
		OnFailure: []*FailureAction{{Name: "retry_once", Type: "retry", RetryLimit: 1}},
	})
	doc.UWS = "1.9.0"
	attempts := 0
	doc.SetRuntime(&mockRuntime{execute: func(context.Context, *Operation) error {
		attempts++
		time.Sleep(30 * time.Millisecond)
		if attempts == 1 {
			return errors.New("try again")
		}
		return nil
	}})

	require.NoError(t, doc.Operations[0].Execute(context.Background(), doc))
	assert.Equal(t, 2, attempts)
}

func TestOperationTimeoutCanBeRetried(t *testing.T) {
	timeout := 0.01
	doc := testDocument(&Operation{
		OperationID: "leaf",
		OperationExecutionFields: OperationExecutionFields{
			Timeout: &timeout,
		},
		OnFailure: []*FailureAction{{Name: "retry_once", Type: "retry", RetryLimit: 1}},
	})
	doc.UWS = "1.9.0"
	attempts := 0
	doc.SetRuntime(&mockRuntime{execute: func(ctx context.Context, _ *Operation) error {
		attempts++
		<-ctx.Done()
		return ctx.Err()
	}})

	require.ErrorIs(t, doc.Operations[0].Execute(context.Background(), doc), context.DeadlineExceeded)
	assert.Equal(t, 2, attempts)
}

func TestNestedLoopsUseFullIterationPathInRecordKeys(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "leaf"})
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeLoop,
		StructuralFields: StructuralFields{
			Items: "outer",
		},
		Steps: []*Step{{
			StepID: "inner",
			Type:   WorkflowTypeLoop,
			StructuralFields: StructuralFields{
				Items: "inner",
			},
			Steps: []*Step{{StepID: "leaf_step", OperationRef: "leaf"}},
		}},
	}}
	doc.SetRuntime(&mockRuntime{items: map[string][]any{
		"outer": {"a", "b"},
		"inner": {1, 2},
	}})

	require.NoError(t, doc.Execute(context.Background()))
	records := doc.ExecutionRecords()
	for _, key := range []string{
		"stepop:leaf_step:leaf#iter:0.0",
		"stepop:leaf_step:leaf#iter:0.1",
		"stepop:leaf_step:leaf#iter:1.0",
		"stepop:leaf_step:leaf#iter:1.1",
	} {
		require.Contains(t, records, key)
	}
	assert.Len(t, doc.Runtime.(*mockRuntime).leafs(), 4)
}

func TestOrchestratorForEachAggregatesOutputsAndResults(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "op1",
		OperationExecutionFields: OperationExecutionFields{
			ForEach: "items",
		},
		Outputs: map[string]string{
			"value": "$item",
		},
	})
	runtime := &mockRuntime{
		items: map[string][]any{
			"items": {1, 2, 3},
		},
		eval: func(ctx context.Context, expr string) (any, error) {
			if expr != "$item" {
				return nil, nil
			}
			state, ok := ExecutionContextFromContext(ctx)
			require.True(t, ok)
			require.NotNil(t, state.Iteration)
			return state.Iteration.Item, nil
		},
	}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.Operations[0].Execute(context.Background(), doc))

	records := doc.ExecutionRecords()
	record := records["op:op1"]
	require.Equal(t, "success", record.Status)
	results, ok := record.Result.([]map[string]any)
	require.True(t, ok)
	require.Len(t, results, 3)
	assert.Equal(t, 0, results[0]["index"])
	assert.Equal(t, 1, results[0]["item"])
	assert.Equal(t, []any{1, 2, 3}, record.Outputs["value"])
	require.Contains(t, records, "op:op1#iter:0")
	assert.Equal(t, 1, records["op:op1#iter:0"].Outputs["value"])
}

func TestOrchestratorForEachRecordsControlSignalIterationAsSuccess(t *testing.T) {
	for _, tc := range []struct {
		name   string
		action *SuccessAction
	}{
		{name: "end", action: &SuccessAction{Name: "stop", Type: "end"}},
		{name: "goto", action: &SuccessAction{Name: "jump", Type: "goto", StepID: "target"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := testDocument(
				&Operation{OperationID: "op1", OnSuccess: []*SuccessAction{tc.action}},
				&Operation{OperationID: "op2"},
			)
			doc.Workflows = []*Workflow{{
				WorkflowID: "main",
				Type:       WorkflowTypeSequence,
				Steps: []*Step{
					{
						StepID:       "each",
						OperationRef: "op1",
						StepExecutionFields: StepExecutionFields{
							ForEach: "items",
						},
					},
					{StepID: "target", OperationRef: "op2"},
				},
			}}
			doc.SetRuntime(&mockRuntime{items: map[string][]any{"items": {"one", "two"}}})

			require.NoError(t, doc.Execute(context.Background()))
			record := doc.ExecutionRecords()["step:each#iter:0"]
			require.Equal(t, "success", record.Status)
			require.Empty(t, record.Error)
		})
	}
}

func TestOrchestratorExecuteStepWorkflowReference(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "leaf"})
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "secondary",
			Type:       WorkflowTypeSequence,
			Steps: []*Step{{
				StepID:       "child",
				OperationRef: "leaf",
			}},
		},
		{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps: []*Step{{
				StepID: "call_secondary",
				StepExecutionFields: StepExecutionFields{
					Workflow: "secondary",
				},
			}},
		},
	}
	runtime := &mockRuntime{}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.Execute(context.Background()))
	assert.Equal(t, []string{"leaf"}, runtime.executedLeafs)
	records := doc.ExecutionRecords()
	require.Contains(t, records, "step:call_secondary")
	require.Contains(t, records, workflowCallKey("secondary", "step:call_secondary"))
}

func TestWorkflowCallsFromDifferentStepsUseDistinctInputsAndRecords(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "leaf"}, &Operation{OperationID: "dependency"})
	doc.UWS = "1.5.0"
	doc.Workflows = []*Workflow{
		{
			WorkflowID:              "secondary",
			Type:                    WorkflowTypeSequence,
			WorkflowExecutionFields: WorkflowExecutionFields{DependsOn: []string{"dependency"}},
			Steps: []*Step{
				{StepID: "child", OperationRef: "leaf", Outputs: map[string]string{"value": "input"}},
				{StepID: "join", Type: WorkflowTypeMerge, StepExecutionFields: StepExecutionFields{DependsOn: []string{"child"}}},
			},
			Outputs: map[string]string{"received": "child-value"},
		},
		{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps: []*Step{
				{StepID: "first_call", Inputs: map[string]any{"value": "first"}, StepExecutionFields: StepExecutionFields{Workflow: "secondary"}},
				{StepID: "second_call", Inputs: map[string]any{"value": "second"}, StepExecutionFields: StepExecutionFields{Workflow: "secondary"}},
			},
		},
	}

	seenInputs := map[string][]any{}
	runtime := &mockRuntime{
		execute: func(ctx context.Context, op *Operation) error {
			state, _ := ExecutionContextFromContext(ctx)
			seenInputs[op.OperationID] = append(seenInputs[op.OperationID], state.Inputs["value"])
			return nil
		},
		eval: func(ctx context.Context, expr string) (any, error) {
			state, _ := ExecutionContextFromContext(ctx)
			if expr == "child-value" {
				child, ok := state.Records["step:child"]
				if !ok {
					return nil, errors.New("child step is not visible in the called workflow scope")
				}
				return child.Outputs["value"], nil
			}
			return state.Inputs["value"], nil
		},
	}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.Execute(context.Background()))
	assert.Equal(t, []any{"first", "second"}, seenInputs["dependency"])
	assert.Equal(t, []any{"first", "second"}, seenInputs["leaf"])

	records := doc.ExecutionRecords()
	first := records[workflowCallKey("secondary", "step:first_call")]
	second := records[workflowCallKey("secondary", "step:second_call")]
	require.Equal(t, "workflow:sequence", first.Kind)
	require.Equal(t, "workflow:sequence", second.Kind)
	assert.Equal(t, "first", first.Outputs["received"])
	assert.Equal(t, "second", second.Outputs["received"])
	for _, callKey := range []string{
		workflowCallKey("secondary", "step:first_call"),
		workflowCallKey("secondary", "step:second_call"),
	} {
		join, ok := records[callKey+"::step:join"]
		require.True(t, ok, "missing merge record for workflow invocation %q", callKey)
		merged, ok := join.Result.([]map[string]any)
		require.True(t, ok, "unexpected scoped merge result: %#v", join.Result)
		require.Len(t, merged, 1, "merge must include only its invocation's child records")
	}
}

func TestRecursiveWorkflowCallFailsInsteadOfWaitingOnItself(t *testing.T) {
	doc := testDocument()
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{{
			StepID:              "again",
			StepExecutionFields: StepExecutionFields{Workflow: "main"},
		}},
	}}
	err := NewOrchestrator(doc, &mockRuntime{}).Execute(context.Background())
	require.ErrorContains(t, err, `recursive workflow invocation "main"`)
}

func TestOrchestratorExecuteMerge(t *testing.T) {
	doc := testDocument(
		&Operation{OperationID: "op1"},
		&Operation{OperationID: "op2"},
	)
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "left", OperationRef: "op1"},
			{StepID: "right", OperationRef: "op2"},
			{
				StepID: "join",
				Type:   WorkflowTypeMerge,
				StepExecutionFields: RunnableExecutionFields{
					DependsOn: []string{"left", "right"},
				},
			},
		},
	}}
	runtime := &mockRuntime{}
	doc.SetRuntime(runtime)

	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if len(runtime.executedLeafs) != 2 {
		t.Fatalf("unexpected merge execution: %v", runtime.executedLeafs)
	}
}

func TestOrchestratorExecuteMergeUsesDeclaredDependenciesOnly(t *testing.T) {
	doc := testDocument(
		&Operation{OperationID: "prep"},
		&Operation{OperationID: "left"},
		&Operation{OperationID: "right"},
	)
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "prep", OperationRef: "prep"},
			{StepID: "left", OperationRef: "left"},
			{StepID: "right", OperationRef: "right"},
			{
				StepID: "join",
				Type:   WorkflowTypeMerge,
				StepExecutionFields: RunnableExecutionFields{
					DependsOn: []string{"left", "right"},
				},
			},
		},
	}}
	orch := NewOrchestrator(doc, &mockRuntime{})

	if err := orch.Execute(context.Background()); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	record := orch.records["step:join"]
	results, ok := record.Result.([]map[string]any)
	if !ok {
		t.Fatalf("unexpected merge result type: %#v", record.Result)
	}
	if len(results) != 2 {
		t.Fatalf("expected exactly declared dependencies, got %#v", results)
	}
	if results[0]["id"] != "left" || results[1]["id"] != "right" {
		t.Fatalf("unexpected dependency order in merge result: %#v", results)
	}
}

func TestMergeDependencyRecordsSortsNestedIterationIndexesNumerically(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op"})
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "leaf_step", OperationRef: "op"}},
	}}
	orch := NewOrchestrator(doc, &mockRuntime{})
	stepBase := stepKey("leaf_step")
	operationBase := stepOperationKey("leaf_step", "op")
	for i := 0; i < 12; i++ {
		path := []int{i, 0}
		orch.setRecord(compositeIterationKey(stepBase, path), ExecutionRecord{
			ID: "leaf_step", Kind: "step:operation", Status: "success",
			Outputs: map[string]any{"index": i},
		})
		orch.setRecord(compositeIterationKey(operationBase, path), ExecutionRecord{
			ID: "op", Kind: "operation", Status: "success",
			Outputs: map[string]any{"index": i},
		})
	}

	for _, dependency := range []string{"leaf_step", "op"} {
		merged := orch.mergeDependencyRecords(context.Background(), []string{dependency})
		require.Len(t, merged, 12)
		for want, record := range merged {
			outputs, ok := record["outputs"].(map[string]any)
			require.True(t, ok, "dependency %q returned outputs %#v", dependency, record["outputs"])
			require.Equal(t, want, outputs["index"], "dependency %q order at position %d", dependency, want)
		}
	}
}

func TestOrchestratorColonOperationDependencyDoesNotMatchSuffix(t *testing.T) {
	doc := testDocument(
		&Operation{OperationID: "a:b"},
		&Operation{OperationID: "b"},
	)
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "colon", OperationRef: "a:b"},
			{
				StepID: "join",
				Type:   WorkflowTypeMerge,
				StepExecutionFields: RunnableExecutionFields{
					DependsOn: []string{"b"},
				},
			},
		},
	}}
	runtime := &mockRuntime{}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.Execute(context.Background()))
	assert.Equal(t, []string{"a:b", "b"}, runtime.executedLeafs)

	record := doc.ExecutionRecords()["step:join"]
	results, ok := record.Result.([]map[string]any)
	require.True(t, ok, "unexpected merge result type: %#v", record.Result)
	require.Len(t, results, 1)
	assert.Equal(t, "b", results[0]["id"])
	assert.Equal(t, "operation", results[0]["kind"])
}

func TestDocumentValidateExecutableRejectsAmbiguousIDs(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "op1"})
	doc.Workflows = []*Workflow{{
		WorkflowID: "shared",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "shared", OperationRef: "op1"}},
	}, {
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "root", StepExecutionFields: StepExecutionFields{Workflow: "shared"}}},
	}}

	err := doc.ValidateExecutable()
	if err == nil {
		t.Fatal("expected executable validation error")
	}
	if got := err.Error(); got == "" || got == "shared" {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestDocumentValidateExecutableAllowsNonSimpleCriteria(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "op1",
		SuccessCriteria: []*Criterion{{
			Type:      CriterionRegex,
			Context:   "$response.body",
			Condition: "^ok",
		}},
	})

	require.NoError(t, doc.ValidateExecutable())
}

func TestDocumentValidateExecutableAllowsOutputsAndResults(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "fetch",
		Outputs: map[string]string{
			"body": "$response.body",
		},
	})
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{
			{StepID: "fetch_step", OperationRef: "fetch"},
			{
				StepID: "joined_step",
				Type:   WorkflowTypeMerge,
				StepExecutionFields: RunnableExecutionFields{
					DependsOn: []string{"fetch_step"},
				},
				Outputs: map[string]string{
					"body": "$steps.fetch_step.outputs.body",
				},
			},
		},
	}}
	doc.Results = []*StructuralResult{{
		Name:  "joined_result",
		Kind:  StructuralResultKindMerge,
		From:  "main.joined_step",
		Value: "$steps.joined_step.outputs.body",
	}}

	require.NoError(t, doc.ValidateExecutable())
}
