package expressions

import (
	"encoding/json"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestStrictItemIndexContextDoesNotChangeParse(t *testing.T) {
	for _, text := range []string{"$item", "$item.name", "$index"} {
		if _, err := Parse(text, Context{Version: "1.13.0", Field: Value}); err != nil {
			t.Fatal("ordinary parse changed", err)
		}
		d := &uws1.Document{UWS: "1.13.0", Operations: []*uws1.Operation{{OperationID: "read", Outputs: map[string]string{"value": text}}}}
		if got := CheckPortability(d); len(got) != 1 || got[0].Code != "expression.context" {
			t.Fatalf("outside iteration %s: %+v", text, got)
		}
		d.Operations[0].ForEach = "$variables.items"
		if got := CheckPortability(d); len(got) != 0 {
			t.Fatalf("forEach iteration %s: %+v", text, got)
		}
	}
}

func TestDependenciesUseIncomingIterationContext(t *testing.T) {
	var d uws1.Document
	text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"operations":[{"operationId":"read","x-uws-operation-profile":"fixture","outputs":{"value":"$item","batch":"$batchIndex"}},{"operationId":"noop","x-uws-operation-profile":"fixture"}],"workflows":[{"workflowId":"main","type":"loop","items":"$variables.items","steps":[{"stepId":"dependency","operationRef":"read"},{"stepId":"entry","operationRef":"noop","dependsOn":["dependency"]}]}],"triggers":[{"triggerId":"event","outputs":["ready"],"routes":[{"output":"ready","to":["entry"]}]}]}`
	if err := json.Unmarshal([]byte(text), &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	if got := CheckPortability(&d); len(got) != 2 {
		t.Fatalf("transitive direct trigger context: %+v", got)
	}
	// A dependency runs before its depender creates a forEach iteration.
	d.Triggers = nil
	d.Workflows[0].Type = uws1.WorkflowTypeSequence
	d.Workflows[0].Items = ""
	d.Workflows[0].Steps = []*uws1.Step{{StepID: "each", OperationRef: "noop"}}
	d.Operations[1].ForEach = "$variables.items"
	d.Operations[1].DependsOn = []string{"read"}
	if got := CheckPortability(&d); len(got) != 2 {
		t.Fatalf("future iteration inherited by dependency: %+v", got)
	}
	// A workflow dependency inside a structural loop inherits that loop.
	d.Operations[1].ForEach = ""
	d.Operations[1].DependsOn = nil
	d.Operations[0].Outputs = nil
	d.Workflows[0].Type = uws1.WorkflowTypeLoop
	d.Workflows[0].Items = "$variables.items"
	d.Workflows[0].Steps = []*uws1.Step{{StepID: "join", Type: uws1.WorkflowTypeMerge, StepExecutionFields: uws1.StepExecutionFields{DependsOn: []string{"helper"}}}}
	d.Workflows = append(d.Workflows, &uws1.Workflow{WorkflowID: "helper", Type: uws1.WorkflowTypeSequence, WorkflowExecutionFields: uws1.WorkflowExecutionFields{When: "$index == 0"}})
	if got := CheckPortability(&d); len(got) != 0 {
		t.Fatalf("valid loop dependency context lost: %+v", got)
	}
}

func TestForEachControlsRunBeforeNewIteration(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Operations: []*uws1.Operation{{OperationID: "read", OperationExecutionFields: uws1.OperationExecutionFields{ForEach: "$variables.items", When: "$index == 0", Wait: "$index"}}}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, WorkflowExecutionFields: uws1.WorkflowExecutionFields{ForEach: "$variables.items", When: "$item"}, Steps: []*uws1.Step{{StepID: "read", OperationRef: "read"}}}}}
	// Workflow body iteration is available to its operation, but the workflow's
	// own when still runs outside that newly-created iteration.
	got := CheckPortability(d)
	if len(got) != 1 || got[0].Path != "/workflows/0/when" {
		t.Fatal(got)
	}
	d.Workflows[0].ForEach = ""
	d.Workflows[0].When = ""
	got = CheckPortability(d)
	if len(got) != 2 || got[0].Path != "/operations/0/when" || got[1].Path != "/operations/0/wait" {
		t.Fatal(got)
	}
}

func TestGotoTargetsUseRootIterationContext(t *testing.T) {
	var d uws1.Document
	text := `{"uws":"1.13.0","info":{"title":"fixture","version":"1"},"operations":[{"operationId":"read","x-uws-operation-profile":"fixture","outputs":{"value":"$item","batch":"$batchIndex"}},{"operationId":"jump","x-uws-operation-profile":"fixture","onSuccess":[{"name":"transfer","type":"goto","stepId":"target"}]}],"workflows":[{"workflowId":"main","type":"loop","items":"$variables.items","steps":[{"stepId":"target","operationRef":"read"},{"stepId":"transfer","operationRef":"jump"}]}]}`
	if err := json.Unmarshal([]byte(text), &d); err != nil {
		t.Fatal(err)
	}
	if err := d.Validate(); err != nil {
		t.Fatal("ordinary fixture invalid", err)
	}
	if got := CheckPortability(&d); len(got) != 2 {
		t.Fatalf("root goto retained loop context: %+v", got)
	}
	d.Operations[1].OnSuccess = nil
	d.Operations[1].OnFailure = []*uws1.FailureAction{{Name: "transfer", Type: "goto", StepID: "target"}}
	if got := CheckPortability(&d); len(got) != 2 {
		t.Fatalf("failure goto retained loop context: %+v", got)
	}
}

func TestTriggerWorkflowLoopAndDirectStepContext(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Operations: []*uws1.Operation{{OperationID: "read", Outputs: map[string]string{"value": "$item", "batch": "$batchIndex"}}}, Workflows: []*uws1.Workflow{
		{WorkflowID: "main", Type: uws1.WorkflowTypeSequence},
		{WorkflowID: "event", Type: uws1.WorkflowTypeLoop, StructuralFields: uws1.StructuralFields{Items: "$variables.items"}, Steps: []*uws1.Step{{StepID: "read-step", OperationRef: "read"}}},
	}, Triggers: []*uws1.Trigger{{Routes: []*uws1.TriggerRoute{{TriggerRouteFields: uws1.TriggerRouteFields{To: []string{"event"}}}}}}}
	if got := CheckPortability(d); len(got) != 0 {
		t.Fatalf("trigger loop context lost: %+v", got)
	}
	// Directly routed main steps bypass the main workflow's loop body.
	d.Workflows = []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeLoop, StructuralFields: uws1.StructuralFields{Items: "$variables.items"}, Steps: []*uws1.Step{{StepID: "read-step", OperationRef: "read", Inputs: map[string]any{"value": "$index"}}}}}
	d.Triggers[0].Routes[0].To = []string{"read-step"}
	got := CheckPortability(d)
	if len(got) != 3 {
		t.Fatalf("mixed direct/loop contexts accepted: %+v", got)
	}
	for _, diagnostic := range got {
		if diagnostic.Code != "expression.context" {
			t.Fatal(got)
		}
	}
}

func TestForEachCollectionUsesOuterIterationOnly(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "each", StepExecutionFields: uws1.StepExecutionFields{ForEach: "$item.children"}, Outputs: map[string]string{"value": "$item"}}}}}}
	got := CheckPortability(d)
	if len(got) != 1 || got[0].Path != "/workflows/0/steps/0/forEach" {
		t.Fatal(got)
	}
}
