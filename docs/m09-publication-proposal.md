# M09 source publication proposal

The [qualified exact source](m09-qualification.md) is
`b099f6803277ae94c7e9f1da0904a0140b278f20`, whole review 6/10 passed. Proposed
root and nested codec module version is
`v0.0.0-20261008043726-b099f6803277`. Local exact archives/full selected closures
and a reusable public consumer fixture are prepared. UWS publication authority
has not been granted; M09.5, acceptance and retirement remain incomplete.

Request a fresh **UWS_M09_SOURCE_PUBLICATION** grant limited to normal
fast-forward source/evidence/acceptance-closure publication to:

```text
git@github.com-tabilet:OpenUdon/uws.git
refs/heads/main
```

The concrete initial push head is the evidence-envelope commit containing this
proposal, identified in the owner handoff. Its exact reviewed implementation
ancestor remains b099f6803277ae94c7e9f1da0904a0140b278f20. The local tracking
ref was observed at 0a4597122a7baa7e79e46e79e4ec60dbfffc3720 and is an ancestor
of the proposed history; this is not a fresh remote observation. Resolve the
unchanged configured origin and fresh remote main before acting, and stop for
reconciliation if normal fast-forward publication is no longer available.

After the separately granted normal push, independently observe remote main and
resolve root/codec at the full implementation source through the ordinary
configured Go source. Verify Origin.Hash, codec subdirectory hcl, proposed
versions, archive/GoMod sums and complete selected closures against qualification.
Use an independent ordinary module cache without the owner bootstrap file proxy
or directory replacement. Run the downloaded module-only full suites, vets and
builds and the provided public consumer fixture with both exact modules selected.
The fixture can be copied to public_test.go in a disposable consumer module:

```go
module example.com/uws-m09-ordinary-proof

go 1.25.4

require (
    github.com/OpenUdon/uws v0.0.0-20261008043726-b099f6803277
    github.com/OpenUdon/uws/hcl v0.0.0-20261008043726-b099f6803277
)
```

Only passing independent ordinary proof satisfies M09.5 publication acceptance.
Then record exact publication/consumer evidence, reconcile pending package-local
consumers, consolidate current knowledge and retire M09 through its existing
protocol. Any acceptance-closure head still uses the same normal main-only grant
and is independently observed; do not relabel a preliminary source as accepted.

Every proposed push head uses the retained `[skip ci]` convention documented in
[M08 publication](m08-publication.md). The existing main-push docs workflow invokes
`mkdocs gh-deploy --force` to another ref; suppressing it preserves the requested
main-only scope. Local checks are the evidence; no hosted-CI/docs-deployment
claim, tag, force push, branch deletion or gh-pages action is included.

The confirmed GOAL permits local task commits and no external mutations; consumed
Udon-only and original Stage 11 grants supply no UWS authority. This proposal
requests only the named source/closure publication above. No deployment, host,
provider/API/model/mail action, live ledger, registration opening, installed
consumer migration, browser-pin change or unrelated sibling write is included.
