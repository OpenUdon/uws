package uws1

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type runnableExecution struct {
	key             string
	id              string
	kind            string
	responseID      string
	dependencies    []string
	dependencyScope string
	outputsScope    string
	when            string
	forEach         string
	wait            string
	timeout         *float64
	outputs         map[string]string
	run             func(context.Context) error
}

type forEachExecution struct {
	execKey      string
	baseKey      string
	id           string
	kind         string
	responseID   string
	expression   string
	outputs      map[string]string
	outputsScope string
	run          func(context.Context) error
}

func (o *Orchestrator) executeRunnable(ctx context.Context, spec runnableExecution) error {
	execKey := o.keyForContext(ctx, spec.key)
	return o.executeOnce(ctx, execKey, spec.id, spec.kind, func(ctx context.Context) error {
		dependencyCtx := ctx
		if spec.dependencyScope != "" {
			dependencyCtx = withWorkflowScope(dependencyCtx, spec.dependencyScope)
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := o.executeDependencies(dependencyCtx, spec.dependencies); err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		shouldRun, err := o.evaluateWhen(dependencyCtx, spec.when, execKey, spec.id, spec.kind)
		if err != nil {
			return err
		}
		if !shouldRun {
			return nil
		}
		return executeWithTimeout(ctx, spec.timeout, func(runCtx context.Context) error {
			if spec.wait != "" && !strings.HasSuffix(spec.kind, ":await") && o.Document != nil && supportsUWSVersionAtLeast(o.Document.UWS, 1, 10, 0) {
				if err := o.waitBeforeRun(runCtx, spec.wait, spec.id); err != nil {
					return err
				}
			}
			if spec.forEach != "" {
				return o.executeForEach(runCtx, forEachExecution{
					execKey: execKey, baseKey: spec.key, id: spec.id, kind: spec.kind,
					responseID: spec.responseID, expression: spec.forEach,
					outputs: spec.outputs, outputsScope: spec.outputsScope, run: spec.run,
				})
			}
			if err := spec.run(runCtx); err != nil {
				return err
			}
			outputsCtx := runCtx
			if spec.outputsScope != "" {
				outputsCtx = withWorkflowScope(outputsCtx, spec.outputsScope)
			}
			return o.finalizeOutputs(outputsCtx, execKey, spec.id, spec.kind, spec.responseID, spec.outputs)
		})
	})
}

const maxWaitSeconds = 86400

func (o *Orchestrator) waitBeforeRun(ctx context.Context, expression, id string) error {
	value, err := o.Runtime.EvaluateExpression(ctx, expression)
	if err != nil {
		return fmt.Errorf("evaluating wait duration for %q: %w", id, err)
	}
	seconds, ok := waitSecondsNumber(value)
	if !ok || math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 || seconds > maxWaitSeconds {
		return fmt.Errorf("wait duration for %q must resolve to a finite number from 0 through %d seconds", id, maxWaitSeconds)
	}
	if seconds == 0 {
		return ctx.Err()
	}
	nanoseconds := math.Round(seconds * float64(time.Second))
	if nanoseconds < 1 {
		return ctx.Err()
	}
	timer := time.NewTimer(time.Duration(nanoseconds))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return ctx.Err()
	}
}

func waitSecondsNumber(value any) (float64, bool) {
	switch typed := value.(type) {
	case json.Number:
		seconds, err := typed.Float64()
		return seconds, err == nil
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case int:
		return float64(typed), true
	case int8:
		return float64(typed), true
	case int16:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case uint:
		return float64(typed), true
	case uint8:
		return float64(typed), true
	case uint16:
		return float64(typed), true
	case uint32:
		return float64(typed), true
	case uint64:
		return float64(typed), true
	case uintptr:
		return float64(typed), true
	default:
		return 0, false
	}
}

func executeWithTimeout(ctx context.Context, timeout *float64, run func(context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if timeout == nil {
		return run(ctx)
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, timeoutDuration(*timeout))
	defer cancel()
	err := run(timeoutCtx)
	if timeoutErr := timeoutCtx.Err(); timeoutErr != nil {
		return timeoutErr
	}
	return err
}

func timeoutDuration(seconds float64) time.Duration {
	const maxDuration = time.Duration(1<<63 - 1)
	nanoseconds := seconds * float64(time.Second)
	if nanoseconds >= float64(maxDuration) {
		return maxDuration
	}
	return time.Duration(nanoseconds)
}

// evaluateWhen evaluates a when-expression. Returns false if the runnable
// should be skipped (and writes a "skipped" record). An empty expression is
// treated as truthy.
func (o *Orchestrator) evaluateWhen(ctx context.Context, whenExpr, execKey, id, kind string) (bool, error) {
	if whenExpr == "" {
		return true, nil
	}
	ok, err := o.evaluateTruthy(ctx, whenExpr)
	if err != nil {
		return false, fmt.Errorf("evaluating when condition for %q: %w", id, err)
	}
	if !ok {
		o.setRecord(execKey, ExecutionRecord{ID: id, Kind: kind, Status: "skipped"})
		return false, nil
	}
	return true, nil
}

// finalizeOutputs resolves and stores outputs on the record. No-op when there
// are no output definitions.
func (o *Orchestrator) finalizeOutputs(ctx context.Context, execKey, id, kind, responseID string, outputs map[string]string) error {
	if len(outputs) == 0 {
		return nil
	}
	outputsCtx := o.withRecordContext(ctx)
	resolved, err := o.resolveOutputs(outputsCtx, execKey, id, kind, responseID, outputs)
	if err != nil {
		return err
	}
	o.mu.Lock()
	record := o.records[execKey]
	record.Outputs = resolved
	o.writeRecordLocked(execKey, record)
	o.mu.Unlock()
	return nil
}

// executeForEach iterates a runnable over each item resolved from forEachExpr,
// resolving per-iteration outputs and aggregating them under the parent record.
func (o *Orchestrator) executeForEach(ctx context.Context, spec forEachExecution) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	items, err := o.Runtime.ResolveItems(ctx, spec.expression)
	if err != nil {
		return fmt.Errorf("resolving forEach for %q: %w", spec.id, err)
	}
	iterationResults := make([]map[string]any, 0, len(items))
	aggregatedOutputs := make(map[string][]any)
	for index, item := range items {
		if err := ctx.Err(); err != nil {
			return err
		}
		itemCtx := o.withIterationContext(ctx, item, index, nil, -1)
		itemKey := o.keyForContext(itemCtx, spec.baseKey)
		o.setRecord(itemKey, ExecutionRecord{ID: spec.id, Kind: spec.kind, Status: "running"})
		if err := spec.run(itemCtx); err != nil {
			if isControlSignal(err) {
				o.setRecord(itemKey, ExecutionRecord{ID: spec.id, Kind: spec.kind, Status: "success"})
				return err
			}
			o.setRecord(itemKey, ExecutionRecord{ID: spec.id, Kind: spec.kind, Status: "error", Error: err.Error()})
			return err
		}
		var resolved map[string]any
		if len(spec.outputs) > 0 {
			outputsCtx := itemCtx
			if spec.outputsScope != "" {
				outputsCtx = withWorkflowScope(outputsCtx, spec.outputsScope)
			}
			outputsCtx = o.withRecordContext(outputsCtx)
			resolved, err = o.resolveOutputs(outputsCtx, itemKey, spec.id, spec.kind, spec.responseID, spec.outputs)
			if err != nil {
				o.setRecord(itemKey, ExecutionRecord{ID: spec.id, Kind: spec.kind, Status: "error", Error: err.Error()})
				return err
			}
		}
		o.mu.Lock()
		record := o.records[itemKey]
		if record.Status == "running" {
			record.Status = "success"
		}
		if len(resolved) > 0 {
			record.Outputs = resolved
		}
		o.writeRecordLocked(itemKey, record)
		o.mu.Unlock()
		iterationResults = append(iterationResults, map[string]any{
			"index":   index,
			"item":    item,
			"status":  record.Status,
			"error":   record.Error,
			"result":  record.Result,
			"outputs": cloneMapAny(record.Outputs),
		})
		for name, value := range resolved {
			aggregatedOutputs[name] = append(aggregatedOutputs[name], value)
		}
	}
	o.mu.Lock()
	record := o.records[spec.execKey]
	record.Result = iterationResults
	record.Status = "success"
	if len(aggregatedOutputs) > 0 {
		record.Outputs = make(map[string]any, len(aggregatedOutputs))
		for name, values := range aggregatedOutputs {
			record.Outputs[name] = append([]any(nil), values...)
		}
	}
	o.writeRecordLocked(spec.execKey, record)
	o.mu.Unlock()
	return nil
}

func (o *Orchestrator) executeOnce(ctx context.Context, key, id, kind string, run func(context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	o.mu.Lock()
	if record, ok := o.records[key]; ok && record.Status != "running" {
		cachedErr := o.recordErrors[key]
		o.mu.Unlock()
		return replayedRunnableError(record, cachedErr)
	}
	if ch, ok := o.inFlight[key]; ok {
		o.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
		}
		o.mu.Lock()
		record := o.records[key]
		cachedErr := o.recordErrors[key]
		o.mu.Unlock()
		return replayedRunnableError(record, cachedErr)
	}
	ch := make(chan struct{})
	o.inFlight[key] = ch
	o.writeRecordLocked(key, ExecutionRecord{ID: id, Kind: kind, Status: "running"})
	o.mu.Unlock()

	err := run(o.withRecordContext(ctx))

	o.mu.Lock()
	record := o.records[key]
	switch {
	case err == nil:
		if record.Status == "running" {
			record.Status = "success"
		}
	case isControlSignal(err):
		if record.Status == "running" {
			record.Status = "success"
		}
	default:
		record.Status = "error"
		record.Error = err.Error()
		// Cache the original error so re-entrant dependency calls can return
		// the typed error rather than a generic errors.New(record.Error).
		o.recordErrors[key] = err
	}
	o.writeRecordLocked(key, record)
	delete(o.inFlight, key)
	close(ch)
	o.mu.Unlock()

	return err
}

// replayedRunnableError returns nil if a previously-completed runnable
// succeeded (or terminated via a control signal), the cached typed error if
// one was captured, or a generic errors.New fallback when the cache is empty
// (e.g. records loaded from outside this orchestrator instance).
func replayedRunnableError(record ExecutionRecord, cached error) error {
	if record.Status != "error" {
		return nil
	}
	if cached != nil {
		return cached
	}
	if record.Error != "" {
		return errors.New(record.Error)
	}
	return nil
}

func (o *Orchestrator) executeDependencies(ctx context.Context, deps []string) error {
	for _, dep := range deps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if dep == "" {
			continue
		}
		if members := o.parallelGroups[dep]; len(members) > 0 {
			for _, member := range members {
				if err := o.executeDependency(ctx, member); err != nil {
					return err
				}
			}
			continue
		}
		if err := o.executeDependency(ctx, dep); err != nil {
			return err
		}
	}
	return nil
}

func (o *Orchestrator) executeDependency(ctx context.Context, name string) error {
	if step := o.stepIndex[name]; step != nil {
		return o.ExecuteStep(ctx, step)
	}
	if wf := o.workflowIndex[name]; wf != nil {
		return o.ExecuteWorkflow(ctx, wf)
	}
	if op := o.opIndex[name]; op != nil {
		if satisfied, err := o.waitForExistingOperationInvocation(ctx, op.OperationID); satisfied || err != nil {
			return err
		}
		return o.executeOperationByID(ctx, op.OperationID)
	}
	return fmt.Errorf("uws1: unknown dependency %q", name)
}

func (o *Orchestrator) waitForExistingOperationInvocation(ctx context.Context, operationID string) (bool, error) {
	for {
		o.mu.Lock()
		keys := o.operationInvocationKeysLocked(operationID, workflowScopeFromContext(ctx))
		if len(keys) == 0 {
			o.mu.Unlock()
			return false, nil
		}
		chans := make([]chan struct{}, 0, len(keys))
		for _, key := range keys {
			if ch := o.inFlight[key]; ch != nil {
				chans = append(chans, ch)
			}
		}
		if len(chans) == 0 {
			for _, key := range keys {
				record := o.records[key]
				cachedErr := o.recordErrors[key]
				if err := replayedRunnableError(record, cachedErr); err != nil {
					o.mu.Unlock()
					return true, err
				}
			}
			o.mu.Unlock()
			return true, nil
		}
		o.mu.Unlock()

		for _, ch := range chans {
			select {
			case <-ctx.Done():
				return true, ctx.Err()
			case <-ch:
			}
		}
	}
}

func (o *Orchestrator) entryWorkflow() (*Workflow, error) {
	return requireExecutableEntryWorkflow(o.Document)
}

func (o *Orchestrator) evaluateTruthy(ctx context.Context, expr string) (bool, error) {
	value, err := o.Runtime.EvaluateExpression(ctx, expr)
	if err != nil {
		return false, err
	}
	return o.truthy(value)
}

func workflowCallKey(id, callerKey string) string {
	hash := sha256.Sum256([]byte(callerKey))
	return workflowKey(id) + "#call:" + hex.EncodeToString(hash[:])
}

func operationKey(id string) string { return "op:" + id }
func workflowKey(id string) string  { return "wf:" + id }
func stepKey(id string) string      { return "step:" + id }
func stepOperationKey(stepID, operationID string) string {
	return "stepop:" + stepID + ":" + operationID
}

// compositeKey returns the per-iteration key for a runnable executing inside a
// forEach/loop frame. Iter < 0 means "no iteration suffix".
func compositeKey(base string, iter int) string {
	if iter < 0 {
		return base
	}
	return fmt.Sprintf("%s#iter:%d", base, iter)
}

func compositeIterationKey(base string, path []int) string {
	if len(path) == 0 {
		return base
	}
	if len(path) == 1 {
		return compositeKey(base, path[0])
	}
	parts := make([]string, len(path))
	for i, index := range path {
		parts[i] = fmt.Sprintf("%d", index)
	}
	return base + "#iter:" + strings.Join(parts, ".")
}

// baseFromCompositeKey strips a "#iter:N" suffix; if there is none it returns
// the input. Used by recordKeysForDependencyLocked / setRecord to maintain the
// recordKeysByBase index without re-deriving the format from string slicing.
func baseFromCompositeKey(key string) string {
	if i := strings.Index(key, "#iter:"); i >= 0 {
		return key[:i]
	}
	return key
}
