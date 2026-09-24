# Knowledge History

## 2026-09-23 - Versioned Semantics Lesson Updated

**Source heading.** `tabilet/memory-bank/lessons.md#review-declared-version-semantics-before-freezing-markdown`

**Previous wording.**

> **Scope.** Cross-version semantic changes in a repository that freezes published
> specifications.
>
> **Lesson.** Compare schema, validator, and executor behavior for older declared
> versions before finalizing a new release's prose or freezing its Markdown.
> Record deliberate exceptions and compatibility policy explicitly.
>
> **Rationale.** The current validator applies 1.9.2 reference-step constraints
> to 1.9.1 documents, while the 1.9.1 schema accepts the shape. The published
> Markdown is now digest-protected, so a later semantic correction needs a new
> version or a separately justified amendment rather than an unnoticed edit.
>
> **Evidence.** `versions/1.9.1.json`; `uws1/validation_workflow.go`;
> `schemas/version_immutability_test.go`; `tabilet/docs/history/status-M02.md`.

**Reason.** C02 implemented and verified declared-version gating: 1.9.2
reference-step restrictions now apply only from 1.9.2, 1.5 step inputs only
from 1.5, and criterion pointer index behavior follows the declared version.
The previous wording described a known defect as current truth. Exact schema
selection and SemVer prerelease policy are now documented, with all published
core schemas covered by the compatibility corpus.

**Evidence.** `d4b0766`; `846d707`; `947990b`;
`validation/version_compatibility_test.go`.

**Replacement.** [Review Declared-Version Semantics Before Freezing Markdown](../../memory-bank/lessons.md#review-declared-version-semantics-before-freezing-markdown)

## 2026-09-23 - Release-Surface Lesson Expanded To Core Guides

**Source heading.** `tabilet/memory-bank/lessons.md#release-surfaces-must-move-with-a-profile-release`

**Previous wording.**

> **Scope.** Publishing or advancing a separately versioned profile.
>
> **Lesson.** Reconcile the top-level README, agent instructions, documentation
> home/reference navigation, release changelog, and current memory-bank facts in
> the same release task as the profile artifacts.
>
> **Rationale.** Browser registration 1.2 was fully represented by schemas, Go
> helpers, fixtures, and version documents, while several discovery surfaces
> continued to identify 1.1 as current. Tests of the executable contract did not
> detect that documentation drift.
>
> **Evidence.** `versions/browser-registration.1.2.*`;
> `versions/browser-registration-call.1.2.*`; `README.md`; `AGENTS.md`;
> `mkdocs.yml`; `tabilet/docs/archive-B01.md`.

**Reason.** The second review found that a core 1.10 release also left feature
guides and examples stale or inconsistent with published grammar and behavior.
The profile-only scope was too narrow; this expansion records a current
release-review lesson, not a claim that the guides are already fixed.

**Evidence.** `versions/1.10.0.md`;
`docs/04-Triggers-and-Route-Dispatch.md`;
`docs/05-Structural-Results.md`;
`docs/06-Success-Criteria-and-Actions.md`;
`docs/07-Execution-Model.md`; `tabilet/docs/history/status-C04.md`.

**Replacement.** [Release Surfaces Must Move With A Versioned Release](../../memory-bank/lessons.md#release-surfaces-must-move-with-a-versioned-release)

## 2026-09-23 - Current Browser Profile Updated

**Source heading.** `tabilet/memory-bank/product.md#current-contract-surface`

**Previous wording.**

> UWS 1.10.0 is the current core contract. Browser 1.8 is the current browser
> capability profile; browser authentication/call 1.1 is current for sign-in;
> browser registration/call 1.2 is current for reviewed registration
> verification; and registration input 1.0 is the private envelope format.
> Runtime Supplement 1.0 remains the public metadata floor for common non-HTTP
> extension operations. UWS 1.10 adds versioned portable execution semantics and
> expression-addressable-name restrictions while preserving the 1.x wire
> vocabulary; earlier published contracts remain accepted according to their
> version gates. Ansible support is historical UWS 1.6 material only.

**Reason.** B02 published the separately versioned, opt-in Browser 1.9 profile.
Browser 1.8 remains immutable and is still the empty-lookup compatibility
default, but it is no longer the latest published browser capability contract.

**Evidence.** `versions/browser.1.9.json`; `versions/browser.1.9.md`;
`schemas/schema.go`; `schemas/version_immutability_test.go`.

**Replacement.** [Current Contract Surface](../../memory-bank/product.md#current-contract-surface)
