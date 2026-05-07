package uws1

import (
	"fmt"
	"strconv"
)

func (t *Trigger) validate(path string, idx *documentIndex, result *ValidationResult) {
	if t.TriggerID == "" {
		result.addError(path+".triggerId", "is required")
	}
	outputs := make(map[string]bool, len(t.Outputs))
	for i, name := range t.Outputs {
		outputPath := fmt.Sprintf("%s.outputs[%d]", path, i)
		if name == "" {
			result.addError(outputPath, "is required")
			continue
		}
		if !outputNamePattern.MatchString(name) {
			result.addError(outputPath, fmt.Sprintf("output name %q is not valid", name))
			continue
		}
		if outputs[name] {
			result.addError(outputPath, fmt.Sprintf("duplicate output %q", name))
			continue
		}
		outputs[name] = true
	}
	if len(t.Routes) > 0 && len(t.Outputs) == 0 {
		result.addError(path+".outputs", "is required when routes is set")
	}
	for i, route := range t.Routes {
		routePath := fmt.Sprintf("%s.routes[%d]", path, i)
		if route == nil {
			result.addError(routePath, "is nil")
			continue
		}
		route.validate(routePath, t.Outputs, outputs, idx, result)
	}
}

func (r *TriggerRoute) validate(path string, outputList []string, outputs map[string]bool, idx *documentIndex, result *ValidationResult) {
	if r.Output == "" {
		result.addError(path+".output", "is required")
	} else if len(outputList) > 0 && !resolveTriggerOutput(r.Output, outputList, outputs) {
		result.addError(path+".output", fmt.Sprintf("%q is not a declared trigger output", r.Output))
	}
	if len(r.To) == 0 {
		result.addError(path+".to", "must contain at least one top-level stepId or workflowId")
	}
	for i, target := range r.To {
		if target == "" {
			result.addError(fmt.Sprintf("%s.to[%d]", path, i), "is required")
		} else if idx.workflows[target] != nil {
			continue
		} else if idx.hasEntryWorkflow && !idx.entryWorkflowSteps[target] {
			result.addError(fmt.Sprintf("%s.to[%d]", path, i), fmt.Sprintf("references unknown top-level stepId or workflowId %q", target))
		}
	}
}

func resolveTriggerOutput(output string, outputList []string, outputs map[string]bool) bool {
	if outputs[output] {
		return true
	}
	if idx, err := strconv.Atoi(output); err == nil && idx >= 0 && idx < len(outputList) {
		return true
	}
	return false
}
