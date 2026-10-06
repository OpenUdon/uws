# Source-neutral binding and flow reference

`github.com/OpenUdon/uws/binding` owns metadata contracts and advisory checks.
Source parsing stays in source tooling; metadata is not approval or execution
permission. The package imports no APItools/provider/credential client.

## Shape table contract

`uws.shape-table.v1` binds each Source ID/kind and exact lowercase SHA-256 to
operation metadata. Source URLs are provenance, never fetch permission. Native
selectors retain ID/reference distinctions, stable operation keys and known
aliases. Browser-profile and runtime-function kinds are reserved metadata;
UWS does not implement their leaf execution or source production.

OperationShape describes an explicit protocol, supported input/output schemas,
security and completeness. Non-HTTP operations carry no fabricated HTTP method,
path or server. Partial HTTP descriptions may omit unknown details; Complete
cannot assert their absence is resolved. Schema.Known distinguishes a complete
JSON Schema projection from absent/partial evidence. Security alternatives are
OR; requirements within one alternative are AND. An empty alternative is known
anonymous access; unknown security is distinct and cannot become anonymous.

ShapeTable.Validate/Marshal/ParseTable check closed structure, bounded data,
source identity consistency, selector/input/output uniqueness and schema JSON.
Serialization sorts copied source/operation slices while preserving declaration
and security alternative order. No caller metadata is mutated. Limits are
8 MiB per table, 512 sources, 10000 operations and 256 KiB per schema.

NewResolver owns a private decoded snapshot and returns independent results.
It resolves exact source/selector identities, reports missing or ambiguous
matches and propagates cancellation. It cannot verify a producer's claims
against real source bytes: OpenUdon/P09 and consuming workers independently
reproduce/verify those claims before deriving authority. Successful resolution
therefore does not establish trust or enable an operation.

Binding validation and deterministic flow findings are delivered by C09.2/.3.
Ordinary document validation and frozen published contracts remain unchanged.
