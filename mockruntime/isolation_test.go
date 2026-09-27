package mockruntime

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestMockRuntimeIsolatesParallelLeafResponses(t *testing.T) {
	document := mockRuntimeDocument()
	operation := document.Operations[0]
	operation.Request = map[string]any{"query": map[string]any{"name": "$inputs.name"}}
	operation.SuccessCriteria = []*uws1.Criterion{{Condition: `$response.body.name == $inputs.name`}}
	operation.Outputs = nil
	document.Workflows[0].Type = uws1.WorkflowTypeParallel
	document.Workflows[0].Steps = []*uws1.Step{
		{StepID: "alice_step", OperationRef: operation.OperationID, Inputs: map[string]any{"name": "Alice"}, Outputs: map[string]string{"name": "$response.body.name"}},
		{StepID: "bob_step", OperationRef: operation.OperationID, Inputs: map[string]any{"name": "Bob"}, Outputs: map[string]string{"name": "$response.body.name"}},
	}
	fixtures := fixtureSetForNames(t, operation.OperationID, "Alice", "Bob")
	runtime, err := NewRuntime(document, Options{Fixtures: fixtures})
	if err != nil {
		t.Fatal(err)
	}
	document.SetRuntime(runtime)
	if err := document.Execute(context.Background()); err != nil {
		t.Fatalf("parallel Execute() error = %v", err)
	}
	records := document.ExecutionRecords()
	for stepID, want := range map[string]string{"alice_step": "Alice", "bob_step": "Bob"} {
		if got := records["step:"+stepID].Outputs["name"]; got != want {
			t.Errorf("step %s response = %#v, want %q", stepID, got, want)
		}
	}
	if got := len(runtime.RequestRecords()); got != 2 {
		t.Fatalf("parallel request count = %d, want 2", got)
	}
}

func TestMockRuntimeIsolatesForEachIterationResponses(t *testing.T) {
	document := mockRuntimeDocument()
	operation := document.Operations[0]
	operation.Request = map[string]any{"query": map[string]any{"name": "$item.name"}}
	operation.SuccessCriteria = []*uws1.Criterion{{Condition: `$response.body.name == $item.name`}}
	operation.Outputs = nil
	document.Variables = map[string]any{"people": []any{map[string]any{"name": "Alice"}, map[string]any{"name": "Bob"}}}
	document.Workflows[0].Steps = []*uws1.Step{{
		StepID: "people_step",
		StepExecutionFields: uws1.StepExecutionFields{
			ForEach: "$variables.people",
		},
		OperationRef: operation.OperationID,
		Outputs:      map[string]string{"name": "$response.body.name"},
	}}
	runtime, err := NewRuntime(document, Options{Fixtures: fixtureSetForNames(t, operation.OperationID, "Alice", "Bob")})
	if err != nil {
		t.Fatal(err)
	}
	document.SetRuntime(runtime)
	if err := document.Execute(context.Background()); err != nil {
		t.Fatalf("forEach Execute() error = %v", err)
	}
	got := document.ExecutionRecords()["step:people_step"].Outputs["name"]
	if !reflect.DeepEqual(got, []any{"Alice", "Bob"}) {
		t.Fatalf("forEach response outputs = %#v, want ordered Alice/Bob", got)
	}
	if got := len(runtime.RequestRecords()); got != 2 {
		t.Fatalf("forEach request count = %d, want 2", got)
	}
}

func TestMockRuntimeForEachReadsCompletedOuterStep(t *testing.T) {
	document := mockRuntimeDocument()
	config := document.Operations[0]
	config.SuccessCriteria = nil
	config.Outputs = nil
	document.Operations = append(document.Operations, &uws1.Operation{
		OperationID: "consume", Extensions: config.Extensions,
		Request: map[string]any{"body": map[string]any{"token": "$steps.config.outputs.token"}},
	})
	document.Variables = map[string]any{"items": []any{"one", "two"}}
	document.Workflows[0].Steps = []*uws1.Step{
		{StepID: "config", OperationRef: config.OperationID, Outputs: map[string]string{"token": "$response.body.token"}},
		{StepID: "consume_step", OperationRef: "consume", StepExecutionFields: uws1.StepExecutionFields{ForEach: "$variables.items"}},
	}
	runtime, err := NewRuntime(document, Options{ResponseResolver: ResponseResolverFunc(func(context.Context, *uws1.Operation) (ResponseDefinition, error) {
		return ResponseDefinition{Example: json.RawMessage(`{"statusCode":200,"body":{"token":"ready"}}`)}, nil
	})})
	if err != nil {
		t.Fatal(err)
	}
	document.SetRuntime(runtime)
	if err := document.Execute(context.Background()); err != nil {
		t.Fatalf("forEach could not read completed outer step: %v", err)
	}
	requests := runtime.RequestRecords()
	if len(requests) != 3 {
		t.Fatalf("request count = %d, want config plus two consumers", len(requests))
	}
	for _, request := range requests[1:] {
		if string(request.Request) != `{"body":{"token":"ready"}}` {
			t.Fatalf("consumer request = %s", request.Request)
		}
	}
}

func TestMockRuntimeStepLookupPrefersNearestVisibleIteration(t *testing.T) {
	runtime, err := NewRuntime(mockRuntimeDocument(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	record := func(value string) uws1.ExecutionRecord {
		return uws1.ExecutionRecord{ID: "config", Kind: "step:operation", Status: "success", Outputs: map[string]any{"token": value}}
	}
	state := &uws1.ExecutionContext{
		Current: &uws1.CurrentExecutionContext{Key: "stepop:consume:consume#iter:1.2"},
		Records: map[string]uws1.ExecutionRecord{
			"step:config":                      record("root"),
			"step:config#iter:1":               record("parent"),
			"step:config#iter:1.2":             record("current"),
			"step:config#iter:1.3":             record("sibling"),
			"call:other::step:config#iter:1.2": record("other invocation"),
		},
	}
	lookup := func(want string) {
		t.Helper()
		value, err := runtime.EvaluateExpression(uws1.WithExecutionContext(context.Background(), state), "$steps.config.outputs.token")
		if err != nil || value != want {
			t.Fatalf("lookup = %#v, %v; want %q", value, err, want)
		}
	}
	lookup("current")
	delete(state.Records, "step:config#iter:1.2")
	lookup("parent")
	delete(state.Records, "step:config#iter:1")
	lookup("root")
	delete(state.Records, "step:config")
	if _, err := runtime.EvaluateExpression(uws1.WithExecutionContext(context.Background(), state), "$steps.config.outputs.token"); err == nil || !strings.Contains(err.Error(), "no execution record") {
		t.Fatalf("sibling or other invocation was visible: %v", err)
	}
	state.Records["step:config#iter:1.2"] = record("one")
	state.Records["step:alias#iter:1.2"] = record("two")
	if _, err := runtime.EvaluateExpression(uws1.WithExecutionContext(context.Background(), state), "$steps.config.outputs.token"); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("ambiguous lookup did not fail: %v", err)
	}
	state.WorkflowScope = "call:fetch:step:call#iter:1"
	state.Current.Key = state.WorkflowScope + "::stepop:consume:consume#iter:1.2"
	state.Records = map[string]uws1.ExecutionRecord{
		"step:config#iter:1.2": record("local invocation"),
	}
	lookup("local invocation")
}

func TestMockRuntimeIsolatesWorkflowCallInputsAndResponses(t *testing.T) {
	document := mockRuntimeDocument()
	operation := document.Operations[0]
	operation.Request = map[string]any{"query": map[string]any{"name": "$inputs.name"}}
	operation.SuccessCriteria = []*uws1.Criterion{{Condition: `$response.body.name == $inputs.name`}}
	operation.Outputs = nil
	document.Workflows = []*uws1.Workflow{
		{
			WorkflowID: "main",
			Type:       uws1.WorkflowTypeSequence,
			Steps: []*uws1.Step{
				{StepID: "call_alice", StepExecutionFields: uws1.StepExecutionFields{Workflow: "fetch"}, Inputs: map[string]any{"name": "Alice"}},
				{StepID: "call_bob", StepExecutionFields: uws1.StepExecutionFields{Workflow: "fetch"}, Inputs: map[string]any{"name": "Bob"}},
			},
		},
		{
			WorkflowID: "fetch",
			Type:       uws1.WorkflowTypeSequence,
			Steps:      []*uws1.Step{{StepID: "fetch_step", OperationRef: operation.OperationID, Outputs: map[string]string{"name": "$response.body.name"}}},
			Outputs:    map[string]string{"name": "$steps.fetch_step.outputs.name"},
		},
	}
	runtime, err := NewRuntime(document, Options{Fixtures: fixtureSetForNames(t, operation.OperationID, "Alice", "Bob")})
	if err != nil {
		t.Fatal(err)
	}
	document.SetRuntime(runtime)
	if err := document.Execute(context.Background()); err != nil {
		t.Fatalf("workflow-call Execute() error = %v", err)
	}
	var names []string
	for _, record := range document.ExecutionRecords() {
		if record.ID == "fetch" && record.Kind == "workflow:sequence" {
			if name, ok := record.Outputs["name"].(string); ok {
				names = append(names, name)
			}
		}
	}
	sort.Strings(names)
	if !reflect.DeepEqual(names, []string{"Alice", "Bob"}) {
		t.Fatalf("workflow call outputs = %#v, want isolated Alice/Bob", names)
	}
}

func TestMockRuntimeRepeatsExactFixtureForOperationRetry(t *testing.T) {
	document := mockRuntimeDocument()
	operation := document.Operations[0]
	operation.SuccessCriteria = []*uws1.Criterion{{Condition: "$response.statusCode == 200"}}
	operation.OnFailure = []*uws1.FailureAction{{Name: "retry_once", Type: "retry", RetryLimit: 1}}
	digest, err := RequestDigest(operation.Request)
	if err != nil {
		t.Fatal(err)
	}
	runtime, err := NewRuntime(document, Options{Fixtures: &FixtureSet{Format: FixtureFormatV1, Fixtures: []Fixture{{
		OperationID:   operation.OperationID,
		RequestDigest: digest,
		Provenance:    FixtureProvenance{Kind: "example"},
		Response:      json.RawMessage(`{"statusCode":503,"body":{"retryable":true}}`),
	}}}})
	if err != nil {
		t.Fatal(err)
	}
	document.SetRuntime(runtime)
	if err := document.Execute(context.Background()); err == nil {
		t.Fatal("retry execution with a 503 fixture unexpectedly succeeded")
	}
	requests := runtime.RequestRecords()
	if len(requests) != 2 || requests[0].RequestDigest != digest || requests[1].RequestDigest != digest || requests[0].ResponseKind != ResponseFixture || requests[1].ResponseKind != ResponseFixture {
		t.Fatalf("retry request records = %#v, want two stable fixture replays", requests)
	}
}

func fixtureSetForNames(t *testing.T, operationID string, names ...string) *FixtureSet {
	t.Helper()
	fixtures := &FixtureSet{Format: FixtureFormatV1}
	for _, name := range names {
		digest, err := RequestDigest(map[string]any{"query": map[string]any{"name": name}})
		if err != nil {
			t.Fatal(err)
		}
		fixtures.Fixtures = append(fixtures.Fixtures, Fixture{
			OperationID:   operationID,
			RequestDigest: digest,
			Provenance:    FixtureProvenance{Kind: "example"},
			Response:      json.RawMessage(`{"statusCode":200,"body":{"name":"` + name + `"}}`),
		})
	}
	return fixtures
}
