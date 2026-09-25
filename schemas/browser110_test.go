package schemas

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrowser110CountProfileFixture(t *testing.T) {
	path := filepath.Join("..", "testdata", "browser-profile", "count-outputs-1.10.yaml")
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NoError(t, ValidateBrowserSourceProfile(data))

	profile := readBrowserProfileFixture(t, "count-outputs-1.10.yaml")
	profile["profile"] = "uws.browser.1.9"
	require.Error(t, ValidateBrowserSourceProfile(mustBrowserJSON(t, profile)), "Browser 1.9 must not accept Browser 1.10 count fields")
}

func TestBrowser110CountOutputConstraints(t *testing.T) {
	valid := map[string]any{
		"type":           "integer",
		"source":         "css",
		"selector":       ".card",
		"matchCount":     true,
		"visibility":     "rendered",
		"fallbackReason": "no_a11y_region",
		"validation":     map[string]any{"type": "integer", "minimum": 0, "maximum": browser19MaxSafeInteger},
	}
	require.NoError(t, ValidateBrowserSourceProfile(mustBrowserJSON(t, profileWithBrowser110Output(valid))))

	tests := map[string]func(map[string]any){
		"non-integer output": func(output map[string]any) { output["type"] = "number" },
		"non-css source": func(output map[string]any) {
			output["source"] = "a11y"
			output["locator"] = map[string]any{"role": "status"}
		},
		"missing visibility":          func(output map[string]any) { delete(output, "visibility") },
		"unknown visibility":          func(output map[string]any) { output["visibility"] = "viewport" },
		"attribute extraction":        func(output map[string]any) { output["attribute"] = "href" },
		"non-integer validation":      func(output map[string]any) { output["validation"] = map[string]any{"type": "number", "minimum": 0} },
		"missing validation minimum":  func(output map[string]any) { output["validation"] = map[string]any{"type": "integer"} },
		"negative validation minimum": func(output map[string]any) { output["validation"] = map[string]any{"type": "integer", "minimum": -1} },
		"validation maximum exceeds safe range": func(output map[string]any) {
			output["validation"] = map[string]any{"type": "integer", "minimum": 0, "maximum": browser19MaxSafeInteger + 1}
		},
		"empty selector": func(output map[string]any) { output["selector"] = "" },
		"scope without match count": func(output map[string]any) {
			delete(output, "matchCount")
			delete(output, "visibility")
			delete(output, "validation")
			output["within"] = "main"
		},
		"empty scope selector": func(output map[string]any) { output["within"] = "" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			output := cloneBrowser110Output(valid)
			mutate(output)
			require.Error(t, ValidateBrowserSourceProfile(mustBrowserJSON(t, profileWithBrowser110Output(output))))
		})
	}
}

func TestBrowser110InheritsBrowser19TemplateRules(t *testing.T) {
	profile := browser110TestProfile()
	action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
	action["parameters"] = map[string]any{"type": "object", "properties": map[string]any{
		"id": map[string]any{"type": "integer", "default": int64(browser19MaxSafeInteger + 1)},
	}}
	require.ErrorContains(t, ValidateBrowserSourceProfile(mustBrowserJSON(t, profile)), "safe integer")
}

func profileWithBrowser110Output(output map[string]any) map[string]any {
	profile := browser110TestProfile()
	action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
	action["outputs"] = map[string]any{"count": output}
	return profile
}

func browser110TestProfile() map[string]any {
	profile := validBrowserProfile()
	profile["profile"] = "uws.browser.1.10"
	return profile
}

func cloneBrowser110Output(output map[string]any) map[string]any {
	data, _ := json.Marshal(output)
	var clone map[string]any
	_ = json.Unmarshal(data, &clone)
	return clone
}
