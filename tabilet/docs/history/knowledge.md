# Knowledge History

## 2026-09-23 - Versioned Semantics Lesson Updated

**Source heading.** `tabilet/memory-bank/lessons.md#review-declared-version-semantics-before-freezing-markdown`

**Previous wording.**

> **Scope.** Cross-version semantic changes in a repository that freezes published
> specifications.
>
> **Lesson.** Compare schema, validator, and executor behavior for older declared
> versions before finalizing a new release's prose or freezing its Markdown.
> Record deliberate exceptions and compatibility policy explicitly.
>
> **Rationale.** The current validator applies 1.9.2 reference-step constraints
> to 1.9.1 documents, while the 1.9.1 schema accepts the shape. The published
> Markdown is now digest-protected, so a later semantic correction needs a new
> version or a separately justified amendment rather than an unnoticed edit.
>
> **Evidence.** `versions/1.9.1.json`; `uws1/validation_workflow.go`;
> `schemas/version_immutability_test.go`; `tabilet/docs/history/status-M02.md`.

**Reason.** C02 implemented and verified declared-version gating: 1.9.2
reference-step restrictions now apply only from 1.9.2, 1.5 step inputs only
from 1.5, and criterion pointer index behavior follows the declared version.
The previous wording described a known defect as current truth. Exact schema
selection and SemVer prerelease policy are now documented, with all published
core schemas covered by the compatibility corpus.

**Evidence.** `d4b0766`; `846d707`; `947990b`;
`validation/version_compatibility_test.go`.

**Replacement.** [Review Declared-Version Semantics Before Freezing Markdown](../../memory-bank/lessons.md#review-declared-version-semantics-before-freezing-markdown)

## 2026-09-23 - Release-Surface Lesson Expanded To Core Guides

**Source heading.** `tabilet/memory-bank/lessons.md#release-surfaces-must-move-with-a-profile-release`

**Previous wording.**

> **Scope.** Publishing or advancing a separately versioned profile.
>
> **Lesson.** Reconcile the top-level README, agent instructions, documentation
> home/reference navigation, release changelog, and current memory-bank facts in
> the same release task as the profile artifacts.
>
> **Rationale.** Browser registration 1.2 was fully represented by schemas, Go
> helpers, fixtures, and version documents, while several discovery surfaces
> continued to identify 1.1 as current. Tests of the executable contract did not
> detect that documentation drift.
>
> **Evidence.** `versions/browser-registration.1.2.*`;
> `versions/browser-registration-call.1.2.*`; `README.md`; `AGENTS.md`;
> `mkdocs.yml`; `tabilet/docs/archive-B01.md`.

**Reason.** The second review found that a core 1.10 release also left feature
guides and examples stale or inconsistent with published grammar and behavior.
The profile-only scope was too narrow; this expansion records a current
release-review lesson, not a claim that the guides are already fixed.

**Evidence.** `versions/1.10.0.md`;
`docs/04-Triggers-and-Route-Dispatch.md`;
`docs/05-Structural-Results.md`;
`docs/06-Success-Criteria-and-Actions.md`;
`docs/07-Execution-Model.md`; `tabilet/docs/history/status-C04.md`.

**Replacement.** [Release Surfaces Must Move With A Versioned Release](../../memory-bank/lessons.md#release-surfaces-must-move-with-a-versioned-release)

## 2026-09-23 - Current Browser Profile Updated

**Source heading.** `tabilet/memory-bank/product.md#current-contract-surface`

**Previous wording.**

> UWS 1.10.0 is the current core contract. Browser 1.8 is the current browser
> capability profile; browser authentication/call 1.1 is current for sign-in;
> browser registration/call 1.2 is current for reviewed registration
> verification; and registration input 1.0 is the private envelope format.
> Runtime Supplement 1.0 remains the public metadata floor for common non-HTTP
> extension operations. UWS 1.10 adds versioned portable execution semantics and
> expression-addressable-name restrictions while preserving the 1.x wire
> vocabulary; earlier published contracts remain accepted according to their
> version gates. Ansible support is historical UWS 1.6 material only.

**Reason.** B02 published the separately versioned, opt-in Browser 1.9 profile.
Browser 1.8 remains immutable and is still the empty-lookup compatibility
default, but it is no longer the latest published browser capability contract.

**Evidence.** `versions/browser.1.9.json`; `versions/browser.1.9.md`;
`schemas/schema.go`; `schemas/version_immutability_test.go`.

**Replacement.** [Current Contract Surface](../../memory-bank/product.md#current-contract-surface)

## 2026-09-25 - Browser 1.10 Became Current

**Source heading.** `tabilet/memory-bank/product.md#current-contract-surface`

**Previous wording.**

> UWS 1.11.0 is the current core contract. It preserves the UWS 1.x wire model
> and adds version-gated response-body dot-walks, loop-only `$batchIndex`, numeric
> `wait`/`batchSize` literals, terminal root-scoped `goto`, and corrected
> `forEach` merge records. UWS 1.10 execution semantics remain active for 1.10
> and later; earlier declarations retain their versioned behavior. Browser 1.9 is the current opt-in
> browser capability profile; empty profile lookup retains Browser 1.8 as its
> compatibility default. Browser authentication/call 1.1 is current for sign-in;
> browser registration/call 1.2 is current for reviewed registration
> verification; and registration input 1.0 is the private envelope format.
> Runtime Supplement 1.0 remains the public metadata floor for common non-HTTP
> extension operations. The UWS 1.11 specification, exact-version schema, Go
> validator/executor, and executable conformance corpus are coordinated release
> artifacts; earlier published contracts remain accepted according to their
> version gates. Ansible support is historical UWS 1.6 material only.

**Reason.** M05 published Browser 1.10 as the current opt-in profile while
preserving Browser 1.8 as the empty-lookup compatibility default. The new
profile adds count outputs without changing UWS core or earlier profile bytes.

**Evidence.** `versions/browser.1.10.json`; `versions/browser.1.10.md`;
`schemas/schema.go`; `schemas/version_immutability_test.go`;
`80ee9bfb24a688b5e875dadf9ecacdc65398f1ff`.

**Replacement.** [Current Contract Surface](../../memory-bank/product.md#current-contract-surface)


## 2026-09-27 - UWS 1.12 Became Current

**Source heading.** `tabilet/memory-bank/product.md#current-contract-surface`

**Previous wording.**

> UWS 1.11.0 is the current core contract. It preserves the UWS 1.x wire model
> and adds version-gated response-body dot-walks, loop-only `$batchIndex`, numeric
> `wait`/`batchSize` literals, terminal root-scoped `goto`, and corrected
> `forEach` merge records. UWS 1.10 execution semantics remain active for 1.10
> and later; earlier declarations retain their versioned behavior. Browser 1.10
> is the current opt-in browser capability profile and adds typed CSS selector
> match counts; empty profile lookup retains Browser 1.8 as its compatibility
> default. Browser authentication/call 1.1 is current for sign-in;
> browser registration/call 1.2 is current for reviewed registration
> verification; and registration input 1.0 is the private envelope format.
> Runtime Supplement 1.0 remains the public metadata floor for common non-HTTP
> extension operations. The UWS 1.11 specification, exact-version schema, Go
> validator/executor, and executable conformance corpus are coordinated release
> artifacts; earlier published contracts remain accepted according to their
> version gates. Ansible support is historical UWS 1.6 material only.

**Reason.** C07.3 published UWS 1.12.0 with effect classification and pending-step contracts. The 1.11 schema, specification, and conformance corpus remain immutable; current product truth now identifies the published 1.12 contract and its execution boundary.

**Evidence.** `versions/1.12.0.json`; `versions/1.12.0.md`; `versions/CHANGELOG.md`; `testdata/examples/pending-only.1.12.json`; `uws1/validation_version.go`; `schemas/version_immutability_test.go`.

**Replacement.** [Current Contract Surface](../../memory-bank/product.md#current-contract-surface)

## Stage 11 approved target provenance — 2026-10-06

Source: `AGENTS.md` at `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` (clean). The following literal prior boundary remains evidence for existing behavior. The approved Stage 11 proposal adds a scoped target/transition; it does not claim that code has already moved. Replacement target: [Stage 11](../../../../kinet/docs/stage11.md) and the active milestone specifications. No retired record is changed.

````markdown
## Conversion

`convert/` provides JSON, YAML, and HCL helpers.

Key invariants:

- JSON and YAML preserve `x-*` extensions through the `Extensions` map pattern.
- HCL preserves object-level `x-*` extensions through `extensions { ... }` blocks. JSON and YAML keep extensions flattened as normal `x-*` fields.
- HCL key rewriting preserves `$`-prefixed keys on round-trip. Legacy JSON Schema keys (`$ref`, `$id`, `$schema`, `$defs`, etc.) use the `_`-prefix form in HCL; other `$foo` keys use `__dollar__foo`.
- `MarshalHCL` works on a deep copy and does not mutate the caller's document.
- The `uws1.Document` wire tree and all fields/sub-structs reachable from it should carry `json` and `hcl` tags for parsing.
````

## Stage 11 candidate promotion — 2026-10-06

Literal prior rows from `tabilet/memory-bank/milestone.md` at `a7688f54c68f5a75c7cc95aa2b31cea98b31af41`; user-approved promotion now belongs to the active Stage 11 owners, not a duplicate candidate.

````markdown
| Optional expression portability tooling | Third-review R4: core validation intentionally allows implementation-specific expressions under §5.5; strict grammar checking is not ordinary semantic validation. | A consumer requests an opt-in portability check with a specified interface and compatibility tests. |
````

Replacement: active [Stage 11 milestones](../../memory-bank/milestone.md#stage-11-cross-package-refactoring). Other candidate scope remains deferred.

## Partial conformance candidate promotion — 2026-10-06

Literal prior candidate from milestone.md at a7688f54c68f5a75c7cc95aa2b31cea98b31af41:

````markdown
| Conformance corpus supplement | Third-review R9: Go tests cover behaviors absent from the pinned 1.11 corpus. Changing that frozen corpus would alter published evidence. | Design and approve a separately pinned supplement or a later release corpus with interoperable vectors. |
````

C08/C09 now own the expression/binding subset; the active candidate retains the other corpus gaps. The frozen 1.11 corpus is unchanged. See [the current plan](../../memory-bank/milestone.md#stage-11-cross-package-refactoring).

## C08 reference evaluator boundary — 2026-10-06

Source: architecture.md approved Stage 11 target at `d78189c75f3f67d68a0d1ebabfd176c4d5c81c30`. The parser-only boundary now gains the already approved reference evaluator; no new module dependency or ordinary-validator change. Literal prior wording:

```markdown
C08.1 adds the standard-library-only expressions package and immutable parse results. Source resolution, provider I/O and document admission remain outside parsing; numbers retain JSON lexemes. The old mock path remains unchanged until C08.3.
```

Replacement: architecture.md current C08 parser/evaluator boundary and docs/expression-reference.md. The evaluator consumes UWS-owned in-memory context, performs no provider I/O and requires losslessly prepared caller values.

## C08 mock reference adoption — 2026-10-06

Source: architecture.md at `b40d2226024809cd65859330405ec76c60f97db5`. Literal prior transition wording:

```markdown
The old mock path remains unchanged until C08.3.
```

C08.3 now delegates mock evaluation to the shared reference implementation while preserving historical generic-number/encoded-root compatibility. CheckPortability is opt-in, source-field-only and does not alter ordinary validation. Replacement: current architecture.md and docs/expression-reference.md.

## C08 planning-to-current foundation consolidation — 2026-10-06

Source: product/architecture/tech-stack Stage 11 target clauses at `0411eea6fc84fbd6aa97cef94f53f301260f4844`. Literal prior planning-only wording:

```markdown
Both phases belong to one stage. Current facts below remain the observed implementation; no new acceptance, publication or installed behavior is claimed.
```

Replacement: those current documents now record accepted/published C08 while C09/M08 remain pending. The planned full-stage target and installed Kinet M44 boundary are unchanged.

## M08.1 current core contract advances to 1.13 — 2026-10-06

Source: product.md, Current Contract Surface, at observed c91e9f0afd4e42b7741671e5e9eca53f9d472471 (planning/progress context; full baseline is recorded in M08 status). UWS 1.13 retains these wire/execution rules and adds presentation/deprecation. Replacement: [product current contract](../../memory-bank/product.md#current-contract-surface) and [UWS 1.13](../../../versions/1.13.0.md).

````markdown
UWS 1.12.0 is the current core contract. It adds optional operation `effect`
classification (`read`, `write`, or `unknown`; omission is `unknown`) and
non-executable pending-step contracts with recursive object-root input and
output schemas. Pending-only documents may use an empty `operations` array
without placeholder operations. Structural and semantic validation accepts a
valid pending contract; executable validation rejects any pending step before
runtime methods are invoked, including in an unselected branch. Effect labels
are descriptive and do not grant execution authorization. UWS 1.11 introduced
version-gated response-body dot-walks, loop-only `$batchIndex`, numeric
`wait`/`batchSize` literals, terminal root-scoped `goto`, and corrected
`forEach` merge records. UWS 1.10 execution semantics remain active for 1.10
and later; earlier declarations retain their versioned behavior. Browser 1.10
is the current opt-in browser capability profile and adds typed CSS selector
match counts; empty profile lookup retains Browser 1.8 as its compatibility
default. Browser authentication/call 1.1 is current for sign-in; browser
registration/call 1.2 is current for reviewed registration verification; and
registration input 1.0 is the private envelope format. Runtime Supplement 1.0
remains the public metadata floor for common non-HTTP extension operations.
Mock Fixture Format 1.0 is a separate, inert response-fixture format keyed by a
UWS-local operation ID and SHA-256 over RFC 8785 canonical bytes for the
resolved request-binding object. Repeated exact-key lookups reuse the same
response. Its codec requires explicit redaction for recorded fixtures and has
no automatic response persistence. The public `mockruntime` package uses the
core orchestrator, exact fixture replay, caller-supplied examples or bounded
schema synthesis, expression evaluation, and bounded in-memory would-be
request records; this pure runtime makes no transport calls or automatic file
writes. Its separately enabled `HybridRuntime` delegates only explicitly
declared UWS 1.12+ `read` operations to a caller-provided adapter; writes and
unknown effects remain mocked.
The UWS 1.12 specification, exact-version schema, Go model and validator,
executor, conversion helpers, embedded archive, and protected digests are
coordinated release artifacts. The pinned 1.11 executable conformance corpus
remains unchanged; earlier published contracts remain accepted according to
their version gates. Ansible support is historical UWS 1.6 material only.

````

## 2026-10-06 — M08 review-2 current-fact consolidation

Reason: replace obsolete staged implementation/pending wording and release
references with current implemented boundaries; preserve exact C08/C09 sources.
Evidence baseline: `f9016e6937f7ef0f390d08cad09977b63292b0b3`; includes uncommitted M08.4 qualification and
review fixes. Replacement: [product](../../memory-bank/product.md),
[architecture](../../memory-bank/architecture.md),
[technical stack](../../memory-bank/tech-stack.md) and
[qualification](../../../docs/m08-qualification.md).

### product.md — Stage 11 progress/foundation sections

````markdown
## M08.2 codec implementation progress

The separate uws/hcl module implements deterministic typed views directly from
JSON/YAML bytes, independent exact-value/lexeme Verify and deprecated explicit
Import. Returned provenance binds raw source and codec/view identities; failed
proofs return no view. No file/network loading, variable/function evaluation,
credentials, workflow validation, approval or execution is supplied. Existing
root HCL APIs remain unchanged. M08.3 qualifies the lossless/refusal corpus; M08.4 still owns release acceptance.

## M08.1 presentation contract progress

UWS 1.13's new specification/schema retain the 1.12 wire/execution rules and
add source/codec-bound, independently verified HCL views. Derived views are
neither package input nor authority; explicit deprecated import preserves the
legacy compatibility boundary. Codec implementation/qualification remains
M08.2–M08.4 work. No existing published artifact or installed behavior changed.


## C08 accepted reference foundation

C08 is accepted/retired after review 2 at `0411eea6fc84fbd6aa97cef94f53f301260f4844`, independently observed on origin/main. The [qualification](../../docs/c08-qualification.md) pins the source and supplement; M08 remains pending.

## C09 accepted binding foundation

C09 is accepted/retired after review 3 at `6a267306032edc687a298cefc8bba7019d3ad059`, independently observed on authorized origin/main. [Qualification](../../docs/c09-qualification.md) records the pinned source-neutral contracts and limitations; M08 remains pending.

## Approved Stage 11 product direction — not implemented

C09.3 adds deterministic advisory flow diagnostics for reachability, unreferenced outputs, unknown effects/pending ordering and work bounds. It executes no condition/effect and never grants authority.

C09.2 adds advisory binding validation with compatible/incompatible/indeterminate outcomes. It checks required/literal/expression/output/security metadata without credentials or effects.

C09.1 adds the public source-neutral binding shape/resolver contract. Metadata completeness, exact identity and security OR/AND semantics are explicit; metadata never grants execution.

C08.1 now provides the public expressions parser for the existing core grammar, including version/field/loop gates and value-free errors. C08.2 adds exact scalar comparison, null propagation and workflow/nearest-iteration lookup with shared vectors. C08.3 now makes the mock consume the shared evaluator and exposes opt-in CheckPortability diagnostics. Historical mock extensions remain explicit compatibility adapters; ordinary validation is unchanged.

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. UWS will own opt-in expression portability, source-neutral binding diagnostics and a verified HCL presentation module. Source parsing and product authority remain outside UWS. Existing published versions and ordinary validation compatibility stay intact; HCL input removal is outside Stage 11.
Both phases belong to one stage. C08 acceptance/publication is recorded above; M08 remains pending and claim no installed Kinet behavior. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

````

### product.md — Current Contract Surface

````markdown
## Current Contract Surface

UWS 1.13.0 is the current core specification. It retains 1.12 wire/execution
semantics and adds the separate verified HCL presentation/deprecation contract.
M08.1 supplies the versioned documents; M08.2/M08.3 supply the codec and
byte-pinned conformance. Whole release qualification remains pending in M08.4. Earlier published versions and ordinary HCL APIs
remain available. UWS 1.12 added optional operation `effect`
classification (`read`, `write`, or `unknown`; omission is `unknown`) and
non-executable pending-step contracts with recursive object-root input and
output schemas. Pending-only documents may use an empty `operations` array
without placeholder operations. Structural and semantic validation accepts a
valid pending contract; executable validation rejects any pending step before
runtime methods are invoked, including in an unselected branch. Effect labels
are descriptive and do not grant execution authorization. UWS 1.11 introduced
version-gated response-body dot-walks, loop-only `$batchIndex`, numeric
`wait`/`batchSize` literals, terminal root-scoped `goto`, and corrected
`forEach` merge records. UWS 1.10 execution semantics remain active for 1.10
and later; earlier declarations retain their versioned behavior. Browser 1.10
is the current opt-in browser capability profile and adds typed CSS selector
match counts; empty profile lookup retains Browser 1.8 as its compatibility
default. Browser authentication/call 1.1 is current for sign-in; browser
registration/call 1.2 is current for reviewed registration verification; and
registration input 1.0 is the private envelope format. Runtime Supplement 1.0
remains the public metadata floor for common non-HTTP extension operations.
Mock Fixture Format 1.0 is a separate, inert response-fixture format keyed by a
UWS-local operation ID and SHA-256 over RFC 8785 canonical bytes for the
resolved request-binding object. Repeated exact-key lookups reuse the same
response. Its codec requires explicit redaction for recorded fixtures and has
no automatic response persistence. The public `mockruntime` package uses the
core orchestrator, exact fixture replay, caller-supplied examples or bounded
schema synthesis, expression evaluation, and bounded in-memory would-be
request records; this pure runtime makes no transport calls or automatic file
writes. Its separately enabled `HybridRuntime` delegates only explicitly
declared UWS 1.12+ `read` operations to a caller-provided adapter; writes and
unknown effects remain mocked.
The UWS 1.12 specification, exact-version schema, Go model and validator,
executor, conversion helpers, embedded archive, and protected digests are
coordinated release artifacts. The pinned 1.11 executable conformance corpus
remains unchanged; earlier published contracts remain accepted according to
their version gates. Ansible support is historical UWS 1.6 material only.

The MCP tool-call design is an OpenUdon experiment proposal, not an adopted UWS
contract or implementation.

````

### architecture.md — Stage 11 progress/foundation sections

````markdown
## M08.2 codec implementation progress

The separate `hcl/` Go module now supplies byte-based Render/Verify/deprecated
Import. Public model tags drive typed block/label/extension mappings without
legacy custom decoding. Inert HCL AST parsing reconstructs JSON numeric token
ranges exactly and refuses evaluation, unknown/ambiguous mappings or failed
source/codec/view provenance. The root module and legacy HCL APIs are unchanged;
standalone nested checks and initial privacy/refusal fixtures pass. Full corpus
is qualified by M08.3; whole release qualification remains M08.4.

## M08.1 presentation contract progress

UWS 1.13 is distributed as a new exact-version schema/specification, with no
new wire field or execution rule. The separate planned uws/hcl module owns
source-byte Render/Verify/deprecated Import; this row publishes its contract,
not its implementation. Independent exact-value/lexeme verification precedes
view exposure. Existing root HCL APIs, core dependencies and browser input
remain retained. New-version schema/model/archive/navigation are synchronized.


## C08 accepted reference foundation

C08 is accepted/retired after review 2 at `0411eea6fc84fbd6aa97cef94f53f301260f4844`, independently observed on origin/main. The [qualification](../../docs/c08-qualification.md) pins the source and supplement; M08 remains pending.

## C09 accepted binding foundation

C09 is accepted/retired after review 3 at `6a267306032edc687a298cefc8bba7019d3ad059`, independently observed on authorized origin/main. [Qualification](../../docs/c09-qualification.md) records the pinned source-neutral contracts and limitations; M08 remains pending.

## Approved Stage 11 architecture target — not implemented

C09.3 analyzes the UWS structural graph, core output references and declared effects into sorted value-free code/path findings with an explicit truncation flag. Conditions/profile semantics remain opaque; document mutation and effect dispatch are excluded.

C09.2 validates neutral binding projections and reviewed expression schemas, refuses external schema loads, checks exact resolver/source/selector identity and preserves security OR/AND alternatives. Unknown/unproved compatibility remains indeterminate; ordinary UWS validation is untouched.

C09.1 binding types and strict bounded serialization/resolver snapshots are separate from source parsers. Known/unknown schemas/security and protocol-specific fields remain honest; no APItools/provider import or network access is added.

C08 provides the expressions parser plus a reference evaluator over UWS-owned in-memory document/execution snapshots. Parsing uses only the standard library; evaluation reuses uws1 state and the existing internal strict-JSON helper. It performs no source/provider I/O or document admission. Callers must preserve json.Number in value projections before constructing snapshots; ordinary UWS custom decoding remains compatible and may round any-field numbers. The mock now delegates to the shared evaluator, with small historical generic-number/encoded-root adapters. CheckPortability is opt-in, value-free and limited to known core fields; source-function/browser templates and non-simple queries remain profile-owned. Known invocation scopes prevent declaration-site false positives for batch expressions. Ordinary validators and published-version bytes are unchanged.

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. UWS will own opt-in expression portability, source-neutral binding diagnostics and a verified HCL presentation module. Source parsing and product authority remain outside UWS. Existing published versions and ordinary validation compatibility stay intact; HCL input removal is outside Stage 11.
Both phases belong to one stage. C08 acceptance/publication is recorded above; M08 remains pending and claim no installed Kinet behavior. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

````

### tech-stack.md — Stage 11 progress/foundation sections

````markdown
## M08.2 codec implementation progress

Module `github.com/OpenUdon/uws/hcl` lives in `hcl/`, using existing pinned
HashiCorp HCL `v2.24.0`, cty `v1.17.0`, YAML `v3.0.1` and accepted root UWS C09
`v0.0.0-20261006181058-6a267306032e`. No operator go.work/root dependency edits.
Run nested tests/races/vet with GOWORK=off, GOPROXY=off. Bounds are 8 MiB source/
view, 100,000 nodes/depth 100; hard process resource isolation stays external.
The module graph still retains Horizon/HCL compatibility. Publication and exact
final consumer closure remain M08.4 acceptance work.

## M08.1 presentation contract progress

Exact 1.13 schema/specification membership is added while all older published
hashes remain unchanged. The core declared-version map and latest conformance
references select 1.13; prior version gates/wire fields stay intact. The embedded
archive is regenerated with go generate ./schemas. Existing dependency versions
and legacy HCL APIs are unchanged; the nested codec module is implemented and tested separately.


## C08 accepted reference foundation

C08 is accepted/retired after review 2 at `0411eea6fc84fbd6aa97cef94f53f301260f4844`, independently observed on origin/main. The [qualification](../../docs/c08-qualification.md) pins the source and supplement; M08 remains pending.

## C09 accepted binding foundation

C09 is accepted/retired after review 3 at `6a267306032edc687a298cefc8bba7019d3ad059`, independently observed on authorized origin/main. [Qualification](../../docs/c09-qualification.md) records the pinned source-neutral contracts and limitations; M08 remains pending.

## Approved Stage 11 tooling target — not implemented

[Stage 11](../../../kinet/docs/stage11.md) and [local milestones](milestone.md#stage-11-cross-package-refactoring) define the approved target. Exact published module revisions, frozen build closures and standalone verification are acceptance gates. go test ./...; go test -race ./...; go vet ./...; schema/conformance and published-version immutability checks; mkdocs build --strict; git diff --check. Run the separate codec module checks once it exists.
Both phases belong to one stage. C08 acceptance/publication is recorded above; M08 remains pending and claim no installed Kinet behavior. The installed Kinet M44 service remains unchanged, and Stage 12 owns the browser-dependent removal gates.

````
