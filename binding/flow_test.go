package binding

import (
	"context"
	"encoding/json"
	"github.com/OpenUdon/uws/uws1"
	"reflect"
	"strings"
	"testing"
)

func hasCode(r FlowReport, code, path string) bool {
	for _, f := range r.Findings {
		if f.Code == code && (path == "" || f.Path == path) {
			return true
		}
	}
	return false
}
func TestFlowReachabilityPendingEffectsAndPrivacy(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", Variables: map[string]any{"private": "private-canary"}, Operations: []*uws1.Operation{{OperationID: "write", Effect: uws1.OperationEffectWrite, Outputs: map[string]string{"unused": "$response.body"}}, {OperationID: "unreachable", Effect: uws1.OperationEffectRead}}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "pending", Pending: &uws1.PendingStep{}}, {StepID: "send", OperationRef: "write"}}}}}
	before, _ := json.Marshal(d)
	r, err := AnalyzeFlow(t.Context(), d)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ code, path string }{{"flow.pending", "/workflows/0/steps/0/pending"}, {"flow.pending_before_effect", "/workflows/0/steps/1"}, {"flow.unreachable", "/operations/1"}, {"flow.output_unreferenced", "/operations/0/outputs/unused"}} {
		if !hasCode(r, c.code, c.path) {
			t.Fatalf("missing %+v: %+v", c, r)
		}
	}
	after, _ := json.Marshal(d)
	if string(before) != string(after) {
		t.Fatal("analysis mutated document")
	}
	b, _ := json.Marshal(r)
	if strings.Contains(string(b), "private-canary") {
		t.Fatal("value leaked")
	}
	r2, _ := AnalyzeFlow(t.Context(), d)
	if !reflect.DeepEqual(r, r2) {
		t.Fatal("nondeterministic report")
	}
}
func TestFlowReferencesKeepUsedOutputsAndUnknownEffects(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", Operations: []*uws1.Operation{{OperationID: "read", Outputs: map[string]string{"value": "$response.body"}}, {OperationID: "write", Effect: uws1.OperationEffectWrite}}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "fetch", OperationRef: "read"}, {StepID: "consume", OperationRef: "write", Inputs: map[string]any{"n": "$steps.fetch.outputs.value"}}}}}}
	r, err := AnalyzeFlow(t.Context(), d)
	if err != nil {
		t.Fatal(err)
	}
	if hasCode(r, "flow.output_unreferenced", "/operations/0/outputs/value") {
		t.Fatal("used output classified unused", r)
	}
	if !hasCode(r, "flow.effect_unknown", "/workflows/0/steps/0") {
		t.Fatal("unknown effect inferred", r)
	}
}
func TestFlowCyclesUnboundedWorkAndIgnoredMergeChildren(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "recursive", StepExecutionFields: uws1.StepExecutionFields{Workflow: "main"}}, {StepID: "wait", Type: uws1.WorkflowTypeAwait, StepExecutionFields: uws1.StepExecutionFields{Wait: "$inputs.ready"}}, {StepID: "loop", Type: uws1.WorkflowTypeLoop, StructuralFields: uws1.StructuralFields{Items: "$inputs.items"}}, {StepID: "merge", Type: uws1.WorkflowTypeMerge, Steps: []*uws1.Step{{StepID: "ignored"}}}}}}}
	r, err := AnalyzeFlow(t.Context(), d)
	if err != nil {
		t.Fatal(err)
	}
	if !hasCode(r, "flow.cycle", "") || !hasCode(r, "flow.loop_unbounded", "/workflows/0/steps/1/wait") || !hasCode(r, "flow.loop_unbounded", "/workflows/0/steps/2/items") || !hasCode(r, "flow.unreachable", "/workflows/0/steps/3/steps/0") {
		t.Fatal(r)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := AnalyzeFlow(ctx, d); err != context.Canceled {
		t.Fatal(err)
	}
}
func TestStaticLiteralLoopAndBranchPendingOrder(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", Variables: map[string]any{"items": []any{1, 2}}, Operations: []*uws1.Operation{{OperationID: "write", Effect: uws1.OperationEffectWrite}}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeLoop, StructuralFields: uws1.StructuralFields{Items: "$variables.items"}, Steps: []*uws1.Step{{StepID: "switch", Type: uws1.WorkflowTypeSwitch, Cases: []*uws1.Case{{Steps: []*uws1.Step{{StepID: "pending", Pending: &uws1.PendingStep{}}, {StepID: "write", OperationRef: "write"}}}}}}}}}
	r, err := AnalyzeFlow(t.Context(), d)
	if err != nil {
		t.Fatal(err)
	}
	if hasCode(r, "flow.loop_unbounded", "/workflows/0/items") || !hasCode(r, "flow.pending_before_effect", "/workflows/0/steps/0/cases/0/steps/1") {
		t.Fatal(r)
	}
}
