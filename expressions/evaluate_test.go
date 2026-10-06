package expressions

import (
	"context"
	"encoding/json"
	"github.com/OpenUdon/uws/uws1"
	"math"
	"strings"
	"testing"
)

func TestExactNumbersComponentsAndNullPropagation(t *testing.T) {
	e, err := NewEvaluator(&uws1.Document{UWS: "1.12.0", Variables: map[string]any{"n": json.Number("900719925474099312345"), "precise": json.Number("1.2300"), "shared": "top", "object": map[string]any{"present": nil}}, Components: &uws1.Components{Variables: map[string]any{"shared": "component", "only": "available"}}})
	if err != nil {
		t.Fatal(err)
	}
	for expression, want := range map[string]any{
		"$variables.n == 900719925474099312345": true,
		"$variables.n > 900719925474099312344":  true,
		"$variables.precise == 1.23":            true,
		"$variables.shared":                     "top", "$variables.only": "available",
		"$variables.object.missing": nil, "$variables.object.present": nil,
	} {
		got, err := e.Evaluate(context.Background(), expression, Value)
		if err != nil || got != want {
			t.Fatalf("%s: %v %v", expression, got, err)
		}
	}
	for _, expression := range []string{"$variables.precise < null", "$variables.n == \"private-canary\"", "$variables.n < 1e10001"} {
		_, err := e.Evaluate(context.Background(), expression, Value)
		if err == nil || strings.Contains(err.Error(), "private-canary") {
			t.Fatalf("invalid/unsafe comparison error %v", err)
		}
	}
}
func TestReferenceWorkflowAndNearestIterationVisibility(t *testing.T) {
	e, _ := NewEvaluator(&uws1.Document{UWS: "1.12.0"})
	rec := func(v string) uws1.ExecutionRecord {
		return uws1.ExecutionRecord{ID: "fetch", Kind: "step:operation", Status: "success", Outputs: map[string]any{"result": v}}
	}
	state := &uws1.ExecutionContext{Current: &uws1.CurrentExecutionContext{Key: "step:consume#iter:1.2"}, Records: map[string]uws1.ExecutionRecord{"step:fetch": rec("root"), "step:fetch#iter:1": rec("parent"), "step:fetch#iter:2": rec("sibling"), "other::step:fetch#iter:1.2": rec("foreign")}}
	ctx := uws1.WithExecutionContext(context.Background(), state)
	got, err := e.Evaluate(ctx, "$steps.fetch.outputs.result", Value)
	if err != nil || got != "parent" {
		t.Fatalf("nearest parent: %v %v", got, err)
	}
	delete(state.Records, "step:fetch#iter:1")
	got, err = e.Evaluate(ctx, "$steps.fetch.outputs.result", Value)
	if err != nil || got != "root" {
		t.Fatalf("root: %v %v", got, err)
	}
	delete(state.Records, "step:fetch")
	state.WorkflowScope = "current-call"
	if _, err = e.Evaluate(ctx, "$steps.fetch.outputs.result", Value); err == nil {
		t.Fatal("sibling/foreign invocation became visible")
	}
	state.Current.Key = "foreign-call::op:response"
	state.Records["foreign-call::op:response"] = uws1.ExecutionRecord{Result: map[string]any{"body": "foreign"}}
	if _, err = e.Evaluate(ctx, "$response.body", Value); err == nil {
		t.Fatal("foreign current response became visible")
	}
	state.Current.Key = "current-call::op:response"
	state.Records["op:response"] = uws1.ExecutionRecord{Result: map[string]any{"body": "local"}}
	got, err = e.Evaluate(ctx, "$response.body", Value)
	if err != nil || got != "local" {
		t.Fatalf("local response: %v %v", got, err)
	}
}
func TestReferenceCancellationCyclesAndJSONValueRefusal(t *testing.T) {
	e, _ := NewEvaluator(&uws1.Document{UWS: "1.12.0", Variables: map[string]any{"bad": math.NaN(), "hex": json.Number("0x10")}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := e.Evaluate(ctx, "$inputs", Value); err != context.Canceled {
		t.Fatal(err)
	}
	if _, err := e.Evaluate(nil, "$inputs", Value); err == nil {
		t.Fatal("nil context")
	}
	if _, err := e.Evaluate(context.Background(), "$variables.bad", Value); err == nil {
		t.Fatal("nonfinite result")
	}
	if _, err := e.Evaluate(context.Background(), "$variables.hex == 16", Value); err == nil {
		t.Fatal("non-JSON number accepted through a comparison")
	}
	state := &uws1.ExecutionContext{Inputs: map[string]any{"cycle": "$inputs.cycle"}}
	if _, err := e.Evaluate(uws1.WithExecutionContext(context.Background(), state), "$inputs.cycle", Value); err == nil {
		t.Fatal("cyclic input")
	}
}
