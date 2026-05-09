# Runtimes

UWS core defines workflow orchestration. Concrete runtimes bind at execution time and own leaf work:
HTTP calls, expression evaluation, item resolution, credentials, clients, provider configuration,
process management, storage, and product-specific behavior.

## Runtime Supplement 1.0

The public `uws.runtime.1.0` supplement is a small metadata profile for non-HTTP leaf operations.
It is used on extension-owned operations:

```yaml
operationId: render_message
x-uws-operation-profile: uws.runtime.1.0
x-uws-runtime:
  type: fnct
  function: render_message
  arguments:
    - template: "Ticket {{ticket_id}} is ready"
```

`x-uws-runtime.type` is required. Valid runtime type selectors are:

`ssh`, `cmd`, `fnct`, `fileio`, `sql`, `s3`, `smtp`, `dns`, `ldaps`, `scp`, `sftp`, and `llm`.

HTTP is intentionally not a runtime type. HTTP/OpenAPI operations use core UWS operation binding
fields: `sourceDescription` plus `openapiOperationId` or `openapiOperationRef`.

## Payload Fields

| Field | Purpose |
|-------|---------|
| `type` | Required non-HTTP runtime selector. |
| `command` | Command text for command-like runtimes. |
| `workingDir` | Working directory for command-like runtimes. |
| `function` | Function name for function runtimes. |
| `workflow` | Nested workflow reference. |
| `arguments` | Runtime-owned argument values. |

The supplement does not standardize credentials, hosts, provider selection, client behavior,
security configuration, result schemas, or execution side effects. Those belong to the bound runtime
or product-owned extension profiles.

## References

- [Runtime Supplement 1.0](https://github.com/OpenUdon/uws/blob/main/versions/runtime.1.0.md)
- [Runtime Supplement JSON Schema](https://github.com/OpenUdon/uws/blob/main/versions/runtime.1.0.json)
- [Extension Profiles](08-Extension-Profiles.md)
- [Execution Model](07-Execution-Model.md)
