  uws = "1.1.0"
  info {
    title       = "Incident Response Automation"
    summary     = "End-to-end security incident enrichment, containment, notification, and evidence packaging."
    description = "A large fixture built from Go structs first. It exercises OpenAPI-bound operations, all public runtime supplement selectors, structural workflows, triggers, results, components, and extension preservation across JSON, YAML, and HCL."
    version     = "2026.05.08"
    extensions {
      x-owner = "security-automation"
      x-slo-minutes = 15
    }
  }
  sourceDescription "incident_api" {
    url  = "./openapi/incident-api.yaml"
    type = "openapi"
    extensions {
      x-service-tier = "gold"
      x-contact = "secops-api@example.test"
    }
  }
  sourceDescription "crm_api" {
    url  = "./openapi/crm-api.yaml"
    type = "openapi"
    extensions {
      x-contact = "customer-ops@example.test"
      x-service-tier = "silver"
    }
  }
  variables {
    severityPolicy {
      critical = "page"
      high = "ticket"
    }
    dryRun = false
    environment = "production"
    regions = [
      "us-east-1",
      "eu-west-1"
    ]
  }
  operation "fetch_ticket" {
    sourceDescription  = "incident_api"
    openapiOperationId = "getIncident"
    description        = "Fetch the incident record and current triage state."
    when               = "$inputs.incidentId != ''"
    forEach            = "$variables.regions"
    wait               = "$signals.incident_api_ready"
    parallelGroup      = "api_fetch_group"
    outputs = {
      severity = "$response.body.severity"
      ticket   = "$response.body"
    }
    request {
      path {
        incidentId = "$inputs.incidentId"
        tenantId = "$inputs.tenantId"
      }
      query {
        depth = "full"
        include = [
          "timeline",
          "assets"
        ]
      }
    }
    timeout = 20
    successCriterion {
      condition = "$response.statusCode == 200"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.severity"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_fetch_ticket" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_fetch_ticket"
      }
    }
    onFailure "manual_fetch_review" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "manual_fetch_review"
      }
    }
    onSuccess "continue_to_parallel" {
      type       = "goto"
      workflowId = "wf_parallel"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "continue_to_parallel"
      }
    }
    onSuccess "ticket_loaded" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "ticket_loaded"
      }
    }
    extensions {
      x-kind = "operation"
      x-name = "fetch-ticket"
    }
  }
  operation "load_customer" {
    sourceDescription   = "crm_api"
    openapiOperationRef = "#/paths/~1customers~1{customerId}/get"
    description         = "Load customer account context used for notification decisions."
    when                = "$steps.step_collect_context.outputs.customerId != ''"
    forEach             = "$variables.regions"
    wait                = "$signals.crm_ready"
    parallelGroup       = "api_fetch_group"
    outputs = {
      customer = "$response.body"
      segment  = "$response.body.segment"
    }
    request {
      path {
        customerId = "$steps.step_collect_context.outputs.customerId"
        tenantId = "$inputs.tenantId"
      }
      header {
        X-Workflow = "incident-response"
        X-Trace-ID = "$context.traceId"
      }
    }
    timeout = 25
    successCriterion {
      condition = "$response.statusCode == 200"
      type      = "simple"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    successCriterion {
      condition = "enterprise|strategic"
      type      = "regex"
      context   = "$response.body.segment"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_load_customer" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_load_customer"
      }
    }
    onFailure "manual_customer_review" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "manual_customer_review"
      }
    }
    onSuccess "continue_to_notification" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "continue_to_notification"
      }
    }
    onSuccess "customer_loaded" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "customer_loaded"
      }
    }
    extensions {
      x-name = "load-customer"
      x-kind = "operation"
    }
  }
  operation "run_ssh_primary" {
    description   = "Run ssh primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "ssh"
        variant = "primary"
      }
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "ssh"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_ssh_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "retry_run_ssh_primary"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_ssh_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_ssh_primary"
      }
    }
    onSuccess "notify_run_ssh_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_ssh_primary"
      }
    }
    onSuccess "complete_run_ssh_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_ssh_primary"
      }
    }
    extensions {
      x-uws-runtime {
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "ssh"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "sudo systemctl status incident-agent --no-pager"
        function = "execute_ssh_primary"
        type = "ssh"
        workflow = "runtime/ssh-primary.uws.hcl"
        workingDir = "/srv/incident-response/ssh"
      }
      x-name = "run_ssh_primary"
      x-kind = "runtime-operation"
      x-uws-operation-profile = "uws.runtime.1.0"
    }
  }
  operation "run_ssh_fallback" {
    description   = "Run ssh fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "ssh"
        variant = "fallback"
      }
      header {
        X-Runtime-Type = "ssh"
        X-Run-Variant = "fallback"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_ssh_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_ssh_fallback"
      }
    }
    onFailure "review_run_ssh_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "review_run_ssh_fallback"
        x-kind = "failure-action"
      }
    }
    onSuccess "notify_run_ssh_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_ssh_fallback"
      }
    }
    onSuccess "complete_run_ssh_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "complete_run_ssh_fallback"
        x-kind = "success-action"
      }
    }
    extensions {
      x-uws-runtime {
        workflow = "runtime/ssh-fallback.uws.hcl"
        workingDir = "/srv/incident-response/ssh"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "ssh"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "sudo systemctl status incident-agent --no-pager"
        function = "execute_ssh_fallback"
        type = "ssh"
      }
      x-name = "run_ssh_fallback"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-kind = "runtime-operation"
    }
  }
  operation "run_cmd_primary" {
    description   = "Run cmd primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "cmd"
        variant = "primary"
      }
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "cmd"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_cmd_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_cmd_primary"
      }
    }
    onFailure "review_run_cmd_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_cmd_primary"
      }
    }
    onSuccess "notify_run_cmd_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_cmd_primary"
      }
    }
    onSuccess "complete_run_cmd_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "complete_run_cmd_primary"
        x-kind = "success-action"
      }
    }
    extensions {
      x-uws-operation-profile = "uws.runtime.1.0"
      x-name = "run_cmd_primary"
      x-kind = "runtime-operation"
      x-uws-runtime {
        arguments = [
          {
            runtime = "cmd"
            incidentId = "$inputs.incidentId"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "incidentctl preflight --format=json --variant=primary"
        function = "execute_cmd_primary"
        type = "cmd"
        workflow = "runtime/cmd-primary.uws.hcl"
        workingDir = "/srv/incident-response/cmd"
      }
    }
  }
  operation "run_cmd_fallback" {
    description   = "Run cmd fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "cmd"
        variant = "fallback"
      }
      header {
        X-Runtime-Type = "cmd"
        X-Run-Variant = "fallback"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_cmd_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_cmd_fallback"
      }
    }
    onFailure "review_run_cmd_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_cmd_fallback"
      }
    }
    onSuccess "notify_run_cmd_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "notify_run_cmd_fallback"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_cmd_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "complete_run_cmd_fallback"
        x-kind = "success-action"
      }
    }
    extensions {
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        workingDir = "/srv/incident-response/cmd"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "cmd"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "incidentctl preflight --format=json --variant=fallback"
        function = "execute_cmd_fallback"
        type = "cmd"
        workflow = "runtime/cmd-fallback.uws.hcl"
      }
      x-kind = "runtime-operation"
      x-name = "run_cmd_fallback"
    }
  }
  operation "run_fnct_primary" {
    description   = "Run fnct primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "fnct"
      }
      body {
        incidentId = "$inputs.incidentId"
        runtime = "fnct"
        variant = "primary"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_fnct_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_fnct_primary"
      }
    }
    onFailure "review_run_fnct_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_fnct_primary"
      }
    }
    onSuccess "notify_run_fnct_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "notify_run_fnct_primary"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_fnct_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_fnct_primary"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-uws-runtime {
        command = "fnct task primary"
        function = "render_incident_brief_primary"
        type = "fnct"
        workflow = "runtime/fnct-primary.uws.hcl"
        workingDir = "/srv/incident-response/fnct"
        arguments = [
          {
            runtime = "fnct"
            incidentId = "$inputs.incidentId"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
      }
      x-name = "run_fnct_primary"
      x-uws-operation-profile = "uws.runtime.1.0"
    }
  }
  operation "run_fnct_fallback" {
    description   = "Run fnct fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "fnct"
        variant = "fallback"
      }
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "fnct"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_fnct_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "retry_run_fnct_fallback"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_fnct_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "review_run_fnct_fallback"
        x-kind = "failure-action"
      }
    }
    onSuccess "notify_run_fnct_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_fnct_fallback"
      }
    }
    onSuccess "complete_run_fnct_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_fnct_fallback"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        type = "fnct"
        workflow = "runtime/fnct-fallback.uws.hcl"
        workingDir = "/srv/incident-response/fnct"
        arguments = [
          {
            runtime = "fnct"
            incidentId = "$inputs.incidentId"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "fnct task fallback"
        function = "render_incident_brief_fallback"
      }
      x-name = "run_fnct_fallback"
    }
  }
  operation "run_fileio_primary" {
    description   = "Run fileio primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        runtime = "fileio"
        variant = "primary"
        incidentId = "$inputs.incidentId"
      }
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "fileio"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_fileio_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "retry_run_fileio_primary"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_fileio_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_fileio_primary"
      }
    }
    onSuccess "notify_run_fileio_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_fileio_primary"
      }
    }
    onSuccess "complete_run_fileio_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "complete_run_fileio_primary"
        x-kind = "success-action"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-name = "run_fileio_primary"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        type = "fileio"
        workflow = "runtime/fileio-primary.uws.hcl"
        workingDir = "/srv/incident-response/fileio"
        arguments = [
          {
            runtime = "fileio"
            incidentId = "$inputs.incidentId"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "fileio task primary"
        function = "write_evidence_bundle_primary"
      }
    }
  }
  operation "run_fileio_fallback" {
    description   = "Run fileio fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        runtime = "fileio"
        variant = "fallback"
        incidentId = "$inputs.incidentId"
      }
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "fileio"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_fileio_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_fileio_fallback"
      }
    }
    onFailure "review_run_fileio_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_fileio_fallback"
      }
    }
    onSuccess "notify_run_fileio_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_fileio_fallback"
      }
    }
    onSuccess "complete_run_fileio_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "complete_run_fileio_fallback"
        x-kind = "success-action"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-name = "run_fileio_fallback"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "fileio"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "fileio task fallback"
        function = "write_evidence_bundle_fallback"
        type = "fileio"
        workflow = "runtime/fileio-fallback.uws.hcl"
        workingDir = "/srv/incident-response/fileio"
      }
    }
  }
  operation "run_sql_primary" {
    description   = "Run sql primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        variant = "primary"
        incidentId = "$inputs.incidentId"
        runtime = "sql"
      }
      header {
        X-Runtime-Type = "sql"
        X-Run-Variant = "primary"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_sql_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "retry_run_sql_primary"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_sql_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "review_run_sql_primary"
        x-kind = "failure-action"
      }
    }
    onSuccess "notify_run_sql_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_sql_primary"
      }
    }
    onSuccess "complete_run_sql_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_sql_primary"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-name = "run_sql_primary"
      x-uws-runtime {
        function = "execute_sql_primary"
        type = "sql"
        workflow = "runtime/sql-primary.uws.hcl"
        workingDir = "/srv/incident-response/sql"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "sql"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "SELECT id, severity FROM incidents WHERE id = :incident_id"
      }
    }
  }
  operation "run_sql_fallback" {
    description   = "Run sql fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "sql"
        variant = "fallback"
      }
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "sql"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_sql_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_sql_fallback"
      }
    }
    onFailure "review_run_sql_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_sql_fallback"
      }
    }
    onSuccess "notify_run_sql_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "notify_run_sql_fallback"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_sql_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_sql_fallback"
      }
    }
    extensions {
      x-uws-runtime {
        workingDir = "/srv/incident-response/sql"
        arguments = [
          {
            runtime = "sql"
            incidentId = "$inputs.incidentId"
          },
          {
            variant = "fallback"
            region = "$context.region"
          }
        ]
        command = "SELECT id, severity FROM incidents WHERE id = :incident_id"
        function = "execute_sql_fallback"
        type = "sql"
        workflow = "runtime/sql-fallback.uws.hcl"
      }
      x-kind = "runtime-operation"
      x-name = "run_sql_fallback"
      x-uws-operation-profile = "uws.runtime.1.0"
    }
  }
  operation "run_s3_primary" {
    description   = "Run s3 primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Runtime-Type = "s3"
        X-Run-Variant = "primary"
      }
      body {
        incidentId = "$inputs.incidentId"
        runtime = "s3"
        variant = "primary"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_s3_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_s3_primary"
      }
    }
    onFailure "review_run_s3_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "review_run_s3_primary"
        x-kind = "failure-action"
      }
    }
    onSuccess "notify_run_s3_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "notify_run_s3_primary"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_s3_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_s3_primary"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-name = "run_s3_primary"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        type = "s3"
        workflow = "runtime/s3-primary.uws.hcl"
        workingDir = "/srv/incident-response/s3"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "s3"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "s3 task primary"
        function = "archive_evidence_s3_primary"
      }
    }
  }
  operation "run_s3_fallback" {
    description   = "Run s3 fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Runtime-Type = "s3"
        X-Run-Variant = "fallback"
      }
      body {
        incidentId = "$inputs.incidentId"
        runtime = "s3"
        variant = "fallback"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_s3_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "retry_run_s3_fallback"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_s3_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "review_run_s3_fallback"
        x-kind = "failure-action"
      }
    }
    onSuccess "notify_run_s3_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_s3_fallback"
      }
    }
    onSuccess "complete_run_s3_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_s3_fallback"
      }
    }
    extensions {
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        workflow = "runtime/s3-fallback.uws.hcl"
        workingDir = "/srv/incident-response/s3"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "s3"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "s3 task fallback"
        function = "archive_evidence_s3_fallback"
        type = "s3"
      }
      x-kind = "runtime-operation"
      x-name = "run_s3_fallback"
    }
  }
  operation "run_smtp_primary" {
    description   = "Run smtp primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        runtime = "smtp"
        variant = "primary"
        incidentId = "$inputs.incidentId"
      }
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "smtp"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_smtp_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_smtp_primary"
      }
    }
    onFailure "review_run_smtp_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_smtp_primary"
      }
    }
    onSuccess "notify_run_smtp_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_smtp_primary"
      }
    }
    onSuccess "complete_run_smtp_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "complete_run_smtp_primary"
        x-kind = "success-action"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-name = "run_smtp_primary"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        command = "smtp task primary"
        function = "send_customer_notice_primary"
        type = "smtp"
        workflow = "runtime/smtp-primary.uws.hcl"
        workingDir = "/srv/incident-response/smtp"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "smtp"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
      }
    }
  }
  operation "run_smtp_fallback" {
    description   = "Run smtp fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        variant = "fallback"
        incidentId = "$inputs.incidentId"
        runtime = "smtp"
      }
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "smtp"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_smtp_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "retry_run_smtp_fallback"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_smtp_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_smtp_fallback"
      }
    }
    onSuccess "notify_run_smtp_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "notify_run_smtp_fallback"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_smtp_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "complete_run_smtp_fallback"
        x-kind = "success-action"
      }
    }
    extensions {
      x-uws-runtime {
        type = "smtp"
        workflow = "runtime/smtp-fallback.uws.hcl"
        workingDir = "/srv/incident-response/smtp"
        arguments = [
          {
            runtime = "smtp"
            incidentId = "$inputs.incidentId"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "smtp task fallback"
        function = "send_customer_notice_fallback"
      }
      x-kind = "runtime-operation"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-name = "run_smtp_fallback"
    }
  }
  operation "run_dns_primary" {
    description   = "Run dns primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "dns"
      }
      body {
        incidentId = "$inputs.incidentId"
        runtime = "dns"
        variant = "primary"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_dns_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_dns_primary"
      }
    }
    onFailure "review_run_dns_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_dns_primary"
      }
    }
    onSuccess "notify_run_dns_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "notify_run_dns_primary"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_dns_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_dns_primary"
      }
    }
    extensions {
      x-name = "run_dns_primary"
      x-kind = "runtime-operation"
      x-uws-runtime {
        workflow = "runtime/dns-primary.uws.hcl"
        workingDir = "/srv/incident-response/dns"
        arguments = [
          {
            runtime = "dns"
            incidentId = "$inputs.incidentId"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "dig +short suspicious.example.test"
        function = "execute_dns_primary"
        type = "dns"
      }
      x-uws-operation-profile = "uws.runtime.1.0"
    }
  }
  operation "run_dns_fallback" {
    description   = "Run dns fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "dns"
      }
      body {
        incidentId = "$inputs.incidentId"
        runtime = "dns"
        variant = "fallback"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_dns_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_dns_fallback"
      }
    }
    onFailure "review_run_dns_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_dns_fallback"
      }
    }
    onSuccess "notify_run_dns_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_dns_fallback"
      }
    }
    onSuccess "complete_run_dns_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "complete_run_dns_fallback"
        x-kind = "success-action"
      }
    }
    extensions {
      x-uws-operation-profile = "uws.runtime.1.0"
      x-kind = "runtime-operation"
      x-name = "run_dns_fallback"
      x-uws-runtime {
        workflow = "runtime/dns-fallback.uws.hcl"
        workingDir = "/srv/incident-response/dns"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "dns"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "dig +short suspicious.example.test"
        function = "execute_dns_fallback"
        type = "dns"
      }
    }
  }
  operation "run_ldaps_primary" {
    description   = "Run ldaps primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "ldaps"
      }
      body {
        variant = "primary"
        incidentId = "$inputs.incidentId"
        runtime = "ldaps"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_ldaps_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "retry_run_ldaps_primary"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_ldaps_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_ldaps_primary"
      }
    }
    onSuccess "notify_run_ldaps_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_ldaps_primary"
      }
    }
    onSuccess "complete_run_ldaps_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "complete_run_ldaps_primary"
        x-kind = "success-action"
      }
    }
    extensions {
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        type = "ldaps"
        workflow = "runtime/ldaps-primary.uws.hcl"
        workingDir = "/srv/incident-response/ldaps"
        arguments = [
          {
            runtime = "ldaps"
            incidentId = "$inputs.incidentId"
          },
          {
            variant = "primary"
            region = "$context.region"
          }
        ]
        command = "ldaps task primary"
        function = "lookup_owner_directory_primary"
      }
      x-kind = "runtime-operation"
      x-name = "run_ldaps_primary"
    }
  }
  operation "run_ldaps_fallback" {
    description   = "Run ldaps fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "ldaps"
      }
      body {
        variant = "fallback"
        incidentId = "$inputs.incidentId"
        runtime = "ldaps"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_ldaps_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_ldaps_fallback"
      }
    }
    onFailure "review_run_ldaps_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_ldaps_fallback"
      }
    }
    onSuccess "notify_run_ldaps_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_ldaps_fallback"
      }
    }
    onSuccess "complete_run_ldaps_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_ldaps_fallback"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-name = "run_ldaps_fallback"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        workflow = "runtime/ldaps-fallback.uws.hcl"
        workingDir = "/srv/incident-response/ldaps"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "ldaps"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "ldaps task fallback"
        function = "lookup_owner_directory_fallback"
        type = "ldaps"
      }
    }
  }
  operation "run_scp_primary" {
    description   = "Run scp primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "scp"
      }
      body {
        incidentId = "$inputs.incidentId"
        runtime = "scp"
        variant = "primary"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_scp_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "retry_run_scp_primary"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_scp_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_scp_primary"
      }
    }
    onSuccess "notify_run_scp_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "notify_run_scp_primary"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_scp_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_scp_primary"
      }
    }
    extensions {
      x-name = "run_scp_primary"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-kind = "runtime-operation"
      x-uws-runtime {
        command = "scp evidence.tar.gz evidence-vault:/incoming/primary"
        function = "execute_scp_primary"
        type = "scp"
        workflow = "runtime/scp-primary.uws.hcl"
        workingDir = "/srv/incident-response/scp"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "scp"
          },
          {
            variant = "primary"
            region = "$context.region"
          }
        ]
      }
    }
  }
  operation "run_scp_fallback" {
    description   = "Run scp fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "scp"
        variant = "fallback"
      }
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "scp"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_scp_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_scp_fallback"
      }
    }
    onFailure "review_run_scp_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_scp_fallback"
      }
    }
    onSuccess "notify_run_scp_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_scp_fallback"
      }
    }
    onSuccess "complete_run_scp_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "complete_run_scp_fallback"
        x-kind = "success-action"
      }
    }
    extensions {
      x-uws-operation-profile = "uws.runtime.1.0"
      x-uws-runtime {
        function = "execute_scp_fallback"
        type = "scp"
        workflow = "runtime/scp-fallback.uws.hcl"
        workingDir = "/srv/incident-response/scp"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "scp"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "scp evidence.tar.gz evidence-vault:/incoming/fallback"
      }
      x-name = "run_scp_fallback"
      x-kind = "runtime-operation"
    }
  }
  operation "run_sftp_primary" {
    description   = "Run sftp primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        incidentId = "$inputs.incidentId"
        runtime = "sftp"
        variant = "primary"
      }
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "sftp"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_sftp_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_sftp_primary"
      }
    }
    onFailure "review_run_sftp_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "review_run_sftp_primary"
        x-kind = "failure-action"
      }
    }
    onSuccess "notify_run_sftp_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_sftp_primary"
      }
    }
    onSuccess "complete_run_sftp_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "complete_run_sftp_primary"
        x-kind = "success-action"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-name = "run_sftp_primary"
      x-uws-runtime {
        arguments = [
          {
            runtime = "sftp"
            incidentId = "$inputs.incidentId"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "put evidence.json /incoming/primary/evidence.json"
        function = "execute_sftp_primary"
        type = "sftp"
        workflow = "runtime/sftp-primary.uws.hcl"
        workingDir = "/srv/incident-response/sftp"
      }
      x-uws-operation-profile = "uws.runtime.1.0"
    }
  }
  operation "run_sftp_fallback" {
    description   = "Run sftp fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        variant = "fallback"
        incidentId = "$inputs.incidentId"
        runtime = "sftp"
      }
      header {
        X-Runtime-Type = "sftp"
        X-Run-Variant = "fallback"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_sftp_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_sftp_fallback"
      }
    }
    onFailure "review_run_sftp_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_sftp_fallback"
      }
    }
    onSuccess "notify_run_sftp_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "notify_run_sftp_fallback"
        x-kind = "success-action"
      }
    }
    onSuccess "complete_run_sftp_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_sftp_fallback"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-uws-runtime {
        type = "sftp"
        workflow = "runtime/sftp-fallback.uws.hcl"
        workingDir = "/srv/incident-response/sftp"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "sftp"
          },
          {
            region = "$context.region"
            variant = "fallback"
          }
        ]
        command = "put evidence.json /incoming/fallback/evidence.json"
        function = "execute_sftp_fallback"
      }
      x-name = "run_sftp_fallback"
      x-uws-operation-profile = "uws.runtime.1.0"
    }
  }
  operation "run_llm_primary" {
    description   = "Run llm primary runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_primary_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      body {
        runtime = "llm"
        variant = "primary"
        incidentId = "$inputs.incidentId"
      }
      header {
        X-Run-Variant = "primary"
        X-Runtime-Type = "llm"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    onFailure "retry_run_llm_primary" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "retry_run_llm_primary"
      }
    }
    onFailure "review_run_llm_primary" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-name = "review_run_llm_primary"
        x-kind = "failure-action"
      }
    }
    onSuccess "notify_run_llm_primary" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_llm_primary"
      }
    }
    onSuccess "complete_run_llm_primary" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_llm_primary"
      }
    }
    extensions {
      x-uws-runtime {
        type = "llm"
        workflow = "runtime/llm-primary.uws.hcl"
        workingDir = "/srv/incident-response/llm"
        arguments = [
          {
            incidentId = "$inputs.incidentId"
            runtime = "llm"
          },
          {
            region = "$context.region"
            variant = "primary"
          }
        ]
        command = "llm task primary"
        function = "summarize_incident_primary"
      }
      x-kind = "runtime-operation"
      x-name = "run_llm_primary"
      x-uws-operation-profile = "uws.runtime.1.0"
    }
  }
  operation "run_llm_fallback" {
    description   = "Run llm fallback runtime task for incident response."
    dependsOn     = ["fetch_ticket", "load_customer"]
    when          = "$steps.step_collect_context.outputs.enabled == true"
    forEach       = "$variables.regions"
    wait          = "$signals.runtime_slot_available"
    parallelGroup = "runtime_fallback_group"
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
    request {
      header {
        X-Run-Variant = "fallback"
        X-Runtime-Type = "llm"
      }
      body {
        variant = "fallback"
        incidentId = "$inputs.incidentId"
        runtime = "llm"
      }
    }
    timeout = 30
    successCriterion {
      condition = "$response.statusCode < 500"
      type      = "simple"
      extensions {
        x-evidence = "runtime-observed"
        x-owner = "quality-gate"
      }
    }
    successCriterion {
      condition = "$.ok"
      type      = "jsonpath"
      context   = "$response.body"
      extensions {
        x-owner = "quality-gate"
        x-evidence = "runtime-observed"
      }
    }
    onFailure "retry_run_llm_fallback" {
      type       = "retry"
      retryAfter = 5
      retryLimit = 2
      criterion {
        condition = "$error.transient == true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "timeout|rate limit"
        type      = "regex"
        context   = "$error.message"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-name = "retry_run_llm_fallback"
        x-kind = "failure-action"
      }
    }
    onFailure "review_run_llm_fallback" {
      type       = "goto"
      workflowId = "wf_manual_review"
      criterion {
        condition = "$error.transient != true"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.risk"
        type      = "jsonpath"
        context   = "$error.details"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "failure-action"
        x-name = "review_run_llm_fallback"
      }
    }
    onSuccess "notify_run_llm_fallback" {
      type       = "goto"
      workflowId = "wf_notify_operators"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      criterion {
        condition = "$.ok"
        type      = "jsonpath"
        context   = "$response.body"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "notify_run_llm_fallback"
      }
    }
    onSuccess "complete_run_llm_fallback" {
      type = "end"
      criterion {
        condition = "$response.statusCode < 300"
        type      = "simple"
        extensions {
          x-owner = "quality-gate"
          x-evidence = "runtime-observed"
        }
      }
      criterion {
        condition = "accepted|complete"
        type      = "regex"
        context   = "$response.body.status"
        extensions {
          x-evidence = "runtime-observed"
          x-owner = "quality-gate"
        }
      }
      extensions {
        x-kind = "success-action"
        x-name = "complete_run_llm_fallback"
      }
    }
    extensions {
      x-kind = "runtime-operation"
      x-uws-operation-profile = "uws.runtime.1.0"
      x-name = "run_llm_fallback"
      x-uws-runtime {
        workingDir = "/srv/incident-response/llm"
        arguments = [
          {
            runtime = "llm"
            incidentId = "$inputs.incidentId"
          },
          {
            variant = "fallback"
            region = "$context.region"
          }
        ]
        command = "llm task fallback"
        function = "summarize_incident_fallback"
        type = "llm"
        workflow = "runtime/llm-fallback.uws.hcl"
      }
    }
  }
  workflow "main" {
    type        = "sequence"
    description = "Coordinate enrichment, runtime checks, branching, containment, and notification."
    dependsOn   = ["fetch_ticket", "load_customer"]
    when        = "$inputs.incidentId != ''"
    forEach     = "$variables.regions"
    wait        = "$signals.start"
    outputs = {
      decision = "$steps.step_decide_path.outputs.selectedPath"
      incident = "$steps.step_collect_context.outputs.incident"
    }
    inputs {
      type     = "object"
      format   = "uws-incidentInput"
      _ref     = "#/components/schemas/incidentInput"
      required = ["incidentId", "severity"]
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-property"
          x-name = "incidentInput-incident"
        }
      }
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-kind = "schema-property"
          x-name = "incidentInput-severity"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "incidentInput"
          x-kind = "schema-items"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "incidentInput-tenant"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        properties "source" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "incidentInput-source"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "host" {
          type = "string"
        }
        properties "ip" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "incidentInput-host"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "incidentInput-user"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "caseId" {
          type = "string"
        }
        properties "ticket" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "incidentInput-ticket"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "asset" {
          type = "string"
        }
        properties "owner" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "incidentInput-asset"
        }
      }
      extensions {
        x-name = "incidentInput"
        x-kind = "schema"
      }
    }
    idempotency {
      key        = "$inputs.tenantId + ':' + $inputs.incidentId"
      onConflict = "returnPrevious"
      ttl = 900
      extensions {
        x-kind = "idempotency"
        x-name = "main"
      }
    }
    timeout = 120
    step "step_collect_context" {
      description   = "Execute fetch_ticket and expose normalized outputs."
      operationRef  = "fetch_ticket"
      dependsOn     = ["run_cmd_primary", "run_fnct_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_collect_context_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_collect_context"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-name = "step_collect_context"
        x-kind = "step"
      }
    }
    step "step_parallel_checks" {
      description   = "Invoke workflow wf_parallel."
      dependsOn     = ["step_collect_context", "load_customer"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.workflow_step_ready"
      workflow      = "wf_parallel"
      parallelGroup = "steps_step_parallel_checks_group"
      outputs = {
        audit  = "$workflow.outputs.audit"
        result = "$workflow.outputs.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          workflow = "wf_parallel"
        }
        meta {
          stepId = "step_parallel_checks"
          requestedBy = "$context.actor"
        }
      }
      timeout = 40
      extensions {
        x-kind = "step"
        x-name = "step_parallel_checks"
      }
    }
    step "step_local_switch" {
      type          = "switch"
      description   = "Choose a local fast path before invoking the top-level switch workflow."
      dependsOn     = ["step_collect_context", "step_parallel_checks"]
      when          = "$inputs.severity != 'low'"
      forEach       = "$variables.regions"
      wait          = "$signals.local_switch_ready"
      parallelGroup = "steps_step_local_switch_group"
      outputs = {
        audit    = "$selected.steps[0].outputs.audit"
        selected = "$selected.case.name"
      }
      body {
        input {
          step = "step_local_switch"
          incidentId = "$inputs.incidentId"
        }
        meta {
          owner = "secops"
          decision = "local"
        }
      }
      timeout = 50
      case "local_switch_critical" {
        when = "$inputs.severity == 'critical'"
        body {
          decision {
            when = "$inputs.severity == 'critical'"
            name = "local_switch_critical"
          }
          policy {
            requiresApproval = true
            notifyCustomer = true
          }
        }
        step "step_local_switch_critical_a" {
          description   = "Execute run_cmd_primary and expose normalized outputs."
          operationRef  = "run_cmd_primary"
          dependsOn     = ["run_llm_primary", "run_sql_primary"]
          when          = "$inputs.enabled != false"
          forEach       = "$variables.regions"
          wait          = "$signals.step_ready"
          parallelGroup = "steps_step_local_switch_critical_a_group"
          outputs = {
            audit  = "$response.body.auditId"
            result = "$response.body.result"
          }
          body {
            input {
              incidentId = "$inputs.incidentId"
              source = "step_local_switch_critical_a"
            }
            meta {
              attempt = 1
              correlationId = "$context.correlationId"
            }
          }
          timeout = 30
          extensions {
            x-kind = "step"
            x-name = "step_local_switch_critical_a"
          }
        }
        step "step_local_switch_critical_b" {
          description   = "Execute run_fileio_primary and expose normalized outputs."
          operationRef  = "run_fileio_primary"
          dependsOn     = ["run_llm_primary", "run_sql_primary"]
          when          = "$inputs.enabled != false"
          forEach       = "$variables.regions"
          wait          = "$signals.step_ready"
          parallelGroup = "steps_step_local_switch_critical_b_group"
          outputs = {
            audit  = "$response.body.auditId"
            result = "$response.body.result"
          }
          body {
            input {
              incidentId = "$inputs.incidentId"
              source = "step_local_switch_critical_b"
            }
            meta {
              attempt = 1
              correlationId = "$context.correlationId"
            }
          }
          timeout = 30
          extensions {
            x-kind = "step"
            x-name = "step_local_switch_critical_b"
          }
        }
        extensions {
          x-kind = "case"
          x-name = "local_switch_critical"
        }
      }
      case "local_switch_customer" {
        when = "$steps.step_collect_context.outputs.segment == 'strategic'"
        body {
          decision {
            name = "local_switch_customer"
            when = "$steps.step_collect_context.outputs.segment == 'strategic'"
          }
          policy {
            notifyCustomer = true
            requiresApproval = true
          }
        }
        step "step_local_switch_customer_a" {
          description   = "Execute run_smtp_primary and expose normalized outputs."
          operationRef  = "run_smtp_primary"
          dependsOn     = ["run_llm_primary", "run_sql_primary"]
          when          = "$inputs.enabled != false"
          forEach       = "$variables.regions"
          wait          = "$signals.step_ready"
          parallelGroup = "steps_step_local_switch_customer_a_group"
          outputs = {
            audit  = "$response.body.auditId"
            result = "$response.body.result"
          }
          body {
            input {
              incidentId = "$inputs.incidentId"
              source = "step_local_switch_customer_a"
            }
            meta {
              correlationId = "$context.correlationId"
              attempt = 1
            }
          }
          timeout = 30
          extensions {
            x-kind = "step"
            x-name = "step_local_switch_customer_a"
          }
        }
        step "step_local_switch_customer_b" {
          description   = "Execute run_llm_primary and expose normalized outputs."
          operationRef  = "run_llm_primary"
          dependsOn     = ["run_llm_primary", "run_sql_primary"]
          when          = "$inputs.enabled != false"
          forEach       = "$variables.regions"
          wait          = "$signals.step_ready"
          parallelGroup = "steps_step_local_switch_customer_b_group"
          outputs = {
            audit  = "$response.body.auditId"
            result = "$response.body.result"
          }
          body {
            input {
              source = "step_local_switch_customer_b"
              incidentId = "$inputs.incidentId"
            }
            meta {
              attempt = 1
              correlationId = "$context.correlationId"
            }
          }
          timeout = 30
          extensions {
            x-kind = "step"
            x-name = "step_local_switch_customer_b"
          }
        }
        extensions {
          x-kind = "case"
          x-name = "local_switch_customer"
        }
      }
      default "step_local_switch_default_a" {
        description   = "Execute run_fnct_primary and expose normalized outputs."
        operationRef  = "run_fnct_primary"
        dependsOn     = ["step_collect_context", "step_parallel_checks"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_local_switch_default_a_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          meta {
            correlationId = "$context.correlationId"
            attempt = 1
          }
          input {
            incidentId = "$inputs.incidentId"
            source = "step_local_switch_default_a"
          }
        }
        timeout = 30
        extensions {
          x-kind = "step"
          x-name = "step_local_switch_default_a"
        }
      }
      default "step_local_switch_default_b" {
        description   = "Execute run_dns_primary and expose normalized outputs."
        operationRef  = "run_dns_primary"
        dependsOn     = ["step_collect_context", "step_parallel_checks"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_local_switch_default_b_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          meta {
            correlationId = "$context.correlationId"
            attempt = 1
          }
          input {
            incidentId = "$inputs.incidentId"
            source = "step_local_switch_default_b"
          }
        }
        timeout = 30
        extensions {
          x-kind = "step"
          x-name = "step_local_switch_default_b"
        }
      }
      extensions {
        x-name = "step_local_switch"
        x-kind = "step"
      }
    }
    step "step_local_loop" {
      type          = "loop"
      description   = "Loop over suspect indicators and write local evidence before remote upload."
      dependsOn     = ["step_collect_context", "step_parallel_checks"]
      when          = "$inputs.indicators != []"
      forEach       = "$inputs.indicators"
      wait          = "$signals.local_loop_ready"
      parallelGroup = "steps_step_local_loop_group"
      items         = "$inputs.indicators"
      mode          = "parallel"
      batchSize     = "2"
      outputs = {
        archive  = "$steps.step_local_loop_s3.outputs.audit"
        evidence = "$steps.step_local_loop_fileio.outputs.result"
      }
      body {
        input {
          step = "step_local_loop"
          indicators = "$inputs.indicators"
        }
        meta {
          batch = "local"
          owner = "forensics"
        }
      }
      timeout = 70
      step "step_local_loop_fileio" {
        description   = "Execute run_fileio_primary and expose normalized outputs."
        operationRef  = "run_fileio_primary"
        dependsOn     = ["step_collect_context", "step_parallel_checks"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_local_loop_fileio_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          meta {
            correlationId = "$context.correlationId"
            attempt = 1
          }
          input {
            incidentId = "$inputs.incidentId"
            source = "step_local_loop_fileio"
          }
        }
        timeout = 30
        extensions {
          x-kind = "step"
          x-name = "step_local_loop_fileio"
        }
      }
      step "step_local_loop_s3" {
        description   = "Execute run_s3_primary and expose normalized outputs."
        operationRef  = "run_s3_primary"
        dependsOn     = ["step_collect_context", "step_parallel_checks"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_local_loop_s3_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          input {
            incidentId = "$inputs.incidentId"
            source = "step_local_loop_s3"
          }
          meta {
            attempt = 1
            correlationId = "$context.correlationId"
          }
        }
        timeout = 30
        extensions {
          x-name = "step_local_loop_s3"
          x-kind = "step"
        }
      }
      extensions {
        x-kind = "step"
        x-name = "step_local_loop"
      }
    }
    step "step_decide_path" {
      description   = "Invoke workflow wf_switch."
      dependsOn     = ["step_local_switch", "step_local_loop"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.workflow_step_ready"
      workflow      = "wf_switch"
      parallelGroup = "steps_step_decide_path_group"
      outputs = {
        audit  = "$workflow.outputs.audit"
        result = "$workflow.outputs.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          workflow = "wf_switch"
        }
        meta {
          requestedBy = "$context.actor"
          stepId = "step_decide_path"
        }
      }
      timeout = 40
      extensions {
        x-name = "step_decide_path"
        x-kind = "step"
      }
    }
    step "step_notify" {
      description   = "Invoke workflow wf_notify_operators."
      dependsOn     = ["step_decide_path", "run_llm_fallback"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.workflow_step_ready"
      workflow      = "wf_notify_operators"
      parallelGroup = "steps_step_notify_group"
      outputs = {
        audit  = "$workflow.outputs.audit"
        result = "$workflow.outputs.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          workflow = "wf_notify_operators"
        }
        meta {
          requestedBy = "$context.actor"
          stepId = "step_notify"
        }
      }
      timeout = 40
      extensions {
        x-name = "step_notify"
        x-kind = "step"
      }
    }
    extensions {
      x-name = "main"
      x-kind = "workflow"
    }
  }
  workflow "wf_parallel" {
    type        = "parallel"
    description = "Run independent runtime checks in parallel."
    dependsOn   = ["fetch_ticket", "load_customer"]
    when        = "$inputs.severity in ['critical','high']"
    forEach     = "$variables.regions"
    wait        = "$signals.parallel_ready"
    outputs = {
      health = "$steps.step_ssh_primary.outputs.result"
      risk   = "$steps.step_llm_primary.outputs.result"
    }
    inputs {
      type     = "object"
      format   = "uws-parallelInput"
      _ref     = "#/components/schemas/parallelInput"
      required = ["incidentId", "severity"]
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-property"
          x-name = "parallelInput-incident"
        }
      }
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-kind = "schema-property"
          x-name = "parallelInput-severity"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-items"
          x-name = "parallelInput"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "parallelInput-tenant"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        properties "source" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "parallelInput-source"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "host" {
          type = "string"
        }
        properties "ip" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "parallelInput-host"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "parallelInput-user"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "ticket" {
          type = "string"
        }
        properties "caseId" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "parallelInput-ticket"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "asset" {
          type = "string"
        }
        properties "owner" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "parallelInput-asset"
        }
      }
      extensions {
        x-kind = "schema"
        x-name = "parallelInput"
      }
    }
    timeout = 90
    step "step_ssh_primary" {
      description   = "Execute run_ssh_primary and expose normalized outputs."
      operationRef  = "run_ssh_primary"
      dependsOn     = ["fetch_ticket", "load_customer"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_ssh_primary_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
        input {
          incidentId = "$inputs.incidentId"
          source = "step_ssh_primary"
        }
      }
      timeout = 30
      extensions {
        x-name = "step_ssh_primary"
        x-kind = "step"
      }
    }
    step "step_sql_primary" {
      description   = "Execute run_sql_primary and expose normalized outputs."
      operationRef  = "run_sql_primary"
      dependsOn     = ["fetch_ticket", "load_customer"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_sql_primary_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_sql_primary"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_sql_primary"
      }
    }
    step "step_dns_primary" {
      description   = "Execute run_dns_primary and expose normalized outputs."
      operationRef  = "run_dns_primary"
      dependsOn     = ["fetch_ticket", "load_customer"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_dns_primary_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_dns_primary"
        }
        meta {
          correlationId = "$context.correlationId"
          attempt = 1
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_dns_primary"
      }
    }
    step "step_llm_primary" {
      description   = "Execute run_llm_primary and expose normalized outputs."
      operationRef  = "run_llm_primary"
      dependsOn     = ["fetch_ticket", "load_customer"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_llm_primary_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_llm_primary"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_llm_primary"
      }
    }
    extensions {
      x-kind = "workflow"
      x-name = "parallel"
    }
  }
  workflow "wf_switch" {
    type        = "switch"
    description = "Choose containment path based on severity and customer tier."
    dependsOn   = ["run_llm_primary", "run_sql_primary"]
    when        = "$inputs.severity != 'low'"
    forEach     = "$variables.regions"
    wait        = "$signals.decision_ready"
    outputs = {
      nextAction   = "$selected.steps[0].outputs.result"
      selectedPath = "$selected.case.name"
    }
    inputs {
      type     = "object"
      format   = "uws-switchInput"
      _ref     = "#/components/schemas/switchInput"
      required = ["incidentId", "severity"]
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-property"
          x-name = "switchInput-incident"
        }
      }
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-kind = "schema-property"
          x-name = "switchInput-severity"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-items"
          x-name = "switchInput"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "switchInput-tenant"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "source" {
          type = "string"
        }
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        extensions {
          x-name = "switchInput-source"
          x-kind = "schema-allOf"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "host" {
          type = "string"
        }
        properties "ip" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "switchInput-host"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "switchInput-user"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "caseId" {
          type = "string"
        }
        properties "ticket" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "switchInput-ticket"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "asset" {
          type = "string"
        }
        properties "owner" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "switchInput-asset"
        }
      }
      extensions {
        x-kind = "schema"
        x-name = "switchInput"
      }
    }
    timeout = 45
    case "case_critical" {
      when = "$inputs.severity == 'critical'"
      body {
        decision {
          when = "$inputs.severity == 'critical'"
          name = "case_critical"
        }
        policy {
          requiresApproval = true
          notifyCustomer = true
        }
      }
      step "step_case_critical_a" {
        description   = "Execute run_smtp_primary and expose normalized outputs."
        operationRef  = "run_smtp_primary"
        dependsOn     = ["run_llm_primary", "run_sql_primary"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_case_critical_a_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          input {
            incidentId = "$inputs.incidentId"
            source = "step_case_critical_a"
          }
          meta {
            attempt = 1
            correlationId = "$context.correlationId"
          }
        }
        timeout = 30
        extensions {
          x-kind = "step"
          x-name = "step_case_critical_a"
        }
      }
      step "step_case_critical_b" {
        description   = "Execute run_s3_primary and expose normalized outputs."
        operationRef  = "run_s3_primary"
        dependsOn     = ["run_llm_primary", "run_sql_primary"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_case_critical_b_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          input {
            incidentId = "$inputs.incidentId"
            source = "step_case_critical_b"
          }
          meta {
            correlationId = "$context.correlationId"
            attempt = 1
          }
        }
        timeout = 30
        extensions {
          x-name = "step_case_critical_b"
          x-kind = "step"
        }
      }
      extensions {
        x-kind = "case"
        x-name = "case_critical"
      }
    }
    case "case_enterprise" {
      when = "$steps.step_collect_context.outputs.segment == 'enterprise'"
      body {
        decision {
          name = "case_enterprise"
          when = "$steps.step_collect_context.outputs.segment == 'enterprise'"
        }
        policy {
          requiresApproval = true
          notifyCustomer = true
        }
      }
      step "step_case_enterprise_a" {
        description   = "Execute run_ldaps_primary and expose normalized outputs."
        operationRef  = "run_ldaps_primary"
        dependsOn     = ["run_llm_primary", "run_sql_primary"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_case_enterprise_a_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          input {
            incidentId = "$inputs.incidentId"
            source = "step_case_enterprise_a"
          }
          meta {
            attempt = 1
            correlationId = "$context.correlationId"
          }
        }
        timeout = 30
        extensions {
          x-kind = "step"
          x-name = "step_case_enterprise_a"
        }
      }
      step "step_case_enterprise_b" {
        description   = "Execute run_llm_primary and expose normalized outputs."
        operationRef  = "run_llm_primary"
        dependsOn     = ["run_llm_primary", "run_sql_primary"]
        when          = "$inputs.enabled != false"
        forEach       = "$variables.regions"
        wait          = "$signals.step_ready"
        parallelGroup = "steps_step_case_enterprise_b_group"
        outputs = {
          audit  = "$response.body.auditId"
          result = "$response.body.result"
        }
        body {
          input {
            incidentId = "$inputs.incidentId"
            source = "step_case_enterprise_b"
          }
          meta {
            attempt = 1
            correlationId = "$context.correlationId"
          }
        }
        timeout = 30
        extensions {
          x-name = "step_case_enterprise_b"
          x-kind = "step"
        }
      }
      extensions {
        x-kind = "case"
        x-name = "case_enterprise"
      }
    }
    default "step_switch_default_notify" {
      description   = "Execute run_smtp_fallback and expose normalized outputs."
      operationRef  = "run_smtp_fallback"
      dependsOn     = ["run_llm_primary", "run_sql_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_switch_default_notify_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
        input {
          source = "step_switch_default_notify"
          incidentId = "$inputs.incidentId"
        }
      }
      timeout = 30
      extensions {
        x-name = "step_switch_default_notify"
        x-kind = "step"
      }
    }
    default "step_switch_default_archive" {
      description   = "Execute run_s3_fallback and expose normalized outputs."
      operationRef  = "run_s3_fallback"
      dependsOn     = ["run_llm_primary", "run_sql_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_switch_default_archive_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_switch_default_archive"
        }
        meta {
          correlationId = "$context.correlationId"
          attempt = 1
        }
      }
      timeout = 30
      extensions {
        x-name = "step_switch_default_archive"
        x-kind = "step"
      }
    }
    extensions {
      x-kind = "workflow"
      x-name = "switch"
    }
  }
  workflow "wf_loop" {
    type        = "loop"
    description = "Apply containment across affected hosts and evidence locations."
    dependsOn   = ["run_fileio_primary", "run_s3_primary"]
    when        = "$inputs.assets != []"
    forEach     = "$inputs.assets"
    wait        = "$signals.loop_ready"
    items       = "$inputs.assets"
    mode        = "parallel"
    batchSize   = "2"
    outputs = {
      containmentResults = "$steps.*.outputs.result"
      evidenceIds        = "$steps.step_loop_scp.outputs.audit"
    }
    inputs {
      type     = "object"
      format   = "uws-loopInput"
      _ref     = "#/components/schemas/loopInput"
      required = ["incidentId", "severity"]
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "loopInput-incident"
          x-kind = "schema-property"
        }
      }
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-kind = "schema-property"
          x-name = "loopInput-severity"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "loopInput"
          x-kind = "schema-items"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "loopInput-tenant"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        properties "source" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "loopInput-source"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "ip" {
          type = "string"
        }
        properties "host" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "loopInput-host"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "loopInput-user"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "caseId" {
          type = "string"
        }
        properties "ticket" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "loopInput-ticket"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "owner" {
          type = "string"
        }
        properties "asset" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "loopInput-asset"
        }
      }
      extensions {
        x-name = "loopInput"
        x-kind = "schema"
      }
    }
    timeout = 180
    step "step_loop_scp" {
      description   = "Execute run_scp_primary and expose normalized outputs."
      operationRef  = "run_scp_primary"
      dependsOn     = ["run_fileio_primary", "run_s3_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_loop_scp_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
        input {
          incidentId = "$inputs.incidentId"
          source = "step_loop_scp"
        }
      }
      timeout = 30
      extensions {
        x-name = "step_loop_scp"
        x-kind = "step"
      }
    }
    step "step_loop_sftp" {
      description   = "Execute run_sftp_primary and expose normalized outputs."
      operationRef  = "run_sftp_primary"
      dependsOn     = ["run_fileio_primary", "run_s3_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_loop_sftp_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          source = "step_loop_sftp"
          incidentId = "$inputs.incidentId"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_loop_sftp"
      }
    }
    step "step_loop_cmd" {
      description   = "Execute run_cmd_primary and expose normalized outputs."
      operationRef  = "run_cmd_primary"
      dependsOn     = ["run_fileio_primary", "run_s3_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_loop_cmd_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          source = "step_loop_cmd"
          incidentId = "$inputs.incidentId"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_loop_cmd"
      }
    }
    extensions {
      x-name = "loop"
      x-kind = "workflow"
    }
  }
  workflow "wf_await" {
    type        = "await"
    description = "Wait for external approval before applying high-risk actions."
    dependsOn   = ["run_smtp_primary", "run_llm_primary"]
    when        = "$inputs.requiresApproval == true"
    forEach     = "$variables.regions"
    wait        = "$signals.approval_received"
    outputs = {
      approval = "$signals.approval_received"
      summary  = "$steps.step_await_llm.outputs.result"
    }
    inputs {
      type     = "object"
      format   = "uws-awaitInput"
      _ref     = "#/components/schemas/awaitInput"
      required = ["incidentId", "severity"]
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-name = "awaitInput-severity"
          x-kind = "schema-property"
        }
      }
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-property"
          x-name = "awaitInput-incident"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "awaitInput"
          x-kind = "schema-items"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "awaitInput-tenant"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "source" {
          type = "string"
        }
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        extensions {
          x-name = "awaitInput-source"
          x-kind = "schema-allOf"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "host" {
          type = "string"
        }
        properties "ip" {
          type = "string"
        }
        extensions {
          x-name = "awaitInput-host"
          x-kind = "schema-oneOf"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "awaitInput-user"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "caseId" {
          type = "string"
        }
        properties "ticket" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "awaitInput-ticket"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "asset" {
          type = "string"
        }
        properties "owner" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "awaitInput-asset"
        }
      }
      extensions {
        x-name = "awaitInput"
        x-kind = "schema"
      }
    }
    timeout = 600
    step "step_await_smtp" {
      description   = "Execute run_smtp_fallback and expose normalized outputs."
      operationRef  = "run_smtp_fallback"
      dependsOn     = ["run_smtp_primary", "run_llm_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_await_smtp_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_await_smtp"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-name = "step_await_smtp"
        x-kind = "step"
      }
    }
    step "step_await_llm" {
      description   = "Execute run_llm_fallback and expose normalized outputs."
      operationRef  = "run_llm_fallback"
      dependsOn     = ["run_smtp_primary", "run_llm_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_await_llm_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_await_llm"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-name = "step_await_llm"
        x-kind = "step"
      }
    }
    extensions {
      x-name = "await"
      x-kind = "workflow"
    }
  }
  workflow "wf_merge" {
    type        = "merge"
    description = "Merge all branch outputs into a single handoff summary."
    dependsOn   = ["wf_switch", "wf_loop"]
    when        = "$workflows.wf_switch.outputs.selectedPath != ''"
    forEach     = "$variables.regions"
    wait        = "$signals.merge_ready"
    outputs = {
      archive = "$steps.step_merge_s3.outputs.audit"
      summary = "$steps.step_merge_fileio.outputs.result"
    }
    inputs {
      type     = "object"
      format   = "uws-mergeInput"
      _ref     = "#/components/schemas/mergeInput"
      required = ["incidentId", "severity"]
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-name = "mergeInput-severity"
          x-kind = "schema-property"
        }
      }
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "mergeInput-incident"
          x-kind = "schema-property"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "mergeInput"
          x-kind = "schema-items"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-name = "mergeInput-tenant"
          x-kind = "schema-allOf"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        properties "source" {
          type = "string"
        }
        extensions {
          x-name = "mergeInput-source"
          x-kind = "schema-allOf"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "host" {
          type = "string"
        }
        properties "ip" {
          type = "string"
        }
        extensions {
          x-name = "mergeInput-host"
          x-kind = "schema-oneOf"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-name = "mergeInput-user"
          x-kind = "schema-oneOf"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "caseId" {
          type = "string"
        }
        properties "ticket" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "mergeInput-ticket"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "asset" {
          type = "string"
        }
        properties "owner" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "mergeInput-asset"
        }
      }
      extensions {
        x-name = "mergeInput"
        x-kind = "schema"
      }
    }
    timeout = 60
    step "step_merge_fileio" {
      description   = "Execute run_fileio_fallback and expose normalized outputs."
      operationRef  = "run_fileio_fallback"
      dependsOn     = ["wf_switch", "wf_loop"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_merge_fileio_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          source = "step_merge_fileio"
          incidentId = "$inputs.incidentId"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_merge_fileio"
      }
    }
    step "step_merge_s3" {
      description   = "Execute run_s3_fallback and expose normalized outputs."
      operationRef  = "run_s3_fallback"
      dependsOn     = ["wf_switch", "wf_loop"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_merge_s3_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_merge_s3"
        }
        meta {
          correlationId = "$context.correlationId"
          attempt = 1
        }
      }
      timeout = 30
      extensions {
        x-name = "step_merge_s3"
        x-kind = "step"
      }
    }
    extensions {
      x-kind = "workflow"
      x-name = "merge"
    }
  }
  workflow "wf_manual_review" {
    type        = "sequence"
    description = "Collect evidence for human review when automated policy blocks execution."
    dependsOn   = ["fetch_ticket", "load_customer"]
    when        = "$inputs.requiresManualReview == true"
    forEach     = "$variables.regions"
    wait        = "$signals.reviewer_available"
    outputs = {
      notification  = "$steps.step_manual_smtp.outputs.audit"
      reviewPackage = "$steps.step_manual_fnct.outputs.result"
    }
    inputs {
      type     = "object"
      format   = "uws-manualReviewInput"
      _ref     = "#/components/schemas/manualReviewInput"
      required = ["incidentId", "severity"]
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-property"
          x-name = "manualReviewInput-incident"
        }
      }
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-name = "manualReviewInput-severity"
          x-kind = "schema-property"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "manualReviewInput"
          x-kind = "schema-items"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-name = "manualReviewInput-tenant"
          x-kind = "schema-allOf"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        properties "source" {
          type = "string"
        }
        extensions {
          x-name = "manualReviewInput-source"
          x-kind = "schema-allOf"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "host" {
          type = "string"
        }
        properties "ip" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "manualReviewInput-host"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "manualReviewInput-user"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "ticket" {
          type = "string"
        }
        properties "caseId" {
          type = "string"
        }
        extensions {
          x-name = "manualReviewInput-ticket"
          x-kind = "schema-anyOf"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "asset" {
          type = "string"
        }
        properties "owner" {
          type = "string"
        }
        extensions {
          x-name = "manualReviewInput-asset"
          x-kind = "schema-anyOf"
        }
      }
      extensions {
        x-name = "manualReviewInput"
        x-kind = "schema"
      }
    }
    timeout = 300
    step "step_manual_fnct" {
      description   = "Execute run_fnct_fallback and expose normalized outputs."
      operationRef  = "run_fnct_fallback"
      dependsOn     = ["fetch_ticket", "load_customer"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_manual_fnct_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_manual_fnct"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_manual_fnct"
      }
    }
    step "step_manual_smtp" {
      description   = "Execute run_smtp_fallback and expose normalized outputs."
      operationRef  = "run_smtp_fallback"
      dependsOn     = ["fetch_ticket", "load_customer"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_manual_smtp_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          incidentId = "$inputs.incidentId"
          source = "step_manual_smtp"
        }
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
      }
      timeout = 30
      extensions {
        x-name = "step_manual_smtp"
        x-kind = "step"
      }
    }
    extensions {
      x-name = "manual-review"
      x-kind = "workflow"
    }
  }
  workflow "wf_notify_operators" {
    type        = "sequence"
    description = "Notify operators and stakeholders after evidence is prepared."
    dependsOn   = ["run_smtp_primary", "run_llm_primary"]
    when        = "$inputs.suppressNotifications != true"
    forEach     = "$variables.regions"
    wait        = "$signals.notification_window"
    outputs = {
      brief     = "$steps.step_notify_llm.outputs.result"
      messageId = "$steps.step_notify_smtp.outputs.audit"
    }
    inputs {
      type     = "object"
      format   = "uws-notifyInput"
      _ref     = "#/components/schemas/notifyInput"
      required = ["incidentId", "severity"]
      properties "incidentId" {
        type   = "string"
        format = "uuid"
        extensions {
          x-name = "notifyInput-incident"
          x-kind = "schema-property"
        }
      }
      properties "severity" {
        type   = "string"
        format = "enum"
        extensions {
          x-name = "notifyInput-severity"
          x-kind = "schema-property"
        }
      }
      items {
        type   = "string"
        format = "uuid"
        extensions {
          x-kind = "schema-items"
          x-name = "notifyInput"
        }
      }
      allOf {
        type     = "object"
        required = ["tenantId", "region"]
        properties "region" {
          type = "string"
        }
        properties "tenantId" {
          type = "string"
        }
        extensions {
          x-name = "notifyInput-tenant"
          x-kind = "schema-allOf"
        }
      }
      allOf {
        type     = "object"
        required = ["source", "priority"]
        properties "priority" {
          type   = "integer"
          format = "int32"
        }
        properties "source" {
          type = "string"
        }
        extensions {
          x-kind = "schema-allOf"
          x-name = "notifyInput-source"
        }
      }
      oneOf {
        type     = "object"
        required = ["host", "ip"]
        properties "host" {
          type = "string"
        }
        properties "ip" {
          type = "string"
        }
        extensions {
          x-name = "notifyInput-host"
          x-kind = "schema-oneOf"
        }
      }
      oneOf {
        type     = "object"
        required = ["user", "email"]
        properties "email" {
          type   = "string"
          format = "email"
        }
        properties "user" {
          type = "string"
        }
        extensions {
          x-kind = "schema-oneOf"
          x-name = "notifyInput-user"
        }
      }
      anyOf {
        type     = "object"
        required = ["ticket", "caseId"]
        properties "caseId" {
          type = "string"
        }
        properties "ticket" {
          type = "string"
        }
        extensions {
          x-kind = "schema-anyOf"
          x-name = "notifyInput-ticket"
        }
      }
      anyOf {
        type     = "object"
        required = ["asset", "owner"]
        properties "asset" {
          type = "string"
        }
        properties "owner" {
          type = "string"
        }
        extensions {
          x-name = "notifyInput-asset"
          x-kind = "schema-anyOf"
        }
      }
      extensions {
        x-name = "notifyInput"
        x-kind = "schema"
      }
    }
    timeout = 80
    step "step_notify_smtp" {
      description   = "Execute run_smtp_primary and expose normalized outputs."
      operationRef  = "run_smtp_primary"
      dependsOn     = ["run_smtp_primary", "run_llm_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_notify_smtp_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        input {
          source = "step_notify_smtp"
          incidentId = "$inputs.incidentId"
        }
        meta {
          correlationId = "$context.correlationId"
          attempt = 1
        }
      }
      timeout = 30
      extensions {
        x-name = "step_notify_smtp"
        x-kind = "step"
      }
    }
    step "step_notify_llm" {
      description   = "Execute run_llm_primary and expose normalized outputs."
      operationRef  = "run_llm_primary"
      dependsOn     = ["run_smtp_primary", "run_llm_primary"]
      when          = "$inputs.enabled != false"
      forEach       = "$variables.regions"
      wait          = "$signals.step_ready"
      parallelGroup = "steps_step_notify_llm_group"
      outputs = {
        audit  = "$response.body.auditId"
        result = "$response.body.result"
      }
      body {
        meta {
          attempt = 1
          correlationId = "$context.correlationId"
        }
        input {
          incidentId = "$inputs.incidentId"
          source = "step_notify_llm"
        }
      }
      timeout = 30
      extensions {
        x-kind = "step"
        x-name = "step_notify_llm"
      }
    }
    extensions {
      x-kind = "workflow"
      x-name = "notify"
    }
  }
  trigger "alert_webhook" {
    path           = "/webhooks/security/alerts"
    methods        = ["POST", "PUT"]
    authentication = "bearer"
    outputs        = ["alert.received", "alert.replayed"]
    options {
      acceptedContentTypes = [
        "application/json",
        "application/problem+json"
      ]
      deduplicateWindowSeconds = 300
    }
    route {
      output = "alert.received"
      to     = ["main", "wf_parallel"]
      extensions {
        x-name = "primary"
        x-kind = "route"
      }
    }
    route {
      output = "alert.replayed"
      to     = ["step_collect_context", "wf_switch"]
      extensions {
        x-kind = "route"
        x-name = "replay"
      }
    }
    extensions {
      x-kind = "trigger"
      x-name = "alert"
    }
  }
  trigger "operator_command" {
    path           = "/commands/security/remediate"
    methods        = ["POST", "PATCH"]
    authentication = "signed"
    outputs        = ["command.approved", "command.denied"]
    options {
      allowedTeams = [
        "secops",
        "platform"
      ]
      requiresMFA = true
    }
    route {
      output = "command.approved"
      to     = ["wf_loop", "wf_await"]
      extensions {
        x-kind = "route"
        x-name = "approved"
      }
    }
    route {
      output = "command.denied"
      to     = ["wf_manual_review", "step_notify"]
      extensions {
        x-kind = "route"
        x-name = "denied"
      }
    }
    extensions {
      x-name = "operator"
      x-kind = "trigger"
    }
  }
  result "decision.branch" {
    kind  = "switch"
    from  = "wf_switch"
    value = "$workflows.wf_switch.outputs.selectedPath"
    extensions {
      x-name = "switch"
      x-kind = "result"
    }
  }
  result "containment.loop" {
    kind  = "loop"
    from  = "wf_loop"
    value = "$workflows.wf_loop.outputs.containmentResults"
    extensions {
      x-name = "loop"
      x-kind = "result"
    }
  }
  result "merge.summary" {
    kind  = "merge"
    from  = "wf_merge"
    value = "$workflows.wf_merge.outputs.summary"
    extensions {
      x-kind = "result"
      x-name = "merge"
    }
  }
  components {
    variables {
      notificationLists {
        operator = [
          "secops@example.test",
          "platform@example.test"
        ]
        executive = [
          "ciso@example.test",
          "cio@example.test"
        ]
      }
      priorityThreshold = 80
    }
    extensions {
      x-name = "shared"
      x-kind = "components"
    }
  }
  extensions {
    x-generator = "testdata/big/main.go"
    x-purpose = "large-runtime-roundtrip-fixture"
  }