package schemas

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"unicode"
)

const browser19MaxSafeInteger = 9007199254740991

type browser19TemplateKind uint8

const (
	browser19Parameter browser19TemplateKind = iota
	browser19LiteralOpen
	browser19LiteralClose
)

type browser19TemplateToken struct {
	start int
	end   int
	name  string
	kind  browser19TemplateKind
}

func validateBrowser19Templates(root map[string]any) error {
	if err := validateBrowser19IntegerDefaults(root); err != nil {
		return err
	}
	return validateBrowserTemplates(root, parseBrowser19TemplateNames, validateBrowser19NavigateTemplate, validateBrowser19Text)
}

func validateBrowser19IntegerDefaults(root map[string]any) error {
	actions, _ := root["actions"].(map[string]any)
	for actionName, rawAction := range actions {
		action, _ := rawAction.(map[string]any)
		parameters, _ := action["parameters"].(map[string]any)
		properties, _ := parameters["properties"].(map[string]any)
		for name, rawProperty := range properties {
			property, _ := rawProperty.(map[string]any)
			typeName, _ := property["type"].(string)
			if typeName != "integer" {
				continue
			}
			if value, ok := property["default"]; ok && !isBrowser19SafeInteger(value) {
				return fmt.Errorf("actions.%s.parameters.properties.%s.default: integer default must be a safe integer", actionName, name)
			}
		}
	}
	return nil
}

func isBrowser19SafeInteger(value any) bool {
	switch number := value.(type) {
	case int:
		return number >= -browser19MaxSafeInteger && number <= browser19MaxSafeInteger
	case int64:
		return number >= -browser19MaxSafeInteger && number <= browser19MaxSafeInteger
	case uint:
		return number <= browser19MaxSafeInteger
	case uint64:
		return number <= browser19MaxSafeInteger
	case float64:
		return !math.IsNaN(number) && !math.IsInf(number, 0) && math.Trunc(number) == number && number >= -browser19MaxSafeInteger && number <= browser19MaxSafeInteger
	case float32:
		converted := float64(number)
		return !math.IsNaN(converted) && !math.IsInf(converted, 0) && math.Trunc(converted) == converted && converted >= -browser19MaxSafeInteger && converted <= browser19MaxSafeInteger
	default:
		return false
	}
}

func parseBrowser19TemplateNames(value string) ([]string, error) {
	tokens, err := parseBrowser19TemplateTokens(value)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if token.kind == browser19Parameter {
			names = append(names, token.name)
		}
	}
	return names, nil
}

func parseBrowser19TemplateTokens(value string) ([]browser19TemplateToken, error) {
	var tokens []browser19TemplateToken
	for index := 0; index < len(value); {
		switch {
		case strings.HasPrefix(value[index:], "{{{{"):
			tokens = append(tokens, browser19TemplateToken{start: index, end: index + 4, kind: browser19LiteralOpen})
			index += 4
		case strings.HasPrefix(value[index:], "}}}}"):
			tokens = append(tokens, browser19TemplateToken{start: index, end: index + 4, kind: browser19LiteralClose})
			index += 4
		case strings.HasPrefix(value[index:], "{{"):
			endRelative := strings.Index(value[index+2:], "}}")
			if endRelative < 0 {
				return nil, fmt.Errorf("unmatched template opening braces")
			}
			end := index + 2 + endRelative
			name := value[index+2 : end]
			if strings.Contains(name, "{") || strings.Contains(name, "}") || !browserTemplateName.MatchString(name) {
				return nil, fmt.Errorf("invalid template placeholder %q", name)
			}
			tokens = append(tokens, browser19TemplateToken{start: index, end: end + 2, name: name, kind: browser19Parameter})
			index = end + 2
		case strings.HasPrefix(value[index:], "}}"):
			return nil, fmt.Errorf("unmatched template closing braces")
		default:
			index++
		}
	}
	return tokens, nil
}

func validateBrowser19Text(path, value string) error {
	if !strings.HasSuffix(path, ".type_text.value") && !strings.HasSuffix(path, ".confirmationPolicy.prompt") {
		return nil
	}
	for _, character := range value {
		if unicode.IsControl(character) || character == '\u2028' || character == '\u2029' || isBrowser19BidiControl(character) {
			return fmt.Errorf("contains a prohibited control, line-separator, or bidi-control character")
		}
	}
	return nil
}

func isBrowser19BidiControl(character rune) bool {
	switch character {
	case '\u061c', '\u200e', '\u200f':
		return true
	}
	return character >= '\u202a' && character <= '\u202e' || character >= '\u2066' && character <= '\u2069'
}

func validateBrowser19NavigateTemplate(target string) error {
	tokens, err := parseBrowser19TemplateTokens(target)
	if err != nil {
		return err
	}
	var marked strings.Builder
	type tokenMarker struct {
		value string
	}
	var markers []tokenMarker
	last := 0
	for _, token := range tokens {
		marked.WriteString(target[last:token.start])
		marker := fmt.Sprintf("UWSBROWSER19TOKEN%dEND", len(markers))
		markers = append(markers, tokenMarker{value: marker})
		marked.WriteString(marker)
		last = token.end
	}
	marked.WriteString(target[last:])
	parsed, err := url.Parse(marked.String())
	if err != nil {
		return fmt.Errorf("navigate URL template is malformed: %w", err)
	}
	for _, marker := range markers {
		userinfo := ""
		if parsed.User != nil {
			userinfo = parsed.User.String()
		}
		if strings.Contains(parsed.Scheme, marker.value) || strings.Contains(parsed.Host, marker.value) || strings.Contains(userinfo, marker.value) || strings.Contains(parsed.Fragment, marker.value) {
			return fmt.Errorf("navigate placeholders and brace escapes are forbidden in scheme, authority, or fragment")
		}
	}
	for _, pair := range strings.Split(parsed.RawQuery, "&") {
		key, value, hasValue := strings.Cut(pair, "=")
		for _, marker := range markers {
			if strings.Contains(key, marker.value) {
				return fmt.Errorf("navigate placeholders and brace escapes are forbidden in query parameter names")
			}
			if !hasValue && strings.Contains(value, marker.value) {
				return fmt.Errorf("navigate placeholders and brace escapes are permitted only in query parameter values")
			}
		}
	}
	for _, segment := range strings.Split(parsed.EscapedPath(), "/") {
		decoded, decodeErr := url.PathUnescape(segment)
		if decodeErr != nil {
			return fmt.Errorf("navigate path contains malformed percent encoding")
		}
		if decoded == "." || decoded == ".." {
			return fmt.Errorf("navigate path must not contain dot segments")
		}
	}
	return nil
}
