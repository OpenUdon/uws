package strictjson

import "testing"

func TestValidateSingleValue(t *testing.T) {
	for _, input := range []string{
		`{}`,
		`{"nested":[true,null,{"x":1}]}`,
		`{"unicode":"\ud83d\ude00"}`,
		`{"a":1,"\u0061":2}`,
		`{"a":{"x":1,"x":2}}`,
		`{} {}`,
		`{"a":1,}`,
		`{"unpaired":"\ud800"}`,
		`{"unpaired":"\udc00"}`,
		`{"badPair":"\ud800\u0041"}`,
	} {
		t.Run(input, func(t *testing.T) {
			err := ValidateSingleValue([]byte(input))
			wantError := input == `{"a":1,"\u0061":2}` || input == `{"a":{"x":1,"x":2}}` || input == `{} {}` || input == `{"a":1,}` || input == `{"unpaired":"\ud800"}` || input == `{"unpaired":"\udc00"}` || input == `{"badPair":"\ud800\u0041"}`
			if (err != nil) != wantError {
				t.Fatalf("ValidateSingleValue(%s) error = %v, wantError %t", input, err, wantError)
			}
		})
	}
	if err := ValidateSingleValue([]byte{'{', '"', 'x', '"', ':', '"', 0xff, '"', '}'}); err == nil {
		t.Fatal("invalid UTF-8 was accepted")
	}
}
