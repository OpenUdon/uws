package contenttrust

import "testing"

func TestParseNormativeExpressionsExactly(t *testing.T) {
	valid := []string{
		`$inputs.name`,
		`$inputs.name == "contains == operator"`,
		`$inputs.name != "contains == operator"`,
		`$response.body#/items/0/name`,
		`$steps.read.outputs.body.path`,
		`$trigger.enabled == true`,
	}
	for _, expression := range valid {
		if _, ok := parseExpression(expression); !ok {
			t.Errorf("normative expression rejected: %s", expression)
		}
	}

	invalid := []string{
		`$inputs.name  == "extra space"`,
		`$inputs.name == "trailing" `,
		`$response.body#/bad token`,
		`$response.body#/bad~2escape`,
		`$steps.read.outputs.name.with.dot?`,
		`prefix $inputs.name`,
	}
	for _, expression := range invalid {
		if _, ok := parseExpression(expression); ok {
			t.Errorf("non-normative expression accepted: %s", expression)
		}
	}
}
