# M05 History

Milestone: `M05`
Outcome: `completed`
Retired: `2026-09-25`
Source status: `tabilet/memory-bank/status-M05.md`
Source specification: `tabilet/memory-bank/milestone.md#M05`
Evidence: `80ee9bfb24a688b5e875dadf9ecacdc65398f1ff`
Worktree: `clean`
Review: `passed`
Review iterations: `1`
Verification: `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`; focused `go test ./schemas -count=1`; `go vet ./schemas`; published JSON and Markdown immutability verification. The initial strict documentation build found an out-of-tree relative link, fixed in `b3fe55e3f6076333b957d94a90ce419fdc8fb017`, and the strict build then passed.
Consolidated into: [product contract](../../memory-bank/product.md#current-contract-surface), [profile versioning procedure](../../../docs/future-source-profiles.md#adding-a-browser-profile-version), [Browser 1.10 profile](../../../versions/browser.1.10.md), and [README](../../../README.md); no new reusable lesson.

## Final milestone specification

````markdown
## M05 — Browser 1.10 selector match-count profile

Specify and publish a separately versioned Browser 1.10 profile with an
explicit CSS-selector match-count output using `matchCount: true`. It returns
only the exact nonnegative safe-integer count; zero and multiple matches are
valid, and text or attributes are not exposed. An optional `within` CSS
selector scopes descendant matches and MUST resolve to exactly one root;
omission scopes to the selected browser context. `visibility` is required for
count outputs and is either `all` or `rendered`. `all` counts every matching
connected DOM element. `rendered` excludes candidates with no positive-area
client rectangle or with `display: none`, `visibility: hidden|collapse`, or
`content-visibility: hidden` on the candidate or an ancestor. Both modes count
elements below the viewport; clipping, occlusion, and opacity do not change a
count. Missing or ambiguous roots, malformed selectors, invalid responses,
and counts outside the declared validation schema fail closed without exposing
page content. Preserve Browser 1.9 and older profiles unchanged. Keep UWS core
at 1.11 unless a generic binding change is proven necessary. Add schema,
examples and conformance vectors for both visibility modes, scoped roots,
valid zero/one/many counts, and missing, invalid, ambiguous or over-bound
consumer handling. Maintain the Browser profile versioning checklist
alongside the profile docs. Complete focused and full published-source
verification and bounded review; no runtime or real-target browser operation
is included.

Status: [M05](status-M05.md).
````

## Final status

````markdown
# Status M05 — Browser 1.10 selector match-count profile

**State:** Complete. Browser 1.10 was published with reviewed source commit `80ee9bfb24a688b5e875dadf9ecacdc65398f1ff`.

**Goal.** Define and publish Browser 1.10 as a separately versioned profile that
returns a CSS-selector match count as a typed nonnegative integer without
exposing page text or element attributes.

**Compatibility.** Browser 1.9 and older profiles remain immutable. UWS core
stays at 1.11 unless a generic profile-binding change proves necessary. No
browser runtime integration or live target action is included.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| M05.1 Document the Browser profile versioning procedure | [x] | Added the maintenance checklist in [future-source-profiles.md](../../docs/future-source-profiles.md#adding-a-browser-profile-version) and linked it from README; committed in `72d50ab`. |
| M05.2 Define Browser 1.10 count semantics | [x] | Defined `matchCount: true` for exact CSS match counts; optional `within` must resolve to exactly one root and scopes descendant matches; required `visibility` is `all` or `rendered` (rendered excludes hidden/zero-area elements but ignores viewport, clipping, occlusion, and opacity). Zero and multiple matches are valid counts. Results are nonnegative safe integers; malformed selectors, missing/ambiguous roots, invalid responses, and values outside bounds fail closed. `git diff --check` passed. |
| M05.3 Update profile schema, guide, examples and conformance vectors | [x] | Added Browser 1.10 JSON/Markdown, exact integer and CSS/root/visibility constraints, synthetic count vectors and fixture; embedded schema dispatch inherits 1.9 context/template checks. Updated current docs and immutable digests; `go generate ./schemas`, `go test ./schemas -count=1`, `go vet ./schemas`, and `git diff --check` passed. Prior profile digest checks passed. |
| M05.4 Verify and review the published specification | [x] | Bounded review iteration 1 passed with no P1/P2 findings. `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed; focused schema tests and immutability digests also passed. The first strict docs build found an out-of-tree relative link; fixed in a separate commit and rebuilt successfully. Reviewed source was committed as `b3fe55e3f6076333b957d94a90ce419fdc8fb017` and published with the verification record at `80ee9bfb24a688b5e875dadf9ecacdc65398f1ff`. |
````
