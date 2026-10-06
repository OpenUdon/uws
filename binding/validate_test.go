package binding

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func requestFixture() Request {
	return Request{Binding: Binding{Source: sourceFixture(), SelectorKind: "id", SelectorValue: "read"}, Inputs: []BoundInput{{Location: "query", Name: "city", Value: "Toronto"}}, Security: []SecurityBinding{{Scheme: "bearer", CredentialSlot: "owner-token"}}}
}
func TestLiteralRequiredAndIndeterminateSchemas(t *testing.T) {
	table := tableFixture()
	resolver, _ := NewResolver(table)
	r, err := ValidateBinding(t.Context(), resolver, requestFixture())
	if err != nil || r.Outcome != Compatible {
		t.Fatalf("valid literal: %+v %v", r, err)
	}
	req := requestFixture()
	req.Inputs = nil
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Incompatible || r.Diagnostics[0].Code != "binding.input_required" {
		t.Fatal(r)
	}
	req = requestFixture()
	req.Inputs[0].Value = 123
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Incompatible {
		t.Fatal(r)
	}
	table.Operations[0].Inputs[0].Schema = Schema{Known: false, JSON: json.RawMessage(`{"type":"string"}`)}
	resolver, _ = NewResolver(table)
	r, _ = ValidateBinding(t.Context(), resolver, requestFixture())
	if r.Outcome != Indeterminate {
		t.Fatal("partial schema became compatible", r)
	}
	table.Operations[0].Inputs[0].Schema = Schema{Known: true, JSON: json.RawMessage(`{"$ref":"https://must-not-fetch.invalid/schema"}`)}
	resolver, _ = NewResolver(table)
	r, _ = ValidateBinding(t.Context(), resolver, requestFixture())
	if r.Outcome != Indeterminate {
		t.Fatal("unavailable ref became valid", r)
	}
}
func TestKnownExpressionTypesWithoutConstraintGuessing(t *testing.T) {
	resolver, _ := NewResolver(tableFixture())
	req := requestFixture()
	req.Inputs[0].Value = "$inputs.city"
	r, _ := ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Indeterminate {
		t.Fatal(r)
	}
	req.ExpressionTypes = map[string]Schema{"$inputs.city": {Known: true, JSON: json.RawMessage(`{"type":"string","maxLength":20}`)}}
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Compatible {
		t.Fatal(r)
	}
	req.ExpressionTypes["$inputs.city"] = Schema{Known: true, JSON: json.RawMessage(`{"type":"integer"}`)}
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Incompatible {
		t.Fatal(r)
	}
	table := tableFixture()
	table.Operations[0].Inputs[0].Schema = Schema{Known: true, JSON: json.RawMessage(`{"type":"string","maxLength":5}`)}
	resolver, _ = NewResolver(table)
	req.ExpressionTypes["$inputs.city"] = Schema{Known: true, JSON: json.RawMessage(`{"type":"string"}`)}
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Indeterminate {
		t.Fatal("unproved constraint containment", r)
	}
}
func TestSecurityAlternativesAreOrAndRequirementsAreAnd(t *testing.T) {
	resolver, _ := NewResolver(tableFixture())
	req := requestFixture()
	req.Security = []SecurityBinding{{Scheme: "key", CredentialSlot: "owner-key"}}
	r, _ := ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Incompatible {
		t.Fatal("partial AND satisfied", r)
	}
	req.Security = append(req.Security, SecurityBinding{Scheme: "oauth", CredentialSlot: "owner-oauth", Scopes: []string{"weather.read"}})
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Compatible {
		t.Fatal(r)
	}
	req.Security[1].Scopes = nil
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Incompatible {
		t.Fatal("missing scope satisfied", r)
	}
	table := tableFixture()
	table.Operations[0].Security = Security{Known: false}
	resolver, _ = NewResolver(table)
	r, _ = ValidateBinding(t.Context(), resolver, requestFixture())
	if r.Outcome != Indeterminate {
		t.Fatal("unknown security became anonymous", r)
	}
	table.Operations[0].Security = Security{Known: true, Alternatives: []SecurityAlternative{{}}}
	resolver, _ = NewResolver(table)
	req = requestFixture()
	req.Security = nil
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Compatible {
		t.Fatal("known anonymous refused", r)
	}
}
func TestResponseFieldReferencesAndValueFreeDiagnostics(t *testing.T) {
	table := tableFixture()
	table.Operations[0].Outputs = []Output{{Location: "body", Name: "response", Schema: Schema{Known: true, JSON: json.RawMessage(`{"type":"object","additionalProperties":false,"properties":{"items":{"type":"array","items":{"type":"object","additionalProperties":false,"properties":{"name":{"type":"string"}}}}}}`)}}}
	resolver, _ := NewResolver(table)
	req := requestFixture()
	req.OutputReferences = []OutputReference{{Location: "body", Name: "response", Pointer: "#/items/0/name"}}
	r, _ := ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Compatible {
		t.Fatal(r)
	}
	req.OutputReferences[0].Pointer = "#/items/00/name"
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Incompatible {
		t.Fatal("noncanonical index accepted", r)
	}
	req.OutputReferences[0].Pointer = "#/missing"
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Incompatible {
		t.Fatal("missing known field accepted", r)
	}
	req = requestFixture()
	req.Inputs[0].Value = map[string]any{"token": "private-canary"}
	r, _ = ValidateBinding(t.Context(), resolver, req)
	data, _ := json.Marshal(r)
	if strings.Contains(string(data), "private-canary") {
		t.Fatal("literal exposed")
	}
}

type misleadingResolver struct{ shape OperationShape }

func (r misleadingResolver) Resolve(context.Context, Binding) (Resolution, error) {
	return Resolution{Status: Resolved, Shape: &r.shape}, nil
}
func TestCustomResolverWrongSelectorAndCancellation(t *testing.T) {
	op := shapeFixture()
	op.Selector.Value = "wrong"
	op.Aliases = nil
	r, err := ValidateBinding(t.Context(), misleadingResolver{op}, requestFixture())
	if err != nil || r.Outcome != Indeterminate {
		t.Fatal("wrong-selector resolver trusted", r, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	resolver, _ := NewResolver(tableFixture())
	if _, err := ValidateBinding(ctx, resolver, requestFixture()); err != context.Canceled {
		t.Fatal(err)
	}
}

func TestUnavailableOutputSchemasNeverClaimCompatibility(t *testing.T) {
	table := tableFixture()
	table.Operations[0].Outputs = []Output{{Location: "body", Name: "response", Schema: Schema{Known: true, JSON: json.RawMessage(`{"$ref":"https://must-not-fetch.invalid/schema"}`)}}}
	resolver, _ := NewResolver(table)
	req := requestFixture()
	req.OutputReferences = []OutputReference{{Location: "body", Name: "response"}}
	r, _ := ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Indeterminate {
		t.Fatal(r)
	}
	table.Operations[0].Complete = false
	table.Operations[0].Outputs = nil
	resolver, _ = NewResolver(table)
	r, _ = ValidateBinding(t.Context(), resolver, req)
	if r.Outcome != Indeterminate {
		t.Fatal("incomplete missing field became definitive", r)
	}
}
