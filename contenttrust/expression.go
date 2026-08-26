package contenttrust

import (
	"encoding/json"
	"regexp"
	"strings"
)

type expressionReference struct {
	root string
	id   string
	name string
}

type parsedExpression struct {
	references []expressionReference
	condition  bool
}

var (
	idSegmentPattern  = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	jsonNumberPattern = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$`)
)

func parseExpression(raw string) (parsedExpression, bool) {
	if ref, ok := parseSourceExpression(raw); ok {
		return parsedExpression{references: []expressionReference{ref}}, true
	}
	for _, op := range []string{" == ", " != ", " <= ", " >= ", " < ", " > "} {
		index := strings.Index(raw, op)
		if index < 0 {
			continue
		}
		left, ok := parseSourceExpression(raw[:index])
		if !ok {
			// The token can occur inside the right-hand JSON string of a
			// different comparison operator. Keep looking for the operator whose
			// left operand is an exact source expression.
			continue
		}
		parsed := parsedExpression{references: []expressionReference{left}, condition: true}
		rightOperand := raw[index+len(op):]
		if right, ok := parseSourceExpression(rightOperand); ok {
			parsed.references = append(parsed.references, right)
			return parsed, true
		}
		if isJSONScalarLiteral(rightOperand) {
			return parsed, true
		}
		return parsedExpression{}, false
	}
	return parsedExpression{}, false
}

func parseSourceExpression(raw string) (expressionReference, bool) {
	var ref expressionReference
	switch {
	case raw == "$response.statusCode":
		ref.root = "response"
		ref.name = "statusCode"
		return ref, true
	case raw == "$response.body":
		ref.root = "response"
		ref.name = "body"
		return ref, true
	case strings.HasPrefix(raw, "$response.body#"):
		if !validJSONPointerFragment(strings.TrimPrefix(raw, "$response.body")) {
			return expressionReference{}, false
		}
		ref.root = "response"
		ref.name = "body"
		return ref, true
	case strings.HasPrefix(raw, "$response.headers."):
		name := strings.TrimPrefix(raw, "$response.headers.")
		if !idSegmentPattern.MatchString(name) {
			return expressionReference{}, false
		}
		ref.root = "response"
		ref.name = "headers"
		return ref, true
	case strings.HasPrefix(raw, "$outputs."):
		parts := strings.Split(strings.TrimPrefix(raw, "$outputs."), ".")
		if !validSegments(parts) {
			return expressionReference{}, false
		}
		ref.root = "outputs"
		ref.name = parts[0]
		return ref, true
	case strings.HasPrefix(raw, "$steps."):
		parts := strings.Split(strings.TrimPrefix(raw, "$steps."), ".")
		if len(parts) < 3 || parts[1] != "outputs" || !validSegments(parts) {
			return expressionReference{}, false
		}
		ref.root = "steps"
		ref.id = parts[0]
		ref.name = parts[2]
		return ref, true
	case strings.HasPrefix(raw, "$variables."):
		parts := strings.Split(strings.TrimPrefix(raw, "$variables."), ".")
		if !validSegments(parts) {
			return expressionReference{}, false
		}
		ref.root = "variables"
		ref.name = parts[0]
		return ref, true
	case raw == "$trigger":
		ref.root = "trigger"
		return ref, true
	case strings.HasPrefix(raw, "$trigger."):
		if !validSegments(strings.Split(strings.TrimPrefix(raw, "$trigger."), ".")) {
			return expressionReference{}, false
		}
		ref.root = "trigger"
		return ref, true
	case raw == "$inputs":
		ref.root = "inputs"
		return ref, true
	case strings.HasPrefix(raw, "$inputs."):
		parts := strings.Split(strings.TrimPrefix(raw, "$inputs."), ".")
		if !validSegments(parts) {
			return expressionReference{}, false
		}
		ref.root = "inputs"
		ref.name = parts[0]
		return ref, true
	case raw == "$item":
		ref.root = "item"
		return ref, true
	case strings.HasPrefix(raw, "$item."):
		if !validSegments(strings.Split(strings.TrimPrefix(raw, "$item."), ".")) {
			return expressionReference{}, false
		}
		ref.root = "item"
		return ref, true
	case raw == "$index":
		ref.root = "index"
		return ref, true
	default:
		return expressionReference{}, false
	}
}

func validSegments(parts []string) bool {
	if len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		if !idSegmentPattern.MatchString(part) {
			return false
		}
	}
	return true
}

func validJSONPointerFragment(fragment string) bool {
	if fragment == "#" {
		return true
	}
	if !strings.HasPrefix(fragment, "#/") {
		return false
	}
	for i := 1; i < len(fragment); i++ {
		switch fragment[i] {
		case '/':
			continue
		case '~':
			if i+1 >= len(fragment) || (fragment[i+1] != '0' && fragment[i+1] != '1') {
				return false
			}
			i++
		case '%':
			if i+2 >= len(fragment) || !isHex(fragment[i+1]) || !isHex(fragment[i+2]) {
				return false
			}
			i += 2
		default:
			if !isPointerUnreserved(fragment[i]) {
				return false
			}
		}
	}
	return true
}

func isPointerUnreserved(value byte) bool {
	return value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' || value >= '0' && value <= '9' || strings.ContainsRune("-_.", rune(value))
}

func isHex(value byte) bool {
	return value >= '0' && value <= '9' || value >= 'a' && value <= 'f' || value >= 'A' && value <= 'F'
}

func isJSONScalarLiteral(raw string) bool {
	if strings.TrimSpace(raw) != raw {
		return false
	}
	if raw == "true" || raw == "false" || raw == "null" || jsonNumberPattern.MatchString(raw) {
		return true
	}
	if !strings.HasPrefix(raw, `"`) {
		return false
	}
	var value any
	return json.Unmarshal([]byte(raw), &value) == nil
}
