package runtimes

import (
	"encoding/json"
	"fmt"
)

const (
	// ProfileName is the public operation profile marker for this supplement.
	ProfileName = "uws.runtime.1.0"

	// ExtensionRuntime is the operation-level runtime metadata extension key.
	ExtensionRuntime = "x-uws-runtime"
)

const (
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
	case RuntimeTypeSSH, RuntimeTypeCmd, RuntimeTypeFnct,
		RuntimeTypeFileIO, RuntimeTypeSQL, RuntimeTypeS3, RuntimeTypeSMTP,
		RuntimeTypeDNS, RuntimeTypeLDAPS, RuntimeTypeSCP, RuntimeTypeSFTP,
		RuntimeTypeLLM:
		return true
	default:
		return false
	}
}

// OperationRuntime is the typed payload for x-uws-runtime.
type OperationRuntime struct {
	Type       string `json:"type,omitempty" hcl:"type,optional"`
	Command    string `json:"command,omitempty" hcl:"command,optional"`
	WorkingDir string `json:"workingDir,omitempty" hcl:"workingDir,optional"`
	Function   string `json:"function,omitempty" hcl:"function,optional"`
	Workflow   string `json:"workflow,omitempty" hcl:"workflow,optional"`
	Arguments  []any  `json:"arguments,omitempty" hcl:"arguments,optional"`
}

// ReadOperationExtension decodes x-uws-runtime from an extension map.
func ReadOperationExtension(extensions map[string]any) (*OperationRuntime, bool, error) {
	if len(extensions) == 0 {
		return nil, false, nil
	}
	value, ok := extensions[ExtensionRuntime]
	if !ok {
		return nil, false, nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return nil, false, fmt.Errorf("marshal %s extension: %w", ExtensionRuntime, err)
	}
	var out OperationRuntime
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, false, fmt.Errorf("unmarshal %s extension: %w", ExtensionRuntime, err)
	}
	return &out, true, nil
}

// SetOperationExtension encodes x-uws-runtime into an extension map.
func SetOperationExtension(dst *map[string]any, value *OperationRuntime) error {
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
	(*dst)[ExtensionRuntime] = generic
	return nil
}
