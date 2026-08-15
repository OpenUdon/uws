package versions

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"golang.org/x/mod/module"
	"gopkg.in/yaml.v3"
)

const uwsModulePath = "github.com/OpenUdon/uws"

//go:embed *.json
var embeddedSchemas embed.FS

var (
	browserSchemaOnce  sync.Once
	browserSchema      *jsonschema.Schema
	browserSchemaErr   error
	authSchemaOnce     sync.Once
	authSchema         *jsonschema.Schema
	authSchemaErr      error
	authCallSchemaOnce sync.Once
	authCallSchema     *jsonschema.Schema
	authCallSchemaErr  error
)

const maxBrowserAuthenticationProfileBytes = 1 << 20

// PathForVersion returns the best local schema path for a UWS document version.
func PathForVersion(anchorDir, version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		version = "1.0.0"
	}
	return pathForSchemaName(anchorDir, version+".json")
}

// PathForRuntimeSupplement returns the best local schema path for a runtime supplement profile.
func PathForRuntimeSupplement(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, runtimeSupplementSchemaName(profile))
}

// PathForAnsibleModuleCallSupplement returns the best local schema path for an
// Ansible module-call supplement profile.
func PathForAnsibleModuleCallSupplement(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, ansibleModuleCallSupplementSchemaName(profile))
}

// PathForAnsibleArgspec returns the best local schema path for an Ansible
// argspec document.
func PathForAnsibleArgspec(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, familySchemaName(profile, "ansible", "1.0"))
}

// PathForAnsibleSourceProfile returns the best local schema path for an
// Ansible argspec document.
//
// Deprecated: use PathForAnsibleArgspec. Ansible argspecs are not source
// profiles as of UWS 1.7.
func PathForAnsibleSourceProfile(anchorDir, profile string) string {
	return PathForAnsibleArgspec(anchorDir, profile)
}

// PathForBrowserSourceProfile returns the best local schema path for a browser
// source profile.
func PathForBrowserSourceProfile(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, familySchemaName(profile, "browser", "1.5"))
}

// BrowserSourceProfileSchema returns an independent copy of the embedded
// browser-profile JSON Schema selected by profile. An empty profile selects
// the current uws.browser.1.5 contract.
func BrowserSourceProfileSchema(profile string) ([]byte, error) {
	name := familySchemaName(profile, "browser", "1.5")
	data, err := embeddedSchemas.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("load browser source profile schema %q: %w", profile, err)
	}
	return append([]byte(nil), data...), nil
}

// PathForBrowserAuthenticationProfile returns the best local schema path for
// a portable browser-authentication recipe.
func PathForBrowserAuthenticationProfile(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, familySchemaName(profile, "browser-authentication", "1.0"))
}

// BrowserAuthenticationProfileSchema returns an independent copy of the
// embedded uws.browser-authentication.1.0 JSON Schema.
func BrowserAuthenticationProfileSchema(profile string) ([]byte, error) {
	name := familySchemaName(profile, "browser-authentication", "1.0")
	data, err := embeddedSchemas.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("load browser authentication profile schema %q: %w", profile, err)
	}
	return append([]byte(nil), data...), nil
}

// PathForBrowserAuthenticationCallSupplement returns the best local schema
// path for the browser-authentication-call operation supplement.
func PathForBrowserAuthenticationCallSupplement(anchorDir, profile string) string {
	return pathForSchemaName(anchorDir, familySchemaName(profile, "browser-authentication-call", "1.0"))
}

// BrowserAuthenticationCallSupplementSchema returns an independent copy of
// the embedded uws.browser-authentication-call.1.0 JSON Schema.
func BrowserAuthenticationCallSupplementSchema(profile string) ([]byte, error) {
	name := familySchemaName(profile, "browser-authentication-call", "1.0")
	data, err := embeddedSchemas.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("load browser authentication call schema %q: %w", profile, err)
	}
	return append([]byte(nil), data...), nil
}

// ValidateBrowserSourceProfile validates one JSON or YAML browser-profile
// document against the embedded uws.browser.1.5 schema. It validates portable
// wire shape only; freshness, review evidence, registry lifecycle, sessions,
// and execution policy remain downstream responsibilities.
func ValidateBrowserSourceProfile(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("browser source profile document is empty")
	}
	value, err := decodeSingleJSONOrYAMLDocument(data)
	if err != nil {
		return fmt.Errorf("decode browser source profile: %w", err)
	}
	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(value))
	if err != nil {
		return fmt.Errorf("decode browser source profile as JSON: %w", err)
	}
	schema, err := compiledBrowserSourceProfileSchema()
	if err != nil {
		return err
	}
	if err := schema.Validate(document); err != nil {
		return fmt.Errorf("validate browser source profile: %w", err)
	}
	return nil
}

// ValidateBrowserAuthenticationProfile validates one portable, secret-free
// browser authentication recipe. In addition to JSON Schema validation it
// enforces exact safe origins and the 1 MiB document bound.
func ValidateBrowserAuthenticationProfile(data []byte) error {
	if len(data) > maxBrowserAuthenticationProfileBytes {
		return fmt.Errorf("browser authentication profile exceeds %d bytes", maxBrowserAuthenticationProfileBytes)
	}
	value, document, err := decodeSchemaDocument(data, "browser authentication profile")
	if err != nil {
		return err
	}
	schema, err := compiledBrowserAuthenticationProfileSchema()
	if err != nil {
		return err
	}
	if err := schema.Validate(document); err != nil {
		return fmt.Errorf("validate browser authentication profile: %w", err)
	}
	root, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("browser authentication profile must be an object")
	}
	info, _ := root["info"].(map[string]any)
	declaredOrigins := make(map[string]struct{})
	originCount := 0
	for _, field := range []string{"applicationOrigins", "authenticationOrigins"} {
		origins, _ := info[field].([]any)
		originCount += len(origins)
		for i, raw := range origins {
			origin, _ := raw.(string)
			if err := validateAuthenticationOrigin(origin); err != nil {
				return fmt.Errorf("info.%s[%d]: %w", field, i, err)
			}
			declaredOrigins[canonicalAuthenticationOrigin(origin)] = struct{}{}
		}
	}
	if originCount > 32 {
		return fmt.Errorf("info origins: combined applicationOrigins and authenticationOrigins exceed 32")
	}
	credentialSlots, _ := root["credentialSlots"].(map[string]any)
	flows, _ := root["flows"].(map[string]any)
	for name, raw := range flows {
		flow, _ := raw.(map[string]any)
		effects, _ := flow["effects"].([]any)
		hasMFAEffect := containsString(effects, "sends_mfa_challenge")
		hasChallenge := false
		sequence, _ := flow["sequence"].([]any)
		for i, rawStep := range sequence {
			step, _ := rawStep.(map[string]any)
			if rawNavigate, ok := step["navigate"]; ok {
				navigate, _ := rawNavigate.(string)
				if err := validateAuthenticationTarget(navigate, declaredOrigins); err != nil {
					return fmt.Errorf("flows.%s.sequence[%d].navigate: %w", name, i, err)
				}
			}
			if rawType, ok := step["type_credential"]; ok {
				typeStep, _ := rawType.(map[string]any)
				slot, _ := typeStep["slot"].(string)
				if _, ok := credentialSlots[slot]; !ok {
					return fmt.Errorf("flows.%s.sequence[%d].type_credential.slot: undeclared credential slot %q", name, i, slot)
				}
			}
			if rawChallenge, ok := step["challenge"]; ok {
				hasChallenge = true
				challenge, _ := rawChallenge.(map[string]any)
				if slot, _ := challenge["slot"].(string); slot != "" {
					rawSlot, ok := credentialSlots[slot]
					if !ok {
						return fmt.Errorf("flows.%s.sequence[%d].challenge.slot: undeclared credential slot %q", name, i, slot)
					}
					slotDef, _ := rawSlot.(map[string]any)
					if kind, _ := slotDef["kind"].(string); kind != "totp_seed" {
						return fmt.Errorf("flows.%s.sequence[%d].challenge.slot: TOTP requires a totp_seed slot", name, i)
					}
				}
			}
		}
		if hasChallenge != hasMFAEffect {
			return fmt.Errorf("flows.%s.effects: sends_mfa_challenge must be present exactly when the flow has a challenge step", name)
		}
		success, _ := flow["success"].(map[string]any)
		origin, _ := success["origin"].(string)
		if err := validateAuthenticationOrigin(origin); err != nil {
			return fmt.Errorf("flows.%s.success.origin: %w", name, err)
		}
		if _, ok := declaredOrigins[canonicalAuthenticationOrigin(origin)]; !ok {
			return fmt.Errorf("flows.%s.success.origin: origin is not declared by info", name)
		}
	}
	return nil
}

// ValidateBrowserAuthenticationCallSupplement validates the extension payload
// envelope used by an explicit authentication operation.
func ValidateBrowserAuthenticationCallSupplement(data []byte) error {
	value, document, err := decodeSchemaDocument(data, "browser authentication call")
	if err != nil {
		return err
	}
	schema, err := compiledBrowserAuthenticationCallSupplementSchema()
	if err != nil {
		return err
	}
	if err := schema.Validate(document); err != nil {
		return fmt.Errorf("validate browser authentication call: %w", err)
	}
	root, _ := value.(map[string]any)
	call, _ := root["x-uws-browser-authentication"].(map[string]any)
	profilePath, _ := call["profile"].(string)
	clean := path.Clean(profilePath)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || path.IsAbs(clean) {
		return fmt.Errorf("x-uws-browser-authentication.profile must be a safe relative path")
	}
	return nil
}

func decodeSchemaDocument(data []byte, label string) (any, any, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, nil, fmt.Errorf("%s document is empty", label)
	}
	encoded, err := decodeSingleJSONOrYAMLDocument(data)
	if err != nil {
		return nil, nil, fmt.Errorf("decode %s: %w", label, err)
	}
	var value any
	if err := json.Unmarshal(encoded, &value); err != nil {
		return nil, nil, fmt.Errorf("decode %s value: %w", label, err)
	}
	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(encoded))
	if err != nil {
		return nil, nil, fmt.Errorf("decode %s as JSON: %w", label, err)
	}
	return value, document, nil
}

func validateAuthenticationOrigin(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid origin: %w", err)
	}
	if parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || (parsed.Path != "" && parsed.Path != "/") {
		return fmt.Errorf("must be an exact origin without credentials, path, query, or fragment")
	}
	host := parsed.Hostname()
	loopback := host == "localhost"
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		loopback = true
	}
	if parsed.Scheme != "https" && !(parsed.Scheme == "http" && loopback) {
		return fmt.Errorf("must use https (http is allowed only for loopback)")
	}
	return nil
}

func canonicalAuthenticationOrigin(raw string) string {
	parsed, _ := url.Parse(raw)
	host := strings.ToLower(parsed.Hostname())
	port := parsed.Port()
	if (parsed.Scheme == "https" && port == "443") || (parsed.Scheme == "http" && port == "80") {
		port = ""
	}
	if port != "" {
		host = net.JoinHostPort(host, port)
	} else if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return strings.ToLower(parsed.Scheme) + "://" + host
}

func validateAuthenticationTarget(raw string, origins map[string]struct{}) error {
	parsed, err := url.Parse(raw)
	if err != nil || !parsed.IsAbs() || parsed.Host == "" {
		return fmt.Errorf("must be an absolute URL")
	}
	if parsed.User != nil || (parsed.Scheme != "https" && !(parsed.Scheme == "http" && isLoopbackHost(parsed.Hostname()))) {
		return fmt.Errorf("must use a declared safe origin")
	}
	if _, ok := origins[canonicalAuthenticationOrigin(raw)]; !ok {
		return fmt.Errorf("target origin is not declared by info")
	}
	return nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func containsString(values []any, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func compiledBrowserAuthenticationProfileSchema() (*jsonschema.Schema, error) {
	authSchemaOnce.Do(func() {
		authSchema, authSchemaErr = compileEmbeddedSchema("browser-authentication.1.0.json", BrowserAuthenticationProfileSchema)
	})
	return authSchema, authSchemaErr
}

func compiledBrowserAuthenticationCallSupplementSchema() (*jsonschema.Schema, error) {
	authCallSchemaOnce.Do(func() {
		authCallSchema, authCallSchemaErr = compileEmbeddedSchema("browser-authentication-call.1.0.json", BrowserAuthenticationCallSupplementSchema)
	})
	return authCallSchema, authCallSchemaErr
}

func compileEmbeddedSchema(name string, load func(string) ([]byte, error)) (*jsonschema.Schema, error) {
	data, err := load("")
	if err != nil {
		return nil, err
	}
	document, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode embedded %s schema: %w", name, err)
	}
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(name, document); err != nil {
		return nil, fmt.Errorf("register embedded %s schema: %w", name, err)
	}
	schema, err := compiler.Compile(name)
	if err != nil {
		return nil, fmt.Errorf("compile embedded %s schema: %w", name, err)
	}
	return schema, nil
}

func decodeSingleJSONOrYAMLDocument(data []byte) ([]byte, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("multiple YAML documents are not supported")
		}
		return nil, err
	}
	if value == nil {
		return nil, fmt.Errorf("document is empty")
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func compiledBrowserSourceProfileSchema() (*jsonschema.Schema, error) {
	browserSchemaOnce.Do(func() {
		data, err := BrowserSourceProfileSchema("")
		if err != nil {
			browserSchemaErr = err
			return
		}
		document, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
		if err != nil {
			browserSchemaErr = fmt.Errorf("decode embedded browser source profile schema: %w", err)
			return
		}
		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource("browser.1.5.json", document); err != nil {
			browserSchemaErr = fmt.Errorf("register embedded browser source profile schema: %w", err)
			return
		}
		browserSchema, browserSchemaErr = compiler.Compile("browser.1.5.json")
		if browserSchemaErr != nil {
			browserSchemaErr = fmt.Errorf("compile embedded browser source profile schema: %w", browserSchemaErr)
		}
	})
	return browserSchema, browserSchemaErr
}

func runtimeSupplementSchemaName(profile string) string {
	profile = strings.TrimSpace(profile)
	if profile == "" {
		return "runtime.1.0.json"
	}
	profile = strings.TrimSuffix(profile, ".json")
	profile = strings.TrimPrefix(profile, "uws.")
	if !strings.HasPrefix(profile, "runtime.") {
		profile = "runtime." + profile
	}
	return profile + ".json"
}

func ansibleModuleCallSupplementSchemaName(profile string) string {
	profile = strings.TrimSpace(profile)
	if profile == "" {
		return "ansible-module-call.1.0.json"
	}
	profile = strings.TrimSuffix(profile, ".json")
	profile = strings.TrimPrefix(profile, "uws.")
	if !strings.HasPrefix(profile, "ansible-module-call.") {
		profile = "ansible-module-call." + profile
	}
	return profile + ".json"
}

func familySchemaName(profile, name, defaultVersion string) string {
	profile = strings.TrimSpace(profile)
	if profile == "" {
		profile = defaultVersion
	}
	profile = strings.TrimSuffix(profile, ".json")
	profile = strings.TrimPrefix(profile, "uws.")
	if !strings.HasPrefix(profile, name+".") {
		profile = name + "." + profile
	}
	return profile + ".json"
}

func pathForSchemaName(anchorDir, name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "" || name == "." {
		name = "1.0.0.json"
	}
	if dir := strings.TrimSpace(os.Getenv("UWS_SCHEMA_DIR")); dir != "" {
		return filepath.Join(dir, name)
	}
	if dir := strings.TrimSpace(os.Getenv("OPENUDON_UWS_SCHEMA_DIR")); dir != "" {
		return filepath.Join(dir, name)
	}
	if path, ok := packageSchemaPath(name); ok {
		return path
	}
	if path, ok := moduleCacheSchemaPath(name); ok {
		return path
	}
	if path, ok := embeddedSchemaPath(name); ok {
		return path
	}
	return filepath.Join(anchorDir, "..", "uws", "versions", name)
}

func packageSchemaPath(name string) (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", false
	}
	path := filepath.Join(filepath.Dir(file), name)
	if _, err := os.Stat(path); err == nil {
		return path, true
	}
	return "", false
}

func moduleCacheSchemaPath(name string) (string, bool) {
	version, ok := uwsModuleVersion()
	if !ok {
		return "", false
	}
	path, err := escapedModuleCachePath(uwsModulePath, version)
	if err != nil {
		return "", false
	}
	schema := filepath.Join(path, "versions", name)
	if _, err := os.Stat(schema); err == nil {
		return schema, true
	}
	return "", false
}

func uwsModuleVersion() (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, dep := range info.Deps {
		if dep.Path != uwsModulePath {
			continue
		}
		if dep.Version != "" {
			return dep.Version, true
		}
		if dep.Replace != nil && dep.Replace.Version != "" {
			return dep.Replace.Version, true
		}
	}
	return "", false
}

func escapedModuleCachePath(path, version string) (string, error) {
	escapedPath, err := module.EscapePath(path)
	if err != nil {
		return "", err
	}
	escapedVersion, err := module.EscapeVersion(version)
	if err != nil {
		return "", err
	}
	return filepath.Join(moduleCacheDir(), escapedPath+"@"+escapedVersion), nil
}

func moduleCacheDir() string {
	if dir := strings.TrimSpace(os.Getenv("GOMODCACHE")); dir != "" {
		return dir
	}
	gopath := strings.TrimSpace(os.Getenv("GOPATH"))
	if gopath == "" {
		if home, err := os.UserHomeDir(); err == nil {
			gopath = filepath.Join(home, "go")
		}
	}
	if gopath == "" {
		return ""
	}
	first := filepath.SplitList(gopath)[0]
	if first == "" {
		return ""
	}
	return filepath.Join(first, "pkg", "mod")
}

func embeddedSchemaPath(name string) (string, bool) {
	data, err := embeddedSchemas.ReadFile(filepath.ToSlash(name))
	if err != nil {
		return "", false
	}
	sum := sha256.Sum256(data)
	dir := filepath.Join(os.TempDir(), "uws-schema", fmt.Sprintf("%x", sum[:8]))
	path := filepath.Join(dir, name)
	if existing, err := os.ReadFile(path); err == nil && bytes.Equal(existing, data) {
		return path, true
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", false
	}
	tmp, err := os.CreateTemp(dir, name+".*.tmp")
	if err != nil {
		return "", false
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return "", false
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return "", false
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		os.Remove(tmpName)
		return "", false
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return "", false
	}
	return path, true
}
