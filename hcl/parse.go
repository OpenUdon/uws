package hcl

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"

	"github.com/OpenUdon/uws/uws1"
	hashihcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type viewReader struct {
	data   []byte
	budget workBudget
}

func hclValue(ctx context.Context, data []byte) (map[string]any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(data) == 0 || len(data) > MaxViewBytes {
		return nil, ErrCodec
	}
	if err := preflightHCL(ctx, data); err != nil {
		return nil, err
	}
	file, diags := hclsyntax.ParseConfig(data, "view.hcl", hashihcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, ErrCodec
	}
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, ErrCodec
	}
	reader := viewReader{data: data, budget: workBudget{ctx: ctx}}
	return reader.body(body, reflect.TypeOf(uws1.Document{}), 0)
}

// LexConfig is the dependency's iterative finite-state scanner, not its
// recursive expression parser. Tokens distinguish inert comments/string text
// from syntax, including escaped template markers and heredoc contents.
func preflightHCL(ctx context.Context, data []byte) error {
	tokens, diags := hclsyntax.LexConfig(data, "view.hcl", hashihcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return ErrCodec
	}
	stack := make([]hclsyntax.TokenType, 0, maxDepth)
	unary := 0
	for _, token := range tokens {
		if err := ctx.Err(); err != nil {
			return err
		}
		var end hclsyntax.TokenType
		switch token.Type {
		case hclsyntax.TokenOBrace:
			end = hclsyntax.TokenCBrace
		case hclsyntax.TokenOBrack:
			end = hclsyntax.TokenCBrack
		case hclsyntax.TokenOParen:
			end = hclsyntax.TokenCParen
		case hclsyntax.TokenOQuote:
			end = hclsyntax.TokenCQuote
		case hclsyntax.TokenOHeredoc:
			end = hclsyntax.TokenCHeredoc
		case hclsyntax.TokenTemplateInterp, hclsyntax.TokenTemplateControl:
			// These are already outside the inert subset. Refuse before the
			// parser can recurse through interpolations or template directives.
			return ErrCodec
		case hclsyntax.TokenCBrace, hclsyntax.TokenCBrack, hclsyntax.TokenCParen, hclsyntax.TokenCQuote, hclsyntax.TokenCHeredoc:
			if len(stack) == 0 || stack[len(stack)-1] != token.Type {
				return ErrCodec
			}
			stack = stack[:len(stack)-1]
		}
		if end != 0 {
			if len(stack) == maxDepth {
				return ErrCodec
			}
			stack = append(stack, end)
		}
		if token.Type == hclsyntax.TokenMinus || token.Type == hclsyntax.TokenBang {
			unary++
			if unary > maxDepth {
				return ErrCodec
			}
		} else if token.Type != hclsyntax.TokenComment && token.Type != hclsyntax.TokenNewline {
			unary = 0
		}
	}
	if len(stack) != 0 {
		return ErrCodec
	}
	return ctx.Err()
}

func (r *viewReader) body(body *hclsyntax.Body, kind reflect.Type, depth int) (map[string]any, error) {
	if err := r.budget.take(depth); err != nil {
		return nil, err
	}
	definitions := map[string]field{}
	for _, f := range fields(kind) {
		definitions[f.hclName] = f
	}
	out := map[string]any{}
	keys := make([]string, 0, len(body.Attributes))
	for key := range body.Attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, name := range keys {
		f, exists := definitions[name]
		if !exists {
			return nil, ErrCodec
		}
		value, err := r.expression(body.Attributes[name].Expr, depth+1)
		if err != nil {
			return nil, err
		}
		if f.mode == "block" && value != nil {
			if list, ok := value.([]any); !ok || len(list) != 0 || indirect(f.kind).Kind() != reflect.Slice {
				return nil, ErrCodec
			}
		}
		value, err = r.typed(value, f.kind, depth+1)
		if err != nil {
			return nil, err
		}
		out[f.jsonName] = value
	}
	seenExtensions := false
	for _, block := range body.Blocks {
		if block.Type == "extensions" {
			if seenExtensions || len(block.Labels) != 0 || len(block.Body.Blocks) > 0 {
				return nil, ErrCodec
			}
			seenExtensions = true
			for name, attr := range block.Body.Attributes {
				key := decodeKey(name)
				if !strings.HasPrefix(key, "x-") {
					return nil, ErrCodec
				}
				if _, exists := out[key]; exists {
					return nil, ErrCodec
				}
				value, err := r.expression(attr.Expr, depth+1)
				if err != nil {
					return nil, err
				}
				value, err = r.typed(value, nil, depth+1)
				if err != nil {
					return nil, err
				}
				out[key] = value
			}
			continue
		}
		f, exists := definitions[block.Type]
		if !exists || f.mode != "block" {
			return nil, ErrCodec
		}
		childKind := indirect(f.kind)
		list := childKind.Kind() == reflect.Slice
		if list {
			childKind = childKind.Elem()
		}
		child, err := r.body(block.Body, childKind, depth+1)
		if err != nil {
			return nil, err
		}
		labels := []field{}
		for _, candidate := range fields(childKind) {
			if candidate.mode == "label" {
				labels = append(labels, candidate)
			}
		}
		if len(labels) != len(block.Labels) {
			return nil, ErrCodec
		}
		for i, label := range labels {
			if _, exists := child[label.jsonName]; exists {
				return nil, ErrCodec
			}
			child[label.jsonName] = block.Labels[i]
		}
		if list {
			if old, present := out[f.jsonName]; present {
				items, ok := old.([]any)
				if !ok || len(items) == 0 {
					return nil, ErrCodec
				}
				out[f.jsonName] = append(items, child)
			} else {
				out[f.jsonName] = []any{child}
			}
		} else {
			if _, present := out[f.jsonName]; present {
				return nil, ErrCodec
			}
			out[f.jsonName] = child
		}
	}
	return out, nil
}

func (r *viewReader) expression(expr hclsyntax.Expression, depth int) (any, error) {
	if err := r.budget.take(depth); err != nil {
		return nil, err
	}
	switch value := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		if value.Val.IsNull() {
			return nil, nil
		}
		switch value.Val.Type() {
		case cty.Bool:
			return value.Val.True(), nil
		case cty.String:
			return value.Val.AsString(), nil
		case cty.Number:
			return r.number(expr)
		}
	case *hclsyntax.UnaryOpExpr:
		if value.Op == hclsyntax.OpNegate {
			if child, ok := value.Val.(*hclsyntax.LiteralValueExpr); ok && child.Val.Type() == cty.Number {
				return r.number(expr)
			}
		}
	case *hclsyntax.TemplateExpr:
		var out strings.Builder
		for _, part := range value.Parts {
			literal, ok := part.(*hclsyntax.LiteralValueExpr)
			if !ok || literal.Val.Type() != cty.String || literal.Val.IsNull() {
				return nil, ErrCodec
			}
			out.WriteString(literal.Val.AsString())
		}
		return out.String(), nil
	case *hclsyntax.TupleConsExpr:
		out := []any{}
		for _, child := range value.Exprs {
			item, err := r.expression(child, depth+1)
			if err != nil {
				return nil, err
			}
			out = append(out, item)
		}
		return out, nil
	case *hclsyntax.ObjectConsExpr:
		out := map[string]any{}
		for _, item := range value.Items {
			key, err := r.objectKey(item.KeyExpr, depth+1)
			if err != nil {
				return nil, err
			}
			if _, duplicate := out[key]; duplicate {
				return nil, ErrCodec
			}
			child, err := r.expression(item.ValueExpr, depth+1)
			if err != nil {
				return nil, err
			}
			out[key] = child
		}
		return out, nil
	}
	// No variable traversal, function, arithmetic, template interpolation,
	// comprehension or filesystem/network evaluation is admitted.
	return nil, ErrCodec
}

func (r *viewReader) number(expr hclsyntax.Expression) (any, error) {
	span := expr.Range()
	if span.Start.Byte < 0 || span.End.Byte > len(r.data) {
		return nil, ErrCodec
	}
	text := string(r.data[span.Start.Byte:span.End.Byte])
	if !numberPattern.MatchString(text) {
		return nil, ErrCodec
	}
	return json.Number(text), nil
}

func (r *viewReader) objectKey(expr hclsyntax.Expression, depth int) (string, error) {
	if key, ok := expr.(*hclsyntax.ObjectConsKeyExpr); ok {
		if !key.ForceNonLiteral {
			if traversal, ok := key.Wrapped.(*hclsyntax.ScopeTraversalExpr); ok && len(traversal.Traversal) == 1 {
				if root, ok := traversal.Traversal[0].(hashihcl.TraverseRoot); ok {
					return root.Name, nil
				}
			}
		}
		expr = key.Wrapped
	}
	value, err := r.expression(expr, depth+1)
	if err != nil {
		return "", err
	}
	key, ok := value.(string)
	if !ok {
		return "", ErrCodec
	}
	return key, nil
}

func (r *viewReader) typed(value any, kind reflect.Type, depth int) (any, error) {
	if err := r.budget.take(depth); err != nil {
		return nil, err
	}
	if kind != nil {
		kind = indirect(kind)
		if value != nil {
			switch kind.Kind() {
			case reflect.Struct, reflect.Map:
				if _, ok := value.(map[string]any); !ok {
					return nil, ErrCodec
				}
			case reflect.Slice:
				if _, ok := value.([]any); !ok {
					return nil, ErrCodec
				}
			}
		}
	}
	switch object := value.(type) {
	case map[string]any:
		out := map[string]any{}
		if kind != nil && kind.Kind() == reflect.Struct {
			definitions := map[string]field{}
			for _, f := range fields(kind) {
				definitions[f.hclName] = f
			}
			for _, key := range sortedKeys(object) {
				if key == "extensions" {
					extensions, ok := object[key].(map[string]any)
					if !ok {
						return nil, ErrCodec
					}
					for name, v := range extensions {
						name = decodeKey(name)
						if !strings.HasPrefix(name, "x-") {
							return nil, ErrCodec
						}
						if _, duplicate := out[name]; duplicate {
							return nil, ErrCodec
						}
						child, err := r.typed(v, nil, depth+1)
						if err != nil {
							return nil, err
						}
						out[name] = child
					}
					continue
				}
				f, exists := definitions[key]
				if !exists {
					return nil, ErrCodec
				}
				child, err := r.typed(object[key], f.kind, depth+1)
				if err != nil {
					return nil, err
				}
				if _, duplicate := out[f.jsonName]; duplicate {
					return nil, ErrCodec
				}
				out[f.jsonName] = child
			}
			return out, nil
		}
		var childKind reflect.Type
		if kind != nil && kind.Kind() == reflect.Map {
			childKind = kind.Elem()
		}
		for key, v := range object {
			key = decodeKey(key)
			if _, duplicate := out[key]; duplicate {
				return nil, ErrCodec
			}
			child, err := r.typed(v, childKind, depth+1)
			if err != nil {
				return nil, err
			}
			out[key] = child
		}
		return out, nil
	case []any:
		out := []any{}
		var childKind reflect.Type
		if kind != nil && kind.Kind() == reflect.Slice {
			childKind = kind.Elem()
		}
		for _, v := range object {
			child, err := r.typed(v, childKind, depth+1)
			if err != nil {
				return nil, err
			}
			out = append(out, child)
		}
		return out, nil
	}
	return value, nil
}
