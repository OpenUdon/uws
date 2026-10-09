package binding

import (
	"bytes"
	"encoding/json"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/OpenUdon/uws/internal/strictjson"
	"github.com/OpenUdon/uws/schemas"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

// BrowserOperationShape is inert native-profile metadata. Exact source and
// selected-subtree reproduction belongs to source tooling, never this package.
// Credential and registration slots contain declarations, never private values.
type BrowserOperationShape struct {
	ProfileVersion         string                       `json:"profile_version"`
	CallKind               string                       `json:"call_kind"`
	SelectedSHA256         string                       `json:"selected_sha256"`
	Origins                []string                     `json:"origins"`
	Effects                json.RawMessage              `json:"effects"`
	ConfirmationPolicy     json.RawMessage              `json:"confirmation_policy,omitempty"`
	AuthenticationRequired bool                         `json:"authentication_required"`
	CredentialSlots        []CredentialSlotShape        `json:"credential_slots"`
	RegistrationInputSlots []RegistrationInputSlotShape `json:"registration_input_slots,omitempty"`
}

type CredentialSlotShape struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Required bool   `json:"required"`
}

type RegistrationInputSlotShape struct {
	Name      string          `json:"name"`
	Kind      string          `json:"kind"`
	Required  bool            `json:"required"`
	Schema    json.RawMessage `json:"schema"`
	Condition json.RawMessage `json:"condition,omitempty"`
}

// A present empty optional slot inventory stays present. Nil represents absence.
func (b BrowserOperationShape) MarshalJSON() ([]byte, error) {
	type wire BrowserOperationShape
	data, err := json.Marshal(wire(b))
	if err != nil || b.RegistrationInputSlots == nil {
		return data, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	fields["registration_input_slots"], err = json.Marshal(b.RegistrationInputSlots)
	if err != nil {
		return nil, err
	}
	return json.Marshal(fields)
}

// Required false is evidence, not the zero value of an omitted field. Strict
// decoders retain that distinction while RawMessages preserve native numbers,
// schema extensions and optional metadata without a float round trip.
func (b *BrowserOperationShape) UnmarshalJSON(data []byte) error {
	type wire BrowserOperationShape
	var w wire
	if err := closedBrowserJSON(data, &w, "profile_version", "call_kind", "selected_sha256", "origins", "effects", "authentication_required", "credential_slots"); err != nil {
		return err
	}
	*b = BrowserOperationShape(w)
	return nil
}
func (s *CredentialSlotShape) UnmarshalJSON(data []byte) error {
	type wire CredentialSlotShape
	var w wire
	if err := closedBrowserJSON(data, &w, "name", "kind", "required"); err != nil {
		return err
	}
	*s = CredentialSlotShape(w)
	return nil
}
func (s *RegistrationInputSlotShape) UnmarshalJSON(data []byte) error {
	type wire RegistrationInputSlotShape
	var w wire
	if err := closedBrowserJSON(data, &w, "name", "kind", "required", "schema"); err != nil {
		return err
	}
	*s = RegistrationInputSlotShape(w)
	return nil
}
func closedBrowserJSON(data []byte, target any, required ...string) error {
	if strictjson.ValidateSingleValue(data) != nil {
		return ErrTable
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(data, &fields) != nil || fields == nil {
		return ErrTable
	}
	for _, raw := range fields {
		if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return ErrTable
		}
	}
	for _, key := range required {
		if raw, ok := fields[key]; !ok || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return ErrTable
		}
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if d.Decode(target) != nil {
		return ErrTable
	}
	return nil
}

var browserIdentifier = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{0,127}$`)
var browserActionKey = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func browserShapeValid(op OperationShape) bool {
	if op.Source.Kind != KindBrowser {
		return op.Browser == nil
	}
	b := op.Browser
	if b == nil || op.Protocol != "browser" || op.Selector.Kind != "id" || op.Selector.Value != op.Selector.Key || len(op.Aliases) != 0 || !browserActionKey.MatchString(op.Selector.Key) || !sourceValid(Source{ID: "selected", Kind: KindBrowser, SHA256: b.SelectedSHA256}) {
		return false
	}
	if len(b.Origins) == 0 || len(b.Origins) > 64 || b.CredentialSlots == nil || len(b.CredentialSlots) > 128 || len(b.RegistrationInputSlots) > 128 {
		return false
	}
	for i, origin := range b.Origins {
		if !canonicalBrowserOrigin(origin) || i > 0 && b.Origins[i-1] >= origin {
			return false
		}
	}
	var schemaBytes []byte
	var effectsPath string
	var err error
	switch b.CallKind {
	case "action":
		if b.ProfileVersion != "uws.browser.1.5" && b.ProfileVersion != "uws.browser.1.6" && b.ProfileVersion != "uws.browser.1.7" && b.ProfileVersion != "uws.browser.1.8" && b.ProfileVersion != "uws.browser.1.9" && b.ProfileVersion != "uws.browser.1.10" || len(b.CredentialSlots) != 0 || len(b.RegistrationInputSlots) != 0 {
			return false
		}
		schemaBytes, err = schemas.BrowserSourceProfileSchema(b.ProfileVersion)
		effectsPath = "/$defs/action/properties/sideEffects"
	case "authentication":
		if b.ProfileVersion != "uws.browser-authentication.1.0" && b.ProfileVersion != "uws.browser-authentication.1.1" || len(b.RegistrationInputSlots) != 0 || len(b.ConfirmationPolicy) != 0 || !browserIdentifier.MatchString(op.Selector.Key) {
			return false
		}
		schemaBytes, err = schemas.BrowserAuthenticationProfileSchema(b.ProfileVersion)
		effectsPath = "/$defs/flow/properties/effects"
	case "registration":
		if b.ProfileVersion != "uws.browser-registration.1.0" && b.ProfileVersion != "uws.browser-registration.1.1" && b.ProfileVersion != "uws.browser-registration.1.2" || b.ProfileVersion == "uws.browser-registration.1.0" && len(b.RegistrationInputSlots) > 0 || !browserIdentifier.MatchString(op.Selector.Key) {
			return false
		}
		schemaBytes, err = schemas.BrowserRegistrationProfileSchema(b.ProfileVersion)
		effectsPath = "/$defs/flow/properties/effects"
	default:
		return false
	}
	if err != nil || !browserFragmentValid(b.ProfileVersion, schemaBytes, effectsPath, b.Effects) {
		return false
	}
	if b.CallKind != "authentication" {
		if !browserFragmentValid(b.ProfileVersion, schemaBytes, "/$defs/confirmation-policy", b.ConfirmationPolicy) {
			return false
		}
		var effects []string
		var policy struct {
			Required bool `json:"required"`
		}
		if json.Unmarshal(b.Effects, &effects) != nil || json.Unmarshal(b.ConfirmationPolicy, &policy) != nil {
			return false
		}
		for _, effect := range effects {
			if effect != "read_only" && !policy.Required {
				return false
			}
			if effect == "read_only" && len(effects) != 1 {
				return false
			}
		}
	}
	seen := map[string]bool{}
	for _, slot := range b.CredentialSlots {
		if !browserIdentifier.MatchString(slot.Name) || seen[slot.Name] || slot.Kind != "identifier" && slot.Kind != "password" && (slot.Kind != "totp_seed" || b.CallKind != "authentication") {
			return false
		}
		seen[slot.Name] = true
	}
	credentials := seen
	seen = map[string]bool{}
	for _, slot := range b.RegistrationInputSlots {
		if b.CallKind != "registration" || !browserIdentifier.MatchString(slot.Name) || seen[slot.Name] || credentials[slot.Name] || !schemaValid(Schema{Known: true, JSON: slot.Schema}) || slot.Kind != "string" && slot.Kind != "boolean" && slot.Kind != "integer" && slot.Kind != "number" {
			return false
		}
		seen[slot.Name] = true
		var schema map[string]json.RawMessage
		var kind string
		if json.Unmarshal(slot.Schema, &schema) != nil || json.Unmarshal(schema["type"], &kind) != nil || kind != slot.Kind {
			return false
		}
		if raw, present := schema["required"]; present {
			var required bool
			if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) || json.Unmarshal(raw, &required) != nil || required != slot.Required {
				return false
			}
		}
		if raw, present := schema["requiredWhen"]; present {
			a, ea := decode(raw)
			b, eb := decode(slot.Condition)
			x, _ := json.Marshal(a)
			y, _ := json.Marshal(b)
			if ea != nil || eb != nil || !bytes.Equal(x, y) {
				return false
			}
		}
		if len(slot.Condition) != 0 {
			var condition struct {
				Slot   string          `json:"slot"`
				Equals json.RawMessage `json:"equals"`
			}
			if slot.Required || closedBrowserJSON(slot.Condition, &condition, "slot", "equals") != nil || !browserIdentifier.MatchString(condition.Slot) {
				return false
			}
			v, e := decode(condition.Equals)
			if e != nil {
				return false
			}
			switch v.(type) {
			case string, bool, json.Number:
			default:
				return false
			}
		}
	}
	for _, slot := range b.RegistrationInputSlots {
		if len(slot.Condition) == 0 {
			continue
		}
		var condition struct {
			Slot   string          `json:"slot"`
			Equals json.RawMessage `json:"equals"`
		}
		if json.Unmarshal(slot.Condition, &condition) != nil {
			return false
		}
		var controller *RegistrationInputSlotShape
		for i := range b.RegistrationInputSlots {
			if b.RegistrationInputSlots[i].Name == condition.Slot {
				controller = &b.RegistrationInputSlots[i]
				break
			}
		}
		if controller == nil {
			if op.Complete {
				return false
			}
			continue
		}
		if !controller.Required || len(controller.Condition) != 0 {
			return false
		}
		// Native slot declarations may contain scalar schema annotations plus
		// required/requiredWhen metadata. Preserve their bytes; remove only native
		// declaration keys from this temporary JSON Schema validation projection.
		var schema map[string]json.RawMessage
		if json.Unmarshal(controller.Schema, &schema) != nil || len(schema["enum"]) == 0 {
			return false
		}
		delete(schema, "required")
		delete(schema, "requiredWhen")
		delete(schema, "label")
		projection, err := json.Marshal(schema)
		if err != nil {
			return false
		}
		value, err := decode(condition.Equals)
		if err != nil || validateLiteral(Schema{Known: true, JSON: projection}, value) != Compatible {
			return false
		}
	}
	for _, in := range op.Inputs {
		if in.Location != "body" {
			return false
		}
	}
	for _, out := range op.Outputs {
		if out.Location != "body" {
			return false
		}
	}
	return true
}

func canonicalBrowserOrigin(raw string) bool {
	if !text(raw, 4096, true) {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "http" && u.Scheme != "https" || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Host != strings.ToLower(u.Host) {
		return false
	}
	if u.Hostname() == "" || strings.ContainsAny(u.Hostname(), "\\ *%") || strings.Contains(u.Hostname(), ":") && net.ParseIP(u.Hostname()) == nil {
		return false
	}
	if p := u.Port(); p != "" {
		n, e := strconv.Atoi(p)
		if e != nil || n < 1 || n > 65535 || strconv.Itoa(n) != p || u.Scheme == "http" && p == "80" || u.Scheme == "https" && p == "443" {
			return false
		}
	} else if strings.HasSuffix(u.Host, ":") {
		return false
	}
	return u.String() == raw
}

var browserFragments sync.Map

// Fragments come only from immutable embedded native schemas. External refs
// refuse; input metadata never chooses a resource loader.
func browserFragmentValid(profile string, data []byte, path string, raw json.RawMessage) bool {
	if len(raw) == 0 || len(raw) > MaxSchemaBytes || strictjson.ValidateSingleValue(raw) != nil {
		return false
	}
	key := profile + path
	var compiled *jsonschema.Schema
	if v, ok := browserFragments.Load(key); ok {
		compiled = v.(*jsonschema.Schema)
	} else {
		resource, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			return false
		}
		c := jsonschema.NewCompiler()
		c.UseLoader(denyLoader{})
		c.AssertFormat()
		if c.AddResource("https://uws.invalid/browser", resource) != nil {
			return false
		}
		compiled, err = c.Compile("https://uws.invalid/browser#" + path)
		if err != nil {
			return false
		}
		browserFragments.Store(key, compiled)
	}
	value, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
	return err == nil && compiled.Validate(value) == nil
}
