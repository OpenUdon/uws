# M08 — Verified HCL presentation

**Stage:** Kinet STG-11, Phase A. **Owner:** UWS.
**State:** Confirmed Stage 11 execution; M08.1 complete, M08.2–M08.4 pending; whole review 0/10 not started.
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
| M08.2 — Add render verify and import APIs | `[ ]` | Add the separate uws/hcl module with deterministic Render, lossless Verify and a deprecated public Import. Bind provenance to exact source bytes and codec revision; preserve existing typed key/extension mappings. |
| M08.3 — Prove lossless presentation | `[ ]` | Test large integers, precise decimals, exponent notation, required lexeme preservation, strings, keys, extensions and malformed HCL. Never use float64 or JCS equality as the losslessness oracle; refuse a misleading view. |
| M08.4 — Qualify and publish both modules | `[ ]` | Verify old APIs and immutable version hashes, schema/code/docs/archive parity, nested-module builds and codec round trips. Publish exact root/codec sources with named authority and record consumer handoffs. |

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

- Review iteration: **0/10**; not started.
- Closing-review findings: none; the whole-milestone review has not started. Approved intake requirements above remain pending.
- Accepted revision: not available.
- Published revision / artifact evidence: not available.
- Verification: pending implementation; no test result is claimed by this planning record.

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
