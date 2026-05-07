package uws1

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math"
	"sort"
	"time"
)

const defaultAwaitPollInterval = 200 * time.Millisecond

// AwaitTimeoutError reports that an await construct exceeded the executor's
// configured internal timeout.
type AwaitTimeoutError struct {
	Timeout time.Duration
}

func (e *AwaitTimeoutError) Error() string {
	if e == nil || e.Timeout <= 0 {
		return "uws1: await timed out"
	}
	return fmt.Sprintf("uws1: await timed out after %s", e.Timeout)
}

func (e *AwaitTimeoutError) Is(target error) bool {
	_, ok := target.(*AwaitTimeoutError)
	return ok
}

func (o *Orchestrator) executeStructural(ctx context.Context, typeName string, deps []string, steps []*Step, cases []*Case, defaultSteps []*Step, itemsExpr, mode, batchSizeExpr, waitExpr, key string) error {
	switch typeName {
	case WorkflowTypeSequence:
		return o.executeSteps(ctx, steps)
	case WorkflowTypeParallel:
		return o.executeStepsParallel(ctx, steps)
	case WorkflowTypeSwitch:
		return o.executeSwitch(ctx, cases, defaultSteps)
	case WorkflowTypeMerge:
		return o.executeMerge(ctx, deps, key)
	case WorkflowTypeLoop:
		return o.executeLoop(ctx, steps, itemsExpr, batchSizeExpr, key)
	case WorkflowTypeAwait:
		return o.executeAwait(ctx, steps, waitExpr)
	default:
		return fmt.Errorf("uws1: unsupported workflow type %q", typeName)
	}
}

// executeSteps executes a list of steps sequentially.
func (o *Orchestrator) executeSteps(ctx context.Context, steps []*Step) error {
	for _, step := range steps {
		if err := o.ExecuteStep(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

// executeStepsParallel executes a list of steps concurrently.
//
// Control signals (*endSignal, *gotoSignal) raised by a parallel branch are
// converted into errors. Control-flow semantics across siblings would be
// undefined — there is no obvious answer to "goto X happens while sibling Y
// is mid-flight" — so we explicitly reject them rather than letting the
// errgroup race decide which signal (or real error) survives.
func (o *Orchestrator) executeStepsParallel(ctx context.Context, steps []*Step) error {
	group, groupCtx := errgroup.WithContext(ctx)
	for _, step := range steps {
		step := step
		group.Go(func() error {
			err := o.ExecuteStep(groupCtx, step)
			if isControlSignal(err) {
				return fmt.Errorf("uws1: control signal %q is not allowed inside a parallel workflow", err.Error())
			}
			return err
		})
	}
	return group.Wait()
}

// executeSwitch executes a switch construct.
func (o *Orchestrator) executeSwitch(ctx context.Context, cases []*Case, defaultSteps []*Step) error {
	for _, c := range cases {
		if c == nil {
			continue
		}
		if c.When == "" {
			return o.executeSteps(ctx, c.Steps)
		}
		matched, err := o.evaluateTruthy(ctx, c.When)
		if err != nil {
			return fmt.Errorf("evaluating switch case %q condition: %w", c.Name, err)
		}
		if matched {
			return o.executeSteps(ctx, c.Steps)
		}
	}
	if len(defaultSteps) > 0 {
		return o.executeSteps(ctx, defaultSteps)
	}
	return nil
}

func (o *Orchestrator) executeMerge(ctx context.Context, deps []string, key string) error {
	result := o.mergeDependencyRecords(deps)
	o.mu.Lock()
	record := o.records[key]
	record.Result = result
	record.Status = "success"
	o.writeRecordLocked(key, record)
	o.mu.Unlock()
	return nil
}

func (o *Orchestrator) mergeDependencyRecords(deps []string) []map[string]any {
	if len(deps) == 0 {
		return nil
	}
	ordered := make([]string, 0, len(deps))
	seen := make(map[string]struct{})
	for _, dep := range deps {
		if dep == "" {
			continue
		}
		if members := o.parallelGroups[dep]; len(members) > 0 {
			for _, member := range members {
				if _, ok := seen[member]; ok {
					continue
				}
				ordered = append(ordered, member)
				seen[member] = struct{}{}
			}
			continue
		}
		if _, ok := seen[dep]; ok {
			continue
		}
		ordered = append(ordered, dep)
		seen[dep] = struct{}{}
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	out := make([]map[string]any, 0, len(ordered))
	for _, dep := range ordered {
		keys := o.recordKeysForDependencyLocked(dep)
		for _, key := range keys {
			record := o.records[key]
			out = append(out, map[string]any{
				"id":      record.ID,
				"kind":    record.Kind,
				"status":  record.Status,
				"error":   record.Error,
				"result":  record.Result,
				"outputs": record.Outputs,
			})
		}
	}
	return out
}

// recordKeysForDependencyLocked returns the record keys associated with a
// single dependency name, including any per-iteration suffixes added by
// keyForContext. It uses o.recordKeysByBase rather than scanning all records.
//
// Precondition: callers must have already executed (and waited for) the
// dependency. Merge constructs reach this through executeDependencies →
// executeOnce, whose inFlight channel ensures the depended-on runnable's
// close(ch) fires before the merge proceeds; there is no race against
// concurrently writing branches.
func (o *Orchestrator) recordKeysForDependencyLocked(dep string) []string {
	var base string
	switch {
	case o.stepIndex[dep] != nil:
		base = stepKey(dep)
	case o.workflowIndex[dep] != nil:
		base = workflowKey(dep)
	case o.opIndex[dep] != nil:
		base = operationKey(dep)
	default:
		return nil
	}
	set := o.recordKeysByBase[base]
	if len(set) == 0 {
		return nil
	}
	matches := make([]string, 0, len(set))
	for key := range set {
		matches = append(matches, key)
	}
	sort.Strings(matches)
	return matches
}

// executeLoop executes a loop construct.
func (o *Orchestrator) executeLoop(ctx context.Context, steps []*Step, itemsExpr, batchSizeExpr, key string) error {
	items, err := o.Runtime.ResolveItems(ctx, itemsExpr)
	if err != nil {
		return fmt.Errorf("resolving loop items: %w", err)
	}
	batchSize, err := o.resolveBatchSize(ctx, batchSizeExpr)
	if err != nil {
		return err
	}
	if batchSize <= 0 {
		batchSize = len(items)
		if batchSize == 0 {
			batchSize = 1
		}
	}

	var results []map[string]any
	for batchIndex, start := 0, 0; start < len(items); batchIndex, start = batchIndex+1, start+batchSize {
		end := start + batchSize
		if end > len(items) {
			end = len(items)
		}
		batch := append([]any(nil), items[start:end]...)
		for i, item := range batch {
			itemCtx := o.withIterationContext(ctx, item, start+i, batch, batchIndex)
			if err := o.executeSteps(itemCtx, steps); err != nil {
				return err
			}
			results = append(results, map[string]any{
				"index":      start + i,
				"batchIndex": batchIndex,
				"item":       item,
			})
		}
	}
	o.mu.Lock()
	record := o.records[key]
	record.Result = results
	record.Status = "success"
	o.writeRecordLocked(key, record)
	o.mu.Unlock()
	return nil
}

func (o *Orchestrator) resolveBatchSize(ctx context.Context, batchSizeExpr string) (int, error) {
	if batchSizeExpr == "" {
		return 0, nil
	}
	value, err := o.Runtime.EvaluateExpression(ctx, batchSizeExpr)
	if err != nil {
		return 0, fmt.Errorf("evaluating batchSize: %w", err)
	}
	switch typed := value.(type) {
	case int:
		if typed <= 0 {
			return 0, fmt.Errorf("batchSize must resolve to a positive integer")
		}
		return typed, nil
	case int64:
		if typed <= 0 {
			return 0, fmt.Errorf("batchSize must resolve to a positive integer")
		}
		return int(typed), nil
	case float64:
		if typed <= 0 || math.Trunc(typed) != typed {
			return 0, fmt.Errorf("batchSize must resolve to a positive integer")
		}
		return int(typed), nil
	default:
		return 0, fmt.Errorf("batchSize must resolve to a positive integer")
	}
}

func (o *Orchestrator) executeAwait(ctx context.Context, steps []*Step, waitExpr string) error {
	pollInterval := defaultAwaitPollInterval
	if o != nil && o.Document != nil && o.Document.ExecutionOptions.AwaitPollInterval > 0 {
		pollInterval = o.Document.ExecutionOptions.AwaitPollInterval
	}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	var timeout <-chan time.Time
	if o != nil && o.Document != nil && o.Document.ExecutionOptions.AwaitTimeout > 0 {
		timer := time.NewTimer(o.Document.ExecutionOptions.AwaitTimeout)
		defer timer.Stop()
		timeout = timer.C
	}
	for {
		ok, err := o.evaluateTruthy(ctx, waitExpr)
		if err != nil {
			return fmt.Errorf("evaluating await wait expression: %w", err)
		}
		if ok {
			return o.executeSteps(ctx, steps)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return &AwaitTimeoutError{Timeout: o.Document.ExecutionOptions.AwaitTimeout}
		case <-ticker.C:
		}
	}
}
