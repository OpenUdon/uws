package uws1

import "testing"

func TestUWS15AcceptsExtensionOwnedAnsibleModuleCall(t *testing.T) {
	doc := &Document{
		UWS:  "1.5.0",
		Info: &Info{Title: "ansible compat", Version: "1.0.0"},
		Operations: []*Operation{{
			OperationID: "install_nginx",
			Request: map[string]any{"body": map[string]any{
				"name":  "nginx",
				"state": "present",
			}},
			Outputs: map[string]string{"changed": "$response.body.changed"},
			SuccessCriteria: []*Criterion{{
				Condition: "$response.body.failed != true",
			}},
			Extensions: map[string]any{
				ExtensionOperationProfile: "uws.ansible-module-call.1.0",
				"x-uws-ansible-module": map[string]any{
					"module": "ansible.builtin.apt",
					"argspec": map[string]any{
						"sourceId":   "builtin",
						"url":        "./ansible-builtin.argspec.json",
						"collection": "ansible.builtin",
					},
				},
			},
		}},
		Workflows: []*Workflow{{
			WorkflowID: "main",
			Type:       WorkflowTypeSequence,
			Steps: []*Step{{
				StepID:       "install",
				OperationRef: "install_nginx",
			}},
		}},
	}

	if err := doc.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}
