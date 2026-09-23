# Status M02 - Published Markdown Immutability Guard

**Acceptance.** The focused published-document immutability test, full tests,
race tests, vet, strict MkDocs build, and diff check pass; every
`versions/*.md` file except the mutable changelog has exact SHA-256 and
membership coverage; no protected Markdown, versioned JSON, or embedded archive
is changed merely to satisfy the guard.

**Dependencies.** [M01](../docs/history/status-M01.md) completed and accepted;
the final corrected Markdown bytes are ready to hash.

**Downstream impacts.** none in the active horizon.

**Review gate.** 0 of 10 iterations started; no persisted findings.

| Item | State | Notes |
|---|---|---|
| Freeze published Markdown artifacts | `[x]` | Completed 2026-09-23. The shared immutability test now enforces exact names and SHA-256 bytes for all 31 non-changelog `versions/*.md` files while retaining JSON coverage; agent guidance and architecture/technical current truth are synchronized. The focused test, full tests, race tests, vet, strict MkDocs build, scope check, and `git diff --check` passed; no protected Markdown, JSON, or embedded archive changed. |
