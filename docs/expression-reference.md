# UWS core expression reference

`github.com/OpenUdon/uws/expressions` exposes the existing core grammar from
[UWS 1.12 sections 5.6–5.7](../versions/1.12.0.md#56-formal-grammar).
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
