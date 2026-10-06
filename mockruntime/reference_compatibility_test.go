package mockruntime

import (
	"context"
	"encoding/json"
	"github.com/OpenUdon/uws/uws1"
	"testing"
)

func TestReferenceAdoptionPreservesHistoricalMockExtensions(t *testing.T) {
	doc := &uws1.Document{UWS: "1.10.0", Variables: map[string]any{"literal": "$response.body#%2Fx"}}
	r, err := NewRuntime(doc, Options{})
	if err != nil {
		t.Fatal(err)
	}
	state := &uws1.ExecutionContext{Current: &uws1.CurrentExecutionContext{Key: "op:read"}, Records: map[string]uws1.ExecutionRecord{"op:read": {Result: map[string]any{"body": map[string]any{"x": "value"}}}}}
	ctx := uws1.WithExecutionContext(context.Background(), state)
	got, err := r.EvaluateExpression(ctx, "1.2300")
	if err != nil || got != json.Number("1.2300") {
		t.Fatalf("generic numeric compatibility: %v %v", got, err)
	}
	got, err = r.EvaluateExpression(ctx, "$response.body#%2Fx")
	if err != nil || got != "value" {
		t.Fatalf("encoded root compatibility: %v %v", got, err)
	}
	got, err = r.EvaluateExpression(ctx, "$variables.literal == \"$response.body#%2Fx\"")
	if err != nil || got != true {
		t.Fatalf("literal was rewritten: %v %v", got, err)
	}
}
func TestReferenceAdoptionKeepsUnshadowedComponentVariables(t *testing.T) {
	r, err := NewRuntime(&uws1.Document{UWS: "1.12.0", Variables: map[string]any{"shared": "top"}, Components: &uws1.Components{Variables: map[string]any{"shared": "component", "only": "available"}}}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := r.EvaluateExpression(context.Background(), "$variables.only")
	if err != nil || got != "available" {
		t.Fatalf("component-only declaration lost: %v %v", got, err)
	}
}
