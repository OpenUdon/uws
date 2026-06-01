package versions

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
)

func TestBrowserProfileSchemaValidReviewedProfiles(t *testing.T) {
	schema := compileBrowserProfileSchema(t)

	require.NoError(t, validateBrowserProfile(t, schema, validBrowserProfile()))

	withStateChange := validBrowserProfile()
	action := withStateChange["actions"].(map[string]any)["read_status"].(map[string]any)
	action["sequence"] = []any{
		map[string]any{"click": map[string]any{
			"locator":  map[string]any{"role": "button", "name": "Approve"},
			"wait_for": map[string]any{"navigation": "network_idle"},
		}},
	}
	action["outputs"] = map[string]any{
		"pr_number": map[string]any{
			"type":           "string",
			"source":         "css",
			"selector":       "[data-testid='pr-number']",
			"fallbackReason": "no_structured_data",
			"validation":     map[string]any{"pattern": "^#[0-9]+$"},
		},
	}
	action["sideEffects"] = []any{"state_change"}
	action["confirmationPolicy"] = map[string]any{"required": true, "prompt": "Approve this pull request?"}
	require.NoError(t, validateBrowserProfile(t, schema, withStateChange))
}

func TestBrowserProfileSchemaRejectsMissingSafetyMetadata(t *testing.T) {
	schema := compileBrowserProfileSchema(t)

	for _, field := range []string{"observationKind", "evidence", "confidence", "expiresAfter", "verification"} {
		t.Run(field, func(t *testing.T) {
			profile := validBrowserProfile()
			delete(profile, field)
			require.Error(t, validateBrowserProfile(t, schema, profile))
		})
	}

	profile := validBrowserProfile()
	delete(profile["evidence"].(map[string]any), "learnedAt")
	require.Error(t, validateBrowserProfile(t, schema, profile))

	profile = validBrowserProfile()
	delete(profile["verification"].(map[string]any), "lastVerifiedAt")
	require.Error(t, validateBrowserProfile(t, schema, profile))
}

func TestBrowserProfileSchemaRejectsUnsafeOutputs(t *testing.T) {
	schema := compileBrowserProfileSchema(t)

	t.Run("css output requires validation", func(t *testing.T) {
		profile := profileWithOutput(map[string]any{
			"type":           "string",
			"source":         "css",
			"selector":       ".status",
			"fallbackReason": "no_a11y_region",
		})
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})

	t.Run("css output forbids locator", func(t *testing.T) {
		profile := profileWithOutput(map[string]any{
			"type":           "string",
			"source":         "css",
			"selector":       ".status",
			"fallbackReason": "no_a11y_region",
			"validation":     map[string]any{"type": "string"},
			"locator":        map[string]any{"role": "status"},
		})
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})

	t.Run("a11y output forbids selector", func(t *testing.T) {
		profile := profileWithOutput(map[string]any{
			"type":     "string",
			"source":   "a11y",
			"locator":  map[string]any{"role": "status"},
			"selector": ".status",
		})
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})

	for _, source := range []string{"jsonld", "microdata"} {
		t.Run(source+" output forbids locator", func(t *testing.T) {
			profile := profileWithOutput(map[string]any{
				"type":    "string",
				"source":  source,
				"locator": map[string]any{"role": "status"},
			})
			require.Error(t, validateBrowserProfile(t, schema, profile))
		})
	}

	t.Run("css output allows attribute", func(t *testing.T) {
		profile := profileWithOutput(map[string]any{
			"type":           "string",
			"source":         "css",
			"selector":       ".titleline > a",
			"attribute":      "href",
			"fallbackReason": "no_structured_data",
			"validation":     map[string]any{"type": "string"},
		})
		require.NoError(t, validateBrowserProfile(t, schema, profile))
	})

	for _, source := range []string{"a11y", "jsonld", "microdata"} {
		t.Run(source+" output forbids attribute", func(t *testing.T) {
			out := map[string]any{
				"type":      "string",
				"source":    source,
				"attribute": "href",
			}
			if source == "a11y" {
				out["locator"] = map[string]any{"role": "link"}
			}
			profile := profileWithOutput(out)
			require.Error(t, validateBrowserProfile(t, schema, profile))
		})
	}
}

func TestBrowserProfileSchemaRejectsUnsafeSideEffects(t *testing.T) {
	schema := compileBrowserProfileSchema(t)

	t.Run("confirmation policy required", func(t *testing.T) {
		profile := validBrowserProfile()
		action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
		delete(action, "confirmationPolicy")
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})

	t.Run("read only is exclusive", func(t *testing.T) {
		profile := validBrowserProfile()
		action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
		action["sideEffects"] = []any{"read_only", "state_change"}
		action["confirmationPolicy"] = map[string]any{"required": true}
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})

	t.Run("duplicate side effects rejected", func(t *testing.T) {
		profile := validBrowserProfile()
		action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
		action["sideEffects"] = []any{"state_change", "state_change"}
		action["confirmationPolicy"] = map[string]any{"required": true}
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})

	t.Run("side effectful actions require confirmation", func(t *testing.T) {
		profile := validBrowserProfile()
		action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
		action["sideEffects"] = []any{"state_change"}
		action["confirmationPolicy"] = map[string]any{"required": false}
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})
}

func TestBrowserProfileSchemaRejectsEmptyDurations(t *testing.T) {
	schema := compileBrowserProfileSchema(t)

	for _, duration := range []string{"P", "PT"} {
		t.Run(duration, func(t *testing.T) {
			profile := validBrowserProfile()
			profile["expiresAfter"] = duration
			require.Error(t, validateBrowserProfile(t, schema, profile))
		})
	}
}

func TestBrowserProfileSchemaAcceptsExpandedRoles(t *testing.T) {
	schema := compileBrowserProfileSchema(t)

	for _, role := range []string{"heading", "img", "list", "listitem", "combobox", "option", "tab", "table", "row", "cell", "region", "navigation", "article", "switch", "group"} {
		t.Run(role, func(t *testing.T) {
			profile := profileWithOutput(map[string]any{
				"type":    "string",
				"source":  "a11y",
				"locator": map[string]any{"role": role},
			})
			require.NoError(t, validateBrowserProfile(t, schema, profile))
		})
	}

	t.Run("unknown role rejected", func(t *testing.T) {
		profile := profileWithOutput(map[string]any{
			"type":    "string",
			"source":  "a11y",
			"locator": map[string]any{"role": "spreadsheet"},
		})
		require.Error(t, validateBrowserProfile(t, schema, profile))
	})
}

func compileBrowserProfileSchema(t *testing.T) *jsonschema.Schema {
	t.Helper()
	data, err := os.ReadFile("browser.1.5.json")
	require.NoError(t, err)
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	compiler := jsonschema.NewCompiler()
	require.NoError(t, compiler.AddResource("browser.1.5.json", doc))
	schema, err := compiler.Compile("browser.1.5.json")
	require.NoError(t, err)
	return schema
}

func validateBrowserProfile(t *testing.T, schema *jsonschema.Schema, profile map[string]any) error {
	t.Helper()
	data, err := json.Marshal(profile)
	require.NoError(t, err)
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	require.NoError(t, err)
	return schema.Validate(value)
}

func profileWithOutput(output map[string]any) map[string]any {
	profile := validBrowserProfile()
	action := profile["actions"].(map[string]any)["read_status"].(map[string]any)
	action["outputs"] = map[string]any{"status": output}
	return profile
}

func validBrowserProfile() map[string]any {
	return map[string]any{
		"profile": "uws.browser.1.5",
		"info": map[string]any{
			"title":  "Reviewed pull request profile",
			"origin": "https://github.com",
		},
		"observationKind": "accessibility_snapshot",
		"evidence": map[string]any{
			"learnedAt": "2026-05-30T00:00:00Z",
			"source":    "reviewed local profile",
		},
		"confidence":   "high",
		"expiresAfter": "P30D",
		"verification": map[string]any{
			"lastVerifiedAt":   "2026-05-30T00:00:00Z",
			"successfulRuns":   1,
			"uiStabilityScore": 0.95,
		},
		"actions": map[string]any{
			"read_status": map[string]any{
				"sequence": []any{
					map[string]any{"navigate": "/pull/1"},
					map[string]any{"wait_for": map[string]any{"role": "status", "name": "Ready"}},
				},
				"outputs": map[string]any{
					"status": map[string]any{
						"type":    "string",
						"source":  "a11y",
						"locator": map[string]any{"role": "status", "name": "Ready"},
					},
				},
				"sideEffects":        []any{"read_only"},
				"confirmationPolicy": map[string]any{"required": false},
			},
		},
	}
}
