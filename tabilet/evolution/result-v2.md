# Review Remediation State v2

## Revalidation Baseline

The review stated `8382d0f26b3b10870125760643078d1a1a3e31b6`; current
revalidation used full HEAD `e765fed4ba0a241e9481f99fc324c6a8afb9a8a1`.
Untracked memory-bank initialization artifacts were present and preserved;
they were not treated as implementation evidence. The portable source is
“UWS spec and architecture review: 1.9.2 and earlier (2026-09-23)”.

## Current Product State

UWS 1.9.2 is the current core contract, with the separately versioned browser
profiles and Go implementation. M01 and M02 are completed and retired; the
published Markdown immutability guard is active. No issue in this review has
been fixed by the planning reconciliation.

## Confirmed Delivery Gaps

- [M03](../memory-bank/status-M03.md): HCL key spelling has direct-decode
  ambiguity and rejects some ordinary keys; the reported silent JSON-to-HCL
  corruption example was narrower than stated.
- [C01](../memory-bank/status-C01.md): workflow calls can share an execution
  key, and iteration-key lexical ordering affects merges after ten items.
- [C02](../memory-bank/status-C02.md): schema selection and semantic version
  gates are not coherently tied to the declared version.
- [B01](../memory-bank/status-B01.md): browser template substitution lacks
  context-sensitive encoding and scalar formatting rules.
- [C03](../memory-bank/status-C03.md): several portable execution rules,
  conformance boundaries, security considerations, and examples disagree
  with or are absent from the current core prose.

## Active Horizon And Deferred Directions

The [milestone index](../memory-bank/milestone.md) owns the required order
M03 -> C01 -> C02 -> B01 -> C03 and all candidate promotion triggers. C03
depends on the four preceding milestones. No candidate has an execution ID.

## Planning Evidence

The current tree was revalidated against the review, relevant schemas,
validator/executor/converter paths, versioned prose, history, and Git state.
Planning-file structural checks can confirm the recorded graph; they cannot
establish application acceptance before implementation.
