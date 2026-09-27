package convert

import (
	"encoding/json"
	"testing"

	"github.com/OpenUdon/uws/uws1"
	"github.com/stretchr/testify/require"
)

const effectRoundTripJSON = `{"uws":"1.12.0","info":{"title":"Effect round trip","version":"1.0.0"},"operations":[{"operationId":"op","x-uws-operation-profile":"test.profile","effect":"read"}]}`

func requireReadEffect(t *testing.T, data []byte) {
	t.Helper()
	var doc uws1.Document
	require.NoError(t, json.Unmarshal(data, &doc))
	require.Len(t, doc.Operations, 1)
	require.Equal(t, uws1.OperationEffectRead, doc.Operations[0].Effect)
}

func TestOperationEffectConversionRoundTrips(t *testing.T) {
	jsonData := []byte(effectRoundTripJSON)

	hclData, err := JSONToHCL(jsonData)
	require.NoError(t, err)
	require.Contains(t, string(hclData), `effect = "read"`)
	hclJSON, err := HCLToJSON(hclData)
	require.NoError(t, err)
	requireReadEffect(t, hclJSON)

	yamlData, err := JSONToYAML(jsonData)
	require.NoError(t, err)
	requireReadEffect(t, mustYAMLToJSON(t, yamlData))

	yamlSource := []byte(`uws: "1.12.0"
info:
  title: Effect round trip
  version: "1.0.0"
operations:
  - operationId: op
    x-uws-operation-profile: test.profile
    effect: read
`)
	yamlHCL, err := YAMLToHCL(yamlSource)
	require.NoError(t, err)
	yamlHCLJSON, err := HCLToJSON(yamlHCL)
	require.NoError(t, err)
	requireReadEffect(t, yamlHCLJSON)
}

func mustYAMLToJSON(t *testing.T, data []byte) []byte {
	t.Helper()
	jsonData, err := YAMLToJSON(data)
	require.NoError(t, err)
	return jsonData
}
