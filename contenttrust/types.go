// Package contenttrust performs deterministic, advisory integrity analysis of
// UWS 1.9.1 documents. It never evaluates runtime values, changes validation
// results, blocks execution, or mutates the analyzed document.
package contenttrust

import (
	"context"

	"github.com/OpenUdon/uws/uws1"
)

// ChannelKind describes how an operation interprets an input channel.
type ChannelKind string

const (
	ChannelData        ChannelKind = "data"
	ChannelInstruction ChannelKind = "instruction"
	ChannelAuthority   ChannelKind = "authority"
)

// ValueCapability describes the injection-relevant shape of a value. It is
// deliberately separate from provenance: a constrained attacker-controlled
// scalar remains untrusted even though it cannot carry free-form text.
type ValueCapability string

const (
	CapabilityFreeText          ValueCapability = "free_text"
	CapabilityConstrainedScalar ValueCapability = "constrained_scalar"
	CapabilityComposite         ValueCapability = "composite"
	CapabilityUnknown           ValueCapability = "unknown"
)

// Severity is the advisory importance of a finding.
type Severity string

const (
	SeverityHigh    Severity = "high"
	SeverityWarning Severity = "warning"
)

const (
	CodeUntrustedInstruction = "content_trust.untrusted_instruction"
	CodeUntrustedAuthority   = "content_trust.untrusted_authority"
	CodeUntrustedControl     = "content_trust.untrusted_control"
	CodeUnknownProvenance    = "content_trust.unknown_provenance"
	CodeOpaqueExpression     = "content_trust.opaque_expression"
	CodeUnresolvedReference  = "content_trust.unresolved_reference"
	CodeNonDominatingRef     = "content_trust.non_dominating_reference"
	CodeResolverFailure      = "content_trust.resolver_failure"
	CodeResolverConflict     = "content_trust.resolver_conflict"
)

var findingMessages = map[string]string{
	CodeUntrustedInstruction: "untrusted free-form content can reach an instruction channel",
	CodeUntrustedAuthority:   "untrusted content can reach an authority-bearing channel",
	CodeUntrustedControl:     "untrusted content can influence control flow",
	CodeUnknownProvenance:    "content provenance cannot be determined statically",
	CodeOpaqueExpression:     "expression syntax is outside the UWS core grammar and has no resolver-supplied references",
	CodeUnresolvedReference:  "expression reference cannot be resolved statically",
	CodeNonDominatingRef:     "referenced step does not dominate this use",
	CodeResolverFailure:      "operation resolver failed while describing the operation",
	CodeResolverConflict:     "more than one resolver claims the operation",
}

// Reference is a resolver-supplied source expression. Resolvers use this when
// an implementation-specific interpolation language prevents the core scanner
// from discovering references in a channel subtree. Expression must identify
// a UWS source expression; Path optionally identifies the relative location of
// the interpolation within the channel for stable edge reporting. When set,
// Path must be an RFC 6901 pointer relative to the channel value.
type Reference struct {
	Expression string `json:"expression"`
	Path       string `json:"path,omitempty"`
}

// InputChannel identifies an operation-object-relative RFC 6901 pointer and
// how the operation interprets the selected value. A pointer beginning with
// /request addresses the whole operation object; other pointers are resolved
// relative to operation.request for convenience.
type InputChannel struct {
	Path       string      `json:"path"`
	Kind       ChannelKind `json:"kind"`
	References []Reference `json:"references,omitempty"`
}

// ValueContract describes an operation output's value capability and an
// optional per-output provenance-inheritance override.
type ValueContract struct {
	Capability ValueCapability `json:"capability,omitempty"`
	// InheritsInputProvenance overrides the operation-wide setting when non-nil.
	InheritsInputProvenance *bool `json:"inheritsInputProvenance,omitempty"`
}

// OutputContract is an alias for the value contract of an operation output.
type OutputContract = ValueContract

// OperationContract is the static view of one operation supplied by a source
// or extension-profile resolver.
type OperationContract struct {
	Inputs                  []InputChannel           `json:"inputs,omitempty"`
	Outputs                 map[string]ValueContract `json:"outputs,omitempty"`
	DefaultTrust            uws1.ContentTrustLevel   `json:"defaultTrust,omitempty"`
	InheritsInputProvenance bool                     `json:"inheritsInputProvenance,omitempty"`
}

// Resolver describes source- or implementation-profile-owned operation
// semantics. owned is false when the resolver does not recognize operation.
// Resolver errors, nil resolvers, and malformed owned contracts become advisory
// resolver-failure findings; context cancellation is returned directly from
// Analyze. A malformed claim is discarded rather than silently normalized.
type Resolver interface {
	ResolveOperation(ctx context.Context, doc *uws1.Document, operation *uws1.Operation) (owned bool, contract OperationContract, err error)
}

// Edge is one statically recovered flow into a document field. It contains no
// runtime value or document-content excerpt.
type Edge struct {
	From       string                 `json:"from"`
	To         string                 `json:"to"`
	Provenance uws1.ContentTrustLevel `json:"provenance"`
	Capability ValueCapability        `json:"capability"`
}

// Finding is a stable advisory diagnostic. Message is fixed by Code; Path is
// the document location at which the issue is observed.
type Finding struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Path     string   `json:"path"`
	Message  string   `json:"message"`
}

// Report is the deterministic result of an analysis.
type Report struct {
	Edges    []Edge    `json:"edges,omitempty"`
	Findings []Finding `json:"findings,omitempty"`
}
