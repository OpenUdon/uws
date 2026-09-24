# B02 History

Milestone: `B02`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-B02.md`
Source specification: `tabilet/memory-bank/status-B02.md#status-b02---browser-compatibility-and-text-safety`
Evidence: `2b1a5944fea8bee66b7951aa8edfa2e4770add18`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `2`
Verification: focused Browser 1.8/1.9 profile tests; `go test ./schemas -count=1`; `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`; `jq empty versions/browser.1.9.json`; repeated `go generate ./schemas` produced an identical archive
Consolidated into: [Browser 1.9 profile](../../../versions/browser.1.9.md), [Browser 1.9 schema](../../../versions/browser.1.9.json), [Browser profile validation](../../../schemas/browser19_templates.go), [profile dispatch](../../../schemas/schema.go), [current product facts](../../memory-bank/product.md), and completed [C04 release work](status-C04.md)

## Final milestone specification

````markdown
## B02 - Browser Compatibility And Text Safety

**Goal.** Restore valid Browser 1.8 profile acceptance and publish a separately
versioned Browser 1.9 contract for deterministic literal braces and safe text
sinks.

**Scope.** Browser 1.8 validation permits isolated braces in any string,
including `Price {USD}` and regular-expression quantifiers, but rejects every
double-brace template or malformed attempt outside approved sinks. Browser 1.8
published artifacts remain unchanged. Browser 1.9 defines one-pass
`{{{{` / `}}}}` literal-brace escapes, safe path/query placement and URL
encoding, resolved text control/bidi rejection, text-only prompts, value-setting
`type_text` without keypress or form submission, BCP 14 language, current core
references, and a shared safe-integer boundary for parameters and
accessibility-text integer outputs. Browser 1.9 is explicitly opt-in; empty
profile-schema lookup remains Browser 1.8.

**Dependencies.** Completed B01 established Browser 1.8. C04 consumes the
accepted Browser 1.9 reference in UWS 1.11 release surfaces.

**Compatibility.** Browser 1.5–1.8 remain accepted and immutable. The Browser
1.8 validator correction restores behavior permitted by its existing profile
contract. No core schema or wire-model change is part of B02. No browser driver
exists in this repository; concrete runtime behavior is specified but not
claimed as driver-tested.

**Acceptance.** Regression tests cover the original brace/regex rejection,
malformed and misplaced templates, Browser 1.9 escapes and sinks, control/bidi
values, safe-integer boundaries, explicit lookup, empty-lookup compatibility,
historical profile acceptance, and immutable artifact membership. The
generated embedded archive and release documentation match the published
profile. Focused tests, full and race suites, vet, strict MkDocs, diff checks,
and the bounded whole-milestone review gate pass.
````

## Final status

````markdown
# Status B02 - Browser Compatibility And Text Safety

**Goal.** Restore valid Browser 1.8 profile acceptance and publish a separately
versioned Browser 1.9 contract for deterministic literal braces and safe text
sinks. Accepted 2026-09-23; Browser 1.9 is an opt-in published profile.

**Scope and decisions.** Browser 1.8's validator now allows ordinary standalone
braces inside and outside template sinks, including `Price {USD}` in
`type_text.value` and regex quantifiers in schema fields, while rejecting
every `{{` or `}}` sequence outside sinks, including malformed attempts. The
published Browser 1.8 schema/spec bytes remain unchanged. Browser 1.9 defines
unambiguous text-sink escapes (`{{{{` for literal `{{`, `}}}}` for literal
`}}`, with no rescan), percent-encoded navigation literals, prompt control/
bidi handling, and value-setting `type_text` without keypress or form-submit
events. It also supplies BCP 14 language, current core references, and
consistent safe-integer bounds. The published schema, static validator, and
tests cover the contract; concrete supplied-value validation and execution
remain runtime responsibilities.

**Dependencies.** [B01](../docs/history/status-B01.md) completed Browser 1.8.
**Downstream impact.** [C04](status-C04.md) must use the accepted Browser 1.9
reference in its new core release and guides.

**Compatibility and rollout.** Browser 1.5–1.8 remain accepted and immutable;
the Browser 1.8 validator correction restores profiles its published contract
already permits. Browser 1.9 is opt-in. No prior schema, spec, or archive is
rewritten; default lookup and profile-version behavior must be checked.
Rollback is selection of a previously published profile for documents that
do not rely on 1.9 behavior.

**Acceptance and verification.** Tests reproduce the brace/regex rejection
before the fix and pass after it. Browser 1.9 schema, spec, validation helpers,
fixtures, lookup, immutable-hash membership, generated archive, and release
surfaces agree. Exercise malformed placeholders, literal braces, URL and text
sinks, control/bidi values, safe integers, and historical profile compatibility.
Run focused browser/schema tests, `go test ./...`, `go test -race ./...`,
`go vet ./...`, `mkdocs build --strict`, `git diff --check`, and the bounded
whole-milestone review gate described in [milestone.md](milestone.md).

**Review provenance.** Source: “UWS review, pass 2: 1.10.0, browser 1.8 and
the remediation (2026-09-23)” (`uws-review-2.md`), stated baseline `5923cc0`;
revalidated at full HEAD `5923cc0d3b2c852e45aee3771201c642b9daeca9`.
No relevant uncommitted browser code or versioned-document change formed the
evidence. The worktree contained only pre-existing untracked Tabilet GOAL,
archives, and evolution files. Lineage: completed B01 and its Browser 1.8
release; do not reopen its retired status. This is ordinary review intake,
not a review-gate iteration.

**Follow-up plan review.** The pasted “Plan review of second-review
remediation (2026-09-23)” supplied no source priority or baseline; revalidated
at full HEAD `5923cc0d3b2c852e45aee3771201c642b9daeca9` with the
previous reconciliation's uncommitted planning edits and pre-existing
untracked Tabilet evidence present. No browser implementation edit formed
this evidence. F02 below refined N1/B02's first pending row; it is ordinary
intake, not a review-gate pass.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| N1 | P2 | P2 | Confirmed: `schemas/browser_templates.go` rejects standalone braces outside sinks; `testdata/browser-profile/read-only.yaml` has an output validation schema that can contain regex quantifiers. |
| N7 | P2 | P2 | Partially confirmed: `versions/browser.1.8.md` has no control/bidi text-sink policy or value-setting guarantee; an actual browser-driver exploit was not demonstrated. |
| N17 | P3 | P3 | Partially confirmed: `versions/browser.1.8.md` lacks BCP 14 language, contains stale core-version references, and specifies a wider integer range than its output conversion. |
| F02 | not supplied | P2 | Confirmed: `schemas/browser_templates.go` rejects single literal braces within a template sink; Browser 1.8 text only forbids unmatched double braces. Lineage: N1 and completed B01. |

**Review gate.** Iteration 1 of 10 started 2026-09-23. Deep-review of the
Browser 1.8 validator correction and Browser 1.9 profile definition, dispatch,
archive/hash protections, compatibility behavior, test corpus, release
surfaces, and three task commits found a P2 status drift and a P3 wording issue.
The status still called B02 proposed, described publication and profile
dispatch as future work, claimed Browser 1.5–1.7 (not 1.5–1.8) were the
accepted immutable range, and said row 3 had not integrated the helper. The
Browser 1.9 JSON Schema also described all integer outputs as safe integers,
although the normative output bound applies to accessibility-text outputs.
The status and schema description were corrected; the new schema digest was
updated and the embedded archive regenerated. `go test ./schemas -count=1`,
`jq empty versions/browser.1.9.json`, `mkdocs build --strict`, and
`git diff --check` passed after the fixes.

Iteration 2 of 10 started 2026-09-23. Re-review of the full milestone and the
iteration-1 fixes, including profile compatibility, deterministic archive/hash
membership, source references, text validation boundaries, and release-surface
links found one P3 status-note issue: row 1 named a version-selection test that
was later renamed, and row 2 did not qualify its safe-integer output summary
as accessibility-text output. Both notes were corrected, and the current row-1
focused test command passed. `go test ./...`, `go test -race ./...`,
`go vet ./...`, `mkdocs build --strict`, `git diff --check`, and repeated
`go generate ./schemas` all passed. Review passed on iteration 2 of 10 on
2026-09-23: no P1, P2, or higher finding remains. No browser driver exists in
this repository, so no driver-runtime behavior is claimed as tested.

| Item | State | Notes |
|---|---|---|
| Add regression tests and repair Browser 1.8 validator acceptance | `[+]` | Completed 2026-09-23. Tests first reproduced rejection of a regex quantifier, ordinary locator/description braces, and `Price {USD}`; the validator now treats only `{{` and `}}` as template markers, permits isolated braces in all strings, and still rejects valid or malformed double-brace sequences outside approved sinks. Navigation and text sinks accept literal braces. Browser 1.8 schema/spec artifacts remain unchanged. `go test ./schemas -run 'TestBrowser18TemplatesAcceptDeclaredScalarsInApprovedSinks|TestBrowser18TemplatesRejectUnsafeOrAmbiguousCases|TestBrowserSourceProfileVersionSelectionKeepsOptInAndDefaultCompatibility' -count=1`, `go test ./schemas -count=1`, and `git diff --check` passed. Owner: N1/F02. |
| Define and validate Browser 1.9 text-sink contract | `[+]` | Completed 2026-09-23. Added the Browser 1.9 schema/spec, a one-pass brace-escape tokenizer (`{{{{` / `}}}}`) with strict placeholder placement, safe-integer default validation, and static/resolved-default checks for prohibited controls and bidi characters in `type_text.value` and prompts. The normative contract makes navigation escapes `%7B%7B` / `%7D%7D`, retains component-encoded URL substitutions, makes `type_text` set a field value without keypress/Enter/submit events, renders prompts as text only, adopts BCP 14 language, and aligns integer parameters and accessibility-text output values to ±9007199254740991. Added focused positive and negative YAML fixtures and unit cases for escapes, malformed/misplaced templates, controls/bidi, and integer boundaries. `jq empty versions/browser.1.9.json`, `go test ./schemas -run 'TestBrowser19' -count=1`, and `git diff --check` passed. Its helper was integrated by row 3; Browser 1.8 remains unchanged. Owners: N7/N17 and new-profile N1. |
| Integrate and verify Browser 1.9 distribution | `[+]` | Completed 2026-09-23. Added explicit Browser 1.9 lookup and profile validation while preserving empty lookup as Browser 1.8; added positive/negative public-dispatch and context coverage, SHA-256 membership, and regenerated the embedded version archive. Updated README, current docs, navigation, changelog, AGENTS, and the current product fact; preserved the replaced product wording in the knowledge journal. Browser 1.5–1.8 remain accepted and immutable. `go test ./schemas -count=1`, `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed; repeated `go generate ./schemas` produced an identical archive. Whole-milestone review passed on iteration 2 below. |
````
