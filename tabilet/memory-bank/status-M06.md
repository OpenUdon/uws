# Status M06 — Public Simulation Runtime And Fixtures

**State:** Active; M06.1 is complete and M06.2 is next. C07 was accepted
locally at commit `9d092664a6062563e0414527f997a2475aeab003`.

**Specification:** [M06](milestone.md#m06--public-simulation-runtime-and-fixtures).

**Authority and provenance:** User approved C07 -> M06 on 2026-09-26.
[S1 provenance](milestone.md#s1-provenance-and-downstream-boundaries) identifies
the inspected baselines and Kinet working-tree source.

**Prerequisite:** C07 accepted, including its whole-milestone review gate.
The local C07 release commit and review gate satisfy this prerequisite.
**Downstream:** OpenUdon simulation, Udon hybrid authoring, and Kinet's later
W04 adoption remain in their own ledgers.

## C07 Contract Reconciliation And Downstream Handoff

- UWS 1.12.0 defines optional `Operation.Effect`; omission means `unknown`.
  Core never infers effects from HTTP methods, and an effect label does not
  authorize execution. Hybrid delegation requires separate explicit enablement
  and only explicitly `read` operations may be delegated; `write` and
  `unknown` stay mocked.
- A valid pending step may pass structural and semantic validation, but the
  document is not executable. The mock runtime must enter through the normal
  executable validation path, which rejects pending steps before any runtime
  method is called, including in unselected branches.
- Read-only downstream recheck on 2026-09-27 found OpenUdon at
  `8178e7ead454b766cdef4ae09e48d7ca457a9ef8`, with M87.1–M87.3 complete,
  M87.5 active, and its M87 worktree changes uncommitted. The shared
  `purpose`, recursive object-root `inputs`/`outputs`, and `effect` declarations
  match UWS 1.12's `PendingStep` and `ParamSchema`; the existing fixture tests
  the recursive Go JSON round trip. UWS C07.2 is accepted locally at
  `9d092664a6062563e0414527f997a2475aeab003`, but OpenUdon still pins the
  published UWS 1.11 module. Its M87.6 wrapper-mapping fixture should wait until
  it can pin a published UWS 1.12 revision. No OpenUdon files were changed here.
- Kinet remains the consumer of OpenUdon's released step-command contract for
  W03. Its W03 blocker was recorded before the current OpenUdon M87.2 working
  tree appeared; the implementation files are now untracked locally but are
  not a released dependency. W03 adoption remains in Kinet's ledger. Kinet W04
  is a later consumer of the UWS mock runtime and OpenUdon simulation.
- APItools remains the owner of source-backed effect discovery and ranking;
  its observed HEAD `e3b4b6ec343a18c48fa93a971a69930b993203db` has uncommitted
  M77 work, so M06.1 does not assume a released APItools dependency. Udon owns
  credentialed real-read enforcement; its clean HEAD
  `4266ac99610a6fe39e363c75068e8256bdd821f5` retains an M43 compatibility note
  for UWS 1.11, so any 1.12 adoption belongs in its ledger. Browsertools owns
  snapshot simulation/acquisition and was clean at HEAD
  `2cdd788e2f9ed38536fd48d743200e89eff7cf32`; this is a separate S2b track.
  Kinet was clean at HEAD `1f986b88978acda9c7c196af86fabb818ce96959`; W03
  remains blocked on released OpenUdon artifacts/commands, while W04 is the
  later mock-runtime consumer. No sibling worktree is modified by this UWS
  milestone.

## Tasks

Execute in this order with one active row owner. Each row includes contract,
implementation, tests, and evidenced current-truth corrections for its scope.

| Item | State | Notes |
|---|---|---|
| M06.1 Fixture contract | `[+]` | Published `uws.mock-fixtures.1.0` with exact `(operationId, requestDigest)` keys, RFC 8785 canonical request bytes and `sha256:` digest, source-neutral JSON responses, stable repeated replay, provenance, strict no-downgrade parsing, bounded codec, exact schema lookup/embedding, portable vectors, and an explicit redactor-required recorded-fixture constructor. No automatic capture or disk writes. `go generate ./schemas`, focused fixture/schema tests, immutable-document and archive checks, `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed. Review added an encoding preflight so oversized caller fixtures are rejected before JSON encoding. |
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
