package uws1

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	latestUWSSchemaPath     = "../versions/1.1.0.json"
	latestUWSSchemaResource = "versions/1.1.0.json"
	latestUWSSpecPath       = "../versions/1.1.0.md"
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
