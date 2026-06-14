# UWS Ansible Module Source Profile 1.0

The UWS Ansible Module Source Profile is the sub-spec for documents referenced
by UWS 1.6 `sourceDescriptions[].type: ansible-module`. It is inert source
metadata for conversion and review. It does not standardize the Ansible
execution engine, inventory connection behavior, module invocation, connection
plugins, privilege escalation (`become`), forks/serial strategy, check mode,
vault decryption, callback or strategy plugins, or the Python execution
environment. Jinja2 templating is owned by authoring tooling and by modules
that consume template files; it never enters the UWS expression grammar.

UWS core (the main `versions/1.6.0.json` schema) only references this profile
by *type name* and reuses the existing generic `sourceOperationId` /
`sourceOperationRef` selector rules. Validating an argspec document against the
shape below — and validating a UWS operation's `request.body` against a
module's parameter specification — is the job of source-aware conversion and
review tooling.

This profile does not define Ansible playbook control flow. Tooling may lower
playbook constructs into existing UWS workflow objects, but those rules are
converter-owned. For example, a converter can map `when` to `when` or `switch`,
`loop` to `forEach`, `register` field references to operation `outputs`,
`failed_when` to `successCriteria`, and `until`/`retries` to `onFailure:
retry`. The emitted UWS document carries the result of that lowering; the
argspec profile remains only the module leaf contract.

Unlike `browser-profile`, an `ansible-module` source is a **provider-published
contract**: collections ship documented argument specifications
(`ansible-doc --json`) for every module. This profile defines the normalized
document shape those argspecs are carried in.

For UWS 1.5-compatible documents, tooling can still preserve Ansible module
identity by emitting extension-owned operations with
`x-uws-operation-profile: uws.ansible-module-call.1.0` and
`x-uws-ansible-module.module`. That compatibility supplement references the
same argspec documents for review, but it is not a source binding and does not
change the UWS 1.6 `ansible-module` source contract.

## Profile

| Field | Value |
| --- | --- |
| Profile name | `uws.ansible.1.0` |
| JSON Schema | `versions/ansible.1.0.json` |
| Bound from | UWS 1.6+ `sourceDescriptions[].type: ansible-module` |
| Selector | `sourceOperationId` (module FQCN) or `sourceOperationRef` (preferred form: `#/modules/<fqcn>`) |

A minimal argspec document:

```yaml
argspec: uws.ansible.1.0
collection: ansible.builtin
info:
  collectionVersion: "2.17.0"
  source: ansible-doc --json
modules:
  ansible.builtin.apt:
    shortDescription: Manages apt packages
    parameters:
      name:
        type: list
        elements: str
        description: Package name or list of names.
      state:
        type: str
        choices: [absent, build-dep, latest, present, fixed]
        default: present
    returns:
      stdout:
        type: str
        returned: success

  ansible.builtin.service:
    shortDescription: Manage services
    parameters:
      name:
        type: str
        required: true
      state:
        type: str
        choices: [reloaded, restarted, started, stopped]
```

## Document Fields

| Field | Type | Purpose |
| --- | --- | --- |
| `argspec` | string const | REQUIRED. Literal `uws.ansible.1.0`. |
| `collection` | string | REQUIRED. Collection namespace (`ansible.builtin`, `community.postgresql`). |
| `info` | object | Optional provenance: `title`, `collectionVersion`, `extractedAt`, `source`. |
| `modules` | object | REQUIRED. Module entries keyed by FQCN. Each key is the value a UWS operation uses as `sourceOperationId`. |

## Module Entries

| Field | Type | Purpose |
| --- | --- | --- |
| `shortDescription` | string | One-line module summary. |
| `description` | string or string array | Longer documentation. |
| `parameters` | object | REQUIRED. Parameters keyed by name; `request.body` values bind against these names. |
| `returns` | object | Documented return values; UWS `outputs` reference them via `$response.body.<name>`. The conventional `changed`, `failed`, and `msg` keys are implied for every module and need not be declared. |

Parameter fields: `type` (one of `str`, `bool`, `int`, `float`, `list`, `dict`,
`path`, `raw`, `jsonarg`, `json`, `bytes`, `sig`), `required`, `default`,
`choices`, `elements` (for lists), `aliases`, `description`, and `noLog`.

**`noLog` parameters are sensitive.** A UWS document MUST NOT carry literal
values for `noLog` parameters; the bound value MUST be a symbolic credential
binding the runtime resolves at execution time. Source-aware validation MUST
fail closed on a literal `noLog` value.

## Request Binding

Module arguments bind under `request.body` in the UWS operation, keyed by
parameter name (aliases permitted; tooling SHOULD normalize to the canonical
name):

```yaml
operations:
  - operationId: install_nginx
    sourceDescription: builtin
    sourceOperationId: ansible.builtin.apt
    request:
      body:
        name: nginx
        state: present
    outputs:
      changed: $response.body.changed
```

Source-aware tooling validates `request.body` against the module's
`parameters`: unknown parameter names, missing `required` parameters, and
out-of-`choices` literals are validation failures. UWS runtime expressions
(`$inputs.*`, `$steps.*`, iteration scope) are opaque to argspec validation and
are checked only for placement, not value.

## Result Conventions

Ansible modules return JSON with conventional keys. This profile standardizes
how UWS workflows consume them:

| Ansible state | UWS representation |
| --- | --- |
| `ok` (no change) | Successful execution; `$response.body.changed == false`. |
| `changed` | Successful execution; `$response.body.changed == true`. Module operations SHOULD expose `changed: $response.body.changed` in `outputs`. |
| `failed` | Operation failure: `successCriteria: [{condition: $response.body.failed != true}]`, or the runtime's module-failure signaling. |
| `skipped` | Structural: a step whose `when` gate is false does not execute. No additional state is needed. |

## Handler Lowering

Ansible `notify`/handler semantics lower to ordinary UWS steps gated on
`changed` outputs:

**Single notifier** — the handler step is `when`-gated directly:

```yaml
steps:
  - stepId: deploy
    operationRef: deploy_config
  - stepId: restart
    operationRef: restart_nginx
    when: $steps.deploy.outputs.changed == true
```

**Multiple notifiers** — UWS core expressions are single binary comparisons
(no logical OR), so the lowering uses a `switch` step placed after the
notifying steps. Each case gates on one notifier's `changed` output, and every
case routes to a step referencing the *same* handler operation. Because a
`switch` executes at most one matching case, this preserves Ansible's
run-once-at-end, deduplicated handler semantics: the handler executes at most
once per workflow run, after all notifying steps have completed, and only if
at least one reported a change.

```yaml
steps:
  - stepId: install
    operationRef: install_nginx
  - stepId: deploy
    operationRef: deploy_config
  - stepId: restart_nginx_notify
    type: switch
    cases:
      - name: notified_by_install
        when: $steps.install.outputs.changed == true
        steps:
          - stepId: restart_nginx_run_1
            operationRef: restart_nginx
      - name: notified_by_deploy
        when: $steps.deploy.outputs.changed == true
        steps:
          - stepId: restart_nginx_run_2
            operationRef: restart_nginx
```

Handler lowering is a converter convention, not a new source-profile field.
The same handler orchestration can be emitted in UWS 1.5 compatibility output
using `uws.ansible-module-call.1.0` extension-owned operations or in UWS 1.6
output using first-class `ansible-module` source bindings.

## Inventory Posture (Stage 1)

This profile does not define an inventory document. Host fan-out uses existing
UWS constructs: a `loop` over `$inputs.hosts` (or a single host input), with
connection details and credentials resolved per host by the bound runtime.
A typed inventory document and an execution-target extension are reserved for a
future `inventory.1.0` sub-spec, to be considered when multiple runtimes need
interoperable inventory interchange.

## What Stays Runtime-Owned

The following are NOT part of `ansible.1.0` and MUST NOT appear in argspec or
UWS documents:

- Connection plugins, transports, ports, SSH options, and host key policy.
- `become` / `become_user` / privilege-escalation configuration and passwords.
- Vault-encrypted values, vault passwords, or any decrypted secret material.
- Forks, serial, strategy plugins, async/poll, and check-mode invocation.
- Callback plugins, logging configuration, and execution-environment selection.
- Continue-on-error behavior such as `ignore_errors`. UWS 1.6 has no portable
  continue-on-failure primitive, so converters must fail closed or keep the
  behavior runtime-owned until a future core or profile field exists.
- Jinja2 template evaluation. Authoring tooling lowers templates before the
  document is written: static values render inline, dynamic references become
  UWS expressions, and template *files* (e.g. the `template` module's `src`)
  stay files the module consumes.

Argspec documents are inert contract metadata. This profile does not define how
to execute modules, connect to inventory hosts, bind credentials, or interpret
Ansible runtime behavior.
