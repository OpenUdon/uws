package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/OpenUdon/uws/convert"
	"github.com/OpenUdon/uws/runtimes"
	"github.com/OpenUdon/uws/uws1"
)

func main() {
	if err := writeFixture(); err != nil {
		panic(err)
	}
}

func writeFixture() error {
	doc := BuildDocument()
	if err := doc.Validate(); err != nil {
		return fmt.Errorf("validate source document: %w", err)
	}

	jsonData, err := convert.MarshalJSONIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	hclData, err := convert.MarshalHCL(doc)
	if err != nil {
		return fmt.Errorf("marshal HCL: %w", err)
	}
	yamlData, err := convert.MarshalYAML(doc)
	if err != nil {
		return fmt.Errorf("marshal YAML: %w", err)
	}
	if err := verifyRoundTrips(jsonData, hclData, yamlData); err != nil {
		return err
	}

	_, source, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("resolve fixture directory")
	}
	dir := filepath.Dir(source)
	if err := os.WriteFile(filepath.Join(dir, "big.json"), append(jsonData, '\n'), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "big.hcl"), hclData, 0644); err != nil {
		return err
	}
	return nil
}

func verifyRoundTrips(jsonData, hclData, yamlData []byte) error {
	var fromJSON uws1.Document
	if err := convert.UnmarshalJSON(jsonData, &fromJSON); err != nil {
		return fmt.Errorf("unmarshal JSON: %w", err)
	}
	if err := fromJSON.Validate(); err != nil {
		return fmt.Errorf("validate JSON document: %w", err)
	}
	jsonRoundTrip, err := convert.MarshalJSONIndent(&fromJSON, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON round-trip: %w", err)
	}
	if !bytes.Equal(jsonData, jsonRoundTrip) {
		return fmt.Errorf("JSON round-trip changed canonical document")
	}

	var fromYAML uws1.Document
	if err := convert.UnmarshalYAML(yamlData, &fromYAML); err != nil {
		return fmt.Errorf("unmarshal YAML: %w", err)
	}
	if err := fromYAML.Validate(); err != nil {
		return fmt.Errorf("validate YAML document: %w", err)
	}
	yamlRoundTrip, err := convert.MarshalJSONIndent(&fromYAML, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal YAML round-trip as JSON: %w", err)
	}
	if !bytes.Equal(jsonData, yamlRoundTrip) {
		return fmt.Errorf("YAML round-trip changed canonical document")
	}

	var fromHCL uws1.Document
	if err := convert.UnmarshalHCL(hclData, &fromHCL); err != nil {
		return fmt.Errorf("unmarshal HCL: %w", err)
	}
	if err := fromHCL.Validate(); err != nil {
		return fmt.Errorf("validate HCL document: %w", err)
	}
	hclRoundTrip, err := convert.MarshalJSONIndent(&fromHCL, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal HCL round-trip JSON: %w", err)
	}
	if !bytes.Contains(hclRoundTrip, []byte(`"x-uws-runtime"`)) {
		return fmt.Errorf("HCL round-trip lost x-uws-runtime payloads")
	}
	if !bytes.Equal(jsonData, hclRoundTrip) {
		return fmt.Errorf("HCL round-trip changed canonical document")
	}
	return nil
}

// BuildDocument constructs a deliberately large UWS document from Go structs.
// It models a real incident-response workflow that enriches an alert, runs
// local and remote checks, performs storage and directory operations, sends
// notifications, and produces an executive summary.
func BuildDocument() *uws1.Document {
	ops := []*uws1.Operation{
		apiOperationByID(),
		apiOperationByRef(),
	}

	runtimeTypes := []string{
		runtimes.RuntimeTypeSSH,
		runtimes.RuntimeTypeCmd,
		runtimes.RuntimeTypeFnct,
		runtimes.RuntimeTypeFileIO,
		runtimes.RuntimeTypeSQL,
		runtimes.RuntimeTypeS3,
		runtimes.RuntimeTypeSMTP,
		runtimes.RuntimeTypeDNS,
		runtimes.RuntimeTypeLDAPS,
		runtimes.RuntimeTypeSCP,
		runtimes.RuntimeTypeSFTP,
		runtimes.RuntimeTypeLLM,
	}
	for _, typ := range runtimeTypes {
		ops = append(ops,
			runtimeOperation(typ, "primary"),
			runtimeOperation(typ, "fallback"),
		)
	}

	return &uws1.Document{
		UWS: "1.1.0",
		Info: &uws1.Info{
			Title:       "Incident Response Automation",
			Summary:     "End-to-end security incident enrichment, containment, notification, and evidence packaging.",
			Description: "A large fixture built from Go structs first. It exercises OpenAPI-bound operations, all public runtime supplement selectors, structural workflows, triggers, results, components, and extension preservation across JSON, YAML, and HCL.",
			Version:     "2026.05.08",
			Extensions: map[string]any{
				"x-owner":       "security-automation",
				"x-slo-minutes": 15,
			},
		},
		SourceDescriptions: []*uws1.SourceDescription{
			{
				Name: "incident_api",
				URL:  "./openapi/incident-api.yaml",
				Type: uws1.SourceDescriptionTypeOpenAPI,
				Extensions: map[string]any{
					"x-service-tier": "gold",
					"x-contact":      "secops-api@example.test",
				},
			},
			{
				Name: "crm_api",
				URL:  "./openapi/crm-api.yaml",
				Type: uws1.SourceDescriptionTypeOpenAPI,
				Extensions: map[string]any{
					"x-service-tier": "silver",
					"x-contact":      "customer-ops@example.test",
				},
			},
		},
		Variables: map[string]any{
			"regions":        []any{"us-east-1", "eu-west-1"},
			"severityPolicy": map[string]any{"critical": "page", "high": "ticket"},
			"environment":    "production",
			"dryRun":         false,
		},
		Operations: ops,
		Workflows: []*uws1.Workflow{
			mainWorkflow(),
			parallelWorkflow(),
			switchWorkflow(),
			loopWorkflow(),
			awaitWorkflow(),
			mergeWorkflow(),
			manualReviewWorkflow(),
			notifyOperatorsWorkflow(),
		},
		Triggers: []*uws1.Trigger{
			{
				TriggerID: "alert_webhook",
				TriggerFields: uws1.TriggerFields{
					Path:           "/webhooks/security/alerts",
					Methods:        []string{"POST", "PUT"},
					Authentication: "bearer",
				},
				Options: map[string]any{
					"deduplicateWindowSeconds": 300,
					"acceptedContentTypes":     []any{"application/json", "application/problem+json"},
				},
				Outputs: []string{"alert.received", "alert.replayed"},
				Routes: []*uws1.TriggerRoute{
					{TriggerRouteFields: uws1.TriggerRouteFields{Output: "alert.received", To: []string{"main", "wf_parallel"}}, Extensions: ext("route", "primary")},
					{TriggerRouteFields: uws1.TriggerRouteFields{Output: "alert.replayed", To: []string{"step_collect_context", "wf_switch"}}, Extensions: ext("route", "replay")},
				},
				Extensions: ext("trigger", "alert"),
			},
			{
				TriggerID: "operator_command",
				TriggerFields: uws1.TriggerFields{
					Path:           "/commands/security/remediate",
					Methods:        []string{"POST", "PATCH"},
					Authentication: "signed",
				},
				Options: map[string]any{
					"allowedTeams": []any{"secops", "platform"},
					"requiresMFA":  true,
				},
				Outputs: []string{"command.approved", "command.denied"},
				Routes: []*uws1.TriggerRoute{
					{TriggerRouteFields: uws1.TriggerRouteFields{Output: "command.approved", To: []string{"wf_loop", "wf_await"}}, Extensions: ext("route", "approved")},
					{TriggerRouteFields: uws1.TriggerRouteFields{Output: "command.denied", To: []string{"wf_manual_review", "step_notify"}}, Extensions: ext("route", "denied")},
				},
				Extensions: ext("trigger", "operator"),
			},
		},
		Results: []*uws1.StructuralResult{
			{Name: "decision.branch", Kind: uws1.StructuralResultKindSwitch, From: "wf_switch", Value: "$workflows.wf_switch.outputs.selectedPath", Extensions: ext("result", "switch")},
			{Name: "containment.loop", Kind: uws1.StructuralResultKindLoop, From: "wf_loop", Value: "$workflows.wf_loop.outputs.containmentResults", Extensions: ext("result", "loop")},
			{Name: "merge.summary", Kind: uws1.StructuralResultKindMerge, From: "wf_merge", Value: "$workflows.wf_merge.outputs.summary", Extensions: ext("result", "merge")},
		},
		Components: &uws1.Components{
			Variables: map[string]any{
				"priorityThreshold": 80,
				"notificationLists": map[string]any{
					"executive": []any{"ciso@example.test", "cio@example.test"},
					"operator":  []any{"secops@example.test", "platform@example.test"},
				},
			},
			Extensions: ext("components", "shared"),
		},
		Extensions: map[string]any{
			"x-generator": "testdata/big/main.go",
			"x-purpose":   "large-runtime-roundtrip-fixture",
		},
	}
}

func apiOperationByID() *uws1.Operation {
	timeout := 20.0
	return &uws1.Operation{
		OperationID:        "fetch_ticket",
		SourceDescription:  "incident_api",
		OpenAPIOperationID: "getIncident",
		Description:        "Fetch the incident record and current triage state.",
		Request: map[string]any{
			"path": map[string]any{"incidentId": "$inputs.incidentId", "tenantId": "$inputs.tenantId"},
			"query": map[string]any{
				"include": []any{"timeline", "assets"},
				"depth":   "full",
			},
		},
		OperationExecutionFields: uws1.OperationExecutionFields{
			When:          "$inputs.incidentId != ''",
			ForEach:       "$variables.regions",
			Wait:          "$signals.incident_api_ready",
			Timeout:       &timeout,
			ParallelGroup: "api_fetch_group",
		},
		SuccessCriteria: []*uws1.Criterion{
			criteria("$response.statusCode == 200", uws1.CriterionSimple, ""),
			criteria("$.severity", uws1.CriterionJSONPath, "$response.body"),
		},
		OnFailure: []*uws1.FailureAction{
			failureRetry("retry_fetch_ticket"),
			failureGoto("manual_fetch_review", "wf_manual_review"),
		},
		OnSuccess: []*uws1.SuccessAction{
			successGoto("continue_to_parallel", "wf_parallel"),
			successEnd("ticket_loaded"),
		},
		Outputs: map[string]string{
			"ticket":   "$response.body",
			"severity": "$response.body.severity",
		},
		Extensions: ext("operation", "fetch-ticket"),
	}
}

func apiOperationByRef() *uws1.Operation {
	timeout := 25.0
	return &uws1.Operation{
		OperationID:         "load_customer",
		SourceDescription:   "crm_api",
		OpenAPIOperationRef: "#/paths/~1customers~1{customerId}/get",
		Description:         "Load customer account context used for notification decisions.",
		Request: map[string]any{
			"path":   map[string]any{"customerId": "$steps.step_collect_context.outputs.customerId", "tenantId": "$inputs.tenantId"},
			"header": map[string]any{"X-Trace-ID": "$context.traceId", "X-Workflow": "incident-response"},
		},
		OperationExecutionFields: uws1.OperationExecutionFields{
			When:          "$steps.step_collect_context.outputs.customerId != ''",
			ForEach:       "$variables.regions",
			Wait:          "$signals.crm_ready",
			Timeout:       &timeout,
			ParallelGroup: "api_fetch_group",
		},
		SuccessCriteria: []*uws1.Criterion{
			criteria("$response.statusCode == 200", uws1.CriterionSimple, ""),
			criteria("enterprise|strategic", uws1.CriterionRegex, "$response.body.segment"),
		},
		OnFailure: []*uws1.FailureAction{
			failureRetry("retry_load_customer"),
			failureGoto("manual_customer_review", "wf_manual_review"),
		},
		OnSuccess: []*uws1.SuccessAction{
			successGoto("continue_to_notification", "wf_notify_operators"),
			successEnd("customer_loaded"),
		},
		Outputs: map[string]string{
			"customer": "$response.body",
			"segment":  "$response.body.segment",
		},
		Extensions: ext("operation", "load-customer"),
	}
}

func runtimeOperation(runtimeType, variant string) *uws1.Operation {
	timeout := 30.0
	id := runtimeOperationID(runtimeType, variant)
	rt := &runtimes.OperationRuntime{
		Type:       runtimeType,
		Command:    runtimeCommand(runtimeType, variant),
		WorkingDir: "/srv/incident-response/" + runtimeType,
		Function:   runtimeFunction(runtimeType, variant),
		Workflow:   "runtime/" + runtimeType + "-" + variant + ".uws.hcl",
		Arguments: []any{
			map[string]any{"incidentId": "$inputs.incidentId", "runtime": runtimeType},
			map[string]any{"variant": variant, "region": "$context.region"},
		},
	}
	extensions := ext("runtime-operation", id)
	if err := runtimes.SetOperationExtension(&extensions, rt); err != nil {
		panic(err)
	}
	extensions[uws1.ExtensionOperationProfile] = runtimes.ProfileName

	return &uws1.Operation{
		OperationID: id,
		Description: fmt.Sprintf("Run %s %s runtime task for incident response.", runtimeType, variant),
		Request:     runtimeRequest(runtimeType, variant),
		SuccessCriteria: []*uws1.Criterion{
			criteria("$response.statusCode < 500", uws1.CriterionSimple, ""),
			criteria("$.ok", uws1.CriterionJSONPath, "$response.body"),
		},
		OnFailure: []*uws1.FailureAction{
			failureRetry("retry_" + id),
			failureGoto("review_"+id, "wf_manual_review"),
		},
		OnSuccess: []*uws1.SuccessAction{
			successGoto("notify_"+id, "wf_notify_operators"),
			successEnd("complete_" + id),
		},
		Outputs: map[string]string{
			"result": "$response.body.result",
			"audit":  "$response.body.auditId",
		},
		OperationExecutionFields: uws1.OperationExecutionFields{
			DependsOn:     []string{"fetch_ticket", "load_customer"},
			When:          "$steps.step_collect_context.outputs.enabled == true",
			ForEach:       "$variables.regions",
			Wait:          "$signals.runtime_slot_available",
			Timeout:       &timeout,
			ParallelGroup: "runtime_" + variant + "_group",
		},
		Extensions: extensions,
	}
}

func runtimeOperationID(runtimeType, variant string) string {
	return "run_" + runtimeType + "_" + variant
}

func runtimeCommand(runtimeType, variant string) string {
	switch runtimeType {
	case runtimes.RuntimeTypeSSH:
		return "sudo systemctl status incident-agent --no-pager"
	case runtimes.RuntimeTypeCmd:
		return "incidentctl preflight --format=json --variant=" + variant
	case runtimes.RuntimeTypeSQL:
		return "SELECT id, severity FROM incidents WHERE id = :incident_id"
	case runtimes.RuntimeTypeDNS:
		return "dig +short suspicious.example.test"
	case runtimes.RuntimeTypeSCP:
		return "scp evidence.tar.gz evidence-vault:/incoming/" + variant
	case runtimes.RuntimeTypeSFTP:
		return "put evidence.json /incoming/" + variant + "/evidence.json"
	default:
		return runtimeType + " task " + variant
	}
}

func runtimeFunction(runtimeType, variant string) string {
	switch runtimeType {
	case runtimes.RuntimeTypeFnct:
		return "render_incident_brief_" + variant
	case runtimes.RuntimeTypeFileIO:
		return "write_evidence_bundle_" + variant
	case runtimes.RuntimeTypeS3:
		return "archive_evidence_s3_" + variant
	case runtimes.RuntimeTypeSMTP:
		return "send_customer_notice_" + variant
	case runtimes.RuntimeTypeLDAPS:
		return "lookup_owner_directory_" + variant
	case runtimes.RuntimeTypeLLM:
		return "summarize_incident_" + variant
	default:
		return "execute_" + runtimeType + "_" + variant
	}
}

func runtimeRequest(runtimeType, variant string) map[string]any {
	return map[string]any{
		"body": map[string]any{
			"incidentId": "$inputs.incidentId",
			"runtime":    runtimeType,
			"variant":    variant,
		},
		"header": map[string]any{
			"X-Runtime-Type": runtimeType,
			"X-Run-Variant":  variant,
		},
	}
}

func mainWorkflow() *uws1.Workflow {
	timeout := 120.0
	ttl := 900.0
	return &uws1.Workflow{
		WorkflowID:  "main",
		Type:        uws1.WorkflowTypeSequence,
		Description: "Coordinate enrichment, runtime checks, branching, containment, and notification.",
		Inputs:      fullSchema("incidentInput"),
		Idempotency: &uws1.Idempotency{
			Key:        "$inputs.tenantId + ':' + $inputs.incidentId",
			OnConflict: "returnPrevious",
			TTL:        &ttl,
			Extensions: ext("idempotency", "main"),
		},
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"fetch_ticket", "load_customer"},
			When:      "$inputs.incidentId != ''",
			ForEach:   "$variables.regions",
			Wait:      "$signals.start",
			Timeout:   &timeout,
		},
		Steps: []*uws1.Step{
			operationStep("step_collect_context", "fetch_ticket", []string{"run_cmd_primary", "run_fnct_primary"}),
			workflowStep("step_parallel_checks", "wf_parallel", []string{"step_collect_context", "load_customer"}),
			structuralSwitchStep("step_local_switch", []string{"step_collect_context", "step_parallel_checks"}),
			structuralLoopStep("step_local_loop", []string{"step_collect_context", "step_parallel_checks"}),
			workflowStep("step_decide_path", "wf_switch", []string{"step_local_switch", "step_local_loop"}),
			workflowStep("step_notify", "wf_notify_operators", []string{"step_decide_path", "run_llm_fallback"}),
		},
		Outputs: map[string]string{
			"incident": "$steps.step_collect_context.outputs.incident",
			"decision": "$steps.step_decide_path.outputs.selectedPath",
		},
		Extensions: ext("workflow", "main"),
	}
}

func parallelWorkflow() *uws1.Workflow {
	timeout := 90.0
	return &uws1.Workflow{
		WorkflowID:  "wf_parallel",
		Type:        uws1.WorkflowTypeParallel,
		Description: "Run independent runtime checks in parallel.",
		Inputs:      fullSchema("parallelInput"),
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"fetch_ticket", "load_customer"},
			When:      "$inputs.severity in ['critical','high']",
			ForEach:   "$variables.regions",
			Wait:      "$signals.parallel_ready",
			Timeout:   &timeout,
		},
		Steps: []*uws1.Step{
			operationStep("step_ssh_primary", "run_ssh_primary", []string{"fetch_ticket", "load_customer"}),
			operationStep("step_sql_primary", "run_sql_primary", []string{"fetch_ticket", "load_customer"}),
			operationStep("step_dns_primary", "run_dns_primary", []string{"fetch_ticket", "load_customer"}),
			operationStep("step_llm_primary", "run_llm_primary", []string{"fetch_ticket", "load_customer"}),
		},
		Outputs: map[string]string{
			"health": "$steps.step_ssh_primary.outputs.result",
			"risk":   "$steps.step_llm_primary.outputs.result",
		},
		Extensions: ext("workflow", "parallel"),
	}
}

func switchWorkflow() *uws1.Workflow {
	timeout := 45.0
	return &uws1.Workflow{
		WorkflowID:  "wf_switch",
		Type:        uws1.WorkflowTypeSwitch,
		Description: "Choose containment path based on severity and customer tier.",
		Inputs:      fullSchema("switchInput"),
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"run_llm_primary", "run_sql_primary"},
			When:      "$inputs.severity != 'low'",
			ForEach:   "$variables.regions",
			Wait:      "$signals.decision_ready",
			Timeout:   &timeout,
		},
		Cases: []*uws1.Case{
			caseBranch("case_critical", "$inputs.severity == 'critical'", "run_smtp_primary", "run_s3_primary"),
			caseBranch("case_enterprise", "$steps.step_collect_context.outputs.segment == 'enterprise'", "run_ldaps_primary", "run_llm_primary"),
		},
		Default: []*uws1.Step{
			operationStep("step_switch_default_notify", "run_smtp_fallback", []string{"run_llm_primary", "run_sql_primary"}),
			operationStep("step_switch_default_archive", "run_s3_fallback", []string{"run_llm_primary", "run_sql_primary"}),
		},
		Outputs: map[string]string{
			"selectedPath": "$selected.case.name",
			"nextAction":   "$selected.steps[0].outputs.result",
		},
		Extensions: ext("workflow", "switch"),
	}
}

func loopWorkflow() *uws1.Workflow {
	timeout := 180.0
	return &uws1.Workflow{
		WorkflowID:  "wf_loop",
		Type:        uws1.WorkflowTypeLoop,
		Description: "Apply containment across affected hosts and evidence locations.",
		Inputs:      fullSchema("loopInput"),
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"run_fileio_primary", "run_s3_primary"},
			When:      "$inputs.assets != []",
			ForEach:   "$inputs.assets",
			Wait:      "$signals.loop_ready",
			Timeout:   &timeout,
		},
		StructuralFields: uws1.StructuralFields{
			Items:     "$inputs.assets",
			Mode:      "parallel",
			BatchSize: "2",
		},
		Steps: []*uws1.Step{
			operationStep("step_loop_scp", "run_scp_primary", []string{"run_fileio_primary", "run_s3_primary"}),
			operationStep("step_loop_sftp", "run_sftp_primary", []string{"run_fileio_primary", "run_s3_primary"}),
			operationStep("step_loop_cmd", "run_cmd_primary", []string{"run_fileio_primary", "run_s3_primary"}),
		},
		Outputs: map[string]string{
			"containmentResults": "$steps.*.outputs.result",
			"evidenceIds":        "$steps.step_loop_scp.outputs.audit",
		},
		Extensions: ext("workflow", "loop"),
	}
}

func awaitWorkflow() *uws1.Workflow {
	timeout := 600.0
	return &uws1.Workflow{
		WorkflowID:  "wf_await",
		Type:        uws1.WorkflowTypeAwait,
		Description: "Wait for external approval before applying high-risk actions.",
		Inputs:      fullSchema("awaitInput"),
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"run_smtp_primary", "run_llm_primary"},
			When:      "$inputs.requiresApproval == true",
			ForEach:   "$variables.regions",
			Wait:      "$signals.approval_received",
			Timeout:   &timeout,
		},
		Steps: []*uws1.Step{
			operationStep("step_await_smtp", "run_smtp_fallback", []string{"run_smtp_primary", "run_llm_primary"}),
			operationStep("step_await_llm", "run_llm_fallback", []string{"run_smtp_primary", "run_llm_primary"}),
		},
		Outputs: map[string]string{
			"approval": "$signals.approval_received",
			"summary":  "$steps.step_await_llm.outputs.result",
		},
		Extensions: ext("workflow", "await"),
	}
}

func mergeWorkflow() *uws1.Workflow {
	timeout := 60.0
	return &uws1.Workflow{
		WorkflowID:  "wf_merge",
		Type:        uws1.WorkflowTypeMerge,
		Description: "Merge all branch outputs into a single handoff summary.",
		Inputs:      fullSchema("mergeInput"),
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"wf_switch", "wf_loop"},
			When:      "$workflows.wf_switch.outputs.selectedPath != ''",
			ForEach:   "$variables.regions",
			Wait:      "$signals.merge_ready",
			Timeout:   &timeout,
		},
		Steps: []*uws1.Step{
			operationStep("step_merge_fileio", "run_fileio_fallback", []string{"wf_switch", "wf_loop"}),
			operationStep("step_merge_s3", "run_s3_fallback", []string{"wf_switch", "wf_loop"}),
		},
		Outputs: map[string]string{
			"summary": "$steps.step_merge_fileio.outputs.result",
			"archive": "$steps.step_merge_s3.outputs.audit",
		},
		Extensions: ext("workflow", "merge"),
	}
}

func manualReviewWorkflow() *uws1.Workflow {
	timeout := 300.0
	return &uws1.Workflow{
		WorkflowID:  "wf_manual_review",
		Type:        uws1.WorkflowTypeSequence,
		Description: "Collect evidence for human review when automated policy blocks execution.",
		Inputs:      fullSchema("manualReviewInput"),
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"fetch_ticket", "load_customer"},
			When:      "$inputs.requiresManualReview == true",
			ForEach:   "$variables.regions",
			Wait:      "$signals.reviewer_available",
			Timeout:   &timeout,
		},
		Steps: []*uws1.Step{
			operationStep("step_manual_fnct", "run_fnct_fallback", []string{"fetch_ticket", "load_customer"}),
			operationStep("step_manual_smtp", "run_smtp_fallback", []string{"fetch_ticket", "load_customer"}),
		},
		Outputs: map[string]string{
			"reviewPackage": "$steps.step_manual_fnct.outputs.result",
			"notification":  "$steps.step_manual_smtp.outputs.audit",
		},
		Extensions: ext("workflow", "manual-review"),
	}
}

func notifyOperatorsWorkflow() *uws1.Workflow {
	timeout := 80.0
	return &uws1.Workflow{
		WorkflowID:  "wf_notify_operators",
		Type:        uws1.WorkflowTypeSequence,
		Description: "Notify operators and stakeholders after evidence is prepared.",
		Inputs:      fullSchema("notifyInput"),
		WorkflowExecutionFields: uws1.WorkflowExecutionFields{
			DependsOn: []string{"run_smtp_primary", "run_llm_primary"},
			When:      "$inputs.suppressNotifications != true",
			ForEach:   "$variables.regions",
			Wait:      "$signals.notification_window",
			Timeout:   &timeout,
		},
		Steps: []*uws1.Step{
			operationStep("step_notify_smtp", "run_smtp_primary", []string{"run_smtp_primary", "run_llm_primary"}),
			operationStep("step_notify_llm", "run_llm_primary", []string{"run_smtp_primary", "run_llm_primary"}),
		},
		Outputs: map[string]string{
			"messageId": "$steps.step_notify_smtp.outputs.audit",
			"brief":     "$steps.step_notify_llm.outputs.result",
		},
		Extensions: ext("workflow", "notify"),
	}
}

func operationStep(stepID, operationRef string, deps []string) *uws1.Step {
	timeout := 30.0
	return &uws1.Step{
		StepID:       stepID,
		Description:  "Execute " + operationRef + " and expose normalized outputs.",
		OperationRef: operationRef,
		Body: map[string]any{
			"input": map[string]any{"incidentId": "$inputs.incidentId", "source": stepID},
			"meta":  map[string]any{"correlationId": "$context.correlationId", "attempt": 1},
		},
		StepExecutionFields: uws1.StepExecutionFields{
			DependsOn:     deps,
			When:          "$inputs.enabled != false",
			ForEach:       "$variables.regions",
			Wait:          "$signals.step_ready",
			Timeout:       &timeout,
			ParallelGroup: "steps_" + stepID + "_group",
		},
		Outputs: map[string]string{
			"result": "$response.body.result",
			"audit":  "$response.body.auditId",
		},
		Extensions: ext("step", stepID),
	}
}

func workflowStep(stepID, workflow string, deps []string) *uws1.Step {
	timeout := 40.0
	return &uws1.Step{
		StepID:      stepID,
		Description: "Invoke workflow " + workflow + ".",
		Body: map[string]any{
			"input": map[string]any{"incidentId": "$inputs.incidentId", "workflow": workflow},
			"meta":  map[string]any{"requestedBy": "$context.actor", "stepId": stepID},
		},
		StepExecutionFields: uws1.StepExecutionFields{
			DependsOn:     deps,
			When:          "$inputs.enabled != false",
			ForEach:       "$variables.regions",
			Wait:          "$signals.workflow_step_ready",
			Timeout:       &timeout,
			Workflow:      workflow,
			ParallelGroup: "steps_" + stepID + "_group",
		},
		Outputs: map[string]string{
			"result": "$workflow.outputs.result",
			"audit":  "$workflow.outputs.audit",
		},
		Extensions: ext("step", stepID),
	}
}

func structuralSwitchStep(stepID string, deps []string) *uws1.Step {
	timeout := 50.0
	return &uws1.Step{
		StepID:      stepID,
		Type:        uws1.WorkflowTypeSwitch,
		Description: "Choose a local fast path before invoking the top-level switch workflow.",
		Body: map[string]any{
			"input": map[string]any{"incidentId": "$inputs.incidentId", "step": stepID},
			"meta":  map[string]any{"decision": "local", "owner": "secops"},
		},
		StepExecutionFields: uws1.StepExecutionFields{
			DependsOn:     deps,
			When:          "$inputs.severity != 'low'",
			ForEach:       "$variables.regions",
			Wait:          "$signals.local_switch_ready",
			Timeout:       &timeout,
			ParallelGroup: "steps_" + stepID + "_group",
		},
		Cases: []*uws1.Case{
			caseBranch("local_switch_critical", "$inputs.severity == 'critical'", "run_cmd_primary", "run_fileio_primary"),
			caseBranch("local_switch_customer", "$steps.step_collect_context.outputs.segment == 'strategic'", "run_smtp_primary", "run_llm_primary"),
		},
		Default: []*uws1.Step{
			operationStep("step_local_switch_default_a", "run_fnct_primary", []string{"step_collect_context", "step_parallel_checks"}),
			operationStep("step_local_switch_default_b", "run_dns_primary", []string{"step_collect_context", "step_parallel_checks"}),
		},
		Outputs: map[string]string{
			"selected": "$selected.case.name",
			"audit":    "$selected.steps[0].outputs.audit",
		},
		Extensions: ext("step", stepID),
	}
}

func structuralLoopStep(stepID string, deps []string) *uws1.Step {
	timeout := 70.0
	return &uws1.Step{
		StepID:      stepID,
		Type:        uws1.WorkflowTypeLoop,
		Description: "Loop over suspect indicators and write local evidence before remote upload.",
		Body: map[string]any{
			"input": map[string]any{"indicators": "$inputs.indicators", "step": stepID},
			"meta":  map[string]any{"batch": "local", "owner": "forensics"},
		},
		StepExecutionFields: uws1.StepExecutionFields{
			DependsOn:     deps,
			When:          "$inputs.indicators != []",
			ForEach:       "$inputs.indicators",
			Wait:          "$signals.local_loop_ready",
			Timeout:       &timeout,
			ParallelGroup: "steps_" + stepID + "_group",
		},
		StructuralFields: uws1.StructuralFields{
			Items:     "$inputs.indicators",
			Mode:      "parallel",
			BatchSize: "2",
		},
		Steps: []*uws1.Step{
			operationStep("step_local_loop_fileio", "run_fileio_primary", []string{"step_collect_context", "step_parallel_checks"}),
			operationStep("step_local_loop_s3", "run_s3_primary", []string{"step_collect_context", "step_parallel_checks"}),
		},
		Outputs: map[string]string{
			"evidence": "$steps.step_local_loop_fileio.outputs.result",
			"archive":  "$steps.step_local_loop_s3.outputs.audit",
		},
		Extensions: ext("step", stepID),
	}
}

func caseBranch(name, when, opA, opB string) *uws1.Case {
	return &uws1.Case{
		CaseFields: uws1.CaseFields{Name: name, When: when},
		Body: map[string]any{
			"decision": map[string]any{"name": name, "when": when},
			"policy":   map[string]any{"requiresApproval": true, "notifyCustomer": true},
		},
		Steps: []*uws1.Step{
			operationStep("step_"+name+"_a", opA, []string{"run_llm_primary", "run_sql_primary"}),
			operationStep("step_"+name+"_b", opB, []string{"run_llm_primary", "run_sql_primary"}),
		},
		Extensions: ext("case", name),
	}
}

func fullSchema(name string) *uws1.ParamSchema {
	return &uws1.ParamSchema{
		Type:   "object",
		Format: "uws-" + name,
		Ref:    "#/components/schemas/" + name,
		Properties: map[string]*uws1.ParamSchema{
			"incidentId": {Type: "string", Format: "uuid", Extensions: ext("schema-property", name+"-incident")},
			"severity":   {Type: "string", Format: "enum", Extensions: ext("schema-property", name+"-severity")},
		},
		Required: []string{"incidentId", "severity"},
		Items: &uws1.ParamSchema{
			Type:       "string",
			Format:     "uuid",
			Extensions: ext("schema-items", name),
		},
		AllOf: []*uws1.ParamSchema{
			{Type: "object", Properties: map[string]*uws1.ParamSchema{"tenantId": {Type: "string"}, "region": {Type: "string"}}, Required: []string{"tenantId", "region"}, Extensions: ext("schema-allOf", name+"-tenant")},
			{Type: "object", Properties: map[string]*uws1.ParamSchema{"source": {Type: "string"}, "priority": {Type: "integer", Format: "int32"}}, Required: []string{"source", "priority"}, Extensions: ext("schema-allOf", name+"-source")},
		},
		OneOf: []*uws1.ParamSchema{
			{Type: "object", Properties: map[string]*uws1.ParamSchema{"host": {Type: "string"}, "ip": {Type: "string"}}, Required: []string{"host", "ip"}, Extensions: ext("schema-oneOf", name+"-host")},
			{Type: "object", Properties: map[string]*uws1.ParamSchema{"user": {Type: "string"}, "email": {Type: "string", Format: "email"}}, Required: []string{"user", "email"}, Extensions: ext("schema-oneOf", name+"-user")},
		},
		AnyOf: []*uws1.ParamSchema{
			{Type: "object", Properties: map[string]*uws1.ParamSchema{"ticket": {Type: "string"}, "caseId": {Type: "string"}}, Required: []string{"ticket", "caseId"}, Extensions: ext("schema-anyOf", name+"-ticket")},
			{Type: "object", Properties: map[string]*uws1.ParamSchema{"asset": {Type: "string"}, "owner": {Type: "string"}}, Required: []string{"asset", "owner"}, Extensions: ext("schema-anyOf", name+"-asset")},
		},
		Extensions: ext("schema", name),
	}
}

func criteria(condition string, typ uws1.CriterionExpressionType, context string) *uws1.Criterion {
	return &uws1.Criterion{
		Condition: condition,
		Type:      typ,
		Context:   context,
		Extensions: map[string]any{
			"x-evidence": "runtime-observed",
			"x-owner":    "quality-gate",
		},
	}
}

func failureRetry(name string) *uws1.FailureAction {
	return &uws1.FailureAction{
		Name:       name,
		Type:       "retry",
		RetryAfter: 5,
		RetryLimit: 2,
		Criteria: []*uws1.Criterion{
			criteria("$error.transient == true", uws1.CriterionSimple, ""),
			criteria("timeout|rate limit", uws1.CriterionRegex, "$error.message"),
		},
		Extensions: ext("failure-action", name),
	}
}

func failureGoto(name, workflowID string) *uws1.FailureAction {
	return &uws1.FailureAction{
		Name:       name,
		Type:       "goto",
		WorkflowID: workflowID,
		StepID:     "",
		RetryAfter: 0,
		RetryLimit: 0,
		Criteria: []*uws1.Criterion{
			criteria("$error.transient != true", uws1.CriterionSimple, ""),
			criteria("$.risk", uws1.CriterionJSONPath, "$error.details"),
		},
		Extensions: ext("failure-action", name),
	}
}

func successGoto(name, workflowID string) *uws1.SuccessAction {
	return &uws1.SuccessAction{
		Name:       name,
		Type:       "goto",
		WorkflowID: workflowID,
		Criteria: []*uws1.Criterion{
			criteria("$response.statusCode < 300", uws1.CriterionSimple, ""),
			criteria("$.ok", uws1.CriterionJSONPath, "$response.body"),
		},
		Extensions: ext("success-action", name),
	}
}

func successEnd(name string) *uws1.SuccessAction {
	return &uws1.SuccessAction{
		Name: name,
		Type: "end",
		Criteria: []*uws1.Criterion{
			criteria("$response.statusCode < 300", uws1.CriterionSimple, ""),
			criteria("accepted|complete", uws1.CriterionRegex, "$response.body.status"),
		},
		Extensions: ext("success-action", name),
	}
}

func ext(kind, name string) map[string]any {
	return map[string]any{
		"x-kind": kind,
		"x-name": name,
	}
}
