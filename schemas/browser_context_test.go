package schemas_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/schemas"
)

func TestContextProfileFixturesValidate(t *testing.T) {
	browser, err := os.ReadFile(filepath.Join("..", "testdata", "browser-profile", "context-popup-frame.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := schemas.ValidateBrowserSourceProfile(browser); err != nil {
		t.Fatalf("browser context fixture: %v", err)
	}
	authentication, err := os.ReadFile(filepath.Join("..", "testdata", "browser-authentication", "member-popup-frame.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := schemas.ValidateBrowserAuthenticationProfile(authentication); err != nil {
		t.Fatalf("authentication context fixture: %v", err)
	}
}

func TestProfileValidatorsDispatchAndRejectUnknownDiscriminators(t *testing.T) {
	oldBrowser := `{"profile":"uws.browser.1.5","info":{"title":"Old","origin":"https://example.test"},"observationKind":"other","evidence":{"learnedAt":"2026-08-16T00:00:00Z"},"confidence":"high","expiresAfter":"P1D","verification":{"lastVerifiedAt":"2026-08-16T00:00:00Z","successfulRuns":1},"actions":{"read":{"sequence":[{"navigate":"/"}],"sideEffects":["read_only"],"confirmationPolicy":{"required":false}}}}`
	if err := schemas.ValidateBrowserSourceProfile([]byte(oldBrowser)); err != nil {
		t.Fatalf("browser 1.5 compatibility: %v", err)
	}
	unknown := strings.Replace(oldBrowser, "uws.browser.1.5", "uws.browser.9.9", 1)
	if err := schemas.ValidateBrowserSourceProfile([]byte(unknown)); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("unknown browser discriminator error = %v", err)
	}

	oldAuth, err := os.ReadFile(filepath.Join("..", "testdata", "browser-authentication", "member-push.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := schemas.ValidateBrowserAuthenticationProfile(oldAuth); err != nil {
		t.Fatalf("authentication 1.0 compatibility: %v", err)
	}
	unknownAuth := strings.Replace(string(oldAuth), "uws.browser-authentication.1.0", "uws.browser-authentication.9.9", 1)
	if err := schemas.ValidateBrowserAuthenticationProfile([]byte(unknownAuth)); err == nil || !strings.Contains(err.Error(), "unsupported") {
		t.Fatalf("unknown authentication discriminator error = %v", err)
	}
}

func TestContextSemanticsFailClosed(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "browser-profile", "context-popup-frame.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	valid := string(data)
	tests := map[string]string{
		"unknown reference":        strings.Replace(valid, "context: detail_frame", "context: missing_frame", 1),
		"origin outside allowlist": strings.Replace(valid, "origin: https://statements.example.test", "origin: https://other.example.test", 1),
		"target context mismatch":  strings.Replace(valid, "context: main\n      - click:", "context: detail_frame\n      - click:", 1),
		"unsafe path":              strings.Replace(valid, "path: /embedded/detail", "path: /embedded/../secret", 1),
		"missing frame identity":   strings.Replace(valid, "    path: /embedded/detail\n    name: Statement detail\n", "", 1),
		"missing popup binding":    strings.Replace(valid, "          opensContext: statement_popup\n", "", 1),
		"duplicate popup binding":  strings.Replace(valid, "          opensContext: statement_popup\n", "          opensContext: statement_popup\n      - click:\n          locator: {role: link, name: Open statement again}\n          opensContext: statement_popup\n", 1),
		"popup parent cycle":       strings.Replace(valid, "parent: main", "parent: detail_frame", 1),
		"unknown field":            strings.Replace(valid, "kind: popup", "kind: popup\n    selector: body", 1),
	}
	deepContexts := "  depth_three:\n    kind: frame\n    parent: detail_frame\n    origin: https://statements.example.test\n    name: Three\n  depth_four:\n    kind: frame\n    parent: depth_three\n    origin: https://statements.example.test\n    name: Four\n  depth_five:\n    kind: frame\n    parent: depth_four\n    origin: https://statements.example.test\n    name: Five\n"
	tests["depth over four"] = strings.Replace(valid, "actions:\n", deepContexts+"actions:\n", 1)
	for name, document := range tests {
		t.Run(name, func(t *testing.T) {
			if err := schemas.ValidateBrowserSourceProfile([]byte(document)); err == nil {
				t.Fatal("invalid context profile unexpectedly validated")
			}
		})
	}
}

func TestAuthenticationContextOriginAndPathFailClosed(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "browser-authentication", "member-popup-frame.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	valid := string(data)
	tests := map[string]string{
		"success context origin mismatch": strings.Replace(valid, "      path: /dashboard\n", "      path: /dashboard\n      context: idp_popup\n", 1),
		"unsafe success path":             strings.Replace(valid, "path: /dashboard", "path: /dashboard?token=secret", 1),
		"frame origin outside allowlist":  strings.Replace(valid, "origin: https://login.example.test\n    path: /challenge/otp", "origin: https://other.example.test\n    path: /challenge/otp", 1),
		"missing popup open":              strings.Replace(valid, "          opensContext: idp_popup\n", "", 1),
	}
	for name, document := range tests {
		t.Run(name, func(t *testing.T) {
			if err := schemas.ValidateBrowserAuthenticationProfile([]byte(document)); err == nil {
				t.Fatal("invalid authentication context profile unexpectedly validated")
			}
		})
	}
}

func TestAuthenticationCallVersionsRemainAccepted(t *testing.T) {
	data := []byte(`{"x-uws-browser-authentication":{"profile":"browser-authentication/member.yaml","flow":"login","session":"member","credentialBindings":{}}}`)
	for _, profile := range []string{"uws.browser-authentication-call.1.0", "uws.browser-authentication-call.1.1"} {
		if err := schemas.ValidateBrowserAuthenticationCallSupplementForProfile(data, profile); err != nil {
			t.Fatalf("%s: %v", profile, err)
		}
	}
	if err := schemas.ValidateBrowserAuthenticationCallSupplementForProfile(data, "uws.browser-authentication-call.9.9"); err == nil {
		t.Fatal("unknown call profile unexpectedly validated")
	}
}
