package expressions

import (
	"encoding/json"
	"github.com/OpenUdon/uws/uws1"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSharedReferenceConformanceVectors(t *testing.T) {
	file, err := os.Open("../docs/examples/expressions/v1/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var corpus struct {
		Schema string `json:"schema"`
		Cases  []struct {
			Name       string `json:"name"`
			Expression string `json:"expression"`
			Field      string `json:"field"`
			Document   struct {
				UWS        string         `json:"uws"`
				Variables  map[string]any `json:"variables"`
				Components *struct {
					Variables map[string]any `json:"variables"`
				} `json:"components"`
			} `json:"document"`
			State    uws1.ExecutionContext `json:"state"`
			Error    bool                  `json:"expect_error"`
			Expected any                   `json:"expected"`
		} `json:"cases"`
	}
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&corpus); err != nil {
		t.Fatal(err)
	}
	if corpus.Schema != "uws.expression-reference.v1" || len(corpus.Cases) < 18 {
		t.Fatal("missing corpus")
	}
	for _, c := range corpus.Cases {
		t.Run(c.Name, func(t *testing.T) {
			field, ok := map[string]Field{"value": Value, "predicate": Predicate, "wait": Wait, "batchSize": BatchSize}[c.Field]
			if !ok {
				t.Fatal("unknown field")
			}
			document := &uws1.Document{UWS: c.Document.UWS, Variables: c.Document.Variables}
			if c.Document.Components != nil {
				document.Components = &uws1.Components{Variables: c.Document.Components.Variables}
			}
			evaluator, err := NewEvaluator(document)
			if err != nil {
				t.Fatal(err)
			}
			got, err := evaluator.Evaluate(uws1.WithExecutionContext(t.Context(), &c.State), c.Expression, field)
			if c.Error {
				if err == nil {
					t.Fatal("required refusal missing")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			b, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			d := json.NewDecoder(strings.NewReader(string(b)))
			d.UseNumber()
			var normalized any
			if err = d.Decode(&normalized); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(normalized, c.Expected) {
				t.Fatalf("got %v want %v", normalized, c.Expected)
			}
		})
	}
}
