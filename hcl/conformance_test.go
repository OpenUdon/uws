package hcl

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assertVerifiedFixture(t *testing.T, source Source) View {
	t.Helper()
	view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(context.Background(), source, view, Options{Revision: fixtureRevision}); err != nil {
		t.Fatal(err)
	}
	data, err := Import(context.Background(), view.HCL)
	if err != nil {
		t.Fatal(err)
	}
	decode := func(bytes []byte) any {
		d := json.NewDecoder(strings.NewReader(string(bytes)))
		d.UseNumber()
		var value any
		if err := d.Decode(&value); err != nil {
			t.Fatal(err)
		}
		return value
	}
	var want any
	if source.Format == JSON {
		want = decode(source.Bytes)
	} else {
		want, err = sourceValue(context.Background(), source)
		if err != nil {
			t.Fatal(err)
		}
	}
	got := decode(data)
	left, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	right, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(left, right) {
		t.Fatal("independent imported projection differs, including numeric lexemes")
	}
	return view
}

func TestCodecExistingTypedCorpus(t *testing.T) {
	manifestBytes, err := os.ReadFile("testdata/conformance/v1/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Version string `json:"version"`
		Sources []struct {
			Path    string `json:"path"`
			Fixture string `json:"fixture"`
			Format  Format `json:"format"`
			Bytes   int    `json:"bytes"`
			SHA256  string `json:"sha256"`
		} `json:"sources"`
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Version != "uws.hcl-conformance.v1" || len(manifest.Sources) != 6 {
		t.Fatal("unexpected pinned conformance membership")
	}
	for _, entry := range manifest.Sources {
		name := entry.Path
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", "conformance", "v1", entry.Fixture))
			if err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(data)
			if len(data) != entry.Bytes || hex.EncodeToString(sum[:]) != entry.SHA256 {
				t.Fatal("pinned source fixture changed")
			}
			assertVerifiedFixture(t, Source{Format: entry.Format, Bytes: data})
		})
	}
}

func TestCodecExactNumbersStringsAndKeyCorpus(t *testing.T) {
	for _, numbers := range []string{
		`[0,-0,0.00,-0.000,1.2300,1e+03,1E-0003,-1.23456789012345678901234567890123456789e+17]`,
		`[9007199254740993,900719925474099312345678901234567890,0.000000000000000000000001234567890]`,
	} {
		source := Source{Format: JSON, Bytes: []byte(`{"uws":"1.13.0","info":{"title":"numeric","version":"1"},"operations":[],"variables":{"numbers":` + numbers + `}}`)}
		view := assertVerifiedFixture(t, source)
		for _, token := range strings.Split(strings.Trim(numbers, "[]"), ",") {
			if !bytes.Contains(view.HCL, []byte(token)) {
				t.Fatal("numeric token changed", token)
			}
		}
	}
	values := map[string]any{
		"$ref": "a", "_ref": "b", "$custom": "c", "__dollar__custom": "d", "__uws_literal___ref": "e",
		"space key": "space", "中文/🙂": "unicode", "empty": "", "quote": "\"\\\n\t\r",
		"literal": "${function(/private)} %{if x} $${already} %%{already}",
		"nested":  []any{map[string]any{"$defs": map[string]any{"_schema": nil}}, []any{}, map[string]any{}},
	}
	data, err := json.Marshal(map[string]any{"uws": "1.13.0", "info": map[string]any{"title": "strings", "version": "1"}, "operations": []any{}, "variables": values})
	if err != nil {
		t.Fatal(err)
	}
	assertVerifiedFixture(t, Source{Format: JSON, Bytes: data})
}

func TestCodecMalformedInputCorpus(t *testing.T) {
	for name, source := range map[string]Source{
		"duplicate":         {Format: JSON, Bytes: []byte(`{"uws":"1.13.0","uws":"1.12.0"}`)},
		"escaped-duplicate": {Format: JSON, Bytes: []byte(`{"uws":"1.13.0","\u0075ws":"1.12.0"}`)},
		"trailing":          {Format: JSON, Bytes: []byte(`{"uws":"1.13.0"} {}`)},
		"surrogate":         {Format: JSON, Bytes: []byte(`{"variables":{"bad":"\ud800"}}`)},
		"unknown-typed":     {Format: JSON, Bytes: []byte(`{"uws":"1.13.0","unsupported":true}`)},
		"yaml-alias":        {Format: YAML, Bytes: []byte("uws: '1.13.0'\nvariables: &v {a: 1}\ncomponents: {variables: *v}\n")},
		"yaml-merge":        {Format: YAML, Bytes: []byte("variables: {<<: {a: 1}}\n")},
		"yaml-key":          {Format: YAML, Bytes: []byte("variables: {true: 'value'}\n")},
		"yaml-number":       {Format: YAML, Bytes: []byte("variables: {value: .nan}\n")},
		"yaml-trailing":     {Format: YAML, Bytes: []byte("uws: '1.13.0'\n---\nuws: '1.12.0'\n")},
	} {
		t.Run(name, func(t *testing.T) {
			view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
			if !errors.Is(err, ErrCodec) || len(view.HCL) > 0 {
				t.Fatal("malformed/unsupported source exposed a view", err)
			}
		})
	}
	for name, data := range map[string]string{
		"duplicate-attribute": `uws = "1.13.0"
uws = "1.12.0"`,
		"duplicate-object":      `variables = { key = "one", "key" = "two" }`,
		"decoded-key-collision": `variables = { _ref = "one", __dollar__ref = "two" }`,
		"mixed-empty-block": `operation = []
operation "op" { effect = "read" }`,
		"label-attribute": `operation "first" { operationId = "second" }`,
		"unknown-block":   `unsupported { }`,
		"extension-typed": `extensions { uws = "1.13.0" }`,
		"template":        `variables = { value = "%{if secret}hidden%{endif}" }`,
		"for-expression":  `variables = { value = [for v in secret : v] }`,
		"bad-number":      `variables = { value = - 1 }`,
	} {
		t.Run(name, func(t *testing.T) {
			data, err := Import(context.Background(), []byte(data))
			if !errors.Is(err, ErrCodec) || len(data) > 0 {
				t.Fatal("ambiguous/executable HCL imported", err)
			}
		})
	}
}

func TestCodecValueFreeFailuresAndBounds(t *testing.T) {
	canary := "private-known-canary"
	for _, source := range []Source{
		{Format: JSON, Bytes: []byte(`{"unknown":"` + canary + `"}`)},
		{Format: JSON, Bytes: []byte(strings.Repeat(" ", MaxSourceBytes+1))},
		{Format: JSON, Bytes: []byte(`{"variables":{"value":` + strings.Repeat("[", maxDepth+1) + `1` + strings.Repeat("]", maxDepth+1) + `}}`)},
	} {
		view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
		if err == nil || len(view.HCL) > 0 || strings.Contains(err.Error(), canary) {
			t.Fatal("unsafe failure/oversize result")
		}
	}
	if data, err := Import(context.Background(), []byte(strings.Repeat(" ", MaxViewBytes+1))); err == nil || len(data) > 0 {
		t.Fatal("oversized HCL imported")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if data, err := Import(ctx, []byte(`uws = "1.13.0"`)); !errors.Is(err, context.Canceled) || len(data) > 0 {
		t.Fatal("cancelled import exposed values", err)
	}
}
