package schemas

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var browserTemplateName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]*$`)

// validateBrowser18Templates enforces the static half of the browser 1.8
// substitution contract. Concrete parameter values and defaults are checked
// by a runtime before it begins executing an action sequence.
func validateBrowser18Templates(root map[string]any) error {
	actions, _ := root["actions"].(map[string]any)
	for actionName, rawAction := range actions {
		action, _ := rawAction.(map[string]any)
		properties := browserParameterProperties(action["parameters"])
		allowed := make(map[string]bool)
		sequence, _ := action["sequence"].([]any)
		for index, rawStep := range sequence {
			step, _ := rawStep.(map[string]any)
			if rawNavigate, ok := step["navigate"]; ok {
				var target string
				path := fmt.Sprintf("actions.%s.sequence[%d].navigate", actionName, index)
				switch value := rawNavigate.(type) {
				case string:
					target = value
				case map[string]any:
					target, _ = value["url"].(string)
					path += ".url"
				}
				allowed[path] = true
				if err := validateBrowserNavigateTemplate(target); err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
			}
			for _, key := range []string{"type_text", "select_option"} {
				if _, ok := step[key].(map[string]any); ok {
					path := fmt.Sprintf("actions.%s.sequence[%d].%s.value", actionName, index, key)
					allowed[path] = true
				}
			}
		}
		if _, ok := action["confirmationPolicy"].(map[string]any); ok {
			allowed["actions."+actionName+".confirmationPolicy.prompt"] = true
		}
		var walk func(any, string) error
		walk = func(value any, path string) error {
			switch typed := value.(type) {
			case map[string]any:
				for key, child := range typed {
					if err := walk(child, path+"."+key); err != nil {
						return err
					}
				}
			case []any:
				for index, child := range typed {
					if err := walk(child, fmt.Sprintf("%s[%d]", path, index)); err != nil {
						return err
					}
				}
			case string:
				if !strings.ContainsAny(typed, "{}") {
					return nil
				}
				if !allowed[path] {
					return fmt.Errorf("%s: parameter templates are not permitted in this field", path)
				}
				names, err := parseBrowserTemplateNames(typed)
				if err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				for _, name := range names {
					typeName, ok := properties[name]
					if !ok {
						return fmt.Errorf("%s: template parameter %q is not declared in parameters.properties", path, name)
					}
					if typeName == "" {
						return fmt.Errorf("%s: parameter %q must have exactly one scalar type (string, boolean, integer, or number)", path, name)
					}
				}
			}
			return nil
		}
		if err := walk(rawAction, "actions."+actionName); err != nil {
			return err
		}
	}
	return nil
}

func browserParameterProperties(raw any) map[string]string {
	result := make(map[string]string)
	parameters, _ := raw.(map[string]any)
	properties, _ := parameters["properties"].(map[string]any)
	for name, rawProperty := range properties {
		property, _ := rawProperty.(map[string]any)
		typeName, _ := property["type"].(string)
		if _, hasOneOf := property["oneOf"]; hasOneOf {
			typeName = ""
		}
		if _, hasAnyOf := property["anyOf"]; hasAnyOf {
			typeName = ""
		}
		switch typeName {
		case "string", "boolean", "integer", "number":
			result[name] = typeName
		default:
			result[name] = ""
		}
	}
	return result
}

func parseBrowserTemplateNames(value string) ([]string, error) {
	var names []string
	for index := 0; index < len(value); {
		relative := strings.IndexAny(value[index:], "{}")
		if relative < 0 {
			break
		}
		brace := value[index+relative]
		if brace == '}' {
			if index+relative+1 >= len(value) || value[index+relative+1] != '}' {
				return nil, fmt.Errorf("unmatched template closing brace")
			}
			return nil, fmt.Errorf("unmatched template closing braces")
		}
		open := index + relative
		if open+1 >= len(value) || value[open+1] != '{' {
			return nil, fmt.Errorf("unmatched template opening brace")
		}
		end := strings.Index(value[open+2:], "}}")
		if end < 0 {
			return nil, fmt.Errorf("unmatched template opening braces")
		}
		end = open + 2 + end
		name := value[open+2 : end]
		if strings.Contains(name, "{") || strings.Contains(name, "}") || !browserTemplateName.MatchString(name) {
			return nil, fmt.Errorf("invalid template placeholder %q", name)
		}
		names = append(names, name)
		index = end + 2
	}
	return names, nil
}

func validateBrowserNavigateTemplate(target string) error {
	if _, err := parseBrowserTemplateNames(target); err != nil {
		return err
	}
	var markerNames []string
	var replaced strings.Builder
	for index := 0; index < len(target); {
		open := strings.Index(target[index:], "{{")
		if open < 0 {
			replaced.WriteString(target[index:])
			break
		}
		open += index
		replaced.WriteString(target[index:open])
		end := open + 2 + strings.Index(target[open+2:], "}}")
		marker := fmt.Sprintf("UWSPLACEHOLDER%dTOKEN", len(markerNames))
		markerNames = append(markerNames, marker)
		replaced.WriteString(marker)
		index = end + 2
	}
	marked := replaced.String()
	parsed, err := url.Parse(marked)
	if err != nil {
		return fmt.Errorf("navigate URL template is malformed: %w", err)
	}
	for _, marker := range markerNames {
		userinfo := ""
		if parsed.User != nil {
			userinfo = parsed.User.String()
		}
		if strings.Contains(parsed.Scheme, marker) || strings.Contains(parsed.Host, marker) || strings.Contains(userinfo, marker) || strings.Contains(parsed.Fragment, marker) {
			return fmt.Errorf("navigate placeholders are forbidden in scheme, authority, or fragment")
		}
	}
	for _, pair := range strings.Split(parsed.RawQuery, "&") {
		key, value, hasValue := strings.Cut(pair, "=")
		for _, marker := range markerNames {
			if strings.Contains(key, marker) {
				return fmt.Errorf("navigate placeholders are forbidden in query parameter names")
			}
			if !hasValue && strings.Contains(value, marker) {
				return fmt.Errorf("navigate placeholders are permitted only in query parameter values")
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
