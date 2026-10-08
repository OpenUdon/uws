# Technical Stack

## Accepted M09 ordinary modules

Root and nested codec accepted source is
`b099f6803277ae94c7e9f1da0904a0140b278f20`, version
`v0.0.0-20261008043726-b099f6803277`, whole review 7/10. Independent ordinary
module downloads match full source provenance, Go archive/GoMod sums and every
reviewed filename/content. Root and codec module-only full tests/races/vet/build
pass with retained Go 1.26.6, GOWORK=off, GOTOOLCHAIN=local and GOPROXY=off
after acquisition. Complete selected closures contain 52 root / 50 codec
modules. Standalone codec retains C09 root; the 52-module ordinary consumer
explicitly selects both corrected modules and passes races without replacements.
Use owned root-filesystem temporary directories for the retained qualification;
the separate /tmp quota is not inferred from root filesystem free space.
[Publication](../../docs/m09-publication.md) records actual ordinary ZIP hashes,
canonical sums, full closures and proof locations. No declared dependency,
workspace, browser, frozen-consumer or installed pin is changed by this release.

## Stage 11 public foundations

Accepted C08/C09 add no dependency upgrade, source/provider parser or private
runtime import. Exact sources, vector/table hashes and review evidence remain
in [C08 qualification](../../docs/c08-qualification.md) and
[C09 qualification](../../docs/c09-qualification.md).

UWS 1.13 adds exact schema/spec membership and regenerated embedded archive;
all 74 pre-existing immutable version documents and legacy/root APIs remain
unchanged. Root go.mod/sum and operator-owned go.work are unchanged.
The separate `github.com/OpenUdon/uws/hcl` module lives in `hcl/`, using accepted
root C09 `v0.0.0-20261006181058-6a267306032e`, existing HCL `v2.24.0`,
cty `v1.17.0` and YAML `v3.0.1`. The graph retains Horizon/HCL compatibility.

Run root and nested full tests/races/vet independently with GOWORK=off,
GOPROXY=off. Nested tests include a self-contained byte-pinned six-source corpus
and external-package public API test; a module-only copy qualifies packaging.
CI has a distinct nested-module job because root ./... cannot discover it.
Bounds are 8 MiB source/view, 100,000 work nodes and depth 100. CPU/RSS/deadline/
mount/network controls remain with consuming workers; no hard resource claim.
M08 is accepted at `c0b19385a3b034cd45de16726668b9150f0633f2` after review 3;
root/codec source is independently published/resolved.
[qualification](../../docs/m08-qualification.md) records exact artifact hashes.

The [Stage 11 coordinator](../../../kinet/docs/stage11.md) requires published
root/codec source and independently observed closure before adoption. The
installed M44 service and current consumer/browser pins remain unchanged.

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
`browserregistration`, `expressions` and `binding`.
`hcl/` is a separate Go module with its own go.mod/sum and verification job.
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

The main CI workflow runs root and nested full tests, race tests, vet, and the root diff check. Pull
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
