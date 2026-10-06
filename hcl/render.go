package hcl

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
)

type viewWriter struct {
	bytes.Buffer
	budget workBudget
}

func (w *viewWriter) put(text string) error {
	if len(text) > MaxViewBytes-w.Len() {
		return ErrCodec
	}
	w.WriteString(text)
	return nil
}
func sortedKeys(object map[string]any) []string {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (w *viewWriter) body(object map[string]any, kind reflect.Type, label string, depth int) error {
	if err := w.budget.take(depth); err != nil {
		return err
	}
	definitions := fields(kind)
	known := map[string]bool{}
	for _, f := range definitions {
		known[f.jsonName] = true
		value, present := object[f.jsonName]
		if !present || f.jsonName == label {
			continue
		}
		if f.mode == "block" && value != nil {
			blockKind := indirect(f.kind)
			if blockKind.Kind() == reflect.Slice {
				list, ok := value.([]any)
				if !ok {
					return ErrCodec
				}
				for _, item := range list {
					child, ok := item.(map[string]any)
					if !ok {
						return ErrCodec
					}
					if err := w.block(f.hclName, child, blockKind.Elem(), depth+1); err != nil {
						return err
					}
				}
				if len(list) > 0 {
					continue
				}
			} else if blockKind.Kind() == reflect.Struct {
				child, ok := value.(map[string]any)
				if !ok {
					return ErrCodec
				}
				if err := w.block(f.hclName, child, blockKind, depth+1); err != nil {
					return err
				}
				continue
			} else {
				return ErrCodec
			}
		}
		if !identifierPattern.MatchString(f.hclName) {
			return ErrCodec
		}
		if err := w.put(f.hclName + " = "); err != nil {
			return err
		}
		if err := w.value(value, f.kind, depth+1); err != nil {
			return err
		}
		if err := w.put("\n"); err != nil {
			return err
		}
	}
	extensions := map[string]any{}
	for _, key := range sortedKeys(object) {
		if known[key] {
			continue
		}
		if !strings.HasPrefix(key, "x-") {
			return ErrCodec
		}
		extensions[key] = object[key]
	}
	if len(extensions) > 0 {
		if err := w.put("extensions {\n"); err != nil {
			return err
		}
		for _, key := range sortedKeys(extensions) {
			name := encodeKey(key)
			if !identifierPattern.MatchString(name) {
				return ErrCodec
			}
			if err := w.put(name + " = "); err != nil {
				return err
			}
			if err := w.value(extensions[key], nil, depth+1); err != nil {
				return err
			}
			if err := w.put("\n"); err != nil {
				return err
			}
		}
		if err := w.put("}\n"); err != nil {
			return err
		}
	}
	return nil
}

func (w *viewWriter) block(name string, object map[string]any, kind reflect.Type, depth int) error {
	if !identifierPattern.MatchString(name) {
		return ErrCodec
	}
	if err := w.put(name); err != nil {
		return err
	}
	label := ""
	for _, f := range fields(kind) {
		if f.mode != "label" {
			continue
		}
		if label != "" {
			return ErrCodec
		}
		value, ok := object[f.jsonName].(string)
		if !ok {
			return ErrCodec
		}
		label = f.jsonName
		if err := w.put(" " + quote(value)); err != nil {
			return err
		}
	}
	if err := w.put(" {\n"); err != nil {
		return err
	}
	if err := w.body(object, kind, label, depth+1); err != nil {
		return err
	}
	return w.put("}\n")
}

func (w *viewWriter) value(value any, kind reflect.Type, depth int) error {
	if err := w.budget.take(depth); err != nil {
		return err
	}
	if kind != nil {
		kind = indirect(kind)
	}
	switch v := value.(type) {
	case nil:
		return w.put("null")
	case bool:
		if v {
			return w.put("true")
		}
		return w.put("false")
	case string:
		return w.put(quote(v))
	case json.Number:
		if !numberPattern.MatchString(v.String()) {
			return ErrCodec
		}
		return w.put(v.String())
	case []any:
		if err := w.put("["); err != nil {
			return err
		}
		var child reflect.Type
		if kind != nil && kind.Kind() == reflect.Slice {
			child = kind.Elem()
		}
		for i, item := range v {
			if i > 0 {
				if err := w.put(", "); err != nil {
					return err
				}
			}
			if err := w.value(item, child, depth+1); err != nil {
				return err
			}
		}
		return w.put("]")
	case map[string]any:
		if kind != nil && kind.Kind() == reflect.Struct {
			return w.typedObject(v, kind, depth+1)
		}
		var child reflect.Type
		if kind != nil && kind.Kind() == reflect.Map {
			child = kind.Elem()
		}
		if err := w.put("{"); err != nil {
			return err
		}
		for i, key := range sortedKeys(v) {
			if i > 0 {
				if err := w.put(", "); err != nil {
					return err
				}
			}
			if err := w.put(quote(encodeKey(key)) + " = "); err != nil {
				return err
			}
			if err := w.value(v[key], child, depth+1); err != nil {
				return err
			}
		}
		return w.put("}")
	}
	return ErrCodec
}

func (w *viewWriter) typedObject(object map[string]any, kind reflect.Type, depth int) error {
	definitions := fields(kind)
	known := map[string]bool{}
	first := true
	if err := w.put("{"); err != nil {
		return err
	}
	emit := func(name string, value any, child reflect.Type) error {
		if !first {
			if err := w.put(", "); err != nil {
				return err
			}
		}
		first = false
		if err := w.put(quote(name) + " = "); err != nil {
			return err
		}
		return w.value(value, child, depth+1)
	}
	for _, f := range definitions {
		known[f.jsonName] = true
		if value, present := object[f.jsonName]; present {
			if err := emit(f.hclName, value, f.kind); err != nil {
				return err
			}
		}
	}
	extensions := map[string]any{}
	for _, key := range sortedKeys(object) {
		if known[key] {
			continue
		}
		if !strings.HasPrefix(key, "x-") {
			return ErrCodec
		}
		extensions[key] = object[key]
	}
	if len(extensions) > 0 {
		if err := emit("extensions", extensions, nil); err != nil {
			return err
		}
	}
	return w.put("}")
}
