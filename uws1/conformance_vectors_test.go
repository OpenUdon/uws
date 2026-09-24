package uws1

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const uws111ConformanceCorpusSHA256 = "51ea40b3f392125a22a95d515ee5002071dd8f4d98ab49415819b62f61ef26b5"

type conformanceVectorCase struct {
	ID        string          `json:"id"`
	Area      string          `json:"area,omitempty"`
	Operation string          `json:"operation,omitempty"`
	Input     json.RawMessage `json:"input"`
	Expected  json.RawMessage `json:"expected"`
}

type conformanceVectorCorpus struct {
	UWS   string                  `json:"uws"`
	Cases []conformanceVectorCase `json:"cases"`
}

func readConformanceVectors(t *testing.T, path string) conformanceVectorCorpus {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var corpus conformanceVectorCorpus
	require.NoError(t, decoder.Decode(&corpus))
	return corpus
}

func decodeConformanceInput[T any](t *testing.T, raw json.RawMessage) T {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var input T
	require.NoError(t, decoder.Decode(&input))
	return input
}

func requireConformanceJSON(t *testing.T, expected json.RawMessage, actual any) {
	t.Helper()
	data, err := json.Marshal(actual)
	require.NoError(t, err)
	require.JSONEq(t, string(expected), string(data))
}

func TestUWS110LanguageNeutralVectorsExecuteAgainstGo(t *testing.T) {
	corpus := readConformanceVectors(t, "../testdata/conformance/1.10.0.json")
	require.Equal(t, "1.10.0", corpus.UWS)
	require.NotEmpty(t, corpus.Cases)
	seen := make(map[string]bool, len(corpus.Cases))
	for _, vector := range corpus.Cases {
		vector := vector
		t.Run(vector.ID, func(t *testing.T) {
			require.NotEmpty(t, vector.ID)
			require.NotEmpty(t, vector.Area)
			require.False(t, seen[vector.ID], "duplicate vector ID")
			seen[vector.ID] = true
			runUWS110Vector(t, vector)
		})
	}
	require.Len(t, seen, len(corpus.Cases), "every 1.10 vector must be executed exactly once")
}

func runUWS110Vector(t *testing.T, vector conformanceVectorCase) {
	t.Helper()
	switch vector.ID {
	case "truthiness-json-values":
		input := decodeConformanceInput[[]any](t, vector.Input)
		want := decodeConformanceInput[[]bool](t, vector.Expected)
		got := make([]bool, 0, len(input))
		for _, value := range input {
			truthy, err := truthyValueV110(value)
			require.NoError(t, err)
			got = append(got, truthy)
		}
		requireConformanceValue(t, want, got)
	case "expression-addressable-identifiers":
		input := decodeConformanceInput[[]string](t, vector.Input)
		want := decodeConformanceInput[[]bool](t, vector.Expected)
		got := make([]bool, 0, len(input))
		for _, name := range input {
			doc := testDocument(&Operation{OperationID: "noop"})
			doc.UWS = "1.10.0"
			doc.Variables = map[string]any{name: true}
			got = append(got, doc.Validate() == nil)
		}
		requireConformanceValue(t, want, got)
	case "json-pointer-fragment-decoding":
		var input struct {
			Document  map[string]any `json:"document"`
			Fragments []string       `json:"fragments"`
		}
		input = decodeConformanceInput[struct {
			Document  map[string]any `json:"document"`
			Fragments []string       `json:"fragments"`
		}](t, vector.Input)
		want := decodeConformanceInput[[]any](t, vector.Expected)
		orch := &Orchestrator{Document: &Document{UWS: "1.10.0"}}
		got := make([]any, 0, len(input.Fragments))
		for _, fragment := range input.Fragments {
			value, found, err := orch.resolveCriterionPointer(input.Document, fragment, true)
			switch {
			case err != nil:
				got = append(got, "invalid-escape")
			case !found:
				got = append(got, "missing")
			default:
				got = append(got, value)
			}
		}
		requireConformanceValue(t, want, got)
	case "switch-first-match":
		input := decodeConformanceInput[[]any](t, vector.Input)
		require.Len(t, input, 3)
		doc := testDocument(
			&Operation{OperationID: "false-op"},
			&Operation{OperationID: "unguarded-op"},
			&Operation{OperationID: "true-op"},
		)
		doc.UWS = "1.10.0"
		doc.Workflows = []*Workflow{{
			WorkflowID: "main", Type: WorkflowTypeSwitch,
			Cases: []*Case{
				{CaseFields: CaseFields{Name: "false", When: "false"}, Steps: []*Step{{StepID: "false-step", OperationRef: "false-op"}}},
				{CaseFields: CaseFields{Name: "unguarded"}, Steps: []*Step{{StepID: "unguarded-step", OperationRef: "unguarded-op"}}},
				{CaseFields: CaseFields{Name: "true", When: "true"}, Steps: []*Step{{StepID: "true-step", OperationRef: "true-op"}}},
			},
		}}
		runtime := &mockRuntime{expressions: map[string]any{"false": input[0], "true": input[2]}}
		doc.SetRuntime(runtime)
		require.NoError(t, doc.Execute(context.Background()))
		want := decodeConformanceInput[string](t, vector.Expected)
		leafs := runtime.leafs()
		require.Equal(t, []string{"unguarded-op"}, leafs)
		require.Equal(t, want, strings.TrimSuffix(leafs[0], "-op"))
	case "entry-selection":
		input := decodeConformanceInput[struct {
			WorkflowIDs []string `json:"workflowIds"`
		}](t, vector.Input)
		doc := testDocument(&Operation{OperationID: "noop"})
		doc.UWS = "1.10.0"
		for _, id := range input.WorkflowIDs {
			doc.Workflows = append(doc.Workflows, &Workflow{WorkflowID: id, Type: WorkflowTypeSequence})
		}
		orch := NewOrchestrator(doc, &mockRuntime{})
		entry, err := orch.entryWorkflow()
		require.NoError(t, err)
		require.Equal(t, decodeConformanceInput[string](t, vector.Expected), entry.WorkflowID)
	case "non-await-wait-bounds":
		input := decodeConformanceInput[[]json.Number](t, vector.Input)
		want := decodeConformanceInput[[]bool](t, vector.Expected)
		got := make([]bool, 0, len(input))
		for _, value := range input {
			seconds, ok := waitSecondsNumber(value)
			got = append(got, ok && !math.IsNaN(seconds) && !math.IsInf(seconds, 0) && seconds >= 0 && seconds <= maxWaitSeconds)
		}
		requireConformanceValue(t, want, got)
	case "loop-order-and-batches":
		var input struct {
			Items     []any       `json:"items"`
			BatchSize json.Number `json:"batchSize"`
		}
		input = decodeConformanceInput[struct {
			Items     []any       `json:"items"`
			BatchSize json.Number `json:"batchSize"`
		}](t, vector.Input)
		batch, err := input.BatchSize.Int64()
		require.NoError(t, err)
		doc := testDocument(&Operation{OperationID: "noop"})
		doc.UWS = "1.10.0"
		doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeLoop, StructuralFields: StructuralFields{Items: "items", BatchSize: "$inputs.batch"}}}
		doc.SetRuntime(&mockRuntime{items: map[string][]any{"items": input.Items}, expressions: map[string]any{"$inputs.batch": batch}})
		require.NoError(t, doc.Execute(context.Background()))
		got := doc.ExecutionRecords()["wf:main"].Result
		requireConformanceJSON(t, vector.Expected, got)
	case "foreach-output-aggregation":
		input := decodeConformanceInput[struct {
			Items  []any  `json:"items"`
			Output string `json:"output"`
		}](t, vector.Input)
		doc := testDocument(&Operation{OperationID: "each-op", Outputs: map[string]string{"item": input.Output}})
		doc.UWS = "1.10.0"
		doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{{
			StepID: "each", OperationRef: "each-op",
			StepExecutionFields: StepExecutionFields{ForEach: "items"}, Outputs: map[string]string{"item": input.Output},
		}}}}
		runtime := &mockRuntime{
			items: map[string][]any{"items": input.Items},
			eval: func(ctx context.Context, expr string) (any, error) {
				if expr == "$item" {
					state, ok := ExecutionContextFromContext(ctx)
					if ok && state.Iteration != nil {
						return state.Iteration.Item, nil
					}
				}
				return nil, nil
			},
		}
		doc.SetRuntime(runtime)
		require.NoError(t, doc.Execute(context.Background()))
		requireConformanceJSON(t, vector.Expected, doc.ExecutionRecords()["step:each"].Outputs["item"])
	case "merge-keeps-declared-dependency-order":
		input := decodeConformanceInput[struct {
			DependsOn []string       `json:"dependsOn"`
			Records   map[string]any `json:"records"`
		}](t, vector.Input)
		doc := testDocument(&Operation{OperationID: "left-op"}, &Operation{OperationID: "right-op"})
		doc.UWS = "1.10.0"
		doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{
			{StepID: "left", OperationRef: "left-op"},
			{StepID: "right", OperationRef: "right-op"},
			{StepID: "join", Type: WorkflowTypeMerge, StepExecutionFields: StepExecutionFields{DependsOn: input.DependsOn}},
		}}}
		runtime := &mockRuntime{}
		doc.SetRuntime(runtime)
		require.NoError(t, doc.Execute(context.Background()))
		merged, ok := doc.ExecutionRecords()["step:join"].Result.([]map[string]any)
		require.True(t, ok)
		got := make([]any, 0, len(merged))
		for _, record := range merged {
			got = append(got, input.Records[record["id"].(string)])
		}
		requireConformanceJSON(t, vector.Expected, got)
	case "retry-limit-counts-retries-after-initial":
		input := decodeConformanceInput[struct {
			RetryLimit     int `json:"retryLimit"`
			InitialAttempt int `json:"initialAttempt"`
		}](t, vector.Input)
		expected := decodeConformanceInput[struct {
			MaximumAttempts int `json:"maximumAttempts"`
			MaximumRetries  int `json:"maximumRetries"`
		}](t, vector.Expected)
		calls := 0
		doc := testDocument(&Operation{OperationID: "retry-op", OnFailure: []*FailureAction{{Name: "retry", Type: "retry", RetryLimit: input.RetryLimit}}})
		doc.UWS = "1.10.0"
		doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{{StepID: "retry", OperationRef: "retry-op"}}}}
		doc.SetRuntime(&mockRuntime{execute: func(context.Context, *Operation) error {
			calls++
			if calls < expected.MaximumAttempts {
				return errors.New("scripted retry")
			}
			return nil
		}})
		require.NoError(t, doc.Execute(context.Background()))
		require.Equal(t, expected.MaximumAttempts, calls)
		require.Equal(t, expected.MaximumRetries, calls-input.InitialAttempt)
	case "trigger-output-index-must-be-declared":
		input := decodeConformanceInput[struct {
			Outputs       []string `json:"outputs"`
			DispatchIndex int      `json:"dispatchIndex"`
		}](t, vector.Input)
		doc := testDocument(&Operation{OperationID: "fetch"})
		doc.UWS = "1.10.0"
		doc.Workflows = []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: []*Step{{StepID: "fetch-step", OperationRef: "fetch"}}}}
		doc.Triggers = []*Trigger{{TriggerID: "incoming", Outputs: input.Outputs, Routes: []*TriggerRoute{{TriggerRouteFields: TriggerRouteFields{Output: "created", To: []string{"fetch-step"}}}}}}
		err := NewOrchestrator(doc, &mockRuntime{}).ExecuteTrigger(context.Background(), "incoming", input.DispatchIndex, nil)
		want := decodeConformanceInput[string](t, vector.Expected)
		if want == "error" {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	default:
		t.Fatalf("no Go conformance runner for UWS 1.10 vector %q", vector.ID)
	}
}

func requireConformanceValue(t *testing.T, expected, actual any) {
	t.Helper()
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)
	requireConformanceJSON(t, expectedJSON, actual)
}

func TestUWS111TypedConformanceVectorsExecuteAgainstGo(t *testing.T) {
	corpus := readConformanceVectors(t, "../testdata/conformance/1.11.0.json")
	require.Equal(t, "1.11.0", corpus.UWS)
	require.NotEmpty(t, corpus.Cases)
	seen := make(map[string]bool, len(corpus.Cases))
	executed := 0
	for _, vector := range corpus.Cases {
		vector := vector
		require.NotEmpty(t, vector.ID)
		require.NotEmpty(t, vector.Operation)
		require.NotEmpty(t, vector.Input)
		require.NotEmpty(t, vector.Expected)
		require.False(t, seen[vector.ID], "duplicate vector ID %q", vector.ID)
		seen[vector.ID] = true
		switch vector.Operation {
		case "expression.parse":
			continue // The contenttrust package executes these parser vectors.
		case "document.validate", "document.execute":
			t.Run(vector.ID, func(t *testing.T) {
				runUWS111CoreVector(t, vector)
			})
			executed++
		default:
			t.Fatalf("unknown UWS 1.11 vector operation %q", vector.Operation)
		}
	}
	require.Positive(t, executed, "at least one core vector must execute in the uws1 runner")
	require.Len(t, seen, len(corpus.Cases), "every typed 1.11 vector must be classified by a Go runner")
}

func TestUWS111ConformanceCorpusDigestIsPinnedBySpecification(t *testing.T) {
	data, err := os.ReadFile("../testdata/conformance/1.11.0.json")
	require.NoError(t, err)
	digest := fmt.Sprintf("%x", sha256.Sum256(data))
	require.Equal(t, uws111ConformanceCorpusSHA256, digest)

	spec, err := os.ReadFile("../versions/1.11.0.md")
	require.NoError(t, err)
	require.Contains(t, string(spec), "SHA-256 `"+uws111ConformanceCorpusSHA256+"`")
}

type uws111ValidationInput struct {
	Document json.RawMessage `json:"document"`
}

type uws111ExecutionInput struct {
	Document json.RawMessage `json:"document"`
	Runtime  struct {
		Items        map[string][]any `json:"items"`
		Expressions  map[string]any   `json:"expressions"`
		LeafFailures map[string]int   `json:"leafFailures"`
	} `json:"runtime"`
	Observe struct {
		RecordStatuses       []uws111RecordSelector `json:"recordStatuses"`
		MergeResultIDsFrom   *uws111RecordSelector  `json:"mergeResultIDsFrom"`
		LoopBatchIndexesFrom *uws111RecordSelector  `json:"loopBatchIndexesFrom"`
	} `json:"observe"`
}

type uws111RecordSelector struct {
	Kind             string   `json:"kind"`
	ID               string   `json:"id"`
	WorkflowCallPath []string `json:"workflowCallPath,omitempty"`
}

func runUWS111CoreVector(t *testing.T, vector conformanceVectorCase) {
	t.Helper()
	switch vector.Operation {
	case "document.validate":
		input := decodeConformanceInput[uws111ValidationInput](t, vector.Input)
		var doc Document
		require.NoError(t, json.Unmarshal(input.Document, &doc))
		err := doc.Validate()
		got := map[string]any{"valid": err == nil, "errorKind": conformanceErrorKind(err)}
		requireConformanceJSON(t, vector.Expected, got)
	case "document.execute":
		input := decodeConformanceInput[uws111ExecutionInput](t, vector.Input)
		var doc Document
		require.NoError(t, json.Unmarshal(input.Document, &doc))
		runtime := &conformanceScriptRuntime{
			items: input.Runtime.Items, expressions: input.Runtime.Expressions,
			failures: input.Runtime.LeafFailures, leafOrder: make([]string, 0),
		}
		doc.SetRuntime(runtime)
		err := doc.Execute(context.Background())
		records := doc.ExecutionRecords()
		statuses := make([]any, 0, len(input.Observe.RecordStatuses))
		for _, selector := range input.Observe.RecordStatuses {
			key, ok := conformanceRecordKey(selector)
			require.True(t, ok, "Go conformance adapter cannot resolve selector %+v", selector)
			status := any(nil)
			if record, exists := records[key]; exists {
				status = record.Status
			}
			statuses = append(statuses, map[string]any{
				"selector": selector, "status": status,
			})
		}
		mergeIDs := make([]string, 0)
		if selector := input.Observe.MergeResultIDsFrom; selector != nil {
			key, ok := conformanceRecordKey(*selector)
			require.True(t, ok, "Go conformance adapter cannot resolve selector %+v", *selector)
			if record, exists := records[key]; exists {
				if result, ok := record.Result.([]map[string]any); ok {
					for _, row := range result {
						if id, ok := row["id"].(string); ok {
							mergeIDs = append(mergeIDs, id)
						}
					}
				}
			}
		}
		batchIndexes := make([]any, 0)
		if selector := input.Observe.LoopBatchIndexesFrom; selector != nil {
			key, ok := conformanceRecordKey(*selector)
			require.True(t, ok, "Go conformance adapter cannot resolve selector %+v", *selector)
			if record, exists := records[key]; exists {
				if result, ok := record.Result.([]map[string]any); ok {
					for _, row := range result {
						batchIndexes = append(batchIndexes, row["batchIndex"])
					}
				}
			}
		}
		got := map[string]any{
			"errorKind": conformanceErrorKind(err), "leafOrder": runtime.leafOrder,
			"recordStatuses": statuses, "mergeResultIDs": mergeIDs, "loopBatchIndexes": batchIndexes,
		}
		requireConformanceJSON(t, vector.Expected, got)
	default:
		t.Fatalf("unexpected core vector operation %q", vector.Operation)
	}
}

func conformanceRecordKey(selector uws111RecordSelector) (string, bool) {
	switch selector.Kind {
	case "step":
		return "step:" + selector.ID, selector.ID != "" && len(selector.WorkflowCallPath) == 0
	case "workflow":
		return "wf:" + selector.ID, selector.ID != "" && len(selector.WorkflowCallPath) == 0
	default:
		return "", false
	}
}

func conformanceErrorKind(err error) any {
	if err == nil {
		return nil
	}
	message := err.Error()
	switch {
	case strings.Contains(message, "is not a published UWS version"):
		return "unsupported-version"
	case strings.Contains(message, "already completed"):
		return "goto-target-already-completed"
	case strings.Contains(message, "wait duration"):
		return "invalid-wait-duration"
	default:
		return "execution-error"
	}
}

type conformanceScriptRuntime struct {
	items       map[string][]any
	expressions map[string]any
	failures    map[string]int
	leafOrder   []string
}

func (r *conformanceScriptRuntime) ExecuteLeaf(_ context.Context, operation *Operation) error {
	r.leafOrder = append(r.leafOrder, operation.OperationID)
	if r.failures[operation.OperationID] > 0 {
		r.failures[operation.OperationID]--
		return errors.New("scripted leaf failure")
	}
	return nil
}

func (r *conformanceScriptRuntime) EvaluateExpression(ctx context.Context, expression string) (any, error) {
	if expression == "$item" || expression == "$index" || expression == "$batchIndex" {
		if state, ok := ExecutionContextFromContext(ctx); ok && state.Iteration != nil {
			switch expression {
			case "$item":
				return state.Iteration.Item, nil
			case "$index":
				return state.Iteration.Index, nil
			default:
				return state.Iteration.BatchIndex, nil
			}
		}
	}
	if value, ok := r.expressions[expression]; ok {
		return value, nil
	}
	if jsonNumberPattern.MatchString(expression) {
		return json.Number(expression), nil
	}
	return nil, fmt.Errorf("scripted runtime has no value for expression %q", expression)
}

func (r *conformanceScriptRuntime) ResolveItems(_ context.Context, expression string) ([]any, error) {
	items, ok := r.items[expression]
	if !ok {
		return nil, fmt.Errorf("scripted runtime has no items for expression %q", expression)
	}
	return items, nil
}
