# C12 — Registration governance and capability guidance

**Stage/owner:** STG-12; UWS. **State:** C12.1 complete; whole-milestone review1 in progress.
**Approval:** The user selected “Include UWS docs milestone” and “Keep review only”,
then approved “Reconcile human-only registration and narrow M53” on 2026-10-17.
**Specification:** [Owning milestone](milestone.md#c12--registration-governance-and-capability-guidance).

## Provenance and boundaries

BRP human-only registration / Stage12 scope reconciliation, 2026-10-17.
Sources: `brp-uws.md`, `brp-kinet-stage12.md` and conversational assessment;
source priority not supplied. Source baselinea66a6b2/full revalidation
`a66a6b279619eab5e06e1bd7f2faedc15f00597c`, clean before planning writes.
F04 confirmed Lower: published contracts are immutable; F05 partially confirmed
Lower: clarify capability limits without redefining execution contracts. This
is review intake, not closing review1. No Kinet or UWS implementation is delivered.

C12 depends only on accepted frozen contracts. Kinet's source refusal belongs to
M53.19; final artifacts/physical proof remain M53.4. C11's conditional trigger
is absent. No module adoption or publication edge is introduced. Coordinator
alone writes integrated ledgers and shared memory; children receive only narrowed
ordinary documentation scope and frozen inputs. Existing protocol/model/grants,
commit/integration and scheduling policies govern execution. Audit disabled.

## Tasks

| Task | State | Notes |
|---|---|---|
| C12.1 — Clarify registration governance and capability limits | `[+]` | One coherent docs change across README.md, docs/registration-authoring.md and versions/CHANGELOG.md. Explain inert review/governance, human account creation, no capability guarantee from schema validity/humanVerification, and unchanged published obligations for optional runtimes. Keep all schema/specification/code/archive/default/pin bytes unchanged. Verify links/wording/whitespace and existing published-document immutability offline; complete required closure checks once, normal task review and fresh whole review before local acceptance. No new schema, browser or publication. |

## Acceptance and verification

All requirements and exact write boundaries are in the owning specification.
Registration review stays distinct from execution. Do not imply that Kinet's
pending M53.19 guard already exists. No universal automation-impossibility claim.
Run the existing offline schemas immutability test with -count=1/-timeout=2m,
required root closure checks once, and git diff --check. No new tests or physical
fixtures are needed for this documentation-only milestone.

## Review

Task review1 PASS by coordinator, independent of the assigned documentation author: exact three-file diff satisfies the scope, preserves optional runtime obligations and makes no Kinet implementation claim. Whole-milestone review1/10 STARTED; fresh Sol xhigh reviewer required.
No expensive/protected invocation consumed. No implementation, commit,
publication or deployment is performed by reconciliation.

## Current execution

The user resumed the existing goal with “continue” after approving R67 and the
new budget amendments. Serial execution in the primary UWS worktree; coordinator
retains ledger/Git ownership. Captured integration ref `refs/heads/main`, primary
`/home/peter/Workspace/uws`, source/planning baseline `205de61c0b7307a5e659c98a99c2b06941cbcecb`.
Only the assigned three ordinary documentation files may be edited by the child.
Frozen published contracts and source code are unchanged. No concurrent milestone
lease or protected operation is dispatched. Normal model Sol high; whole review
fresh Sol xhigh. Audit disabled.

## Executed verification

All four offline checks passed at unchanged documentation input hashes: published-document immutability18.563s; root go test ./... -count=1 -timeout=2m15.322s; root go vet ./...5.235s; mkdocs build --strict1.751s. Original processes waited, output directory removed. Existing coordinator state records exact argv, hashes and stdout/stderr evidence. git diff --check PASS. No code, schema, normative specification, archive, pin, browser or publication changed; no M53 attempt consumed.
