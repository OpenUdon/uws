# Status M06 — Public Simulation Runtime And Fixtures

**State:** Planned; waiting for C07 acceptance. No implementation or acceptance
verification performed.

**Specification:** [M06](milestone.md#m06--public-simulation-runtime-and-fixtures).

**Authority and provenance:** User approved C07 -> M06 on 2026-09-26.
[S1 provenance](milestone.md#s1-provenance-and-downstream-boundaries) identifies
the inspected baselines and Kinet working-tree source.

**Prerequisite:** C07 accepted, including its whole-milestone review gate.
Reconcile this pending scope against C07's accepted field and validation
contracts before starting M06. **Downstream:** OpenUdon simulation, Udon
hybrid authoring, and Kinet's later W04 adoption, in their own ledgers.

## Tasks

Execute in this order with one active row owner. Each row includes contract,
implementation, tests, and evidenced current-truth corrections for its scope.

| Item | State | Notes |
|---|---|---|
| M06.1 Fixture contract | `[ ]` | Define and publish fixture format 1.0 with responses keyed by local operation ID and request digest. Specify canonical request encoding, digest algorithm, response envelope, repeated-call/replay behavior, provenance, and invalid/unsupported-version handling. Include portable vectors and schema/codec validation; coordinate published artifacts, lookup/distribution, and protected digests. Define explicit redaction controls for export and avoid automatic persistence of credentials or private responses. Resolve format details before freezing publication. |
| M06.2 Public mock runtime | `[ ]` | Add importable `mockruntime` implementing the existing `uws1.Runtime` without breaking it. Use the real orchestrator for fixture replay, deterministic responses from caller-supplied schemas/examples, expression evaluation, item resolution, outputs, and success criteria. Record resolved would-be requests. Specify fixture/example/synthesis selection and missing-fixture behavior; diagnose unsupported synthesis/expressions explicitly. Pure mock execution performs no network calls. Preserve state isolation across loops, workflow calls, retries, and parallel branches. |
| M06.3 Hybrid reads and qualification | `[ ]` | Add an explicitly enabled real-read adapter with response handoff into mock evaluation. Only `read` operations can reach it; writes and unknowns remain mocked. Qualify mixed workflows, live-read values feeding mocked writes, cancellation/error propagation, and per-invocation evidence distinguishing synthesis/replay/live reads. Publish integration examples and downstream handoff guidance; complete full acceptance, current documentation, and bounded review. |

## Contract Boundaries

- Source schemas/examples come through caller-supplied resolvers. Provider
  source parsers, credential selection, and concrete clients remain downstream.
- The public mock uses the existing runtime interface and UWS orchestrator;
  expression/item support must be sufficient for its documented simulations.
  Unsupported expressions produce diagnostics without changing ordinary UWS
  semantic validation's allowance for implementation-specific expressions.
- Would-be request records represent resolved operation requests. Record
  simulated versus delegated behavior explicitly; do not present mock evidence
  as successful live execution. Export redaction must have defined interaction
  with fixture digest matching, so redacted data is not silently mis-replayed.
- Hybrid enablement and classification are separate: declaring `read` alone
  cannot enable a real call. Udon owns real-read enforcement and credentials;
  this milestone provides the adapter contract and verifies routing with test
  delegates. No live provider or installed browser is required for UWS tests.
- OpenUdon owns simulation commands and approval refusal. Browsertools owns
  snapshot simulation. Kinet owns workflow authoring and future W04 adoption.
  No sibling edits or adoption completion are implied by UWS acceptance.

## Acceptance And Verification

- Fixture tests cover deterministic digest vectors, key/order and value
  distinctions, round trips, unknown versions, malformed records, missing
  matches, repeated calls, and documented redaction/export behavior.
- Runtime tests drive real orchestration and success criteria; verify
  synthesized/example/replayed responses, request recording, unsupported
  inputs, loop/call/retry/parallel isolation, and no automatic persistence.
- Instrumented delegates prove pure mock makes zero real calls and hybrid
  sends only explicitly enabled reads. Verify response visibility and mixed
  real-read/mock-write data flow, with distinguishable evidence on each call.
- Run focused fixture/runtime tests and the complete repository checks:

  ```bash
  go test ./...
  go test -race ./...
  go vet ./...
  mkdocs build --strict
  git diff --check
  ```

- Include published-document immutability and embedded-archive checks when
  adding fixture documents. Record actual results in the owning rows.
- Reconcile downstream contracts and document remaining adoption owners.
  Existing candidates stay deferred; this fixture contract does not promote
  the broader interoperability or expression-portability candidates.

## Whole-Milestone Review Gate

Not started; zero iterations consumed. After all rows and acceptance commands
pass, persist iteration 1 before reviewing the complete milestone. Follow the
[review gate](milestone.md#milestone-review-gate): fix every P1/P2 or higher
finding, rerun affected verification, and review the whole milestone again.
The limit is 10 persisted iterations; resume interrupted passes at their saved
number. Acceptance also requires current-truth consolidation and downstream
reconciliation before retirement under the existing policy.
