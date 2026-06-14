# UWS Versions Changelog

This changelog summarizes externally visible changes between published UWS
versioned schemas and specification documents. The versioned `.md` files remain
the normative human-readable specifications.

## Ansible Module Call Supplement 1.0 - 2026-06-14

- Added `uws.ansible-module-call.1.0` for extension-owned Ansible module leaf
  operations in UWS 1.5-compatible documents.
- Defined the slim `x-uws-ansible-module` payload with required module FQCN and
  optional argspec review reference.
- Kept Ansible execution, inventory, credentials, vault material, `become`,
  strategy, async/poll, check mode, and callbacks outside public UWS metadata.

## 1.6.0 - 2026-06-02

- Added `ansible-module` as a first-class `sourceDescription.type` for Ansible
  collection argspec source documents.
- Reused generic `sourceOperationId` and `sourceOperationRef` selectors;
  `sourceOperationId` is the module's fully qualified collection name (FQCN,
  e.g. `ansible.builtin.apt`), and the preferred `sourceOperationRef` form is
  `#/modules/<fqcn>`.
- Published the Ansible module sub-spec separately as
  `versions/ansible.1.0.{json,md}` so UWS core stays thin (mirrors the
  `runtime.1.0` and `browser.1.5` splits): UWS core only references the type
  name and selector rules; the argspec document shape, the `changed` output
  convention, and handler-lowering rules live in the sub-spec.
- Kept connection plugins, privilege escalation, forks/serial strategy, check
  mode, vault material, Jinja2 templating, and the module execution environment
  runtime-owned or authoring-tool-owned.

## 1.5.0 - 2026-05-30

- Added `browser-profile` as a first-class `sourceDescription.type` for browser
  capability profile source documents.
- Reused generic `sourceOperationId` and `sourceOperationRef` selectors;
  `sourceOperationId` identifies an action key in the profile's `actions`
  object, `sourceOperationRef` preferred form is `#/actions/<name>`.
- Published the browser profile sub-spec separately as
  `versions/browser.1.5.{json,md}` so UWS core stays thin (mirrors the
  `runtime.1.0` split): UWS core only references the type name and selector
  rules; locator vocabulary, macro action vocabulary, output extraction methods,
  and safety controls live in the sub-spec.
- Kept browser session contexts, credentials, cookies, rendering, and
  interaction protocols runtime-private.

## 1.4.0 - 2026-05-27

- Added first-class `sourceDescription.type` values for `graphql`, `openrpc`,
  `grpc-protobuf`, and `odata`.
- Required generic `sourceOperationId` or `sourceOperationRef` selectors for
  those source families and kept legacy OpenAPI selectors limited to
  `openapi` sources.
- Kept source parsing, selector resolution, transport, authentication,
  credentials, and runtime invocation source-aware or runtime-owned.

## 1.3.0 - 2026-05-27

- Added `asyncapi` as a first-class `sourceDescription.type` for AsyncAPI source documents.
- Reused generic `sourceOperationId` and `sourceOperationRef` selectors for AsyncAPI operations.
- Defined `sourceOperationId` for AsyncAPI as a root AsyncAPI Operation Object key.
- Defined `sourceOperationRef` compatibility targets for `#/operations/...`, `#/channels/...`,
  and `#/channels/.../messages/...`.
- Kept browser capability profiles outside UWS core as extension-owned operations.

## Runtime Supplement 1.0 - 2026-05-08

- Added `uws.runtime.1.0` as a public runtime metadata supplement.
- Defined the slim `x-uws-runtime` operation payload schema.
- Defined non-HTTP runtime type identifiers: `ssh`, `cmd`, `fnct`, `fileio`, `sql`, `s3`,
  `smtp`, `dns`, `ldaps`, `scp`, `sftp`, and `llm`.
- Kept HTTP/OpenAPI metadata, provider configuration, credentials, security configuration,
  and request/response schemas outside the public runtime supplement.

## 1.2.0 - 2026-05-22

- Added first-class `sourceDescription.type` values for `openapi`, `google-discovery`, and
  `aws-smithy`; omitted type remains `openapi`.
- Added canonical operation selectors `sourceOperationId` and `sourceOperationRef`.
- Kept `openapiOperationId` and `openapiOperationRef` as backward-compatible selectors for
  OpenAPI source descriptions only.
- Required generic selectors for Google Discovery and AWS Smithy sources.

## 1.1.1 - 2026-05-20

- Clarified the long-term source model position: UWS core remains OpenAPI-first.
- Recognized Google Discovery and AWS Smithy as source model families that compliant tooling may
  lower into UWS/OpenAPI-bound operations.
- Reserved first-class native Discovery/Smithy source binding for a future minor version if
  interoperability demands it.

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
