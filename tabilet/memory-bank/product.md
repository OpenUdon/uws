# Product

## C08 accepted reference foundation

C08 is accepted/retired after review 2 at `0411eea6fc84fbd6aa97cef94f53f301260f4844`, independently observed on origin/main. The [qualification](../../docs/c08-qualification.md) pins the source and supplement; C09/M08 remain pending.

## Approved Stage 11 product direction — not implemented

C09.3 adds deterministic advisory flow diagnostics for reachability, unreferenced outputs, unknown effects/pending ordering and work bounds. It executes no condition/effect and never grants authority.

C09.2 adds advisory binding validation with compatible/incompatible/indeterminate outcomes. It checks required/literal/expression/output/security metadata without credentials or effects.

C09.1 adds the public source-neutral binding shape/resolver contract. Metadata completeness, exact identity and security OR/AND semantics are explicit; metadata never grants execution.

C08.1 now provides the public expressions parser for the existing core grammar, including version/field/loop gates and value-free errors. C08.2 adds exact scalar comparison, null propagation and workflow/nearest-iteration lookup with shared vectors. C08.3 now makes the mock consume the shared evaluator and exposes opt-in CheckPortability diagnostics. Historical mock extensions remain explicit compatibility adapters; ordinary validation is unchanged.

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. UWS will own opt-in expression portability, source-neutral binding diagnostics and a verified HCL presentation module. Source parsing and product authority remain outside UWS. Existing published versions and ordinary validation compatibility stay intact; HCL input removal is outside Stage 11.
Both phases belong to one stage. C08 acceptance/publication is recorded above; C09/M08 remain pending and claim no installed Kinet behavior. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

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

UWS 1.12.0 is the current core contract. It adds optional operation `effect`
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
The UWS 1.12 specification, exact-version schema, Go model and validator,
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
