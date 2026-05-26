# Future Source Profiles

This note records candidate source-profile work beyond the current UWS 1.2
contract. It is intentionally roadmap material, not a normative schema change.

The current UWS model remains API-source-first:

- OpenAPI, Google Discovery, and AWS Smithy describe HTTP API contracts.
- UWS describes workflow structure, operation binding, request values, outputs,
  triggers, and control flow.
- Non-HTTP behavior is represented by extension-owned operations and named
  operation profiles.

Two future source families are worth tracking separately:

| Candidate | Theme | Purpose |
| --- | --- | --- |
| `v0.1.3` | AsyncAPI source profiles | Bind UWS workflows to event, message, and subscription contracts when a provider publishes AsyncAPI. |
| `v0.1.4` | Browser capability profiles | Reuse reviewed browser/UI capability maps when no sufficient API source document exists. |

The `v0.1.x` labels are implementation-roadmap labels. They do not change the
current published UWS 1.2 wire contract until a future version explicitly adopts
the feature.

## Candidate v0.1.3: AsyncAPI Source Profiles

AsyncAPI should be the next source-profile candidate because it covers a class
of provider contracts that OpenAPI does not model well: events, subscriptions,
message channels, and asynchronous request/reply flows.

The design goal is to keep the same source-document boundary:

- AsyncAPI owns channels, messages, protocols, bindings, payload schemas, and
  security.
- UWS owns which message operation participates in a workflow, how data flows
  between steps, what trigger starts execution, and how structural control flow
  proceeds.

Possible future source description:

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

AsyncAPI should not be collapsed into a generic runtime function call. It is a
published provider contract, so UWS should treat it as a source profile rather
than an extension-owned runtime action.

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

- no sufficient API source document is available
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

Example profile:

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

A future UWS version could model browser profiles as source descriptions:

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

This mirrors the API-source pattern: the profile owns UI targets and action
capabilities; UWS owns operation selection, request values, dependencies, and
workflow structure.

Before this becomes a source type, current UWS documents can still represent
browser work as extension-owned operations:

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

The extension-owned form is valid today. The future source-profile form would
become useful only after multiple runtimes need portable browser-profile
interchange.

## Safety Requirements

Browser profiles must be stricter than normal read-only scraping flows because
browser actions commonly create side effects.

Recommended requirements:

- Side-effectful actions must declare side effects explicitly.
- Send, create, update, delete, post, upload, purchase, approve, invite, and
  notify actions should require explicit confirmation unless a trusted policy
  has pre-approved the exact action.
- Runtime execution must fail closed when stable targets are missing,
  duplicated, renamed, or visually/semantically inconsistent.
- Profiles should expire and require revalidation.
- Secrets must not be embedded in the profile. Login state and credentials
  remain runtime-private.
- Profiles should retain evidence about the observation method, last review,
  and last successful validation.

## Compatibility Goal

The compatibility goal is not to make UWS a browser automation language. The
goal is to let UWS reference reviewed browser capabilities in the same disciplined
way it references API operations:

- published API source when available
- AsyncAPI for event/message contracts
- browser capability profile only when the workflow target is UI-only or API
  coverage is insufficient
- extension-owned operation profiles while the browser profile contract is still
  experimental

