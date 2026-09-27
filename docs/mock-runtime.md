# Public Mock Runtime

The public `github.com/OpenUdon/uws/mockruntime` package runs a UWS document
through the normal core orchestrator while returning simulated responses for
leaf operations. Its `Runtime` makes no transport calls and performs no file
writes. A caller-supplied `ResponseResolver` is invoked as a local data hook;
keep that hook free of network and persistence behavior when using the runtime
for pure simulation. Its separate `HybridRuntime` can hand eligible reads to
a caller-supplied real-read delegate after explicit opt-in.

## Bind And Run

Create the runtime for a validated document, bind it, and execute normally:

```go
package main

import (
	"context"
	"encoding/json"

	"github.com/OpenUdon/uws/mockruntime"
	"github.com/OpenUdon/uws/uws1"
)

func simulate(ctx context.Context, doc *uws1.Document, fixtures *mockruntime.FixtureSet) error {
	runtime, err := mockruntime.NewRuntime(doc, mockruntime.Options{
		Fixtures: fixtures,
		ResponseResolver: mockruntime.ResponseResolverFunc(
			func(_ context.Context, _ *uws1.Operation) (mockruntime.ResponseDefinition, error) {
				return mockruntime.ResponseDefinition{
					Example: json.RawMessage(`{"statusCode":200,"body":{"name":"Ada"}}`),
				}, nil
			},
		),
	})
	if err != nil {
		return err
	}
	doc.SetRuntime(runtime)
	return doc.Execute(ctx)
}
```

The runtime implements `uws1.Runtime` and the optional additive
`uws1.RuntimeWithResult`. A returned response is stored on the operation
execution record before success criteria run, so `$response` expressions can
inspect it. A direct operation step receives a copy on its step record so its
outputs can also inspect `$response`. The response shape is owned by the
operation's source or extension profile. HTTP-style responses commonly use
`statusCode`, `headers`, and `body`.

## Response Selection

For each resolved request, the runtime canonicalizes the complete request
binding object and computes its fixture digest. Selection is deterministic:

1. An exact `(operationId, requestDigest)` fixture match is replayed. Repeated
   calls, retries, and parallel calls with that key receive the same response.
2. If no exact match exists and generated fallback was explicitly enabled, the
   response resolver is consulted.
3. With no fixture set, the response resolver supplies a response directly.
4. A resolver example is returned as supplied. If it also supplies a schema,
   the example takes precedence. Otherwise a deterministic value is generated
   from the schema.

When a fixture set is supplied, a miss is an error by default, even when a
resolver exists. Set `AllowGeneratedFallback: true` only when the workflow
intentionally permits an example or synthesized response after a fixture miss.
The runtime does not infer responses from HTTP methods, source selectors, or
operation names. Resolver errors, invalid examples, missing response
definitions, and unsupported schemas fail the leaf operation with a diagnostic.

Schema synthesis supports strings, integers, numbers, booleans, nulls, arrays,
and objects. Arrays with an item schema contain one generated item; arrays
without one are empty. Objects include every declared property in sorted key
order. `$ref`, composition (`allOf`, `oneOf`, `anyOf`), `format`, recursive
schemas, and unknown schema types require a caller-supplied example and are
rejected when synthesis is requested. Depth, node count, and response byte
limits (64 schema levels, 10,000 nodes, and 16 MiB per response) bound
generation.

## Expressions And Items

The runtime evaluates the core UWS expression sources against the execution
context: `$response`, `$outputs`, `$steps`, `$variables`, `$trigger`, `$inputs`,
`$item`, `$index`, and loop-only `$batchIndex`, including supported dot-walks
and JSON Pointer response paths. It implements the six core comparisons with
type-preserving JSON comparisons. `ResolveItems` requires an ordered JSON
array. The core orchestrator continues to handle workflow structure, actions,
criteria kinds, and output propagation.

Implementation-specific functions or operators are not guessed. Unsupported
expressions and malformed paths return explicit errors; they do not silently
resolve to success. This runtime does not make ordinary document validation
reject implementation-specific expressions.

## Request Records And Data Handling

`Runtime.RequestRecords()` returns a defensive snapshot of bounded in-memory
records. Each record includes the operation ID, canonical resolved request,
request digest, and response kind (`fixture`, `example`, `synthesized`,
`live-read`, or `unavailable`). The history is retained only by that runtime instance and is
never written automatically. Request nesting is limited to 64 levels;
expression text is limited to 64 KiB, and expression/input nesting is limited
to 32 levels. Request canonicalization accepts at most 1 MiB; retained history
is limited to 10,000 records and 16 MiB of canonical request data per runtime.

Resolved request bindings may contain private values. Treat request records
as sensitive, review them before display or export, and redact them explicitly.
Execution records also retain returned response values, which may be sensitive;
apply the same care to those records. The request digest is an identity key,
not a privacy control. The fixture format's [recorded-response
constructor](mock-fixtures.md) separately requires a caller-supplied redactor;
the runtime does not export or capture response fixtures automatically.

## Hybrid Live Reads

`NewHybridRuntime` wraps a pure `Runtime` with a caller-owned read adapter. It
requires a UWS 1.12.0-or-later document, `AllowLiveReads: true`, and a non-nil
`ReadDelegate`. Construct it only after the document's operations and effect
labels are final:

```go
hybrid, err := mockruntime.NewHybridRuntime(mock, mockruntime.HybridOptions{
    AllowLiveReads: true,
    ReadDelegate: mockruntime.ReadDelegateFunc(
        func(ctx context.Context, operation *uws1.Operation, request map[string]any) (json.RawMessage, error) {
            return localAdapter.ExecuteRead(ctx, operation, request)
        },
    ),
})
if err != nil {
    return err
}
doc.SetRuntime(hybrid)
return doc.Execute(ctx)
```

When enabled, an operation declared on the bound document with exactly
`effect: read` is sent to the delegate. That route is selected before fixture
lookup or response generation, so even an exact read fixture is bypassed.
Operations marked `write`, marked `unknown`, or with no effect stay on the pure
mock path and use its normal fixture/example/synthesis policy. The delegate
receives the request after UWS expression resolution and must return one valid
JSON value of at most 16 MiB. The wrapper passes the execution context through,
propagates delegate errors and cancellation, and returns the response to the
orchestrator for criteria and output evaluation.

Each resolved and canonicalized leaf request adds a request record before
response selection or delegation. `responseKind: live-read`
identifies the selected delegate path, including when that attempt later fails
or observes cancellation. Other calls identify fixture replay, example use,
schema synthesis, or an unavailable response. The ordered records describe
each leaf operation invocation individually, so mock and delegated results
stay distinguishable without treating mock success as evidence that a real
call succeeded. They do not assign one ID to an entire workflow run; callers
that reuse a runtime can compare record snapshots before and after each
execution.

An effect label describes expected behavior and does not authorize access.
`AllowLiveReads` is a separate explicit opt-in; the caller still owns account
selection, credentials, policy, and the delegate's real transport. The UWS
package supplies no live client. Delegates should honor cancellation, avoid
mutating the operation/request, keep private credentials outside request
bindings and digests, and avoid putting sensitive values in returned errors.

## Scope

This package simulates leaf responses using fixtures and caller-provided
examples or schemas. It does not supply source parsers, credentials, provider
clients, or concrete provider transports. The separate
[`HybridRuntime`](#hybrid-live-reads) supplies a narrow handoff to a caller-
owned read adapter.

See [Mock Fixtures](mock-fixtures.md), [Execution Model](07-Execution-Model.md),
and [Mock Fixture Format 1.0](https://github.com/OpenUdon/uws/blob/main/versions/mock-fixtures.1.0.md).
