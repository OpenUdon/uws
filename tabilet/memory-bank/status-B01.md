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

**Review gate.** Ordinary intake; 0 of 10 iterations. Start only after task
acceptance checks.

| Item | State | Notes |
|---|---|---|
| Specify safe substitution and scalar formatting | `[ ]` | Decide context-specific path/query/other encoding, number and Boolean formatting, and fail-closed ambiguity behavior with compatibility examples. Owner: C20. |
| Publish isolated browser profile update | `[ ]` | Add new version/schema/validator fixtures and documentation; prove exact selection and preservation of 1.5–1.7. Owner: C20. |
