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
| Audit existing Tabilet evidence before tracking | `[ ]` | Inspect exact existing GOAL, five archives, and evolution files for expected baseline, hashes, private content, and link resolution; record the accepted set without modifying frozen archive bytes or the GOAL protocol. Stop for a revised plan if integrity fails. |
| Track audited evidence and verify links | `[ ]` | Add only the approved unchanged artifacts during later M04 execution, verify that existing “Archive ID” display labels remain distinct from milestone IDs without renaming either, check clone-visible links and history references, run acceptance commands, and complete the bounded review gate. No tracking or commit is authorized by this planning reconciliation. |
