package mockruntime

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestHybridRuntimeDelegatesOnlyReadEffectsAndMocksTheRest(t *testing.T) {
	document := hybridWorkflowDocument()
	readRequest := map[string]any{"query": map[string]any{"id": "item-1"}}
	writeRequest := map[string]any{"body": map[string]any{"id": json.Number("41")}}
	fixtures := &FixtureSet{Format: FixtureFormatV1}
	for operationID, request := range map[string]any{
		"read_item":      readRequest,
		"write_item":     writeRequest,
		"unknown_item":   map[string]any{},
		"omitted_effect": map[string]any{},
	} {
		digest, err := RequestDigest(request)
		if err != nil {
			t.Fatal(err)
		}
		fixtures.Fixtures = append(fixtures.Fixtures, Fixture{
			OperationID:   operationID,
			RequestDigest: digest,
			Provenance:    FixtureProvenance{Kind: "example"},
			Response:      json.RawMessage(`{"statusCode":200,"body":{"accepted":true}}`),
		})
	}
	mock, err := NewRuntime(document, Options{
		Fixtures:               fixtures,
		AllowGeneratedFallback: true,
		ResponseResolver: ResponseResolverFunc(func(context.Context, *uws1.Operation) (ResponseDefinition, error) {
			return ResponseDefinition{Schema: &uws1.ParamSchema{Type: "object"}}, nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	var delegated []string
	hybrid, err := NewHybridRuntime(mock, HybridOptions{
		AllowLiveReads: true,
		ReadDelegate: ReadDelegateFunc(func(_ context.Context, operation *uws1.Operation, request map[string]any) (json.RawMessage, error) {
			delegated = append(delegated, operation.OperationID)
			if operation.Effect != uws1.OperationEffectRead {
				t.Errorf("delegate received non-read operation %q with effect %q", operation.OperationID, operation.Effect)
			}
			query, ok := request["query"].(map[string]any)
			if !ok {
				t.Errorf("resolved read query = %#v, want object", request["query"])
			} else if got := query["id"]; got != "item-1" {
				t.Errorf("resolved read request id = %#v, want item-1", got)
			}
			return json.RawMessage(`{"statusCode":200,"body":{"id":41}}`), nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	document.SetRuntime(hybrid)
	if err := document.Execute(context.Background()); err != nil {
		t.Fatalf("hybrid Execute() error = %v", err)
	}
	if !reflect.DeepEqual(delegated, []string{"read_item"}) {
		t.Fatalf("delegated operations = %#v, want only read_item", delegated)
	}
	if got := document.ExecutionRecords()["step:write_step"].Outputs["accepted"]; got != true {
		t.Fatalf("mocked write output = %#v, want true", got)
	}

	readCopy := *document.Operations[0]
	if _, err := hybrid.ExecuteLeafWithResult(context.Background(), &readCopy); err != nil {
		t.Fatalf("unbound read copy should use mock fixture: %v", err)
	}
	for _, operation := range document.Operations[2:] {
		if _, err := hybrid.ExecuteLeafWithResult(context.Background(), operation); err != nil {
			t.Fatalf("mocked operation %q: %v", operation.OperationID, err)
		}
	}
	if !reflect.DeepEqual(delegated, []string{"read_item"}) {
		t.Fatalf("unknown or omitted-effect operations reached delegate: %#v", delegated)
	}
	records := hybrid.RequestRecords()
	if len(records) != 6 {
		t.Fatalf("request record count = %d, want six", len(records))
	}
	wantKinds := []ResponseKind{ResponseLiveRead, ResponseFixture, ResponseFixture, ResponseFixture, ResponseFixture, ResponseSynthesized}
	for index, want := range wantKinds {
		if records[index].ResponseKind != want {
			t.Errorf("request %d kind = %q, want %q", index, records[index].ResponseKind, want)
		}
	}
}

func TestHybridRuntimeRequiresExplicitEnablementAndDelegate(t *testing.T) {
	mock, err := NewRuntime(hybridWorkflowDocument(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	delegate := ReadDelegateFunc(func(context.Context, *uws1.Operation, map[string]any) (json.RawMessage, error) {
		t.Fatal("delegate called during disabled hybrid construction")
		return nil, nil
	})
	if _, err := NewHybridRuntime(mock, HybridOptions{ReadDelegate: delegate}); err == nil {
		t.Fatal("hybrid runtime without AllowLiveReads opt-in was accepted")
	}
	if _, err := NewHybridRuntime(mock, HybridOptions{AllowLiveReads: true}); err == nil {
		t.Fatal("hybrid runtime without a delegate was accepted")
	}
	var nilDelegate ReadDelegateFunc
	if _, err := NewHybridRuntime(mock, HybridOptions{AllowLiveReads: true, ReadDelegate: nilDelegate}); err == nil {
		t.Fatal("hybrid runtime accepted a typed nil delegate")
	}
	legacyDocument := singleReadDocument()
	legacyDocument.UWS = "1.11.0"
	legacyMock, err := NewRuntime(legacyDocument, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewHybridRuntime(legacyMock, HybridOptions{AllowLiveReads: true, ReadDelegate: delegate}); err == nil {
		t.Fatal("hybrid runtime accepted a document version without operation effects")
	}
}

func TestHybridRuntimePropagatesCancellationAndDelegateErrors(t *testing.T) {
	t.Run("cancellation", func(t *testing.T) {
		document := singleReadDocument()
		mock, err := NewRuntime(document, Options{})
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		hybrid, err := NewHybridRuntime(mock, HybridOptions{
			AllowLiveReads: true,
			ReadDelegate: ReadDelegateFunc(func(context.Context, *uws1.Operation, map[string]any) (json.RawMessage, error) {
				cancel()
				return json.RawMessage(`{"statusCode":200,"body":{}}`), nil
			}),
		})
		if err != nil {
			t.Fatal(err)
		}
		document.SetRuntime(hybrid)
		if err := document.Execute(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("Execute() error = %v, want context.Canceled", err)
		}
		if got := hybrid.RequestRecords()[0].ResponseKind; got != ResponseLiveRead {
			t.Fatalf("canceled read record kind = %q, want live-read attempt", got)
		}
	})

	t.Run("delegate error", func(t *testing.T) {
		document := singleReadDocument()
		mock, err := NewRuntime(document, Options{})
		if err != nil {
			t.Fatal(err)
		}
		wantErr := errors.New("read transport failed")
		hybrid, err := NewHybridRuntime(mock, HybridOptions{
			AllowLiveReads: true,
			ReadDelegate: ReadDelegateFunc(func(context.Context, *uws1.Operation, map[string]any) (json.RawMessage, error) {
				return nil, wantErr
			}),
		})
		if err != nil {
			t.Fatal(err)
		}
		document.SetRuntime(hybrid)
		if err := document.Execute(context.Background()); !errors.Is(err, wantErr) {
			t.Fatalf("Execute() error = %v, want wrapped delegate error", err)
		}
		if got := hybrid.RequestRecords()[0].ResponseKind; got != ResponseLiveRead {
			t.Fatalf("failed read record kind = %q, want live-read attempt", got)
		}
	})
}

func TestHybridRuntimeRejectsInvalidLiveResponse(t *testing.T) {
	document := singleReadDocument()
	mock, err := NewRuntime(document, Options{})
	if err != nil {
		t.Fatal(err)
	}
	hybrid, err := NewHybridRuntime(mock, HybridOptions{
		AllowLiveReads: true,
		ReadDelegate: ReadDelegateFunc(func(context.Context, *uws1.Operation, map[string]any) (json.RawMessage, error) {
			return json.RawMessage(`{"statusCode":200} {"extra":true}`), nil
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := hybrid.ExecuteLeafWithResult(context.Background(), document.Operations[0]); err == nil {
		t.Fatal("live delegate response with trailing JSON was accepted")
	}
	if got := hybrid.RequestRecords()[0].ResponseKind; got != ResponseLiveRead {
		t.Fatalf("invalid live response record kind = %q, want live-read attempt", got)
	}
}

func hybridWorkflowDocument() *uws1.Document {
	read := &uws1.Operation{
		OperationID: "read_item",
		Effect:      uws1.OperationEffectRead,
		Extensions:  map[string]any{uws1.ExtensionOperationProfile: "test-profile"},
		Request:     map[string]any{"query": map[string]any{"id": "item-1"}},
		SuccessCriteria: []*uws1.Criterion{{
			Condition: "$response.statusCode == 200",
		}},
	}
	write := &uws1.Operation{
		OperationID: "write_item",
		Effect:      uws1.OperationEffectWrite,
		Extensions:  map[string]any{uws1.ExtensionOperationProfile: "test-profile"},
		Request:     map[string]any{"body": map[string]any{"id": "$steps.read_step.outputs.id"}},
		SuccessCriteria: []*uws1.Criterion{{
			Condition: "$response.statusCode == 200",
		}},
	}
	unknown := &uws1.Operation{
		OperationID: "unknown_item",
		Effect:      uws1.OperationEffectUnknown,
		Extensions:  map[string]any{uws1.ExtensionOperationProfile: "test-profile"},
	}
	omitted := &uws1.Operation{
		OperationID: "omitted_effect",
		Extensions:  map[string]any{uws1.ExtensionOperationProfile: "test-profile"},
	}
	synthesized := &uws1.Operation{
		OperationID: "synthesized_item",
		Effect:      uws1.OperationEffectUnknown,
		Extensions:  map[string]any{uws1.ExtensionOperationProfile: "test-profile"},
	}
	return &uws1.Document{
		UWS:        "1.12.0",
		Info:       &uws1.Info{Title: "hybrid mock test", Version: "1.0.0"},
		Operations: []*uws1.Operation{read, write, unknown, omitted, synthesized},
		Workflows: []*uws1.Workflow{{
			WorkflowID: "main",
			Type:       uws1.WorkflowTypeSequence,
			Steps: []*uws1.Step{
				{StepID: "read_step", OperationRef: "read_item", Outputs: map[string]string{"id": "$response.body.id"}},
				{StepID: "write_step", OperationRef: "write_item", Outputs: map[string]string{"accepted": "$response.body.accepted"}},
			},
		}},
	}
}

func singleReadDocument() *uws1.Document {
	document := hybridWorkflowDocument()
	document.Operations = document.Operations[:1]
	document.Workflows[0].Steps = document.Workflows[0].Steps[:1]
	document.Workflows[0].Steps[0].Outputs = nil
	return document
}
