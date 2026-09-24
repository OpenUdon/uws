# Big Fixture

The repository includes a deliberately large fixture under `testdata/big/`.
It is generated from Go structs, validated, serialized to JSON/YAML/HCL, and
round-tripped back to canonical JSON in `testdata/big/main.go`.

The fixture covers:

- OpenAPI-bound operations selected by `openapiOperationId` and `openapiOperationRef`.
- Extension-owned operations using every `uws.runtime.1.0` runtime selector.
- Sequence, parallel, switch, loop, await, and merge workflows.
- Triggers, routes, structural results, criteria, actions, components, and `x-*` extensions.
- HCL extension blocks and JSON/HCL dynamic-key round trips.

`testdata/big/big.json` is a generated test-runtime fixture, not a portable
expression-grammar or conformance fixture. It deliberately retains UWS 1.1.0
while exercising implementation-defined `$signals.*`, `$error.*`,
`$context.*`, `$workflow.*`, `$workflows.*`, and `$selected.*` sources, plus
other test-runtime expression operators. The exact file is excluded from the
core-grammar lint by a file-level legacy-fixture disposition; this does not
make its expressions conforming, and it is not evidence that those
implementation-defined expressions are portable. Core-only examples and
conformance evidence must use syntax accepted by the document's
declared-version specification.

Full fixture files:

- [`testdata/big/big.hcl`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.hcl)
- [`testdata/big/big.json`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.json)
- [`testdata/big/main.go`](https://github.com/OpenUdon/uws/blob/main/testdata/big/main.go)

Feature pages include short excerpts from these files. The excerpts are intentionally incomplete;
they show a feature in context without duplicating the full fixture.
