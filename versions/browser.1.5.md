# UWS Browser Capability Profile 1.5

The UWS Browser Capability Profile is the sub-spec for documents referenced by
UWS 1.5 `sourceDescriptions[].type: browser-profile`. It is wire/spec metadata
only. It does not standardize browser execution, session management,
credentials, rendering, or any specific browser automation protocol
(WebDriver, WebDriver BiDi, Chrome DevTools Protocol, Playwright, Puppeteer).
Those remain runtime-private.

UWS core (the main `versions/1.5.0.json` schema) only references this profile
by *type name* and reuses the existing generic `sourceOperationId` /
`sourceOperationRef` selector rules. Validating a profile document against the
constraints below — locator vocabulary, declarative action vocabulary, output
extraction methods, safety controls — is the job of browser-aware tooling and
runtimes.

`browser-profile` is the **weakest** first-class source type by intent. It
records *learned evidence* of a UI, not a vendor-published API contract.
Tooling SHOULD prefer a published API source (`openapi`, `google-discovery`,
`aws-smithy`, `asyncapi`, `graphql`, `openrpc`, `grpc-protobuf`, `odata`) when
one covers the task. See [`docs/future-source-profiles.md`](../docs/future-source-profiles.md)
"Browser Profile Scope" for the use/don't-use guidance.

## Profile

| Field | Value |
| --- | --- |
| Profile name | `uws.browser.1.5` |
| JSON Schema | `versions/browser.1.5.json` |
| Bound from | UWS 1.5+ `sourceDescriptions[].type: browser-profile` |
| Selector | `sourceOperationId` (action name) or `sourceOperationRef` (preferred form: `#/actions/<name>`) |

A minimal profile document:

```yaml
profile: uws.browser.1.5
info:
  title: GitHub PR Manager
  origin: https://github.com
  loginStateRequired: true
observationKind: accessibility_snapshot
evidence:
  learnedAt: "2026-05-30T00:00:00Z"
  source: reviewed_browser_session
confidence: high
expiresAfter: P30D
verification:
  lastVerifiedAt: "2026-05-30T00:00:00Z"
  successfulRuns: 1
  uiStabilityScore: 0.95
actions:
  approve_pr:
    description: Approves a Pull Request on GitHub.
    parameters:
      type: object
      required: [owner, repo, pr_number, comment]
      properties:
        owner: { type: string }
        repo: { type: string }
        pr_number: { type: integer }
        comment: { type: string }
    sequence:
      - navigate: "/{{owner}}/{{repo}}/pull/{{pr_number}}"
      - click:
          locator: { role: button, name: "Review changes" }
      - type_text:
          locator: { role: textbox, name: "Submit review comment" }
          value: "{{comment}}"
      - check_radio:
          locator: { role: radio, name: "Approve" }
      - click:
          locator: { role: button, name: "Submit review" }
          wait_for: { role: status, name: "Approved" }
    sideEffects: [state_change]
    confirmationPolicy:
      required: true
      prompt: "Approve PR #{{pr_number}} on {{owner}}/{{repo}}?"
```

## Profile Fields

| Field | Type | Purpose |
| --- | --- | --- |
| `profile` | string const | REQUIRED. Literal `uws.browser.1.5`. |
| `info.title` | string | REQUIRED. Human-readable provider/app name. |
| `info.origin` | string or array of strings | REQUIRED. Browser origin allowlist; the runtime MUST reject `navigate` outside it. |
| `info.loginStateRequired` | boolean | Whether execution assumes an authenticated browser session. |
| `observationKind` | enum | REQUIRED. How the profile was learned: `accessibility_snapshot`, `dom_text`, `screenshot_ocr`, `other`. The snapshot format itself stays runtime-owned. |
| `evidence` | object | REQUIRED. `learnedAt` (RFC 3339 timestamp) plus optional `source` (e.g. `reviewed_browser_session`). |
| `confidence` | enum | REQUIRED. `low` \| `medium` \| `high`. Target-stability assessment from review. |
| `expiresAfter` | ISO-8601 duration | REQUIRED. Revalidation window with at least one duration component (e.g. `P30D`). |
| `verification` | object | REQUIRED. `lastVerifiedAt`, `successfulRuns`, and optional `uiStabilityScore` ([0,1]). |
| `actions` | object | REQUIRED. Map from action name (matched by `sourceOperationId`) to action definition. |

## Locators (Accessibility-First)

Locators are objects in the form `{ role, name?, text?, value? }`. They are
**accessibility-only**: runtimes resolve them by querying the browser's
accessibility tree. CSS selectors, XPath, DOM-structure paths, and pixel /
coordinate references are NOT permitted as locator forms.

| Field | Type | Purpose |
| --- | --- | --- |
| `role` | enum | REQUIRED. One of `button`, `link`, `textbox`, `checkbox`, `radio`, `dialog`, `status`, `alert`. |
| `name` | string | Accessible name (label / `aria-label` / computed name). |
| `text` | string | Accessible inner text for `status`/`alert` matching. |
| `value` | string | Current value/state (checked, selected, etc.). |

A runtime that resolves a locator to more than one accessibility-tree element
MUST fail closed rather than guess (the "Strict Fail-Closed on Ambiguity"
safety control in the proposal).

## Declarative Macro Actions

The action vocabulary is a **closed enumeration**. Sequences are arrays of
single-key objects; the key selects the macro and its value carries the
arguments.

| Macro | Shape | Purpose |
| --- | --- | --- |
| `navigate` | string | Relative or absolute URL. MUST be templatable from action `parameters` via `{{name}}`. Runtime MUST reject targets outside `info.origin`. |
| `click` | `{ locator, wait_for? }` | Fire a standard click event on the located element. |
| `type_text` | `{ locator, value, wait_for? }` | Input text. `value` MAY be templated. |
| `check_radio` | `{ locator, wait_for? }` | Check a radio or checkbox element. |
| `uncheck` | `{ locator, wait_for? }` | Uncheck a checkbox or toggle element. |
| `select_option` | `{ locator, value, wait_for? }` | Pick a value option from a drop-down. |
| `wait_for` | locator or `{ navigation }` | Halt until a locator condition or navigation event (`load`, `domcontentloaded`, `network_idle`) fires. |

The set is **closed**: additions or removals require a `browser.x.y` schema
bump (and a corresponding UWS minor-version bump). Vendor-specific extension
fields under `x-` are NOT permitted inside the sequence — that would re-open
the "browser automation language" trap UWS is explicitly avoiding.

## Outputs

Output extraction is declared in an action's `outputs` map. Each output names a
field and how to extract it:

| Field | Type | Purpose |
| --- | --- | --- |
| `type` | string | Declared JSON type for the extracted value (string, integer, boolean, …). |
| `source` | enum | REQUIRED. `a11y` \| `jsonld` \| `microdata` \| `css`. |
| `locator` | locator object | REQUIRED only when `source: a11y`; forbidden for other sources. |
| `selector` | string | REQUIRED only when `source: css`; forbidden for other sources. |
| `fallbackReason` | enum | REQUIRED only when `source: css`; forbidden for other sources. One of `no_a11y_region`, `no_structured_data`, `ambiguous_a11y`, `other`. |
| `validation` | inline JSON Schema | REQUIRED when `source: css`; optional constraints for other extracted values (e.g. `enum`, `pattern`). |

`a11y`, `jsonld`, and `microdata` are the **primary** output methods. `css` is
a **tightly-constrained fallback** that MUST: (a) be paired with a typed
`validation` schema; (b) record a `fallbackReason`; (c) be surfaced as
lower-confidence evidence in runtime reports. CSS is permitted for outputs only
— never as a *locator*.

## Side Effects and Confirmation

Each action MUST declare its `sideEffects` (a non-empty array). The vocabulary
is:

`read_only`, `state_change`, `sends_email`, `creates_record`,
`updates_record`, `deletes_resource`.

Each action MUST declare `confirmationPolicy`. `read_only` MUST be the only
side effect when present. Any side-effectful action MUST set
`confirmationPolicy.required: true`. Runtimes MUST honor `confirmationPolicy`
before committing the action's first state-changing step.

## What Stays Runtime-Owned

The following are NOT part of `browser.1.5` and MUST NOT appear in the
profile document:

- Session tokens, cookies, passwords, OAuth state, or any other secrets.
- Login flows, credential resolution, captcha handling, multi-factor prompts.
- Browser binary/version/channel selection, headless mode, viewport sizing,
  user-agent strings, network throttling, geolocation, timezone, locale.
- WebDriver / WebDriver BiDi / Chrome DevTools Protocol / Playwright /
  Puppeteer command streams, snapshots, traces, or video output.
- DOM snapshot schemas, accessibility-snapshot schemas, screenshot-OCR
  output formats, XPath expressions, computer-vision pipelines.
- Retry policy, rate limiting, parallel-tab orchestration, and other runtime
  execution policy beyond what `wait_for` and the UWS workflow's `when` /
  `successCriteria` already express.

Runtimes bind active session contexts (a pre-authenticated browser profile or
cookie jar) at execution time. Profiles are *learned evidence* and stay
inert until a trusted runtime executes them with bound credentials.

## HCL Representation

UWS workflow operations bind to a browser profile through ordinary
`sourceDescriptions`/`operations` blocks; the profile document itself is the
referenced file. There is no separate HCL representation for the profile —
authors write it in YAML or JSON.

```hcl
sourceDescriptions {
  source_description "github_pr_manager" {
    type = "browser-profile"
    url  = "./profiles/github_pr_manager.yaml"
  }
}

operation "submit_approval" {
  sourceDescription  = "github_pr_manager"
  sourceOperationId  = "approve_pr"
  request {
    body = {
      owner     = "$inputs.owner"
      repo      = "$inputs.repo"
      pr_number = "$inputs.pr_number"
      comment   = "$inputs.comment"
    }
  }
}
```

JSON and YAML use normal `sourceDescriptions` / `operations` fields with
`sourceOperationId` / `sourceOperationRef` per UWS 1.5.
