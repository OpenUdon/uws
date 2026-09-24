# Second-Review Remediation State v3

## Revalidation Baseline

The review stated `5923cc0`; current revalidation used full HEAD
`5923cc0d3b2c852e45aee3771201c642b9daeca9`. Existing untracked Tabilet
GOAL, archive, and evolution files were present; they were preserved and were
not treated as a product implementation. Source: “UWS review, pass 2: 1.10.0,
browser 1.8 and the remediation (2026-09-23)” (`uws-review-2.md`).

## Current Product State

UWS 1.10.0 and Browser 1.8 are published. Earlier remediation milestones
M03/C01/C02/B01/C03 are completed and retired. The second review's confirmed
issues remain open: this reconciliation creates planning records only. No
Browser 1.9 or UWS 1.11 artifact has been published by this work.

## Approved Active Gaps

- [B02](../memory-bank/status-B02.md) owns the Browser 1.8 brace-validation
  defect and the new Browser 1.9 text-sink/safe-integer contract.
- [C04](../memory-bank/status-C04.md) owns exact version admission, UWS 1.11
  goto/merge/expression semantics, executable conformance evidence, genuine
  spec corrections, and guide synchronization.
- [M04](../memory-bank/status-M04.md) owns later audit and tracking of already
  cited Tabilet evidence, preserving archives and GOAL unchanged.

## Deferred Directions And Evidence Boundary

Portable `await` operation reexecution, literal `items`, wider expression
redesign, MCP interoperability, content-trust wire formats, and unmeasured
evaluation-cost work remain unnumbered candidates in the
[milestone index](../memory-bank/milestone.md). The order B02 -> C04 -> M04
is approved, but neither planning checks nor this snapshot prove future
application acceptance. Each status requires its own tests and review gate.
