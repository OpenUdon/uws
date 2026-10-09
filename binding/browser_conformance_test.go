package binding_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/OpenUdon/uws/binding"
	"github.com/OpenUdon/uws/schemas"
)

type browserVector struct {
	Name, Source, SHA256 string
	ProfileVersion       string `json:"profile_version"`
	CallKind             string `json:"call_kind"`
	Selector             string
}

func browserCorpus(t *testing.T) ([]browserVector, string) {
	t.Helper()
	root := "../docs/examples/browser-shapes/v1"
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Version             string
		KinetRevision       string            `json:"kinet_revision"`
		KinetManifestSHA256 string            `json:"kinet_manifest_sha256"`
		FrozenFiles         map[string]string `json:"frozen_files"`
		Vectors             []browserVector
	}
	if json.Unmarshal(data, &manifest) != nil || manifest.Version != "uws.browser-shape-conformance.v1" || manifest.KinetRevision != "ec760d83b6e344e9e9cd7c034f9a92eba99f38b1" || manifest.KinetManifestSHA256 != "0c445a5c90d2c09be561e713c364747f4ab9a3698ea8ab46e7b7b774ccf16bad" || len(manifest.Vectors) != 11 {
		t.Fatal("invalid frozen conformance provenance")
	}
	for path, want := range manifest.FrozenFiles {
		data, err := os.ReadFile(filepath.Join(root, path))
		if err != nil || digest(data) != want {
			t.Fatal("changed frozen fixture", path, err)
		}
	}
	return manifest.Vectors, root
}

func digest(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }

func rawObject(t *testing.T, raw []byte) map[string]json.RawMessage {
	t.Helper()
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil || fields == nil {
		t.Fatal("invalid fixture object")
	}
	return fields
}

func canonical(t *testing.T, raw []byte) []byte {
	t.Helper()
	d := json.NewDecoder(bytes.NewReader(raw))
	d.UseNumber()
	var value any
	if d.Decode(&value) != nil {
		t.Fatal("invalid canonical fixture")
	}
	var out bytes.Buffer
	e := json.NewEncoder(&out)
	e.SetEscapeHTML(false)
	if e.Encode(value) != nil {
		t.Fatal("canonical encoding failed")
	}
	return bytes.TrimSuffix(out.Bytes(), []byte("\n"))
}

// This test-only projection is an independent consumer of native inert files.
// UWS exports no source parser/shape producer and this helper starts no runtime.
func vectorShape(t *testing.T, root string, v browserVector) binding.OperationShape {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, v.Source))
	if err != nil || digest(data) != v.SHA256 {
		t.Fatal("changed native vector", v.Name, err)
	}
	var validate func([]byte) error
	switch v.CallKind {
	case "action":
		validate = schemas.ValidateBrowserSourceProfile
	case "authentication":
		validate = schemas.ValidateBrowserAuthenticationProfile
	case "registration":
		validate = schemas.ValidateBrowserRegistrationProfile
	}
	if validate == nil || validate(data) != nil {
		t.Fatal("native fixture no longer conforms", v.Name, validate(data))
	}
	fields := rawObject(t, data)
	info := rawObject(t, fields["info"])
	inventory := "flows"
	if v.CallKind == "action" {
		inventory = "actions"
	}
	selected := rawObject(t, fields[inventory])[v.Selector]
	leaf := rawObject(t, selected)
	b := &binding.BrowserOperationShape{ProfileVersion: v.ProfileVersion, CallKind: v.CallKind, SelectedSHA256: digest(canonical(t, selected)), CredentialSlots: []binding.CredentialSlotShape{}}
	if v.CallKind == "action" {
		if json.Unmarshal(info["origin"], &b.Origins) != nil {
			var origin string
			if json.Unmarshal(info["origin"], &origin) != nil {
				t.Fatal("invalid action origin")
			}
			b.Origins = []string{origin}
		}
		_ = json.Unmarshal(info["loginStateRequired"], &b.AuthenticationRequired)
		b.Effects = leaf["sideEffects"]
	} else {
		set := map[string]bool{}
		for _, name := range []string{"applicationOrigins", "authenticationOrigins", "registrationOrigins"} {
			var origins []string
			_ = json.Unmarshal(info[name], &origins)
			for _, origin := range origins {
				set[origin] = true
			}
		}
		for origin := range set {
			b.Origins = append(b.Origins, origin)
		}
		sort.Strings(b.Origins)
		b.Effects = leaf["effects"]
		for name, raw := range rawObject(t, fields["credentialSlots"]) {
			var declaration struct{ Kind string }
			if json.Unmarshal(raw, &declaration) != nil {
				t.Fatal("bad credential declaration")
			}
			b.CredentialSlots = append(b.CredentialSlots, binding.CredentialSlotShape{Name: name, Kind: declaration.Kind, Required: true})
		}
		sort.Slice(b.CredentialSlots, func(i, j int) bool { return b.CredentialSlots[i].Name < b.CredentialSlots[j].Name })
	}
	b.ConfirmationPolicy = leaf["confirmationPolicy"]
	if raw, present := fields["inputSlots"]; present {
		for name, schema := range rawObject(t, raw) {
			var declaration struct {
				Type         string
				Required     bool
				RequiredWhen json.RawMessage
			}
			if json.Unmarshal(schema, &declaration) != nil {
				t.Fatal("bad input declaration")
			}
			b.RegistrationInputSlots = append(b.RegistrationInputSlots, binding.RegistrationInputSlotShape{Name: name, Kind: declaration.Type, Required: declaration.Required, Schema: schema, Condition: declaration.RequiredWhen})
		}
		sort.Slice(b.RegistrationInputSlots, func(i, j int) bool { return b.RegistrationInputSlots[i].Name < b.RegistrationInputSlots[j].Name })
	}
	op := binding.OperationShape{Source: binding.Source{ID: v.Name, Kind: binding.KindBrowser, SHA256: digest(data)}, Selector: binding.Selector{Kind: "id", Value: v.Selector, Key: v.Selector}, Protocol: "browser", Complete: true, Browser: b, Security: binding.Security{Known: true, Alternatives: []binding.SecurityAlternative{{Requirements: []binding.SecurityRequirement{}}}}}
	if v.CallKind == "action" {
		if raw, present := leaf["parameters"]; present {
			op.Inputs = []binding.Input{{Location: "body", Name: "body", Schema: binding.Schema{Known: true, JSON: raw}}}
		}
		if raw, present := leaf["outputs"]; present {
			for name, schema := range rawObject(t, raw) {
				op.Outputs = append(op.Outputs, binding.Output{Location: "body", Name: name, Schema: binding.Schema{Known: true, JSON: schema}})
			}
			sort.Slice(op.Outputs, func(i, j int) bool { return op.Outputs[i].Name < op.Outputs[j].Name })
		}
	}
	return op
}

func TestNativeBrowserProfileShapeCorpus(t *testing.T) {
	vectors, root := browserCorpus(t)
	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			op := vectorShape(t, root, v)
			table := binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{op.Source}, Operations: []binding.OperationShape{op}}
			data, err := table.Marshal()
			if err != nil {
				t.Fatal("native browser shape refused", err)
			}
			parsed, err := binding.ParseTable(data)
			if err != nil {
				t.Fatal(err)
			}
			again, err := parsed.Marshal()
			if err != nil || !bytes.Equal(data, again) {
				t.Fatal("native shape not deterministic", err)
			}
			// Complete native schema, extraction metadata and conditions stay intact.
			if !reflect.DeepEqual(op.Inputs, parsed.Operations[0].Inputs) || !reflect.DeepEqual(op.Outputs, parsed.Operations[0].Outputs) {
				for i, in := range op.Inputs {
					if !bytes.Equal(canonical(t, in.Schema.JSON), canonical(t, parsed.Operations[0].Inputs[i].Schema.JSON)) {
						t.Fatal("parameters projection lost fields")
					}
				}
				for i, out := range op.Outputs {
					if !bytes.Equal(canonical(t, out.Schema.JSON), canonical(t, parsed.Operations[0].Outputs[i].Schema.JSON)) {
						t.Fatal("output projection lost extraction metadata")
					}
				}
			}
			for i, slot := range op.Browser.RegistrationInputSlots {
				got := parsed.Operations[0].Browser.RegistrationInputSlots[i]
				if !bytes.Equal(canonical(t, slot.Schema), canonical(t, got.Schema)) {
					t.Fatal("native input schema lost")
				}
				if len(slot.Condition) > 0 && !bytes.Equal(canonical(t, slot.Condition), canonical(t, got.Condition)) {
					t.Fatal("native condition lost")
				}
			}
			r, err := binding.NewResolver(table)
			if err != nil {
				t.Fatal(err)
			}
			b := binding.Binding{Source: op.Source, SelectorKind: "id", SelectorValue: v.Selector}
			resolved, err := r.Resolve(t.Context(), b)
			if err != nil || resolved.Status != binding.Resolved || resolved.Shape.Browser.SelectedSHA256 != op.Browser.SelectedSHA256 {
				t.Fatal("exact native selector failed", err)
			}
			b.Source.SHA256 = strings.Repeat("f", 64)
			resolved, err = r.Resolve(t.Context(), b)
			if err != nil || resolved.Status != binding.Missing {
				t.Fatal("foreign source matched")
			}
			b.Source = op.Source
			b.SelectorKind = "ref"
			resolved, err = r.Resolve(t.Context(), b)
			if err != nil || resolved.Status != binding.Missing {
				t.Fatal("native key expanded into reference")
			}
			op.Complete = false
			for i := range op.Inputs {
				op.Inputs[i].Schema.Known = false
			}
			r, err = binding.NewResolver(binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{op.Source}, Operations: []binding.OperationShape{op}})
			if err != nil {
				t.Fatal(err)
			}
			request := binding.Request{Binding: binding.Binding{Source: op.Source, SelectorKind: "id", SelectorValue: v.Selector}}
			report, err := binding.ValidateBinding(t.Context(), r, request)
			if err != nil || report.Outcome != binding.Indeterminate {
				t.Fatal("incomplete shape became conclusive", report, err)
			}
		})
	}
}

type ambiguousBrowserResolver struct{}

func (ambiguousBrowserResolver) Resolve(context.Context, binding.Binding) (binding.Resolution, error) {
	return binding.Resolution{Status: binding.Ambiguous}, nil
}

func TestAmbiguousBrowserEvidenceNeverResolves(t *testing.T) {
	got, err := binding.ValidateBinding(t.Context(), ambiguousBrowserResolver{}, binding.Request{})
	if err != nil || got.Outcome != binding.Indeterminate || len(got.Diagnostics) != 1 || got.Diagnostics[0].Code != "binding.operation_ambiguous" {
		t.Fatal(got, err)
	}
}

func TestFrozenBrowserPreservationGolden(t *testing.T) {
	_, root := browserCorpus(t)
	data, err := os.ReadFile(filepath.Join(root, "m51-canonical-golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden struct {
		InputJSON    string `json:"input_json"`
		CanonicalHex string `json:"canonical_hex"`
		SHA256       string
	}
	if json.Unmarshal(data, &golden) != nil {
		t.Fatal("bad golden")
	}
	got := canonical(t, []byte(golden.InputJSON))
	if digest(got) != golden.SHA256 || hex.EncodeToString(got) != golden.CanonicalHex {
		t.Fatal("frozen lossless canonical proof changed")
	}
}

func TestFrozenBrowserShapePreservationVectors(t *testing.T) {
	_, root := browserCorpus(t)
	data, err := os.ReadFile(filepath.Join(root, "m51-shape-cases.json"))
	if err != nil {
		t.Fatal(err)
	}
	var corpus struct {
		Vectors []struct {
			Name              string
			ProfileVersion    string `json:"profile_version"`
			RawNumericToken   string `json:"raw_numeric_token"`
			ExpectedAdmission string `json:"expected_admission"`
			Values            json.RawMessage
			Input             string
			Slot              json.RawMessage
			ParametersSchema  json.RawMessage `json:"parameters_schema"`
		}
	}
	if json.Unmarshal(data, &corpus) != nil || len(corpus.Vectors) != 9 {
		t.Fatal("frozen shape vectors changed")
	}
	for _, v := range corpus.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			if v.RawNumericToken != "" {
				op := binding.OperationShape{Source: binding.Source{ID: "fixture", Kind: binding.KindBrowser, SHA256: strings.Repeat("a", 64)}, Selector: binding.Selector{Kind: "id", Value: "read", Key: "read"}, Protocol: "browser", Complete: true, Browser: &binding.BrowserOperationShape{ProfileVersion: v.ProfileVersion, CallKind: "action", SelectedSHA256: strings.Repeat("b", 64), Origins: []string{"https://example.test"}, Effects: json.RawMessage(`["read_only"]`), ConfirmationPolicy: json.RawMessage(`{"required":false}`), CredentialSlots: []binding.CredentialSlotShape{}}, Security: binding.Security{Known: true, Alternatives: []binding.SecurityAlternative{{Requirements: []binding.SecurityRequirement{}}}}, Inputs: []binding.Input{{Location: "body", Name: "n", Schema: binding.Schema{Known: true, JSON: json.RawMessage(`{"type":"integer"}`)}}}}
				r, err := binding.NewResolver(binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{op.Source}, Operations: []binding.OperationShape{op}})
				if err != nil {
					t.Fatal(err)
				}
				report, err := binding.ValidateBinding(t.Context(), r, binding.Request{Binding: binding.Binding{Source: op.Source, SelectorKind: "id", SelectorValue: "read"}, Inputs: []binding.BoundInput{{Location: "body", Name: "n", Value: json.Number(v.RawNumericToken)}}})
				want := binding.Compatible
				if v.ExpectedAdmission == "refused" {
					want = binding.Incompatible
				}
				if err != nil || report.Outcome != want {
					t.Fatal(report, err)
				}
			}
			if len(v.Values) > 0 {
				if !bytes.Equal(canonical(t, v.Values), []byte(`{"count":0,"present":false,"text":""}`)) {
					t.Fatal("frozen presence values changed")
				}
			}
			if v.Input != "" && v.Input != "{{{{literal}}}}" {
				t.Fatal("brace escape fixture changed")
			}
			if len(v.Slot) > 0 {
				slot := rawObject(t, v.Slot)
				condition := rawObject(t, slot["condition"])
				if !bytes.Equal(condition["equals"], []byte(`"fixture-country"`)) {
					t.Fatal("conditional scalar changed")
				}
			}
			if len(v.ParametersSchema) > 0 {
				if !bytes.Contains(canonical(t, v.ParametersSchema), []byte(`"x-fixture-extension":{"flag":false,"zero":0}`)) {
					t.Fatal("complete schema extensions lost")
				}
			}
		})
	}
}

func TestNativeRegistrationMetadataRefusalsAndPartialCondition(t *testing.T) {
	vectors, root := browserCorpus(t)
	v := vectors[len(vectors)-1]
	for name, mutate := range map[string]func(*binding.OperationShape){
		"duplicate": func(op *binding.OperationShape) {
			op.Browser.RegistrationInputSlots = append(op.Browser.RegistrationInputSlots, op.Browser.RegistrationInputSlots[0])
		},
		"wrong kind": func(op *binding.OperationShape) { op.Browser.RegistrationInputSlots[0].Kind = "password" },
		"credential collision": func(op *binding.OperationShape) {
			op.Browser.CredentialSlots = append(op.Browser.CredentialSlots, binding.CredentialSlotShape{Name: op.Browser.RegistrationInputSlots[0].Name, Kind: "password", Required: true})
		},
		"undeclared condition": func(op *binding.OperationShape) {
			for i := range op.Browser.RegistrationInputSlots {
				if len(op.Browser.RegistrationInputSlots[i].Condition) > 0 {
					op.Browser.RegistrationInputSlots[i].Condition = json.RawMessage(`{"slot":"missing","equals":"business"}`)
					return
				}
			}
		},
		"condition private field": func(op *binding.OperationShape) {
			op.Browser.RegistrationInputSlots[0].Condition = json.RawMessage(`{"slot":"missing","equals":"business","token":"private-canary"}`)
		},
	} {
		t.Run(name, func(t *testing.T) {
			op := vectorShape(t, root, v)
			mutate(&op)
			table := binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{op.Source}, Operations: []binding.OperationShape{op}}
			if table.Validate() != binding.ErrTable {
				t.Fatal("malformed native metadata accepted")
			}
		})
	}
	op := vectorShape(t, root, v)
	op.Complete = false
	// A partial inventory can retain an unresolved native condition and its
	// full schema; it cannot become conclusive or grant a private input read.
	for i := range op.Browser.RegistrationInputSlots {
		if len(op.Browser.RegistrationInputSlots[i].Condition) > 0 {
			slot := &op.Browser.RegistrationInputSlots[i]
			slot.Condition = json.RawMessage(`{"slot":"missing","equals":"business"}`)
			schema := rawObject(t, slot.Schema)
			schema["requiredWhen"] = slot.Condition
			slot.Schema, _ = json.Marshal(schema)
		}
	}
	r, err := binding.NewResolver(binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{op.Source}, Operations: []binding.OperationShape{op}})
	if err != nil {
		t.Fatal(err)
	}
	report, err := binding.ValidateBinding(t.Context(), r, binding.Request{Binding: binding.Binding{Source: op.Source, SelectorKind: "id", SelectorValue: v.Selector}})
	if err != nil || report.Outcome != binding.Indeterminate {
		t.Fatal(report, err)
	}
}

func TestCompleteNativeRegistrationSlotDeclarations(t *testing.T) {
	vectors, root := browserCorpus(t)
	for _, v := range vectors {
		if v.ProfileVersion != "uws.browser-registration.1.1" && v.ProfileVersion != "uws.browser-registration.1.2" {
			continue
		}
		t.Run(v.Name, func(t *testing.T) {
			for name, raw := range map[string]string{
				"missing label":                             `{"type":"string","required":false}`,
				"missing requiredness":                      `{"type":"string","label":"Contact"}`,
				"minimal preservation projection":           `{"type":"string"}`,
				"empty label":                               `{"type":"string","label":"","required":false}`,
				"unknown field":                             `{"type":"string","label":"Contact","required":false,"private":"private-canary"}`,
				"schema extension outside native inventory": `{"type":"string","label":"Contact","required":false,"x-extension":false}`,
				"both requiredness branches":                `{"type":"string","label":"Contact","required":false,"requiredWhen":{"slot":"country","equals":"fixture-country"}}`,
				"required flag contradiction":               `{"type":"string","label":"Contact","required":true}`,
				"condition without metadata":                `{"type":"string","label":"Contact","requiredWhen":{"slot":"country","equals":"fixture-country"}}`,
			} {
				t.Run(name, func(t *testing.T) {
					op := vectorShape(t, root, v)
					op.Browser.RegistrationInputSlots = []binding.RegistrationInputSlotShape{{Name: "contact", Kind: "string", Required: false, Schema: json.RawMessage(raw)}}
					table := binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{op.Source}, Operations: []binding.OperationShape{op}}
					if err := table.Validate(); err != binding.ErrTable || strings.Contains(err.Error(), "private-canary") {
						t.Fatal("incomplete or foreign native declaration accepted", err)
					}
				})
			}
			op := vectorShape(t, root, v)
			op.Browser.RegistrationInputSlots = []binding.RegistrationInputSlotShape{{Name: "contact", Kind: "string", Required: false, Schema: json.RawMessage(`{"type":"string","label":"Contact","required":false,"maxLength":0}`)}}
			table := binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{op.Source}, Operations: []binding.OperationShape{op}}
			data, err := table.Marshal()
			if err != nil {
				t.Fatal("complete native declaration refused", err)
			}
			parsed, err := binding.ParseTable(data)
			if err != nil || !bytes.Equal(canonical(t, parsed.Operations[0].Browser.RegistrationInputSlots[0].Schema), canonical(t, op.Browser.RegistrationInputSlots[0].Schema)) {
				t.Fatal("complete raw declaration changed", err)
			}
			// The frozen minimal slot is a preservation vector. Its incomplete
			// declaration can be carried only with explicitly partial evidence.
			op.Complete = false
			op.Browser.RegistrationInputSlots[0].Schema = json.RawMessage(`{"type":"string"}`)
			table.Operations[0] = op
			r, err := binding.NewResolver(table)
			if err != nil {
				t.Fatal("partial projection refused", err)
			}
			report, err := binding.ValidateBinding(t.Context(), r, binding.Request{Binding: binding.Binding{Source: op.Source, SelectorKind: "id", SelectorValue: op.Selector.Value}})
			if err != nil || report.Outcome != binding.Indeterminate || len(report.Diagnostics) == 0 || report.Diagnostics[0].Code != "binding.operation_incomplete" {
				t.Fatal("partial declaration became conclusive", report, err)
			}
		})
	}
}
