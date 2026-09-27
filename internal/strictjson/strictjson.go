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
// JSON string escapes have been decoded.
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
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, isDelim := token.(json.Delim)
	if !isDelim {
		return nil
	}
	switch delim {
	case '{':
		seen := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("JSON object member name is not a string")
			}
			if _, exists := seen[key]; exists {
				return fmt.Errorf("JSON object contains duplicate member %q", key)
			}
			seen[key] = struct{}{}
			if err := consumeValue(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim('}') {
			return fmt.Errorf("malformed JSON object")
		}
	case '[':
		for decoder.More() {
			if err := consumeValue(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim(']') {
			return fmt.Errorf("malformed JSON array")
		}
	default:
		return fmt.Errorf("unexpected JSON delimiter %q", delim)
	}
	return nil
}
