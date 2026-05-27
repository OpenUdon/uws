# Future Source Profiles

This note records source-profile roadmap material around UWS 1.3 and later, and
separates import/advisory tooling features from normative UWS source types. It
is explanatory documentation; the normative wire contract lives in the versioned
specification and JSON Schema.

## Status And Normativity

This document is not part of the normative UWS 1.4 wire contract. The published
UWS 1.4 source description types are `openapi`, `google-discovery`,
`aws-smithy`, `asyncapi`, `graphql`, `openrpc`, `grpc-protobuf`, and `odata`;
missing `sourceDescriptions[].type` still defaults to `openapi`.

The `v0.1.x` labels below are implementation-roadmap labels. They do not define
released UWS versions, JSON Schema changes, Go model changes, or new validator
behavior by themselves. AsyncAPI has graduated into UWS 1.3, and GraphQL,
OpenRPC, gRPC/protobuf, and OData have graduated into UWS 1.4.
`browser-profile` is not a valid UWS 1.4 `sourceDescriptions[].type` value.

Examples in this note are illustrative future shapes unless a section explicitly
labels them as valid today. Current UWS 1.4 documents that need non-core,
runtime-owned work, including browser/UI work, use extension-owned operations
with `x-uws-operation-profile`.

The current UWS model remains source-document-first:

- OpenAPI, Google Discovery, AWS Smithy, AsyncAPI, GraphQL, OpenRPC,
  gRPC/protobuf, and OData describe API, RPC, or event source contracts.
- UWS describes workflow structure, operation binding, request values, outputs,
  triggers, and control flow.
- Non-HTTP behavior is represented by extension-owned operations and named
  operation profiles.

Future source-profile tracks are listed separately:

| Track | Theme | Purpose |
| --- | --- | --- |
| UWS 1.1-compatible tooling | Import and advisory evidence | Use Dropbox Stone, Postman Collection, RAML, and API Blueprint as local source evidence only after conversion to OpenAPI or as review-only metadata. |
| `v0.1.3` | AsyncAPI source profiles | Graduated in UWS 1.3 for event, message, and subscription contracts when a provider publishes AsyncAPI. |
| UWS 1.4 source profiles | Formal API source families | Graduated `graphql`, `openrpc`, `grpc-protobuf`, and `odata` as normative source type values. |
| Later candidate | Browser capability profiles | Reuse reviewed browser/UI capability maps when no sufficient API source document exists. |

Browser capability profiles become UWS source profiles only if a future UWS
version explicitly adopts them.

## UWS 1.1-Compatible Import And Advisory Features

UWS 1.1 remains OpenAPI-bound at the wire level. That makes Dropbox Stone,
Postman Collection, RAML, and API Blueprint useful UWS 1.1-compatible tooling
features only when they are converted or lowered into reviewed OpenAPI, or when
they remain non-executable review evidence next to a UWS package. They are not
UWS 1.1 `sourceDescriptions[].type` values.

The intended posture is:

| Source family | UWS 1.1-compatible posture |
| --- | --- |
| Dropbox Stone | Official Dropbox source metadata and advisory-overlay provenance. Use reviewed Stone-derived OpenAPI-shaped output when available; do not treat Stone itself as an OpenAPI document. |
| Postman Collection | Local client/workspace evidence. Convert only when provenance proves provider ownership and useful API completeness; otherwise keep advisory. |
| RAML | Local API description evidence. Convert to OpenAPI when methods, resources, parameters, bodies, schemas/examples, auth hints, and provenance survive review. |
| API Blueprint | Local documentation-oriented API description evidence. Convert only when Markdown ambiguity and missing schema/auth detail do not make the OpenAPI output misleading. |

This keeps UWS 1.1 useful for import workflows without changing its source
description enum. The executable UWS package still binds to OpenAPI after
conversion, or it treats the artifact as advisory review evidence outside the
wire contract.

## Adopted in UWS 1.4: Source Family Baseline

UWS 1.4 is strict: the adopted source families are `graphql`, `openrpc`,
`grpc-protobuf`, and `odata`. Each source type has a clear operation selector
model, source-owned request/response or message metadata, parser/index evidence
or a bounded tooling path, and a runtime path that can consume the source
semantics without flattening important behavior into a weaker OpenAPI-shaped
approximation.

The preferred evaluation order is:

| Priority | Candidate | UWS posture |
| --- | --- | --- |
| 1 | GraphQL SDL or introspection | Strong first-class candidate. GraphQL has operation, type, variable, query, mutation, and subscription semantics that are not well represented by pretending every operation is REST. |
| 2 | OpenRPC | Strong first-class candidate for JSON-RPC services. Method names, params, and result schemas map cleanly to source operations and request bindings. |
| 3 | gRPC / Protocol Buffers | Strong first-class candidate for service/method/message contracts, especially internal service automation. Runtime execution remains implementation-owned. |
| 4 | OData CSDL | Enterprise candidate for entity sets, actions, functions, query semantics, and metadata-driven application APIs. It should preserve OData semantics instead of hiding them behind generic REST overlays. |

This ordering keeps UWS source-document-first. Provider-owned or protocol-owned
contracts should come before inferred browser/network traces. Import/advisory
formats such as Stone, Postman Collection, RAML, and API Blueprint stay on the
UWS 1.1-compatible tooling path unless a later roadmap explicitly reopens them.
Browser profiles remain useful for UI-only workflows, but they are weaker than
formal API source documents and should not lead the next source-type expansion.

### UWS 1.4 Shape

At a high level, UWS 1.4 is a source-family expansion, not a replacement for
the UWS 1.3 model. The shape is:

- Add narrowly scoped source description types for the graduating formats:
  `graphql`, `openrpc`, `grpc-protobuf`, and `odata`.
- Keep source documents authoritative for protocol-specific semantics: GraphQL
  owns schemas and variables, OpenRPC owns JSON-RPC methods and params,
  protobuf owns services and messages, and OData owns entity/action/query
  metadata.
- Keep UWS operations bound through the existing selector pattern where
  possible: `sourceDescription` plus `sourceOperationId` or
  `sourceOperationRef`, with source-aware tooling validating whether a selector
  resolves to a supported operation/method/action.
- Keep request values in UWS. The source profile describes where values belong
  and what shapes they may take; the workflow supplies inputs, dependencies,
  constants, and data flow.
- Keep runtime transport, authentication, credential lookup, channel setup, and
  side-effect execution outside the UWS wire contract.
- Treat Stone, Postman/RAML/API Blueprint importers, and browser/network
  capture as advisory or conversion evidence unless a later roadmap explicitly
  promotes one of them to a source profile.

This keeps the public contract incremental: UWS gains stronger source bindings
for well-known API families while avoiding a generic "anything captured from a
tool becomes a workflow contract" model.

### Source-Type Notes

GraphQL needs selectors for a named or embedded query, mutation, or
subscription and request bindings for variables. The GraphQL schema owns object
types, input types, enum values, nullability, and field selection constraints;
UWS owns workflow structure, variable values, dependencies, and output routing.

OpenRPC needs selectors for JSON-RPC method names or refs. The OpenRPC
document owns method params, result schemas, errors, servers, and security
metadata when available; UWS owns param values and workflow orchestration.

gRPC/protobuf needs selectors for service/method names or descriptors. The
protobuf descriptor owns services, methods, messages, streaming shape, and field
types; UWS owns request values, step dependencies, and outputs. Runtime
implementations own channel credentials, transport configuration, and invocation.

OData needs selectors for entity-set operations, functions, and actions.
The OData metadata document owns entity types, navigation properties, functions,
actions, and query semantics; UWS owns selected action/query inputs and workflow
structure.

Browser and network capture should record provenance and review evidence only.
It may help produce or validate an API source document, but UWS should not treat
captured requests as a provider contract unless a future profile format defines
portable target, action, safety, expiry, and revalidation semantics.

## Adopted in UWS 1.3: AsyncAPI Source Profiles

AsyncAPI covers a class of provider contracts that OpenAPI does not model well:
events, subscriptions, message channels, and asynchronous request/reply flows.

The design goal is to keep the same source-document boundary:

- AsyncAPI owns channels, messages, protocols, bindings, payload schemas, and
  security.
- UWS owns which message operation participates in a workflow, how data flows
  between steps, what trigger starts execution, and how structural control flow
  proceeds.

UWS 1.3 source description:

```yaml
sourceDescriptions:
  - name: billing_events
    type: asyncapi
    url: ./asyncapi/billing-events.yaml

operations:
  - operationId: wait_for_invoice_paid
    sourceDescription: billing_events
    sourceOperationId: invoicePaid
```

AsyncAPI is not collapsed into a generic runtime function call. It is a
published provider contract, so UWS treats it as a source profile rather than an
extension-owned runtime action.

### AsyncAPI Graduation Result

AsyncAPI graduated to UWS 1.3 with the following settled scope:

- Source ownership remains clean: AsyncAPI owns channels, messages, protocol
  bindings, payload schemas, and security; UWS does not duplicate them.
- UWS operations select AsyncAPI work with existing generic selectors:
  `sourceOperationId` for root AsyncAPI Operation Object keys and
  `sourceOperationRef` for JSON Pointer fragments.
- `sourceOperationRef` is validated by UWS core as a JSON Pointer fragment.
  Source-aware tooling and runtimes validate whether it resolves to a supported
  AsyncAPI target. `#/operations/<name>` is preferred; `#/channels/<name>` and
  `#/channels/<name>/messages/<name>` are compatibility targets when runtimes
  can resolve them unambiguously.
- Existing trigger and `await` constructs remain the structural UWS way to model
  event entry and waits; protocol-level subscription behavior remains runtime
  and source owned.
- Security and credential handling remain source-derived or runtime-private, not
  embedded in UWS core fields.
- Validation rules and conformance fixtures cover source description types,
  selector exclusivity, and invalid mixed OpenAPI/AsyncAPI binding shapes.

## Candidate v0.1.4: Browser Capability Profiles

Browser capability profiles address a different problem: many useful workflow
targets expose a web UI before they expose a complete API specification. A
browser agent can often complete the task, but repeated LLM-driven browsing is
expensive, difficult to review, and risky for side-effectful actions.

The proposed pattern is:

1. First run: an LLM-assisted browser session explores the target UI.
2. The system produces a browser capability profile from the observed DOM,
   accessibility tree, screenshots, and successful/failed actions.
3. A human or policy gate reviews the profile before trusted execution.
4. Later runs use the reviewed profile instead of repeatedly asking the LLM to
   rediscover the same UI.
5. If the observed UI no longer matches the profile, execution stops or falls
   back to a new review flow.

This separates understanding from commitment. The LLM may help understand an
unknown UI, but the committed workflow should execute against a reviewed,
bounded capability description.

### Browser Profile Graduation Criteria

Browser capability profiles should graduate from roadmap candidate to UWS source
profile only after the following design points are settled:

- The profile format is portable across runtimes and records reviewed evidence,
  not a runtime-specific browser command trace.
- Target resolution semantics are defined well enough for runtimes to detect
  missing, duplicated, expired, renamed, or semantically changed targets.
- Side-effect and confirmation policies are explicit and can be reviewed before
  execution.
- Expiration and revalidation rules are part of the profile contract.
- The profile does not carry secrets, credentials, cookies, session state, or
  product-private runtime configuration.
- Multiple runtimes need interoperable browser-capability interchange rather
  than product-owned extension fields.

## Standards Relationship

There is no single standard called a "Vision/DOM Parser" for browser agents.
The practical browser-observation layer is assembled from existing standards and
runtime-specific protocols:

- DOM and HTML define the page document model.
- WAI-ARIA, accessible-name computation, and accessibility API mappings define
  the semantic information that often becomes a useful browser-agent snapshot.
- WebDriver and WebDriver BiDi define standard browser automation protocols.
- Chrome DevTools Protocol, Playwright, and Puppeteer expose practical
  implementation APIs, but their LLM-facing snapshot formats are not UWS
  standards.

UWS should therefore avoid standardizing a specific parser output. The portable
unit should be a reviewed capability profile that records the observation method
and the stable UI targets a runtime may use.

Useful references:

- DOM Standard: https://dom.spec.whatwg.org/
- WAI-ARIA: https://www.w3.org/TR/wai-aria/
- Accessible Name and Description Computation: https://w3c.github.io/accname/
- WebDriver: https://w3c.github.io/webdriver/
- WebDriver BiDi: https://www.w3.org/TR/webdriver-bidi/
- Chrome DevTools Protocol: https://chromedevtools.github.io/devtools-protocol/

## Browser Profile Scope

A browser capability profile is learned evidence, not a provider API contract.
It should be weaker than OpenAPI, Discovery, Smithy, or AsyncAPI.

Use a browser profile when:

- no sufficient API or event source document is available
- the task is only exposed through a web UI
- the UI has stable accessible roles, names, labels, URLs, or selectors
- the profile can be reviewed before side-effectful execution

Do not use a browser profile when:

- an official API contract exists and can support the task
- the browser flow depends on brittle visual guessing only
- the target action cannot be bounded or confirmed
- the action involves regulated or high-risk side effects without stronger
  runtime controls

## Profile Shape

A future profile should record at least:

| Field | Purpose |
| --- | --- |
| `provider` | Human-readable provider or app name. |
| `origin` | Browser origin or allowed origin list. |
| `loginStateRequired` | Whether execution assumes an authenticated browser session. |
| `observationKind` | `accessibility_snapshot`, `dom_text`, `screenshot_ocr`, or another declared method. |
| `actions` | Named reusable UI capabilities exposed by the profile. |
| `inputs` | Required action inputs and validation rules. |
| `outputs` | Observable result fields or completion evidence. |
| `sideEffects` | Declared side effects such as `sends_email`, `creates_record`, or `deletes_resource`. |
| `confirmationPolicy` | Whether runtime confirmation is required before commitment. |
| `targets` | Stable refs, roles, names, selectors, or other runtime-resolvable UI targets. |
| `evidence` | When and how the profile was learned or last validated. |
| `confidence` | `low`, `medium`, or `high` based on target stability and verification. |
| `expiresAfter` | Revalidation window. |

Illustrative future profile:

```yaml
profile: uws.browser.0.1
provider: jinba
origin: https://flow.jinba.io
loginStateRequired: true
observationKind: accessibility_snapshot
actions:
  send_gmail_report:
    description: Run a prepared workflow that sends a report through Gmail.
    inputs:
      recipient_email:
        type: string
        format: email
      report_body:
        type: string
    sideEffects:
      - sends_email
    confirmationPolicy:
      required: true
      prompt: Confirm sending this report through Gmail.
    targets:
      - role: textbox
        name: Recipient email
      - role: button
        name: Run
    outputs:
      sent:
        evidence: success_toast_or_run_status
    evidence:
      learnedAt: 2026-05-25T00:00:00Z
      source: reviewed_browser_session
    confidence: medium
    expiresAfter: P30D
```

## Possible UWS Binding

A future UWS version could model browser profiles as source descriptions. This
shape is not valid UWS 1.4:

```yaml
sourceDescriptions:
  - name: jinba_ui
    type: browser-profile
    url: ./browser-profiles/jinba.yaml

operations:
  - operationId: send_weather_report_with_jinba
    sourceDescription: jinba_ui
    sourceOperationId: send_gmail_report
    request:
      body:
        recipient_email: $inputs.recipient_email
        report_body: $steps.render_report.outputs.body
```

This mirrors the source-document pattern: the profile owns UI targets and action
capabilities; UWS owns operation selection, request values, dependencies, and
workflow structure.

## Valid Today: Extension-Owned Browser Work

Before browser profiles become a source type, current UWS 1.4 documents
represent browser work as extension-owned operations:

```yaml
operationId: send_weather_report_with_browser
x-uws-operation-profile: acme.browser.1.0
x-acme-browser:
  profile: ./browser-profiles/jinba.yaml
  action: send_gmail_report
request:
  body:
    recipient_email: $inputs.recipient_email
    report_body: $steps.render_report.outputs.body
```

This form is valid today because the operation has no source binding. It uses
`x-uws-operation-profile` to name the implementation profile, and product-owned
`x-*` fields carry the browser profile path and action name.

Extension-owned browser operations must not set `sourceDescription`,
`sourceOperationId`, `sourceOperationRef`, `openapiOperationId`, or
`openapiOperationRef`. Those fields remain reserved for source-bound operations.
OpenUdon-like authoring tools should use this extension-owned shape today and
must not emit `sourceDescriptions[].type: browser-profile` until a future UWS
version adopts that source type.

The future source-profile form becomes useful only after multiple runtimes need
portable browser-profile interchange.

## Safety Requirements

Browser profiles must be stricter than normal read-only scraping flows because
browser actions commonly create side effects.

Minimum future requirements:

- Profiles must declare an origin allowlist and runtime execution must fail
  closed outside those origins.
- Profiles must declare whether login state is required, but must not contain
  credentials, session cookies, tokens, or raw login state.
- Profiles must declare the observation method used to learn or validate the UI.
- Profiles must declare named actions, action inputs, observable outputs, stable
  targets, side effects, confirmation policy, review evidence, confidence, and
  expiry.
- Side-effectful actions must declare side effects explicitly.
- Send, create, update, delete, post, upload, purchase, approve, invite, and
  notify actions must require explicit confirmation unless a trusted policy
  has pre-approved the exact action.
- Runtime execution must fail closed when stable targets are missing,
  duplicated, renamed, expired, visually inconsistent, or semantically
  inconsistent.
- Profiles must expire and require revalidation.
- Profiles should retain evidence about the observation method, last review,
  and last successful validation.

## Non-Goals

UWS should not become a browser automation language. Browser capability profiles
are intended to let a workflow reference reviewed UI capabilities, not to encode
browser mechanics directly in UWS.

Future UWS source-profile work should not standardize:

- DOM snapshot, accessibility snapshot, screenshot OCR, CSS selector, XPath, or
  computer-vision output formats.
- Playwright, Puppeteer, Chrome DevTools Protocol, WebDriver, or WebDriver BiDi
  command streams.
- Login flows, credential storage, session management, captcha handling, or
  provider-specific browser client configuration.
- Repeated live LLM browsing loops as the committed execution path.
- Browser profiles as a substitute for stronger provider contracts when
  OpenAPI, Google Discovery, AWS Smithy, AsyncAPI, or another reviewed API/event
  source can support the task.

## Compatibility Goal

The compatibility goal is not to make UWS a browser automation language. The
goal is to let UWS reference reviewed browser capabilities in the same disciplined
way it references API operations:

- published API or event source when available
- AsyncAPI for event/message contracts
- browser capability profile only when the workflow target is UI-only or API
  coverage is insufficient
- extension-owned operation profiles while the browser profile contract is still
  experimental
