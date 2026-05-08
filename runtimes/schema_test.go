package runtimes

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func TestRuntimeSupplementSchemaRuntimeTypes(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)

	for _, typ := range runtimeTypeConstants() {
		value := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"`+typ+`"}}`))
		require.NoError(t, schema.Validate(value), "schema should accept runtime type %q", typ)
	}

	plainLDAP := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"ldap"}}`))
	require.Error(t, schema.Validate(plainLDAP))

	httpRuntime := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime":{"type":"http"}}`))
	require.Error(t, schema.Validate(httpRuntime))
}

func TestRuntimeSupplementSchemaRuntimeTypeEnumMatchesConstants(t *testing.T) {
	var doc struct {
		Defs map[string]struct {
			Enum []string `json:"enum"`
		} `json:"$defs"`
	}
	require.NoError(t, json.Unmarshal(readRuntimeSupplementSchema(t), &doc))

	require.Equal(t, runtimeTypeConstants(), doc.Defs["runtime-type"].Enum)
}

func TestRuntimeSupplementSchemaRequiresType(t *testing.T) {
	schema := compileRuntimeSupplementSchema(t)

	emptyRuntime := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime": {}}`))
	require.Error(t, schema.Validate(emptyRuntime))

	commandOnly := decodeRuntimeJSONValue(t, []byte(`{"x-uws-runtime": {"command": "echo ok"}}`))
	require.Error(t, schema.Validate(commandOnly))
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
	data := readRuntimeSupplementSchema(t)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	compiler := jsonschema.NewCompiler()
	const resource = "https://github.com/OpenUdon/uws/versions/runtime.1.0.json"
	require.NoError(t, compiler.AddResource(resource, doc))
	schema, err := compiler.Compile(resource)
	require.NoError(t, err)
	return schema
}

func readRuntimeSupplementSchema(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "versions", "runtime.1.0.json")
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func runtimeTypeConstants() []string {
	return []string{
		RuntimeTypeSSH, RuntimeTypeCmd, RuntimeTypeFnct,
		RuntimeTypeFileIO, RuntimeTypeSQL, RuntimeTypeS3, RuntimeTypeSMTP,
		RuntimeTypeDNS, RuntimeTypeLDAPS, RuntimeTypeSCP, RuntimeTypeSFTP,
		RuntimeTypeLLM,
	}
}

func decodeRuntimeJSONValue(t *testing.T, data []byte) any {
	t.Helper()
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	return value
}
