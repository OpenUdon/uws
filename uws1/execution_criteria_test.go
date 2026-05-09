package uws1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrchestratorExecutesRegexCriterion(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "fetch",
		SuccessCriteria: []*Criterion{{
			Type:      CriterionRegex,
			Context:   "$response.body",
			Condition: "^ok$",
		}},
	})
	runtime := &mockRuntime{
		eval: func(ctx context.Context, expr string) (any, error) {
			if expr == "$response.body" {
				return "ok", nil
			}
			return nil, nil
		},
	}
	doc.SetRuntime(runtime)

	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "fetch_step", OperationRef: "fetch"}},
	}}
	require.NoError(t, doc.Execute(context.Background()))
	assert.Equal(t, []string{"fetch"}, runtime.executedLeafs)
}

func TestOrchestratorExecutesJSONPathCriterion(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "fetch",
		SuccessCriteria: []*Criterion{{
			Type:      CriterionJSONPath,
			Context:   "$response.body",
			Condition: "#/id",
		}},
	})
	runtime := &mockRuntime{
		eval: func(ctx context.Context, expr string) (any, error) {
			if expr == "$response.body" {
				return map[string]any{"id": "123"}, nil
			}
			return nil, nil
		},
	}
	doc.SetRuntime(runtime)

	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "fetch_step", OperationRef: "fetch"}},
	}}
	require.NoError(t, doc.Execute(context.Background()))
	assert.Equal(t, []string{"fetch"}, runtime.executedLeafs)
}

func TestParseCriterionIndexRejectsMalformedArrayIndexes(t *testing.T) {
	for _, token := range []string{"", "1abc", "-1", "+1", "1.0"} {
		t.Run(token, func(t *testing.T) {
			_, err := parseCriterionIndex(token)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid array index")
		})
	}
}

func TestOrchestratorExecutesXPathCriterion(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "fetch",
		SuccessCriteria: []*Criterion{{
			Type:      CriterionXPath,
			Context:   "$response.body",
			Condition: "count(/root/item[@kind='primary'][text()='123']) = 1",
		}},
	})
	runtime := &mockRuntime{
		eval: func(ctx context.Context, expr string) (any, error) {
			if expr == "$response.body" {
				return `<root><item kind="secondary">nope</item><item kind="primary">123</item></root>`, nil
			}
			return nil, nil
		},
	}
	doc.SetRuntime(runtime)

	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps:      []*Step{{StepID: "fetch_step", OperationRef: "fetch"}},
	}}
	require.NoError(t, doc.Execute(context.Background()))
	assert.Equal(t, []string{"fetch"}, runtime.executedLeafs)
}

func TestOrchestratorCapturesOutputs(t *testing.T) {
	doc := testDocument(&Operation{
		OperationID: "fetch",
		Outputs: map[string]string{
			"body": "$response.body",
		},
	})
	doc.Workflows = []*Workflow{{
		WorkflowID: "main",
		Type:       WorkflowTypeSequence,
		Steps: []*Step{{
			StepID:       "fetch_step",
			OperationRef: "fetch",
			Outputs: map[string]string{
				"copy": "$response.body",
			},
		}},
		Outputs: map[string]string{
			"from_step": "$steps.fetch_step.outputs.copy",
		},
	}}
	runtime := &mockRuntime{
		eval: func(ctx context.Context, expr string) (any, error) {
			state, _ := ExecutionContextFromContext(ctx)
			switch expr {
			case "$response.body":
				return map[string]any{"city": "Toronto"}, nil
			case "$steps.fetch_step.outputs.copy":
				if state == nil {
					return nil, nil
				}
				record, ok := state.Records["step:fetch_step"]
				if !ok {
					return nil, nil
				}
				return record.Outputs["copy"], nil
			default:
				return nil, nil
			}
		},
	}
	doc.SetRuntime(runtime)

	require.NoError(t, doc.Execute(context.Background()))
	records := doc.ExecutionRecords()
	require.Contains(t, records, "op:fetch")
	require.Contains(t, records, "step:fetch_step")
	require.Contains(t, records, "wf:main")
	assert.Equal(t, map[string]any{"city": "Toronto"}, records["op:fetch"].Outputs["body"])
	assert.Equal(t, map[string]any{"city": "Toronto"}, records["step:fetch_step"].Outputs["copy"])
	assert.Equal(t, map[string]any{"city": "Toronto"}, records["wf:main"].Outputs["from_step"])
}
