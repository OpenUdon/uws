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
| M05.2 Define Browser 1.10 count semantics | [x] | Defined `matchCount: true` for exact CSS match counts; optional `within` must resolve to exactly one root and scopes descendant matches; required `visibility` is `all` or `rendered` (rendered excludes hidden/zero-area elements but ignores viewport, clipping, occlusion, and opacity). Zero and multiple matches are valid counts. Results are nonnegative safe integers; malformed selectors, missing/ambiguous roots, invalid responses, and values outside bounds fail closed. `git diff --check` passed. |
| M05.3 Update profile schema, guide, examples and conformance vectors | [x] | Added Browser 1.10 JSON/Markdown, exact integer and CSS/root/visibility constraints, synthetic count vectors and fixture; embedded schema dispatch inherits 1.9 context/template checks. Updated current docs and immutable digests; `go generate ./schemas`, `go test ./schemas -count=1`, `go vet ./schemas`, and `git diff --check` passed. Prior profile digest checks passed. |
| M05.4 Verify and review the published specification | [~] | Full verification and bounded review iteration 1 in progress. Publish exact reviewed source for downstream consumers after acceptance. |
