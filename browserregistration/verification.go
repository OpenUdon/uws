package browserregistration

// HumanVerification binds one reviewed provider widget to the existing submit
// locator's form and its exact POST destination. It contains no selectors,
// scripts, provider keys, responses or challenge content.
type HumanVerification struct {
	Provider      string                   `json:"provider" yaml:"provider"`
	Activation    string                   `json:"activation" yaml:"activation"`
	WidgetBinding string                   `json:"widgetBinding" yaml:"widgetBinding"`
	SubmissionURL string                   `json:"submissionURL" yaml:"submissionURL"`
	Dependencies  VerificationDependencies `json:"dependencies" yaml:"dependencies"`
}

// VerificationDependencies selects an immutable trusted adapter policy. These
// limits may be tightened, never widened, by a runtime's consumer authority.
type VerificationDependencies struct {
	Policy           string `json:"policy" yaml:"policy"`
	MaxRequests      int    `json:"maxRequests" yaml:"maxRequests"`
	MaxResponseBytes int    `json:"maxResponseBytes" yaml:"maxResponseBytes"`
	TimeoutMS        int    `json:"timeoutMs" yaml:"timeoutMs"`
}
