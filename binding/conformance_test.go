package binding

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAcceptedOpenUdonBindingFixture(t *testing.T) {
	root := "../docs/examples/binding/v1"
	manifestBytes, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Revision string                          `json:"upstream_revision"`
		Files    []struct{ Path, SHA256 string } `json:"upstream_files"`
		ShapeSHA string                          `json:"shape_sha256"`
	}
	if json.Unmarshal(manifestBytes, &manifest) != nil || manifest.Revision != "c2f161d762bc9f2217bbf0c34b00cdef64b0f7d0" {
		t.Fatal("invalid pinned corpus")
	}
	for _, f := range manifest.Files {
		data, err := os.ReadFile(filepath.Join(root, f.Path))
		sum := sha256.Sum256(data)
		if err != nil || hex.EncodeToString(sum[:]) != f.SHA256 {
			t.Fatal("changed upstream fixture", f.Path)
		}
	}
	data, err := os.ReadFile(filepath.Join(root, "shapes.json"))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	if hex.EncodeToString(sum[:]) != manifest.ShapeSHA {
		t.Fatal("shape fixture changed")
	}
	table, err := ParseTable(data)
	if err != nil {
		t.Fatal(err)
	}
	var runnable struct {
		Result struct {
			Assessment string
			Ref        struct {
				SourceID string `json:"source_id"`
				SHA      string `json:"source_sha256"`
				Selector string `json:"native_selector"`
				Key      string `json:"operation_key"`
			} `json:"operation_ref"`
			Checks []struct{ Code, Status string }
		}
	}
	data, _ = os.ReadFile(filepath.Join(root, "step-check-runnable.json"))
	if json.Unmarshal(data, &runnable) != nil {
		t.Fatal("bad runnable fixture")
	}
	shape := table.Operations[0]
	if runnable.Result.Assessment != "compatible" || runnable.Result.Ref.SourceID != shape.Source.ID || strings.TrimPrefix(runnable.Result.Ref.SHA, "sha256:") != shape.Source.SHA256 || runnable.Result.Ref.Selector != shape.Selector.Value || runnable.Result.Ref.Key != shape.Selector.Key {
		t.Fatal("upstream operation identity mismatch")
	}
	codes := map[string]bool{}
	for _, check := range runnable.Result.Checks {
		if check.Status == "pass" {
			codes[check.Code] = true
		}
	}
	for _, code := range []string{"operation.exact_match", "mapping.required_inputs", "mapping.output_references", "mapping.input_contract", "authentication.alternative"} {
		if !codes[code] {
			t.Fatal("missing upstream semantic check", code)
		}
	}
	resolver, _ := NewResolver(table)
	req := Request{Binding: Binding{Source: shape.Source, SelectorKind: shape.Selector.Kind, SelectorValue: shape.Selector.Value}, Inputs: []BoundInput{{Location: "query", Name: "page_size", Value: "$inputs.page_size"}}, ExpressionTypes: map[string]Schema{"$inputs.page_size": {Known: true, JSON: json.RawMessage(`{"type":"integer"}`)}}, Security: []SecurityBinding{{Scheme: "apiKeyAuth", CredentialSlot: "project_api_key"}}, OutputReferences: []OutputReference{{Location: "body", Name: "response", Pointer: "#/projects"}}}
	report, err := ValidateBinding(t.Context(), resolver, req)
	if err != nil || report.Outcome != Compatible {
		t.Fatalf("semantic parity: %+v %v", report, err)
	}
	req.Inputs = nil
	report, _ = ValidateBinding(t.Context(), resolver, req)
	if report.Outcome != Incompatible {
		t.Fatal("missing upstream required mapping accepted")
	}
	req.Inputs = []BoundInput{{Location: "query", Name: "page_size", Value: "wrong-type"}}
	report, _ = ValidateBinding(t.Context(), resolver, req)
	if report.Outcome != Incompatible {
		t.Fatal("wrong mapping type accepted")
	}
	req.Inputs = []BoundInput{{Location: "query", Name: "page_size", Value: 10}}
	req.Security = nil
	report, _ = ValidateBinding(t.Context(), resolver, req)
	if report.Outcome != Incompatible {
		t.Fatal("missing authentication alternative accepted")
	}
}
