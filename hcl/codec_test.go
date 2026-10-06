package hcl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

var fixtureRevision = strings.Repeat("a", 40)

func TestRenderVerifyExactTypedView(t *testing.T) {
	source := Source{Format: JSON, Bytes: []byte(`{"uws":"1.13.0","info":{"title":"Fixture ${literal}","version":"1","x-note":{"$ref":"literal"}},"operations":[{"operationId":"read","effect":"read","request":{"body":{"big":900719925474099312345,"precise":1.2300,"exp":1e+03,"minus":-0,"$custom":"value","_ref":"literal","__dollar__foo":"escaped"}},"x-uws-operation-profile":"fixture"}],"workflows":[{"workflowId":"main","type":"sequence","steps":[{"stepId":"one","operationRef":"read"}],"inputs":{"type":"object","properties":{"foo":{"$ref":"#/x","x-note":"safe"}}}}],"variables":{},"x-private":null}`)}
	view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range []string{`operation "read"`, `workflow "main"`, `step "one"`, `900719925474099312345`, `1.2300`, `1e+03`, `-0`, `extensions {`, `__uws_literal___ref`, `__dollar__custom`} {
		if !bytes.Contains(view.HCL, []byte(token)) {
			t.Fatalf("missing typed/lexical mapping %s in %s", token, view.HCL)
		}
	}
	if err := Verify(context.Background(), source, view, Options{Revision: fixtureRevision}); err != nil {
		t.Fatal(err)
	}
	imported, err := Import(context.Background(), view.HCL)
	if err != nil {
		t.Fatal(err)
	}
	var value map[string]any
	d := json.NewDecoder(bytes.NewReader(imported))
	d.UseNumber()
	if err := d.Decode(&value); err != nil {
		t.Fatal(err)
	}
	if value["x-private"] != nil || value["variables"] == nil {
		t.Fatal("lost null/empty shape")
	}
	other, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
	if err != nil || !bytes.Equal(view.HCL, other.HCL) {
		t.Fatal("nondeterministic view", err)
	}
}

func TestRenderYAMLAndEmptyPresence(t *testing.T) {
	for _, source := range []Source{
		{Format: YAML, Bytes: []byte("uws: '1.13.0'\ninfo: {title: Fixture, version: '1'}\noperations: []\nvariables: {value: 1.2300, big: 900719925474099312345, exponent: 1e+03}\ncomponents: null\n")},
		{Format: JSON, Bytes: []byte(`{"uws":"1.13.0","info":{},"operations":[],"workflows":null,"variables":{"empty":[],"map":{},"nothing":null}}`)},
	} {
		view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
		if err != nil {
			t.Fatal(err)
		}
		if err := Verify(context.Background(), source, view, Options{Revision: fixtureRevision}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCodecKeepsLiteralBlockLabels(t *testing.T) {
	source := Source{Format: JSON, Bytes: []byte(`{"uws":"1.13.0","info":{"title":"labels","version":"1"},"operations":[{"operationId":"name ${literal} %{directive}","x-uws-operation-profile":"fixture"}]}`)}
	view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(context.Background(), source, view, Options{Revision: fixtureRevision}); err != nil {
		t.Fatal(err)
	}
}

func TestCodecRefusesTamperingAndExecution(t *testing.T) {
	source := Source{Format: JSON, Bytes: []byte(`{"uws":"1.13.0","info":{"title":"Fixture","version":"1"},"operations":[],"variables":{"value":1.2300}}`)}
	view, err := Render(context.Background(), source, Options{Revision: fixtureRevision})
	if err != nil {
		t.Fatal(err)
	}
	copy := view
	copy.HCL = bytes.ReplaceAll(view.HCL, []byte("1.2300"), []byte("1.23"))
	copy.Provenance.ViewSHA256 = digest(copy.HCL)
	if err := Verify(context.Background(), source, copy, Options{Revision: fixtureRevision}); !errors.Is(err, ErrCodec) {
		t.Fatal("changed numeric lexeme accepted", err)
	}
	changed := source
	changed.Bytes = append(append([]byte(nil), source.Bytes...), ' ')
	if err := Verify(context.Background(), changed, view, Options{Revision: fixtureRevision}); !errors.Is(err, ErrCodec) {
		t.Fatal("stale raw identity accepted", err)
	}
	if err := Verify(context.Background(), source, view, Options{Revision: strings.Repeat("b", 40)}); !errors.Is(err, ErrCodec) {
		t.Fatal("stale codec accepted", err)
	}
	for _, data := range []string{`variables = { secret = file("/private") }`, `variables = { secret = env.SECRET }`, `variables = { secret = "${var.secret}" }`, `variables = { secret = 1+2 }`, `variables = { x = 1, x = 2 }`} {
		if _, err := Import(context.Background(), []byte(data)); !errors.Is(err, ErrCodec) {
			t.Fatal("unsafe/ambiguous HCL accepted", data, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if view, err := Render(ctx, source, Options{Revision: fixtureRevision}); !errors.Is(err, context.Canceled) || len(view.HCL) > 0 {
		t.Fatal("canceled rendering returned view", err)
	}
}
