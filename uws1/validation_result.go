package uws1

import (
	"fmt"
	"strings"
)

func (r *StructuralResult) validate(path string, idx *documentIndex, seenNames map[string]bool, result *ValidationResult) {
	if r.Name == "" {
		result.addError(path+".name", "is required")
	} else {
		if !componentNamePattern.MatchString(r.Name) {
			result.addError(path+".name", fmt.Sprintf("%q is not valid", r.Name))
		}
		if seenNames[r.Name] {
			result.addError(path+".name", fmt.Sprintf("duplicate result name %q", r.Name))
		}
		seenNames[r.Name] = true
	}
	if r.Kind == "" {
		result.addError(path+".kind", "is required")
	} else if !IsStructuralResultKind(r.Kind) {
		result.addError(path+".kind", fmt.Sprintf("%q is not valid", r.Kind))
	}
	if r.From == "" {
		result.addError(path+".from", "is required")
		return
	}
	workflowID, stepID, hasStep := strings.Cut(r.From, ".")
	if workflowID == "" || !constructIDPattern.MatchString(workflowID) || (hasStep && (stepID == "" || !constructIDPattern.MatchString(stepID))) || strings.Contains(stepID, ".") {
		result.addError(path+".from", fmt.Sprintf("%q is not a valid workflowId or workflowId.stepId", r.From))
		return
	}
	if idx.workflows[workflowID] == nil {
		result.addError(path+".from", fmt.Sprintf("references unknown workflowId %q", workflowID))
		return
	}
	var resolvedType string
	if !hasStep {
		resolvedType = idx.workflowTypes[workflowID]
	} else {
		stepTypes, ok := idx.workflowSteps[workflowID]
		if !ok {
			result.addError(path+".from", fmt.Sprintf("references unknown stepId %q in workflow %q", stepID, workflowID))
			return
		}
		stepType, stepFound := stepTypes[stepID]
		if !stepFound {
			result.addError(path+".from", fmt.Sprintf("references unknown stepId %q in workflow %q", stepID, workflowID))
			return
		}
		resolvedType = stepType
		if resolvedType == "" {
			result.addError(path+".from", fmt.Sprintf("references stepId %q in workflow %q, but that step is not a structural construct", stepID, workflowID))
			return
		}
	}
	if r.Kind != "" && resolvedType != "" && resolvedType != r.Kind {
		result.addError(path+".kind", fmt.Sprintf("kind %q does not match %q type %q", r.Kind, r.From, resolvedType))
	}
}

func (c *Components) validate(path string, result *ValidationResult) {
	for name := range c.Variables {
		if !componentNamePattern.MatchString(name) {
			result.addError(path+".variables."+name, fmt.Sprintf("component name %q is not valid", name))
		}
	}
}
