# C08 expression reference qualification

C08 adds the public expressions parser, scoped reference evaluator and optional
core-field portability diagnostics. The mock delegates to this implementation
with explicit adapters for its two historical numeric/pointer extensions.
This is an additive library/conformance publication, not a new UWS version,
provider runtime or hosted permission.

Full offline `go test ./...`, `go test -race ./...`, `go vet ./...` and
`mkdocs build --strict` passed with the retained Go1.26.6 executable over the
unchanged Go1.25.4 module directive. `git diff --check` passed. Existing schema,
archive/conformance and published-version immutability tests passed. No module
or lockfile changed; uws1 and all published version bytes remain untouched.

The [18-case supplement](examples/expressions/v1/cases.json) and its
[manifest](examples/expressions/v1/manifest.json) are separately versioned/pinned
from frozen published contracts. Source/runtime tests additionally cover every
core source, comparison token/whitespace, exact numeric lexemes and large
integers, pointer escapes/indexes, version/field/loop scope, source-to-source
comparison, context cancellation, cyclic inputs, component precedence and
workflow/nearest-iteration visibility. Optional portability tests prove
ordinary validation compatibility, profile/query exclusions, deterministic
bounded diagnostics, typed nested bindings and empty-output refusal.

Consumers prepare exact json.Number value projections before constructing the
UWS model snapshot; the old model decoder's float64 behavior is unchanged.
The evaluator performs no source/provider/credential I/O and never dispatches
a leaf. Input/caller document admission and authorization remain consumer-owned.

Publication is the owner-confirmed STG11_SOURCE_PUBLICATION scope: normal
fast-forward source/closure pushes to the exact UWS origin/main, with fetched
ancestry and independently observed commit reachability. Downstream adoption
waits for whole-milestone acceptance and its final published closure. No tag,
force push, unrelated source, installation or live provider action is included.
