# Status C05 - Exact-Version Admission

**Goal.** Close the fail-open direct-API path for documents declaring a UWS
version without a matching published schema before the Browser 1.9 and UWS
1.11 release work. This is pending implementation, not a delivered fix.

**Scope.** Apply exact declared-version admission consistently to
`Document.Validate()`, `ValidateResult()`, execution entry points, and file
validation. Valid SemVer syntax alone is insufficient: unsupported
prereleases, unpublished patch/minor versions, and future versions fail
closed. Use the existing C02 publication policy; do not silently select the
latest available schema or newest semantics. Correct the misleading
leading-zero version error and test that the latest-schema constant matches
the highest published core SemVer. This admission correction intentionally
applies to every declared version, including older documents; it is a
cross-version bug-fix exception, not a new 1.11-only feature.

**Dependencies.** Completed [C02](../docs/history/status-C02.md) established
exact file-schema selection and SemVer gates; completed
[C03](../docs/history/status-C03.md) published 1.10.0. No B02 dependency.
**Downstream impact.** [C04](status-C04.md) must carry the admitted-version
rule and cross-version exception into UWS 1.11's normative policy and release
notes; C05 may complete independently before Browser 1.9 work.

**Compatibility and rollout.** Preserve all published schemas/specifications
and the behavior of supported declared versions. Previously accepted
unsupported version strings now fail explicitly at every entry path; no
fallback or migration of their semantics is offered. A caller must select an
actually published version. Keep `versions/CHANGELOG.md` and the validation
guide clear about this deliberate cross-version correction. A rollback would
restore the unsafe acceptance path, so compatibility should be evaluated
against supported versions rather than treating fail-open behavior as a
contract.

**Acceptance and verification.** Failing tests first for direct and file
paths using `1.99.0`, an unpublished patch such as `1.10.1` while absent,
`1.0.0-beta.1`, and malformed leading-zero `1.010.0`; supported published
versions continue to pass. Exercise `Validate`, `ValidateResult`, `Execute`
and other relevant entry points, schema-selection parity, and the
latest-published-version guard. Update the mutable changelog and validation
guide for the cross-version exception without editing frozen 1.10.0 files.
Run focused version/validation tests, `go test ./...`, `go test -race ./...`,
`go vet ./...`, `mkdocs build --strict`, `git diff --check`, and the bounded
whole-milestone review gate in [milestone.md](milestone.md).

**Review provenance.** Source: “UWS review, pass 2: 1.10.0, browser 1.8 and
the remediation (2026-09-23)” (`uws-review-2.md`), stated baseline `5923cc0`;
and pasted “Plan review of second-review remediation (2026-09-23)”, source
baseline `not supplied`. Both were revalidated at full HEAD
`5923cc0d3b2c852e45aee3771201c642b9daeca9`. The worktree contained the
previous reconciliation's uncommitted planning edits and pre-existing
untracked Tabilet GOAL, archives, and evolution files; no relevant
uncommitted core implementation or versioned-document change formed the
evidence. Lineage: completed C02/C03 and C04's previously pending admission
row, moved here before execution; retired records remain unchanged. This is
ordinary review intake, not a bounded-gate iteration.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| N3 | P2 | P2 | Confirmed: `uws1/validation.go` and direct execution accept unsupported declared versions while `validation/validation.go` fails exact file selection. |
| N14 | P3 | P3 | Confirmed: `uws1/validation.go` mislabels leading-zero syntax, and `uws1` lacks a highest-published-SemVer assertion for its latest-schema constant. |
| F08 | not supplied | P2 | Confirmed sequencing risk: `tabilet/memory-bank/milestone.md` placed C04 after B02 solely for final Browser 1.9 references, delaying N3's fail-open correction. C05 is the independent owner. |
| F09 | not supplied | P2 | Confirmed: exact-version admission must apply to older declared versions and be documented as an intentional cross-version correction in `versions/CHANGELOG.md`; C04 later reflects it without duplicating this implementation. |

| Item | State | Notes |
|---|---|---|
| Enforce exact declared-version admission and edge-case guards | `[+]` | Completed 2026-09-23. Added a registry for the published core schema versions and made `Document.ValidateResult` reject unsupported prerelease, unpublished patch/minor, and future versions; `Validate` and `Execute` share the rejection. Invalid SemVer now reports that diagnosis without incorrectly naming prerelease syntax. Regression tests cover direct Validate/ValidateResult/Execute, every published version, file validation, exact registry/schema membership, and latest-schema SemVer ordering. `go test ./uws1 ./validation -run 'TestValidate_RequiresPublishedCoreVersion|TestValidate_LeadingZeroVersionDiagnostic|TestLatestUWSSchemaIsHighestPublishedCoreVersion|TestPublishedUWSVersionRegistryMatchesSchemaArtifacts|TestValidateDocumentFileRejectsUnpublishedVersions' -count=1`, `go test ./...`, `go test -race ./...`, `go vet ./...`, and `git diff --check` passed. Supported historical versions remain admitted. Owners: N3/N14/F08. |
| Document compatibility exception and verify distribution | `[+]` | Completed 2026-09-23. Updated `versions/CHANGELOG.md` and `docs/09-Validation.md` to distinguish SemVer syntax from published-version support, state that external schemas cannot establish core-version support, and record the intentional all-version admission exception without retroactively applying later feature rules. Added a file-validation regression test with an exact external schema. `go test ./validation -run 'TestValidateDocumentFileRejectsUnpublishedVersions|TestValidateDocumentFileRejectsUnpublishedVersionWithExternalSchema' -count=1`, `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed. C04 will carry this rule into 1.11 without reopening frozen history. Owner: F09. |
