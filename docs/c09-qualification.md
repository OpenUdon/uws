# C09 binding contract qualification

C09 provides the bounded source-neutral shape/resolver contract, advisory
binding validation and deterministic flow observations. This metadata surface
never parses provider source formats, loads keys or grants execution.

Full offline Go tests, full races, vet, strict documentation and diff checks
passed over the unchanged module/lockfile, core model, generated archive and
published-version bytes. Package dependency inspection confirms no APItools
or private-runtime dependency. Table tests include closed fields, duplicate
JSON keys, forged identity, selector ambiguity, snapshot mutation, partial
schemas/security, and fabricated non-HTTP details.

The [pinned supplement](examples/binding/v1/manifest.json) contains unchanged
accepted OpenUdon step fixtures from revision
`c2f161d762bc9f2217bbf0c34b00cdef64b0f7d0` plus a source-neutral projection.
The source byte digest and selector/key agree with the accepted runnable
check. Its required-input/type, response-field and complete-authentication
findings agree with ValidateBinding. Missing inputs, incompatible values and
missing authentication produce the corresponding negative outcomes. The
separate indeterminate fixture's effect/account/purpose uncertainty remains
outside this metadata check; a compatible binding is not a complete workflow
assessment or permission.

Binding regressions qualify unavailable schema refs without file/network
loading, unproved expression constraint containment, canonical response array
indexes, security OR/AND/scopes, cancellation, wrong custom-resolver selectors
and value-free reports. Flow regressions qualify deterministic ordering,
nonmutation/privacy, possible reachability, unused core outputs, declared
effects and pending ordering, work bounds, merge-child exclusion and cycles.

Publication uses the already confirmed STG11_SOURCE_PUBLICATION UWS origin/main
scope with normal fast-forward pushes and independent observed reachability.
Consumers wait for the whole-milestone review and accepted published closure;
local builds and producer completeness flags are not authority.

Closing review iteration 1 reproduced two supported binding scenarios: nested
expression templates were wrongly treated as literals, and number-to-integer
was wrongly classified disjoint. Recursive template projections now prove the
supported closed-object/finite-array case while retaining indeterminate results
for unsupported constraints/open-object uncertainty. New regressions cover both
fixes, array uniqueness uncertainty and open source object constraints. Full
tests/races/vet passed after correction.

Iteration 2 also qualified equivalent typed Go JSON containers containing
expressions. The value projection is normalized with UseNumber before template
discovery, preserving numeric precision and rejecting invalid/duplicate-key
values. Full tests/races/vet passed again, including the typed-map equivalence
regression. This changes only metadata checks and no source/provider behavior.
