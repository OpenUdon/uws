package binding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/OpenUdon/uws/expressions"
	"github.com/OpenUdon/uws/internal/strictjson"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

type Outcome string

const (
	Compatible    Outcome = "compatible"
	Incompatible  Outcome = "incompatible"
	Indeterminate Outcome = "indeterminate"
)

// Diagnostic carries stable metadata only; values/resolver/schema errors aren't exposed.
type Diagnostic struct {
	Code    string  `json:"code"`
	Path    string  `json:"path"`
	Outcome Outcome `json:"outcome"`
}
type BoundInput struct {
	Location string
	Name     string
	Value    any
}
type SecurityBinding struct {
	Scheme         string
	CredentialSlot string
	Scopes         []string
}
type OutputReference struct {
	Location string
	Name     string
	Pointer  string
}

// Request is a consumer-prepared projection. ExpressionTypes must come from
// independently reviewed input/output contracts, not a model's trust assertion.
type Request struct {
	Binding           Binding
	Inputs            []BoundInput
	Security          []SecurityBinding
	OutputReferences  []OutputReference
	ExpressionTypes   map[string]Schema
	ExpressionContext expressions.Context
}
type Report struct {
	Outcome     Outcome      `json:"outcome"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

func pointer(s string) string { return strings.ReplaceAll(strings.ReplaceAll(s, "~", "~0"), "/", "~1") }
func (r *Report) add(code, path string, outcome Outcome) {
	if len(r.Diagnostics) < 128 {
		r.Diagnostics = append(r.Diagnostics, Diagnostic{Code: code, Path: path, Outcome: outcome})
	}
	if outcome == Incompatible {
		r.Outcome = Incompatible
	} else if outcome == Indeterminate && r.Outcome == Compatible {
		r.Outcome = Indeterminate
	}
}

// ValidateBinding checks metadata, never credentials or runtime permission.
func ValidateBinding(ctx context.Context, resolver Resolver, request Request) (Report, error) {
	report := Report{Outcome: Compatible, Diagnostics: []Diagnostic{}}
	if ctx == nil {
		return report, ErrTable
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	if resolver == nil || len(request.Inputs) > 2048 || len(request.Security) > 128 || len(request.OutputReferences) > 2048 || len(request.ExpressionTypes) > 2048 {
		return report, ErrTable
	}
	resolved, err := resolver.Resolve(ctx, request.Binding)
	if err != nil {
		if ctx.Err() != nil {
			return report, ctx.Err()
		}
		report.add("binding.resolver_failure", "/binding", Indeterminate)
		return report, nil
	}
	if resolved.Status != Resolved || resolved.Shape == nil {
		code := "binding.operation_unsupported"
		state := Indeterminate
		if resolved.Status == Missing {
			code = "binding.operation_missing"
			state = Incompatible
		}
		if resolved.Status == Ambiguous {
			code = "binding.operation_ambiguous"
		}
		report.add(code, "/binding", state)
		return report, nil
	}
	shape := resolved.Shape
	// Validate custom resolver projections too; never accept a malformed claim.
	if (ShapeTable{Version: TableVersion, Sources: []Source{shape.Source}, Operations: []OperationShape{*shape}}).Validate() != nil || shape.Source != request.Binding.Source || !selectorMatches(*shape, request.Binding) {
		report.add("binding.resolver_contract", "/binding", Indeterminate)
		return report, nil
	}
	if !shape.Complete {
		report.add("binding.operation_incomplete", "/binding", Indeterminate)
	}
	supplied := map[string]BoundInput{}
	for i, in := range request.Inputs {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		key := in.Location + "\x00" + in.Name
		path := "/inputs/" + strconv.Itoa(i)
		if !text(in.Location, 128, true) || !text(in.Name, 1024, true) {
			return report, ErrTable
		}
		if _, ok := supplied[key]; ok {
			report.add("binding.input_duplicate", path, Incompatible)
		}
		supplied[key] = in
	}
	declared := map[string]bool{}
	for _, in := range shape.Inputs {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		key := in.Location + "\x00" + in.Name
		declared[key] = true
		path := "/inputs/" + pointer(in.Location) + "/" + pointer(in.Name)
		value, present := supplied[key]
		if !present {
			if in.Required {
				report.add("binding.input_required", path, Incompatible)
			}
			continue
		}
		if !in.Schema.Known {
			report.add("binding.schema_unknown", path, Indeterminate)
			continue
		}
		if expression, isExpression := value.Value.(string); isExpression && (strings.HasPrefix(expression, "$") || strings.HasPrefix(strings.TrimSpace(expression), "expr(")) {
			if _, err := expressions.Parse(expression, request.ExpressionContext); err != nil {
				report.add("binding.expression_unsupported", path, Indeterminate)
				continue
			}
			source, ok := request.ExpressionTypes[expression]
			if !ok || !source.Known || !schemaValid(source) {
				report.add("binding.expression_type_unknown", path, Indeterminate)
				continue
			}
			state := schemaCompatibility(source, in.Schema)
			if state != Compatible {
				code := "binding.expression_type_indeterminate"
				if state == Incompatible {
					code = "binding.expression_type_mismatch"
				}
				report.add(code, path, state)
			}
		} else {
			state := validateBoundValue(in.Schema, value.Value, request.ExpressionTypes, request.ExpressionContext)
			if state != Compatible {
				code := "binding.literal_schema_indeterminate"
				if state == Incompatible {
					code = "binding.literal_schema_mismatch"
				}
				report.add(code, path, state)
			}
		}
	}
	keys := make([]string, 0, len(supplied))
	for key := range supplied {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if !declared[key] {
			in := supplied[key]
			report.add("binding.input_undeclared", "/inputs/"+pointer(in.Location)+"/"+pointer(in.Name), Indeterminate)
		}
	}
	checkSecurity(&report, shape.Security, request.Security)
	for i, ref := range request.OutputReferences {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		path := "/output_references/" + strconv.Itoa(i)
		var output *Output
		for n := range shape.Outputs {
			if shape.Outputs[n].Location == ref.Location && shape.Outputs[n].Name == ref.Name {
				output = &shape.Outputs[n]
				break
			}
		}
		if output == nil {
			state := Incompatible
			if !shape.Complete {
				state = Indeterminate
			}
			report.add("binding.output_missing", path, state)
			continue
		}
		if !output.Schema.Known {
			report.add("binding.output_schema_unknown", path, Indeterminate)
			continue
		}
		state := schemaPath(output.Schema, ref.Pointer)
		if state != Compatible {
			code := "binding.output_field_indeterminate"
			if state == Incompatible {
				code = "binding.output_field_missing"
			}
			report.add(code, path, state)
		}
	}
	return report, nil
}
func checkSecurity(report *Report, security Security, bindings []SecurityBinding) {
	if !security.Known {
		report.add("binding.security_unknown", "/security", Indeterminate)
		return
	}
	byScheme := map[string]SecurityBinding{}
	for i, b := range bindings {
		if !text(b.Scheme, 256, true) || !text(b.CredentialSlot, 256, true) || len(b.Scopes) > 128 {
			report.add("binding.security_symbol_invalid", "/security/"+strconv.Itoa(i), Incompatible)
			continue
		}
		if _, ok := byScheme[b.Scheme]; ok {
			report.add("binding.security_duplicate", "/security/"+strconv.Itoa(i), Incompatible)
		}
		byScheme[b.Scheme] = b
	}
	for _, alternative := range security.Alternatives {
		satisfied := true
		for _, req := range alternative.Requirements {
			bound, ok := byScheme[req.Scheme]
			if !ok {
				satisfied = false
				break
			}
			scopes := map[string]bool{}
			for _, scope := range bound.Scopes {
				scopes[scope] = true
			}
			for _, scope := range req.Scopes {
				if !scopes[scope] {
					satisfied = false
				}
			}
		}
		if satisfied {
			return
		}
	}
	report.add("binding.security_missing", "/security", Incompatible)
}

type denyLoader struct{}

func (denyLoader) Load(string) (any, error) {
	return nil, errors.New("schema external resources are unavailable")
}
func decode(data []byte) (any, error) {
	if len(data) > MaxSchemaBytes || strictjson.ValidateSingleValue(data) != nil {
		return nil, ErrTable
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var v any
	err := d.Decode(&v)
	return v, err
}
func compile(schema Schema) (*jsonschema.Schema, error) {
	if !schema.Known || !schemaValid(schema) {
		return nil, ErrTable
	}
	value, err := decode(schema.JSON)
	if err != nil {
		return nil, err
	}
	c := jsonschema.NewCompiler()
	c.UseLoader(denyLoader{})
	c.AssertFormat()
	if err = c.AddResource("https://uws.invalid/shape", value); err != nil {
		return nil, err
	}
	return c.Compile("https://uws.invalid/shape")
}
func validateLiteral(schema Schema, value any) Outcome {
	compiled, err := compile(schema)
	if err != nil {
		return Indeterminate
	}
	data, err := json.Marshal(value)
	if err != nil || len(data) > MaxTableBytes {
		return Incompatible
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var normalized any
	if err = d.Decode(&normalized); err != nil {
		return Incompatible
	}
	if compiled.Validate(normalized) != nil {
		return Incompatible
	}
	return Compatible
}
func schemaCompatibility(source, target Schema) Outcome {
	s, err := decode(source.JSON)
	if err != nil {
		return Indeterminate
	}
	t, err := decode(target.JSON)
	if err != nil {
		return Indeterminate
	}
	cs, err := compile(source)
	if err != nil {
		return Indeterminate
	}
	ct, err := compile(target)
	if err != nil {
		return Indeterminate
	}
	if reflect.DeepEqual(s, t) {
		return Compatible
	}
	return schemaContainment(s, t, cs, ct, 0)
}

// Child nodes stay attached to the original compiled resource. Recompiling raw
// child maps would lose inherited dialect, local references and resource IDs.
func schemaContainment(s, t any, cs, ct *jsonschema.Schema, depth int) Outcome {
	if cs == nil || ct == nil || depth > 32 {
		return Indeterminate
	}
	if cs.Bool != nil && !*cs.Bool {
		return Compatible
	}
	if ct.Bool != nil {
		if *ct.Bool {
			return Compatible
		}
		return Incompatible
	}
	if cs.Ref != nil || cs.DynamicRef != nil || cs.RecursiveRef != nil {
		return Indeterminate
	}
	sm, ok := s.(map[string]any)
	if !ok {
		return Indeterminate
	}
	if cs.Const != nil {
		if ct.Validate(*cs.Const) == nil {
			return Compatible
		}
		return Incompatible
	}
	if cs.Enum != nil && len(cs.Enum.Values) > 0 {
		state := Compatible
		for _, v := range cs.Enum.Values {
			if ct.Validate(v) != nil {
				state = Indeterminate
			}
		}
		return state
	}
	tm, ok := t.(map[string]any)
	if !ok {
		if t == true {
			return Compatible
		}
		return Indeterminate
	}
	from, to := types(sm["type"]), types(tm["type"])
	if len(from) == 0 || len(to) == 0 {
		return Indeterminate
	}
	subset, overlap := true, false
	for _, a := range from {
		found := false
		for _, b := range to {
			if a == b || a == "integer" && b == "number" {
				found = true
				overlap = true
			}
			if a == "number" && b == "integer" {
				overlap = true
			}
		}
		if !found {
			subset = false
		}
	}
	if !overlap {
		return Incompatible
	}
	// Type containment proves compatibility only when target has no constraints
	// beyond type/annotation. Other containment questions stay indeterminate.
	typeOnly := true
	for key := range tm {
		switch key {
		case "type", "title", "description", "$schema", "$id", "$defs", "definitions", "$comment":
		default:
			typeOnly = false
		}
	}
	if subset && typeOnly {
		return Compatible
	}
	if subset && len(from) == 1 && len(to) == 1 && from[0] == "object" && to[0] == "object" {
		return objectContainment(sm, tm, cs, ct, depth+1)
	}
	if subset && len(from) == 1 && len(to) == 1 && from[0] == "array" && to[0] == "array" {
		return arrayContainment(sm, tm, cs, ct, depth+1)
	}
	return Indeterminate
}
func types(value any) []string {
	if s, ok := value.(string); ok {
		return []string{s}
	}
	if a, ok := value.([]any); ok {
		out := make([]string, 0, len(a))
		for _, x := range a {
			if s, ok := x.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}
func schemaPath(schema Schema, fragment string) Outcome {
	node, err := compile(schema)
	if err != nil {
		return Indeterminate
	}
	if bytes.Equal(bytes.TrimSpace(schema.JSON), []byte("false")) {
		return Incompatible
	}
	if fragment == "" {
		return Compatible
	}
	parts, err := expressions.Pointer(fragment)
	if err != nil {
		return Incompatible
	}
	for _, part := range parts {
		if node.Bool != nil && !*node.Bool {
			return Incompatible
		}
		if node.Ref != nil || node.DynamicRef != nil || node.RecursiveRef != nil || len(node.AnyOf) > 0 || len(node.OneOf) > 0 || len(node.AllOf) > 0 {
			return Indeterminate
		}
		if node.Types == nil {
			return Indeterminate
		}
		kinds := node.Types.ToStrings()
		if len(kinds) != 1 {
			return Indeterminate
		}
		switch kinds[0] {
		case "object":
			// A declared property and every matching pattern constrain the same
			// value. Do not discard that intersection when walking a property.
			for pattern := range node.PatternProperties {
				if pattern.MatchString(part) {
					return Indeterminate
				}
			}
			child, ok := node.Properties[part]
			if !ok {
				if len(node.PatternProperties) > 0 {
					return Indeterminate
				}
				if node.AdditionalProperties == false {
					return Incompatible
				}
				return Indeterminate
			}
			node = child
		case "array":
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || strconv.Itoa(index) != part {
				return Incompatible
			}
			var child *jsonschema.Schema
			if node.DraftVersion >= 2020 {
				if index < len(node.PrefixItems) {
					child = node.PrefixItems[index]
				} else {
					child = node.Items2020
				}
			} else if tuple, ok := node.Items.([]*jsonschema.Schema); ok {
				if index < len(tuple) {
					child = tuple[index]
				} else {
					if node.AdditionalItems == false {
						return Incompatible
					}
					child, _ = node.AdditionalItems.(*jsonschema.Schema)
				}
			} else {
				child, _ = node.Items.(*jsonschema.Schema)
			}
			if child == nil {
				return Indeterminate
			}
			node = child
		default:
			return Incompatible
		}
	}
	if node.Bool != nil && !*node.Bool {
		return Incompatible
	}
	if node.Ref != nil || node.DynamicRef != nil || node.RecursiveRef != nil {
		return Indeterminate
	}
	return Compatible
}

func selectorMatches(shape OperationShape, b Binding) bool {
	if shape.Selector.Kind == b.SelectorKind && shape.Selector.Value == b.SelectorValue {
		return true
	}
	for _, alias := range shape.Aliases {
		if alias.Kind == b.SelectorKind && alias.Value == b.SelectorValue {
			return true
		}
	}
	return false
}
