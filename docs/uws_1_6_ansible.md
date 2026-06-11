# UWS 1.6 Ansible Module Source Design

> Design record for the UWS 1.6 adoption of `ansible-module` as a first-class
> `sourceDescriptions[].type`, with the heavier sub-spec published separately as
> `versions/ansible.1.0.{json,md}` — the same two-layer packaging used by
> `runtime.1.0` and `browser.1.5`. The companion roadmap entry lives in
> [`future-source-profiles.md`](future-source-profiles.md)
> ("Adopted in UWS 1.6: Ansible Module Source Profiles"). Sub-spec versions are
> numbered by their own streams.

## 1. Core Philosophy: Module Argspecs Are Published Operation Contracts

Ansible collections publish machine-readable contracts for every module:
`ansible-doc --json` emits the documented parameter specification (name, type,
required, choices, defaults) and return-value documentation for each module.
That is a provider-published operation contract — exactly the artifact class
UWS already binds to for OpenAPI, Google Discovery, AWS Smithy, AsyncAPI,
GraphQL, OpenRPC, gRPC/protobuf, and OData.

UWS 1.6 therefore enforces the same boundary as every earlier source family:

* **The collection owns the module contract** — parameter schemas, return
  shapes, module semantics, idempotent convergence behavior.
* **UWS owns the workflow overlay** — which modules run, in what order, with
  which argument values, gated by which conditions, feeding which outputs.
* **The bound runtime owns the mechanism** — connection plugins, `become`
  privilege escalation, forks/serial strategy, check mode, vault decryption,
  and the Python execution environment.

A UWS document never restates a module's parameter schema, just as it never
restates an OpenAPI operation's request schema. It names the module and binds
values.

Two properties make Ansible modules an unusually good fit for UWS:

1. **Module invocations are JSON-shaped on both sides.** Arguments are a JSON
   object; results are a JSON object (with conventional keys such as `changed`,
   `failed`, `msg`, `rc`). They map directly onto `request.body` and
   `$response.body.*` with no impedance mismatch.
2. **Playbooks are ordered task lists with convergent leaves.** That is
   imperative control flow over idempotent operations — structurally closer to
   the UWS execution model than any prior source family. `when`, `loop`,
   `dependsOn`, and `register`→`outputs` map one-to-one.

## 2. Document Structure: `ansible-module`

UWS 1.6 introduces `ansible-module` as a first-class `sourceDescriptions[].type`.
The referenced source document is a collection argspec document (the
`ansible-doc --json` output shape, or the sub-spec's normalized equivalent).

```yaml
sourceDescriptions:
  - name: builtin
    type: ansible-module
    url: ./ansible/ansible-builtin.argspec.json

operations:
  - operationId: install_nginx
    sourceDescription: builtin
    sourceOperationId: ansible.builtin.apt     # FQCN selects the module
    request:
      body:
        name: nginx
        state: present
    outputs:
      changed: $response.body.changed
```

Selector rules follow the generic pattern introduced in UWS 1.2 and reused by
every family since:

* `sourceOperationId` is the module's fully qualified collection name (FQCN),
  e.g. `ansible.builtin.apt`, `community.postgresql.postgresql_db`.
* `sourceOperationRef`, when used, MUST be a JSON Pointer fragment; the
  preferred target is `#/modules/<fqcn>`.
* Legacy `openapiOperationId` / `openapiOperationRef` selectors are forbidden on
  `ansible-module` sources (consistent with all non-openapi types since 1.2).
* UWS core validates only the type name and selector shape. Whether an FQCN
  resolves to a real module, and whether `request.body` satisfies the module's
  argspec, is the job of source-aware tooling and runtimes — the same deferral
  that kept GraphQL/OData/protobuf support thin in 1.4.

## 3. Worked Example: A Playbook Three Ways

### 3.1 The original playbook

```yaml
- hosts: webservers
  tasks:
    - name: Install nginx
      ansible.builtin.apt:
        name: nginx
        state: present

    - name: Deploy nginx config
      ansible.builtin.template:
        src: nginx.conf.j2
        dest: /etc/nginx/nginx.conf
      notify: restart nginx

  handlers:
    - name: restart nginx
      ansible.builtin.service:
        name: nginx
        state: restarted
```

### 3.2 Valid today (UWS 1.5, degenerate lowering)

A single-target playbook lowers to valid UWS 1.5 right now through the
`uws.runtime.1.0` extension profile, whose type identifiers (`ssh`, `cmd`,
`fileio`, `scp`, `sftp`) already cover Ansible's transport verbs:

```yaml
uws: "1.5.0"
info:
  title: Nginx Setup (degenerate lowering)
  version: 1.0.0

operations:
  - operationId: install_nginx
    x-uws-operation-profile: uws.runtime.1.0
    x-uws-runtime:
      type: ssh
      command: "apt-get install -y nginx"
    successCriteria:
      - condition: $response.statusCode == 0

  - operationId: deploy_config
    x-uws-operation-profile: uws.runtime.1.0
    x-uws-runtime:
      type: scp
      command: "nginx.conf /etc/nginx/nginx.conf"
    dependsOn: [install_nginx]

  - operationId: restart_nginx
    x-uws-operation-profile: uws.runtime.1.0
    x-uws-runtime:
      type: ssh
      command: "systemctl restart nginx"
    dependsOn: [deploy_config]

workflows:
  - workflowId: main
    type: sequence
    steps:
      - stepId: install
        operationRef: install_nginx
      - stepId: deploy
        operationRef: deploy_config
      - stepId: restart
        operationRef: restart_nginx
```

This works, but it is a *degenerate* lowering: it shells out to `apt-get`
instead of calling the convergent `ansible.builtin.apt` module, it always
restarts instead of restarting only on change, and the template is pre-rendered
out of band. It demonstrates the floor, not the target.

### 3.3 UWS 1.6 form

```yaml
uws: "1.6.0"
info:
  title: Nginx Setup
  version: 1.0.0

sourceDescriptions:
  - name: builtin
    type: ansible-module
    url: ./ansible/ansible-builtin.argspec.json

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

  - operationId: deploy_config
    sourceDescription: builtin
    sourceOperationId: ansible.builtin.template
    dependsOn: [install_nginx]
    request:
      body:
        src: nginx.conf.j2
        dest: /etc/nginx/nginx.conf
    outputs:
      changed: $response.body.changed

  - operationId: restart_nginx
    sourceDescription: builtin
    sourceOperationId: ansible.builtin.service
    dependsOn: [deploy_config]
    request:
      body:
        name: nginx
        state: restarted

workflows:
  - workflowId: main
    type: sequence
    steps:
      - stepId: install
        operationRef: install_nginx
      - stepId: deploy
        operationRef: deploy_config
      - stepId: restart
        operationRef: restart_nginx
        when: $steps.deploy.outputs.changed == true
```

The handler becomes an ordinary step gated on the notifying task's `changed`
output. Where several tasks notify the same handler, the `ansible.1.0` sub-spec
defines the lowering rule: a `switch` step gates each case on one notifier's
`changed` output, and each case references the same handler operation. Because
`switch` executes at most one matching case, this preserves Ansible's
run-once-at-end semantics without adding logical OR to UWS core. The Jinja2
template (`nginx.conf.j2`) stays a file the `template` module consumes; UWS
never parses it.

## 4. Inventory and Execution Targets

UWS has no host axis; an Ansible play's execution matrix is tasks × hosts. The
adopted design takes the two-stage path browser profiles took:

**Stage 1 — no core change.** Inventory is input data. A `loop` construct with
`items: $inputs.hosts` fans tasks out per host; connection details and
credentials stay runtime-private (resolved per host by the bound runtime, in
the same way credential bindings already resolve). This works today but treats
hosts as untyped strings.

**Stage 2 — `inventory.1.0` sub-spec.** Publish
`versions/inventory.1.0.{json,md}` alongside `runtime.1.0` / `browser.1.5`: a
typed inventory document (hosts, groups, group/host variables — **never**
secrets, connection passwords, or vault material) plus a step/operation-level
execution-target extension (`x-uws-target: webservers`). Graduation of the
target field into core happens only if multiple runtimes need interoperable
inventory interchange — the same graduation bar every other profile met.

## 5. Convergence Results and Handlers

Ansible's result model is three-state (`ok` / `changed` / `failed`, plus
`skipped`), and handler semantics depend on `changed`. UWS core stays binary
(`successCriteria`); the `ansible.1.0` sub-spec layers the convention on top:

* Module operations SHOULD expose `changed: $response.body.changed` in
  `outputs`. Module results are JSON; `changed` is already there.
* `failed` maps to UWS failure via
  `successCriteria: [{condition: $response.body.failed != true}]` or the
  runtime's module-failure signaling.
* `skipped` is represented structurally: a step whose `when` gate is false does
  not execute — same semantics, no new state needed.
* Handlers lower to `when`-gated steps (single notifier) or switch-gated steps
  (multiple notifiers), as in §3.3.

A possible future core convenience — an `onChanged` action list parallel to
`onSuccess`/`onFailure` — is explicitly **not** part of UWS 1.6. The
convention covers the need without widening core.

## 6. Safety & Integrity Controls

* **No secrets in the wire format.** Vault-encrypted values, connection
  passwords, and become passwords never appear in a UWS document or inventory
  document. Credential bindings stay symbolic; the runtime resolves them. This
  is the existing UWS posture applied unchanged.
* **Side-effect declaration.** Module operations that mutate hosts (`state:
  absent`, `service: restarted`, user/file deletion) SHOULD carry the same
  side-effect/confirmation metadata conventions used elsewhere in UWS packages;
  check-mode execution (dry run) is a runtime feature, not a wire-format field.
* **Fail closed on contract drift.** If an FQCN does not resolve in the
  referenced argspec document, or `request.body` violates the argspec,
  source-aware tooling rejects the document before execution rather than
  guessing — consistent with every other source family.

## 7. Non-Goals

UWS 1.6 deliberately does NOT standardize:

* **Jinja2.** Ansible's templating language never enters the UWS expression
  grammar. Authoring tooling lowers templates: static values render inline,
  dynamic references become UWS expressions, complex renders stay inside module
  arguments (e.g. the `template` module's `src` file) where the module owns
  evaluation. Precedent: `browser.1.5` keeps `{{param}}` templating
  sub-spec-owned.
* **Connection and execution policy.** Connection plugins, `become`,
  forks/serial/strategy, async/poll, and check mode are runtime-owned, exactly
  as HTTP transport and retries are for API sources.
* **Roles and collection packaging.** Galaxy packaging, role directory layout,
  and dependency resolution stay in the Ansible ecosystem. UWS binds to module
  contracts, not to distribution mechanics.
* **The Ansible execution engine.** UWS does not re-specify task executor
  internals, callback plugins, or strategy plugins. A UWS runtime MAY execute
  modules through `ansible-core`, through a compatible reimplementation, or by
  any mechanism that honors the module contract.

## 8. Graduation Criteria Check

Against the source-family baseline `future-source-profiles.md` applies
(the UWS 1.4 standard):

| Criterion | Status |
|---|---|
| Clear operation selector model | FQCN as `sourceOperationId`; `#/modules/<fqcn>` as `sourceOperationRef`. |
| Source-owned request/response metadata | Collection argspecs own parameter schemas and return docs; UWS never restates them. |
| Parser/index evidence or bounded tooling path | `ansible-doc --json` output is stable, machine-readable, and already shipped by every collection; an apitools-style argspec parser is a bounded effort. |
| Runtime path that consumes source semantics without flattening | Modules execute as modules (convergent, idempotent), not flattened into shell strings — the degenerate `runtime.1.0` lowering remains available as the explicit floor. |
| No secrets / credentials in the wire format | §6; vault and connection material stay runtime-private. |
| Multi-runtime interchange need | Playbook-style automation spans many runtimes (ansible-core, AWX/Controller, custom executors); a neutral binding lets UWS-authoring tools target all of them. |

The adopted implementation follows the established per-release recipe: schema
enum + selector-description enrichment, Go constant + validation switch +
`supports16` version gate, parametrized binding tests, and the
conformance/convert round-trip extensions.
