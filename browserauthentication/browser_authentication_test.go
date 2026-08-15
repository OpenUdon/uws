package browserauthentication

import "testing"

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
