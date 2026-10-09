package binding

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/OpenUdon/uws/internal/strictjson"
)

var ErrTable = errors.New("invalid or bounded shape-table contract")

func text(s string, max int, required bool) bool {
	if required && s == "" {
		return false
	}
	if len(s) > max || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if r < 32 || r == 127 {
			return false
		}
	}
	return true
}
func sourceValid(s Source) bool {
	if !text(s.ID, 256, true) || !text(s.Kind, 128, true) || !text(s.URL, 4096, false) || len(s.SHA256) != 64 || strings.ToLower(s.SHA256) != s.SHA256 {
		return false
	}
	_, err := hex.DecodeString(s.SHA256)
	return err == nil
}
func schemaValid(s Schema) bool {
	if len(s.JSON) > MaxSchemaBytes || s.Known && len(s.JSON) == 0 {
		return false
	}
	if len(s.JSON) == 0 {
		return true
	}
	if strictjson.ValidateSingleValue(s.JSON) != nil {
		return false
	}
	var v any
	if json.Unmarshal(s.JSON, &v) != nil {
		return false
	}
	switch v.(type) {
	case map[string]any, bool:
		return true
	}
	return false
}
func (t ShapeTable) Validate() error {
	if t.Version != TableVersion || len(t.Sources) > MaxSources || len(t.Operations) > MaxOperations {
		return ErrTable
	}
	sources := map[string]Source{}
	for _, s := range t.Sources {
		if !sourceValid(s) {
			return ErrTable
		}
		if _, ok := sources[s.ID]; ok {
			return ErrTable
		}
		sources[s.ID] = s
	}
	keys := map[string]bool{}
	for _, op := range t.Operations {
		src, ok := sources[op.Source.ID]
		if !ok || src != op.Source {
			return ErrTable
		}
		if !text(op.Selector.Value, 4096, true) || !text(op.Selector.Key, 512, true) || (op.Selector.Kind != "id" && op.Selector.Kind != "ref") || !text(op.Protocol, 128, true) || !text(op.Method, 64, false) || !text(op.Path, 4096, false) {
			return ErrTable
		}
		if op.Protocol != "http" && (op.Method != "" || op.Path != "" || len(op.Servers) > 0) {
			return ErrTable
		}
		if !browserShapeValid(op) {
			return ErrTable
		}
		if op.Protocol == "http" && op.Complete && (op.Method == "" || op.Path == "") {
			return ErrTable
		}
		key := op.Source.ID + "\x00" + op.Selector.Key
		if keys[key] {
			return ErrTable
		}
		keys[key] = true
		if len(op.Aliases) > 32 {
			return ErrTable
		}
		aliases := map[string]bool{op.Selector.Kind + "\x00" + op.Selector.Value: true}
		for _, alias := range op.Aliases {
			k := alias.Kind + "\x00" + alias.Value
			if (alias.Kind != "id" && alias.Kind != "ref") || !text(alias.Value, 4096, true) || alias.Key != op.Selector.Key || aliases[k] {
				return ErrTable
			}
			aliases[k] = true
		}
		if len(op.Inputs) > 2048 || len(op.Outputs) > 2048 || len(op.Servers) > 64 || len(op.Security.Alternatives) > 64 {
			return ErrTable
		}
		for _, server := range op.Servers {
			if !text(server, 4096, true) {
				return ErrTable
			}
		}
		seen := map[string]bool{}
		for _, in := range op.Inputs {
			key := in.Location + "\x00" + in.Name
			if !text(in.Location, 128, true) || !text(in.Name, 1024, true) || seen[key] || !schemaValid(in.Schema) {
				return ErrTable
			}
			seen[key] = true
		}
		seen = map[string]bool{}
		for _, out := range op.Outputs {
			key := out.Location + "\x00" + out.Name
			if !text(out.Location, 128, true) || !text(out.Name, 1024, true) || seen[key] || !schemaValid(out.Schema) {
				return ErrTable
			}
			seen[key] = true
		}
		if op.Security.Known && len(op.Security.Alternatives) == 0 {
			return ErrTable
		}
		for _, alt := range op.Security.Alternatives {
			if len(alt.Requirements) > 64 {
				return ErrTable
			}
			seen = map[string]bool{}
			for _, req := range alt.Requirements {
				if !text(req.Scheme, 256, true) || !text(req.Type, 128, true) || !text(req.Location, 128, false) || !text(req.Name, 1024, false) || seen[req.Scheme] || len(req.Scopes) > 128 {
					return ErrTable
				}
				seen[req.Scheme] = true
				scopes := map[string]bool{}
				for _, scope := range req.Scopes {
					if !text(scope, 1024, true) || scopes[scope] {
						return ErrTable
					}
					scopes[scope] = true
				}
			}
		}
	}
	return nil
}

// Marshal emits deterministic source/operation ordering and copies via JSON.
// Security and binding order remain intact; no OR/AND flattening is performed.
func (t ShapeTable) Marshal() ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	// Sort independent slice copies, never mutate the caller's metadata.
	t.Sources = append([]Source(nil), t.Sources...)
	t.Operations = append([]OperationShape(nil), t.Operations...)
	sort.Slice(t.Sources, func(i, j int) bool { return t.Sources[i].ID < t.Sources[j].ID })
	sort.Slice(t.Operations, func(i, j int) bool {
		a, b := t.Operations[i], t.Operations[j]
		if a.Source.ID != b.Source.ID {
			return a.Source.ID < b.Source.ID
		}
		return a.Selector.Key < b.Selector.Key
	})
	data, err := json.Marshal(t)
	if err != nil || len(data) > MaxTableBytes {
		return nil, ErrTable
	}
	return data, nil
}
func ParseTable(data []byte) (ShapeTable, error) {
	var t ShapeTable
	if len(data) > MaxTableBytes || strictjson.ValidateSingleValue(data) != nil {
		return t, ErrTable
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if d.Decode(&t) != nil {
		return ShapeTable{}, ErrTable
	}
	// Browser activation requires transport fields to be absent, including
	// explicit empty/null fields that a generic optional struct would erase.
	var wire struct {
		Operations []map[string]json.RawMessage `json:"operations"`
	}
	if json.Unmarshal(data, &wire) != nil {
		return ShapeTable{}, ErrTable
	}
	for i, op := range t.Operations {
		if op.Source.Kind != KindBrowser {
			continue
		}
		for name := range wire.Operations[i] {
			if strings.EqualFold(name, "method") || strings.EqualFold(name, "path") || strings.EqualFold(name, "servers") {
				return ShapeTable{}, ErrTable
			}
		}
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return ShapeTable{}, ErrTable
	}
	if t.Validate() != nil {
		return ShapeTable{}, ErrTable
	}
	return t, nil
}

// TableResolver owns a private decoded snapshot to prevent post-review mutation.
type TableResolver struct{ table ShapeTable }

func NewResolver(table ShapeTable) (*TableResolver, error) {
	data, err := table.Marshal()
	if err != nil {
		return nil, err
	}
	snapshot, err := ParseTable(data)
	if err != nil {
		return nil, err
	}
	return &TableResolver{table: snapshot}, nil
}
func (r *TableResolver) Resolve(ctx context.Context, b Binding) (Resolution, error) {
	if ctx == nil {
		return Resolution{}, ErrTable
	}
	if err := ctx.Err(); err != nil {
		return Resolution{}, err
	}
	if r == nil || !sourceValid(b.Source) || (b.SelectorKind != "id" && b.SelectorKind != "ref") || !text(b.SelectorValue, 4096, true) {
		return Resolution{}, ErrTable
	}
	var matches []OperationShape
	for _, op := range r.table.Operations {
		if op.Source != b.Source {
			continue
		}
		matched := op.Selector.Kind == b.SelectorKind && op.Selector.Value == b.SelectorValue
		for _, alias := range op.Aliases {
			if alias.Kind == b.SelectorKind && alias.Value == b.SelectorValue {
				matched = true
			}
		}
		if matched {
			matches = append(matches, op)
		}
	}
	if len(matches) == 0 {
		return Resolution{Status: Missing}, nil
	}
	if len(matches) > 1 {
		return Resolution{Status: Ambiguous}, nil
	}
	// Return an independent shape; a consumer can't mutate the resolver snapshot.
	t := ShapeTable{Version: TableVersion, Sources: []Source{matches[0].Source}, Operations: matches}
	data, _ := t.Marshal()
	copy, err := ParseTable(data)
	if err != nil {
		return Resolution{}, err
	}
	return Resolution{Status: Resolved, Shape: &copy.Operations[0]}, nil
}
