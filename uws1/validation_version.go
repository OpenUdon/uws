package uws1

import (
	"fmt"
	"strconv"
	"strings"
)

func (d *Document) validateVersionedFields(result *ValidationResult) {
	supports11 := supportsUWSVersion(d.UWS, 1, 1)
	supports12 := supportsUWSVersion(d.UWS, 1, 2)
	supports13 := supportsUWSVersion(d.UWS, 1, 3)
	supports14 := supportsUWSVersion(d.UWS, 1, 4)
	supports15 := supportsUWSVersion(d.UWS, 1, 5)
	supports16 := supportsUWSVersion(d.UWS, 1, 6)
	supports17 := supportsUWSVersion(d.UWS, 1, 7)
	for i, sd := range d.SourceDescriptions {
		if sd == nil {
			continue
		}
		switch sd.EffectiveType() {
		case SourceDescriptionTypeGoogleDiscovery, SourceDescriptionTypeAWSSmithy:
			if !supports12 {
				result.addError(fmt.Sprintf("sourceDescriptions[%d].type", i), "requires UWS 1.2.0 or later")
			}
		case SourceDescriptionTypeAsyncAPI:
			if !supports13 {
				result.addError(fmt.Sprintf("sourceDescriptions[%d].type", i), "requires UWS 1.3.0 or later")
			}
		case SourceDescriptionTypeGraphQL, SourceDescriptionTypeOpenRPC, SourceDescriptionTypeGRPCProtobuf, SourceDescriptionTypeOData:
			if !supports14 {
				result.addError(fmt.Sprintf("sourceDescriptions[%d].type", i), "requires UWS 1.4.0 or later")
			}
		case SourceDescriptionTypeBrowserProfile:
			if !supports15 {
				result.addError(fmt.Sprintf("sourceDescriptions[%d].type", i), "requires UWS 1.5.0 or later")
			}
		case SourceDescriptionTypeAnsibleModule:
			// ansible-module was introduced in 1.6.0 and withdrawn in 1.7.0.
			// Source descriptions identify pre-existing named operations on a
			// remote target; the control node instead supplies Ansible module
			// implementations, so calls use the operation profile below.
			switch {
			case !supports16:
				result.addError(fmt.Sprintf("sourceDescriptions[%d].type", i), "requires UWS 1.6.0 or later")
			case supports17:
				result.addError(fmt.Sprintf("sourceDescriptions[%d].type", i), "removed in UWS 1.7.0; use the uws.ansible-module-call.1.0 operation profile")
			}
		}
	}
	for i, op := range d.Operations {
		if op == nil {
			continue
		}
		validateVersionedTimeout(op.Timeout, fmt.Sprintf("operations[%d].timeout", i), supports11, result)
		if (op.SourceOperationID != "" || op.SourceOperationRef != "") && !supports12 {
			result.addError(fmt.Sprintf("operations[%d]", i), "sourceOperationId and sourceOperationRef require UWS 1.2.0 or later")
		}
	}
	for i, wf := range d.Workflows {
		if wf == nil {
			continue
		}
		workflowPath := fmt.Sprintf("workflows[%d]", i)
		validateVersionedTimeout(wf.Timeout, workflowPath+".timeout", supports11, result)
		validateVersionedIdempotency(wf.Idempotency, workflowPath+".idempotency", supports11, result)
		validateVersionedStepFields(wf.Steps, workflowPath+".steps", supports11, result)
		validateVersionedCaseStepFields(wf.Cases, workflowPath+".cases", supports11, result)
		validateVersionedStepFields(wf.Default, workflowPath+".default", supports11, result)
	}
}

func validateVersionedStepFields(steps []*Step, path string, supports11 bool, result *ValidationResult) {
	for i, step := range steps {
		if step == nil {
			continue
		}
		stepPath := fmt.Sprintf("%s[%d]", path, i)
		validateVersionedTimeout(step.Timeout, stepPath+".timeout", supports11, result)
		validateVersionedStepFields(step.Steps, stepPath+".steps", supports11, result)
		validateVersionedCaseStepFields(step.Cases, stepPath+".cases", supports11, result)
		validateVersionedStepFields(step.Default, stepPath+".default", supports11, result)
	}
}

func validateVersionedCaseStepFields(cases []*Case, path string, supports11 bool, result *ValidationResult) {
	for i, c := range cases {
		if c == nil {
			continue
		}
		validateVersionedStepFields(c.Steps, fmt.Sprintf("%s[%d].steps", path, i), supports11, result)
	}
}

func validateVersionedTimeout(timeout *float64, path string, supports11 bool, result *ValidationResult) {
	if timeout == nil {
		return
	}
	if !supports11 {
		result.addError(path, "requires UWS 1.1.0 or later")
		return
	}
	if *timeout <= 0 {
		result.addError(path, "must be greater than 0")
	}
}

func validateVersionedIdempotency(idempotency *Idempotency, path string, supports11 bool, result *ValidationResult) {
	if idempotency == nil {
		return
	}
	if !supports11 {
		result.addError(path, "requires UWS 1.1.0 or later")
		return
	}
	if strings.TrimSpace(idempotency.Key) == "" {
		result.addError(path+".key", "is required")
	}
	switch idempotency.OnConflict {
	case "", "reject", "returnPrevious":
	default:
		result.addError(path+".onConflict", fmt.Sprintf("%q is not valid (must be reject or returnPrevious)", idempotency.OnConflict))
	}
	if idempotency.TTL != nil && *idempotency.TTL <= 0 {
		result.addError(path+".ttl", "must be greater than 0")
	}
}

func supportsUWSVersion(version string, major, minor int) bool {
	if !uws1VersionPattern.MatchString(version) {
		return false
	}
	base := version
	if idx := strings.Index(base, "-"); idx >= 0 {
		base = base[:idx]
	}
	parts := strings.Split(base, ".")
	if len(parts) < 2 {
		return false
	}
	gotMajor, err := strconv.Atoi(parts[0])
	if err != nil || gotMajor != major {
		return false
	}
	gotMinor, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	return gotMinor >= minor
}
