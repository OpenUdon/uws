# Quick Start: Author Your First UWS Workflow

This guide is for authors who need to hand-write a UWS document without reading the full specification first. It focuses on the common path: bind existing source operations, put them in a workflow, pass values between steps, and leave concrete credentials to the runtime.

UWS is a workflow overlay. Source documents such as OpenAPI, Google Discovery, AWS Smithy, and AsyncAPI own HTTP methods, paths, channels, messages, schemas, servers, and security. UWS owns operation binding, workflow structure, request values, outputs, triggers, and control flow.

## Start With Three Files

Keep the source contract, workflow, and runtime configuration separate:

```text
openapi/
  support.yaml
google-discovery/
  gmail.json
asyncapi/
  billing-events.yaml
workflow.uws.yaml
runtime-config.json       # runtime-owned, not UWS core
```

The UWS document points to source documents. It does not copy endpoint URLs, channels, schemas, or credentials out of the source document.

## Minimal Workflow

```yaml
uws: "1.3.0"
info:
  title: Support Ticket Workflow
  version: "1.0.0"

sourceDescriptions:
  - name: support_api
    url: ./openapi/support.yaml
    type: openapi

operations:
  - operationId: create_ticket
    sourceDescription: support_api
    sourceOperationId: createTicket
    request:
      body:
        subject: "Example ticket"
        priority: normal
    outputs:
      ticket_id: "$response.body.id"

  - operationId: get_ticket
    sourceDescription: support_api
    sourceOperationId: getTicket
    dependsOn: [create_ticket]
    request:
      path:
        ticketId: "$steps.create.outputs.ticket_id"
    outputs:
      status: "$response.body.status"

workflows:
  - workflowId: main
    type: sequence
    steps:
      - stepId: create
        operationRef: create_ticket
      - stepId: verify
        operationRef: get_ticket
```

This document declares:

- One OpenAPI source named `support_api`.
- Two UWS-local operations, `create_ticket` and `get_ticket`.
- One workflow named `main`.
- Two sequence steps, where `verify` can refer to the output of step `create`.

## Step 1: Declare Sources

Use `sourceDescriptions` for every API or event source document the workflow references:

```yaml
sourceDescriptions:
  - name: support_api
    url: ./openapi/support.yaml
    type: openapi
```

Rules of thumb:

- `name` is the stable local handle used by operations.
- `url` can be a local path or reviewed remote location, depending on the runtime.
- `type` may be `openapi`, `google-discovery`, `aws-smithy`, or `asyncapi`. Missing `type` defaults to `openapi`.

## Step 2: Bind Operations

The preferred operation binding uses `sourceDescription` plus the generic `sourceOperationId`:

```yaml
operations:
  - operationId: create_ticket
    sourceDescription: support_api
    sourceOperationId: createTicket
```

`operationId` is local to the UWS file. `sourceOperationId` is the operation ID from the referenced source document. For AsyncAPI sources, it names a root AsyncAPI Operation Object. OpenAPI sources may also use legacy `openapiOperationId`.

Do not add HTTP method, path, channel, server URL, or security configuration here. Those belong to the source document and the bound runtime.

## Step 3: Add Request Values

Request values are grouped by common source request location:

```yaml
request:
  path:
    ticketId: "$steps.create.outputs.ticket_id"
  query:
    includeHistory: "true"
  header:
    X-Request-Source: review
  body:
    priority: normal
```

Use runtime expressions when a value comes from a previous step, workflow input, trigger, variable, or response.
For AsyncAPI source operations, use `body` for the message payload or operation input values and `header` for AsyncAPI Message Object `headers` values when the source operation defines them.

## Step 4: Name Outputs

Outputs turn response values into stable handles:

```yaml
outputs:
  ticket_id: "$response.body.id"
  status: "$response.body.status"
```

Common expression roots:

- `$response.body`: response body from the current operation.
- `$response.headers`: response headers from the current operation.
- `$steps.<stepId>.outputs.<name>`: output from a previous step in the same workflow.
- `$inputs.<name>`: workflow input.
- `$variables.<name>`: reusable document variable.
- `$trigger`: trigger payload during trigger dispatch.

Keep expressions simple unless your runtime documents extensions beyond UWS core.

## Step 5: Compose Workflows

A sequence workflow runs steps in order:

```yaml
workflows:
  - workflowId: main
    type: sequence
    steps:
      - stepId: create
        operationRef: create_ticket
      - stepId: verify
        operationRef: get_ticket
```

Use the six structural types when needed:

- `sequence`: run steps in declaration order.
- `parallel`: run independent steps concurrently.
- `switch`: choose a branch.
- `loop`: repeat over items.
- `merge`: join control-flow branches.
- `await`: wait until a runtime expression becomes true.

Start with `sequence`. Add other constructs only when the workflow needs them.

## Add A Branch

```yaml
workflows:
  - workflowId: classify_ticket
    type: switch
    cases:
      - name: closed
        when: "$steps.verify.outputs.status == \"closed\""
        steps:
          - stepId: archive
            operationRef: archive_ticket
    default:
      - stepId: notify_owner
        operationRef: notify_owner
```

Switch cases contain branch steps. Keep branch conditions in the core expression grammar when portability matters.

## Add A Trigger

```yaml
triggers:
  - triggerId: new_ticket
    outputs:
      - accepted
    routes:
      - output: accepted
        to: [main]
```

Trigger ingress, authentication, hosting, output classification, and persistence are runtime-owned. UWS only describes declared output labels and route dispatch after a trigger event has been accepted.

## Extension-Owned Operations

If a step is not a source-bound operation, make that explicit with an extension profile:

```yaml
operations:
  - operationId: render_message
    x-uws-operation-profile: uws.runtime.1.0
    x-uws-runtime:
      type: fnct
      function: render_message
      arguments:
        - template: "Ticket {{ticket_id}} is ready"
```

UWS validates the shape and references. The bound runtime decides whether it supports the named extension profile.

When you use the public `uws.runtime.1.0` supplement, `x-uws-runtime.type` is required and must be one of the supplement's non-HTTP runtime selectors. Use source operation bindings for HTTP and event calls; do not model HTTP as `x-uws-runtime.type: http`.

## Validation Checklist

Before handing a workflow to a runtime:

- Every `sourceDescription` name is unique.
- Every source-bound operation names an existing source.
- Every local `operationId`, `workflowId`, and step `stepId` is unique where required.
- Every `operationRef`, `workflow`, trigger route, and dependency target resolves.
- Request values do not contain secrets.
- Credentials and endpoint policy are supplied by the runtime, not embedded in UWS.
- Extension-owned operations declare an explicit `x-uws-operation-profile`.
- `uws.runtime.1.0` operations include `x-uws-runtime.type`; HTTP and event calls remain source-bound.

## Where To Go Next

- [Source Operation Binding](01-OpenAPI-Operation-Binding.md)
- [Six Structural Constructs](02-Six-Structural-Constructs.md)
- [Runtime Expression Grammar](03-Runtime-Expression-Grammar.md)
- [Triggers and Route Dispatch](04-Triggers-and-Route-Dispatch.md)
- [Validation](09-Validation.md)
