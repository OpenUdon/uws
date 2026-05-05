// Package convert provides functions to convert UWS documents between JSON and HCL formats.
package convert

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/genelet/horizon/dethcl"
	"github.com/tabilet/uws/uws1"
	"gopkg.in/yaml.v3"
)

const hclDollarKeyPrefix = "__dollar__"

var legacyDollarKeys = map[string]struct{}{
	"$ref":           {},
	"$id":            {},
	"$schema":        {},
	"$defs":          {},
	"$comment":       {},
	"$vocabulary":    {},
	"$anchor":        {},
	"$dynamicRef":    {},
	"$dynamicAnchor": {},
}

func toHCLKey(key string) string {
	if !strings.HasPrefix(key, "$") {
		return key
	}
	if _, ok := legacyDollarKeys[key]; ok {
		return "_" + key[1:]
	}
	return hclDollarKeyPrefix + key[1:]
}

func fromHCLKey(key string) string {
	if strings.HasPrefix(key, hclDollarKeyPrefix) {
		return "$" + key[len(hclDollarKeyPrefix):]
	}
	if strings.HasPrefix(key, "_") {
		candidate := "$" + key[1:]
		if _, ok := legacyDollarKeys[candidate]; ok {
			return candidate
		}
	}
	return key
}

// transformValue recursively transforms values for HCL compatibility.
func transformValue(v any, toHCL bool) any {
	switch val := v.(type) {
	case string:
		if toHCL {
			return escapeForHCL(val)
		}
		return unescapeFromHCL(val)
	case map[string]any:
		result := make(map[string]any)
		for k, v := range val {
			newKey := k
			if toHCL {
				newKey = toHCLKey(k)
			} else {
				newKey = fromHCLKey(k)
			}
			result[newKey] = transformValue(v, toHCL)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			result[i] = transformValue(item, toHCL)
		}
		return result
	default:
		return v
	}
}

func escapeForHCL(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\x00ESCAPED_N\x00")
	s = strings.ReplaceAll(s, "\\\"", "\x00ESCAPED_Q\x00")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\x00ESCAPED_N\x00", "\\\\n")
	s = strings.ReplaceAll(s, "\x00ESCAPED_Q\x00", "\\\\\"")
	return s
}

func escapeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\x00ESCAPED_N\x00")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\x00ESCAPED_N\x00", "\\\\n")
	return s
}

func unescapeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\\\\n", "\x00ESCAPED_N\x00")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\x00ESCAPED_N\x00", "\\n")
	return s
}

func unescapeFromHCL(s string) string {
	s = strings.ReplaceAll(s, "\\\\n", "\x00ESCAPED_N\x00")
	s = strings.ReplaceAll(s, "\\\\\"", "\x00ESCAPED_Q\x00")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\x00ESCAPED_N\x00", "\\n")
	s = strings.ReplaceAll(s, "\x00ESCAPED_Q\x00", "\\\"")
	return s
}

// transformDocumentForHCL transforms a UWS document's dynamic fields for HCL compatibility.
func transformDocumentForHCL(doc *uws1.Document) {
	transformDocumentDynamicFields(doc, true)
}

// transformDocumentFromHCL transforms a UWS document's dynamic fields back from HCL.
func transformDocumentFromHCL(doc *uws1.Document) {
	transformDocumentDynamicFields(doc, false)
}

func transformDocumentDynamicFields(doc *uws1.Document, toHCL bool) {
	_ = walkDocument(doc, documentWalkHandlers{
		description: func(value *string) {
			if toHCL {
				*value = escapeNewlines(*value)
				return
			}
			*value = unescapeNewlines(*value)
		},
		dynamicMap: func(_ string, value *map[string]any) error {
			*value = transformValue(*value, toHCL).(map[string]any)
			return nil
		},
	})
}

func cloneDocument(doc *uws1.Document) (*uws1.Document, error) {
	if doc == nil {
		return nil, fmt.Errorf("document is nil")
	}
	// This is a conversion-only deep copy. Runtime and execution-only state are
	// intentionally omitted because the UWS JSON tags exclude them from the wire
	// document.
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	var cloned uws1.Document
	if err := json.Unmarshal(data, &cloned); err != nil {
		return nil, err
	}
	return &cloned, nil
}

func validateHCLSerializable(doc *uws1.Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	return walkDocument(doc, documentWalkHandlers{
		extensions: rejectExtensionsForHCL,
		dynamicMap: func(path string, value *map[string]any) error {
			return rejectDynamicExtensionsForHCL(path, *value)
		},
		paramSchema: func(path string, schema *uws1.ParamSchema) error {
			for name := range schema.Properties {
				if strings.HasPrefix(name, "x-") {
					return fmt.Errorf("%s.properties.%s contains x-* extensions; UWS HCL conversion is core-only and cannot preserve extension profiles, use JSON or YAML", path, name)
				}
			}
			return nil
		},
	})
}

func rejectExtensionsForHCL(path string, extensions map[string]any) error {
	if len(extensions) == 0 {
		return nil
	}
	return fmt.Errorf("%s contains x-* extensions; UWS HCL conversion is core-only and cannot preserve extension profiles, use JSON or YAML", path)
}

func rejectDynamicExtensionsForHCL(path string, value any) error {
	switch v := value.(type) {
	case nil:
		return nil
	case map[string]any:
		for key, item := range v {
			childPath := path + "." + key
			if strings.HasPrefix(key, "x-") {
				return fmt.Errorf("%s contains x-* extensions; UWS HCL conversion is core-only and cannot preserve extension profiles, use JSON or YAML", childPath)
			}
			if err := rejectDynamicExtensionsForHCL(childPath, item); err != nil {
				return err
			}
		}
	case map[any]any:
		for key, item := range v {
			keyText := fmt.Sprint(key)
			childPath := path + "." + keyText
			if strings.HasPrefix(keyText, "x-") {
				return fmt.Errorf("%s contains x-* extensions; UWS HCL conversion is core-only and cannot preserve extension profiles, use JSON or YAML", childPath)
			}
			if err := rejectDynamicExtensionsForHCL(childPath, item); err != nil {
				return err
			}
		}
	case []any:
		for i, item := range v {
			if err := rejectDynamicExtensionsForHCL(fmt.Sprintf("%s[%d]", path, i), item); err != nil {
				return err
			}
		}
	}
	return nil
}

func toJSONCompatible(v any) any {
	switch val := v.(type) {
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, item := range val {
			result[k] = toJSONCompatible(item)
		}
		return result
	case map[any]any:
		result := make(map[string]any, len(val))
		for k, item := range val {
			result[fmt.Sprint(k)] = toJSONCompatible(item)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			result[i] = toJSONCompatible(item)
		}
		return result
	default:
		return v
	}
}

// JSONToHCL converts a UWS document from JSON format to HCL format.
func JSONToHCL(jsonData []byte) ([]byte, error) {
	var doc uws1.Document
	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return nil, err
	}
	if err := validateHCLSerializable(&doc); err != nil {
		return nil, err
	}
	transformDocumentForHCL(&doc)
	return dethcl.Marshal(&doc)
}

// JSONToYAML converts a UWS document from JSON format to YAML format.
func JSONToYAML(jsonData []byte) ([]byte, error) {
	var v any
	if err := json.Unmarshal(jsonData, &v); err != nil {
		return nil, err
	}
	return yaml.Marshal(toJSONCompatible(v))
}

// YAMLToJSON converts a UWS document from YAML format to JSON format.
func YAMLToJSON(yamlData []byte) ([]byte, error) {
	var v any
	if err := yaml.Unmarshal(yamlData, &v); err != nil {
		return nil, err
	}
	return json.Marshal(toJSONCompatible(v))
}

// YAMLToJSONIndent converts a UWS document from YAML format to indented JSON format.
func YAMLToJSONIndent(yamlData []byte, prefix, indent string) ([]byte, error) {
	var v any
	if err := yaml.Unmarshal(yamlData, &v); err != nil {
		return nil, err
	}
	return json.MarshalIndent(toJSONCompatible(v), prefix, indent)
}

// YAMLToHCL converts a UWS document from YAML format to HCL format.
func YAMLToHCL(yamlData []byte) ([]byte, error) {
	jsonData, err := YAMLToJSON(yamlData)
	if err != nil {
		return nil, err
	}
	return JSONToHCL(jsonData)
}

// HCLToJSON converts a UWS document from HCL format to JSON format.
func HCLToJSON(hclData []byte) ([]byte, error) {
	var doc uws1.Document
	if err := dethcl.Unmarshal(hclData, &doc); err != nil {
		return nil, err
	}
	transformDocumentFromHCL(&doc)
	return json.Marshal(&doc)
}

// HCLToJSONIndent converts a UWS document from HCL format to indented JSON format.
func HCLToJSONIndent(hclData []byte, prefix, indent string) ([]byte, error) {
	var doc uws1.Document
	if err := dethcl.Unmarshal(hclData, &doc); err != nil {
		return nil, err
	}
	transformDocumentFromHCL(&doc)
	return json.MarshalIndent(&doc, prefix, indent)
}

// HCLToYAML converts a UWS document from HCL format to YAML format.
func HCLToYAML(hclData []byte) ([]byte, error) {
	jsonData, err := HCLToJSON(hclData)
	if err != nil {
		return nil, err
	}
	return JSONToYAML(jsonData)
}

// MarshalHCL marshals a UWS document to HCL format.
func MarshalHCL(doc *uws1.Document) ([]byte, error) {
	cloned, err := cloneDocument(doc)
	if err != nil {
		return nil, err
	}
	if err := validateHCLSerializable(cloned); err != nil {
		return nil, err
	}
	transformDocumentForHCL(cloned)
	return dethcl.Marshal(cloned)
}

// UnmarshalHCL unmarshals HCL data into a UWS document.
func UnmarshalHCL(hclData []byte, doc *uws1.Document) error {
	if err := dethcl.Unmarshal(hclData, doc); err != nil {
		return err
	}
	transformDocumentFromHCL(doc)
	return nil
}

// MarshalJSON marshals a UWS document to JSON format.
func MarshalJSON(doc *uws1.Document) ([]byte, error) {
	return json.Marshal(doc)
}

// MarshalJSONIndent marshals a UWS document to indented JSON format.
func MarshalJSONIndent(doc *uws1.Document, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(doc, prefix, indent)
}

// UnmarshalJSON unmarshals JSON data into a UWS document.
func UnmarshalJSON(jsonData []byte, doc *uws1.Document) error {
	return json.Unmarshal(jsonData, doc)
}

// MarshalYAML marshals a UWS document to YAML format.
func MarshalYAML(doc *uws1.Document) ([]byte, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	return JSONToYAML(jsonData)
}

// UnmarshalYAML unmarshals YAML data into a UWS document.
func UnmarshalYAML(yamlData []byte, doc *uws1.Document) error {
	jsonData, err := YAMLToJSON(yamlData)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, doc)
}
