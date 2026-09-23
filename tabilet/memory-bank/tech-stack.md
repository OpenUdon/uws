# Technical Stack

## Language And Module

- Go 1.25.4
- Module: `github.com/OpenUdon/uws`
- Library and specification distribution; there is no repository-owned service,
  database, worker deployment, or concrete provider runtime.

The main Go packages are `uws1`, `contenttrust`, `convert`, `validation`,
`schemas`, `runtimes`, `browserauthentication`, and `browserregistration`.
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
