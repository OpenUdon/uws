# Status B02 - Browser Compatibility And Text Safety

**Goal.** Restore valid Browser 1.8 profile acceptance and publish a separately
versioned Browser 1.9 contract for deterministic literal braces and safe text
sinks. This is proposed work, not an implemented profile.

**Scope and decisions.** Browser 1.8's validator must allow ordinary standalone
braces inside and outside template sinks, including `Price {USD}` in
`type_text.value` and regex quantifiers in schema fields, while rejecting
every `{{` or `}}` sequence outside sinks, including malformed attempts. Preserve
the published Browser 1.8 schema/spec bytes. Browser 1.9 will define an
unambiguous text-sink escape (`{{{{` for literal `{{`, `}}}}` for literal
`}}`, with no rescan), keep URL literals percent-encoded, specify prompt
control/bidi handling, and make `type_text` set a field value rather than
emitting keypresses or submitting a form. Browser 1.9 also supplies BCP 14
language, current core references, and safe-integer consistency. The
implementation task must settle exact rejection/escaping details in the new
profile and tests before publication.

**Dependencies.** [B01](../docs/history/status-B01.md) completed Browser 1.8.
**Downstream impact.** [C04](status-C04.md) must use the accepted Browser 1.9
reference in its new core release and guides.

**Compatibility and rollout.** Browser 1.5–1.7 remain accepted and immutable;
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
this evidence. F02 below refines N1/B02's first pending row; it is ordinary
intake, not a review-gate pass.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| N1 | P2 | P2 | Confirmed: `schemas/browser_templates.go` rejects standalone braces outside sinks; `testdata/browser-profile/read-only.yaml` has an output validation schema that can contain regex quantifiers. |
| N7 | P2 | P2 | Partially confirmed: `versions/browser.1.8.md` has no control/bidi text-sink policy or value-setting guarantee; an actual browser-driver exploit was not demonstrated. |
| N17 | P3 | P3 | Partially confirmed: `versions/browser.1.8.md` lacks BCP 14 language, contains stale core-version references, and specifies a wider integer range than its output conversion. |
| F02 | not supplied | P2 | Confirmed: `schemas/browser_templates.go` rejects single literal braces within a template sink; Browser 1.8 text only forbids unmatched double braces. Lineage: N1 and completed B01. |

| Item | State | Notes |
|---|---|---|
| Add regression tests and repair Browser 1.8 validator acceptance | `[+]` | Completed 2026-09-23. Tests first reproduced rejection of a regex quantifier, ordinary locator/description braces, and `Price {USD}`; the validator now treats only `{{` and `}}` as template markers, permits isolated braces in all strings, and still rejects valid or malformed double-brace sequences outside approved sinks. Navigation and text sinks accept literal braces. Browser 1.8 schema/spec artifacts remain unchanged. `go test ./schemas -run 'TestBrowser18TemplatesAcceptDeclaredScalarsInApprovedSinks|TestBrowser18TemplatesRejectUnsafeOrAmbiguousCases|TestBrowserSourceProfileVersionSelectionKeepsHistoricalProfiles' -count=1`, `go test ./schemas -count=1`, and `git diff --check` passed. Owner: N1/F02. |
| Define and validate Browser 1.9 text-sink contract | `[+]` | Completed 2026-09-23. Added the Browser 1.9 schema/spec, a one-pass brace-escape tokenizer (`{{{{` / `}}}}`) with strict placeholder placement, safe-integer default validation, and static/resolved-default checks for prohibited controls and bidi characters in `type_text.value` and prompts. The normative contract makes navigation escapes `%7B%7B` / `%7D%7D`, retains component-encoded URL substitutions, makes `type_text` set a field value without keypress/Enter/submit events, renders prompts as text only, adopts BCP 14 language, and aligns integer parameters/output conversions to ±9007199254740991. Added focused positive and negative YAML fixtures and unit cases for escapes, malformed/misplaced templates, controls/bidi, and integer boundaries. `jq empty versions/browser.1.9.json`, `go test ./schemas -run 'TestBrowser19' -count=1`, and `git diff --check` passed. The helper awaits explicit profile dispatch in row 3; Browser 1.8 remains unchanged. Owners: N7/N17 and new-profile N1. |
| Integrate and verify Browser 1.9 distribution | `[+]` | Completed 2026-09-23. Added explicit Browser 1.9 lookup and profile validation while preserving empty lookup as Browser 1.8; added positive/negative public-dispatch and context coverage, SHA-256 membership, and regenerated the embedded version archive. Updated README, current docs, navigation, changelog, AGENTS, and the current product fact; preserved the replaced product wording in the knowledge journal. Browser 1.5–1.8 remain accepted and immutable. `go test ./schemas -count=1`, `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed; repeated `go generate ./schemas` produced an identical archive. The whole-milestone review gate remains pending below. |
