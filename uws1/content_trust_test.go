package uws1

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func contentTrustDocument() *Document {
	return &Document{
		UWS:  "1.9.1",
		Info: &Info{Title: "Trust", Version: "1.0.0"},
		SourceDescriptions: []*SourceDescription{{
			Name: "mail_api", URL: "./mail.yaml", Type: SourceDescriptionTypeOpenAPI,
		}},
		Operations: []*Operation{{
			OperationID: "read_message", SourceDescription: "mail_api", SourceOperationID: "readMessage",
			Outputs: map[string]string{"body": "$response.body", "message_id": "$response.body#/id"},
		}},
		Workflows: []*Workflow{{
			WorkflowID: "summarize", Type: WorkflowTypeSequence,
			Inputs: &ParamSchema{Type: "object", Properties: map[string]*ParamSchema{"locale": {Type: "string"}}},
		}},
		Triggers: []*Trigger{{TriggerID: "incoming_mail"}},
		ContentTrust: &ContentTrust{
			SourceDescriptions: map[string]ContentTrustLevel{"mail_api": ContentTrustUntrusted},
			Operations: map[string]*OperationContentTrust{
				"read_message": {Default: ContentTrustUntrusted, Outputs: map[string]ContentTrustLevel{"message_id": ContentTrustTrusted, "body": ContentTrustUntrusted}},
			},
			Triggers: map[string]ContentTrustLevel{"incoming_mail": ContentTrustUntrusted},
			Workflows: map[string]*WorkflowContentTrust{
				"summarize": {Default: ContentTrustUnknown, Inputs: map[string]ContentTrustLevel{"locale": ContentTrustTrusted}},
			},
		},
	}
}

func TestContentTrustValidation(t *testing.T) {
	if err := contentTrustDocument().Validate(); err != nil {
		t.Fatalf("valid contentTrust rejected: %v", err)
	}

	tests := []struct {
		name string
		edit func(*Document)
		want string
	}{
		{"version gate", func(doc *Document) { doc.UWS = "1.9.0" }, "contentTrust requires UWS 1.9.1 or later"},
		{"empty root", func(doc *Document) { doc.ContentTrust = &ContentTrust{} }, "must contain at least one declaration"},
		{"unknown source", func(doc *Document) { doc.ContentTrust.SourceDescriptions["missing"] = ContentTrustTrusted }, "unknown sourceDescription"},
		{"unknown operation", func(doc *Document) {
			doc.ContentTrust.Operations["missing"] = &OperationContentTrust{Default: ContentTrustTrusted}
		}, "unknown operationId"},
		{"unknown output", func(doc *Document) {
			doc.ContentTrust.Operations["read_message"].Outputs["missing"] = ContentTrustTrusted
		}, "undeclared operation output"},
		{"empty operation", func(doc *Document) { doc.ContentTrust.Operations["read_message"] = &OperationContentTrust{} }, "must contain default, outputs"},
		{"unknown trigger", func(doc *Document) { doc.ContentTrust.Triggers["missing"] = ContentTrustTrusted }, "unknown triggerId"},
		{"unknown workflow", func(doc *Document) {
			doc.ContentTrust.Workflows["missing"] = &WorkflowContentTrust{Default: ContentTrustTrusted}
		}, "unknown workflowId"},
		{"unknown input", func(doc *Document) { doc.ContentTrust.Workflows["summarize"].Inputs["missing"] = ContentTrustTrusted }, "undeclared workflow input"},
		{"invalid level", func(doc *Document) { doc.ContentTrust.Triggers["incoming_mail"] = "safe" }, "must be unknown, trusted, or untrusted"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			doc := contentTrustDocument()
			test.edit(doc)
			if err := doc.Validate(); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Validate() error = %v, want substring %q", err, test.want)
			}
		})
	}

	doc := contentTrustDocument()
	doc.UWS = "1.10.0"
	if err := doc.Validate(); err != nil {
		t.Fatalf("later compatible version rejected: %v", err)
	}
}

func TestContentTrustJSONAndHCLRoundTrip(t *testing.T) {
	doc := contentTrustDocument()
	doc.ContentTrust.Extensions = map[string]any{"x-review": "security"}
	doc.ContentTrust.Operations["read_message"].Extensions = map[string]any{"x-owner": "mail"}
	doc.ContentTrust.Workflows["summarize"].Extensions = map[string]any{"x-owner": "ai"}

	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	var fromJSON Document
	if err := json.Unmarshal(data, &fromJSON); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(doc.ContentTrust, fromJSON.ContentTrust) {
		t.Fatalf("JSON contentTrust mismatch:\nwant %#v\n got %#v", doc.ContentTrust, fromJSON.ContentTrust)
	}

	hclData, err := doc.MarshalHCL()
	if err != nil {
		t.Fatalf("MarshalHCL: %v", err)
	}
	var fromHCL Document
	if err := fromHCL.UnmarshalHCL(hclData); err != nil {
		t.Fatalf("UnmarshalHCL: %v\n%s", err, hclData)
	}
	if !reflect.DeepEqual(doc.ContentTrust, fromHCL.ContentTrust) {
		t.Fatalf("HCL contentTrust mismatch:\nwant %#v\n got %#v\n%s", doc.ContentTrust, fromHCL.ContentTrust, hclData)
	}
	hclAgain, err := fromHCL.MarshalHCL()
	if err != nil {
		t.Fatal(err)
	}
	if string(hclData) != string(hclAgain) {
		t.Fatalf("HCL output is not deterministic:\n%s\n---\n%s", hclData, hclAgain)
	}
}

func TestContentTrustRejectsUnknownJSONFields(t *testing.T) {
	var trust ContentTrust
	if err := json.Unmarshal([]byte(`{"operations":{},"typo":true}`), &trust); err == nil {
		t.Fatal("expected unknown contentTrust field to be rejected")
	}
	var operationTrust OperationContentTrust
	if err := json.Unmarshal([]byte(`{"default":"trusted","typo":true}`), &operationTrust); err == nil {
		t.Fatal("expected unknown operationContentTrust field to be rejected")
	}
}

func TestContentTrustJSONSchema(t *testing.T) {
	schema := compileUWSSchema(t)
	valid := []byte(`{
  "uws":"1.9.1",
  "info":{"title":"Trust","version":"1.0.0"},
  "operations":[{
    "operationId":"read",
    "x-uws-operation-profile":"test.runtime.1",
    "outputs":{"body":"$response.body"}
  }],
  "contentTrust":{"operations":{"read":{"outputs":{"body":"untrusted"}}}}
}`)
	if err := schema.Validate(decodeJSONValue(t, valid)); err != nil {
		t.Fatalf("1.9.1 contentTrust rejected by schema: %v", err)
	}

	for name, data := range map[string][]byte{
		"bad level": []byte(`{
  "uws":"1.9.1","info":{"title":"Trust","version":"1"},
  "operations":[{"operationId":"read","x-uws-operation-profile":"test.runtime.1"}],
  "contentTrust":{"operations":{"read":{"default":"safe"}}}
}`),
		"empty root": []byte(`{
  "uws":"1.9.1","info":{"title":"Trust","version":"1"},
  "operations":[{"operationId":"read","x-uws-operation-profile":"test.runtime.1"}],
  "contentTrust":{}
}`),
		"empty outputs": []byte(`{
  "uws":"1.9.1","info":{"title":"Trust","version":"1"},
  "operations":[{"operationId":"read","x-uws-operation-profile":"test.runtime.1"}],
  "contentTrust":{"operations":{"read":{"outputs":{}}}}
}`),
	} {
		t.Run(name, func(t *testing.T) {
			if err := schema.Validate(decodeJSONValue(t, data)); err == nil {
				t.Fatal("malformed contentTrust was accepted by schema")
			}
		})
	}
}
