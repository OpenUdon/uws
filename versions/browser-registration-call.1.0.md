# UWS Browser Registration Call Supplement 1.0

`uws.browser-registration-call.1.0` represents one explicitly approved
account-registration mutation as an extension-owned UWS operation. The
normative schema is
[`browser-registration-call.1.0.json`](browser-registration-call.1.0.json).

The operation selects this contract with:

```yaml
x-uws-operation-profile: uws.browser-registration-call.1.0
x-uws-browser-registration:
  profile: browser-registration/dedicated-test-user.yaml
  flow: create_dedicated_test_user
  credentialBindings:
    identifier: dedicated_test_identifier
    password: dedicated_test_password
  approval: register_dedicated_test_user
  duplicatePrevention: operator_attestation
  onDuplicate: fail
  ambiguousOutcome: stop_without_retry
  cleanupDisposition: delete_separately
```

All names are symbolic. The envelope MUST NOT contain account identifiers,
passwords, verification values, cookies, browser storage, or session handles.
`profile` is a safe package-relative path. `flow` selects one exact flow, and
`credentialBindings` maps every declared flow slot to a runtime-owned symbolic
credential name.

`approval` identifies the exact operator approval required immediately before
the selected flow's one `submit` step. A trusted runtime MUST independently
bind that approval to the immutable operation and profile bytes; the string is
not itself proof of approval.

The three fixed controls have no permissive alternatives:

- `duplicatePrevention: operator_attestation` requires a private pre-run
  attestation for the exact candidate;
- `onDuplicate: fail` prevents treating an existing account as success; and
- `ambiguousOutcome: stop_without_retry` prevents blind or automatic retries.

`cleanupDisposition` records the pre-run decision to `delete_separately` or
`retain_dedicated_test_identity`. It does not authorize or perform cleanup.

Registration establishes no named browser session. A later operation that
needs an authenticated session uses a separate
`uws.browser-authentication-call` operation for an existing identity.
