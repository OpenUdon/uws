# UWS Ansible Module Call Supplement 1.0

The UWS Ansible Module Call Supplement defines a small public extension profile
for extension-owned Ansible module leaf operations. It is intended for UWS
1.5-compatible documents, where `sourceDescriptions[].type: ansible-module`
does not exist yet.

This supplement is wire/spec metadata only. It does not standardize the
Ansible execution engine, inventory connections, connection plugins,
credentials, vault material, privilege escalation, forks/strategy, check mode,
callback plugins, or the Python execution environment.

## Profile

| Field | Value |
| --- | --- |
| Profile name | `uws.ansible-module-call.1.0` |
| Operation extension | `x-uws-ansible-module` |
| JSON Schema | `versions/ansible-module-call.1.0.json` |

An extension-owned operation using this supplement sets:

```yaml
operationId: install_nginx
x-uws-operation-profile: uws.ansible-module-call.1.0
x-uws-ansible-module:
  module: ansible.builtin.apt
  argspec:
    sourceId: builtin
    url: ./ansible/ansible-builtin.argspec.json
    collection: ansible.builtin
request:
  body:
    name: nginx
    state: present
outputs:
  changed: $response.body.changed
successCriteria:
  - condition: $response.body.failed != true
```

## Operation Payload

`x-uws-ansible-module` is an object with these fields:

| Field | Type | Purpose |
| --- | --- | --- |
| `module` | string | REQUIRED. Ansible module FQCN, such as `ansible.builtin.apt`. |
| `argspec` | object | Optional review reference to the argspec source used by conversion tooling. |

`argspec`, when present, may carry `sourceId`, `url`, and `collection`.

Module arguments stay in normal UWS `request.body`. Module result values stay
under `$response.body.*`. This mirrors the UWS 1.6 `ansible-module` source-bound
form while keeping the operation extension-owned for older UWS documents.

## Boundary

This supplement only identifies the Ansible module leaf operation and its
review-time argspec reference. The bound runtime owns module invocation and all
Ansible execution behavior.

The following are not part of this supplement and MUST remain runtime-private or
conversion diagnostics:

- inventory parsing and host connection behavior
- SSH options, connection plugins, and execution environments
- `become` / privilege escalation and all related secrets
- vault-encrypted or decrypted secret material
- forks, serial, strategy plugins, async/poll, check mode, and callbacks
- arbitrary Jinja2 evaluation

UWS core still owns workflow orchestration around the leaf: `sequence`,
`switch`, `loop`, `await`, `when`, `forEach`, `outputs`, `successCriteria`, and
`onFailure`.

## HCL Representation

In HCL, extension fields are represented inside an `extensions` block:

```hcl
operation "install_nginx" {
  request = {
    body = {
      name  = "nginx"
      state = "present"
    }
  }

  extensions {
    x-uws-operation-profile = "uws.ansible-module-call.1.0"
    x-uws-ansible-module {
      module = "ansible.builtin.apt"
      argspec {
        sourceId   = "builtin"
        url        = "./ansible/ansible-builtin.argspec.json"
        collection = "ansible.builtin"
      }
    }
  }
}
```

JSON and YAML use normal flattened `x-*` fields.
