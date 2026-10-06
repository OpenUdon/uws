// Package expressions implements the published UWS core expression grammar.
// It performs no source I/O, credential lookup or provider operation.
package expressions

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// MaxBytes retains the existing mock expression input ceiling.
const MaxBytes = 64 << 10

var (
	ErrSyntax  = errors.New("expression does not match the UWS core grammar")
	ErrVersion = errors.New("expression feature is unavailable in the declared version")
	ErrContext = errors.New("expression feature is unavailable in this context")
	ErrLimit   = errors.New("expression exceeds the reference parser byte limit")
)

// Field determines whether a bare number is a published field-specific form.
type Field uint8

const (
	Value Field = iota
	Predicate
	Wait
	BatchSize
)

// Context supplies declared-version and structural-loop information. An empty
// version selects the current 1.12 grammar; document admission is caller-owned.
type Context struct {
	Version string
	Field   Field
	InLoop  bool
}

// Expression is an immutable parse result. Literal numbers preserve their JSON
// lexemes; source lookup and workflow/iteration visibility belong to evaluation.
type Expression struct {
	left     string
	operator string
	right    string
	literal  any
	number   bool
}

func (e *Expression) Source() string        { return e.left }
func (e *Expression) Operator() string      { return e.operator }
func (e *Expression) OperandSource() string { return e.right }
func (e *Expression) Literal() any          { return e.literal }
func (e *Expression) IsNumber() bool        { return e.number }

// Parse accepts one source or one binary comparison. Bare numeric literals are
// available only in non-await wait/batch-size fields from UWS 1.11 onward.
// Errors contain no input values, making them suitable for bounded diagnostics.
func Parse(text string, scope Context) (*Expression, error) {
	if len(text) > MaxBytes {
		return nil, ErrLimit
	}
	if text == "" || strings.TrimSpace(text) != text || scope.Field > BatchSize {
		return nil, ErrSyntax
	}
	if scope.Version == "" {
		scope.Version = "1.12.0"
	}
	minor, _, ok := version(scope.Version)
	if !ok {
		return nil, ErrVersion
	}
	if !strings.HasPrefix(text, "$") {
		scalar, err := parseScalar(text)
		if err != nil {
			return nil, ErrSyntax
		}
		if _, ok := scalar.(json.Number); !ok {
			return nil, ErrSyntax
		}
		if scope.Field != Wait && scope.Field != BatchSize {
			return nil, ErrContext
		}
		if minor < 11 {
			return nil, ErrVersion
		}
		return &Expression{literal: scalar, number: true}, nil
	}
	left, rest, compared := strings.Cut(text, " ")
	if err := source(left, scope); err != nil {
		return nil, err
	}
	result := &Expression{left: left}
	if compared && (scope.Field == Wait || scope.Field == BatchSize) {
		return nil, ErrContext
	}
	if !compared {
		return result, nil
	}
	for _, operator := range []string{"==", "!=", "<=", ">=", "<", ">"} {
		if !strings.HasPrefix(rest, operator+" ") {
			continue
		}
		operand := rest[len(operator)+1:]
		if operand == "" {
			return nil, ErrSyntax
		}
		result.operator = operator
		if strings.HasPrefix(operand, "$") {
			if err := source(operand, scope); err != nil {
				return nil, err
			}
			result.right = operand
		} else {
			value, err := parseScalar(operand)
			if err != nil {
				return nil, ErrSyntax
			}
			result.literal = value
		}
		return result, nil
	}
	return nil, ErrSyntax
}

func source(text string, scope Context) error {
	switch {
	case text == "$response.statusCode", text == "$response.body", text == "$inputs", text == "$trigger", text == "$item", text == "$index":
		return nil
	case text == "$batchIndex":
		minor, _, _ := version(scope.Version)
		if minor < 11 {
			return ErrVersion
		}
		if !scope.InLoop {
			return ErrContext
		}
		return nil
	case strings.HasPrefix(text, "$response.headers."):
		if segment(strings.TrimPrefix(text, "$response.headers.")) {
			return nil
		}
	case strings.HasPrefix(text, "$response.body#"):
		_, err := Pointer(strings.TrimPrefix(text, "$response.body"))
		return err
	case strings.HasPrefix(text, "$response.body."):
		minor, _, _ := version(scope.Version)
		if minor < 11 {
			return ErrVersion
		}
		if path(strings.TrimPrefix(text, "$response.body.")) {
			return nil
		}
	case strings.HasPrefix(text, "$steps."):
		parts := strings.Split(strings.TrimPrefix(text, "$steps."), ".")
		if len(parts) >= 3 && segment(parts[0]) && parts[1] == "outputs" && segment(parts[2]) && (len(parts) == 3 || path(strings.Join(parts[3:], "."))) {
			return nil
		}
	default:
		for _, prefix := range []string{"$outputs.", "$variables.", "$inputs.", "$trigger.", "$item."} {
			if strings.HasPrefix(text, prefix) && path(strings.TrimPrefix(text, prefix)) {
				return nil
			}
		}
	}
	return ErrSyntax
}

func segment(text string) bool {
	if text == "" {
		return false
	}
	for i := 0; i < len(text); i++ {
		c := text[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}
func path(text string) bool {
	for _, part := range strings.Split(text, ".") {
		if !segment(part) {
			return false
		}
	}
	return true
}
func version(text string) (int, int, bool) {
	if len(text) > 32 {
		return 0, 0, false
	}
	parts := strings.Split(text, ".")
	if len(parts) != 3 || parts[0] != "1" {
		return 0, 0, false
	}
	for _, part := range parts {
		if len(part) == 0 || (len(part) > 1 && part[0] == '0') {
			return 0, 0, false
		}
		for i := range part {
			if part[i] < '0' || part[i] > '9' {
				return 0, 0, false
			}
		}
	}
	minor, e1 := strconv.Atoi(parts[1])
	patch, e2 := strconv.Atoi(parts[2])
	return minor, patch, e1 == nil && e2 == nil
}
func parseScalar(text string) (any, error) {
	if text == "" || strings.TrimSpace(text) != text || !json.Valid([]byte(text)) {
		return nil, ErrSyntax
	}
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, ErrSyntax
	}
	switch value.(type) {
	case nil, string, bool, json.Number:
		return value, nil
	}
	return nil, ErrSyntax
}

// Pointer parses the published fragment syntax. Percent decoding precedes token
// splitting; only ~0/~1 are admitted and evaluation enforces canonical indexes.
func Pointer(text string) ([]string, error) {
	if len(text) > MaxBytes {
		return nil, ErrLimit
	}
	if !strings.HasPrefix(text, "#") {
		return nil, ErrSyntax
	}
	fragment := text[1:]
	if fragment != "" && !strings.HasPrefix(fragment, "/") {
		return nil, ErrSyntax
	}
	for i := 0; i < len(fragment); i++ {
		c := fragment[i]
		switch {
		case c == '/' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.':
		case c == '%':
			if i+2 >= len(fragment) || !hex(fragment[i+1]) || !hex(fragment[i+2]) {
				return nil, ErrSyntax
			}
			i += 2
		case c == '~':
			if i+1 >= len(fragment) || (fragment[i+1] != '0' && fragment[i+1] != '1') {
				return nil, ErrSyntax
			}
			i++
		default:
			return nil, ErrSyntax
		}
	}
	decoded, err := url.PathUnescape(fragment)
	if err != nil {
		return nil, ErrSyntax
	}
	if decoded == "" {
		return nil, nil
	}
	parts := strings.Split(decoded[1:], "/")
	for i, part := range parts {
		var out strings.Builder
		for j := 0; j < len(part); j++ {
			if part[j] != '~' {
				out.WriteByte(part[j])
				continue
			}
			if j+1 >= len(part) {
				return nil, ErrSyntax
			}
			j++
			switch part[j] {
			case '0':
				out.WriteByte('~')
			case '1':
				out.WriteByte('/')
			default:
				return nil, ErrSyntax
			}
		}
		parts[i] = out.String()
	}
	return parts, nil
}
func hex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
