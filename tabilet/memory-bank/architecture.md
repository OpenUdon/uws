# Architecture

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
