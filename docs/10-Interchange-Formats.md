# Feature 10: Interchange Formats

← [Validation](09-Validation.md) | [Home →](index.md)

---

UWS documents are valid JSON, YAML, or canonical HCL. The `convert` package in `github.com/OpenUdon/uws` moves documents between all three formats with round-trip guarantees.

## Three Formats, One Document

| Format | Best for | Extensions preserved |
|--------|----------|---------------------|
| JSON | Machine interchange, API responses, LLM output | ✓ |
| YAML | Human authoring, configuration files | ✓ |
| HCL | Canonical authoring for runtime tooling | ✓ via `extensions { ... }` blocks |

## Example 1: The Same Operation in All Three Formats

The following is semantically identical in JSON, YAML, and HCL:

**JSON**
```json
{
  "operationId": "list_pets",
  "sourceDescription": "petstore_api",
  "openapiOperationId": "listPets",
  "request": {
    "query": { "limit": 10, "status": "available" },
    "header": { "X-Trace-Id": "abc-123" }
  },
  "outputs": {
    "firstPet": "$response.body#/0",
    "total": "$response.body.total"
  }
}
```

**YAML**
```yaml
operationId: list_pets
sourceDescription: petstore_api
openapiOperationId: listPets
request:
  query:
    limit: 10
    status: available
  header:
    X-Trace-Id: abc-123
outputs:
  firstPet: $response.body#/0
  total: $response.body.total
```

**HCL**
```hcl
operation "list_pets" {
  sourceDescription  = "petstore_api"
  openapiOperationId = "listPets"

  request = {
    query  = { limit = 10, status = "available" }
    header = { "X-Trace-Id" = "abc-123" }
  }

  outputs = {
    firstPet = "$response.body#/0"
    total    = "$response.body.total"
  }
}
```

## Example 2: A Full Document in All Three Formats

A complete minimal document with a workflow:

**YAML** (most readable for authoring)
```yaml
uws: "1.1.0"
info:
  title: Pet Workflow
  version: 1.0.0
sourceDescriptions:
  - name: petstore_api
    url: ./petstore.yaml
    type: openapi
operations:
  - operationId: list_pets
    sourceDescription: petstore_api
    openapiOperationId: listPets
    request:
      query:
        limit: 5
    outputs:
      first: $response.body#/0/id
workflows:
  - workflowId: main
    type: sequence
    steps:
      - stepId: fetch
        operationRef: list_pets
```

**HCL** (canonical authoring form)
```hcl
uws  = "1.1.0"

info {
  title   = "Pet Workflow"
  version = "1.0.0"
}

sourceDescription "petstore_api" {
  url  = "./petstore.yaml"
  type = "openapi"
}

operation "list_pets" {
  sourceDescription  = "petstore_api"
  openapiOperationId = "listPets"
  request = { query = { limit = 5 } }
  outputs = { first = "$response.body#/0/id" }
}

workflow "main" {
  type = "sequence"
  step "fetch" {
    operationRef = "list_pets"
  }
}
```

## Example 3: Switch Case in HCL

Structural constructs translate naturally:

**YAML**
```yaml
workflows:
  - workflowId: route
    type: switch
    cases:
      - name: premium
        when: $outputs.tier == "premium"
        steps:
          - stepId: fast_track
            operationRef: express_process
    default:
      - stepId: standard
        operationRef: normal_process
```

**HCL**
```hcl
workflow "route" {
  type = "switch"

  case "premium" {
    when = "$outputs.tier == \"premium\""
    step "fast_track" {
      operationRef = "express_process"
    }
  }

  default {
    step "standard" {
      operationRef = "normal_process"
    }
  }
}
```

## The `convert` Package

All conversion helpers live in `github.com/OpenUdon/uws/convert`:

```go
// Between byte slices (raw format conversion)
jsonOut, _ := convert.YAMLToJSON(yamlData)
yamlOut, _ := convert.JSONToYAML(jsonData)
hclOut,  _ := convert.JSONToHCL(jsonData)
jsonOut, _ = convert.HCLToJSON(hclData)

// Marshal a Document struct to bytes
jsonBytes, _ := convert.MarshalJSON(doc)
yamlBytes, _ := convert.MarshalYAML(doc)
hclBytes,  _ := convert.MarshalHCL(doc)

// Unmarshal bytes into a Document struct
convert.UnmarshalJSON(jsonData, &doc)
convert.UnmarshalYAML(yamlData, &doc)
convert.UnmarshalHCL(hclData, &doc)
```

## Example 4: Format Conversion in a Go Program

Read a YAML workflow, validate it, and write it back as JSON for machine consumption:

```go
package main

import (
    "log"
    "os"

    "github.com/OpenUdon/uws/convert"
    "github.com/OpenUdon/uws/uws1"
)

func main() {
    yamlData, _ := os.ReadFile("workflow.uws.yaml")

    var doc uws1.Document
    if err := convert.UnmarshalYAML(yamlData, &doc); err != nil {
        log.Fatal(err)
    }
    if err := doc.Validate(); err != nil {
        log.Fatal(err)
    }

    jsonData, err := convert.MarshalJSONIndent(&doc, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    os.WriteFile("workflow.uws.json", jsonData, 0644)
}
```

## HCL Extensions

HCL represents object-level `x-*` extension fields inside an `extensions` block:

```hcl
operation "render" {
  extensions {
    x-uws-operation-profile = "uws.runtime.1.0"
    x-uws-runtime {
      type     = "fnct"
      function = "identity"
    }
  }
}
```

When converted to JSON or YAML, these fields flatten back to normal `x-*`
properties on the owning object.

The same rule applies to the public runtime supplement. A JSON/YAML
`x-uws-runtime` object becomes an HCL block inside `extensions`, and conversion
back to JSON/YAML restores the flattened extension field. The runtime supplement
payload itself still follows its schema: `type` is required, and HTTP/OpenAPI
calls are represented by core OpenAPI binding fields rather than
`x-uws-runtime`.

## `$`-Key Rewriting for HCL

JSON Schema keys like `$ref`, `$id`, `$defs` are not valid HCL identifiers. The package rewrites them symmetrically in both directions:

| JSON / YAML key | HCL key |
|-----------------|---------|
| `$ref` | `_ref` |
| `$id` | `_id` |
| `$schema` | `_schema` |
| `$defs` | `_defs` |
| `$customKey` | `__dollar__customKey` |

**A `ParamSchema` with `$ref` in YAML and HCL:**

```yaml
# YAML
inputs:
  _ref: "#/components/schemas/OrderInput"
```

```hcl
# HCL — $ref becomes _ref
inputs = { _ref = "#/components/schemas/OrderInput" }
```

Round-tripping through HCL → JSON restores `$ref` exactly.

## Round-Trip Guarantee

For any core-only (extension-free) UWS document:

```
JSON → HCL → JSON  produces a structurally identical document
YAML → HCL → YAML  produces a structurally identical document
```

`MarshalHCL` works on a deep copy — the caller's document is never mutated during conversion.

## From The Big Fixture

The large fixture keeps JSON and HCL versions of the same generated document.
This compact JSON excerpt corresponds to the HCL operation examples on the
feature pages:

```json
{
  "operationId": "fetch_ticket",
  "sourceDescription": "incident_api",
  "openapiOperationId": "getIncident",
  "request": {
    "path": {
      "incidentId": "$inputs.incidentId",
      "tenantId": "$inputs.tenantId"
    },
    "query": {
      "depth": "full",
      "include": ["timeline", "assets"]
    }
  },
  "outputs": {
    "severity": "$response.body.severity",
    "ticket": "$response.body"
  }
}
```

Full context: [`testdata/big/big.json`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.json) and [`testdata/big/big.hcl`](https://github.com/OpenUdon/uws/blob/main/testdata/big/big.hcl).

---

← [Validation](09-Validation.md) | [Home →](index.md)
