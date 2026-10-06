package uws1

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func compileSchemaFile(t *testing.T, path string) *jsonschema.Schema {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	var metadata struct {
		ID string `json:"$id"`
	}
	require.NoError(t, json.Unmarshal(data, &metadata))
	require.NotEmpty(t, metadata.ID)
	compiler := jsonschema.NewCompiler()
	require.NoError(t, compiler.AddResource(metadata.ID, doc))
	schema, err := compiler.Compile(metadata.ID)
	require.NoError(t, err)
	return schema
}

func effectDocument(effect string, includeEffect bool) []byte {
	operation := map[string]any{
		"operationId":             "effect_operation",
		"x-uws-operation-profile": "test.profile",
	}
	if includeEffect {
		operation["effect"] = effect
	}
	doc := map[string]any{
		"uws":        "1.12.0",
		"info":       map[string]any{"title": "Effect contract", "version": "1.0.0"},
		"operations": []any{operation},
	}
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	return data
}

func TestOperationEffect_PublishedSchemaAndVersionGate(t *testing.T) {
	current := compileUWSSchema(t)
	legacySchema := compileSchemaFile(t, "../versions/1.11.0.json")

	for _, effect := range []OperationEffect{OperationEffectRead, OperationEffectWrite, OperationEffectUnknown} {
		t.Run(string(effect), func(t *testing.T) {
			data := effectDocument(string(effect), true)
			require.NoError(t, current.Validate(decodeJSONValue(t, data)))
			require.Error(t, legacySchema.Validate(decodeJSONValue(t, data)), "1.11 must not admit the 1.12 field")

			var operation Operation
			require.NoError(t, json.Unmarshal([]byte(`{"operationId":"op","effect":"`+string(effect)+`"}`), &operation))
			assert.Equal(t, effect, operation.Effect)
			encoded, err := json.Marshal(operation)
			require.NoError(t, err)
			assert.Contains(t, string(encoded), `"effect":"`+string(effect)+`"`)

			for _, version := range []string{"1.12.0", "1.12.1-rc.1"} {
				result := &ValidationResult{}
				validateOperationEffect(effect, "operations[0].effect", version, result)
				assert.True(t, result.Valid(), "%s should support effect %q", version, effect)
			}

			result := &ValidationResult{}
			validateOperationEffect(effect, "operations[0].effect", "1.11.0", result)
			require.ErrorContains(t, result, "requires UWS 1.12.0 or later")
		})
	}

	require.NoError(t, current.Validate(decodeJSONValue(t, effectDocument("", false))), "effect is optional")
	require.Error(t, current.Validate(decodeJSONValue(t, effectDocument("maybe", true))), "unknown effect values are rejected")

	result := &ValidationResult{}
	validateOperationEffect("maybe", "operations[0].effect", "1.12.0", result)
	require.ErrorContains(t, result, `"maybe" is not valid`)

	legacy := validDocument()
	legacy.Operations[0].Effect = OperationEffectRead
	require.ErrorContains(t, legacy.Validate(), "operations[0].effect requires UWS 1.12.0 or later")

	declared := validDocument()
	declared.UWS = "1.12.0"
	declared.Operations[0].Effect = OperationEffectRead
	assert.NoError(t, declared.Validate())
}

func TestSchemaParity_Published112Fields(t *testing.T) {
	data, err := os.ReadFile("../versions/1.12.0.json")
	require.NoError(t, err)
	var schema map[string]any
	require.NoError(t, json.Unmarshal(data, &schema))
	version := schemaVersionFromID(t, schema)
	require.Equal(t, "1.12.0", version)

	var operationEntry, stepEntry, pendingStepEntry schemaParityEntry
	for _, entry := range schemaParityEntries() {
		switch entry.defName {
		case "operation-object":
			operationEntry = entry
		case "step-object":
			stepEntry = entry
		case "pending-step-object":
			pendingStepEntry = entry
		}
	}
	for _, entry := range []schemaParityEntry{operationEntry, stepEntry, pendingStepEntry} {
		require.NotEmpty(t, entry.label)
		props := dropExtensionKeys(schemaPropertyNames(t, schema, entry.defName))
		assert.ElementsMatch(t, knownFieldsForSchemaVersion(entry, version), props)
	}
	assert.Contains(t, dropExtensionKeys(schemaPropertyNames(t, schema, operationEntry.defName)), "effect")
	assert.Contains(t, dropExtensionKeys(schemaPropertyNames(t, schema, stepEntry.defName)), "pending")
}
