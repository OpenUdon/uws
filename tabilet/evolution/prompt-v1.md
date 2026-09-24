# Initial Direction v1

## Origin

The user requested memory-bank initialization after approving and completing a
broad archive preflight for the UWS repository. During discovery, the user
accepted the recommended first delivery outcome.

## Delivery Boundary

The boundary is the tracked UWS specification and Go-module distribution:
versioned contracts, Go packages, fixtures, documentation, and repository
verification. Ignored local idea material, generated site output, downstream
runtime repositories, provider infrastructure, and external services are out
of scope.

## Requested Outcome

Produce one coherent current-release documentation surface, then add automated
immutability protection for published Markdown artifacts. The outcome must
surface browser registration 1.2 consistently, preserve established version
selection defaults, repair two stale editorial statements, add missing release
history and navigation, and prevent later undocumented Markdown drift.

## Confirmed Decisions

- Use one maintenance lane and two dependency-ordered milestones.
- Make no JSON Schema, UWS wire, Go API, validation, execution, or embedded
  archive change for documentation coherence.
- Permit only meaning-preserving editorial corrections to published Markdown.
- Protect every `versions/*.md` document except the intentionally mutable
  changelog after those corrections are accepted.
- Use one commit per completed task row in future execution.
- Perform no external mutations.
- Preserve every verified archive as frozen baseline evidence.

## Acceptance Boundary

Acceptance uses the Go tests and race tests, vet, strict MkDocs build, diff
check, focused published-document immutability test, targeted stale-claim
searches, and manual confirmation that public schemas and behavior did not
change.
