# Archive C01 - Core Workflow Contract And Execution

**Context.** The UWS 1.x workflow language, its Go document model, semantic
validation, and runtime-independent orchestration semantics.

**Baseline.** 8382d0f26b3b10870125760643078d1a1a3e31b6

**Coverage.** verified

**Supersedes.** none

## Scope And Responsibilities

UWS is a compact workflow overlay for operations already described by API,
RPC, event, or browser-profile source documents. Source documents remain
authoritative for methods, paths, channels, messages, schemas, servers,
protocol details, and security. UWS owns local operation binding, request
values, dependencies, workflow structure, triggers, outputs, structural
results, runtime expressions, and control flow.

The current core release is UWS 1.9.2. It supports nine current source types
and extension-owned operations. It retains compatibility gates for earlier
1.x documents, including the historical `ansible-module` source type only for
UWS 1.6 documents. Concrete transports, credentials, provider configuration,
leaf implementations, and persistence backends are outside this context.

## Domain And Workflows

A document contains metadata, optional source descriptions and variables,
operations, optional workflows and triggers, optional structural results and
components, and optional content-trust declarations. An operation is either
source-bound through exactly one compatible selector or extension-owned through
`x-uws-operation-profile`.

Workflows and structural steps use six closed structural types: `sequence`,
`parallel`, `switch`, `loop`, `merge`, and `await`. Steps may invoke an
operation, invoke another workflow, or own structural content; reference steps
cannot also declare child blocks. Dependencies, conditions, iteration,
timeouts, criteria, retry/goto/end actions, outputs, results, and trigger routes
form the executable graph. Runtime expressions provide the portable data-flow
vocabulary for responses, outputs, steps, variables, triggers, inputs, items,
and indexes.

Execution validates the document and executable graph before starting. A
document-level run resolves a single entry workflow; trigger dispatch instead
selects declared routes and executes their workflow or top-level-step targets
through the same orchestration rules.

## System Shape

`uws1` supplies the wire types, extension-preserving JSON behavior, HCL hooks,
semantic validator, execution indices, and orchestrator. `Document.Execute`,
workflow/step entry points, and `Document.DispatchTrigger` create an
`Orchestrator` after validation. The orchestrator owns dependency scheduling,
structural execution, action handling, output propagation, trigger routing,
and execution records.

The bound `Runtime` interface is deliberately narrow: it executes a leaf
operation, evaluates a runtime expression, and resolves iteration items.
Execution context exposes current inputs, trigger state, iteration state, and
record snapshots without adding those values to the portable wire document.
Runtime rebinding, execution, and record reads on the same document are not
internally synchronized.

## Contracts And Dependencies

The normative contract is coordinated across `versions/1.9.2.json`,
`versions/1.9.2.md`, and the `uws1` model and validator. JSON Schema owns shape
checks; `Document.Validate` and `ValidateResult` own semantic checks such as
identifier uniqueness, selector compatibility, reference integrity, action
rules, structural constraints, and version gating. `ValidateExecutable` adds
requirements imposed by the built-in orchestrator, including unambiguous
executable identifiers and a resolvable entry workflow.

Objects that accept extensions preserve arbitrary `x-*` fields while rejecting
unknown non-extension fields. The `x-uws-` prefix is reserved for UWS-owned
contracts. JSON and HCL tags on the reachable document tree support the
separate interchange context.

The orchestrator depends on XPath, JSONPath, HCL decoding, synchronization, and
context cancellation libraries. It does not depend on concrete runtime or
transport packages.

## Operations And Verification

Timeout and cancellation failures terminate the active construct. Dependency
cycles, unresolved references, malformed selectors, invalid structural field
combinations, and unsupported executable criteria fail validation or execution
before affected leaf work proceeds. Execution records are in-memory snapshots;
UWS does not define a serialized history store, credential store, idempotency
store, or retry-replay backend.

At this baseline, `go test ./...`, `go vet ./...`, `mkdocs build --strict`, and
`git diff --check` pass. CI also runs the full suite with the race detector.

## Evidence

| Claim | Repository Evidence |
|---|---|
| UWS is a source-document-first workflow overlay with deliberately narrow scope. | `README.md`; `versions/1.9.2.md` sections 2 and 2.1 |
| The current structural and wire contract is UWS 1.9.2. | `versions/1.9.2.json`; `versions/1.9.2.md`; `versions/CHANGELOG.md` |
| The Go root model carries the portable document plus runtime-only execution state. | `uws1/document.go`; `uws1/operation.go`; `uws1/workflow.go`; `uws1/trigger.go` |
| Core orchestration and runtime responsibilities are explicitly separated. | `uws1/execution.go`; `uws1/execution_structural.go`; `versions/1.9.2.md` section 7 |
| Validation is layered into schema, semantic, executable, and entrypoint checks. | `uws1/validation.go`; `uws1/executable_validation.go`; `uws1/schema_conformance_test.go` |
| Version-specific source and field behavior is enforced without rewriting older contracts. | `uws1/validation_version.go`; `schemas/version_immutability_test.go` |
| Core semantics have focused tests for execution, triggers, validation, schema parity, and serialization. | `uws1/*_test.go`; `testdata/sample.uws.json`; `testdata/invalid/` |

## Observed Gaps

This repository defines no concrete source or extension runtime, durable
execution-record store, idempotency backend, or credential system. Those are
intentional executor-owned boundaries, so end-to-end provider invocation and
persistence behavior cannot be verified from this repository alone.
