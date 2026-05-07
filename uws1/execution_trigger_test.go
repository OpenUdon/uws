package uws1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentDispatchTriggerExecutesTopLevelStepTargets(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "fetch"}, &Operation{OperationID: "save"})
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "secondary",
			Type:       WorkflowTypeSequence,
			Steps:      []*Step{{StepID: "secondary_fetch", OperationRef: "fetch"}},
		},
		{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps: []*Step{
				{StepID: "fetch_step", OperationRef: "fetch"},
				{StepID: "save_step", OperationRef: "save", StepExecutionFields: StepExecutionFields{DependsOn: []string{"fetch_step"}}},
			},
		},
	}
	doc.Triggers = []*Trigger{{
		TriggerID: "incoming",
		Outputs:   []string{"primary"},
		Routes: []*TriggerRoute{
			{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"save_step"}}},
		},
	}}
	runtime := &mockRuntime{
		eval: func(ctx context.Context, expr string) (any, error) {
			switch expr {
			case "$trigger.kind":
				state, ok := ExecutionContextFromContext(ctx)
				require.True(t, ok)
				require.NotNil(t, state.Trigger)
				payload, _ := state.Trigger.Payload.(map[string]any)
				return payload["kind"], nil
			default:
				return nil, nil
			}
		},
	}
	doc.Operations[1].Outputs = map[string]string{"kind": "$trigger.kind"}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.DispatchTrigger(context.Background(), "incoming", 0, map[string]any{"kind": "webhook"}))
	assert.Equal(t, []string{"fetch", "save"}, runtime.executedLeafs)
	records := doc.ExecutionRecords()
	require.Contains(t, records, "step:save_step")
	assert.Equal(t, "webhook", records["op:save"].Outputs["kind"])
}

func TestDocumentDispatchTriggerExecutesWorkflowTargets(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "fetch"})
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "secondary",
			Type:       WorkflowTypeSequence,
			Steps:      []*Step{{StepID: "secondary_fetch", OperationRef: "fetch"}},
		},
		{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps:      []*Step{{StepID: "root", StepExecutionFields: StepExecutionFields{Workflow: "secondary"}}},
		},
	}
	doc.Triggers = []*Trigger{{
		TriggerID: "incoming",
		Outputs:   []string{"primary"},
		Routes: []*TriggerRoute{
			{TriggerRouteFields: TriggerRouteFields{Output: "0", To: []string{"secondary"}}},
		},
	}}
	runtime := &mockRuntime{}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.DispatchTrigger(context.Background(), "incoming", 0, map[string]any{"ok": true}))
	assert.Equal(t, []string{"fetch"}, runtime.executedLeafs)
	require.Contains(t, doc.ExecutionRecords(), "wf:secondary")
}

func TestDocumentDispatchTriggerRejectsUnknownTarget(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "fetch"})
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "fetch_step", OperationRef: "fetch"}},
	}}
	doc.Triggers = []*Trigger{{
		TriggerID: "incoming",
		Outputs:   []string{"primary"},
		Routes: []*TriggerRoute{
			{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"fetch"}}},
		},
	}}
	doc.SetRuntime(&mockRuntime{})

	err := doc.DispatchTrigger(context.Background(), "incoming", 0, nil)
	require.ErrorContains(t, err, "top-level stepId or workflowId")
}
