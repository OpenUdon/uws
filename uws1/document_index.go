package uws1

import "fmt"

type documentIndex struct {
	operations           map[string]*Operation
	workflows            map[string]*Workflow
	workflowTypes        map[string]string
	workflowSteps        map[string]map[string]string
	entryWorkflowID      string
	entryWorkflow        *Workflow
	hasEntryWorkflow     bool
	entryWorkflowSteps   map[string]bool
	steps                map[string]*Step
	topLevelSteps        map[string]*Step
	triggers             map[string]bool
	parallelGroups       map[string]bool
	parallelGroupMembers map[string][]string
	sourceDescriptions   map[string]bool
	dependencies         map[string][]string
}

func newDocumentIndex() *documentIndex {
	return &documentIndex{
		operations:           make(map[string]*Operation),
		workflows:            make(map[string]*Workflow),
		workflowTypes:        make(map[string]string),
		workflowSteps:        make(map[string]map[string]string),
		entryWorkflowSteps:   make(map[string]bool),
		steps:                make(map[string]*Step),
		topLevelSteps:        make(map[string]*Step),
		triggers:             make(map[string]bool),
		parallelGroups:       make(map[string]bool),
		parallelGroupMembers: make(map[string][]string),
		sourceDescriptions:   make(map[string]bool),
		dependencies:         make(map[string][]string),
	}
}

func buildDocumentIndex(d *Document, result *ValidationResult) *documentIndex {
	idx := newDocumentIndex()
	if d == nil {
		return idx
	}
	if entry := selectSemanticEntryWorkflow(d); entry != nil {
		idx.entryWorkflowID = entry.WorkflowID
		idx.entryWorkflow = entry
		idx.hasEntryWorkflow = true
		for _, step := range entry.Steps {
			if step != nil && step.StepID != "" {
				idx.entryWorkflowSteps[step.StepID] = true
				idx.topLevelSteps[step.StepID] = step
			}
		}
	}
	idx.collectDocumentIDs(d, result)
	return idx
}

func (idx *documentIndex) collectDocumentIDs(d *Document, result *ValidationResult) {
	for i, sd := range d.SourceDescriptions {
		path := fmt.Sprintf("sourceDescriptions[%d]", i)
		if sd == nil {
			addIndexError(result, path, "is nil")
			continue
		}
		if sd.Name != "" {
			if idx.sourceDescriptions[sd.Name] {
				addIndexError(result, path+".name", fmt.Sprintf("duplicate sourceDescription name %q", sd.Name))
			}
			idx.sourceDescriptions[sd.Name] = true
		}
	}

	for i, op := range d.Operations {
		path := fmt.Sprintf("operations[%d]", i)
		if op == nil {
			addIndexError(result, path, "is nil")
			continue
		}
		if op.OperationID != "" {
			if idx.operations[op.OperationID] != nil {
				addIndexError(result, path+".operationId", fmt.Sprintf("duplicate operationId %q", op.OperationID))
			}
			idx.operations[op.OperationID] = op
			if len(op.DependsOn) > 0 {
				idx.dependencies[op.OperationID] = append(idx.dependencies[op.OperationID], op.DependsOn...)
			}
		}
		if op.ParallelGroup != "" {
			if idx.hasExecutableID(op.ParallelGroup) {
				addIndexError(result, path+".parallelGroup", fmt.Sprintf("parallelGroup %q collides with an executable identifier", op.ParallelGroup))
			}
			idx.parallelGroups[op.ParallelGroup] = true
			if op.OperationID != "" {
				idx.parallelGroupMembers[op.ParallelGroup] = append(idx.parallelGroupMembers[op.ParallelGroup], op.OperationID)
			}
		}
	}

	for i, wf := range d.Workflows {
		path := fmt.Sprintf("workflows[%d]", i)
		if wf == nil {
			addIndexError(result, path, "is nil")
			continue
		}
		if wf.WorkflowID != "" {
			if idx.workflows[wf.WorkflowID] != nil {
				addIndexError(result, path+".workflowId", fmt.Sprintf("duplicate workflowId %q", wf.WorkflowID))
			}
			idx.workflows[wf.WorkflowID] = wf
			idx.workflowTypes[wf.WorkflowID] = wf.Type
			idx.workflowSteps[wf.WorkflowID] = make(map[string]string)
			idx.collectWorkflowStepTypes(wf.WorkflowID, wf.Steps)
			idx.collectWorkflowCaseStepTypes(wf.WorkflowID, wf.Cases)
			idx.collectWorkflowStepTypes(wf.WorkflowID, wf.Default)
			if len(wf.DependsOn) > 0 {
				idx.dependencies[wf.WorkflowID] = append(idx.dependencies[wf.WorkflowID], wf.DependsOn...)
			}
		}
		idx.collectStepIDs(wf.Steps, path+".steps", result)
		idx.collectCaseStepIDs(wf.Cases, path+".cases", result)
		idx.collectStepIDs(wf.Default, path+".default", result)
	}

	for i, trigger := range d.Triggers {
		path := fmt.Sprintf("triggers[%d]", i)
		if trigger == nil {
			addIndexError(result, path, "is nil")
			continue
		}
		if trigger.TriggerID != "" {
			if idx.triggers[trigger.TriggerID] {
				addIndexError(result, path+".triggerId", fmt.Sprintf("duplicate triggerId %q", trigger.TriggerID))
			}
			idx.triggers[trigger.TriggerID] = true
		}
	}
}

func selectSemanticEntryWorkflow(d *Document) *Workflow {
	if d == nil {
		return nil
	}
	for _, wf := range d.Workflows {
		if wf != nil && wf.WorkflowID == "main" {
			return wf
		}
	}
	var only *Workflow
	for _, wf := range d.Workflows {
		if wf == nil {
			continue
		}
		if only != nil {
			return nil
		}
		only = wf
	}
	return only
}

// collectWorkflowStepTypes populates the workflow to stepID to structural-type
// index used when resolving results[].from references. Nil and unnamed steps
// are skipped here; collectStepIDs runs on the same tree and reports nil steps.
func (idx *documentIndex) collectWorkflowStepTypes(workflowID string, steps []*Step) {
	_ = walkStepTree("", steps, stepTreeWalkHandlers{
		step: func(_ string, step *Step) error {
			if step.StepID != "" {
				idx.workflowSteps[workflowID][step.StepID] = step.Type
			}
			return nil
		},
	})
}

func (idx *documentIndex) collectWorkflowCaseStepTypes(workflowID string, cases []*Case) {
	_ = walkCaseTree("", cases, stepTreeWalkHandlers{
		step: func(_ string, step *Step) error {
			if step.StepID != "" {
				idx.workflowSteps[workflowID][step.StepID] = step.Type
			}
			return nil
		},
	})
}

func (idx *documentIndex) collectStepIDs(steps []*Step, path string, result *ValidationResult) {
	_ = walkStepTree(path, steps, stepTreeWalkHandlers{
		nilStep: func(stepPath string) error {
			addIndexError(result, stepPath, "is nil")
			return nil
		},
		step: func(stepPath string, step *Step) error {
			idx.collectStepID(step, stepPath, result)
			return nil
		},
		nilCase: func(casePath string) error {
			addIndexError(result, casePath, "is nil")
			return nil
		},
	})
}

func (idx *documentIndex) collectCaseStepIDs(cases []*Case, path string, result *ValidationResult) {
	_ = walkCaseTree(path, cases, stepTreeWalkHandlers{
		nilStep: func(stepPath string) error {
			addIndexError(result, stepPath, "is nil")
			return nil
		},
		step: func(stepPath string, step *Step) error {
			idx.collectStepID(step, stepPath, result)
			return nil
		},
		nilCase: func(casePath string) error {
			addIndexError(result, casePath, "is nil")
			return nil
		},
	})
}

func (idx *documentIndex) collectStepID(step *Step, path string, result *ValidationResult) {
	if step.StepID != "" {
		if idx.steps[step.StepID] != nil {
			addIndexError(result, path+".stepId", fmt.Sprintf("duplicate stepId %q", step.StepID))
		}
		idx.steps[step.StepID] = step
		if len(step.DependsOn) > 0 {
			idx.dependencies[step.StepID] = append(idx.dependencies[step.StepID], step.DependsOn...)
		}
	}
	if step.ParallelGroup != "" {
		if idx.hasExecutableID(step.ParallelGroup) {
			addIndexError(result, path+".parallelGroup", fmt.Sprintf("parallelGroup %q collides with an executable identifier", step.ParallelGroup))
		}
		idx.parallelGroups[step.ParallelGroup] = true
		if step.StepID != "" {
			idx.parallelGroupMembers[step.ParallelGroup] = append(idx.parallelGroupMembers[step.ParallelGroup], step.StepID)
		}
	}
}

func (idx *documentIndex) hasExecutableID(name string) bool {
	return idx.operations[name] != nil || idx.workflows[name] != nil || idx.steps[name] != nil
}

func addIndexError(result *ValidationResult, path, message string) {
	if result != nil {
		result.addError(path, message)
	}
}
