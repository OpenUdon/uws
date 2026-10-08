package hcl

import (
	"context"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"testing"
)

func TestLexicalBudget(t *testing.T) {
	ctx := context.Background()
	for _, delimiters := range [][2]string{{"(", ")"}, {"[", "]"}, {"{x=", "}"}} {
		for _, n := range []int{100, 101} {
			data := []byte("uws = " + strings.Repeat(delimiters[0], n) + "0" + strings.Repeat(delimiters[1], n))
			if (preflightHCL(ctx, data) != nil) != (n > 100) {
				t.Fatalf("%s depth %d", delimiters[0], n)
			}
		}
	}
	for _, data := range []string{
		"# " + strings.Repeat("([{", 200) + "\nuws = \"1.13.0\"\n",
		"/* " + strings.Repeat("([{", 200) + " */\nuws = \"1.13.0\"\n",
		"uws = \"" + strings.Repeat("([{", 200) + "$${literal} %%{literal}\"\n",
		"uws = <<END\n" + strings.Repeat("([{", 200) + "\nEND\n",
	} {
		if err := preflightHCL(ctx, []byte(data)); err != nil {
			t.Fatal(err)
		}
	}
	for _, data := range []string{`uws = "${1}"`, `uws = "%{if true}x%{endif}"`} {
		if preflightHCL(ctx, []byte(data)) == nil {
			t.Fatal("active template accepted")
		}
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if preflightHCL(cancelled, []byte(`uws = "1.13.0"`)) == nil {
		t.Fatal("cancellation ignored")
	}
}

func TestHostileDepthSubprocess(t *testing.T) {
	if os.Getenv("UWS_HCL_DEPTH_CHILD") == "1" {
		debug.SetMaxStack(1 << 20)
		for _, data := range []string{
			"uws = " + strings.Repeat("(", 20000) + "1" + strings.Repeat(")", 20000),
			"uws = " + strings.Repeat("[", 20000) + "1" + strings.Repeat("]", 20000),
			"uws = " + strings.Repeat("{x=", 20000) + "1" + strings.Repeat("}", 20000),
			"uws = \"${" + strings.Repeat("(", 20000) + "1" + strings.Repeat(")", 20000) + "}\"",
			"uws = \"" + strings.Repeat("%{if true}", 20000) + "x" + strings.Repeat("%{endif}", 20000) + "\"",
			"uws = " + strings.Repeat("true ? 0 : ", 20000) + "0\n",
			"uws = " + strings.Repeat("true ? ", 20000) + "0" + strings.Repeat(" : 0", 20000) + "\n",
		} {
			if _, err := Import(context.Background(), []byte(data)); err == nil {
				t.Fatal("hostile input accepted")
			}
		}
		source := Source{Format: JSON, Bytes: []byte(`{"variables":{"deep":` + strings.Repeat("[", 20000) + "0" + strings.Repeat("]", 20000) + "}}")}
		if _, err := Render(context.Background(), source, Options{Revision: fixtureRevision}); err == nil {
			t.Fatal("hostile source JSON accepted")
		}
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestHostileDepthSubprocess$")
	cmd.Env = append(os.Environ(), "UWS_HCL_DEPTH_CHILD=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("child: %v\n%s", err, out)
	}
}
