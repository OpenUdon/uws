package uws1

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validDocument() *Document {
	return &Document{
		UWS: "1.0.0",
		Info: &Info{
			Title:   "Test",
			Version: "1.0.0",
		},
		SourceDescriptions: []*SourceDescription{
			{Name: "api", URL: "./openapi.yaml", Type: SourceDescriptionTypeOpenAPI},
		},
		Operations: []*Operation{
			{
				OperationID:        "get_data",
				SourceDescription:  "api",
				OpenAPIOperationID: "getData",
			},
		},
	}
}

func TestValidate_Valid(t *testing.T) {
	doc := validDocument()
	assert.NoError(t, doc.Validate())
}

func TestValidate_OpenAPIOperationRefValid(t *testing.T) {
	doc := validDocument()
	doc.Operations[0].OpenAPIOperationID = ""
	doc.Operations[0].OpenAPIOperationRef = "#/paths/~1data/get"
	assert.NoError(t, doc.Validate())
}

func TestValidate_UWS11TimeoutAndIdempotency(t *testing.T) {
	timeout := 30.5
	ttl := 86400.0
	doc := validDocument()
	doc.UWS = "1.1.0"
	doc.Operations[0].Timeout = &timeout
	doc.Workflows = []*Workflow{{
		WorkflowID:  "main",
		Type:        WorkflowTypeSequence,
		Idempotency: &Idempotency{Key: "$variables.requestId", OnConflict: "returnPrevious", TTL: &ttl},
		WorkflowExecutionFields: WorkflowExecutionFields{
			Timeout: &timeout,
		},
		Steps: []*Step{{
			StepID:       "fetch",
			OperationRef: "get_data",
			StepExecutionFields: StepExecutionFields{
				Timeout: &timeout,
			},
		}},
	}}
	assert.NoError(t, doc.Validate())
}

func TestValidate_UWS11FieldsRequireVersionAndPositiveValues(t *testing.T) {
	timeout := 0.0
	ttl := -1.0
	doc := validDocument()
	doc.Operations[0].Timeout = &timeout
	doc.Workflows = []*Workflow{{
		WorkflowID:  "main",
		Type:        WorkflowTypeSequence,
		Idempotency: &Idempotency{Key: "$variables.requestId"},
		Steps: []*Step{{
			StepID:       "fetch",
			OperationRef: "get_data",
			StepExecutionFields: StepExecutionFields{
				Timeout: &timeout,
			},
		}},
	}}
	err := doc.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "operations[0].timeout requires UWS 1.1.0 or later")
	assert.Contains(t, err.Error(), "workflows[0].idempotency requires UWS 1.1.0 or later")
	assert.Contains(t, err.Error(), "workflows[0].steps[0].timeout requires UWS 1.1.0 or later")

	doc.UWS = "1.1.0"
	doc.Workflows[0].Idempotency = &Idempotency{OnConflict: "replace", TTL: &ttl}
	err = doc.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "operations[0].timeout must be greater than 0")
	assert.Contains(t, err.Error(), "workflows[0].idempotency.key is required")
	assert.Contains(t, err.Error(), `workflows[0].idempotency.onConflict "replace" is not valid`)
	assert.Contains(t, err.Error(), "workflows[0].idempotency.ttl must be greater than 0")
	assert.Contains(t, err.Error(), "workflows[0].steps[0].timeout must be greater than 0")
}

func TestOperationBindingHelpers(t *testing.T) {
	var nilOp *Operation
	assert.False(t, nilOp.HasOpenAPIBinding())
	assert.False(t, nilOp.HasCompleteOpenAPIBinding())
	assert.Empty(t, nilOp.ExtensionProfile())
	assert.False(t, nilOp.IsExtensionOwned())

	openAPIBound := &Operation{
		SourceDescription:  "api",
		OpenAPIOperationID: "getData",
	}
	assert.True(t, openAPIBound.HasOpenAPIBinding())
	assert.True(t, openAPIBound.HasCompleteOpenAPIBinding())
	assert.Empty(t, openAPIBound.ExtensionProfile())
	assert.False(t, openAPIBound.IsExtensionOwned())

	partialOpenAPIBinding := &Operation{OpenAPIOperationID: "getData"}
	assert.True(t, partialOpenAPIBinding.HasOpenAPIBinding())
	assert.False(t, partialOpenAPIBinding.HasCompleteOpenAPIBinding())
	assert.False(t, partialOpenAPIBinding.IsExtensionOwned())

	conflictingOpenAPIBinding := &Operation{
		SourceDescription:   "api",
		OpenAPIOperationID:  "getData",
		OpenAPIOperationRef: "#/paths/~1data/get",
	}
	assert.True(t, conflictingOpenAPIBinding.HasOpenAPIBinding())
	assert.False(t, conflictingOpenAPIBinding.HasCompleteOpenAPIBinding())

	extensionOwned := &Operation{
		Extensions: map[string]any{ExtensionOperationProfile: " udon "},
	}
	assert.False(t, extensionOwned.HasOpenAPIBinding())
	assert.False(t, extensionOwned.HasCompleteOpenAPIBinding())
	assert.Equal(t, "udon", extensionOwned.ExtensionProfile())
	assert.True(t, extensionOwned.IsExtensionOwned())

	whitespaceProfile := &Operation{
		Extensions: map[string]any{ExtensionOperationProfile: "   "},
	}
	assert.Empty(t, whitespaceProfile.ExtensionProfile())
	assert.False(t, whitespaceProfile.IsExtensionOwned())

	nonStringProfile := &Operation{
		Extensions: map[string]any{ExtensionOperationProfile: 1},
	}
	assert.Empty(t, nonStringProfile.ExtensionProfile())
	assert.False(t, nonStringProfile.IsExtensionOwned())
}

func TestValidate_MissingRootFields(t *testing.T) {
	doc := validDocument()
	doc.UWS = ""
	doc.Info = nil
	doc.SourceDescriptions = nil
	doc.Operations = nil

	err := doc.Validate()
	assert.ErrorContains(t, err, "uws version is required")
	assert.ErrorContains(t, err, "info is required")
	assert.ErrorContains(t, err, "operations at least one operation is required")
}

func TestValidate_BadVersionPattern(t *testing.T) {
	doc := validDocument()
	doc.UWS = "2.0.0"
	assert.ErrorContains(t, doc.Validate(), "does not match pattern")
}

func TestValidate_InfoRequiredFields(t *testing.T) {
	doc := validDocument()
	doc.Info.Title = ""
	doc.Info.Version = ""

	err := doc.Validate()
	assert.ErrorContains(t, err, "info.title is required")
	assert.ErrorContains(t, err, "info.version is required")
}

func TestValidate_OperationBindingRules(t *testing.T) {
	t.Run("missing operationId", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OperationID = ""
		assert.ErrorContains(t, doc.Validate(), "operationId is required")
	})

	t.Run("missing sourceDescription", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].SourceDescription = ""
		assert.ErrorContains(t, doc.Validate(), "sourceDescription is required for OpenAPI-bound operations")
	})

	t.Run("unknown sourceDescription", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].SourceDescription = "missing"
		assert.ErrorContains(t, doc.Validate(), `references unknown sourceDescription "missing"`)
	})

	t.Run("missing OpenAPI binding", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OpenAPIOperationID = ""
		assert.ErrorContains(t, doc.Validate(), "requires exactly one of openapiOperationId or openapiOperationRef for OpenAPI-bound operations")
	})

	t.Run("conflicting OpenAPI bindings", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OpenAPIOperationRef = "#/paths/~1data/get"
		assert.ErrorContains(t, doc.Validate(), "cannot specify both openapiOperationId and openapiOperationRef")
	})

	t.Run("extension-owned operation", func(t *testing.T) {
		doc := &Document{
			UWS:  "1.0.0",
			Info: &Info{Title: "Extension", Version: "1.0.0"},
			Operations: []*Operation{
				{
					OperationID: "build_email",
					Extensions: map[string]any{
						ExtensionOperationProfile: "udon",
						"x-udon-runtime":          map[string]any{"type": "fnct", "function": "mail_raw"},
					},
				},
			},
		}
		assert.NoError(t, doc.Validate())
	})

	t.Run("missing binding and extension", func(t *testing.T) {
		doc := &Document{
			UWS:        "1.0.0",
			Info:       &Info{Title: "Invalid", Version: "1.0.0"},
			Operations: []*Operation{{OperationID: "op"}},
		}
		assert.ErrorContains(t, doc.Validate(), "requires an OpenAPI binding or x-uws-operation-profile")
	})

	t.Run("extension-owned operation requires profile marker", func(t *testing.T) {
		doc := &Document{
			UWS:  "1.0.0",
			Info: &Info{Title: "Extension", Version: "1.0.0"},
			Operations: []*Operation{
				{
					OperationID: "build_email",
					Extensions: map[string]any{
						"x-udon-runtime": map[string]any{"type": "fnct", "function": "mail_raw"},
					},
				},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "requires an OpenAPI binding or x-uws-operation-profile")
	})

	t.Run("extension-owned operation requires non-whitespace profile marker", func(t *testing.T) {
		doc := &Document{
			UWS:  "1.0.0",
			Info: &Info{Title: "Extension", Version: "1.0.0"},
			Operations: []*Operation{
				{
					OperationID: "build_email",
					Extensions: map[string]any{
						ExtensionOperationProfile: "   ",
						"x-udon-runtime":          map[string]any{"type": "fnct", "function": "mail_raw"},
					},
				},
			},
		}
		assert.ErrorContains(t, doc.Validate(), "requires an OpenAPI binding or x-uws-operation-profile")
	})

	t.Run("non pointer operation ref", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OpenAPIOperationID = ""
		doc.Operations[0].OpenAPIOperationRef = "operation://getData"
		assert.ErrorContains(t, doc.Validate(), "must be a JSON Pointer fragment")
	})

	t.Run("standard request binding keys", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].Request = map[string]any{
			"path":   map[string]any{"id": "123"},
			"query":  map[string]any{"limit": 10},
			"header": map[string]any{"X-Test": "ok"},
			"cookie": map[string]any{"session": "abc"},
			"body":   map[string]any{"name": "widget"},
			"x-test": map[string]any{"ok": true},
		}
		assert.NoError(t, doc.Validate())
	})

	t.Run("unknown request binding key", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].Request = map[string]any{"limit": 10}
		assert.ErrorContains(t, doc.Validate(), "is not a standard request binding key")
	})

	t.Run("request parameter sections must be objects", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].Request = map[string]any{"query": "limit=10"}
		assert.ErrorContains(t, doc.Validate(), "request.query must be an object")
	})
}

func TestValidate_DuplicateIDs(t *testing.T) {
	doc := validDocument()
	doc.SourceDescriptions = append(doc.SourceDescriptions, &SourceDescription{Name: "api", URL: "./other.yaml"})
	doc.Operations = append(doc.Operations, &Operation{
		OperationID:        "get_data",
		SourceDescription:  "api",
		OpenAPIOperationID: "getOtherData",
	})
	doc.Workflows = []*Workflow{
		{WorkflowID: "wf", Type: "parallel"},
		{WorkflowID: "wf", Type: "switch"},
	}
	doc.Triggers = []*Trigger{
		{TriggerID: "t1"},
		{TriggerID: "t1"},
	}

	err := doc.Validate()
	assert.ErrorContains(t, err, "duplicate sourceDescription name")
	assert.ErrorContains(t, err, "duplicate operationId")
	assert.ErrorContains(t, err, "duplicate workflowId")
	assert.ErrorContains(t, err, "duplicate triggerId")
}

func TestValidate_SourceDescriptions(t *testing.T) {
	t.Run("missing name", func(t *testing.T) {
		doc := validDocument()
		doc.SourceDescriptions[0].Name = ""
		assert.ErrorContains(t, doc.Validate(), "sourceDescriptions[0].name is required")
	})

	t.Run("missing url", func(t *testing.T) {
		doc := validDocument()
		doc.SourceDescriptions[0].URL = ""
		assert.ErrorContains(t, doc.Validate(), "sourceDescriptions[0].url is required")
	})

	t.Run("invalid name", func(t *testing.T) {
		doc := validDocument()
		doc.SourceDescriptions[0].Name = "bad name"
		assert.ErrorContains(t, doc.Validate(), "must match pattern")
	})

	t.Run("invalid type", func(t *testing.T) {
		doc := validDocument()
		doc.SourceDescriptions[0].Type = "arazzo"
		assert.ErrorContains(t, doc.Validate(), "must be openapi")
	})

	t.Run("omitted type", func(t *testing.T) {
		doc := validDocument()
		doc.SourceDescriptions[0].Type = ""
		assert.NoError(t, doc.Validate())
	})
}

func TestValidate_SourceDescriptionsRequiredWhenBound(t *testing.T) {
	t.Run("missing top-level array with bound op", func(t *testing.T) {
		doc := validDocument()
		doc.SourceDescriptions = nil
		err := doc.Validate()
		require.Error(t, err)
		assert.ErrorContains(t, err, "sourceDescriptions is required when any operation declares sourceDescription")
	})

	t.Run("empty top-level array with bound op", func(t *testing.T) {
		doc := validDocument()
		doc.SourceDescriptions = []*SourceDescription{}
		err := doc.Validate()
		require.Error(t, err)
		assert.ErrorContains(t, err, "sourceDescriptions is required when any operation declares sourceDescription")
	})

	t.Run("extension-owned op needs no sourceDescriptions", func(t *testing.T) {
		doc := &Document{
			UWS:  "1.0.0",
			Info: &Info{Title: "Ext", Version: "1.0.0"},
			Operations: []*Operation{
				{
					OperationID: "do_thing",
					Extensions:  map[string]any{ExtensionOperationProfile: "udon"},
				},
			},
		}
		assert.NoError(t, doc.Validate())
	})
}


func TestCriterionUnmarshalRejectsExplicitEmptyType(t *testing.T) {
	var criterion Criterion
	require.ErrorContains(t, json.Unmarshal([]byte(`{"condition":"true","type":""}`), &criterion), "criterion.type must be omitted")

	require.NoError(t, json.Unmarshal([]byte(`{"condition":"true"}`), &criterion))
	assert.Empty(t, criterion.Type)
}



func TestValidate_OpenAPIOperationRefRequiresPointer(t *testing.T) {
	doc := validDocument()
	doc.Operations[0].OpenAPIOperationID = ""
	doc.Operations[0].OpenAPIOperationRef = "paths/~1data/get"
	assert.ErrorContains(t, doc.Validate(), "must be a JSON Pointer fragment beginning with #/")
}

func TestValidate_ComponentsVariables(t *testing.T) {
	doc := validDocument()
	doc.Components = &Components{
		Variables: map[string]any{
			"ok_name":  true,
			"bad name": false,
		},
	}

	assert.ErrorContains(t, doc.Validate(), "component name")
}

func TestValidationResult_ErrorFormat(t *testing.T) {
	var nilResult *ValidationResult
	assert.True(t, nilResult.Valid())
	assert.Empty(t, nilResult.Error())

	empty := &ValidationResult{}
	assert.True(t, empty.Valid())
	assert.Empty(t, empty.Error())

	one := &ValidationResult{
		Errors: []ValidationError{
			{Path: "operations[0]", Message: "is invalid"},
		},
	}
	assert.False(t, one.Valid())
	assert.Equal(t, "operations[0] is invalid", one.Error())

	many := &ValidationResult{
		Errors: []ValidationError{
			{Path: "info.title", Message: "is required"},
			{Path: "uws", Message: "must match pattern"},
		},
	}
	assert.Equal(t, "info.title is required; uws must match pattern", many.Error())
}

func TestValidateResult_AccumulatesErrors(t *testing.T) {
	doc := &Document{
		UWS:  "2.0.0",
		Info: &Info{},
		SourceDescriptions: []*SourceDescription{
			{Name: "api", URL: "./openapi.yaml", Type: SourceDescriptionTypeOpenAPI},
		},
		Operations: []*Operation{
			{OperationID: "op", SourceDescription: "api"},
			{OperationID: "op", SourceDescription: "missing", OpenAPIOperationRef: "#/paths/~1x/get"},
		},
	}

	result := doc.ValidateResult()
	assert.False(t, result.Valid())
	want := []ValidationError{
		{Path: "uws", Message: `version "2.0.0" does not match pattern 1.x.x`},
		{Path: "info.title", Message: "is required"},
		{Path: "info.version", Message: "is required"},
		{Path: "operations[1].operationId", Message: `duplicate operationId "op"`},
		{Path: "operations[0]", Message: "requires exactly one of openapiOperationId or openapiOperationRef for OpenAPI-bound operations"},
		{Path: "operations[1].sourceDescription", Message: `references unknown sourceDescription "missing"`},
	}
	assert.Equal(t, want, result.Errors)
}

func TestValidateResult_StructuredErrorShape(t *testing.T) {
	// Each case exercises one distinct error path and asserts the exact
	// ValidationError tuple rather than substring-matching the flattened
	// string.
	t.Run("duplicate workflowId", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "w", Type: "sequence"},
			{WorkflowID: "w", Type: "sequence"},
		}
		result := doc.ValidateResult()
		assert.Contains(t, result.Errors, ValidationError{
			Path:    "workflows[1].workflowId",
			Message: `duplicate workflowId "w"`,
		})
	})

	t.Run("unknown trigger route output", func(t *testing.T) {
		doc := validDocument()
		doc.Triggers = []*Trigger{
			{
				TriggerID: "hook",
				Outputs:   []string{"ok"},
				Routes: []*TriggerRoute{
					{TriggerRouteFields: TriggerRouteFields{Output: "missing", To: []string{"get_data"}}},
				},
			},
		}
		result := doc.ValidateResult()
		assert.Contains(t, result.Errors, ValidationError{
			Path:    "triggers[0].routes[0].output",
			Message: `"missing" is not a declared trigger output`,
		})
	})

	t.Run("merge workflow missing dependsOn", func(t *testing.T) {
		doc := validDocument()
		doc.Workflows = []*Workflow{
			{WorkflowID: "m", Type: "merge"},
		}
		result := doc.ValidateResult()
		assert.Contains(t, result.Errors, ValidationError{
			Path:    "workflows[0].dependsOn",
			Message: "is required and must name at least one upstream construct for merge",
		})
	})

	t.Run("criterion regex without context", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].SuccessCriteria = []*Criterion{
			{Condition: "foo", Type: CriterionRegex},
		}
		result := doc.ValidateResult()
		assert.Contains(t, result.Errors, ValidationError{
			Path:    "operations[0].successCriteria[0].context",
			Message: "is required when type is regex, jsonpath, or xpath",
		})
	})

	t.Run("failure action retry without limit", func(t *testing.T) {
		doc := validDocument()
		doc.Operations[0].OnFailure = []*FailureAction{
			{Name: "r", Type: "retry"},
		}
		result := doc.ValidateResult()
		assert.Contains(t, result.Errors, ValidationError{
			Path:    "operations[0].onFailure[0]",
			Message: "retry requires retryLimit > 0",
		})
	})
}

func TestValidate_DependencyCycle_TwoNodes(t *testing.T) {
	doc := validDocument()
	doc.Operations = []*Operation{
		{
			OperationID:        "a",
			SourceDescription:  "api",
			OpenAPIOperationID: "getA",
			OperationExecutionFields: OperationExecutionFields{
				DependsOn: []string{"b"},
			},
		},
		{
			OperationID:        "b",
			SourceDescription:  "api",
			OpenAPIOperationID: "getB",
			OperationExecutionFields: OperationExecutionFields{
				DependsOn: []string{"a"},
			},
		},
	}

	err := doc.Validate()
	assert.ErrorContains(t, err, "cycle detected")
	assert.ErrorContains(t, err, "a -> b -> a")
}

func TestValidate_DependencyCycle_SelfLoop(t *testing.T) {
	doc := validDocument()
	doc.Operations[0].DependsOn = []string{"get_data"}

	err := doc.Validate()
	assert.ErrorContains(t, err, "cycle detected: get_data -> get_data")
}

func TestValidate_DependencyCycle_ThreeNodes(t *testing.T) {
	doc := validDocument()
	doc.Operations = []*Operation{
		{OperationID: "a", SourceDescription: "api", OpenAPIOperationID: "getA", OperationExecutionFields: OperationExecutionFields{DependsOn: []string{"b"}}},
		{OperationID: "b", SourceDescription: "api", OpenAPIOperationID: "getB", OperationExecutionFields: OperationExecutionFields{DependsOn: []string{"c"}}},
		{OperationID: "c", SourceDescription: "api", OpenAPIOperationID: "getC", OperationExecutionFields: OperationExecutionFields{DependsOn: []string{"a"}}},
	}

	err := doc.Validate()
	assert.ErrorContains(t, err, "cycle detected")
	assert.ErrorContains(t, err, "a -> b -> c -> a")
}

func TestValidate_DependencyCycle_ThroughParallelGroup(t *testing.T) {
	doc := validDocument()
	doc.Operations = []*Operation{
		{
			OperationID:        "a",
			SourceDescription:  "api",
			OpenAPIOperationID: "getA",
			OperationExecutionFields: OperationExecutionFields{
				ParallelGroup: "fanout",
				DependsOn:     []string{"b"},
			},
		},
		{
			OperationID:        "b",
			SourceDescription:  "api",
			OpenAPIOperationID: "getB",
			OperationExecutionFields: OperationExecutionFields{
				DependsOn: []string{"fanout"},
			},
		},
	}

	err := doc.Validate()
	assert.ErrorContains(t, err, "cycle detected")
}

func TestValidate_DependencyCycle_AcyclicIsFine(t *testing.T) {
	doc := validDocument()
	doc.Operations = []*Operation{
		{OperationID: "a", SourceDescription: "api", OpenAPIOperationID: "getA"},
		{OperationID: "b", SourceDescription: "api", OpenAPIOperationID: "getB", OperationExecutionFields: OperationExecutionFields{DependsOn: []string{"a"}}},
		{OperationID: "c", SourceDescription: "api", OpenAPIOperationID: "getC", OperationExecutionFields: OperationExecutionFields{DependsOn: []string{"a", "b"}}},
	}

	assert.NoError(t, doc.Validate())
}

func TestValidate_DependencyCycle_ReportedOnce(t *testing.T) {
	doc := validDocument()
	doc.Operations = []*Operation{
		{OperationID: "a", SourceDescription: "api", OpenAPIOperationID: "getA", OperationExecutionFields: OperationExecutionFields{DependsOn: []string{"b"}}},
		{OperationID: "b", SourceDescription: "api", OpenAPIOperationID: "getB", OperationExecutionFields: OperationExecutionFields{DependsOn: []string{"a"}}},
	}

	result := doc.ValidateResult()
	cycleErrors := 0
	for _, e := range result.Errors {
		if e.Path == "dependsOn" {
			cycleErrors++
		}
	}
	assert.Equal(t, 1, cycleErrors)
}

func TestValidate_ParamSchema_RequiredMustExistInProperties(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "wf",
			Type:       "sequence",
			Inputs: &ParamSchema{
				Type: "object",
				Properties: map[string]*ParamSchema{
					"limit": {Type: "integer"},
				},
				Required: []string{"missing"},
			},
		},
	}
	err := doc.Validate()
	assert.ErrorContains(t, err, `workflows[0].inputs.required[0]`)
	assert.ErrorContains(t, err, `references unknown property "missing"`)
}

func TestOperationRejectsWorkflowField(t *testing.T) {
	var op Operation
	err := json.Unmarshal([]byte(`{
		"operationId":"op",
		"sourceDescription":"api",
		"openapiOperationId":"getOp",
		"workflow":"child"
	}`), &op)
	assert.ErrorContains(t, err, "not defined by UWS core")
}

func TestValidate_ParamSchema_DuplicateRequired(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "wf",
			Type:       "sequence",
			Inputs: &ParamSchema{
				Properties: map[string]*ParamSchema{"a": {Type: "string"}},
				Required:   []string{"a", "a"},
			},
		},
	}
	err := doc.Validate()
	assert.ErrorContains(t, err, "duplicate required entry")
}

func TestValidate_ParamSchema_NilNestedSchema(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "wf",
			Type:       "sequence",
			Inputs: &ParamSchema{
				Type:  "array",
				AllOf: []*ParamSchema{nil},
			},
		},
	}
	err := doc.Validate()
	assert.ErrorContains(t, err, `workflows[0].inputs.allOf[0]`)
	assert.ErrorContains(t, err, "is nil")
}

func TestValidate_ParamSchema_RecursesIntoItems(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "wf",
			Type:       "sequence",
			Inputs: &ParamSchema{
				Type: "array",
				Items: &ParamSchema{
					Properties: map[string]*ParamSchema{"a": {Type: "string"}},
					Required:   []string{"b"},
				},
			},
		},
	}
	err := doc.Validate()
	assert.ErrorContains(t, err, `workflows[0].inputs.items.required[0]`)
	assert.ErrorContains(t, err, `references unknown property "b"`)
}

func TestValidate_Variables_AcceptsOpenShape(t *testing.T) {
	doc := validDocument()
	doc.Variables = map[string]any{
		"a-string": "hello",
		"a-number": 42,
		"a-bool":   true,
		"a-null":   nil,
		"a-list":   []any{1, "two", map[string]any{"nested": true}},
		"a-obj":    map[string]any{"deep": map[string]any{"deeper": []any{1, 2}}},
	}
	assert.NoError(t, doc.Validate())
}

func TestValidate_Components_Variables_KeyPatternEnforced(t *testing.T) {
	doc := validDocument()
	doc.Components = &Components{
		Variables: map[string]any{
			"valid.key": "ok",
			"bad key":   "nope",
		},
	}
	err := doc.Validate()
	assert.ErrorContains(t, err, `component name "bad key"`)
}

func TestValidate_TriggerOptions_AcceptsOpenShape(t *testing.T) {
	doc := validDocument()
	doc.Triggers = []*Trigger{
		{
			TriggerID: "webhook",
			Options: map[string]any{
				"string": "a",
				"int":    1,
				"nested": map[string]any{"list": []any{"x"}},
			},
		},
	}
	assert.NoError(t, doc.Validate())
}

func TestValidate_ParamSchema_ValidSchemaPasses(t *testing.T) {
	doc := validDocument()
	doc.Workflows = []*Workflow{
		{
			WorkflowID: "wf",
			Type:       "sequence",
			Inputs: &ParamSchema{
				Type: "object",
				Properties: map[string]*ParamSchema{
					"limit": {Type: "integer"},
					"name":  {Type: "string"},
				},
				Required: []string{"limit"},
			},
		},
	}
	assert.NoError(t, doc.Validate())
}

func TestValidate_InvalidFixtures(t *testing.T) {
	cases := []struct {
		name    string
		file    string
		wantErr ValidationError
	}{
		{
			name:    "bad uws version",
			file:    "bad_uws_version.json",
			wantErr: ValidationError{Path: "uws", Message: `version "2.0.0" does not match pattern 1.x.x`},
		},
		{
			name:    "dependency cycle",
			file:    "dependency_cycle.json",
			wantErr: ValidationError{Path: "dependsOn", Message: "cycle detected: a -> b -> a"},
		},
		{
			name:    "duplicate operation id",
			file:    "duplicate_operation_id.json",
			wantErr: ValidationError{Path: "operations[1].operationId", Message: `duplicate operationId "op"`},
		},
		{
			name:    "merge without dependsOn",
			file:    "merge_without_dependson.json",
			wantErr: ValidationError{Path: "workflows[0].dependsOn", Message: "is required and must name at least one upstream construct for merge"},
		},
		{
			name:    "missing openapi binding",
			file:    "missing_openapi_binding.json",
			wantErr: ValidationError{Path: "operations[0]", Message: "requires exactly one of openapiOperationId or openapiOperationRef for OpenAPI-bound operations"},
		},
		{
			name:    "regex criterion without context",
			file:    "regex_criterion_no_context.json",
			wantErr: ValidationError{Path: "operations[0].successCriteria[0].context", Message: "is required when type is regex, jsonpath, or xpath"},
		},
		{
			name:    "retry action without retryLimit",
			file:    "retry_without_limit.json",
			wantErr: ValidationError{Path: "operations[0].onFailure[0]", Message: "retry requires retryLimit > 0"},
		},
		{
			name:    "unknown sourceDescription",
			file:    "unknown_source_description.json",
			wantErr: ValidationError{Path: "operations[0].sourceDescription", Message: `references unknown sourceDescription "missing"`},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "testdata", "invalid", tc.file))
			require.NoError(t, err)
			var doc Document
			require.NoError(t, json.Unmarshal(data, &doc))

			result := doc.ValidateResult()
			assert.Contains(t, result.Errors, tc.wantErr,
				"expected error %+v in results %+v", tc.wantErr, result.Errors)
		})
	}
}
