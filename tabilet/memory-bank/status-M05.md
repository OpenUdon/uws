# Status M05 — Browser 1.10 selector match-count profile

**State:** Active. Approved planning; specification work is pending.

**Goal.** Define and publish Browser 1.10 as a separately versioned profile that
returns a CSS-selector match count as a typed nonnegative integer without
exposing page text or element attributes.

**Compatibility.** Browser 1.9 and older profiles remain immutable. UWS core
stays at 1.11 unless a generic profile-binding change proves necessary. No
browser runtime integration or live target action is included.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| M05.1 Document the Browser profile versioning procedure | [x] | Added the maintenance checklist in [future-source-profiles.md](../../docs/future-source-profiles.md#adding-a-browser-profile-version) and linked it from README. |
| M05.2 Define Browser 1.10 count semantics | [ ] | Specify selector match-count output, integer bounds, zero and multiple-match behavior, and closed failure semantics. |
| M05.3 Update profile schema, guide, examples and conformance vectors | [ ] | Preserve existing profile bytes, add Browser 1.10 coverage, and keep the versioning checklist current. |
| M05.4 Verify and review the published specification | [ ] | Run focused and full UWS checks and bounded review; publish exact reviewed source for downstream consumers. |
