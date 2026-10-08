package binding

import (
	"encoding/json"
	"testing"

	"github.com/OpenUdon/uws/expressions"
)

func knownSchema(text string) Schema { return Schema{Known: true, JSON: json.RawMessage(text)} }

func TestPatternSourceDoesNotProveClosedContainment(t *testing.T) {
	source := knownSchema(`{"type":"object","patternProperties":{"^x":{"type":"string"}},"additionalProperties":false}`)
	target := knownSchema(`{"type":"object","additionalProperties":false}`)
	if got := schemaCompatibility(source, target); got != Indeterminate {
		t.Fatal(got)
	}
}

func TestEffective2019ArrayConstraintsRemainIndeterminate(t *testing.T) {
	target := knownSchema(`{"$schema":"https://json-schema.org/draft/2019-09/schema","type":"array","unevaluatedItems":false}`)
	value := []any{"$inputs.marker"}
	types := map[string]Schema{"$inputs.marker": knownSchema(`{"type":"string"}`)}
	if got := validateLiteral(target, []any{"value"}); got != Incompatible {
		t.Fatal(got)
	}
	if got := validateBoundValue(target, value, types, expressions.Context{Version: "1.13.0", Field: expressions.Value}); got != Indeterminate {
		t.Fatal(got)
	}
}

func TestIgnoredReferenceSiblingDoesNotProveType(t *testing.T) {
	source := knownSchema(`{"$schema":"http://json-schema.org/draft-07/schema#","definitions":{"n":{"type":"integer"}},"$ref":"#/definitions/n","type":"string"}`)
	target := knownSchema(`{"type":"string"}`)
	if validateLiteral(source, 1) != Compatible || validateLiteral(target, 1) != Incompatible {
		t.Fatal("dialect fixture")
	}
	if got := schemaCompatibility(source, target); got != Indeterminate {
		t.Fatal("ignored type proved", got)
	}
}

func TestSchemaPathsFalseNullableAndPatterns(t *testing.T) {
	for _, tc := range []struct {
		schema, path string
		want         Outcome
	}{
		{`{"type":"object","properties":{"x":false}}`, "#/x", Incompatible},
		{`{"type":"object","properties":{"x":{"type":"string"}},"patternProperties":{"^x":false},"additionalProperties":false}`, "#/x", Indeterminate},
		{`{"type":["object","null"],"properties":{"x":{"type":"string"}},"additionalProperties":false}`, "#/x", Indeterminate},
		{`{"type":"object","patternProperties":{"^x":{"type":"string"}},"additionalProperties":false}`, "#/xyz", Indeterminate},
		{`{"$schema":"http://json-schema.org/draft-07/schema#","type":"array","items":[false],"additionalItems":false}`, "#/0", Incompatible},
		{`{"$schema":"http://json-schema.org/draft-07/schema#","type":"array","prefixItems":[false]}`, "#/0", Indeterminate},
	} {
		if got := schemaPath(knownSchema(tc.schema), tc.path); got != tc.want {
			t.Fatalf("%s at %s: %s", tc.schema, tc.path, got)
		}
	}
}

func TestNestedDraftAndReferenceContext(t *testing.T) {
	for _, target := range []string{
		`{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"x":{"type":"array","items":[{"type":"integer"}],"additionalItems":false}},"required":["x"],"additionalProperties":false}`,
		`{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"x":{"type":"array","prefixItems":[false]}},"required":["x"],"additionalProperties":false}`,
		`{"$schema":"http://json-schema.org/draft-07/schema#","definitions":{"integer":{"type":"integer"}},"type":"object","properties":{"x":{"type":"array","items":[{"$ref":"#/definitions/integer"}],"additionalItems":false}},"required":["x"],"additionalProperties":false}`,
	} {
		literal := map[string]any{"x": []any{json.Number("1")}}
		if got := validateLiteral(knownSchema(target), literal); got != Compatible {
			t.Fatal("whole literal", got)
		}
		// An expression-containing template also has a constant sibling. This
		// proves the sibling using the child compiled in the original resource.
		value := map[string]any{"x": []any{1}, "marker": "$inputs.marker"}
		// Keep the extra expression in a target-declared property.
		var m map[string]any
		_ = json.Unmarshal([]byte(target), &m)
		m["properties"].(map[string]any)["marker"] = map[string]any{"type": "string"}
		data, _ := json.Marshal(m)
		types := map[string]Schema{"$inputs.marker": knownSchema(`{"type":"string"}`)}
		if got := validateBoundValue(Schema{Known: true, JSON: data}, value, types, expressions.Context{Version: "1.13.0", Field: expressions.Value}); got != Compatible {
			t.Fatalf("nested proof %s: %s", target, got)
		}
	}
}
