package uws1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_WorkflowAndStepReferences(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "wf",
			Type:       "parallel",
			WorkflowExecutionFields: WorkflowExecutionFields{
				DependsOn: []string{"missing_dependency"},
			},
			Steps: []*Step{
				{
					StepID:       "step",
					Type:         "not-a-step-type",
					OperationRef: "missing_operation",
					StepExecutionFields: StepExecutionFields{
						DependsOn: []string{"missing_step"},
					},
				},
			},
		},
	}

	err := doc.Validate()
	assert.ErrorContains(t, err, "references unknown dependency")
	assert.ErrorContains(t, err, "references unknown operationId")
	assert.ErrorContains(t, err, "not-a-step-type")
}

func TestValidate_WorkflowAndStepIDsRejectDots(t *testing.T) {
	t.Run("workflowId", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "daily.v1", Type: "sequence"},
		}
		assert.ErrorContains(t, doc.Validate(), "workflowId must match pattern")
	})

	t.Run("stepId", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "wf",
				Type:       "sequence",
				Steps:      []*Step{{StepID: "fetch.user", OperationRef: "get_data"}},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "stepId must match pattern")
	})
}

func TestValidate_SequenceWorkflowOperationStep(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "main",
			Type:       "sequence",
			Steps: []*Step{
				{
					StepID:       "get_data",
					OperationRef: "get_data",
				},
			},
		},
	}

	assert.NoError(t, doc.Validate())
}

func TestValidate_GotoStepID(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "handler",
			Type:       "parallel",
			Steps:      []*Step{{StepID: "fallback_step", Type: "sequence"}},
		},
	}
	doc.Operations[0].OnFailure = []*FailureAction{
		{Name: "action", Type: "goto", StepID: "fallback_step"},
	}

	assert.NoError(t, doc.Validate())
}

func TestValidate_GotoRejectsBothTargets(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "handler",
			Type:       "parallel",
			Steps:      []*Step{{StepID: "fallback_step", Type: "sequence"}},
		},
	}
	doc.Operations[0].OnFailure = []*FailureAction{
		{Name: "action", Type: "goto", WorkflowID: "handler", StepID: "fallback_step"},
	}

	assert.ErrorContains(t, doc.Validate(), "goto cannot specify both workflowId and stepId")
}

func TestValidate_TriggerRoutes(t *testing.T) {
	t.Run("valid named output", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Outputs:   []string{"primary"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"main"}}},
				},
			},
		}
		doc.Workflows = []*Workflow{{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps:      []*Step{{StepID: "get_data", OperationRef: "get_data"}},
		}}
		assert.NoError(t, doc.Validate())
	})

	t.Run("valid decimal index output", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Outputs:   []string{"primary", "secondary"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "1", To: []string{"step_a"}}},
				},
			},
		}
		doc.Workflows = []*Workflow{{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps:      []*Step{{StepID: "step_a", OperationRef: "get_data"}},
		}}
		assert.NoError(t, doc.Validate())
	})

	t.Run("unknown target", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Outputs:   []string{"primary"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"missing"}}},
				},
			},
		}
		doc.Workflows = []*Workflow{{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps:      []*Step{{StepID: "get_data", OperationRef: "get_data"}},
		}}
		assert.ErrorContains(t, doc.Validate(), "references unknown top-level stepId or workflowId")
	})

	t.Run("empty target list", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Outputs:   []string{"primary"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "primary"}},
				},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "must contain at least one top-level stepId or workflowId")
	})

	t.Run("routes without outputs declaration", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"main"}}},
				},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "outputs is required when routes is set")
	})

	t.Run("route output not declared", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Outputs:   []string{"primary"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "other", To: []string{"main"}}},
				},
			},
		}
		assert.ErrorContains(t, doc.Validate(), `"other" is not a declared trigger output`)
	})

	t.Run("index out of range", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Outputs:   []string{"primary"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "5", To: []string{"main"}}},
				},
			},
		}
		assert.ErrorContains(t, doc.Validate(), `"5" is not a declared trigger output`)
	})

	t.Run("duplicate outputs rejected", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "webhook",
				Outputs:   []string{"primary", "primary"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"main"}}},
				},
			},
		}
		assert.ErrorContains(t, doc.Validate(), `duplicate output "primary"`)
	})

	t.Run("rejects non-top-level step target", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps: []*Step{{
				StepID: "outer",
				Type:   WorkflowTypeSequence,
				Steps:  []*Step{{StepID: "nested", OperationRef: "get_data"}},
			}},
		}}
		doc.Triggers = []*Trigger{{
			TriggerID: "webhook",
			Outputs:   []string{"primary"},
			Routes: []*TriggerRoute{
				{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"nested"}}},
			},
		}}
		assert.ErrorContains(t, doc.Validate(), "references unknown top-level stepId or workflowId")
	})

	t.Run("ambiguous entry workflow defers step-target validation to executable layer", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "alpha",
				Type:       WorkflowTypeSequence,
				Steps:      []*Step{{StepID: "start", OperationRef: "get_data"}},
			},
			{
				WorkflowID: "beta",
				Type:       WorkflowTypeSequence,
				Steps:      []*Step{{StepID: "other", OperationRef: "get_data"}},
			},
		}
		doc.Triggers = []*Trigger{{
			TriggerID: "webhook",
			Outputs:   []string{"primary"},
			Routes: []*TriggerRoute{
				{TriggerRouteFields: TriggerRouteFields{Output: "primary", To: []string{"start"}}},
			},
		}}

		err := doc.Validate()
		require.NoError(t, err)
		require.ErrorContains(t, doc.ValidateExecutable(), `multiple workflows require an explicit "main" entry workflow`)
	})
}

func TestValidate_StructuralTypeConstraints(t *testing.T) {
	t.Run("loop requires items", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "loop_wf", Type: WorkflowTypeLoop},
		}
		assert.ErrorContains(t, doc.Validate(), "items is required for loop")
	})

	t.Run("loop with items is valid", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "loop_wf", Type: WorkflowTypeLoop, StructuralFields: StructuralFields{Items: "$variables.tags"}},
		}
		assert.NoError(t, doc.Validate())
	})

	t.Run("await requires wait", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "await_wf", Type: WorkflowTypeAwait},
		}
		assert.ErrorContains(t, doc.Validate(), "wait is required for await")
	})

	t.Run("await with wait is valid", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "await_wf", Type: WorkflowTypeAwait, WorkflowExecutionFields: WorkflowExecutionFields{Wait: "$response.statusCode == 200"}},
		}
		assert.NoError(t, doc.Validate())
	})

	t.Run("switch rejects items", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "sw", Type: WorkflowTypeSwitch, StructuralFields: StructuralFields{Items: "$variables.tags"}},
		}
		assert.ErrorContains(t, doc.Validate(), "items is not valid on switch")
	})

	t.Run("sequence rejects default", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "wf",
				Type:       WorkflowTypeSequence,
				Default:    []*Step{{StepID: "fallback", OperationRef: "get_data"}},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "default is not valid on sequence")
	})

	t.Run("parallel rejects cases", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "wf",
				Type:       WorkflowTypeParallel,
				Cases:      []*Case{{CaseFields: CaseFields{Name: "a"}}},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "cases is not valid on parallel")
	})

	t.Run("switch allows cases and default", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "wf",
				Type:       WorkflowTypeSwitch,
				Cases:      []*Case{{CaseFields: CaseFields{Name: "premium"}}},
				Default:    []*Step{{StepID: "fallback", OperationRef: "get_data"}},
			},
		}
		assert.NoError(t, doc.Validate())
	})

	t.Run("step with loop type requires items", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "wf",
				Type:       WorkflowTypeParallel,
				Steps:      []*Step{{StepID: "loop_step", Type: WorkflowTypeLoop}},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "items is required for loop")
	})

	t.Run("merge workflow requires dependsOn", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "merge_wf", Type: WorkflowTypeMerge},
		}
		assert.ErrorContains(t, doc.Validate(), "dependsOn is required and must name at least one upstream construct for merge")
	})

	t.Run("merge step requires dependsOn", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "wf",
				Type:       WorkflowTypeParallel,
				Steps:      []*Step{{StepID: "merge_step", Type: WorkflowTypeMerge}},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "dependsOn is required and must name at least one upstream construct for merge")
	})
}

func TestValidate_StructuralResult(t *testing.T) {
	baseDoc := func() *Document {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{
				WorkflowID: "wf_merge",
				Type:       WorkflowTypeMerge,
				WorkflowExecutionFields: WorkflowExecutionFields{
					DependsOn: []string{"get_data"},
				},
			},
			{
				WorkflowID: "wf_parallel",
				Type:       WorkflowTypeParallel,
				Steps: []*Step{{StepID: "merge_step", Type: WorkflowTypeMerge, StepExecutionFields: StepExecutionFields{
					DependsOn: []string{"get_data"},
				}}},
			},
		}
		return doc
	}

	t.Run("valid workflow reference", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: StructuralResultKindMerge, From: "wf_merge"},
		}
		assert.NoError(t, doc.Validate())
	})

	t.Run("valid step reference", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: StructuralResultKindMerge, From: "wf_parallel.merge_step"},
		}
		assert.NoError(t, doc.Validate())
	})

	t.Run("missing from", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: StructuralResultKindMerge},
		}
		assert.ErrorContains(t, doc.Validate(), "from is required")
	})

	t.Run("missing name", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Kind: "merge", From: "wf_merge"},
		}
		assert.ErrorContains(t, doc.Validate(), "name is required")
	})

	t.Run("missing kind", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", From: "wf_merge"},
		}
		assert.ErrorContains(t, doc.Validate(), "kind is required")
	})

	t.Run("invalid kind", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: "parallel", From: "wf_merge"},
		}
		assert.ErrorContains(t, doc.Validate(), `"parallel" is not valid`)
	})

	t.Run("unknown workflow", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: "merge", From: "missing_wf"},
		}
		assert.ErrorContains(t, doc.Validate(), `references unknown workflowId "missing_wf"`)
	})

	t.Run("invalid from shape", func(t *testing.T) {
		for _, from := range []string{"a..b", "a.b.c", ".step", "wf."} {
			doc := baseDoc()
			doc.Results = []*StructuralResult{
				{Name: "out", Kind: "merge", From: from},
			}
			assert.ErrorContains(t, doc.Validate(), "is not a valid workflowId or workflowId.stepId")
		}
	})

	t.Run("unknown step in workflow", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: "merge", From: "wf_parallel.missing_step"},
		}
		assert.ErrorContains(t, doc.Validate(), `references unknown stepId "missing_step"`)
	})

	t.Run("operation step is not structural result source", func(t *testing.T) {
		doc := baseDoc()
		doc.Workflows[1].Steps = append(doc.Workflows[1].Steps, &Step{StepID: "fetch", OperationRef: "get_data"})
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: "merge", From: "wf_parallel.fetch"},
		}
		assert.ErrorContains(t, doc.Validate(), "is not a structural construct")
	})

	t.Run("kind mismatch with workflow type", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: "loop", From: "wf_merge"},
		}
		assert.ErrorContains(t, doc.Validate(), `kind "loop" does not match`)
	})

	t.Run("duplicate name", func(t *testing.T) {
		doc := baseDoc()
		doc.Results = []*StructuralResult{
			{Name: "out", Kind: "merge", From: "wf_merge"},
			{Name: "out", Kind: "merge", From: "wf_parallel.merge_step"},
		}
		assert.ErrorContains(t, doc.Validate(), `duplicate result name "out"`)
	})
}
