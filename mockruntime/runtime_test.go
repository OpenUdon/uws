package mockruntime

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestMockRuntimeExecutesThroughOrchestratorWithFixtureAndResponseExpressions(t *testing.T) {
	document := mockRuntimeDocument()
	request := map[string]any{"query": map[string]any{"name": "Ada"}}
	digest, err := RequestDigest(request)
	if err != nil {
		t.Fatal(err)
	}
	fixtures := &FixtureSet{Format: FixtureFormatV1, Fixtures: []Fixture{{
		OperationID:   "read_person",
		RequestDigest: digest,
		Provenance:    FixtureProvenance{Kind: "example"},
		Response:      json.RawMessage(`{"statusCode":200,"headers":{"content-type":"application/json"},"body":{"id":12,"name":"Ada"}}`),
	}}}
	runtime, err := NewRuntime(document, Options{Fixtures: fixtures})
	if err != nil {
		t.Fatal(err)
	}
	document.SetRuntime(runtime)
	if err := document.Execute(context.Background()); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	requests := runtime.RequestRecords()
	if len(requests) != 1 {
		t.Fatalf("request records = %#v, want one resolved request", requests)
	}
	if requests[0].RequestDigest != digest || requests[0].ResponseKind != ResponseFixture || !json.Valid(requests[0].Request) {
		t.Fatalf("request record = %#v", requests[0])
	}
	var recordedRequest map[string]any
	if err := json.Unmarshal(requests[0].Request, &recordedRequest); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(recordedRequest, request) {
		t.Fatalf("resolved request = %#v, want %#v", recordedRequest, request)
	}
	records := document.ExecutionRecords()
	operation := records["stepop:read_person_step:read_person"]
	if operation.Status != "success" || !json.Valid(operation.Result.(json.RawMessage)) {
		t.Fatalf("operation record = %#v", operation)
	}
	if operation.Outputs["person_id"] != json.Number("12") {
		t.Fatalf("operation output = %#v, want id 12", operation.Outputs)
	}
	step := records["step:read_person_step"]
	if step.Outputs["person_name"] != "Ada" {
		t.Fatalf("step output = %#v, want Ada", step.Outputs)
	}
}

func TestMockRuntimeResolvesAndRecordsWouldBeRequests(t *testing.T) {
	document := mockRuntimeDocument()
	document.Variables = map[string]any{"person": "Ada"}
	resolverCalls := 0
	runtime, err := NewRuntime(document, Options{ResponseResolver: ResponseResolverFunc(func(_ context.Context, operation *uws1.Operation) (ResponseDefinition, error) {
		resolverCalls++
		if operation.OperationID != "read_person" {
			t.Fatalf("resolver operation = %q", operation.OperationID)
		}
		return ResponseDefinition{Example: json.RawMessage(`{"statusCode":200,"body":{"name":"Ada"}}`)}, nil
	})})
	if err != nil {
		t.Fatal(err)
	}
	operation := &uws1.Operation{OperationID: "read_person", Request: map[string]any{
		"query": map[string]any{"name": "$variables.person"},
		"body":  map[string]any{"id": "$inputs.item.id"},
	}}
	ctx := uws1.WithExecutionContext(context.Background(), &uws1.ExecutionContext{Inputs: map[string]any{"item": map[string]any{"id": json.Number("7")}}})
	response, err := runtime.ExecuteLeafWithResult(ctx, operation)
	if err != nil {
		t.Fatal(err)
	}
	if string(response.(json.RawMessage)) != `{"statusCode":200,"body":{"name":"Ada"}}` || resolverCalls != 1 {
		t.Fatalf("response = %s, resolver calls = %d", response, resolverCalls)
	}
	requests := runtime.RequestRecords()
	if len(requests) != 1 || requests[0].ResponseKind != ResponseExample {
		t.Fatalf("request records = %#v", requests)
	}
	wantDigest, err := RequestDigest(map[string]any{
		"query": map[string]any{"name": "Ada"},
		"body":  map[string]any{"id": json.Number("7")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if requests[0].RequestDigest != wantDigest {
		t.Fatalf("request digest = %s, want %s", requests[0].RequestDigest, wantDigest)
	}
	requests[0].Request[0] = ' '
	if !json.Valid(runtime.RequestRecords()[0].Request) {
		t.Fatal("mutating a returned request record changed runtime history")
	}
}

func TestMockRuntimeFixtureMissNeedsExplicitGeneratedFallback(t *testing.T) {
	document := mockRuntimeDocument()
	digest, err := RequestDigest(map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	fixtures := &FixtureSet{Format: FixtureFormatV1, Fixtures: []Fixture{{
		OperationID:   "other_operation",
		RequestDigest: digest,
		Provenance:    FixtureProvenance{Kind: "synthetic"},
		Response:      json.RawMessage(`null`),
	}}}
	resolver := ResponseResolverFunc(func(context.Context, *uws1.Operation) (ResponseDefinition, error) {
		return ResponseDefinition{Example: json.RawMessage(`{"generated":true}`)}, nil
	})
	operation := &uws1.Operation{OperationID: "read_person"}

	strictRuntime, err := NewRuntime(document, Options{Fixtures: fixtures, ResponseResolver: resolver})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := strictRuntime.ExecuteLeafWithResult(context.Background(), operation); err == nil {
		t.Fatal("fixture miss silently fell back to a generated example")
	}
	if records := strictRuntime.RequestRecords(); len(records) != 1 || records[0].ResponseKind != ResponseUnavailable {
		t.Fatalf("fixture-miss request record = %#v", records)
	}

	fallbackRuntime, err := NewRuntime(document, Options{Fixtures: fixtures, ResponseResolver: resolver, AllowGeneratedFallback: true})
	if err != nil {
		t.Fatal(err)
	}
	response, err := fallbackRuntime.ExecuteLeafWithResult(context.Background(), operation)
	if err != nil {
		t.Fatal(err)
	}
	if string(response.(json.RawMessage)) != `{"generated":true}` || fallbackRuntime.RequestRecords()[0].ResponseKind != ResponseExample {
		t.Fatalf("fallback response = %s; request = %#v", response, fallbackRuntime.RequestRecords()[0])
	}
}

func TestMockRuntimeRejectsExcessivelyNestedRequests(t *testing.T) {
	document := mockRuntimeDocument()
	runtime, err := NewRuntime(document, Options{})
	if err != nil {
		t.Fatal(err)
	}
	var value any = "end"
	for i := 0; i <= maxRequestValueDepth; i++ {
		value = map[string]any{"nested": value}
	}
	if _, err := runtime.resolveRequest(context.Background(), map[string]any{"body": value}); err == nil {
		t.Fatalf("request deeper than %d was accepted", maxRequestValueDepth)
	}
}

func TestMockRuntimeUsesExampleBeforeSchemaAndSynthesizesSupportedSchemas(t *testing.T) {
	document := mockRuntimeDocument()
	schema := &uws1.ParamSchema{Type: "object", Properties: map[string]*uws1.ParamSchema{
		"id":   {Type: "integer"},
		"tags": {Type: "array", Items: &uws1.ParamSchema{Type: "string"}},
	}, Required: []string{"id"}}
	operation := &uws1.Operation{OperationID: "read_person"}
	runtime, err := NewRuntime(document, Options{ResponseResolver: ResponseResolverFunc(func(context.Context, *uws1.Operation) (ResponseDefinition, error) {
		return ResponseDefinition{Schema: schema}, nil
	})})
	if err != nil {
		t.Fatal(err)
	}
	response, err := runtime.ExecuteLeafWithResult(context.Background(), operation)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(response.(json.RawMessage)), `{"id":0,"tags":[""]}`; got != want {
		t.Fatalf("synthesized response = %s, want %s", got, want)
	}
	if runtime.RequestRecords()[0].ResponseKind != ResponseSynthesized {
		t.Fatalf("response kind = %q, want synthesized", runtime.RequestRecords()[0].ResponseKind)
	}
	unsupported, err := NewRuntime(document, Options{ResponseResolver: ResponseResolverFunc(func(context.Context, *uws1.Operation) (ResponseDefinition, error) {
		return ResponseDefinition{Schema: &uws1.ParamSchema{Ref: "#/components/schemas/Person"}}, nil
	})})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := unsupported.ExecuteLeafWithResult(context.Background(), operation); err == nil {
		t.Fatal("$ref synthesis was accepted without a resolved schema or example")
	}
}

func mockRuntimeDocument() *uws1.Document {
	return &uws1.Document{
		UWS:  "1.12.0",
		Info: &uws1.Info{Title: "mock runtime test", Version: "1.0.0"},
		Operations: []*uws1.Operation{{
			OperationID: "read_person",
			Extensions:  map[string]any{uws1.ExtensionOperationProfile: "test-profile"},
			Request:     map[string]any{"query": map[string]any{"name": "Ada"}},
			SuccessCriteria: []*uws1.Criterion{{
				Condition: "$response.statusCode == 200",
			}},
			Outputs: map[string]string{"person_id": "$response.body.id"},
		}},
		Workflows: []*uws1.Workflow{{
			WorkflowID: "main",
			Type:       uws1.WorkflowTypeSequence,
			Steps: []*uws1.Step{{
				StepID:       "read_person_step",
				OperationRef: "read_person",
				Outputs:      map[string]string{"person_name": "$response.body.name"},
			}},
		}},
	}
}
