package binding

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/expressions"
	"github.com/OpenUdon/uws/uws1"
)

func TestBrowserBindingCoreBodyTemplatesAndPartialEvidence(t *testing.T) {
	table := browserTableFixture()
	table.Operations[0].Inputs[0].Schema.JSON = json.RawMessage(`{"type":"object","required":["text"],"properties":{"text":{"type":"string"}},"additionalProperties":false}`)
	r, err := NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{Binding: Binding{Source: table.Sources[0], SelectorKind: "id", SelectorValue: "read"}, Inputs: []BoundInput{{Location: "body", Name: "body", Value: map[string]any{"text": "$inputs.text"}}}, ExpressionContext: expressions.Context{Version: "1.13.0", Field: expressions.Value}, ExpressionTypes: map[string]Schema{"$inputs.text": {Known: true, JSON: json.RawMessage(`{"type":"string"}`)}}}
	got, err := ValidateBinding(t.Context(), r, req)
	if err != nil || got.Outcome != Compatible {
		t.Fatal(got, err)
	}
	req.ExpressionTypes = nil
	got, _ = ValidateBinding(t.Context(), r, req)
	if got.Outcome != Indeterminate {
		t.Fatal("unproved browser type accepted", got)
	}
	req.ExpressionTypes = map[string]Schema{"$inputs.text": {Known: true, JSON: json.RawMessage(`{"type":"integer"}`)}}
	got, _ = ValidateBinding(t.Context(), r, req)
	if got.Outcome != Incompatible {
		t.Fatal("wrong browser template type accepted", got)
	}
	req.Inputs[0].Value = map[string]any{"text": "{{{{literal}}}}"}
	got, _ = ValidateBinding(t.Context(), r, req)
	if got.Outcome != Compatible {
		t.Fatal("native braces reinterpreted", got)
	}
	req.Inputs[0].Location = "query"
	got, _ = ValidateBinding(t.Context(), r, req)
	if got.Outcome != Incompatible || got.Diagnostics[0].Code != "binding.browser_input_location" {
		t.Fatal(got)
	}
	table.Operations[0].Complete = false
	table.Operations[0].Inputs[0].Schema.Known = false
	r, _ = NewResolver(table)
	req.Inputs[0].Location = "body"
	req.OutputReferences = []OutputReference{{Location: "body", Name: "unknown"}}
	got, _ = ValidateBinding(t.Context(), r, req)
	if got.Outcome != Indeterminate {
		t.Fatal("incomplete evidence upgraded", got)
	}
	data, _ := json.Marshal(got)
	if strings.Contains(string(data), "literal") || strings.Contains(string(data), "$inputs") {
		t.Fatal("binding values leaked")
	}
}

func TestBrowserCredentialDeclarationsRemainSymbolic(t *testing.T) {
	table := browserTableFixture()
	op := &table.Operations[0]
	op.Inputs = nil
	op.Outputs = nil
	op.Browser.CallKind = "authentication"
	op.Browser.ProfileVersion = "uws.browser-authentication.1.1"
	op.Browser.Effects = json.RawMessage(`["establishes_session"]`)
	op.Browser.ConfirmationPolicy = nil
	op.Browser.CredentialSlots = []CredentialSlotShape{{Name: "password", Kind: "password", Required: true}, {Name: "identifier", Kind: "identifier", Required: false}}
	r, err := NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	req := Request{Binding: Binding{Source: table.Sources[0], SelectorKind: "id", SelectorValue: "read"}}
	got, _ := ValidateBinding(t.Context(), r, req)
	if got.Outcome != Incompatible || got.Diagnostics[0].Code != "binding.browser_credential_missing" {
		t.Fatal(got)
	}
	req.Security = []SecurityBinding{{Scheme: "password", CredentialSlot: "host-password"}}
	got, _ = ValidateBinding(t.Context(), r, req)
	if got.Outcome != Compatible {
		t.Fatal(got)
	}
	op.Complete = false
	r, _ = NewResolver(table)
	req.Security = nil
	got, _ = ValidateBinding(t.Context(), r, req)
	if got.Outcome != Indeterminate {
		t.Fatal("partial credential declaration became conclusive", got)
	}
	data, _ := json.Marshal(got)
	if strings.Contains(string(data), "host-password") {
		t.Fatal("host binding leaked")
	}
}

func TestBrowserFlowUsesBodyReferencesWithoutProfileInterpretation(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", SourceDescriptions: []*uws1.SourceDescription{{Name: "browser", Type: KindBrowser}}, Operations: []*uws1.Operation{
		{OperationID: "fetch", Effect: uws1.OperationEffectRead},
		{OperationID: "consume", SourceDescription: "browser", SourceOperationID: "read", Effect: uws1.OperationEffectUnknown, Request: map[string]any{"body": map[string]string{"text": "$steps.fetch.outputs.body"}, "query": map[string]any{"ignored": "$steps.fetch.outputs.query"}, "x-template": "$steps.fetch.outputs.extension"}, Extensions: map[string]any{"x-profile": "$steps.fetch.outputs.profile"}},
	}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "fetch", OperationRef: "fetch", Outputs: map[string]string{"body": "$response.body", "query": "$response.body", "extension": "$response.body", "profile": "$response.body"}}, {StepID: "consume", OperationRef: "consume"}}}}}
	before, _ := json.Marshal(d)
	a, err := AnalyzeFlow(t.Context(), d)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := AnalyzeFlow(t.Context(), d)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("nondeterministic browser flow")
	}
	if hasCode(a, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/body") {
		t.Fatal("browser body reference skipped", a)
	}
	for _, name := range []string{"query", "extension", "profile"} {
		if !hasCode(a, "flow.output_unreferenced", "/workflows/0/steps/0/outputs/"+name) {
			t.Fatal("profile interpreted", name, a)
		}
	}
	if !hasCode(a, "flow.effect_unknown", "/workflows/0/steps/1") {
		t.Fatal("effect authority invented", a)
	}
	after, _ := json.Marshal(d)
	if string(before) != string(after) {
		t.Fatal("flow changed browser document")
	}
}
