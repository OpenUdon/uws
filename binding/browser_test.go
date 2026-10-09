package binding

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func browserShapeFixture() OperationShape {
	return OperationShape{
		Source:   Source{ID: "browser", Kind: KindBrowser, SHA256: strings.Repeat("a", 64)},
		Selector: Selector{Kind: "id", Value: "read", Key: "read"}, Protocol: "browser", Complete: true,
		Security: Security{Known: true, Alternatives: []SecurityAlternative{{Requirements: []SecurityRequirement{}}}},
		Browser:  &BrowserOperationShape{ProfileVersion: "uws.browser.1.10", CallKind: "action", SelectedSHA256: strings.Repeat("b", 64), Origins: []string{"https://example.test"}, Effects: json.RawMessage(`["read_only"]`), ConfirmationPolicy: json.RawMessage(`{"required":false,"prompt":""}`), CredentialSlots: []CredentialSlotShape{}},
		Inputs:   []Input{{Location: "body", Name: "body", Required: true, Schema: Schema{Known: true, JSON: json.RawMessage(`{"type":"object","properties":{"flag":{"const":false},"count":{"const":0},"text":{"const":""}},"x-keep":{"wide":9223372036854775807,"lexeme":1.00}}`)}}},
		Outputs:  []Output{{Location: "body", Name: "present", Schema: Schema{Known: true, JSON: json.RawMessage(`{"type":"boolean","x-presence":true}`)}}},
	}
}

func browserTableFixture() ShapeTable {
	op := browserShapeFixture()
	return ShapeTable{Version: TableVersion, Sources: []Source{op.Source}, Operations: []OperationShape{op}}
}

func TestBrowserShapeLosslessOptionalAndSnapshot(t *testing.T) {
	table := browserTableFixture()
	table.Operations[0].Browser.RegistrationInputSlots = []RegistrationInputSlotShape{}
	data, err := table.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{`"authentication_required":false`, `"credential_slots":[]`, `"registration_input_slots":[]`, `9223372036854775807`, `1.00`, `"const":false`, `"const":0`, `"const":""`, `"prompt":""`} {
		if !bytes.Contains(data, []byte(fragment)) {
			t.Fatal("lost native value", fragment)
		}
	}
	parsed, err := ParseTable(data)
	if err != nil {
		t.Fatal(err)
	}
	again, err := parsed.Marshal()
	if err != nil || !bytes.Equal(data, again) {
		t.Fatal("unstable roundtrip", err)
	}
	r, err := NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	table.Operations[0].Browser.Origins[0] = "https://changed.test"
	b := Binding{Source: parsed.Sources[0], SelectorKind: "id", SelectorValue: "read"}
	for range 2 {
		got, err := r.Resolve(t.Context(), b)
		if err != nil || got.Shape.Browser.Origins[0] != "https://example.test" {
			t.Fatal("mutable browser snapshot", err)
		}
		got.Shape.Browser.Origins[0] = "https://also-changed.test"
	}
	parsed.Operations[0].Browser.RegistrationInputSlots = nil
	data, _ = parsed.Marshal()
	if bytes.Contains(data, []byte(`"registration_input_slots"`)) {
		t.Fatal("invented optional metadata")
	}
}

func TestBrowserShapeRejectsMalformedNativeMetadata(t *testing.T) {
	cases := map[string]func(*OperationShape){
		"protocol":          func(o *OperationShape) { o.Protocol = "http" },
		"method":            func(o *OperationShape) { o.Method = "GET" },
		"path":              func(o *OperationShape) { o.Path = "/" },
		"server":            func(o *OperationShape) { o.Servers = []string{"https://example.test"} },
		"reference":         func(o *OperationShape) { o.Selector.Kind = "ref" },
		"alias":             func(o *OperationShape) { o.Aliases = []Selector{{Kind: "id", Value: "alias", Key: "read"}} },
		"foreign key":       func(o *OperationShape) { o.Selector.Key = "other" },
		"missing metadata":  func(o *OperationShape) { o.Browser = nil },
		"unknown profile":   func(o *OperationShape) { o.Browser.ProfileVersion = "uws.browser.1.11" },
		"family mismatch":   func(o *OperationShape) { o.Browser.CallKind = "registration" },
		"selected digest":   func(o *OperationShape) { o.Browser.SelectedSHA256 = strings.Repeat("B", 64) },
		"default port":      func(o *OperationShape) { o.Browser.Origins = []string{"https://example.test:443"} },
		"origin path":       func(o *OperationShape) { o.Browser.Origins = []string{"https://example.test/"} },
		"origin wildcard":   func(o *OperationShape) { o.Browser.Origins = []string{"https://*.test"} },
		"unsorted origins":  func(o *OperationShape) { o.Browser.Origins = []string{"https://z.test", "https://a.test"} },
		"duplicate origins": func(o *OperationShape) { o.Browser.Origins = []string{"https://example.test", "https://example.test"} },
		"unknown effect":    func(o *OperationShape) { o.Browser.Effects = json.RawMessage(`["unknown"]`) },
		"read mixed with write": func(o *OperationShape) {
			o.Browser.Effects = json.RawMessage(`["read_only","submits_form"]`)
			o.Browser.ConfirmationPolicy = json.RawMessage(`{"required":true}`)
		},
		"unconfirmed write": func(o *OperationShape) { o.Browser.Effects = json.RawMessage(`["submits_form"]`) },
		"private policy field": func(o *OperationShape) {
			o.Browser.ConfirmationPolicy = json.RawMessage(`{"required":false,"token":"private-canary"}`)
		},
		"action credentials": func(o *OperationShape) {
			o.Browser.CredentialSlots = []CredentialSlotShape{{Name: "password", Kind: "password", Required: true}}
		},
		"query input":   func(o *OperationShape) { o.Inputs[0].Location = "query" },
		"header output": func(o *OperationShape) { o.Outputs[0].Location = "header" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			table := browserTableFixture()
			mutate(&table.Operations[0])
			if err := table.Validate(); err != ErrTable || strings.Contains(err.Error(), "private-canary") {
				t.Fatal("malformed shape accepted", err)
			}
		})
	}
	data, _ := browserTableFixture().Marshal()
	for _, bad := range []string{
		strings.Replace(string(data), `"authentication_required":false`, `"authentication_required":null`, 1),
		strings.Replace(string(data), `"authentication_required":false,`, ``, 1),
		strings.Replace(string(data), `"credential_slots":[]`, `"credential_slots":null`, 1),
		strings.Replace(string(data), `"call_kind":"action"`, `"call_kind":"action","call_kind":"action"`, 1),
		strings.Replace(string(data), `"credential_slots":[]`, `"credential_slots":[],"private":"private-canary"`, 1),
	} {
		if _, err := ParseTable([]byte(bad)); err != ErrTable {
			t.Fatal("open or incomplete wire accepted", err)
		}
	}
}

func TestBrowserAdditionPreservesExistingShapeBytes(t *testing.T) {
	data, err := tableFixture().Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte(`"browser"`)) {
		t.Fatal("non-browser metadata changed")
	}
	var table ShapeTable
	if json.Unmarshal(data, &table) != nil {
		t.Fatal("old shape cannot decode")
	}
	again, err := table.Marshal()
	if err != nil || !bytes.Equal(data, again) {
		t.Fatal("old v1 bytes changed", err)
	}
}
