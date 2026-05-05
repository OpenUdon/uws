package uws1

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

func (o *Orchestrator) withRecordContext(ctx context.Context) context.Context {
	state, ok := ExecutionContextFromContext(ctx)
	if !ok || state == nil {
		state = &ExecutionContext{}
	}
	state.Iteration = cloneIteration(state.Iteration)
	state.Trigger = cloneTriggerContext(state.Trigger)
	state.Records = o.snapshotRecords()
	state.Current = cloneCurrentExecution(state.Current)
	return WithExecutionContext(ctx, state)
}

func (o *Orchestrator) withIterationContext(ctx context.Context, item any, index int, batch []any, batchIndex int) context.Context {
	state, ok := ExecutionContextFromContext(ctx)
	if !ok || state == nil {
		state = &ExecutionContext{}
	} else {
		state = &ExecutionContext{
			Iteration: cloneIteration(state.Iteration),
			Trigger:   cloneTriggerContext(state.Trigger),
			Records:   state.Records,
			Current:   cloneCurrentExecution(state.Current),
		}
	}
	state.Iteration = &IterationContext{
		Item:       item,
		Index:      index,
		Batch:      append([]any(nil), batch...),
		BatchIndex: batchIndex,
	}
	if state.Records == nil {
		state.Records = o.snapshotRecords()
	}
	return WithExecutionContext(ctx, state)
}

func cloneIteration(iteration *IterationContext) *IterationContext {
	if iteration == nil {
		return nil
	}
	return &IterationContext{
		Item:       iteration.Item,
		Index:      iteration.Index,
		Batch:      append([]any(nil), iteration.Batch...),
		BatchIndex: iteration.BatchIndex,
	}
}

func (o *Orchestrator) snapshotRecords() map[string]ExecutionRecord {
	o.mu.Lock()
	defer o.mu.Unlock()
	out := make(map[string]ExecutionRecord, len(o.records))
	for key, record := range o.records {
		out[key] = cloneExecutionRecord(record)
	}
	return out
}

func (o *Orchestrator) setRecord(key string, record ExecutionRecord) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.records[key] = cloneExecutionRecord(record)
}

func (o *Orchestrator) withCurrentExecutionContext(ctx context.Context, key, id, kind, responseID string, outputs map[string]any) context.Context {
	state, ok := ExecutionContextFromContext(ctx)
	if !ok || state == nil {
		state = &ExecutionContext{}
	} else {
		state = &ExecutionContext{
			Iteration: cloneIteration(state.Iteration),
			Trigger:   cloneTriggerContext(state.Trigger),
			Records:   state.Records,
			Current:   cloneCurrentExecution(state.Current),
		}
	}
	state.Current = &CurrentExecutionContext{
		Key:        key,
		ID:         id,
		Kind:       kind,
		ResponseID: responseID,
		Outputs:    cloneMapAny(outputs),
	}
	if state.Records == nil {
		state.Records = o.snapshotRecords()
	}
	return WithExecutionContext(ctx, state)
}

func (o *Orchestrator) resolveOutputs(ctx context.Context, key, id, kind, responseID string, definitions map[string]string) (map[string]any, error) {
	if len(definitions) == 0 {
		return nil, nil
	}
	names := make([]string, 0, len(definitions))
	for name := range definitions {
		if strings.TrimSpace(name) == "" {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	if len(names) == 0 {
		return nil, nil
	}
	resolved := make(map[string]any, len(names))
	for _, name := range names {
		expr := strings.TrimSpace(definitions[name])
		if expr == "" {
			continue
		}
		exprCtx := o.withCurrentExecutionContext(ctx, key, id, kind, responseID, resolved)
		value, err := o.Runtime.EvaluateExpression(exprCtx, expr)
		if err != nil {
			return nil, fmt.Errorf("evaluating output %q for %s %q: %w", name, kind, id, err)
		}
		resolved[name] = value
	}
	return resolved, nil
}

func cloneMapAny(src map[string]any) map[string]any {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]any, len(src))
	for key, value := range src {
		out[key] = value
	}
	return out
}

func (o *Orchestrator) keyForContext(ctx context.Context, key string) string {
	state, ok := ExecutionContextFromContext(ctx)
	if !ok || state == nil || state.Iteration == nil {
		return key
	}
	return fmt.Sprintf("%s#iter:%d", key, state.Iteration.Index)
}
