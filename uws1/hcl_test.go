package uws1

import (
	"bytes"
	"strings"
	"testing"

	"github.com/genelet/horizon/dethcl"
)

func TestDocumentHCLInterfacesDecodeNativeHCL(t *testing.T) {
	hclData := []byte(`
uws = "1.0.0"
variables = {
  __dollar__root = "kept"
  _ref = "#/components/schemas/Root"
}

info {
  title = "HCL Document"
  summary = "Line 1\\nLine 2"
  version = "1.0.0"
  extensions {
    x-owner = "team"
  }
}

operation "op1" {
  sourceDescription = "api"
  openapiOperationId = "getOp"
  request = {
    body = {
      __dollar__request = "value"
      _ref = "#/components/schemas/Body"
    }
  }
  extensions {
    x-uws-operation-profile = "uws.runtime.1.0"
    x-uws-runtime = {
      type = "fnct"
      function = "identity"
    }
  }
}
`)

	var doc Document
	if err := dethcl.Unmarshal(hclData, &doc); err != nil {
		t.Fatalf("dethcl.Unmarshal failed: %v", err)
	}

	if got := doc.Variables["$root"]; got != "kept" {
		t.Fatalf("expected restored $root variable, got %#v", doc.Variables)
	}
	if got := doc.Variables["$ref"]; got != "#/components/schemas/Root" {
		t.Fatalf("expected restored $ref variable, got %#v", doc.Variables)
	}
	if doc.Info == nil || doc.Info.Summary != "Line 1\nLine 2" {
		t.Fatalf("expected summary newline unescape, got %#v", doc.Info)
	}
	if doc.Info.Extensions["x-owner"] != "team" {
		t.Fatalf("expected info extension, got %#v", doc.Info.Extensions)
	}
	body, ok := doc.Operations[0].Request["body"].(map[string]any)
	if !ok {
		t.Fatalf("expected request body map, got %#v", doc.Operations[0].Request["body"])
	}
	if body["$request"] != "value" || body["$ref"] != "#/components/schemas/Body" {
		t.Fatalf("expected restored request keys, got %#v", body)
	}
	runtime, ok := doc.Operations[0].Extensions["x-uws-runtime"].(map[string]any)
	if !ok || runtime["type"] != "fnct" || runtime["function"] != "identity" {
		t.Fatalf("expected runtime extension payload, got %#v", doc.Operations[0].Extensions)
	}
}

func TestDocumentHCLInterfacesDecodeEmbeddedDocumentBlock(t *testing.T) {
	type envelope struct {
		Name     string   `hcl:"name"`
		Document Document `hcl:"document,block"`
	}

	hclData := []byte(`
name = "package"

document {
  uws = "1.0.0"
  variables = {
    __dollar__root = "kept"
  }
  info {
    title = "Wrapped"
    version = "1.0.0"
  }
  operation "op1" {
    sourceDescription = "api"
    openapiOperationId = "getOp"
  }
}
`)

	var got envelope
	if err := dethcl.Unmarshal(hclData, &got); err != nil {
		t.Fatalf("dethcl.Unmarshal wrapper failed: %v", err)
	}
	if got.Name != "package" {
		t.Fatalf("expected wrapper name, got %q", got.Name)
	}
	if got.Document.Info == nil || got.Document.Info.Title != "Wrapped" {
		t.Fatalf("expected embedded document info, got %#v", got.Document.Info)
	}
	if got.Document.Variables["$root"] != "kept" {
		t.Fatalf("expected embedded document key restoration, got %#v", got.Document.Variables)
	}
}

func TestDocumentHCLInterfacesDecodeNestedStructsDirectly(t *testing.T) {
	operationHCL := []byte(`
sourceDescription = "api"
openapiOperationId = "getOp"
description = "Line 1\\nLine 2"
request = {
  body = {
    __dollar__request = "value"
    _ref = "#/components/schemas/Body"
  }
}

successCriterion {
  condition = "true"
  extensions {
    x-owner = "team"
  }
}

extensions {
  x-uws-runtime = {
    type = "fnct"
    __dollar__expr = "$inputs.name"
  }
}
`)

	var op Operation
	if err := dethcl.Unmarshal(operationHCL, &op, "op1"); err != nil {
		t.Fatalf("dethcl.Unmarshal operation failed: %v", err)
	}
	if op.OperationID != "op1" {
		t.Fatalf("expected label to populate operationId, got %q", op.OperationID)
	}
	if op.Description != "Line 1\nLine 2" {
		t.Fatalf("expected operation description newline unescape, got %q", op.Description)
	}
	body, ok := op.Request["body"].(map[string]any)
	if !ok {
		t.Fatalf("expected request body map, got %#v", op.Request["body"])
	}
	if body["$request"] != "value" || body["$ref"] != "#/components/schemas/Body" {
		t.Fatalf("expected operation request key restoration, got %#v", body)
	}
	runtime, ok := op.Extensions["x-uws-runtime"].(map[string]any)
	if !ok || runtime["type"] != "fnct" || runtime["$expr"] != "$inputs.name" {
		t.Fatalf("expected operation extension key restoration, got %#v", op.Extensions)
	}
	if len(op.SuccessCriteria) != 1 || op.SuccessCriteria[0].Extensions["x-owner"] != "team" {
		t.Fatalf("expected criterion extension from child unmarshaller, got %#v", op.SuccessCriteria)
	}
}

func TestStepHCLRestoresInputDollarKeysDirectly(t *testing.T) {
	stepHCL := []byte(`
operationRef = "op1"
inputs = {
  __dollar__expr = "$inputs.name"
  nested = {
    _ref = "#/components/schemas/Input"
  }
}
`)

	var step Step
	if err := dethcl.Unmarshal(stepHCL, &step, "s1"); err != nil {
		t.Fatalf("dethcl.Unmarshal step failed: %v", err)
	}
	if step.StepID != "s1" {
		t.Fatalf("expected label to populate stepId, got %q", step.StepID)
	}
	if got := step.Inputs["$expr"]; got != "$inputs.name" {
		t.Fatalf("expected restored step input $expr, got %#v", step.Inputs)
	}
	nested, ok := step.Inputs["nested"].(map[string]any)
	if !ok || nested["$ref"] != "#/components/schemas/Input" {
		t.Fatalf("expected restored nested step input $ref, got %#v", step.Inputs["nested"])
	}
}

func TestDocumentHCLInterfacesMarshalThroughDethcl(t *testing.T) {
	doc := &Document{
		UWS: "1.0.0",
		Info: &Info{
			Title:   "Marshal",
			Summary: "Line 1\nLine 2",
			Version: "1.0.0",
		},
		Variables: map[string]any{"$root": "kept"},
		Operations: []*Operation{
			{
				OperationID:        "op1",
				SourceDescription:  "api",
				OpenAPIOperationID: "getOp",
				Request:            map[string]any{"body": map[string]any{"$request": "value"}},
				Extensions:         map[string]any{"x-owner": "team"},
			},
		},
	}

	hclData, err := dethcl.Marshal(doc)
	if err != nil {
		t.Fatalf("dethcl.Marshal failed: %v", err)
	}
	hcl := string(hclData)
	for _, want := range []string{"__dollar__root", "__dollar__request", "extensions", "x-owner", `Line 1\\nLine 2`} {
		if !strings.Contains(hcl, want) {
			t.Fatalf("HCL output missing %q:\n%s", want, hcl)
		}
	}

	if _, ok := doc.Variables["$root"]; !ok {
		t.Fatalf("MarshalHCL mutated original variable keys: %#v", doc.Variables)
	}
	body := doc.Operations[0].Request["body"].(map[string]any)
	if _, ok := body["$request"]; !ok {
		t.Fatalf("MarshalHCL mutated original request keys: %#v", body)
	}
}

func TestDocumentMarshalHCLIsDeterministic(t *testing.T) {
	doc := &Document{
		UWS: "1.7.0",
		Info: &Info{
			Title:      "Deterministic HCL",
			Version:    "1.0.0",
			Extensions: map[string]any{"x-zulu": 26, "x-alpha": 1},
		},
		Variables: map[string]any{"zulu": "last", "alpha": "first"},
		Operations: []*Operation{{
			OperationID: "op1",
			Request: map[string]any{"body": map[string]any{
				"zulu":  26,
				"alpha": map[string]any{"two": 2, "one": 1},
			}},
			Extensions: map[string]any{
				ExtensionOperationProfile: "test.profile.1.0",
				"x-zulu":                  map[string]any{"two": 2, "one": 1},
				"x-alpha":                 true,
			},
		}},
	}

	first, err := doc.MarshalHCL()
	if err != nil {
		t.Fatalf("MarshalHCL failed: %v", err)
	}
	for i := 0; i < 100; i++ {
		got, err := doc.MarshalHCL()
		if err != nil {
			t.Fatalf("MarshalHCL failed: %v", err)
		}
		if !bytes.Equal(got, first) {
			t.Fatalf("MarshalHCL output changed between calls:\nfirst:\n%s\ngot:\n%s", first, got)
		}
	}

	hcl := string(first)
	for before, after := range map[string]string{
		`alpha = "first"`: `zulu = "last"`,
		"one = 1":         "two = 2",
	} {
		beforeIndex, afterIndex := strings.Index(hcl, before), strings.Index(hcl, after)
		if beforeIndex < 0 || afterIndex < 0 || beforeIndex >= afterIndex {
			t.Fatalf("expected %q before %q in deterministic HCL:\n%s", before, after, hcl)
		}
	}
	if alphaIndex, zuluIndex := strings.Index(hcl, "x-alpha = true"), strings.LastIndex(hcl, "x-zulu"); alphaIndex < 0 || zuluIndex < 0 || alphaIndex >= zuluIndex {
		t.Fatalf("expected operation extension x-alpha before x-zulu in deterministic HCL:\n%s", hcl)
	}
}
