# Status M04 - Reproducible Memory-Bank Evidence

**Goal.** Make the evidence already cited by the current memory bank
reproducible in a clone, while preserving frozen archive identities and the
user-owned GOAL protocol. This is later execution work, not tracking or
committing those files during reconciliation.

**Dependencies and order.** No product dependency; run after completed
[C04](../docs/history/status-C04.md) in the remaining approved order M04.
Completed historical M01–M03, B02, C05, and C04 remain retired.
**Downstream impacts.** None active.

**Scope and decision.** Audit the five existing, untracked, verified archive
files, `tabilet/GOAL.md`, and `tabilet/evolution/` for baseline integrity,
private content, and link targets. During M04 execution, track those unchanged
approved artifacts as evidence after that audit. Archive IDs and milestone
IDs occupy separate namespaces: `architecture.md` already uses “Archive C01”
and similar editorial display labels, without renaming archive files or
retired records. M04 must verify that distinction while making the evidence
clone-visible. The empty local `ansiblemodulecall/` directory
is outside this milestone and requires no repository action.

**Compatibility and rollback.** Do not alter archive bytes, GOAL protocol,
evolution v1/v2 history, or retired statuses. If an artifact fails the audit,
stop the tracking row and obtain a revised evidence disposition; do not silently
rewrite or drop current citations. Existing links and history remain valid.

**Acceptance and verification.** Confirm exact planned artifact set, hashes,
links, and no private data; track the audited unchanged files in a later task
commit; distinguish archive labels from milestone IDs in current prose. Check
`git diff --check`, `mkdocs build --strict`, history/current-memory links,
`go test ./...`, `go test -race ./...`, `go vet ./...`, and the bounded
whole-milestone review gate. Product tests are regression checks, not evidence
that this planning phase implemented product fixes.

**Review provenance.** Source: “UWS review, pass 2: 1.10.0, browser 1.8 and
the remediation (2026-09-23)” (`uws-review-2.md`), stated baseline `5923cc0`;
revalidated at full HEAD `5923cc0d3b2c852e45aee3771201c642b9daeca9`.
Relevant worktree state: five untracked `tabilet/docs/archive-*.md` files,
untracked `tabilet/GOAL.md`, and untracked `tabilet/evolution/`, all pre-existing
user content; no tracked edits. Lineage: M01–M03, C01–C03 and B01 history
cite archived evidence; none is reopened. This is ordinary intake, not a
review-gate iteration.

**Follow-up plan review.** The pasted “Plan review of second-review
remediation (2026-09-23)” supplied no source priority or baseline; revalidated
at full HEAD `5923cc0d3b2c852e45aee3771201c642b9daeca9` with the
previous reconciliation's uncommitted planning edits and pre-existing
untracked Tabilet evidence present. F10 concerns wording in this pending
status, not delivered evidence tracking. Ordinary intake, not a gate pass.

| Finding | Source priority | Local severity | Current evidence and disposition |
|---|---|---|---|
| N19 | P3 | P2 | Partially confirmed: `tabilet/memory-bank/architecture.md`, `tabilet/memory-bank/lessons.md`, and `tabilet/docs/history/index.md` link or cite untracked Tabilet evidence, which a clone cannot reproduce. Archive/status IDs occupy separate namespaces; architecture display labels were clarified editorially, while tracking remains pending. Empty local directory is not a tracked repository issue. |
| F10 | not supplied | P3 | Confirmed editorial distinction: `tabilet/memory-bank/architecture.md` already labels its archive links “Archive C01” and so on, while the former M04 row still planned to clarify them. The substantive N19 evidence-tracking gap remains. Lineage: N19/M04. |

| Item | State | Notes |
|---|---|---|
| Audit existing Tabilet evidence before tracking | `[x]` | Read-only audit completed 2026-09-24 against HEAD `34267299e33f7fae2cf4ecc3e956a6530806263f`. The exact 14-file set and SHA-256 digests are recorded below. All five archive baselines (`8382d0f26b3b10870125760643078d1a1a3e31b6`) resolve to an existing commit. Targeted scans found no credential/private-key/API-token values, email addresses, or absolute local paths. Link audit found 10 literal relative links in `evolution/result-v2.md` through `result-v4.md` that still point to retired `memory-bank/status-*.md` paths; the corresponding status records exist under `docs/history/`. Per this status's stop rule, evidence tracking is held pending a revised disposition. |
| Track audited evidence and verify links | `[!]` | Blocked: do not stage or commit the 14 artifacts until the user decides how to treat the 10 stale links in the preserved evolution records. Owner: user. Unblock by approving a revised evidence disposition (link-only edits, preserve-and-document the historical links, or exclude the affected evolution files) and updating this status/acceptance accordingly. |

## M04 Evidence Audit (2026-09-24 UTC)

The planned evidence set is exactly one `tabilet/GOAL.md`, five frozen archive
files, and eight `tabilet/evolution/` files. Hashes below were observed
read-only at audit baseline `34267299e33f7fae2cf4ecc3e956a6530806263f`.
All archives declare baseline
`8382d0f26b3b10870125760643078d1a1a3e31b6` (`2026-09-20`, “Publish UWS 1.9.2
content-trust corrections”), which is present in Git. The private-content
scan found no high-confidence secrets, private keys, API tokens, email
addresses, or absolute local filesystem paths; UWS review source filenames in
evolution v3 are provenance references, not embedded review contents.

| Artifact | SHA-256 |
|---|---|
| `tabilet/GOAL.md` | `b312ffc76c1727a2fc5403ae7de9ad5f2392cdfb3608ca6263b467bd020847e3` |
| `tabilet/docs/archive-B01.md` | `330b5b16436052f4e6ed1e300e243e6e3680e28d3a3bf12a0b999daf38670714` |
| `tabilet/docs/archive-C01.md` | `b59e3b2d2ee6e6c8e47066fbd30492ea724624b8cede96ab3276e9b9f102dcaa` |
| `tabilet/docs/archive-M01.md` | `2764ec9a3631e204d341ca4fa11d042466d082a4fccd7b99b06fcfd30eecd43d` |
| `tabilet/docs/archive-T01.md` | `7c00a5b6a65f43c3bc7e00aa64a0d077d52290e8feb4c02829abad75ed61d662` |
| `tabilet/docs/archive-X01.md` | `7159441b5132521696f7a8d7be4f4761c16121760e7445cf413131749d851ad8` |
| `tabilet/evolution/prompt-v1.md` | `8080cdedf72160a2a5cb8d3340ee22bef2473d3a6858ccab0c1d3987900ca61f` |
| `tabilet/evolution/prompt-v2.md` | `bdd2c52f908ffa9d1b30edb357be08158557ea3e6796f7893fe48232ea21708a` |
| `tabilet/evolution/prompt-v3.md` | `d1e09315655adb7fe73fbcbacd899ae55866ce38bfcfaa40af86d194c99b0879` |
| `tabilet/evolution/prompt-v4.md` | `e39e1f9ce692af84d13a2bedc30eab35af3d76faf07b14743c1e22523761676f` |
| `tabilet/evolution/result-v1.md` | `75160e316e707b66ffb621e9a6c326cfb15d05751b67e45db35ac2ba9050671d` |
| `tabilet/evolution/result-v2.md` | `737d66bdb913a5d0386397353a2c41de6b7e8b72ec92004476a4cec3486ce5f4` |
| `tabilet/evolution/result-v3.md` | `41fda40419b31cce9c874753c35814f2477587854d77b99c7fbb0e8723347291` |
| `tabilet/evolution/result-v4.md` | `c9cce7ec59dd56280389e2b9cbedb7faab3c6a83edb11656166eb30e83d6dbce` |

The 10 unresolved literal link targets are in `evolution/result-v2.md`
(M03, C01, C02, B01, C03), `result-v3.md` (B02, C04), and `result-v4.md`
(C05, B02, C04). All corresponding history records exist, but the preserved
links target the removed active-status paths. No artifact was staged or
modified during this audit.
