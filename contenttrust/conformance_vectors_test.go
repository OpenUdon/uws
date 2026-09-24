package contenttrust

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUWS111ExpressionConformanceVectorsExecuteAgainstParser(t *testing.T) {
	data, err := os.ReadFile("../testdata/conformance/1.11.0.json")
	require.NoError(t, err)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var corpus struct {
		UWS   string `json:"uws"`
		Cases []struct {
			ID        string          `json:"id"`
			Operation string          `json:"operation"`
			Input     json.RawMessage `json:"input"`
			Expected  json.RawMessage `json:"expected"`
		} `json:"cases"`
	}
	require.NoError(t, decoder.Decode(&corpus))
	require.Equal(t, "1.11.0", corpus.UWS)

	executed := 0
	for _, vector := range corpus.Cases {
		vector := vector
		switch vector.Operation {
		case "document.validate", "document.execute":
			continue // The uws1 package executes these vectors against the core.
		case "expression.parse":
		default:
			t.Fatalf("unknown UWS 1.11 vector operation %q", vector.Operation)
		}
		t.Run(vector.ID, func(t *testing.T) {
			var input struct {
				DeclaredVersion string `json:"declaredVersion"`
				Field           string `json:"field"`
				Expression      string `json:"expression"`
				LoopContext     bool   `json:"loopContext"`
			}
			inputDecoder := json.NewDecoder(bytes.NewReader(vector.Input))
			require.NoError(t, inputDecoder.Decode(&input))
			allowNumeric := input.Field == "wait" || input.Field == "batchSize"
			_, valid := parseExpressionValue(input.Expression, input.DeclaredVersion, input.LoopContext, allowNumeric)
			var expected struct {
				Valid bool `json:"valid"`
			}
			require.NoError(t, json.Unmarshal(vector.Expected, &expected))
			require.Equal(t, expected.Valid, valid)
		})
		executed++
	}
	require.Positive(t, executed, "the corpus must contain expression vectors")
}
