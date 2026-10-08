package binding

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/expressions"
	"github.com/OpenUdon/uws/uws1"
)

type flowReferenceRuntime struct{ *expressions.Evaluator }

func (r *flowReferenceRuntime) ExecuteLeaf(context.Context, *uws1.Operation) error { return nil }
func (r *flowReferenceRuntime) EvaluateExpression(ctx context.Context, text string) (any, error) {
	field := expressions.Value
	if !strings.HasPrefix(text, "$") {
		field = expressions.Wait
	}
	return r.Evaluate(ctx, text, field)
}

func TestGotoUsesGloballyIndexedStepAndRootContext(t *testing.T) {
	var d uws1.Document
	text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"variables":{"items":[1]},"operations":[{"operationId":"read","x-uws-operation-profile":"fixture","outputs":{"batch":"$batchIndex"}},{"operationId":"jump","x-uws-operation-profile":"fixture","onSuccess":[{"name":"transfer","type":"goto","stepId":"target"}]}],"workflows":[{"workflowId":"main","type":"loop","items":"$variables.items","steps":[{"stepId":"call","workflow":"helper"},{"stepId":"transfer","operationRef":"jump"}]},{"workflowId":"helper","type":"sequence","steps":[{"stepId":"target","operationRef":"read"}]}]}`
	if err := json.Unmarshal([]byte(text), &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	if err := d.ValidateExecutable(); err != nil {
		t.Fatal("executable fixture invalid", err)
	}
	r, err := AnalyzeFlow(t.Context(), &d)
	if err != nil || hasCode(r, "flow.reference_missing", "") {
		t.Fatal("global target lost", r, err)
	}
	if got := expressions.CheckPortability(&d); len(got) != 1 || got[0].Path != "/operations/0/outputs/batch" || got[0].Code != "expression.context" {
		t.Fatal("root goto retained loop context", got)
	}
	e, err := expressions.NewEvaluator(&d)
	if err != nil {
		t.Fatal(err)
	}
	d.SetRuntime(&flowReferenceRuntime{e})
	if err := d.Execute(t.Context()); !errors.Is(err, expressions.ErrContext) {
		t.Fatal("root reference execution did not reject batchIndex", err)
	}
	d.Operations[1].OnSuccess = nil
	d.Operations[1].OnFailure = []*uws1.FailureAction{{Name: "transfer", Type: "goto", StepID: "target"}}
	if got := expressions.CheckPortability(&d); len(got) != 1 || got[0].Code != "expression.context" {
		t.Fatal("failure target lost", got)
	}
}

func TestGotoDoesNotInheritHelperOutputFrame(t *testing.T) {
	var d uws1.Document
	text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"variables":{"value":"fixture"},"operations":[{"operationId":"noop","x-uws-operation-profile":"fixture"},{"operationId":"read","x-uws-operation-profile":"fixture","outputs":{"value":"$steps.fetch.outputs.value"}},{"operationId":"jump","x-uws-operation-profile":"fixture","onSuccess":[{"name":"transfer","type":"goto","stepId":"target"}]}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"targets","type":"merge","dependsOn":["noop"],"steps":[{"stepId":"target","operationRef":"read"}]},{"stepId":"call","workflow":"helper"}]},{"workflowId":"helper","type":"sequence","steps":[{"stepId":"fetch","operationRef":"noop","outputs":{"value":"$variables.value"}},{"stepId":"transfer","operationRef":"jump"}]}]}`
	if err := json.Unmarshal([]byte(text), &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	if err := d.ValidateExecutable(); err != nil {
		t.Fatal("executable fixture invalid", err)
	}
	r, err := AnalyzeFlow(t.Context(), &d)
	if err != nil || !hasCode(r, "flow.output_unreferenced", "/workflows/1/steps/0/outputs/value") {
		t.Fatal("helper output attributed to root target", r, err)
	}
	e, err := expressions.NewEvaluator(&d)
	if err != nil {
		t.Fatal(err)
	}
	d.SetRuntime(&flowReferenceRuntime{e})
	if err := d.Execute(t.Context()); err == nil || !strings.Contains(err.Error(), "fetch") {
		t.Fatal("root target resolved helper output", err)
	}
}

func TestWorkflowDependencyUsesRootOrChildFrame(t *testing.T) {
	for _, nested := range []bool{false, true} {
		var d uws1.Document
		text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"variables":{"value":"fixture"},"operations":[{"operationId":"noop","x-uws-operation-profile":"fixture"},{"operationId":"consume","x-uws-operation-profile":"fixture","outputs":{"value":"$steps.fetch.outputs.value"}}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"fetch","operationRef":"noop","outputs":{"value":"$variables.value"}},{"stepId":"join","type":"merge","dependsOn":["reader"]}]},{"workflowId":"reader","type":"sequence","steps":[{"stepId":"read","operationRef":"consume"}]}]}`
		if err := json.Unmarshal([]byte(text), &d); err != nil {
			t.Fatal(err)
		}
		if nested {
			d.Workflows[0].WorkflowID = "caller"
			d.Workflows = append(d.Workflows, &uws1.Workflow{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "call", StepExecutionFields: uws1.StepExecutionFields{Workflow: "caller"}}}})
		}
		if err := d.Validate(); err != nil {
			t.Fatal("ordinary fixture invalid", err)
		}
		if err := d.ValidateExecutable(); err != nil {
			t.Fatal("executable fixture invalid", err)
		}
		r, err := AnalyzeFlow(t.Context(), &d)
		if err != nil || hasCode(r, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/value") != nested {
			t.Fatal("workflow dependency lost its invocation frame", nested, r, err)
		}
		e, err := expressions.NewEvaluator(&d)
		if err != nil {
			t.Fatal(err)
		}
		d.SetRuntime(&flowReferenceRuntime{e})
		err = d.Execute(t.Context())
		if !nested && err != nil || nested && (err == nil || !strings.Contains(err.Error(), "fetch")) {
			t.Fatal("reference execution disagrees with invocation frame", nested, err)
		}
	}
}

func TestStepReferenceDoesNotConsumeOperationOutput(t *testing.T) {
	for _, sameID := range []bool{false, true} {
		var d uws1.Document
		text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"variables":{"used":"step","unused":"operation"},"operations":[{"operationId":"produce","x-uws-operation-profile":"fixture","outputs":{"value":"$variables.unused"}},{"operationId":"consume","x-uws-operation-profile":"fixture","outputs":{"found":"$steps.fetch.outputs.value"}}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"fetch","operationRef":"produce","outputs":{"value":"$variables.used"}},{"stepId":"read","operationRef":"consume"}]}]}`
		if err := json.Unmarshal([]byte(text), &d); err != nil {
			t.Fatal(err)
		}
		if sameID {
			d.Workflows[0].Steps[0].StepID = "produce"
			d.Operations[1].Outputs["found"] = "$steps.produce.outputs.value"
		}
		if err := d.Validate(); err != nil {
			t.Fatal("ordinary fixture invalid", err)
		}
		if err := d.ValidateExecutable(); err != nil {
			t.Fatal("executable fixture invalid", err)
		}
		e, err := expressions.NewEvaluator(&d)
		if err != nil {
			t.Fatal(err)
		}
		d.SetRuntime(&flowReferenceRuntime{e})
		if err := d.Execute(t.Context()); err != nil {
			t.Fatal(err)
		}
		if d.ExecutionRecords()["stepop:read:consume"].Outputs["found"] != "step" {
			t.Fatal("consumer did not read the step output")
		}
		r, err := AnalyzeFlow(t.Context(), &d)
		if err != nil || !hasCode(r, "flow.output_unreferenced", "/operations/0/outputs/value") || hasCode(r, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/value") {
			t.Fatal("distinct step/operation output ownership lost", r, err)
		}
	}
}

func TestWorkflowControlsUseIncomingRecordsAndOutputsUseBody(t *testing.T) {
	for _, control := range []string{"forEach", "wait", "when", "items", "case", "await", "outputs"} {
		var d uws1.Document
		text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"variables":{"items":[1],"zero":0,"ready":true},"operations":[{"operationId":"noop","x-uws-operation-profile":"fixture"}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"fetch","operationRef":"noop","outputs":{"items":"$variables.items","zero":"$variables.zero","ready":"$variables.ready"}},{"stepId":"call","workflow":"helper"}]},{"workflowId":"helper","type":"sequence","steps":[{"stepId":"read","operationRef":"noop"}]}]}`
		if err := json.Unmarshal([]byte(text), &d); err != nil {
			t.Fatal(err)
		}
		helper := d.Workflows[1]
		output := "ready"
		switch control {
		case "forEach":
			helper.ForEach = "$steps.fetch.outputs.items"
			output = "items"
		case "wait":
			helper.Wait = "$steps.fetch.outputs.zero"
			output = "zero"
		case "when":
			helper.When = "$steps.fetch.outputs.ready"
		case "items":
			helper.Type = uws1.WorkflowTypeLoop
			helper.Items = "$steps.fetch.outputs.items"
			output = "items"
		case "case":
			helper.Type = uws1.WorkflowTypeSwitch
			helper.Cases = []*uws1.Case{{CaseFields: uws1.CaseFields{Name: "ready", When: "$steps.fetch.outputs.ready"}, Steps: helper.Steps}}
			helper.Steps = nil
		case "await":
			helper.Type = uws1.WorkflowTypeAwait
			helper.Wait = "$steps.fetch.outputs.ready"
		case "outputs":
			helper.Outputs = map[string]string{"ready": "$steps.fetch.outputs.ready"}
		}
		if err := d.Validate(); err != nil {
			t.Fatal("ordinary fixture invalid", control, err)
		}
		if err := d.ValidateExecutable(); err != nil {
			t.Fatal("executable fixture invalid", control, err)
		}
		r, err := AnalyzeFlow(t.Context(), &d)
		if err != nil || hasCode(r, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/"+output) != (control == "outputs") {
			t.Fatal("record snapshot scope lost", control, r, err)
		}
		e, err := expressions.NewEvaluator(&d)
		if err != nil {
			t.Fatal(err)
		}
		d.SetRuntime(&flowReferenceRuntime{e})
		err = d.Execute(t.Context())
		if control != "outputs" && err != nil || control == "outputs" && (err == nil || !strings.Contains(err.Error(), "fetch")) {
			t.Fatal("native record scope differs", control, err)
		}
	}
}

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
	text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"sourceDescriptions":[{"name":"api","type":"openapi","url":"fixture.json"}],"operations":[{"operationId":"read","sourceDescription":"api","sourceOperationId":"read","outputs":{"value":"$response.body"}},{"operationId":"send","sourceDescription":"api","sourceOperationId":"send","request":{"body":"$steps.fetch.outputs.value"}}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"fetch","operationRef":"read","outputs":{"value":"$response.body"}},{"stepId":"deliver","operationRef":"send"}]}]}`
	if err := json.Unmarshal([]byte(text), &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	r, _ := AnalyzeFlow(t.Context(), &d)
	if hasCode(r, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/value") || !hasCode(r, "flow.output_unreferenced", "/operations/0/outputs/value") {
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
