package browserauthentication

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOperationExtensionRoundTrip(t *testing.T) {
	want := &OperationAuthentication{
		Profile: "./browser-authentication/member.yaml",
		Flow:    "member_login",
		Session: "member_portal",
		CredentialBindings: map[string]string{
			"username": "member_username",
			"password": "member_password",
		},
	}
	var extensions map[string]any
	if err := SetAuthenticationExtension(&extensions, want); err != nil {
		t.Fatal(err)
	}
	got, ok, err := ReadAuthenticationExtension(extensions)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || got.Profile != want.Profile || got.Session != want.Session || got.CredentialBindings["password"] != "member_password" {
		t.Fatalf("round trip = %#v, %v", got, ok)
	}
}

func TestSessionExtensionRoundTrip(t *testing.T) {
	var extensions map[string]any
	if err := SetSessionExtension(&extensions, &OperationSession{Session: "member_portal"}); err != nil {
		t.Fatal(err)
	}
	got, ok, err := ReadSessionExtension(extensions)
	if err != nil || !ok || got.Session != "member_portal" {
		t.Fatalf("round trip = %#v, %v, %v", got, ok, err)
	}
}

func TestContextProfileWireRoundTrip(t *testing.T) {
	want := Profile{
		Profile: ContextProfileName,
		Contexts: map[string]Context{
			"idp_popup": {Kind: "popup", Parent: "main", Origin: "https://login.example.test"},
		},
		Flows: map[string]Flow{
			"login": {
				Sequence: []Step{
					{NavigateTarget: &NavigateStep{URL: "https://login.example.test/start", Context: "idp_popup"}},
					{Click: &ClickStep{Locator: Locator{Role: "button", Name: "Continue"}, Context: "main", OpensContext: "idp_popup"}},
				},
				Success: SuccessCondition{Origin: "https://members.example.test", Path: "/dashboard", Locator: Locator{Role: "heading", Name: "Dashboard"}},
			},
		},
	}
	encoded, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got Profile
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	step := got.Flows["login"].Sequence[0]
	if step.NavigateTarget == nil || step.NavigateTarget.Context != "idp_popup" || step.Navigate != "" {
		t.Fatalf("JSON navigation round trip = %#v", step)
	}
	yamlData, err := yaml.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var yamlGot Profile
	if err := yaml.Unmarshal(yamlData, &yamlGot); err != nil {
		t.Fatal(err)
	}
	if yamlGot.Flows["login"].Sequence[0].NavigateTarget == nil || yamlGot.Flows["login"].Sequence[0].NavigateTarget.URL == "" {
		t.Fatalf("YAML navigation round trip = %#v", yamlGot.Flows["login"].Sequence[0])
	}
}

func TestHistoricalNavigateStringWireRemainsStable(t *testing.T) {
	step := Step{Navigate: "https://members.example.test/login"}
	encoded, err := json.Marshal(step)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != `{"navigate":"https://members.example.test/login"}` {
		t.Fatalf("encoded historical step = %s", encoded)
	}
	var got Step
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if got.Navigate != step.Navigate || got.NavigateTarget != nil {
		t.Fatalf("historical navigation round trip = %#v", got)
	}
}
