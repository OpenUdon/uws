# S1 Authoring Contracts And Simulation Direction v6

## Origin

On 2026-09-26 the user requested planning for S1 from Kinet `docs/icot.md` §7
UWS, then approved the complete C07 -> M06 planning update. Effect
classification comes first, followed by pending steps, a public mock runtime
with a real-read hybrid hook, and a versioned fixture format.

UWS was inspected at clean full HEAD
`1d5535ec75d5693a66bcced5bd98f5c4c824fb2a`. The Kinet source was working-tree
evidence based on HEAD `9719a364904743c0d2cb8b3f47cb4e50ed6322e8`, with
uncommitted edits to `docs/icot.md`; its SHA-256 was
`7d2903d0ef70f046e9b711933749791271e77c53622e258b47a9fb4a6c610fe6`.
That commit is not claimed to contain the supplied document's current bytes.

## Approved Direction

Open C07 for ordered effect and pending-step work, then a single combined
UWS 1.12.0 publication. Earlier versions remain immutable. Pending drafts
pass structural and semantic validation but fail executable validation,
including pending-only drafts without a dummy operation.

Open M06 after C07 acceptance for fixture format 1.0 and public `mockruntime`
implementation of the existing `uws1.Runtime`. Use the real orchestrator,
caller-supplied schema/example data, fixture replay, and would-be request
records. Pure mock makes no network calls. Explicit hybrid enablement permits
only classified reads to reach a real-read adapter; write and unknown remain
mocked. Preserve evidence provenance and explicit redaction controls.

## Boundary

This material direction adds portable draft contracts and simulation tooling
while keeping source parsing, concrete clients, credentials, approval policy,
and browser snapshot simulation with their downstream owners. S1 supports
Kinet's later W04 work and does not block W03. Existing candidates remain
deferred. The [milestone index](../memory-bank/milestone.md),
[C07](../memory-bank/status-C07.md), and [M06](../memory-bank/status-M06.md)
own implementation and acceptance. This approval authorizes planning records;
execution requires a separate request.
