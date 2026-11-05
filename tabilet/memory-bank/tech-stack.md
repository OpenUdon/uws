# Technical Stack

## Accepted C10 browser-profile shapes

C10 is accepted after whole review1 at implementation
`f01a2542410c583d0ea909dadd8f17527cc0d27d`, ordinary root module
`v0.0.0-20261009220914-f01a2542410c`. [Publication proof](../../docs/c10-publication.md)
records independent source/archive checks and all seven frozen consumer builds.
Native source production, actual consumer adoption, browser/runtime containment,
credential policy and execution authority remain with their owners. Existing
published profile/schema bytes, root/HCL declarations and actual pins are unchanged.


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

Every test, check and qualification follows [test-policy.md](test-policy.md).
Verification standard: `standard`; this upgrade adds no duration budgets or
instrumentation exclusions. Existing package limits and exact-artifact
qualification requirements remain binding.

`REVIEW_MAX_ITERATIONS: 10` is the project review ceiling. Resolve each run's
positive limit independently of commit policy and use the minimum of the
human-approved request, this project cap and any milestone cap. Raising this
project ceiling needs a separate explicit human-approved rule change; an
upgrade, repair or resume never resets consumed iterations.

**Row checks.** Before completing a row or starting its task review or re-review,
run `GOWORK=off GOPROXY=off GOTOOLCHAIN=local go build ./...` and
`GOWORK=off GOPROXY=off GOTOOLCHAIN=local go vet ./...` separately from the root
and `hcl/`, then run Go tests for the affected whole packages in their owning
module. Add direct consumers for changes to public contracts, persisted formats
or relied-upon behavior. Add the strict documentation build for documentation
site changes. Add race checks on changed concurrent units and consumers of a
changed concurrency contract. Earlier narrow edit-loop checks do not replace
this row check. Toolchain test-cache results are permitted when the cache sees
all inputs; record cached results as cached. Record the command, full source
revision plus a saved diff/hash for uncommitted work, PASS or FAIL, and one log
link in ignored scratch, including qualified artifact identity where applicable.

**Required gates.** Run the following Go commands separately from the repository
root and `hcl/`; root `./...` does not discover the nested module. Use the
installed local toolchain and cached dependencies; missing inputs stop the
affected check rather than triggering a download.

```bash
GOWORK=off GOPROXY=off GOTOOLCHAIN=local go build ./...
GOWORK=off GOPROXY=off GOTOOLCHAIN=local go vet ./...
GOWORK=off GOPROXY=off GOTOOLCHAIN=local go test -count=1 ./...
GOWORK=off GOPROXY=off GOTOOLCHAIN=local go test -race -count=1 ./...
```

Also run these repository-root gates:

```bash
mkdocs build --strict
git diff --check
```

Run gates cheapest first, immediately before the whole-milestone review.
`-count=1` bypasses Go's test-result cache for every gate run, including reruns;
compiler and dependency caches stay. Before later review iterations or closure,
reuse a passed class O result only when nothing the gate reads has changed;
otherwise rerun it. Record the command, full source revision plus saved diff/hash
for uncommitted work, artifact identity where qualified, PASS or FAIL and one
ignored-scratch log link; a reused result states its inputs are unchanged.
Race checks keep full project-code instrumentation. Run heavy instrumented
targets serially on a shared host under the stock capacity recommendation;
measured changes to concurrency or instrumentation scope are recorded here.

**Stage qualification.** None scheduled for UWS: its retained C11 and C13 rows
are withdrawn. Existing accepted exact-artifact and standalone-module proofs
remain evidence at their recorded contexts. A future approved stage requiring
end-to-end qualification assigns that required gate to its final milestone,
which depends on the other stage milestones and uses real pinned artifacts.

**Class R checks.** None scheduled. A later named real-effect check requires
separate human authority, scope and supervision; this policy grants none.

**Existing other-language tooling.** The Python-based MkDocs documentation build
retains its documented role. No other-language Go test runner or test-evidence
generator is approved.

Focused edit-loop commands include:

```bash
go test ./uws1 -run TestSchemaConformance
go test ./convert -run TestRoundtrip
go test ./schemas -run TestPublishedVersionDocumentsAreImmutable
```

These narrow checks do not replace row checks or required gates. The existing
CI workflow runs root and nested full tests, race tests, vet and the root diff
check; pull requests also run the strict documentation build. The commands above
specify local gate freshness without changing CI configuration.

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

C10 adds no dependency or tool upgrade. The actual public root ZIP has396 files
and native archive sum `h1:TTjAn0++TdGHY0vahbXw3XsICQaavBytoqTSHu0K6o0=`;
GoMod sum `h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`. Public ZIP
compression differs from the provisional container, while all source bytes and
Go checksums match. Go1.26.6/offline isolated cache checks qualify root and all
seven disposable consumer units; retained unchanged HCL checks remain separate.
