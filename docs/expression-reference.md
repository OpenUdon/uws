# UWS core expression reference

`github.com/OpenUdon/uws/expressions` exposes the existing core grammar from
[UWS 1.12 sections 5.6–5.7](https://github.com/OpenUdon/uws/blob/main/versions/1.12.0.md#56-formal-grammar).
It adds no expression language feature and performs no source/provider I/O.
Published specification/schema/archive bytes remain unchanged.

`Parse(text, Context)` produces an immutable expression with source/operator,
optional source operand, or JSON scalar literal accessors. JSON-number lexemes
remain `json.Number`; parsing never converts them to float64. Errors expose
only stable syntax/version/context/limit categories, not literal values.

Context names the document version, field kind and whether a structural loop
is in scope. Empty version selects the current 1.12 grammar, independently of
document admission. Version syntax is canonical major-1 semver; the caller
still admits only its supported published document versions. Response-body
dot paths and loop-only `$batchIndex` require 1.11+. Bare numeric text is allowed
only for non-await wait and batchSize from 1.11 onward. Those fields accept a
source or number, not a comparison; await.wait uses predicate context.

The parser accepts the nine published expression source families, dot paths,
response-body JSON Pointer fragments and a single binary comparison with one
ASCII space around its operator. Comparison operands may be another source or
a JSON scalar string/number/boolean/null. Parentheses, boolean connectives,
arithmetic and function calls remain implementation-profile extensions.
Runtime-function templates and browser-profile strings are not core-expression
fields and must never be rewritten or scanned as such by a consumer.

`Pointer` checks the canonical fragment grammar, then percent-decodes before
splitting tokens and applies only `~0`/`~1`. An initial slash is required for a
nonempty canonical pointer fragment; percent-encoded root delimiters are not
a canonical grammar form. Pointer evaluation will enforce canonical array
indexes; dot-walk numeric segments retain the published dot-path behavior.

Source and expression input is bounded to the existing 64 KiB mock ceiling.
Parsing does not resolve source values, grant execution or establish source
shape correctness. Evaluation and optional document portability checking are
separate C08 rows; ordinary validation and the existing mock remain unchanged
at the parser milestone boundary.

## Scoped reference evaluation

`NewEvaluator(document)` uses an admitted immutable caller snapshot of the
UWS model and `Evaluator.Evaluate(ctx, text, field)` resolves only in-memory
orchestrator state. It never dispatches a leaf, fetches a source or discovers
credentials. Current workflow records are frame-local; completed steps resolve
from the nearest matching iteration then its enclosing iterations. Qualified
child/foreign invocation records are excluded. Current responses must belong
to the current invocation. Component variables remain visible when unshadowed;
top-level declarations win on overlapping names.

Input indirections retain the existing nesting ceiling of 32. Numeric
comparisons use exact rationals with the existing 256-byte / ±10000 exponent
resource limits. Nonfinite, unsupported or invalid JSON operands refuse;
missing path children propagate null. JSON Pointer array indexes are canonical,
while ordinary dot-walk indexes preserve their existing behavior. ResolveItems
returns a source-ordered independent slice. Context cancellation is checked
before evaluation and recursive expression resolution.

Callers must preserve numbers before constructing the snapshot. The legacy
UWS model's custom JSON unmarshalling historically converts some `any` fields
to float64; setting UseNumber on an outer decoder does not override a custom
UnmarshalJSON method. UseNumber on the source-value projection, then populate
Variables/Components/input/result values with json.Number, preserves exact
integers and numeric lexemes. The shared fixture loader demonstrates this
boundary without changing the ordinary model decoder's compatibility.

The separately versioned corpus `docs/examples/expressions/v1/cases.json`
contains 18 source/version/context/visibility/number vectors and is consumed by
the reference tests. C08.4 pins it as a conformance supplement; it does not
edit any published schema or earlier executable contract.

## Optional strict portability and mock adoption

`CheckPortability(document)` returns deterministic code/path diagnostics for
core expression fields only. Ordinary Document.Validate is unchanged. Core
controls, outputs/structural-result values, simple criteria, standard source-bound request values and step input
bindings use the parser; regex/JSONPath/XPath queries, extensions, opaque bodies,
trigger options, function templates and native browser profile placeholders
remain owned by their profiles. Source-bound browser `request.body` values use
the same core reference grammar and invocation scopes as other source payloads;
browser path/query/header/cookie fields are not interpreted. Authentication and
registration credential/session/input declarations remain opaque. Legacy expr
wrappers receive an explicit diagnostic. Traversal
is depth bounded; diagnostics are capped at 128 and carry no expression values.
Root x-* request extensions remain profile-owned; identically named keys inside
ordinary body/query payloads are still checked as data bindings.
Known loop invocation contexts are tracked so a reusable operation's batch
reference is checked in its actual call scope, not merely its declaration.
Trigger-route workflows and direct top-level step routes contribute their
actual entry contexts. Direct steps bypass the enclosing workflow loop body.
Explicit goto targets also contribute fresh root contexts.
The strict checker also diagnoses `$item`/`$index` outside known iterations;
forEach creates item/index context after dependencies, when, non-await wait
and its collection expression are checked in the incoming context,
while its separate batch-index behavior remains unchanged. Ordinary Parse
continues accepting the declared grammar independently of availability.

The mock now delegates source lookup, comparison and item resolution to the
shared evaluator. Its historical generic numeric-evaluation and encoded-root
pointer extensions remain in a small explicit compatibility adapter. Quoted
comparison literals are never rewritten. The strict checker diagnoses those
forms in their real core field context. Unshadowed component variables now
follow the existing published precedence rule rather than being discarded
when unrelated top-level variables exist. No provider behavior, schema,
ordinary-validator rule or published-version byte changes in this adoption.
