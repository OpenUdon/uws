package convert

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/OpenUdon/uws/uws1"
	"github.com/stretchr/testify/require"
)

func TestPendingStepConversionRoundTrips(t *testing.T) {
	data, err := os.ReadFile("../testdata/candidate/pending-only.1.12.json")
	require.NoError(t, err)
	var doc uws1.Document
	require.NoError(t, json.Unmarshal(data, &doc))
	pending := doc.Workflows[0].Steps[0].Pending
	pending.Purpose += "\nThe output is an object schema."
	pending.Extensions = map[string]any{"x-pending-note": map[string]any{"$expr": "literal"}}
	want := *pending

	jsonData, err := json.Marshal(&doc)
	require.NoError(t, err)
	assertPending := func(encoded []byte) {
		t.Helper()
		var got uws1.Document
		require.NoError(t, json.Unmarshal(encoded, &got))
		require.True(t, reflect.DeepEqual(&want, got.Workflows[0].Steps[0].Pending), "pending field-set changed across conversion:\nwant %#v\n got %#v", &want, got.Workflows[0].Steps[0].Pending)
	}

	yamlData, err := JSONToYAML(jsonData)
	require.NoError(t, err)
	assertPending(mustYAMLToJSON(t, yamlData))

	hclData, err := JSONToHCL(jsonData)
	require.NoError(t, err)
	require.Contains(t, string(hclData), "pending {")
	require.Contains(t, string(hclData), "inputs {")
	require.Contains(t, string(hclData), `effect  = "read"`)
	hclJSON, err := HCLToJSON(hclData)
	require.NoError(t, err)
	assertPending(hclJSON)
	var hclDoc uws1.Document
	require.NoError(t, json.Unmarshal(hclJSON, &hclDoc))
	require.NotNil(t, hclDoc.Operations)
	require.Empty(t, hclDoc.Operations, "pending-only HCL round trips to an explicit empty operations array")
}
