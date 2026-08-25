package schemas_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/schemas"
)

func TestBrowserRegistrationProfileCanonicalFixture(t *testing.T) {
	data := readBrowserRegistrationFixture(t, "dedicated-test-user.yaml")
	if err := schemas.ValidateBrowserRegistrationProfile(data); err != nil {
		t.Fatalf("ValidateBrowserRegistrationProfile() error = %v", err)
	}
}

func TestBrowserRegistrationCallCanonicalFixture(t *testing.T) {
	data := readBrowserRegistrationFixture(t, "dedicated-test-user-call.json")
	if err := schemas.ValidateBrowserRegistrationCallSupplement(data); err != nil {
		t.Fatalf("ValidateBrowserRegistrationCallSupplement() error = %v", err)
	}
}

func TestBrowserRegistrationRejectsUnsafeOrRetryableDocuments(t *testing.T) {
	valid := string(readBrowserRegistrationFixture(t, "dedicated-test-user.yaml"))
	tests := map[string]string{
		"inline credential field":    strings.Replace(valid, "kind: password", "kind: password\n    value: {}", 1),
		"unsafe origin":              strings.Replace(valid, "https://app.example.test", "http://app.example.test", 1),
		"undeclared slot":            strings.Replace(valid, "slot: password", "slot: missing", 1),
		"missing submit":             strings.Replace(valid, "      - submit:\n          locator: {role: button, name: Register}\n", "", 1),
		"duplicate submit":           strings.Replace(valid, "      - submit:\n", "      - submit:\n          locator: {role: button}\n      - submit:\n", 1),
		"missing human effect":       strings.Replace(valid, ", requires_human_verification", "", 1),
		"success outside origins":    strings.Replace(valid, "origin: https://app.example.test", "origin: https://other.example.test", 1),
		"navigation outside origins": strings.Replace(valid, "navigate: https://app.example.test/register", "navigate: https://other.example.test/register", 1),
		"relative navigation":        strings.Replace(valid, "navigate: https://app.example.test/register", "navigate: /register", 1),
		"confirmation false":         strings.Replace(valid, "required: true", "required: false", 1),
		"trailing document":          valid + "\n---\nprofile: uws.browser-registration.1.0\n",
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			if err := schemas.ValidateBrowserRegistrationProfile([]byte(data)); err == nil {
				t.Fatal("invalid profile unexpectedly validated")
			}
		})
	}
}

func TestBrowserRegistrationCallRejectsUnsafeOrWeakenedControls(t *testing.T) {
	valid := string(readBrowserRegistrationFixture(t, "dedicated-test-user-call.json"))
	tests := map[string]string{
		"traversal":         strings.Replace(valid, "browser-registration/dedicated-test-user.yaml", "../secret.yaml", 1),
		"inline values":     strings.Replace(valid, `"approval":`, `"credentialValues":{},"approval":`, 1),
		"noncanonical path": strings.Replace(valid, "browser-registration/dedicated-test-user.yaml", "browser-registration//dedicated-test-user.yaml", 1),
		"blind retry":       strings.Replace(valid, "stop_without_retry", "retry", 1),
		"duplicate allowed": strings.Replace(valid, `"onDuplicate": "fail"`, `"onDuplicate": "continue"`, 1),
		"missing approval":  strings.Replace(valid, "    \"approval\": \"register_dedicated_test_user\",\n", "", 1),
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			if err := schemas.ValidateBrowserRegistrationCallSupplement([]byte(data)); err == nil {
				t.Fatal("invalid call unexpectedly validated")
			}
		})
	}
}

func TestBrowserRegistrationSchemaBytesAreIndependent(t *testing.T) {
	first, err := schemas.BrowserRegistrationProfileSchema("")
	if err != nil {
		t.Fatal(err)
	}
	second, err := schemas.BrowserRegistrationProfileSchema("uws.browser-registration.1.0")
	if err != nil {
		t.Fatal(err)
	}
	first[0] = 'x'
	if second[0] == 'x' {
		t.Fatal("schema callers share mutable backing bytes")
	}
	first, err = schemas.BrowserRegistrationCallSupplementSchema("")
	if err != nil {
		t.Fatal(err)
	}
	second, err = schemas.BrowserRegistrationCallSupplementSchema("uws.browser-registration-call.1.0")
	if err != nil {
		t.Fatal(err)
	}
	first[0] = 'x'
	if second[0] == 'x' {
		t.Fatal("call schema callers share mutable backing bytes")
	}
}

func readBrowserRegistrationFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "testdata", "browser-registration", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}
