package schemas

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
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

func TestBrowserSourceProfileVersionSelectionKeepsOptInAndDefaultCompatibility(t *testing.T) {
	defaultSchema, err := BrowserSourceProfileSchema("")
	require.NoError(t, err)
	var defaultDoc map[string]any
	require.NoError(t, json.Unmarshal(defaultSchema, &defaultDoc))
	require.Equal(t, "uws.browser.1.8", defaultDoc["properties"].(map[string]any)["profile"].(map[string]any)["const"])
	for _, version := range []string{"1.5", "1.6", "1.7", "1.8", "1.9"} {
		data, err := BrowserSourceProfileSchema(version)
		require.NoError(t, err)
		require.NotEmpty(t, data)
	}
	for _, version := range []string{"1.5", "1.6", "1.7", "1.8"} {
		profile := validBrowserProfile()
		profile["profile"] = "uws.browser." + version
		require.NoError(t, ValidateBrowserSourceProfile(mustBrowserJSON(t, profile)), "browser %s", version)
	}
	require.NoError(t, ValidateBrowserSourceProfile(mustBrowserJSON(t, readBrowser19TemplateFixture(t, "text-sinks-1.9.yaml"))))
}

func TestBrowser19TemplatesAcceptEscapesAndSafeText(t *testing.T) {
	profile := browser19TestProfile()
	action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
	action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{
		"id":      map[string]any{"type": "integer", "default": int64(browser19MaxSafeInteger)},
		"comment": map[string]any{"type": "string", "default": "reviewed text"},
		"choice":  map[string]any{"type": "boolean"},
	}}
	action["sequence"] = []any{
		map[string]any{"navigate": "/records/{{id}}?sample={{{{tag}}}}"},
		map[string]any{"type_text": map[string]any{"locator": map[string]any{"role": "textbox", "name": "Price {USD}"}, "value": "Price {USD}: {{{{literal}}}} {{comment}}"}},
		map[string]any{"select_option": map[string]any{"locator": map[string]any{"role": "combobox"}, "value": "{{choice}}"}},
	}
	action["confirmationPolicy"] = map[string]any{"required": false, "prompt": "Review {{id}}: {{{{record}}}}"}
	require.NoError(t, validateBrowser19Templates(profile))

	tokens, err := parseBrowser19TemplateTokens("a {{{{name}}}} {{id}}")
	require.NoError(t, err)
	require.Equal(t, []browser19TemplateKind{browser19LiteralOpen, browser19LiteralClose, browser19Parameter}, []browser19TemplateKind{tokens[0].kind, tokens[1].kind, tokens[2].kind})
	require.Equal(t, "id", tokens[2].name)
}

func TestBrowser19TextSinkFixtures(t *testing.T) {
	valid := readBrowser19TemplateFixture(t, "text-sinks-1.9.yaml")
	require.NoError(t, validateBrowser19Templates(valid))
	for _, name := range []string{"text-sinks-1.9-invalid-control.yaml", "text-sinks-1.9-invalid-template.yaml"} {
		t.Run(name, func(t *testing.T) {
			fixture := readBrowser19TemplateFixture(t, name)
			require.Error(t, validateBrowser19Templates(fixture))
		})
	}
}

func TestBrowser19TemplatesRejectMalformedOrMisplacedTokens(t *testing.T) {
	tests := map[string]func(map[string]any){
		"unmatched opening": func(action map[string]any) {
			action["sequence"] = []any{map[string]any{"type_text": map[string]any{
				"locator": map[string]any{"role": "textbox"}, "value": "{{id",
			}}}
		},
		"unmatched closing": func(action map[string]any) {
			action["sequence"] = []any{map[string]any{"type_text": map[string]any{
				"locator": map[string]any{"role": "textbox"}, "value": "id}}",
			}}}
		},
		"nested placeholder": func(action map[string]any) {
			action["sequence"] = []any{map[string]any{"type_text": map[string]any{
				"locator": map[string]any{"role": "textbox"}, "value": "{{outer{{inner}}}}",
			}}}
		},
		"valid template outside sink": func(action map[string]any) { action["description"] = "{{id}}" },
		"escape outside sink":         func(action map[string]any) { action["description"] = "{{{{literal}}}}" },
		"escape in query name": func(action map[string]any) {
			action["sequence"] = []any{map[string]any{"navigate": "/?{{{{key}}}}=value"}}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			profile := browser19TestProfile()
			action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
			mutate(action)
			require.Error(t, validateBrowser19Templates(profile))
		})
	}
}

func TestBrowser19RejectsUnsafeTypeTextAndPrompts(t *testing.T) {
	tests := map[string]func(map[string]any){
		"type_text newline": func(action map[string]any) {
			action["sequence"] = []any{map[string]any{"type_text": map[string]any{"locator": map[string]any{"role": "textbox"}, "value": "first\nsecond"}}}
		},
		"type_text bidi control": func(action map[string]any) {
			action["sequence"] = []any{map[string]any{"type_text": map[string]any{"locator": map[string]any{"role": "textbox"}, "value": "approved\u202Etxt"}}}
		},
		"prompt newline": func(action map[string]any) {
			action["confirmationPolicy"] = map[string]any{"required": true, "prompt": "approve\nnow"}
		},
		"prompt bidi control": func(action map[string]any) {
			action["confirmationPolicy"] = map[string]any{"required": true, "prompt": "approve\u2066now"}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			profile := browser19TestProfile()
			action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
			mutate(action)
			require.Error(t, validateBrowser19Templates(profile))
		})
	}

	profile := browser19TestProfile()
	action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
	action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{
		"comment": map[string]any{"type": "string", "default": "line\nbreak"},
	}}
	action["sequence"] = []any{map[string]any{"type_text": map[string]any{"locator": map[string]any{"role": "textbox"}, "value": "{{comment}}"}}}
	require.ErrorContains(t, validateBrowser19Templates(profile), "default for parameter")
}

func TestBrowser19EnforcesSafeIntegerBoundaries(t *testing.T) {
	for _, value := range []any{int64(-browser19MaxSafeInteger), int64(browser19MaxSafeInteger), float64(-browser19MaxSafeInteger), float64(browser19MaxSafeInteger)} {
		require.True(t, isBrowser19SafeInteger(value), "value %#v", value)
	}
	for _, value := range []any{int64(browser19MaxSafeInteger + 1), int64(-browser19MaxSafeInteger - 1), float64(1.5), math.Inf(1)} {
		require.False(t, isBrowser19SafeInteger(value), "value %#v", value)
	}

	profile := browser19TestProfile()
	action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
	action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{
		"id": map[string]any{"type": "integer", "default": int64(browser19MaxSafeInteger + 1)},
	}}
	require.ErrorContains(t, validateBrowser19Templates(profile), "safe integer")
}

func browser18TestProfile() map[string]any {
	profile := validBrowserProfile()
	profile["profile"] = "uws.browser.1.8"
	return profile
}

func browser19TestProfile() map[string]any {
	profile := validBrowserProfile()
	profile["profile"] = "uws.browser.1.9"
	return profile
}

func readBrowser19TemplateFixture(t *testing.T, name string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "testdata", "browser-profile", name))
	require.NoError(t, err)
	encoded, err := decodeSingleJSONOrYAMLDocument(data)
	require.NoError(t, err)
	var profile map[string]any
	require.NoError(t, json.Unmarshal(encoded, &profile))
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
