package contenttrust

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/OpenUdon/uws/uws1"
)

type valueState struct {
	provenance uws1.ContentTrustLevel
	capability ValueCapability
	from       string
}

type operationResolution struct {
	contract OperationContract
	resolved bool
}

type orderPosition struct {
	container  string
	index      int
	sequential bool
}

type stepMetadata struct {
	workflow  string
	positions []orderPosition
	contexts  map[string]bool
	dependsOn []string
}

type evalEnvironment struct {
	workflow       string
	currentStep    string
	inputs         map[string]valueState
	outputs        map[string]valueState
	response       valueState
	item           valueState
	atWorkflowExit bool
}

type analyzer struct {
	ctx context.Context
	doc *uws1.Document

	operations       map[string]*uws1.Operation
	operationPaths   map[string]string
	workflows        map[string]*uws1.Workflow
	workflowPaths    map[string]string
	steps            map[string]*uws1.Step
	stepPaths        map[string]string
	stepMeta         map[string]stepMetadata
	parallelGroups   map[string][]string
	operationResolve map[string]operationResolution
	referencedOps    map[string]bool

	variables      map[string]valueState
	workflowInputs map[string]map[string]valueState
	workflowOutput map[string]map[string]valueState
	stepOutputs    map[string]map[string]valueState

	trigger valueState
	changed bool
	emit    bool

	edges       []Edge
	edgeKeys    map[string]bool
	findings    []Finding
	findingKeys map[string]bool
}

// Analyze validates the document and performs deterministic advisory content-
// trust analysis. Resolver diagnostics are findings; malformed UWS documents
// and context cancellation are returned as errors.
func Analyze(ctx context.Context, doc *uws1.Document, resolvers ...Resolver) (*Report, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, fmt.Errorf("contenttrust: document is required")
	}
	if err := doc.Validate(); err != nil {
		return nil, fmt.Errorf("contenttrust: validate document: %w", err)
	}
	a := newAnalyzer(ctx, doc)
	a.resolveOperations(resolvers)
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Propagation is monotonic. Revisit the finite graph until no provenance or
	// capability state changes, then make one reporting pass over the fixed point.
	limit := 4 + len(a.steps)*2 + len(a.workflows)*2
	for i := 0; i < limit; i++ {
		a.changed = false
		if err := a.runPass(); err != nil {
			return nil, err
		}
		if !a.changed {
			break
		}
	}
	a.emit = true
	if err := a.runPass(); err != nil {
		return nil, err
	}

	sort.Slice(a.edges, func(i, j int) bool {
		left, right := a.edges[i], a.edges[j]
		if left.To != right.To {
			return left.To < right.To
		}
		if left.From != right.From {
			return left.From < right.From
		}
		if left.Provenance != right.Provenance {
			return left.Provenance < right.Provenance
		}
		return left.Capability < right.Capability
	})
	sort.Slice(a.findings, func(i, j int) bool {
		left, right := a.findings[i], a.findings[j]
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		return left.Severity < right.Severity
	})
	return &Report{Edges: a.edges, Findings: a.findings}, nil
}

func newAnalyzer(ctx context.Context, doc *uws1.Document) *analyzer {
	a := &analyzer{
		ctx:              ctx,
		doc:              doc,
		operations:       make(map[string]*uws1.Operation),
		operationPaths:   make(map[string]string),
		workflows:        make(map[string]*uws1.Workflow),
		workflowPaths:    make(map[string]string),
		steps:            make(map[string]*uws1.Step),
		stepPaths:        make(map[string]string),
		stepMeta:         make(map[string]stepMetadata),
		parallelGroups:   make(map[string][]string),
		operationResolve: make(map[string]operationResolution),
		referencedOps:    make(map[string]bool),
		variables:        make(map[string]valueState),
		workflowInputs:   make(map[string]map[string]valueState),
		workflowOutput:   make(map[string]map[string]valueState),
		stepOutputs:      make(map[string]map[string]valueState),
		edgeKeys:         make(map[string]bool),
		findingKeys:      make(map[string]bool),
	}
	for i, operation := range doc.Operations {
		if operation == nil {
			continue
		}
		a.operations[operation.OperationID] = operation
		a.operationPaths[operation.OperationID] = fmt.Sprintf("operations[%d]", i)
	}
	for i, workflow := range doc.Workflows {
		if workflow == nil {
			continue
		}
		a.workflows[workflow.WorkflowID] = workflow
		a.workflowPaths[workflow.WorkflowID] = fmt.Sprintf("workflows[%d]", i)
		rootContexts := make(map[string]bool)
		if workflow.When != "" || workflow.ForEach != "" {
			rootContexts[workflow.WorkflowID+":conditional"] = true
		}
		if workflow.Type == uws1.WorkflowTypeLoop {
			rootContexts[workflow.WorkflowID+":loop"] = true
		}
		a.collectStepList(workflow.WorkflowID, workflow.Steps, a.workflowPaths[workflow.WorkflowID]+".steps", workflow.WorkflowID+":steps", workflow.Type != uws1.WorkflowTypeParallel, nil, rootContexts)
		a.collectCaseLists(workflow.WorkflowID, workflow.Cases, a.workflowPaths[workflow.WorkflowID]+".cases", workflow.WorkflowID+":cases", nil, rootContexts)
		defaultContexts := cloneContexts(rootContexts)
		defaultContexts[workflow.WorkflowID+":default"] = true
		a.collectStepList(workflow.WorkflowID, workflow.Default, a.workflowPaths[workflow.WorkflowID]+".default", workflow.WorkflowID+":default", true, nil, defaultContexts)
	}
	for name, raw := range doc.Variables {
		a.variables[name] = literalState(raw, "variables."+name)
	}
	if doc.Components != nil {
		for name, raw := range doc.Components.Variables {
			if _, exists := a.variables[name]; !exists {
				a.variables[name] = literalState(raw, "components.variables."+name)
			}
		}
	}
	a.initializeWorkflowInputs()
	a.trigger = a.documentTriggerState()
	return a
}

func (a *analyzer) collectStepList(workflow string, steps []*uws1.Step, path, container string, sequential bool, positions []orderPosition, contexts map[string]bool) {
	for i, step := range steps {
		if step == nil {
			continue
		}
		stepPath := fmt.Sprintf("%s[%d]", path, i)
		stepPositions := appendCopy(positions, orderPosition{container: container, index: i, sequential: sequential})
		stepContexts := cloneContexts(contexts)
		if step.When != "" || step.ForEach != "" {
			stepContexts[step.StepID+":conditional"] = true
		}
		a.steps[step.StepID] = step
		if step.OperationRef != "" {
			a.referencedOps[step.OperationRef] = true
		}
		a.stepPaths[step.StepID] = stepPath
		a.stepMeta[step.StepID] = stepMetadata{workflow: workflow, positions: stepPositions, contexts: stepContexts, dependsOn: append([]string(nil), step.DependsOn...)}
		if step.ParallelGroup != "" {
			a.parallelGroups[step.ParallelGroup] = append(a.parallelGroups[step.ParallelGroup], step.StepID)
		}

		childContexts := cloneContexts(stepContexts)
		switch step.Type {
		case uws1.WorkflowTypeSwitch:
			a.collectCaseLists(workflow, step.Cases, stepPath+".cases", step.StepID+":cases", stepPositions, childContexts)
			defaultContexts := cloneContexts(childContexts)
			defaultContexts[step.StepID+":default"] = true
			a.collectStepList(workflow, step.Default, stepPath+".default", step.StepID+":default", true, stepPositions, defaultContexts)
		case uws1.WorkflowTypeLoop:
			childContexts[step.StepID+":loop"] = true
			a.collectStepList(workflow, step.Steps, stepPath+".steps", step.StepID+":steps", true, stepPositions, childContexts)
		default:
			childSequential := step.Type != uws1.WorkflowTypeParallel
			a.collectStepList(workflow, step.Steps, stepPath+".steps", step.StepID+":steps", childSequential, stepPositions, childContexts)
			a.collectCaseLists(workflow, step.Cases, stepPath+".cases", step.StepID+":cases", stepPositions, childContexts)
			a.collectStepList(workflow, step.Default, stepPath+".default", step.StepID+":default", true, stepPositions, childContexts)
		}
	}
}

func (a *analyzer) collectCaseLists(workflow string, cases []*uws1.Case, path, container string, positions []orderPosition, contexts map[string]bool) {
	for i, c := range cases {
		if c == nil {
			continue
		}
		caseContexts := cloneContexts(contexts)
		caseContexts[fmt.Sprintf("%s:%d", container, i)] = true
		a.collectStepList(workflow, c.Steps, fmt.Sprintf("%s[%d].steps", path, i), fmt.Sprintf("%s:%d", container, i), true, positions, caseContexts)
	}
}

func appendCopy[T any](values []T, value T) []T {
	out := append([]T(nil), values...)
	return append(out, value)
}

func cloneContexts(values map[string]bool) map[string]bool {
	out := make(map[string]bool, len(values)+1)
	for key, value := range values {
		out[key] = value
	}
	return out
}

func (a *analyzer) initializeWorkflowInputs() {
	for id, workflow := range a.workflows {
		inputs := make(map[string]valueState)
		if workflow.Inputs != nil {
			for name, schema := range workflow.Inputs.Properties {
				inputs[name] = valueState{
					provenance: a.workflowInputTrust(id, name),
					capability: capabilityForSchema(schema),
					from:       a.workflowPaths[id] + ".inputs.properties." + name,
				}
			}
		}
		a.workflowInputs[id] = inputs
	}
}

func (a *analyzer) documentTriggerState() valueState {
	state := valueState{provenance: uws1.ContentTrustUntrusted, capability: CapabilityUnknown, from: "$trigger"}
	if len(a.doc.Triggers) == 0 {
		return state
	}
	first := true
	for _, trigger := range a.doc.Triggers {
		if trigger == nil {
			continue
		}
		level := uws1.ContentTrustUntrusted
		if a.doc.ContentTrust != nil {
			if declared, ok := a.doc.ContentTrust.Triggers[trigger.TriggerID]; ok {
				level = declared
			}
		}
		candidate := valueState{provenance: level, capability: CapabilityUnknown, from: "$trigger"}
		if first {
			state = candidate
			first = false
		} else {
			state = joinState(state, candidate)
		}
	}
	return state
}

func (a *analyzer) resolveOperations(resolvers []Resolver) {
	for _, id := range sortedOperationIDs(a.operations) {
		operation := a.operations[id]
		claims := make([]OperationContract, 0, 1)
		for _, resolver := range resolvers {
			if err := a.ctx.Err(); err != nil {
				return
			}
			if nilResolver(resolver) {
				a.addFinding(CodeResolverFailure, a.operationPaths[id])
				continue
			}
			owned, contract, err := resolver.ResolveOperation(a.ctx, a.doc, operation)
			if err != nil {
				if a.ctx.Err() != nil {
					return
				}
				a.addFinding(CodeResolverFailure, a.operationPaths[id])
				continue
			}
			if owned {
				normalized, ok := validateContract(operation, contract)
				if !ok {
					a.addFinding(CodeResolverFailure, a.operationPaths[id])
					continue
				}
				claims = append(claims, normalized)
			}
		}
		switch len(claims) {
		case 0:
			a.operationResolve[id] = operationResolution{}
		case 1:
			a.operationResolve[id] = operationResolution{contract: claims[0], resolved: true}
		default:
			a.addFinding(CodeResolverConflict, a.operationPaths[id])
			a.operationResolve[id] = operationResolution{}
		}
	}
}

func nilResolver(resolver Resolver) bool {
	if resolver == nil {
		return true
	}
	value := reflect.ValueOf(resolver)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func validateContract(operation *uws1.Operation, contract OperationContract) (OperationContract, bool) {
	normalized := OperationContract{
		Inputs:                  make([]InputChannel, 0, len(contract.Inputs)),
		Outputs:                 make(map[string]ValueContract, len(contract.Outputs)),
		DefaultTrust:            contract.DefaultTrust,
		InheritsInputProvenance: contract.InheritsInputProvenance,
	}
	if normalized.DefaultTrust == "" {
		normalized.DefaultTrust = uws1.ContentTrustUnknown
	} else if !validTrust(normalized.DefaultTrust) {
		return OperationContract{}, false
	}
	for _, channel := range contract.Inputs {
		switch channel.Kind {
		case ChannelData, ChannelInstruction, ChannelAuthority:
		default:
			return OperationContract{}, false
		}
		if !validRelativeJSONPointer(channel.Path) {
			return OperationContract{}, false
		}
		if _, ok := operationValueAtPointer(operation, channel.Path); !ok {
			return OperationContract{}, false
		}
		normalizedChannel := InputChannel{
			Path:       channel.Path,
			Kind:       channel.Kind,
			References: append([]Reference(nil), channel.References...),
		}
		for _, reference := range normalizedChannel.References {
			if !validRelativeJSONPointer(reference.Path) {
				return OperationContract{}, false
			}
		}
		normalized.Inputs = append(normalized.Inputs, normalizedChannel)
	}
	for name, output := range contract.Outputs {
		if operation.Outputs == nil {
			return OperationContract{}, false
		}
		if _, ok := operation.Outputs[name]; !ok {
			return OperationContract{}, false
		}
		if output.Capability == "" {
			output.Capability = CapabilityUnknown
		} else if !validCapability(output.Capability) {
			return OperationContract{}, false
		}
		normalized.Outputs[name] = output
	}
	return normalized, true
}

func validRelativeJSONPointer(pointer string) bool {
	if pointer == "" || pointer == "/" {
		return true
	}
	if !strings.HasPrefix(pointer, "/") {
		return false
	}
	for i := 0; i < len(pointer); i++ {
		if pointer[i] != '~' {
			continue
		}
		if i+1 >= len(pointer) || (pointer[i+1] != '0' && pointer[i+1] != '1') {
			return false
		}
		i++
	}
	return true
}

func validCapability(capability ValueCapability) bool {
	switch capability {
	case CapabilityFreeText, CapabilityConstrainedScalar, CapabilityComposite, CapabilityUnknown:
		return true
	default:
		return false
	}
}

func (a *analyzer) runPass() error {
	for _, id := range sortedOperationIDs(a.operations) {
		if err := a.ctx.Err(); err != nil {
			return err
		}
		// Analyze operation declarations once with unknown external inputs. Calls
		// from steps are analyzed below with their actual bound provenance.
		if !a.referencedOps[id] {
			a.analyzeOperation(id, evalEnvironment{inputs: map[string]valueState{}}, "")
		}
	}
	for _, id := range sortedWorkflowIDs(a.workflows) {
		if err := a.ctx.Err(); err != nil {
			return err
		}
		a.analyzeWorkflow(id)
	}
	a.analyzeStructuralResults()
	return nil
}

func (a *analyzer) analyzeStructuralResults() {
	for i, result := range a.doc.Results {
		if result == nil || result.Value == "" {
			continue
		}
		workflowID := result.From
		if index := strings.IndexByte(workflowID, '.'); index >= 0 {
			workflowID = workflowID[:index]
		}
		env := evalEnvironment{
			workflow:       workflowID,
			inputs:         a.workflowInputs[workflowID],
			atWorkflowExit: true,
		}
		a.evalString(result.Value, fmt.Sprintf("results[%d].value", i), env)
	}
}

func (a *analyzer) analyzeWorkflow(id string) map[string]valueState {
	workflow := a.workflows[id]
	if workflow == nil {
		return nil
	}
	path := a.workflowPaths[id]
	env := evalEnvironment{workflow: id, inputs: a.workflowInputs[id]}
	a.analyzeControl(workflow.When, path+".when", env)
	a.analyzeControl(workflow.ForEach, path+".forEach", env)
	a.analyzeControl(workflow.Wait, path+".wait", env)
	a.analyzeControl(workflow.Items, path+".items", env)
	a.analyzeControl(workflow.BatchSize, path+".batchSize", env)
	if workflow.Idempotency != nil {
		a.analyzeControl(workflow.Idempotency.Key, path+".idempotency.key", env)
	}
	a.analyzeSteps(workflow.Steps, path+".steps", env)
	a.analyzeCases(workflow.Cases, path+".cases", env)
	a.analyzeSteps(workflow.Default, path+".default", env)

	outputs := make(map[string]valueState, len(workflow.Outputs))
	env.atWorkflowExit = true
	for _, name := range sortedStringKeys(workflow.Outputs) {
		state := a.evalString(workflow.Outputs[name], path+".outputs."+name, env)
		state.from = path + ".outputs." + name
		outputs[name] = state
		a.setNestedState(a.workflowOutput, id, name, state)
	}
	return outputs
}

func (a *analyzer) analyzeSteps(steps []*uws1.Step, path string, env evalEnvironment) {
	for i, step := range steps {
		if step == nil {
			continue
		}
		a.analyzeStep(step, fmt.Sprintf("%s[%d]", path, i), env)
	}
}

func (a *analyzer) analyzeCases(cases []*uws1.Case, path string, env evalEnvironment) {
	for i, c := range cases {
		if c == nil {
			continue
		}
		casePath := fmt.Sprintf("%s[%d]", path, i)
		a.analyzeControl(c.When, casePath+".when", env)
		a.scanValue(c.Body, casePath+".body", env)
		a.analyzeSteps(c.Steps, casePath+".steps", env)
	}
}

func (a *analyzer) analyzeStep(step *uws1.Step, path string, parent evalEnvironment) {
	env := parent
	env.currentStep = step.StepID
	a.analyzeControl(step.When, path+".when", env)
	item := a.analyzeControl(step.ForEach, path+".forEach", env)
	if item.from != "" {
		env.item = item
	}
	a.analyzeControl(step.Wait, path+".wait", env)
	loopItem := a.analyzeControl(step.Items, path+".items", env)
	if loopItem.from != "" {
		env.item = loopItem
	}
	a.analyzeControl(step.BatchSize, path+".batchSize", env)
	a.scanValue(step.Body, path+".body", env)

	boundInputs := make(map[string]valueState, len(step.Inputs))
	for _, name := range sortedAnyKeys(step.Inputs) {
		boundInputs[name] = a.evalValue(step.Inputs[name], path+".inputs."+name, env)
	}

	localOutputs := map[string]valueState{}
	response := unknownState(a.operationPaths[step.OperationRef] + ".response")
	if step.OperationRef != "" {
		callEnv := env
		callEnv.inputs = boundInputs
		localOutputs, response = a.analyzeOperation(step.OperationRef, callEnv, step.StepID)
	} else if step.Workflow != "" {
		for name, state := range boundInputs {
			a.setNestedState(a.workflowInputs, step.Workflow, name, state)
		}
		localOutputs = cloneStateMap(a.workflowOutput[step.Workflow])
	} else {
		a.analyzeSteps(step.Steps, path+".steps", env)
		a.analyzeCases(step.Cases, path+".cases", env)
		a.analyzeSteps(step.Default, path+".default", env)
	}

	outputEnv := env
	outputEnv.inputs = boundInputs
	outputEnv.outputs = localOutputs
	outputEnv.response = response
	for _, name := range sortedStringKeys(step.Outputs) {
		state := a.evalString(step.Outputs[name], path+".outputs."+name, outputEnv)
		state.from = path + ".outputs." + name
		a.setNestedState(a.stepOutputs, step.StepID, name, state)
	}
}

func (a *analyzer) analyzeOperation(id string, env evalEnvironment, callerStep string) (map[string]valueState, valueState) {
	operation := a.operations[id]
	if operation == nil {
		return nil, unknownState("operations.unknown.response")
	}
	path := a.operationPaths[id]
	env.currentStep = callerStep
	a.analyzeControl(operation.When, path+".when", env)
	a.analyzeControl(operation.ForEach, path+".forEach", env)
	a.analyzeControl(operation.Wait, path+".wait", env)

	resolution := a.operationResolve[id]
	inputState := trustedState(path + ".request")
	if resolution.resolved {
		for _, channel := range resolution.contract.Inputs {
			channelPath := pointerDocumentPath(path, channel.Path)
			var state valueState
			if len(channel.References) > 0 {
				state = trustedState(channelPath)
				for _, supplied := range channel.References {
					target := channelPath + pointerSuffix(supplied.Path)
					parsed, ok := parseExpression(supplied.Expression)
					if !ok {
						a.addFinding(CodeOpaqueExpression, target)
						state = joinState(state, unknownState(target))
						continue
					}
					state = joinState(state, a.evalParsed(parsed, target, env))
				}
			} else if raw, ok := operationValueAtPointer(operation, channel.Path); ok {
				state = a.evalValue(raw, channelPath, env)
			} else {
				a.addFinding(CodeResolverFailure, path)
				state = unknownState(channelPath)
			}
			inputState = joinState(inputState, state)
			a.analyzeChannel(channel.Kind, state, channelPath)
		}
	} else {
		// Core request flow remains statically visible, but without a resolver it
		// is data rather than an instruction or authority sink.
		inputState = a.scanValue(operation.Request, path+".request", env)
	}

	response := valueState{
		provenance: a.operationOutputTrust(operation, "", resolution),
		capability: CapabilityUnknown,
		from:       path + ".response",
	}
	if resolution.contract.InheritsInputProvenance {
		response.provenance = joinTrust(response.provenance, inputState.provenance)
	}

	outputs := make(map[string]valueState, len(operation.Outputs))
	outputEnv := env
	outputEnv.response = response
	for _, name := range sortedStringKeys(operation.Outputs) {
		state := a.evalString(operation.Outputs[name], path+".outputs."+name, outputEnv)
		state.provenance = a.operationOutputTrust(operation, name, resolution)
		if output, ok := resolution.contract.Outputs[name]; ok {
			state.capability = output.Capability
			inherits := resolution.contract.InheritsInputProvenance
			if output.InheritsInputProvenance != nil {
				inherits = *output.InheritsInputProvenance
			}
			if inherits {
				state.provenance = joinTrust(state.provenance, inputState.provenance)
			}
		} else if resolution.contract.InheritsInputProvenance {
			state.provenance = joinTrust(state.provenance, inputState.provenance)
		}
		state.from = path + ".outputs." + name
		outputs[name] = state
	}

	criteriaEnv := env
	criteriaEnv.outputs = outputs
	criteriaEnv.response = response
	for i, criterion := range operation.SuccessCriteria {
		a.analyzeCriterion(criterion, fmt.Sprintf("%s.successCriteria[%d]", path, i), criteriaEnv)
	}
	for i, action := range operation.OnFailure {
		if action == nil {
			continue
		}
		for j, criterion := range action.Criteria {
			a.analyzeCriterion(criterion, fmt.Sprintf("%s.onFailure[%d].criteria[%d]", path, i, j), criteriaEnv)
		}
	}
	for i, action := range operation.OnSuccess {
		if action == nil {
			continue
		}
		for j, criterion := range action.Criteria {
			a.analyzeCriterion(criterion, fmt.Sprintf("%s.onSuccess[%d].criteria[%d]", path, i, j), criteriaEnv)
		}
	}
	return outputs, response
}

func (a *analyzer) analyzeCriterion(criterion *uws1.Criterion, path string, env evalEnvironment) {
	if criterion == nil {
		return
	}
	a.analyzeControl(criterion.Condition, path+".condition", env)
	a.analyzeControl(criterion.Context, path+".context", env)
}

func (a *analyzer) analyzeControl(raw, path string, env evalEnvironment) valueState {
	if raw == "" {
		return valueState{}
	}
	state := a.evalString(raw, path, env)
	switch state.provenance {
	case uws1.ContentTrustUntrusted:
		a.addFinding(CodeUntrustedControl, path)
	case uws1.ContentTrustUnknown:
		a.addFinding(CodeUnknownProvenance, path)
	}
	return state
}

func (a *analyzer) analyzeChannel(kind ChannelKind, state valueState, path string) {
	switch {
	case state.provenance == uws1.ContentTrustUnknown:
		a.addFinding(CodeUnknownProvenance, path)
	case kind == ChannelInstruction && state.provenance == uws1.ContentTrustUntrusted && state.capability != CapabilityConstrainedScalar:
		a.addFinding(CodeUntrustedInstruction, path)
	case kind == ChannelAuthority && state.provenance == uws1.ContentTrustUntrusted:
		a.addFinding(CodeUntrustedAuthority, path)
	}
}

func (a *analyzer) evalValue(raw any, path string, env evalEnvironment) valueState {
	switch value := raw.(type) {
	case string:
		return a.evalString(value, path, env)
	case map[string]any:
		state := trustedState(path)
		state.capability = CapabilityComposite
		for _, key := range sortedAnyKeys(value) {
			state = joinState(state, a.evalValue(value[key], path+"."+key, env))
		}
		state.capability = CapabilityComposite
		return state
	case []any:
		state := trustedState(path)
		state.capability = CapabilityComposite
		for i, item := range value {
			state = joinState(state, a.evalValue(item, fmt.Sprintf("%s[%d]", path, i), env))
		}
		state.capability = CapabilityComposite
		return state
	default:
		if value == nil {
			return literalState(nil, path)
		}
		reflected := reflect.ValueOf(value)
		if reflected.Kind() == reflect.String {
			return a.evalString(reflected.String(), path, env)
		}
		switch reflected.Kind() {
		case reflect.Map, reflect.Array, reflect.Slice, reflect.Struct, reflect.Pointer, reflect.Interface:
			normalized, ok := normalizeWireValue(value)
			if !ok {
				a.addFinding(CodeUnknownProvenance, path)
				return unknownState(path)
			}
			return a.evalValue(normalized, path, env)
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return literalState(value, path)
		default:
			a.addFinding(CodeUnknownProvenance, path)
			return unknownState(path)
		}
	}
}

func normalizeWireValue(raw any) (any, bool) {
	data, err := json.Marshal(raw)
	if err != nil {
		return nil, false
	}
	var normalized any
	if err := json.Unmarshal(data, &normalized); err != nil {
		return nil, false
	}
	return normalized, true
}

func (a *analyzer) scanValue(raw any, path string, env evalEnvironment) valueState {
	if raw == nil {
		return trustedState(path)
	}
	return a.evalValue(raw, path, env)
}

func (a *analyzer) evalString(raw, path string, env evalEnvironment) valueState {
	parsed, ok := parseExpression(raw)
	if ok {
		return a.evalParsed(parsed, path, env)
	}
	if strings.Contains(raw, "$") {
		a.addFinding(CodeOpaqueExpression, path)
		return unknownState(path)
	}
	return valueState{provenance: uws1.ContentTrustTrusted, capability: CapabilityFreeText, from: path}
}

func (a *analyzer) evalParsed(parsed parsedExpression, path string, env evalEnvironment) valueState {
	var state valueState
	for _, reference := range parsed.references {
		resolved, ok := a.resolveReference(reference, env, path)
		if !ok {
			a.addFinding(CodeUnresolvedReference, path)
			resolved = unknownState(referenceLabel(reference, env))
		}
		state = joinState(state, resolved)
		a.addEdge(resolved, path)
	}
	if state.from == "" {
		state = trustedState(path)
	}
	if parsed.condition {
		state.capability = CapabilityConstrainedScalar
	}
	return state
}

func (a *analyzer) resolveReference(ref expressionReference, env evalEnvironment, path string) (valueState, bool) {
	switch ref.root {
	case "variables":
		state, ok := a.variables[ref.name]
		return state, ok
	case "trigger":
		return a.trigger, true
	case "inputs":
		if ref.name == "" {
			if len(env.inputs) == 0 {
				return unknownState(a.workflowPaths[env.workflow] + ".inputs"), false
			}
			return joinMapStates(env.inputs, a.workflowPaths[env.workflow]+".inputs"), true
		}
		state, ok := env.inputs[ref.name]
		return state, ok
	case "outputs":
		state, ok := env.outputs[ref.name]
		return state, ok
	case "steps":
		outputs, exists := a.stepOutputs[ref.id]
		state, ok := outputs[ref.name]
		if !exists || !ok {
			return unknownState("$steps." + ref.id + ".outputs." + ref.name), false
		}
		if env.currentStep != "" {
			if !a.dominates(ref.id, env.currentStep) {
				a.addFinding(CodeNonDominatingRef, path)
			}
		} else if !a.dominatesWorkflowExit(ref.id, env) {
			a.addFinding(CodeNonDominatingRef, path)
		}
		return state, true
	case "response":
		state := env.response
		if state.from == "" {
			return unknownState("$response"), false
		}
		if ref.name == "statusCode" {
			state.capability = CapabilityConstrainedScalar
		} else if ref.name == "headers" {
			state.capability = CapabilityFreeText
		}
		return state, true
	case "item":
		if env.item.from == "" {
			return unknownState("$item"), false
		}
		return env.item, true
	case "index":
		return valueState{provenance: uws1.ContentTrustTrusted, capability: CapabilityConstrainedScalar, from: "$index"}, true
	default:
		return unknownState(referenceLabel(ref, env)), false
	}
}

func (a *analyzer) dominatesWorkflowExit(source string, env evalEnvironment) bool {
	if !env.atWorkflowExit || env.workflow == "" {
		return false
	}
	meta, ok := a.stepMeta[source]
	return ok && meta.workflow == env.workflow && len(meta.contexts) == 0
}

func referenceLabel(ref expressionReference, env evalEnvironment) string {
	switch ref.root {
	case "variables":
		return "$variables." + ref.name
	case "inputs":
		if ref.name == "" {
			return "$inputs"
		}
		return "$inputs." + ref.name
	case "outputs":
		return "$outputs." + ref.name
	case "steps":
		return "$steps." + ref.id + ".outputs." + ref.name
	case "response":
		return "$response." + ref.name
	case "trigger":
		return "$trigger"
	case "item":
		return "$item"
	case "index":
		return "$index"
	default:
		return "$unknown"
	}
}

func (a *analyzer) dominates(source, target string) bool {
	if source == target {
		return false
	}
	sourceMeta, sourceOK := a.stepMeta[source]
	targetMeta, targetOK := a.stepMeta[target]
	if !sourceOK || !targetOK || sourceMeta.workflow != targetMeta.workflow {
		return false
	}
	for context := range sourceMeta.contexts {
		if !targetMeta.contexts[context] {
			return false
		}
	}
	if a.explicitlyDependsOn(target, source, make(map[string]bool)) {
		return true
	}
	limit := len(sourceMeta.positions)
	if len(targetMeta.positions) < limit {
		limit = len(targetMeta.positions)
	}
	for i := 0; i < limit; i++ {
		left, right := sourceMeta.positions[i], targetMeta.positions[i]
		if left.container != right.container {
			return false
		}
		if left.index != right.index {
			return left.sequential && right.sequential && left.index < right.index
		}
	}
	return false
}

func (a *analyzer) explicitlyDependsOn(target, source string, seen map[string]bool) bool {
	if seen[target] {
		return false
	}
	seen[target] = true
	meta, ok := a.stepMeta[target]
	if !ok {
		return false
	}
	for _, dependency := range meta.dependsOn {
		if a.dependencyCoversSource(dependency, source, seen) {
			return true
		}
	}
	return false
}

func (a *analyzer) dependencyCoversSource(dependency, source string, seen map[string]bool) bool {
	if dependency == source || a.isStepAncestor(dependency, source) {
		return true
	}
	if members := a.parallelGroups[dependency]; len(members) > 0 {
		for _, member := range members {
			if member == source || a.isStepAncestor(member, source) || a.explicitlyDependsOn(member, source, seen) {
				return true
			}
		}
	}
	return a.explicitlyDependsOn(dependency, source, seen)
}

func (a *analyzer) isStepAncestor(ancestor, descendant string) bool {
	ancestorMeta, ancestorOK := a.stepMeta[ancestor]
	descendantMeta, descendantOK := a.stepMeta[descendant]
	if !ancestorOK || !descendantOK || ancestorMeta.workflow != descendantMeta.workflow || len(ancestorMeta.positions) >= len(descendantMeta.positions) {
		return false
	}
	for i, position := range ancestorMeta.positions {
		if position != descendantMeta.positions[i] {
			return false
		}
	}
	return true
}

func (a *analyzer) operationOutputTrust(operation *uws1.Operation, output string, resolution operationResolution) uws1.ContentTrustLevel {
	if a.doc.ContentTrust != nil {
		if declaration := a.doc.ContentTrust.Operations[operation.OperationID]; declaration != nil {
			if output != "" {
				if level, ok := declaration.Outputs[output]; ok {
					return level
				}
			}
			if declaration.Default != "" {
				return declaration.Default
			}
		}
		if operation.SourceDescription != "" {
			if level, ok := a.doc.ContentTrust.SourceDescriptions[operation.SourceDescription]; ok {
				return level
			}
		}
	}
	if resolution.resolved && validTrust(resolution.contract.DefaultTrust) {
		return resolution.contract.DefaultTrust
	}
	return uws1.ContentTrustUnknown
}

func (a *analyzer) workflowInputTrust(workflow, input string) uws1.ContentTrustLevel {
	if a.doc.ContentTrust != nil {
		if declaration := a.doc.ContentTrust.Workflows[workflow]; declaration != nil {
			if level, ok := declaration.Inputs[input]; ok {
				return level
			}
			if declaration.Default != "" {
				return declaration.Default
			}
		}
	}
	return uws1.ContentTrustUnknown
}

func (a *analyzer) addEdge(state valueState, to string) {
	if !a.emit || state.from == "" {
		return
	}
	edge := Edge{From: state.from, To: to, Provenance: normalizeTrust(state.provenance), Capability: normalizeCapability(state.capability)}
	key := edge.From + "\x00" + edge.To + "\x00" + string(edge.Provenance) + "\x00" + string(edge.Capability)
	if a.edgeKeys[key] {
		return
	}
	a.edgeKeys[key] = true
	a.edges = append(a.edges, edge)
}

func (a *analyzer) addFinding(code, path string) {
	if !a.emit && code != CodeResolverFailure && code != CodeResolverConflict {
		return
	}
	severity := SeverityWarning
	if code == CodeUntrustedInstruction || code == CodeUntrustedAuthority {
		severity = SeverityHigh
	}
	key := path + "\x00" + code
	if a.findingKeys[key] {
		return
	}
	a.findingKeys[key] = true
	a.findings = append(a.findings, Finding{Code: code, Severity: severity, Path: path, Message: findingMessages[code]})
}

func (a *analyzer) setNestedState(target map[string]map[string]valueState, outer, inner string, state valueState) {
	if target[outer] == nil {
		target[outer] = make(map[string]valueState)
	}
	previous, exists := target[outer][inner]
	if !exists {
		target[outer][inner] = state
		a.changed = true
		return
	}
	joined := joinState(previous, state)
	joined.from = previous.from
	if joined != previous {
		target[outer][inner] = joined
		a.changed = true
	}
}

func operationValueAtPointer(operation *uws1.Operation, pointer string) (any, bool) {
	data, err := json.Marshal(operation)
	if err != nil {
		return nil, false
	}
	var root any
	if json.Unmarshal(data, &root) != nil {
		return nil, false
	}
	if pointer == "" || pointer == "/" {
		pointer = "/request"
	} else if !strings.HasPrefix(pointer, "/request") && !strings.HasPrefix(pointer, "/x-") {
		pointer = "/request" + pointer
	}
	return resolveJSONPointer(root, pointer)
}

func resolveJSONPointer(root any, pointer string) (any, bool) {
	if pointer == "" {
		return root, true
	}
	if !strings.HasPrefix(pointer, "/") {
		return nil, false
	}
	current := root
	for _, raw := range strings.Split(pointer[1:], "/") {
		token := strings.ReplaceAll(strings.ReplaceAll(raw, "~1", "/"), "~0", "~")
		switch typed := current.(type) {
		case map[string]any:
			value, ok := typed[token]
			if !ok {
				return nil, false
			}
			current = value
		case []any:
			index, err := strconv.Atoi(token)
			if err != nil || index < 0 || index >= len(typed) {
				return nil, false
			}
			current = typed[index]
		default:
			return nil, false
		}
	}
	return current, true
}

func pointerDocumentPath(base, pointer string) string {
	if pointer == "" || pointer == "/" {
		return base + ".request"
	}
	if !strings.HasPrefix(pointer, "/request") && !strings.HasPrefix(pointer, "/x-") {
		pointer = "/request" + pointer
	}
	return base + pointerSuffix(pointer)
}

func pointerSuffix(pointer string) string {
	if pointer == "" {
		return ""
	}
	var result strings.Builder
	for _, raw := range strings.Split(strings.TrimPrefix(pointer, "/"), "/") {
		token := strings.ReplaceAll(strings.ReplaceAll(raw, "~1", "/"), "~0", "~")
		if _, err := strconv.Atoi(token); err == nil {
			result.WriteString("[")
			result.WriteString(token)
			result.WriteString("]")
		} else {
			result.WriteString(".")
			result.WriteString(token)
		}
	}
	return result.String()
}

func literalState(raw any, from string) valueState {
	state := valueState{provenance: uws1.ContentTrustTrusted, from: from}
	if raw == nil {
		state.capability = CapabilityConstrainedScalar
		return state
	}
	value := reflect.ValueOf(raw)
	switch value.Kind() {
	case reflect.String:
		state.capability = CapabilityFreeText
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		state.capability = CapabilityConstrainedScalar
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Struct:
		state.capability = CapabilityComposite
	default:
		state.capability = CapabilityUnknown
	}
	return state
}

func trustedState(from string) valueState {
	return valueState{provenance: uws1.ContentTrustTrusted, capability: CapabilityUnknown, from: from}
}

func unknownState(from string) valueState {
	return valueState{provenance: uws1.ContentTrustUnknown, capability: CapabilityUnknown, from: from}
}

func joinState(left, right valueState) valueState {
	if left.from == "" {
		return right
	}
	if right.from == "" {
		return left
	}
	left.provenance = joinTrust(left.provenance, right.provenance)
	left.capability = joinCapability(left.capability, right.capability)
	return left
}

func joinTrust(left, right uws1.ContentTrustLevel) uws1.ContentTrustLevel {
	left, right = normalizeTrust(left), normalizeTrust(right)
	if left == uws1.ContentTrustUntrusted || right == uws1.ContentTrustUntrusted {
		return uws1.ContentTrustUntrusted
	}
	if left == uws1.ContentTrustUnknown || right == uws1.ContentTrustUnknown {
		return uws1.ContentTrustUnknown
	}
	return uws1.ContentTrustTrusted
}

func joinCapability(left, right ValueCapability) ValueCapability {
	left, right = normalizeCapability(left), normalizeCapability(right)
	if left == right {
		return left
	}
	return CapabilityUnknown
}

func normalizeTrust(level uws1.ContentTrustLevel) uws1.ContentTrustLevel {
	if validTrust(level) {
		return level
	}
	return uws1.ContentTrustUnknown
}

func validTrust(level uws1.ContentTrustLevel) bool {
	switch level {
	case uws1.ContentTrustUnknown, uws1.ContentTrustTrusted, uws1.ContentTrustUntrusted:
		return true
	default:
		return false
	}
}

func normalizeCapability(capability ValueCapability) ValueCapability {
	if validCapability(capability) {
		return capability
	}
	return CapabilityUnknown
}

func joinMapStates(values map[string]valueState, from string) valueState {
	state := trustedState(from)
	state.capability = CapabilityComposite
	for _, key := range sortedStateKeys(values) {
		state = joinState(state, values[key])
	}
	state.capability = CapabilityComposite
	return state
}

func capabilityForSchema(schema *uws1.ParamSchema) ValueCapability {
	if schema == nil {
		return CapabilityUnknown
	}
	switch schema.Type {
	case "string":
		return CapabilityFreeText
	case "boolean", "integer", "number", "null":
		return CapabilityConstrainedScalar
	case "object", "array":
		return CapabilityComposite
	default:
		return CapabilityUnknown
	}
}

func cloneStateMap(input map[string]valueState) map[string]valueState {
	if input == nil {
		return nil
	}
	out := make(map[string]valueState, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}

func sortedOperationIDs(values map[string]*uws1.Operation) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedWorkflowIDs(values map[string]*uws1.Workflow) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedStringKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedAnyKeys(values map[string]any) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedStateKeys(values map[string]valueState) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
