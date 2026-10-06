# Verified HCL conformance supplement

The manifest pins six existing, unchanged source artifacts by byte length and
SHA-256. Byte-identical copies live inside this nested module so downloaded
module tests need no parent checkout; manifest paths retain root-fixture lineage. The big model-first fixture exercises operations, source descriptions,
all structural workflow forms, triggers/routes/results/components, content trust,
criteria/actions/idempotency/pending schemas and extensions. Smaller fixtures
cover pending-only and retained browser/authentication/registration declarations.
These are inert metadata tests, not browser execution or new profile authority.

Tests independently decode original/imported JSON with UseNumber and compare
their complete value trees and exact numeric token strings. The dedicated token
corpus covers large integers, precise decimals, exponents/casing, trailing zeros,
signed zero and string/number distinction. Strings include Unicode, quotes,
controls, backslashes and literal template markers. Keys include dollars,
legacy key/escape-prefix collisions, spaces and Unicode; empty/null/absence and
array order remain distinct. No float64 or JCS equality is an oracle.

Malformed sources, duplicate/escaped-duplicate keys, YAML aliases/merges/tags,
unsupported numeric spellings, ambiguous HCL attributes/blocks/labels/decoded
keys, runtime expressions/functions/traversals/comprehensions, failed provenance,
oversized/deep inputs and cancellation refuse with no partial view/value excerpts.
The public external-package API test supplies independent consumer compilation.

Run standalone nested full tests/races/vet with GOWORK=off, GOPROXY=off. The
supplement does not alter the frozen 1.11 corpus or existing version/profile bytes.
