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
| Freeze published Markdown artifacts | `[ ]` | Hash the final M01 bytes, enforce exact non-changelog Markdown membership, and update agent guidance plus architecture/technical current truth in the same commit. |
