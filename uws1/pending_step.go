package uws1

import (
	"fmt"
	"strings"
)

// PendingStep declares the purpose, input and output schemas, and expected
// effect of a step that is not yet bound to an executable operation. The
// schemas are contracts, not runtime input values or output expressions.
type PendingStep struct {
	Purpose    string          `json:"purpose" yaml:"purpose" hcl:"purpose,optional"`
	Inputs     *ParamSchema    `json:"inputs" yaml:"inputs" hcl:"inputs,block"`
	Outputs    *ParamSchema    `json:"outputs" yaml:"outputs" hcl:"outputs,block"`
	Effect     OperationEffect `json:"effect" yaml:"effect" hcl:"effect,optional"`
	Extensions map[string]any  `json:"-" yaml:"-" hcl:"extensions,block"`
}

type pendingStepAlias PendingStep

var pendingStepKnownFields = []string{"purpose", "inputs", "outputs", "effect"}

func (p *PendingStep) UnmarshalJSON(data []byte) error {
	var alias pendingStepAlias
	_, extensions, err := unmarshalCoreWithExtensions(data, "pendingStep", pendingStepKnownFields, &alias)
	if err != nil {
		return err
	}
	*p = PendingStep(alias)
	p.Extensions = extensions
	return nil
}

func (p PendingStep) MarshalJSON() ([]byte, error) {
	alias := pendingStepAlias(p)
	return marshalWithExtensions(&alias, p.Extensions)
}

func (p *PendingStep) validate(path, version string, result *ValidationResult) {
	if p == nil {
		return
	}
	if !supportsUWSVersionAtLeast(version, 1, 12, 0) {
		result.addError(path, "requires UWS 1.12.0 or later")
	}
	if strings.TrimSpace(p.Purpose) == "" {
		result.addError(path+".purpose", "is required")
	}
	validatePendingFieldSet(p.Inputs, path+".inputs", result)
	validatePendingFieldSet(p.Outputs, path+".outputs", result)
	if p.Effect == "" {
		result.addError(path+".effect", "is required")
	} else {
		validateOperationEffect(p.Effect, path+".effect", "1.12.0", result)
	}
}

func validatePendingFieldSet(schema *ParamSchema, path string, result *ValidationResult) {
	if schema == nil {
		result.addError(path, "is required")
		return
	}
	if schema.Type != "object" {
		result.addError(path+".type", fmt.Sprintf("must be %q for a pending-step field set", "object"))
	}
	schema.validate(path, result)
}

func documentHasPendingSteps(d *Document) bool {
	return firstPendingStepInDocument(d) != nil
}

func firstPendingStepInDocument(d *Document) *Step {
	if d == nil {
		return nil
	}
	for _, workflow := range d.Workflows {
		if step := firstPendingStepInWorkflow(workflow); step != nil {
			return step
		}
	}
	return nil
}

func firstPendingStepInWorkflow(workflow *Workflow) *Step {
	if workflow == nil {
		return nil
	}
	if step := firstPendingStepInSteps(workflow.Steps); step != nil {
		return step
	}
	if step := firstPendingStepInCases(workflow.Cases); step != nil {
		return step
	}
	return firstPendingStepInSteps(workflow.Default)
}

func firstPendingStepInSteps(steps []*Step) *Step {
	var found *Step
	_ = walkStepTree("", steps, stepTreeWalkHandlers{
		step: func(_ string, step *Step) error {
			if found == nil && step.Pending != nil {
				found = step
			}
			return nil
		},
	})
	return found
}

func firstPendingStepInCases(cases []*Case) *Step {
	var found *Step
	_ = walkCaseTree("", cases, stepTreeWalkHandlers{
		step: func(_ string, step *Step) error {
			if found == nil && step.Pending != nil {
				found = step
			}
			return nil
		},
	})
	return found
}

func pendingStepExecutionError(step *Step) error {
	return fmt.Errorf("uws1: step %q is pending and cannot be executed", step.StepID)
}
