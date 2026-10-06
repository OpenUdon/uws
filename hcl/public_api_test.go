package hcl_test

import (
	"context"
	"strings"
	"testing"

	viewhcl "github.com/OpenUdon/uws/hcl"
)

func TestPublicLosslessViewAPI(t *testing.T) {
	source := viewhcl.Source{Format: viewhcl.JSON, Bytes: []byte(`{"uws":"1.13.0","info":{"title":"public API","version":"1"},"operations":[],"variables":{"number":900719925474099312345,"lexeme":1.2300,"text":"1.2300"}}`)}
	options := viewhcl.Options{Revision: strings.Repeat("a", 40)}
	view, err := viewhcl.Render(context.Background(), source, options)
	if err != nil {
		t.Fatal(err)
	}
	if err := viewhcl.Verify(context.Background(), source, view, options); err != nil {
		t.Fatal(err)
	}
	if view.Provenance.Version != viewhcl.ContractVersion || view.Provenance.Format != viewhcl.JSON {
		t.Fatal("public provenance contract")
	}
}
