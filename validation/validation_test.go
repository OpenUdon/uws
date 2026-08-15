package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/convert"
	"github.com/OpenUdon/uws/schemas"
	"github.com/OpenUdon/uws/uws1"
)

func TestValidateFileAcceptsJSONAndYAML(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(schema, []byte(`{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["name"],
  "properties": {"name": {"type": "string"}}
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	jsonDoc := filepath.Join(dir, "doc.json")
	if err := os.WriteFile(jsonDoc, []byte(`{"name":"uws"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	yamlDoc := filepath.Join(dir, "doc.yaml")
	if err := os.WriteFile(yamlDoc, []byte("name: uws\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, doc := range []string{jsonDoc, yamlDoc} {
		if err := ValidateFile(schema, doc); err != nil {
			t.Fatalf("ValidateFile(%s) returned error: %v", doc, err)
		}
	}
}

func TestValidateFileRejectsInvalidAndUnsupportedDocuments(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(schema, []byte(`{"type":"object","required":["name"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	invalid := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(invalid, []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateFile(schema, invalid); err == nil {
		t.Fatalf("expected schema validation failure")
	}
	unsupported := filepath.Join(dir, "doc.txt")
	if err := os.WriteFile(unsupported, []byte(`{"name":"uws"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ValidateFile(schema, unsupported); err == nil || !strings.Contains(err.Error(), "unsupported document extension") {
		t.Fatalf("expected unsupported extension error, got %v", err)
	}
}

func TestValidateFileAcceptsVersionedUWSDocuments(t *testing.T) {
	for _, tc := range []struct {
		name    string
		version string
		body    string
	}{
		{
			name:    "typed source",
			version: "1.2.0",
			body: `{
  "uws": "1.2.0",
  "info": {"title": "typed source", "version": "1.0.0"},
  "sourceDescriptions": [
    {"name": "gmail", "url": "google-discovery/gmail.json", "type": "google-discovery"}
  ],
  "operations": [
    {"operationId": "send", "sourceDescription": "gmail", "sourceOperationId": "gmail_users_messages_send"}
  ],
  "workflows": [
    {"workflowId": "main", "type": "sequence", "steps": [{"stepId": "send", "operationRef": "send"}]}
  ]
}`,
		},
		{
			name:    "async source",
			version: "1.3.0",
			body: `{
  "uws": "1.3.0",
  "info": {"title": "async source", "version": "1.0.0"},
  "sourceDescriptions": [
    {"name": "billing_events", "url": "asyncapi/billing-events.yaml", "type": "asyncapi"}
  ],
  "operations": [
    {"operationId": "publish_invoice", "sourceDescription": "billing_events", "sourceOperationId": "publishInvoice"}
  ],
  "workflows": [
    {"workflowId": "main", "type": "sequence", "steps": [{"stepId": "publish_invoice", "operationRef": "publish_invoice"}]}
  ]
}`,
		},
		{
			name:    "graphql source",
			version: "1.4.0",
			body: `{
  "uws": "1.4.0",
  "info": {"title": "graphql source", "version": "1.0.0"},
  "sourceDescriptions": [
    {"name": "catalog", "url": "graphql/catalog.graphql", "type": "graphql"}
  ],
  "operations": [
    {"operationId": "query_product", "sourceDescription": "catalog", "sourceOperationId": "ProductQuery"}
  ],
  "workflows": [
    {"workflowId": "main", "type": "sequence", "steps": [{"stepId": "query_product", "operationRef": "query_product"}]}
  ]
}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			doc := filepath.Join(dir, "workflow.uws.json")
			if err := os.WriteFile(doc, []byte(tc.body), 0o644); err != nil {
				t.Fatal(err)
			}
			if err := ValidateFile(schemas.PathForVersion(dir, tc.version), doc); err != nil {
				t.Fatalf("ValidateFile returned error: %v", err)
			}
		})
	}
}

func TestLoadDocumentFileReadsJSONYAMLAndHCL(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"workflow.uws.json": `{"uws":"1.0.0","info":{"title":"json","version":"1.0.0"},"operations":[{"operationId":"op","x-uws-operation-profile":"uws.runtime.1.0"}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"op","operationRef":"op"}]}]}`,
		"workflow.uws.yaml": "uws: 1.0.0\ninfo:\n  title: yaml\n  version: 1.0.0\noperations:\n  - operationId: op\n    x-uws-operation-profile: uws.runtime.1.0\nworkflows:\n  - workflowId: main\n    type: sequence\n    steps:\n      - stepId: op\n        operationRef: op\n",
		"workflow.uws.hcl": `uws = "1.0.0"
info {
  title = "hcl"
  version = "1.0.0"
}
operation "op" {
  extensions {
    x-uws-operation-profile = "uws.runtime.1.0"
  }
}
workflow "main" {
  type = "sequence"
  step "op" {
    operationRef = "op"
  }
}
`,
	}
	for name, body := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		doc, err := LoadDocumentFile(path)
		if err != nil {
			t.Fatalf("LoadDocumentFile(%s) returned error: %v", name, err)
		}
		if doc == nil || strings.TrimSpace(doc.UWS) == "" {
			t.Fatalf("LoadDocumentFile(%s) returned empty document", name)
		}
	}
}

func TestValidateDocumentFileRunsSchemaAndSemanticValidation(t *testing.T) {
	dir := t.TempDir()
	valid := filepath.Join(dir, "workflow.uws.yaml")
	if err := os.WriteFile(valid, []byte("uws: 1.0.0\ninfo:\n  title: valid\n  version: 1.0.0\noperations:\n  - operationId: op\n    x-uws-operation-profile: uws.runtime.1.0\nworkflows:\n  - workflowId: main\n    type: sequence\n    steps:\n      - stepId: op\n        operationRef: op\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateDocumentFile(valid); err != nil {
		t.Fatalf("ValidateDocumentFile(valid) returned error: %v", err)
	}

	schemaInvalid := filepath.Join(dir, "schema-invalid.uws.json")
	if err := os.WriteFile(schemaInvalid, []byte(`{"uws":"1.0.0","info":{"title":"bad","version":"1.0.0"},"operations":"nope"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateDocumentFile(schemaInvalid); err == nil {
		t.Fatalf("expected schema validation failure")
	}

	semanticInvalid := filepath.Join(dir, "semantic-invalid.uws.json")
	if err := os.WriteFile(semanticInvalid, []byte(`{"uws":"1.0.0","info":{"title":"bad","version":"1.0.0"},"operations":[{"operationId":"dup","x-uws-operation-profile":"uws.runtime.1.0"},{"operationId":"dup","x-uws-operation-profile":"uws.runtime.1.0"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateDocumentFile(semanticInvalid); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("expected semantic duplicate failure, got %v", err)
	}
}

func TestValidateDocumentFileAcceptsOperationRefStepFromJSONMarshal(t *testing.T) {
	dir := t.TempDir()
	doc := &uws1.Document{
		UWS:  "1.4.0",
		Info: &uws1.Info{Title: "marshaled", Version: "1.0.0"},
		Operations: []*uws1.Operation{{
			OperationID: "op",
			Extensions:  map[string]any{uws1.ExtensionOperationProfile: "uws.runtime.1.0"},
		}},
		Workflows: []*uws1.Workflow{{
			WorkflowID: "main",
			Type:       uws1.WorkflowTypeSequence,
			Steps: []*uws1.Step{{
				StepID:       "op",
				OperationRef: "op",
			}},
		}},
	}
	data, err := convert.MarshalJSONIndent(doc, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "workflow.uws.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ValidateDocumentFile(path); err != nil {
		t.Fatalf("ValidateDocumentFile returned error: %v\n%s", err, data)
	}
}

func TestCollectArtifactFilesReturnsSortedUWSArtifacts(t *testing.T) {
	dir := t.TempDir()
	for _, rel := range []string{"b.uws.yaml", "nested/a.uws.json", "nested/c.uws.yml", "nested/ignore.yaml", "workflow.uws.hcl"} {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("{}"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	files, err := CollectArtifactFiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(dir, "b.uws.yaml"),
		filepath.Join(dir, "nested", "a.uws.json"),
		filepath.Join(dir, "nested", "c.uws.yml"),
	}
	if strings.Join(files, "\n") != strings.Join(want, "\n") {
		t.Fatalf("CollectArtifactFiles = %#v, want %#v", files, want)
	}
	for _, rel := range []string{"workflow.uws.json", "workflow.uws.yaml", "workflow.uws.yml"} {
		if !IsArtifactFile(rel) {
			t.Fatalf("IsArtifactFile(%q) = false", rel)
		}
	}
	if IsArtifactFile("workflow.uws.hcl") || IsArtifactFile("workflow.yaml") {
		t.Fatalf("IsArtifactFile accepted non-artifact path")
	}
}
