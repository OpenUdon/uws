package ansiblemodulecall

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func TestAnsibleModuleCallSupplementSchemaAcceptsMinimalPayload(t *testing.T) {
	schema := compileAnsibleModuleCallSupplementSchema(t)
	value := decodeAnsibleModuleCallJSONValue(t, []byte(`{
		"x-uws-ansible-module": {
			"module": "ansible.builtin.apt"
		}
	}`))
	require.NoError(t, schema.Validate(value))
}

func TestAnsibleModuleCallSupplementSchemaAcceptsArgspecReference(t *testing.T) {
	schema := compileAnsibleModuleCallSupplementSchema(t)
	value := decodeAnsibleModuleCallJSONValue(t, []byte(`{
		"x-uws-ansible-module": {
			"module": "community.postgresql.postgresql_db",
			"argspec": {
				"sourceId": "postgres",
				"url": "./postgres.argspec.json",
				"collection": "community.postgresql"
			}
		}
	}`))
	require.NoError(t, schema.Validate(value))
}

func TestAnsibleModuleCallSupplementSchemaRejectsMissingModule(t *testing.T) {
	schema := compileAnsibleModuleCallSupplementSchema(t)
	value := decodeAnsibleModuleCallJSONValue(t, []byte(`{"x-uws-ansible-module": {}}`))
	require.Error(t, schema.Validate(value))
}

func TestAnsibleModuleCallSupplementSchemaRejectsNonCanonicalFQCNs(t *testing.T) {
	schema := compileAnsibleModuleCallSupplementSchema(t)
	for name, module := range map[string]string{
		"mixed case":     "Ansible.Builtin.Apt",
		"four segments":  "ansible.builtin.apt.extra",
		"two segments":   "ansible.apt",
		"short name":     "apt",
		"hyphen segment": "ansible.builtin.apt-get",
	} {
		value := decodeAnsibleModuleCallJSONValue(t, []byte(`{
			"x-uws-ansible-module": {
				"module": "`+module+`"
			}
		}`))
		require.Error(t, schema.Validate(value), "module %s (%s) should be rejected", module, name)
	}
}

func TestAnsibleModuleCallSupplementSchemaRejectsUnknownFields(t *testing.T) {
	schema := compileAnsibleModuleCallSupplementSchema(t)
	value := decodeAnsibleModuleCallJSONValue(t, []byte(`{
		"x-uws-ansible-module": {
			"module": "ansible.builtin.apt",
			"connection": "ssh"
		}
	}`))
	require.Error(t, schema.Validate(value))
}

func TestAnsibleModuleCallSupplementSchemaRejectsUnknownTopLevelFields(t *testing.T) {
	schema := compileAnsibleModuleCallSupplementSchema(t)
	value := decodeAnsibleModuleCallJSONValue(t, []byte(`{
		"x-uws-ansible-module": {
			"module": "ansible.builtin.apt"
		},
		"x-uws-runtime-config": {}
	}`))
	require.Error(t, schema.Validate(value))
}

func TestAnsibleModuleCallSupplementSchemaProfileNameMatchesConstant(t *testing.T) {
	require.Equal(t, "uws.ansible-module-call.1.0", ProfileName)
}

func compileAnsibleModuleCallSupplementSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	data := readAnsibleModuleCallSupplementSchema(t)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	compiler := jsonschema.NewCompiler()
	const resource = "https://github.com/OpenUdon/uws/versions/ansible-module-call.1.0.json"
	require.NoError(t, compiler.AddResource(resource, doc))
	schema, err := compiler.Compile(resource)
	require.NoError(t, err)
	return schema
}

func readAnsibleModuleCallSupplementSchema(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "versions", "ansible-module-call.1.0.json")
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func decodeAnsibleModuleCallJSONValue(t *testing.T, data []byte) any {
	t.Helper()
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	return value
}

func TestAnsibleModuleCallSupplementSchemaRequiredTopLevelExtension(t *testing.T) {
	var doc struct {
		Required []string `json:"required"`
	}
	require.NoError(t, json.Unmarshal(readAnsibleModuleCallSupplementSchema(t), &doc))
	require.Equal(t, []string{ExtensionAnsibleModule}, doc.Required)
}
