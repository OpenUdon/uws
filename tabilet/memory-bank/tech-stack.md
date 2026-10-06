# Technical Stack

## M08.1 presentation contract progress

Exact 1.13 schema/specification membership is added while all older published
hashes remain unchanged. The core declared-version map and latest conformance
references select 1.13; prior version gates/wire fields stay intact. The embedded
archive is regenerated with go generate ./schemas. Existing dependency versions
and legacy HCL APIs are unchanged; the nested codec module is still pending.


## C08 accepted reference foundation

C08 is accepted/retired after review 2 at `0411eea6fc84fbd6aa97cef94f53f301260f4844`, independently observed on origin/main. The [qualification](../../docs/c08-qualification.md) pins the source and supplement; M08 remains pending.

## C09 accepted binding foundation

C09 is accepted/retired after review 3 at `6a267306032edc687a298cefc8bba7019d3ad059`, independently observed on authorized origin/main. [Qualification](../../docs/c09-qualification.md) records the pinned source-neutral contracts and limitations; M08 remains pending.

## Approved Stage 11 tooling target — not implemented

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. Exact published module revisions, frozen build closures and standalone verification are acceptance gates. go test ./...; go test -race ./...; go vet ./...; schema/conformance and published-version immutability checks; mkdocs build --strict; git diff --check. Run the separate codec module checks once it exists.
Both phases belong to one stage. C08 acceptance/publication is recorded above; M08 remains pending and claim no installed Kinet behavior. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

## Language And Module

C09.3 flow checks use existing UWS types/core expression parsing; no dependency change. Binding/flow race/vet regressions cover determinism, privacy, cancellation, cycles, merge-child exclusion, pending/write order and static loop arrays.

C09.2 reuses the existing jsonschema/v6 dependency with a refusing URL loader and format assertions; no new dependency/version. Binding race/vet checks pass, including privacy, forged resolver, schema-ref and constraint-indeterminacy cases.

C09.1 adds binding using the standard library plus existing internal strict-JSON support. No module/lockfile changes; shape-table race/vet tests pass.

C08 adds expressions with standard-library parsing and a reference evaluator using existing uws1 context/strict-JSON support. No go.mod/go.sum or published-version bytes change. Parser/evaluator/mock race/vet and published-version immutability checks pass; the new 18-case reference corpus uses lossless value projections. C08.3 replaces the duplicate mock expression implementation with delegation and adds opt-in portability tests; no dependency version changed.

- Go 1.25.4
- Module: `github.com/OpenUdon/uws`
- Library and specification distribution; there is no repository-owned service,
  database, worker deployment, or concrete provider runtime.

The main Go packages are `uws1`, `contenttrust`, `convert`, `validation`,
`schemas`, `runtimes`, `mockruntime`, `browserauthentication`, and
`browserregistration`.
`internal/generateversionarchive` is a repository generator rather than a
public package.

## Formats And Libraries

- JSON Schema Draft 2020-12 documents define structural contracts.
- JSON and YAML are portable serializations; HCL is an authoring form with
  reversible extension and dollar-key handling.
- `jsonschema/v6` provides runtime schema compilation and validation.
- `yaml.v3` provides YAML parsing and emission.
- Horizon/dethcl and HashiCorp HCL dependencies support HCL conversion.
- XPath and JSONPath dependencies implement non-simple criteria.
- `x/sync/errgroup` supports concurrent orchestration.
- `github.com/gowebpki/jcs` implements RFC 8785 request canonicalization for
  Mock Fixture Format 1.0.
- `mockruntime` builds the public pure mock and explicit hybrid adapter on the
  existing orchestrator; it adds no provider, network-client, or persistence
  dependencies. The hybrid path calls only a caller-supplied read delegate.

Dependency versions are governed by `go.mod` and `go.sum`. Provider SDKs,
browser drivers, credential stores, and concrete runtime clients do not belong
in this repository unless a separately approved delivery outcome changes that
boundary.

## Versioned Artifacts

`versions/` is document-only and contains current, compatible, and historical
contracts. `schemas/version_immutability_test.go` enforces exact membership and
SHA-256 bytes for every published JSON document and every Markdown document
except the intentionally mutable `versions/CHANGELOG.md`.

After changing a JSON document, regenerate the deterministic embedded archive:

```bash
go generate ./schemas
```

Schema, specification, Go model, semantic validation, profile helpers,
fixtures, and the embedded archive must remain coordinated. Earlier published
artifacts are not edited for new semantics; create a new version instead.

## Documentation

- MkDocs with Material 9.7.7
- Source directory: `docs/`
- Configuration: `mkdocs.yml`
- Deployment: GitHub Pages through `.github/workflows/deploy-docs.yml`

Install the documentation dependency with:

```bash
python -m pip install -r docs/requirements.txt
```

## Verification Commands

```bash
go test ./...
go test -race ./...
go vet ./...
mkdocs build --strict
git diff --check
```

Focused commands include:

```bash
go test ./uws1 -run TestSchemaConformance
go test ./convert -run TestRoundtrip
go test ./schemas -run TestPublishedVersionDocumentsAreImmutable
```

The main CI workflow runs the full tests, race tests, vet, and diff check. Pull
requests also run the strict documentation build.

## Change Discipline

- Use `apply_patch` for hand edits and preserve unrelated worktree changes.
- Keep one active status-row owner and one task commit per completed row.
- Make schema/spec/code changes together when the public contract changes.
- Keep implementation, tests, fixtures, generated artifacts, and corrections
  to current memory-bank facts in the same task row.
- Treat browser authentication, registration, private inputs, and content-trust
  analysis as security-sensitive, fail-closed surfaces.
- Verify published-version membership and hashes before claiming immutability.
- Do not modify frozen archive files; create a successor archive when a new
  repository baseline is materially needed.
