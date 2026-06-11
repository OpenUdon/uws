# UWS

<p align="center">
  <img src="docs/assets/uws2.png" alt="UWS logo" width="520">
</p>

UWS is the Udon Workflow Specification Go package. It defines the UWS 1.x document model, JSON Schema, validation helpers, and JSON/YAML/HCL conversion helpers.

UWS is a workflow **overlay** over source documents — OpenAPI, AsyncAPI, GraphQL, OpenRPC, Protocol Buffers, OData, and the like. The source document owns the operations: methods, paths, channels, messages, schemas, servers, and security are all defined server-side and authoritative. UWS adds only what the source cannot express: operation binding, workflow structure, request values, outputs, triggers, and control flow.

This is what distinguishes UWS from full client-side workflow tools such as Arazzo and IaC engines such as OpenTofu and Terraform. Arazzo describes full client-side action sequences and treats each step as a bespoke client action. OpenTofu and Terraform act as full client-side workflow engines for infrastructure: each resource and provider call is described in the client configuration and resolved against a provider plugin at apply time. Neither approach assumes that the underlying operations are already defined by a server contract. UWS takes the opposite position: server actions are pre-defined by the source document, and UWS workflows reference those operations by ID rather than re-describing them. The result is a much smaller overlay: UWS does not duplicate request/response shapes, does not redeclare endpoints, and does not encode anything the source document already specifies.

UWS 1.6 keeps OpenAPI compatibility and adds first-class `ansible-module` source descriptions alongside `openapi`, `google-discovery`, `aws-smithy`, `asyncapi`, `graphql`, `openrpc`, `grpc-protobuf`, `odata`, and `browser-profile`. The Ansible module sub-spec is published separately as `versions/ansible.1.0.{json,md}`, joining `runtime.1.0` and `browser.1.5`. Missing `sourceDescription.type` still defaults to `openapi`; legacy `openapiOperationId` and `openapiOperationRef` remain valid for OpenAPI sources.

### Version highlights

| Version | Adds |
|---|---|
| **1.0** | Initial spec: OpenAPI-bound operations, workflow structure, request binding, structural control flow (sequence/parallel/switch/loop/merge/await), triggers, results, success criteria, success/failure actions, runtime expressions, and the `x-uws-` extension prefix with `x-uws-operation-profile`. |
| **1.1** | Portable `timeout` on operations/workflows/steps; workflow-level `idempotency` metadata for run de-duplication. |
| **1.2** | First-class `sourceDescription.type` for `openapi`, `google-discovery`, `aws-smithy`; canonical `sourceOperationId` / `sourceOperationRef` selectors. Legacy `openapiOperationId` / `openapiOperationRef` kept for OpenAPI sources. |
| **1.3** | First-class `asyncapi` source type; AsyncAPI operation selector rules including `#/operations/...`, `#/channels/...`, and `#/channels/.../messages/...` ref forms. |
| **1.4** | First-class `graphql`, `openrpc`, `grpc-protobuf`, and `odata` source types; generic selectors required for those families. |
| **1.5** | First-class `browser-profile` source type; capability profile sub-spec published separately as `versions/browser.1.5.{json,md}`. |
| **1.6** | First-class `ansible-module` source type (FQCN as `sourceOperationId`, `#/modules/<fqcn>` refs); argspec sub-spec published separately as `versions/ansible.1.0.{json,md}`. |

See [`versions/CHANGELOG.md`](versions/CHANGELOG.md) for the full changelog.

Non-source runtimes such as command execution, function calls, file I/O, SSH, SQL, browser automation, or LLM calls are extension-profile concerns represented with `x-*` fields, not UWS core service types. Operations without a source binding are extension-owned and require `x-uws-operation-profile` to name the implementation profile that can execute them. The optional `uws.runtime.1.0` supplement standardizes a small `x-uws-runtime` selector payload for those extension-owned operations.


[![GoDoc](https://godoc.org/github.com/OpenUdon/uws?status.svg)](https://godoc.org/github.com/OpenUdon/uws)


## Documentation

- **Docs site**: [openudon.github.io/uws](https://openudon.github.io/uws/)
- Human-readable specification: [versions/1.6.0.md](versions/1.6.0.md)
- Runtime supplement: [versions/runtime.1.0.md](versions/runtime.1.0.md)
- Runtime supplement schema: [versions/runtime.1.0.json](versions/runtime.1.0.json)
- Browser profile supplement: [versions/browser.1.5.md](versions/browser.1.5.md) / [versions/browser.1.5.json](versions/browser.1.5.json)
- Ansible module supplement: [versions/ansible.1.0.md](versions/ansible.1.0.md) / [versions/ansible.1.0.json](versions/ansible.1.0.json)
- UWS 1.6 Ansible design note: [docs/uws_1_6_ansible.md](docs/uws_1_6_ansible.md)
- JSON Schema: [versions/1.6.0.json](versions/1.6.0.json)

## Packages

- `uws1` contains the UWS 1.x Go model, structural vocabulary, and structural validation.
- `convert` converts UWS documents between JSON, YAML, and the HCL authoring form.
- `runtimes` contains the public `uws.runtime.1.0` supplement constants, wire structs, and extension helpers.
- `versions/1.6.0.md` is the human-readable UWS 1.6 specification.
- `versions/1.6.0.json` is the JSON Schema for UWS 1.6 documents.
- `versions/browser.1.5.md` / `versions/browser.1.5.json` publish the browser capability profile sub-spec referenced by `sourceDescriptions[].type: browser-profile`.
- `versions/ansible.1.0.md` / `versions/ansible.1.0.json` publish the Ansible module sub-spec referenced by `sourceDescriptions[].type: ansible-module`.

## Validation

Use `(*uws1.Document).Validate()` when an `error` is enough, or `ValidateResult()` when callers need all path-tagged validation errors.

```go
result := doc.ValidateResult()
if !result.Valid() {
    return result
}
```

Validation checks required root fields, source operation bindings, extension-owned operation profiles, duplicate identifiers, standard request-binding keys, known structural types, selected reference integrity, action/criterion rules, and trigger routes.

`versions/1.6.0.json` provides structural JSON Schema validation. Use the Go validator for semantic checks such as duplicate identifiers and reference integrity.

The separate `versions/runtime.1.0.json` schema validates the public runtime supplement payload. It requires `x-uws-runtime.type`, accepts only the non-HTTP runtime identifiers defined by the supplement, and rejects HTTP/API/event source metadata because HTTP and event calls are represented by core source operation binding fields.

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

Large round-trip fixtures under `testdata/big/` exercise the HCL/JSON converter with runtime supplement metadata and multi-file source references.

## Development

```bash
go test ./...
go vet ./...
```

## License

Apache License 2.0. See [LICENSE](LICENSE).
