# M08 — Verified HCL presentation

**Stage:** Kinet STG-11, Phase A. **Owner:** UWS.
**State:** Confirmed Stage 11 execution; M08.1–M08.3 complete, M08.4 in progress; whole review 3/10 passed; source publication pending.
**Source baseline:** `a7688f54c68f5a75c7cc95aa2b31cea98b31af41` (clean at planning).
**Coordinator:** [Stage 11 contract](../../../kinet/docs/stage11.md); the package-local milestone/status owns acceptance.

## Dependencies and handoff

[APItools:M82](../../../apitools/tabilet/docs/history/status-M82.md).
The serial predecessor is a scheduling gate; direct contract and regression impacts are also listed. Every prerequisite must pass its whole review, and required publication must be independently verified before adoption. Record exact accepted/published sources and fixture/build hashes; no Stage 11 acceptance or future pin is claimed yet.

**Downstream:** [Udon:M47](../../../udon/tabilet/memory-bank/status-M47.md), [Udon:M48](../../../udon/tabilet/memory-bank/status-M48.md), [Kinet:M46](../../../kinet/tabilet/memory-bank/status-M46.md), [Kinet:W17](../../../kinet/tabilet/memory-bank/status-W17.md). Reconcile every affected consumer against the accepted prerequisite revision before advancing.

## Tasks

| Item | State | Notes |
|---|---|---|
| M08.1 — Specify the additive 1.13 transition | `[+]` | Add exact 1.13 schema/specification and the verified HCL presentation/deprecation contract; retain 1.12 wire/execution rules and all prior artifact hashes. Synchronize declared-version admission, latest schema/conformance references, archive, discovery/navigation and current facts. Full standalone tests/races/vet, strict docs, protected versions and consumer regressions passed. Separate codec implementation and whole acceptance/publication remain M08.2–M08.4. |
| M08.2 — Add render verify and import APIs | `[+]` | Separate github.com/OpenUdon/uws/hcl implements deterministic byte-based Render, independent lossless Verify and deprecated inert Import. Public model tags preserve typed blocks/labels/extensions; numeric AST token ranges preserve exact lexemes without custom model decoding/float64/JCS. Exact source/codec/view provenance, tampering/evaluation refusal, literal labels/templates and empty/null/key collisions pass standalone nested tests/races/vet. Root legacy APIs/dependencies and operator go.work remain unchanged; M08.3/.4 own full corpus/release qualification. |
| M08.3 — Prove lossless presentation | `[+]` | Self-contained six-source byte-pinned typed corpus, exact number/string/key/container cases and malformed/ambiguous/evaluation/resource/cancellation refusals pass. Independent UseNumber JSON projection and exact numeric tokens are the losslessness oracle; external-package consumer compiles. Nested standalone/copy-only tests, races/vet, strict docs and diff checks pass. |
| M08.4 — Qualify and publish both modules | `[~]` | Verify old APIs and immutable version hashes, schema/code/docs/archive parity, nested-module builds and codec round trips. Publish exact root/codec sources with named authority and record consumer handoffs. |

## Acceptance and verification

A deterministic HCL view can be proved against the approved document. Existing HCL readers remain available; removal and core dependency extraction are deferred, not falsely claimed complete.

go test ./...; go test -race ./...; go vet ./...; schema/conformance and published-version immutability checks; mkdocs build --strict; git diff --check. Run the separate codec module checks once it exists.
Use only disposable roots and fixtures. Preserve published schemas/wires, historic evidence, current runtime capability restrictions and the installed M44 service. Changed v3/package/worker identities require fresh approval; they do not preserve old grants.

## Execution policy

One execution owner, serial execution and task commits under the later confirmed goal. Planning authorizes no code execution, commit, publication or external operation. Source publication requires separately named authority; a status marker or local build is not publication. Consumers must record exact accepted and published prerequisites before adoption. Default checks are offline, credential-free and model-free. No deployment, live ledger migration, real API/model/mail action or registration change.

## Reconciled consumer contract — 2026-10-06

Kinet:W17 consumes verified presentation via M46; legacy packaged HCL remains a distinct artifact.

## Accepted shape producer prerequisite — 2026-10-06

APItools:M82 is accepted after whole review 3 at exact public source
`54583f9b2f452b7cc522360c5aeeff29ca22f96c`, independently observed on the
unchanged authorized origin/main. The configured registry resolves
`v0.0.0-20261006210844-54583f9b2f45` to that full origin hash.
[Contract](../../../apitools/docs/operation-shapes.md) and
[qualification](../../../apitools/docs/m82-qualification.md) define the additive
BuildOperationShapeTable/VerifyOperationShapeTable APIs over accepted UWS C09
`6a267306032edc687a298cefc8bba7019d3ad059`. The eight-family/twelve-operation
fixture is 9,829 bytes, SHA-256
`dc20d2287322a4b20d5d92b1a0d1ec8036f1bec853b642d7df8797ed096c3ba1`.

Consumers must independently reproduce exact source/shape claims, preserve
native selectors/protocols and OR-of-AND security, and keep partial dialect/
wire/presence/auth evidence unknown. Source URLs are sanitized provenance only;
no fetching, credential resolution or execution is supplied. Bounds are 32
sources, 20 MiB each/64 MiB total raw, 10,000 operations, 8 MiB table/aggregate
projected schemas, 256 KiB per schema and 100,000 projection nodes total/10,000
per schema/depth 50. Incremental expansion/serialization checks refuse without
partial tables; hard CPU/RSS/deadline/mount/network controls remain with workers.
The source tooling still carries the public UWS Horizon/HCL closure; no HCL-free
or private runtime claim is made. Current consumer/browser pins remain unchanged
until this milestone's explicit adoption; Retirement closure `fae9982e42d6b16fe7a5ebfd342a016613a62adb` was independently
observed on authorized APItools origin/main, satisfying the publication gate.

## Persisted review

- Review iteration: **3/10**; passed after a full milestone review with no remaining P1/P2 or higher finding.
- Closing-review findings, iteration 1 (both fixed): **P2 M08-R1-F01** — malformed source block containers (for example info as a nonempty array) reach reflection Elem on a non-container and panic, violating value-free refusal. **P2 M08-R1-F02** — explicit YAML tags on mapping/sequence/scalar nodes are discarded by the JSON projection and can expose a misleading verified view, contrary to the documented unsupported-tag refusal. Both reproduced in hcl/review_test.go; persist before correction. No legacy/root contract finding. Required fixes and next full review remain pending.
- Accepted revision: not available.
- Published revision / artifact evidence: not available.
- Verification: root/nested offline full tests/races/vet, schema/model/archive/version immutability, nested module-only copy, nested compatible staticcheck, strict docs, gofmt and diff checks passed; exact public module/source verification remains M08.4.

After all tasks finish, perform the whole-milestone review with persisted iteration/finding state and fix every P1/P2 before acceptance. Resume an interrupted pass at the same counter. Consolidate current facts, reconcile downstream work and retire under this package’s normal procedure.

## Execution evidence — M08.1, 2026-10-06

Started after exact accepted/public APItools M82 source and retirement closure
were independently verified. UWS baseline before this row was
`c91e9f0afd4e42b7741671e5e9eca53f9d472471`. New 1.13 schema/specification bytes
retain all 1.12 wire fields, validation/execution rules and historical feature
floors; only the release identity/description and additive presentation contract
change. The new schema is 57,162 bytes, SHA-256
`ce118674a63cbaceb46bfbe898d85568aa9281175c1361097c044d7eb6270490`;
new specification is 117,145 bytes, SHA-256
`ca731dbf91fc6fe28563554f88fc48433d65174d66623bb3ad5e41091534e089`.
Existing immutable version documents were not edited. Updated 1.13 membership/
digests, declared-version map and embedded archive passed parity/immutability;
the explicitly named 1.12 parity test continues reading the immutable 1.12 file.

Passed with installed Go1.26.6: standalone/offline full `go test ./...`,
`go test -race ./...`, `go vet ./...`, focused schema/model/version suites,
`mkdocs build --strict`, gofmt and `git diff --check`. Current workspace full
Udon passed. OpenUdon's first concurrent full run had a registration-worker
terminal-diagnostic timing failure (`worker_exit` instead of `duplicate_code`);
the isolated fixture passed three fresh repeats and a subsequent full suite
passed without any sibling changes. This failure/recheck is retained rather
than silently counted as an initial pass. All sibling source/pins stay unchanged.
No dependency upgrade, source publication, deployed docs/service, live/provider
action or legacy HCL removal was performed. Evolution decision: no bump;
this implements the approved Stage 11 additive transition. M08.2 must now
implement the separate byte-based Render/Verify/deprecated Import module;
M08.3/.4 still own qualification, whole review and exact publication.

## Execution evidence — M08.2, 2026-10-06

The separate public codec module consumes exact accepted root UWS C09
`v0.0.0-20261006181058-6a267306032e` for public model tags and strict JSON
helpers. Existing HCL/cty/YAML versions are pinned; standalone graph inspection
retains Horizon/HCL compatibility and supplies no private runtime/source parser/
credential/provider dependency. The operator-owned go.work and root go.mod/sum
are unchanged. Returned provenance names exact source bytes, codec contract/full
revision and HCL bytes; workers must independently bind revision to actual code.

Rendering walks byte-decoded values and model field tags, not custom model JSON
decoders. Independent inert AST reconstruction preserves native typed blocks/
labels, extensions/dollar-key escaping, empty/null containers and exact JSON/YAML
numeric lexemes. Functions, traversals, interpolation, arithmetic, comprehensions,
ambiguous keys/fields, unsupported mappings and malformed inputs refuse. Errors
have stable value-free messages. Bounds are 8 MiB source/view, 100,000 work nodes,
depth 100, checked writer growth and cooperative cancellation; no hard CPU/RSS/
deadline claim is made. No view returns unless Verify proves the projection.

Passed: standalone nested full tests/vet/races with GOWORK=off/GOPROXY=off;
initial exact typed/number/key/extension/empty/null/literal-label fixtures,
tampered lexeme/source/codec and unsafe HCL/cancellation negatives; root full
offline regression tests; strict MkDocs; gofmt and git diff --check. Exact declared
module metadata acquisition was under the confirmed dependency-read scope;
checks remain offline afterward. No source publication, legacy API removal,
dependency upgrade, workflow/package mutation or live authority was exercised.
M08.3 must broaden the lossless/malformed/resource conformance corpus before
M08.4's whole review/publication; no milestone acceptance is claimed here.

## Execution evidence — M08.3, 2026-10-06

The nested conformance supplement pins six byte-identical existing root fixtures
without changing their originals: sample, the complete big typed corpus,
pending-only and three retained browser declarations. Copies inside the module
allow registry/downloaded-module tests without a parent checkout. Manifest SHA-256
is `97408edfe8d7d8147bd0950ac88c00e1757c8564c84992a93e0ec70b19fd5fd2`.
The separate JSON oracle uses standard-library UseNumber and exact token strings;
large integers, precise/trailing-zero decimals, exponent case/sign/padding and
signed zero remain distinct. Unicode/control/template-marker strings, dollar/
escape-prefix/key collisions, null/empty/absence and ordered arrays round-trip.
YAML fixtures and explicit lexical cases also pass; no float64/JCS oracle is used.

Malformed/escaped duplicate sources, unsupported YAML forms, ambiguous typed
blocks/attributes/labels/decoded keys, executable HCL, tampered provenance,
oversized/deep inputs and cancellation refuse without partial values or private
input excerpts. The external-package public API compiles and verifies. Passed:
standalone nested full tests/races/vet, full tests from a disposable copy containing
only the nested module, strict MkDocs, gofmt and git diff --check. Production
codec bytes, root code/schema/archive, immutable versions and fixture originals
were unchanged in this task. M08.4 still owns whole review, both-module release
qualification and exact authorized source publication; no acceptance is claimed.

## Closing review — iteration 1, 2026-10-06

P2 M08-R1-F01 reproduced a reflection panic from an unsupported single-block
array source. Render now checks typed block container kinds before using Elem;
wrong arrays/maps/scalars refuse with no partial view. P2 M08-R1-F02 reproduced
tag loss on YAML maps, sequences and scalars; all explicit tags, including
explicit tagged keys, now refuse and implicit container tags are checked.
Targeted regressions reproduced both failures before correction; full nested
tests/races/vet and compatible staticcheck passed after the principal fixes.
The expanded tagged-key regression also passes. Source limits are checked
before Verify hashes bytes. Full iteration 2 is now in progress, not accepted.

## Closing review — iteration 2, 2026-10-06

No additional production-code defect found; both iteration-1 refusal fixes and
full nested tests/races/vet passed, including tagged keys and module-only copy.
**P2 M08-R2-F01**: current release discovery/current-truth surfaces still link
UWS 1.12 as the specification/schema and retain obsolete codec-pending wording,
contradicting 1.13 admission and the implemented M08.2/M08.3 module. Consolidate
current facts and repair latest discovery links before final qualification;
preserve superseded sections literally in the knowledge journal. Historical
1.12 examples/contracts and frozen records remain unchanged. This finding is
persisted before correction; whole review iteration 3 remains required.

## Closing review — iteration 3, 2026-10-06

The complete additive contract, root admission/parity/archive changes, codec
implementation/corpus/refusal fixes, module/build boundaries and consolidated
current release/docs/CI surfaces were reviewed again. No remaining P1/P2 or
higher finding; the review gate passes at 3/10. No lower-priority finding is
carried without an owner. All required local checks pass. Udon/OpenUdon consumer
root regressions passed as recorded in M08.1; no later root behavior change
requires repeating those suites. Evolution decision: no new version; this
implements the approved Stage 11 boundary.

M08.4 remains in progress until qualified root/codec source is committed,
published under the exact named main-only grant, independently resolved/checked,
and consumers are reconciled. Publication commits suppress push workflows with
[skip ci] because the existing docs workflow force-deploys another ref outside
the source-only grant. Local checks are evidence; hosted CI/deployed docs are
not claimed. Exact acceptance/source/closure hashes are recorded after observing
them; local qualification alone is not publication or overall goal completion.
