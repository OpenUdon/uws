package hcl

import (
	"context"
	"errors"
	"testing"
)

func TestCodecRefusesWrongBlockContainerWithoutPanic(t *testing.T) {
	for _, text := range []string{
		`{"info":[{}]}`,
		`{"info":[]}`,
		`{"info":"invalid"}`,
		`{"operations":{"operationId":"one"}}`,
		`{"workflows":[{"workflowId":"one","steps":{"stepId":"one"}}]}`,
	} {
		t.Run(text, func(t *testing.T) {
			defer func() {
				if recover() != nil {
					t.Error("unsupported typed container panicked instead of refusing")
				}
			}()
			view, err := Render(context.Background(), Source{Format: JSON, Bytes: []byte(text)}, Options{Revision: fixtureRevision})
			if !errors.Is(err, ErrCodec) || len(view.HCL) != 0 {
				t.Error("unsupported typed block container was admitted")
			}
		})
	}
}

func TestCodecRefusesExplicitYAMLTags(t *testing.T) {
	for _, text := range []string{
		"variables: !private {value: 1}\n",
		"variables: {value: !private [1, 2]}\n",
		"variables: {value: !!str 1}\n",
		"variables: {!!str value: 1}\n",
		"!private {variables: {value: 1}}\n",
	} {
		view, err := Render(context.Background(), Source{Format: YAML, Bytes: []byte(text)}, Options{Revision: fixtureRevision})
		if !errors.Is(err, ErrCodec) || len(view.HCL) != 0 {
			t.Error("explicit/unsupported YAML tag was silently discarded")
		}
	}
}
