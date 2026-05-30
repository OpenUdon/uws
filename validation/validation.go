package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/OpenUdon/uws/convert"
	"github.com/OpenUdon/uws/uws1"
	"github.com/OpenUdon/uws/versions"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

// ValidateFile validates one JSON or YAML document against the supplied JSON Schema.
func ValidateFile(schemaPath, documentPath string) error {
	schema, err := compileSchema(schemaPath)
	if err != nil {
		return err
	}
	value, err := loadSchemaValue(documentPath)
	if err != nil {
		return err
	}
	if err := schema.Validate(value); err != nil {
		return fmt.Errorf("%s: %w", documentPath, err)
	}
	return nil
}

// LoadDocumentFile loads a UWS document from JSON, YAML, or HCL based on extension.
func LoadDocumentFile(path string) (*uws1.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	doc, err := parseDocument(path, data)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// ValidateDocumentFile validates a UWS document with the versioned schema and semantic validator.
func ValidateDocumentFile(path string) (*uws1.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	doc, err := parseDocument(path, data)
	if err != nil {
		return nil, err
	}
	version := strings.TrimSpace(doc.UWS)
	if version == "" {
		version = "1.0.0"
	}
	schemaPath := versions.PathForVersion(filepath.Dir(path), version)
	if isRawSchemaDocument(path) {
		if err := ValidateFile(schemaPath, path); err != nil {
			return nil, err
		}
	} else {
		if err := validateDocumentAgainstSchema(schemaPath, doc, path); err != nil {
			return nil, err
		}
	}
	if err := doc.Validate(); err != nil {
		return nil, fmt.Errorf("validate UWS document %s: %w", path, err)
	}
	return doc, nil
}

// IsArtifactFile reports whether path names a JSON or YAML UWS artifact.
func IsArtifactFile(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".uws.json") || strings.HasSuffix(lower, ".uws.yaml") || strings.HasSuffix(lower, ".uws.yml")
}

// CollectArtifactFiles returns sorted JSON and YAML UWS artifact paths under root.
func CollectArtifactFiles(root string) ([]string, error) {
	var files []string
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if IsArtifactFile(path) {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	slices.Sort(files)
	return files, nil
}

func compileSchema(path string) (*jsonschema.Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read schema %s: %w", path, err)
	}
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse schema %s: %w", path, err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(path, doc); err != nil {
		return nil, fmt.Errorf("add schema resource %s: %w", path, err)
	}
	schema, err := compiler.Compile(path)
	if err != nil {
		return nil, fmt.Errorf("compile schema %s: %w", path, err)
	}
	return schema, nil
}

func loadSchemaValue(path string) (any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read document %s: %w", path, err)
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		value, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("parse JSON document %s: %w", path, err)
		}
		return value, nil
	case ".yaml", ".yml":
		var value any
		if err := yaml.Unmarshal(data, &value); err != nil {
			return nil, fmt.Errorf("parse YAML document %s: %w", path, err)
		}
		return normalizeYAML(value), nil
	default:
		return nil, fmt.Errorf("unsupported document extension %q", filepath.Ext(path))
	}
}

func parseDocument(path string, data []byte) (*uws1.Document, error) {
	var doc uws1.Document
	var err error
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		err = convert.UnmarshalJSON(data, &doc)
	case ".yaml", ".yml":
		err = convert.UnmarshalYAML(data, &doc)
	case ".hcl":
		err = convert.UnmarshalHCL(data, &doc)
	default:
		err = fmt.Errorf("unsupported UWS document extension %q", filepath.Ext(path))
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func validateDocumentAgainstSchema(schemaPath string, doc *uws1.Document, documentPath string) error {
	schema, err := compileSchema(schemaPath)
	if err != nil {
		return err
	}
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document %s for schema validation: %w", documentPath, err)
	}
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("parse marshaled document %s for schema validation: %w", documentPath, err)
	}
	if err := schema.Validate(value); err != nil {
		return fmt.Errorf("%s: %w", documentPath, err)
	}
	return nil
}

func isRawSchemaDocument(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json", ".yaml", ".yml":
		return true
	default:
		return false
	}
}

func normalizeYAML(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, val := range typed {
			out[key] = normalizeYAML(val)
		}
		return out
	case []any:
		for i, val := range typed {
			typed[i] = normalizeYAML(val)
		}
		return typed
	case json.Number:
		return typed
	default:
		return typed
	}
}
