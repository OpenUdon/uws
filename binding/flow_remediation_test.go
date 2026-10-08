package binding

import (
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestFlowDependenciesUseKindAndWorkflow(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Operations: []*uws1.Operation{
		{OperationID: "same"},
		{OperationID: "read", OperationExecutionFields: uws1.OperationExecutionFields{DependsOn: []string{"same"}}},
	}, Workflows: []*uws1.Workflow{
		{WorkflowID: "same", Type: uws1.WorkflowTypeSequence},
		{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "same"}, {StepID: "run", OperationRef: "read", StepExecutionFields: uws1.StepExecutionFields{DependsOn: []string{"same"}}}}},
		{WorkflowID: "other", Steps: []*uws1.Step{{StepID: "same"}}},
	}}
	r, err := AnalyzeFlow(t.Context(), d)
	if err != nil {
		t.Fatal(err)
	}
	if hasCode(r, "flow.reference_ambiguous", "") || hasCode(r, "flow.reference_missing", "") {
		t.Fatal("cross-kind ambiguity", r)
	}
	if hasCode(r, "flow.unreachable", "/operations/0") || !hasCode(r, "flow.unreachable", "/workflows/0") || !hasCode(r, "flow.unreachable", "/workflows/2/steps/0") {
		t.Fatal("wrong namespace reached", r)
	}
}

func TestFlowOutputReferencesNeverFallBackToOtherWorkflow(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Workflows: []*uws1.Workflow{
		{WorkflowID: "main", Steps: []*uws1.Step{{StepID: "consume", Inputs: map[string]any{"value": "$steps.foreign.outputs.value"}}}},
		{WorkflowID: "other", Steps: []*uws1.Step{{StepID: "foreign", Outputs: map[string]string{"value": "$inputs.value"}}}},
	}}
	r, _ := AnalyzeFlow(t.Context(), d)
	if !hasCode(r, "flow.output_unreferenced", "/workflows/1/steps/0/outputs/value") {
		t.Fatal("foreign output marked used", r)
	}
	// Duplicate definitions in the same workflow remain genuine ambiguity.
	d.Workflows[0].Steps = []*uws1.Step{{StepID: "same"}, {StepID: "same"}, {StepID: "use", StepExecutionFields: uws1.StepExecutionFields{DependsOn: []string{"same"}}}}
	r, _ = AnalyzeFlow(t.Context(), d)
	if !hasCode(r, "flow.reference_ambiguous", "/workflows/0/steps/2/dependsOn") {
		t.Fatal(r)
	}
}
