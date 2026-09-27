// Package mockruntime defines the portable response-fixture format and its
// request-key helpers. It does not execute operations or persist responses.
package mockruntime

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/OpenUdon/uws/internal/strictjson"
	"github.com/OpenUdon/uws/schemas"
	"github.com/gowebpki/jcs"
)

const (
	// FixtureFormatV1 is the exact discriminator for fixture format 1.0.
	FixtureFormatV1 = "uws.mock-fixtures.1.0"

	// MaxFixtureSetBytes is the maximum encoded size accepted for a fixture set.
	MaxFixtureSetBytes = 16 << 20

	// MaxRequestBytes bounds request canonicalization before JCS processing.
	MaxRequestBytes = 1 << 20
)

var (
	operationIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	digestPattern      = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	pointerPattern     = regexp.MustCompile(`^(?:/(?:[^~/]|~[01])*)*$`)
)

// FixtureSet is one versioned collection of exact operation/request matches.
type FixtureSet struct {
	Format   string    `json:"format"`
	Fixtures []Fixture `json:"fixtures"`
}

// Fixture provides a source-neutral JSON response for one exact key pair.
type Fixture struct {
	OperationID   string            `json:"operationId"`
	RequestDigest string            `json:"requestDigest"`
	Provenance    FixtureProvenance `json:"provenance"`
	Response      json.RawMessage   `json:"response"`
}

// FixtureProvenance identifies the response origin. Recorded responses must
// explicitly identify that redaction ran and list the JSON Pointer locations
// it removed or replaced.
type FixtureProvenance struct {
	Kind             string   `json:"kind"`
	Redacted         bool     `json:"-"`
	RedactedPointers []string `json:"-"`
}

type fixtureProvenanceWire struct {
	Kind             string    `json:"kind"`
	Redacted         *bool     `json:"redacted,omitempty"`
	RedactedPointers *[]string `json:"redactedPointers,omitempty"`
}

func (p FixtureProvenance) MarshalJSON() ([]byte, error) {
	if p.Kind != "recorded" {
		return json.Marshal(struct {
			Kind string `json:"kind"`
		}{Kind: p.Kind})
	}
	redacted := p.Redacted
	pointers := p.RedactedPointers
	return json.Marshal(fixtureProvenanceWire{
		Kind:             p.Kind,
		Redacted:         &redacted,
		RedactedPointers: &pointers,
	})
}

func (p *FixtureProvenance) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var wire fixtureProvenanceWire
	if err := decoder.Decode(&wire); err != nil {
		return err
	}
	p.Kind = wire.Kind
	if wire.Redacted != nil {
		p.Redacted = *wire.Redacted
	}
	if wire.RedactedPointers != nil {
		p.RedactedPointers = *wire.RedactedPointers
	}
	return nil
}

// ResponseRedactor is an explicit caller-owned policy for converting a
// supplied response into portable fixture data. It returns the redacted JSON
// value and pointers to locations removed or replaced by that policy.
type ResponseRedactor func(response json.RawMessage) (redacted json.RawMessage, redactedPointers []string, err error)

// DecodeFixtures validates and decodes one strict JSON fixture set. Unsupported
// format versions are rejected without downgrade.
func DecodeFixtures(data []byte) (*FixtureSet, error) {
	if len(data) > MaxFixtureSetBytes {
		return nil, fmt.Errorf("mock fixture document exceeds %d bytes", MaxFixtureSetBytes)
	}
	if err := schemas.ValidateMockFixtures(data); err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var set FixtureSet
	if err := decoder.Decode(&set); err != nil {
		return nil, fmt.Errorf("decode mock fixture document: %w", err)
	}
	if err := validateFixtureSemantics(&set); err != nil {
		return nil, err
	}
	return &set, nil
}

// EncodeFixtures validates and encodes one fixture set. It performs no file
// writes; the caller decides whether and where to persist the returned bytes.
func EncodeFixtures(set *FixtureSet) ([]byte, error) {
	if set == nil {
		return nil, fmt.Errorf("mock fixture set is required")
	}
	if err := checkFixtureEncodingBudget(set); err != nil {
		return nil, err
	}
	if err := validateFixtureSemantics(set); err != nil {
		return nil, err
	}
	data, err := json.Marshal(set)
	if err != nil {
		return nil, fmt.Errorf("encode mock fixture document: %w", err)
	}
	data = append(data, '\n')
	if len(data) > MaxFixtureSetBytes {
		return nil, fmt.Errorf("mock fixture document exceeds %d bytes", MaxFixtureSetBytes)
	}
	if err := schemas.ValidateMockFixtures(data); err != nil {
		return nil, err
	}
	return data, nil
}

// checkFixtureEncodingBudget bounds the caller-provided data inspected and
// copied by EncodeFixtures before JSON encoding. The final encoded-size check
// remains authoritative because JSON escaping can expand the output.
func checkFixtureEncodingBudget(set *FixtureSet) error {
	total := 0
	add := func(size int) error {
		if size > MaxFixtureSetBytes-total {
			return fmt.Errorf("mock fixture document exceeds %d bytes", MaxFixtureSetBytes)
		}
		total += size
		return nil
	}
	if err := add(len(set.Format)); err != nil {
		return err
	}
	for i, fixture := range set.Fixtures {
		for _, size := range []int{
			len(fixture.OperationID),
			len(fixture.RequestDigest),
			len(fixture.Provenance.Kind),
			len(fixture.Response),
		} {
			if err := add(size); err != nil {
				return fmt.Errorf("fixtures[%d]: %w", i, err)
			}
		}
		for _, pointer := range fixture.Provenance.RedactedPointers {
			if err := add(len(pointer)); err != nil {
				return fmt.Errorf("fixtures[%d].provenance.redactedPointers: %w", i, err)
			}
		}
	}
	return nil
}

// CanonicalizeRequest returns the RFC 8785 canonical UTF-8 bytes for a
// resolved UWS request-binding object. A nil request is the empty object.
func CanonicalizeRequest(request any) ([]byte, error) {
	var encoded []byte
	var err error
	if request == nil {
		encoded = []byte("{}")
	} else {
		encoded, err = json.Marshal(request)
		if err != nil {
			return nil, fmt.Errorf("encode resolved request: %w", err)
		}
	}
	if len(encoded) > MaxRequestBytes {
		return nil, fmt.Errorf("resolved request exceeds %d bytes", MaxRequestBytes)
	}
	if err := strictjson.ValidateSingleValue(encoded); err != nil {
		return nil, fmt.Errorf("validate resolved request JSON: %w", err)
	}
	trimmed := bytes.TrimSpace(encoded)
	if len(trimmed) == 0 || trimmed[0] != '{' {
		return nil, fmt.Errorf("resolved request must be a JSON object")
	}
	canonical, err := jcs.Transform(trimmed)
	if err != nil {
		return nil, fmt.Errorf("canonicalize resolved request using RFC 8785: %w", err)
	}
	return canonical, nil
}

// RequestDigest returns "sha256:" followed by the lowercase SHA-256 digest of
// CanonicalizeRequest's bytes.
func RequestDigest(request any) (string, error) {
	canonical, err := CanonicalizeRequest(request)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// LookupResponse returns a copy of the exact response for (operationID,
// requestDigest). The same fixture is returned on every lookup; fixtures are
// never consumed. Duplicate keys in an in-memory set are rejected as
// ambiguous even though DecodeFixtures and EncodeFixtures already reject them.
func (set *FixtureSet) LookupResponse(operationID, requestDigest string) (json.RawMessage, bool, error) {
	if set == nil {
		return nil, false, fmt.Errorf("mock fixture set is required")
	}
	var match json.RawMessage
	for _, fixture := range set.Fixtures {
		if fixture.OperationID != operationID || fixture.RequestDigest != requestDigest {
			continue
		}
		if match != nil {
			return nil, false, fmt.Errorf("duplicate mock fixture key for operation %q", operationID)
		}
		match = bytes.Clone(fixture.Response)
	}
	return match, match != nil, nil
}

// NewRecordedFixture creates a recorded fixture only after an explicit
// redactor has processed the caller-supplied response. It does not persist the
// fixture or retain the original response.
func NewRecordedFixture(operationID, requestDigest string, response json.RawMessage, redact ResponseRedactor) (Fixture, error) {
	if err := validateFixtureKey(operationID, requestDigest); err != nil {
		return Fixture{}, err
	}
	if redact == nil {
		return Fixture{}, fmt.Errorf("an explicit response redactor is required for recorded fixtures")
	}
	if len(response) == 0 || len(response) > MaxFixtureSetBytes {
		return Fixture{}, fmt.Errorf("supplied response must be non-empty JSON no larger than %d bytes", MaxFixtureSetBytes)
	}
	if err := strictjson.ValidateSingleValue(response); err != nil {
		return Fixture{}, fmt.Errorf("validate supplied response JSON: %w", err)
	}
	redacted, pointers, err := redact(bytes.Clone(response))
	if err != nil {
		return Fixture{}, fmt.Errorf("redact supplied response: %w", err)
	}
	if len(redacted) == 0 || len(redacted) > MaxFixtureSetBytes {
		return Fixture{}, fmt.Errorf("redacted response must be non-empty JSON no larger than %d bytes", MaxFixtureSetBytes)
	}
	if err := strictjson.ValidateSingleValue(redacted); err != nil {
		return Fixture{}, fmt.Errorf("validate redacted response JSON: %w", err)
	}
	if len(pointers) > 256 {
		return Fixture{}, fmt.Errorf("redaction produced %d JSON Pointers; maximum is 256", len(pointers))
	}
	seen := make(map[string]struct{}, len(pointers))
	for i, pointer := range pointers {
		if utf8.RuneCountInString(pointer) > 512 || !pointerPattern.MatchString(pointer) {
			return Fixture{}, fmt.Errorf("redactedPointers[%d] is not a valid bounded JSON Pointer", i)
		}
		if _, exists := seen[pointer]; exists {
			return Fixture{}, fmt.Errorf("redaction returned duplicate JSON Pointer %q", pointer)
		}
		seen[pointer] = struct{}{}
	}
	return Fixture{
		OperationID:   operationID,
		RequestDigest: requestDigest,
		Provenance: FixtureProvenance{
			Kind:             "recorded",
			Redacted:         true,
			RedactedPointers: append([]string{}, pointers...),
		},
		Response: bytes.Clone(redacted),
	}, nil
}

func validateFixtureSemantics(set *FixtureSet) error {
	if set.Format != FixtureFormatV1 {
		return fmt.Errorf("unsupported mock fixture format %q", set.Format)
	}
	if len(set.Fixtures) == 0 || len(set.Fixtures) > 10000 {
		return fmt.Errorf("mock fixture set must contain between 1 and 10000 fixtures")
	}
	seen := make(map[string]struct{}, len(set.Fixtures))
	for i, fixture := range set.Fixtures {
		if err := validateFixtureKey(fixture.OperationID, fixture.RequestDigest); err != nil {
			return fmt.Errorf("fixtures[%d]: %w", i, err)
		}
		key := fixture.OperationID + "\x00" + fixture.RequestDigest
		if _, exists := seen[key]; exists {
			return fmt.Errorf("fixtures[%d]: duplicate operationId/requestDigest pair", i)
		}
		seen[key] = struct{}{}
		if !json.Valid(fixture.Response) || len(fixture.Response) == 0 {
			return fmt.Errorf("fixtures[%d].response must contain one JSON value", i)
		}
		if err := strictjson.ValidateSingleValue(fixture.Response); err != nil {
			return fmt.Errorf("fixtures[%d].response: %w", i, err)
		}
		switch fixture.Provenance.Kind {
		case "synthetic", "example":
			if fixture.Provenance.Redacted || fixture.Provenance.RedactedPointers != nil {
				return fmt.Errorf("fixtures[%d]: non-recorded provenance must omit redaction fields", i)
			}
		case "recorded":
			if !fixture.Provenance.Redacted || fixture.Provenance.RedactedPointers == nil {
				return fmt.Errorf("fixtures[%d]: recorded provenance requires explicit redaction metadata", i)
			}
		default:
			return fmt.Errorf("fixtures[%d].provenance.kind is unsupported", i)
		}
		redacted := make(map[string]struct{}, len(fixture.Provenance.RedactedPointers))
		for j, pointer := range fixture.Provenance.RedactedPointers {
			if utf8.RuneCountInString(pointer) > 512 || !pointerPattern.MatchString(pointer) {
				return fmt.Errorf("fixtures[%d].provenance.redactedPointers[%d] is invalid", i, j)
			}
			if _, exists := redacted[pointer]; exists {
				return fmt.Errorf("fixtures[%d].provenance.redactedPointers has duplicate pointer %q", i, pointer)
			}
			redacted[pointer] = struct{}{}
		}
	}
	return nil
}

func validateFixtureKey(operationID, requestDigest string) error {
	if !operationIDPattern.MatchString(operationID) {
		return fmt.Errorf("operationId must use UWS identifier characters")
	}
	if !digestPattern.MatchString(requestDigest) {
		return fmt.Errorf("requestDigest must be sha256: followed by 64 lowercase hexadecimal characters")
	}
	return nil
}
