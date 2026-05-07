package uws1

import "fmt"

var validFailureActionTypes = map[string]bool{
	"end": true, "goto": true, "retry": true,
}

var validSuccessActionTypes = map[string]bool{
	"end": true, "goto": true,
}

func validateFailureActions(actions []*FailureAction, path string, idx *documentIndex, result *ValidationResult) {
	seenNames := make(map[string]bool, len(actions))
	for i, a := range actions {
		actionPath := fmt.Sprintf("%s[%d]", path, i)
		if a == nil {
			result.addError(actionPath, "is nil")
			continue
		}
		validateCommonAction(a.Name, a.Type, a.WorkflowID, a.StepID, actionPath, validFailureActionTypes, "end, goto, or retry", idx, result)
		if a.Name != "" {
			if seenNames[a.Name] {
				result.addError(actionPath+".name", fmt.Sprintf("duplicate action name %q", a.Name))
			}
			seenNames[a.Name] = true
		}
		if a.Type == "retry" && a.RetryLimit <= 0 {
			result.addError(actionPath, "retry requires retryLimit > 0")
		}
		if a.RetryAfter < 0 {
			result.addError(actionPath+".retryAfter", "must be non-negative")
		}
		validateCriteria(a.Criteria, actionPath+".criteria", result)
	}
}

func validateSuccessActions(actions []*SuccessAction, path string, idx *documentIndex, result *ValidationResult) {
	seenNames := make(map[string]bool, len(actions))
	for i, a := range actions {
		actionPath := fmt.Sprintf("%s[%d]", path, i)
		if a == nil {
			result.addError(actionPath, "is nil")
			continue
		}
		validateCommonAction(a.Name, a.Type, a.WorkflowID, a.StepID, actionPath, validSuccessActionTypes, "end or goto", idx, result)
		if a.Name != "" {
			if seenNames[a.Name] {
				result.addError(actionPath+".name", fmt.Sprintf("duplicate action name %q", a.Name))
			}
			seenNames[a.Name] = true
		}
		validateCriteria(a.Criteria, actionPath+".criteria", result)
	}
}

func validateCommonAction(name, actionType, workflowID, stepID, path string, validTypes map[string]bool, typeList string, idx *documentIndex, result *ValidationResult) {
	if name == "" {
		result.addError(path+".name", "is required")
	}
	if actionType == "" {
		result.addError(path+".type", "is required")
	} else if !validTypes[actionType] {
		result.addError(path+".type", fmt.Sprintf("%q is not valid (must be %s)", actionType, typeList))
	}
	validateGotoTarget(actionType, workflowID, stepID, path, idx, result)
}

func validateGotoTarget(actionType, workflowID, stepID, path string, idx *documentIndex, result *ValidationResult) {
	if actionType != "goto" {
		if workflowID != "" || stepID != "" {
			result.addError(path, "workflowId/stepId are only valid for goto actions")
		}
		return
	}
	if workflowID == "" && stepID == "" {
		result.addError(path, "goto requires workflowId or stepId")
		return
	}
	if workflowID != "" && stepID != "" {
		result.addError(path, "goto cannot specify both workflowId and stepId")
	}
	if workflowID != "" && idx.workflows[workflowID] == nil {
		result.addError(path+".workflowId", fmt.Sprintf("references unknown workflowId %q", workflowID))
	}
	if stepID != "" && idx.steps[stepID] == nil {
		result.addError(path+".stepId", fmt.Sprintf("references unknown stepId %q", stepID))
	}
}
