# C06 History

Milestone: `C06`
Outcome: `completed`
Retired: `2026-09-24`
Source status: `tabilet/memory-bank/status-C06.md`
Source specification: `tabilet/memory-bank/status-C06.md#status-c06---uws-111-specification-amendment`
Evidence: `7c0effb87305b005a38cc4a0899f30626f03aceb`
Worktree: `includes uncommitted changes`
Review: `passed`
Review iterations: `1`
Verification: `go test ./...`; `go test -race ./...`; `go vet ./...`; `mkdocs build --strict`; `git diff --check`; clause and pinned-target review
Consolidated into: [UWS 1.11 specification](../../../versions/1.11.0.md), [release changelog](../../../versions/CHANGELOG.md), and [profile compatibility lesson](../../memory-bank/lessons.md)

## Final milestone specification

````markdown
# Status C06 - UWS 1.11 Specification Amendment

**Goal.** Make the published UWS 1.11 normative text internally consistent
with its existing loop, criterion, validation, and execution contracts.

**Scope.** Amend `versions/1.11.0.md` only: §5.1 permits `$item` and `$index`
in `forEach` and `loop`; §5.2 describes `Criterion.context` as the input and
`Criterion.condition` as the JSONPath query or pointer; normative validation
and content-trust requirements are language-neutral; §11.2 profile links use
the reviewed full commit. Record the amendment in `versions/CHANGELOG.md` and
update only the `1.11.0.md` digest. Do not change schema, wire shape, Go API,
runtime behavior, earlier published artifacts, or the pinned conformance
corpus. Browser 1.9 wording remains a candidate.

**Dependencies and lineage.** Completed [C04](../docs/history/status-C04.md)
published UWS 1.11. C06 has no pending prerequisite or active downstream
consumer.

**Compatibility.** The amendment states already implemented behavior and
existing semantic rules; there is no migration or new UWS version. Rollback
would restore contradictory normative text, so the amended digest and
changelog provide the publication record.

**Acceptance and verification.** Cross-check the corrected clauses against
§§4.5.6.3, 4.5.11, and 5.7 and the current implementation. Confirm only the
`1.11.0.md` protected digest changes. Run `mkdocs build --strict`,
`go test ./...`, `go test -race ./...`, `go vet ./...`, and `git diff --check`,
then complete the bounded whole-milestone review gate in
[milestone.md](milestone.md).

**Review provenance.** Source: “Third-review intake: one core remediation
milestone”, findings R1/R2/R5/R11; source priority P3 for each, local severity
P2 for R1/R2 and P3 for R5/R11. Review baseline and revalidation HEAD:
`7c0effb87305b005a38cc4a0899f30626f03aceb`; the worktree was clean.
R1/R2 evidence: `versions/1.11.0.md` §§5.1–5.2 conflict with §§4.5.6.3,
4.5.11, and 5.7 and with `uws1/execution_criteria.go`. R5 evidence:
`versions/1.11.0.md` §§4.6–4.7 name Go APIs in normative requirements. R11
evidence: §11.2 links to moving `blob/main` profile documents. These are
confirmed; C06 repairs completed C04 without reopening its historical record.
This intake is not review-gate iteration 1.

| State | Task | Evidence |
|---|---|---|
| `[+]` | Amend the four 1.11 wording issues, changelog, and protected digest; verify and review C06. | `versions/1.11.0.md`, `versions/CHANGELOG.md`, and the one protected digest updated. `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed. |

**Review gate.** Passed at iteration 1. The whole-milestone review checked
the amended §§5.1–5.2 against §§4.5.6.3, 4.5.11, and 5.7 and the current
criterion implementation; confirmed the former Go API names and moving
profile links are absent, each pinned target exists at the reviewed commit,
and the `1.11.0.md` SHA-256 matches the sole changed protected digest. No
P1/P2 issue remains. The documented acceptance commands passed.
````

## Final status

````markdown
# Status C06 - UWS 1.11 Specification Amendment

**Goal.** Make the published UWS 1.11 normative text internally consistent
with its existing loop, criterion, validation, and execution contracts.

**Scope.** Amend `versions/1.11.0.md` only: §5.1 permits `$item` and `$index`
in `forEach` and `loop`; §5.2 describes `Criterion.context` as the input and
`Criterion.condition` as the JSONPath query or pointer; normative validation
and content-trust requirements are language-neutral; §11.2 profile links use
the reviewed full commit. Record the amendment in `versions/CHANGELOG.md` and
update only the `1.11.0.md` digest. Do not change schema, wire shape, Go API,
runtime behavior, earlier published artifacts, or the pinned conformance
corpus. Browser 1.9 wording remains a candidate.

**Dependencies and lineage.** Completed [C04](../docs/history/status-C04.md)
published UWS 1.11. C06 has no pending prerequisite or active downstream
consumer.

**Compatibility.** The amendment states already implemented behavior and
existing semantic rules; there is no migration or new UWS version. Rollback
would restore contradictory normative text, so the amended digest and
changelog provide the publication record.

**Acceptance and verification.** Cross-check the corrected clauses against
§§4.5.6.3, 4.5.11, and 5.7 and the current implementation. Confirm only the
`1.11.0.md` protected digest changes. Run `mkdocs build --strict`,
`go test ./...`, `go test -race ./...`, `go vet ./...`, and `git diff --check`,
then complete the bounded whole-milestone review gate in
[milestone.md](milestone.md).

**Review provenance.** Source: “Third-review intake: one core remediation
milestone”, findings R1/R2/R5/R11; source priority P3 for each, local severity
P2 for R1/R2 and P3 for R5/R11. Review baseline and revalidation HEAD:
`7c0effb87305b005a38cc4a0899f30626f03aceb`; the worktree was clean.
R1/R2 evidence: `versions/1.11.0.md` §§5.1–5.2 conflict with §§4.5.6.3,
4.5.11, and 5.7 and with `uws1/execution_criteria.go`. R5 evidence:
`versions/1.11.0.md` §§4.6–4.7 name Go APIs in normative requirements. R11
evidence: §11.2 links to moving `blob/main` profile documents. These are
confirmed; C06 repairs completed C04 without reopening its historical record.
This intake is not review-gate iteration 1.

| State | Task | Evidence |
|---|---|---|
| `[+]` | Amend the four 1.11 wording issues, changelog, and protected digest; verify and review C06. | `versions/1.11.0.md`, `versions/CHANGELOG.md`, and the one protected digest updated. `go test ./...`, `go test -race ./...`, `go vet ./...`, `mkdocs build --strict`, and `git diff --check` passed. |

**Review gate.** Passed at iteration 1. The whole-milestone review checked
the amended §§5.1–5.2 against §§4.5.6.3, 4.5.11, and 5.7 and the current
criterion implementation; confirmed the former Go API names and moving
profile links are absent, each pinned target exists at the reviewed commit,
and the `1.11.0.md` SHA-256 matches the sole changed protected digest. No
P1/P2 issue remains. The documented acceptance commands passed.
````
