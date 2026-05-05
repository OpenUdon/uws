package convert

import (
	"fmt"

	"github.com/tabilet/uws/uws1"
)

type documentWalkHandlers struct {
	description func(value *string)
	dynamicMap  func(path string, value *map[string]any) error
	extensions  func(path string, extensions map[string]any) error
	paramSchema func(path string, schema *uws1.ParamSchema) error
}

func walkDocument(doc *uws1.Document, h documentWalkHandlers) error {
	if doc == nil {
		return nil
	}
	if err := walkExtensions("document", doc.Extensions, h); err != nil {
		return err
	}
	if err := walkDynamicMap("variables", &doc.Variables, h); err != nil {
		return err
	}
	if doc.Info != nil {
		if err := walkExtensions("info", doc.Info.Extensions, h); err != nil {
			return err
		}
		walkDescription(&doc.Info.Description, h)
		walkDescription(&doc.Info.Summary, h)
	}
	for i, source := range doc.SourceDescriptions {
		if source == nil {
			continue
		}
		if err := walkExtensions(fmt.Sprintf("sourceDescriptions[%d]", i), source.Extensions, h); err != nil {
			return err
		}
	}
	for i, op := range doc.Operations {
		if err := walkOperation(fmt.Sprintf("operations[%d]", i), op, h); err != nil {
			return err
		}
	}
	for i, wf := range doc.Workflows {
		if err := walkWorkflow(fmt.Sprintf("workflows[%d]", i), wf, h); err != nil {
			return err
		}
	}
	for i, trigger := range doc.Triggers {
		if err := walkTrigger(fmt.Sprintf("triggers[%d]", i), trigger, h); err != nil {
			return err
		}
	}
	for i, result := range doc.Results {
		if result == nil {
			continue
		}
		if err := walkExtensions(fmt.Sprintf("results[%d]", i), result.Extensions, h); err != nil {
			return err
		}
	}
	if doc.Components != nil {
		if err := walkExtensions("components", doc.Components.Extensions, h); err != nil {
			return err
		}
		if err := walkDynamicMap("components.variables", &doc.Components.Variables, h); err != nil {
			return err
		}
	}
	return nil
}

func walkOperation(path string, op *uws1.Operation, h documentWalkHandlers) error {
	if op == nil {
		return nil
	}
	if err := walkExtensions(path, op.Extensions, h); err != nil {
		return err
	}
	walkDescription(&op.Description, h)
	if err := walkDynamicMap(path+".request", &op.Request, h); err != nil {
		return err
	}
	for i, criterion := range op.SuccessCriteria {
		if err := walkCriterion(fmt.Sprintf("%s.successCriteria[%d]", path, i), criterion, h); err != nil {
			return err
		}
	}
	for i, action := range op.OnFailure {
		if err := walkFailureAction(fmt.Sprintf("%s.onFailure[%d]", path, i), action, h); err != nil {
			return err
		}
	}
	for i, action := range op.OnSuccess {
		if err := walkSuccessAction(fmt.Sprintf("%s.onSuccess[%d]", path, i), action, h); err != nil {
			return err
		}
	}
	return nil
}

func walkWorkflow(path string, wf *uws1.Workflow, h documentWalkHandlers) error {
	if wf == nil {
		return nil
	}
	if err := walkExtensions(path, wf.Extensions, h); err != nil {
		return err
	}
	walkDescription(&wf.Description, h)
	if err := walkParamSchema(path+".inputs", wf.Inputs, h); err != nil {
		return err
	}
	if err := walkSteps(path+".steps", wf.Steps, h); err != nil {
		return err
	}
	if err := walkCases(path+".cases", wf.Cases, h); err != nil {
		return err
	}
	return walkSteps(path+".default", wf.Default, h)
}

func walkSteps(path string, steps []*uws1.Step, h documentWalkHandlers) error {
	for i, step := range steps {
		if err := walkStep(fmt.Sprintf("%s[%d]", path, i), step, h); err != nil {
			return err
		}
	}
	return nil
}

func walkStep(path string, step *uws1.Step, h documentWalkHandlers) error {
	if step == nil {
		return nil
	}
	if err := walkExtensions(path, step.Extensions, h); err != nil {
		return err
	}
	walkDescription(&step.Description, h)
	if err := walkDynamicMap(path+".body", &step.Body, h); err != nil {
		return err
	}
	if err := walkSteps(path+".steps", step.Steps, h); err != nil {
		return err
	}
	if err := walkCases(path+".cases", step.Cases, h); err != nil {
		return err
	}
	return walkSteps(path+".default", step.Default, h)
}

func walkCases(path string, cases []*uws1.Case, h documentWalkHandlers) error {
	for i, c := range cases {
		if err := walkCase(fmt.Sprintf("%s[%d]", path, i), c, h); err != nil {
			return err
		}
	}
	return nil
}

func walkCase(path string, c *uws1.Case, h documentWalkHandlers) error {
	if c == nil {
		return nil
	}
	if err := walkExtensions(path, c.Extensions, h); err != nil {
		return err
	}
	if err := walkDynamicMap(path+".body", &c.Body, h); err != nil {
		return err
	}
	return walkSteps(path+".steps", c.Steps, h)
}

func walkTrigger(path string, trigger *uws1.Trigger, h documentWalkHandlers) error {
	if trigger == nil {
		return nil
	}
	if err := walkExtensions(path, trigger.Extensions, h); err != nil {
		return err
	}
	if err := walkDynamicMap(path+".options", &trigger.Options, h); err != nil {
		return err
	}
	for i, route := range trigger.Routes {
		if route == nil {
			continue
		}
		if err := walkExtensions(fmt.Sprintf("%s.routes[%d]", path, i), route.Extensions, h); err != nil {
			return err
		}
	}
	return nil
}

func walkCriterion(path string, criterion *uws1.Criterion, h documentWalkHandlers) error {
	if criterion == nil {
		return nil
	}
	return walkExtensions(path, criterion.Extensions, h)
}

func walkFailureAction(path string, action *uws1.FailureAction, h documentWalkHandlers) error {
	if action == nil {
		return nil
	}
	if err := walkExtensions(path, action.Extensions, h); err != nil {
		return err
	}
	for i, criterion := range action.Criteria {
		if err := walkCriterion(fmt.Sprintf("%s.criteria[%d]", path, i), criterion, h); err != nil {
			return err
		}
	}
	return nil
}

func walkSuccessAction(path string, action *uws1.SuccessAction, h documentWalkHandlers) error {
	if action == nil {
		return nil
	}
	if err := walkExtensions(path, action.Extensions, h); err != nil {
		return err
	}
	for i, criterion := range action.Criteria {
		if err := walkCriterion(fmt.Sprintf("%s.criteria[%d]", path, i), criterion, h); err != nil {
			return err
		}
	}
	return nil
}

func walkParamSchema(path string, schema *uws1.ParamSchema, h documentWalkHandlers) error {
	if schema == nil {
		return nil
	}
	if err := walkExtensions(path, schema.Extensions, h); err != nil {
		return err
	}
	if h.paramSchema != nil {
		if err := h.paramSchema(path, schema); err != nil {
			return err
		}
	}
	for name, child := range schema.Properties {
		if err := walkParamSchema(path+".properties."+name, child, h); err != nil {
			return err
		}
	}
	if err := walkParamSchema(path+".items", schema.Items, h); err != nil {
		return err
	}
	for i, child := range schema.AllOf {
		if err := walkParamSchema(fmt.Sprintf("%s.allOf[%d]", path, i), child, h); err != nil {
			return err
		}
	}
	for i, child := range schema.OneOf {
		if err := walkParamSchema(fmt.Sprintf("%s.oneOf[%d]", path, i), child, h); err != nil {
			return err
		}
	}
	for i, child := range schema.AnyOf {
		if err := walkParamSchema(fmt.Sprintf("%s.anyOf[%d]", path, i), child, h); err != nil {
			return err
		}
	}
	return nil
}

func walkDescription(value *string, h documentWalkHandlers) {
	if h.description != nil {
		h.description(value)
	}
}

func walkDynamicMap(path string, value *map[string]any, h documentWalkHandlers) error {
	if h.dynamicMap != nil && *value != nil {
		return h.dynamicMap(path, value)
	}
	return nil
}

func walkExtensions(path string, extensions map[string]any, h documentWalkHandlers) error {
	if h.extensions != nil {
		return h.extensions(path, extensions)
	}
	return nil
}
