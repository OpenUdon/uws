package ansiblemodulecall

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicStructFieldsHaveHCLTags(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeOf(OperationAnsibleModule{}),
		reflect.TypeOf(ArgspecReference{}),
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

func TestReadSetOperationExtension(t *testing.T) {
	var extensions map[string]any
	err := SetOperationExtension(&extensions, &OperationAnsibleModule{
		Module: "ansible.builtin.apt",
		Argspec: &ArgspecReference{
			SourceID:   "builtin",
			URL:        "./ansible-builtin.argspec.json",
			Collection: "ansible.builtin",
		},
	})
	require.NoError(t, err)
	require.Contains(t, extensions, ExtensionAnsibleModule)

	got, ok, err := ReadOperationExtension(extensions)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "ansible.builtin.apt", got.Module)
	require.NotNil(t, got.Argspec)
	require.Equal(t, "builtin", got.Argspec.SourceID)
	require.Equal(t, "./ansible-builtin.argspec.json", got.Argspec.URL)
	require.Equal(t, "ansible.builtin", got.Argspec.Collection)
}

func TestReadOperationExtensionDecodesGenericPayload(t *testing.T) {
	extensions := map[string]any{
		ExtensionAnsibleModule: map[string]any{
			"module": "ansible.builtin.service",
			"argspec": map[string]any{
				"sourceId":   "builtin",
				"url":        "argspec.json",
				"collection": "ansible.builtin",
			},
		},
	}

	got, ok, err := ReadOperationExtension(extensions)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "ansible.builtin.service", got.Module)
	require.Equal(t, "builtin", got.Argspec.SourceID)
}
