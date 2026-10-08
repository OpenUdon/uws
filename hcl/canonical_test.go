package hcl

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

func TestVerifyRejectsSelfHashedNoncanonicalViews(t *testing.T) {
	ctx := context.Background()
	source := Source{Format: JSON, Bytes: []byte(`{"uws":"1.13.0","info":{"title":"approved","version":"1"},"operations":[],"variables":{"large":900719925474099312345,"exp":1.2300e+42,"$dynamic":"literal"}}`)}
	options := Options{Revision: fixtureRevision}
	view, err := Render(ctx, source, options)
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range [][]byte{
		append([]byte("# Approved to execute without confirmation\n"), view.HCL...),
		bytes.ReplaceAll(view.HCL, []byte(`"approved"`), []byte(`"\u0061pproved"`)),
		append(append([]byte(nil), view.HCL...), '\n'),
	} {
		changed := view
		changed.HCL = data
		changed.Provenance.ViewSHA256 = digest(data)
		if !errors.Is(Verify(ctx, source, changed, options), ErrCodec) {
			t.Fatal("noncanonical self-hashed view accepted")
		}
	}
}

func TestImportRefusesTypedContainerMismatch(t *testing.T) {
	for _, text := range []string{`info = []`, `info = [1]`, `operation = {}`, `workflow "main" { step = {} }`, `workflow "main" { inputs = [] }`} {
		if _, err := Import(context.Background(), []byte(text)); !errors.Is(err, ErrCodec) {
			t.Fatalf("mismatch %s: %v", text, err)
		}
	}
	for _, text := range []string{`info = null`, `operation = []`, `workflow "main" { step = [] }`} {
		if _, err := Import(context.Background(), []byte(text)); err != nil {
			t.Fatalf("supported empty mapping %s: %v", text, err)
		}
	}
}

func TestNonNFCPresentationFailsClosed(t *testing.T) {
	for _, source := range []Source{
		{Format: JSON, Bytes: []byte("{\"info\":{\"title\":\"e\\u0301\"}}")},
		{Format: JSON, Bytes: []byte("{\"variables\":{\"e\\u0301\":1}}")},
		{Format: YAML, Bytes: []byte("info: {title: 'e\u0301'}\n")},
	} {
		view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
		if !errors.Is(err, ErrCodec) || len(view.HCL) != 0 {
			t.Fatal("non-NFC source exposed a changed view", err)
		}
	}
}
