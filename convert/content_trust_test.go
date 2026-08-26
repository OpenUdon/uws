package convert

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestContentTrustRoundTripsAcrossFormats(t *testing.T) {
	doc := &uws1.Document{
		UWS:  "1.9.1",
		Info: &uws1.Info{Title: "Trust interchange", Version: "1.0.0"},
		Operations: []*uws1.Operation{{
			OperationID: "op", Outputs: map[string]string{"text": "$response.body"},
			Extensions: map[string]any{uws1.ExtensionOperationProfile: "test.runtime.1"},
		}},
		Workflows: []*uws1.Workflow{{
			WorkflowID: "main", Type: uws1.WorkflowTypeSequence,
			Inputs: &uws1.ParamSchema{Type: "object", Properties: map[string]*uws1.ParamSchema{"locale": {Type: "string"}}},
		}},
		Triggers: []*uws1.Trigger{{TriggerID: "event"}},
		ContentTrust: &uws1.ContentTrust{
			Operations: map[string]*uws1.OperationContentTrust{"op": {
				Default: uws1.ContentTrustUnknown,
				Outputs: map[string]uws1.ContentTrustLevel{"text": uws1.ContentTrustUntrusted},
			}},
			Triggers: map[string]uws1.ContentTrustLevel{"event": uws1.ContentTrustUntrusted},
			Workflows: map[string]*uws1.WorkflowContentTrust{"main": {
				Default: uws1.ContentTrustUnknown,
				Inputs:  map[string]uws1.ContentTrustLevel{"locale": uws1.ContentTrustTrusted},
			}},
		},
	}

	jsonData, err := MarshalJSON(doc)
	if err != nil {
		t.Fatal(err)
	}
	yamlData, err := JSONToYAML(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	hclData, err := JSONToHCL(jsonData)
	if err != nil {
		t.Fatal(err)
	}

	for name, decode := range map[string]func() ([]byte, error){
		"yaml": func() ([]byte, error) { return YAMLToJSON(yamlData) },
		"hcl":  func() ([]byte, error) { return HCLToJSON(hclData) },
	} {
		t.Run(name, func(t *testing.T) {
			data, err := decode()
			if err != nil {
				t.Fatal(err)
			}
			var got uws1.Document
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(doc.ContentTrust, got.ContentTrust) {
				t.Fatalf("contentTrust mismatch:\nwant %#v\n got %#v", doc.ContentTrust, got.ContentTrust)
			}
		})
	}

	hclAgain, err := JSONToHCL(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if string(hclData) != string(hclAgain) {
		t.Fatal("contentTrust HCL output is not deterministic")
	}
}
