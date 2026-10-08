package binding

import (
	"encoding/json"
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
	}}
	r, err := AnalyzeFlow(t.Context(), d)
	if err != nil {
		t.Fatal(err)
	}
	if hasCode(r, "flow.reference_ambiguous", "") || hasCode(r, "flow.reference_missing", "") {
		t.Fatal("cross-kind ambiguity", r)
	}
	if !hasCode(r, "flow.unreachable", "/operations/0") || !hasCode(r, "flow.unreachable", "/workflows/0") || hasCode(r, "flow.unreachable", "/workflows/1/steps/0") {
		t.Fatal("wrong namespace reached", r)
	}
}

func TestOperationReferencesUseInvocationWorkflow(t *testing.T) {
	var d uws1.Document
	text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"sourceDescriptions":[{"name":"api","type":"openapi","url":"fixture.json"}],"operations":[{"operationId":"read","sourceDescription":"api","sourceOperationId":"read","outputs":{"value":"$response.body"}},{"operationId":"send","sourceDescription":"api","sourceOperationId":"send","request":{"body":"$steps.fetch.outputs.value"}}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"fetch","operationRef":"read"},{"stepId":"deliver","operationRef":"send"}]}]}`
	if err := json.Unmarshal([]byte(text), &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	r, _ := AnalyzeFlow(t.Context(), &d)
	if hasCode(r, "flow.output_unreferenced", "/operations/0/outputs/value") {
		t.Fatal(r)
	}
}

func TestCrossDeclarationDependencyKeepsCallerFrame(t *testing.T) {
	var d uws1.Document
	text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"variables":{"value":"resolved"},"sourceDescriptions":[{"name":"api","type":"openapi","url":"fixture.json"}],"operations":[{"operationId":"read","sourceDescription":"api","sourceOperationId":"read"},{"operationId":"send","sourceDescription":"api","sourceOperationId":"send","request":{"body":"$steps.fetch.outputs.value"}}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"fetch","operationRef":"read","outputs":{"value":"$variables.value"}},{"stepId":"join","type":"merge","dependsOn":["foreign"]}]},{"workflowId":"helper","type":"sequence","steps":[{"stepId":"foreign","operationRef":"send"}]}]}`
	if err := json.Unmarshal([]byte(text), &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	r, _ := AnalyzeFlow(t.Context(), &d)
	if hasCode(r, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/value") {
		t.Fatal(r)
	}
}

func TestFlowRetainsGenericDependenciesAndParallelGroups(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Info: &uws1.Info{Title: "fixture", Version: "1"}, Operations: []*uws1.Operation{
		{OperationID: "first", Extensions: map[string]any{uws1.ExtensionOperationProfile: "fixture"}, OperationExecutionFields: uws1.OperationExecutionFields{ParallelGroup: "workers"}},
		{OperationID: "second", Extensions: map[string]any{uws1.ExtensionOperationProfile: "fixture"}, OperationExecutionFields: uws1.OperationExecutionFields{ParallelGroup: "workers"}},
	}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, WorkflowExecutionFields: uws1.WorkflowExecutionFields{DependsOn: []string{"workers"}}}}}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	r, _ := AnalyzeFlow(t.Context(), d)
	if hasCode(r, "flow.reference_missing", "") || hasCode(r, "flow.unreachable", "/operations/0") || hasCode(r, "flow.unreachable", "/operations/1") {
		t.Fatal(r)
	}
	d.Workflows[0].DependsOn = []string{"first"}
	r, _ = AnalyzeFlow(t.Context(), d)
	if hasCode(r, "flow.reference_missing", "") || hasCode(r, "flow.unreachable", "/operations/0") {
		t.Fatal(r)
	}
}

func TestParallelGroupMemberUsesGenericNamePrecedence(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Info: &uws1.Info{Title: "fixture", Version: "1"}, Operations: []*uws1.Operation{{OperationID: "read", Extensions: map[string]any{uws1.ExtensionOperationProfile: "fixture"}, OperationExecutionFields: uws1.OperationExecutionFields{ParallelGroup: "workers"}}}, Workflows: []*uws1.Workflow{
		{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, WorkflowExecutionFields: uws1.WorkflowExecutionFields{DependsOn: []string{"workers"}}},
		{WorkflowID: "helper", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "read", OperationRef: "read"}}},
	}}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	if err := d.ValidateExecutable(); err != nil {
		t.Fatal("executable fixture invalid", err)
	}
	r, _ := AnalyzeFlow(t.Context(), d)
	if hasCode(r, "flow.unreachable", "/workflows/1/steps/0") || hasCode(r, "flow.reference_missing", "") {
		t.Fatal(r)
	}
}

func TestFlowResultsUseExplicitWorkflowOwner(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Workflows: []*uws1.Workflow{{WorkflowID: "main", Steps: []*uws1.Step{{StepID: "join", Type: uws1.WorkflowTypeMerge, Outputs: map[string]string{"value": "$inputs.value"}}}}}, Results: []*uws1.StructuralResult{{From: "main.join", Value: "$steps.join.outputs.value"}}}
	r, _ := AnalyzeFlow(t.Context(), d)
	if hasCode(r, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/value") {
		t.Fatal(r)
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
