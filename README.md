# UWS

UWS is the Udon Workflow Specification Go package. It defines the UWS 1.x document model, JSON Schema, validation helpers, and JSON/YAML/HCL conversion helpers.

UWS is similar in role to Arazzo and complements OpenAPI and AsyncAPI, but it is a smaller workflow overlay for OpenAPI-backed HTTP operations only. OpenAPI owns methods, paths, schemas, servers, and security. UWS owns operation binding, workflow structure, request values, outputs, triggers, and control flow.

Non-OpenAPI runtimes such as command execution, function calls, file I/O, SSH, SQL, or LLM calls are extension-profile concerns represented with `x-*` fields, not UWS core service types. Operations without an OpenAPI binding are extension-owned and require `x-uws-operation-profile` to name the implementation profile that can execute them. The optional `uws.runtime.1.0` supplement standardizes a small `x-uws-runtime` selector payload for those extension-owned operations.


[![GoDoc](https://godoc.org/github.com/OpenUdon/uws?status.svg)](https://godoc.org/github.com/OpenUdon/uws)


## Documentation

- **Docs site**: [openudon.github.io/uws](https://openudon.github.io/uws/)
- Human-readable specification: [versions/1.1.0.md](versions/1.1.0.md)
- Runtime supplement: [versions/runtime.1.0.md](versions/runtime.1.0.md)
- Runtime supplement schema: [versions/runtime.1.0.json](versions/runtime.1.0.json)
- JSON Schema: [versions/1.1.0.json](versions/1.1.0.json)

## Packages

- `uws1` contains the UWS 1.x Go model, structural vocabulary, and structural validation.
- `convert` converts UWS documents between JSON, YAML, and the HCL authoring form.
- `runtimes` contains the public `uws.runtime.1.0` supplement constants, wire structs, and extension helpers.
- `versions/1.1.0.md` is the human-readable UWS 1.1 specification.
- `versions/1.1.0.json` is the JSON Schema for UWS 1.1 documents.

## Validation

Use `(*uws1.Document).Validate()` when an `error` is enough, or `ValidateResult()` when callers need all path-tagged validation errors.

```go
result := doc.ValidateResult()
if !result.Valid() {
    return result
}
```

Validation checks required root fields, OpenAPI operation bindings, extension-owned operation profiles, duplicate identifiers, standard request-binding keys, known structural types, selected reference integrity, action/criterion rules, and trigger routes.

`versions/1.1.0.json` provides structural JSON Schema validation. Use the Go validator for semantic checks such as duplicate identifiers and reference integrity.

The separate `versions/runtime.1.0.json` schema validates the public runtime supplement payload. It requires `x-uws-runtime.type`, accepts only the non-HTTP runtime identifiers defined by the supplement, and rejects HTTP/OpenAPI metadata because HTTP calls are represented by core OpenAPI operation binding fields.

## Execution

UWS 1.x defines a bound-runtime execution model. UWS core owns orchestration and structural execution semantics; the bound runtime owns leaf execution plus the evaluation services needed for expressions and iterative constructs.

At a high level:

- `Document.Execute(ctx)` executes the document through the orchestrator
- `Document.DispatchTrigger(ctx, triggerID, output, payload)` dispatches a trigger event into the same execution model
- `Document.ExecutionRecords()` exposes the accumulated execution snapshot
- `Runtime` is responsible for leaf execution, expression evaluation, and item resolution

Execution requires a bound runtime and a document that passes validation for execution. Trigger dispatch resolves outputs by label or decimal index and routes only to declared workflows or top-level entry-workflow steps.

## Interchange

The `convert` package provides JSON, YAML, and HCL helpers such as `JSONToHCL`, `HCLToJSON`, and `MarshalYAML`. `MarshalHCL` works on a deep copy and does not mutate the caller-owned document.

HCL conversion preserves dynamic map keys such as `$ref` through reversible key rewriting. JSON and YAML preserve `x-*` extensions through the JSON extension model; HCL represents object-level extensions with `extensions { ... }` blocks and flattens them back to `x-*` fields when converting to JSON or YAML.

Large round-trip fixtures under `testdata/big/` exercise the HCL/JSON converter with runtime supplement metadata and multi-file OpenAPI references.

## Development

```bash
go test ./...
go vet ./...
```

## License

Apache License 2.0. See [LICENSE](LICENSE).
