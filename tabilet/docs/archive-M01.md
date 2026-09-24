# Archive M01 - Interchange, Validation, And Distribution Tooling

**Context.** Cross-cutting Go APIs and repository automation that serialize,
validate, package, test, document, and distribute UWS contracts.

**Baseline.** 8382d0f26b3b10870125760643078d1a1a3e31b6

**Coverage.** verified

**Supersedes.** none

## Scope And Responsibilities

This context turns versioned UWS documents and Go models into consumable
library and documentation artifacts. It owns JSON/YAML/HCL conversion, document
loading, structural-plus-semantic validation, schema lookup and embedding,
deterministic generation of the embedded version archive, conformance fixtures,
CI, and the MkDocs site.

It does not own the meaning of core workflow fields, browser profile semantics,
runtime implementations, deployment of an executor, or storage of workflow
state. Those contracts are inputs to this tooling.

## Domain And Workflows

The `convert` package translates between JSON, YAML, and HCL while preserving
the UWS document model. JSON and YAML retain flattened `x-*` fields. HCL uses
`extensions` blocks and reversible rewriting for dollar-prefixed keys.
Conversion rejects ambiguous or irreversible inputs, preserves numeric lexemes
where required, and marshals HCL from a deep copy rather than mutating the
caller's document.

The `validation` package loads JSON, YAML, or HCL, resolves the schema matching
the document's declared version, applies JSON Schema validation, and then runs
the Go semantic validator. It also discovers `.uws.json`, `.uws.yaml`, and
`.uws.yml` artifacts beneath a directory.

`versions/` is the document source of truth. The schema generator sorts its
JSON members, assigns fixed ZIP metadata, and produces
`schemas/version-documents.zip`, which is embedded into the Go package. Schema
lookup supports repository, configured-directory, module-cache, embedded, and
sibling-package use cases. Profile validators add discriminator and semantic
checks beyond raw JSON Schema.

## System Shape

`convert` depends on `uws1`; `validation` composes `convert`, `schemas`,
`uws1`, and the JSON Schema validator; `schemas` is otherwise independent of
the core Go model and embeds version documents. The internal generator has one
directional responsibility from `versions/*.json` to the embedded ZIP.

Repository tests bridge schema, specification, model tags, known-field lists,
semantic rules, profile validators, fixtures, and immutable document bytes.
MkDocs builds user-facing guides from `docs/` and links normative reference
documents in `versions/` through the repository.

## Contracts And Dependencies

Published JSON contracts are frozen by a filename-to-SHA-256 manifest. Adding a
new JSON contract requires adding its digest, while changing or removing an
existing member fails tests. A separate test requires the embedded archive to
match every current version JSON document byte-for-byte.

Schema, prose, Go validation, tags, and known-field lists are expected to stay
coordinated. The public validation facade asserts formats for UWS documents and
reports schema or semantic failures with file context. Profile schema accessors
return independent byte slices and validators reject unknown discriminators.

## Operations And Verification

The main CI job runs `go test ./...`, `go test -race ./...`, `go vet ./...`,
and `git diff --check`. Pull requests also install pinned documentation
requirements and run `mkdocs build --strict`. A separate workflow deploys the
documentation site from `main` when `docs/**` or `mkdocs.yml` changes.

At this baseline, the ordinary test suite, vet, strict MkDocs build, and diff
check all pass. The strict build reports that `registration-authoring.md` exists
outside configured navigation. Generated site output is ignored and is not a
versioned product artifact.

## Evidence

| Claim | Repository Evidence |
|---|---|
| JSON, YAML, and HCL conversions preserve the supported document model and extensions. | `convert/convert.go`; `convert/convert_test.go`; `convert/roundtrip_test.go` |
| File validation composes versioned JSON Schema and semantic validation. | `validation/validation.go`; `validation/validation_test.go` |
| Version documents are embedded through a deterministic generated archive. | `schemas/schema.go`; `internal/generateversionarchive/main.go`; `schemas/schema_test.go` |
| Published JSON contracts are byte-frozen. | `schemas/version_immutability_test.go`; `versions/` |
| Schema/model/spec parity has explicit automated coverage. | `uws1/schema_conformance_test.go`; `uws1/schema_parity_test.go`; `uws1/schema_artifacts_test.go` |
| CI verifies tests, race behavior, vet, formatting, and strict documentation. | `.github/workflows/test.yml` |
| MkDocs publishes the documentation site from the main branch. | `mkdocs.yml`; `.github/workflows/deploy-docs.yml` |

## Observed Gaps

Automated immutability checks cover published JSON documents but not their
Markdown specification counterparts. The documentation deployment trigger does
not include `versions/**`, and `docs/registration-authoring.md` is absent from
MkDocs navigation. Those facts do not break the current strict build.
