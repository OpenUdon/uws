// Package convert provides functions to convert UWS documents between JSON and HCL formats.
package convert

import (
	"encoding/json"
	"fmt"

	"github.com/OpenUdon/uws/uws1"
	"gopkg.in/yaml.v3"
)

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
	return doc.MarshalHCL()
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
	if err := doc.UnmarshalHCL(hclData); err != nil {
		return nil, err
	}
	return json.Marshal(&doc)
}

// HCLToJSONIndent converts a UWS document from HCL format to indented JSON format.
func HCLToJSONIndent(hclData []byte, prefix, indent string) ([]byte, error) {
	var doc uws1.Document
	if err := doc.UnmarshalHCL(hclData); err != nil {
		return nil, err
	}
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
	return doc.MarshalHCL()
}

// UnmarshalHCL unmarshals HCL data into a UWS document.
func UnmarshalHCL(hclData []byte, doc *uws1.Document) error {
	return doc.UnmarshalHCL(hclData)
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
