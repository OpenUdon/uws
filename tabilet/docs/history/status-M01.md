# M01 History

Milestone: `M01`
Outcome: `completed`
Retired: `2026-09-23`
Source status: `tabilet/memory-bank/status-M01.md`
Source specification: `tabilet/memory-bank/milestone.md#m01---current-release-documentation-coherence`
Evidence: `7a8cc0609db434be6b823e2794180fe39cdd6fa7`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `1`
Verification: `mkdocs build --strict`; `go test ./...`; focused registration tests; targeted stale-claim and scope searches; `git diff --check`
Consolidated into: [README.md](../../../README.md), [AGENTS.md](../../../AGENTS.md), [documentation home](../../../docs/index.md), [registration authoring](../../../docs/registration-authoring.md), [architecture](../../memory-bank/architecture.md)

## Final milestone specification

````markdown
## M01 - Current Release Documentation Coherence

**Goal.** Make every current release surface accurately describe browser
registration 1.2 and historical compatibility, and remove two stale statements
from published supplements without changing wire or execution semantics.

**Scope.** Reconcile `README.md`, `AGENTS.md`, `docs/index.md`,
`docs/registration-authoring.md`, `mkdocs.yml`, and `versions/CHANGELOG.md`;
then apply narrow editorial corrections to `versions/browser.1.7.md` and
`versions/runtime.1.0.md`. Update current memory-bank facts in the same row
that invalidates them. Preserve registration 1.0 schema/call lookup defaults,
the 1.1 binding-helper default, every JSON document, Go API behavior, and the
embedded archive.

**Acceptance.** `mkdocs build --strict`, `go test ./...`, and
`git diff --check` pass. Registration authoring is in navigation; current
reference surfaces and release history describe registration 1.1 and 1.2;
searches find neither the browser profile's stale "current 1.9.0" claim nor
present-tense support for an Ansible source binding in Runtime Supplement 1.0.
Review confirms that the milestone changes no JSON document, generated archive,
Go API, validation behavior, or execution behavior.

**Dependencies.** none.

**Downstream impacts.** M02 hashes the corrected final Markdown bytes and must
not start before M01 passes acceptance and review.

**Task breakdown.** First reconcile registration release surfaces and their
current-memory facts. Then correct the two published supplement statements and
their current-memory facts. Each task is one commit.
````

## Final status

````markdown
# Status M01 - Current Release Documentation Coherence

**Acceptance.** `mkdocs build --strict`, `go test ./...`, and
`git diff --check` pass; registration 1.1/1.2 release and reference surfaces are
coherent; targeted searches find no stale browser 1.9.0 or present-tense
Ansible support claims; no JSON, generated archive, Go API, validation, or
execution behavior changes.

**Dependencies.** none.

**Downstream impacts.** M02 must hash the final accepted Markdown bytes.

**Review gate.** Passed on iteration 1 of 10 on 2026-09-23. The whole-milestone review found no P1, P2, or higher-severity issue; one non-blocking Runtime prose sentence was polished for grammatical clarity.

| Item | State | Notes |
|---|---|---|
| Reconcile registration 1.2 release surfaces | `[x]` | Completed 2026-09-23. README, agent guidance, docs home/authoring/navigation, and changelog now surface 1.1/1.2 while retaining the 1.0 lookup/call defaults and 1.1 binding-helper default. `mkdocs build --strict`, `go test ./...`, targeted searches, and `git diff --check` passed; no JSON, Go, or generated archive changed. |
| Correct stale published supplement wording | `[x]` | Completed 2026-09-23. Browser 1.7 now uses the time-stable UWS 1.9+ reference; Runtime Supplement 1.0 describes Ansible binding only as historical UWS 1.6 behavior and states that 1.7+ does not support it. `mkdocs build --strict`, `go test ./...`, targeted searches, and `git diff --check` passed; no JSON, Go, or generated archive changed. |
````
