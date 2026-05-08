# UWS Versions Changelog

This changelog summarizes externally visible changes between published UWS
versioned schemas and specification documents. The versioned `.md` files remain
the normative human-readable specifications.

## Runtime Supplement 1.0 - 2026-05-08

- Added `uws.runtime.1.0` as a public runtime metadata supplement.
- Defined the slim `x-uws-runtime` operation payload schema.
- Defined non-HTTP runtime type identifiers: `ssh`, `cmd`, `fnct`, `fileio`, `sql`, `s3`,
  `smtp`, `dns`, `ldaps`, `scp`, `sftp`, and `llm`.
- Kept HTTP/OpenAPI metadata, provider configuration, credentials, security configuration,
  and request/response schemas outside the public runtime supplement.

## 1.1.0 - 2026-04-28

- Added portable `timeout` fields on Operation, Workflow, and Step objects.
- Added workflow-level `idempotency` metadata for logical workflow-run de-duplication.
- Added validation requirements for positive timeout values, required idempotency keys,
  `onConflict` enum values, and positive `ttl` values.
- Clarified that idempotency storage, retry replay protection, and timeout enforcement details
  are executor responsibilities outside the serialized wire format.

## 1.0.0 - 2026-04-26

- Initial UWS 1.0.0 specification and JSON Schema.
- Defined OpenAPI-bound operations, workflow structure, request binding, structural control flow,
  triggers, results, success criteria, failure/success actions, runtime expressions, and extension
  profiles.
- Reserved the `x-uws-` extension prefix and defined `x-uws-operation-profile`.
