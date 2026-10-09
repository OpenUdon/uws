# C10 — Browser-profile shape activation

**Stage:** Kinet STG-12, Phase A. **Owner:** UWS.
**State:** Candidate pre-review corrections in progress in retained lease
`goal/C10`; required publication remains unperformed. Closing review 0/10.
**Source baseline:** `989e3f2c88cac5c0f5a2911dfe04c36a61e43126` (clean at planning).
**Coordinator:** [Stage 12 contract](../../../kinet/docs/stage12.md). This
package-local milestone and status own acceptance. Planning was approved on
2026-10-09 (Kinet R65). It authorizes these planning files only.

**Parallel planning provenance.** Stage 12 parallel execution review,
2026-10-09; source priority and separate review baseline not supplied. Current
revalidation `989e3f2c88cac5c0f5a2911dfe04c36a61e43126`, including uncommitted planning files. F01 (confirmed local P2) adopts scoped workflow ownership; F02 (confirmed Lower) replaces strict ordering with readiness; F04 (confirmed local P2) requires frozen consumer checks.
Evidence: the old Kinet goal/launch rules, owning agent rules, existing pending
dependencies and sibling consumer checks; Udon `go.mod`/`pkg/execute` import no
OpenUdon code. This approved intake changes no row state or review counter.

## Dispatch and lease boundaries

**Depends on.** Kinet:M51. All required prerequisites must have
accepted closure at exact revisions; sibling producer adoption also needs
independently verified publication. A priority position never supplies authority.

**Downstream impacts.** Browsertools:M33, OpenUdon:P10, Udon:M53, Kinet:M56, Kinet:W20.

**Write set.** The owning `uws/` package's implementation, tests, ordinary
documentation, manifests and qualification outputs only as required by this
milestone's existing scope, plus `tabilet/memory-bank/status-C10.md` in its
assigned worktree. Excludes `AGENTS.md`, `tabilet/GOAL.md`, shared memory-bank
files, other statuses, evolution, stages, history/knowledge, the package audit
database/sidecars, coordination docs and launch input. The coordinator alone applies shared-memory and closure
changes serially; no child writes a sibling repository or user ledger.

**Contracts read.** Immutable exact prerequisite artifacts listed above, the
M51 native-owner-reviewed contract/fixtures when applicable, the assigned
package baseline and frozen shared-memory/consumer snapshots captured at
dispatch. Cross-package checks use read-only exact snapshots or approved
published module inputs, never changing sibling checkouts. Record full source,
artifact and fixture hashes in the later execution brief; contract drift pauses
affected leases for coordinator reconciliation. Existing no-workspace/no-directory
substitution requirements for ordinary published adoption remain in force.

**Parallel-safe.** yes. Eligible only under the explicit Stage 12 lease opt-in, with no dependency path or bidirectional read/write conflict against any running lease.
At most one live milestone per package. All tests use private lease ports,
disposable stores/caches/browser profiles and unique output directories.

## Dependencies and handoff

**Upstream.** [Kinet:M51](../../../kinet/tabilet/docs/history/status-M51.md)
browser consumer contract. Build on accepted C09: `KindBrowser` is reserved in
`binding/types.go:15`, and browser templates are skipped in `binding/flow.go:311`
and `expressions/portability.go:43`.

**Downstream.**

- [Browsertools:M33](../../../browsertools/tabilet/memory-bank/status-M33.md)
- [OpenUdon:P10](../../../openudon/tabilet/memory-bank/status-P10.md)
- [Udon:M53](../../../udon/tabilet/memory-bank/status-M53.md)
- [Kinet:W20](../../../kinet/tabilet/memory-bank/status-W20.md)

Consumers adopt only the exact accepted and independently published revision.

## Tasks

| Item | State | Notes |
|---|---|---|
| C10.1 — Activate the browser-profile shape kind | `[x]` | Additive exact browser metadata, native profile fragment checks, closed lossless decoding, no HTTP fields/aliases, immutable resolver snapshots and unchanged non-browser table bytes. Focused tests/race/vet passed offline with Go1.26.6. |
| C10.2 — Browser binding and flow rules | `[x]` | Browser body templates use core binding/reference grammar and actual invocation scopes. Stable body-location/required-symbolic-credential diagnostics; partial/unknown evidence remains indeterminate. Native placeholders, extensions and private declarations stay opaque. Focused tests/race/vet passed. |
| C10.3 — Conformance vectors | `[x]` | Eleven native action/authentication/registration profile vectors, frozen complete M51 corpus/golden, exact native integer literal/default/symbolic guards, full conditional input schemas, tampered/ambiguous/partial evidence and stable value-free refusals. Focused tests/race/vet passed offline. |
| C10.4 — Qualify and publish | `[~]` | Apply the two persisted candidate findings and repeat affected local qualification. Required publication remains unperformed and separately gated. Retain lease; no accepted closure/adoption. |

## Acceptance and verification

**Acceptance.**

- Browser shapes validate and bind with stable codes and deterministic fixtures.
- Shape production stays outside UWS.
- No published profile, schema or `uws.shape-table.v1` byte changes.
- No diagnostic grants execution.

**Verification.**

- `go test ./...`, `go test -race ./...` and `go vet ./...`.
- `GOWORK=off go test ./...`, also run from `hcl/`.
- Schema, conformance and published-version immutability checks.
- `mkdocs build --strict` and `git diff --check`.
- Use only disposable roots and fixtures.

## Execution policy

One coordinator owns the integrated ledgers, shared memory and serialized
integration/closure. Serial execution remains the default. Concurrent leases
require this milestone's declared safety, frozen inputs, a complete explicit
Kinet goal request and the Stage 12 agent-rule opt-in. Each lease has one
in-progress row and one assigned milestone; its persisted review count survives
resume/rebase. Commit policy comes from that later request. Source publication,
deployment and live operations retain separate named authority. Planning and
status markers grant none; audit stays disabled.

## Review

Whole-milestone review: 0/10, not started.

## Candidate pre-review findings — 2026-10-09

Coordinator-approved correction request resumes the same local-only goal and
assignment with effective AUTHORIZATION_GRANTS empty. Independent read-only
candidate inspection by `/root/m16_candidate_review` of exact checkpoint
`23d445227ac09eddbed6c92bc2ba4b27f589ca99` found 0P1/2P2/0P3. No reviewer tests
were run. This candidate inspection is not the closing review; required
publication remains unperformed and the formal counter stays0. Findings are
persisted before fixes; the old checkpoint and qualification evidence are retained.

- **C10-P2-1 — fixed and focused verification passed.** Complete registration input
  slots currently pass generic schema validation with `{type:string}` despite
  native1.1/1.2 requiring type, label and exactly one required/requiredWhen branch
  with closed fields. Validate complete raw declarations against the exact
  embedded input-slot fragment, align Required/Condition and preserve all bytes.
  Keep partial evidence explicitly indeterminate. Add missing-label/requiredness,
  unknown-field and contradictory-declaration regressions. Owner: C10 executor.
  Complete declarations now validate against the exact embedded native input-slot
  fragment. Required/Condition alignment and full raw values are retained;
  minimal partial declarations remain indeterminate. Go1.26.6 focused binding
  tests/race/vet and diff checks pass offline. Final packaging remains pending.
- **C10-P2-2 — confirmed; correction in progress.** Whole-body typed integer
  proof checks target properties but ignores open or extra source properties;
  a bounded n inside an open object can incorrectly become Compatible while
  an extra integer remains unsafe. Require closed source objects or supported
  complete proof of all possible extras; otherwise retain Indeterminate. Add
  closed/open/declared-extra controls without schema rewriting. Owner: C10 executor.

## C10.4 local candidate handoff — 2026-10-09

Final qualified implementation is `23d445227ac09eddbed6c92bc2ba4b27f589ca99`.
Captured integration ref remains `refs/heads/main` at
`4429c07bab4616ab46d9a151d50fb26b1370b1fe`; local rebase was up to date.
Root and separate unchanged HCL full tests/races/vet, version immutability,
strict docs and diff checks pass offline with retained Go1.26.6. Final candidate
root ZIP-only tests/races/vet pass. All396 packaged source filenames/content
match the committed candidate. Root/codec module declarations, sums, embedded
schemas and published versions are unchanged.

Provisional private file-proxy version
`v0.0.0-20261009210335-23d445227ac0` has Go archive sum
`h1:BLltWADAzCV3SkugyuZwCW/x6iVOUma5P0e+gHM3l/U=` and GoMod sum
`h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`. Frozen APItools,
Browsertools, OpenUdon, Udon and Kinet author/exec copies build after disposable
module-requirement adaptation, with GOWORKoff and no directory replacement.
The unchanged Kinet host also builds and consumes no UWS module. Existing Udon
Docker module-version replacement is retained. Every frozen input hash matches
its dispatch manifest. Ordinary [qualification](../../docs/c10-qualification.md),
[complete manifest](../../docs/c10-local-candidate-qualification.json) and
[pending publication handoff](../../docs/c10-publication-handoff.md) retain exact
source identities, artifact/verification hashes and seven module closures.
Those documents and this final blocked status are uncommitted C10.4 work.

Cache-isolation audit caught Go adding the two owned provisional versions to
the retained module-cache version list through a borrowed symlink. Only those
introduced list entries were restored; every other entry was preserved. Private
version lists are now regular copies, a subsequent read recheck passes, and no
candidate artifact/source directory remains in the retained cache. No frozen
repository or retained source-module bytes changed. This repaired tooling side
effect supplies no publication proof.

**Required blocker.** Effective AUTHORIZATION_GRANTS remain empty and
EXTERNAL_MUTATIONS remains none. No publication action is authorized or performed.
Owner: coordinator/human for scoped policy reconciliation and exact authority;
then owning assigned executor for publication and independent ordinary-module
resolution. Proposed observed destination is origin
`git@github.com-tabilet:OpenUdon/uws.git`, `refs/heads/main`, fast-forward only.
The local proxy cannot satisfy this gate. C10.4 remains uncompleted and closing
review remains0 until required action/verification are supplied. The branch and
worktree are retained; children do not integrate, close, retire or select a new
milestone. No browser, provider/model/mail, SSH command, deployment, real ledger
or audit operation occurred.

## C10.1 execution evidence — 2026-10-09

`OperationShape.Browser` is optional and leaves non-browser v1 serialization
unchanged. Protocol/selector identities exactly match frozen M51 public-contract
metadata. Browser/authentication/registration metadata uses the immutable native
schemas, bounded strict JSON and RawMessage preservation. Required false values
and present empty optional arrays survive roundtrip. Tests cover exact numeric
lexemes/schema extensions, forged fields, native effects/policy, canonical origins,
selector/HTTP refusal and private resolver snapshot independence.

`GOWORK=off`, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off`:
Go1.26.6 `go test -p 2 ./binding`, `go test -p 2 -race ./binding`,
`go vet -p 2 ./binding` and `git diff --check` passed. Initial compilation under
owned `/tmp/uws-c10-go-tmp` hit its user quota; the completed checks use private
lease resources for GOCACHE/GOTMPDIR under `/dev/shm`. No installed/downloaded
dependency, browser operation, shared cache deletion or publication occurred.

## Authorized isolated execution — 2026-10-09

**C10.4 local qualification correction.** Initial committed candidate
`1b404c59998223d14cd1d4161076fa9c51e9b48d` passed root/codec full tests/races/vet,
strict docs, module-only ZIP tests/races/vet and all seven frozen consumer builds.
Further decoder checks found browser HTTP fields could be supplied explicitly
empty/null and erased by generic optional structs. Browser-only raw-presence
refusal now enforces the frozen absent-field contract without changing generic
parsing. Browser metadata also refuses case aliases, invalid bracket/Unicode
origins and contradictory native conditional declarations. Focused tests/race/vet
and complete root tests passed after correction. Final candidate packaging and
consumer checks are repeated below before handoff. This is local task
qualification, not the whole-milestone closing review; that counter remains0
while required publication is unperformed.

**C10.3 evidence.** The self-contained `docs/examples/browser-shapes/v1` corpus
pins eleven fixture-only native sources (Browser1.5–1.10, authentication1.0/1.1,
registration1.0–1.2), every exact M51 manifest-linked JSON artifact and the
canonical golden. External-package consumers validate native sources, exact
source/selector resolution, full parameter/output/extraction metadata and native
conditional slot schema/condition preservation. Test-only projections start no
browser and implement no public source producer. Negative cases cover foreign
sources, unsupported ref expansion, malformed/duplicate/wrong-kind/colliding
slots and ambiguous or incomplete evidence.

Native action integer binding checks implicit ranges using exact bounded
rationals without rewriting stored schemas: signed64 in Browser1.8 and safe
integer ±9007199254740991 in Browser1.9/1.10. Proved out-of-range literals/defaults
are incompatible; symbolic integer const/enum/range proof is required and unknown
types/ranges remain indeterminate. Frozen unsafe-integer refusal, wide1.8 values,
numeric lexemes, presence false/zero/empty, brace escapes and extensions are
covered. Complete native registration slot declarations retain label/required/
requiredWhen; only a temporary condition-value validation projection removes
those native annotation keys. Offline Go1.26.6 `go test -p 2 ./binding
./expressions ./schemas`, corresponding `-race`, `go vet -p 2` and
`git diff --check` passed using private resources and disabled module networking.

**C10.2 evidence.** Source-bound browser `request.body` participates in binding,
deterministic prior-step output-use observations and strict core portability.
Body references retain actual loop/iteration context. Native profile curly
placeholders, root request/profile extensions, authentication/registration
credential/session/input declarations and browser transport-like request fields
remain opaque. Required native credential slot Name maps to existing consumer
`SecurityBinding.Scheme`; `CredentialSlot` is only the symbolic host slot name.
Missing declarations under partial shapes remain indeterminate. No added browser
wire field, runtime readiness proof or authority. Unit regressions retain
typed maps, absent/mismatched expression proof, body-only locations, declared
unknown effect, deterministic/privacy/no-mutation properties. Offline Go1.26.6
`go test -p 2 ./binding ./expressions`, corresponding `-race`, `go vet -p 2`
and `git diff --check` passed with the isolated resources above.

Coordinator dispatch follows the human-approved local-only Kinet goal with
`STATUS_PRIORITY`, `PARALLELISM: 3`, `INTEGRATION: local-rebase-ff`,
`COMMIT_POLICY: task`, `EXTERNAL_MUTATIONS: none` and effective
`AUTHORIZATION_GRANTS: {}`. Governing Kinet `tabilet/GOAL.md` SHA256
`cad15b1c175094505d380c581f577bde7ced4109ffda40b4b39d5ffe6ae88113`.
This assignment owns only C10 implementation/tests/ordinary docs and this status;
the coordinator owns shared memory, integration, closure and retirement.

Lease `/home/peter/Workspace/uws.goal/C10`, branch `goal/C10`, captured
integration ref `refs/heads/main`, primary `/home/peter/Workspace/uws`.
Original goal base `0034148e9ba26cc33295dc4e9247ee589cdb75f5`; dispatch base
`4429c07bab4616ab46d9a151d50fb26b1370b1fe`. Frozen Kinet accepted input
`ec760d83b6e344e9e9cd7c034f9a92eba99f38b1`; unchanged M51 fixture manifest
SHA256 `0c445a5c90d2c09be561e713c364747f4ab9a3698ea8ab46e7b7b774ccf16bad`.
Immutable consumer inputs and lease-private caches/outputs are under the
coordinator's `/dev/shm/kinet-stage12-leases-f8qot0xc` resource root.
No publication, browser, install, live target or audit authority is supplied.

## Frozen M51 producer input checkpoint — 2026-10-09

Kinet:M51.3 completed at local commit `589b844628b358167fe0a0137eded11c1028875c`. The [native-owner-reviewed contract](../../../kinet/docs/stage12-browser-contract.md), [host ABI](../../../kinet/docs/stage12-browser-host.go.txt), [public declarations](../../../kinet/docs/stage12-browser-public.go.txt) and [fixture manifest](../../../kinet/fixtures/stage12-browser-v1/manifest.json) SHA256 `0c445a5c90d2c09be561e713c364747f4ab9a3698ea8ab46e7b7b774ccf16bad` are exact producer inputs. Both native/public contract reviews pass with zero P1/P2/P3. This is input review, not this producer's implementation, publication or closing review. Kinet:M51 is now accepted locally after all five rows, review1, clean source-built checks and consumed fixture qualification; only that complete retired dependency satisfies its dispatch gate. Rows and persisted review counters remain unchanged.

Use the exact source/extraction/identity maps, per-leaf old driver protocols, complete-plan credential lease, private registration inputs, automatic TOTP versus claimed push continuations, original admitted deadline and separate bounded teardown. The supported consumer profile has one durable session binding per execution and permits other fresh named contexts. M16's immutable host-private save plan preserves v2–v11; candidate creation precedes Join and encrypted host acceptance follows Join/current-generation checks. Preserve report-v5 uncertainty and independently reproduce canonical/golden digests and positive/negative host witnesses. No current artifact claims real conformance or adoption.

**Accepted local M51 prerequisite.** [Retired M51](../../../kinet/tabilet/docs/history/status-M51.md) at observed task/review evidence `33980aafdc614424f61c32a2d79171689c7952bc` (reviewed implementation874c89e, declared coordinator closure) qualifies the unchanged frozen consumer manifest0c445a5c…. Browser qualification1+3 is consumed. This producer remains pending with its original row count/review counter; exact accepted publication and its own browser authority, where required, remain separate.
