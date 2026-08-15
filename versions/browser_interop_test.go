package versions_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/convert"
	"github.com/OpenUdon/uws/uws1"
	"github.com/OpenUdon/uws/versions"
	"gopkg.in/yaml.v3"
)

func TestValidateBrowserSourceProfileCanonicalFixtures(t *testing.T) {
	for _, name := range []string{"read-only.yaml", "confirmed-side-effect.yaml"} {
		t.Run(name, func(t *testing.T) {
			data := readBrowserFixture(t, name)
			if err := versions.ValidateBrowserSourceProfile(data); err != nil {
				t.Fatalf("ValidateBrowserSourceProfile() error = %v", err)
			}

			var value any
			if err := yaml.Unmarshal(data, &value); err != nil {
				t.Fatal(err)
			}
			jsonData, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			if err := versions.ValidateBrowserSourceProfile(jsonData); err != nil {
				t.Fatalf("JSON ValidateBrowserSourceProfile() error = %v", err)
			}
		})
	}
}

func TestCanonicalBrowserFixturesCoverClosedMacroVocabulary(t *testing.T) {
	data := readBrowserFixture(t, "confirmed-side-effect.yaml")
	var profile map[string]any
	if err := yaml.Unmarshal(data, &profile); err != nil {
		t.Fatal(err)
	}
	actions := profile["actions"].(map[string]any)
	action := actions["update_record"].(map[string]any)
	steps := action["sequence"].([]any)
	var macros []string
	for _, raw := range steps {
		for name := range raw.(map[string]any) {
			macros = append(macros, name)
		}
	}
	slices.Sort(macros)
	macros = slices.Compact(macros)
	want := []string{"check_radio", "click", "navigate", "select_option", "type_text", "uncheck", "wait_for"}
	if !slices.Equal(macros, want) {
		t.Fatalf("macro coverage = %#v, want %#v", macros, want)
	}
}

func TestCanonicalBrowserWorkflowBindingsValidate(t *testing.T) {
	for _, name := range []string{"read-only.uws.yaml", "confirmed-side-effect.uws.yaml"} {
		t.Run(name, func(t *testing.T) {
			data := readBrowserFixture(t, name)
			var doc uws1.Document
			if err := convert.UnmarshalYAML(data, &doc); err != nil {
				t.Fatalf("UnmarshalYAML() error = %v", err)
			}
			if err := doc.Validate(); err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
		})
	}
}

func TestValidateBrowserSourceProfileRejectsInvalidAndMultipleDocuments(t *testing.T) {
	invalid := strings.Replace(string(readBrowserFixture(t, "read-only.yaml")), "read_only", "unknown_effect", 1)
	if err := versions.ValidateBrowserSourceProfile([]byte(invalid)); err == nil {
		t.Fatal("invalid side effect unexpectedly validated")
	}
	if err := versions.ValidateBrowserSourceProfile([]byte("profile: uws.browser.1.5\n---\nprofile: uws.browser.1.5\n")); err == nil || !strings.Contains(err.Error(), "multiple YAML documents") {
		t.Fatalf("multiple documents error = %v", err)
	}
	if err := versions.ValidateBrowserSourceProfile(nil); err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("empty document error = %v", err)
	}
}

func TestBrowserSourceProfileSchemaReturnsIndependentBytes(t *testing.T) {
	first, err := versions.BrowserSourceProfileSchema("")
	if err != nil {
		t.Fatal(err)
	}
	second, err := versions.BrowserSourceProfileSchema("uws.browser.1.5")
	if err != nil {
		t.Fatal(err)
	}
	first[0] = 'x'
	if second[0] == 'x' {
		t.Fatal("schema callers share mutable backing bytes")
	}
	if _, err := versions.BrowserSourceProfileSchema("browser.9.9"); err == nil {
		t.Fatal("unsupported browser profile schema unexpectedly loaded")
	}
}

func readBrowserFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "testdata", "browser-profile", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}
