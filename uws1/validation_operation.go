package uws1

import (
	"fmt"
	"reflect"
	"strings"
)

var standardRequestKeys = map[string]bool{
	"path":   true,
	"query":  true,
	"header": true,
	"cookie": true,
	"body":   true,
}

func (op *Operation) validate(path string, idx *documentIndex, result *ValidationResult) {
	if op.OperationID == "" {
		result.addError(path+".operationId", "is required")
	}

	hasSource := op.SourceDescription != ""
	hasSourceOperationID := op.SourceOperationID != ""
	hasSourceOperationRef := op.SourceOperationRef != ""
	hasOpenAPIOperationID := op.OpenAPIOperationID != ""
	hasOpenAPIOperationRef := op.OpenAPIOperationRef != ""
	hasGenericSelector := hasSourceOperationID || hasSourceOperationRef
	hasLegacySelector := hasOpenAPIOperationID || hasOpenAPIOperationRef
	switch {
	case hasSourceOperationID && hasSourceOperationRef:
		result.addError(path, "cannot specify both sourceOperationId and sourceOperationRef")
	case hasOpenAPIOperationID && hasOpenAPIOperationRef:
		result.addError(path, "cannot specify both openapiOperationId and openapiOperationRef")
	case hasGenericSelector && hasLegacySelector:
		result.addError(path, "cannot mix sourceOperationId/sourceOperationRef with openapiOperationId/openapiOperationRef")
	case op.HasSourceBinding():
		if !hasSource {
			if hasLegacySelector && !hasGenericSelector {
				result.addError(path+".sourceDescription", "is required for OpenAPI-bound operations")
			} else {
				result.addError(path+".sourceDescription", "is required for source-bound operations")
			}
		} else if _, ok := idx.sourceDescriptions[op.SourceDescription]; !ok {
			result.addError(path+".sourceDescription", fmt.Sprintf("references unknown sourceDescription %q", op.SourceDescription))
		} else {
			sourceType := idx.sourceDescriptions[op.SourceDescription]
			switch sourceType {
			case "", SourceDescriptionTypeOpenAPI:
				if !hasGenericSelector && !hasLegacySelector {
					result.addError(path, "requires exactly one of sourceOperationId, sourceOperationRef, openapiOperationId, or openapiOperationRef for OpenAPI-bound operations")
				}
			case SourceDescriptionTypeGoogleDiscovery, SourceDescriptionTypeAWSSmithy, SourceDescriptionTypeAsyncAPI, SourceDescriptionTypeGraphQL, SourceDescriptionTypeOpenRPC, SourceDescriptionTypeGRPCProtobuf, SourceDescriptionTypeOData, SourceDescriptionTypeBrowserProfile:
				if hasLegacySelector {
					result.addError(path, fmt.Sprintf("%s sourceDescriptions require sourceOperationId or sourceOperationRef, not openapiOperationId or openapiOperationRef", sourceType))
				}
				if !hasGenericSelector {
					result.addError(path, fmt.Sprintf("%s sourceDescriptions require exactly one of sourceOperationId or sourceOperationRef", sourceType))
				}
			}
		}
		if hasSourceOperationRef && !strings.HasPrefix(op.SourceOperationRef, "#/") {
			result.addError(path+".sourceOperationRef", "must be a JSON Pointer fragment beginning with #/")
		}
		if hasOpenAPIOperationRef && !strings.HasPrefix(op.OpenAPIOperationRef, "#/") {
			result.addError(path+".openapiOperationRef", "must be a JSON Pointer fragment beginning with #/")
		}
	case !op.IsExtensionOwned():
		result.addError(path, "requires a source binding or x-uws-operation-profile for extension-owned operations")
	}
	validateRequest(op.Request, path+".request", result)
	validateDependencyList(op.DependsOn, path+".dependsOn", idx, result)
	validateCriteria(op.SuccessCriteria, path+".successCriteria", result)
	validateFailureActions(op.OnFailure, path+".onFailure", idx, result)
	validateSuccessActions(op.OnSuccess, path+".onSuccess", idx, result)
	validateOutputs(op.Outputs, path+".outputs", result)
}

// validateRequest enforces request-binding shape rules. Body is intentionally
// unconstrained because payload shape is operation-specific (matches the
// schema's request-binding-object, which leaves body open).
func validateRequest(request map[string]any, path string, result *ValidationResult) {
	for key, value := range request {
		if strings.HasPrefix(key, "x-") {
			continue
		}
		if !standardRequestKeys[key] {
			result.addError(path+"."+key, "is not a standard request binding key; use path, query, header, cookie, body, or x-*")
			continue
		}
		switch key {
		case "path", "query", "header", "cookie":
			if !isObjectValue(value) {
				result.addError(path+"."+key, "must be an object")
			}
		}
	}
}

func isObjectValue(value any) bool {
	if value == nil {
		return false
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return false
	}
	if rv.Type().Key().Kind() != reflect.String {
		return false
	}
	return true
}

func validateDependencyList(deps []string, path string, idx *documentIndex, result *ValidationResult) {
	for i, dep := range deps {
		if dep == "" {
			result.addError(fmt.Sprintf("%s[%d]", path, i), "is required")
			continue
		}
		if idx.operations[dep] == nil && idx.workflows[dep] == nil && idx.steps[dep] == nil && !idx.parallelGroups[dep] {
			result.addError(fmt.Sprintf("%s[%d]", path, i), fmt.Sprintf("references unknown dependency %q", dep))
		}
	}
}

func validateOutputs(outputs map[string]string, path string, result *ValidationResult) {
	for key := range outputs {
		if !outputNamePattern.MatchString(key) {
			result.addError(path+"."+key, fmt.Sprintf("output name %q is not valid", key))
		}
	}
}

var validCriterionTypes = map[CriterionExpressionType]bool{
	CriterionSimple:   true,
	CriterionRegex:    true,
	CriterionJSONPath: true,
	CriterionXPath:    true,
}

func validateCriteria(criteria []*Criterion, path string, result *ValidationResult) {
	for i, c := range criteria {
		criterionPath := fmt.Sprintf("%s[%d]", path, i)
		if c == nil {
			result.addError(criterionPath, "is nil")
			continue
		}
		if c.Condition == "" {
			result.addError(criterionPath+".condition", "is required")
		}
		if c.Type != "" && !validCriterionTypes[c.Type] {
			result.addError(criterionPath+".type", fmt.Sprintf("%q is not valid (must be simple, regex, jsonpath, or xpath)", c.Type))
		}
		if c.Type != "" && c.Type != CriterionSimple && c.Context == "" {
			result.addError(criterionPath+".context", "is required when type is regex, jsonpath, or xpath")
		}
	}
}
