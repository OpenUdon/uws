package uws1

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// schemaParityEntry binds a Go type that carries x-* extensions to the $def
// that describes it in the latest versions/1.x JSON Schema, and to the package-level knownFields list its
// UnmarshalJSON uses to reject unknown properties. Adding a new such type is
// the one place that must be kept in sync by hand; the test below catches
// every other drift direction automatically.
//
// defName == "" means the schema root (the Document struct).
type schemaParityEntry struct {
	label       string
	defName     string
	goType      reflect.Type
	knownFields []string
}

func schemaParityEntries() []schemaParityEntry {
	return []schemaParityEntry{
		{label: "Document", defName: "", goType: reflect.TypeOf(Document{}), knownFields: documentKnownFields},
		{label: "Info", defName: "info", goType: reflect.TypeOf(Info{}), knownFields: infoKnownFields},
		{label: "SourceDescription", defName: "source-description-object", goType: reflect.TypeOf(SourceDescription{}), knownFields: sourceDescriptionKnownFields},
		{label: "Operation", defName: "operation-object", goType: reflect.TypeOf(Operation{}), knownFields: operationKnownFields},
		{label: "Workflow", defName: "workflow-object", goType: reflect.TypeOf(Workflow{}), knownFields: workflowKnownFields},
		{label: "Step", defName: "step-object", goType: reflect.TypeOf(Step{}), knownFields: stepKnownFields},
		{label: "Case", defName: "case-object", goType: reflect.TypeOf(Case{}), knownFields: caseKnownFields},
		{label: "Idempotency", defName: "idempotency-object", goType: reflect.TypeOf(Idempotency{}), knownFields: idempotencyKnownFields},
		{label: "Trigger", defName: "trigger-object", goType: reflect.TypeOf(Trigger{}), knownFields: triggerKnownFields},
		{label: "TriggerRoute", defName: "trigger-route-object", goType: reflect.TypeOf(TriggerRoute{}), knownFields: triggerRouteKnownFields},
		{label: "ParamSchema", defName: "param-schema-object", goType: reflect.TypeOf(ParamSchema{}), knownFields: paramSchemaKnownFields},
		{label: "StructuralResult", defName: "structural-result-object", goType: reflect.TypeOf(StructuralResult{}), knownFields: structuralResultKnownFields},
		{label: "Components", defName: "components-object", goType: reflect.TypeOf(Components{}), knownFields: componentsKnownFields},
		{label: "Criterion", defName: "criterion-object", goType: reflect.TypeOf(Criterion{}), knownFields: criterionKnownFields},
		{label: "FailureAction", defName: "failure-action-object", goType: reflect.TypeOf(FailureAction{}), knownFields: failureActionKnownFields},
		{label: "SuccessAction", defName: "success-action-object", goType: reflect.TypeOf(SuccessAction{}), knownFields: successActionKnownFields},
	}
}

// TestSchemaParity_StructTagsMatchKnownFields ensures every type with an
// Extensions map keeps its struct JSON tags and its knownFields list in sync.
// A mismatch means rejectUnknownFields would either reject valid documents or
// silently accept invalid ones.
func TestSchemaParity_StructTagsMatchKnownFields(t *testing.T) {
	for _, entry := range schemaParityEntries() {
		t.Run(entry.label, func(t *testing.T) {
			gotTags := jsonFieldTags(t, entry.goType)
			assert.ElementsMatch(t, entry.knownFields, gotTags,
				"struct %s JSON tags do not match its knownFields list", entry.label)

			require.Contains(t, namedFields(entry.goType), "Extensions",
				"type %s is in parity list but has no Extensions field", entry.label)
		})
	}
}

// TestSchemaParity_KnownFieldsMatchSchema ensures every parity-tracked type's
// knownFields list is exactly the set of non-extension properties declared by
// its $def in the latest schema. A drift in either direction fails this test: a new
// schema property without a Go equivalent, or a Go known field without a
// schema property.
func TestSchemaParity_KnownFieldsMatchSchema(t *testing.T) {
	schema := loadSchemaDoc(t)

	for _, entry := range schemaParityEntries() {
		t.Run(entry.label, func(t *testing.T) {
			schemaProps := schemaPropertyNames(t, schema, entry.defName)
			nonExtensionSchemaProps := dropExtensionKeys(schemaProps)

			assert.ElementsMatch(t, entry.knownFields, nonExtensionSchemaProps,
				"%s knownFields diverge from schema %q properties", entry.label, entry.defName)
		})
	}
}

// TestSchemaParity_DefCoverageIsExhaustive fails when the latest schema grows a $def
// that no parity entry tracks. This is the tripwire for adding a new type
// without wiring it through the extension machinery.
func TestSchemaParity_DefCoverageIsExhaustive(t *testing.T) {
	schema := loadSchemaDoc(t)
	defs, ok := schema["$defs"].(map[string]any)
	require.True(t, ok, "schema $defs is not an object")

	tracked := map[string]bool{}
	for _, entry := range schemaParityEntries() {
		if entry.defName != "" {
			tracked[entry.defName] = true
		}
	}
	// Meta defs that describe JSON Schema plumbing, not UWS object shapes.
	tracked["specification-extensions"] = true
	tracked["structural-type-constraints"] = true
	// request-binding-object is a bag of free-form locations backed by
	// map[string]any on Operation, not a dedicated Go type.
	tracked["request-binding-object"] = true

	var untracked []string
	for name := range defs {
		if !tracked[name] {
			untracked = append(untracked, name)
		}
	}
	sort.Strings(untracked)
	assert.Empty(t, untracked,
		"latest schema declares $defs that no schemaParityEntries covers: %v", untracked)
}

func TestSchemaParity_SpecFixedFieldsMatchSchema(t *testing.T) {
	schema := loadSchemaDoc(t)
	spec := loadSpecMarkdown(t)
	specTables := map[string]string{
		"":                          "Document Object",
		"info":                      "Info Object",
		"source-description-object": "Source Description Object",
		"operation-object":          "Operation Object",
		"request-binding-object":    "Request Binding Object",
		"workflow-object":           "Workflow Object",
		"step-object":               "Step Object",
		"case-object":               "Case Object",
		"trigger-object":            "Trigger Object",
		"trigger-route-object":      "Trigger Route Object",
		"criterion-object":          "Criterion Object",
		"failure-action-object":     "Failure Action Object",
		"success-action-object":     "Success Action Object",
		"structural-result-object":  "Structural Result Object",
		"components-object":         "Components Object",
		"param-schema-object":       "ParamSchema Object",
		"idempotency-object":        "Idempotency Object",
	}

	for defName, heading := range specTables {
		t.Run(defLabel(defName), func(t *testing.T) {
			schemaProps := dropExtensionKeys(schemaPropertyNames(t, schema, defName))
			specFields := dropExtensionKeys(specFixedFieldNames(t, spec, heading))
			assert.ElementsMatch(t, schemaProps, specFields,
				"spec fixed fields for %q diverge from schema properties", heading)
		})
	}
}

func jsonFieldTags(t *testing.T, typ reflect.Type) []string {
	t.Helper()
	if typ.Kind() != reflect.Struct {
		t.Fatalf("type %s is not a struct", typ)
	}
	var tags []string
	var collect func(reflect.Type)
	collect = func(current reflect.Type) {
		for i := 0; i < current.NumField(); i++ {
			field := current.Field(i)
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				collect(field.Type)
				continue
			}
			tag, ok := field.Tag.Lookup("json")
			if !ok {
				continue
			}
			name, _, _ := strings.Cut(tag, ",")
			if name == "" || name == "-" {
				continue
			}
			tags = append(tags, name)
		}
	}
	collect(typ)
	return tags
}

func namedFields(typ reflect.Type) []string {
	names := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		names = append(names, typ.Field(i).Name)
	}
	return names
}

func specFixedFieldNames(t *testing.T, spec, heading string) []string {
	t.Helper()
	lines := strings.Split(spec, "\n")
	headingLine := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "#") && strings.Contains(line, heading) {
			headingLine = i
			break
		}
	}
	require.NotEqual(t, -1, headingLine, "spec heading %q not found", heading)

	tableStart := -1
	for i := headingLine + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if i > headingLine+1 && strings.HasPrefix(line, "#### ") && strings.Contains(line, " Object") {
			break
		}
		if line == "| Field Name | Type | Description |" {
			tableStart = i
			break
		}
	}
	require.NotEqual(t, -1, tableStart, "fixed fields table for %q not found", heading)

	var fields []string
	for i := tableStart + 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "|") {
			break
		}
		cells := strings.Split(line, "|")
		if len(cells) < 4 {
			continue
		}
		field := strings.TrimSpace(cells[1])
		field = strings.Trim(field, "`")
		if field != "" {
			fields = append(fields, field)
		}
	}
	sort.Strings(fields)
	return fields
}
