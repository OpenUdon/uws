package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var publishedCoreVersions = []string{
	"1.0.0", "1.1.0", "1.1.1", "1.2.0", "1.3.0", "1.4.0",
	"1.5.0", "1.6.0", "1.7.0", "1.8.0", "1.9.0", "1.9.1", "1.9.2",
	"1.10.0",
}

func TestPublishedCoreSchemaAndSemanticCompatibilityCorpus(t *testing.T) {
	for _, version := range publishedCoreVersions {
		t.Run(version, func(t *testing.T) {
			if err := validateCompatibilityFixture(t, version, nil); err != nil {
				t.Fatalf("valid baseline fixture failed schema or semantic validation: %v", err)
			}
		})
	}
}

func TestVersionedCompatibilityBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		mutate    func(map[string]any)
		wantError string
	}{
		{
			name:    "timeout before 1.1 is rejected",
			version: "1.0.0",
			mutate: func(doc map[string]any) {
				doc["operations"].([]any)[0].(map[string]any)["timeout"] = 5
			},
			wantError: "timeout",
		},
		{
			name:    "timeout at 1.1 is accepted",
			version: "1.1.0",
			mutate: func(doc map[string]any) {
				doc["operations"].([]any)[0].(map[string]any)["timeout"] = 5
			},
		},
		{
			name:    "step inputs before 1.5 are rejected",
			version: "1.4.0",
			mutate: func(doc map[string]any) {
				doc["workflows"].([]any)[0].(map[string]any)["steps"].([]any)[0].(map[string]any)["inputs"] = map[string]any{"id": "123"}
			},
			wantError: "inputs",
		},
		{
			name:    "step inputs at 1.5 are accepted",
			version: "1.5.0",
			mutate: func(doc map[string]any) {
				doc["workflows"].([]any)[0].(map[string]any)["steps"].([]any)[0].(map[string]any)["inputs"] = map[string]any{"id": "123"}
			},
		},
		{
			name:    "content trust before 1.9.1 is rejected",
			version: "1.9.0",
			mutate: func(doc map[string]any) {
				doc["contentTrust"] = map[string]any{"operations": map[string]any{"op": map[string]any{"default": "untrusted"}}}
			},
			wantError: "contentTrust",
		},
		{
			name:    "content trust at 1.9.1 is accepted",
			version: "1.9.1",
			mutate: func(doc map[string]any) {
				doc["contentTrust"] = map[string]any{"operations": map[string]any{"op": map[string]any{"default": "untrusted"}}}
			},
		},
		{
			name:    "reference child blocks at 1.9.1 retain prior behavior",
			version: "1.9.1",
			mutate: func(doc map[string]any) {
				step := doc["workflows"].([]any)[0].(map[string]any)["steps"].([]any)[0].(map[string]any)
				step["steps"] = []any{map[string]any{"stepId": "nested", "operationRef": "op"}}
			},
		},
		{
			name:    "reference child blocks at 1.9.2 are rejected",
			version: "1.9.2",
			mutate: func(doc map[string]any) {
				step := doc["workflows"].([]any)[0].(map[string]any)["steps"].([]any)[0].(map[string]any)
				step["steps"] = []any{map[string]any{"stepId": "nested", "operationRef": "op"}}
			},
			wantError: "steps",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateCompatibilityFixture(t, test.version, test.mutate)
			if test.wantError == "" {
				if err != nil {
					t.Fatalf("expected acceptance, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(test.wantError)) {
				t.Fatalf("expected error containing %q, got %v", test.wantError, err)
			}
		})
	}
}

func TestUnpublishedVersionsRequireTheirExactSchemaArtifact(t *testing.T) {
	for _, version := range []string{"1.9.2-rc.1", "1.9.3"} {
		t.Run(version, func(t *testing.T) {
			err := validateCompatibilityFixture(t, version, nil)
			if err == nil || !strings.Contains(err.Error(), version+".json") {
				t.Fatalf("expected missing exact schema for %q, got %v", version, err)
			}
		})
	}
}

func validateCompatibilityFixture(t *testing.T, version string, mutate func(map[string]any)) error {
	t.Helper()
	doc := map[string]any{
		"uws":                version,
		"info":               map[string]any{"title": "version compatibility", "version": "1.0.0"},
		"sourceDescriptions": []any{map[string]any{"name": "api", "url": "./api.yaml", "type": "openapi"}},
		"operations":         []any{map[string]any{"operationId": "op", "sourceDescription": "api", "openapiOperationId": "get"}},
		"workflows":          []any{map[string]any{"workflowId": "main", "type": "sequence", "steps": []any{map[string]any{"stepId": "call", "operationRef": "op"}}}},
	}
	if mutate != nil {
		mutate(doc)
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	path := filepath.Join(t.TempDir(), fmt.Sprintf("compat-%s.uws.json", strings.ReplaceAll(version, ".", "_")))
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	validated, err := ValidateDocumentFile(path)
	if err != nil {
		return err
	}
	if validated == nil || validated.UWS != version {
		return fmt.Errorf("validated version = %v, want %q", validated, version)
	}
	return nil
}
