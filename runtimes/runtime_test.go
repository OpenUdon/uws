package runtimes

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicRuntimeStructFieldsHaveHCLTags(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeOf(ConfigRuntime{}),
		reflect.TypeOf(OperationRuntime{}),
		reflect.TypeOf(Provider{}),
		reflect.TypeOf(SecurityRequirement{}),
		reflect.TypeOf(SecurityScheme{}),
		reflect.TypeOf(OAuthFlows{}),
		reflect.TypeOf(OAuthFlow{}),
	} {
		t.Run(typ.Name(), func(t *testing.T) {
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				if !field.IsExported() {
					continue
				}
				require.NotEmpty(t, field.Tag.Get("hcl"), "%s.%s must have an hcl tag", typ.Name(), field.Name)
			}
		})
	}
}

func TestIsRuntimeType(t *testing.T) {
	for _, typ := range []string{
		RuntimeTypeHTTP, RuntimeTypeSSH, RuntimeTypeCmd, RuntimeTypeFnct,
		RuntimeTypeFileIO, RuntimeTypeSQL, RuntimeTypeS3, RuntimeTypeSMTP,
		RuntimeTypeDNS, RuntimeTypeLDAPS, RuntimeTypeSCP, RuntimeTypeSFTP,
		RuntimeTypeLLM,
	} {
		require.True(t, IsRuntimeType(typ), "IsRuntimeType(%q)", typ)
	}
	require.False(t, IsRuntimeType("ldap"))
	require.False(t, IsRuntimeType("unknown"))
}

func TestReadSetOperationExtension(t *testing.T) {
	required := true
	status := 201
	var extensions map[string]any
	err := SetOperationExtension(&extensions, &OperationRuntime{
		Type:               RuntimeTypeFnct,
		Function:           "identity",
		PayloadRequired:    &required,
		ResponseStatusCode: &status,
		Arguments:          []any{"input"},
	})
	require.NoError(t, err)
	require.Contains(t, extensions, ExtensionRuntime)

	got, ok, err := ReadOperationExtension(extensions)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, RuntimeTypeFnct, got.Type)
	require.Equal(t, "identity", got.Function)
	require.NotNil(t, got.PayloadRequired)
	require.True(t, *got.PayloadRequired)
	require.NotNil(t, got.ResponseStatusCode)
	require.Equal(t, 201, *got.ResponseStatusCode)
	require.Equal(t, []any{"input"}, got.Arguments)
}

func TestReadConfigExtension(t *testing.T) {
	var extensions map[string]any
	err := SetConfigExtension(&extensions, &ConfigRuntime{
		Provider: &Provider{Name: "api", ServerURL: "https://api.example.test"},
		Security: []*SecurityRequirement{
			{Name: "oauth", Scopes: []string{"read"}, Scheme: &SecurityScheme{Type: "oauth2"}},
		},
	})
	require.NoError(t, err)

	got, ok, err := ReadConfigExtension(extensions)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "api", got.Provider.Name)
	require.Equal(t, "oauth", got.Security[0].Name)

	data, err := json.Marshal(extensions[ExtensionRuntimeConfig])
	require.NoError(t, err)
	require.Contains(t, string(data), "serverUrl")
}
