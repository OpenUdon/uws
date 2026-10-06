package expressions

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestPublishedGrammarSourcesAndComparisons(t *testing.T) {
	sources := []string{"$response.statusCode", "$response.headers.Content-Type", "$response.body", "$response.body#/a~1b/%7E0", "$response.body.items.00.name", "$outputs.result.value", "$steps.fetch.outputs.result.value", "$variables.name", "$trigger.event", "$inputs", "$inputs.name", "$item.name", "$index", "$batchIndex"}
	for _, text := range sources {
		for _, suffix := range []string{"", " == null", " != false", " <= 900719925474099312345", " >= 1e+09", " < -1.2300", " > $variables.threshold", " == \"a >= b\""} {
			t.Run(text+suffix, func(t *testing.T) {
				p, err := Parse(text+suffix, Context{Version: "1.12.0", InLoop: true})
				if err != nil || p.Source() != text {
					t.Fatalf("parse: %v %#v", err, p)
				}
			})
		}
	}
	p, err := Parse("$variables.number == 1.2300", Context{})
	if err != nil || p.Literal() != json.Number("1.2300") {
		t.Fatalf("numeric lexeme lost: %v %#v", err, p)
	}
	p, err = Parse("1e+09", Context{Field: Wait})
	if err != nil || !p.IsNumber() || p.Literal() != json.Number("1e+09") {
		t.Fatalf("numeric field parse: %v %#v", err, p)
	}
}
func TestVersionAndFieldContexts(t *testing.T) {
	for _, c := range []struct {
		text    string
		context Context
		want    error
	}{
		{"$batchIndex", Context{Version: "1.10.0", InLoop: true}, ErrVersion},
		{"$batchIndex", Context{}, ErrContext},
		{"$response.body.name", Context{Version: "1.10.0"}, ErrVersion},
		{"$variables.x == $response.body.name", Context{Version: "1.10.0"}, ErrVersion},
		{"10", Context{Version: "1.10.0", Field: Wait}, ErrVersion},
		{"10", Context{Field: Predicate}, ErrContext},
		{"10", Context{}, ErrContext},
		{"$variables.n > 0", Context{Field: Wait}, ErrContext},
		{"$variables.n > 0", Context{Field: BatchSize}, ErrContext},
		{"$inputs", Context{Version: "01.12.0"}, ErrVersion},
		{"$inputs", Context{Version: "2.0.0"}, ErrVersion},
		{"$inputs", Context{Version: "1.12"}, ErrVersion},
	} {
		_, err := Parse(c.text, c.context)
		if !errors.Is(err, c.want) {
			t.Errorf("%q context=%+v got %v want %v", c.text, c.context, err, c.want)
		}
	}
	for _, field := range []Field{Wait, BatchSize} {
		if _, err := Parse("2.5", Context{Version: "1.11.0", Field: field}); err != nil {
			t.Fatal(err)
		}
	}
}
func TestRejectExtensionsAndMalformedSyntaxWithoutValues(t *testing.T) {
	for _, text := range []string{"", " $inputs", "$inputs ", "$variables", "$variables.x ==  true", "$variables.x  == true", "$variables.x==true", "$variables.x == []", "$variables.x == {}", "$variables.x == \"private-canary\" trailing", "$variables.x == 01", "($variables.x)", "$variables.x && true", "expr($variables.x)", "$steps.x.value", "$response.headers.a.b", "$inputs.a..b", "$response.body#%2Fa", "$response.body#/bad~2", "$response.body#/bad%q0", "$response.body#/plain space"} {
		_, err := Parse(text, Context{InLoop: true})
		if err == nil {
			t.Errorf("accepted %q", text)
		} else if strings.Contains(err.Error(), "private-canary") {
			t.Fatal("error exposed literal")
		}
	}
	if _, err := Parse(strings.Repeat("a", MaxBytes+1), Context{}); !errors.Is(err, ErrLimit) {
		t.Fatal(err)
	}
}
func TestPointerDecodedTokensAndEscapes(t *testing.T) {
	for _, c := range []struct {
		input string
		want  []string
	}{
		{"#", nil}, {"#/", []string{""}}, {"#/a~1b/~0/%E2%98%83", []string{"a/b", "~", "☃"}},
		{"#/%2F/x", []string{"", "", "x"}}, {"#/%7E1", []string{"/"}},
	} {
		got, err := Pointer(c.input)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Fatalf("%s: %v %v", c.input, got, err)
		}
	}
	for _, s := range []string{"/a", "#%2Fa", "#/a~", "#/a%", "#/%7E2", "#/é"} {
		if _, err := Pointer(s); err == nil {
			t.Errorf("accepted %q", s)
		}
	}
}
