# Review Remediation Direction v2

## Origin

The user approved reconciliation of “UWS spec and architecture review: 1.9.2
and earlier (2026-09-23)” against current HEAD
`e765fed4ba0a241e9481f99fc324c6a8afb9a8a1`, after the M01/M02
documentation-and-immutability direction was completed. The review stated
baseline `8382d0f26b3b10870125760643078d1a1a3e31b6`.

## Delivery Boundary

The approved planning boundary remains this repository's versioned UWS and
browser contracts, Go model/validator/orchestrator, conversion helpers,
fixtures, and documentation. This evolution snapshot records direction only;
no finding is implemented, no historical version is edited, and no downstream
runtime or external service is changed.

## Requested Outcome

Prioritize HCL key interoperability (M03), workflow execution correctness
(C01), declared-version compatibility (C02), and browser template safety
(B01), then publish a portable execution contract (C03) informed by those
outcomes. Keep MCP, content-trust wire formats, broader expression/trigger
redesign, and other interoperability proposals unnumbered until their stated
promotion triggers are met.

## Confirmed Decisions

- Execute the required order M03 -> C01 -> C02 -> B01 -> C03; only C03 has
  dependencies, on all preceding milestones.
- Gate 1.9.2 semantics by the document's declared version, rather than
  reclassifying them as retroactive errata; define prerelease/schema policy.
- Make non-`await` wait orchestrator-owned bounded numeric seconds while
  `await` remains predicate-based, subject to older-version compatibility
  tests before implementation.
- Version changed browser template semantics separately and preserve browser
  1.5–1.7.
- Choose any future core version under C02's policy, not by assuming a
  text-only 1.9.3.
- Keep review intake planning-only. Later task execution follows the project's
  task commit policy; this reconciliation makes no commit or external mutation.

## Acceptance Boundary

Each milestone's status file owns its test, compatibility, documentation, and
review gate. A completed plan is not an implemented fix. Published history and
frozen archives remain unchanged.
