package uws1

import "context"

type executionContextKey struct{}
type iterationPathContextKey struct{}

// ExecutionContext carries runtime-only orchestration state into runtime hooks.
type ExecutionContext struct {
	Iteration *IterationContext
	Trigger   *TriggerExecutionContext
	Inputs    map[string]any
	Records   map[string]ExecutionRecord
	Current   *CurrentExecutionContext
}

// IterationContext describes the current orchestrator-owned iteration scope.
type IterationContext struct {
	Item       any
	Index      int
	Batch      []any
	BatchIndex int
}

// TriggerExecutionContext describes the trigger event currently being routed.
type TriggerExecutionContext struct {
	ID         string
	Output     int
	OutputName string
	Payload    any
}

// ExecutionRecord is the orchestrator-owned summary of one construct execution.
type ExecutionRecord struct {
	ID      string
	Kind    string
	Status  string
	Error   string
	Result  any
	Outputs map[string]any
}

// CurrentExecutionContext describes the construct currently being evaluated.
type CurrentExecutionContext struct {
	Key        string
	ID         string
	Kind       string
	ResponseID string
	Outputs    map[string]any
}

// WithExecutionContext returns a new context carrying the given execution state.
func WithExecutionContext(ctx context.Context, state *ExecutionContext) context.Context {
	return context.WithValue(ctx, executionContextKey{}, state)
}

// ExecutionContextFromContext returns the current execution state, if any.
func ExecutionContextFromContext(ctx context.Context) (*ExecutionContext, bool) {
	state, ok := ctx.Value(executionContextKey{}).(*ExecutionContext)
	return state, ok
}

func cloneTriggerContext(trigger *TriggerExecutionContext) *TriggerExecutionContext {
	if trigger == nil {
		return nil
	}
	return &TriggerExecutionContext{
		ID:         trigger.ID,
		Output:     trigger.Output,
		OutputName: trigger.OutputName,
		Payload:    trigger.Payload,
	}
}

// cloneExecutionContext returns an independent execution-state wrapper for a
// child context. Values carried inside Item, Payload, and record Results remain
// runtime-owned; all orchestrator-owned pointers and maps are copied.
func cloneExecutionContext(state *ExecutionContext) *ExecutionContext {
	if state == nil {
		return &ExecutionContext{}
	}
	return &ExecutionContext{
		Iteration: cloneIteration(state.Iteration),
		Trigger:   cloneTriggerContext(state.Trigger),
		Inputs:    cloneInputs(state.Inputs),
		Records:   cloneExecutionRecords(state.Records),
		Current:   cloneCurrentExecution(state.Current),
	}
}

func iterationPathFromContext(ctx context.Context) []int {
	path, _ := ctx.Value(iterationPathContextKey{}).([]int)
	return append([]int(nil), path...)
}

func withIterationPath(ctx context.Context, path []int) context.Context {
	return context.WithValue(ctx, iterationPathContextKey{}, append([]int(nil), path...))
}
