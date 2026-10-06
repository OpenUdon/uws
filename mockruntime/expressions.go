package mockruntime

import (
	"context"
	"fmt"
	"github.com/OpenUdon/uws/expressions"
	"strconv"
	"strings"
)

const maxExpressionBytes = expressions.MaxBytes

// EvaluateExpression uses the shared reference evaluator. Two historical mock
// extensions remain compatible: generic numeric evaluation and an encoded root
// delimiter in response-body pointers. Strict portability diagnoses those forms
// in their actual core field context without narrowing ordinary validation.
func (r *Runtime) EvaluateExpression(ctx context.Context, text string) (any, error) {
	if r == nil || r.document == nil {
		return nil, fmt.Errorf("mock runtime is not initialized")
	}
	if ctx == nil {
		return nil, fmt.Errorf("mock expression evaluation requires a context")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(text, "$") {
		parsed, err := expressions.Parse(text, expressions.Context{Field: expressions.Wait})
		if err != nil {
			return nil, err
		}
		return parsed.Literal(), nil
	}
	evaluator, err := expressions.NewEvaluator(r.document)
	if err != nil {
		return nil, err
	}
	return evaluator.Evaluate(ctx, legacyPointerExpression(text), expressions.Value)
}

func (r *Runtime) ResolveItems(ctx context.Context, text string) ([]any, error) {
	if r == nil || r.document == nil {
		return nil, fmt.Errorf("mock runtime is not initialized")
	}
	evaluator, err := expressions.NewEvaluator(r.document)
	if err != nil {
		return nil, err
	}
	return evaluator.ResolveItems(ctx, legacyPointerExpression(text))
}

func legacyPointerExpression(text string) string {
	normalize := func(source string) string {
		const prefix = "$response.body#"
		if strings.HasPrefix(source, prefix) && len(source) >= len(prefix)+3 && strings.EqualFold(source[len(prefix):len(prefix)+3], "%2f") {
			return prefix + "/" + source[len(prefix)+3:]
		}
		return source
	}
	left, rest, compared := strings.Cut(text, " ")
	left = normalize(left)
	if !compared {
		return left
	}
	for _, operator := range []string{"==", "!=", "<=", ">=", "<", ">"} {
		if strings.HasPrefix(rest, operator+" ") {
			right := rest[len(operator)+1:]
			if strings.HasPrefix(right, "$") {
				right = normalize(right)
			}
			return left + " " + operator + " " + right
		}
	}
	return text
}

func versionAtLeast(version string, minor, patch int) bool {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}
	major, err1 := strconv.Atoi(parts[0])
	gotMinor, err2 := strconv.Atoi(parts[1])
	gotPatch, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	if major != 1 {
		return major > 1
	}
	if gotMinor != minor {
		return gotMinor > minor
	}
	return gotPatch >= patch
}
