package binding

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/expressions"
)

func integerBrowserRequest(t *testing.T, profile string, target string, value any, proof Schema) (Report, error) {
	t.Helper()
	table := browserTableFixture()
	table.Operations[0].Browser.ProfileVersion = profile
	table.Operations[0].Inputs = []Input{{Location: "body", Name: "integer", Schema: Schema{Known: true, JSON: json.RawMessage(target)}}}
	r, err := NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{Binding: Binding{Source: table.Sources[0], SelectorKind: "id", SelectorValue: "read"}, ExpressionContext: expressions.Context{Version: "1.13.0", Field: expressions.Value}}
	if value != nil {
		req.Inputs = []BoundInput{{Location: "body", Name: "integer", Value: value}}
	}
	if proof.Known {
		req.ExpressionTypes = map[string]Schema{"$inputs.value": proof}
	}
	return ValidateBinding(t.Context(), r, req)
}

func TestBrowserNativeIntegerLiteralRangesAndDefaults(t *testing.T) {
	for _, c := range []struct {
		profile, value string
		want           Outcome
	}{
		{"uws.browser.1.8", "9223372036854775807", Compatible}, {"uws.browser.1.8", "-9223372036854775808", Compatible},
		{"uws.browser.1.8", "9223372036854775808", Incompatible}, {"uws.browser.1.8", "-9223372036854775809", Incompatible},
		{"uws.browser.1.9", "9007199254740991", Compatible}, {"uws.browser.1.9", "-9007199254740991", Compatible},
		{"uws.browser.1.9", "9007199254740992", Incompatible}, {"uws.browser.1.10", "-9007199254740992", Incompatible},
		{"uws.browser.1.10", "9007199254740991e0", Compatible}, {"uws.browser.1.10", "1.00", Compatible},
		{"uws.browser.1.10", "-0", Compatible}, {"uws.browser.1.10", "1.5", Incompatible},
		{"uws.browser.1.7", "9223372036854775808", Compatible},
	} {
		t.Run(c.profile+"/"+c.value, func(t *testing.T) {
			r, err := integerBrowserRequest(t, c.profile, `{"type":"integer"}`, json.Number(c.value), Schema{})
			if err != nil || r.Outcome != c.want {
				t.Fatal(r, err)
			}
			if c.want == Incompatible && (len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "binding.browser_integer_range") {
				t.Fatal("missing native value-free code", r)
			}
		})
	}
	for _, target := range []string{`{"type":"integer","default":9007199254740992}`, `{"type":"object","properties":{"n":{"type":"integer","default":9007199254740992}}}`} {
		r, err := integerBrowserRequest(t, "uws.browser.1.9", target, nil, Schema{})
		if err != nil || r.Outcome != Incompatible {
			t.Fatal("unsafe omitted-input default accepted", r, err)
		}
	}
}

func TestBrowserNativeIntegerSymbolicProofAndUnknownTypes(t *testing.T) {
	for _, c := range []struct {
		proof string
		want  Outcome
	}{
		{`{"type":"integer"}`, Indeterminate},
		{`{"type":"integer","minimum":-9007199254740991,"maximum":9007199254740991}`, Compatible},
		{`{"const":9007199254740992}`, Incompatible},
		{`{"const":9223372036854775807}`, Incompatible},
		{`{"const":1.00}`, Compatible},
		{`{"enum":[0,9007199254740992]}`, Indeterminate},
		{`{"enum":[9007199254740992,9007199254740993]}`, Incompatible},
	} {
		t.Run(c.proof, func(t *testing.T) {
			r, err := integerBrowserRequest(t, "uws.browser.1.9", `{"type":"integer"}`, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(c.proof)})
			if err != nil || r.Outcome != c.want {
				t.Fatal(r, err)
			}
		})
	}
	r, err := integerBrowserRequest(t, "uws.browser.1.9", `{"type":"integer"}`, "$inputs.value", Schema{})
	if err != nil || r.Outcome != Indeterminate {
		t.Fatal("untyped symbol accepted", r, err)
	}
	r, err = integerBrowserRequest(t, "uws.browser.1.9", `{}`, json.Number("9007199254740992"), Schema{})
	if err != nil || r.Outcome != Indeterminate {
		t.Fatal("unknown native parameter type guessed", r, err)
	}
	r, err = integerBrowserRequest(t, "uws.browser.1.8", `{"type":"integer"}`, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(`{"const":9223372036854775807}`)})
	if err != nil || r.Outcome != Compatible {
		t.Fatal("wide exact proof lost", r, err)
	}
}

func TestBrowserNativeIntegerWholeBodyReferencesAndLiteralContainers(t *testing.T) {
	target := `{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"]}`
	r, err := integerBrowserRequest(t, "uws.browser.1.9", target, map[string]any{"n": json.Number("9007199254740992")}, Schema{})
	if err != nil || r.Outcome != Incompatible {
		t.Fatal(r, err)
	}
	r, err = integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(`{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"]}`)})
	if err != nil || r.Outcome != Indeterminate {
		t.Fatal("whole-body integer proof bypass", r, err)
	}
}

func TestBrowserNativeIntegerWholeBodyProofRequiresCompletePropertyInventory(t *testing.T) {
	target := `{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"]}`
	for name, proof := range map[string]string{
		"closed":                          `{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`,
		"open omitted":                    `{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1}},"required":["n"]}`,
		"open explicit":                   `{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":true}`,
		"open bounded extras unsupported": `{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":{"type":"integer","minimum":0,"maximum":1}}`,
		"closed declared unsafe extra":    `{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1},"extra":{"const":9007199254740992}},"required":["n"],"additionalProperties":false}`,
		"closed declared safe extra lacks native type": `{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1},"extra":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`,
		"pattern extras": `{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false,"patternProperties":{"^extra":{"type":"integer"}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			r, err := integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(proof)})
			want := Indeterminate
			if name == "closed" {
				want = Compatible
			}
			if err != nil || r.Outcome != want {
				t.Fatal("incomplete property proof accepted", r, err)
			}
			if want == Indeterminate && (len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "binding.browser_integer_range") {
				t.Fatal("missing stable native range diagnostic", r)
			}
		})
	}
	// Declared extras can be proved when the complete closed source and native
	// target both describe them; no schema bytes are patched to force closure.
	target = `{"type":"object","properties":{"n":{"type":"integer"},"extra":{"type":"integer"}},"required":["n"]}`
	proof := Schema{Known: true, JSON: json.RawMessage(`{"type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1},"extra":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`)}
	r, err := integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", proof)
	if err != nil || r.Outcome != Compatible {
		t.Fatal("complete declared property proof refused", r, err)
	}
}

func TestBrowserNativeIntegerReferenceProofDoesNotTrustIgnoredSiblings(t *testing.T) {
	for name, raw := range map[string]string{
		"direct draft07":         `{"$schema":"http://json-schema.org/draft-07/schema#","definitions":{"unsafe":{"const":9007199254740992}},"$ref":"#/definitions/unsafe","type":"integer","minimum":0,"maximum":1}`,
		"const sibling unproved": `{"$schema":"http://json-schema.org/draft-07/schema#","definitions":{"unsafe":{"const":9007199254740992}},"$ref":"#/definitions/unsafe","type":"integer","const":0}`,
		"inherited draft07":      `{"$schema":"http://json-schema.org/draft-07/schema#","definitions":{"unsafe":{"const":9007199254740992}},"type":"object","properties":{"n":{"$ref":"#/definitions/unsafe","type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`,
		"ancestor draft07":       `{"$schema":"http://json-schema.org/draft-07/schema#","definitions":{"unsafe":{"const":{"n":9007199254740992}}},"$ref":"#/definitions/unsafe","type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`,
	} {
		t.Run(name, func(t *testing.T) {
			schema := Schema{Known: true, JSON: json.RawMessage(raw)}
			// Primary regressions use an effective unsafe compiler witness.
			// Const siblings separately exercise unproved reference context;
			// ordinary equality-based containment is intentionally unchanged.
			compiled, err := compile(schema)
			if err != nil {
				t.Fatal(err)
			}
			value := any(json.Number("9007199254740992"))
			if name == "inherited draft07" || name == "ancestor draft07" {
				value = map[string]any{"n": json.Number("9007199254740992")}
			}
			if name != "const sibling unproved" && compiled.Validate(value) != nil {
				t.Fatal("regression does not exercise ignored Draft07 siblings")
			}
			if schemaCompatibility(schema, schema) != Compatible {
				t.Fatal("generic identical-schema compatibility changed")
			}
			r, err := integerBrowserRequest(t, "uws.browser.1.9", raw, "$inputs.value", schema)
			if err != nil || r.Outcome != Indeterminate || len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "binding.browser_integer_range" {
				t.Fatal("reference siblings proved native safety", r, err)
			}
			plainTarget := `{"type":"integer"}`
			if name == "inherited draft07" || name == "ancestor draft07" {
				plainTarget = `{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"],"additionalProperties":false}`
			}
			// Also isolate the source-proof guard: a reference-free target must
			// not mask a missing native diagnostic behind generic containment.
			r, err = integerBrowserRequest(t, "uws.browser.1.9", plainTarget, "$inputs.value", schema)
			if err != nil || r.Outcome != Indeterminate || len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "binding.browser_integer_range" {
				t.Fatal("reference-bearing source proof read raw siblings", r, err)
			}
			if string(schema.JSON) != raw {
				t.Fatal("numeric proof rewrote raw source schema")
			}
		})
	}
}

func TestBrowserNativeIntegerUnprovedDynamicAndRecursiveReferences(t *testing.T) {
	for _, keyword := range []string{"$dynamicRef", "$recursiveRef"} {
		for _, nested := range []bool{false, true} {
			raw := `{"$schema":"http://json-schema.org/draft-07/schema#","` + keyword + `":"#unknown","type":"integer","minimum":0,"maximum":1}`
			if nested {
				raw = `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"n":{"` + keyword + `":"#unknown","type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`
			}
			schema := Schema{Known: true, JSON: json.RawMessage(raw)}
			if schemaCompatibility(schema, schema) != Compatible {
				t.Fatal("generic equality control failed")
			}
			r, err := integerBrowserRequest(t, "uws.browser.1.9", raw, "$inputs.value", schema)
			if err != nil || r.Outcome != Indeterminate {
				t.Fatal("unproved reference context became positive", keyword, nested, r, err)
			}
		}
	}
	// Reference-shaped data and property names are not schema references.
	raw := `{"type":"object","properties":{"n":{"type":"integer"},"$ref":{"type":"string"}},"required":["n","$ref"],"additionalProperties":false}`
	proof := Schema{Known: true, JSON: json.RawMessage(`{"const":{"n":1,"$ref":"literal data"}}`)}
	r, err := integerBrowserRequest(t, "uws.browser.1.9", raw, "$inputs.value", proof)
	if err != nil || r.Outcome != Compatible {
		t.Fatal("literal reference-shaped data reinterpreted", r, err)
	}
	for _, raw := range []string{`{"$schema":"http://json-schema.org/draft-07/schema#","type":"integer","minimum":0,"maximum":1}`, `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"n":{"type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`} {
		r, err := integerBrowserRequest(t, "uws.browser.1.9", raw, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(raw)})
		if err != nil || r.Outcome != Compatible {
			t.Fatal("reference-free Draft07 safe control refused", r, err)
		}
	}
}

func TestBrowserNativeIntegerUnsupportedDialectDoesNotTrustConst(t *testing.T) {
	for name, raw := range map[string]string{
		"direct draft04":    `{"$schema":"http://json-schema.org/draft-04/schema#","type":"integer","const":0}`,
		"inherited draft04": `{"$schema":"http://json-schema.org/draft-04/schema#","type":"object","properties":{"n":{"type":"integer","const":0}},"required":["n"],"additionalProperties":false}`,
	} {
		t.Run(name, func(t *testing.T) {
			proof := Schema{Known: true, JSON: json.RawMessage(raw)}
			compiled, err := compile(proof)
			if err != nil {
				t.Fatal(err)
			}
			target := `{"type":"integer"}`
			value := any(json.Number("9007199254740992"))
			if name == "inherited draft04" {
				target = `{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"],"additionalProperties":false}`
				value = map[string]any{"n": json.Number("9007199254740992")}
			}
			if compiled.Validate(value) != nil {
				t.Fatal("Draft04 witness does not exercise ignored const")
			}
			r, err := integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", proof)
			if err != nil || r.Outcome != Indeterminate || len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "binding.browser_integer_range" {
				t.Fatal("unsupported dialect proved safety", r, err)
			}
			if schemaCompatibility(proof, proof) != Compatible {
				t.Fatal("generic compatibility changed")
			}
			if string(proof.JSON) != raw {
				t.Fatal("proof schema rewritten")
			}
		})
	}
}

func TestBrowserNativeIntegerSupportedSubsetCoversEveryProofPath(t *testing.T) {
	for _, dialect := range []string{"", "http://json-schema.org/draft-06/schema#", "http://json-schema.org/draft-07/schema#", "https://json-schema.org/draft/2019-09/schema", "https://json-schema.org/draft/2020-12/schema"} {
		for name, raw := range map[string]string{"const": `{"const":0}`, "enum": `{"enum":[0,1]}`, "range": `{"type":"integer","minimum":0,"maximum":1}`, "object": `{"type":"object","properties":{"n":{"const":0}},"required":["n"],"additionalProperties":false}`} {
			t.Run(dialect+"/"+name, func(t *testing.T) {
				if dialect != "" {
					raw = `{"$schema":"` + dialect + `",` + raw[1:]
				}
				target := `{"type":"integer"}`
				if name == "object" {
					target = `{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"],"additionalProperties":false}`
				}
				r, err := integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(raw)})
				if err != nil || r.Outcome != Compatible {
					t.Fatal("supported proof control refused", r, err)
				}
			})
		}
	}
	for _, raw := range []string{
		`{"$schema":"https://unknown.invalid/schema","type":"integer","const":0}`,
		`{"type":"object","properties":{"n":{"$schema":"http://json-schema.org/draft-04/schema#","type":"integer","const":0}},"required":["n"],"additionalProperties":false}`,
		`{"type":"object","properties":{"n":{"$schema":"https://unknown.invalid/schema","type":"integer","minimum":0,"maximum":1}},"required":["n"],"additionalProperties":false}`,
		`{"type":"integer","const":0,"allOf":[{"$ref":"#unknown"}]}`,
	} {
		target := `{"type":"integer"}`
		if strings.Contains(raw, `"properties"`) {
			target = `{"type":"object","properties":{"n":{"type":"integer"}},"required":["n"],"additionalProperties":false}`
		}
		r, err := integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(raw)})
		if err != nil || r.Outcome != Indeterminate || len(r.Diagnostics) == 0 || r.Diagnostics[0].Code != "binding.browser_integer_range" {
			t.Fatal("unproved tree context accepted", r, err)
		}
	}
	// Literal schema-like data and names do not select a dialect or loader.
	target := `{"type":"object","properties":{"$schema":{"type":"string"},"n":{"type":"integer"}},"required":["$schema","n"],"additionalProperties":false}`
	proof := Schema{Known: true, JSON: json.RawMessage(`{"const":{"$schema":"http://json-schema.org/draft-04/schema#","n":0}}`)}
	r, err := integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", proof)
	if err != nil || r.Outcome != Compatible {
		t.Fatal("dialect-shaped data reinterpreted", r, err)
	}
}

func TestBrowserNativeIntegerSymbolicArrayNeedsEffectiveItemProof(t *testing.T) {
	raw := `{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"array","prefixItems":[{"type":"integer"}],"items":{"type":"integer","minimum":0,"maximum":1}}`
	schema := Schema{Known: true, JSON: json.RawMessage(raw)}
	compiled, err := compile(schema)
	if err != nil {
		t.Fatal(err)
	}
	if compiled.Validate([]any{json.Number("9007199254740992")}) != nil {
		t.Fatal("prefix witness does not exercise items gap")
	}
	if schemaCompatibility(schema, schema) != Compatible {
		t.Fatal("generic compatibility changed")
	}
	r, err := integerBrowserRequest(t, "uws.browser.1.9", raw, "$inputs.value", schema)
	if err != nil || r.Outcome != Indeterminate {
		t.Fatal("typed items proof ignored unsafe prefix", r, err)
	}
	target := `{"type":"array","items":{"type":"integer"}}`
	r, err = integerBrowserRequest(t, "uws.browser.1.9", target, "$inputs.value", Schema{Known: true, JSON: json.RawMessage(`{"const":[0,1]}`)})
	if err != nil || r.Outcome != Compatible {
		t.Fatal("finite constant array control refused", r, err)
	}
}
