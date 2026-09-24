# UWS Versions Changelog

This changelog summarizes externally visible changes between published UWS
versioned schemas and specification documents. The versioned `.md` files remain
the normative human-readable specifications. Purely editorial clarifications to
already-published artifacts may land without an entry; any change that adjusts
the meaning or scope of a published schema or sub-spec is recorded as an
"Amended" note under the affected release.

## Version and Schema Selection Policy

- A published UWS document selects the schema whose filename exactly matches
  its declared `uws` version. The latest schema is not a fallback for older,
  prerelease, or unpublished versions.
- Core semantic validation and execution accept only declared versions for
  which this distribution publishes an exact UWS core schema. An external
  schema file for an unpublished or prerelease version does not make that
  version supported; document validation fails closed. This repository
  publishes no UWS prerelease schemas.
- Minimum-version feature gates use SemVer precedence among published stable
  versions. A gate applies to a later published stable version unless its
  release contract is version-specific. Later published rules are not applied
  retroactively to an earlier declared version. Prerelease identifiers must
  use valid SemVer syntax, but no prerelease UWS core schema is published or
  supported even though the schema's lexical pattern is broader.
- This policy clarifies version selection and validator behavior; it does not
  amend the bytes or normative requirements of earlier published artifacts.

## UWS 1.11.0 Contract Correction - 2026-09-24

- Published the UWS 1.11 core schema and corrected human-readable contract.
  Added opt-in response-body dot-walk, loop-only `$batchIndex`, complete JSON
  numeric literals for non-`await` `wait` and `batchSize`, root-scope terminal
  `goto`, and `forEach` merge selection that uses iteration records instead of
  duplicating their parent aggregate. The executor also accepts JSON-number
  runtime values for `batchSize` only when they are finite positive integers;
  literal grammar remains version-gated. The 1.11 grammar and execution
  behavior are gated by the declared, exactly published UWS version.
- Corrected the 1.10 prose on step forms, expression-addressable variable
  names, loop-only `items`, criteria contexts, JSON Pointer tilde escapes,
  header-name matching, criterion pointer normalization, version gates, and
  separately versioned browser/profile references. Added the complete
  published `x-uws-*` field inventory and pinned the normative standards
  references used by this specification.
- Exact-version admission remains a deliberate cross-version validator
  correction as documented below. Other implemented bug fixes that repair
  existing behavior across declared versions are: per-call workflow execution
  record isolation, exact JSON-number `batchSize` evaluation when the value is
  a positive integer, deterministic merge dependency/member/iteration
  ordering, and rejection of trigger output indexes outside the declared
  output list. These corrections do not enable later UWS 1.11 expression or
  control-flow features on older documents.
- UWS 1.10.0 and every earlier published core schema/specification remain
  immutable. Browser 1.9 remains a separate opt-in profile; no core wire field
  or Go API is added by this release.

## Cross-Version Validator Correction - 2026-09-23

- Direct UWS semantic validation and execution now reject syntactically valid
  prerelease, future, or unpublished patch/minor versions even when an exact
  external schema file is available. File validation continues to select only
  the exact declared schema and also rejects versions outside the published
  core-version set.
- This admission guard intentionally applies to every declared version,
  including older documents. It is an explicit bug-fix exception to the
  general promise that an older declaration retains prior behavior; it does
  not apply later published feature semantics retroactively. The `1.0.0-beta.1`
  string is valid SemVer syntax but is not a supported UWS version because no
  matching core schema is published.

## Browser 1.9 - 2026-09-23

- Added opt-in `uws.browser.1.9`, retaining the Browser 1.8 action, context,
  output, and scalar-substitution contract while defining deterministic
  `{{{{` / `}}}}` literal-brace escapes and one-pass substitution.
- Kept navigation substitutions component-encoded and encoded escaped braces
  as URI octets. Defined control/bidi rejection for `type_text.value` and
  confirmation prompts, text-only prompt rendering, and value-setting
  `type_text` behavior without keyboard/Enter/form-submit events.
- Aligned integer parameter inputs with the existing accessibility-output
  safe-integer range (`-9007199254740991` through `9007199254740991`). Added
  BCP 14 language and current UWS 1.10.0 core references.
- Preserved Browser 1.5–1.8 schemas/specs and behavior. Browser 1.9 requires an
  explicit profile selector; empty schema lookup retains Browser 1.8 as its
  compatibility default. No UWS core schema or wire-model change was made.

## Browser 1.8 - 2026-09-23

- Added `uws.browser.1.8` with single-pass `{{name}}` parameter templates in
  navigation path/query values, text/selection values, and confirmation
  prompts. URI components use RFC 3986 encoding; text sinks use literal scalar
  formatting. Ambiguous, malformed, non-scalar, and unsafe placements fail
  closed.
- Retained browser 1.5–1.7 artifacts and behavior unchanged. Browser 1.8
  requires UWS 1.9 and does not change the UWS core schema or wire model.

## 1.10.0 Portable Execution Contract - 2026-09-23

- Published version-gated truthiness, criterion, JSON Pointer, wait, entry,
  action, trigger, and structural-result execution semantics without changing
  the UWS 1.x object vocabulary.
- Restricted expression-addressable declaration names in UWS 1.10 documents
  to the grammar's identifier characters; earlier UWS versions retain their
  existing name rules and execution behavior.
- Defined ordered loop, `forEach`, merge, action, switch, and trigger behavior;
  rejected trigger dispatch output indexes outside the declared output list.
- Added security guidance for exact source validation, trigger ingress,
  credential handling, egress, budgets, retries, and advisory content-trust.
- Preserved every previously published core schema/specification and all
  separately versioned profiles.

## 1.9.2 Compatible Corrections - 2026-09-19

- Rejected `steps`, `cases`, and `default` properties on operation- and
  workflow-reference steps, including explicitly present empty arrays.
- Required canonical decimal array-index tokens in JSON Pointer criteria.
- Aligned advisory content-trust analysis with core lexical output evaluation,
  retained untrusted provenance through operation-output aliases, and stopped
  treating ignored merge children as executable or dominating.
- Matched trigger reachability to runtime namespace rules: trigger targets use
  workflow-first dispatch, generic dependencies use step/workflow/operation
  precedence, and `operationRef` and `workflow` remain typed edges.
- Preserved UWS 1.9.1 and every earlier published artifact byte-for-byte.

## Browser Registration 1.2 - 2026-09-13

- Added `uws.browser-registration.1.2` and
  `uws.browser-registration-call.1.2`, retaining the 1.1 typed private inputs,
  checkpoints, and one-attempt controls while requiring a reviewed
  human-verification policy.
- Bounded supported Turnstile, reCAPTCHA v2, and hCaptcha adapters by exact
  provider policy, request/response budgets, frame ancestry, application POST,
  and nonrenewable operation deadline. Human challenges remain human-operated.
- Required explicit 1.2 selection and support before browser execution. Kept
  UWS core, the private input envelope, earlier registration documents, 1.0
  schema/call defaults, and the 1.1 binding-helper default unchanged.

## Browser Registration 1.1 - 2026-09-07

- Added `uws.browser-registration.1.1` and
  `uws.browser-registration-call.1.1` with typed private form fields, explicit
  input checkpoints, and a symbolic runtime-owned input binding.
- Added `uws.browser-registration-input.1.0` plus browser-free template and
  update validation helpers. Filled input documents and their digests remain
  private and outside packages, source control, prompts, logs, and reports.
- Retained registration profile/call 1.0 and the existing 1.0 schema and call
  supplement defaults; the unversioned binding helper selects 1.1.

## 1.9.1 Content Trust - 2026-08-26

- Added the optional root `contentTrust` registry with source-description,
  operation default/output, trigger, and workflow default/input provenance
  declarations. Existing operation, workflow, and step output shapes are
  unchanged.
- Added patch-aware semantic gating and identifier/output/input integrity
  checks. Malformed declarations are invalid; every 1.9.0 and earlier document
  retains its validation and execution behavior.
- Added deterministic advisory analysis in package `contenttrust`, with
  resolver-defined data/instruction/authority channels, separate provenance and
  value capability, control-flow-aware expression references, and stable
  findings that contain no runtime values or content excerpts.
- Kept findings outside `ValidationResult` and all execution paths. Reserved
  the remaining 1.9 patch line for evidence-driven compatible refinements;
  breaking defaults, enforcement, or wire changes remain UWS 2.0 work.
- Preserved UWS 1.9.0 and every existing profile artifact byte-for-byte.

## Browser Registration 1.0 - 2026-08-25

- Added the separate `uws.browser-registration.1.0` inert, secret-free account-
  creation recipe and `uws.browser-registration-call.1.0` operation envelope.
- Required symbolic credential bindings, exact safe origins, exactly one
  explicitly approved submit, fixed fail-on-duplicate and stop-without-retry
  controls, and a preselected cleanup disposition.
- Kept account values, verification responses, browser state, captures,
  retries, runtime execution, and cleanup outside portable artifacts. Existing
  UWS core, browser, and browser-authentication contracts are unchanged.

## 1.9.0 and Typed Browser Outputs - 2026-08-17

- Published UWS 1.9.0 with the UWS core document shape unchanged from 1.8.
- Added `uws.browser.1.7`, retaining browser 1.6 contexts while defining
  locale-free conversion of trimmed accessibility text to string, safe
  integer, finite JSON number, or lowercase Boolean values.
- Required browser 1.7 accessibility outputs to remain scalar, kept presence
  as Boolean match behavior without text access, and failed closed on empty,
  noncanonical, non-finite, or out-of-range text.
- Kept `uws.browser-authentication.1.1` unchanged and retained every earlier
  UWS and browser schema as an immutable compatibility contract.

## 1.8.0 and Browser Context Profiles - 2026-08-16

- Published UWS 1.8.0 with the UWS core document shape unchanged from 1.7.
- Added `uws.browser.1.6` and `uws.browser-authentication.1.1` optional bounded
  popup/frame context graphs, context-qualified actions/waits/outputs/success,
  the navigate object form, explicit `opensContext`, and exact success paths.
- Added the wire-compatible `uws.browser-authentication-call.1.1` supplement.
- Made Go profile validators dispatch on exact document discriminators,
  compile each embedded schema once, and reject unknown versions, fields,
  references, cycles, excessive depth, unsafe paths, missing/duplicate popup
  bindings, and context/origin mismatches.
- Kept every UWS 1.7, browser 1.5, browser-authentication 1.0, and call 1.0
  document immutable and accepted. Main-only authors SHOULD emit the oldest
  sufficient profile version.

## Browser Authentication 1.0 - 2026-08-15

- Added the additive `uws.browser-authentication.1.0` secret-free sign-in
  recipe schema and validation helpers.
- Added `uws.browser-authentication-call.1.0`,
  `x-uws-browser-authentication`, and `x-uws-browser-session` for explicit
  named-session establishment and use.
- Kept UWS 1.7 core and `uws.browser.1.5` unchanged.

## Ansible Support Retirement - 2026-08-15

- Retired the UWS-owned `uws.ansible-module-call.1.0` supplement and its typed
  Go helpers. UWS 1.7 defines no Ansible source or operation profile.
- Removed the Ansible-specific schema accessors while retaining
  `versions/1.6.0.{json,md}`, the version-gated 1.6 Go model behavior, and
  `versions/ansible.1.0.{json,md}` as historical material.
- Ramen localizes its static conversion-only schemas and wire helpers under its
  `internal/ansibleconvert` package with Ramen-owned identifiers; it is not required to retain
  or accept the removed UWS-named contracts.
- The removed implementation and module-call documents remain available from
  repository history at commit `a68a209`.
- Moved Go schema access and browser-profile validation from package `versions`
  to package `schemas`; `versions/` now contains documents only.

## 1.7.0 - 2026-08-08

- Removed the `ansible-module` `sourceDescription.type` introduced in 1.6.0.
  UWS 1.7 defines nine source description type values.
- Rationale: first-class source descriptions name operations that already exist
  on the remote target. Provider-published API descriptions are the strongest
  form; an author-asserted browser profile is weaker but still names actions the
  target UI implements. An Ansible host does not expose a collection module as
  a pre-existing operation: the control node supplies and runs that
  implementation, so the argspec is a client-side library manifest.
- UWS 1.7 defines no replacement Ansible operation profile. The briefly
  published module-call supplement was retired on 2026-08-15.
- `versions/ansible.1.0.{json,md}` remains only as historical UWS 1.6 material;
  it is not a UWS 1.7 source profile.
- Documents declaring `uws: 1.6.0` are unaffected; `versions/1.6.0.json` still
  accepts `ansible-module`. Validators MUST reject the type on documents
  declaring 1.7.0 or later.

## Go Schema Accessors - 2026-07-25

- Added `versions.PathForAnsibleSourceProfile` and
  `versions.PathForBrowserSourceProfile` so source-aware tooling can resolve
  the published profile schemas through the supported environment, package,
  module-cache, embedded-schema, and sibling fallback lookup sequence.
- The accessors accept bare versions, profile-prefixed names, optional `uws.`
  prefixes, and optional `.json` suffixes. They default to `ansible.1.0` and
  `browser.1.5`, respectively.
- No schema wire format changed.

Withdrawn 2026-08-15:

- Schema access moved to package `schemas`; the Ansible-specific accessor was
  removed, while browser-profile lookup and validation remain supported there.

## Ansible Module Call Supplement 1.0 (withdrawn) - 2026-06-14

> Withdrawn on 2026-08-15. This entry is retained as release history, not as a
> current UWS operation profile.

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

Amended 2026-08-08:

- Repositioned the supplement as the supported Ansible module-call
  representation at every UWS version after UWS 1.7 withdrew the temporary
  source type. Updated schema metadata only; the supplement wire shape is
  unchanged.

Withdrawn 2026-08-15:

- Removed the supplement from the current tree. UWS 1.7 no longer standardizes
  Ansible automation. The historical schema, prose, and Go helper package are
  available at commit `a68a209`.

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
