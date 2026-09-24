# Plan-Review Contract Refinement v4

## Origin

The user approved the pasted “Plan review of second-review remediation
(2026-09-23)” after [direction v3](prompt-v3.md). That review supplied no
separate source priority or baseline. Revalidation used full HEAD
`5923cc0d3b2c852e45aee3771201c642b9daeca9` with the prior planning
edits and existing untracked Tabilet evidence present. This is ordinary
planning intake, not an implementation or milestone review-gate pass.

## Approved Direction

Repair the direct-version fail-open path first in new C05, independently of
Browser 1.9. Then run B02, C04, and M04. C04 incorporates C05's exact-version
admission rule and B02's Browser 1.9 reference into a new UWS 1.11 release.
M04 remains product-independent and last in execution order. All four statuses
are required, with downstream impacts C05 -> C04 and B02 -> C04.

## Refined Contract Choices

- A UWS 1.11 `goto` from nested execution unwinds the entire top-level run.
  Resolve a globally declared target in root scope after unwinding, not to a
  prior caller-scoped invocation. Terminal or in-flight root target records
  fail explicitly; a distinct scoped call record does not count. There is no
  continuation after the target.
- A UWS 1.11 merge over `forEach` uses iteration records when present,
  otherwise the parent record for skipped or zero-item dependencies. Ordinary
  dependency failure still aborts merge.
- UWS 1.11 admits loop-context `$batchIndex` alongside the already approved
  response dot-walk and numeric-field literal grammar. Lint examples by their
  declared version; keep legacy conversion fixtures explicitly identified and
  add conforming 1.11 fixtures.
- The new 1.11 publication audits an explicit clause/reference checklist,
  including replacement of the unsupported prerelease example, without
  rewriting frozen 1.10.0 artifacts or changing historically accurate prose.

## Boundary

The [milestone index](../memory-bank/milestone.md) and linked status files own
task rows, dependencies, verification, and review gates. This evolution pair
records approved direction only. It does not implement a finding, track the
existing Tabilet evidence, edit GOAL or archives, or authorize a commit.
