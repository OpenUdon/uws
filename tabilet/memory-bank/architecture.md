# Architecture

## Accepted C10 browser-profile shapes

C10 is accepted after whole review1 at implementation
`f01a2542410c583d0ea909dadd8f17527cc0d27d`, ordinary root module
`v0.0.0-20261009220914-f01a2542410c`. [Publication proof](../../docs/c10-publication.md)
records independent source/archive checks and all seven frozen consumer builds.
Native source production, actual consumer adoption, browser/runtime containment,
credential policy and execution authority remain with their owners. Existing
published profile/schema bytes, root/HCL declarations and actual pins are unchanged.


## M09 parser remediation

Accepted source `b099f6803277ae94c7e9f1da0904a0140b278f20` is independently
published as root/codec `v0.0.0-20261008043726-b099f6803277`, whole review 7.
[Ordinary module proof](../../docs/m09-publication.md) distinguishes exact
reviewed file contents and Go checksums from ZIP-container compression.

Strict JSON prevalidation now uses a bounded iterative container stack (depth
100) before shape-table/model decoding. The separate codec preflights HCL with
the retained dependency's iterative lexer before recursive parsing, counting
syntactic delimiters and refusing active template constructs outside its inert
subset. Comments, escaped template markers and literal text remain inert.
Worker process resource/cancellation custody remains consumer-owned.

M09 binding proofs retain compiled child nodes from the original schema resource
and dialect. Pattern-based closure and nullable/pattern output paths remain
indeterminate where no proof exists; false leaves are incompatible. Draft-07
tuples use their effective items/additionalItems semantics, preserving offline
local references and ignoring future array keywords under that dialect.
Path proofs preserve containing cardinality and required predecessor constraints;
unproved constraints remain indeterminate and native literal witnesses prove
possible values. Mixed templates preserve empty arrays as exact const schemas.

M09 flow uses explicit identities while preserving generic dependency group
barriers and step→workflow→operation precedence. Step outputs use explicit
workflow/result ownership and known operation invocation owners, without a
foreign-workflow fallback. Terminal goto selects the exact globally indexed
target in the root invocation, independently of dependency caller frames.
Root workflow dependencies share unqualified records; nested workflow
dependencies and explicit workflow calls create child invocation frames.
Step and operation outputs stay separate even when their names or IDs match.
Workflow controls retain incoming record snapshots, separately from refreshed
body/output frames.
Strict portability records loop/iteration availability
separately, including transitive trigger/dependency contexts and pre-forEach
controls. These opt-in diagnostics do not change ordinary parsing or validation.

## Stage 11 public foundations

Accepted C08 expressions use standard-library parsing and UWS-owned in-memory
execution snapshots. UseNumber projections preserve exact values before custom
legacy model decoding; mock delegates to the reference evaluator and retains
small explicit compatibility adapters. CheckPortability is opt-in and scans
only core-owned fields; profile/function/browser templates remain opaque.
Ordinary validation is unchanged.

Accepted C09 binding types, strict bounded shape-table serialization and
resolver snapshots are source-neutral. Advisory validation checks exact
source/selector identities and honest known/unknown schemas/security with
OR-of-AND semantics. Unsupported containment remains indeterminate; external
refs refuse. Flow produces deterministic value-free graph/effect observations
without conditions, mutation or effect dispatch. Source parsers/providers stay
outside UWS. Exact accepted/public sources and fixtures are in
[C08 qualification](../../docs/c08-qualification.md) and
[C09 qualification](../../docs/c09-qualification.md).

UWS 1.13 schema/spec/model/archive retain 1.12 wire/execution rules. The separate
`hcl/` module reads public model tags to render typed blocks/labels/extensions
directly from lossless byte-decoded values, avoiding legacy custom decoding.
An independent inert AST walk reconstructs exact numeric token ranges and
complete values; source/codec/view provenance is checked before exposure.
Functions, traversals, interpolation, ambiguous fields/keys, unsupported YAML
forms/tags and malformed mappings refuse. The module supplies no admission,
approval or execution; consumers bind its caller-supplied revision to actual
worker/module closure and own process limits/privacy. Existing root HCL APIs
and compatibility dependencies remain; no HCL-free core claim is made.
M08 is accepted at `c0b19385a3b034cd45de16726668b9150f0633f2` after review 3;
root/codec source is independently published/resolved. see
[qualification](../../docs/m08-qualification.md).

Both [Stage 11 phases](../../../kinet/docs/stage11.md) share exact upstream
reconciliation and one serial owner. Installed M44 and frozen consumer/browser
paths are unchanged; Stage 12 owns later browser/input-removal qualification.

This file records current observed system structure. Archive documents are
frozen repository snapshots; this summary should evolve when the implementation
or public contracts change.

## System Overview

UWS separates portable orchestration from operation contracts and concrete
execution:

```text
source/profile documents
          |
          v
UWS document -> schema + semantic validation -> core orchestrator
                                                   |
                                                   v
                                             bound runtime
                                                   |
                                                   v
                                          provider/leaf operation

UWS document + profile resolvers -> advisory content-trust report
```

The versioned JSON Schema and specification define the wire contract. `uws1`
implements the Go model, semantic validation, and orchestrator. Concrete
runtimes are injected at execution time and do not enter the serialized
document. Separately versioned profile schemas add browser and extension
contracts without expanding core into a transport or automation engine.

## Repository Ownership

| Area | Responsibility |
|---|---|
| `versions/` | Normative and historical version documents and separately versioned profile contracts. |
| `uws1/` | Core Go wire model, extension handling, semantic/executable validation, execution context, and orchestration. |
| `contenttrust/` | Deterministic advisory provenance and capability analysis plus the resolver API. |
| `browserauthentication/` | Inert authentication profile and operation-extension wire types. |
| `browserregistration/` | Inert registration, private-input declaration, and verification-policy wire types. |
| `runtimes/` | Runtime Supplement 1.0 constants, typed payload, and extension helpers. |
| `mockruntime/` | Versioned fixture codec, RFC 8785 request digest, exact fixture lookup, recorded-response redaction helper, and public orchestrator-backed mock runtime with bounded in-memory request records. |
| `convert/` | JSON, YAML, and HCL interchange with extension and dynamic-key preservation. |
| `expressions/` | Core expression parser, exact-value reference evaluator and opt-in portability checks. |
| `binding/` | Source-neutral shape/resolver contracts and advisory binding/flow diagnostics. |
| `hcl/` (separate module) | Byte-based verified HCL presentation and deprecated inert import. |
| `validation/` | File loading plus coordinated schema and semantic validation. |
| `schemas/` | Schema lookup, embedded version archive, profile validators, and browser cross-document checks. |
| `internal/generateversionarchive/` | Deterministic generation of the embedded JSON-document ZIP. |
| `testdata/` | Canonical, adversarial, conversion, and profile fixtures. |
| `docs/`, `mkdocs.yml` | User guides, roadmap/design records, and documentation-site navigation. |
| `.github/workflows/` | Go verification and documentation deployment automation. |

## Execution And Analysis Flow

`Document.Execute` and related entry points require a bound `Runtime`. They run
ordinary semantic validation, executable validation, and entrypoint checks
before constructing an `Orchestrator`. The orchestrator indexes operations,
workflows, steps, and parallel groups; resolves dependencies; executes the six
structural constructs; applies actions; and accumulates records. It delegates
leaf execution, expression evaluation, and item resolution. A runtime may
additionally implement `RuntimeWithResult` to hand a JSON-compatible leaf
response back before success criteria and outputs are evaluated; the base
`Runtime` interface is unchanged.

`contenttrust.Analyze` is an explicit parallel path. It validates the document,
combines core expression recovery with consumer-supplied operation resolvers,
propagates provenance and capability to a fixed point, and returns stable
advisory edges and findings. It does not participate in execution or ordinary
validation results.

Browser profile validation is similarly explicit. Core UWS validates the
browser source binding or extension marker; browser-aware tooling selects the
exact separate profile/call schema and runs the additional context, origin,
input, verification, and cross-document checks before starting a browser.

## Artifact And Validation Flow

`versions/*.json` are the source documents for schema distribution. The
generator produces a deterministic ZIP embedded by `schemas`. Consumers can
locate repository, configured, module-cache, embedded, or sibling schemas.
`validation` loads a document through `convert`, applies the selected schema,
then calls the semantic validator.

`mock-fixtures.1.0` has its own exact JSON Schema, embedded lookup helper, and
protected JSON/Markdown digests. The `mockruntime` codec validates against that
schema and applies additional exact-key uniqueness checks. Request keys hash
RFC 8785 canonical JSON for the resolved request object; fixture replay always
reuses a matched response and performs no storage or network activity.

The public `mockruntime.Runtime` implements `uws1.Runtime` and the optional
`uws1.RuntimeWithResult` response handoff. It resolves request expressions,
records bounded canonical would-be requests in memory, and selects exact
fixtures or caller-supplied examples/schemas. Supported schema synthesis is
deterministic and bounded; unsupported schema features and expressions return
errors. The pure runtime makes no transport calls or file writes. A
caller-provided resolver remains caller-owned and should use local data for
pure simulation. The separate `HybridRuntime` requires explicit enablement and
passes only declared UWS 1.12+ `read` operations to a caller-owned delegate;
writes and unknown effects stay on the mock runtime.
For `$steps` expressions, the mock runtime reads successful step records from
the current iteration, then enclosing iterations in the same workflow call.
Sibling iterations and other workflow invocations are excluded; duplicate
matches at the nearest visible level remain ambiguous.

Schema conformance and parity tests connect the latest core schema to Go rules,
tags, known fields, and specification tables. SHA-256 fixtures enforce exact
membership and bytes for every published JSON document and every non-changelog
Markdown document; archive tests ensure embedded JSON bytes match their
sources. CI adds race testing, vet, diff checks, and strict documentation
builds.

## Public Contracts And Dependencies

- Core public APIs are the `uws1.Document` model, validation methods, execution
  entry points, `Runtime`, additive `RuntimeWithResult`, `Orchestrator`, and
  execution records.
- `convert`, `validation`, and `schemas` are the public interchange and schema
  integration surfaces.
- `contenttrust.Analyze` and `contenttrust.Resolver` are the advisory analysis
  surfaces.
- Browser authentication/registration packages and `runtimes` expose inert
  profile wire types and extension helpers, not executors.
- JSON Schema, YAML, HCL, XPath, JSONPath, synchronization, and module lookup
  libraries are implementation dependencies; provider SDKs and browser drivers
  are intentionally absent.

## Operational Boundaries

The repository has no production service, database, state store, credential
store, worker fleet, or provider deployment. Its deployed operational surface
is the MkDocs site; its released behavior is principally the Go module and
versioned documents. Persistence, retries across process failure, runtime
authorization, source fetching, and live provider compatibility remain
consumer responsibilities.

Known baseline discrepancies are preserved in the matching archives rather
than silently normalized. Some historical downstream compatibility claims
require external repositories.

## Memory Bank Ownership

`product.md` owns current product terminology, user workflows, relationships,
and invariants. This file owns current system layout, data flow, public
contracts, and archive registration. `tech-stack.md` owns toolchain and
verification commands, `lessons.md` retains only reusable evidenced lessons,
and `milestone.md` plus its linked status files own active delivery work.

Active milestone specifications and status records are mutable execution state.
After acceptance and bounded review they are consolidated into current truth
and retired through the history protocol in `milestone.md`. Candidate
directions remain unnumbered until a promotion proposal is approved. Verified
archives below are a separate, immutable evidence namespace and never become
active status work.

## Archive Baselines

Archive files are frozen repository snapshots. Current product and system truth
lives in this memory bank; use each archive only at its recorded baseline.

| Archive Lane | Context |
|---|---|
| C | Core workflow contract and execution |
| T | Advisory content-trust analysis |
| B | Browser capability and account lifecycle profiles |
| X | Source admission and extension profiles |
| M | Cross-cutting interchange, validation, and distribution tooling |

| Archive | Context | Baseline | Coverage | Supersedes |
|---|---|---|---|---|
| [Archive C01](../docs/archive-C01.md) | Core workflow contract and execution | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [Archive T01](../docs/archive-T01.md) | Advisory content-trust analysis | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [Archive B01](../docs/archive-B01.md) | Browser capability and account lifecycle profiles | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [Archive X01](../docs/archive-X01.md) | Source admission and extension profiles | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [Archive M01](../docs/archive-M01.md) | Interchange, validation, and distribution tooling | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |

C10 optional `OperationShape.Browser` preserves complete native metadata and raw
schema numeric/value/presence information. Browser-only strict decoding refuses
HTTP fields and aliases, retaining old non-browser table serialization. Binding,
flow and portability scan core references in browser request bodies while native
profile/private templates stay opaque. Signed64 (1.8) and safe-integer (1.9/1.10)
proof uses exact rationals and a rooted supported-dialect/reference-free subset.
Unproved references, dialects, open objects and symbolic arrays stay indeterminate.
Source-neutral diagnostics are value-free metadata, never execution readiness.
