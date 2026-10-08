package binding

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/OpenUdon/uws/expressions"
	"github.com/OpenUdon/uws/uws1"
)

// FlowFinding is an advisory code/path projection, never a runtime outcome.
type FlowFinding struct {
	Code string `json:"code"`
	Path string `json:"path"`
}
type FlowReport struct {
	Findings  []FlowFinding `json:"findings"`
	Truncated bool          `json:"truncated"`
}
type flowNode struct {
	key, path, workflow                 string
	children, dependencies, expressions []string
	transfers                           []string
	controlContexts                     map[string]bool
	outputs                             map[string]string
	op                                  *uws1.Operation
	step                                *uws1.Step
	wf                                  *uws1.Workflow
}
type flowAnalyzer struct {
	ctx         context.Context
	doc         *uws1.Document
	nodes       map[string]*flowNode
	ambiguous   map[string]bool
	stepsByName map[string][]string
	groups      map[string][]string
	contexts    map[string]map[string]bool
	reached     map[string]bool
	used        map[string]bool
	findings    map[string]FlowFinding
}

// AnalyzeFlow observes possible control flow. Conditions aren't executed;
// profile-only references remain opaque. No finding authorizes or changes a run.
func AnalyzeFlow(ctx context.Context, document *uws1.Document) (FlowReport, error) {
	report := FlowReport{Findings: []FlowFinding{}}
	if ctx == nil || document == nil {
		return report, errors.New("flow analysis requires a document and context")
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	a := flowAnalyzer{ctx: ctx, doc: document, nodes: map[string]*flowNode{}, ambiguous: map[string]bool{}, stepsByName: map[string][]string{}, groups: map[string][]string{}, contexts: map[string]map[string]bool{}, reached: map[string]bool{}, used: map[string]bool{}, findings: map[string]FlowFinding{}}
	for i, op := range document.Operations {
		if op != nil {
			key := "op:" + op.OperationID
			n := &flowNode{key: key, path: "/operations/" + strconv.Itoa(i), op: op, outputs: op.Outputs, dependencies: append([]string(nil), op.DependsOn...)}
			n.expressions = append(n.expressions, op.When, op.ForEach, op.Wait)
			a.add(n, op.OperationID)
			a.operationExpressions(n)
		}
	}
	for i, w := range document.Workflows {
		if w != nil {
			key := "wf:" + w.WorkflowID
			n := &flowNode{key: key, path: "/workflows/" + strconv.Itoa(i), workflow: w.WorkflowID, wf: w, outputs: w.Outputs, dependencies: append([]string(nil), w.DependsOn...), expressions: []string{w.When, w.ForEach, w.Wait, w.Items, w.BatchSize}}
			a.add(n, w.WorkflowID)
			a.steps(n, w.Steps, n.path+"/steps", w.Type, 0)
			a.steps(n, w.Default, n.path+"/default", w.Type, 0)
			for j, c := range w.Cases {
				if c != nil {
					n.expressions = append(n.expressions, c.When)
					a.steps(n, c.Steps, n.path+"/cases/"+strconv.Itoa(j)+"/steps", uws1.WorkflowTypeSequence, 0)
				}
			}
		}
	}
	for _, r := range document.Results {
		if r != nil {
			workflow, step, hasStep := strings.Cut(r.From, ".")
			current := "wf:" + workflow
			if hasStep {
				current = "step:" + workflow + ":" + step
			}
			if a.nodes[current] != nil {
				a.references(r.Value, workflow, current)
			}
		}
	}
	entry := "wf:main"
	if a.nodes[entry] == nil && len(document.Workflows) == 1 && document.Workflows[0] != nil {
		entry = "wf:" + document.Workflows[0].WorkflowID
	}
	if a.nodes[entry] != nil {
		a.visitWorkflow(entry, map[string]bool{}, "", false)
	} else {
		a.find("flow.entry_indeterminate", "/workflows")
	}
	for _, trigger := range document.Triggers {
		if trigger != nil {
			for _, route := range trigger.Routes {
				if route != nil {
					for _, target := range route.To {
						ref := "wf:" + target
						if a.nodes[ref] == nil {
							ref = "step:" + strings.TrimPrefix(entry, "wf:") + ":" + target
							valid := false
							if root := a.nodes[entry]; root != nil && root.wf != nil {
								for _, s := range root.wf.Steps {
									if s != nil && s.StepID == target {
										valid = true
									}
								}
							}
							if !valid {
								a.find("flow.reference_missing", "/triggers")
								continue
							}
						}
						a.visitReference(ref, "/triggers", map[string]bool{}, "")
					}
				}
			}
		}
	}
	keys := make([]string, 0, len(a.nodes))
	for key := range a.nodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		n := a.nodes[key]
		if !a.reached[key] {
			a.find("flow.unreachable", n.path)
		}
		scopes := []string{n.workflow}
		if len(a.contexts[key]) > 0 {
			scopes = nil
			for workflow := range a.contexts[key] {
				scopes = append(scopes, workflow)
			}
			sort.Strings(scopes)
		}
		expressionScopes := scopes
		if n.wf != nil && len(n.controlContexts) > 0 {
			expressionScopes = nil
			for workflow := range n.controlContexts {
				expressionScopes = append(expressionScopes, workflow)
			}
			sort.Strings(expressionScopes)
		}
		for _, workflow := range expressionScopes {
			for _, expression := range n.expressions {
				a.references(expression, workflow, n.key)
			}
		}
		for _, workflow := range scopes {
			for _, expression := range n.outputs {
				a.references(expression, workflow, n.key)
			}
		}
	}
	for _, key := range keys {
		n := a.nodes[key]
		names := make([]string, 0, len(n.outputs))
		for name := range n.outputs {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			if !a.used[n.key+"\x00"+name] {
				a.find("flow.output_unreferenced", n.path+"/outputs/"+pointer(name))
			}
		}
	}
	ordered := make([]FlowFinding, 0, len(a.findings))
	for _, finding := range a.findings {
		ordered = append(ordered, finding)
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Path != ordered[j].Path {
			return ordered[i].Path < ordered[j].Path
		}
		return ordered[i].Code < ordered[j].Code
	})
	if len(ordered) > 128 {
		ordered = ordered[:128]
		report.Truncated = true
	}
	report.Findings = ordered
	return report, nil
}
func (a *flowAnalyzer) add(n *flowNode, name string) {
	if _, exists := a.nodes[n.key]; exists {
		a.ambiguous[n.key] = true
		a.find("flow.ambiguous_definition", n.path)
		return
	}
	a.nodes[n.key] = n
	group := ""
	if n.op != nil {
		group = n.op.ParallelGroup
	}
	if n.step != nil {
		a.stepsByName[name] = append(a.stepsByName[name], n.key)
		group = n.step.ParallelGroup
	}
	if group != "" {
		a.groups[group] = append(a.groups[group], name)
	}
}
func (a *flowAnalyzer) find(code, path string) {
	a.findings[code+"\x00"+path] = FlowFinding{Code: code, Path: path}
}
func (a *flowAnalyzer) steps(parent *flowNode, steps []*uws1.Step, path, kind string, depth int) {
	if depth > 32 {
		a.find("flow.depth", path)
		return
	}
	pending := false
	for i, s := range steps {
		if s == nil {
			continue
		}
		base := path + "/" + strconv.Itoa(i)
		key := "step:" + parent.workflow + ":" + s.StepID
		n := &flowNode{key: key, path: base, workflow: parent.workflow, step: s, outputs: s.Outputs, dependencies: append([]string(nil), s.DependsOn...), expressions: []string{s.When, s.ForEach, s.Wait, s.Items, s.BatchSize}}
		a.add(n, s.StepID)
		if kind != uws1.WorkflowTypeMerge {
			parent.children = append(parent.children, key)
		}
		if s.Pending != nil {
			a.find("flow.pending", base+"/pending")
			pending = true
		}
		if s.OperationRef != "" {
			n.children = append(n.children, "op:"+s.OperationRef)
			op := a.nodes["op:"+s.OperationRef]
			if op != nil {
				if op.op.Effect == uws1.OperationEffectUnknown || op.op.Effect == "" {
					a.find("flow.effect_unknown", base)
				}
				if pending && op.op.Effect == uws1.OperationEffectWrite && kind == uws1.WorkflowTypeSequence {
					a.find("flow.pending_before_effect", base)
				}
			}
		}
		if s.Workflow != "" {
			n.children = append(n.children, "wf:"+s.Workflow)
		}
		a.valueExpressions(s.Inputs, &n.expressions, 0)
		if s.ForEach != "" && !a.staticItems(s.ForEach) {
			a.find("flow.loop_unbounded", base+"/forEach")
		}
		if s.Type == uws1.WorkflowTypeAwait && s.Timeout == nil {
			a.find("flow.loop_unbounded", base+"/wait")
		}
		if s.Type == uws1.WorkflowTypeLoop && !a.staticItems(s.Items) {
			a.find("flow.loop_unbounded", base+"/items")
		}
		a.steps(n, s.Steps, base+"/steps", s.Type, depth+1)
		a.steps(n, s.Default, base+"/default", s.Type, depth+1)
		for j, c := range s.Cases {
			if c != nil {
				n.expressions = append(n.expressions, c.When)
				a.steps(n, c.Steps, base+"/cases/"+strconv.Itoa(j)+"/steps", uws1.WorkflowTypeSequence, depth+1)
			}
		}
	}
	if parent.wf != nil {
		if parent.wf.ForEach != "" && !a.staticItems(parent.wf.ForEach) {
			a.find("flow.loop_unbounded", parent.path+"/forEach")
		}
		if parent.wf.Type == uws1.WorkflowTypeAwait && parent.wf.Timeout == nil {
			a.find("flow.loop_unbounded", parent.path+"/wait")
		}
		if parent.wf.Type == uws1.WorkflowTypeLoop && !a.staticItems(parent.wf.Items) {
			a.find("flow.loop_unbounded", parent.path+"/items")
		}
	}
}
func (a *flowAnalyzer) staticItems(text string) bool {
	if !strings.HasPrefix(text, "$variables.") {
		return false
	}
	name := strings.TrimPrefix(text, "$variables.")
	if strings.Contains(name, ".") {
		return false
	}
	value, ok := a.doc.Variables[name]
	if !ok && a.doc.Components != nil {
		value, ok = a.doc.Components.Variables[name]
	}
	if !ok {
		return false
	}
	_, ok = value.([]any)
	return ok
}
func (a *flowAnalyzer) operationExpressions(n *flowNode) {
	op := n.op
	if op.ForEach != "" && !a.staticItems(op.ForEach) {
		a.find("flow.loop_unbounded", n.path+"/forEach")
	}
	browser := false
	for _, source := range a.doc.SourceDescriptions {
		if source != nil && source.Name == op.SourceDescription && source.Type == "browser-profile" {
			browser = true
		}
	}
	if op.HasSourceBinding() && !browser {
		for _, key := range []string{"path", "query", "header", "cookie", "body"} {
			a.valueExpressions(op.Request[key], &n.expressions, 0)
		}
	}
	criteria := func(list []*uws1.Criterion) {
		for _, c := range list {
			if c != nil {
				n.expressions = append(n.expressions, c.Context)
				if c.Type == "" || c.Type == uws1.CriterionSimple {
					n.expressions = append(n.expressions, c.Condition)
				}
			}
		}
	}
	criteria(op.SuccessCriteria)
	for _, action := range op.OnSuccess {
		if action != nil {
			criteria(action.Criteria)
			if action.Type == "goto" {
				if action.WorkflowID != "" {
					n.transfers = append(n.transfers, "wf:"+action.WorkflowID)
				}
				if action.StepID != "" {
					n.transfers = append(n.transfers, "step:"+action.StepID)
				}
			}
		}
	}
	for _, action := range op.OnFailure {
		if action != nil {
			criteria(action.Criteria)
			if action.Type == "goto" {
				if action.WorkflowID != "" {
					n.transfers = append(n.transfers, "wf:"+action.WorkflowID)
				}
				if action.StepID != "" {
					n.transfers = append(n.transfers, "step:"+action.StepID)
				}
			}
			if action.Type == "retry" && action.RetryLimit <= 0 {
				a.find("flow.loop_unbounded", n.path+"/onFailure")
			}
		}
	}
}
func (a *flowAnalyzer) valueExpressions(value any, out *[]string, depth int) {
	if depth > 32 {
		return
	}
	switch v := value.(type) {
	case string:
		if strings.HasPrefix(v, "$") {
			*out = append(*out, v)
		}
	case map[string]any:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			a.valueExpressions(v[key], out, depth+1)
		}
	case []any:
		for _, x := range v {
			a.valueExpressions(x, out, depth+1)
		}
	}
}
func (a *flowAnalyzer) visitReference(ref, path string, stack map[string]bool, workflow string) {
	if a.ambiguous[ref] {
		a.find("flow.reference_ambiguous", path)
		return
	}
	if n := a.nodes[ref]; n != nil {
		if n.wf != nil {
			a.visitWorkflow(ref, stack, workflow, false)
			return
		}
		a.visit(ref, stack, workflow)
		return
	}
	// Generic dependsOn retains the published group barrier and unqualified
	// step/workflow/operation precedence. Cross-kind names are not ambiguity.
	if members := a.groups[ref]; len(members) > 0 {
		for _, member := range members {
			a.visitReference(member, path, stack, workflow)
		}
		return
	}
	if steps := a.stepsByName[ref]; len(steps) > 0 {
		if len(steps) != 1 || a.ambiguous[steps[0]] {
			a.find("flow.reference_ambiguous", path)
			return
		}
		a.visit(steps[0], stack, workflow)
		return
	}
	if a.nodes["wf:"+ref] != nil {
		a.visitReference("wf:"+ref, path, stack, workflow)
		return
	}
	if a.nodes["op:"+ref] != nil {
		a.visitReference("op:"+ref, path, stack, workflow)
		return
	}
	a.find("flow.reference_missing", path)
}

func (a *flowAnalyzer) visitTransfer(ref, path string, stack map[string]bool) {
	if id, step := strings.CutPrefix(ref, "step:"); step {
		targets := a.stepsByName[id]
		if len(targets) == 0 {
			a.find("flow.reference_missing", path)
			return
		}
		if len(targets) != 1 || a.ambiguous[targets[0]] {
			a.find("flow.reference_ambiguous", path)
			return
		}
		ref = targets[0]
	}
	if a.nodes[ref] == nil {
		a.find("flow.reference_missing", path)
		return
	}
	// Terminal goto uses the global exact target in the root invocation,
	// independently of the caller frame retained by generic dependencies.
	if a.ambiguous[ref] {
		a.find("flow.reference_ambiguous", path)
		return
	}
	if a.nodes[ref].wf != nil {
		a.visitWorkflow(ref, stack, "", false)
	} else {
		a.visit(ref, stack, "")
	}
}
func (a *flowAnalyzer) visitWorkflow(key string, stack map[string]bool, incoming string, call bool) {
	n := a.nodes[key]
	if n.controlContexts == nil {
		n.controlContexts = map[string]bool{}
	}
	// Native controls retain executeOnce's incoming record snapshot. Child
	// runnables and final outputs refresh records in the body's invocation frame.
	n.controlContexts[incoming] = true
	workflow := incoming
	if call || incoming != "" {
		workflow = n.workflow
	}
	a.visit(key, stack, workflow)
}
func (a *flowAnalyzer) visit(key string, stack map[string]bool, workflow string) {
	n := a.nodes[key]
	if n == nil {
		a.find("flow.reference_missing", "/")
		return
	}
	if stack[key] {
		a.find("flow.cycle", n.path)
		return
	}
	if a.contexts[key] == nil {
		a.contexts[key] = map[string]bool{}
	}
	if a.contexts[key][workflow] {
		return
	}
	a.contexts[key][workflow] = true
	a.reached[key] = true
	stack[key] = true
	for _, ref := range n.dependencies {
		a.visitReference(ref, n.path+"/dependsOn", stack, workflow)
	}
	for _, child := range n.children {
		if target := a.nodes[child]; target != nil && target.wf != nil {
			a.visitWorkflow(child, stack, workflow, true)
		} else {
			a.visit(child, stack, workflow)
		}
	}
	for _, target := range n.transfers {
		a.visitTransfer(target, n.path+"/goto", stack)
	}
	delete(stack, key)
}
func (a *flowAnalyzer) references(text, workflow, current string) {
	if text == "" {
		return
	}
	field := expressions.Value
	if !strings.HasPrefix(text, "$") {
		field = expressions.Wait
	}
	parsed, err := expressions.Parse(text, expressions.Context{Version: a.doc.UWS, Field: field, InLoop: true})
	if err != nil {
		return
	}
	for _, source := range []string{parsed.Source(), parsed.OperandSource()} {
		if strings.HasPrefix(source, "$steps.") {
			parts := strings.Split(strings.TrimPrefix(source, "$steps."), ".")
			if len(parts) < 3 || parts[1] != "outputs" {
				continue
			}
			key := "step:" + workflow + ":" + parts[0]
			if a.nodes[key] == nil || len(a.contexts[key]) > 0 && !a.contexts[key][workflow] {
				// A cross-declaration dependency may execute in this exact caller
				// frame. Only recorded invocation evidence can supply that owner.
				matches := []string{}
				for _, candidate := range a.stepsByName[parts[0]] {
					if a.contexts[candidate][workflow] {
						matches = append(matches, candidate)
					}
				}
				if len(matches) != 1 {
					continue
				}
				key = matches[0]
			}
			if a.ambiguous[key] {
				continue
			}
			a.used[key+"\x00"+parts[2]] = true
		} else if strings.HasPrefix(source, "$outputs.") {
			name := strings.Split(strings.TrimPrefix(source, "$outputs."), ".")[0]
			a.used[current+"\x00"+name] = true
		}
	}
}
