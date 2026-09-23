# C02 History

Milestone: `C02`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-C02.md`
Source specification: `tabilet/memory-bank/milestone.md#c02---version-compatibility`
Evidence: `947990b1edbff0345c560c6e4a33468d8824bd82`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `1`
Verification: `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`
Consolidated into: [version selection guidance](../../../docs/09-Validation.md), [version-gate implementation](../../../uws1/validation_version.go), [compatibility corpus](../../../validation/version_compatibility_test.go), [updated lesson](../../memory-bank/lessons.md#review-declared-version-semantics-before-freezing-markdown)

## Final milestone specification

````markdown
## C02 - Version Compatibility

**Goal.** Align schema selection, semantic validation, and execution rules
with a document's declared UWS version.

**Scope.** Decide and document a policy for published releases, prereleases,
and unpublished versions. Gate the 1.9.2 reference-step and criterion-pointer
rules, plus 1.5 step inputs, to declared versions rather than applying later
rules retroactively. Test a cross-version schema/semantic corpus with explicit
exceptions; reconcile validation layers and release wording. This plan adopts
version gating, not a retroactive erratum. The eventual release number is
chosen under that policy, not assumed to be 1.9.3.

**Acceptance.** Older published documents retain their declared-version
behavior; prerelease and unpublished-version handling is explicit; immutable
historical artifacts remain unchanged. Focused schema/semantic parity,
validation, full and race tests, vet, and diff check pass.

**Dependencies.** None. **Downstream impact.** C03 uses this policy for any
new versioned specification and compatibility claims.
````

## Final status

````markdown
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

**Review gate.** Iteration 1 of 10 passed 2026-09-23. Whole-milestone review
confirmed exact schema lookup, SemVer prerelease ordering, declared-version
gates, cross-version coverage, and unchanged immutable artifacts. No P1, P2,
or higher-severity finding remains.

| Item | State | Notes |
|---|---|---|
| Decide version and schema selection policy | `[+]` | Completed 2026-09-23. Published documents use the exact matching schema artifact; prerelease and unpublished versions require a matching artifact and never fall back to the latest published schema. Feature gates use SemVer precedence, including prerelease ordering, and malformed prerelease identifiers are rejected. Documented in the validation guide and changelog without modifying immutable specifications or schemas. `go test ./uws1 ./schemas ./validation` and `git diff --check` passed. Owners: A5, C13, D12. |
| Gate semantic and execution rules by declared version | `[+]` | Completed 2026-09-23. Semantic validation now rejects child blocks on reference steps only from 1.9.2 and step-local inputs before 1.5. JSON Pointer criterion evaluation rejects noncanonical array-index tokens from 1.9.2 while preserving earlier numeric-index behavior. Regression tests cover 1.9.1/1.9.2, 1.4/1.5, and execution of a leading-zero pointer index under the older version. Existing tests using step inputs now declare UWS 1.5. `go test ./uws1 -count=1` and `git diff --check` passed. Owner: A4. |
| Add cross-version compatibility corpus | `[+]` | Completed 2026-09-23. Added exact-schema plus semantic-validation baseline fixtures for all thirteen published UWS 1.x versions, with explicit boundary cases for 1.1 operation timeouts, 1.5 step inputs, 1.9.1 content trust, and 1.9.2 reference-step blocks. Confirmed unpublished stable and prerelease versions fail when their exact schema is unavailable. Intentional exception: the temporary Ansible source kind remains valid at 1.6 and is rejected from 1.7; all frozen schemas retain their exact historical behavior. `go test ./validation -run 'TestPublishedCoreSchemaAndSemanticCompatibilityCorpus|TestVersionedCompatibilityBoundaries|TestUnpublishedVersionsRequireTheirExactSchemaArtifact' -count=1` and `git diff --check` passed. Owner: E5. |
````
