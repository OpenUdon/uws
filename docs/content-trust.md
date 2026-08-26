# Content Trust: Integrity of Data Flowing Through a Workflow

This guide explains the additive UWS 1.9.1 `contentTrust` registry and the
advisory analyzer in `github.com/OpenUdon/uws/contenttrust`. The normative wire
contract is [versions/1.9.1.md](https://github.com/OpenUdon/uws/blob/main/versions/1.9.1.md);
this page focuses on why the feature exists and how authors and resolvers use it.

## The gap

UWS already had a strong **confidentiality model**. UWS 1.9.1 adds the first
portable **integrity model** for content flowing through a workflow.

The confidentiality side is well developed and spread across several profiles:
`noLog` sensitive parameters, symbolic `credentialSlots`, credential bindings
resolved only by a trusted runtime, and the explicit "MUST NOT contain" lists in
the browser authentication and registration profiles. Every one of those rules
governs secret material flowing *outward* — into documents, prompts, artifacts,
logs, reports, command arguments, and workflow outputs.

Those controls do not describe untrusted data flowing *inward*.

The core Security Considerations section states that tooling "SHOULD treat
documents from untrusted sources as untrusted input." That is a statement about
the trustworthiness of the *document*. It says nothing about the trustworthiness
of *content that operations return at runtime*, which is a different problem with
a different blast radius.

## Why this matters more in UWS than in a single API call

A standalone API call returns data to a caller who decides what to do with it.
UWS chains operations. A value read in one step reaches later steps through
runtime expressions:

```yaml
uws: 1.9.1
info: {title: Mail triage, version: 1.0.0}
sourceDescriptions:
  - name: mail_api
    url: ./mail.openapi.yaml
    type: openapi
operations:
  - operationId: read_message
    sourceDescription: mail_api
    sourceOperationId: getMessage
    outputs:
      body: $response.body#/payload/text

  - operationId: summarize
    x-uws-operation-profile: uws.runtime.1.0
    x-uws-runtime:
      type: llm
    request:
      body:
        data: $inputs.message
    outputs:
      destination: $response.body#/destination

  - operationId: post_summary
    sourceDescription: mail_api
    sourceOperationId: postSummary
    request:
      body:
        mailbox: $inputs.mailbox

workflows:
  - workflowId: main
    type: sequence
    steps:
      - stepId: read
        operationRef: read_message
        outputs:
          body: $outputs.body
      - stepId: model
        operationRef: summarize
        inputs:
          message: $steps.read.outputs.body
        outputs:
          destination: $outputs.destination
      - stepId: publish
        operationRef: post_summary
        inputs:
          mailbox: $steps.model.outputs.destination

contentTrust:
  sourceDescriptions:
    mail_api: untrusted
  operations:
    read_message:
      outputs:
        body: untrusted
```

Three properties combine to make this sharp:

1. **Operations chain by design.** `$steps.<id>.outputs.<name>` feeds a later
   `request.body`, `when`, `items`, or `successCriteria`.
2. **Profiles know channel semantics.** A runtime resolver can classify the
   model's `data` field as data and the posting operation's `mailbox` field as
   authority-bearing. LLM data is allowed, but model output inherits its
   provenance and later authority-bearing use remains visible.
3. **The dangerous channels are ordinary.** Inbound email bodies, CRM notes,
   support tickets, webhook payloads, pull-request descriptions, and text
   extracted from a web page are all normal workflow inputs and all
   attacker-influenceable.

A reviewed document can therefore route attacker-controlled text through a
model and into a side effect. UWS 1.9.1 makes that route reviewable without
treating every request body or every LLM data field as a dangerous sink.

## The opportunity

Because a UWS workflow is a **static document**, its core data flow is
recoverable before anything executes. Sequence order, explicit dependencies,
branches, loops, same-document workflow calls, request bindings, inputs, and
operation, workflow, and step outputs are all visible.

The boundary is exact: extension-profile flow is known only when a resolver
describes its channels and references. Without a resolver, implementation-
specific interpolation remains opaque and provenance across that boundary is
`unknown`. The analyzer reports this uncertainty instead of pretending the
whole graph is known.

Integrity analysis is explicit and advisory, not part of ordinary validation.
`Document.Validate` rejects malformed declarations, but analyzer findings do
not enter `ValidationResult`, fail execution, or alter executor results.

This is not available to systems where a model decides the call sequence at
runtime; there is no graph to analyse until execution is already underway. It is
one of the most concrete payoffs of the reviewable-document design.

## The additive wire shape

Operation, workflow, and step `outputs` remain `Map[string, string]`. Enriching
an output entry into an object would break a shape published since UWS 1.0, so
1.9.1 adds a parallel root registry:

```yaml
contentTrust:
  sourceDescriptions:
    mail_api: untrusted
  operations:
    read_message:
      default: untrusted
      outputs:
        message_id: trusted
        body: untrusted
  triggers:
    incoming_mail: untrusted
  workflows:
    summarize:
      default: unknown
      inputs:
        locale: trusted
```

Every key must resolve to a declaration in the same document. Operation output
keys must exist in that operation's `outputs`; workflow input keys must exist in
the workflow's `inputs.properties`. Empty registry and declaration objects are
invalid. `contentTrust` requires `uws: 1.9.1` or later.

Operation output precedence is:

1. per-output declaration;
2. operation `default`;
3. source-description declaration;
4. resolver-provided default;
5. `unknown`.

Trigger payloads default to `untrusted`. External workflow-entry inputs default
to `unknown`. Document literals, `variables`, and `components.variables` are
trusted only because analysis assumes the document itself has been reviewed.
An internal workflow call passes its computed provenance into the callee; a
callee's external-entry declaration cannot upgrade it.

## Provenance is not capability

The analyzer tracks two independent facts:

| Dimension | Values | Question |
| --- | --- | --- |
| Provenance | `trusted`, `untrusted`, `unknown` | Who can influence this value? |
| Capability | `free_text`, `constrained_scalar`, `composite`, `unknown` | What can the value carry? |

Boolean, numeric, and enum narrowing removes free-text injection capability. It
does **not** clear untrusted provenance. An attacker-controlled boolean cannot
carry prose into an instruction channel, but it can still choose a branch or
select an authority-bearing destination. Browser-aware resolvers therefore
classify string outputs as free text, numeric/boolean/presence/inline-enum
outputs as constrained scalars, and structured outputs as composite while
retaining their browser-derived untrusted provenance.

## Channels, not blanket sinks

A resolver describes how an operation interprets inputs:

- `data` is content to process. Untrusted LLM data is permitted, and a model
  output that inherits input provenance remains untrusted.
- `instruction` changes how an interpreter or model behaves. Untrusted
  free-form instruction content produces a high-severity finding.
- `authority` selects or authorizes a side effect, command, SQL target,
  destination, account, or equivalent capability. Untrusted influence produces
  a high-severity finding even when the value is constrained.

This is why the example's mail body may enter the model's data channel, while
the later use of derived model output in an authority channel is reported. UWS
does not label every request body or every LLM input as a sink.

An explicitly untrusted instruction channel is reported directly:

```yaml
operations:
  - operationId: instructed_model
    x-uws-operation-profile: uws.runtime.1.0
    x-uws-runtime:
      type: llm
    request:
      body:
        system: $inputs.external_system_text
```

The runtime-profile resolver classifies `request.body.system` as `instruction`.
If a step binds untrusted mail, browser, or trigger text to
`external_system_text`, analysis emits
`content_trust.untrusted_instruction`.

## Reference and control-flow analysis

The analyzer extracts exact expressions from request bindings, step inputs,
controls, workflow idempotency keys, criteria, structural-result values, and
operation, workflow, and step outputs. It understands the normative grammar
rather than searching for arbitrary substrings.

Declaration order inside a `sequence` establishes happens-before. Authors do
not need a redundant `dependsOn` merely to read an earlier sibling's output.
Explicit dependencies, parallel entry/exit, switch branches, loops, nested
structural steps, and same-document workflow calls are also included. A
branch-only producer does not dominate a consumer outside the branch, and a
loop-body producer does not dominate a later consumer because the loop may have
zero iterations.

An extension resolver can provide references for its own interpolation syntax.
Without those references, the analyzer scans the resolver-declared channel
subtree only for exact UWS expressions. Ambiguous syntax is reported as
`content_trust.opaque_expression`; content beyond an unresolved extension
boundary stays `unknown`.

Resolver contracts fail closed. Channel and reference paths must use RFC 6901
pointer syntax, every channel must resolve in the operation, channel kinds and
non-empty trust/capability labels must be recognized, and output contracts may
name only declared operation outputs. A malformed claim is discarded and
reported as `content_trust.resolver_failure`; if no valid resolver claim
remains, core request scanning still runs.

## Findings and policy

`contenttrust.Analyze` returns deterministic edges and findings. Reports contain
document paths, provenance/capability labels, stable codes, and fixed messages;
they contain no runtime values or content excerpts.

High-severity codes are:

- `content_trust.untrusted_instruction`
- `content_trust.untrusted_authority`

Warning codes are:

- `content_trust.untrusted_control`
- `content_trust.unknown_provenance`
- `content_trust.opaque_expression`
- `content_trust.unresolved_reference`
- `content_trust.non_dominating_reference`
- `content_trust.resolver_failure`
- `content_trust.resolver_conflict`

All severities are advisory in 1.9.1. Applications may present findings or
apply an explicit stricter policy, but UWS validation and execution behavior do
not change. UWS 1.9.2 through 1.9.9 are reserved for compatible refinements
supported by analyzer evidence; mandatory enforcement, incompatible defaults,
or wire restructuring require UWS 2.0.
