# Browser Capability Distribution Goal

This document is the UWS-owned milestone for the cross-repository browser
capability goal. UWS owns only the portable wire contract and conformance
fixtures. Capture caches, registry metadata, network transport, authoring,
credentials, sessions, and browser execution remain downstream concerns.

## UWS-B01 — Browser Profile Interoperability Baseline

State: Complete; superseded by `uws.browser.1.7` for UWS 1.9 consumers.

This ledger records the browser 1.5 interoperability milestone. Browser 1.5
remains immutable and accepted, while the current capability profile is browser
1.7 with bounded contexts and portable scalar accessibility-text conversion.

### Goal

Prove that `uws.browser.1.5` is sufficient as the portable source contract for
content-addressed distribution and trusted runtime consumption without adding
registry or runtime fields to the UWS document.

### Task Ledger

| Item | State | Notes |
|---|---|---|
| Publish canonical conformance fixtures | `[+]` | Added read-only and confirmed side-effect profiles plus both UWS selector forms; fixtures cover every macro, both wait forms, all output sources, defaults, origins, expiry, login state, and confirmation. |
| Add reusable schema validation entrypoint | `[+]` | `BrowserSourceProfileSchema` returns independent embedded bytes and `ValidateBrowserSourceProfile` validates one JSON or YAML document against `browser.1.5`. |
| Lock registry/runtime separation | `[+]` | The normative supplement now keeps digest, identity/version, lifecycle, storage, signatures, access, sessions, and driver data outside the closed profile schema. |
| Run downstream compatibility | `[+]` | Focused/full UWS tests and vet pass; current Browsertools, OpenUdon, and Udon full suites pass against the unchanged `browser.1.5` wire. |

### Acceptance

- The schema/spec/code contract remains `uws.browser.1.5`; no version bump is
  made unless implementation proves a missing portable semantic.
- Canonical fixtures cover the complete closed macro vocabulary and both
  read-only and side-effectful safety postures.
- Consumers can validate profile bytes through a stable UWS-owned helper.
- Registry metadata and runtime-private session/driver data never enter UWS.

### Verification

```bash
go test ./schemas ./uws1 ./convert
go test ./...
go vet ./...
git diff --check
(cd ../browsertools && go test ./...)
(cd ../openudon && go test ./...)
(cd ../udon && go test ./...)
```

### Downstream Impacts

- Evidence A01 uses generic artifact descriptors and lifecycle assessments; it
  does not import UWS.
- Browsertools E02/P02/M19 validate and package these exact fixtures.
- Udon M26 lowers and executes only schema-valid reviewed profiles.
- OpenUdon A01 materializes profiles without redefining their wire format.
