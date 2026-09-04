# Proposal: MCP Tool Calls as an OpenUdon Experiment

**Status: non-normative proposal. Not adopted or implemented.** This note does
not allocate a UWS version, add a source type, define a schema or Go API, or
reserve a new `x-uws-*` field. The `openudon.*` and `x-openudon-*` names below
describe a possible product-owned experiment; they do not commit OpenUdon to
implementing it.

The experiment asks a narrow question: can a reviewed snapshot of an MCP tool
definition support the same fail-closed workflow binding expected of other UWS
operation contracts? It deliberately keeps MCP out of UWS core while producing
evidence that could justify later standardization.

## The actual MCP gap

The stable MCP 2025-11-25 protocol already exposes a machine-readable operation
catalog. A client can exhaust the paginated `tools/list` method and capture each
complete Tool object, including its name, description, `inputSchema`, optional
`outputSchema`, annotations, and other metadata. Once captured in a deterministic
document, that catalog is reviewable, hashable, and diffable. MCP is therefore
not inherently non-diffable.

The missing piece is a standardized, fetchable, freshness-bound tool manifest.
The [stable tools protocol](https://modelcontextprotocol.io/specification/2025-11-25/server/tools)
discovers tools from a live initialized server. It does not define a portable
artifact that identifies the executable server capability, pins the reviewed
Tool object, sets an expiry, and specifies the comparison that must pass before
`tools/call`. The optional `notifications/tools/list_changed` signal helps an
active client invalidate a cache, but it is not a durable freshness contract and
its absence does not prove that a catalog is unchanged.

The MCP Registry does not fill that gap. Registry `server.json` describes server
identity, packages or remote endpoints, installation, and configuration; it does
not publish the server's Tool objects. The
[Registry publishing example](https://modelcontextprotocol.io/registry/quickstart)
therefore complements tool discovery rather than replacing it.

This distinction changes the adoption question. Stage 1 does not need to invent
an operation catalog. It needs to test whether a product-owned reviewed snapshot
can bind a captured catalog to a live capability and detect staleness or drift.

## Why this is neither a source type nor the runtime floor

MCP tools name invocable operations that the MCP server implements. This differs
from client-shipped code such as an Ansible module: the MCP client does not copy
the tool implementation to a substrate that lacked it. A local stdio server uses
a process boundary rather than a network boundary, but the server still exposes
and implements the named operation.

The live catalog can also be materialized into a reviewable artifact without
discarding its schemas. Flattening a tool into `uws.runtime.1.0` as a generic
function or command selector would lose that contract. The runtime supplement is
the floor for work whose only reviewable artifact is a selector or command; an
MCP Tool object is richer.

MCP is nevertheless not ready for `sourceDescriptions[].type`. A source type
needs portable binding and freshness semantics, not merely a locally invented
snapshot. The provisional snapshot below is the experiment that would test
those semantics before any public UWS contract is proposed.

## Provisional Stage 1 interface

Stage 1 uses the product-owned operation profile
`openudon.mcp-tool-call.0.1`. An extension-owned UWS operation has no source
binding and uses the existing `x-uws-operation-profile` hook with exactly one
new product field, `x-openudon-mcp`:

```yaml
- operationId: render_angles
  x-uws-operation-profile: openudon.mcp-tool-call.0.1
  x-openudon-mcp:
    manifest: ./mcp/blender-farm.tools.json
    manifestSha256: 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
    server: blender_farm
    tool: render_turntable
  request:
    body:
      source: $inputs.source_uri
      angles: 6
      resolution: 2048
  outputs:
    images: $response.body#/renderedUris
```

The four binding fields are:

| Field | Provisional meaning |
| --- | --- |
| `manifest` | URI reference to a separately reviewed `openudon.mcp-tool-manifest.0.1` JSON snapshot. |
| `manifestSha256` | Lowercase 64-character SHA-256 digest of the snapshot's deterministic canonical JSON. A mismatch fails before server connection. |
| `server` | Product-owned symbolic server key in the snapshot and runtime configuration. Connection details and credentials remain runtime-private. |
| `tool` | Exact, case-sensitive MCP Tool name selected within that server's captured catalog. |

The profile does not define server discovery, transport, subprocess spawning,
authorization, credential storage, retries, or connection lifecycle. Those
remain runtime-owned. The observed `serverInfo` is review evidence and a useful
drift signal; it is not assumed to be a globally unique server identity.

## Reviewed tool-manifest snapshot

Stage 1 pairs the call profile with the separate product-owned snapshot profile
`openudon.mcp-tool-manifest.0.1`. This is a reviewed artifact, not a new UWS
source document and not an MCP standard.

An illustrative shape is:

```yaml
profile: openudon.mcp-tool-manifest.0.1
protocolVersion: "2025-11-25"
capturedAt: "2026-08-30T12:00:00Z"
expiresAt: "2026-09-06T12:00:00Z"
servers:
  blender_farm:
    serverInfo:
      name: example-blender-server
      version: 0.8.2
    toolPages:
      - tools:
          - name: render_turntable
            description: Render reviewed views from an addressable Blender asset.
            inputSchema:
              type: object
              properties:
                source: {type: string, format: uri}
                angles: {type: integer, minimum: 1, maximum: 36}
                resolution: {type: integer, enum: [1024, 2048, 4096]}
              required: [source, angles, resolution]
            outputSchema:
              type: object
              properties:
                renderedUris:
                  type: array
                  items: {type: string, format: uri}
              required: [renderedUris]
    reviews:
      render_turntable:
        reviewer: release-review@example.test
        reviewedAt: "2026-08-30T13:00:00Z"
        assertions:
          headless: true
          independentCalls: true
          modelBacked: false
          sideEffects: [creates_render_artifacts]
          concurrency: safe
          referenceBasedIO: true
        inputChannels:
          /source: data
          /angles: data
          /resolution: data
```

This is illustrative rather than a schema commitment. A real Stage 1
implementation should settle the exact spelling only in its owning repository.
The behavioral requirements are more important than the provisional
serialization.

### Capture and hashing

A conforming experiment would apply these rules:

1. Initialize one MCP connection, record the negotiated `protocolVersion` and
   complete observed `serverInfo`, and exhaust every `tools/list` page by
   following `nextCursor` until it is absent. A first-page-only capture is
   invalid.
2. Retain every captured page result without JSON-RPC request IDs or framing.
   At minimum, every reviewed Tool object must be retained complete and
   unprojected. A selected tool name must occur exactly once across the pages.
3. Record UTC capture and expiry instants. Expiry is a reviewer decision with a
   product policy maximum; neither an MCP version nor an annotation silently
   extends it.
4. Store the reviewer identity, review time, required assertions, and reviewed
   input-channel classifications alongside each selected tool.
5. Serialize the complete snapshot as JSON, reject duplicate object member
   names, canonicalize it with
   [RFC 8785 JSON Canonicalization Scheme](https://www.rfc-editor.org/rfc/rfc8785),
   and compute SHA-256 over those canonical UTF-8 bytes. `manifestSha256` pins
   that digest.

Canonicalization makes the captured artifact deterministic; it does not make
the underlying operation deterministic. Array order and any complete Tool
metadata remain part of the reviewed value, so unexplained changes fail closed
rather than being normalized away.

### Freshness and drift

Before every `tools/call`, the runtime must:

1. resolve and hash the manifest, then reject a digest mismatch;
2. reject a snapshot whose `expiresAt` has passed or whose capture/review data
   is incomplete;
3. initialize the configured server using the recorded protocol version,
   compare the complete observed `serverInfo` with the snapshot, and reject a
   negotiated-version or `serverInfo` mismatch;
4. exhaust the live `tools/list` pagination and locate the selected name exactly
   once;
5. canonicalize and compare the complete live Tool object with the complete
   reviewed Tool object; and
6. invoke the tool only if the comparison is exact and all reviewer assertions
   are acceptable to local policy.

A missing tool, duplicate name, changed schema, changed description, changed
annotation, changed execution metadata, or other Tool-object drift fails the
operation before invocation. This exact comparison is intentionally
conservative for an experiment; evidence may later show which metadata can be
separated from the executable contract.

If the server advertises `tools.listChanged` and sends
`notifications/tools/list_changed`, the runtime must invalidate the reviewed
snapshot immediately. Re-listing to the same value does not revive that review:
a new capture and review are required. A notification received after the
pre-call comparison but before dispatch prevents dispatch; one received during
a call invalidates the binding and prevents its result from being exposed as
reviewed output. Because an in-flight side effect may already have occurred,
that outcome is indeterminate and must not be retried automatically.

These checks close the Stage 1 time-of-review gap as far as the protocol permits,
but they do not claim an atomic transaction between `tools/list` and
`tools/call`. That residual race is evidence relevant to any later proposal.

## Reviewer assertions define the qualifying subset

MCP annotations are retained as part of the Tool object, but the stable
specification describes them as hints rather than faithful guarantees. Stage 1
therefore never converts `readOnlyHint`, `destructiveHint`, `idempotentHint`, or
`openWorldHint` directly into trusted policy. A reviewer must make independent
assertions for each selected tool:

| Assertion | Stage 1 requirement |
| --- | --- |
| `headless` | Must be `true`: the call needs no GUI or contemporaneous human interaction. |
| `independentCalls` | Must be `true`: semantics do not depend on a prior call, ambient conversation, current document, session, or server-minted handle. |
| `modelBacked` | Must state `true` or `false`; omission or `unknown` fails review. |
| `sideEffects` | Must enumerate reviewed side-effect classes, including an explicit empty list when none are expected. |
| `concurrency` | Must be `safe` or `serialize`; `serialize` requires the runtime to serialize calls in the asserted scope. Unknown behavior fails review. |
| `referenceBasedIO` | Must be `true`: inputs and structured outputs use JSON values or addressable references, not local paths, live handles, or process-local objects. |

These assertions describe suitability, not truth guaranteed by MCP. Runtime
policy may reject a reviewed tool even when all assertions are present.

Session- or handle-dependent tools are outside Stage 1 even if a server happens
to make them work in one connection. The workflow must own ordering and data
flow; an undeclared ambient server state must not become a second orchestration
graph.

## Invocation and output semantics

For a call that passes the freshness gate:

- The evaluated UWS `request.body` object becomes MCP `tools/call`
  `arguments`. An absent body becomes an empty object. The runtime validates it
  against the reviewed `inputSchema` before dispatch.
- The `tool` binding supplies the MCP `name`. Other UWS request locations are
  not mapped by Stage 1.
- A JSON-RPC or other protocol-level error fails the UWS operation. A successful
  MCP response with `isError: true` also fails it; error content is diagnostic
  only and never becomes `$response.body`.
- A UWS operation that declares outputs requires the reviewed Tool object to
  contain `outputSchema`, requires the call result to contain
  `structuredContent`, and validates that object against `outputSchema`.
  Missing or non-conforming structured content fails the operation.
- Validated `structuredContent` becomes `$response.body`, so ordinary UWS output
  expressions select from it.
- MCP `content` blocks remain runtime diagnostics. Text, image, audio,
  embedded-resource, and resource-link blocks are never projected into
  structured UWS outputs. `_meta` is likewise not exposed as `$response.body`.

This intentionally rejects the tempting fallback of parsing JSON-looking text.
Only the server-declared output contract and conforming structured result cross
the reviewed boundary.

## Excluded MCP capabilities

Stage 1 is one synchronous, independent `tools/call`. It excludes:

- task-augmented calls, including a Tool whose `execution.taskSupport` requires
  tasks;
- server requests for sampling or elicitation during discovery or invocation;
- prompts and resources as invocation or output surfaces;
- resource reads, subscriptions, embedded resources, and resource links as a
  substitute for structured output; and
- any tool whose meaning depends on an MCP session, previous call, current
  server state, or minted handle.

The [2025-11-25 task facility](https://modelcontextprotocol.io/specification/2025-11-25/basic/utilities/tasks)
is experimental and is not evidence for Stage 1 support. Later MCP draft work
may change sessions, list stability, or tool execution, but draft behavior is
not treated here as stable protocol fact.

## UWS 1.9.1 content trust

The Stage 1 resolver must integrate with the existing advisory
[UWS 1.9.1 content-trust model](content-trust.md); it must not invent a parallel
trust system.

Each selected tool review classifies every argument channel that can receive a
UWS-bound value as `data`, `instruction`, or `authority`. Keys are RFC 6901
pointers relative to MCP `arguments`; a dynamic or otherwise unclassified bound
path makes the resolver fail closed. `data` is content to process,
`instruction` changes interpreter or model behavior, and `authority` selects or
authorizes a side effect, destination, account, command, or equivalent
capability.

The resolver supplies these conservative contracts to `contenttrust.Analyze`:

- MCP-derived operation outputs default to `untrusted` provenance.
- Derived outputs inherit input provenance. Explicit UWS `contentTrust`
  declarations retain their normal precedence, but an untrusted input is not
  laundered by derivation.
- Output value capability is derived from the selected location in
  `outputSchema`: unconstrained strings are `free_text`; Boolean, numeric,
  `null`, and closed-enum scalars are `constrained_scalar`; objects and arrays
  are `composite`; missing or ambiguous schema information is `unknown`.
- The resolver reports the reviewed input-channel map and per-output value
  contracts; it does not inspect runtime arguments or results.

Analysis remains deterministic, advisory, and value-free. Findings contain
stable paths, labels, codes, and fixed messages—not MCP result values, excerpts,
credentials, or server error details. Stage 1 runtime policy may choose to block
on findings, but ordinary UWS validation and execution semantics do not change.

## Model-backed tools and determinism

UWS does not prohibit model-backed operations. `uws.runtime.1.0` already permits
a disclosed `llm` execution surface. The MCP concern is hidden nondeterminism:
without review metadata, a model-backed tool can look identical to a stable
deterministic operation even though it makes a fresh decision on every call.

Stage 1 therefore requires the `modelBacked` assertion rather than requiring it
to be false. A declared model-backed tool remains eligible when local policy
accepts its reviewability, side effects, and output contract. The experiment
should record whether retries, concurrency, and repeated identical arguments
produce materially different results; that evidence can inform later policy
without pretending nondeterminism is forbidden by UWS.

## Staged adoption

**Stage 1 — OpenUdon-owned experiment.** Test
`openudon.mcp-tool-call.0.1`, `x-openudon-mcp`, and the separate
`openudon.mcp-tool-manifest.0.1` snapshot in a consuming product repository. No
UWS schema, version, reserved extension, Go API, or public interoperability
commitment follows from this note.

**Stage 2 — public operation supplement.** Consider a UWS-family supplement
only after Stage 1 produces real workflow evidence and a second independent
consumer needs interoperable exchange. That proposal would need to settle the
manifest schema, server identity, freshness window, canonical comparison,
notification race, result mapping, and resolver behavior from implementation
experience.

**Stage 3 — possible source type.** Consider MCP as a source family only after
portable manifest and freshness semantics exist and multiple runtimes
demonstrate demand. An upstream MCP tool manifest is preferred because the
implementer can publish it. A proven UWS-defined reviewed snapshot could also
be evaluated in theory, but only after it has portable binding and drift
semantics rather than one product's conventions.

At every stage, source documents continue to own operations when a stronger
provider contract already exists. An MCP wrapper around an OpenAPI, GraphQL,
SQL, object-storage, mail, or other already modeled interface is not a reason to
replace that contract with MCP.
