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
