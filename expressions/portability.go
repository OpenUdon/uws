package expressions

import (
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/OpenUdon/uws/uws1"
)

// Diagnostic identifies a core field and stable reason without expression data.
type Diagnostic struct {
	Code string `json:"code"`
	Path string `json:"path"`
}

// CheckPortability is opt-in. It examines only UWS core expression fields;
// profile bodies, extensions, trigger options and function/browser strings are
// never reinterpreted. It does not call or change ordinary document validation.
func CheckPortability(document *uws1.Document) []Diagnostic {
	if document == nil {
		return []Diagnostic{{Code: "expression.document", Path: "/"}}
	}
	c := checker{version: document.UWS}
	opScopes, workflowScopes := invocationScopes(document)
	for i, op := range document.Operations {
		if op == nil {
			continue
		}
		base := "/operations/" + strconv.Itoa(i)
		loop := opScopes[op.OperationID] == 2
		c.check(op.ForEach, base+"/forEach", Value, loop)
		if op.ForEach != "" {
			loop = false
		}
		c.check(op.When, base+"/when", Predicate, loop)
		c.check(op.Wait, base+"/wait", Wait, loop)
		// Core source request bindings are expressions only when explicitly marked
		// with a source or legacy wrapper. Extension-owned request templates aren't.
		browser := false
		for _, source := range document.SourceDescriptions {
			if source != nil && source.Name == op.SourceDescription && source.Type == "browser-profile" {
				browser = true
			}
		}
		if op.HasSourceBinding() && !browser {
			c.values(op.Request, base+"/request", loop, 0)
		}
		c.outputs(op.Outputs, base+"/outputs", loop)
		c.criteria(op.SuccessCriteria, base+"/successCriteria", loop)
		for n, a := range op.OnSuccess {
			if a != nil {
				c.criteria(a.Criteria, base+"/onSuccess/"+strconv.Itoa(n)+"/criteria", loop)
			}
		}
		for n, a := range op.OnFailure {
			if a != nil {
				c.criteria(a.Criteria, base+"/onFailure/"+strconv.Itoa(n)+"/criteria", loop)
			}
		}
	}
	for i, w := range document.Workflows {
		if w == nil {
			continue
		}
		base := "/workflows/" + strconv.Itoa(i)
		loop := workflowScopes[w.WorkflowID] == 2
		c.check(w.ForEach, base+"/forEach", Value, loop)
		if w.ForEach != "" {
			loop = false
		}
		c.check(w.When, base+"/when", Predicate, loop)
		c.structural(w.Type, w.Wait, w.Items, w.BatchSize, base, loop)
		bodyLoop := loop || w.Type == uws1.WorkflowTypeLoop
		c.outputs(w.Outputs, base+"/outputs", loop)
		c.steps(w.Steps, base+"/steps", bodyLoop, 0)
		c.steps(w.Default, base+"/default", bodyLoop, 0)
		c.cases(w.Cases, base+"/cases", bodyLoop, 0)
	}
	return c.diagnostics
}

type checker struct {
	version     string
	diagnostics []Diagnostic
}

func (c *checker) add(code, path string) {
	if len(c.diagnostics) < 128 {
		c.diagnostics = append(c.diagnostics, Diagnostic{Code: code, Path: path})
	}
}
func (c *checker) check(text, path string, field Field, loop bool) {
	if text == "" {
		return
	}
	if strings.HasPrefix(strings.TrimSpace(text), "expr(") {
		c.add("expression.legacy-wrapper", path)
		return
	}
	_, err := Parse(text, Context{Version: c.version, Field: field, InLoop: loop})
	if err == nil {
		return
	}
	code := "expression.syntax"
	switch {
	case errors.Is(err, ErrVersion):
		code = "expression.version"
	case errors.Is(err, ErrContext):
		code = "expression.context"
	case errors.Is(err, ErrLimit):
		code = "expression.limit"
	}
	c.add(code, path)
}
func (c *checker) outputs(values map[string]string, path string, loop bool) {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if values[key] == "" {
			c.add("expression.syntax", path+"/"+pointerKey(key))
			continue
		}
		c.check(values[key], path+"/"+pointerKey(key), Value, loop)
	}
}
func pointerKey(key string) string {
	return strings.ReplaceAll(strings.ReplaceAll(key, "~", "~0"), "/", "~1")
}
func (c *checker) values(value any, path string, loop bool, depth int) {
	if depth > maxExpressionDepth {
		c.add("expression.depth", path)
		return
	}
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "$") || strings.HasPrefix(strings.TrimSpace(v), "expr(") {
			c.check(v, path, Value, loop)
		}
	case map[string]any:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			c.values(v[key], path+"/"+pointerKey(key), loop, depth+1)
		}
	case []any:
		for i, item := range v {
			c.values(item, path+"/"+strconv.Itoa(i), loop, depth+1)
		}
	default:
		if value == nil {
			return
		}
		kind := reflect.ValueOf(value).Kind()
		if kind == reflect.String {
			c.values(reflect.ValueOf(value).String(), path, loop, depth+1)
			return
		}
		if kind == reflect.Map || kind == reflect.Slice || kind == reflect.Array || kind == reflect.Struct || kind == reflect.Pointer {
			normalized, err := decodeJSONValue(value)
			if err != nil {
				c.add("expression.value", path)
				return
			}
			c.values(normalized, path, loop, depth+1)
		}
	}
}
func (c *checker) criteria(criteria []*uws1.Criterion, path string, loop bool) {
	for i, criterion := range criteria {
		if criterion == nil {
			continue
		}
		base := path + "/" + strconv.Itoa(i)
		c.check(criterion.Context, base+"/context", Value, loop)
		if criterion.Type == "" || criterion.Type == uws1.CriterionSimple {
			c.check(criterion.Condition, base+"/condition", Predicate, loop)
		}
	}
}
func (c *checker) structural(kind, wait, items, batch, path string, loop bool) {
	field := Wait
	if kind == uws1.WorkflowTypeAwait {
		field = Predicate
	}
	c.check(wait, path+"/wait", field, loop)
	c.check(items, path+"/items", Value, loop)
	c.check(batch, path+"/batchSize", BatchSize, loop)
}
func (c *checker) steps(steps []*uws1.Step, path string, loop bool, depth int) {
	if depth > maxExpressionDepth {
		c.add("expression.depth", path)
		return
	}
	for i, step := range steps {
		if step == nil {
			continue
		}
		base := path + "/" + strconv.Itoa(i)
		c.check(step.ForEach, base+"/forEach", Value, loop)
		effective := loop
		if step.ForEach != "" {
			effective = false
		}
		c.check(step.When, base+"/when", Predicate, effective)
		c.structural(step.Type, step.Wait, step.Items, step.BatchSize, base, effective)
		childLoop := effective || step.Type == uws1.WorkflowTypeLoop
		c.values(step.Inputs, base+"/inputs", effective, 0)
		c.outputs(step.Outputs, base+"/outputs", effective)
		c.steps(step.Steps, base+"/steps", childLoop, depth+1)
		c.steps(step.Default, base+"/default", childLoop, depth+1)
		c.cases(step.Cases, base+"/cases", childLoop, depth+1)
	}
}
func (c *checker) cases(cases []*uws1.Case, path string, loop bool, depth int) {
	if depth > maxExpressionDepth {
		c.add("expression.depth", path)
		return
	}
	for i, branch := range cases {
		if branch != nil {
			base := path + "/" + strconv.Itoa(i)
			c.check(branch.When, base+"/when", Predicate, loop)
			c.steps(branch.Steps, base+"/steps", loop, depth+1)
		}
	}
}

// invocationScopes records known invocation contexts, not authorization or
// reachability guarantees. A reusable operation must be valid in every known
// call context; unused declarations are checked conservatively outside a loop.
func invocationScopes(d *uws1.Document) (map[string]uint8, map[string]uint8) {
	operations := map[string]uint8{}
	workflows := map[string]uint8{}
	byID := map[string]*uws1.Workflow{}
	opByID := map[string]*uws1.Operation{}
	for _, w := range d.Workflows {
		if w != nil {
			byID[w.WorkflowID] = w
		}
	}
	for _, op := range d.Operations {
		if op != nil {
			opByID[op.OperationID] = op
		}
	}
	bit := func(loop bool) uint8 {
		if loop {
			return 2
		}
		return 1
	}
	var visitWorkflow func(string, bool, int)
	var visitSteps func([]*uws1.Step, bool, int)
	var visitOperation func(string, bool, int)
	visitOperation = func(id string, loop bool, depth int) {
		if depth > maxExpressionDepth || operations[id]&bit(loop) != 0 {
			return
		}
		operations[id] |= bit(loop)
		if op := opByID[id]; op != nil {
			effective := loop
			if op.ForEach != "" {
				effective = false
			}
			for _, dependency := range op.DependsOn {
				if opByID[dependency] != nil {
					visitOperation(dependency, effective, depth+1)
				}
			}
		}
	}
	visitSteps = func(steps []*uws1.Step, loop bool, depth int) {
		if depth > maxExpressionDepth {
			return
		}
		for _, s := range steps {
			if s == nil {
				continue
			}
			effective := loop
			if s.ForEach != "" {
				effective = false
			}
			if s.OperationRef != "" {
				visitOperation(s.OperationRef, effective, depth+1)
			}
			if s.Workflow != "" {
				visitWorkflow(s.Workflow, effective, depth+1)
			}
			body := effective || s.Type == uws1.WorkflowTypeLoop
			visitSteps(s.Steps, body, depth+1)
			visitSteps(s.Default, body, depth+1)
			for _, branch := range s.Cases {
				if branch != nil {
					visitSteps(branch.Steps, body, depth+1)
				}
			}
		}
	}
	visitWorkflow = func(id string, loop bool, depth int) {
		if depth > maxExpressionDepth || workflows[id]&bit(loop) != 0 {
			return
		}
		workflows[id] |= bit(loop)
		w := byID[id]
		if w == nil {
			return
		}
		effective := loop
		if w.ForEach != "" {
			effective = false
		}
		body := effective || w.Type == uws1.WorkflowTypeLoop
		visitSteps(w.Steps, body, depth+1)
		visitSteps(w.Default, body, depth+1)
		for _, branch := range w.Cases {
			if branch != nil {
				visitSteps(branch.Steps, body, depth+1)
			}
		}
	}
	entry := "main"
	if byID[entry] == nil && len(byID) == 1 {
		for id := range byID {
			entry = id
		}
	}
	if byID[entry] != nil {
		visitWorkflow(entry, false, 0)
	}
	return operations, workflows
}
