# Mock Fixtures

[UWS Mock Fixtures 1.0](https://github.com/OpenUdon/uws/blob/main/versions/mock-fixtures.1.0.md)
is a versioned JSON format for deterministic leaf-operation response replay.
It is independent of the UWS core version and does not select a transport,
source adapter, or live runtime.

Each fixture is keyed by the UWS-local `operationId` and SHA-256 digest of the
resolved request-binding object. Request keys use RFC 8785 JSON canonical bytes;
the same key always replays the same response, including on repeated,
concurrent, or retried calls. A missing exact key is an error. There are no
wildcard or method-derived matches.

The Go package `github.com/OpenUdon/uws/mockruntime` provides:

- `DecodeFixtures` and `EncodeFixtures` for bounded, strict schema-validated
  fixture JSON;
- `CanonicalizeRequest` and `RequestDigest` for the portable request key;
- `FixtureSet.LookupResponse` for exact response lookup; and
- `NewRecordedFixture`, which requires a caller-supplied redactor and does not
  write files or retain the original response; and
- `NewRuntime`, an orchestrator-backed pure mock runtime. See the
  [mock runtime guide](mock-runtime.md) for response selection, expression
  support, request records, and limitations.

The fixture response is one JSON value whose internal shape stays with the
source or profile adapter. Recorded fixtures require `redacted: true` and a
JSON Pointer list. A redactor may return an empty pointer list when it ran and
found nothing to remove. Credentials, cookies, private responses, and runtime
handles must never be exported. The request digest is not a privacy control;
request-bound secrets are excluded from portable fixtures, and runtime-private
credentials are added after request-key computation.

The fixture codec is independent of execution. `NewRuntime` uses exact fixture
matches first, then caller-supplied examples or schemas when no fixture set is
supplied (or generated fallback is explicitly enabled). It never writes
fixture files or automatically captures responses.

Portable request digest vectors: [`testdata/mock-fixtures/1.0/request-digest-vectors.json`](https://github.com/OpenUdon/uws/blob/main/testdata/mock-fixtures/1.0/request-digest-vectors.json).
