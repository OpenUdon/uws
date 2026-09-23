# B01 History

Milestone: `B01`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-B01.md`
Source specification: `tabilet/docs/history/status-B01.md#final-milestone-specification`
Evidence: `84eb1395b09eed392bfe136a2101c796f4650037`
Worktree: `includes uncommitted retirement reconciliation and pre-existing untracked memory-bank initialization content`
Review: `passed`
Review iterations: `1`
Verification: `go test ./schemas`; `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`
Consolidated into: [browser 1.8 profile](../../../versions/browser.1.8.md), [browser 1.8 schema](../../../versions/browser.1.8.json), [template validator](../../../schemas/browser_templates.go), [validator fixtures](../../../schemas/browser_templates_test.go), [release surfaces](../../../README.md), [downstream contract](../../memory-bank/status-C03.md)

## Final milestone specification

````markdown
## B01 - Browser Template Safety

**Goal.** Define safe, portable `{{param}}` substitution for browser bindings.

**Scope.** Specify context-sensitive escaping and numeric/Boolean formatting;
reject ambiguous or unsafe substitutions. Publish any changed binding under a
new browser profile version, with schema, validator, and fixtures. Preserve
browser 1.5–1.7 and exact profile selection. The repository has no browser
driver, so acceptance does not claim runtime exploit coverage.

**Acceptance.** Tests demonstrate encoded-safe and fail-closed path/query
cases, scalar formatting, version isolation, and validator behavior. Focused
profile tests, full and race tests, vet, strict MkDocs build, and diff check
pass.

**Dependencies.** None; lineage to the completed browser documentation work
in [M01](../docs/history/status-M01.md) is historical only. **Downstream
impact.** C03 must cross-reference the accepted browser profile, not copy its
contract into core.
````

## Final status

````markdown
# Status B01 - Browser Template Safety

**Acceptance.** Browser parameter substitution has context-sensitive safe
encoding and defined numeric/Boolean formatting; ambiguous or unsafe cases
fail closed. Any changed semantics receive a new browser profile version,
schema, validator, and fixtures with exact version selection. Browser 1.5–1.7
stay immutable and accepted. Focused profile tests, full tests, race tests,
vet, strict MkDocs build, and diff check pass. This repo contains no browser
driver, so no runtime exploit or driver-coverage claim is made.

**Dependencies.** None; completed [M01](../docs/history/status-M01.md) is
historical browser-documentation lineage only. **Downstream impacts.** C03
cross-references the accepted profile without importing its semantics.

**Review provenance.** “UWS spec and architecture review: 1.9.2 and earlier
(2026-09-23)”; stated baseline
`8382d0f26b3b10870125760643078d1a1a3e31b6`; revalidated at full HEAD
`e765fed4ba0a241e9481f99fc324c6a8afb9a8a1` with untracked
memory-bank initialization content present. No relevant uncommitted browser
profile change formed the evidence.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| C20 | P2 | P1 | Confirmed contract gap in `versions/browser.1.7.md` (template substitution and parameter typing); same-origin checks do not define path/query escaping. |

**Review gate.** Ordinary intake; iteration 1 of 10 passed after reviewing the
complete browser 1.8 profile publication, validator, tests, compatibility
surface, and documentation diff. No P1/P2 findings; no review fixes were needed.

| Item | State | Notes |
|---|---|---|
| Specify safe substitution and scalar formatting | `[+]` | Decision completed 2026-09-23 for browser 1.8. Placeholders are exactly `{{name}}` with names matching `[A-Za-z][A-Za-z0-9_-]*`; they substitute once and inserted text is never rescanned. Allowed sinks are navigate URL path segments, query values only, type_text/select_option values, and confirmation prompts. URL scheme/authority, query names, fragments, locators, waits, and all other fields are non-templatable. Path/query values percent-encode UTF-8 bytes using RFC 3986 unreserved characters only and uppercase hex; query structure is not interpolated. Dot-segment navigation is rejected, including a substituted path segment equal to `.` or `..`. Human text sinks use literal UTF-8 string values. Booleans format as lowercase JSON tokens; integers use canonical signed 64-bit base-10; numbers use RFC 8785 finite binary64 serialization; null, arrays, objects, non-finite/out-of-range numbers, malformed braces, unknown parameters, and non-scalar parameter schemas fail closed. Parameter values/defaults are validated against their declared schema before any sequence step. Browser 1.5–1.7 retain their historical behavior. Owner: C20. |
| Publish isolated browser profile update | `[+]` | Added immutable `browser.1.8.{json,md}`, static validation of placeholder syntax/placement/declarations/URL safety, accepted and rejected fixtures, latest-profile schema selection and embedded archive, SHA-256 membership, and current discovery pointers. Full tests/race/vet, strict MkDocs, and diff checks pass. Browser 1.5–1.7 artifacts remain unchanged and are accepted by exact selection tests. No browser driver exists here; concrete substitution/runtime behavior is specified but not claimed as driver-tested. Owner: C20. |
````

## Outcome

Browser 1.8 is an isolated profile release requiring UWS 1.9. It does not add a
UWS core schema, wire field, or core Go API. Its profile validator performs
static template-placement, parameter-declaration, scalar-type, URL-component,
and dot-segment checks; concrete parameter/default validation and substitution
remain runtime responsibilities. Browser 1.5–1.7 remain byte-for-byte
unchanged and selectable. No browser-driver implementation is present in this
repository.
