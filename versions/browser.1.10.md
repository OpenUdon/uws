# UWS Browser Capability Profile 1.10

The UWS Browser Capability Profile is the sub-spec for documents referenced by
UWS 1.9 and later `sourceDescriptions[].type: browser-profile`. It is wire/spec metadata
only. It does not standardize browser execution, session management,
credentials, rendering, or any specific browser automation protocol
(WebDriver, WebDriver BiDi, Chrome DevTools Protocol, Playwright, Puppeteer).
Those remain runtime-private.

UWS 1.9+ core schemas only reference this profile by *type name* and reuse the
existing generic `sourceOperationId` /
`sourceOperationRef` selector rules. UWS 1.11.0 is the current core release;
Browser 1.10 does not require a core schema change. Browser 1.10 retains the
Browser 1.9 capability contract and adds typed CSS selector match counts.
Browser-aware tooling validates profile structure and template references;
runtimes validate concrete inputs, perform substitution, and implement the
count semantics below.

## Conformance Language

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **NOT RECOMMENDED**, **MAY**, and
**OPTIONAL** in this document are to be interpreted as described in BCP 14
(RFC 2119 and RFC 8174) only when they appear in all capitals.

`browser-profile` is the **weakest** first-class source type by intent. It
records *learned evidence* of a UI, not a vendor-published API contract.
Tooling SHOULD prefer a published API source (`openapi`, `google-discovery`,
`aws-smithy`, `asyncapi`, `graphql`, `openrpc`, `grpc-protobuf`, `odata`) when
one covers the task. See [`docs/future-source-profiles.md`](../docs/future-source-profiles.md)
"Browser Profile Scope" for the use/don't-use guidance.

## Profile

| Field | Value |
| --- | --- |
| Profile name | `uws.browser.1.10` |
| JSON Schema | `versions/browser.1.10.json` |
| Bound from | UWS 1.9+ `sourceDescriptions[].type: browser-profile` |
| Selector | `sourceOperationId` (action name) or `sourceOperationRef` (preferred form: `#/actions/<name>`) |

A minimal profile document:

```yaml
profile: uws.browser.1.10
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
        pr_number: { type: integer, minimum: 1, maximum: 9007199254740991 }
        comment: { type: string }
    sequence:
      - navigate: "/{{owner}}/{{repo}}/pull/{{pr_number}}"
      - click:
          locator: { role: button, name: "Review changes" }
      - type_text:
          locator: { role: textbox, name: "Submit review comment" }
          value: "{{{{literal}}}} {{comment}}"
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
| `profile` | string const | REQUIRED. Literal `uws.browser.1.10`. |
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

Browser 1.9 retains Browser 1.8's portable popup/frame and action contracts.
Omitted `context` means the implicit
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
| `navigate` | string or `{ url, context? }` | Relative or absolute URL. Browser 1.9 templates and literal-brace escapes are allowed only in path segments and query values; runtimes component-encode substitutions and escaped braces and reject unsafe dot segments. Runtime MUST reject targets outside `info.origin`, context/origin mismatches, and a relative target when `info.origin` contains more than one origin. |
| `click` | `{ locator, wait_for?, context?, opensContext? }` | Fire a standard click event on the located element and optionally bind exactly one declared popup. |
| `type_text` | `{ locator, value, wait_for?, context? }` | Set the exact value of an editable field. Browser 1.9 permits scalar substitution and literal-brace escapes in `value`; it does not synthesize keypresses or submit forms. |
| `check_radio` | `{ locator, wait_for?, context? }` | Check a radio or checkbox element. |
| `uncheck` | `{ locator, wait_for?, context? }` | Uncheck a checkbox or toggle element. |
| `select_option` | `{ locator, value, wait_for?, context? }` | Pick a value option from a drop-down. Browser 1.9 permits scalar substitution and literal-brace escapes in `value`. |
| `wait_for` | locator, `{ locator, context }`, or `{ navigation }` | Halt until a locator condition in one context or navigation event (`load`, `domcontentloaded`, `network_idle`) fires. |

The set is **closed**: additions or removals require a `browser.x.y` schema
bump. They require a UWS core version change only if the generic core binding
shape changes. Vendor-specific extension
fields under `x-` are NOT permitted inside the sequence — that would re-open
the "browser automation language" trap UWS is explicitly avoiding.

## Outputs

Output extraction is declared in an action's `outputs` map. Each output names a
field and how to extract it:

| Field | Type | Purpose |
| --- | --- | --- |
| `type` | enum | REQUIRED. Declared JSON type. Accessibility-text outputs accept only `string`, `integer`, `number`, or `boolean`; `array`, `object`, and `null` remain available only to the structured-data and constrained CSS extraction modes. |
| `source` | enum | REQUIRED. `a11y` \| `jsonld` \| `microdata` \| `css`. |
| `locator` | locator object | REQUIRED only when `source: a11y`; forbidden for other sources. |
| `context` | context ID | Optional execution context for the extraction; omission means `main`. |
| `presence` | boolean | When `true` and `source: a11y`, the runtime returns `true`/`false` based on whether the locator resolves rather than extracting accessible text. Only valid with `source: a11y`. |
| `property` | string | The itemprop (microdata) or JSON-LD property name to extract. When absent, the output key name is used. Only valid with `source: microdata` or `source: jsonld`. |
| `attribute` | string | DOM attribute name to extract from the matched element (e.g. `href`, `src`, `datetime`) instead of its text content. Only valid with `source: css`. |
| `matchCount` | `true` | When set, return the exact CSS selector match count instead of reading text or an attribute. Requires `type: integer`, `source: css`, and `visibility`. |
| `within` | CSS selector | Optional unique scope root for a count output. It MUST match exactly one connected element in the selected context. The output's `selector` is evaluated against that root's descendants; the root itself is not included. |
| `visibility` | `all` \| `rendered` | REQUIRED for count outputs. Selects whether all matching connected DOM elements or only rendered matches are counted. |
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

For every non-presence `a11y` output, the runtime MUST read the matched
element's accessible text, trim Unicode whitespace at both edges, and convert
the resulting text according to the declared `type`:

| Declared type | Normative conversion |
| --- | --- |
| `string` | Return the trimmed text unchanged. |
| `integer` | Accept only `-?(0|[1-9][0-9]*)`, then return the integer only when it is within JavaScript's safe-integer range (`-9007199254740991` through `9007199254740991`). |
| `number` | Accept only strict JSON-number syntax (`-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?`) and return the IEEE-754 binary64 result only when it is finite. |
| `boolean` | Accept only the exact lowercase text `true` or `false`. |

Empty text and noncanonical forms MUST fail closed. This includes leading `+`,
leading zeroes other than the single digit `0`, commas, currency symbols,
localized numbers, `NaN`, and infinity spellings. Conversion failure MUST use
the runtime's closed invalid-response result and MUST NOT disclose the page
text. Accessibility-text outputs with declared type `array`, `object`, or
`null` are invalid Browser 1.9 documents. Integer accessibility outputs use
the same inclusive safe-integer range required for integer parameters below.

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

### CSS Match-Count Outputs

Browser 1.10 count outputs are selected explicitly with `matchCount: true`.
They MUST use `source: css`, `type: integer`, a non-empty `selector`, a
`fallbackReason`, and a `validation` schema whose `type` is `integer` and whose
`minimum` is at least zero. A declared `maximum`, when present, MUST be no
greater than `9007199254740991`. Runtimes MUST also enforce the inclusive
JavaScript safe-integer range from zero through `9007199254740991`, even when
the validation schema omits `maximum`. A count outside either range fails
closed; runtimes MUST NOT clamp or round it.

The selector counts matching descendants of the selected output `context`;
omitting `context` selects `main`. Without `within`, the scope is the selected
context's document. With `within`, the runtime resolves that non-empty CSS
selector in the selected context and MUST fail closed unless it matches exactly
one connected element. It then counts matching descendants of that root; the
root itself is never included, even if it matches the output selector. Zero,
one, and multiple matches are ordinary exact results.

`visibility: all` counts every connected DOM element matching the selector,
regardless of its CSS visibility or box size. `visibility: rendered` counts a
matching element only when it has at least one client rectangle with positive
width and height, and neither the element nor an ancestor has computed
`display: none`, `visibility: hidden` or `visibility: collapse`, or
`content-visibility: hidden`. Elements below the viewport remain eligible;
viewport intersection, clipping, occlusion by another element, and opacity do
not affect the count. An element with `display: contents` and no own positive-
area client rectangle is not counted in rendered mode.

Malformed selectors, a missing or ambiguous scope root, unsupported count
results, non-integer or negative results, unsafe integers, and results rejected
by the declared validation schema MUST fail closed using the runtime's invalid-
response result. Failure reports MUST NOT include page text or attributes.
Count outputs MUST NOT set `attribute`; CSS remains an output-only fallback and
is never a locator. `within` and `visibility` are invalid unless
`matchCount: true` is present.

The following vectors define portable result behavior; the DOM descriptions
are synthetic and contain no page text:

| Vector | Scope and matching descendants | Mode | Expected result |
| --- | --- | --- | --- |
| `count-zero` | No element matches `.result` | `all` | integer `0` |
| `count-one` | One connected element matches `.result` | `all` | integer `1` |
| `count-many-all` | Three elements match; one has `display: none` | `all` | integer `3` |
| `count-many-rendered` | The same three; two have positive-area boxes, including one below the viewport | `rendered` | integer `2` |
| `count-scope-excludes-root` | Unique `.group` root also matches `.result`; no descendant matches `.result` | `all` with `within: .group` | integer `0` |
| `count-root-missing` | `within` matches no element | either | closed invalid-response |
| `count-root-ambiguous` | `within` matches two elements | either | closed invalid-response |
| `count-invalid` | Selector is malformed or result is negative, fractional, nonnumeric, unsafe, or outside validation bounds | either | closed invalid-response |

The complete synthetic profile example is
[`testdata/browser-profile/count-outputs-1.10.yaml`](../testdata/browser-profile/count-outputs-1.10.yaml).

## Parameters, Safe Substitution, and Text Sinks

Each action's `parameters` field is an inline JSON Schema. A parameter
placeholder MUST use the exact form `{{name}}`, where `name` matches
`[A-Za-z][A-Za-z0-9_-]*` and names a declared scalar property. Nested
substitution, runtime expressions, whitespace in braces, malformed pairs,
and undeclared names are invalid.

Browser 1.9 adds two literal-brace escapes. A scanner MUST recognize escapes
before parameter placeholders, left to right:

| Token | Result in a text sink |
| --- | --- |
| `{{{{` | Literal `{{` |
| `}}}}` | Literal `}}` |
| `{{name}}` | The canonical scalar text for the declared parameter |

Single unmatched braces (`{` or `}`) are ordinary literal characters.
Four-brace escapes do not nest. Any remaining unmatched `{{` or `}}`, invalid
or nested placeholder, or template/escape outside an approved sink MUST fail
validation. Substitution is one pass: inserted parameter values are literal
data and MUST NOT be rescanned as a placeholder or escape.

Templates and brace escapes are permitted only in these positions:

| Sink | Placement | Substitution |
| --- | --- | --- |
| `navigate` URL | Within a path segment, or within a query value after `=` and before the next `&` or `#`. | Percent-encode a parameter's UTF-8 bytes as an RFC 3986 component, preserving only ASCII letters, digits, `-`, `.`, `_`, and `~`; use uppercase hex escapes. A literal-brace escape emits `%7B%7B` or `%7D%7D` respectively. |
| `type_text.value`, `select_option.value` | Any position in the value string. | Insert canonical scalar text without URI or HTML encoding. Apply `{{{{` / `}}}}` escapes to the resulting literal text. |
| `confirmationPolicy.prompt` | Any position in the prompt. | Insert canonical scalar text without URI encoding. Apply brace escapes and render the result as text only, never HTML or Markdown. |

Templates and escapes MUST NOT appear in a URL scheme, authority, query name,
fragment, locator, wait condition, context identifier, output declaration, or
any other field. Query delimiters are static: a placeholder or escape cannot
introduce a key, separator, or additional query parameter. A navigation
escape is allowed only in a path segment or query value and is serialized as
uppercase-percent-encoded URI octets, never raw braces. Malformed or misplaced
tokens MUST fail validation rather than be interpreted differently by
runtimes.

After substitution, runtimes MUST parse the navigation target and enforce the
declared exact-origin allowlist. Relative targets remain invalid when
`info.origin` contains multiple origins. Dot path segments are forbidden,
including literal or percent-encoded `.` / `..` segments and a substituted
path segment whose final value is `.` or `..`. Percent-encoding a slash
prevents it from becoming a path separator.

Every referenced parameter MUST have one scalar JSON Schema type. Supported
types and canonical text are:

| JSON Schema type | Canonical text |
| --- | --- |
| `string` | The exact UTF-8 string, without Unicode normalization. |
| `boolean` | Lowercase JSON token `true` or `false`. |
| `integer` | A JSON integer from `-9007199254740991` through `9007199254740991`, rendered in canonical base-10 form: `0` or an optional minus followed by a non-zero digit and digits. No plus sign or leading zero is allowed. |
| `number` | RFC 8785 JSON Canonicalization Scheme number serialization for a finite IEEE-754 binary64 value. Non-finite or out-of-range values MUST fail. |

The integer bounds are inclusive and are the same safe-integer range used for
accessibility-text integer outputs. Runtimes MUST validate the exact integer
before conversion; rounding, truncation, or accepting a value outside this
range is forbidden. This rule applies to supplied values and schema defaults.

For `type_text.value` and `confirmationPolicy.prompt`, runtimes MUST reject the
resolved text before executing an action if it contains Unicode General
Category `Cc` controls (U+0000–U+001F or U+007F–U+009F), line/paragraph
separators U+2028 or U+2029, or a bidirectional control: U+061C, U+200E,
U+200F, U+202A–U+202E, or U+2066–U+2069. The check applies to static profile
text, defaults, and supplied string values after one-pass substitution. A
runtime MUST NOT silently strip or normalize these characters. These
restrictions prevent control-key injection and multiline or directional
spoofing; they do not change `select_option.value` matching semantics.

For `type_text`, the runtime MUST resolve an editable field and set its value
to the resolved text. It MUST NOT synthesize keydown, keypress, keyup, Enter,
or other keyboard events, and MUST NOT submit a form as part of this macro.
The runtime MAY dispatch ordinary `input` or `change` events needed to expose
the new field value. Site behavior triggered by those events remains governed
by the action's declared `sideEffects` and confirmation policy.

Null, arrays, objects, union/nullable types, and values that do not validate
against the declared parameter schema MUST fail before sequence execution.
Schema `default` values MUST also satisfy their parameter schema and the
safe-integer/text rules above. When a caller omits a parameter, a valid default
is used; if it is absent and the parameter is required, the runtime MUST fail
the action before executing any sequence step.

For example, a string `a/b?c` in `/records/{{id}}` becomes
`/records/a%2Fb%3Fc`; the same string in `?q={{query}}` becomes
`?q=a%2Fb%3Fc`. In `type_text.value` it remains `a/b?c`. The text
`{{{{name}}}}` becomes the literal `{{name}}`; an inserted value that itself
contains `{{other}}` remains literal and is not rescanned. The number `1.0` is
inserted as `1`, and Boolean false as `false`.

Browser 1.5 through 1.8 retain their original substitution behavior. A
consumer MUST select Browser 1.9 or later explicitly to claim this contract.

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

The following are NOT part of `browser.1.10` and MUST NOT appear in the
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
`uws.browser.1.10` bytes.

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
`sourceOperationId` / `sourceOperationRef` per UWS 1.9 and later. See the
[current UWS 1.11.0 specification](1.11.0.md) for the core binding contract.

## Compatibility

`uws.browser.1.5` through `uws.browser.1.9` remain immutable and accepted.
Browser 1.10 is opt-in; consumers MUST select its exact discriminator to claim
the match-count contract. Browser 1.10 also retains Browser 1.9's escaped
template and safe-text behavior. A main-only profile that does not need
context-qualified steps or outputs SHOULD continue to use browser 1.5 and may
be bound from UWS 1.7. Context-aware profiles without typed accessibility
conversion MAY continue to use browser 1.6 with UWS 1.8. Browser 1.7 and later
use the generic browser-profile binding introduced in UWS 1.9; Browser 1.10 is
also usable with the current UWS 1.11.0 core. The concrete browser-driver
protocol remains a downstream runtime implementation detail.
Browser-aware validators select the schema from the document's exact `profile`
discriminator and reject unknown versions. Authentication remains the separate
`uws.browser-authentication.1.1` contract.

See the [Browser profile versioning procedure](../docs/future-source-profiles.md#adding-a-browser-profile-version)
before adding a later profile version. Published profile files remain
immutable; new semantics require a new discriminator and schema.
