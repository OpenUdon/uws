# Feature 2: Six Structural Constructs

← [API Source Operation Binding](01-OpenAPI-Operation-Binding.md) | [Next: Runtime Expression Grammar →](03-Runtime-Expression-Grammar.md)

---

Operations are the leaves. Workflows and steps compose them using six structural control-flow constructs. Each workflow declares exactly one `type`; nested steps may also declare a structural `type` to form composite control flow.

This guide describes current UWS 1.11 behavior unless it identifies an older
version explicitly. Version-gated behavior is selected by the document's
declared UWS version.

## `sequence`

Steps execute in declaration order. Each step completes before the next begins. Use `sequence` for any pipeline where outputs flow from one step to the next.

```yaml
workflowId: checkout
type: sequence
steps:
  - stepId: validate_cart
    operationRef: validate_order
  - stepId: charge_card
    operationRef: charge_payment
    dependsOn: [validate_cart]
  - stepId: send_receipt
    operationRef: send_email
    dependsOn: [charge_card]
```

`dependsOn` within a sequence adds explicit cross-step dependencies beyond declaration order. `items` MUST NOT be set.

**Adding a cross-workflow dependency from another workflow:**

```yaml
workflowId: post_checkout
type: sequence
dependsOn: [checkout]        # waits for the entire checkout workflow first
steps:
  - stepId: update_inventory
    operationRef: decrement_stock
```

## `parallel`

Steps execute concurrently, subject only to `dependsOn` relationships. Use `parallel` when independent calls can overlap.

```yaml
workflowId: enrichment
type: parallel
steps:
  - stepId: fetch_weather
    operationRef: get_weather
  - stepId: fetch_stocks
    operationRef: get_stocks
  - stepId: fetch_news
    operationRef: get_headlines
```

All three operations fire simultaneously. The construct completes when every step terminates.

**`parallelGroup` as a dependency barrier:**

```yaml
workflowId: validate_and_merge
type: parallel
steps:
  - stepId: check_format
    operationRef: validate_format
    parallelGroup: validators     # member of the "validators" group

  - stepId: check_range
    operationRef: validate_range
    parallelGroup: validators     # also a member

  - stepId: aggregate
    operationRef: merge_results
    dependsOn: [validators]       # waits for ALL members of "validators"
```

`dependsOn: [validators]` waits for every step in the `validators` group. Membership in a group does not create additional ordering among members themselves — they still run concurrently. From UWS 1.10, a branch failure cancels sibling branch contexts. A `goto` or `end` control signal from a branch fails the parallel construct because cross-branch control flow is undefined.

## `switch`

Cases are evaluated in declaration order. The first truthy `when` runs; a case without `when` is an unconditional match at its position. If no case matches, `default` runs if present. No later case runs.

```yaml
workflowId: route_event
type: switch
cases:
  - name: new_user
    when: $trigger.body.event == "signup"
    steps:
      - stepId: welcome
        operationRef: send_welcome_email

  - name: returning_user
    when: $trigger.body.event == "login"
    steps:
      - stepId: log_access
        operationRef: record_login

default:
  - stepId: fallback_log
    operationRef: log_unknown_event
```

`items` MUST NOT be set on a `switch`.

**`switch` inside a `sequence` step:**

```yaml
workflowId: process_order
type: sequence
steps:
  - stepId: fetch_order
    operationRef: get_order
  - stepId: route
    type: switch
    cases:
      - name: express
        when: $steps.fetch_order.outputs.tier == "express"
        steps:
          - stepId: expedite
            operationRef: fast_ship
      - name: standard
        when: $steps.fetch_order.outputs.tier == "standard"
        steps:
          - stepId: queue
            operationRef: standard_ship
```

## `loop`

Iterates over an ordered JSON array. `items` is REQUIRED and must resolve to a JSON array at runtime. Each element is bound into the iteration scope as `$item`, with a zero-based `$index`; nested steps execute sequentially in source order.

```yaml
workflowId: notify_all
type: loop
items: $outputs.subscribers
steps:
  - stepId: send_one
    operationRef: send_notification
```

**`loop` with `batchSize` — processing in fixed groups:**

```yaml
workflowId: bulk_import
type: loop
items: $outputs.records           # e.g. array of 1000 records
batchSize: "50"                   # UWS 1.11 numeric literal: 50 items per batch
steps:
  - stepId: import_batch
    operationRef: bulk_upsert
```

`batchSize` MUST resolve to a positive integer. In UWS 1.11 and later, a complete JSON-number literal such as `"50"` is allowed; the numeric-literal grammar is not available to earlier declarations. Batches and iterations execute sequentially; `$batchIndex` is the zero-based batch number and is available only inside a `loop`. A successful `loop` produces an ordered array of `{index, batchIndex, item}` records, including an empty array when there are no items. `cases` and `default` MUST NOT be set on a `loop`.

## `merge`

Combines the outputs of multiple upstream constructs named by `dependsOn` into a single structural result. Use `merge` after a `parallel` to collect independent results.

```yaml
workflowId: gather_data
type: parallel
steps:
  - stepId: prices
    operationRef: get_prices
    parallelGroup: fetchers
  - stepId: ratings
    operationRef: get_ratings
    parallelGroup: fetchers
  - stepId: combine
    type: merge
    dependsOn: [fetchers]
    outputs:
      prices:  $steps.prices.outputs.data
      ratings: $steps.ratings.outputs.data
```

The corresponding result declaration:

```yaml
results:
  - name: market_data
    kind: merge
    from: gather_data.combine
    value: $steps.combine.outputs
```

`dependsOn` is REQUIRED and MUST name at least one construct. `items` MUST NOT be set.

## `await`

Evaluates its `wait` predicate immediately, then polls it at the executor's configured interval (200 ms by default). Nested steps run once after the predicate becomes truthy. An `await` does not re-execute already-completed steps or operations, so it cannot portably poll a status endpoint by placing that endpoint in a preceding sequence step.

```yaml
workflowId: wait_for_external_signal
type: await
wait: $inputs.signal_ready == true
timeout: 300
```

The example requires a bound runtime or profile to supply an input whose value can change between predicate evaluations; UWS core does not define how an external signal is refreshed. Polling a remote job status requires an implementation-specific runtime/profile or separately scheduled invocations. `timeout` has been available on operations, workflows, and steps since UWS 1.1. A serialized timeout bounds the await; an executor-owned timeout MAY apply when it is absent. Context cancellation also stops polling. For non-`await` constructs, UWS 1.10 defines `wait` as a cancellable delay expression resolving to a finite number of seconds from 0 to 86,400, evaluated once before the body. In UWS 1.11 and later, a complete JSON-number literal may be used for this delay. `cases`, `default`, and `items` MUST NOT be set on `await`. See the [UWS 1.11 execution contract](https://github.com/OpenUdon/uws/blob/main/versions/1.11.0.md#78-uws-110-and-111-portable-execution-semantics).

## Field Constraints Summary

| Type | `items` | `wait` | `cases`/`default` | `dependsOn` |
|------|---------|--------|-------------------|-------------|
| `sequence` | MUST NOT | optional | MUST NOT | optional |
| `parallel` | MUST NOT | optional | MUST NOT | optional |
| `switch` | MUST NOT | optional | allowed | optional |
| `loop` | REQUIRED | optional | MUST NOT | optional |
| `merge` | MUST NOT | optional | MUST NOT | REQUIRED (≥1) |
| `await` | MUST NOT | REQUIRED | MUST NOT | optional |

The validator enforces every constraint above before the runtime sees the document.

## Composing Types: Nested Steps

A top-level `workflow` always declares a `type`. A nested `step` may declare a `type` to become an inline structural construct, use `operationRef` to call an operation, or use `workflow` to invoke a top-level workflow by ID.

```yaml
# A sequence whose third step is itself a loop
workflowId: process_batch
type: sequence
steps:
  - stepId: fetch_all
    operationRef: list_items
    outputs:
      items: $response.body.items

  - stepId: validate_all
    operationRef: validate_batch

  - stepId: process_each
    type: loop
    items: $steps.fetch_all.outputs.items
    steps:
      - stepId: handle_item
        operationRef: process_single_item
```

## From The Big Fixture

The large fixture exercises structural workflows and nested workflow calls. This
excerpt shows a sequence workflow step invoking another workflow:

```hcl
workflow "main" {
  type        = "sequence"
  description = "Coordinate enrichment, runtime checks, branching, containment, and notification."
  dependsOn   = ["fetch_ticket", "load_customer"]
  outputs = {
    decision = "$steps.step_decide_path.outputs.selectedPath"
    incident = "$steps.step_collect_context.outputs.incident"
  }

  step "step_collect_context" {
    operationRef = "fetch_ticket"
    dependsOn    = ["run_cmd_primary", "run_fnct_primary"]
    outputs = {
      audit  = "$response.body.auditId"
      result = "$response.body.result"
    }
  }

  step "step_parallel_checks" {
    workflow  = "wf_parallel"
    dependsOn = ["step_collect_context", "load_customer"]
  }
}
```

Full context: [`testdata/big/big.hcl`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.hcl).

---

← [API Source Operation Binding](01-OpenAPI-Operation-Binding.md) | [Next: Runtime Expression Grammar →](03-Runtime-Expression-Grammar.md)
