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
| Config extension | `x-uws-runtime-config` |
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

`http`, `ssh`, `cmd`, `fnct`, `fileio`, `sql`, `s3`, `smtp`, `dns`, `ldaps`,
`scp`, `sftp`, and `llm`.

The spelling is exact. `ldaps` is defined; plain `ldap` is not.

## Operation Runtime Payload

`x-uws-runtime` is an object with these fields:

| Field | Type | Purpose |
| --- | --- | --- |
| `type` | runtime type string | Runtime selector. |
| `isJson` | boolean | Runtime-owned JSON payload hint. |
| `host` | string | Host or remote endpoint hint. |
| `method` | string | HTTP method or runtime method hint. |
| `path` | string | Runtime path hint. |
| `payloadRequired` | boolean | Request payload required marker. |
| `requestMediaType` | string | Preferred request media type. |
| `responseMediaType` | string | Preferred response media type. |
| `responseStatusCode` | integer | Expected HTTP response status code. |
| `command` | string | Command text for command-like runtimes. |
| `workingDir` | string | Working directory for command-like runtimes. |
| `function` | string | Function name for function runtimes. |
| `workflow` | string | Nested workflow reference. |
| `arguments` | array | Runtime-owned argument values. |
| `provider` | object | Provider selection metadata. |
| `security` | array | Security requirement metadata. |
| `queryPars` | parameter schema | Query parameter schema. |
| `pathPars` | parameter schema | Path parameter schema. |
| `headerPars` | parameter schema | Header parameter schema. |
| `cookiePars` | parameter schema | Cookie parameter schema. |
| `payloadPars` | parameter schema | Request payload schema. |
| `responseBody` | parameter schema | Response body schema. |
| `responseHeaders` | parameter schema | Response header schema. |

The payload shape is intentionally descriptive. A conforming UWS parser
preserves it. A bound runtime decides whether it can execute it.

## Runtime Config Payload

`x-uws-runtime-config` is an object with optional `provider` and `security`
fields. It is suitable for document-level, component-level, or other
extension-bearing object metadata when tooling wants to keep runtime defaults
near the UWS artifact without making them UWS core semantics.

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
