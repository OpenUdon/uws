package uws1

import (
	"fmt"
	"strings"
)

func (w *Workflow) validate(path string, idx *documentIndex, result *ValidationResult) {
	if w.WorkflowID == "" {
		result.addError(path+".workflowId", "is required")
	} else if !constructIDPattern.MatchString(w.WorkflowID) {
		result.addError(path+".workflowId", fmt.Sprintf("must match pattern ^[A-Za-z0-9_-]+$; got %s", w.WorkflowID))
	}
	if w.Type == "" {
		result.addError(path+".type", "is required")
	} else if !IsWorkflowType(w.Type) {
		result.addError(path+".type", fmt.Sprintf("%q is not valid", w.Type))
	} else {
		validateStructuralTypeFields(w.Type, w.Items, w.Wait, len(w.Cases) > 0, len(w.Default) > 0, path, result)
		if RequiresDependsOnForMerge(w.Type) && len(w.DependsOn) == 0 {
			result.addError(path+".dependsOn", "is required and must name at least one upstream construct for merge")
		}
	}
	validateDependencyList(w.DependsOn, path+".dependsOn", idx, result)
	validateOutputs(w.Outputs, path+".outputs", result)
	w.Inputs.validate(path+".inputs", result)
	validateSteps(w.Steps, path+".steps", idx, result)
	validateCases(w.Cases, path+".cases", idx, result)
	validateSteps(w.Default, path+".default", idx, result)
}

// validateStructuralTypeFields enforces §4.5.6.3 constraints on a workflow or
// step that declares a structural type. The caller passes the relevant fields;
// empty strings indicate the field is unset.
func validateStructuralTypeFields(typeName, items, wait string, hasCases, hasDefault bool, path string, result *ValidationResult) {
	trimmedItems := strings.TrimSpace(items)
	if RequiresItems(typeName) {
		if trimmedItems == "" {
			result.addError(path+".items", fmt.Sprintf("is required for %s", typeName))
		}
	} else if trimmedItems != "" {
		result.addError(path+".items", fmt.Sprintf("is not valid on %s", typeName))
	}
	if RequiresWait(typeName) && strings.TrimSpace(wait) == "" {
		result.addError(path+".wait", fmt.Sprintf("is required for %s", typeName))
	}
	if hasCases && !AllowsCases(typeName) {
		result.addError(path+".cases", fmt.Sprintf("is not valid on %s", typeName))
	}
	if hasDefault && !AllowsDefault(typeName) {
		result.addError(path+".default", fmt.Sprintf("is not valid on %s", typeName))
	}
}

func validateSteps(steps []*Step, path string, idx *documentIndex, result *ValidationResult) {
	for i, step := range steps {
		stepPath := fmt.Sprintf("%s[%d]", path, i)
		if step != nil {
			step.validate(stepPath, idx, result)
		}
	}
}

func (s *Step) validate(path string, idx *documentIndex, result *ValidationResult) {
	if s.StepID == "" {
		result.addError(path+".stepId", "is required")
	} else if !constructIDPattern.MatchString(s.StepID) {
		result.addError(path+".stepId", fmt.Sprintf("must match pattern ^[A-Za-z0-9_-]+$; got %s", s.StepID))
	}
	hasOperationRef := s.OperationRef != ""
	hasWorkflow := s.Workflow != ""
	hasType := s.Type != ""
	// Mirror the schema-level oneOf intent: a step is exactly one of an
	// operation reference, a workflow reference, or a structural construct.
	// Two or more of these together are ambiguous and rejected here.
	switch {
	case hasOperationRef && hasWorkflow:
		result.addError(path, "cannot specify both operationRef and workflow")
	case hasOperationRef && hasType:
		result.addError(path, "operationRef cannot be combined with structural type")
	case hasWorkflow && hasType:
		result.addError(path, "workflow cannot be combined with structural type")
	}
	if hasType {
		if !IsWorkflowType(s.Type) {
			result.addError(path+".type", fmt.Sprintf("%q is not valid", s.Type))
		} else {
			validateStructuralTypeFields(s.Type, s.Items, s.Wait, len(s.Cases) > 0, len(s.Default) > 0, path, result)
			if RequiresDependsOnForMerge(s.Type) && len(s.DependsOn) == 0 {
				result.addError(path+".dependsOn", "is required and must name at least one upstream construct for merge")
			}
		}
	}
	if hasOperationRef && idx.operations[s.OperationRef] == nil {
		result.addError(path+".operationRef", fmt.Sprintf("references unknown operationId %q", s.OperationRef))
	}
	isWorkflowReference := hasWorkflow && !hasOperationRef && !hasType
	if isWorkflowReference && idx.workflows[s.Workflow] == nil {
		result.addError(path+".workflow", fmt.Sprintf("references unknown workflowId %q", s.Workflow))
	}
	if isWorkflowReference && (len(s.Steps) > 0 || len(s.Cases) > 0 || len(s.Default) > 0) {
		result.addError(path, "workflow-reference steps cannot also declare structural type or nested child blocks")
	}
	validateDependencyList(s.DependsOn, path+".dependsOn", idx, result)
	validateOutputs(s.Outputs, path+".outputs", result)
	validateSteps(s.Steps, path+".steps", idx, result)
	validateCases(s.Cases, path+".cases", idx, result)
	validateSteps(s.Default, path+".default", idx, result)
}

func validateCases(cases []*Case, path string, idx *documentIndex, result *ValidationResult) {
	for i, c := range cases {
		casePath := fmt.Sprintf("%s[%d]", path, i)
		if c == nil {
			continue
		}
		if c.Name == "" {
			result.addError(casePath+".name", "is required")
		}
		validateSteps(c.Steps, casePath+".steps", idx, result)
	}
}
