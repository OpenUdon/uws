package mockruntime

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRequestDigestVectors(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "testdata", "mock-fixtures", "1.0", "request-digest-vectors.json"))
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		Name        string `json:"name"`
		RequestJSON string `json:"requestJson"`
		Canonical   string `json:"canonical"`
		Digest      string `json:"digest"`
	}
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	if len(vectors) < 3 {
		t.Fatalf("portable request vector count = %d, want at least 3", len(vectors))
	}
	for _, vector := range vectors {
		t.Run(vector.Name, func(t *testing.T) {
			request := json.RawMessage(vector.RequestJSON)
			canonical, err := CanonicalizeRequest(request)
			if err != nil {
				t.Fatal(err)
			}
			if got := string(canonical); got != vector.Canonical {
				t.Fatalf("canonical request = %q, want %q", got, vector.Canonical)
			}
			digest, err := RequestDigest(request)
			if err != nil {
				t.Fatal(err)
			}
			if digest != vector.Digest {
				t.Fatalf("request digest = %q, want %q", digest, vector.Digest)
			}
		})
	}
}

func TestCanonicalizeRequestKeysValuesAndRoot(t *testing.T) {
	first := json.RawMessage(`{"query":{"q":"dogs"},"body":{"page":1}}`)
	second := json.RawMessage(` { "body" : {"page":1}, "query" : {"q":"dogs"} } `)
	third := json.RawMessage(`{"body":{"page":2},"query":{"q":"dogs"}}`)
	fourth := json.RawMessage(`{"body":{"page":1},"query":{"q":"dogs"},"path":{"ids":[1,2]}}`)
	arrayForward := json.RawMessage(`{"body":{"ids":[1,2]}}`)
	arrayReverse := json.RawMessage(`{"body":{"ids":[2,1]}}`)
	firstDigest, err := RequestDigest(first)
	if err != nil {
		t.Fatal(err)
	}
	for name, request := range map[string]json.RawMessage{"reordered": second} {
		digest, err := RequestDigest(request)
		if err != nil {
			t.Fatal(err)
		}
		if digest != firstDigest {
			t.Fatalf("%s digest = %s, want %s", name, digest, firstDigest)
		}
	}
	for name, request := range map[string]json.RawMessage{"different value": third, "different array": fourth} {
		digest, err := RequestDigest(request)
		if err != nil {
			t.Fatal(err)
		}
		if digest == firstDigest {
			t.Fatalf("%s retained the same request digest", name)
		}
	}
	forwardDigest, err := RequestDigest(arrayForward)
	if err != nil {
		t.Fatal(err)
	}
	reverseDigest, err := RequestDigest(arrayReverse)
	if err != nil {
		t.Fatal(err)
	}
	if forwardDigest == reverseDigest {
		t.Fatal("array element order did not affect the request digest")
	}
	empty, err := RequestDigest(nil)
	if err != nil {
		t.Fatal(err)
	}
	emptyMap, err := RequestDigest(map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if empty != emptyMap {
		t.Fatalf("absent request digest = %s, empty-object digest = %s", empty, emptyMap)
	}
	for _, request := range []any{[]any{}, `not an object`, json.RawMessage(`{"a":1,"a":2}`)} {
		if _, err := CanonicalizeRequest(request); err == nil {
			t.Fatalf("invalid request root or duplicate key accepted: %#v", request)
		}
	}
}

func TestFixtureSetRoundTripAndStableRepeatedLookup(t *testing.T) {
	digest, err := RequestDigest(map[string]any{"query": map[string]any{"q": "dogs"}})
	if err != nil {
		t.Fatal(err)
	}
	set := &FixtureSet{
		Format: FixtureFormatV1,
		Fixtures: []Fixture{{
			OperationID:   "search",
			RequestDigest: digest,
			Provenance:    FixtureProvenance{Kind: "example"},
			Response:      json.RawMessage(`{"body":{"items":["a"]},"statusCode":200}`),
		}},
	}
	encoded, err := EncodeFixtures(set)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeFixtures(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Format != set.Format || len(decoded.Fixtures) != len(set.Fixtures) || decoded.Fixtures[0].OperationID != set.Fixtures[0].OperationID || decoded.Fixtures[0].RequestDigest != set.Fixtures[0].RequestDigest || !reflect.DeepEqual(decoded.Fixtures[0].Provenance, set.Fixtures[0].Provenance) || !sameJSON(decoded.Fixtures[0].Response, set.Fixtures[0].Response) {
		t.Fatalf("decoded fixture set differs:\ngot  %#v\nwant %#v", decoded, set)
	}
	first, found, err := decoded.LookupResponse("search", digest)
	if err != nil || !found {
		t.Fatalf("first lookup: found=%t err=%v", found, err)
	}
	first[0] = ' '
	second, found, err := decoded.LookupResponse("search", digest)
	if err != nil || !found {
		t.Fatalf("repeat lookup: found=%t err=%v", found, err)
	}
	if !sameJSON(second, set.Fixtures[0].Response) {
		t.Fatalf("repeated lookup response = %s, want unchanged %s", second, set.Fixtures[0].Response)
	}
	if _, found, err := decoded.LookupResponse("other", digest); err != nil || found {
		t.Fatalf("wrong operation lookup: found=%t err=%v", found, err)
	}
	if _, found, err := decoded.LookupResponse("search", strings.Repeat("0", 64)); err != nil || found {
		t.Fatalf("wrong digest lookup: found=%t err=%v", found, err)
	}
}

func TestFixtureSetRejectsInvalidAndAmbiguousDocuments(t *testing.T) {
	digest, err := RequestDigest(nil)
	if err != nil {
		t.Fatal(err)
	}
	fixture := Fixture{
		OperationID:   "get_item",
		RequestDigest: digest,
		Provenance:    FixtureProvenance{Kind: "synthetic"},
		Response:      json.RawMessage(`null`),
	}
	valid := &FixtureSet{Format: FixtureFormatV1, Fixtures: []Fixture{fixture}}
	for name, set := range map[string]*FixtureSet{
		"duplicate pair":  {Format: FixtureFormatV1, Fixtures: []Fixture{fixture, fixture}},
		"unknown version": {Format: "uws.mock-fixtures.2.0", Fixtures: []Fixture{fixture}},
		"empty":           {Format: FixtureFormatV1},
		"bad operation":   {Format: FixtureFormatV1, Fixtures: []Fixture{{OperationID: "bad.id", RequestDigest: digest, Provenance: FixtureProvenance{Kind: "example"}, Response: json.RawMessage(`null`)}}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := EncodeFixtures(set); err == nil {
				t.Fatal("invalid fixture set was encoded")
			}
		})
	}
	encoded, err := EncodeFixtures(valid)
	if err != nil {
		t.Fatal(err)
	}
	for name, malformed := range map[string][]byte{
		"duplicate key": []byte(strings.Replace(string(encoded), `"format":"uws.mock-fixtures.1.0",`, `"format":"uws.mock-fixtures.1.0","format":"uws.mock-fixtures.1.0",`, 1)),
		"extra field":   []byte(strings.Replace(string(encoded), `"format":"uws.mock-fixtures.1.0",`, `"format":"uws.mock-fixtures.1.0","extra":true,`, 1)),
		"trailing json": append(bytes.Clone(encoded), []byte(` {}`)...),
		"response dup":  []byte(`{"format":"uws.mock-fixtures.1.0","fixtures":[{"operationId":"x","requestDigest":"` + digest + `","provenance":{"kind":"example"},"response":{"a":1,"a":2}}]}`),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := DecodeFixtures(malformed); err == nil {
				t.Fatal("malformed fixture set was decoded")
			}
		})
	}
}

func TestEncodeFixturesRejectsOversizedInputBeforeJSONValidation(t *testing.T) {
	set := &FixtureSet{
		Format: FixtureFormatV1,
		Fixtures: []Fixture{{
			OperationID:   "read_item",
			RequestDigest: "sha256:" + strings.Repeat("0", 64),
			Provenance:    FixtureProvenance{Kind: "example"},
			Response:      bytes.Repeat([]byte{'x'}, MaxFixtureSetBytes),
		}},
	}
	if _, err := EncodeFixtures(set); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("oversized input error = %v, want size-limit error", err)
	}
}

func TestRecordedFixtureRequiresExplicitRedaction(t *testing.T) {
	digest, err := RequestDigest(nil)
	if err != nil {
		t.Fatal(err)
	}
	source := json.RawMessage(`{"token":"private-value","value":4}`)
	if _, err := NewRecordedFixture("read_item", digest, source, nil); err == nil {
		t.Fatal("recorded fixture was created without a redactor")
	}
	called := false
	fixture, err := NewRecordedFixture("read_item", digest, source, func(response json.RawMessage) (json.RawMessage, []string, error) {
		called = true
		if !bytes.Equal(response, source) {
			t.Fatalf("redactor input = %s, want %s", response, source)
		}
		return json.RawMessage(`{"value":4}`), []string{"/token"}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !called || fixture.Provenance.Kind != "recorded" || !fixture.Provenance.Redacted || !reflect.DeepEqual(fixture.Provenance.RedactedPointers, []string{"/token"}) {
		t.Fatalf("recorded fixture metadata = %#v; called=%t", fixture.Provenance, called)
	}
	if bytes.Contains(fixture.Response, []byte("private-value")) {
		t.Fatal("recorded fixture retained the original private value")
	}
	encoded, err := EncodeFixtures(&FixtureSet{Format: FixtureFormatV1, Fixtures: []Fixture{fixture}})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(encoded, []byte(`"redactedPointers":[`)) || !bytes.Contains(encoded, []byte(`"/token"`)) {
		t.Fatalf("recorded fixture omitted its redaction pointer list: %s", encoded)
	}
	decoded, err := DecodeFixtures(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(decoded.Fixtures[0].Provenance, fixture.Provenance) {
		t.Fatalf("recorded provenance round trip = %#v, want %#v", decoded.Fixtures[0].Provenance, fixture.Provenance)
	}
	emptyPointerFixture, err := NewRecordedFixture("read_item", digest, json.RawMessage(`{"value":4}`), func(response json.RawMessage) (json.RawMessage, []string, error) {
		return response, []string{}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	emptyPointerBytes, err := EncodeFixtures(&FixtureSet{Format: FixtureFormatV1, Fixtures: []Fixture{emptyPointerFixture}})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(emptyPointerBytes, []byte(`"redactedPointers":[]`)) {
		t.Fatalf("recorded fixture did not encode an empty redaction list: %s", emptyPointerBytes)
	}
	if _, err := NewRecordedFixture("read_item", digest, source, func(json.RawMessage) (json.RawMessage, []string, error) {
		return json.RawMessage(`{}`), []string{"/bad~2pointer"}, nil
	}); err == nil {
		t.Fatal("invalid redaction pointer was accepted")
	}
}

func sameJSON(first, second []byte) bool {
	var compactFirst, compactSecond bytes.Buffer
	if json.Compact(&compactFirst, first) != nil || json.Compact(&compactSecond, second) != nil {
		return false
	}
	return bytes.Equal(compactFirst.Bytes(), compactSecond.Bytes())
}
