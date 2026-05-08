package runtimes

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicRuntimeStructFieldsHaveHCLTags(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeOf(OperationRuntime{}),
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
		RuntimeTypeSSH, RuntimeTypeCmd, RuntimeTypeFnct,
		RuntimeTypeFileIO, RuntimeTypeSQL, RuntimeTypeS3, RuntimeTypeSMTP,
		RuntimeTypeDNS, RuntimeTypeLDAPS, RuntimeTypeSCP, RuntimeTypeSFTP,
		RuntimeTypeLLM,
	} {
		require.True(t, IsRuntimeType(typ), "IsRuntimeType(%q)", typ)
	}
	require.False(t, IsRuntimeType("http"))
	require.False(t, IsRuntimeType("ldap"))
	require.False(t, IsRuntimeType("unknown"))
}

func TestReadSetOperationExtension(t *testing.T) {
	var extensions map[string]any
	err := SetOperationExtension(&extensions, &OperationRuntime{
		Type:       RuntimeTypeCmd,
		Command:    "echo ok",
		WorkingDir: "/tmp",
		Function:   "identity",
		Workflow:   "child",
		Arguments:  []any{"input"},
	})
	require.NoError(t, err)
	require.Contains(t, extensions, ExtensionRuntime)

	got, ok, err := ReadOperationExtension(extensions)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, RuntimeTypeCmd, got.Type)
	require.Equal(t, "echo ok", got.Command)
	require.Equal(t, "/tmp", got.WorkingDir)
	require.Equal(t, "identity", got.Function)
	require.Equal(t, "child", got.Workflow)
	require.Equal(t, []any{"input"}, got.Arguments)
}

func TestReadOperationExtensionDecodesSlimPayload(t *testing.T) {
	extensions := map[string]any{
		ExtensionRuntime: map[string]any{
			"type":       RuntimeTypeLLM,
			"function":   "summarize",
			"workflow":   "review",
			"arguments":  []any{map[string]any{"prompt": "$inputs.text"}},
			"workingDir": "/tmp/work",
		},
	}

	got, ok, err := ReadOperationExtension(extensions)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, RuntimeTypeLLM, got.Type)
	require.Equal(t, "summarize", got.Function)
	require.Equal(t, "review", got.Workflow)
	require.Equal(t, "/tmp/work", got.WorkingDir)
	require.Equal(t, []any{map[string]any{"prompt": "$inputs.text"}}, got.Arguments)
}
