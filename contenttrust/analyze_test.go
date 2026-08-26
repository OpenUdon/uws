package contenttrust

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

type resolverFunc func(context.Context, *uws1.Document, *uws1.Operation) (bool, OperationContract, error)

func (f resolverFunc) ResolveOperation(ctx context.Context, doc *uws1.Document, operation *uws1.Operation) (bool, OperationContract, error) {
	return f(ctx, doc, operation)
}

func operationProfile() map[string]any {
	return map[string]any{uws1.ExtensionOperationProfile: "test.runtime.1"}
}

func pipelineDocument() *uws1.Document {
	return &uws1.Document{
		UWS:  "1.9.1",
		Info: &uws1.Info{Title: "Pipeline", Version: "1.0.0"},
		Operations: []*uws1.Operation{
			{OperationID: "read", Outputs: map[string]string{"body": "$response.body"}, Extensions: operationProfile()},
			{OperationID: "model", Request: map[string]any{"body": map[string]any{"prompt": "$inputs.prompt"}}, Outputs: map[string]string{"summary": "$response.body"}, Extensions: operationProfile()},
			{OperationID: "send", Request: map[string]any{"body": map[string]any{"target": "$inputs.target"}}, Extensions: operationProfile()},
		},
		Workflows: []*uws1.Workflow{{
			WorkflowID: "main", Type: uws1.WorkflowTypeSequence,
			Steps: []*uws1.Step{
				{StepID: "read_step", OperationRef: "read", Outputs: map[string]string{"body": "$outputs.body"}},
				{StepID: "model_step", OperationRef: "model", Inputs: map[string]any{"prompt": "$steps.read_step.outputs.body"}, Outputs: map[string]string{"summary": "$outputs.summary"}},
				{StepID: "send_step", OperationRef: "send", Inputs: map[string]any{"target": "$steps.model_step.outputs.summary"}},
			},
		}},
		ContentTrust: &uws1.ContentTrust{Operations: map[string]*uws1.OperationContentTrust{
			"read": {Outputs: map[string]uws1.ContentTrustLevel{"body": uws1.ContentTrustUntrusted}},
		}},
	}
}

func pipelineResolver(modelKind ChannelKind, readCapability ValueCapability) Resolver {
	return resolverFunc(func(_ context.Context, _ *uws1.Document, operation *uws1.Operation) (bool, OperationContract, error) {
		switch operation.OperationID {
		case "read":
			return true, OperationContract{Outputs: map[string]OutputContract{"body": {Capability: readCapability}}}, nil
		case "model":
			return true, OperationContract{
				Inputs:                  []InputChannel{{Path: "/request/body/prompt", Kind: modelKind}},
				Outputs:                 map[string]OutputContract{"summary": {Capability: CapabilityFreeText}},
				InheritsInputProvenance: true,
			}, nil
		case "send":
			return true, OperationContract{Inputs: []InputChannel{{Path: "/request/body/target", Kind: ChannelAuthority}}}, nil
		default:
			return false, OperationContract{}, nil
		}
	})
}

func TestAnalyzeLLMDataPropagationToAuthority(t *testing.T) {
	report, err := Analyze(context.Background(), pipelineDocument(), pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if hasFinding(report, CodeUntrustedInstruction) {
		t.Fatalf("LLM data channel was treated as an instruction sink: %#v", report.Findings)
	}
	finding := findFinding(report, CodeUntrustedAuthority)
	if finding == nil || finding.Severity != SeverityHigh || finding.Path != "operations[2].request.body.target" {
		t.Fatalf("authority finding = %#v, all findings %#v", finding, report.Findings)
	}
	if !hasEdge(report, "workflows[0].steps[0].outputs.body", "workflows[0].steps[1].inputs.prompt", uws1.ContentTrustUntrusted) {
		t.Fatalf("missing mail-to-model edge: %#v", report.Edges)
	}
	if !hasEdge(report, "workflows[0].steps[1].outputs.summary", "workflows[0].steps[2].inputs.target", uws1.ContentTrustUntrusted) {
		t.Fatalf("missing model-to-authority edge: %#v", report.Edges)
	}
}

func TestAnalyzeInstructionAndConstrainedScalar(t *testing.T) {
	report, err := Analyze(context.Background(), pipelineDocument(), pipelineResolver(ChannelInstruction, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeUntrustedInstruction) {
		t.Fatalf("expected instruction finding: %#v", report.Findings)
	}

	report, err = Analyze(context.Background(), pipelineDocument(), pipelineResolver(ChannelInstruction, CapabilityConstrainedScalar))
	if err != nil {
		t.Fatal(err)
	}
	if hasFinding(report, CodeUntrustedInstruction) {
		t.Fatalf("constrained scalar retained free-text injection capability: %#v", report.Findings)
	}
	if !hasFinding(report, CodeUntrustedAuthority) {
		t.Fatalf("narrowing incorrectly cleared provenance: %#v", report.Findings)
	}
}

func TestInheritedProvenanceDoesNotWidenOutputCapability(t *testing.T) {
	doc := pipelineDocument()
	doc.Workflows[0].Steps[2].Inputs = map[string]any{"target": "$steps.model_step.outputs.summary"}
	resolver := resolverFunc(func(_ context.Context, _ *uws1.Document, operation *uws1.Operation) (bool, OperationContract, error) {
		switch operation.OperationID {
		case "read":
			return true, OperationContract{Outputs: map[string]OutputContract{"body": {Capability: CapabilityFreeText}}}, nil
		case "model":
			return true, OperationContract{
				Inputs:                  []InputChannel{{Path: "/request/body/prompt", Kind: ChannelData}},
				Outputs:                 map[string]OutputContract{"summary": {Capability: CapabilityConstrainedScalar}},
				InheritsInputProvenance: true,
			}, nil
		case "send":
			return true, OperationContract{Inputs: []InputChannel{{Path: "/request/body/target", Kind: ChannelInstruction}}}, nil
		default:
			return false, OperationContract{}, nil
		}
	})

	report, err := Analyze(context.Background(), doc, resolver)
	if err != nil {
		t.Fatal(err)
	}
	if hasFinding(report, CodeUntrustedInstruction) {
		t.Fatalf("provenance inheritance widened a constrained output capability: %#v", report.Findings)
	}
	if !hasEdgeWithCapability(report, "workflows[0].steps[1].outputs.summary", "workflows[0].steps[2].inputs.target", uws1.ContentTrustUntrusted, CapabilityConstrainedScalar) {
		t.Fatalf("derived constrained output did not retain provenance and capability: %#v", report.Edges)
	}
}

func TestAnalyzeTriggerControlAndUnknownWorkflowInput(t *testing.T) {
	doc := pipelineDocument()
	doc.Workflows[0].Steps[1].When = `$trigger.enabled == true`
	doc.Workflows[0].Inputs = &uws1.ParamSchema{Type: "object", Properties: map[string]*uws1.ParamSchema{"external": {Type: "string"}}}
	doc.Workflows[0].When = `$inputs.external == "enabled"`
	report, err := Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeUntrustedControl) {
		t.Fatalf("expected trigger control finding: %#v", report.Findings)
	}
	if !hasFinding(report, CodeUnknownProvenance) {
		t.Fatalf("expected unknown external-input finding: %#v", report.Findings)
	}
}

func TestAnalyzeIdempotencyAndStructuralResultExpressions(t *testing.T) {
	doc := pipelineDocument()
	doc.Workflows[0].Type = uws1.WorkflowTypeLoop
	doc.Workflows[0].Items = "$trigger.items"
	doc.Workflows[0].Idempotency = &uws1.Idempotency{Key: "$trigger.request_id"}
	doc.Results = []*uws1.StructuralResult{{
		Name:  "summary",
		Kind:  uws1.WorkflowTypeLoop,
		From:  "main",
		Value: "$steps.model_step.outputs.summary",
	}}

	report, err := Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if finding := findFinding(report, CodeUntrustedControl); finding == nil || finding.Path != "workflows[0].idempotency.key" {
		t.Fatalf("untrusted idempotency key finding = %#v, all findings %#v", finding, report.Findings)
	}
	if !hasEdge(report, "workflows[0].steps[1].outputs.summary", "results[0].value", uws1.ContentTrustUntrusted) {
		t.Fatalf("missing structural-result edge: %#v", report.Edges)
	}
}

func TestAnalyzeImplicitSequenceAndDominance(t *testing.T) {
	report, err := Analyze(context.Background(), pipelineDocument(), pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if hasFinding(report, CodeNonDominatingRef) {
		t.Fatalf("implicit sequence ordering was ignored: %#v", report.Findings)
	}

	doc := pipelineDocument()
	steps := doc.Workflows[0].Steps
	doc.Workflows[0].Steps = []*uws1.Step{steps[1], steps[0], steps[2]}
	report, err = Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeNonDominatingRef) {
		t.Fatalf("forward reference did not produce dominance finding: %#v", report.Findings)
	}
}

func TestAnalyzeParallelWorkflowRequiresHappensBefore(t *testing.T) {
	doc := pipelineDocument()
	doc.Workflows[0].Type = uws1.WorkflowTypeParallel
	report, err := Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeNonDominatingRef) {
		t.Fatalf("parallel declaration order incorrectly established dominance: %#v", report.Findings)
	}

	doc.Workflows[0].Steps[1].DependsOn = []string{"read_step"}
	report, err = Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if findingAtPath(report, CodeNonDominatingRef, "workflows[0].steps[1].inputs.prompt") {
		t.Fatalf("direct dependency did not establish dominance: %#v", report.Findings)
	}

	doc.Workflows[0].Steps[0].ParallelGroup = "mail_ready"
	doc.Workflows[0].Steps[1].DependsOn = []string{"mail_ready"}
	report, err = Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if findingAtPath(report, CodeNonDominatingRef, "workflows[0].steps[1].inputs.prompt") {
		t.Fatalf("parallel-group dependency did not establish dominance: %#v", report.Findings)
	}
}

func TestAnalyzeBranchLoopAndNestedWorkflow(t *testing.T) {
	doc := pipelineDocument()
	doc.Workflows[0].Steps = []*uws1.Step{
		{StepID: "choice", Type: uws1.WorkflowTypeSwitch, Cases: []*uws1.Case{{
			CaseFields: uws1.CaseFields{Name: "mail", When: `$trigger.kind == "mail"`},
			Steps:      []*uws1.Step{{StepID: "read_step", OperationRef: "read", Outputs: map[string]string{"body": "$outputs.body"}}},
		}}},
		{StepID: "model_step", OperationRef: "model", Inputs: map[string]any{"prompt": "$steps.read_step.outputs.body"}, Outputs: map[string]string{"summary": "$outputs.summary"}},
	}
	report, err := Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeNonDominatingRef) {
		t.Fatalf("branch-only producer incorrectly dominated consumer: %#v", report.Findings)
	}

	doc = pipelineDocument()
	doc.Workflows[0].Steps = []*uws1.Step{{
		StepID: "each", Type: uws1.WorkflowTypeLoop,
		StructuralFields: uws1.StructuralFields{Items: "$trigger.items"},
		Steps: []*uws1.Step{
			{StepID: "read_step", OperationRef: "read", Outputs: map[string]string{"body": "$outputs.body"}},
			{StepID: "model_step", OperationRef: "model", Inputs: map[string]any{"prompt": "$steps.read_step.outputs.body"}},
		},
	}}
	report, err = Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if hasFinding(report, CodeNonDominatingRef) {
		t.Fatalf("earlier step in same loop iteration did not dominate: %#v", report.Findings)
	}

	doc = pipelineDocument()
	doc.Workflows = []*uws1.Workflow{
		{
			WorkflowID: "main", Type: uws1.WorkflowTypeSequence,
			Steps: []*uws1.Step{{StepID: "call_child", StepExecutionFields: uws1.StepExecutionFields{Workflow: "child"}, Inputs: map[string]any{"prompt": "$trigger.body"}}},
		},
		{
			WorkflowID: "child", Type: uws1.WorkflowTypeSequence,
			Inputs: &uws1.ParamSchema{Type: "object", Properties: map[string]*uws1.ParamSchema{"prompt": {Type: "string"}}},
			Steps:  []*uws1.Step{{StepID: "model_step", OperationRef: "model", Inputs: map[string]any{"prompt": "$inputs.prompt"}}},
		},
	}
	doc.ContentTrust.Workflows = map[string]*uws1.WorkflowContentTrust{
		"child": {Inputs: map[string]uws1.ContentTrustLevel{"prompt": uws1.ContentTrustTrusted}},
	}
	report, err = Analyze(context.Background(), doc, pipelineResolver(ChannelInstruction, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeUntrustedInstruction) {
		t.Fatalf("workflow entry declaration laundered internally passed trigger content: %#v", report.Findings)
	}
}

func TestAnalyzeOpaqueAndResolverReferences(t *testing.T) {
	doc := pipelineDocument()
	doc.Operations[1].Request["body"] = map[string]any{"prompt": "prefix ${trigger.body}"}
	resolver := pipelineResolver(ChannelInstruction, CapabilityFreeText)
	report, err := Analyze(context.Background(), doc, resolver)
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeOpaqueExpression) {
		t.Fatalf("expected opaque-expression finding: %#v", report.Findings)
	}

	resolver = resolverFunc(func(_ context.Context, _ *uws1.Document, operation *uws1.Operation) (bool, OperationContract, error) {
		if operation.OperationID != "model" {
			return pipelineResolver(ChannelData, CapabilityFreeText).ResolveOperation(context.Background(), doc, operation)
		}
		return true, OperationContract{Inputs: []InputChannel{{
			Path: "/request/body/prompt", Kind: ChannelInstruction,
			References: []Reference{{Expression: "$trigger", Path: "/interpolation/0"}},
		}}}, nil
	})
	report, err = Analyze(context.Background(), doc, resolver)
	if err != nil {
		t.Fatal(err)
	}
	if hasFinding(report, CodeOpaqueExpression) {
		t.Fatalf("resolver-supplied references did not close opaque boundary: %#v", report.Findings)
	}
	if !hasFinding(report, CodeUntrustedInstruction) {
		t.Fatalf("resolver reference did not propagate trigger provenance: %#v", report.Findings)
	}
}

func TestAnalyzeResolverDiagnosticsCancellationAndStableReport(t *testing.T) {
	doc := pipelineDocument()
	doc.Variables = map[string]any{"reviewed_note": "runtime-content-must-not-appear"}
	failing := resolverFunc(func(context.Context, *uws1.Document, *uws1.Operation) (bool, OperationContract, error) {
		return false, OperationContract{}, errors.New("sensitive resolver detail")
	})
	report, err := Analyze(context.Background(), doc, failing)
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeResolverFailure) {
		t.Fatalf("expected resolver failure: %#v", report.Findings)
	}
	encoded, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "sensitive resolver detail") || strings.Contains(string(encoded), "runtime-content-must-not-appear") {
		t.Fatalf("report leaked resolver or document content: %s", encoded)
	}

	claim := resolverFunc(func(context.Context, *uws1.Document, *uws1.Operation) (bool, OperationContract, error) {
		return true, OperationContract{}, nil
	})
	report, err = Analyze(context.Background(), doc, claim, claim)
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeResolverConflict) {
		t.Fatalf("expected resolver conflict: %#v", report.Findings)
	}

	first, err := Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	second, err := Analyze(context.Background(), doc, pipelineResolver(ChannelData, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("reports are not stable:\n%#v\n%#v", first, second)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if report, err := Analyze(canceled, doc); !errors.Is(err, context.Canceled) || report != nil {
		t.Fatalf("canceled Analyze = (%#v, %v), want (nil, context.Canceled)", report, err)
	}
}

func TestAnalyzeProgrammaticContainersMatchWireDocuments(t *testing.T) {
	typed := pipelineDocument()
	typed.Operations[1].Request = map[string]any{
		"body": map[string]string{"prompt": "$inputs.prompt"},
	}
	typed.Workflows[0].Steps[1].Inputs = map[string]any{
		"prompt": map[string]string{"nested": "$steps.read_step.outputs.body"},
	}

	data, err := json.Marshal(typed)
	if err != nil {
		t.Fatal(err)
	}
	var wire uws1.Document
	if err := json.Unmarshal(data, &wire); err != nil {
		t.Fatal(err)
	}

	typedReport, err := Analyze(context.Background(), typed)
	if err != nil {
		t.Fatal(err)
	}
	wireReport, err := Analyze(context.Background(), &wire)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(typedReport, wireReport) {
		t.Fatalf("programmatic and wire reports differ:\n%#v\n%#v", typedReport, wireReport)
	}
	if !hasEdge(typedReport, "workflows[0].steps[0].outputs.body", "workflows[0].steps[1].inputs.prompt.nested", uws1.ContentTrustUntrusted) {
		t.Fatalf("typed nested map hid an expression edge: %#v", typedReport.Edges)
	}
}

func TestAnalyzeUnsupportedProgrammaticValueIsUnknown(t *testing.T) {
	doc := pipelineDocument()
	doc.Workflows[0].Steps[1].Inputs = map[string]any{"prompt": func() {}}
	report, err := Analyze(context.Background(), doc)
	if err != nil {
		t.Fatal(err)
	}
	if !findingAtPath(report, CodeUnknownProvenance, "workflows[0].steps[1].inputs.prompt") {
		t.Fatalf("unsupported value was not reported as unknown: %#v", report.Findings)
	}
}

func TestAnalyzeRejectsMalformedResolverContracts(t *testing.T) {
	tests := map[string]OperationContract{
		"channel kind": {
			Inputs: []InputChannel{{Path: "/request/body/prompt", Kind: ChannelKind("code")}},
		},
		"channel path syntax": {
			Inputs: []InputChannel{{Path: "request/body/prompt", Kind: ChannelData}},
		},
		"channel path escape": {
			Inputs: []InputChannel{{Path: "/request/body/~2", Kind: ChannelData}},
		},
		"missing channel": {
			Inputs: []InputChannel{{Path: "/request/body/missing", Kind: ChannelData}},
		},
		"reference path": {
			Inputs: []InputChannel{{Path: "/request/body/prompt", Kind: ChannelData, References: []Reference{{Expression: "$trigger", Path: "interpolation/0"}}}},
		},
		"unknown output": {
			Outputs: map[string]OutputContract{"missing": {Capability: CapabilityFreeText}},
		},
		"default trust": {
			DefaultTrust: uws1.ContentTrustLevel("reviewed"),
		},
		"output capability": {
			Outputs: map[string]OutputContract{"summary": {Capability: ValueCapability("text")}},
		},
	}
	for name, contract := range tests {
		t.Run(name, func(t *testing.T) {
			resolver := resolverFunc(func(_ context.Context, _ *uws1.Document, operation *uws1.Operation) (bool, OperationContract, error) {
				if operation.OperationID != "model" {
					return false, OperationContract{}, nil
				}
				return true, contract, nil
			})
			report, err := Analyze(context.Background(), pipelineDocument(), resolver)
			if err != nil {
				t.Fatal(err)
			}
			if !findingAtPath(report, CodeResolverFailure, "operations[1]") {
				t.Fatalf("malformed contract was accepted: %#v", report.Findings)
			}
			if hasFinding(report, CodeResolverConflict) {
				t.Fatalf("malformed contract participated in conflict resolution: %#v", report.Findings)
			}
		})
	}
}

func TestAnalyzeResolverContractFailureFallsBackAndComposes(t *testing.T) {
	doc := pipelineDocument()
	doc.Operations[1].Request = map[string]any{"body": map[string]any{"prompt": "$trigger.body"}}
	invalid := resolverFunc(func(_ context.Context, _ *uws1.Document, operation *uws1.Operation) (bool, OperationContract, error) {
		if operation.OperationID != "model" {
			return false, OperationContract{}, nil
		}
		return true, OperationContract{Inputs: []InputChannel{{Path: "/request/body/prompt", Kind: ChannelKind("invalid")}}}, nil
	})

	report, err := Analyze(context.Background(), doc, invalid)
	if err != nil {
		t.Fatal(err)
	}
	if !findingAtPath(report, CodeResolverFailure, "operations[1]") {
		t.Fatalf("expected invalid-contract finding: %#v", report.Findings)
	}
	if !hasEdge(report, "$trigger", "operations[1].request.body.prompt", uws1.ContentTrustUntrusted) {
		t.Fatalf("invalid resolver claim suppressed core scanning: %#v", report.Edges)
	}

	report, err = Analyze(context.Background(), pipelineDocument(), invalid, pipelineResolver(ChannelInstruction, CapabilityFreeText))
	if err != nil {
		t.Fatal(err)
	}
	if !findingAtPath(report, CodeResolverFailure, "operations[1]") || !hasFinding(report, CodeUntrustedInstruction) {
		t.Fatalf("valid claim was not used alongside malformed claim: %#v", report.Findings)
	}
	if hasFinding(report, CodeResolverConflict) {
		t.Fatalf("malformed claim caused a resolver conflict: %#v", report.Findings)
	}
}

func TestAnalyzeNilResolversAndEmptyOptionalContracts(t *testing.T) {
	var typedNil resolverFunc
	report, err := Analyze(context.Background(), pipelineDocument(), nil, typedNil)
	if err != nil {
		t.Fatal(err)
	}
	if !hasFinding(report, CodeResolverFailure) {
		t.Fatalf("nil resolvers did not produce advisory failures: %#v", report.Findings)
	}

	emptyOptional := resolverFunc(func(_ context.Context, _ *uws1.Document, operation *uws1.Operation) (bool, OperationContract, error) {
		if operation.OperationID != "model" {
			return false, OperationContract{}, nil
		}
		return true, OperationContract{
			Inputs:  []InputChannel{{Path: "/request/body/prompt", Kind: ChannelData}},
			Outputs: map[string]OutputContract{"summary": {}},
		}, nil
	})
	report, err = Analyze(context.Background(), pipelineDocument(), emptyOptional)
	if err != nil {
		t.Fatal(err)
	}
	if findingAtPath(report, CodeResolverFailure, "operations[1]") {
		t.Fatalf("empty optional trust or capability was rejected: %#v", report.Findings)
	}
}

func TestOperationTrustPrecedence(t *testing.T) {
	doc := &uws1.Document{
		UWS:  "1.9.1",
		Info: &uws1.Info{Title: "Precedence", Version: "1"},
		SourceDescriptions: []*uws1.SourceDescription{{
			Name: "api", URL: "./api.yaml", Type: uws1.SourceDescriptionTypeOpenAPI,
		}},
		Operations: []*uws1.Operation{
			{OperationID: "declared", SourceDescription: "api", SourceOperationID: "declared", Outputs: map[string]string{"specific": "$response.body", "fallback": "$response.body"}},
			{OperationID: "source", SourceDescription: "api", SourceOperationID: "source", Outputs: map[string]string{"value": "$response.body"}},
			{OperationID: "resolver", Extensions: operationProfile(), Outputs: map[string]string{"value": "$response.body"}},
		},
		ContentTrust: &uws1.ContentTrust{
			SourceDescriptions: map[string]uws1.ContentTrustLevel{"api": uws1.ContentTrustUntrusted},
			Operations: map[string]*uws1.OperationContentTrust{"declared": {
				Default: uws1.ContentTrustUnknown,
				Outputs: map[string]uws1.ContentTrustLevel{"specific": uws1.ContentTrustTrusted},
			}},
		},
	}
	a := newAnalyzer(context.Background(), doc)
	resolverResolution := operationResolution{resolved: true, contract: OperationContract{DefaultTrust: uws1.ContentTrustTrusted}}
	if got := a.operationOutputTrust(doc.Operations[0], "specific", resolverResolution); got != uws1.ContentTrustTrusted {
		t.Fatalf("per-output trust = %q", got)
	}
	if got := a.operationOutputTrust(doc.Operations[0], "fallback", resolverResolution); got != uws1.ContentTrustUnknown {
		t.Fatalf("operation default trust = %q", got)
	}
	if got := a.operationOutputTrust(doc.Operations[1], "value", resolverResolution); got != uws1.ContentTrustUntrusted {
		t.Fatalf("source trust = %q", got)
	}
	if got := a.operationOutputTrust(doc.Operations[2], "value", resolverResolution); got != uws1.ContentTrustTrusted {
		t.Fatalf("resolver default trust = %q", got)
	}
	if got := a.operationOutputTrust(doc.Operations[2], "value", operationResolution{}); got != uws1.ContentTrustUnknown {
		t.Fatalf("implicit trust = %q", got)
	}
}

func hasFinding(report *Report, code string) bool { return findFinding(report, code) != nil }

func findFinding(report *Report, code string) *Finding {
	for i := range report.Findings {
		if report.Findings[i].Code == code {
			return &report.Findings[i]
		}
	}
	return nil
}

func findingAtPath(report *Report, code, path string) bool {
	for _, finding := range report.Findings {
		if finding.Code == code && finding.Path == path {
			return true
		}
	}
	return false
}

func hasEdge(report *Report, from, to string, provenance uws1.ContentTrustLevel) bool {
	for _, edge := range report.Edges {
		if edge.From == from && edge.To == to && edge.Provenance == provenance {
			return true
		}
	}
	return false
}

func hasEdgeWithCapability(report *Report, from, to string, provenance uws1.ContentTrustLevel, capability ValueCapability) bool {
	for _, edge := range report.Edges {
		if edge.From == from && edge.To == to && edge.Provenance == provenance && edge.Capability == capability {
			return true
		}
	}
	return false
}
