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
		RuntimeTypeSSH, RuntimeTypeCmd, RuntimeTypeFnct,
		RuntimeTypeFileIO, RuntimeTypeSQL, RuntimeTypeS3, RuntimeTypeSMTP,
		RuntimeTypeDNS, RuntimeTypeLDAPS, RuntimeTypeSCP, RuntimeTypeSFTP,
		RuntimeTypeLLM,
	} {
		value := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"`+typ+`"}}`))
		require.NoError(t, schema.Validate(value), "schema should accept runtime type %q", typ)
	}

	plainLDAP := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"ldap"}}`))
	require.Error(t, schema.Validate(plainLDAP))

	httpRuntime := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"http"}}`))
	require.Error(t, schema.Validate(httpRuntime))
}

func TestRuntimeSupplementSchemaSlimPayloadFields(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)
	value := decodeRuntimeJSONValue(t, []byte(`{
		"x-uws-runtime": {
			"type": "cmd",
			"command": "echo ok",
			"workingDir": "/tmp",
			"function": "identity",
			"workflow": "child.hcl",
			"arguments": [{"id": "$inputs.id"}]
		}
	}`))
	require.NoError(t, schema.Validate(value))
}

func TestRuntimeSupplementSchemaAcceptsFnctAndLLMSelectors(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)

	fnct := decodeRuntimeJSONValue(t, []byte(`{
		"x-uws-runtime": {
			"type": "fnct",
			"function": "render"
		}
	}`))
	require.NoError(t, schema.Validate(fnct))

	llm := decodeRuntimeJSONValue(t, []byte(`{
		"x-uws-runtime": {
			"type": "llm",
			"function": "summarize",
			"arguments": [{"prompt": "$inputs.text"}]
		}
	}`))
	require.NoError(t, schema.Validate(llm))
}

func TestRuntimeSupplementSchemaRejectsRemovedFields(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)

	for _, field := range []string{
		"isJson",
		"host",
		"method",
		"path",
		"payloadRequired",
		"requestMediaType",
		"responseMediaType",
		"responseStatusCode",
		"provider",
		"security",
		"queryPars",
		"pathPars",
		"headerPars",
		"cookiePars",
		"payloadPars",
		"responseBody",
		"responseHeaders",
	} {
		t.Run(field, func(t *testing.T) {
			value := decodeRuntimeJSONValue(t, []byte(`{
				"x-uws-runtime": {
					"type": "cmd",
					"`+field+`": {}
				}
			}`))
			require.Error(t, schema.Validate(value))
		})
	}
}

func TestRuntimeSupplementSchemaRejectsRuntimeConfig(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)
	value := decodeRuntimeJSONValue(t, []byte(`{
		"x-uws-runtime": {
			"type": "fnct"
		},
		"x-uws-runtime-config": {
			"provider": {"name": "api"}
		}
	}`))
	require.Error(t, schema.Validate(value))
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
