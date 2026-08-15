# UWS Browser Authentication Profile 1.0

`uws.browser-authentication.1.0` is an additive, inert profile for reviewed
browser sign-in recipes. It does not revise `uws.browser.1.5`; an existing
browser capability profile and workflow remain valid and unchanged.

The profile can describe navigation to an application and identity provider,
credential entry through symbolic slots, an explicit MFA handoff, and a
reviewed success condition. A trusted private driver resolves credentials and
retains the resulting browser session.

## Boundaries

The document MUST NOT contain usernames, passwords, OTP seeds or values,
cookies, local/session storage, OAuth state, screenshots, DOM captures, private
keys, or session handles. Slot and session names are identifiers, not values.
The profile does not identify a credential store or browser implementation.

The first version supports sign-in only. Enrollment, recovery, password
changes, consent grants, logout, account creation, and CAPTCHA bypass are out
of scope. A CAPTCHA or an unknown challenge fails closed.

## Shape

```yaml
profile: uws.browser-authentication.1.0
info:
  title: Example member login
  applicationOrigins: [https://members.example.test]
  authenticationOrigins: [https://login.example.test]
observationKind: accessibility_snapshot
evidence:
  learnedAt: "2026-08-15T00:00:00Z"
  source: reviewed_synthetic_fixture
confidence: high
expiresAfter: P30D
verification:
  lastVerifiedAt: "2026-08-15T00:00:00Z"
  successfulRuns: 1
credentialSlots:
  username: {kind: identifier}
  password: {kind: password}
flows:
  member_login_push:
    sequence:
      - navigate: https://members.example.test/
      - type_credential:
          locator: {role: textbox, name: Username}
          slot: username
      - type_credential:
          locator: {role: textbox, name: Password}
          slot: password
      - click:
          locator: {role: button, name: Sign in}
      - challenge: {kind: push}
      - wait_for:
          locator: {role: heading, name: Member dashboard}
    effects: [establishes_session, sends_mfa_challenge]
    success:
      origin: https://members.example.test
      locator: {role: heading, name: Member dashboard}
```

Origins are exact HTTPS origins. HTTP is permitted only for loopback synthetic
fixtures. Navigation targets are absolute and must remain within the combined
application/authentication origin allowlist. Success must be proven on a
declared origin.

`credentialSlots` maps symbolic identifiers to `identifier`, `password`, or
`totp_seed`. Every referenced slot must be declared. Drivers receive a separate
runtime-private mapping from portable slot names to credential bindings.

## Closed Step Vocabulary

Each flow contains at most 256 single-key steps:

| Step | Purpose |
| --- | --- |
| `navigate` | Open an absolute URL on a declared origin. |
| `type_credential` | Fill one accessibility locator from one declared slot. |
| `click` | Activate one accessibility locator. |
| `challenge` | Perform one explicitly selected MFA handoff. |
| `wait_for` | Wait for one accessibility locator. |

Locators use the same accessibility-first role/name/text/value vocabulary as
`uws.browser.1.5`. CSS, XPath, coordinates, arbitrary JavaScript, and browser
implementation commands are forbidden.

Challenge kinds are `push`, `push_number_match`, `totp`, `sms_otp`,
`email_otp`, `voice_otp`, `passkey`, and `security_key`. TOTP uses a declared
`totp_seed` slot and an OTP input locator. SMS/email/voice OTP uses a locator
and obtains the one-time response from the human challenge channel. Push and
WebAuthn challenges have no locator or credential slot. Alternative MFA
methods are separate named flows; the driver never guesses among them.

Every flow declares `establishes_session`. A flow with a challenge step also
declares `sends_mfa_challenge`. Runtime approval is required before credential
submission or MFA initiation. Email OTP is included for compatibility but is
not a phishing-resistant authenticator.

## Bounds and Lifecycle

- Maximum encoded document size: 1 MiB.
- Maximum combined declared origins: 32.
- Maximum flows: 100.
- Maximum steps per flow: 256.
- Evidence, confidence, expiry, and verification metadata are mandatory.

Freshness, review evidence, package digests, runtime approvals, credential
resolution, and session persistence remain downstream responsibilities.
