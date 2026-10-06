package expressions

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestEvaluateExpressionSourcesAndComparisons(t *testing.T) {
	document := &uws1.Document{
		UWS: "1.12.0",
		Variables: map[string]any{
			"name":        "Ada",
			"items":       []any{"first", "second"},
			"typed_items": []string{"go", "lang"},
			"typed_map":   map[string]string{"name": "typed"},
		},
		Components: &uws1.Components{Variables: map[string]any{"name": "shadowed"}},
	}
	runtime, err := newEvaluatorFixture(document, struct{}{})
	if err != nil {
		t.Fatal(err)
	}
	state := &uws1.ExecutionContext{
		Iteration: &uws1.IterationContext{Item: map[string]any{"id": "item-2"}, Index: 1, BatchIndex: 0},
		Trigger:   &uws1.TriggerExecutionContext{Payload: map[string]any{"event": map[string]any{"id": "evt-4"}}},
		Inputs: map[string]any{
			"name":    "$variables.name",
			"payload": map[string]any{"city": "Toronto"},
			"typed":   map[string]string{"name": "$variables.name"},
		},
		Records: map[string]uws1.ExecutionRecord{
			"step:fetch": {ID: "fetch", Kind: "step:operation", Status: "success", Outputs: map[string]any{"item": map[string]any{"name": "widget", "tags": []any{"a", "b"}}}},
			"op:current": {ID: "current", Kind: "operation", Status: "running", Result: json.RawMessage(`{"statusCode":200,"headers":{"Content-Type":"application/json"},"body":{"items":[{"name":"widget"}]}}`)},
		},
		Current: &uws1.CurrentExecutionContext{Key: "op:current", Outputs: map[string]any{"partial": map[string]any{"value": 7}}},
	}
	ctx := uws1.WithExecutionContext(context.Background(), state)
	for expression, want := range map[string]any{
		"$response.statusCode":                       json.Number("200"),
		"$response.headers.content-type":             "application/json",
		"$response.body.items.0.name":                "widget",
		"$response.body.items.00.name":               "widget",
		"$response.body#/items/0/name":               "widget",
		"$response.body#/items/00/name":              nil,
		"$response.body#/items/%2B0/name":            nil,
		"$outputs.partial.value":                     7,
		"$steps.fetch.outputs.item.tags.1":           "b",
		"$variables.name":                            "Ada",
		"$variables.typed_items.00":                  "go",
		"$variables.typed_map.name":                  "typed",
		"$inputs.name":                               "Ada",
		"$inputs.payload.city":                       "Toronto",
		"$inputs.typed.name":                         "Ada",
		"$trigger.event.id":                          "evt-4",
		"$item.id":                                   "item-2",
		"$index":                                     1,
		"$batchIndex":                                0,
		"$variables.name == \"Ada\"":                 true,
		"$response.statusCode >= 200":                true,
		"$response.body#/items/0/name == \"widget\"": true,
		"1.25": json.Number("1.25"),
	} {
		t.Run(expression, func(t *testing.T) {
			got, err := runtime.EvaluateExpression(ctx, expression)
			if err != nil {
				t.Fatalf("EvaluateExpression(%q): %v", expression, err)
			}
			if !jsonValuesEqual(got, want) {
				t.Fatalf("EvaluateExpression(%q) = %#v, want %#v", expression, got, want)
			}
		})
	}
	items, err := runtime.ResolveItems(ctx, "$variables.items")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0] != "first" || items[1] != "second" {
		t.Fatalf("ResolveItems = %#v, want source order", items)
	}
	items, err = runtime.ResolveItems(ctx, "$variables.typed_items")
	if err != nil || len(items) != 2 || items[0] != "go" || items[1] != "lang" {
		t.Fatalf("ResolveItems on typed Go slice = %#v, %v", items, err)
	}
}

func TestEvaluateExpressionRejectsUnsupportedOrInvalidForms(t *testing.T) {
	document := &uws1.Document{UWS: "1.12.0"}
	runtime, err := newEvaluatorFixture(document, struct{}{})
	if err != nil {
		t.Fatal(err)
	}
	state := &uws1.ExecutionContext{
		Iteration: &uws1.IterationContext{Item: "value", Index: 0, BatchIndex: -1},
		Records:   map[string]uws1.ExecutionRecord{"op:current": {ID: "current", Kind: "operation", Result: json.RawMessage(`{"body":{"value":1}}`)}},
		Current:   &uws1.CurrentExecutionContext{Key: "op:current"},
	}
	ctx := uws1.WithExecutionContext(context.Background(), state)
	for _, expression := range []string{
		"length($inputs)",
		"$response.body#/value~2bad",
		"$response.body#/bad path",
		"$response.headers.content.type",
		"$batchIndex",
		"$response.statusCode < true",
		"$response.statusCode == \"200\"",
		"$inputs..name",
	} {
		t.Run(expression, func(t *testing.T) {
			if _, err := runtime.EvaluateExpression(ctx, expression); err == nil {
				t.Fatalf("unsupported expression %q was accepted", expression)
			}
		})
	}
	if _, err := runtime.EvaluateExpression(ctx, "$variables."+strings.Repeat("a", MaxBytes)); err == nil {
		t.Fatal("oversized expression was accepted")
	}
	legacyRuntime, err := newEvaluatorFixture(&uws1.Document{UWS: "1.10.0"}, struct{}{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := legacyRuntime.EvaluateExpression(ctx, "$response.body.value"); err == nil {
		t.Fatal("UWS 1.11 response dot-walk was accepted for UWS 1.10")
	}
	if _, err := runtime.EvaluateExpression(context.Background(), "$item"); err == nil || !strings.Contains(err.Error(), "iteration") {
		t.Fatalf("$item without iteration error = %v", err)
	}
}

func jsonValuesEqual(first, second any) bool {
	firstJSON, err := json.Marshal(first)
	if err != nil {
		return false
	}
	secondJSON, err := json.Marshal(second)
	return err == nil && string(firstJSON) == string(secondJSON)
}

type evaluatorFixture struct{ *Evaluator }

func newEvaluatorFixture(document *uws1.Document, _ struct{}) (*evaluatorFixture, error) {
	e, err := NewEvaluator(document)
	if err != nil {
		return nil, err
	}
	return &evaluatorFixture{e}, nil
}
func (e *evaluatorFixture) EvaluateExpression(ctx context.Context, text string) (any, error) {
	field := Value
	if !strings.HasPrefix(text, "$") {
		field = Wait
	}
	return e.Evaluate(ctx, text, field)
}
