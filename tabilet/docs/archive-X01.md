# Archive X01 - Source Admission And Extension Profiles

**Context.** The boundary between first-class source contracts,
extension-owned operations, runtime metadata, conversion inputs, data files,
and product-owned experiments.

**Baseline.** 8382d0f26b3b10870125760643078d1a1a3e31b6

**Coverage.** verified

**Supersedes.** none

## Scope And Responsibilities

UWS core recognizes source families only when a portable operation contract can
name operations, be reviewed before execution, bind to executable capability
with drift detection, and describe an operation implemented by the remote
party. The current source types are OpenAPI, Google Discovery, AWS Smithy,
AsyncAPI, GraphQL, OpenRPC, gRPC/protobuf, OData, and browser profiles.

Non-source behavior remains extension-owned. The public Runtime Supplement 1.0
offers a small metadata vocabulary for common non-HTTP leaf surfaces, while
product profiles may define their own `x-*` payloads. This context does not own
concrete execution, credentials, transport clients, source parsers, or provider
configuration.

## Domain And Workflows

Source admission separates four questions: whether a document names invocable
operations, whether it can materialize a reviewable operation contract, whether
that contract can be pinned and revalidated against live capability, and
whether the remote party actually implements it. The failing criterion routes
a candidate toward a runtime selector, converter, advisory/capture tooling, or
product-owned operation-profile experiment rather than automatically creating
a public source type.

The runtime supplement names `ssh`, `cmd`, `fnct`, `fileio`, `sql`, `s3`,
`smtp`, `dns`, `ldaps`, `scp`, `sftp`, and `llm`. Its payload carries only a
type and optional command, working directory, function, workflow, or arguments.
HTTP, API, RPC, event, and browser-profile work remains source-bound.

The MCP proposal is explicitly non-normative and unimplemented. Its Stage 1
design uses the product-owned `openudon.mcp-tool-call.0.1` profile,
`x-openudon-mcp`, and a separately reviewed
`openudon.mcp-tool-manifest.0.1` snapshot. It proposes complete paginated
capture, canonical hashing, expiry, per-call relisting and exact Tool-object
comparison, notification invalidation, structured-content-only output mapping,
reviewer assertions, and content-trust resolver integration. Public
standardization is deferred until implementation evidence and a second
independent consumer exist.

## System Shape

Core operations carry source selectors or the generic
`x-uws-operation-profile` marker. The semantic validator enforces exclusivity
and source-family/version compatibility but does not validate arbitrary
product-owned extension payloads. Unknown `x-*` fields survive round trips and
unknown profiles fail only when a runtime tries to execute them.

`runtimes` supplies constants, the typed `x-uws-runtime` payload, and read/set
helpers. The separate runtime schema validates that small payload. Admission
criteria, historical decisions, and candidate experiments live in
non-normative documentation until a versioned schema and specification adopt
them.

## Contracts And Dependencies

`x-uws-` is reserved for UWS-owned fields and supplements; third parties use
their own prefixes. An implementation profile is authoritative for its own
payload, while UWS core preserves but does not interpret it. Profile-specific
static flow semantics require a `contenttrust.Resolver` from the consumer.

Ansible demonstrates the remote-capability boundary: its module source type was
published in UWS 1.6 and withdrawn in 1.7 because the control node supplies the
implementation. Historical 1.6 artifacts remain published, but UWS 1.9 has no
Ansible replacement. Files and spreadsheets are data rather than operation
catalogs; opaque scripts and macros have no operation contract and remain
runtime work.

## Operations And Verification

The runtime supplement schema and Go constants are checked for enum parity,
required type selection, slim payload shape, and rejected removed/configuration
fields. Core tests cover source-family version gates, selector compatibility,
unknown profile preservation, and extension round trips.

The MCP design's identity, freshness, drift, notification, output, and reviewer
rules are proposal text only. No schema, Go API, runtime, snapshot, or execution
evidence for that experiment exists in this repository at the baseline.

## Evidence

| Claim | Repository Evidence |
|---|---|
| Current first-class source types and their version gates are explicit. | `versions/1.9.2.md`; `uws1/source_description.go`; `uws1/validation_version.go` |
| Admission distinguishes reviewable materialization from executable binding and drift detection. | `docs/future-source-profiles.md` section "Source Type Admission Criteria" |
| Extension-owned operations use a generic profile hook and preserved profile payloads. | `uws1/operation.go`; `uws1/validation_operation.go`; `docs/08-Extension-Profiles.md` |
| Runtime Supplement 1.0 is metadata rather than an execution API. | `versions/runtime.1.0.json`; `versions/runtime.1.0.md`; `runtimes/runtime.go` |
| The MCP interface is an OpenUdon-owned Stage 1 proposal with no UWS wire change. | `docs/mcp-tool-calls.md` |
| Ansible is retained only as historical UWS 1.6 material. | `versions/CHANGELOG.md`; `docs/uws_1_6_ansible.md`; `versions/ansible.1.0.*` |

## Observed Gaps

`versions/runtime.1.0.md` still discusses `ansible-module` as a source-bound
case in present-tense runtime guidance, while UWS 1.7 and later explicitly
withdraw that source type. The MCP experiment has no implementation or
multi-consumer evidence, as its proposal status states. No concrete runtime or
profile-specific content-trust resolver is included in this repository.
