package schemas

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrowser18TemplatesAcceptDeclaredScalarsInApprovedSinks(t *testing.T) {
	profile := browser18TestProfile()
	action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
	action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{
		"id": map[string]any{"type": "string", "pattern": `^[A-Z]{2,5}$`}, "query": map[string]any{"type": "string"},
		"comment": map[string]any{"type": "string"}, "choice": map[string]any{"type": "boolean"},
	}}
	action["description"] = "Count {records} and use the reviewed regex field."
	action["sequence"] = []any{
		map[string]any{"navigate": "/records/{{id}}?q={{query}}"},
		map[string]any{"type_text": map[string]any{"locator": map[string]any{"role": "textbox", "name": "Price {USD}"}, "value": "{{comment}}"}},
		map[string]any{"select_option": map[string]any{"locator": map[string]any{"role": "combobox"}, "value": "{{choice}}"}},
		map[string]any{"click": map[string]any{"locator": map[string]any{"role": "button", "name": "Record {ready}"}}},
		map[string]any{"type_text": map[string]any{"locator": map[string]any{"role": "textbox"}, "value": "Price {USD}"}},
	}
	action["confirmationPolicy"] = map[string]any{"required": false, "prompt": "Review {record} {{id}}?"}
	require.NoError(t, validateBrowser18Fixture(t, profile))

	action["sequence"] = []any{map[string]any{"navigate": map[string]any{"url": "/records/{{id}}?q={{query}}", "context": "main"}}}
	require.NoError(t, validateBrowser18Fixture(t, profile))

	action["sequence"] = []any{map[string]any{"navigate": "/records/{literal}"}}
	require.NoError(t, validateBrowser18Fixture(t, profile))
}

func TestBrowser18TemplatesRejectUnsafeOrAmbiguousCases(t *testing.T) {
	tests := map[string]func(map[string]any){
		"scheme":                         func(action map[string]any) { setBrowser18Navigate(action, "https://{{id}}.example.test/") },
		"authority":                      func(action map[string]any) { setBrowser18Navigate(action, "//{{id}}/records") },
		"query name":                     func(action map[string]any) { setBrowser18Navigate(action, "/?{{id}}=value") },
		"fragment":                       func(action map[string]any) { setBrowser18Navigate(action, "/records#{{id}}") },
		"dot segment":                    func(action map[string]any) { setBrowser18Navigate(action, "/records/%2e%2e/private") },
		"malformed braces":               func(action map[string]any) { setBrowser18Navigate(action, "/records/{{id") },
		"undeclared parameter":           func(action map[string]any) { setBrowser18Navigate(action, "/records/{{missing}}") },
		"template outside sink":          func(action map[string]any) { action["description"] = "Record {{id}}" },
		"unmatched opening outside sink": func(action map[string]any) { action["description"] = "Record {{id" },
		"unmatched closing outside sink": func(action map[string]any) { action["description"] = "Record id}}" },
		"non scalar parameter": func(action map[string]any) {
			setBrowser18Navigate(action, "/records/{{id}}")
			action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "array"}}}
		},
		"parameter union": func(action map[string]any) {
			setBrowser18Navigate(action, "/records/{{id}}")
			action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": []any{"string", "null"}}}}
		},
		"forbidden locator": func(action map[string]any) {
			action["sequence"] = []any{map[string]any{"click": map[string]any{"locator": map[string]any{"role": "button", "name": "Record {{id}}"}}}}
			action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "string"}}}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			profile := browser18TestProfile()
			action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
			mutate(action)
			require.Error(t, validateBrowser18Fixture(t, profile))
		})
	}
}

func TestBrowserSourceProfileVersionSelectionKeepsHistoricalProfiles(t *testing.T) {
	latest, err := BrowserSourceProfileSchema("")
	require.NoError(t, err)
	var latestDoc map[string]any
	require.NoError(t, json.Unmarshal(latest, &latestDoc))
	require.Equal(t, "uws.browser.1.8", latestDoc["properties"].(map[string]any)["profile"].(map[string]any)["const"])
	for _, version := range []string{"1.5", "1.6", "1.7", "1.8"} {
		data, err := BrowserSourceProfileSchema(version)
		require.NoError(t, err)
		require.NotEmpty(t, data)
	}
	for _, version := range []string{"1.5", "1.6", "1.7", "1.8"} {
		profile := validBrowserProfile()
		profile["profile"] = "uws.browser." + version
		require.NoError(t, ValidateBrowserSourceProfile(mustBrowserJSON(t, profile)), "browser %s", version)
	}
}

func browser18TestProfile() map[string]any {
	profile := validBrowserProfile()
	profile["profile"] = "uws.browser.1.8"
	return profile
}

func setBrowser18Navigate(action map[string]any, value string) {
	action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{"id": map[string]any{"type": "string"}}}
	action["sequence"] = []any{map[string]any{"navigate": value}}
}

func validateBrowser18Fixture(t *testing.T, profile map[string]any) error {
	t.Helper()
	return ValidateBrowserSourceProfile(mustBrowserJSON(t, profile))
}

func mustBrowserJSON(t *testing.T, profile map[string]any) []byte {
	t.Helper()
	data, err := json.Marshal(profile)
	require.NoError(t, err)
	return data
}
