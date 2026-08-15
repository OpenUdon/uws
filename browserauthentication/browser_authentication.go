// Package browserauthentication defines the portable UWS browser
// authentication profile and operation-extension wire types.
//
// The package contains inert metadata only. Credential values, OTP values,
// cookies, browser storage, private keys, and live session handles remain
// runtime-private and must never be placed in these structures.
package browserauthentication

import (
	"encoding/json"
	"fmt"
)

const (
	// ProfileName is the portable browser sign-in recipe profile.
	ProfileName = "uws.browser-authentication.1.0"

	// CallProfileName identifies extension-owned UWS authentication operations.
	CallProfileName = "uws.browser-authentication-call.1.0"

	// ExtensionAuthentication is the operation-level authentication-call key.
	ExtensionAuthentication = "x-uws-browser-authentication"

	// ExtensionSession is the operation-level named browser-session key.
	ExtensionSession = "x-uws-browser-session"
)

// Profile is one portable, secret-free browser sign-in recipe document.
type Profile struct {
	Profile         string                    `json:"profile" yaml:"profile"`
	Info            Info                      `json:"info" yaml:"info"`
	ObservationKind string                    `json:"observationKind" yaml:"observationKind"`
	Evidence        Evidence                  `json:"evidence" yaml:"evidence"`
	Confidence      string                    `json:"confidence" yaml:"confidence"`
	ExpiresAfter    string                    `json:"expiresAfter" yaml:"expiresAfter"`
	Verification    Verification              `json:"verification" yaml:"verification"`
	CredentialSlots map[string]CredentialSlot `json:"credentialSlots" yaml:"credentialSlots"`
	Flows           map[string]Flow           `json:"flows" yaml:"flows"`
}

// Info identifies the application and identity-provider origins covered by a
// profile. Origins are exact origins, not URL-prefix allowlists.
type Info struct {
	Title                 string   `json:"title" yaml:"title"`
	Provider              string   `json:"provider,omitempty" yaml:"provider,omitempty"`
	ApplicationOrigins    []string `json:"applicationOrigins" yaml:"applicationOrigins"`
	AuthenticationOrigins []string `json:"authenticationOrigins" yaml:"authenticationOrigins"`
}

// Evidence records when and from what reviewed observation the recipe was
// learned. The source is a non-secret provenance label, not captured content.
type Evidence struct {
	LearnedAt string `json:"learnedAt" yaml:"learnedAt"`
	Source    string `json:"source,omitempty" yaml:"source,omitempty"`
}

// Verification records the latest successful review of the recipe.
type Verification struct {
	LastVerifiedAt   string   `json:"lastVerifiedAt" yaml:"lastVerifiedAt"`
	SuccessfulRuns   int      `json:"successfulRuns" yaml:"successfulRuns"`
	UIStabilityScore *float64 `json:"uiStabilityScore,omitempty" yaml:"uiStabilityScore,omitempty"`
}

// CredentialSlot declares the class of a symbolic driver-resolved credential.
// It never carries a credential value.
type CredentialSlot struct {
	Kind string `json:"kind" yaml:"kind"`
}

// Flow is one explicitly selected sign-in alternative. MFA alternatives are
// separate flows so a runtime never guesses which challenge to use.
type Flow struct {
	Description string           `json:"description,omitempty" yaml:"description,omitempty"`
	Sequence    []Step           `json:"sequence" yaml:"sequence"`
	Effects     []string         `json:"effects" yaml:"effects"`
	Success     SuccessCondition `json:"success" yaml:"success"`
}

// Step is a closed declarative sign-in macro union. Exactly one field is set.
type Step struct {
	Navigate       string              `json:"navigate,omitempty" yaml:"navigate,omitempty"`
	TypeCredential *TypeCredentialStep `json:"type_credential,omitempty" yaml:"type_credential,omitempty"`
	Click          *ClickStep          `json:"click,omitempty" yaml:"click,omitempty"`
	Challenge      *ChallengeStep      `json:"challenge,omitempty" yaml:"challenge,omitempty"`
	WaitFor        *WaitForCondition   `json:"wait_for,omitempty" yaml:"wait_for,omitempty"`
}

// Locator is an accessibility-tree locator. CSS, XPath, coordinates, and
// arbitrary scripts are intentionally absent.
type Locator struct {
	Role  string `json:"role" yaml:"role"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Text  string `json:"text,omitempty" yaml:"text,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
}

// TypeCredentialStep fills a reviewed locator from a symbolic slot.
type TypeCredentialStep struct {
	Locator Locator `json:"locator" yaml:"locator"`
	Slot    string  `json:"slot" yaml:"slot"`
}

// ClickStep activates a reviewed accessibility locator.
type ClickStep struct {
	Locator Locator `json:"locator" yaml:"locator"`
}

// ChallengeStep performs one explicitly selected MFA handoff. Locator is
// required by OTP challenges and omitted for push/WebAuthn challenges.
type ChallengeStep struct {
	Kind    string   `json:"kind" yaml:"kind"`
	Locator *Locator `json:"locator,omitempty" yaml:"locator,omitempty"`
	Slot    string   `json:"slot,omitempty" yaml:"slot,omitempty"`
}

// WaitForCondition waits for one reviewed locator.
type WaitForCondition struct {
	Locator Locator `json:"locator" yaml:"locator"`
}

// SuccessCondition proves that sign-in established the named session.
type SuccessCondition struct {
	Origin  string  `json:"origin" yaml:"origin"`
	Locator Locator `json:"locator" yaml:"locator"`
}

// OperationAuthentication is the typed x-uws-browser-authentication payload.
type OperationAuthentication struct {
	Profile            string            `json:"profile" hcl:"profile"`
	Flow               string            `json:"flow" hcl:"flow"`
	Session            string            `json:"session" hcl:"session"`
	CredentialBindings map[string]string `json:"credentialBindings" hcl:"credentialBindings"`
}

// OperationSession selects a named session established by an authentication
// operation in the same workflow execution.
type OperationSession struct {
	Session string `json:"session" hcl:"session"`
}

// ReadAuthenticationExtension decodes x-uws-browser-authentication.
func ReadAuthenticationExtension(extensions map[string]any) (*OperationAuthentication, bool, error) {
	var out OperationAuthentication
	ok, err := readExtension(extensions, ExtensionAuthentication, &out)
	return &out, ok, err
}

// ReadSessionExtension decodes x-uws-browser-session.
func ReadSessionExtension(extensions map[string]any) (*OperationSession, bool, error) {
	var out OperationSession
	ok, err := readExtension(extensions, ExtensionSession, &out)
	return &out, ok, err
}

// SetAuthenticationExtension encodes x-uws-browser-authentication.
func SetAuthenticationExtension(dst *map[string]any, value *OperationAuthentication) error {
	return setExtension(dst, ExtensionAuthentication, value)
}

// SetSessionExtension encodes x-uws-browser-session.
func SetSessionExtension(dst *map[string]any, value *OperationSession) error {
	return setExtension(dst, ExtensionSession, value)
}

func readExtension(extensions map[string]any, key string, out any) (bool, error) {
	if len(extensions) == 0 {
		return false, nil
	}
	value, ok := extensions[key]
	if !ok {
		return false, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("marshal %s extension: %w", key, err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return false, fmt.Errorf("unmarshal %s extension: %w", key, err)
	}
	return true, nil
}

func setExtension(dst *map[string]any, key string, value any) error {
	if value == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var generic any
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}
	if *dst == nil {
		*dst = make(map[string]any)
	}
	(*dst)[key] = generic
	return nil
}
