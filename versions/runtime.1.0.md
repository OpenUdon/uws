# UWS Runtime Supplement 1.0

The UWS Runtime Supplement defines a small public extension profile for
runtime-owned operation metadata. It is wire/spec metadata only. It does not
standardize execution behavior, credentials, clients, process management, or
provider implementations.

## Profile

| Field | Value |
| --- | --- |
| Profile name | `uws.runtime.1.0` |
| Operation extension | `x-uws-runtime` |
| JSON Schema | `versions/runtime.1.0.json` |

An extension-owned operation using this supplement sets:

```yaml
operationId: render
x-uws-operation-profile: uws.runtime.1.0
x-uws-runtime:
  type: fnct
  function: identity
```

## Runtime Types

The supplement defines these runtime type identifiers:

`ssh`, `cmd`, `fnct`, `fileio`, `sql`, `s3`, `smtp`, `dns`, `ldaps`, `scp`,
`sftp`, and `llm`.

HTTP and OpenAPI-bound calls are represented by core UWS operation binding
fields, not by `x-uws-runtime`. A payload that assigns `type: http` in
`x-uws-runtime` is invalid. The spelling is exact. `ldaps` is defined; plain
`ldap` is not.

## Operation Runtime Payload

`x-uws-runtime` is an object with these fields:

| Field | Type | Purpose |
| --- | --- | --- |
| `type` | runtime type string | Non-HTTP runtime selector. |
| `command` | string | Command text for command-like runtimes. |
| `workingDir` | string | Working directory for command-like runtimes. |
| `function` | string | Function name for function runtimes. |
| `workflow` | string | Nested workflow reference. |
| `arguments` | array | Runtime-owned argument values. |

The payload shape is intentionally small. It selects a non-HTTP invocation
surface without standardizing runtime behavior. A bound runtime decides whether
it can execute the selected type and how to interpret the selector fields.

HTTP/OpenAPI operation metadata is not part of `x-uws-runtime`. HTTP method,
path, server, request/response schemas, and operation security requirements
belong in the referenced OpenAPI document and core UWS OpenAPI binding fields.

Runtime-specific credentials, provider selection, client defaults, connection
pools, security material, and other execution configuration belong in
runtime-private configuration or in a product-owned extension profile, not in
public `x-uws-*` fields.

For runtime types such as `ssh`, `cmd`, `fnct`, `fileio`, `sql`, `s3`, `smtp`,
`dns`, `ldaps`, `scp`, `sftp`, and `llm`, authentication and security shapes
vary by implementation and are difficult to standardize portably. When a
runtime needs invocation-specific auth hints, they are runtime-owned argument
values or private runtime configuration, not a public `x-uws-runtime-config`
contract.

## HCL Representation

In HCL, extension fields are represented inside an `extensions` block on any
UWS object that supports specification extensions:

```hcl
operation "render" {
  extensions {
    x-uws-operation-profile = "uws.runtime.1.0"
    x-uws-runtime {
      type     = "fnct"
      function = "identity"
    }
  }
}
```

JSON and YAML use normal flattened `x-*` fields.
