package expressions

import (
	"encoding/json"
	"github.com/OpenUdon/uws/uws1"
	"reflect"
	"testing"
)

func TestPortabilityIsOptInAndDoesNotScanProfileStrings(t *testing.T) {
	doc := &uws1.Document{UWS: "1.12.0", Info: &uws1.Info{Title: "fixture", Version: "1"}, Operations: []*uws1.Operation{{OperationID: "read", Extensions: map[string]any{uws1.ExtensionOperationProfile: "mock-fixture", "x-function-template": "expr($private)"}, OperationExecutionFields: uws1.OperationExecutionFields{When: "length($inputs) > 0"}, Request: map[string]any{"x-template": "expr($function-template)"}}}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "one", OperationRef: "read", Body: map[string]any{"profile": "expr($browser-template)"}}}}}}
	before, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := doc.Validate(); err != nil {
		t.Fatalf("ordinary validation tightened: %v", err)
	}
	got := CheckPortability(doc)
	if !reflect.DeepEqual(got, []Diagnostic{{Code: "expression.syntax", Path: "/operations/0/when"}}) {
		t.Fatalf("unexpected profile scan: %+v", got)
	}
	after, _ := json.Marshal(doc)
	if string(before) != string(after) {
		t.Fatal("checker mutated document")
	}
	doc.Operations[0].When = "expr($inputs.value)"
	got = CheckPortability(doc)
	if len(got) != 1 || got[0].Code != "expression.legacy-wrapper" {
		t.Fatal(got)
	}
}
func TestPortabilityUsesLoopInvocationScopeAndFieldKinds(t *testing.T) {
	doc := &uws1.Document{UWS: "1.12.0", Operations: []*uws1.Operation{{OperationID: "read", Outputs: map[string]string{"batch": "$batchIndex"}}}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeLoop, StructuralFields: uws1.StructuralFields{Items: "$variables.items", BatchSize: "2"}, Steps: []*uws1.Step{{StepID: "one", OperationRef: "read"}}}}}
	if got := CheckPortability(doc); len(got) != 0 {
		t.Fatalf("loop invocation rejected: %+v", got)
	}
	doc.Workflows[0].Type = uws1.WorkflowTypeSequence
	doc.Workflows[0].BatchSize = ""
	doc.Workflows[0].Items = ""
	got := CheckPortability(doc)
	if len(got) != 1 || got[0].Code != "expression.context" {
		t.Fatalf("batch outside loop: %+v", got)
	}
	doc.Operations[0].Outputs = nil
	doc.Workflows[0].Type = uws1.WorkflowTypeAwait
	doc.Workflows[0].Wait = "2"
	got = CheckPortability(doc)
	if len(got) != 1 || got[0].Code != "expression.context" {
		t.Fatalf("await numeric accepted: %+v", got)
	}
}
func TestNonSimpleConditionsAndNestedLegacyInputs(t *testing.T) {
	doc := &uws1.Document{UWS: "1.12.0", SourceDescriptions: []*uws1.SourceDescription{{Name: "api", Type: "openapi"}}, Operations: []*uws1.Operation{{OperationID: "read", SourceDescription: "api", SourceOperationID: "read", Request: map[string]any{"query": map[string]any{"params": []any{"expr($private)", "ordinary literal"}}}, SuccessCriteria: []*uws1.Criterion{{Type: uws1.CriterionRegex, Context: "$response.body", Condition: "^$profile-template"}, {Type: uws1.CriterionJSONPath, Context: "$response.body", Condition: "$.items[?(@.secret)]"}}}}}
	got := CheckPortability(doc)
	if !reflect.DeepEqual(got, []Diagnostic{{Code: "expression.legacy-wrapper", Path: "/operations/0/request/query/params/0"}}) {
		t.Fatalf("noncore query scan: %+v", got)
	}
}
func TestPortabilityBoundsCyclesAndDeterministicOrder(t *testing.T) {
	step := &uws1.Step{StepID: "cycle"}
	step.Steps = []*uws1.Step{step}
	doc := &uws1.Document{UWS: "1.12.0", Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{step}}}, Operations: []*uws1.Operation{{Outputs: map[string]string{"b": "expr($private)", "a": "expr($private)"}}}}
	a := CheckPortability(doc)
	b := CheckPortability(doc)
	if !reflect.DeepEqual(a, b) || len(a) < 3 || a[0].Path != "/operations/0/outputs/a" {
		t.Fatalf("unstable/unbounded diagnostics: %+v", a)
	}
}

func TestEmptyOutputIsNotAValidCoreExpression(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", Operations: []*uws1.Operation{{Outputs: map[string]string{"missing": ""}}}}
	got := CheckPortability(d)
	if len(got) != 1 || got[0].Code != "expression.syntax" {
		t.Fatal(got)
	}
}

func TestTypedNestedCoreValuesDoNotBypassPortability(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", SourceDescriptions: []*uws1.SourceDescription{{Name: "api", Type: "openapi"}}, Operations: []*uws1.Operation{{SourceDescription: "api", SourceOperationID: "read", Request: map[string]any{"query": map[string]string{"n": "expr($private)"}}}}}
	got := CheckPortability(d)
	if len(got) != 1 || got[0].Code != "expression.legacy-wrapper" {
		t.Fatal(got)
	}
}

func TestRequestExtensionsRemainOwnedWhilePayloadBindingsAreChecked(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", SourceDescriptions: []*uws1.SourceDescription{{Name: "api", Type: "openapi"}}, Operations: []*uws1.Operation{{SourceDescription: "api", SourceOperationID: "read", Request: map[string]any{"x-profile-template": "expr($private)", "body": map[string]any{"x-user-field": "expr($private)"}}}}}
	got := CheckPortability(d)
	want := []Diagnostic{{Code: "expression.legacy-wrapper", Path: "/operations/0/request/body/x-user-field"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("extension ownership/payload coverage: %+v", got)
	}
}

func TestStructuralResultExpressionsAreCheckedWithoutInspectingExtensions(t *testing.T) {
	d := &uws1.Document{UWS: "1.12.0", Results: []*uws1.StructuralResult{{Name: "result", Kind: "merge", From: "main", Value: "expr($private)", Extensions: map[string]any{"x-profile": "expr($profile)"}}}}
	got := CheckPortability(d)
	if !reflect.DeepEqual(got, []Diagnostic{{Code: "expression.legacy-wrapper", Path: "/results/0/value"}}) {
		t.Fatal(got)
	}
}
