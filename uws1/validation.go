package uws1

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents one UWS validation error.
type ValidationError struct {
	Path    string
	Message string
}

// ValidationResult accumulates all validation errors found in a document.
type ValidationResult struct {
	Errors []ValidationError
}

// Valid reports whether validation found no errors.
func (r *ValidationResult) Valid() bool {
	return r == nil || len(r.Errors) == 0
}

// Error returns a compact, path-tagged summary of all validation errors.
func (r *ValidationResult) Error() string {
	if r.Valid() {
		return ""
	}
	msgs := make([]string, 0, len(r.Errors))
	for _, err := range r.Errors {
		msgs = append(msgs, fmt.Sprintf("%s %s", err.Path, err.Message))
	}
	return strings.Join(msgs, "; ")
}

func (r *ValidationResult) addError(path, message string) {
	r.Errors = append(r.Errors, ValidationError{Path: path, Message: message})
}

var (
	// uws1VersionPattern matches valid UWS 1.x version strings (1.MAJOR.MINOR
	// with an optional pre-release suffix). Naming is anchored to the major
	// version intentionally; future major versions need their own pattern.
	uws1VersionPattern = regexp.MustCompile(`^1\.\d+\.\d+(-.+)?$`)
	constructIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	// dottedNamePattern is the shared pattern for component names, output
	// names, and trigger output names. Allowing dots distinguishes these from
	// constructIDPattern, which forbids them.
	dottedNamePattern    = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	componentNamePattern = dottedNamePattern
	outputNamePattern    = dottedNamePattern
)

// Validate runs the semantic validation layer and returns the first error as a
// single error, or nil if the document passes.
//
// Validate assumes the document has already been checked against the matching versions/1.x JSON Schema via
// a JSON Schema validator. The schema layer enforces structural shape (required
// fields, enum values, patterns, uniqueness of array items). Validate layers
// semantic rules on top: duplicate identifiers, reference integrity, binding
// mutual exclusivity, structural-type field constraints, and dependsOn cycles.
// Callers that invoke Validate without the schema pre-pass bypass the
// structural checks.
//
// Use ValidateResult when callers need path-tagged errors instead of a single
// collapsed error.
func (d *Document) Validate() error {
	result := d.ValidateResult()
	if result.Valid() {
		return nil
	}
	return result
}

// ValidateResult runs the semantic validation layer and returns every error it
// finds, each tagged with a structured Path. See Validate for the layering
// contract between this method and the versions/1.x JSON Schema pre-pass.
func (d *Document) ValidateResult() *ValidationResult {
	result := &ValidationResult{}
	if d == nil {
		result.addError("document", "is required")
		return result
	}

	if d.UWS == "" {
		result.addError("uws", "version is required")
	} else if !uws1VersionPattern.MatchString(d.UWS) {
		result.addError("uws", fmt.Sprintf("version %q does not match pattern 1.x.x", d.UWS))
	}
	if d.Info == nil {
		result.addError("info", "is required")
	} else {
		d.Info.validate("info", result)
	}
	if len(d.Operations) == 0 {
		result.addError("operations", "at least one operation is required")
	}
	d.validateTopLevelSourceDescriptions(result)
	d.validateVersionedFields(result)

	idx := buildDocumentIndex(d, result)
	d.validateDocumentReferences(idx, result)
	detectDependencyCycles(idx, result)
	return result
}

// validateTopLevelSourceDescriptions mirrors the schema's allOf rule that
// requires sourceDescriptions whenever any operation declares
// sourceDescription. Without this check the per-operation "references unknown
// sourceDescription" diagnostic still fires but points at the wrong field.
func (d *Document) validateTopLevelSourceDescriptions(result *ValidationResult) {
	if len(d.SourceDescriptions) > 0 {
		return
	}
	for i, op := range d.Operations {
		if op == nil {
			continue
		}
		if op.SourceDescription != "" {
			result.addError("sourceDescriptions", fmt.Sprintf("is required when any operation declares sourceDescription; operations[%d].sourceDescription is %q", i, op.SourceDescription))
			return
		}
	}
}

func (d *Document) validateDocumentReferences(idx *documentIndex, result *ValidationResult) {
	for i, sd := range d.SourceDescriptions {
		if sd != nil {
			sd.validate(fmt.Sprintf("sourceDescriptions[%d]", i), result)
		}
	}
	for i, op := range d.Operations {
		if op != nil {
			op.validate(fmt.Sprintf("operations[%d]", i), idx, result)
		}
	}
	for i, wf := range d.Workflows {
		if wf != nil {
			wf.validate(fmt.Sprintf("workflows[%d]", i), idx, result)
		}
	}
	for i, trigger := range d.Triggers {
		path := fmt.Sprintf("triggers[%d]", i)
		if trigger == nil {
			result.addError(path, "is nil")
			continue
		}
		trigger.validate(path, idx, result)
	}
	seenResultNames := make(map[string]bool)
	for i, resultDecl := range d.Results {
		resultPath := fmt.Sprintf("results[%d]", i)
		if resultDecl == nil {
			result.addError(resultPath, "is nil")
			continue
		}
		resultDecl.validate(resultPath, idx, seenResultNames, result)
	}
	if d.Components != nil {
		d.Components.validate("components", result)
	}
}

func (i *Info) validate(path string, result *ValidationResult) {
	if i.Title == "" {
		result.addError(path+".title", "is required")
	}
	if i.Version == "" {
		result.addError(path+".version", "is required")
	}
}

func (s *SourceDescription) validate(path string, result *ValidationResult) {
	if s.Name == "" {
		result.addError(path+".name", "is required")
	} else if !sourceDescriptionNamePattern.MatchString(s.Name) {
		result.addError(path+".name", fmt.Sprintf("must match pattern ^[A-Za-z0-9_-]+$; got %s", s.Name))
	}
	if s.URL == "" {
		result.addError(path+".url", "is required")
	}
	if s.Type != "" && s.Type != SourceDescriptionTypeOpenAPI {
		result.addError(path+".type", fmt.Sprintf("%q is not valid (must be openapi)", s.Type))
	}
}
