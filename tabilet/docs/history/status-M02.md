# M02 History

Milestone: `M02`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-M02.md`
Source specification: `tabilet/memory-bank/milestone.md#m02---published-markdown-immutability-guard`
Evidence: `917f95bd4888d1c441b24b821e137b7a46a3d598`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `1`
Verification: `go test ./schemas -run TestPublishedVersionDocumentsAreImmutable -count=1`; `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; scope/count checks; `git diff --check`
Consolidated into: [AGENTS.md](../../../AGENTS.md), [immutability test](../../../schemas/version_immutability_test.go), [architecture](../../memory-bank/architecture.md), [technical stack](../../memory-bank/tech-stack.md)

## Final milestone specification

````markdown
## M02 - Published Markdown Immutability Guard

**Goal.** Detect changes, removals, and unregistered additions among published
Markdown documents under `versions/`, while leaving the release changelog
intentionally mutable.

**Scope.** Extend `schemas/version_immutability_test.go` with SHA-256 and exact
membership coverage for every `versions/*.md` file except
`versions/CHANGELOG.md`, using the final M01 bytes. Update `AGENTS.md`,
`architecture.md`, and `tech-stack.md` in the same task to describe the new
current rule. Do not modify the protected Markdown documents merely to satisfy
their hashes.

**Acceptance.** `go test ./schemas -run
TestPublishedVersionDocumentsAreImmutable`, `go test ./...`, `go test -race
./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` pass.
The focused test proves exact non-changelog Markdown membership and SHA-256
bytes. Review confirms that no versioned JSON document or embedded archive
changed.

**Dependencies.** [M01](../docs/history/status-M01.md) completed and accepted.

**Downstream impacts.** none in the active horizon.

**Task breakdown.** Add the Markdown digest/membership guard and reconcile the
current instructions and memory-bank facts in one commit.
````

## Final status

````markdown
# Status M02 - Published Markdown Immutability Guard

**Acceptance.** The focused published-document immutability test, full tests,
race tests, vet, strict MkDocs build, and diff check pass; every
`versions/*.md` file except the mutable changelog has exact SHA-256 and
membership coverage; no protected Markdown, versioned JSON, or embedded archive
is changed merely to satisfy the guard.

**Dependencies.** [M01](../docs/history/status-M01.md) completed and accepted;
the final corrected Markdown bytes are ready to hash.

**Downstream impacts.** none in the active horizon.

**Review gate.** Passed on iteration 1 of 10 on 2026-09-23. The whole-milestone review found no P1, P2, or higher-severity issue.

| Item | State | Notes |
|---|---|---|
| Freeze published Markdown artifacts | `[x]` | Completed 2026-09-23. The shared immutability test now enforces exact names and SHA-256 bytes for all 31 non-changelog `versions/*.md` files while retaining JSON coverage; agent guidance and architecture/technical current truth are synchronized. The focused test, full tests, race tests, vet, strict MkDocs build, scope check, and `git diff --check` passed; no protected Markdown, JSON, or embedded archive changed. |
````
