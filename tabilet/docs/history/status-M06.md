# Retired milestone M06 — Public Simulation Runtime And Fixtures

Milestone: `M06`
Outcome: `completed`
Retired: `2026-09-27`
Source status: `tabilet/memory-bank/status-M06.md`
Source specification: `tabilet/memory-bank/milestone.md#m06--public-simulation-runtime-and-fixtures`
Evidence: `5a1ccbf6af2add432ba6142112439b6eb8d9cc83`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `1`
Verification: `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check` (all passed on 2026-09-27)
Consolidated into: [product contract](../../memory-bank/product.md#current-contract-surface), [architecture](../../memory-bank/architecture.md), [technical stack](../../memory-bank/tech-stack.md), [mock runtime guide](../../../docs/mock-runtime.md), and [fixture guide](../../../docs/mock-fixtures.md).

## Final milestone specification

````markdown
### M06 — Public Simulation Runtime And Fixtures

Following C07 acceptance, publish fixture format 1.0 for responses keyed by
local operation ID and request digest. Define canonical request encoding, digest
calculation, response envelopes, repeated-call behavior, provenance, and
invalid/unsupported-version handling. Add an importable `mockruntime` package
implementing the existing `uws1.Runtime` interface, using the real orchestrator
with fixture replay, deterministic schema/example responses, expression/item
evaluation, and resolved would-be request records. Caller-supplied resolvers
provide source data; provider parsers remain downstream.

Pure mock execution makes no network calls. An explicitly enabled real-read
adapter supplies response handoff into mock evaluation; only `read` operations
may reach it. Writes and unknown effects remain mocked. Isolate state across
loops, calls, retries, and parallel branches; diagnose unsupported synthesis
and expressions rather than inventing success. Distinguish synthesized,
replayed, and live-read evidence. Fixture export has explicit redaction
controls and no automatic persistence of credentials or private responses.
Status and acceptance: [M06](status-M06.md).
````

## Final status document

````markdown
# Status M06 — Public Simulation Runtime And Fixtures

**State:** Complete. M06.1–M06.3 and the whole-milestone review passed in
iteration 1. C07 was accepted locally at commit
`9d092664a6062563e0414527f997a2475aeab003`.

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
- Read-only downstream recheck on 2026-09-27 found OpenUdon at HEAD
  `8178e7ead454b766cdef4ae09e48d7ca457a9ef8` with uncommitted M87 work.
  M87.1–M87.3 and M87.5 are marked complete; M87.4 awaits a compatible
  published APItools M77 API, and M87.6 depends on M87.4. The M87.1 shared
  `purpose`, recursive object-root `inputs`/`outputs`, and `effect` declarations
  agree with UWS 1.12's pending-step field shape. Its current fixture verifies
  recursive Go JSON round-tripping; the wrapper-level lossless mapping fixture
  remains correctly assigned to M87.6. OpenUdon still pins the published UWS
  1.11 module. UWS C07.2 is accepted locally at
  `9d092664a6062563e0414527f997a2475aeab003`; M87.6 can pin UWS 1.12 only after
  that revision is published. M87.1 is unblocked and no UWS change is needed for
  the downstream contract. No OpenUdon files were changed here.
- APItools remains the owner of source-backed effect discovery and ranking.
  Its HEAD `e3b4b6ec343a18c48fa93a971a69930b993203db` has uncommitted M77 work;
  M77.1–M77.5 are marked complete locally, while M77.6 is blocked awaiting the
  active OpenUdon session's handoff. No compatible published M77 dependency is
  available to OpenUdon M87.4, so M06.1 does not assume a released APItools
  dependency. Udon owns credentialed real-read enforcement; its clean HEAD
  `4266ac99610a6fe39e363c75068e8256bdd821f5` retains an M43 compatibility note
  for UWS 1.11, so any 1.12 adoption belongs in its ledger. Browsertools owns
  snapshot simulation/acquisition and was clean at HEAD
  `2cdd788e2f9ed38536fd48d743200e89eff7cf32`; this is a separate S2b track.
- Kinet was clean at HEAD `1f986b88978acda9c7c196af86fabb818ce96959`. Its W03
  consumer row remains blocked on published OpenUdon command contracts,
  fixtures, and commands; uncommitted OpenUdon work is not a released
  dependency. M06 does not block W03. W04 is the later mock-runtime consumer.
  No sibling worktree is modified by this UWS milestone.

## Tasks

Execute in this order with one active row owner. Each row includes contract,
implementation, tests, and evidenced current-truth corrections for its scope.

| Item | State | Notes |
|---|---|---|
| M06.1 Fixture contract | `[+]` | Published `uws.mock-fixtures.1.0` with exact `(operationId, requestDigest)` keys, RFC 8785 canonical request bytes and `sha256:` digest, source-neutral JSON responses, stable repeated replay, provenance, strict no-downgrade parsing, bounded codec, exact schema lookup/embedding, portable vectors, and an explicit redactor-required recorded-fixture constructor. No automatic capture or disk writes. `go generate ./schemas`, focused fixture/schema tests, immutable-document and archive checks, `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed. Review added an encoding preflight so oversized caller fixtures are rejected before JSON encoding. |
| M06.2 Public mock runtime | `[+]` | Added public `mockruntime.Runtime` over the core orchestrator with exact fixture replay, caller-supplied example responses, and bounded deterministic schema synthesis. Fixture misses fail unless generated fallback is explicitly enabled; unsupported response schemas and expressions return errors. Added source, response, input, output, item, comparison, and request resolution with UWS version gates and explicit nesting/size limits. Added additive `uws1.RuntimeWithResult`; the existing `Runtime` interface is unchanged, and response results reach success criteria and operation/step outputs. Resolved would-be requests are retained only in bounded memory. The package itself makes no network calls or file writes; a caller resolver must use local data for pure simulation. Parallel, `forEach`, nested workflow-call, retry, criteria, output, no-result-runtime-regression, and defensive response-copy tests pass. Full acceptance passed: `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check`. |
| M06.3 Hybrid reads and qualification | `[+]` | Added `HybridRuntime`, requiring UWS 1.12+, explicit `AllowLiveReads`, and a caller delegate. It delegates only declared operation objects with exact `effect: read`; read fixtures are bypassed after opt-in, while `write`, `unknown`, and omitted effects stay on the pure mock selection path. Delegates receive resolved requests and return bounded strict JSON responses through `RuntimeWithResult`; context cancellation and delegate errors propagate. Each leaf request records `live-read` before delegation, including failed/canceled attempts, alongside fixture/example/synthesis evidence. Tests cover live-read output feeding a fixture-backed mocked write, ineligible effects, disabled/missing/typed-nil delegates, old UWS versions, unbound operation copies, invalid responses, cancellation, and delegate errors. The guide documents explicit scope, request/data handling, and downstream ownership. Full acceptance passed: `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check`. No real provider or browser was used. |

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

**Review state:** Passed.
**Review iterations started:** 1 of at most 10.
**Review result:** The full milestone review found no P1, P2, or higher-severity
findings. Review covered fixture decoding/encoding, strict JSON and canonical
request digests, response resolution and synthesis bounds, expression and item
evaluation, request-history limits and privacy handling, hybrid read eligibility
and explicit enablement, response-result propagation through the orchestrator,
parallel/loop/call/retry isolation, compatibility with existing runtimes, and
the public docs and tests. No review fix was required.
**Verification after review:** On implementation HEAD
`5a1ccbf6af2add432ba6142112439b6eb8d9cc83`, `go test ./...`,
`go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and
`git diff --check` all passed on 2026-09-27.

The review gate, full verification, current-truth consolidation, and downstream
reconciliation are complete. The final specification and status are retained
together in the M06 retirement record.
````
