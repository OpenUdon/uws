package expressions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenUdon/uws/internal/strictjson"
	"github.com/OpenUdon/uws/uws1"
	"math/big"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const maxExpressionDepth = 32

// Evaluator resolves only orchestrator-owned in-memory context. The document
// must be an admitted, immutable caller-owned snapshot; no leaf is dispatched.
type Evaluator struct{ document *uws1.Document }

func NewEvaluator(document *uws1.Document) (*Evaluator, error) {
	if document == nil {
		return nil, errors.New("reference evaluator requires a document")
	}
	declared := document.UWS
	if declared == "" {
		declared = "1.12.0"
	}
	if _, _, ok := version(declared); !ok {
		return nil, ErrVersion
	}
	return &Evaluator{document: document}, nil
}

// Evaluate applies the grammar and field/version gates, then resolves sources
// in the current workflow invocation and nearest enclosing iteration.
func (r *Evaluator) Evaluate(ctx context.Context, text string, field Field) (any, error) {
	value, err := r.evaluateExpression(ctx, text, field, 0)
	if err != nil {
		return nil, err
	}
	// Reject non-finite/unsupported nested values without changing native scalar
	// types or the existing typed-map/slice behavior of source lookup.
	if _, err := decodeJSONValue(value); err != nil {
		return nil, errors.New("expression result is not a JSON-compatible value")
	}
	return value, nil
}

func (r *Evaluator) evaluateExpression(ctx context.Context, text string, field Field, depth int) (any, error) {
	if r == nil || r.document == nil {
		return nil, errors.New("reference evaluator is not initialized")
	}
	if ctx == nil {
		return nil, ErrContext
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if depth > maxExpressionDepth {
		return nil, errors.New("reference expression/input nesting exceeds 32")
	}
	state, _ := uws1.ExecutionContextFromContext(ctx)
	inLoop := state != nil && state.Iteration != nil && state.Iteration.BatchIndex >= 0
	parsed, err := Parse(text, Context{Version: r.document.UWS, Field: field, InLoop: inLoop})
	if err != nil {
		return nil, err
	}
	if parsed.number {
		return parsed.literal, nil
	}
	value, err := r.evaluateSource(ctx, parsed.left, depth)
	if err != nil || parsed.operator == "" {
		return value, err
	}
	other := parsed.literal
	if parsed.right != "" {
		other, err = r.evaluateSource(ctx, parsed.right, depth)
	}
	if err != nil {
		return nil, err
	}
	for _, operand := range []any{value, other} {
		if _, err := decodeJSONValue(operand); err != nil {
			return nil, errors.New("comparison operand is not a JSON-compatible value")
		}
	}
	return compareExpressionValues(value, other, parsed.operator)
}

// ResolveItems retains source order and returns an independent slice.
func (r *Evaluator) ResolveItems(ctx context.Context, text string) ([]any, error) {
	value, err := r.Evaluate(ctx, text, Value)
	if err != nil {
		return nil, err
	}
	normalized, err := decodeJSONValue(value)
	if err != nil {
		return nil, err
	}
	items, ok := normalized.([]any)
	if !ok {
		return nil, errors.New("items source must resolve to an array")
	}
	return append([]any(nil), items...), nil
}

func (r *Evaluator) evaluateSource(ctx context.Context, source string, depth int) (any, error) {
	state, _ := uws1.ExecutionContextFromContext(ctx)
	switch {
	case source == "$response.statusCode":
		response, err := currentResponse(state)
		if err != nil {
			return nil, err
		}
		return dotWalk(response, []string{"statusCode"}), nil
	case strings.HasPrefix(source, "$response.headers."):
		response, err := currentResponse(state)
		if err != nil {
			return nil, err
		}
		name := strings.TrimPrefix(source, "$response.headers.")
		return lookupResponseHeader(response, name)
	case source == "$response.body":
		response, err := currentResponse(state)
		if err != nil {
			return nil, err
		}
		return dotWalk(response, []string{"body"}), nil
	case strings.HasPrefix(source, "$response.body#"):
		response, err := currentResponse(state)
		if err != nil {
			return nil, err
		}
		body := dotWalk(response, []string{"body"})
		segments, err := Pointer(strings.TrimPrefix(source, "$response.body"))
		if err != nil {
			return nil, err
		}
		return dotWalkPointer(body, segments), nil
	case strings.HasPrefix(source, "$response.body."):
		response, err := currentResponse(state)
		if err != nil {
			return nil, err
		}
		return dotWalk(dotWalk(response, []string{"body"}), strings.Split(strings.TrimPrefix(source, "$response.body."), ".")), nil
	case strings.HasPrefix(source, "$outputs."):
		if state == nil || state.Current == nil {
			return nil, fmt.Errorf("$outputs requires a current execution context")
		}
		parts := strings.Split(strings.TrimPrefix(source, "$outputs."), ".")
		value, ok := state.Current.Outputs[parts[0]]
		if !ok {
			return nil, fmt.Errorf("$outputs.%s is unresolved in the current output scope", parts[0])
		}
		return dotWalk(value, parts[1:]), nil
	case strings.HasPrefix(source, "$steps."):
		if state == nil {
			return nil, fmt.Errorf("$steps requires an execution context")
		}
		parts := strings.Split(strings.TrimPrefix(source, "$steps."), ".")
		stepID, outputName := parts[0], parts[2]
		var matched *uws1.ExecutionRecord
		// A child iteration can read a completed step in its own iteration or
		// an enclosing iteration. Search nearest first, but never enter a
		// sibling iteration or another workflow invocation.
		for iteration := iterationSuffix(state.Current); ; iteration = parentIterationSuffix(iteration) {
			for key, record := range state.Records {
				if record.ID != stepID || !strings.HasPrefix(record.Kind, "step:") || record.Status != "success" || iterationSuffixFromKey(key) != iteration {
					continue
				}
				if strings.Contains(key, "::") {
					continue
				}
				if matched != nil {
					return nil, fmt.Errorf("$steps.%s is ambiguous in the current execution scope", stepID)
				}
				recordCopy := record
				matched = &recordCopy
			}
			if matched != nil || iteration == "" {
				break
			}
		}
		if matched == nil {
			return nil, fmt.Errorf("$steps.%s has no execution record in the current scope", stepID)
		}
		value, ok := matched.Outputs[outputName]
		if !ok {
			return nil, fmt.Errorf("$steps.%s.outputs.%s is unresolved", stepID, outputName)
		}
		return dotWalk(value, parts[3:]), nil
	case source == "$variables" || strings.HasPrefix(source, "$variables."):
		variables := map[string]any{}
		if r.document.Components != nil {
			for name, value := range r.document.Components.Variables {
				variables[name] = value
			}
		}
		for name, value := range r.document.Variables {
			variables[name] = value
		}
		if source == "$variables" {
			return variables, nil
		}
		return resolveRootPath(variables, strings.TrimPrefix(source, "$variables."), "$variables")
	case source == "$inputs" || strings.HasPrefix(source, "$inputs."):
		inputs := map[string]any{}
		if state != nil && state.Inputs != nil {
			inputs = state.Inputs
		}
		if source == "$inputs" {
			return r.resolveInputValue(ctx, inputs, depth+1)
		}
		parts := strings.Split(strings.TrimPrefix(source, "$inputs."), ".")
		value, ok := inputs[parts[0]]
		if !ok {
			return nil, fmt.Errorf("$inputs.%s is unresolved", parts[0])
		}
		value, err := r.resolveInputValue(ctx, value, depth+1)
		if err != nil {
			return nil, err
		}
		return dotWalk(value, parts[1:]), nil
	case source == "$trigger" || strings.HasPrefix(source, "$trigger."):
		if state == nil || state.Trigger == nil {
			return nil, fmt.Errorf("$trigger is unavailable outside trigger dispatch")
		}
		if source == "$trigger" {
			return state.Trigger.Payload, nil
		}
		return dotWalk(state.Trigger.Payload, strings.Split(strings.TrimPrefix(source, "$trigger."), ".")), nil
	case source == "$item" || strings.HasPrefix(source, "$item."):
		if state == nil || state.Iteration == nil {
			return nil, fmt.Errorf("$item is unavailable outside an iteration")
		}
		if source == "$item" {
			return state.Iteration.Item, nil
		}
		return dotWalk(state.Iteration.Item, strings.Split(strings.TrimPrefix(source, "$item."), ".")), nil
	case source == "$index":
		if state == nil || state.Iteration == nil {
			return nil, fmt.Errorf("$index is unavailable outside an iteration")
		}
		return state.Iteration.Index, nil
	case source == "$batchIndex":
		if state == nil || state.Iteration == nil || state.Iteration.BatchIndex < 0 {
			return nil, fmt.Errorf("$batchIndex is unavailable outside a structural loop")
		}
		return state.Iteration.BatchIndex, nil
	}
	return nil, fmt.Errorf("unsupported reference expression source %q", source)
}

func (r *Evaluator) resolveInputValue(ctx context.Context, value any, depth int) (any, error) {
	if depth > maxExpressionDepth {
		return nil, fmt.Errorf("reference expression/input nesting exceeds %d", maxExpressionDepth)
	}
	switch typed := value.(type) {
	case string:
		if strings.HasPrefix(typed, "$") {
			return r.evaluateExpression(ctx, typed, Value, depth)
		}
		return typed, nil
	case []any:
		resolved := make([]any, len(typed))
		for i, item := range typed {
			value, err := r.resolveInputValue(ctx, item, depth+1)
			if err != nil {
				return nil, fmt.Errorf("input array item %d: %w", i, err)
			}
			resolved[i] = value
		}
		return resolved, nil
	case map[string]any:
		resolved := make(map[string]any, len(typed))
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value, err := r.resolveInputValue(ctx, typed[key], depth+1)
			if err != nil {
				return nil, fmt.Errorf("input property %q: %w", key, err)
			}
			resolved[key] = value
		}
		return resolved, nil
	default:
		if typed != nil {
			switch reflect.ValueOf(typed).Kind() {
			case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct, reflect.Pointer:
				normalized, err := decodeJSONValue(typed)
				if err != nil {
					return nil, fmt.Errorf("normalize input value: %w", err)
				}
				return r.resolveInputValue(ctx, normalized, depth+1)
			}
		}
		return typed, nil
	}
}

func currentResponse(state *uws1.ExecutionContext) (any, error) {
	if state == nil || state.Current == nil || state.Current.Key == "" {
		return nil, fmt.Errorf("$response is unavailable without a current operation")
	}
	key := state.Current.Key
	if strings.Contains(key, "::") {
		prefix := state.WorkflowScope + "::"
		if state.WorkflowScope == "" || !strings.HasPrefix(key, prefix) {
			return nil, ErrContext
		}
		key = strings.TrimPrefix(key, prefix)
		if strings.Contains(key, "::") {
			return nil, ErrContext
		}
	}
	record, ok := state.Records[key]
	if !ok || record.Result == nil {
		return nil, fmt.Errorf("$response has no recorded leaf result for the current operation")
	}
	value, err := decodeJSONValue(record.Result)
	if err != nil {
		return nil, fmt.Errorf("decode current operation response: %w", err)
	}
	return value, nil
}

func lookupResponseHeader(response any, name string) (any, error) {
	headers := dotWalk(response, []string{"headers"})
	object, ok := headers.(map[string]any)
	if !ok {
		return nil, nil
	}
	var value any
	found := false
	for key, candidate := range object {
		if strings.EqualFold(key, name) {
			if found {
				return nil, fmt.Errorf("response contains ambiguous case-insensitive header %q", name)
			}
			value, found = candidate, true
		}
	}
	return value, nil
}

func resolveRootPath(root map[string]any, path, source string) (any, error) {
	parts := strings.Split(path, ".")
	value, ok := root[parts[0]]
	if !ok {
		return nil, fmt.Errorf("%s.%s is unresolved", source, parts[0])
	}
	return dotWalk(value, parts[1:]), nil
}

func dotWalk(value any, segments []string) any {
	return dotWalkWithCanonicalArrayIndex(value, segments, false)
}

func dotWalkPointer(value any, segments []string) any {
	return dotWalkWithCanonicalArrayIndex(value, segments, true)
}

func dotWalkWithCanonicalArrayIndex(value any, segments []string, canonical bool) any {
	for segmentIndex, segment := range segments {
		reflected := reflect.ValueOf(value)
		for reflected.IsValid() && (reflected.Kind() == reflect.Interface || reflected.Kind() == reflect.Pointer) {
			if reflected.IsNil() {
				return nil
			}
			reflected = reflected.Elem()
		}
		if !reflected.IsValid() {
			return nil
		}
		switch reflected.Kind() {
		case reflect.Map:
			if reflected.Type().Key().Kind() != reflect.String {
				return nil
			}
			key := reflect.ValueOf(segment).Convert(reflected.Type().Key())
			child := reflected.MapIndex(key)
			if !child.IsValid() {
				return nil
			}
			value = child.Interface()
		case reflect.Slice, reflect.Array:
			if raw, ok := value.(json.RawMessage); ok {
				decoded, err := decodeJSONValue(raw)
				if err != nil {
					return nil
				}
				return dotWalkWithCanonicalArrayIndex(decoded, segments[segmentIndex:], canonical)
			}
			if reflected.Kind() == reflect.Slice && reflected.Type().Elem().Kind() == reflect.Uint8 {
				return nil
			}
			index, err := strconv.Atoi(segment)
			if err != nil || index < 0 || index >= reflected.Len() || (canonical && strconv.Itoa(index) != segment) {
				return nil
			}
			value = reflected.Index(index).Interface()
		case reflect.Struct:
			decoded, err := decodeJSONValue(value)
			if err != nil {
				return nil
			}
			return dotWalkWithCanonicalArrayIndex(decoded, segments[segmentIndex:], canonical)
		default:
			return nil
		}
	}
	return value
}

func compareExpressionValues(left, right any, operator string) (bool, error) {
	leftKind := jsonValueKind(left)
	rightKind := jsonValueKind(right)
	if leftKind == "unsupported" || rightKind == "unsupported" {
		return false, fmt.Errorf("comparison uses a non-JSON value")
	}
	if leftKind != rightKind {
		return false, fmt.Errorf("comparison operands have different JSON types %s and %s", leftKind, rightKind)
	}
	if operator == "==" || operator == "!=" {
		var equal bool
		switch leftKind {
		case "string":
			equal = left.(string) == right.(string)
		case "number":
			comparison, err := compareNumbers(left, right)
			if err != nil {
				return false, err
			}
			equal = comparison == 0
		case "boolean":
			equal = left.(bool) == right.(bool)
		case "null":
			equal = true
		}
		if operator == "!=" {
			return !equal, nil
		}
		return equal, nil
	}
	var comparison int
	switch leftKind {
	case "string":
		comparison = strings.Compare(left.(string), right.(string))
	case "number":
		var err error
		comparison, err = compareNumbers(left, right)
		if err != nil {
			return false, err
		}
	default:
		return false, fmt.Errorf("ordering comparison %s is unsupported for %s values", operator, leftKind)
	}
	switch operator {
	case "<":
		return comparison < 0, nil
	case "<=":
		return comparison <= 0, nil
	case ">":
		return comparison > 0, nil
	case ">=":
		return comparison >= 0, nil
	default:
		return false, fmt.Errorf("unsupported comparison operator %q", operator)
	}
}

func jsonValueKind(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case json.Number, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	default:
		return "unsupported"
	}
}

func compareNumbers(left, right any) (int, error) {
	leftRat, err := numberRational(left)
	if err != nil {
		return 0, err
	}
	rightRat, err := numberRational(right)
	if err != nil {
		return 0, err
	}
	return leftRat.Cmp(rightRat), nil
}

func numberRational(value any) (*big.Rat, error) {
	var text string
	switch typed := value.(type) {
	case json.Number:
		text = typed.String()
	case float32:
		text = strconv.FormatFloat(float64(typed), 'g', -1, 32)
	case float64:
		text = strconv.FormatFloat(typed, 'g', -1, 64)
	case int:
		text = strconv.FormatInt(int64(typed), 10)
	case int8:
		text = strconv.FormatInt(int64(typed), 10)
	case int16:
		text = strconv.FormatInt(int64(typed), 10)
	case int32:
		text = strconv.FormatInt(int64(typed), 10)
	case int64:
		text = strconv.FormatInt(typed, 10)
	case uint:
		text = strconv.FormatUint(uint64(typed), 10)
	case uint8:
		text = strconv.FormatUint(uint64(typed), 10)
	case uint16:
		text = strconv.FormatUint(uint64(typed), 10)
	case uint32:
		text = strconv.FormatUint(uint64(typed), 10)
	case uint64:
		text = strconv.FormatUint(typed, 10)
	case uintptr:
		text = strconv.FormatUint(uint64(typed), 10)
	default:
		return nil, fmt.Errorf("unsupported JSON number %T", value)
	}
	if len(text) > 256 {
		return nil, fmt.Errorf("JSON number exceeds the reference comparison size limit")
	}
	if exponentIndex := strings.IndexAny(text, "eE"); exponentIndex >= 0 {
		exponent, err := strconv.Atoi(text[exponentIndex+1:])
		if err != nil || exponent < -10000 || exponent > 10000 {
			return nil, fmt.Errorf("JSON number exponent exceeds the reference comparison limit")
		}
	}
	rat, ok := new(big.Rat).SetString(text)
	if !ok {
		return nil, errors.New("invalid JSON number")
	}
	return rat, nil
}

func decodeJSONValue(value any) (any, error) {
	var encoded []byte
	switch typed := value.(type) {
	case json.RawMessage:
		encoded = bytes.Clone(typed)
	default:
		var err error
		encoded, err = json.Marshal(value)
		if err != nil {
			return nil, err
		}
	}
	if err := strictjson.ValidateSingleValue(encoded); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}

func iterationSuffix(current *uws1.CurrentExecutionContext) string {
	if current == nil {
		return ""
	}
	return iterationSuffixFromKey(current.Key)
}

func iterationSuffixFromKey(key string) string {
	if index := strings.LastIndex(key, "#iter:"); index >= 0 {
		return key[index:]
	}
	return ""
}

func parentIterationSuffix(iteration string) string {
	if index := strings.LastIndex(iteration, "."); index >= 0 {
		return iteration[:index]
	}
	return ""
}
