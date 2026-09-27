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

const candidate112SchemaID = "https://github.com/OpenUdon/uws/versions/1.12.0.json"

func compileCandidate112Schema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	data, err := os.ReadFile("../testdata/candidate/1.12.0.json")
	require.NoError(t, err)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	compiler := jsonschema.NewCompiler()
	require.NoError(t, compiler.AddResource(candidate112SchemaID, doc))
	schema, err := compiler.Compile(candidate112SchemaID)
	require.NoError(t, err)
	return schema
}

func candidateEffectDocument(effect string, includeEffect bool) []byte {
	operation := map[string]any{
		"operationId":             "candidate_operation",
		"x-uws-operation-profile": "test.profile",
	}
	if includeEffect {
		operation["effect"] = effect
	}
	doc := map[string]any{
		"uws":        "1.12.0",
		"info":       map[string]any{"title": "Effect candidate", "version": "1.0.0"},
		"operations": []any{operation},
	}
	data, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	return data
}

func TestOperationEffect_CandidateSchemaAndVersionGate(t *testing.T) {
	candidate := compileCandidate112Schema(t)
	current := compileUWSSchema(t)

	for _, effect := range []OperationEffect{OperationEffectRead, OperationEffectWrite, OperationEffectUnknown} {
		t.Run(string(effect), func(t *testing.T) {
			data := candidateEffectDocument(string(effect), true)
			require.NoError(t, candidate.Validate(decodeJSONValue(t, data)))
			require.Error(t, current.Validate(decodeJSONValue(t, data)), "1.11 must not admit the 1.12 field")

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

	require.NoError(t, candidate.Validate(decodeJSONValue(t, candidateEffectDocument("", false))), "effect is optional")
	require.Error(t, candidate.Validate(decodeJSONValue(t, candidateEffectDocument("maybe", true))), "unknown effect values are rejected")

	result := &ValidationResult{}
	validateOperationEffect("maybe", "operations[0].effect", "1.12.0", result)
	require.ErrorContains(t, result, `"maybe" is not valid`)

	legacy := validDocument()
	legacy.Operations[0].Effect = OperationEffectRead
	require.ErrorContains(t, legacy.Validate(), "operations[0].effect requires UWS 1.12.0 or later")

	declared := validDocument()
	declared.UWS = "1.12.0"
	declared.Operations[0].Effect = OperationEffectRead
	assert.ErrorContains(t, declared.Validate(), `version "1.12.0" is not a published UWS version`)
	assert.NotContains(t, declared.Validate().Error(), "operations[0].effect")
}

func TestSchemaParity_Candidate112OperationFields(t *testing.T) {
	data, err := os.ReadFile("../testdata/candidate/1.12.0.json")
	require.NoError(t, err)
	var schema map[string]any
	require.NoError(t, json.Unmarshal(data, &schema))
	version := schemaVersionFromID(t, schema)
	require.Equal(t, "1.12.0", version)

	var operationEntry schemaParityEntry
	for _, entry := range schemaParityEntries() {
		if entry.defName == "operation-object" {
			operationEntry = entry
			break
		}
	}
	require.NotEmpty(t, operationEntry.label)
	props := dropExtensionKeys(schemaPropertyNames(t, schema, operationEntry.defName))
	assert.ElementsMatch(t, knownFieldsForSchemaVersion(operationEntry, version), props)
	assert.Contains(t, props, "effect")
}
