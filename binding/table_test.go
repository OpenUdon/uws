package binding

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func sourceFixture() Source {
	return Source{ID: "weather", Kind: "openapi", SHA256: strings.Repeat("a", 64), URL: "weather.json"}
}
func shapeFixture() OperationShape {
	return OperationShape{Source: sourceFixture(), Selector: Selector{Kind: "id", Value: "read", Key: "read"}, Aliases: []Selector{{Kind: "ref", Value: "#/paths/~1weather/get", Key: "read"}}, Protocol: "http", Method: "GET", Path: "/weather", Complete: true, Inputs: []Input{{Location: "query", Name: "city", Required: true, Schema: Schema{Known: true, JSON: json.RawMessage(`{"type":"string"}`)}}}, Security: Security{Known: true, Alternatives: []SecurityAlternative{{Requirements: []SecurityRequirement{{Scheme: "key", Type: "apiKey", Location: "query", Name: "appid"}, {Scheme: "oauth", Type: "oauth2", Scopes: []string{"weather.read"}}}}, {Requirements: []SecurityRequirement{{Scheme: "bearer", Type: "http"}}}}}}
}
func tableFixture() ShapeTable {
	return ShapeTable{Version: TableVersion, Sources: []Source{sourceFixture()}, Operations: []OperationShape{shapeFixture()}}
}
func TestTableRoundTripAndResolverSnapshot(t *testing.T) {
	table := tableFixture()
	data, err := table.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseTable(data)
	if err != nil || len(parsed.Operations[0].Security.Alternatives) != 2 || len(parsed.Operations[0].Security.Alternatives[0].Requirements) != 2 {
		t.Fatalf("security alternatives lost: %v", err)
	}
	resolver, err := NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	table.Operations[0].Inputs[0].Name = "mutated"
	for _, selector := range []Selector{shapeFixture().Selector, shapeFixture().Aliases[0]} {
		result, err := resolver.Resolve(t.Context(), Binding{Source: sourceFixture(), SelectorKind: selector.Kind, SelectorValue: selector.Value})
		if err != nil || result.Status != Resolved || result.Shape.Inputs[0].Name != "city" {
			t.Fatalf("alias/snapshot: %+v %v", result, err)
		}
		result.Shape.Inputs[0].Name = "also-mutated"
	}
	b := Binding{Source: sourceFixture(), SelectorKind: "id", SelectorValue: "read"}
	b.Source.SHA256 = strings.Repeat("b", 64)
	result, err := resolver.Resolve(t.Context(), b)
	if err != nil || result.Status != Missing {
		t.Fatal("forged source matched")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := resolver.Resolve(ctx, b); err != context.Canceled {
		t.Fatal(err)
	}
}
func TestAmbiguousSelectorsAndUnknownSchemasRemainExplicit(t *testing.T) {
	table := tableFixture()
	other := shapeFixture()
	other.Selector.Key = "other"
	other.Aliases = nil
	table.Operations = append(table.Operations, other)
	resolver, err := NewResolver(table)
	if err != nil {
		t.Fatal(err)
	}
	result, err := resolver.Resolve(t.Context(), Binding{Source: sourceFixture(), SelectorKind: "id", SelectorValue: "read"})
	if err != nil || result.Status != Ambiguous || result.Shape != nil {
		t.Fatal(result, err)
	}
	unknown := tableFixture()
	unknown.Operations[0].Inputs[0].Schema = Schema{Known: false, JSON: json.RawMessage(`{"type":"string"}`)}
	data, err := unknown.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	p, err := ParseTable(data)
	if err != nil || p.Operations[0].Inputs[0].Schema.Known {
		t.Fatal("partial schema became known")
	}
}
func TestClosedTableAndProtocolRefusals(t *testing.T) {
	good, err := tableFixture().Marshal()
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range [][]byte{append(append([]byte(nil), good...), good...), []byte(`{"version":"uws.shape-table.v1","version":"uws.shape-table.v1","sources":[],"operations":[]}`), []byte(`{"version":"uws.shape-table.v1","sources":[],"operations":[],"token":"private-canary"}`)} {
		if _, err := ParseTable(bad); err == nil || strings.Contains(err.Error(), "private-canary") {
			t.Fatal("unsafe contract accepted", err)
		}
	}
	table := tableFixture()
	table.Operations[0].Protocol = "json-rpc"
	if table.Validate() == nil {
		t.Fatal("fabricated HTTP metadata on RPC accepted")
	}
	table.Operations[0].Method = ""
	table.Operations[0].Path = ""
	if table.Validate() != nil {
		t.Fatal("honest non-HTTP metadata refused")
	}
	table.Operations[0].Security = Security{Known: true}
	if table.Validate() == nil {
		t.Fatal("missing security evidence became anonymous")
	}
	table.Operations[0].Security = Security{Known: false}
	if table.Validate() != nil {
		t.Fatal("unknown security refused structurally")
	}
	table.Operations[0].Source.SHA256 = strings.Repeat("b", 64)
	if table.Validate() == nil {
		t.Fatal("source identity mismatch accepted")
	}
}
func TestMarshalDeterminismWithoutMutation(t *testing.T) {
	a := tableFixture()
	op := shapeFixture()
	op.Selector.Key = "a"
	op.Selector.Value = "a"
	op.Aliases = nil
	a.Operations = append(a.Operations, op)
	b := a
	b.Operations = []OperationShape{a.Operations[1], a.Operations[0]}
	x, e1 := a.Marshal()
	y, e2 := b.Marshal()
	if e1 != nil || e2 != nil || string(x) != string(y) || a.Operations[0].Selector.Key != "read" {
		t.Fatal("unstable/mutating serialization")
	}
}
