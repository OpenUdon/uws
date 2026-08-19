// Package convert provides functions to convert UWS documents between JSON and HCL formats.
package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/OpenUdon/uws/uws1"
	"gopkg.in/yaml.v3"
)

func yamlNodeToJSONCompatible(node *yaml.Node, path string) (any, error) {
	switch node.Kind {
	case 0:
		return nil, nil
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil, nil
		}
		return yamlNodeToJSONCompatible(node.Content[0], path)
	case yaml.AliasNode:
		if node.Alias == nil {
			return nil, fmt.Errorf("%s: YAML alias has no target", path)
		}
		return yamlNodeToJSONCompatible(node.Alias, path)
	case yaml.MappingNode:
		result := make(map[string]any, len(node.Content)/2)
		seen := make(map[string]struct{}, len(node.Content)/2)
		var mergeNode *yaml.Node
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if keyNode.ShortTag() == "!!merge" {
				if mergeNode != nil {
					return nil, fmt.Errorf("%s: YAML merge key is defined more than once", path)
				}
				mergeNode = node.Content[i+1]
				continue
			}
			if keyNode.Kind != yaml.ScalarNode || keyNode.ShortTag() != "!!str" {
				return nil, fmt.Errorf("%s: YAML mapping key %q has tag %q; only string keys are supported", path, keyNode.Value, keyNode.ShortTag())
			}
			key := keyNode.Value
			if _, exists := seen[key]; exists {
				return nil, fmt.Errorf("%s: YAML mapping key %q is defined more than once", path, key)
			}
			seen[key] = struct{}{}
			converted, err := yamlNodeToJSONCompatible(node.Content[i+1], path+"."+key)
			if err != nil {
				return nil, err
			}
			result[key] = converted
		}
		if mergeNode != nil {
			merged, err := yamlMergeMappings(mergeNode, path+".<<")
			if err != nil {
				return nil, err
			}
			for _, mapping := range merged {
				for key, value := range mapping {
					if _, exists := result[key]; !exists {
						result[key] = value
					}
				}
			}
		}
		return result, nil
	case yaml.SequenceNode:
		result := make([]any, len(node.Content))
		for i, item := range node.Content {
			converted, err := yamlNodeToJSONCompatible(item, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			result[i] = converted
		}
		return result, nil
	case yaml.ScalarNode:
		if node.ShortTag() == "!!int" || node.ShortTag() == "!!float" {
			// yaml.v3 decodes integer scalars into fixed-width Go values. Retain
			// JSON-compatible numeric text so values outside uint64 remain exact.
			if json.Valid([]byte(node.Value)) {
				return json.Number(node.Value), nil
			}
		}
		var value any
		if err := node.Decode(&value); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		return value, nil
	default:
		return nil, fmt.Errorf("%s: unsupported YAML node kind %d", path, node.Kind)
	}
}

func yamlMergeMappings(node *yaml.Node, path string) ([]map[string]any, error) {
	if node.Kind == yaml.AliasNode {
		if node.Alias == nil {
			return nil, fmt.Errorf("%s: YAML alias has no target", path)
		}
		node = node.Alias
	}
	if node.Kind == yaml.SequenceNode {
		var result []map[string]any
		for i, item := range node.Content {
			mappings, err := yamlMergeMappings(item, fmt.Sprintf("%s[%d]", path, i))
			if err != nil {
				return nil, err
			}
			result = append(result, mappings...)
		}
		return result, nil
	}
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s: YAML map merge requires a map or sequence of maps", path)
	}
	converted, err := yamlNodeToJSONCompatible(node, path)
	if err != nil {
		return nil, err
	}
	mapping, ok := converted.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: YAML map merge requires a map or sequence of maps", path)
	}
	return []map[string]any{mapping}, nil
}

func yamlToJSONCompatible(yamlData []byte) (any, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(yamlData, &node); err != nil {
		return nil, err
	}
	return yamlNodeToJSONCompatible(&node, "document")
}

func jsonValueToYAMLNode(value any) (*yaml.Node, error) {
	switch typed := value.(type) {
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}, nil
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: fmt.Sprint(typed)}, nil
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: typed}, nil
	case json.Number:
		tag := "!!int"
		if strings.ContainsAny(string(typed), ".eE") {
			tag = "!!float"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: string(typed)}, nil
	case []any:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range typed {
			child, err := jsonValueToYAMLNode(item)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, child)
		}
		return node, nil
	case map[string]any:
		node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			child, err := jsonValueToYAMLNode(typed[key])
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
				child,
			)
		}
		return node, nil
	default:
		return nil, fmt.Errorf("cannot convert JSON value of type %T to YAML", value)
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
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("multiple JSON values are not supported")
		}
		return nil, err
	}
	node, err := jsonValueToYAMLNode(v)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(node)
}

// YAMLToJSON converts a UWS document from YAML format to JSON format.
func YAMLToJSON(yamlData []byte) ([]byte, error) {
	compatible, err := yamlToJSONCompatible(yamlData)
	if err != nil {
		return nil, err
	}
	return json.Marshal(compatible)
}

// YAMLToJSONIndent converts a UWS document from YAML format to indented JSON format.
func YAMLToJSONIndent(yamlData []byte, prefix, indent string) ([]byte, error) {
	compatible, err := yamlToJSONCompatible(yamlData)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(compatible, prefix, indent)
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
