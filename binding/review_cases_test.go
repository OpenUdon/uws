package binding

import (
	"encoding/json"
	"testing"
)

func TestReviewNestedExpressionBinding(t *testing.T) {
	table := tableFixture()
	table.Operations[0].Inputs = []Input{{Location: "body", Name: "request", Required: true, Schema: Schema{Known: true, JSON: json.RawMessage(`{"type":"object","properties":{"count":{"type":"integer"}},"required":["count"],"additionalProperties":false}`)}}}
	resolver, _ := NewResolver(table)
	req := requestFixture()
	req.Inputs = []BoundInput{{Location: "body", Name: "request", Value: map[string]any{"count": "$inputs.count"}}}
	req.ExpressionTypes = map[string]Schema{"$inputs.count": {Known: true, JSON: json.RawMessage(`{"type":"integer"}`)}}
	r, err := ValidateBinding(t.Context(), resolver, req)
	if err != nil || r.Outcome != Compatible {
		t.Fatalf("nested expression falsely rejected: %+v %v", r, err)
	}
}
func TestReviewNumberIntegerPartialOverlap(t *testing.T) {
	table := tableFixture()
	table.Operations[0].Inputs[0].Schema = Schema{Known: true, JSON: json.RawMessage(`{"type":"integer"}`)}
	resolver, _ := NewResolver(table)
	req := requestFixture()
	req.Inputs[0].Value = "$inputs.n"
	req.ExpressionTypes = map[string]Schema{"$inputs.n": {Known: true, JSON: json.RawMessage(`{"type":"number"}`)}}
	r, err := ValidateBinding(t.Context(), resolver, req)
	if err != nil || r.Outcome != Indeterminate {
		t.Fatalf("partial number/integer overlap classified disjoint: %+v %v", r, err)
	}
}

func TestReviewNestedArraysAndUnsupportedConstraints(t *testing.T) {
	table := tableFixture()
	table.Operations[0].Inputs = []Input{{Location: "body", Name: "request", Schema: Schema{Known: true, JSON: json.RawMessage(`{"type":"array","items":{"type":"integer"},"minItems":1,"maxItems":2}`)}}}
	resolver, _ := NewResolver(table)
	req := requestFixture()
	req.Inputs = []BoundInput{{Location: "body", Name: "request", Value: []any{"$inputs.n", 2}}}
	req.ExpressionTypes = map[string]Schema{"$inputs.n": {Known: true, JSON: json.RawMessage(`{"type":"integer"}`)}}
	r, _ := ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Compatible {
		t.Fatal(r)
	}
	table.Operations[0].Inputs[0].Schema = Schema{Known: true, JSON: json.RawMessage(`{"type":"array","items":{"type":"integer"},"uniqueItems":true}`)}
	resolver, _ = NewResolver(table)
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Indeterminate {
		t.Fatal("unproved uniqueness asserted", r)
	}
}
func TestReviewOpenObjectSourceDoesNotProveOptionalTargetPropertyType(t *testing.T) {
	table := tableFixture()
	table.Operations[0].Inputs[0].Schema = Schema{Known: true, JSON: json.RawMessage(`{"type":"object","properties":{"n":{"type":"integer"}}}`)}
	resolver, _ := NewResolver(table)
	req := requestFixture()
	req.Inputs[0].Value = "$inputs.object"
	req.ExpressionTypes = map[string]Schema{"$inputs.object": {Known: true, JSON: json.RawMessage(`{"type":"object"}`)}}
	r, _ := ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Indeterminate {
		t.Fatal("open object constraints assumed", r)
	}
}
