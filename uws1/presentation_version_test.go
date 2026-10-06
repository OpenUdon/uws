package uws1

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUWS113Retains112WireAndExecutionGates(t *testing.T) {
	load := func(name string) map[string]any {
		data, err := os.ReadFile("../versions/" + name + ".json")
		require.NoError(t, err)
		var value map[string]any
		require.NoError(t, json.Unmarshal(data, &value))
		delete(value, "$id")
		delete(value, "description")
		return value
	}
	require.Equal(t, load("1.12.0"), load("1.13.0"), "presentation release must not alter the existing wire/schema constraints")
	for _, version := range []string{"1.12.0", "1.13.0"} {
		doc := validDocument()
		doc.UWS = version
		doc.Operations[0].Effect = OperationEffectRead
		require.NoError(t, doc.Validate())
		doc.Operations[0].Effect = "unsupported"
		require.Error(t, doc.Validate(), "effect gates must remain active")
	}
}
