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
	opScopes, workflowScopes, stepScopes := invocationScopes(document)
	c.stepScopes = stepScopes
	for i, op := range document.Operations {
		if op == nil {
			continue
		}
		base := "/operations/" + strconv.Itoa(i)
		loop := onlyLoop(opScopes[op.OperationID])
		c.iteration = onlyIteration(opScopes[op.OperationID])
		c.check(op.ForEach, base+"/forEach", Value, loop)
		c.check(op.When, base+"/when", Predicate, loop)
		c.check(op.Wait, base+"/wait", Wait, loop)
		if op.ForEach != "" {
			loop = false
			c.iteration = true
		}
		// Core source request bindings are expressions only when explicitly marked
		// with a source or legacy wrapper. Extension-owned request templates aren't.
		browser := false
		for _, source := range document.SourceDescriptions {
			if source != nil && source.Name == op.SourceDescription && source.Type == "browser-profile" {
				browser = true
			}
		}
		if op.HasSourceBinding() && !browser {
			for _, binding := range []string{"path", "query", "header", "cookie", "body"} {
				if value, ok := op.Request[binding]; ok {
					c.values(value, base+"/request/"+binding, loop, 0)
				}
			}
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
		loop := onlyLoop(workflowScopes[w.WorkflowID])
		c.iteration = onlyIteration(workflowScopes[w.WorkflowID])
		c.check(w.ForEach, base+"/forEach", Value, loop)
		c.check(w.When, base+"/when", Predicate, loop)
		bodyWait := w.Wait
		if w.Type != uws1.WorkflowTypeAwait {
			c.check(w.Wait, base+"/wait", Wait, loop)
			bodyWait = ""
		}
		if w.ForEach != "" {
			loop = false
			c.iteration = true
		}
		c.structural(w.Type, bodyWait, w.Items, w.BatchSize, base, loop)
		bodyLoop := loop || w.Type == uws1.WorkflowTypeLoop
		c.outputs(w.Outputs, base+"/outputs", loop)
		c.iteration = c.iteration || w.Type == uws1.WorkflowTypeLoop
		c.steps(w.Steps, base+"/steps", bodyLoop, 0)
		c.steps(w.Default, base+"/default", bodyLoop, 0)
		c.cases(w.Cases, base+"/cases", bodyLoop, 0)
	}
	for i, result := range document.Results {
		if result != nil {
			c.iteration = result.Kind == uws1.WorkflowTypeLoop
			c.check(result.Value, "/results/"+strconv.Itoa(i)+"/value", Value, result.Kind == uws1.WorkflowTypeLoop)
		}
	}
	return c.diagnostics
}

type checker struct {
	version     string
	iteration   bool
	stepScopes  map[*uws1.Step]uint8
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
	parsed, err := Parse(text, Context{Version: c.version, Field: field, InLoop: loop})
	if err == nil {
		for _, source := range []string{parsed.Source(), parsed.OperandSource()} {
			if !c.iteration && (source == "$item" || source == "$index" || strings.HasPrefix(source, "$item.")) {
				c.add("expression.context", path)
				break
			}
		}
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
		outerIteration := c.iteration
		stepLoop := loop
		if scope := c.stepScopes[step]; scope != 0 {
			stepLoop = onlyLoop(scope)
			c.iteration = onlyIteration(scope)
		}
		c.check(step.ForEach, base+"/forEach", Value, stepLoop)
		c.check(step.When, base+"/when", Predicate, stepLoop)
		bodyWait := step.Wait
		if step.Type != uws1.WorkflowTypeAwait {
			c.check(step.Wait, base+"/wait", Wait, stepLoop)
			bodyWait = ""
		}
		effective := stepLoop
		if step.ForEach != "" {
			effective = false
			c.iteration = true
		}
		c.structural(step.Type, bodyWait, step.Items, step.BatchSize, base, effective)
		childLoop := effective || step.Type == uws1.WorkflowTypeLoop
		c.values(step.Inputs, base+"/inputs", effective, 0)
		c.outputs(step.Outputs, base+"/outputs", effective)
		c.iteration = c.iteration || step.Type == uws1.WorkflowTypeLoop
		c.steps(step.Steps, base+"/steps", childLoop, depth+1)
		c.steps(step.Default, base+"/default", childLoop, depth+1)
		c.cases(step.Cases, base+"/cases", childLoop, depth+1)
		c.iteration = outerIteration
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
func invocationScopes(d *uws1.Document) (map[string]uint8, map[string]uint8, map[*uws1.Step]uint8) {
	operations := map[string]uint8{}
	workflows := map[string]uint8{}
	stepScopes := map[*uws1.Step]uint8{}
	byID := map[string]*uws1.Workflow{}
	opByID := map[string]*uws1.Operation{}
	stepsByID := map[string]*uws1.Step{}
	groups := map[string][]string{}
	indexed := map[*uws1.Step]bool{}
	var indexSteps func([]*uws1.Step, int)
	indexSteps = func(steps []*uws1.Step, depth int) {
		if depth > maxExpressionDepth {
			return
		}
		for _, s := range steps {
			if s == nil || indexed[s] {
				continue
			}
			indexed[s] = true
			stepsByID[s.StepID] = s
			if s.ParallelGroup != "" {
				groups[s.ParallelGroup] = append(groups[s.ParallelGroup], s.StepID)
			}
			indexSteps(s.Steps, depth+1)
			indexSteps(s.Default, depth+1)
			for _, branch := range s.Cases {
				if branch != nil {
					indexSteps(branch.Steps, depth+1)
				}
			}
		}
	}
	for _, w := range d.Workflows {
		if w != nil {
			byID[w.WorkflowID] = w
			indexSteps(w.Steps, 0)
			indexSteps(w.Default, 0)
			for _, branch := range w.Cases {
				if branch != nil {
					indexSteps(branch.Steps, 0)
				}
			}
		}
	}
	for _, op := range d.Operations {
		if op != nil {
			opByID[op.OperationID] = op
			if op.ParallelGroup != "" {
				groups[op.ParallelGroup] = append(groups[op.ParallelGroup], op.OperationID)
			}
		}
	}
	bit := func(loop, iteration bool) uint8 {
		index := 0
		if loop {
			index++
		}
		if iteration {
			index += 2
		}
		return 1 << index
	}
	var visitWorkflow func(string, bool, bool, int)
	var visitSteps func([]*uws1.Step, bool, bool, int)
	var visitOperation func(string, bool, bool, int)
	var visitDependency func(string, bool, bool, int)
	visitDependency = func(id string, loop, iteration bool, depth int) {
		if depth > maxExpressionDepth {
			return
		}
		if members := groups[id]; len(members) > 0 {
			for _, member := range members {
				visitDependency(member, loop, iteration, depth+1)
			}
		} else if s := stepsByID[id]; s != nil {
			visitSteps([]*uws1.Step{s}, loop, iteration, depth+1)
		} else if byID[id] != nil {
			visitWorkflow(id, loop, iteration, depth+1)
		} else if opByID[id] != nil {
			visitOperation(id, loop, iteration, depth+1)
		}
	}
	visitOperation = func(id string, loop, iteration bool, depth int) {
		if depth > maxExpressionDepth || operations[id]&bit(loop, iteration) != 0 {
			return
		}
		operations[id] |= bit(loop, iteration)
		if op := opByID[id]; op != nil {
			for _, dependency := range op.DependsOn {
				visitDependency(dependency, loop, iteration, depth+1)
			}
			// Terminal goto unwinds the caller and invokes exact root targets.
			visitTarget := func(workflow, step string) {
				if workflow != "" {
					visitWorkflow(workflow, false, false, depth+1)
				}
				if step != "" {
					entry := byID["main"]
					if entry == nil && len(byID) == 1 {
						for _, w := range byID {
							entry = w
						}
					}
					if entry != nil {
						for _, s := range entry.Steps {
							if s != nil && s.StepID == step {
								visitSteps([]*uws1.Step{s}, false, false, depth+1)
							}
						}
					}
				}
			}
			for _, action := range op.OnSuccess {
				if action != nil && action.Type == "goto" {
					visitTarget(action.WorkflowID, action.StepID)
				}
			}
			for _, action := range op.OnFailure {
				if action != nil && action.Type == "goto" {
					visitTarget(action.WorkflowID, action.StepID)
				}
			}
		}
	}
	visitSteps = func(steps []*uws1.Step, loop, iteration bool, depth int) {
		if depth > maxExpressionDepth {
			return
		}
		for _, s := range steps {
			if s == nil || stepScopes[s]&bit(loop, iteration) != 0 {
				continue
			}
			stepScopes[s] |= bit(loop, iteration)
			for _, dependency := range s.DependsOn {
				visitDependency(dependency, loop, iteration, depth+1)
			}
			effective := loop
			bodyIteration := iteration
			if s.ForEach != "" {
				effective = false
				bodyIteration = true
			}
			if s.OperationRef != "" {
				visitOperation(s.OperationRef, effective, bodyIteration, depth+1)
			}
			if s.Workflow != "" {
				visitWorkflow(s.Workflow, effective, bodyIteration, depth+1)
			}
			body := effective || s.Type == uws1.WorkflowTypeLoop
			bodyIteration = bodyIteration || s.Type == uws1.WorkflowTypeLoop
			visitSteps(s.Steps, body, bodyIteration, depth+1)
			visitSteps(s.Default, body, bodyIteration, depth+1)
			for _, branch := range s.Cases {
				if branch != nil {
					visitSteps(branch.Steps, body, bodyIteration, depth+1)
				}
			}
		}
	}
	visitWorkflow = func(id string, loop, iteration bool, depth int) {
		if depth > maxExpressionDepth || workflows[id]&bit(loop, iteration) != 0 {
			return
		}
		workflows[id] |= bit(loop, iteration)
		w := byID[id]
		if w == nil {
			return
		}
		for _, dependency := range w.DependsOn {
			visitDependency(dependency, loop, iteration, depth+1)
		}
		effective := loop
		if w.ForEach != "" {
			effective = false
			iteration = true
		}
		body := effective || w.Type == uws1.WorkflowTypeLoop
		iteration = iteration || w.Type == uws1.WorkflowTypeLoop
		visitSteps(w.Steps, body, iteration, depth+1)
		visitSteps(w.Default, body, iteration, depth+1)
		for _, branch := range w.Cases {
			if branch != nil {
				visitSteps(branch.Steps, body, iteration, depth+1)
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
		visitWorkflow(entry, false, false, 0)
	}
	// Trigger workflows start without inherited iteration state. Direct route
	// steps are top-level entry steps and likewise bypass enclosing loop bodies.
	for _, trigger := range d.Triggers {
		if trigger == nil {
			continue
		}
		for _, route := range trigger.Routes {
			if route == nil {
				continue
			}
			for _, target := range route.To {
				if byID[target] != nil {
					visitWorkflow(target, false, false, 0)
					continue
				}
				if root := byID[entry]; root != nil {
					for _, s := range root.Steps {
						if s != nil && s.StepID == target {
							visitSteps([]*uws1.Step{s}, false, false, 0)
						}
					}
				}
			}
		}
	}
	return operations, workflows, stepScopes
}

func onlyLoop(scopes uint8) bool      { return scopes != 0 && scopes&5 == 0 }
func onlyIteration(scopes uint8) bool { return scopes != 0 && scopes&3 == 0 }
