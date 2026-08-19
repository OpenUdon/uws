package schemas_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/browserauthentication"
	"github.com/OpenUdon/uws/convert"
	"github.com/OpenUdon/uws/schemas"
	"github.com/OpenUdon/uws/uws1"
	"gopkg.in/yaml.v3"
)

func TestBrowserAuthenticationProfileCanonicalFixture(t *testing.T) {
	data := readBrowserAuthenticationFixture(t, "member-push.yaml")
	if err := schemas.ValidateBrowserAuthenticationProfile(data); err != nil {
		t.Fatalf("ValidateBrowserAuthenticationProfile() error = %v", err)
	}
}

func TestBrowserAuthenticationCallWorkflowFixture(t *testing.T) {
	data := readBrowserAuthenticationFixture(t, "member-call.uws.yaml")
	var doc uws1.Document
	if err := convert.UnmarshalYAML(data, &doc); err != nil {
		t.Fatalf("UnmarshalYAML() error = %v", err)
	}
	if err := doc.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	call, err := convertOperationExtensionToYAML(doc.Operations[0].Extensions, "x-uws-browser-authentication")
	if err != nil {
		t.Fatal(err)
	}
	if err := schemas.ValidateBrowserAuthenticationCallSupplement(call); err != nil {
		t.Fatalf("ValidateBrowserAuthenticationCallSupplement() error = %v", err)
	}
}

func TestBrowserAuthenticationProfileRejectsUnsafeAndInconsistentDocuments(t *testing.T) {
	valid := string(readBrowserAuthenticationFixture(t, "member-push.yaml"))
	tests := map[string]string{
		"inline secret":             strings.Replace(valid, "kind: password", "kind: password\n    value: hunter2", 1),
		"TOTP seed typed into page": strings.Replace(valid, "kind: password", "kind: totp_seed", 1),
		"unsafe origin":             strings.Replace(valid, "https://login.example.test", "http://login.example.test", 1),
		"undeclared slot":           strings.Replace(valid, "slot: password", "slot: missing", 1),
		"missing MFA effect":        strings.Replace(valid, "effects: [establishes_session, sends_mfa_challenge]", "effects: [establishes_session]", 1),
		"success outside origins":   strings.Replace(valid, "origin: https://members.example.test", "origin: https://other.example.test", 1),
		"navigate outside origins":  strings.Replace(valid, "navigate: https://members.example.test/", "navigate: https://other.example.test/", 1),
		"relative navigation":       strings.Replace(valid, "navigate: https://members.example.test/", "navigate: /login", 1),
		"trailing document":         valid + "\n---\nprofile: uws.browser-authentication.1.0\n",
		"malformed IPv6 origin":     strings.Replace(valid, "https://login.example.test", "https://[::1", 1),
		"malformed URL port":        strings.Replace(valid, "https://login.example.test", "https://login.example.test:bad", 1),
		"invalid learnedAt format":  strings.Replace(valid, `learnedAt: "2026-08-15T00:00:00Z"`, "learnedAt: yesterday", 1),
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			if err := schemas.ValidateBrowserAuthenticationProfile([]byte(data)); err == nil {
				t.Fatal("invalid profile unexpectedly validated")
			}
		})
	}
}

func TestBrowserAuthenticationVerificationPreservesExplicitZeroStabilityScore(t *testing.T) {
	for name, data := range map[string]string{
		"explicit zero": `{"lastVerifiedAt":"2026-08-15T00:00:00Z","successfulRuns":1,"uiStabilityScore":0}`,
		"omitted":       `{"lastVerifiedAt":"2026-08-15T00:00:00Z","successfulRuns":1}`,
	} {
		t.Run(name, func(t *testing.T) {
			var value browserauthentication.Verification
			if err := json.Unmarshal([]byte(data), &value); err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			hasScore := strings.Contains(string(encoded), `"uiStabilityScore"`)
			if want := name == "explicit zero"; hasScore != want {
				t.Fatalf("encoded verification = %s, score present = %t, want %t", encoded, hasScore, want)
			}
		})
	}
}

func TestBrowserAuthenticationProfileBoundAndIndependentSchemaBytes(t *testing.T) {
	oversized := make([]byte, (1<<20)+1)
	if err := schemas.ValidateBrowserAuthenticationProfile(oversized); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("oversized error = %v", err)
	}
	first, err := schemas.BrowserAuthenticationProfileSchema("")
	if err != nil {
		t.Fatal(err)
	}
	second, err := schemas.BrowserAuthenticationProfileSchema("uws.browser-authentication.1.1")
	if err != nil {
		t.Fatal(err)
	}
	first[0] = 'x'
	if second[0] == 'x' {
		t.Fatal("schema callers share mutable backing bytes")
	}
}

func TestBrowserAuthenticationCallRejectsUnsafeProfilePathAndInlineFields(t *testing.T) {
	for name, data := range map[string]string{
		"traversal":       `{"x-uws-browser-authentication":{"profile":"../secret.yaml","flow":"login","session":"member","credentialBindings":{}}}`,
		"inline password": `{"x-uws-browser-authentication":{"profile":"login.yaml","flow":"login","session":"member","credentialBindings":{},"password":"secret"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			if err := schemas.ValidateBrowserAuthenticationCallSupplement([]byte(data)); err == nil {
				t.Fatal("invalid call unexpectedly validated")
			}
		})
	}
}

func readBrowserAuthenticationFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "testdata", "browser-authentication", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func convertOperationExtensionToYAML(extensions map[string]any, key string) ([]byte, error) {
	value := map[string]any{key: extensions[key]}
	return yaml.Marshal(value)
}
