# Initial State v1

## Baseline

The inspected Git baseline is
`8382d0f26b3b10870125760643078d1a1a3e31b6`. The archive preflight recorded
five verified contexts at that commit: core execution, content trust, browser
profiles, source/extension governance, and cross-cutting distribution tooling.
The memory-bank initialization files themselves are uncommitted worktree
content and are not claimed by that commit.

## Current Product State

UWS 1.9.2 is the current core contract. Browser 1.7, browser authentication
1.1, browser registration 1.2, registration input 1.0, and Runtime Supplement
1.0 are the current related contracts. The Go model, validators, content-trust
analyzer, conversion helpers, embedded schemas, profile helpers, fixtures, CI,
and documentation build are present and passing.

## Confirmed Delivery Gap

Browser registration 1.2 exists in schemas, Go support, fixtures, and version
documents, but several top-level discovery surfaces still foreground 1.1.
Registration 1.1/1.2 release history is incomplete, registration authoring is
not in MkDocs navigation, browser 1.7 names 1.9.0 as the current core schema,
and Runtime Supplement prose discusses withdrawn Ansible binding in present
tense. JSON artifacts have digest protection; Markdown artifacts do not.

## Active Horizon

M01 first reconciles current release documentation and applies the two narrow
editorial corrections. M02 then hashes the final published Markdown bytes and
enforces exact membership. MCP standardization, concrete content-trust
resolvers, new browser/profile versions, and UWS 2.0 changes remain unnumbered
candidate directions.

## Initialization Evidence

Before initialization, `go test ./...`, `go test -race ./...`, `go vet ./...`,
`mkdocs build --strict`, `git diff --check`, and archive-structure checks passed.
The strict build reported the known registration-authoring navigation omission
that M01 owns.
