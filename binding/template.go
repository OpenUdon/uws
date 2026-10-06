package binding

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/OpenUdon/uws/expressions"
)

func expressionValue(value any) (string, bool) {
	s, ok := value.(string)
	return s, ok && (strings.HasPrefix(s, "$") || strings.HasPrefix(strings.TrimSpace(s), "expr("))
}
func containsExpression(value any, depth int) bool {
	if depth > 32 {
		return true
	}
	if _, ok := expressionValue(value); ok {
		return true
	}
	switch v := value.(type) {
	case map[string]any:
		for _, x := range v {
			if containsExpression(x, depth+1) {
				return true
			}
		}
	case []any:
		for _, x := range v {
			if containsExpression(x, depth+1) {
				return true
			}
		}
	}
	return false
}
func validateBoundValue(target Schema, value any, types map[string]Schema, scope expressions.Context) Outcome {
	if !containsExpression(value, 0) {
		return validateLiteral(target, value)
	}
	source, ok := templateSchema(value, types, scope, 0)
	if !ok {
		return Indeterminate
	}
	data, err := json.Marshal(source)
	if err != nil || len(data) > MaxSchemaBytes {
		return Indeterminate
	}
	return schemaCompatibility(Schema{Known: true, JSON: data}, target)
}
func templateSchema(value any, types map[string]Schema, scope expressions.Context, depth int) (any, bool) {
	if depth > 32 {
		return nil, false
	}
	if text, ok := expressionValue(value); ok {
		if _, err := expressions.Parse(text, scope); err != nil {
			return nil, false
		}
		s, ok := types[text]
		if !ok || !s.Known || !schemaValid(s) {
			return nil, false
		}
		v, err := decode(s.JSON)
		return v, err == nil
	}
	switch v := value.(type) {
	case map[string]any:
		props := map[string]any{}
		required := make([]string, 0, len(v))
		for key := range v {
			required = append(required, key)
		}
		sort.Strings(required)
		for _, key := range required {
			child, ok := templateSchema(v[key], types, scope, depth+1)
			if !ok {
				return nil, false
			}
			props[key] = child
		}
		return map[string]any{"type": "object", "properties": props, "required": required, "additionalProperties": false, "minProperties": len(v), "maxProperties": len(v)}, true
	case []any:
		prefix := make([]any, len(v))
		for i, x := range v {
			child, ok := templateSchema(x, types, scope, depth+1)
			if !ok {
				return nil, false
			}
			prefix[i] = child
		}
		return map[string]any{"type": "array", "prefixItems": prefix, "minItems": len(v), "maxItems": len(v), "items": false}, true
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return nil, false
		}
		literal, err := decode(data)
		if err != nil {
			return nil, false
		}
		return map[string]any{"const": literal}, true
	}
}
func constraintKeysKnown(schema map[string]any, allowed ...string) bool {
	known := map[string]bool{"type": true, "title": true, "description": true, "$schema": true}
	for _, key := range allowed {
		known[key] = true
	}
	for key := range schema {
		if !known[key] {
			return false
		}
	}
	return true
}
func schemaFrom(value any) Schema {
	data, err := json.Marshal(value)
	return Schema{Known: err == nil, JSON: data}
}
func minmaxContains(source, target map[string]any, minKey, maxKey string) bool {
	number := func(v any) (int, bool) {
		switch n := v.(type) {
		case json.Number:
			i, err := strconv.Atoi(n.String())
			return i, err == nil
		case int:
			return n, true
		}
		return 0, false
	}
	if v, ok := target[minKey]; ok {
		to, known := number(v)
		from, provided := number(source[minKey])
		if !known || !provided || from < to {
			return false
		}
	}
	if v, ok := target[maxKey]; ok {
		to, known := number(v)
		from, provided := number(source[maxKey])
		if !known || !provided || from > to {
			return false
		}
	}
	return true
}
func objectContainment(source, target map[string]any) Outcome {
	if !constraintKeysKnown(target, "properties", "required", "additionalProperties", "minProperties", "maxProperties") {
		return Indeterminate
	}
	if !minmaxContains(source, target, "minProperties", "maxProperties") {
		return Indeterminate
	}
	sp, _ := source["properties"].(map[string]any)
	tp, _ := target["properties"].(map[string]any)
	required := map[string]bool{}
	if list, ok := source["required"].([]any); ok {
		for _, v := range list {
			if s, ok := v.(string); ok {
				required[s] = true
			}
		}
	}
	if list, ok := target["required"].([]any); ok {
		for _, v := range list {
			key, ok := v.(string)
			if !ok || !required[key] {
				return Indeterminate
			}
		}
	}
	// Unknown source extra properties cannot be proved acceptable when target
	// forbids or constrains them. A constructed template is a closed object.
	targetExtra, hasExtra := target["additionalProperties"]
	sourceClosed := source["additionalProperties"] == false
	if !sourceClosed {
		for key := range tp {
			if _, declared := sp[key]; !declared {
				return Indeterminate
			}
		}
	}
	if hasExtra && targetExtra != true && !sourceClosed {
		return Indeterminate
	}
	state := Compatible
	for key, from := range sp {
		to, declared := tp[key]
		if !declared {
			if !hasExtra || targetExtra == true {
				continue
			}
			if targetExtra == false {
				return Incompatible
			}
			to = targetExtra
		}
		child := schemaCompatibility(schemaFrom(from), schemaFrom(to))
		if child == Incompatible {
			return Incompatible
		}
		if child == Indeterminate {
			state = Indeterminate
		}
	}
	return state
}
func arrayContainment(source, target map[string]any) Outcome {
	if !constraintKeysKnown(target, "prefixItems", "items", "minItems", "maxItems") {
		return Indeterminate
	}
	if !minmaxContains(source, target, "minItems", "maxItems") {
		return Indeterminate
	}
	// General array variance/length relationships stay indeterminate; the finite
	// constructed-template case has exact length and a closed prefix list.
	if source["items"] != false {
		return Indeterminate
	}
	prefix, ok := source["prefixItems"].([]any)
	if !ok {
		return Indeterminate
	}
	targetPrefix, _ := target["prefixItems"].([]any)
	state := Compatible
	for i, from := range prefix {
		var to any
		if i < len(targetPrefix) {
			to = targetPrefix[i]
		} else {
			var provided bool
			to, provided = target["items"]
			if !provided || to == true {
				continue
			}
			if to == false {
				return Incompatible
			}
		}
		child := schemaCompatibility(schemaFrom(from), schemaFrom(to))
		if child == Incompatible {
			return Incompatible
		}
		if child == Indeterminate {
			state = Indeterminate
		}
	}
	return state
}
