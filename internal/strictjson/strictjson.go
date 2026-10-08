// Package strictjson provides bounded, unambiguous JSON document checks used
// by versioned portable artifacts.
package strictjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

// ValidateSingleValue rejects invalid UTF-8, duplicate object member names,
// malformed JSON, and trailing values. Duplicate names are compared after
// JSON string escapes have been decoded. Containers are limited to depth 100.
func ValidateSingleValue(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("JSON document is empty")
	}
	if !utf8.Valid(data) {
		return fmt.Errorf("JSON document is not valid UTF-8")
	}
	if err := validateSurrogateEscapes(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := consumeValue(decoder); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return fmt.Errorf("JSON document contains trailing values")
		}
		return err
	}
	return nil
}

func validateSurrogateEscapes(data []byte) error {
	for i := 0; i < len(data); i++ {
		if data[i] != '"' {
			continue
		}
	stringScan:
		for i++; i < len(data); i++ {
			switch data[i] {
			case '"':
				break stringScan
			case '\\':
				i++
				if i >= len(data) || data[i] != 'u' {
					continue
				}
				if i+4 >= len(data) {
					return fmt.Errorf("truncated JSON Unicode escape")
				}
				value, err := parseHexCodeUnit(data[i+1 : i+5])
				if err != nil {
					return err
				}
				i += 4
				if value >= 0xdc00 && value <= 0xdfff {
					return fmt.Errorf("JSON string contains an unpaired low surrogate")
				}
				if value < 0xd800 || value > 0xdbff {
					continue
				}
				if i+6 >= len(data) || data[i+1] != '\\' || data[i+2] != 'u' {
					return fmt.Errorf("JSON string contains an unpaired high surrogate")
				}
				low, err := parseHexCodeUnit(data[i+3 : i+7])
				if err != nil {
					return err
				}
				if low < 0xdc00 || low > 0xdfff {
					return fmt.Errorf("JSON string contains an unpaired high surrogate")
				}
				i += 6
			}
		}
	}
	return nil
}

func parseHexCodeUnit(data []byte) (uint16, error) {
	if len(data) != 4 {
		return 0, fmt.Errorf("invalid JSON Unicode escape")
	}
	var value uint16
	for _, char := range data {
		value <<= 4
		switch {
		case char >= '0' && char <= '9':
			value |= uint16(char - '0')
		case char >= 'a' && char <= 'f':
			value |= uint16(char-'a') + 10
		case char >= 'A' && char <= 'F':
			value |= uint16(char-'A') + 10
		default:
			return 0, fmt.Errorf("invalid JSON Unicode escape")
		}
	}
	return value, nil
}

func consumeValue(decoder *json.Decoder) error {
	// Decoder.Token scans values without recursively decoding containers. Keep
	// our own bounded stack for duplicate member checks and object key state.
	type frame struct {
		end  json.Delim
		seen map[string]struct{}
		key  bool
	}
	stack := make([]frame, 0, 100)
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		if len(stack) > 0 {
			parent := &stack[len(stack)-1]
			if end, ok := token.(json.Delim); ok && end == parent.end {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 {
					return nil
				}
				continue
			}
			if parent.seen != nil && parent.key {
				key, ok := token.(string)
				if !ok {
					return fmt.Errorf("JSON object member name is not a string")
				}
				if _, exists := parent.seen[key]; exists {
					return fmt.Errorf("JSON object contains duplicate member %q", key)
				}
				parent.seen[key] = struct{}{}
				parent.key = false
				continue
			}
			parent.key = true
		}
		if delim, ok := token.(json.Delim); ok {
			if len(stack) == 100 {
				return fmt.Errorf("JSON nesting exceeds depth 100")
			}
			switch delim {
			case '{':
				stack = append(stack, frame{end: '}', seen: map[string]struct{}{}, key: true})
			case '[':
				stack = append(stack, frame{end: ']'})
			default:
				return fmt.Errorf("unexpected JSON delimiter %q", delim)
			}
		} else if len(stack) == 0 {
			return nil
		}
	}
}
