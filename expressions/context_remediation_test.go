package expressions

import (
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
