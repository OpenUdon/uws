# UWS Browser Authentication Call Supplement 1.0

`uws.browser-authentication-call.1.0` identifies an explicit extension-owned
operation that establishes a named browser session by executing one reviewed
`uws.browser-authentication.1.0` flow.

```yaml
operationId: authenticate_member
x-uws-operation-profile: uws.browser-authentication-call.1.0
x-uws-browser-authentication:
  profile: ./browser-authentication/member.yaml
  flow: member_login_push
  session: member_portal
  credentialBindings:
    username: member_username
    password: member_password
timeout: 120
```

`profile` is a safe package-relative path. `flow` selects exactly one named
alternative. `session` is a portable execution-local name. Each
`credentialBindings` key names a profile slot and each value is an opaque
runtime binding identifier—not a credential value.

Authentication operations MUST set a UWS `timeout` from 1 through 600 seconds;
authoring tools SHOULD recommend 120 seconds. They always require explicit
runtime authentication approval.

Protected `uws.browser.1.5` operations select the resulting session with:

```yaml
x-uws-browser-session:
  session: member_portal
```

The named session exists only inside one workflow execution unless the private
driver is explicitly configured with an operator-owned backing binding. No
session state enters UWS. A session-expired result does not cause implicit
reauthentication or action retry; a workflow may invoke another explicit,
approved authentication operation.
