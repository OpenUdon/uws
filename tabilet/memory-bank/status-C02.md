# Status C02 - Version Compatibility

**Acceptance.** Declared published versions select documented schema and
semantic behavior; prerelease and unpublished-version policy is explicit.
Later 1.9.2 rules and 1.5 step inputs do not apply retroactively. A
cross-version schema/semantic corpus documents intentional exceptions and
passes alongside focused validation tests, full tests, race tests, vet, and
diff check. Earlier published artifacts remain immutable.

**Dependencies.** None. **Downstream impacts.** C03 follows the resulting
version policy and does not assume a text-only 1.9.3 or retroactive erratum.

**Review provenance.** “UWS spec and architecture review: 1.9.2 and earlier
(2026-09-23)”; stated baseline
`8382d0f26b3b10870125760643078d1a1a3e31b6`; revalidated at full HEAD
`e765fed4ba0a241e9481f99fc324c6a8afb9a8a1` with untracked
memory-bank initialization content present. No relevant uncommitted
schema/validator change formed the evidence. The approved decision is to gate
1.9.2 rules by declared version, not to declare them retroactive errata.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| A4 | P1 | P1 | Confirmed in `uws1/validation_workflow.go`, `uws1/execution_criteria.go`, and `versions/1.9.1.json`; newer rules reach older declared documents. |
| A5 | P1 | P2 | Partially confirmed in `uws1/validation_version.go` and `schemas/schema.go`: numeric gating ignores SemVer prerelease ordering; schema lookup has no stated prerelease policy. |
| C13 | P2 | P2 | Confirmed: `schemas/schema.go`, versioned schemas, and `versions/1.9.2.md` do not give a portable declared-version-to-schema rule. |
| D12 | P3 | P2 | Confirmed policy inconsistency between `versions/1.9.2.md` §1.1 and release history in `versions/CHANGELOG.md`. |
| E5 | Recommendation | P2 | Confirmed test gap: current version fixtures and `schemas/version_immutability_test.go` do not form a cross-version schema/semantic parity corpus. |

**Review gate.** Ordinary intake; 0 of 10 iterations. Start only after task
acceptance checks.

| Item | State | Notes |
|---|---|---|
| Decide version and schema selection policy | `[+]` | Completed 2026-09-23. Published documents use the exact matching schema artifact; prerelease and unpublished versions require a matching artifact and never fall back to the latest published schema. Feature gates use SemVer precedence, including prerelease ordering, and malformed prerelease identifiers are rejected. Documented in the validation guide and changelog without modifying immutable specifications or schemas. `go test ./uws1 ./schemas ./validation` and `git diff --check` passed. Owners: A5, C13, D12. |
| Gate semantic and execution rules by declared version | `[ ]` | Isolate 1.9.2 reference-step and canonical-pointer rules and 1.5 step inputs; add older-version regressions. Owner: A4. |
| Add cross-version compatibility corpus | `[ ]` | Compare each published 1.x schema and semantic validator with declared-version fixtures; list intended exceptions and test release/prerelease behavior. Owner: E5. |
