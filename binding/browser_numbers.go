package binding

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
)

// Browser parameter ranges are implicit native profile constraints. Keep the
// complete schema unchanged and inspect exact numbers only in this advisory pass.
func browserIntegerBounds(profile string) (*big.Rat, *big.Rat) {
	lo, hi := "", ""
	switch profile {
	case "uws.browser.1.8":
		lo, hi = "-9223372036854775808", "9223372036854775807"
	case "uws.browser.1.9", "uws.browser.1.10":
		lo, hi = "-9007199254740991", "9007199254740991"
	default:
		return nil, nil
	}
	a, _ := new(big.Rat).SetString(lo)
	b, _ := new(big.Rat).SetString(hi)
	return a, b
}

func browserNumber(value any) (*big.Rat, bool) {
	n, ok := value.(json.Number)
	if !ok || len(n.String()) > 256 {
		return nil, false
	}
	if i := strings.LastIndexAny(n.String(), "eE"); i >= 0 {
		exponent, err := strconv.ParseInt(n.String()[i+1:], 10, 32)
		if err != nil || exponent < -10000 || exponent > 10000 {
			return nil, false
		}
	}
	r, ok := new(big.Rat).SetString(n.String())
	return r, ok
}

func browserIntegerLiteral(value any, lo, hi *big.Rat) Outcome {
	n, ok := browserNumber(value)
	if !ok {
		return Indeterminate
	}
	if !n.IsInt() || n.Cmp(lo) < 0 || n.Cmp(hi) > 0 {
		return Incompatible
	}
	return Compatible
}

func browserIntegerInput(schema Schema, value any, present bool, proofs map[string]Schema, lo, hi *big.Rat) Outcome {
	if !schema.Known {
		return Indeterminate
	}
	target, err := decode(schema.JSON)
	if err != nil {
		return Indeterminate
	}
	if present {
		data, err := json.Marshal(value)
		if err != nil {
			return Incompatible
		}
		value, err = decode(data)
		if err != nil {
			return Indeterminate
		}
	}
	defaults := browserIntegerValue(target, nil, false, proofs, lo, hi, 0)
	state := browserIntegerValue(target, value, present, proofs, lo, hi, 0)
	if defaults == Incompatible || defaults == Indeterminate && state == Compatible {
		return defaults
	}
	return state
}

func browserIntegerValue(target, value any, present bool, proofs map[string]Schema, lo, hi *big.Rat, depth int) Outcome {
	if depth > 32 {
		return Indeterminate
	}
	m, ok := target.(map[string]any)
	if !ok {
		return Indeterminate
	}
	state := Compatible
	join := func(next Outcome) {
		if next == Incompatible || next == Indeterminate && state == Compatible {
			state = next
		}
	}
	if text, isExpression := expressionValue(value); present && isExpression {
		proof, ok := proofs[text]
		if !ok || !proof.Known {
			return Indeterminate
		}
		source, err := decode(proof.JSON)
		if err != nil {
			return Indeterminate
		}
		return browserIntegerType(target, source, lo, hi, depth+1)
	}
	kind, _ := m["type"].(string)
	switch kind {
	case "integer":
		if defaultValue, exists := m["default"]; exists {
			join(browserIntegerLiteral(defaultValue, lo, hi))
		}
		if present {
			join(browserIntegerLiteral(value, lo, hi))
		}
	case "object":
		properties, _ := m["properties"].(map[string]any)
		values, _ := value.(map[string]any)
		for name, child := range properties {
			v, exists := values[name]
			join(browserIntegerValue(child, v, present && exists, proofs, lo, hi, depth+1))
		}
		if present {
			for name := range values {
				if _, declared := properties[name]; !declared {
					join(Indeterminate)
				}
			}
		}
	case "array":
		if values, ok := value.([]any); present && ok {
			for _, v := range values {
				join(browserIntegerValue(m["items"], v, true, proofs, lo, hi, depth+1))
			}
		}
	case "string", "boolean", "number", "null":
	default:
		return Indeterminate
	}
	return state
}

func browserIntegerType(target, source any, lo, hi *big.Rat, depth int) Outcome {
	if depth > 32 {
		return Indeterminate
	}
	t, ok := target.(map[string]any)
	if !ok {
		return Indeterminate
	}
	s, ok := source.(map[string]any)
	if !ok {
		return Indeterminate
	}
	if constant, exists := s["const"]; exists {
		return browserIntegerValue(target, constant, true, nil, lo, hi, depth+1)
	}
	kind, _ := t["type"].(string)
	if kind == "integer" {
		if constant, exists := s["const"]; exists {
			return browserIntegerLiteral(constant, lo, hi)
		}
		if values, ok := s["enum"].([]any); ok && len(values) > 0 {
			allSafe, allUnsafe := true, true
			for _, v := range values {
				state := browserIntegerLiteral(v, lo, hi)
				allSafe = allSafe && state == Compatible
				allUnsafe = allUnsafe && state == Incompatible
			}
			if allSafe {
				return Compatible
			}
			if allUnsafe {
				return Incompatible
			}
			return Indeterminate
		}
		if s["type"] != "integer" {
			return Indeterminate
		}
		minimum, lok := browserNumber(s["minimum"])
		maximum, hok := browserNumber(s["maximum"])
		if lok && hok && minimum.Cmp(lo) >= 0 && maximum.Cmp(hi) <= 0 && minimum.Cmp(maximum) <= 0 {
			return Compatible
		}
		return Indeterminate
	}
	if kind == "object" {
		props, _ := t["properties"].(map[string]any)
		sourceProps, _ := s["properties"].(map[string]any)
		// A proof for named properties does not cover arbitrary extra values.
		// Only finite closed inventories are supported here; pattern/open
		// properties need additional proof and remain indeterminate.
		patterns, _ := s["patternProperties"].(map[string]any)
		if s["type"] != "object" || s["additionalProperties"] != false || len(patterns) != 0 {
			return Indeterminate
		}
		for name := range sourceProps {
			if _, declared := props[name]; !declared {
				return Indeterminate
			}
		}
		state := Compatible
		for name, child := range props {
			next := browserIntegerType(child, sourceProps[name], lo, hi, depth+1)
			if next == Incompatible || next == Indeterminate && state == Compatible {
				state = next
			}
		}
		return state
	}
	if kind == "array" {
		return browserIntegerType(t["items"], s["items"], lo, hi, depth+1)
	}
	if kind == "string" || kind == "boolean" || kind == "number" || kind == "null" {
		return Compatible
	}
	return Indeterminate
}
