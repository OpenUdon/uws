package runtimes

import (
	"encoding/json"
	"fmt"

	"github.com/OpenUdon/uws/uws1"
)

const (
	// ProfileName is the public operation profile marker for this supplement.
	ProfileName = "uws.runtime.1.0"

	// ExtensionRuntime is the operation-level runtime metadata extension key.
	ExtensionRuntime = "x-uws-runtime"

	// ExtensionRuntimeConfig is the document/component-level runtime config key.
	ExtensionRuntimeConfig = "x-uws-runtime-config"
)

const (
	RuntimeTypeHTTP   = "http"
	RuntimeTypeSSH    = "ssh"
	RuntimeTypeCmd    = "cmd"
	RuntimeTypeFnct   = "fnct"
	RuntimeTypeFileIO = "fileio"
	RuntimeTypeSQL    = "sql"
	RuntimeTypeS3     = "s3"
	RuntimeTypeSMTP   = "smtp"
	RuntimeTypeDNS    = "dns"
	RuntimeTypeLDAPS  = "ldaps"
	RuntimeTypeSCP    = "scp"
	RuntimeTypeSFTP   = "sftp"
	RuntimeTypeLLM    = "llm"
)

// IsRuntimeType reports whether typeName names a runtime type defined by the
// UWS runtime supplement.
func IsRuntimeType(typeName string) bool {
	switch typeName {
	case RuntimeTypeHTTP, RuntimeTypeSSH, RuntimeTypeCmd, RuntimeTypeFnct,
		RuntimeTypeFileIO, RuntimeTypeSQL, RuntimeTypeS3, RuntimeTypeSMTP,
		RuntimeTypeDNS, RuntimeTypeLDAPS, RuntimeTypeSCP, RuntimeTypeSFTP,
		RuntimeTypeLLM:
		return true
	default:
		return false
	}
}

// ConfigRuntime is the typed payload for x-uws-runtime-config.
type ConfigRuntime struct {
	Provider *Provider              `json:"provider,omitempty" hcl:"provider,block"`
	Security []*SecurityRequirement `json:"security,omitempty" hcl:"security,block"`
}

// OperationRuntime is the typed payload for x-uws-runtime.
type OperationRuntime struct {
	Type               string                 `json:"type,omitempty" hcl:"type,optional"`
	IsJSON             bool                   `json:"isJson,omitempty" hcl:"isJson,optional"`
	Host               string                 `json:"host,omitempty" hcl:"host,optional"`
	Method             string                 `json:"method,omitempty" hcl:"method,optional"`
	Path               string                 `json:"path,omitempty" hcl:"path,optional"`
	PayloadRequired    *bool                  `json:"payloadRequired,omitempty" hcl:"payloadRequired,optional"`
	RequestMediaType   string                 `json:"requestMediaType,omitempty" hcl:"requestMediaType,optional"`
	ResponseMediaType  string                 `json:"responseMediaType,omitempty" hcl:"responseMediaType,optional"`
	ResponseStatusCode *int                   `json:"responseStatusCode,omitempty" hcl:"responseStatusCode,optional"`
	Command            string                 `json:"command,omitempty" hcl:"command,optional"`
	WorkingDir         string                 `json:"workingDir,omitempty" hcl:"workingDir,optional"`
	Function           string                 `json:"function,omitempty" hcl:"function,optional"`
	Workflow           string                 `json:"workflow,omitempty" hcl:"workflow,optional"`
	Arguments          []any                  `json:"arguments,omitempty" hcl:"arguments,optional"`
	Provider           *Provider              `json:"provider,omitempty" hcl:"provider,block"`
	Security           []*SecurityRequirement `json:"security,omitempty" hcl:"security,block"`
	QueryPars          *uws1.ParamSchema      `json:"queryPars,omitempty" hcl:"queryPars,block"`
	PathPars           *uws1.ParamSchema      `json:"pathPars,omitempty" hcl:"pathPars,block"`
	HeaderPars         *uws1.ParamSchema      `json:"headerPars,omitempty" hcl:"headerPars,block"`
	CookiePars         *uws1.ParamSchema      `json:"cookiePars,omitempty" hcl:"cookiePars,block"`
	PayloadPars        *uws1.ParamSchema      `json:"payloadPars,omitempty" hcl:"payloadPars,block"`
	ResponseBody       *uws1.ParamSchema      `json:"responseBody,omitempty" hcl:"responseBody,block"`
	ResponseHeaders    *uws1.ParamSchema      `json:"responseHeaders,omitempty" hcl:"responseHeaders,block"`
}

// Provider describes runtime provider selection metadata.
type Provider struct {
	Name       string         `json:"name,omitempty" hcl:"name,optional"`
	ServerURL  string         `json:"serverUrl,omitempty" hcl:"serverUrl,optional"`
	Appendices map[string]any `json:"appendices,omitempty" hcl:"appendices,optional"`
}

// SecurityRequirement describes a runtime security requirement binding.
type SecurityRequirement struct {
	Name       string          `json:"name,omitempty" hcl:"name,optional"`
	Scopes     []string        `json:"scopes,omitempty" hcl:"scopes,optional"`
	Scheme     *SecurityScheme `json:"scheme,omitempty" hcl:"scheme,block"`
	Initialize string          `json:"initialize,omitempty" hcl:"initialize,optional"`
	DataFile   string          `json:"dataFile,omitempty" hcl:"dataFile,optional"`
}

// SecurityScheme mirrors the OpenAPI security scheme fields commonly needed by
// runtime metadata.
type SecurityScheme struct {
	Type        string      `json:"type,omitempty" hcl:"type,optional"`
	Name        string      `json:"name,omitempty" hcl:"name,optional"`
	In          string      `json:"in,omitempty" hcl:"in,optional"`
	Scheme      string      `json:"scheme,omitempty" hcl:"scheme,optional"`
	Description string      `json:"description,omitempty" hcl:"description,optional"`
	Flows       *OAuthFlows `json:"flows,omitempty" hcl:"flows,block"`
}

// OAuthFlows contains OAuth flow metadata nested under SecurityScheme.
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" hcl:"implicit,block"`
	Password          *OAuthFlow `json:"password,omitempty" hcl:"password,block"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" hcl:"clientCredentials,block"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" hcl:"authorizationCode,block"`
}

// OAuthFlow contains OAuth endpoint and scope metadata.
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" hcl:"authorizationUrl,optional"`
	TokenURL         string            `json:"tokenUrl,omitempty" hcl:"tokenUrl,optional"`
	RefreshURL       string            `json:"refreshUrl,omitempty" hcl:"refreshUrl,optional"`
	Scopes           map[string]string `json:"scopes,omitempty" hcl:"scopes,optional"`
}

// ReadOperationExtension decodes x-uws-runtime from an extension map.
func ReadOperationExtension(extensions map[string]any) (*OperationRuntime, bool, error) {
	return readExtension[OperationRuntime](extensions, ExtensionRuntime)
}

// SetOperationExtension encodes x-uws-runtime into an extension map.
func SetOperationExtension(dst *map[string]any, value *OperationRuntime) error {
	return setExtension(dst, ExtensionRuntime, value)
}

// ReadConfigExtension decodes x-uws-runtime-config from an extension map.
func ReadConfigExtension(extensions map[string]any) (*ConfigRuntime, bool, error) {
	return readExtension[ConfigRuntime](extensions, ExtensionRuntimeConfig)
}

// SetConfigExtension encodes x-uws-runtime-config into an extension map.
func SetConfigExtension(dst *map[string]any, value *ConfigRuntime) error {
	return setExtension(dst, ExtensionRuntimeConfig, value)
}

func setExtension(dst *map[string]any, key string, value any) error {
	if value == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var generic any
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}
	if *dst == nil {
		*dst = make(map[string]any)
	}
	(*dst)[key] = generic
	return nil
}

func readExtension[T any](extensions map[string]any, key string) (*T, bool, error) {
	if len(extensions) == 0 {
		return nil, false, nil
	}
	value, ok := extensions[key]
	if !ok {
		return nil, false, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, false, fmt.Errorf("marshal %s extension: %w", key, err)
	}
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, false, fmt.Errorf("unmarshal %s extension: %w", key, err)
	}
	return &out, true, nil
}
