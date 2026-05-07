package uws1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_CriteriaAndActions(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{{WorkflowID: "next", Type: "parallel"}}
	doc.Operations[0].SuccessCriteria = []*Criterion{
		{Condition: "$response.statusCode == 200"},
		{Condition: "^ok", Type: CriterionRegex, Context: "$response.body"},
	}
	doc.Operations[0].OnFailure = []*FailureAction{
		{Name: "retry", Type: "retry", RetryAfter: 5, RetryLimit: 3},
		{Name: "abort", Type: "end"},
	}
	doc.Operations[0].OnSuccess = []*SuccessAction{
		{Name: "continue", Type: "end"},
		{Name: "route", Type: "goto", WorkflowID: "next"},
	}
	assert.NoError(t, doc.Validate())
}

func TestValidate_CriteriaAndActionsErrors(t *testing.T) {
	doc := validDocument()
	doc.Operations[0].SuccessCriteria = []*Criterion{
		{Type: CriterionSimple},
		{Condition: "^ok", Type: CriterionRegex},
		{Condition: "test", Type: "invalid"},
	}
	doc.Operations[0].OnFailure = []*FailureAction{
		{Name: "retry", Type: "retry"},
		{Name: "bad", Type: "skip"},
		{Name: "goto", Type: "goto"},
	}
	doc.Operations[0].OnSuccess = []*SuccessAction{
		{Name: "bad", Type: "retry"},
		{Name: "goto", Type: "goto"},
	}

	err := doc.Validate()
	assert.ErrorContains(t, err, "condition is required")
	assert.ErrorContains(t, err, "context is required")
	assert.ErrorContains(t, err, "retry requires retryLimit > 0")
	assert.ErrorContains(t, err, "goto requires workflowId or stepId")
	assert.ErrorContains(t, err, "must be end")
}

func TestValidate_ActionTargetsOnlyAllowedForGoto(t *testing.T) {
	t.Run("failure end", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OnFailure = []*FailureAction{
			{Name: "stop", Type: "end", WorkflowID: "main"},
		}
		assert.ErrorContains(t, doc.Validate(), "workflowId/stepId are only valid for goto actions")
	})

	t.Run("failure retry", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OnFailure = []*FailureAction{
			{Name: "retry", Type: "retry", RetryLimit: 1, StepID: "step"},
		}
		assert.ErrorContains(t, doc.Validate(), "workflowId/stepId are only valid for goto actions")
	})

	t.Run("success end", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OnSuccess = []*SuccessAction{
			{Name: "done", Type: "end", StepID: "step"},
		}
		assert.ErrorContains(t, doc.Validate(), "workflowId/stepId are only valid for goto actions")
	})
}
