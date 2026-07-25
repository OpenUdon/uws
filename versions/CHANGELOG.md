# UWS Versions Changelog

This changelog summarizes externally visible changes between published UWS
versioned schemas and specification documents. The versioned `.md` files remain
the normative human-readable specifications. Purely editorial clarifications to
already-published artifacts may land without an entry; any change that adjusts
the meaning or scope of a published schema or sub-spec is recorded as an
"Amended" note under the affected release.

## Go Schema Accessors - 2026-07-25

- Added `versions.PathForAnsibleSourceProfile` and
  `versions.PathForBrowserSourceProfile` so source-aware tooling can resolve
  the published profile schemas through the supported environment, package,
  module-cache, embedded-schema, and sibling fallback lookup sequence.
- The accessors accept bare versions, profile-prefixed names, optional `uws.`
  prefixes, and optional `.json` suffixes. They default to `ansible.1.0` and
  `browser.1.5`, respectively.
- No schema wire format changed.

## Ansible Module Call Supplement 1.0 - 2026-06-14

- Added `uws.ansible-module-call.1.0` for extension-owned Ansible module leaf
  operations in UWS 1.5-compatible documents.
- Defined the slim `x-uws-ansible-module` payload with required module FQCN and
  optional argspec review reference.
- Kept Ansible execution, inventory, credentials, vault material, `become`,
  strategy, async/poll, check mode, and callbacks outside public UWS metadata.

Amended 2026-07-23:

- Tightened the module FQCN contract to the canonical lowercase
  `namespace.collection.module` form, matching the module-key pattern in the
  `uws.ansible.1.0` argspec schema.

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

Amended 2026-06-13:

- Rescoped the `ansible.1.0` sub-spec to inert conversion/review metadata:
  argspec documents carry the module leaf contract only, and validating or
  resolving them is the job of source-aware conversion and review tooling, not
  runtimes. The `1.6.0.json` selector descriptions were updated to match
  ("inert argspec document"; runtime validation language removed).

Amended 2026-06-14:

- Reclassified handler lowering (and other playbook control-flow mappings) as
  converter-owned conventions documented in `ansible.1.0.md`, not
  source-profile fields.

Amended 2026-07-23:

- Required every module key in a `uws.ansible.1.0` argspec document to begin
  with the declared collection value; source-aware tooling enforces this
  cross-field constraint.

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

Amended 2026-05-29:

- Corrected the `1.4.0.json` top-level description, which incorrectly said
  "v1.3.x" at release.

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
