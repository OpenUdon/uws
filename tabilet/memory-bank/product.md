# Product

## Stage 11 public foundations

C08 and C09 are accepted/published: shared expression parsing/reference
evaluation and opt-in portability; source-neutral binding shapes, advisory
binding validation and deterministic flow observations. Their exact sources
remain `0411eea6fc84fbd6aa97cef94f53f301260f4844` (review 2) and
`6a267306032edc687a298cefc8bba7019d3ad059` (review 3). See the
[C08 qualification](../../docs/c08-qualification.md) and
[C09 qualification](../../docs/c09-qualification.md).

UWS 1.13 adds verified HCL presentation/deprecation while retaining all 1.12
wire/execution rules. The separate `github.com/OpenUdon/uws/hcl` module renders
from exact JSON/YAML bytes, independently verifies complete values and numeric
lexemes, and exposes deprecated inert Import. Failed proof suppresses the view.
A view is derived metadata; it is never source/package input or authority.
Provider source parsing, workflow policy, credential handling and runtime authorization
remain consumer-owned. Existing root HCL APIs and browser inputs remain retained.
M08 is accepted at `c0b19385a3b034cd45de16726668b9150f0633f2` after review 3;
root/codec source is independently published/resolved.
[qualification](../../docs/m08-qualification.md) records exact artifacts and limits.

The [coordinated Stage 11 contract](../../../kinet/docs/stage11.md) covers both
phases. Local acceptance supplies no installed Kinet behavior or live authority;
the M44 service and frozen consumer/browser pins remain unchanged. Stage 12 owns
the separately qualified browser/input-removal boundary.

This file summarizes current observed product truth. The linked archives are
frozen evidence for commit `8382d0f26b3b10870125760643078d1a1a3e31b6`, not a
substitute for this evolving summary.

## Purpose

The Udon Workflow Specification is a compact, language-neutral workflow overlay
for API-, RPC-, event-, and reviewed browser-profile operations. It lets authors
compose operations that another source document already defines without
copying transport, schema, server, or security metadata into the workflow.

The repository publishes the normative UWS documents, Go data and execution
model, semantic and schema validators, interchange helpers, advisory integrity
analysis, separately versioned profiles, and mock-fixture data contracts. It is
a library and specification distribution, not a hosted workflow service or a
complete provider runtime.

## Users And Primary Workflows

- Workflow authors bind reviewed source operations, provide request values,
  compose structural workflows, route triggers, and name outputs and results.
- Tooling authors parse, convert, validate, inspect, package, and review UWS
  documents while preserving extension fields and version semantics.
- Runtime implementers bind concrete leaf execution, expression evaluation,
  item resolution, credentials, provider clients, and optional profile support
  to the shared core orchestrator.
- Profile reviewers publish or assess browser and other operation contracts,
  including provenance, freshness, side effects, confirmation, and safe input
  or output boundaries.
- Security consumers may run deterministic content-trust analysis and apply
  their own policy to its advisory findings before execution.

## Domain Model

| Concept | Meaning And Relationships |
|---|---|
| UWS document | One versioned workflow package with required metadata and at least one operation. It may reference source descriptions and contain workflows, triggers, results, components, and content-trust declarations. Runtime state is never serialized into it. |
| Source description | A uniquely named reference to one provider- or author-owned operation contract. Many local operations may bind to it, but it remains authoritative for protocol and payload semantics. |
| Operation | One uniquely identified executable unit. It is either bound to exactly one source operation through one compatible selector or owned by exactly one named extension profile. |
| Workflow and step | A workflow owns a structural graph of steps. A step invokes one operation, invokes one workflow, or owns one of six structural forms; reference and structural roles cannot be combined. |
| Trigger and route | A trigger is an external entry declaration with named outputs. Its routes target declared workflows or top-level steps and enter the same orchestrator used by ordinary execution. |
| Runtime expression | A portable reference from request, control, result, or output fields into the current execution context. It carries values between graph nodes without redefining source schemas. |
| Bound runtime | Transient execution authority for leaf work, expression evaluation, and iteration items. One runtime is bound to a document at execution time and remains outside the wire contract. |
| Content trust | Optional reviewed provenance declarations plus resolver-supplied channel contracts. Analysis produces advisory edges and findings without changing validation or execution. |
| Profile artifact | A separately versioned contract for behavior kept outside core, such as browser capability, authentication, registration, or runtime metadata. Exact discriminators select compatibility; live secrets and handles stay private. |

Source documents own operations and protocol contracts. Workflows arrange local
operation identities and nested workflow calls through dependencies and the six
structural constructs. Triggers provide routed entry, runtime expressions carry
values through the graph, and a bound runtime performs leaf work while the core
orchestrator owns portable graph semantics.

Browser capability profiles are first-class but weaker, author-asserted source
contracts. Browser authentication and registration remain separate
extension-owned lifecycles because they bind private credentials, sessions,
inputs, approvals, and mutation authority at execution time. Content-trust
declarations and resolver contracts add a static integrity view without
changing execution.

## Product Invariants

- Source contracts remain authoritative; UWS does not redeclare their methods,
  schemas, protocols, servers, or security.
- Structural orchestration is shared core behavior; transports and leaf
  implementations remain runtime-owned.
- Every operation is either source-bound through one valid selector or
  extension-owned through `x-uws-operation-profile`.
- JSON Schema covers structure and Go validation covers semantic relationships;
  executable checks add orchestrator-specific requirements.
- Unknown `x-*` fields are preserved, while `x-uws-*` remains reserved for
  UWS-owned contracts.
- Published versions remain available and immutable. Compatibility is selected
  by explicit document and profile versions rather than silent downgrade.
- Content-trust findings are deterministic, advisory, and value-free; policy
  enforcement belongs to consumers.
- Browser credentials, private registration values, cookies, session handles,
  captures, and verification responses never enter portable profile artifacts.
  Recorded mock fixtures require an explicit redactor and never persist
  responses automatically.
- Ambiguous browser targeting, unsafe origins, uncertain mutations, and stale
  or unsupported profile bindings fail closed.

## Current Contract Surface

UWS 1.13.0 is the current core specification. It retains 1.12 wire/execution
semantics and adds the separate verified HCL presentation/deprecation contract.
The separate codec and byte-pinned conformance supplement are implemented;
M08 qualifies and publishes the final root/codec release. Earlier published versions and
ordinary HCL APIs remain available. UWS 1.12 added optional operation `effect`
classification (`read`, `write`, or `unknown`; omission is `unknown`) and
non-executable pending-step contracts with recursive object-root input and
output schemas. Pending-only documents may use an empty `operations` array
without placeholder operations. Structural and semantic validation accepts a
valid pending contract; executable validation rejects any pending step before
runtime methods are invoked, including in an unselected branch. Effect labels
are descriptive and do not grant execution authorization. UWS 1.11 introduced
version-gated response-body dot-walks, loop-only `$batchIndex`, numeric
`wait`/`batchSize` literals, terminal root-scoped `goto`, and corrected
`forEach` merge records. UWS 1.10 execution semantics remain active for 1.10
and later; earlier declarations retain their versioned behavior. Browser 1.10
is the current opt-in browser capability profile and adds typed CSS selector
match counts; empty profile lookup retains Browser 1.8 as its compatibility
default. Browser authentication/call 1.1 is current for sign-in; browser
registration/call 1.2 is current for reviewed registration verification; and
registration input 1.0 is the private envelope format. Runtime Supplement 1.0
remains the public metadata floor for common non-HTTP extension operations.
Mock Fixture Format 1.0 is a separate, inert response-fixture format keyed by a
UWS-local operation ID and SHA-256 over RFC 8785 canonical bytes for the
resolved request-binding object. Repeated exact-key lookups reuse the same
response. Its codec requires explicit redaction for recorded fixtures and has
no automatic response persistence. The public `mockruntime` package uses the
core orchestrator, exact fixture replay, caller-supplied examples or bounded
schema synthesis, expression evaluation, and bounded in-memory would-be
request records; this pure runtime makes no transport calls or automatic file
writes. Its separately enabled `HybridRuntime` delegates only explicitly
declared UWS 1.12+ `read` operations to a caller-provided adapter; writes and
unknown effects remain mocked.
The UWS 1.13 specification, exact-version schema, Go model and validator,
executor, conversion helpers, embedded archive, and protected digests are
coordinated release artifacts. The pinned 1.11 executable conformance corpus
remains unchanged; earlier published contracts remain accepted according to
their version gates. Ansible support is historical UWS 1.6 material only.

The MCP tool-call design is an OpenUdon experiment proposal, not an adopted UWS
contract or implementation.

## Non-Goals

UWS does not provide source parsers, concrete provider clients, credential or
secret storage, durable execution history, idempotency persistence, registry
transport, browser drivers, session persistence, discovery crawlers, or a
general-purpose automation language. Files such as spreadsheets and PDFs are
workflow data, not operation sources; opaque commands, scripts, and macros use
runtime-owned execution unless a stronger operation contract exists.

## Baseline Evidence

- [Core workflow contract and execution](../docs/archive-C01.md)
- [Advisory content-trust analysis](../docs/archive-T01.md)
- [Browser capability and account lifecycle profiles](../docs/archive-B01.md)
- [Source admission and extension profiles](../docs/archive-X01.md)
- [Interchange, validation, and distribution tooling](../docs/archive-M01.md)
