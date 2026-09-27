package uws1

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadPendingOnly112(t *testing.T) (*Document, []byte) {
	t.Helper()
	data, err := os.ReadFile("../testdata/examples/pending-only.1.12.json")
	require.NoError(t, err)
	var doc Document
	require.NoError(t, json.Unmarshal(data, &doc))
	return &doc, data
}

func TestPendingStep_PublishedSchemaAndSemanticValidation(t *testing.T) {
	published := compileUWSSchema(t)
	legacySchema := compileSchemaFile(t, "../versions/1.11.0.json")
	doc, data := loadPendingOnly112(t)

	require.NoError(t, published.Validate(decodeJSONValue(t, data)))
	legacyValue := decodeJSONValue(t, data).(map[string]any)
	legacyValue["uws"] = "1.11.0"
	legacyValue["operations"] = []any{map[string]any{
		"operationId":             "bound",
		"x-uws-operation-profile": "test.profile",
	}}
	require.Error(t, legacySchema.Validate(legacyValue), "published 1.11 must not admit the pending field")
	require.NoError(t, doc.Validate(), "published pending-only documents need no placeholder operation")

	withoutPending, withoutPendingData := loadPendingOnly112(t)
	withoutPending.Workflows[0].Steps[0].Pending = nil
	assert.ErrorContains(t, withoutPending.Validate(), "operations at least one operation is required")
	withoutPendingData, err := json.Marshal(withoutPending)
	require.NoError(t, err)
	require.NoError(t, published.Validate(decodeJSONValue(t, withoutPendingData)))

	legacy, _ := loadPendingOnly112(t)
	legacy.UWS = "1.11.0"
	legacy.Operations = []*Operation{{OperationID: "bound", Extensions: map[string]any{ExtensionOperationProfile: "test.profile"}}}
	assert.ErrorContains(t, legacy.Validate(), "workflows[0].steps[0].pending requires UWS 1.12.0 or later")
}

func TestPendingStep_PublishedSchemaRejectsMixedForms(t *testing.T) {
	published := compileUWSSchema(t)
	_, data := loadPendingOnly112(t)

	for name, mutate := range map[string]func(map[string]any, map[string]any){
		"operation reference": func(step, _ map[string]any) { step["operationRef"] = "unbound" },
		"step inputs":         func(step, _ map[string]any) { step["inputs"] = map[string]any{"page_size": 10} },
		"structural type":     func(step, _ map[string]any) { step["type"] = "sequence" },
		"child steps":         func(step, _ map[string]any) { step["steps"] = []any{} },
		"output expressions":  func(step, _ map[string]any) { step["outputs"] = map[string]any{} },
		"missing purpose":     func(_, pending map[string]any) { delete(pending, "purpose") },
		"missing inputs":      func(_, pending map[string]any) { delete(pending, "inputs") },
		"wrong input root":    func(_, pending map[string]any) { pending["inputs"].(map[string]any)["type"] = "array" },
		"invalid effect":      func(_, pending map[string]any) { pending["effect"] = "maybe" },
		"unknown pending field": func(_, pending map[string]any) {
			pending["typo"] = true
		},
		"purpose too long": func(_, pending map[string]any) {
			pending["purpose"] = strings.Repeat("p", 2049)
		},
	} {
		t.Run(name, func(t *testing.T) {
			value := decodeJSONValue(t, data).(map[string]any)
			step := value["workflows"].([]any)[0].(map[string]any)["steps"].([]any)[0].(map[string]any)
			pending := step["pending"].(map[string]any)
			mutate(step, pending)
			require.Error(t, published.Validate(value), "published schema must reject mixed pending/executable form")
		})
	}
}

func TestPendingStep_SemanticValidationAndReferences(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*PendingStep, *Step)
		want   string
	}{
		{
			name:   "blank purpose",
			mutate: func(p *PendingStep, _ *Step) { p.Purpose = " \t " },
			want:   "pending.purpose is required",
		},
		{
			name:   "missing inputs",
			mutate: func(p *PendingStep, _ *Step) { p.Inputs = nil },
			want:   "pending.inputs is required",
		},
		{
			name:   "non-object outputs",
			mutate: func(p *PendingStep, _ *Step) { p.Outputs.Type = "array" },
			want:   `pending.outputs.type must be "object"`,
		},
		{
			name:   "invalid effect",
			mutate: func(p *PendingStep, _ *Step) { p.Effect = "maybe" },
			want:   `pending.effect "maybe" is not valid`,
		},
		{
			name:   "effect required",
			mutate: func(p *PendingStep, _ *Step) { p.Effect = "" },
			want:   "pending.effect is required",
		},
		{
			name:   "mixed operation reference",
			mutate: func(_ *PendingStep, s *Step) { s.OperationRef = "bound" },
			want:   "pending contract cannot be combined with executable or structural fields",
		},
		{
			name:   "broken dependency reference",
			mutate: func(_ *PendingStep, s *Step) { s.DependsOn = []string{"missing"} },
			want:   `dependsOn[0] references unknown dependency "missing"`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doc, _ := loadPendingOnly112(t)
			pending := doc.Workflows[0].Steps[0].Pending
			tc.mutate(pending, doc.Workflows[0].Steps[0])
			assert.ErrorContains(t, doc.Validate(), tc.want)
		})
	}
}

func TestPendingStep_UnknownFieldsAndExtensions(t *testing.T) {
	var pending PendingStep
	require.Error(t, json.Unmarshal([]byte(`{"purpose":"p","inputs":{"type":"object"},"outputs":{"type":"object"},"effect":"read","typo":true}`), &pending))

	doc, _ := loadPendingOnly112(t)
	doc.Workflows[0].Steps[0].Pending.Extensions = map[string]any{"x-openudon-contract": "v1"}
	encoded, err := json.Marshal(doc)
	require.NoError(t, err)
	var decoded Document
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, "v1", decoded.Workflows[0].Steps[0].Pending.Extensions["x-openudon-contract"])
}

func TestPendingStep_ExecutableValidationRejectsEveryBranch(t *testing.T) {
	makePending := func(id string) *Step {
		return &Step{StepID: id, Pending: &PendingStep{
			Purpose: "pending work",
			Inputs:  &ParamSchema{Type: "object"},
			Outputs: &ParamSchema{Type: "object"},
			Effect:  OperationEffectUnknown,
		}}
	}

	cases := map[string][]*Step{
		"top-level":       {makePending("pending")},
		"nested":          {{StepID: "container", Type: WorkflowTypeSequence, Steps: []*Step{makePending("pending")}}},
		"unselected case": {{StepID: "choice", Type: WorkflowTypeSwitch, Cases: []*Case{{CaseFields: CaseFields{Name: "other"}, Steps: []*Step{makePending("pending")}}}}},
		"default branch":  {{StepID: "choice", Type: WorkflowTypeSwitch, Default: []*Step{makePending("pending")}}},
	}
	for name, steps := range cases {
		t.Run(name, func(t *testing.T) {
			doc := &Document{Workflows: []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence, Steps: steps}}}
			require.ErrorContains(t, doc.ValidateExecutable(), `step "pending" is pending and cannot be executed`)
		})
	}
}

func TestPendingStep_ExecutionRejectsBeforeRuntimeCalls(t *testing.T) {
	doc, _ := loadPendingOnly112(t)
	runtime := &pendingRuntimeSpy{}
	doc.SetRuntime(runtime)

	require.ErrorContains(t, doc.ValidateExecutable(), `step "list_projects" is pending and cannot be executed`)
	require.ErrorContains(t, NewOrchestrator(doc, runtime).Execute(t.Context()), `step "list_projects" is pending and cannot be executed`)
	require.ErrorContains(t, NewOrchestrator(doc, runtime).ExecuteWorkflow(t.Context(), doc.Workflows[0]), `step "list_projects" is pending and cannot be executed`)
	require.ErrorContains(t, NewOrchestrator(doc, runtime).ExecuteStep(t.Context(), doc.Workflows[0].Steps[0]), `step "list_projects" is pending and cannot be executed`)
	require.ErrorContains(t, NewOrchestrator(doc, runtime).ExecuteTrigger(t.Context(), "missing", 0, nil), `step "list_projects" is pending and cannot be executed`)
	require.ErrorContains(t, doc.Execute(t.Context()), `step "list_projects" is pending and cannot be executed`)

	otherWorkflow := &Workflow{WorkflowID: "other", Type: WorkflowTypeSequence, Steps: []*Step{{StepID: "elsewhere", Pending: doc.Workflows[0].Steps[0].Pending}}}
	entryDocument := &Document{Workflows: []*Workflow{{WorkflowID: "main", Type: WorkflowTypeSequence}, otherWorkflow}}
	require.ErrorContains(t, NewOrchestrator(entryDocument, runtime).Execute(t.Context()), `step "elsewhere" is pending and cannot be executed`)

	legacyDoc := validDocument()
	legacyDoc.UWS = "1.11.0"
	legacyDoc.SetRuntime(runtime)
	orphanStep := &Step{StepID: "orphan", Pending: doc.Workflows[0].Steps[0].Pending}
	require.ErrorContains(t, orphanStep.Execute(t.Context(), legacyDoc), `step "orphan" is pending and cannot be executed`)
	orphanWorkflow := &Workflow{WorkflowID: "orphan_workflow", Type: WorkflowTypeSequence, Steps: []*Step{orphanStep}}
	require.ErrorContains(t, orphanWorkflow.Execute(t.Context(), legacyDoc), `step "orphan" is pending and cannot be executed`)
	assert.Zero(t, runtime.calls)
}

type pendingRuntimeSpy struct{ calls int }

func (r *pendingRuntimeSpy) ExecuteLeaf(context.Context, *Operation) error {
	r.calls++
	return nil
}

func (r *pendingRuntimeSpy) EvaluateExpression(context.Context, string) (any, error) {
	r.calls++
	return nil, nil
}

func (r *pendingRuntimeSpy) ResolveItems(context.Context, string) ([]any, error) {
	r.calls++
	return nil, nil
}
