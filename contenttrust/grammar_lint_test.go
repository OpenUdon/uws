package contenttrust

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type legacyGrammarAllowance struct {
	value  string
	reason string
}

func TestExpressionGrammarDocsAndVersionedFixtures(t *testing.T) {
	docPage, err := os.ReadFile("../docs/03-Runtime-Expression-Grammar.md")
	if err != nil {
		t.Fatal(err)
	}
	blocks := yamlCodeBlocks(t, docPage)
	if len(blocks) == 0 {
		t.Fatal("expression grammar page has no YAML examples")
	}
	docAllowances := map[string]legacyGrammarAllowance{
		"docs/03/yaml[11]/examples/1/when": {
			value:  "$response.statusCode == 200 && $response.body.count > 0",
			reason: "explicitly labeled implementation-extension example in docs/03 §Implementation Extensions",
		},
	}
	seenDocAllowances := make(map[string]bool)
	for i, block := range blocks {
		var example any
		if err := yaml.Unmarshal(block, &example); err != nil {
			t.Fatalf("docs/03 YAML example %d: %v", i, err)
		}
		lintExpressions(t, example, "1.11.0", false, fmt.Sprintf("docs/03/yaml[%d]", i), docAllowances, seenDocAllowances, false)
	}
	for path := range docAllowances {
		if !seenDocAllowances[path] {
			t.Errorf("stale docs/03 expression allowance at %s", path)
		}
	}
	for i, block := range fencedCodeBlocks(t, docPage, "hcl") {
		lintHCLExpressionStrings(t, block, fmt.Sprintf("docs/03/hcl[%d]", i))
	}

	fixture, err := os.ReadFile("../testdata/grammar/1.11.0.json")
	if err != nil {
		t.Fatal(err)
	}
	version, value := decodeVersionedJSON(t, fixture, "testdata/grammar/1.11.0.json")
	if version != "1.11.0" {
		t.Fatalf("grammar fixture declares UWS %q, want 1.11.0", version)
	}
	lintExpressions(t, value, version, false, "testdata/grammar/1.11.0.json", nil, nil, false)

	legacyAllowances := map[string]legacyGrammarAllowance{
		"/operations/0/outputs/petList": {
			value:  "$response.body.items",
			reason: "legacy UWS 1.0 response-body dot-walk; this sample is not a grammar-conformance fixture",
		},
		"/operations/1/when": {
			value:  "length(list_pets.received_body.items) < 100",
			reason: "implementation-defined UWS 1.0 expression retained by the round-trip sample",
		},
		"/workflows/0/steps/0/outputs/isValid": {
			value:  "$response.body.valid",
			reason: "legacy UWS 1.0 response-body dot-walk; this sample is not a grammar-conformance fixture",
		},
		"/results/0/value": {
			value:  "$steps.merge_validation.outputs",
			reason: "legacy whole-step-output reference retained by the UWS 1.0 round-trip sample",
		},
	}
	sample, err := os.ReadFile("../testdata/sample.uws.json")
	if err != nil {
		t.Fatal(err)
	}
	sampleVersion, sampleValue := decodeVersionedJSON(t, sample, "testdata/sample.uws.json")
	if sampleVersion != "1.0.0" {
		t.Fatalf("sample fixture declares UWS %q, want 1.0.0", sampleVersion)
	}
	seenAllowances := make(map[string]bool)
	lintExpressions(t, sampleValue, sampleVersion, false, "", legacyAllowances, seenAllowances, false)
	for path, allowance := range legacyAllowances {
		if !seenAllowances[path] {
			t.Errorf("stale legacy expression allowance at %s (%s)", path, allowance.reason)
		}
		if !bytes.Contains(docPage, []byte(path)) {
			t.Errorf("legacy expression allowance is not documented in docs/03: %s", path)
		}
	}

	// This exact generated file is a test-runtime expression fixture rather
	// than a portable grammar/conformance fixture. Its file-level exception is
	// documented next to the fixture and kept separate from the sample's
	// pointer-level legacy exceptions above.
	const bigFixturePath = "../testdata/big/big.json"
	bigData, err := os.ReadFile(bigFixturePath)
	if err != nil {
		t.Fatal(err)
	}
	bigVersion, _ := decodeVersionedJSON(t, bigData, bigFixturePath)
	if bigVersion != "1.1.0" {
		t.Fatalf("legacy big fixture version changed unexpectedly to %q", bigVersion)
	}
	bigDocs, err := os.ReadFile("../docs/big-fixture.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"`testdata/big/big.json`", "not a portable", "not evidence that those", "implementation-defined expressions are portable"} {
		if !bytes.Contains(bigDocs, []byte(required)) {
			t.Errorf("big fixture file-level grammar disposition does not document %q", required)
		}
	}
}

func decodeVersionedJSON(t *testing.T, data []byte, path string) (string, any) {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	root, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s root is %T, want object", path, value)
	}
	version, ok := root["uws"].(string)
	if !ok {
		t.Fatalf("%s has no string uws version", path)
	}
	return version, value
}

func yamlCodeBlocks(t *testing.T, markdown []byte) [][]byte {
	t.Helper()
	return fencedCodeBlocks(t, markdown, "yaml")
}

func fencedCodeBlocks(t *testing.T, markdown []byte, wantedLanguage string) [][]byte {
	t.Helper()
	var blocks [][]byte
	var current bytes.Buffer
	var scanner = bufio.NewScanner(bytes.NewReader(markdown))
	inFence, isYAML := false, false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if !inFence {
				language := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "```"))
				isYAML = language == wantedLanguage || wantedLanguage == "yaml" && language == "yml"
				inFence = true
				current.Reset()
				continue
			}
			if isYAML {
				blocks = append(blocks, append([]byte(nil), current.Bytes()...))
			}
			inFence, isYAML = false, false
			current.Reset()
			continue
		}
		if inFence && isYAML {
			current.WriteString(line)
			current.WriteByte('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	return blocks
}

func lintHCLExpressionStrings(t *testing.T, block []byte, blockPath string) {
	t.Helper()
	quoted := regexp.MustCompile(`"(?:[^"\\]|\\.)*"`)
	for lineNumber, line := range strings.Split(string(block), "\n") {
		for _, raw := range quoted.FindAllString(line, -1) {
			value, err := strconv.Unquote(raw)
			if err != nil || !strings.Contains(value, "$") {
				continue
			}
			if _, ok := parseExpressionForVersion(value, "1.11.0", false); !ok {
				t.Errorf("expression at %s:%d is not valid under UWS 1.11 grammar: %q", blockPath, lineNumber+1, value)
			}
		}
	}
}

func lintExpressions(t *testing.T, value any, version string, loopContext bool, path string, allowances map[string]legacyGrammarAllowance, seen map[string]bool, awaitWait bool) {
	t.Helper()
	switch typed := value.(type) {
	case map[string]any:
		isAwait := typed["type"] == "await"
		childLoop := loopContext
		if typed["type"] == "loop" {
			childLoop = true
		} else if forEach, ok := typed["forEach"].(string); ok && forEach != "" {
			childLoop = false
		}
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fieldLoop := childLoop
			fieldAwaitWait := isAwait && key == "wait"
			switch key {
			case "when", "forEach", "wait", "items", "batchSize", "condition", "context":
				fieldLoop = loopContext
			}
			lintExpressions(t, typed[key], version, fieldLoop, jsonPointer(path, key), allowances, seen, fieldAwaitWait)
		}
	case []any:
		for index, item := range typed {
			lintExpressions(t, item, version, loopContext, jsonPointer(path, fmt.Sprint(index)), allowances, seen, awaitWait)
		}
	case string:
		key := path[strings.LastIndex(path, "/")+1:]
		isFieldExpression := key == "when" || key == "forEach" || key == "wait" || key == "items" || key == "batchSize" || key == "condition" || key == "context" || key == "value"
		if !strings.Contains(typed, "$") && !isFieldExpression {
			return
		}
		constructType := ""
		if awaitWait {
			constructType = "await"
		}
		allowNumeric := numericLiteralAllowedForField(key, constructType)
		if _, ok := parseExpressionValue(typed, version, loopContext, allowNumeric); ok {
			if allowance, exists := allowances[path]; exists {
				t.Errorf("legacy expression at %s is now core grammar; remove stale allowance (%s)", path, allowance.reason)
				seen[path] = true
			}
			return
		}
		if allowance, exists := allowances[path]; exists {
			seen[path] = typed == allowance.value
			if typed != allowance.value {
				t.Errorf("legacy expression at %s changed from allowlisted %q to %q", path, allowance.value, typed)
			}
			return
		}
		t.Errorf("expression at %s is not valid under UWS %s grammar: %q", path, version, typed)
	}
}

func jsonPointer(path, token string) string {
	token = strings.ReplaceAll(token, "~", "~0")
	token = strings.ReplaceAll(token, "/", "~1")
	return path + "/" + token
}
