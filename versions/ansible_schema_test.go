package versions

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func TestAnsibleArgspecSchemaValidDocuments(t *testing.T) {
	schema := compileAnsibleArgspecSchema(t)

	require.NoError(t, validateAnsibleArgspec(t, schema, validAnsibleArgspec()))

	withReturns := validAnsibleArgspec()
	module := withReturns["modules"].(map[string]any)["ansible.builtin.apt"].(map[string]any)
	module["returns"] = map[string]any{
		"stdout": map[string]any{"type": "str", "returned": "success"},
	}
	require.NoError(t, validateAnsibleArgspec(t, schema, withReturns))
}

func TestAnsibleArgspecSchemaRejectsInvalidDocuments(t *testing.T) {
	schema := compileAnsibleArgspecSchema(t)

	wrongProfile := validAnsibleArgspec()
	wrongProfile["argspec"] = "uws.ansible.2.0"
	require.Error(t, validateAnsibleArgspec(t, schema, wrongProfile), "argspec const must be enforced")

	badFQCN := validAnsibleArgspec()
	badFQCN["modules"] = map[string]any{
		"apt": map[string]any{"parameters": map[string]any{}},
	}
	require.Error(t, validateAnsibleArgspec(t, schema, badFQCN), "module keys must be FQCNs")

	missingParameters := validAnsibleArgspec()
	missingParameters["modules"] = map[string]any{
		"ansible.builtin.apt": map[string]any{"shortDescription": "no parameters"},
	}
	require.Error(t, validateAnsibleArgspec(t, schema, missingParameters), "parameters is required")

	badType := validAnsibleArgspec()
	badType["modules"].(map[string]any)["ansible.builtin.apt"].(map[string]any)["parameters"] = map[string]any{
		"name": map[string]any{"type": "stringly"},
	}
	require.Error(t, validateAnsibleArgspec(t, schema, badType), "parameter type enum must be enforced")

	extraField := validAnsibleArgspec()
	extraField["secrets"] = map[string]any{"token": "nope"}
	require.Error(t, validateAnsibleArgspec(t, schema, extraField), "unknown top-level fields must be rejected")
}

func compileAnsibleArgspecSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	data, err := os.ReadFile("ansible.1.0.json")
	require.NoError(t, err)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	compiler := jsonschema.NewCompiler()
	require.NoError(t, compiler.AddResource("ansible.1.0.json", doc))
	schema, err := compiler.Compile("ansible.1.0.json")
	require.NoError(t, err)
	return schema
}

func validateAnsibleArgspec(t *testing.T, schema *jsonschema.Schema, document map[string]any) error {
	t.Helper()
	data, err := json.Marshal(document)
	require.NoError(t, err)
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	return schema.Validate(value)
}

func validAnsibleArgspec() map[string]any {
	return map[string]any{
		"argspec":    "uws.ansible.1.0",
		"collection": "ansible.builtin",
		"info": map[string]any{
			"source": "ansible-doc --json",
		},
		"modules": map[string]any{
			"ansible.builtin.apt": map[string]any{
				"shortDescription": "Manages apt packages",
				"parameters": map[string]any{
					"name":  map[string]any{"type": "list", "elements": "str"},
					"state": map[string]any{"type": "str", "choices": []any{"absent", "present", "latest"}},
				},
			},
			"ansible.builtin.service": map[string]any{
				"parameters": map[string]any{
					"name":  map[string]any{"type": "str", "required": true},
					"state": map[string]any{"type": "str", "choices": []any{"restarted", "started", "stopped"}},
				},
			},
		},
	}
}
