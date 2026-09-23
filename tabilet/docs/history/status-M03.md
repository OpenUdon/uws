# M03 History

Milestone: `M03`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-M03.md`
Source specification: `tabilet/memory-bank/milestone.md#m03---hcl-key-interoperability`
Evidence: `e765fed4ba0a241e9481f99fc324c6a8afb9a8a1`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `1`
Verification: `go test ./uws1 ./convert`; `go test ./...`; `go test -race ./...`; `go vet ./...`; `git diff --check`
Consolidated into: [HCL key mapping](../../../uws1/hcl.go), [HCL regression tests](../../../uws1/hcl_test.go), [conversion regression test](../../../convert/convert_test.go), [HCL unmarshaler guidance](../../../uws1/unmarshaler.md), [downstream contract](../../memory-bank/status-C03.md)

## Final milestone specification

````markdown
## M03 - HCL Key Interoperability

**Goal.** Preserve dynamic JSON key identity across HCL conversion or reject
unrepresentable keys explicitly; remove direct-HCL decoding ambiguity.

**Scope.** Define an injective spelling for ordinary underscore-prefixed and
dollar-prefixed keys in dynamic maps, including `variables`, nested
`request.body`, and extensions. Preserve existing HCL decoding where possible;
make any unavoidable compatibility change explicit. Do not silently reinterpret
ordinary `_id`, `_ref`, or `__dollar__x` as dollar keys.

**Acceptance.** Focused conversion and direct-HCL regression cases cover
`_id`, `_ref`, `__dollar__x`, `$id`, and nested forms, proving preservation or
explicit failure. `go test ./...`, `go test -race ./...`, `go vet ./...`, and
`git diff --check` pass. Historical contracts remain frozen.

**Dependencies.** None; lineage to completed [M02](../docs/history/status-M02.md)
is historical only. **Downstream impact.** C03 must describe only the accepted
mapping and decide whether a portable HCL appendix belongs in a later release.
````

## Final status

````markdown
# Status M03 - HCL Key Interoperability

**Acceptance.** Dynamic key spelling is injective across JSON/HCL conversion
or fails explicitly; focused conversion and direct-HCL tests cover `_id`,
`_ref`, `__dollar__x`, `$id`, and nested request, variable, and extension
maps. Full tests, race tests, vet, and diff check pass. Existing HCL decoding
is preserved where possible, with any incompatibility documented.

**Dependencies.** None. Completed [M02](../docs/history/status-M02.md) is
historical release-guard lineage, not a prerequisite.

**Downstream impacts.** C03 records the accepted conversion semantics:
`__uws_literal__` escapes ordinary keys that collide with legacy dollar-key
spellings, while legacy dollar decoding remains supported. A normative
cross-runtime HCL mapping remains a separate candidate.

**Review provenance.** “UWS spec and architecture review: 1.9.2 and earlier
(2026-09-23)”; stated baseline
`8382d0f26b3b10870125760643078d1a1a3e31b6`; revalidated at full HEAD
`e765fed4ba0a241e9481f99fc324c6a8afb9a8a1` with untracked
memory-bank initialization content present. No relevant uncommitted converter
change formed the evidence. At that baseline, the reported JSON-to-HCL silent
corruption example was not reproduced because `MarshalHCL` rejected ambiguous
keys before encoding; direct HCL decoding used legacy dollar aliases for raw
`_ref` and `__dollar__name` keys. M03 adds a literal escape while retaining
those legacy interpretations.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| A1 | P0 | P1 | Partially confirmed in `uws1/hcl.go` (key transforms and direct decode) and `convert/convert_test.go`; preflight rejection narrows the reported round-trip claim. |

**Review gate.** Passed on iteration 1 of 10 on 2026-09-23. Review confirmed
that literal underscore/prefix keys are escaped before legacy aliases, dollar
keys retain their prior spellings, and `fromHCLKey` reverses the new escape
without colliding with either legacy representation. Focused, full, race, vet,
and diff checks pass; no P1, P2, or higher finding remains.

| Item | State | Notes |
|---|---|---|
| Define and implement injective HCL key representation | `[+]` | Completed 2026-09-23. Added `__uws_literal__` escaping for ordinary keys colliding with legacy dollar spellings, retained legacy dollar decoding, and covered `_id`, `_ref`, `__dollar__x`, escape-prefix, `$id`, nested request, variable, and extension maps. `go test ./uws1 ./convert`, `go test ./...`, `go test -race ./...`, `go vet ./...`, and `git diff --check` passed. Owner: A1. |
````
