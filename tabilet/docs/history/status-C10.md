# Retired Milestone C10

Milestone: C10
Outcome: completed
Retired: 2026-10-09
Source status: tabilet/memory-bank/status-C10.md
Source specification: tabilet/memory-bank/milestone.md#c10--browser-profile-shape-activation
Evidence: e6e537d89f4ef7b00cb352d2f2def3411413aaf0
Worktree: includes uncommitted changes
Review: passed
Review iterations: 1
Verification: Root/unchanged HCL tests, races and vet; immutable profiles; strict docs/diff; independent Git/public module verification; all396 source files/native checksums; public module tests/races/vet; seven frozen consumer builds/closures; integrated binding/expression/conformance tests and full vet.
Consolidated into: [product](../../memory-bank/product.md), [architecture](../../memory-bank/architecture.md), [tech-stack](../../memory-bank/tech-stack.md), [lessons](../../memory-bank/lessons.md); exact public proof in [publication](../../../docs/c10-publication.md).

## Complete milestone specification

``````markdown
### C10 — Browser-profile shape activation

**Stage/owner.** STG-12 Phase A; UWS. Dispatch priority position 3/19; not a review severity.

**Dependencies.** The Kinet:M51 consumer contract and accepted C09.

**Scope.**

- Activate the reserved `browser-profile` ShapeTable kind.
- Add binding, flow and strict-portability rules for browser templates.
- Add Browser 1.5–1.10 conformance vectors.
- Qualify consumers and publish.

**Acceptance.** Browser shapes validate and bind with stable codes. Shape
production stays outside UWS. Published bytes are unchanged, and no diagnostic
grants execution.

**Verification.** `go test ./...`, race and vet checks, `GOWORK=off` module checks
including `hcl/`, conformance and immutability checks, `mkdocs build --strict` and
`git diff --check`.

**Downstream.** Browsertools:M33, OpenUdon:P10, Udon:M53, Kinet:M56 and Kinet:W20.

**Tasks/review.** [status-C10.md](status-C10.md), 4 rows, whole review 1/10 passed.

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
``````

## Complete status

``````markdown
# C10 — Browser-profile shape activation

**Stage:** Kinet STG-12, Phase A. **Owner:** UWS.
**State:** Accepted/retired after all four rows, whole review1, independent
source/module publication and serialized integration verification.
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
| C10.4 — Qualify and publish | `[x]` | Reviewed checkpoint e0ee25b published by normal fast-forward; independent fresh Git fetch and public root module retrieval reproduce exact f01 source and approved Go checksums/all396 files. Public module tests/races/vet and seven frozen consumer builds/closures pass offline. Formal review/accepted closure remain pending. |

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

Whole-milestone review: 1/10, passed; 0 unresolved P1/P2/P3.

## Approved publication resumption — 2026-10-09

The human approved the authorization grants and policy exceptions in Kinet
`suggested.txt` at SHA256
`f7f453d1ca3fe94068475b3f440aea89c5178f054640ac0b4aaa67c18836f75a`
and requested resumption with audit disabled. This records conversation approval
provenance; this status supplies no authority itself. Coordinator-only grants
`stg12-uws-c10-push-main`, `stg12-uws-c10-publication-checkpoint` and
`stg12-uws-c10-published-module-verification` now apply within their exact scopes.
All other external mutations remain forbidden. Historical empty-grant notes
below describe earlier checkpoints.

Implementation remains exactly `f01a2542410c583d0ea909dadd8f17527cc0d27d`.
Fresh independent candidate inspection reports 0P1/0P2/0P3, without dynamic
tests; retained root/ZIP/consumer qualification is complete. The explicit
exception permits one reviewed qualification/evidence/status publication
checkpoint before C10.4 completes. Required publication proof, task completion
and the formal bounded closing gate remain unperformed. No dependent may adopt
until accepted closure and independently verified exact publication.

Read-only remote resolution on resumption found origin
`git@github.com-tabilet:OpenUdon/uws.git`, `refs/heads/main` at
`989e3f2c88cac5c0f5a2911dfe04c36a61e43126`. Captured local integration ref is
still `refs/heads/main`, primary `/home/peter/Workspace/uws`, tip
`4429c07bab4616ab46d9a151d50fb26b1370b1fe`, with a clean worktree. Original
goal integration base remains `0034148e9ba26cc33295dc4e9247ee589cdb75f5`.
Only fast-forward publication is authorized; no force, arbitrary SSH, browser,
deployment or audit operation is included.

## Supported-subset candidate finding — 2026-10-09

**Limited-subset qualified handoff.** Source
`f01a2542410c583d0ea909dadd8f17527cc0d27d`, provisional private proxy version
`v0.0.0-20261009220914-f01a2542410c`, Go archive sum
`h1:TTjAn0++TdGHY0vahbXw3XsICQaavBytoqTSHu0K6o0=`, unchanged GoMod sum
`h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`. Focused/full root and
ZIP-only tests/races/vet and all seven frozen consumer builds pass at bounded
concurrency. All396 ZIP files match committed source; frozen input manifests,
private version lists and retained cache/source boundaries are unchanged.
Generic compatibility code is unchanged. Prior standalone HCL source/dependency
evidence remains valid. Captured refs/heads/main/base remain
`4429c07bab4616ab46d9a151d50fb26b1370b1fe`.
Final qualification/handoff/manifest and this blocked status remain uncommitted
C10.4 work; strict docs/diff checks are repeated after final documentation edits.
No fresh independent pass is claimed. Formal closing review remains0 and source
publication remains unperformed under the unchanged local-only empty grants.

Independent read-only inspection of checkpoint
`ca11bf9a02bac4c36ac5bd8b1dad2137870580dc` reports 0P1/1P2/0P3 and confirms
earlier fixes. The same local-only assignment and empty grants apply; source
publication remains unperformed and formal closing review stays0. Persisted
before code fixes:

- **C10-P2-4 — fixed and focused verification passed.** Native proof trusts const
  under Draft04, where that keyword is ignored. Known `{type:integer,const:0}`
  proof declared Draft04 can falsely prove safety although its compiled schema
  permits9007199254740992. Inherited Draft04 body/property proof has the same
  defect. Define a supported numeric-proof dialect subset and preserve inherited
  context or refuse unsupported trees before interpreting const/enum/range or
  descending. Review every proof path together with reference handling. Keep
  generic compatibility and raw schemas unchanged. Add compiled unsafe direct/
  inherited Draft04 witnesses and supported/default/unknown-dialect controls.
  Owner: C10 executor.
  Whole schema trees now preflight before numeric defaults/literals or typed
  proof. Only the pinned2020 default and canonical Draft06/07/2019-09/2020-12
  declarations support the shared const/enum/range/object subset. Unsupported
  declared dialects or references anywhere in schema-valued positions refuse
  the whole tree before descent, retaining inherited context. Literal schema-like
  data and property names stay data. Source typed arrays remain indeterminate
  without complete tuple/prefix semantics; finite constant arrays remain supported.
  Compiled unsafe Draft04 direct/inherited and modern prefix-array witnesses,
  all supported/default keyword-path controls, unknown/nested dialects and prior
  reference/data controls pass focused tests/race/vet at -p2/GOMAXPROCS2.
  Generic compatibility and raw schemas remain unchanged.

## Final candidate pre-review finding — 2026-10-09

**Final corrected handoff.** Qualified implementation
`ca11bf9a02bac4c36ac5bd8b1dad2137870580dc`, provisional private proxy version
`v0.0.0-20261009215557-ca11bf9a02ba`, Go archive sum
`h1:72aD22ilIAGa93/YDfuGp6WP0kZKaGOb0/Cg5QQQHqc=`, unchanged GoMod sum
`h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`. Full root and ZIP-only
tests/races/vet and all seven frozen consumer builds pass with bounded compiler
concurrency. All396 root ZIP filenames/content match committed source; frozen
inputs and private mutable-cache boundary are unchanged. Generic compatibility
code is unchanged. Prior unchanged standalone HCL evidence remains valid.
Captured integration ref/base remain refs/heads/main and
`4429c07bab4616ab46d9a151d50fb26b1370b1fe`. The final ordinary qualification,
manifest/handoff and this blocked status remain uncommitted C10.4 work.
Strict docs/diff checks are repeated after final edits. No fresh independent
pass is claimed; formal closing review remains0 and publication is unperformed.

Read-only recheck of exact corrected checkpoint
`7b1e936e6cf7fa78c63276e585ed98a796feb8f7` reports 0P1/1P2/0P3; prior two P2
findings are fixed. Coordinator independently confirmed the remaining path;
the reviewer ran no dynamic tests. The same local-only goal/assignment and
empty grants govern this correction. Required publication remains unperformed
and formal closing review stays0. The finding is persisted before code fixes.

- **C10-P2-3 — fixed and focused verification passed.** Numeric proof reads raw
  min/max siblings even when Draft07 ignores them beside a reference. Identical
  target/known proof can therefore appear safe while its effective reference
  selects the unsafe integer9007199254740992. Preserve the generic equal-schema
  compatibility rule. Native numeric proof must retain effective reference,
  dialect/base/root context or return Indeterminate before reading unsupported
  reference siblings. Cover direct/inherited Draft07 and unproved dynamic/recursive
  references, including nested body proof, with negative and safe controls.
  Preserve raw schema bytes. Owner: C10 executor.
  Native numeric analysis now returns Indeterminate for reference-bearing target
  or source schema nodes before reading const/enum/range siblings or descending.
  It interprets no reference without compiled dialect/base/root context. Tests
  exercise effective unsafe direct/inherited/ancestor Draft07 references, isolated
  source proof, unproved dynamic/recursive references, generic equality remaining
  compatible, literal reference-shaped data and safe reference-free controls.
  Focused full binding tests/race/vet and diff checks pass with -p2/GOMAXPROCS2.
  The inactive private Go build cache was cleared by the coordinator between
  handoffs for quota recovery; source/module/fixture/proof inputs were preserved.

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
  tests/race/vet and diff checks pass offline. Final packaging was pending at
  that focused checkpoint; the corrected handoff below records its result.
- **C10-P2-2 — fixed and focused verification passed.** Whole-body typed integer
  proof checks target properties but ignores open or extra source properties;
  a bounded n inside an open object can incorrectly become Compatible while
  an extra integer remains unsafe. Require closed source objects or supported
  complete proof of all possible extras; otherwise retain Indeterminate. Add
  closed/open/declared-extra controls without schema rewriting. Owner: C10 executor.
  Positive whole-body proof now requires an explicitly closed source object,
  no unsupported pattern properties and a native target contract for every
  declared source property. Open/unknown extra inventories remain indeterminate,
  including bounded additional-property schemas not proved by this implementation.
  Closed fully described property controls stay compatible. Focused binding
  tests/race/vet and diff checks pass offline; stored schemas are unchanged.

## Corrected C10.4 candidate handoff — 2026-10-09

Corrected source is `7b1e936e6cf7fa78c63276e585ed98a796feb8f7`, after verified
corrective commits for the two findings. Complete registration declarations now
use exact embedded native1.1/1.2 input-slot fragments; minimal preservation-only
projections cannot assert completeness. Whole-object integer proof cannot ignore
open, patterned or undeclared source properties. Closed completely covered
controls stay compatible and unproved extras remain indeterminate.

Affected root full tests/races/vet and candidate ZIP-only full tests/races/vet
pass. Standalone HCL source/declarations/pinned dependency closure are unchanged,
so their prior successful full checks remain valid. All seven frozen consumer
builds pass on provisional private proxy version
`v0.0.0-20261009213144-7b1e936e6cf7`, Go archive sum
`h1:TJd0T/K36np0eGodBJNjk06qgVshbyTCeErCeT0/V/k=`, unchanged GoMod sum
`h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`. All396 ordinary root ZIP
filenames/content match the corrected committed snapshot. Frozen repository
hashes are unchanged; private mutable version lists remain regular files and
the retained cache contains none of the three provisional candidates.

Captured integration ref/base remain `refs/heads/main` and
`4429c07bab4616ab46d9a151d50fb26b1370b1fe`; local rebase is up to date.
Ordinary qualification/handoff/manifest documents now name the corrected source
and remain uncommitted with this final C10.4 status. Strict docs/diff checks are
repeated after these documentation updates. Independent pre-review did not run
tests; fixes and local qualification do not claim a fresh independent pass.
Closing review remains0 and required publication remains unperformed.

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

**Approved checkpoint inspection.** Independent read-only review of the four
qualification/status files reports 0P1/0P2/0P3; all396 source-file hashes,
10 qualification-log hashes and seven consumer evidence records match.
Strict MkDocs and diff checks pass after the authorization-note edits; docs log
SHA256 `3e7b2260221b16d4607174a68c5a2b081b89a44db77964eda339d8148937ce43`.
The qualified implementation is unchanged. Formal closing review remains0.

## C10.4 independently published qualification — 2026-10-09

[Public proof](../../docs/c10-publication.md) and its complete manifest bind
source f01a2542410c583d0ea909dadd8f17527cc0d27d to checkpoint
e0ee25b72e6d0a6ddfdcc753ef9ad8a1ff7c009a, fresh remote fetch/public Go
retrieval, native checksums, all396 files and seven offline consumer closures.
All required automated checks pass; formal whole review remains0 pending.
The first local retrieval refused equivalent URI encoding before upstream
contact; fresh successful retrieval is independently proved. The first public
race link hit quota; its failed log is retained, two inactive owned derived Go
caches were cleared, and affected checks pass afterward. Module/source/artifact
evidence is preserved. No real ledger, browser, deployment or audit action.

## Authoritative whole-milestone review iteration1 — started

Persisted before review, after all four task rows and required automated
verification/publication proof passed. Full diff is against captured integration
base `4429c07bab4616ab46d9a151d50fb26b1370b1fe`; local rebase is up to date.
Reviewed task/evidence tip `a35e7b23be525030e8c1347e8d2a484fb7b9f4e0`.
Implementation remains exact published f01a2542410c583d0ea909dadd8f17527cc0d27d.
Review covers complete code/contracts/tests/failure/security/compatibility and
local/public qualification evidence. The read-only reviewer receives no
mutation authority. No acceptance, integration or retirement is claimed yet.

**Iteration1 finding C10-R1-P3-1 — persisted before correction.** The current
publication manifest points successful public-module-tests at its former path;
quota recovery moved that log into quota-interrupted-offline-pass/offline-logs.
The exact recorded checksum is retained there. This is an evidence locator
defect, with no missing verification/source defect. Owner: coordinator; correct
the current manifest locator while preserving the historical pass records.

C10-R1-P3-1 corrected: current manifest now resolves the retained successful
log and records its original execution path/relocation reason. Exact log
SHA256 is unchanged. Historical pass records are preserved.

**Iteration1 result — passed.** Independent full review finds no remaining
P1/P2/P3. All prior source fixes, complete native metadata/body/numeric proof,
compatibility/contracts and actual public qualification are reviewed against
4429c07bab4616ab46d9a151d50fb26b1370b1fe. All396 files, approved checksums,
17 successful checks and seven consumer closures match. One P3 locator was
persisted and corrected without changing evidence bytes; 0P1/0P2/1P3 corrected.
The reviewer performed no edits, tests or protected actions. Coordinator
integration, current-truth consolidation, downstream reconciliation and
retirement remain separate closure work.

## Coordinator closure — 2026-10-09

Reviewed task/evidence e6e537d89f4ef7b00cb352d2f2def3411413aaf0 is fast-forward integrated
on captured refs/heads/main. Integrated binding/expression/conformance tests
and full vet pass with the same frozen inputs; whole review1 has no unresolved
findings. Current UWS truth/lessons are consolidated and consumers are reconciled
to exact published implementation f01a2542410c583d0ea909dadd8f17527cc0d27d,
version v0.0.0-20261009220914-f01a2542410c. This observed evidence commit
precedes the included coordinator consolidation/retirement changes; no earlier
commit is claimed to contain them. Final closure is published by the coordinator
under the approved same-source envelope and independently resolved before
downstream dispatch. No consumer pin, source/profile/version bytes, native
producer, driver, runtime, containment, deployment or audit activation changes.
``````
