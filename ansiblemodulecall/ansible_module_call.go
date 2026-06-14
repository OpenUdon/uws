package ansiblemodulecall

import (
	"encoding/json"
	"fmt"
)

const (
	// ProfileName is the public operation profile marker for this supplement.
	ProfileName = "uws.ansible-module-call.1.0"

	// ExtensionAnsibleModule is the operation-level Ansible module metadata key.
	ExtensionAnsibleModule = "x-uws-ansible-module"
)

// OperationAnsibleModule is the typed payload for x-uws-ansible-module.
type OperationAnsibleModule struct {
	Module  string            `json:"module,omitempty" hcl:"module,optional"`
	Argspec *ArgspecReference `json:"argspec,omitempty" hcl:"argspec,block"`
}

// ArgspecReference identifies the argspec source used for review validation.
type ArgspecReference struct {
	SourceID   string `json:"sourceId,omitempty" hcl:"sourceId,optional"`
	URL        string `json:"url,omitempty" hcl:"url,optional"`
	Collection string `json:"collection,omitempty" hcl:"collection,optional"`
}

// ReadOperationExtension decodes x-uws-ansible-module from an extension map.
func ReadOperationExtension(extensions map[string]any) (*OperationAnsibleModule, bool, error) {
	if len(extensions) == 0 {
		return nil, false, nil
	}
	value, ok := extensions[ExtensionAnsibleModule]
	if !ok {
		return nil, false, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, false, fmt.Errorf("marshal %s extension: %w", ExtensionAnsibleModule, err)
	}
	var out OperationAnsibleModule
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, false, fmt.Errorf("unmarshal %s extension: %w", ExtensionAnsibleModule, err)
	}
	return &out, true, nil
}

// SetOperationExtension encodes x-uws-ansible-module into an extension map.
func SetOperationExtension(dst *map[string]any, value *OperationAnsibleModule) error {
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
	(*dst)[ExtensionAnsibleModule] = generic
	return nil
}
