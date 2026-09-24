# Archive B01 - Browser Capability And Account Lifecycle Profiles

**Context.** Separately versioned, fail-closed browser capability,
authentication, registration, and private registration-input contracts.

**Baseline.** 8382d0f26b3b10870125760643078d1a1a3e31b6

**Coverage.** verified

**Supersedes.** none

## Scope And Responsibilities

The browser family is deliberately separate from UWS core. Browser 1.7 is the
current capability profile and binds through the first-class `browser-profile`
source type. Browser authentication 1.1 and its call supplement describe
reviewed sign-in and named-session establishment. Browser registration 1.2 and
its call supplement describe one explicitly approved account-creation attempt
with reviewed human-verification authority. Registration input 1.0 defines the
owner-private input envelope used by registration 1.1 and 1.2 flows.

Earlier browser 1.5/1.6, authentication 1.0, and registration 1.0/1.1 artifacts
remain accepted. The profiles do not standardize a browser driver, live
session, credential store, capture format, discovery system, or registry.

## Domain And Workflows

Browser capability profiles are author-asserted evidence of actions implemented
by a target UI. They declare exact origins, observation and review metadata,
bounded page/frame contexts, accessibility-first locators, a closed macro
vocabulary, typed outputs, side effects, and confirmation policy. Browser 1.7
adds locale-free scalar conversion for accessibility-text outputs while
retaining fail-closed ambiguity and conversion rules.

Authentication profiles contain symbolic credential slots, explicit sign-in
steps, selected MFA handoffs, and an exact success condition. Credential values,
cookies, OAuth state, and session handles remain private. Call envelopes select
one reviewed flow and give the resulting execution-local session a symbolic
name; expired sessions do not cause implicit reauthentication or retry.

Registration profiles contain symbolic credentials and typed input slots,
explicit checkpoints, exactly one approved submit, duplicate prevention,
ambiguous-outcome handling, cleanup disposition, and a success condition.
Version 1.2 requires a reviewed provider verification policy with bounded
network, frame, request, byte, and timeout authority. Human challenges remain
human-operated. Filled input envelopes bind the exact public profile digest,
selected flow, and monotonically increasing private revision and must remain
outside packages, Git, prompts, logs, and ordinary diagnostics.

## System Shape

The normative JSON and Markdown profile documents live under `versions/`.
`browserauthentication` and `browserregistration` provide inert Go wire types
and extension encoding helpers. `schemas` embeds profile schemas, dispatches on
exact discriminators, validates profile and call documents, applies additional
context/origin and registration-binding rules, creates empty private-input
templates, and validates checkpoint updates.

Browser capability operations use ordinary core source bindings. Authentication
and registration calls are extension-owned operations selected through exact
`x-uws-operation-profile` values. The runtime binds credentials, inputs,
sessions, approvals, browser state, verification adapters, deadlines, and
persistent attempt identity only at execution time.

## Contracts And Dependencies

Profile schemas are closed and versioned independently. Validators reject
unknown discriminators rather than guessing. Newer features require explicit
selection: browser 1.7 requires UWS 1.9, context-aware authentication requires
1.1, typed registration inputs require registration/call 1.1, and reviewed
verification requires registration/call 1.2.

Origin matching is exact, unsafe locator forms are excluded, context graphs are
bounded and acyclic, popup bindings are explicit, and side-effectful operations
require confirmation. Registration call/profile/input bindings are checked
together so a consumer cannot strip newer fields or silently downgrade review
authority.

## Operations And Verification

Schema and semantic tests cover canonical fixtures, closed macro vocabulary,
unsafe outputs and side effects, exact origins and contexts, union decoding,
unknown fields, typed private-input updates, revision and checkpoint rules,
profile/call mismatches, human-verification budgets, and downgrade rejection.
The full test and vet suites pass at this baseline.

The profiles consistently fail closed on ambiguity, stale or unsupported
versions, missing bindings, unsafe origins, unknown challenges, invalid private
input, and indeterminate submission outcomes. They do not authorize automatic
retry of authentication or registration mutations.

## Evidence

| Claim | Repository Evidence |
|---|---|
| Browser 1.7 is a separately versioned, author-asserted capability contract. | `versions/browser.1.7.md`; `versions/browser.1.7.json`; `docs/future-source-profiles.md` |
| Authentication 1.1 keeps credentials and sessions runtime-private. | `versions/browser-authentication.1.1.md`; `versions/browser-authentication-call.1.1.md`; `browserauthentication/browser_authentication.go` |
| Registration 1.2 adds bounded reviewed verification without weakening earlier one-attempt controls. | `versions/browser-registration.1.2.md`; `versions/browser-registration-call.1.2.md`; `browserregistration/verification.go` |
| Registration input values are digest-bound private snapshots with staged update rules. | `versions/browser-registration-input.1.0.md`; `browserregistration/private_inputs.go`; `schemas/registration_inputs.go` |
| Profile schema selection and cross-document checks fail closed. | `schemas/schema.go`; `schemas/registration_inputs.go`; `schemas/registration_verification.go` |
| Canonical and adversarial fixtures cover the profile family. | `testdata/browser-profile/`; `testdata/browser-authentication/`; `testdata/browser-registration/`; `schemas/*browser*_test.go`; `schemas/registration_*_test.go` |

## Observed Gaps

The repository's current artifacts include registration profile/call 1.2, but
`README.md`, `AGENTS.md`, and the MkDocs reference navigation primarily surface
registration 1.1. `docs/registration-authoring.md` is also omitted from the
configured navigation, as reported by the strict documentation build.

`versions/browser.1.7.md` calls `versions/1.9.0.json` the current core schema
even though UWS 1.9.2 is current at this baseline. Historical downstream-suite
claims in `docs/browser-capability-goal.md` cannot be independently rechecked
without the Browsertools, OpenUdon, and Udon repositories.
