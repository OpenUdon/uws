package mockruntime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/OpenUdon/uws/internal/strictjson"
	"github.com/OpenUdon/uws/uws1"
)

const (
	maxRequestValueDepth = 64

	// MaxRequestRecords bounds the in-memory would-be request history per
	// Runtime instance.
	MaxRequestRecords = 10000

	// MaxRequestLogBytes bounds the canonical request bytes retained by one
	// Runtime instance.
	MaxRequestLogBytes = 16 << 20
)

// ResponseResolver supplies source-owned response examples or schemas. It
// must be safe for concurrent calls when its Runtime is used by parallel UWS
// branches.
type ResponseResolver interface {
	ResolveResponse(ctx context.Context, operation *uws1.Operation) (ResponseDefinition, error)
}

// ResponseResolverFunc adapts a function to ResponseResolver.
type ResponseResolverFunc func(context.Context, *uws1.Operation) (ResponseDefinition, error)

func (f ResponseResolverFunc) ResolveResponse(ctx context.Context, operation *uws1.Operation) (ResponseDefinition, error) {
	return f(ctx, operation)
}

// ResponseDefinition carries a caller-supplied JSON example or the schema used
// for deterministic synthesis. Example takes precedence when both are present.
type ResponseDefinition struct {
	Example json.RawMessage
	Schema  *uws1.ParamSchema
}

// Options configure fixture replay and caller-supplied response generation.
type Options struct {
	// Fixtures, when supplied, are checked before ResponseResolver. A missing
	// exact key is an error unless AllowGeneratedFallback is explicitly true.
	Fixtures *FixtureSet
	// ResponseResolver supplies a JSON example or schema when no fixture is
	// selected.
	ResponseResolver ResponseResolver
	// AllowGeneratedFallback allows the resolver to supply a response after an
	// exact fixture miss. It never changes an exact fixture match.
	AllowGeneratedFallback bool
}

// ResponseKind identifies how one leaf response was selected.
type ResponseKind string

const (
	ResponseFixture     ResponseKind = "fixture"
	ResponseExample     ResponseKind = "example"
	ResponseSynthesized ResponseKind = "synthesized"
	ResponseUnavailable ResponseKind = "unavailable"
)

// RequestRecord is an in-memory record of one resolved would-be request.
// Request contains canonical JSON and may include private request-bound data;
// callers must review and redact it before any export. It is never persisted by
// this package.
type RequestRecord struct {
	OperationID   string          `json:"operationId"`
	RequestDigest string          `json:"requestDigest"`
	Request       json.RawMessage `json:"request"`
	ResponseKind  ResponseKind    `json:"responseKind"`
}

// Runtime implements uws1.Runtime and uws1.RuntimeWithResult. It never makes
// network calls or writes to disk.
type Runtime struct {
	document               *uws1.Document
	fixtures               *FixtureSet
	responseResolver       ResponseResolver
	allowGeneratedFallback bool

	mu             sync.RWMutex
	requests       []RequestRecord
	requestLogSize int
}

// NewRuntime returns a mock runtime bound to document. Call document.SetRuntime
// with the result before executing the document.
func NewRuntime(document *uws1.Document, options Options) (*Runtime, error) {
	if document == nil {
		return nil, fmt.Errorf("mock runtime requires a UWS document")
	}
	runtime := &Runtime{
		document:               document,
		responseResolver:       options.ResponseResolver,
		allowGeneratedFallback: options.AllowGeneratedFallback,
	}
	if options.Fixtures != nil {
		encoded, err := EncodeFixtures(options.Fixtures)
		if err != nil {
			return nil, fmt.Errorf("validate mock runtime fixtures: %w", err)
		}
		runtime.fixtures, err = DecodeFixtures(encoded)
		if err != nil {
			return nil, fmt.Errorf("copy mock runtime fixtures: %w", err)
		}
	}
	return runtime, nil
}

// ExecuteLeaf implements the existing Runtime interface. The orchestrator
// uses ExecuteLeafWithResult when available so the response also enters the
// current execution record.
func (r *Runtime) ExecuteLeaf(ctx context.Context, operation *uws1.Operation) error {
	_, err := r.ExecuteLeafWithResult(ctx, operation)
	return err
}

// ExecuteLeafWithResult resolves request expressions, records the would-be
// request, and selects an exact fixture or caller-supplied generated response.
func (r *Runtime) ExecuteLeafWithResult(ctx context.Context, operation *uws1.Operation) (any, error) {
	if r == nil || r.document == nil {
		return nil, fmt.Errorf("mock runtime is not initialized")
	}
	if ctx == nil {
		return nil, fmt.Errorf("mock runtime requires a context")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if operation == nil {
		return nil, fmt.Errorf("mock runtime requires an operation")
	}
	request, err := r.resolveRequest(ctx, operation.Request)
	if err != nil {
		return nil, fmt.Errorf("resolve request for operation %q: %w", operation.OperationID, err)
	}
	canonical, err := CanonicalizeRequest(request)
	if err != nil {
		return nil, fmt.Errorf("canonicalize request for operation %q: %w", operation.OperationID, err)
	}
	digest, err := RequestDigest(request)
	if err != nil {
		return nil, fmt.Errorf("digest request for operation %q: %w", operation.OperationID, err)
	}
	requestIndex, err := r.recordRequest(RequestRecord{
		OperationID:   operation.OperationID,
		RequestDigest: digest,
		Request:       bytes.Clone(canonical),
		ResponseKind:  ResponseUnavailable,
	})
	if err != nil {
		return nil, err
	}
	response, kind, err := r.selectResponse(ctx, operation, digest)
	r.finishRequest(requestIndex, kind)
	if err != nil {
		return nil, err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return json.RawMessage(bytes.Clone(response)), nil
}

// RequestRecords returns a defensive snapshot of this runtime's in-memory
// would-be request history. Records span executions made with the same
// Runtime instance and are not automatically persisted.
func (r *Runtime) RequestRecords() []RequestRecord {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]RequestRecord, len(r.requests))
	for i, record := range r.requests {
		result[i] = record
		result[i].Request = bytes.Clone(record.Request)
	}
	return result
}

func (r *Runtime) resolveRequest(ctx context.Context, request map[string]any) (map[string]any, error) {
	if request == nil {
		return map[string]any{}, nil
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode request bindings: %w", err)
	}
	if err := strictjson.ValidateSingleValue(encoded); err != nil {
		return nil, fmt.Errorf("validate request bindings: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode request bindings: %w", err)
	}
	resolved, err := r.resolveRequestValue(ctx, value, 0)
	if err != nil {
		return nil, err
	}
	object, ok := resolved.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("request bindings must resolve to a JSON object")
	}
	return object, nil
}

func (r *Runtime) resolveRequestValue(ctx context.Context, value any, depth int) (any, error) {
	if depth > maxRequestValueDepth {
		return nil, fmt.Errorf("mock request/input nesting exceeds %d", maxRequestValueDepth)
	}
	switch typed := value.(type) {
	case string:
		if len(typed) > 0 && typed[0] == '$' {
			return r.EvaluateExpression(ctx, typed)
		}
		return typed, nil
	case []any:
		resolved := make([]any, len(typed))
		for i, item := range typed {
			value, err := r.resolveRequestValue(ctx, item, depth+1)
			if err != nil {
				return nil, fmt.Errorf("request array item %d: %w", i, err)
			}
			resolved[i] = value
		}
		return resolved, nil
	case map[string]any:
		resolved := make(map[string]any, len(typed))
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value, err := r.resolveRequestValue(ctx, typed[key], depth+1)
			if err != nil {
				return nil, fmt.Errorf("request property %q: %w", key, err)
			}
			resolved[key] = value
		}
		return resolved, nil
	default:
		return typed, nil
	}
}

func (r *Runtime) selectResponse(ctx context.Context, operation *uws1.Operation, digest string) (json.RawMessage, ResponseKind, error) {
	if r.fixtures != nil {
		response, found, err := r.fixtures.LookupResponse(operation.OperationID, digest)
		if err != nil {
			return nil, ResponseUnavailable, err
		}
		if found {
			return response, ResponseFixture, nil
		}
		if !r.allowGeneratedFallback {
			return nil, ResponseUnavailable, fmt.Errorf("no mock fixture matches operation %q and request digest %q", operation.OperationID, digest)
		}
	}
	if r.responseResolver == nil {
		return nil, ResponseUnavailable, fmt.Errorf("no response fixture or response resolver is available for operation %q", operation.OperationID)
	}
	definition, err := r.responseResolver.ResolveResponse(ctx, operation)
	if err != nil {
		return nil, ResponseUnavailable, fmt.Errorf("resolve response definition for operation %q: %w", operation.OperationID, err)
	}
	if len(definition.Example) > 0 {
		if len(definition.Example) > MaxFixtureSetBytes {
			return nil, ResponseUnavailable, fmt.Errorf("response example for operation %q exceeds %d bytes", operation.OperationID, MaxFixtureSetBytes)
		}
		if err := strictjson.ValidateSingleValue(definition.Example); err != nil {
			return nil, ResponseUnavailable, fmt.Errorf("response example for operation %q is invalid: %w", operation.OperationID, err)
		}
		return bytes.Clone(definition.Example), ResponseExample, nil
	}
	if definition.Schema == nil {
		return nil, ResponseUnavailable, fmt.Errorf("response resolver for operation %q supplied neither an example nor a schema", operation.OperationID)
	}
	value, err := synthesizeResponse(definition.Schema)
	if err != nil {
		return nil, ResponseUnavailable, fmt.Errorf("synthesize response for operation %q: %w", operation.OperationID, err)
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, ResponseUnavailable, fmt.Errorf("encode synthesized response for operation %q: %w", operation.OperationID, err)
	}
	if len(encoded) > MaxFixtureSetBytes {
		return nil, ResponseUnavailable, fmt.Errorf("synthesized response for operation %q exceeds %d bytes", operation.OperationID, MaxFixtureSetBytes)
	}
	return encoded, ResponseSynthesized, nil
}

func (r *Runtime) recordRequest(record RequestRecord) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.requests) >= MaxRequestRecords {
		return 0, fmt.Errorf("mock runtime request history exceeds %d records", MaxRequestRecords)
	}
	if len(record.Request) > MaxRequestLogBytes-r.requestLogSize {
		return 0, fmt.Errorf("mock runtime request history exceeds %d bytes", MaxRequestLogBytes)
	}
	r.requestLogSize += len(record.Request)
	r.requests = append(r.requests, record)
	return len(r.requests) - 1, nil
}

func (r *Runtime) finishRequest(index int, kind ResponseKind) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if index >= 0 && index < len(r.requests) {
		r.requests[index].ResponseKind = kind
	}
}
