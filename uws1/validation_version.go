package uws1

import (
	"fmt"
	"strconv"
	"strings"
)

func (d *Document) validateVersionedFields(result *ValidationResult) {
	supports11 := supportsUWSVersion(d.UWS, 1, 1)
	for i, op := range d.Operations {
		if op == nil {
			continue
		}
		validateVersionedTimeout(op.Timeout, fmt.Sprintf("operations[%d].timeout", i), supports11, result)
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
	if !versionPattern.MatchString(version) {
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
