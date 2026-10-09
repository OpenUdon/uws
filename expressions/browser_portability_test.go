package expressions

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/uws1"
)

func TestBrowserBodyPortabilityUsesCoreScopeWithoutScanningNativeTemplates(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", SourceDescriptions: []*uws1.SourceDescription{{Name: "browser", Type: "browser-profile"}}, Operations: []*uws1.Operation{{OperationID: "read", SourceDescription: "browser", SourceOperationID: "read", Request: map[string]any{"body": map[string]any{"a": "$inputs.text", "b": "expr(private-canary)", "c": "{{{{literal}}}}", "d": "$batchIndex", "e": "$index"}, "query": "expr(private-canary)", "x-template": "expr(private-canary)"}, Extensions: map[string]any{"x-profile": "expr(private-canary)"}}}, Workflows: []*uws1.Workflow{{WorkflowID: "main", Type: uws1.WorkflowTypeSequence, Steps: []*uws1.Step{{StepID: "read", OperationRef: "read"}}}}}
	before, _ := json.Marshal(d)
	want := []Diagnostic{{Code: "expression.legacy-wrapper", Path: "/operations/0/request/body/b"}, {Code: "expression.context", Path: "/operations/0/request/body/d"}, {Code: "expression.context", Path: "/operations/0/request/body/e"}}
	got := CheckPortability(d)
	if !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
	data, _ := json.Marshal(got)
	if strings.Contains(string(data), "private-canary") {
		t.Fatal("value leaked")
	}
	after, _ := json.Marshal(d)
	if string(before) != string(after) {
		t.Fatal("portability changed document")
	}
	d.Workflows[0].Type = uws1.WorkflowTypeLoop
	d.Workflows[0].Items = "$variables.items"
	got = CheckPortability(d)
	if !reflect.DeepEqual(got, want[:1]) {
		t.Fatal("browser invocation loop scope lost", got)
	}
	d.Operations[0].Request["body"] = map[string]string{"typed": "expr(private-canary)"}
	got = CheckPortability(d)
	if !reflect.DeepEqual(got, []Diagnostic{{Code: "expression.legacy-wrapper", Path: "/operations/0/request/body/typed"}}) {
		t.Fatal(got)
	}
}

func TestBrowserAuthenticationAndRegistrationCallsStayOpaque(t *testing.T) {
	d := &uws1.Document{UWS: "1.13.0", Operations: []*uws1.Operation{{OperationID: "authenticate", Request: map[string]any{"body": "expr(private-canary)"}, Extensions: map[string]any{uws1.ExtensionOperationProfile: "browser-authentication", "x-uws-browser-authentication": map[string]any{"credentialBindings": map[string]any{"password": "expr(private-canary)"}}}}, {OperationID: "register", Extensions: map[string]any{uws1.ExtensionOperationProfile: "browser-registration", "x-uws-browser-registration": map[string]any{"inputBinding": "expr(private-canary)"}}}}}
	if got := CheckPortability(d); len(got) != 0 {
		t.Fatal("private declarations reinterpreted", got)
	}
}
