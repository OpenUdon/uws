# Lessons

Read only the topics relevant to the current change. Product terminology,
architecture contracts, commands, and active work remain in their dedicated
memory-bank files.

## State Profile Compatibility As A Minimum Core Version

**Scope.** Separately versioned profiles referenced by a core specification.

**Lesson.** State the minimum compatible core version in a profile document,
not whichever core release is current when that profile is published. Check
those references when a later core version ships.

**Rationale.** Browser 1.9 remains opt-in and compatible with UWS 1.9 and
later, but its prose still calls UWS 1.10.0 current after UWS 1.11.0 shipped.

**Evidence.** `versions/browser.1.9.md` §§1 and 8;
`versions/1.11.0.md`; third-review R6.

## Release Surfaces Must Move With A Versioned Release

**Scope.** Publishing or advancing a core version or separately versioned
profile.

**Lesson.** Reconcile the top-level README, agent instructions, documentation
home/reference navigation, feature guides and examples, release changelog,
and current memory-bank facts in the same release task as the versioned
artifacts. Check examples against the published grammar and actual behavior,
not only against document structure.

**Rationale.** Browser registration 1.2 was fully represented by schemas, Go
helpers, fixtures, and version documents, while several discovery surfaces
continued to identify 1.1 as current. Tests of the executable contract did not
detect that documentation drift. The UWS 1.10 release similarly left feature
guides with stale or contradictory execution examples; passing Go and schema
tests did not detect that mismatch.

**Evidence.** `versions/browser-registration.1.2.*`;
`versions/browser-registration-call.1.2.*`; `README.md`; `AGENTS.md`;
`mkdocs.yml`; `tabilet/docs/archive-B01.md`; `versions/1.10.0.md`;
`docs/04-Triggers-and-Route-Dispatch.md`; `docs/05-Structural-Results.md`;
`docs/06-Success-Criteria-and-Actions.md`; `docs/07-Execution-Model.md`.

## Separate Editorial Correction From Semantic Amendment

**Scope.** Corrections to an already published Markdown contract.

**Lesson.** A wording correction may repair a stale description without a new
contract version only when it leaves wire shape and meaning unchanged. Any
scope or semantic change needs the changelog's amendment treatment or a new
version before immutable bytes are updated.

**Rationale.** Freezing inaccurate prose preserves the error, while silently
changing normative meaning defeats versioned compatibility. The repository's
changelog already draws this boundary explicitly.

**Evidence.** `versions/CHANGELOG.md`; `versions/browser.1.7.md`;
`versions/runtime.1.0.md`; `tabilet/docs/archive-M01.md`.

## Review Declared-Version Semantics Before Freezing Markdown

**Scope.** Cross-version semantic changes in a repository that freezes published
specifications.

**Lesson.** Compare exact schema selection and semantic/execution behavior
against each declared SemVer before finalizing a release or freezing its
Markdown. Gate features by SemVer precedence, require a matching schema for
prerelease and unpublished versions, and record deliberate exceptions
explicitly.

**Rationale.** The current validator applies 1.9.2 reference-step constraints
only to documents declaring 1.9.2 or later; older versions retain their
declared behavior. Published Markdown and schemas are digest-protected, so
compatibility fixes must be version-gated instead of silently editing history.

**Evidence.** `versions/CHANGELOG.md`; `docs/09-Validation.md`;
`validation/version_compatibility_test.go`; `tabilet/docs/history/status-C02.md`.
