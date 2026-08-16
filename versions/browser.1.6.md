# UWS Browser Capability Profile 1.6

The UWS Browser Capability Profile is the sub-spec for documents referenced by
UWS 1.8 `sourceDescriptions[].type: browser-profile`. It is wire/spec metadata
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
| Profile name | `uws.browser.1.6` |
| JSON Schema | `versions/browser.1.6.json` |
| Bound from | UWS 1.8+ `sourceDescriptions[].type: browser-profile` |
| Selector | `sourceOperationId` (action name) or `sourceOperationRef` (preferred form: `#/actions/<name>`) |

A minimal profile document:

```yaml
profile: uws.browser.1.6
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
| `profile` | string const | REQUIRED. Literal `uws.browser.1.6`. |
| `info.title` | string | REQUIRED. Human-readable provider/app name. |
| `info.origin` | string or array of strings | REQUIRED. Browser origin allowlist; the runtime MUST reject `navigate` outside it. A profile with multiple origins MUST use absolute `navigate` targets because no implicit base origin is defined. |
| `info.loginStateRequired` | boolean | Whether execution assumes an authenticated browser session. |
| `observationKind` | enum | REQUIRED. How the profile was learned: `accessibility_snapshot`, `dom_text`, `screenshot_ocr`, `other`. The snapshot format itself stays runtime-owned. |
| `evidence` | object | REQUIRED. `learnedAt` (RFC 3339 timestamp) plus optional `source` (e.g. `reviewed_browser_session`). |
| `confidence` | enum | REQUIRED. `low` \| `medium` \| `high`. Target-stability assessment from review. |
| `expiresAfter` | ISO-8601 duration | REQUIRED. Revalidation window with at least one duration component (e.g. `P30D`). |
| `verification` | object | REQUIRED. `lastVerifiedAt`, `successfulRuns`, and optional `uiStabilityScore` ([0,1]). |
| `contexts` | object | Optional bounded popup/frame graph. `main` is implicit and cannot be redeclared. |
| `actions` | object | REQUIRED. Map from action name (matched by `sourceOperationId`) to action definition. |

## Portable Contexts

Browser 1.6 adds portable popup and frame relationships while retaining every
browser 1.5 main-page action shape. Omitted `context` means the implicit
`main` page.

```yaml
contexts:
  idp_popup:
    kind: popup
    parent: main
    origin: https://login.example.com
  login_frame:
    kind: frame
    parent: main
    origin: https://login.example.com
    path: /embedded/login
    name: Login
```

Context IDs are bounded identifiers. The graph MUST be acyclic and at most
four contexts deep. Every context origin MUST be an exact origin in
`info.origin`. A frame declares an exact clean path or non-empty frame name and
MUST resolve uniquely at runtime. A popup declares neither path nor name and
MUST be opened by exactly one explicit click whose `opensContext` names it.
The popup's declared `parent` MUST equal that click's execution context.
Automatic, missing, duplicate, or multiple popups fail closed.

## Locators (Accessibility-First)

Locators are objects in the form `{ role, name?, text?, value? }`. They are
**accessibility-only**: runtimes resolve them by querying the browser's
accessibility tree. CSS selectors, XPath, DOM-structure paths, and pixel /
coordinate references are NOT permitted as locator forms.

| Field | Type | Purpose |
| --- | --- | --- |
| `role` | enum | REQUIRED. One of `button`, `link`, `textbox`, `checkbox`, `radio`, `dialog`, `status`, `alert`, `heading`, `img`, `list`, `listitem`, `combobox`, `option`, `menu`, `menuitem`, `tab`, `tabpanel`, `table`, `row`, `cell`, `region`, `navigation`, `article`, `form`, `search`, `switch`, `group`. |
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
| `navigate` | string or `{ url, context? }` | Relative or absolute URL. MUST be templatable from action `parameters` via `{{name}}`. Runtime MUST reject targets outside `info.origin`, context/origin mismatches, and a relative target when `info.origin` contains more than one origin. |
| `click` | `{ locator, wait_for?, context?, opensContext? }` | Fire a standard click event on the located element and optionally bind exactly one declared popup. |
| `type_text` | `{ locator, value, wait_for?, context? }` | Input text. `value` MAY be templated. |
| `check_radio` | `{ locator, wait_for?, context? }` | Check a radio or checkbox element. |
| `uncheck` | `{ locator, wait_for?, context? }` | Uncheck a checkbox or toggle element. |
| `select_option` | `{ locator, value, wait_for?, context? }` | Pick a value option from a drop-down. |
| `wait_for` | locator, `{ locator, context }`, or `{ navigation }` | Halt until a locator condition in one context or navigation event (`load`, `domcontentloaded`, `network_idle`) fires. |

The set is **closed**: additions or removals require a `browser.x.y` schema
bump (and a corresponding UWS minor-version bump). Vendor-specific extension
fields under `x-` are NOT permitted inside the sequence — that would re-open
the "browser automation language" trap UWS is explicitly avoiding.

## Outputs

Output extraction is declared in an action's `outputs` map. Each output names a
field and how to extract it:

| Field | Type | Purpose |
| --- | --- | --- |
| `type` | enum | REQUIRED. Declared JSON type: `string`, `integer`, `number`, `boolean`, `array`, `object`, or `null`. |
| `source` | enum | REQUIRED. `a11y` \| `jsonld` \| `microdata` \| `css`. |
| `locator` | locator object | REQUIRED only when `source: a11y`; forbidden for other sources. |
| `context` | context ID | Optional execution context for the extraction; omission means `main`. |
| `presence` | boolean | When `true` and `source: a11y`, the runtime returns `true`/`false` based on whether the locator resolves rather than extracting accessible text. Only valid with `source: a11y`. |
| `property` | string | The itemprop (microdata) or JSON-LD property name to extract. When absent, the output key name is used. Only valid with `source: microdata` or `source: jsonld`. |
| `attribute` | string | DOM attribute name to extract from the matched element (e.g. `href`, `src`, `datetime`) instead of its text content. Only valid with `source: css`. |
| `selector` | string | REQUIRED only when `source: css`; forbidden for other sources. |
| `fallbackReason` | enum | REQUIRED only when `source: css`; forbidden for other sources. One of `no_a11y_region`, `no_structured_data`, `ambiguous_a11y`, `other`. |
| `validation` | inline JSON Schema | REQUIRED when `source: css`; optional constraints for other extracted values (e.g. `enum`, `pattern`). |

`a11y`, `jsonld`, and `microdata` are the **primary** output methods. `css` is
a **tightly-constrained fallback** that MUST: (a) be paired with a typed
`validation` schema; (b) record a `fallbackReason`; (c) be surfaced as
lower-confidence evidence in runtime reports. CSS is permitted for outputs only
— never as a *locator*.

When `presence: true` is set on an `a11y` output, `type` MUST be `boolean`.
Runtimes MUST NOT extract accessible text for presence-typed outputs — the sole
result is whether the locator matched at least one accessibility-tree element.

When `property` is set on a `microdata` or `jsonld` output, it names the
itemprop attribute value or JSON-LD property key to extract. When absent,
tooling SHOULD match the output key name to an itemprop or property of the same
name; runtimes that cannot match fall back to returning `null`.

When `attribute` is set on a `css` output, the runtime extracts the named DOM
attribute of the matched element (for example `href` from a link or `src` from
an image) instead of the element's text content. `attribute` is valid only with
`source: css`; the output still MUST carry `fallbackReason` and a typed
`validation` schema like any other css output. Attribute extraction does not
apply to `a11y`, `jsonld`, or `microdata` outputs — accessible state is exposed
through the locator's `value`, and structured-data values through `property`.

## Parameters and Default Substitution

Each action's `parameters` field is an inline JSON Schema describing the
templatable inputs. Parameter values are substituted into sequence steps via
`{{paramName}}` placeholders.

If a parameter carries a JSON Schema `default`, runtimes MUST substitute that
default into `{{paramName}}` placeholders when the caller omits the parameter.
A parameter that is neither supplied by the caller nor has a `default` MUST
cause the runtime to fail the action with a missing-parameter error before any
sequence steps execute.

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

The following are NOT part of `browser.1.6` and MUST NOT appear in the
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

Runtimes bind active sessions at execution time. Portable `contexts` describe
page/frame topology only; they never contain a live browser context, cookie,
storage state, or resumption handle. Profiles are *learned evidence* and stay
inert until a trusted runtime executes them with bound credentials.

## What Stays Distribution-Owned

Content-addressed distribution does not change the portable profile wire. A
catalog, package, or registry MAY bind a profile to external metadata such as
an artifact digest, catalog identity and release version, media type, byte
size, provenance, license, publication time, expiry assessment, revocation, or
supersession. That metadata is an envelope around the exact profile bytes and
MUST NOT be inserted into this closed schema.

Likewise, registry URLs, storage paths, pull-request authors, reviewer
identities, signatures, membership, access policy, and transport credentials
remain distribution concerns. Consumers validate the profile against this
schema, verify the external digest/lifecycle envelope according to their own
policy, and bind sessions only at trusted execution time. A profile copied
between a local package and a static registry therefore retains identical
`uws.browser.1.6` bytes.

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
`sourceOperationId` / `sourceOperationRef` per UWS 1.8.

## Compatibility

`uws.browser.1.5` remains immutable and accepted. A main-only profile that does
not need context-qualified steps or outputs SHOULD continue to use browser 1.5
and may be bound from UWS 1.7. Browser 1.6 requires UWS 1.8. Browser-aware
validators select the schema from the document's exact `profile`
discriminator and reject unknown versions.
