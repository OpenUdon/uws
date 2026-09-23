package uws1

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpressionAddressableNamesAreVersionGated(t *testing.T) {
	makeDocument := func(version string) *Document {
		doc := testDocument(&Operation{
			OperationID: "fetch",
			Outputs:     map[string]string{"first.name": "$response.body"},
		})
		doc.UWS = version
		doc.Variables = map[string]any{"tenant.name": "acme"}
		doc.Components = &Components{Variables: map[string]any{"shared.name": "default"}}
		doc.Workflows = []*Workflow{{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Inputs: &ParamSchema{Type: "object", Properties: map[string]*ParamSchema{
				"tenant.name": {Type: "string"},
			}},
			Steps: []*Step{{StepID: "fetch_step", OperationRef: "fetch", Inputs: map[string]any{"input.name": "value"}}},
		}}
		return doc
	}

	legacy := makeDocument("1.9.2")
	require.NoError(t, legacy.Validate())

	current := makeDocument("1.10.0")
	err := current.Validate()
	require.Error(t, err)
	for _, name := range []string{"first.name", "tenant.name", "shared.name", "input.name"} {
		require.True(t, strings.Contains(err.Error(), name), "expected %q in %v", name, err)
	}

	encoded, marshalErr := json.Marshal(current)
	require.NoError(t, marshalErr)
	require.Error(t, compileUWSSchema(t).Validate(decodeJSONValue(t, encoded)), "1.10 schema must reject names the expression grammar cannot address")
}

func TestWorkflowReferenceIsALiteralWorkflowID(t *testing.T) {
	doc := testDocument(&Operation{OperationID: "noop"})
	doc.UWS = "1.10.0"
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "call", StepExecutionFields: StepExecutionFields{Workflow: "$variables.target"}}},
	}, {
		WorkflowID: "target",
		Type:       WorkflowTypeSequence,
	}}
	require.ErrorContains(t, doc.Validate(), "references unknown workflowId \"$variables.target\"")
}
