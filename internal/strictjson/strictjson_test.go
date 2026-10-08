package strictjson

import (
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"testing"
)

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

func TestNestingBudget(t *testing.T) {
	for _, n := range []int{99, 100, 101} {
		data := []byte(strings.Repeat("[", n) + "0" + strings.Repeat("]", n))
		if (ValidateSingleValue(data) != nil) != (n > 100) {
			t.Fatalf("depth %d", n)
		}
	}
	if err := ValidateSingleValue([]byte(`{"x":"[[[{{{", "n":9007199254740993e+42}`)); err != nil {
		t.Fatal(err)
	}
}

func TestHostileDepthSubprocess(t *testing.T) {
	if os.Getenv("UWS_STRICTJSON_DEPTH_CHILD") == "1" {
		debug.SetMaxStack(1 << 20)
		data := []byte(strings.Repeat("[", 20000) + "0" + strings.Repeat("]", 20000))
		if ValidateSingleValue(data) == nil {
			t.Fatal("hostile depth accepted")
		}
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestHostileDepthSubprocess$")
	cmd.Env = append(os.Environ(), "UWS_STRICTJSON_DEPTH_CHILD=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("child: %v\n%s", err, out)
	}
}
