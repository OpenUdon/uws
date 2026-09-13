package schemas

import "fmt"

// Called only after the closed 1.2 schema has validated every member.
func validateRegistrationVerification(root map[string]any, origins map[string]struct{}) error {
	for name, raw := range root["flows"].(map[string]any) {
		flow := raw.(map[string]any)
		verification := flow["humanVerification"].(map[string]any)
		dependencies := verification["dependencies"].(map[string]any)
		if dependencies["policy"] != verification["provider"].(string)+".v1" {
			return fmt.Errorf("flows.%s: verification provider/policy mismatch", name)
		}
		if err := validateAuthenticationTarget(verification["submissionURL"].(string), origins); err != nil {
			return fmt.Errorf("flows.%s: invalid verification submission destination: %w", name, err)
		}
		submitted := false
		for _, rawStep := range flow["sequence"].([]any) {
			step := rawStep.(map[string]any)
			if checkpoint, ok := step["human_checkpoint"].(map[string]any); ok && checkpoint["kind"] == "captcha" {
				return fmt.Errorf("flows.%s: verification readiness replaces CAPTCHA Continue", name)
			}
			if submitted {
				if _, ok := step["wait_for"]; !ok {
					return fmt.Errorf("flows.%s: only observation is allowed after submit", name)
				}
			}
			if _, ok := step["submit"]; ok {
				submitted = true
			}
		}
	}
	return nil
}
