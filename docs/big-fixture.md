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

Full fixture files:

- [`testdata/big/big.hcl`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.hcl)
- [`testdata/big/big.json`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.json)
- [`testdata/big/main.go`](https://github.com/OpenUdon/uws/blob/main/testdata/big/main.go)

Feature pages include short excerpts from these files. The excerpts are intentionally incomplete;
they show a feature in context without duplicating the full fixture.
