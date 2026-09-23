# AGENTS.md

## Commands

```bash
go test ./...                              # full test suite
go test ./uws1 -run TestSchemaConformance  # run a single test by name
go test ./convert -run TestRoundtrip       # run tests in one package
go vet ./...
```

Module path: `github.com/OpenUdon/uws` (Go 1.25.4).

## Memory Bank Workflow

Before substantial work, read this file, then
`tabilet/memory-bank/product.md`, `tabilet/memory-bank/architecture.md`,
`tabilet/memory-bank/tech-stack.md`, and
`tabilet/memory-bank/milestone.md`. Read the active milestone's linked status
file and only the relevant topics in `tabilet/memory-bank/lessons.md`.

Verified archive files under `tabilet/docs/` are frozen evidence at their
recorded baselines. Current product and system truth lives in `product.md` and
`architecture.md`; read an archive only when its historical baseline is
relevant, and never update it for later changes.

Work one pending status row at a time, keep implementation, tests, and current-
truth corrections in the same row, and follow the milestone dependency order.
One execution owner controls the active ledger across sessions and launchers;
native todos, goal state, or session completion do not replace milestone
acceptance or authorize concurrent ledger writers. Future task execution uses
one commit per completed row unless the governing request says otherwise.

Resolve missing information through safe inspection first. If required files,
resources, commands, permissions, or user answers are unavailable, stop the
affected step and report what is missing while continuing independent in-scope
work. Never invent evidence, bypass permissions, or infer approval from silence
or process exit. A write-gated workflow remains read-only until approval.

Milestone completion requires its persisted review gate and verification, not
only terminal task markers. Review iteration limits persist across sessions and
reviewers; resume an interrupted pass at its recorded iteration before choosing
new work. New engineering reviews enter through the intake procedure in
`tabilet/memory-bank/milestone.md`.

## Architecture

UWS is a workflow overlay for API- and event-source-backed operations. Source documents such as OpenAPI, Google Discovery, AWS Smithy, and AsyncAPI own methods, paths, channels, messages, schemas, servers, and security; UWS owns operation binding, workflow structure, request values, outputs, triggers, and control flow. Non-HTTP runtimes (command exec, SSH, SQL, LLM, browser automation, etc.) are extension-profile concerns expressed via `x-*` fields and `x-uws-operation-profile`, not built-in service types.

The coordinated artifacts must stay in sync:

1. `versions/1.9.2.json` — the latest canonical JSON Schema for UWS 1.x documents. All earlier schemas, including 1.9.1, remain published and immutable. `versions/browser.1.7.json` adds portable scalar accessibility-text conversion while browser 1.5/1.6 remain accepted. `versions/browser-authentication.1.1.json` and its call supplement retain the sign-in context contract. `versions/ansible.1.0.json` is historical UWS 1.6 material; UWS 1.9 does not support Ansible.
2. `uws1/` — the Go model and semantic validator.
3. `versions/1.9.2.md` — the latest human-readable spec. Earlier numbered specifications, `versions/arazzo.md`, `versions/article.md`, and `ideas/terraform.md` are historical or comparison documents. `versions/browser.1.7.md` is the latest browser capability profile and `versions/ansible.1.0.md` is retained only for historical UWS 1.6 documents.
4. `schemas/` — Go lookup and profile-validation helpers plus the generated embedded document archive. `versions/` is document-only; regenerate the archive with `go generate ./schemas` after changing a JSON document.

Browser registration is a separate extension: `versions/browser-registration.1.2.*`
and its 1.2 call supplement retain typed private inputs and input checkpoints
and add a reviewed human-verification contract.
`versions/browser-registration-input.1.0.*` describes the private data envelope;
filled instances never belong in packages or Git. Keep `browserregistration/`,
the profile/input validators in `schemas/`, fixtures and the embedded archive in
sync. Registration 1.0/1.1 documents, 1.0 schema/call lookup defaults, and the
1.1 binding-helper default remain unchanged; newer callers select 1.2
explicitly.

## Execution Model

This repo uses a bound-runtime execution model. Structural orchestration lives in `uws1`; concrete engines bind a runtime implementation at execution time.

### Runtime lives on the base document

The base `Document` carries a `Runtime` interface reference. The runtime provides only leaf execution plus expression/item evaluation.

```go
type Document struct {
    Runtime Runtime
}
```

### Structural execution stays in UWS core

`Document.Execute()` constructs an `Orchestrator`, and the orchestrator walks workflows, steps, and dependencies. Leaf operations are delegated to the bound runtime.

```go
func (d *Document) Execute(ctx context.Context) error {
    orch := NewOrchestrator(d, d.Runtime)
    return orch.Execute(ctx)
}
```

### Binding happens at execution time

The specialized engine binds its runtime implementation to the document before execution.

```go
func ExecuteWithRuntime(ctx context.Context, doc *uws1.Document, rt Runtime) error {
    doc.SetRuntime(rt)
    return doc.Execute(ctx)
}
```

Benefits of this model:

- Structural orchestration logic (`loop`, `switch`, `parallel`, `merge`, `await`) is defined once in `uws1`.
- `uws1` has zero knowledge of any concrete engine's specific runtimes.
- All UWS-compliant executors share the same orchestration behavior and bind their own runtime implementations.

## Schema / Spec / Code Sync

`uws1/schema_conformance_test.go` is the bridge between the schema and the Go validator. It reads the latest `versions/1.x.y.json` schema and asserts that the schema's required fields, enums, patterns, and related rule coverage match the Go-side validation rules.

When changing validation rules:

1. update the latest `versions/1.x.y.json`
2. update `uws1/validation.go`
3. update the matching `versions/1.x.y.md` when the public contract changed
4. make the schema conformance and parity tests pass

## Validation Layering

- The versioned JSON Schema covers structural and shape checks.
- `(*Document).Validate()` / `ValidateResult()` in `uws1/validation.go` cover semantic checks the schema cannot: duplicate identifiers, source binding rules, reference integrity across operations/workflows/steps/triggers/parallel groups/sourceDescriptions, action and criterion rules, trigger routes, and standard request-binding keys.
- `contenttrust.Analyze()` is an explicit advisory pass. Its findings do not enter ordinary validation or execution; profile-specific input/output semantics are supplied by resolvers.
- Use `Validate()` when a single `error` is enough.
- Use `ValidateResult()` when callers need path-tagged errors.

## Conversion

`convert/` provides JSON, YAML, and HCL helpers.

Key invariants:

- JSON and YAML preserve `x-*` extensions through the `Extensions` map pattern.
- HCL preserves object-level `x-*` extensions through `extensions { ... }` blocks. JSON and YAML keep extensions flattened as normal `x-*` fields.
- HCL key rewriting preserves `$`-prefixed keys on round-trip. Legacy JSON Schema keys (`$ref`, `$id`, `$schema`, `$defs`, etc.) use the `_`-prefix form in HCL; other `$foo` keys use `__dollar__foo`.
- `MarshalHCL` works on a deep copy and does not mutate the caller's document.
- The `uws1.Document` wire tree and all fields/sub-structs reachable from it should carry `json` and `hcl` tags for parsing.

## Extension Pattern

Every type that accepts `x-*` fields follows the same recipe:

1. an `Extensions map[string]any` field with `json:"-" yaml:"-" hcl:"-"`
2. a `knownFields` list of schema-owned keys
3. `UnmarshalJSON` using a type alias plus `rejectUnknownFields` and `extractExtensions`
4. `MarshalJSON` calling `marshalWithExtensions`

When adding a field to such a type, update both the struct and its `knownFields` list or the unmarshaller will reject valid documents.

## References

- [Peter Bi, "Achieving Full Object Inheritance in Go", Medium, 2024](https://medium.com/@peterbi_91340/implement-true-inheritance-in-go-ff6243bfd7a8)
