package uws1

import "fmt"

type documentHCLWalkHandlers struct {
	description func(value *string)
	dynamicMap  func(path string, value *map[string]any) error
	extensions  func(path string, extensions map[string]any) error
	paramSchema func(path string, schema *ParamSchema) error
}

func walkDocumentHCL(doc *Document, h documentHCLWalkHandlers) error {
	if doc == nil {
		return nil
	}
	if err := walkHCLExtensions("document", doc.Extensions, h); err != nil {
		return err
	}
	if err := walkHCLDynamicMap("variables", &doc.Variables, h); err != nil {
		return err
	}
	if doc.Info != nil {
		if err := walkHCLExtensions("info", doc.Info.Extensions, h); err != nil {
			return err
		}
		walkHCLDescription(&doc.Info.Description, h)
		walkHCLDescription(&doc.Info.Summary, h)
	}
	for i, source := range doc.SourceDescriptions {
		if source == nil {
			continue
		}
		if err := walkHCLExtensions(fmt.Sprintf("sourceDescriptions[%d]", i), source.Extensions, h); err != nil {
			return err
		}
	}
	for i, op := range doc.Operations {
		if err := walkHCLOperation(fmt.Sprintf("operations[%d]", i), op, h); err != nil {
			return err
		}
	}
	for i, wf := range doc.Workflows {
		if err := walkHCLWorkflow(fmt.Sprintf("workflows[%d]", i), wf, h); err != nil {
			return err
		}
	}
	for i, trigger := range doc.Triggers {
		if err := walkHCLTrigger(fmt.Sprintf("triggers[%d]", i), trigger, h); err != nil {
			return err
		}
	}
	for i, result := range doc.Results {
		if result == nil {
			continue
		}
		if err := walkHCLExtensions(fmt.Sprintf("results[%d]", i), result.Extensions, h); err != nil {
			return err
		}
	}
	if doc.Components != nil {
		if err := walkHCLExtensions("components", doc.Components.Extensions, h); err != nil {
			return err
		}
		if err := walkHCLDynamicMap("components.variables", &doc.Components.Variables, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLOperation(path string, op *Operation, h documentHCLWalkHandlers) error {
	if op == nil {
		return nil
	}
	if err := walkHCLExtensions(path, op.Extensions, h); err != nil {
		return err
	}
	walkHCLDescription(&op.Description, h)
	if err := walkHCLDynamicMap(path+".request", &op.Request, h); err != nil {
		return err
	}
	for i, criterion := range op.SuccessCriteria {
		if err := walkHCLCriterion(fmt.Sprintf("%s.successCriteria[%d]", path, i), criterion, h); err != nil {
			return err
		}
	}
	for i, action := range op.OnFailure {
		if err := walkHCLFailureAction(fmt.Sprintf("%s.onFailure[%d]", path, i), action, h); err != nil {
			return err
		}
	}
	for i, action := range op.OnSuccess {
		if err := walkHCLSuccessAction(fmt.Sprintf("%s.onSuccess[%d]", path, i), action, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLWorkflow(path string, wf *Workflow, h documentHCLWalkHandlers) error {
	if wf == nil {
		return nil
	}
	if err := walkHCLExtensions(path, wf.Extensions, h); err != nil {
		return err
	}
	walkHCLDescription(&wf.Description, h)
	if err := walkHCLParamSchema(path+".inputs", wf.Inputs, h); err != nil {
		return err
	}
	if wf.Idempotency != nil {
		if err := walkHCLExtensions(path+".idempotency", wf.Idempotency.Extensions, h); err != nil {
			return err
		}
	}
	if err := walkHCLSteps(path+".steps", wf.Steps, h); err != nil {
		return err
	}
	if err := walkHCLCases(path+".cases", wf.Cases, h); err != nil {
		return err
	}
	return walkHCLSteps(path+".default", wf.Default, h)
}

func walkHCLSteps(path string, steps []*Step, h documentHCLWalkHandlers) error {
	for i, step := range steps {
		if err := walkHCLStep(fmt.Sprintf("%s[%d]", path, i), step, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLStep(path string, step *Step, h documentHCLWalkHandlers) error {
	if step == nil {
		return nil
	}
	if err := walkHCLExtensions(path, step.Extensions, h); err != nil {
		return err
	}
	walkHCLDescription(&step.Description, h)
	if err := walkHCLDynamicMap(path+".body", &step.Body, h); err != nil {
		return err
	}
	if err := walkHCLSteps(path+".steps", step.Steps, h); err != nil {
		return err
	}
	if err := walkHCLCases(path+".cases", step.Cases, h); err != nil {
		return err
	}
	return walkHCLSteps(path+".default", step.Default, h)
}

func walkHCLCases(path string, cases []*Case, h documentHCLWalkHandlers) error {
	for i, c := range cases {
		if err := walkHCLCase(fmt.Sprintf("%s[%d]", path, i), c, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLCase(path string, c *Case, h documentHCLWalkHandlers) error {
	if c == nil {
		return nil
	}
	if err := walkHCLExtensions(path, c.Extensions, h); err != nil {
		return err
	}
	if err := walkHCLDynamicMap(path+".body", &c.Body, h); err != nil {
		return err
	}
	return walkHCLSteps(path+".steps", c.Steps, h)
}

func walkHCLTrigger(path string, trigger *Trigger, h documentHCLWalkHandlers) error {
	if trigger == nil {
		return nil
	}
	if err := walkHCLExtensions(path, trigger.Extensions, h); err != nil {
		return err
	}
	if err := walkHCLDynamicMap(path+".options", &trigger.Options, h); err != nil {
		return err
	}
	for i, route := range trigger.Routes {
		if route == nil {
			continue
		}
		if err := walkHCLExtensions(fmt.Sprintf("%s.routes[%d]", path, i), route.Extensions, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLCriterion(path string, criterion *Criterion, h documentHCLWalkHandlers) error {
	if criterion == nil {
		return nil
	}
	return walkHCLExtensions(path, criterion.Extensions, h)
}

func walkHCLFailureAction(path string, action *FailureAction, h documentHCLWalkHandlers) error {
	if action == nil {
		return nil
	}
	if err := walkHCLExtensions(path, action.Extensions, h); err != nil {
		return err
	}
	for i, criterion := range action.Criteria {
		if err := walkHCLCriterion(fmt.Sprintf("%s.criteria[%d]", path, i), criterion, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLSuccessAction(path string, action *SuccessAction, h documentHCLWalkHandlers) error {
	if action == nil {
		return nil
	}
	if err := walkHCLExtensions(path, action.Extensions, h); err != nil {
		return err
	}
	for i, criterion := range action.Criteria {
		if err := walkHCLCriterion(fmt.Sprintf("%s.criteria[%d]", path, i), criterion, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLParamSchema(path string, schema *ParamSchema, h documentHCLWalkHandlers) error {
	if schema == nil {
		return nil
	}
	if err := walkHCLExtensions(path, schema.Extensions, h); err != nil {
		return err
	}
	if h.paramSchema != nil {
		if err := h.paramSchema(path, schema); err != nil {
			return err
		}
	}
	for name, child := range schema.Properties {
		if err := walkHCLParamSchema(path+".properties."+name, child, h); err != nil {
			return err
		}
	}
	if err := walkHCLParamSchema(path+".items", schema.Items, h); err != nil {
		return err
	}
	for i, child := range schema.AllOf {
		if err := walkHCLParamSchema(fmt.Sprintf("%s.allOf[%d]", path, i), child, h); err != nil {
			return err
		}
	}
	for i, child := range schema.OneOf {
		if err := walkHCLParamSchema(fmt.Sprintf("%s.oneOf[%d]", path, i), child, h); err != nil {
			return err
		}
	}
	for i, child := range schema.AnyOf {
		if err := walkHCLParamSchema(fmt.Sprintf("%s.anyOf[%d]", path, i), child, h); err != nil {
			return err
		}
	}
	return nil
}

func walkHCLDescription(value *string, h documentHCLWalkHandlers) {
	if h.description != nil {
		h.description(value)
	}
}

func walkHCLDynamicMap(path string, value *map[string]any, h documentHCLWalkHandlers) error {
	if h.dynamicMap != nil && *value != nil {
		return h.dynamicMap(path, value)
	}
	return nil
}

func walkHCLExtensions(path string, extensions map[string]any, h documentHCLWalkHandlers) error {
	if h.extensions != nil {
		return h.extensions(path, extensions)
	}
	return nil
}
