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

## Advisory binding validation

ValidateBinding accepts a resolver plus a consumer-prepared projection of bound
inputs, symbolic credential slots/scopes, response references and independently
reviewed expression types. Compatible means the metadata constraints checked
are proved; it is never credential readiness or execution permission.
Incompatible identifies a concrete mismatch; missing/partial/unsupported type,
source or schema evidence remains indeterminate. Reports contain stable code,
path and outcome metadata, never input values or resolver/schema error text.

Required input presence and literal schemas are checked with the existing
JSON-Schema implementation and an external-resource loader that always refuses.
No file/network ref is fetched. Exact schema equality, finite enum/const values
or simple proven type containment can establish expression compatibility;
unproved constraint containment stays indeterminate. Nested body/array templates
are projected as closed schemas using reviewed expression types and literal
const leaves; structural containment is proved only for the supported constraint
subset. Partial number/integer overlap is indeterminate, not disjoint.
Response field pointers
use canonical indexes; unavailable schemas and partial missing fields don't
become compatible by omission. Security alternatives preserve OR/AND structure
and require reviewed symbolic slots plus required scopes, without loading keys.

## Deterministic advisory flow

AnalyzeFlow reports possible structural reachability from the main/sole entry
workflow and declared trigger routes, including dependencies, workflow calls
and explicit goto targets. Conditions are not evaluated. Merge children are
not treated as executed. Missing/ambiguous references and cycles remain visible.
Known core references mark step/operation outputs used; output_unreferenced
means no recognized core reference, not proof that an opaque profile cannot
consume the value. Findings are sorted by path/code, capped at 128 with an
explicit Truncated flag, and contain no document values or content excerpts.

Sequence/branch observations show pending steps before subsequent declared
writes and unknown effects without inferring effects from HTTP methods. Await
without a timeout, unbounded retries and loops/forEach without a statically
known item array are work-bound warnings, not claims of an infinite runtime
loop. The analysis does not dispatch effects, alter validation, grant authority
or decide consumer policy. Caller context cancellation is propagated.
