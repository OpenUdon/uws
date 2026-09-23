package uws1

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/theory/jsonpath"
)

func (o *Orchestrator) criteriaMatchAll(ctx context.Context, criteria []*Criterion) (bool, error) {
	if len(criteria) == 0 {
		return true, nil
	}
	for _, criterion := range criteria {
		if criterion == nil {
			continue
		}
		ok, err := o.evaluateCriterion(ctx, criterion)
		if err != nil {
			return false, fmt.Errorf("evaluating criterion %q (%s): %w", criterion.Condition, criterion.Type, err)
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (o *Orchestrator) evaluateCriterion(ctx context.Context, criterion *Criterion) (bool, error) {
	if criterion == nil {
		return true, nil
	}
	switch criterion.Type {
	case "", CriterionSimple:
		return o.evaluateTruthy(ctx, criterion.Condition)
	case CriterionRegex:
		return o.evaluateRegexCriterion(ctx, criterion)
	case CriterionJSONPath:
		return o.evaluateJSONPathCriterion(ctx, criterion)
	case CriterionXPath:
		return o.evaluateXPathCriterion(ctx, criterion)
	default:
		return false, fmt.Errorf("uws1: criterion type %q is not executable by the core orchestrator", criterion.Type)
	}
}

func (o *Orchestrator) evaluateRegexCriterion(ctx context.Context, criterion *Criterion) (bool, error) {
	source, err := o.Runtime.EvaluateExpression(ctx, criterion.Context)
	if err != nil {
		return false, fmt.Errorf("evaluating regex context %q: %w", criterion.Context, err)
	}
	pattern := strings.TrimSpace(criterion.Condition)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("compile regex %q: %w", pattern, err)
	}
	text, err := criterionSourceString(source, supportsUWSVersionAtLeast(o.documentVersion(), 1, 10, 0))
	if err != nil {
		return false, err
	}
	return re.MatchString(text), nil
}

func (o *Orchestrator) evaluateJSONPathCriterion(ctx context.Context, criterion *Criterion) (bool, error) {
	source, err := o.Runtime.EvaluateExpression(ctx, criterion.Context)
	if err != nil {
		return false, fmt.Errorf("evaluating jsonpath context %q: %w", criterion.Context, err)
	}
	target := normalizeCriterionTarget(criterion.Context, criterion.Condition)
	canonicalIndexes := o.Document != nil && supportsUWSVersionAtLeast(o.Document.UWS, 1, 9, 2)
	switch {
	case target == "":
		return o.truthy(source)
	case strings.HasPrefix(target, "#"):
		value, found, err := o.resolveCriterionPointer(source, target, canonicalIndexes)
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
		return o.truthy(value)
	case strings.HasPrefix(target, "/"):
		value, found, err := o.resolveCriterionPointer(source, "#"+target, canonicalIndexes)
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
		return o.truthy(value)
	default:
		path, err := jsonpath.Parse(target)
		if err != nil {
			return false, fmt.Errorf("parse jsonpath %q: %w", target, err)
		}
		for _, value := range path.Select(source) {
			matched, err := o.truthy(value)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
		return false, nil
	}
}

func (o *Orchestrator) evaluateXPathCriterion(ctx context.Context, criterion *Criterion) (bool, error) {
	source, err := o.Runtime.EvaluateExpression(ctx, criterion.Context)
	if err != nil {
		return false, fmt.Errorf("evaluating xpath context %q: %w", criterion.Context, err)
	}
	xmlText, err := criterionSourceString(source, supportsUWSVersionAtLeast(o.documentVersion(), 1, 10, 0))
	if err != nil {
		return false, err
	}
	root, err := xmlquery.Parse(strings.NewReader(strings.TrimSpace(xmlText)))
	if err != nil {
		return false, fmt.Errorf("parse xpath XML: %w", err)
	}
	target := normalizeCriterionTarget(criterion.Context, criterion.Condition)
	if target == "" {
		return truthyValue(xmlText), nil
	}
	compiled, err := xpath.Compile(target)
	if err != nil {
		return false, fmt.Errorf("compile xpath %q: %w", target, err)
	}
	value := compiled.Evaluate(xmlquery.CreateXPathNavigator(root))
	return truthyXPathValueVersioned(value, supportsUWSVersionAtLeast(o.documentVersion(), 1, 10, 0)), nil
}

func normalizeCriterionTarget(contextExpr, condition string) string {
	contextExpr = strings.TrimSpace(contextExpr)
	condition = strings.TrimSpace(condition)
	if contextExpr != "" && strings.HasPrefix(condition, contextExpr) {
		return strings.TrimSpace(strings.TrimPrefix(condition, contextExpr))
	}
	return condition
}

func criterionSourceString(value any, requireString bool) (string, error) {
	switch typed := value.(type) {
	case string:
		return typed, nil
	case []byte:
		if requireString && !utf8.Valid(typed) {
			return "", fmt.Errorf("criterion context bytes are not valid UTF-8 text")
		}
		return string(typed), nil
	default:
		if requireString {
			return "", fmt.Errorf("criterion context must resolve to a string")
		}
		if value == nil {
			return "", nil
		}
		if data, err := json.Marshal(typed); err == nil {
			return string(data), nil
		}
		return fmt.Sprint(typed), nil
	}
}

func (o *Orchestrator) documentVersion() string {
	if o == nil || o.Document == nil {
		return ""
	}
	return o.Document.UWS
}

func (o *Orchestrator) truthy(value any) (bool, error) {
	if supportsUWSVersionAtLeast(o.documentVersion(), 1, 10, 0) {
		return truthyValueV110(value)
	}
	return truthyValue(value), nil
}

func (o *Orchestrator) truthyValue(value any) bool {
	truthy, err := o.truthy(value)
	return err == nil && truthy
}

func truthyValueV110(value any) (bool, error) {
	if value == nil {
		return false, nil
	}
	if number, ok := value.(json.Number); ok {
		rational, valid := new(big.Rat).SetString(number.String())
		if !valid {
			return false, fmt.Errorf("runtime expression resolved to an invalid JSON number")
		}
		return rational.Sign() != 0, nil
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Bool:
		return reflected.Bool(), nil
	case reflect.String:
		return reflected.Len() != 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflected.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflected.Uint() != 0, nil
	case reflect.Float32, reflect.Float64:
		number := reflected.Float()
		if math.IsNaN(number) || math.IsInf(number, 0) {
			return false, fmt.Errorf("runtime expression resolved to a non-finite number")
		}
		return number != 0, nil
	case reflect.Slice, reflect.Array, reflect.Map:
		return reflected.Len() != 0, nil
	}
	return false, fmt.Errorf("runtime expression resolved to unsupported non-JSON value %T", value)
}

func (o *Orchestrator) resolveCriterionPointer(root any, pointer string, canonicalIndexes bool) (any, bool, error) {
	if !supportsUWSVersionAtLeast(o.documentVersion(), 1, 10, 0) {
		value, err := resolveCriterionJSONPointerVersioned(root, pointer, canonicalIndexes)
		return value, value != nil, err
	}
	return resolveCriterionJSONPointerV110(root, pointer, canonicalIndexes)
}

func resolveCriterionJSONPointerV110(root any, pointer string, canonicalIndexes bool) (any, bool, error) {
	if pointer == "#" {
		return root, true, nil
	}
	if !strings.HasPrefix(pointer, "#/") {
		return nil, false, fmt.Errorf("invalid JSON pointer fragment %q", pointer)
	}
	path, err := url.PathUnescape(strings.TrimPrefix(pointer, "#"))
	if err != nil {
		return nil, false, fmt.Errorf("invalid percent encoding in JSON pointer fragment %q", pointer)
	}
	if !utf8.ValidString(path) {
		return nil, false, fmt.Errorf("JSON pointer fragment %q is not valid UTF-8", pointer)
	}
	current := root
	for _, rawToken := range strings.Split(strings.TrimPrefix(path, "/"), "/") {
		token, err := decodeCriterionPointerToken(rawToken)
		if err != nil {
			return nil, false, err
		}
		switch typed := current.(type) {
		case map[string]any:
			var found bool
			current, found = typed[token]
			if !found {
				return nil, false, nil
			}
		case []any:
			index, err := parseCriterionIndexVersioned(token, canonicalIndexes)
			if err != nil {
				return nil, false, err
			}
			if index >= len(typed) {
				return nil, false, nil
			}
			current = typed[index]
		default:
			return nil, false, nil
		}
	}
	return current, true, nil
}

func decodeCriterionPointerToken(token string) (string, error) {
	for index := 0; index < len(token); index++ {
		if token[index] != '~' {
			continue
		}
		if index+1 >= len(token) || (token[index+1] != '0' && token[index+1] != '1') {
			return "", fmt.Errorf("invalid JSON pointer escape in token %q", token)
		}
		index++
	}
	return strings.ReplaceAll(strings.ReplaceAll(token, "~1", "/"), "~0", "~"), nil
}

func truthyValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case bool:
		return typed
	case string:
		return typed != ""
	case int:
		return typed != 0
	case int64:
		return typed != 0
	case float64:
		return typed != 0
	case []any:
		return len(typed) > 0
	case map[string]any:
		return len(typed) > 0
	default:
		return true
	}
}

func resolveCriterionJSONPointer(root any, pointer string) (any, error) {
	return resolveCriterionJSONPointerVersioned(root, pointer, true)
}

func resolveCriterionJSONPointerVersioned(root any, pointer string, canonicalIndexes bool) (any, error) {
	if pointer == "" || pointer == "#" {
		return root, nil
	}
	if !strings.HasPrefix(pointer, "#") {
		return nil, fmt.Errorf("invalid JSON pointer %q", pointer)
	}
	path := strings.TrimPrefix(pointer, "#")
	if path == "" {
		return root, nil
	}
	if !strings.HasPrefix(path, "/") {
		return nil, fmt.Errorf("invalid JSON pointer %q", pointer)
	}
	current := root
	for _, rawToken := range strings.Split(path[1:], "/") {
		token := strings.ReplaceAll(strings.ReplaceAll(rawToken, "~1", "/"), "~0", "~")
		switch typed := current.(type) {
		case map[string]any:
			current = typed[token]
		case []any:
			index, err := parseCriterionIndexVersioned(token, canonicalIndexes)
			if err != nil {
				return nil, err
			}
			if index < 0 || index >= len(typed) {
				return nil, nil
			}
			current = typed[index]
		default:
			return nil, nil
		}
	}
	return current, nil
}

func parseCriterionIndex(token string) (int, error) {
	return parseCriterionIndexVersioned(token, true)
}

func parseCriterionIndexVersioned(token string, canonical bool) (int, error) {
	if token == "" || (canonical && len(token) > 1 && token[0] == '0') {
		return 0, fmt.Errorf("invalid array index %q", token)
	}
	for _, ch := range token {
		if ch < '0' || ch > '9' {
			return 0, fmt.Errorf("invalid array index %q", token)
		}
	}
	index, err := strconv.Atoi(token)
	if err != nil {
		return 0, fmt.Errorf("invalid array index %q", token)
	}
	return index, nil
}

func truthyXPathValue(value any) bool {
	switch typed := value.(type) {
	case nil:
		return false
	case bool:
		return typed
	case string:
		return typed != ""
	case float64:
		return typed != 0
	case *xpath.NodeIterator:
		if typed == nil {
			return false
		}
		return typed.MoveNext()
	default:
		return truthyValue(value)
	}
}

func truthyXPathValueVersioned(value any, uws110 bool) bool {
	if !uws110 {
		return truthyXPathValue(value)
	}
	switch typed := value.(type) {
	case float64:
		return typed != 0 && !math.IsNaN(typed)
	default:
		return truthyXPathValue(value)
	}
}
