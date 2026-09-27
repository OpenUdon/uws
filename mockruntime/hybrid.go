package mockruntime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/OpenUdon/uws/internal/strictjson"
	"github.com/OpenUdon/uws/uws1"
)

// ReadDelegate executes one explicitly classified read using caller-owned
// transport and credentials. Implementations must treat operation and request
// as read-only and must be safe for concurrent calls when used with parallel
// UWS branches.
type ReadDelegate interface {
	ExecuteRead(ctx context.Context, operation *uws1.Operation, request map[string]any) (json.RawMessage, error)
}

// ReadDelegateFunc adapts a function to ReadDelegate.
type ReadDelegateFunc func(context.Context, *uws1.Operation, map[string]any) (json.RawMessage, error)

func (f ReadDelegateFunc) ExecuteRead(ctx context.Context, operation *uws1.Operation, request map[string]any) (json.RawMessage, error) {
	return f(ctx, operation, request)
}

// HybridOptions require an explicit opt-in before any read delegate can run.
type HybridOptions struct {
	AllowLiveReads bool
	ReadDelegate   ReadDelegate
}

// HybridRuntime delegates explicitly classified read operations to a caller-
// supplied adapter and keeps writes, unknown effects, and omitted effects on
// the pure mock runtime.
type HybridRuntime struct {
	mock           *Runtime
	delegate       ReadDelegate
	readOperations map[*uws1.Operation]struct{}
}

// NewHybridRuntime creates a runtime that may perform live reads. Both the
// explicit AllowLiveReads opt-in and a read delegate are required. Every
// operation with an effect other than exactly "read" stays mocked.
func NewHybridRuntime(mock *Runtime, options HybridOptions) (*HybridRuntime, error) {
	if mock == nil || mock.document == nil {
		return nil, fmt.Errorf("hybrid runtime requires an initialized mock runtime")
	}
	if !versionAtLeast(mock.document.UWS, 12, 0) {
		return nil, fmt.Errorf("hybrid live reads require a UWS 1.12.0-or-later document with operation effects")
	}
	if !options.AllowLiveReads {
		return nil, fmt.Errorf("hybrid runtime requires explicit AllowLiveReads opt-in")
	}
	if nilReadDelegate(options.ReadDelegate) {
		return nil, fmt.Errorf("hybrid runtime requires a read delegate")
	}
	reads := make(map[*uws1.Operation]struct{})
	for _, operation := range mock.document.Operations {
		if operation != nil && operation.Effect == uws1.OperationEffectRead {
			reads[operation] = struct{}{}
		}
	}
	return &HybridRuntime{mock: mock, delegate: options.ReadDelegate, readOperations: reads}, nil
}

func nilReadDelegate(delegate ReadDelegate) bool {
	if delegate == nil {
		return true
	}
	value := reflect.ValueOf(delegate)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

// ExecuteLeaf implements uws1.Runtime. The orchestrator uses
// ExecuteLeafWithResult when available to expose successful delegated
// responses to criteria and outputs.
func (r *HybridRuntime) ExecuteLeaf(ctx context.Context, operation *uws1.Operation) error {
	_, err := r.ExecuteLeafWithResult(ctx, operation)
	return err
}

// ExecuteLeafWithResult sends only effect="read" operations to the enabled
// delegate. Writes, unknown effects, and omitted effects use the pure mock.
func (r *HybridRuntime) ExecuteLeafWithResult(ctx context.Context, operation *uws1.Operation) (any, error) {
	if r == nil || r.mock == nil || r.delegate == nil {
		return nil, fmt.Errorf("hybrid runtime is not initialized")
	}
	if operation == nil {
		return nil, fmt.Errorf("hybrid runtime requires an operation")
	}
	if operation.Effect != uws1.OperationEffectRead {
		return r.mock.ExecuteLeafWithResult(ctx, operation)
	}
	if _, declaredRead := r.readOperations[operation]; !declaredRead {
		return r.mock.ExecuteLeafWithResult(ctx, operation)
	}
	request, digest, requestIndex, err := r.mock.prepareRequest(ctx, operation)
	if err != nil {
		return nil, err
	}
	// The record identifies the selected execution path even when the delegate
	// later returns an error or observes cancellation.
	r.mock.finishRequest(requestIndex, ResponseLiveRead)
	response, err := r.delegate.ExecuteRead(ctx, operation, request)
	if err != nil {
		return nil, fmt.Errorf("delegate live read for operation %q (request %s): %w", operation.OperationID, digest, err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(response) > MaxFixtureSetBytes {
		return nil, fmt.Errorf("live response for operation %q exceeds %d bytes", operation.OperationID, MaxFixtureSetBytes)
	}
	if err := strictjson.ValidateSingleValue(response); err != nil {
		return nil, fmt.Errorf("live response for operation %q is invalid: %w", operation.OperationID, err)
	}
	return json.RawMessage(bytes.Clone(response)), nil
}

// EvaluateExpression uses the same mock expression evaluator for both
// simulated and delegated responses.
func (r *HybridRuntime) EvaluateExpression(ctx context.Context, expression string) (any, error) {
	if r == nil || r.mock == nil {
		return nil, fmt.Errorf("hybrid runtime is not initialized")
	}
	return r.mock.EvaluateExpression(ctx, expression)
}

// ResolveItems uses the same ordered item resolution as the pure mock runtime.
func (r *HybridRuntime) ResolveItems(ctx context.Context, expression string) ([]any, error) {
	if r == nil || r.mock == nil {
		return nil, fmt.Errorf("hybrid runtime is not initialized")
	}
	return r.mock.ResolveItems(ctx, expression)
}

// RequestRecords returns the mock runtime's defensive snapshot of simulated
// and delegated would-be requests.
func (r *HybridRuntime) RequestRecords() []RequestRecord {
	if r == nil || r.mock == nil {
		return nil
	}
	return r.mock.RequestRecords()
}
