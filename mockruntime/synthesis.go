package mockruntime

import (
	"fmt"
	"sort"

	"github.com/OpenUdon/uws/uws1"
)

const (
	maxSynthesisDepth = 64
	maxSynthesisNodes = 10000
)

type synthesisState struct {
	nodes  int
	active map[*uws1.ParamSchema]bool
}

func synthesizeResponse(schema *uws1.ParamSchema) (any, error) {
	state := &synthesisState{active: make(map[*uws1.ParamSchema]bool)}
	return state.value(schema, 0)
}

func (s *synthesisState) value(schema *uws1.ParamSchema, depth int) (any, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is missing")
	}
	if depth > maxSynthesisDepth {
		return nil, fmt.Errorf("schema exceeds synthesis depth %d", maxSynthesisDepth)
	}
	s.nodes++
	if s.nodes > maxSynthesisNodes {
		return nil, fmt.Errorf("schema exceeds synthesis node limit %d", maxSynthesisNodes)
	}
	if s.active[schema] {
		return nil, fmt.Errorf("recursive schema references cannot be synthesized")
	}
	if schema.Ref != "" {
		return nil, fmt.Errorf("$ref schemas require a caller-supplied example")
	}
	if len(schema.AllOf) > 0 || len(schema.OneOf) > 0 || len(schema.AnyOf) > 0 {
		return nil, fmt.Errorf("composed schemas require a caller-supplied example")
	}
	if schema.Format != "" {
		return nil, fmt.Errorf("schema format %q requires a caller-supplied example", schema.Format)
	}
	s.active[schema] = true
	defer delete(s.active, schema)

	switch schema.Type {
	case "string":
		return "", nil
	case "integer", "number":
		return 0, nil
	case "boolean":
		return false, nil
	case "null":
		return nil, nil
	case "array":
		if schema.Items == nil {
			return []any{}, nil
		}
		item, err := s.value(schema.Items, depth+1)
		if err != nil {
			return nil, fmt.Errorf("array item schema: %w", err)
		}
		return []any{item}, nil
	case "object":
		properties := make(map[string]*uws1.ParamSchema, len(schema.Properties))
		for name, property := range schema.Properties {
			properties[name] = property
		}
		for _, required := range schema.Required {
			if _, ok := properties[required]; !ok {
				return nil, fmt.Errorf("required property %q has no declared schema", required)
			}
		}
		names := make([]string, 0, len(properties))
		for name := range properties {
			names = append(names, name)
		}
		sort.Strings(names)
		result := make(map[string]any, len(names))
		for _, name := range names {
			value, err := s.value(properties[name], depth+1)
			if err != nil {
				return nil, fmt.Errorf("property %q: %w", name, err)
			}
			result[name] = value
		}
		return result, nil
	default:
		return nil, fmt.Errorf("schema type %q is unsupported for deterministic synthesis", schema.Type)
	}
}
