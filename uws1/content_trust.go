package uws1

import (
	"fmt"
	"sort"
)

// ContentTrustLevel describes the reviewed integrity provenance of content.
// It is intentionally independent of the value's runtime type or capability.
type ContentTrustLevel string

const (
	ContentTrustUnknown   ContentTrustLevel = "unknown"
	ContentTrustTrusted   ContentTrustLevel = "trusted"
	ContentTrustUntrusted ContentTrustLevel = "untrusted"
)

// ContentTrust contains optional, reviewable provenance declarations for
// source descriptions, operations, triggers, and workflow entry inputs.
type ContentTrust struct {
	SourceDescriptions map[string]ContentTrustLevel      `json:"sourceDescriptions,omitempty" yaml:"sourceDescriptions,omitempty" hcl:"sourceDescriptions,optional"`
	Operations         map[string]*OperationContentTrust `json:"operations,omitempty" yaml:"operations,omitempty" hcl:"operations,optional"`
	Triggers           map[string]ContentTrustLevel      `json:"triggers,omitempty" yaml:"triggers,omitempty" hcl:"triggers,optional"`
	Workflows          map[string]*WorkflowContentTrust  `json:"workflows,omitempty" yaml:"workflows,omitempty" hcl:"workflows,optional"`
	Extensions         map[string]any                    `json:"-" yaml:"-" hcl:"extensions,block"`
}

// OperationContentTrust declares an operation-wide default and optional
// overrides for named operation outputs.
type OperationContentTrust struct {
	Default    ContentTrustLevel            `json:"default,omitempty" yaml:"default,omitempty" hcl:"default,optional"`
	Outputs    map[string]ContentTrustLevel `json:"outputs,omitempty" yaml:"outputs,omitempty" hcl:"outputs,optional"`
	Extensions map[string]any               `json:"-" yaml:"-" hcl:"extensions,block"`
}

// WorkflowContentTrust declares the provenance of externally supplied inputs
// to a workflow. A same-document workflow call retains the provenance of the
// values passed by the caller instead of applying this entry declaration.
type WorkflowContentTrust struct {
	Default    ContentTrustLevel            `json:"default,omitempty" yaml:"default,omitempty" hcl:"default,optional"`
	Inputs     map[string]ContentTrustLevel `json:"inputs,omitempty" yaml:"inputs,omitempty" hcl:"inputs,optional"`
	Extensions map[string]any               `json:"-" yaml:"-" hcl:"extensions,block"`
}

type contentTrustAlias ContentTrust
type operationContentTrustAlias OperationContentTrust
type workflowContentTrustAlias WorkflowContentTrust

var contentTrustKnownFields = []string{"sourceDescriptions", "operations", "triggers", "workflows"}
var operationContentTrustKnownFields = []string{"default", "outputs"}
var workflowContentTrustKnownFields = []string{"default", "inputs"}

func (c *ContentTrust) UnmarshalJSON(data []byte) error {
	var alias contentTrustAlias
	_, extensions, err := unmarshalCoreWithExtensions(data, "contentTrust", contentTrustKnownFields, &alias)
	if err != nil {
		return err
	}
	*c = ContentTrust(alias)
	c.Extensions = extensions
	return nil
}

func (c ContentTrust) MarshalJSON() ([]byte, error) {
	alias := contentTrustAlias(c)
	return marshalWithExtensions(&alias, c.Extensions)
}

func (c *OperationContentTrust) UnmarshalJSON(data []byte) error {
	var alias operationContentTrustAlias
	_, extensions, err := unmarshalCoreWithExtensions(data, "operationContentTrust", operationContentTrustKnownFields, &alias)
	if err != nil {
		return err
	}
	*c = OperationContentTrust(alias)
	c.Extensions = extensions
	return nil
}

func (c OperationContentTrust) MarshalJSON() ([]byte, error) {
	alias := operationContentTrustAlias(c)
	return marshalWithExtensions(&alias, c.Extensions)
}

func (c *WorkflowContentTrust) UnmarshalJSON(data []byte) error {
	var alias workflowContentTrustAlias
	_, extensions, err := unmarshalCoreWithExtensions(data, "workflowContentTrust", workflowContentTrustKnownFields, &alias)
	if err != nil {
		return err
	}
	*c = WorkflowContentTrust(alias)
	c.Extensions = extensions
	return nil
}

func (c WorkflowContentTrust) MarshalJSON() ([]byte, error) {
	alias := workflowContentTrustAlias(c)
	return marshalWithExtensions(&alias, c.Extensions)
}

func isContentTrustLevel(level ContentTrustLevel) bool {
	switch level {
	case ContentTrustUnknown, ContentTrustTrusted, ContentTrustUntrusted:
		return true
	default:
		return false
	}
}

func validateContentTrustLevel(level ContentTrustLevel, path string, result *ValidationResult) {
	if !isContentTrustLevel(level) {
		result.addError(path, fmt.Sprintf("%q is not valid (must be unknown, trusted, or untrusted)", level))
	}
}

func (d *Document) validateContentTrust(idx *documentIndex, result *ValidationResult) {
	trust := d.ContentTrust
	if trust == nil {
		return
	}
	if !supportsUWSVersionAtLeast(d.UWS, 1, 9, 1) {
		result.addError("contentTrust", "requires UWS 1.9.1 or later")
	}
	if len(trust.SourceDescriptions) == 0 && len(trust.Operations) == 0 && len(trust.Triggers) == 0 && len(trust.Workflows) == 0 && len(trust.Extensions) == 0 {
		result.addError("contentTrust", "must contain at least one declaration or extension")
	}
	for _, name := range sortedContentTrustKeys(trust.SourceDescriptions) {
		level := trust.SourceDescriptions[name]
		path := "contentTrust.sourceDescriptions." + name
		if _, ok := idx.sourceDescriptions[name]; !ok {
			result.addError(path, fmt.Sprintf("references unknown sourceDescription %q", name))
		}
		validateContentTrustLevel(level, path, result)
	}
	for _, id := range sortedContentTrustKeys(trust.Operations) {
		declaration := trust.Operations[id]
		path := "contentTrust.operations." + id
		op := idx.operations[id]
		if op == nil {
			result.addError(path, fmt.Sprintf("references unknown operationId %q", id))
		}
		if declaration == nil {
			result.addError(path, "is nil")
			continue
		}
		if declaration.Default == "" && len(declaration.Outputs) == 0 && len(declaration.Extensions) == 0 {
			result.addError(path, "must contain default, outputs, or an extension")
		}
		if declaration.Default != "" {
			validateContentTrustLevel(declaration.Default, path+".default", result)
		}
		for _, output := range sortedContentTrustKeys(declaration.Outputs) {
			level := declaration.Outputs[output]
			outputPath := path + ".outputs." + output
			if op != nil {
				if _, ok := op.Outputs[output]; !ok {
					result.addError(outputPath, fmt.Sprintf("references undeclared operation output %q", output))
				}
			}
			validateContentTrustLevel(level, outputPath, result)
		}
	}
	for _, id := range sortedContentTrustKeys(trust.Triggers) {
		level := trust.Triggers[id]
		path := "contentTrust.triggers." + id
		if !idx.triggers[id] {
			result.addError(path, fmt.Sprintf("references unknown triggerId %q", id))
		}
		validateContentTrustLevel(level, path, result)
	}
	for _, id := range sortedContentTrustKeys(trust.Workflows) {
		declaration := trust.Workflows[id]
		path := "contentTrust.workflows." + id
		workflow := idx.workflows[id]
		if workflow == nil {
			result.addError(path, fmt.Sprintf("references unknown workflowId %q", id))
		}
		if declaration == nil {
			result.addError(path, "is nil")
			continue
		}
		if declaration.Default == "" && len(declaration.Inputs) == 0 && len(declaration.Extensions) == 0 {
			result.addError(path, "must contain default, inputs, or an extension")
		}
		if declaration.Default != "" {
			validateContentTrustLevel(declaration.Default, path+".default", result)
		}
		for _, input := range sortedContentTrustKeys(declaration.Inputs) {
			level := declaration.Inputs[input]
			inputPath := path + ".inputs." + input
			if workflow != nil {
				if workflow.Inputs == nil || workflow.Inputs.Properties == nil {
					result.addError(inputPath, fmt.Sprintf("references undeclared workflow input %q", input))
				} else if _, ok := workflow.Inputs.Properties[input]; !ok {
					result.addError(inputPath, fmt.Sprintf("references undeclared workflow input %q", input))
				}
			}
			validateContentTrustLevel(level, inputPath, result)
		}
	}
}

func sortedContentTrustKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
