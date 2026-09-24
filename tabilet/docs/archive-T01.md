# Archive T01 - Advisory Content-Trust Analysis

**Context.** Static integrity provenance and value-capability analysis for UWS
documents, including profile-owned channel semantics supplied by resolvers.

**Baseline.** 8382d0f26b3b10870125760643078d1a1a3e31b6

**Coverage.** verified

**Supersedes.** none

## Scope And Responsibilities

UWS 1.9.1 introduced optional `contentTrust` declarations and the
`contenttrust` Go package. UWS 1.9.2 retains the wire shape and corrects static
analysis around lexical outputs, inherited provenance, structural execution,
dominance, canonical pointer indexes, and trigger reachability.

The context describes integrity of data entering and moving through a workflow.
It does not replace confidentiality controls, inspect runtime values, enforce a
security policy, change validation outcomes for otherwise valid documents, or
block execution.

## Domain And Workflows

Provenance is one of `trusted`, `untrusted`, or `unknown`. Value capability is
tracked separately as `free_text`, `constrained_scalar`, `composite`, or
`unknown`. This separation preserves attacker influence even when schema
constraints prevent a value from carrying arbitrary text.

Profile resolvers classify operation input channels as `data`, `instruction`,
or `authority`; describe output capabilities and provenance inheritance; and
may supply references recovered from profile-specific interpolation. Untrusted
free text reaching instructions and untrusted content reaching authority are
high-severity findings. Untrusted controls, unknown or unresolved flow,
non-dominating references, opaque expressions, and resolver failures or
conflicts are warnings.

The analyzer recovers exact core expressions from request bindings, inputs,
controls, criteria, idempotency keys, outputs, and structural results. It models
sequence order, explicit dependencies, branches, loops, nested structures,
workflow calls, and trigger reachability. Propagation repeats to an actual
fixed point and then emits sorted, deterministic edges and findings.

## System Shape

`contenttrust.Analyze` first runs ordinary document validation, indexes the
workflow graph, resolves operation contracts, and initializes declared or
default provenance. Resolver failures become findings unless analysis itself
cannot continue, while context cancellation and malformed UWS documents are
returned as errors.

Resolver contracts are checked before use. Channel paths and reference paths
must be valid RFC 6901 pointers, array indexes must be canonical decimal
tokens, output contracts may name only declared outputs, and claimed labels
must be recognized. Invalid resolver claims are discarded, with core request
scanning retained as a conservative fallback.

## Contracts And Dependencies

The root wire registry is available only to documents declaring UWS 1.9.1 or
later. Its identifiers and per-output or per-input keys must resolve to the
same document. Operation output trust follows declared output, operation
default, source default, resolver default, then `unknown`; inherited untrusted
provenance cannot be upgraded by an alias or output declaration.

`contenttrust.Resolver` is the extension point for source- and profile-aware
semantics. Reports expose stable codes, severities, paths, and fixed messages,
but never document excerpts, runtime values, credentials, or server errors.
Consumers own any policy that interprets the advisory report.

## Operations And Verification

The analyzer is deterministic and value-free. External trigger payloads default
to untrusted, external workflow inputs default to unknown, and reviewed document
literals are treated as trusted. Unknown extension interpolation stays unknown
or produces an opaque-expression finding instead of being guessed.

The test suite covers resolver diagnostics, fixed-point propagation, sequence
and parallel dominance, branch and loop behavior, workflow calls, aliases,
merge semantics, trigger scoping, namespace precedence, canonical pointer
indexes, stable reporting, and operation trust precedence. The full repository
test and vet suites pass at this baseline.

## Evidence

| Claim | Repository Evidence |
|---|---|
| Content trust is additive, advisory, and separate from ordinary validation and execution. | `docs/content-trust.md`; `versions/1.9.2.md` sections 4.5.18-4.7; `contenttrust/types.go` |
| Provenance and value capability are independent analysis dimensions. | `contenttrust/types.go`; `docs/content-trust.md` section "Provenance is not capability" |
| Profile semantics enter through validated resolver contracts. | `contenttrust/types.go`; `contenttrust/analyze.go`; `contenttrust/analyze_test.go` |
| Analysis uses a fixed point and emits deterministic sorted results. | `contenttrust/analyze.go`; `contenttrust/analyze_test.go` |
| Core declarations are version-gated and semantically cross-referenced. | `uws1/content_trust.go`; `uws1/content_trust_test.go`; `uws1/validation.go` |
| UWS 1.9.2 corrections are locked by regression coverage. | `versions/CHANGELOG.md`; `contenttrust/analyze_test.go`; `uws1/schema_conformance_test.go` |

## Observed Gaps

The repository publishes the resolver interface but no concrete resolver for a
source family, runtime supplement, or browser profile. Resolver selection and
report-enforcement policy therefore remain consumer-owned and cannot be
verified here. Runtime-only data flow is intentionally outside this static
analysis.
