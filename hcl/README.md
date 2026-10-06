# UWS verified HCL views

This separate module renders inert typed HCL from exact JSON/YAML bytes and
independently verifies the source projection and provenance. It does not load
source references, validate workflow semantics, approve packages or execute.

```go
source := hcl.Source{Format: hcl.JSON, Bytes: approvedJSON}
options := hcl.Options{Revision: exactCodecCommit} // full lowercase Git SHA
view, err := hcl.Render(ctx, source, options)
// Render returns no view unless independent verification passes.
err = hcl.Verify(ctx, source, view, options)
```

Deprecated `Import(ctx, hclBytes)` reconstructs JSON from the inert supported
typed mapping. Its output remains a proposal requiring consumer validation and
fresh confirmation/approval. Variables/functions/interpolation/comprehensions,
ambiguous fields/keys, unsupported block mappings and malformed inputs refuse.
Import never calls HCL expression evaluation with a runtime context.

Typed block/label names are derived from the public UWS model tags without
decoding values through custom model unmarshaling. Root/owned extensions use
`extensions { ... }`; inline typed object values use an `extensions` object.
Dynamic dollar keys retain `_ref`/`_id`/`_schema`/`_defs` and `__dollar__foo`
mapping, with `__uws_literal__` escaping for literal collisions. Null or empty
repeated blocks use an explicit attribute with the same mapped HCL field name,
such as `operation = []`; this is a derived-view mapping, never legacy package
input. Unsupported extension attribute names refuse rather than losing keys.

Exact JSON numeric tokens are emitted directly and independently reconstructed
from inert numeric AST source ranges, including large integers, decimals,
exponents and signed zero. JSON/YAML source decoding and HCL verification never
use float64 or JCS equality. YAML aliases/merges/tags, non-string keys and numeric
spellings outside JSON numeric syntax refuse. JSON duplicate members/trailing
values/invalid Unicode refuse. Strings remain literal, including template markers.

Bounds: 8 MiB source and HCL, 100,000 work nodes, nesting depth 100. Writer growth
is checked before append and context cancellation is propagated. Parser/library
work is cooperative; consumers own hard CPU/RSS/deadline/mount/network isolation
and private-value filtering. Errors are stable and contain no input excerpts.

The standalone module currently consumes exact accepted root UWS C09 for public
model tags and strict JSON helpers. That closure still contains Horizon/HashiCorp
HCL; no core-dependency removal is claimed. Final root/codec publication and
consumer pins remain M08.4 work. The operator-owned go.work is unchanged.

Run `GOWORK=off GOPROXY=off go test ./...`, its `-race` variant and `go vet ./...`
from this module. Root-module tests do not discover a nested module automatically.
