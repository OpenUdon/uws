# Status M01 - Current Release Documentation Coherence

**Acceptance.** `mkdocs build --strict`, `go test ./...`, and
`git diff --check` pass; registration 1.1/1.2 release and reference surfaces are
coherent; targeted searches find no stale browser 1.9.0 or present-tense
Ansible support claims; no JSON, generated archive, Go API, validation, or
execution behavior changes.

**Dependencies.** none.

**Downstream impacts.** M02 must hash the final accepted Markdown bytes.

**Review gate.** 0 of 10 iterations started; no persisted findings.

| Item | State | Notes |
|---|---|---|
| Reconcile registration 1.2 release surfaces | `[x]` | Completed 2026-09-23. README, agent guidance, docs home/authoring/navigation, and changelog now surface 1.1/1.2 while retaining the 1.0 lookup/call defaults and 1.1 binding-helper default. `mkdocs build --strict`, `go test ./...`, targeted searches, and `git diff --check` passed; no JSON, Go, or generated archive changed. |
| Correct stale published supplement wording | `[x]` | Completed 2026-09-23. Browser 1.7 now uses the time-stable UWS 1.9+ reference; Runtime Supplement 1.0 describes Ansible binding only as historical UWS 1.6 behavior and states that 1.7+ does not support it. `mkdocs build --strict`, `go test ./...`, targeted searches, and `git diff --check` passed; no JSON, Go, or generated archive changed. |
