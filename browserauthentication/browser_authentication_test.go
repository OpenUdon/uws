package browserauthentication

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestStepUnionRequiresExactlyOneArm(t *testing.T) {
	for _, tc := range []struct {
		name string
		json string
		yaml string
	}{
		{name: "empty", json: `{}`, yaml: `{}`},
		{
			name: "multiple",
			json: `{"navigate":"https://example.test","click":{"locator":{"role":"button"}}}`,
			yaml: "navigate: https://example.test\nclick:\n  locator:\n    role: button\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var jsonStep Step
			require.ErrorContains(t, json.Unmarshal([]byte(tc.json), &jsonStep), "exactly one")
			var yamlStep Step
			require.ErrorContains(t, yaml.Unmarshal([]byte(tc.yaml), &yamlStep), "exactly one")
		})
	}

	_, err := json.Marshal(Step{})
	require.ErrorContains(t, err, "exactly one")
	_, err = yaml.Marshal(Step{Navigate: "https://example.test", Click: &ClickStep{}})
	require.ErrorContains(t, err, "exactly one")
	_, err = json.Marshal(Step{Navigate: "https://example.test", NavigateTarget: &NavigateStep{URL: "https://example.test"}})
	require.ErrorContains(t, err, "both navigate forms")

	var withNullNavigation Step
	require.NoError(t, yaml.Unmarshal([]byte("navigate: null\nclick:\n  locator:\n    role: button\n"), &withNullNavigation))
	require.NotNil(t, withNullNavigation.Click)
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

func TestStepReuseClearsPreviousNavigationUnionArm(t *testing.T) {
	for _, test := range []struct {
		name   string
		decode func([]byte, *Step) error
		first  []byte
		second []byte
	}{
		{
			name: "json object to string",
			decode: func(data []byte, step *Step) error {
				return json.Unmarshal(data, step)
			},
			first:  []byte(`{"navigate":{"url":"https://login.example.test/start","context":"idp_popup"}}`),
			second: []byte(`{"navigate":"https://members.example.test/dashboard"}`),
		},
		{
			name: "yaml object to click",
			decode: func(data []byte, step *Step) error {
				return yaml.Unmarshal(data, step)
			},
			first:  []byte("navigate:\n  url: https://login.example.test/start\n  context: idp_popup\n"),
			second: []byte("click:\n  locator:\n    role: button\n    name: Continue\n"),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var step Step
			if err := test.decode(test.first, &step); err != nil {
				t.Fatal(err)
			}
			if err := test.decode(test.second, &step); err != nil {
				t.Fatal(err)
			}
			if step.NavigateTarget != nil {
				t.Fatalf("reused step retained object navigation: %#v", step)
			}
			if test.name == "json object to string" && step.Navigate != "https://members.example.test/dashboard" {
				t.Fatalf("reused step navigation = %#v", step)
			}
			if test.name == "yaml object to click" && (step.Navigate != "" || step.Click == nil) {
				t.Fatalf("reused step click = %#v", step)
			}
		})
	}
}
