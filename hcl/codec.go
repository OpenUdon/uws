// Package hcl renders and independently verifies inert UWS HCL views from
// exact JSON/YAML bytes. It supplies no validation, approval or execution.
package hcl

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/OpenUdon/uws/internal/strictjson"
	"github.com/OpenUdon/uws/uws1"
	"gopkg.in/yaml.v3"
)

const ContractVersion = "uws.hcl-view.v1"
const MaxSourceBytes = 8 << 20
const MaxViewBytes = 8 << 20
const maxNodes = 100000
const maxDepth = 100

var ErrCodec = errors.New("invalid, unsupported or bounded HCL presentation contract")

type Format string

const (
	JSON Format = "json"
	YAML Format = "yaml"
)

type Source struct {
	Format Format
	Bytes  []byte
}

// Options identifies the exact consuming codec source revision. Revision is
// caller-supplied provenance, not proof of which binary ran; workers bind that.
type Options struct{ Revision string }

type Provenance struct {
	Version       string `json:"version"`
	Format        Format `json:"format"`
	SourceSHA256  string `json:"source_sha256"`
	CodecRevision string `json:"codec_revision"`
	ViewSHA256    string `json:"view_sha256"`
}

type View struct {
	HCL        []byte     `json:"hcl"`
	Provenance Provenance `json:"provenance"`
}

// Render returns no view unless independent inert parsing proves the complete
// source projection, including exact numeric lexemes, matches its HCL bytes.
func Render(ctx context.Context, source Source, options Options) (View, error) {
	ctx = codecContext(ctx)
	if !revisionValid(options.Revision) {
		return View{}, ErrCodec
	}
	value, err := sourceValue(ctx, source)
	if err != nil {
		return View{}, err
	}
	writer := viewWriter{budget: workBudget{ctx: ctx}}
	if err := writer.body(value, reflect.TypeOf(uws1.Document{}), "", 0); err != nil {
		return View{}, err
	}
	view := View{HCL: append([]byte(nil), writer.Bytes()...), Provenance: Provenance{Version: ContractVersion, Format: source.Format, SourceSHA256: digest(source.Bytes), CodecRevision: options.Revision, ViewSHA256: digest(writer.Bytes())}}
	if err := Verify(ctx, source, view, options); err != nil {
		return View{}, err
	}
	return view, nil
}

// Verify checks exact source/codec/view provenance and independently reconstructs
// all values. Hashes alone cannot prove correspondence or grant authority.
func Verify(ctx context.Context, source Source, view View, options Options) error {
	ctx = codecContext(ctx)
	if err := ctx.Err(); err != nil {
		return err
	}
	if !revisionValid(options.Revision) || len(view.HCL) > MaxViewBytes || view.Provenance != (Provenance{Version: ContractVersion, Format: source.Format, SourceSHA256: digest(source.Bytes), CodecRevision: options.Revision, ViewSHA256: digest(view.HCL)}) {
		return ErrCodec
	}
	want, err := sourceValue(ctx, source)
	if err != nil {
		return err
	}
	got, err := hclValue(ctx, view.HCL)
	if err != nil {
		return err
	}
	// json.Number is compared as a lexical string; no float/JCS normalization.
	if !reflect.DeepEqual(want, got) {
		return ErrCodec
	}
	return ctx.Err()
}

// Import reconstructs JSON from the supported inert typed HCL subset. It never
// evaluates functions/variables or supplies semantic validation/authority.
// Deprecated: author new documents as JSON/YAML and use verified views.
func Import(ctx context.Context, data []byte) ([]byte, error) {
	ctx = codecContext(ctx)
	value, err := hclValue(ctx, data)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(value)
	if err != nil || len(out) > MaxSourceBytes {
		return nil, ErrCodec
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func codecContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func revisionValid(s string) bool {
	if len(s) != 40 || strings.ToLower(s) != s {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

type workBudget struct {
	ctx   context.Context
	nodes int
}

func (b *workBudget) take(depth int) error {
	if err := b.ctx.Err(); err != nil {
		return err
	}
	b.nodes++
	if b.nodes > maxNodes || depth > maxDepth {
		return ErrCodec
	}
	return nil
}

func sourceValue(ctx context.Context, source Source) (map[string]any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(source.Bytes) == 0 || len(source.Bytes) > MaxSourceBytes || !utf8.Valid(source.Bytes) {
		return nil, ErrCodec
	}
	b := workBudget{ctx: ctx}
	var value any
	var err error
	switch source.Format {
	case JSON:
		d := json.NewDecoder(bytes.NewReader(source.Bytes))
		d.UseNumber()
		value, err = jsonValue(d, &b, 0)
		if err == nil {
			if _, e := d.Token(); e != io.EOF {
				err = ErrCodec
			}
			if err == nil && strictjson.ValidateSingleValue(source.Bytes) != nil {
				err = ErrCodec
			}
		}
	case YAML:
		var n yaml.Node
		d := yaml.NewDecoder(bytes.NewReader(source.Bytes))
		if d.Decode(&n) != nil {
			return nil, ErrCodec
		}
		var extra yaml.Node
		if d.Decode(&extra) != io.EOF {
			return nil, ErrCodec
		}
		value, err = yamlValue(&n, &b, 0)
	default:
		return nil, ErrCodec
	}
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, ErrCodec
	}
	object, ok := value.(map[string]any)
	if !ok {
		return nil, ErrCodec
	}
	return object, nil
}

func jsonValue(d *json.Decoder, b *workBudget, depth int) (any, error) {
	if err := b.take(depth); err != nil {
		return nil, err
	}
	token, err := d.Token()
	if err != nil {
		return nil, ErrCodec
	}
	delim, container := token.(json.Delim)
	if !container {
		return token, nil
	}
	switch delim {
	case '{':
		out := map[string]any{}
		for d.More() {
			token, err := d.Token()
			if err != nil {
				return nil, ErrCodec
			}
			key, ok := token.(string)
			if !ok {
				return nil, ErrCodec
			}
			if _, exists := out[key]; exists {
				return nil, ErrCodec
			}
			v, err := jsonValue(d, b, depth+1)
			if err != nil {
				return nil, err
			}
			out[key] = v
		}
		if token, err := d.Token(); err != nil || token != json.Delim('}') {
			return nil, ErrCodec
		}
		return out, nil
	case '[':
		out := []any{}
		for d.More() {
			v, err := jsonValue(d, b, depth+1)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		if token, err := d.Token(); err != nil || token != json.Delim(']') {
			return nil, ErrCodec
		}
		return out, nil
	}
	return nil, ErrCodec
}

func yamlValue(n *yaml.Node, b *workBudget, depth int) (any, error) {
	if err := b.take(depth); err != nil {
		return nil, err
	}
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 1 {
			return yamlValue(n.Content[0], b, depth+1)
		}
	case yaml.MappingNode:
		out := map[string]any{}
		for i := 0; i < len(n.Content); i += 2 {
			key := n.Content[i]
			if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
				return nil, ErrCodec
			}
			if _, exists := out[key.Value]; exists {
				return nil, ErrCodec
			}
			v, err := yamlValue(n.Content[i+1], b, depth+1)
			if err != nil {
				return nil, err
			}
			out[key.Value] = v
		}
		return out, nil
	case yaml.SequenceNode:
		out := []any{}
		for _, child := range n.Content {
			v, err := yamlValue(child, b, depth+1)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	case yaml.ScalarNode:
		switch n.Tag {
		case "!!str":
			return n.Value, nil
		case "!!null":
			return nil, nil
		case "!!bool":
			return strings.EqualFold(n.Value, "true"), nil
		case "!!int", "!!float":
			if numberPattern.MatchString(n.Value) {
				return json.Number(n.Value), nil
			}
		}
	}
	// Aliases, merges, tags and alternate numeric spellings are unsupported.
	return nil, ErrCodec
}

var numberPattern = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$`)
var identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)

var legacyDollarKeys = map[string]bool{"$ref": true, "$id": true, "$schema": true, "$defs": true, "$comment": true, "$vocabulary": true, "$anchor": true, "$dynamicRef": true, "$dynamicAnchor": true}

func encodeKey(key string) string {
	if strings.HasPrefix(key, "$") {
		if legacyDollarKeys[key] {
			return "_" + key[1:]
		}
		return "__dollar__" + key[1:]
	}
	if strings.HasPrefix(key, "__uws_literal__") || strings.HasPrefix(key, "__dollar__") || legacyDollarKeys["$"+strings.TrimPrefix(key, "_")] && strings.HasPrefix(key, "_") {
		return "__uws_literal__" + key
	}
	return key
}
func decodeKey(key string) string {
	if strings.HasPrefix(key, "__uws_literal__") {
		return key[len("__uws_literal__"):]
	}
	if strings.HasPrefix(key, "__dollar__") {
		return "$" + key[len("__dollar__"):]
	}
	if strings.HasPrefix(key, "_") && legacyDollarKeys["$"+key[1:]] {
		return "$" + key[1:]
	}
	return key
}

func quote(s string) string {
	data, _ := json.Marshal(s)
	return strings.ReplaceAll(strings.ReplaceAll(string(data), "${", "$${"), "%{", "%%{")
}

// Field mappings are read from the public model tags without decoding source
// values through its legacy custom unmarshaling/float64 paths.
type field struct {
	jsonName, hclName, mode string
	kind                    reflect.Type
}

func fields(kind reflect.Type) []field {
	for kind.Kind() == reflect.Pointer {
		kind = kind.Elem()
	}
	if kind.Kind() != reflect.Struct {
		return nil
	}
	var out []field
	for i := 0; i < kind.NumField(); i++ {
		f := kind.Field(i)
		if f.Anonymous {
			out = append(out, fields(f.Type)...)
			continue
		}
		j := strings.Split(f.Tag.Get("json"), ",")[0]
		h := strings.Split(f.Tag.Get("hcl"), ",")
		if j == "" || j == "-" || len(h) == 0 || h[0] == "-" || h[0] == "" {
			continue
		}
		mode := ""
		if len(h) > 1 {
			mode = h[1]
		}
		out = append(out, field{j, h[0], mode, f.Type})
	}
	return out
}
func indirect(kind reflect.Type) reflect.Type {
	for kind.Kind() == reflect.Pointer {
		kind = kind.Elem()
	}
	return kind
}
