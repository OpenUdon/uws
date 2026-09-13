package schemas_test

import (
	"encoding/json"
	"testing"

	"github.com/OpenUdon/uws/browserregistration"
	"github.com/OpenUdon/uws/schemas"
	"gopkg.in/yaml.v3"
)

func TestVerificationContractsAndRoundTrips(t *testing.T) {
	for _, provider := range []string{"turnstile", "recaptcha_v2", "hcaptcha"} {
		for _, activation := range []string{"before_approval", "approved_submit"} {
			root := registrationObject(t, readBrowserRegistrationFixture(t, "verification-form.json"))
			for _, raw := range root["flows"].(map[string]any) {
				v := raw.(map[string]any)["humanVerification"].(map[string]any)
				v["provider"], v["activation"] = provider, activation
				v["dependencies"].(map[string]any)["policy"] = provider + ".v1"
			}
			data := registrationJSON(t, root)
			if err := schemas.ValidateBrowserRegistrationProfile(data); err != nil {
				t.Fatal(err)
			}
			call := readBrowserRegistrationFixture(t, "verification-form-call.json")
			if err := schemas.ValidateBrowserRegistrationCallBindingForProfile(data, call, browserregistration.CallProfileNameV12); err != nil {
				t.Fatal(err)
			}
			if err := schemas.ValidateBrowserRegistrationCallBinding(data, call); err == nil {
				t.Fatal("old binding API accepted 1.2")
			}
			if _, err := schemas.BrowserRegistrationInputTemplate(data, "advertiser"); err != nil {
				t.Fatal(err)
			}
			var model browserregistration.Profile
			if err := json.Unmarshal(data, &model); err != nil {
				t.Fatal(err)
			}
			yamlData, err := yaml.Marshal(model)
			if err != nil {
				t.Fatal(err)
			}
			if err := schemas.ValidateBrowserRegistrationProfile(yamlData); err != nil {
				t.Fatal(err)
			}
			var restored browserregistration.Profile
			if err := yaml.Unmarshal(yamlData, &restored); err != nil {
				t.Fatal(err)
			}
			v := restored.Flows["advertiser"].HumanVerification
			if v == nil || v.Provider != provider || v.Activation != activation || v.Dependencies.Policy != provider+".v1" {
				t.Fatal("verification metadata lost")
			}
		}
	}
}

func TestVerificationRejectsUnsafeAndDowngradedAuthority(t *testing.T) {
	for name, mutate := range map[string]func(map[string]any, map[string]any){
		"provider":        func(_, v map[string]any) { v["provider"] = "recaptcha_v3" },
		"policy mismatch": func(_, v map[string]any) { v["dependencies"].(map[string]any)["policy"] = "hcaptcha.v1" },
		"budget":          func(_, v map[string]any) { v["dependencies"].(map[string]any)["maxRequests"] = 257 },
		"bytes":           func(_, v map[string]any) { v["dependencies"].(map[string]any)["maxResponseBytes"] = 33554433 },
		"deadline":        func(_, v map[string]any) { v["dependencies"].(map[string]any)["timeoutMs"] = 120001 },
		"zero":            func(_, v map[string]any) { v["dependencies"].(map[string]any)["timeoutMs"] = 0 },
		"destination":     func(_, v map[string]any) { v["submissionURL"] = "https://evil.example/register" },
		"selector":        func(_, v map[string]any) { v["selector"] = "#widget" },
		"token":           func(_, v map[string]any) { v["token"] = "TOKEN_CANARY" },
		"downgrade":       func(r, _ map[string]any) { r["profile"] = browserregistration.ProfileNameV11 },
		"future version":  func(r, _ map[string]any) { r["profile"] = "uws.browser-registration.1.3" },
		"captcha continue": func(r, _ map[string]any) {
			f := r["flows"].(map[string]any)["advertiser"].(map[string]any)
			f["sequence"] = append(f["sequence"].([]any), map[string]any{"human_checkpoint": map[string]any{"kind": "captcha"}})
		},
	} {
		t.Run(name, func(t *testing.T) {
			r := registrationObject(t, readBrowserRegistrationFixture(t, "verification-form.json"))
			v := r["flows"].(map[string]any)["advertiser"].(map[string]any)["humanVerification"].(map[string]any)
			mutate(r, v)
			if err := schemas.ValidateBrowserRegistrationProfile(registrationJSON(t, r)); err == nil {
				t.Fatal("unsafe profile accepted")
			}
		})
	}
	oldProfile := registrationForm(t)
	call := readBrowserRegistrationFixture(t, "verification-form-call.json")
	if err := schemas.ValidateBrowserRegistrationCallBindingForProfile(oldProfile, call, browserregistration.CallProfileNameV12); err == nil {
		t.Fatal("call accepted wrong profile version")
	}
}
