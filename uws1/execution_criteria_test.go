package uws1

import (
	"context"
	"encoding/json"
	"math"
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

func TestUWS110TruthinessCoversJSONNumbersAndPreservesCollections(t *testing.T) {
	legacy := &Orchestrator{Document: &Document{UWS: "1.9.2"}}
	current := &Orchestrator{Document: &Document{UWS: "1.10.0"}}
	for _, value := range []any{json.Number("0"), json.Number("0.0"), int8(0), uint32(0), uintptr(0), float32(0)} {
		matched, err := current.truthy(value)
		require.NoError(t, err)
		assert.False(t, matched, "zero value %T(%v) is false", value, value)
	}
	for _, value := range []json.Number{"0e10", "-0.000E+999999999"} {
		matched, err := current.truthy(value)
		require.NoError(t, err)
		assert.False(t, matched, "zero significand remains false regardless of exponent %q", value)
	}
	for _, value := range []any{json.Number("-2.5e1"), int16(-1), uint64(1), uintptr(1), float32(0.5)} {
		matched, err := current.truthy(value)
		require.NoError(t, err)
		assert.True(t, matched, "nonzero value %T(%v) is true", value, value)
	}
	for _, value := range []any{nil, false, "", []any{}, map[string]any{}} {
		matched, err := current.truthy(value)
		require.NoError(t, err)
		assert.False(t, matched, "empty/null value %#v is false", value)
	}
	matched, err := legacy.truthy(json.Number("0"))
	require.NoError(t, err)
	assert.True(t, matched, "older versions retain the historical unhandled-number truthiness")
	_, err = current.truthy(math.NaN())
	require.ErrorContains(t, err, "non-finite")
}

func TestUWS110TruthinessRejectsNonJSONCompositeValues(t *testing.T) {
	current := &Orchestrator{Document: &Document{UWS: "1.10.0"}}
	for _, value := range []any{
		map[bool]string{true: "not a JSON object"},
		map[string]any{"nested": math.Inf(1)},
		map[string]any{"nested": json.Number("1e+")},
		[]any{make(chan int)},
		struct{ Value string }{Value: "not a JSON scalar or collection"},
	} {
		_, err := current.truthy(value)
		require.Error(t, err, "unsupported value %#v must not be truthy", value)
	}
}

func TestUWS110TruthinessHandlesVeryLargeJSONNumberExponentsWithoutExpansion(t *testing.T) {
	current := &Orchestrator{Document: &Document{UWS: "1.10.0"}}
	matched, err := current.truthy(json.Number("1e999999999999999999999999999999999999999999999999"))
	require.NoError(t, err)
	assert.True(t, matched)
}

func TestUWS110JSONPointerDecodesURIFragmentsAndValidatesEscapes(t *testing.T) {
	root := map[string]any{"c% d": true, "a/b": true, "tilde~key": true}
	runtime := &mockRuntime{eval: func(context.Context, string) (any, error) { return root, nil }}
	current := &Orchestrator{Document: &Document{UWS: "1.10.0"}, Runtime: runtime}
	for _, pointer := range []string{"#/c%25%20d", "#/a~1b", "#/tilde~0key"} {
		matched, err := current.evaluateCriterion(context.Background(), &Criterion{
			Type: CriterionJSONPath, Context: "$response.body", Condition: pointer,
		})
		require.NoError(t, err)
		assert.True(t, matched, "pointer %s resolves", pointer)
	}
	_, err := current.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "#/bad%2",
	})
	require.ErrorContains(t, err, "percent encoding")
	_, err = current.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "#/bad~2escape",
	})
	require.ErrorContains(t, err, "JSON pointer escape")

	legacy := &Orchestrator{Document: &Document{UWS: "1.9.2"}, Runtime: runtime}
	matched, err := legacy.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "#/c%25%20d",
	})
	require.NoError(t, err)
	assert.False(t, matched, "older versions retain raw-fragment lookup behavior")
}

func TestUWS110XPathUsesXPathNumberTruthiness(t *testing.T) {
	criterion := &Criterion{Type: CriterionXPath, Context: "$response.body", Condition: "number('not-a-number')"}
	runtime := &mockRuntime{eval: func(context.Context, string) (any, error) { return "<root/>", nil }}
	legacy := &Orchestrator{Document: &Document{UWS: "1.9.2"}, Runtime: runtime}
	matched, err := legacy.evaluateCriterion(context.Background(), criterion)
	require.NoError(t, err)
	assert.True(t, matched, "preserve historical NaN conversion before 1.10")
	current := &Orchestrator{Document: &Document{UWS: "1.10.0"}, Runtime: runtime}
	matched, err = current.evaluateCriterion(context.Background(), criterion)
	require.NoError(t, err)
	assert.False(t, matched, "XPath 1.0 converts NaN to false")
}

func TestUWS110RegexCriterionRequiresStringContext(t *testing.T) {
	criterion := &Criterion{Type: CriterionRegex, Context: "$response.statusCode", Condition: "^200$"}
	runtime := &mockRuntime{eval: func(context.Context, string) (any, error) { return 200, nil }}
	legacy := &Orchestrator{Document: &Document{UWS: "1.9.2"}, Runtime: runtime}
	matched, err := legacy.evaluateCriterion(context.Background(), criterion)
	require.NoError(t, err)
	assert.True(t, matched, "older versions keep JSON serialization coercion")
	current := &Orchestrator{Document: &Document{UWS: "1.10.0"}, Runtime: runtime}
	_, err = current.evaluateCriterion(context.Background(), criterion)
	require.ErrorContains(t, err, "must resolve to a string")
}

func TestRegexCriterionUsesSubstringMatching(t *testing.T) {
	criterion := &Criterion{Type: CriterionRegex, Context: "$response.body.status", Condition: "ready"}
	runtime := &mockRuntime{eval: func(context.Context, string) (any, error) { return "not-ready-yet", nil }}
	for _, version := range []string{"1.9.2", "1.10.0"} {
		orch := &Orchestrator{Document: &Document{UWS: version}, Runtime: runtime}
		matched, err := orch.evaluateCriterion(context.Background(), criterion)
		require.NoError(t, err)
		assert.True(t, matched, "regex is a search unless explicitly anchored")
	}
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

func TestJSONPathCriterionUsesRFC9535SelectionTruthiness(t *testing.T) {
	orch := &Orchestrator{Runtime: &mockRuntime{eval: func(context.Context, string) (any, error) {
		return map[string]any{
			"items": []any{
				map[string]any{"active": false},
				map[string]any{"active": true},
			},
		}, nil
	}}}

	matched, err := orch.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "$.items[*].active",
	})
	require.NoError(t, err)
	assert.True(t, matched)

	matched, err = orch.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "$.items[?@.active == false].active",
	})
	require.NoError(t, err)
	assert.False(t, matched)
}

func TestJSONPathCriterionRejectsMalformedPath(t *testing.T) {
	orch := &Orchestrator{Runtime: &mockRuntime{eval: func(context.Context, string) (any, error) {
		return map[string]any{"items": []any{}}, nil
	}}}

	_, err := orch.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "$[",
	})
	require.ErrorContains(t, err, "parse jsonpath")
}

func TestParseCriterionIndexRejectsMalformedArrayIndexes(t *testing.T) {
	for _, token := range []string{"", "01", "1abc", "-1", "-0", "+1", "1.0"} {
		t.Run(token, func(t *testing.T) {
			_, err := parseCriterionIndex(token)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid array index")
		})
	}
}

func TestCriterionJSONPointerArrayIndexIsVersionGated(t *testing.T) {
	root := map[string]any{"items": []any{"zero", "one"}}
	runtime := &mockRuntime{eval: func(context.Context, string) (any, error) {
		return root, nil
	}}

	older := &Orchestrator{Document: &Document{UWS: "1.9.1"}, Runtime: runtime}
	matched, err := older.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "#/items/01",
	})
	require.NoError(t, err)
	assert.True(t, matched, "older declared versions retain the pre-1.9.2 numeric index behavior")

	current := &Orchestrator{Document: &Document{UWS: "1.9.2"}, Runtime: runtime}
	_, err = current.evaluateCriterion(context.Background(), &Criterion{
		Type: CriterionJSONPath, Context: "$response.body", Condition: "#/items/01",
	})
	require.ErrorContains(t, err, "invalid array index")
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
	require.Contains(t, records, "stepop:fetch_step:fetch")
	require.Contains(t, records, "step:fetch_step")
	require.Contains(t, records, "wf:main")
	assert.Equal(t, map[string]any{"city": "Toronto"}, records["stepop:fetch_step:fetch"].Outputs["body"])
	assert.Equal(t, map[string]any{"city": "Toronto"}, records["step:fetch_step"].Outputs["copy"])
	assert.Equal(t, map[string]any{"city": "Toronto"}, records["wf:main"].Outputs["from_step"])
}
