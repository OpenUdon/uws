package contenttrust

import (
	"testing"
)

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

func TestParseUWS111ExpressionsAndFieldScopedNumericLiterals(t *testing.T) {
	for _, expression := range []string{
		`$response.body.customer.name`,
		`$response.body.customer.name == "Ada"`,
	} {
		if _, ok := parseExpressionForVersion(expression, "1.11.0", false); !ok {
			t.Errorf("UWS 1.11 response dot-walk rejected: %s", expression)
		}
		if _, ok := parseExpressionForVersion(expression, "1.10.0", false); ok {
			t.Errorf("UWS 1.10 accepted UWS 1.11 response dot-walk: %s", expression)
		}
	}

	if _, ok := parseExpressionForVersion(`$batchIndex`, "1.11.0", true); !ok {
		t.Error("UWS 1.11 loop batch index rejected")
	}
	for _, tc := range []struct {
		version string
		inLoop  bool
	}{
		{version: "1.11.0"},
		{version: "1.10.0", inLoop: true},
	} {
		if _, ok := parseExpressionForVersion(`$batchIndex`, tc.version, tc.inLoop); ok {
			t.Errorf("$batchIndex accepted for version %s inLoop=%v", tc.version, tc.inLoop)
		}
	}

	validNumbers := []string{"0", "-0", "12", "-12", "0.25", "-12.5", "1e3", "1E+3", "-2.5e-4"}
	for _, number := range validNumbers {
		if _, ok := parseExpressionValue(number, "1.11.0", false, true); !ok {
			t.Errorf("valid JSON number rejected for numeric field: %s", number)
		}
		if _, ok := parseExpressionValue(number, "1.11.0", false, false); ok {
			t.Errorf("bare JSON number accepted outside numeric field: %s", number)
		}
	}
	for _, number := range []string{"+1", ".5", "01", "1.", "1e", "--1", "NaN", "Infinity"} {
		if _, ok := parseExpressionValue(number, "1.11.0", false, true); ok {
			t.Errorf("invalid JSON number accepted for numeric field: %s", number)
		}
	}
}

func TestNumericWaitLiteralIsNotValidForAwaitPredicate(t *testing.T) {
	if _, ok := parseExpressionValue("30", "1.11.0", false, numericLiteralAllowedForField("wait", "await")); ok {
		t.Fatal("bare numeric await predicate was accepted as a delay")
	}
	if _, ok := parseExpressionValue("30", "1.11.0", false, numericLiteralAllowedForField("wait", "sequence")); !ok {
		t.Fatal("bare numeric non-await wait literal was rejected")
	}
	if _, ok := parseExpressionValue("2", "1.11.0", false, numericLiteralAllowedForField("batchSize", "loop")); !ok {
		t.Fatal("bare numeric batchSize literal was rejected")
	}
}
