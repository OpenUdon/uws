// Package binding describes source-neutral operation shapes and advisory UWS
// binding/flow checks. It never parses source formats or grants execution.
package binding

import (
	"context"
	"encoding/json"
)

const TableVersion = "uws.shape-table.v1"
const MaxTableBytes = 8 << 20
const MaxSchemaBytes = 256 << 10
const MaxOperations = 10000
const MaxSources = 512
const KindBrowser = "browser-profile"
const KindFunction = "runtime-function"

// Source binds metadata to exact source bytes. URL is provenance, never fetch
// permission. Digest is a lowercase, unprefixed SHA-256; producer IDs aren't trust.
type Source struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	SHA256 string `json:"sha256"`
	URL    string `json:"url,omitempty"`
}

// Selector preserves the source-native ID or reference and stable operation key.
type Selector struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
	Key   string `json:"key"`
}

// Schema carries a JSON Schema projection and explicit completeness. Unknown
// projections remain indeterminate even if their partial JSON happens to fit.
type Schema struct {
	Known bool            `json:"known"`
	JSON  json.RawMessage `json:"json,omitempty"`
}
type Input struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Schema   Schema `json:"schema"`
}
type Output struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Schema   Schema `json:"schema"`
}

// Security alternatives are OR; requirements within each alternative are AND.
// An empty alternative means anonymous. Unknown is distinct from anonymous.
type Security struct {
	Known        bool                  `json:"known"`
	Alternatives []SecurityAlternative `json:"alternatives,omitempty"`
}
type SecurityAlternative struct {
	Requirements []SecurityRequirement `json:"requirements"`
}
type SecurityRequirement struct {
	Scheme   string   `json:"scheme"`
	Type     string   `json:"type"`
	Location string   `json:"location,omitempty"`
	Name     string   `json:"name,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
}

// OperationShape is metadata only. HTTP details are absent for non-HTTP leaves;
// source tooling must never invent a method/server from an RPC selector.
type OperationShape struct {
	Source   Source                 `json:"source"`
	Selector Selector               `json:"selector"`
	Aliases  []Selector             `json:"aliases,omitempty"`
	Protocol string                 `json:"protocol"`
	Method   string                 `json:"method,omitempty"`
	Path     string                 `json:"path,omitempty"`
	Servers  []string               `json:"servers,omitempty"`
	Inputs   []Input                `json:"inputs,omitempty"`
	Outputs  []Output               `json:"outputs,omitempty"`
	Security Security               `json:"security"`
	Complete bool                   `json:"complete"`
	Browser  *BrowserOperationShape `json:"browser,omitempty"`
}

// ShapeTable is untrusted until independently reproduced against source bytes.
// Serialization validation checks structure/identity, not producer correctness.
type ShapeTable struct {
	Version    string           `json:"version"`
	Sources    []Source         `json:"sources"`
	Operations []OperationShape `json:"operations"`
}
type Binding struct {
	Source        Source `json:"source"`
	SelectorKind  string `json:"selector_kind"`
	SelectorValue string `json:"selector_value"`
}
type Status string

const (
	Resolved    Status = "resolved"
	Missing     Status = "missing"
	Ambiguous   Status = "ambiguous"
	Unsupported Status = "unsupported"
)

type Resolution struct {
	Status Status
	Shape  *OperationShape
}

// Resolver is source-neutral. It receives identity, not a file/network loader.
type Resolver interface {
	Resolve(context.Context, Binding) (Resolution, error)
}
