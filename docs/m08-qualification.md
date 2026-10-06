# M08 verified HCL presentation qualification

UWS 1.13 adds the independently verified HCL presentation/deprecation contract
without altering UWS 1.12 wire/execution rules. The separate public module
`github.com/OpenUdon/uws/hcl` implements byte-based Render, independent Verify
and deprecated inert Import. Existing root conversion/model HCL APIs, ordinary
HCL input, browser compatibility and all older version bytes remain retained.

Exact release artifacts:

| Artifact | Bytes | SHA-256 |
|---|---:|---|
| versions/1.13.0.json | 57162 | ce118674a63cbaceb46bfbe898d85568aa9281175c1361097c044d7eb6270490 |
| versions/1.13.0.md | 117145 | ca731dbf91fc6fe28563554f88fc48433d65174d66623bb3ad5e41091534e089 |
| schemas/version-documents.zip | 136097 | e84046ebd0ee8742ac1e4ced9239eb854cbcc59296902f7b4b6e1e7269822efa |
| hcl/testdata/conformance/v1/manifest.json | 1703 | 97408edfe8d7d8147bd0950ac88c00e1757c8564c84992a93e0ec70b19fd5fd2 |
| hcl/go.mod | 1400 | a2a3515570eee58de7fd98edfe40f9bd0f214d540258fb594dd0d2809ef4de4a |
| hcl/go.sum | 12972 | 3d8f959da36092e62a340386482a9b74444235c01854815aa4d3e14bfd91cd7c |

The codec's standalone module graph consumes accepted root C09
`v0.0.0-20261006181058-6a267306032e` for public model tags/strict JSON helpers.
The 1.13 root and codec accepted source is
`c0b19385a3b034cd45de16726668b9150f0633f2` after whole review 3. Both resolve
`v0.0.0-20261006224744-c0b19385a3b0` to that exact source. Workers bind caller-supplied codec revision to their actual executable
and complete module/build closure. Horizon/HCL remain in the graph; there is
no HCL-free core or private runtime claim. Root go.mod/sum and operator-owned
go.work are unchanged; no dependency upgrade is included.

The self-contained supplement pins six byte-identical existing typed sources.
Tests exercise the full model-first corpus and inert retained browser metadata,
exact large integers/decimals/exponents/signed zero, Unicode/control/template
strings, dollar/escape-prefix collisions, extensions, empty/null/absence and
array order. Independent UseNumber projection and exact numeric strings are the
oracle; no float64/JCS equality is used. Public external-package compilation
and a disposable module-only copy qualify packaging without parent fixtures.

Malformed sources/HCL, duplicate/ambiguous keys/attributes/blocks/labels,
unsupported YAML aliases/merges/explicit tags, functions/traversals/interpolation/
comprehensions and failed provenance refuse without partial output or source
excerpts. Limits are 8 MiB source/view, 100,000 work nodes and depth 100, checked
writer growth and cooperative cancellation. Parser work is cooperative;
consumers own hard CPU/RSS/deadline limits, mounts/network and private-value
filtering. Views are derived metadata; Import returns a proposal requiring
consumer validation and fresh authority. No source loading, credentials,
workflow execution or file write is supplied.

Offline standalone root/nested full tests, full races and vet passed with
Go 1.26.6. Root parity/conformance/declared-version/archive checks passed; all
74 pre-existing immutable version documents retain their exact bytes. Strict
MkDocs, gofmt and git diff --check passed. The compatible staticcheck reports
no nested-module diagnostics. Root and nested modules have separate CI jobs.
Udon/OpenUdon root regressions passed as recorded in M08.1. Downloaded root/codec
archive full suites and an independent public consumer passed without go.work
or replacements. The [51-module closure](m08-module-closure.json) hashes to
`36a8a7a1152f9f49af5064294eca756bf9b00f4632dfa1d215112157991d1e96`;
[publication](m08-publication.md) identifies exact module sums/closure evidence.

Closing review 1 reproduced and fixed a reflection panic on wrong typed block
containers and silent YAML tag loss, including tags on keys. Regressions fail
before correction and pass afterward. Review 2 found stale release discovery
and current-fact references; latest links and consolidated current documents
are repaired, with literal old memory-bank wording preserved in the knowledge
journal. Whole review 3 found no remaining P1/P2 or higher issue. Exact source
publication and module resolution are independently verified; downstream
contracts are reconciled before retirement closure publication/next execution.

The confirmed STG11_SOURCE_PUBLICATION grant permits normal fast-forward UWS
origin/main source/closure publication. This repository's push workflow also
auto-deploys docs using a force push to another ref; that action is outside the
grant. Publication commits use `[skip ci]` to suppress push workflows, following
[GitHub's documented skip mechanism](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/skip-workflow-runs).
Required checks are established locally; no hosted CI/deployed-docs result is
claimed. No installation, live ledger, provider/API/model/mail or registration
action is included.
