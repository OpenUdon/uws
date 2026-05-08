package runtimes

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func TestRuntimeSupplementSchemaRuntimeTypes(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)

	for _, typ := range []string{
		RuntimeTypeHTTP, RuntimeTypeSSH, RuntimeTypeCmd, RuntimeTypeFnct,
		RuntimeTypeFileIO, RuntimeTypeSQL, RuntimeTypeS3, RuntimeTypeSMTP,
		RuntimeTypeDNS, RuntimeTypeLDAPS, RuntimeTypeSCP, RuntimeTypeSFTP,
		RuntimeTypeLLM,
	} {
		value := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"`+typ+`"}}`))
		require.NoError(t, schema.Validate(value), "schema should accept runtime type %q", typ)
	}

	plainLDAP := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"ldap"}}`))
	require.Error(t, schema.Validate(plainLDAP))
}

func TestRuntimeSupplementSchemaPayloadFields(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)
	value := decodeRuntimeJSONValue(t, []byte(`{
		"x-uws-runtime": {
			"type": "http",
			"isJson": true,
			"host": "api.example.test",
			"method": "POST",
			"path": "/items",
			"payloadRequired": true,
			"requestMediaType": "application/json",
			"responseMediaType": "application/json",
			"responseStatusCode": 201,
			"command": "echo ok",
			"workingDir": "/tmp",
			"function": "identity",
			"workflow": "child.hcl",
			"arguments": [{"id": "$inputs.id"}],
			"provider": {"name": "api", "serverUrl": "https://api.example.test"},
			"security": [{"name": "api_key", "scheme": {"type": "apiKey", "name": "X-API-Key", "in": "header"}}],
			"queryPars": {"type": "object", "properties": {"limit": {"type": "integer"}}},
			"pathPars": {"type": "object"},
			"headerPars": {"type": "object"},
			"cookiePars": {"type": "object"},
			"payloadPars": {"type": "object"},
			"responseBody": {"type": "object"},
			"responseHeaders": {"type": "object"}
		},
		"x-uws-runtime-config": {
			"provider": {"name": "api"},
			"security": [{"name": "oauth", "scopes": ["read"]}]
		}
	}`))
	require.NoError(t, schema.Validate(value))
}

func compileRuntimeSupplementSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	path := filepath.Join("..", "versions", "runtime.1.0.json")
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	compiler := jsonschema.NewCompiler()
	const resource = "https://github.com/OpenUdon/uws/versions/runtime.1.0.json"
	require.NoError(t, compiler.AddResource(resource, doc))
	schema, err := compiler.Compile(resource)
	require.NoError(t, err)
	return schema
}

func decodeRuntimeJSONValue(t *testing.T, data []byte) any {
	t.Helper()
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	return value
}
