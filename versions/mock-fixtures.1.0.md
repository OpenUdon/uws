# UWS Mock Fixtures 1.0

**Format:** `uws.mock-fixtures.1.0`

**Schema:** [`mock-fixtures.1.0.json`](mock-fixtures.1.0.json)

**Go package:** `mockruntime`

This document defines an inert, portable collection of leaf-operation response
fixtures. It does not define a new UWS core version, operation profile, source
protocol, or execution behavior. A mock runtime owns how a matched response is
made visible to runtime expressions. Source and profile adapters own the
response value's meaning and shape.

## Document

A fixture document is one UTF-8 JSON object with exactly these top-level
fields:

| Field | Requirement | Meaning |
|---|---|---|
| `format` | REQUIRED; exact value `uws.mock-fixtures.1.0` | Selects this format version. |
| `fixtures` | REQUIRED; 1–10,000 entries | Response fixtures. |

The document is limited to 16 MiB. Unknown fields, duplicate object member
names, trailing JSON values, invalid UTF-8, empty fixture lists, malformed
records, and unsupported format versions are errors. Readers MUST reject an
unknown version; they MUST NOT downgrade to another format or ignore unknown
fields.

Each fixture has this shape:

```json
{
  "operationId": "list_projects",
  "requestDigest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
  "provenance": { "kind": "example" },
  "response": {
    "statusCode": 200,
    "headers": { "content-type": "application/json" },
    "body": { "projects": [] }
  }
}
```

`operationId` is the UWS-local operation ID and follows the UWS executable-ID
characters `A–Z`, `a–z`, `0–9`, `_`, and `-`. `response` is one JSON value,
including JSON `null`; it is carried as response data rather than interpreted
by this format. An HTTP adapter commonly uses an object with `statusCode`,
`headers`, and `body`, while another source or profile may define another
response value. Fixture format readers preserve that value.

## Request key and digest

A fixture key is the exact pair `(operationId, requestDigest)`. The digest is
computed from the resolved UWS request-binding object passed to the leaf
runtime, after runtime-expression resolution and before source-specific
transport encoding. The object includes all standard request locations
(`path`, `query`, `header`, `cookie`, and `body`) and any request-binding
extensions. It excludes operation metadata, the operation effect, source
identity, headers or defaults added by a source adapter, and transport bytes.
An absent request is the empty object `{}`.

Canonical request bytes use the [JSON Canonicalization Scheme (JCS), RFC 8785](https://www.rfc-editor.org/rfc/rfc8785.html).
The digest is lowercase hexadecimal SHA-256 over those canonical UTF-8 bytes,
prefixed exactly as `sha256:`. Object member order and insignificant JSON
whitespace therefore do not affect a key; array order and JSON values do.
Implementations MUST reject a non-object request or one that cannot be
canonicalized; its encoded JSON form MUST NOT exceed 1 MiB. JCS uses
ECMAScript-compatible IEEE-754 binary64 number
serialization. Authors who need exact large integers or decimal precision
MUST represent those values as strings.

The `mockruntime.CanonicalizeRequest` and `mockruntime.RequestDigest` helpers
implement this rule. The published digest vectors are under
`testdata/mock-fixtures/1.0/`.

A SHA-256 request digest is a matching key, not a privacy control; low-entropy
values may be guessed from it. Request-bound values, including explicit header
or cookie values, MUST NOT contain credentials or private data when they are
used in portable fixtures. Runtime-private credentials are added by an adapter
after the request key is computed and are excluded from the digest.

## Matching and replay

There MUST be at most one fixture for an exact key pair. A duplicate pair makes
the whole document invalid; list order does not resolve ambiguity. Lookup uses
both the local operation ID and the digest. A different operation ID or
request value does not match. There are no wildcard keys, fallback fixtures,
or method-based aliases. A missing exact match is an explicit runtime error.

Repeated invocations with the same pair reuse the same fixture response. A
fixture is not consumed and does not advance a cursor. This makes replay
deterministic under loops, retries, and parallel execution. Format 1.0 cannot
represent occurrence-specific responses for an identical pair; such behavior
requires a later format version or a distinct request key.

## Provenance and export safety

`provenance.kind` is REQUIRED and is one of:

| Kind | Meaning |
|---|---|
| `synthetic` | Produced by a mock/synthesis source. |
| `example` | Supplied from an explicit schema or example source. |
| `recorded` | Derived from a response explicitly supplied for fixture export. |

Recorded fixtures MUST set `redacted: true` and provide `redactedPointers`, an
array of at most 256 unique RFC 6901 JSON Pointers, each at most 512 characters,
naming fields removed or replaced by the explicit redaction step. The empty
pointer denotes the entire response. The array MAY be empty if the redaction
policy ran and found no private fields. Pointers contain locations only; they
MUST NOT include the removed value, replacement secret, credential, or other
private data.
Synthetic and example fixtures MUST omit both redaction fields.

Redaction applies only to the recorded response and MUST NOT rewrite
`requestDigest`. The digest continues to identify the original resolved request
object. Replacing request values changes the digest and will not match the
original request. If a digest derived from the actual request is not acceptable
for export, the caller MUST omit that fixture.

The codec performs no network calls, response capture, credential lookup, or
file writes. It never persists a runtime response automatically. Creating a
recorded fixture requires an explicit caller-supplied redactor; the caller
must then explicitly encode or store the returned fixture. Redaction policy
and private-data classification remain caller responsibilities. A fixture
document containing credentials, tokens, cookies, private response data, or
environment-specific handles MUST NOT be exported or committed.

## Compatibility

This format is versioned independently of UWS core and operation profiles.
UWS 1.12's pending-step declarations and effect labels do not alter fixture
matching. The fixture `operationId` names a local executable operation; a
pending step has no executable operation and cannot match a fixture.
