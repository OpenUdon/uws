# Browser registration profile 1.2

This additive extension retains registration 1.1 typed private inputs and all
one-attempt controls. UWS core and published registration 1.0/1.1 are unchanged.
Every 1.2 flow requires `humanVerification`. Consumers must explicitly support
1.2 before starting a browser; never strip this field to downgrade a recipe.

The descriptor has exactly `provider` (turnstile, recaptcha_v2 or hcaptcha),
`activation` (before_approval or approved_submit), `widgetBinding`
(single_in_submit_form), `submissionURL`, and `dependencies`. The existing
unique submit locator identifies the reviewed form; the adapter must prove
exact POST action/destination and exactly one supported provider widget bound
to that form. Ambiguity, missing widgets and unsupported integrations stop.
No CSS, JavaScript, site keys, secrets, tokens or challenge content are portable.
CAPTCHA Continue checkpoints are forbidden in 1.2; ordinary consent remains
separate. Only wait_for observation may follow the one submit macro.

Dependencies name an immutable adapter policy: turnstile.v1, recaptcha_v2.v1 or
hcaptcha.v1, matching the provider. Required positive integer limits have maxima
of maxRequests=256, maxResponseBytes=33554432 and timeoutMs=120000. Consumer
policy may tighten them. The verification phase is additionally bounded by the
nonrenewable operation deadline. These budgets never authorize an application
mutation. Policies grant HTTPS GET/HEAD resources and verification POSTs only:

| Policy | Domain/path scope |
| --- | --- |
| turnstile.v1 | challenges.cloudflare.com /turnstile/ and /cdn-cgi/challenge-platform/ |
| recaptcha_v2.v1 | www.google.com and www.recaptcha.net /recaptcha/; www.gstatic.com /recaptcha/ resources only |
| hcaptcha.v1 | hcaptcha.com and DNS-label subdomains; provider resources and verification endpoints |

Provider frames must descend from the reviewed application main frame through
approved provider frames; no provider top-level navigation, application frames,
popups, downloads, service workers or persistent channels are authorized.
Enforce every redirect and nested/out-of-process frame. Application navigation
remains reviewed and exact; a provider policy never widens application origins
or the reviewed POST destination. The hCaptcha DNS family is deliberate because
asset subdomains change; suffix-confusion domains and other ports are excluded.

For before_approval, wait until the trusted widget probe reports ready, present
final approval automatically, recheck readiness at the POST boundary, then
release at most one approved application POST. For approved_submit, obtain
approval covering one trigger plus one submission, activate the control once,
allow human challenge handling and admit a possibly delayed application POST
only after current readiness. An early or second POST is blocked. Expiration,
cancellation, invalid approval or timeout stops without click replay, reload,
reset or automatic registration retry. Provider POSTs have a separate count.

Probes return only loading, awaiting_interaction, ready, expired, failed or
unsupported. Response values stay inside the browser. Client readiness means a
currently usable widget response exists; only backend validation establishes
server acceptance. Diagnostic verification uses the same trusted adapter with
no private inputs and zero application mutations, and produces no account claim
or reusable response. Human challenges are never solved by the agent.

Provider references: [Turnstile configuration](https://developers.cloudflare.com/turnstile/get-started/client-side-rendering/widget-configurations/),
[reCAPTCHA v2](https://developers.google.com/recaptcha/docs/display),
[invisible reCAPTCHA](https://developers.google.com/recaptcha/docs/invisible),
[hCaptcha integration and CSP](https://docs.hcaptcha.com/).
Production acceptance remains separate from synthetic/test-key integration;
[Cloudflare documents automation rejection](https://developers.cloudflare.com/turnstile/troubleshooting/testing/).
