package uws1

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Runtime defines the interface that specialized executors must implement
// to provide leaf operation execution and expression evaluation for a
// UWS document.
type Runtime interface {
	// ExecuteLeaf executes a single leaf operation.
	ExecuteLeaf(ctx context.Context, op *Operation) error

	// EvaluateExpression evaluates a UWS runtime expression against the
	// current execution context.
	EvaluateExpression(ctx context.Context, expr string) (any, error)

	// ResolveItems resolves the items/forEach expression for iterative constructs.
	ResolveItems(ctx context.Context, itemsExpr string) ([]any, error)
}

// Orchestrator provides the abstract orchestration logic for walking the
// workflow graph and managing structural state transitions.
type Orchestrator struct {
	Document *Document
	Runtime  Runtime

	opIndex           map[string]*Operation
	workflowIndex     map[string]*Workflow
	stepIndex         map[string]*Step
	topLevelStepIndex map[string]*Step
	parallelGroups    map[string][]string
	mu                sync.Mutex
	records           map[string]ExecutionRecord
	// recordKeysByBase indexes record keys by their base portion (the part
	// before any "#iter:N" suffix added by keyForContext). It exists to keep
	// recordKeysForDependencyLocked O(iterations-of-dep) instead of O(all
	// records). Maintain it from every write to o.records.
	recordKeysByBase map[string]map[string]struct{}
	// recordErrors caches the original error returned by a runnable so a
	// re-entrant dependency call receives the typed error rather than
	// errors.New(record.Error). Keyed identically to o.records.
	recordErrors map[string]error
	inFlight     map[string]chan struct{}
}

// NewOrchestrator creates a new Orchestrator for the given document and runtime.
func NewOrchestrator(doc *Document, runtime Runtime) *Orchestrator {
	o := &Orchestrator{
		Document:          doc,
		Runtime:           runtime,
		opIndex:           make(map[string]*Operation),
		workflowIndex:     make(map[string]*Workflow),
		stepIndex:         make(map[string]*Step),
		topLevelStepIndex: make(map[string]*Step),
		parallelGroups:    make(map[string][]string),
		records:           make(map[string]ExecutionRecord),
		recordKeysByBase:  make(map[string]map[string]struct{}),
		recordErrors:      make(map[string]error),
		inFlight:          make(map[string]chan struct{}),
	}
	if doc != nil {
		idx := buildDocumentIndex(doc, nil)
		for id, op := range idx.operations {
			o.opIndex[id] = op
		}
		for id, wf := range idx.workflows {
			o.workflowIndex[id] = wf
		}
		for id, step := range idx.steps {
			o.stepIndex[id] = step
		}
		for id, step := range idx.topLevelSteps {
			o.topLevelStepIndex[id] = step
		}
		for id, members := range idx.parallelGroupMembers {
			o.parallelGroups[id] = append([]string(nil), members...)
		}
	}
	return o
}

// Execute executes the main workflow of the document.
func (o *Orchestrator) Execute(ctx context.Context) error {
	if o.Document == nil {
		return nil
	}
	start := func(ctx context.Context) error {
		wf, err := o.entryWorkflow()
		if err != nil {
			return err
		}
		return o.ExecuteWorkflow(ctx, wf)
	}
	err := o.executeWithSignals(ctx, start)
	o.Document.setExecutionRecords(o.snapshotRecords())
	return err
}

func (o *Orchestrator) executeWithSignals(ctx context.Context, start func(context.Context) error) error {
	for {
		err := start(ctx)
		switch typed := err.(type) {
		case nil:
			return nil
		case *endSignal:
			return nil
		case *gotoSignal:
			target, targetErr := o.gotoTarget(ctx, typed)
			if targetErr != nil {
				return targetErr
			}
			start = target
		default:
			return err
		}
	}
}

// ExecuteWorkflow executes a structural workflow.
func (o *Orchestrator) ExecuteWorkflow(ctx context.Context, wf *Workflow) error {
	if wf == nil {
		return nil
	}
	return o.executeRunnable(ctx, runnableExecution{
		key: workflowKey(wf.WorkflowID), id: wf.WorkflowID,
		kind: "workflow:" + wf.Type, responseID: wf.WorkflowID,
		dependencies: wf.DependsOn, when: wf.When, forEach: wf.ForEach,
		timeout: wf.Timeout, outputs: wf.Outputs,
		run: func(ctx context.Context) error {
			return o.executeStructural(ctx, structuralExecution{
				typeName: wf.Type, dependencies: wf.DependsOn, steps: wf.Steps,
				cases: wf.Cases, defaultSteps: wf.Default, items: wf.Items,
				batchSize: wf.BatchSize, wait: wf.Wait, key: workflowKey(wf.WorkflowID),
				useDefaultAwaitTimeout: wf.Timeout == nil,
			})
		},
	})
}

// ExecuteStep executes a single step.
func (o *Orchestrator) ExecuteStep(ctx context.Context, step *Step) error {
	if step == nil {
		return nil
	}
	responseID := step.StepID
	if strings.TrimSpace(step.OperationRef) != "" {
		responseID = strings.TrimSpace(step.OperationRef)
	} else if strings.TrimSpace(step.Workflow) != "" {
		responseID = strings.TrimSpace(step.Workflow)
	}
	return o.executeRunnable(ctx, runnableExecution{
		key: stepKey(step.StepID), id: step.StepID,
		kind: "step:" + step.Type, responseID: responseID,
		dependencies: step.DependsOn, when: step.When, forEach: step.ForEach,
		timeout: step.Timeout, outputs: step.Outputs,
		run: func(ctx context.Context) error {
			if step.Inputs != nil {
				ctx = withInputsContext(ctx, step.Inputs)
			}
			if step.OperationRef != "" {
				return o.executeOperationByIDForStep(ctx, step.OperationRef, step.StepID)
			}
			if step.Workflow != "" {
				return o.executeWorkflowByID(ctx, step.Workflow)
			}
			return o.executeStructural(ctx, structuralExecution{
				typeName: step.Type, dependencies: step.DependsOn, steps: step.Steps,
				cases: step.Cases, defaultSteps: step.Default, items: step.Items,
				batchSize: step.BatchSize, wait: step.Wait, key: stepKey(step.StepID),
				useDefaultAwaitTimeout: step.Timeout == nil,
			})
		},
	})
}

func (o *Orchestrator) executeWorkflowByID(ctx context.Context, workflowID string) error {
	wf := o.workflowIndex[workflowID]
	if wf == nil {
		return fmt.Errorf("uws1: workflow %q not found", workflowID)
	}
	return o.ExecuteWorkflow(ctx, wf)
}

func (o *Orchestrator) executeOperationByID(ctx context.Context, operationID string) error {
	return o.executeOperationByIDWithKey(ctx, operationID, operationKey(operationID))
}

func (o *Orchestrator) executeOperationByIDForStep(ctx context.Context, operationID, stepID string) error {
	return o.executeOperationByIDWithKey(ctx, operationID, stepOperationKey(stepID, operationID))
}

func (o *Orchestrator) executeOperationByIDWithKey(ctx context.Context, operationID, key string) error {
	op := o.opIndex[operationID]
	if op == nil {
		return fmt.Errorf("uws1: operation %q not found", operationID)
	}
	return o.executeRunnable(ctx, runnableExecution{
		key: key, id: op.OperationID, kind: "operation", responseID: op.OperationID,
		dependencies: op.DependsOn, when: op.When, forEach: op.ForEach,
		outputs: op.Outputs,
		run: func(ctx context.Context) error {
			return o.executeOperation(ctx, op, key)
		},
	})
}
