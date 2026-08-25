package browserregistration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestOperationExtensionRoundTrip(t *testing.T) {
	want := &OperationRegistration{
		Profile: "browser-registration/dedicated.yaml", Flow: "create_test_user",
		CredentialBindings: map[string]string{"identifier": "test_identifier", "password": "test_password"},
		Approval:           "register_test_user", DuplicatePrevention: "operator_attestation",
		OnDuplicate: "fail", AmbiguousOutcome: "stop_without_retry", CleanupDisposition: "delete_separately",
	}
	var extensions map[string]any
	require.NoError(t, SetRegistrationExtension(&extensions, want))
	got, ok, err := ReadRegistrationExtension(extensions)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, want, got)
}

func TestStepUnionRequiresExactlyOneArm(t *testing.T) {
	for _, data := range []string{
		`{}`,
		`{"navigate":"https://example.test/register","submit":{"locator":{"role":"button"}}}`,
	} {
		var step Step
		require.ErrorContains(t, json.Unmarshal([]byte(data), &step), "exactly one")
		require.ErrorContains(t, yaml.Unmarshal([]byte(data), &step), "exactly one")
	}
	_, err := json.Marshal(Step{})
	require.ErrorContains(t, err, "exactly one")
}

func TestVerificationPreservesExplicitZeroStabilityScore(t *testing.T) {
	var value Verification
	require.NoError(t, json.Unmarshal([]byte(`{"lastVerifiedAt":"2026-08-25T00:00:00Z","uiStabilityScore":0}`), &value))
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"uiStabilityScore":0`)
}
