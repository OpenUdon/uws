# Second-Review Remediation Direction v3

## Origin

The user approved reconciliation of “UWS review, pass 2: 1.10.0, browser 1.8
and the remediation (2026-09-23)” (`uws-review-2.md`), stated baseline
`5923cc0`, against full HEAD
`5923cc0d3b2c852e45aee3771201c642b9daeca9`. The previous direction in
[prompt-v2](prompt-v2.md) led to completed historical M03/C01/C02/B01/C03;
this is a new remediation direction, not a reopening of them.

## Approved Delivery Boundary

Use this repository's separately versioned Browser profile and core UWS
contracts, Go validation/orchestration, conformance fixtures, guides, and
memory-bank evidence. Planning alone implements nothing. Preserve published
historical artifacts, frozen archives, retired records, and the existing GOAL
protocol. Do not add a concrete browser driver or downstream runtime.

## Requested Outcome And Order

Complete B02 -> C04 -> M04. B02 fixes Browser 1.8's overbroad validator and
defines opt-in Browser 1.9 text-sink safety. C04 publishes an opt-in UWS 1.11
contract and corrects direct-version admission, goto/merge behavior, grammar,
executable vectors, and guides. M04 later audits and tracks the existing
Tabilet evidence; it has no product dependency, but runs last. B02 affects
C04's release references.

## Confirmed Contract Choices

- Browser 1.8 allows standalone braces outside template sinks but rejects
  double-brace sequences there, including malformed sequences. Browser 1.9
  defines deterministic literal-brace escape, URL encoding, prompt controls,
  and value-setting `type_text` semantics.
- UWS 1.11 `goto` is terminal transfer and fails on an already-completed
  target. There is no continuation or cached replay.
- UWS 1.11 admits `$response.body.<path>` and complete JSON-number literals
  for numeric `wait` and `batchSize`; literal `items` remains deferred.
- UWS 1.11 merge over `forEach` includes iteration records without counting
  the parent aggregate twice. Historical behavior is version-gated.
- `await` guide claims are corrected now; portable operation reexecution is
  deferred to a separately approved 2.0 direction.
- The frozen 1.10.0 spec is not rewritten for semantic changes. Genuine
  contradictions and omissions are corrected in a new 1.11 release, while
  historically accurate references remain.

## Acceptance Boundary

The linked [milestone index](../memory-bank/milestone.md) and status files own
task rows, verification, compatibility, and review gates. M04 tracks existing
unchanged evidence only after its later audit, not during this reconciliation.
No code, specification, archive, or GOAL edit is authorized by this snapshot.
