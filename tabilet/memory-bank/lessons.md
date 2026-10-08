# Lessons

Read only the topics relevant to the current change. Product terminology,
architecture contracts, commands, and active work remain in their dedicated
memory-bank files.

## State Profile Compatibility As A Minimum Core Version

**Scope.** Separately versioned profiles referenced by a core specification.

**Lesson.** State the minimum compatible core version in a profile document,
not whichever core release is current when that profile is published. Check
those references when a later core version ships.

**Rationale.** Browser 1.9 remains opt-in and compatible with UWS 1.9 and
later, but its prose still calls UWS 1.10.0 current after UWS 1.11.0 shipped.

**Evidence.** `versions/browser.1.9.md` §§1 and 8;
`versions/1.11.0.md`; third-review R6.

## Release Surfaces Must Move With A Versioned Release

**Scope.** Publishing or advancing a core version or separately versioned
profile.

**Lesson.** Reconcile the top-level README, agent instructions, documentation
home/reference navigation, feature guides and examples, release changelog,
and current memory-bank facts in the same release task as the versioned
artifacts. Check examples against the published grammar and actual behavior,
not only against document structure.

**Rationale.** Browser registration 1.2 was fully represented by schemas, Go
helpers, fixtures, and version documents, while several discovery surfaces
continued to identify 1.1 as current. Tests of the executable contract did not
detect that documentation drift. The UWS 1.10 release similarly left feature
guides with stale or contradictory execution examples; passing Go and schema
tests did not detect that mismatch.

**Evidence.** `versions/browser-registration.1.2.*`;
`versions/browser-registration-call.1.2.*`; `README.md`; `AGENTS.md`;
`mkdocs.yml`; `tabilet/docs/archive-B01.md`; `versions/1.10.0.md`;
`docs/04-Triggers-and-Route-Dispatch.md`; `docs/05-Structural-Results.md`;
`docs/06-Success-Criteria-and-Actions.md`; `docs/07-Execution-Model.md`.

## Separate Editorial Correction From Semantic Amendment

**Scope.** Corrections to an already published Markdown contract.

**Lesson.** A wording correction may repair a stale description without a new
contract version only when it leaves wire shape and meaning unchanged. Any
scope or semantic change needs the changelog's amendment treatment or a new
version before immutable bytes are updated.

**Rationale.** Freezing inaccurate prose preserves the error, while silently
changing normative meaning defeats versioned compatibility. The repository's
changelog already draws this boundary explicitly.

**Evidence.** `versions/CHANGELOG.md`; `versions/browser.1.7.md`;
`versions/runtime.1.0.md`; `tabilet/docs/archive-M01.md`.

## Review Declared-Version Semantics Before Freezing Markdown

**Scope.** Cross-version semantic changes in a repository that freezes published
specifications.

**Lesson.** Compare exact schema selection and semantic/execution behavior
against each declared SemVer before finalizing a release or freezing its
Markdown. Gate features by SemVer precedence, require a matching schema for
prerelease and unpublished versions, and record deliberate exceptions
explicitly.

**Rationale.** The current validator applies 1.9.2 reference-step constraints
only to documents declaring 1.9.2 or later; older versions retain their
declared behavior. Published Markdown and schemas are digest-protected, so
compatibility fixes must be version-gated instead of silently editing history.

**Evidence.** `versions/CHANGELOG.md`; `docs/09-Validation.md`;
`validation/version_compatibility_test.go`; `tabilet/docs/history/status-C02.md`.

## Preserve exact values across custom JSON model decoders

An outer UseNumber decoder does not override nested custom UnmarshalJSON
methods. Prepare a lossless source-value projection before constructing a
model snapshot when exact numeric comparison is required. Ordinary legacy
model decoding can remain compatible while new worker paths use json.Number.
C08's large-integer corpus exposed the rounding boundary; explicit projection
made the 18-case shared corpus pass. See docs/expression-reference.md.

## Keep portability field ownership explicit

An opt-in core checker must not scan profile-owned x-* request fields, regex/
JSONPath/XPath queries or function/browser templates. Standard payload fields
remain bindings even when a nested data key happens to start x-*. C08 review
R1-F01/F02 qualified this distinction and the complete structural-result
expression inventory without tightening ordinary validation.

## Preserve indeterminate schema evidence when proving bindings

A reviewed source type does not prove an arbitrary target constraint. Literal
values may be validated, but expression templates require schema containment;
unsupported constraints and partial number/integer overlap stay indeterminate.
Normalize typed JSON-compatible containers with UseNumber before discovering
expressions. C09 review R1/R2 fixtures reproduced nested and typed-body false
rejections and an incorrect disjoint numeric classification; the fixed checks
pass without key/provider loading or value-bearing diagnostics.

Containment must also preserve the original schema dialect and reference
context through children. The Stage 11 intake probe found a nested draft-07
array incompatible/indeterminate while literal validation against its whole
schema was compatible; the specific dependencies example did not reproduce.
Pattern properties and nullable paths are not proof of absence, and a false
leaf cannot be reported compatible. Accepted [M09.3](../docs/history/status-M09.md)
corrects the confirmed cases at b099f6803277ae94c7e9f1da0904a0140b278f20,
whole review 7, without turning the unsupported example into a new claim.

Output paths must also retain containing cardinality and required property/item
predecessor constraints. A native literal witness proves a value exists; failed
samples prove only indeterminacy. False or exhausted finite const/enum schemas
can prove absence. Mixed templates need valid exact constant schemas for empty
arrays rather than prefixItems:[], which fails the 2020-12 metaschema.
M09 reviews 4–7 reproduce the public false-compatibility and false-indeterminacy
cases with bounded fixtures. Accepted source
b099f6803277ae94c7e9f1da0904a0140b278f20 passes independent ordinary root/codec
and consumer proof; [publication](../../docs/m09-publication.md) records it.

## Distinguish record snapshots from invocation metadata

Changing a workflow-scope field does not necessarily refresh its record snapshot.
Native workflow controls consume executeOnce's incoming records, while child
runnables and final outputs refresh records in the body's frame. Test the actual
reference execution for positive and negative scope cases before deriving flow
ownership from a context label. Root generic workflow dependencies share root
records; explicit/nested calls create child frames, and terminal goto uses the
global exact target at root. Step and operation outputs are separate records
even when their IDs or output names match. M09 reviews 3–6 and its bounded
ordinary/executable-valid pure-runtime regressions provide the evidence.

## Ship nested-module conformance fixtures with the module

Go module downloads exclude parent-module fixture directories. Keep byte-pinned
fixture copies inside a nested module and preserve the original paths as lineage.
Prove the tests in a disposable module-only copy before release. For presentation
code, use an independent UseNumber value-tree oracle and exact numeric token
strings; rendered text or floating-point equality cannot prove losslessness.
M08.3's six-source manifest and exact-number/key corpus provide the evidence.

## Enforce parser budgets before calling a recursive dependency

**Applies when** a public inert parser accepts bounded but untrusted bytes.

**Lesson.** A byte cap and a later semantic depth check do not bound the parser
that runs first. Lexical preflight must account for comments, strings and
interpolation; strict JSON prevalidation must not recursively traverse without
a bound. Canonical presentation verification also needs exact rendered bytes:
semantic equality alone accepts comments that claim authority.

**Evidence.** At intake source 0a4597122a7baa7e79e46e79e4ec60dbfffc3720,
hcl/parse.go parsed before its later depth budget and strict JSON traversed
recursively. The added-comment probe passed semantic Verify; NFD failed closed
because cty normalizes it. Accepted [M09.1/M09.2](../docs/history/status-M09.md)
add lexical/iterative depth-100 preflight and exact canonical-byte verification.
Unsupported conditional chains refuse before recursive parsing. Low-stack child
regressions and independent ordinary root/codec suites pass at
b099f6803277ae94c7e9f1da0904a0140b278f20, whole review 7. NFD remains a documented
fail-closed limitation; frozen versions/corpora and numeric-lexeme proof remain.
