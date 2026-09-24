package uws1

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"
)

const (
	latestUWSSchemaPath     = "../versions/1.11.0.json"
	latestUWSSchemaResource = "versions/1.11.0.json"
	latestUWSSpecPath       = "../versions/1.11.0.md"
)

func loadSchemaDoc(t *testing.T) map[string]any {
	t.Helper()
	data, err := os.ReadFile(latestUWSSchemaPath)
	require.NoError(t, err)
	var schema map[string]any
	require.NoError(t, json.Unmarshal(data, &schema))
	return schema
}

func loadSpecMarkdown(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(latestUWSSpecPath)
	require.NoError(t, err)
	return string(data)
}

// schemaPropertyNames returns the property names of a $def. When defName is
// empty, it returns the root document's properties.
func schemaPropertyNames(t *testing.T, schema map[string]any, defName string) []string {
	t.Helper()
	obj := schema
	if defName != "" {
		defs, ok := schema["$defs"].(map[string]any)
		require.True(t, ok, "schema has no $defs")
		entry, ok := defs[defName].(map[string]any)
		require.True(t, ok, "schema $defs has no %q", defName)
		obj = entry
	}
	props, ok := obj["properties"].(map[string]any)
	require.Truef(t, ok, "schema %q has no properties object", fmt.Sprintf("$defs/%s", defName))
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func dropExtensionKeys(keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if strings.HasPrefix(k, "x-") {
			continue
		}
		out = append(out, k)
	}
	return out
}

func TestUWS110LanguageNeutralConformanceVectors(t *testing.T) {
	data, err := os.ReadFile("../testdata/conformance/1.10.0.json")
	require.NoError(t, err)
	var fixture struct {
		UWS   string `json:"uws"`
		Cases []struct {
			ID       string          `json:"id"`
			Area     string          `json:"area"`
			Input    json.RawMessage `json:"input"`
			Expected json.RawMessage `json:"expected"`
		} `json:"cases"`
	}
	require.NoError(t, json.Unmarshal(data, &fixture))
	require.Equal(t, "1.10.0", fixture.UWS)
	require.NotEmpty(t, fixture.Cases)
	seen := make(map[string]bool, len(fixture.Cases))
	for _, testCase := range fixture.Cases {
		require.NotEmpty(t, testCase.ID)
		require.NotEmpty(t, testCase.Area)
		require.NotEmpty(t, testCase.Input)
		require.NotEmpty(t, testCase.Expected)
		require.False(t, seen[testCase.ID], "duplicate conformance case %q", testCase.ID)
		seen[testCase.ID] = true
	}
}

func TestUWS111GrammarFixturePassesSchemaAndSemanticValidation(t *testing.T) {
	data, err := os.ReadFile("../testdata/grammar/1.11.0.json")
	require.NoError(t, err)

	var doc Document
	require.NoError(t, json.Unmarshal(data, &doc))
	require.NoError(t, doc.Validate())
	require.NoError(t, compileUWSSchema(t).Validate(decodeJSONValue(t, data)))
}

func TestLatestUWSSchemaIsHighestPublishedCoreVersion(t *testing.T) {
	paths, err := filepath.Glob("../versions/1.*.json")
	require.NoError(t, err)
	require.NotEmpty(t, paths)

	var latest string
	for _, path := range paths {
		version := strings.TrimSuffix(filepath.Base(path), ".json")
		if !semver.IsValid("v" + version) {
			continue
		}
		if latest == "" || semver.Compare("v"+version, "v"+latest) > 0 {
			latest = version
		}
	}
	require.NotEmpty(t, latest)
	require.Equal(t, latest, strings.TrimSuffix(filepath.Base(latestUWSSchemaPath), ".json"))
}

func TestLatestUWSSchemaIDMatchesPublishedPath(t *testing.T) {
	schema := loadSchemaDoc(t)
	require.Equal(t, "https://github.com/OpenUdon/uws/versions/1.11.0.json", schema["$id"])
	require.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"])
}

func TestPublishedUWSVersionRegistryMatchesSchemaArtifacts(t *testing.T) {
	paths, err := filepath.Glob("../versions/1.*.json")
	require.NoError(t, err)

	artifacts := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		version := strings.TrimSuffix(filepath.Base(path), ".json")
		if semver.IsValid("v" + version) {
			artifacts[version] = struct{}{}
		}
	}
	require.NotEmpty(t, artifacts)
	require.Equal(t, artifacts, publishedUWSVersions,
		"publishedUWSVersions must match the core schema artifacts exactly")
}
