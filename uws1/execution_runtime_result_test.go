package uws1

import (
	"context"
	"encoding/json"
	"testing"
)

func TestRuntimeWithResultFeedsResponseExpressionsAndStepRecord(t *testing.T) {
	doc := &Document{
		UWS:  "1.12.0",
		Info: &Info{Title: "mock result", Version: "1.0.0"},
		Operations: []*Operation{{
			OperationID: "read_item",
			Extensions:  map[string]any{ExtensionOperationProfile: "test"},
			SuccessCriteria: []*Criterion{{
				Condition: "$response.statusCode == 200",
			}},
		}},
		Workflows: []*Workflow{{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps: []*Step{{
				StepID:       "read_step",
				OperationRef: "read_item",
				Outputs:      map[string]string{"name": "$response.body.name"},
			}},
		}},
	}
	runtime := &resultTestRuntime{response: json.RawMessage(`{"statusCode":200,"headers":{"content-type":"application/json"},"body":{"name":"fixture"}}`)}
	doc.SetRuntime(runtime)
	if err := doc.Execute(context.Background()); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !runtime.sawCurrent {
		t.Fatal("leaf runtime did not receive a current execution context")
	}
	if !runtime.sawResponse {
		t.Fatal("success criteria did not see the leaf response in execution records")
	}
	records := doc.ExecutionRecords()
	step, ok := records["step:read_step"]
	if !ok {
		t.Fatalf("missing step execution record: %#v", records)
	}
	if string(step.Result.(json.RawMessage)) != string(runtime.response) {
		t.Fatalf("step result = %s, want %s", step.Result, runtime.response)
	}
	if step.Outputs["name"] != "fixture" {
		t.Fatalf("step output name = %#v, want fixture", step.Outputs["name"])
	}
	runtime.response[0] = ' '
	if got := string(doc.ExecutionRecords()["step:read_step"].Result.(json.RawMessage)); got != `{"statusCode":200,"headers":{"content-type":"application/json"},"body":{"name":"fixture"}}` {
		t.Fatalf("mutating runtime response changed the stored result: %s", got)
	}
}

type resultTestRuntime struct {
	response    json.RawMessage
	sawCurrent  bool
	sawResponse bool
}

func (r *resultTestRuntime) ExecuteLeaf(context.Context, *Operation) error { return nil }

func (r *resultTestRuntime) ExecuteLeafWithResult(ctx context.Context, _ *Operation) (any, error) {
	state, _ := ExecutionContextFromContext(ctx)
	r.sawCurrent = state != nil && state.Current != nil && state.Current.ID == "read_item"
	return r.response, nil
}

func (r *resultTestRuntime) EvaluateExpression(ctx context.Context, expression string) (any, error) {
	state, _ := ExecutionContextFromContext(ctx)
	if state != nil && state.Current != nil && state.Current.Key != "" {
		if record, ok := state.Records[state.Current.Key]; ok {
			r.sawResponse = record.Result != nil
		}
	}
	if expression == "$response.statusCode == 200" {
		return true, nil
	}
	if expression == "$response.body.name" {
		return "fixture", nil
	}
	return nil, nil
}

func (*resultTestRuntime) ResolveItems(context.Context, string) ([]any, error) { return nil, nil }
