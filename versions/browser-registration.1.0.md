# UWS Browser Registration Profile 1.0

## Scope

`uws.browser-registration.1.0` is an additive, inert profile for one reviewed
browser account-creation recipe. It is separate from
`uws.browser-authentication`: registration mutates remote account state, while
authentication establishes a browser session for an existing identity.

The normative JSON Schema is
[`browser-registration.1.0.json`](browser-registration.1.0.json). This document
defines the behavioral requirements that accompany the schema.

## Safety boundary

A profile contains only symbolic credential slots, exact origins, declarative
accessibility locators, fixed macro steps, provenance, freshness, effects, and
a success predicate. It MUST NOT contain an account identifier, password,
verification value or token, cookie, browser storage, request or response
content, page capture, live session handle, or browser protocol transcript.

Credential values are supplied only by a human in the browser or by a trusted
runtime resolving the call's symbolic bindings. A runtime MUST NOT put those
values in prompts, artifacts, logs, reports, command arguments, or workflow
outputs.

Registration is an explicitly approved mutation. Every flow:

- MUST include `creates_account` in `effects`;
- MUST set `confirmationPolicy.required: true`;
- MUST contain exactly one `submit` step; and
- MUST obtain the call's exact approval immediately before executing that
  `submit` step.

The approval covers only the exact profile bytes, selected flow, symbolic
bindings, origin inventory, duplicate policy, ambiguity policy, and cleanup
disposition in the call. It is not reusable after any of those facts change.

## Exact origins and navigation

`info.applicationOrigins` and `info.registrationOrigins` are exact origins,
not URL prefixes. Origins MUST use HTTPS, except that loopback HTTP is allowed
for provider-free conformance tests. User information, fragments, paths,
queries, and non-default empty ports are not origins.

Every `navigate` URL and the `success.origin` MUST use a declared origin.
Navigation is absolute. A redirect to an undeclared origin fails closed before
the runtime follows or interacts with it. Adding an origin requires a new
human review and approval.

## Credential slots

`credentialSlots` is a map from symbolic slot names to one of two kinds:
`identifier` or `password`. Slot names are portable metadata; values are not.
Every `type_credential.slot` MUST resolve to a declared slot. A registration
profile has no slot for email links, OTP values, CAPTCHA responses, cookies, or
session state.

## Steps

Each sequence item selects exactly one closed macro:

- `navigate`: open one absolute URL on a declared origin;
- `type_credential`: fill an accessibility locator from a symbolic slot;
- `click`: activate a reviewed non-submission control;
- `submit`: activate the single account-creation control after approval;
- `human_checkpoint`: pause for normal human handling of CAPTCHA, email
  verification, MFA, consent, or another control; or
- `wait_for`: wait for one reviewed accessibility locator.

CSS, XPath, coordinates, arbitrary JavaScript, ordinary literal text input,
uploads, downloads, and automatic challenge solving are absent. `click` MUST
NOT be used to disguise an account-creation submission; that interaction is
represented only by `submit`.

A `human_checkpoint` never carries the human response. The runtime MUST pause,
must not infer or automate the response, and may continue only after the human
has completed the ordinary control and explicitly resumes. CAPTCHA, email
verification, MFA, consent, and rate limits MUST NOT be bypassed. An
unsupported control or rate limit fails the run closed.

`requires_human_verification` MUST appear in `effects` exactly when the
sequence contains a `human_checkpoint`. `sends_verification` declares that the
reviewed registration sends an email, SMS, or comparable verification notice.

## Exactly-once outcome

The registration call fixes duplicate prevention to operator attestation,
duplicate handling to failure, and ambiguous outcomes to stop without retry.
Neither the profile nor a runtime provides an automatic retry policy. After a
timeout, transport loss, unexpected redirect, duplicate response, or other
ambiguous post-submit outcome, the operator reconciles the exact account state
privately before any separately approved future action.

Before execution, the operator MUST attest outside the portable artifacts that
the exact candidate is a dedicated test identity and is not already registered.
The operator MUST also select one call-level cleanup disposition: separately
delete the exact account, or retain it as a dedicated test identity. Actual
cleanup is a separately approved operation and is not performed by this
profile.

## Success and evidence

`success` proves only that the reviewed registration outcome is present at the
declared origin and optional exact clean path. It does not expose an account
identifier, page value, verification token, cookie, or session. A successful
registration profile does not establish a reusable browser session; a later
sign-in uses `uws.browser-authentication`.

`evidence` and `verification` are authoring metadata. Registration profiles do
not carry a run or review counter: neither repeated account creation nor a
self-asserted review count is portable evidence. Source labels and descriptions
MUST remain non-secret and free of personal information. Raw captures and page
content remain private and outside the profile.

## Runtime-owned behavior

The trusted runtime owns credential resolution, browser implementation,
network containment, redirect interception, approval presentation, private
human interaction, deadlines, and value-free execution evidence. Registration
support is optional. A runtime that does not implement this exact profile and
call version MUST reject it without executing any step.
