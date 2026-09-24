# Feature 5: Structural Results

← [Triggers and Route Dispatch](04-Triggers-and-Route-Dispatch.md) | [Next: Success Criteria and Actions →](06-Success-Criteria-and-Actions.md)

---

A structural result declares a named mapping to a workflow or step whose `type` is `switch`, `merge`, or `loop`. It does not create a result store or make an execution record portable. The optional `value` is an expression for a consumer to evaluate in the construct's execution context; its evaluation and exposure remain implementation-defined. UWS 1.11 does define the portable execution-record result shapes for `loop`, `forEach`, and `merge`, described below.

## Structural Result Object Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | REQUIRED | Unique result name within `results[]`. MUST match `^[a-zA-Z0-9._-]+$` |
| `kind` | REQUIRED | One of `switch`, `merge`, `loop`. MUST equal the `type` of the referenced construct |
| `from` | REQUIRED | `<workflowId>` or `<workflowId>.<stepId>` of the emitting construct |
| `value` | optional | Runtime expression for a consumer to evaluate against the construct's execution context; evaluation and exposure are implementation-defined |

## The `from` Field

`from` identifies the emitting construct in one of two forms:

- **`<workflowId>`** — a top-level workflow whose `type` is `switch`, `merge`, or `loop`.
- **`<workflowId>.<stepId>`** — a step within a named workflow whose `type` is `switch`, `merge`, or `loop`.

The validator resolves `from` to a real workflow or step, then checks that the referenced `type` matches `kind`. A mismatch produces a structured error. The `name` is a declaration label; UWS core does not guarantee persistence, materialization, or a particular lookup API for it.

## Why Only Three Kinds?

Only `switch`, `merge`, and `loop` are permitted `kind` values in `results[]`, but they do not all synthesize the same sort of result:

- **`switch`** — selects one branch in declaration order. It does not synthesize a portable selected-branch result; declare outputs explicitly when consumers need values.
- **`merge`** — its execution record contains an ordered array of records from the declared dependencies.
- **`loop`** — its execution record contains an ordered array of iteration metadata (`index`, `batchIndex`, `item`), not the nested operations' output values.

`sequence`, `parallel`, and `await` are not eligible `results[].kind` values. Step and workflow `outputs` remain explicit expression mappings. A `forEach` operation or step is not a structural-result kind; its parent execution record has an ordered per-item result array and exposes each declared output as an ordered array of per-iteration values.

## Portable Result Shapes in UWS 1.11

- A successful `loop` record has an ordered `result` array of `{index, batchIndex, item}` objects. Batching is sequential and does not imply parallel execution.
- A successful `forEach` parent record has an ordered `result` array of per-item objects: `{index, item, status, error, result, outputs}`. Each declared output is also present on the parent record as an ordered array of values.
- A successful `merge` record has an ordered `result` array of `{id, kind, status, error, result, outputs}` records. Dependencies retain declaration order; parallel groups expand in member declaration order. For a `forEach` dependency, iteration records replace the parent aggregate when present. If it was skipped or produced zero items, the parent record is included once. A failed dependency aborts the merge.

These are execution-record results, not values automatically extracted or stored by `results[]`. Execution-record keys, persistence, and result declaration lookup remain implementation-defined.

## Example 1: `merge` Result — Combining Parallel Checks

Two validation steps run in parallel; a merge step collects their records; a named result declaration provides a consumer expression for the combined step outputs.

```yaml
workflows:
  - workflowId: validate_order
    type: parallel
    steps:
      - stepId: check_inventory
        operationRef: validate_stock
        parallelGroup: validators
        outputs:
          ok:      $response.body.available
          message: $response.body.reason

      - stepId: check_credit
        operationRef: validate_payment
        parallelGroup: validators
        outputs:
          ok:      $response.body.approved
          limit:   $response.body.creditLimit

      - stepId: combine_checks
        type: merge
        dependsOn: [validators]
        outputs:
          inventory_ok: $steps.check_inventory.outputs.ok
          credit_ok:    $steps.check_credit.outputs.ok
          credit_limit: $steps.check_credit.outputs.limit

results:
  - name: order_validation
    kind: merge
    from: validate_order.combine_checks
    value: $steps.combine_checks.outputs
```

`order_validation` is a declaration that associates a label with the merge construct and a consumer expression. It does not itself extract, save, or publish the value; a consuming implementation evaluates `value` in the appropriate execution context.

## Example 2: `loop` Result — Ordered Iteration Metadata

A loop processes each item in an array. Its execution-record result contains ordered iteration metadata; the declared step outputs are separately aggregated on their parent record.

```yaml
workflows:
  - workflowId: import_records
    type: loop
    items: $outputs.pending_records
    steps:
      - stepId: upsert_record
        operationRef: create_or_update
        outputs:
          record_id: $response.body.id
          created:   $response.body.created

results:
  - name: import_summary
    kind: loop
    from: import_records
    value: $steps.upsert_record.outputs
```

The `loop` record's portable `result` contains `{index, batchIndex, item}` metadata. The nested step's declared outputs are separately available as ordered per-iteration arrays on its parent record; the optional `value` expression selects those outputs for a consumer. It does not redefine the loop record's result shape.

## Example 3: A `switch` Does Not Synthesize a Branch Result

A `switch` selects one processing path, but UWS core does not synthesize a portable value that identifies the selected case or combines branch outputs. A result declaration can name the construct, but a portable value must come from explicitly declared outputs and an expression the consumer can evaluate. Do not assume that an unexecuted branch's step output exists.

```yaml
workflows:
  - workflowId: classify_event
    type: switch
    cases:
      - name: high_value
        when: $trigger.amount >= 1000
        steps:
          - stepId: premium_process
            operationRef: handle_premium_order
            outputs:
              tier:     "premium"
              discount: $response.body.appliedDiscount

      - name: standard
        when: $trigger.amount < 1000
        steps:
          - stepId: normal_process
            operationRef: handle_standard_order
            outputs:
              tier:     "standard"
              discount: "0"

results:
  - name: classification_result
    kind: switch
    from: classify_event
```

This declaration has no `value`; its exposure is implementation-defined. Add an explicit, valid output mapping and consumer expression when the workflow needs portable data from the switch.

## Example 4: `from` Pointing at a Top-Level Workflow

When the emitting construct is itself a top-level workflow (not a nested step):

```yaml
workflows:
  - workflowId: batch_process
    type: loop
    items: $outputs.job_queue
    steps:
      - stepId: run_job
        operationRef: execute_job

results:
  - name: batch_results
    kind: loop
    from: batch_process       # top-level workflow, no ".stepId"
```

## Name Uniqueness

Result names MUST be unique within `results[]`. The validator catches duplicates:

```yaml
results:
  - name: my_result
    kind: merge
    from: wf1.step1
  - name: my_result     # ← duplicate
    kind: loop
    from: wf2
# error: results[1].name: duplicate result name "my_result"
```

## Validator: Kind/Type Mismatch

```yaml
workflows:
  - workflowId: my_loop
    type: loop
    items: $outputs.items
    steps:
      - stepId: process
        operationRef: do_work

results:
  - name: bad_result
    kind: merge          # ← wrong: loop construct, but kind is merge
    from: my_loop
# error: results[0].kind: kind "merge" does not match "my_loop" type "loop"
```

## From The Big Fixture

The large HCL fixture declares structural-result metadata from top-level
workflows. The excerpts below omit that fixture's `$workflows...` `value`
expressions; those are fixture-specific and not part of the portable UWS 1.11
expression grammar:

```hcl
result "decision.branch" {
  kind  = "switch"
  from  = "wf_switch"
}

result "containment.loop" {
  kind  = "loop"
  from  = "wf_loop"
}

result "merge.summary" {
  kind  = "merge"
  from  = "wf_merge"
}
```

Full context: [`testdata/big/big.hcl`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.hcl).

---

← [Triggers and Route Dispatch](04-Triggers-and-Route-Dispatch.md) | [Next: Success Criteria and Actions →](06-Success-Criteria-and-Actions.md)
