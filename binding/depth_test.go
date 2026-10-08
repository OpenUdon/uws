package binding

import (
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"testing"
)

func TestTableHostileDepthSubprocess(t *testing.T) {
	if os.Getenv("UWS_TABLE_DEPTH_CHILD") == "1" {
		debug.SetMaxStack(1 << 20)
		data := []byte(`{"version":"` + TableVersion + `","sources":[],"operations":[],"unknown":` + strings.Repeat("[", 20000) + "0" + strings.Repeat("]", 20000) + "}")
		if _, err := ParseTable(data); err == nil {
			t.Fatal("hostile depth accepted")
		}
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestTableHostileDepthSubprocess$")
	cmd.Env = append(os.Environ(), "UWS_TABLE_DEPTH_CHILD=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("child: %v\n%s", err, out)
	}
}
