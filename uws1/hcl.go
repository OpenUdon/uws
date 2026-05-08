package uws1

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/genelet/horizon/dethcl"
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

var (
	_ dethcl.Unmarshaler = (*Document)(nil)
	_ dethcl.Marshaler   = (*Document)(nil)

	_ dethcl.Unmarshaler = (*Info)(nil)
	_ dethcl.Unmarshaler = (*SourceDescription)(nil)
	_ dethcl.Unmarshaler = (*Operation)(nil)
	_ dethcl.Unmarshaler = (*Workflow)(nil)
	_ dethcl.Unmarshaler = (*Step)(nil)
	_ dethcl.Unmarshaler = (*Case)(nil)
	_ dethcl.Unmarshaler = (*Trigger)(nil)
	_ dethcl.Unmarshaler = (*TriggerRoute)(nil)
	_ dethcl.Unmarshaler = (*StructuralResult)(nil)
	_ dethcl.Unmarshaler = (*Components)(nil)
	_ dethcl.Unmarshaler = (*ParamSchema)(nil)
	_ dethcl.Unmarshaler = (*Criterion)(nil)
	_ dethcl.Unmarshaler = (*SuccessAction)(nil)
	_ dethcl.Unmarshaler = (*FailureAction)(nil)
	_ dethcl.Unmarshaler = (*Idempotency)(nil)
)

type documentHCLAlias Document

// UnmarshalHCL unmarshals native UWS HCL into a document.
func (d *Document) UnmarshalHCL(data []byte, labels ...string) error {
	if d == nil {
		return fmt.Errorf("document is nil")
	}
	alias := documentHCLAlias(*d)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*d = Document(alias)
	transformDynamicMapFromHCL(&d.Variables)
	transformExtensionsFromHCL(d.Extensions)
	return nil
}

// MarshalHCL marshals a UWS document to native UWS HCL.
func (d *Document) MarshalHCL() ([]byte, error) {
	cloned, err := cloneDocumentForHCL(d)
	if err != nil {
		return nil, err
	}
	if err := validateHCLSerializable(cloned); err != nil {
		return nil, err
	}
	transformDocumentForHCL(cloned)
	alias := documentHCLAlias(*cloned)
	return dethcl.Marshal(&alias)
}

func cloneDocumentForHCL(doc *Document) (*Document, error) {
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
	var cloned Document
	if err := json.Unmarshal(data, &cloned); err != nil {
		return nil, err
	}
	return &cloned, nil
}

func (i *Info) UnmarshalHCL(data []byte, labels ...string) error {
	alias := infoAlias(*i)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*i = Info(alias)
	transformDescriptionFromHCL(&i.Description)
	transformDescriptionFromHCL(&i.Summary)
	transformExtensionsFromHCL(i.Extensions)
	return nil
}

func (s *SourceDescription) UnmarshalHCL(data []byte, labels ...string) error {
	alias := sourceDescriptionAlias(*s)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*s = SourceDescription(alias)
	transformExtensionsFromHCL(s.Extensions)
	return nil
}

func (o *Operation) UnmarshalHCL(data []byte, labels ...string) error {
	alias := operationAlias(*o)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*o = Operation(alias)
	transformDescriptionFromHCL(&o.Description)
	transformDynamicMapFromHCL(&o.Request)
	transformExtensionsFromHCL(o.Extensions)
	return nil
}

func (w *Workflow) UnmarshalHCL(data []byte, labels ...string) error {
	alias := workflowAlias(*w)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*w = Workflow(alias)
	transformDescriptionFromHCL(&w.Description)
	transformExtensionsFromHCL(w.Extensions)
	return nil
}

func (s *Step) UnmarshalHCL(data []byte, labels ...string) error {
	alias := stepAlias(*s)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*s = Step(alias)
	transformDescriptionFromHCL(&s.Description)
	transformDynamicMapFromHCL(&s.Body)
	transformExtensionsFromHCL(s.Extensions)
	return nil
}

func (c *Case) UnmarshalHCL(data []byte, labels ...string) error {
	alias := caseAlias(*c)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*c = Case(alias)
	transformDynamicMapFromHCL(&c.Body)
	transformExtensionsFromHCL(c.Extensions)
	return nil
}

func (t *Trigger) UnmarshalHCL(data []byte, labels ...string) error {
	alias := triggerAlias(*t)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*t = Trigger(alias)
	transformDynamicMapFromHCL(&t.Options)
	transformExtensionsFromHCL(t.Extensions)
	return nil
}

func (t *TriggerRoute) UnmarshalHCL(data []byte, labels ...string) error {
	alias := triggerRouteAlias(*t)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*t = TriggerRoute(alias)
	transformExtensionsFromHCL(t.Extensions)
	return nil
}

func (s *StructuralResult) UnmarshalHCL(data []byte, labels ...string) error {
	alias := structuralResultAlias(*s)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*s = StructuralResult(alias)
	transformExtensionsFromHCL(s.Extensions)
	return nil
}

func (c *Components) UnmarshalHCL(data []byte, labels ...string) error {
	alias := componentsAlias(*c)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*c = Components(alias)
	transformDynamicMapFromHCL(&c.Variables)
	transformExtensionsFromHCL(c.Extensions)
	return nil
}

func (p *ParamSchema) UnmarshalHCL(data []byte, labels ...string) error {
	alias := paramSchemaAlias(*p)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*p = ParamSchema(alias)
	transformExtensionsFromHCL(p.Extensions)
	return nil
}

func (c *Criterion) UnmarshalHCL(data []byte, labels ...string) error {
	alias := criterionAlias(*c)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*c = Criterion(alias)
	transformExtensionsFromHCL(c.Extensions)
	return nil
}

func (s *SuccessAction) UnmarshalHCL(data []byte, labels ...string) error {
	alias := successActionAlias(*s)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*s = SuccessAction(alias)
	transformExtensionsFromHCL(s.Extensions)
	return nil
}

func (f *FailureAction) UnmarshalHCL(data []byte, labels ...string) error {
	alias := failureActionAlias(*f)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*f = FailureAction(alias)
	transformExtensionsFromHCL(f.Extensions)
	return nil
}

func (i *Idempotency) UnmarshalHCL(data []byte, labels ...string) error {
	alias := idempotencyAlias(*i)
	if err := dethcl.Unmarshal(data, &alias, labels...); err != nil {
		return err
	}
	*i = Idempotency(alias)
	transformExtensionsFromHCL(i.Extensions)
	return nil
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

func transformValueForHCL(v any, toHCL bool) any {
	switch val := v.(type) {
	case string:
		if toHCL {
			return escapeForHCL(val)
		}
		return unescapeFromHCL(val)
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			newKey := k
			if toHCL {
				newKey = toHCLKey(k)
			} else {
				newKey = fromHCLKey(k)
			}
			result[newKey] = transformValueForHCL(v, toHCL)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			result[i] = transformValueForHCL(item, toHCL)
		}
		return result
	default:
		return v
	}
}

func transformDescriptionFromHCL(value *string) {
	*value = unescapeNewlines(*value)
}

func transformDynamicMapFromHCL(value *map[string]any) {
	if value != nil && *value != nil {
		*value = transformValueForHCL(*value, false).(map[string]any)
	}
}

func transformExtensionsFromHCL(extensions map[string]any) {
	if extensions != nil {
		transformed := transformValueForHCL(extensions, false).(map[string]any)
		clear(extensions)
		for key, value := range transformed {
			extensions[key] = value
		}
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

func transformDocumentForHCL(doc *Document) {
	transformDocumentDynamicFields(doc, true)
}

func transformDocumentFromHCL(doc *Document) {
	transformDocumentDynamicFields(doc, false)
}

func transformDocumentDynamicFields(doc *Document, toHCL bool) {
	_ = walkDocumentHCL(doc, documentHCLWalkHandlers{
		description: func(value *string) {
			if toHCL {
				*value = escapeNewlines(*value)
				return
			}
			*value = unescapeNewlines(*value)
		},
		dynamicMap: func(_ string, value *map[string]any) error {
			*value = transformValueForHCL(*value, toHCL).(map[string]any)
			return nil
		},
		extensions: func(_ string, extensions map[string]any) error {
			if extensions != nil {
				transformed := transformValueForHCL(extensions, toHCL).(map[string]any)
				clear(extensions)
				for key, value := range transformed {
					extensions[key] = value
				}
			}
			return nil
		},
	})
}

func validateHCLSerializable(doc *Document) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}
	return walkDocumentHCL(doc, documentHCLWalkHandlers{
		dynamicMap: func(path string, value *map[string]any) error {
			return rejectDynamicExtensionsForHCL(path, *value)
		},
		paramSchema: func(path string, schema *ParamSchema) error {
			for name := range schema.Properties {
				if strings.HasPrefix(name, "x-") {
					return fmt.Errorf("%s.properties.%s contains x-* dynamic keys; place metadata on the owning UWS object with an extensions block", path, name)
				}
			}
			return nil
		},
	})
}

func rejectDynamicExtensionsForHCL(path string, value any) error {
	switch v := value.(type) {
	case nil:
		return nil
	case map[string]any:
		for key, item := range v {
			childPath := path + "." + key
			if strings.HasPrefix(key, "x-") {
				return fmt.Errorf("%s contains x-* dynamic keys; place metadata on the owning UWS object with an extensions block", childPath)
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
				return fmt.Errorf("%s contains x-* dynamic keys; place metadata on the owning UWS object with an extensions block", childPath)
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
