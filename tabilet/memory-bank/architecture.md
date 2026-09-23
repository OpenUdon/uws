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
only leaf execution, expression evaluation, and item resolution.

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

Schema conformance and parity tests connect the latest core schema to Go rules,
tags, known fields, and specification tables. SHA-256 fixtures freeze every
published JSON document, and archive tests ensure embedded bytes match their
sources. CI adds race testing, vet, diff checks, and strict documentation
builds.

## Public Contracts And Dependencies

- Core public APIs are the `uws1.Document` model, validation methods, execution
  entry points, `Runtime`, `Orchestrator`, and execution records.
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
than silently normalized: one browser profile sentence still calls 1.9.0 the
current core schema, runtime prose retains present-tense Ansible wording,
Markdown specifications lack the JSON documents' automated digest guard, and
some historical downstream compatibility claims require external repositories.

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
| [C01](../docs/archive-C01.md) | Core workflow contract and execution | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [T01](../docs/archive-T01.md) | Advisory content-trust analysis | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [B01](../docs/archive-B01.md) | Browser capability and account lifecycle profiles | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [X01](../docs/archive-X01.md) | Source admission and extension profiles | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
| [M01](../docs/archive-M01.md) | Interchange, validation, and distribution tooling | `8382d0f26b3b10870125760643078d1a1a3e31b6` | verified | none |
