# C13 — Versioned HCL-input retirement and replacement projection

**Stage/owner:** STG-13; UWS.
**Planning approval:** Owner-approved full rescope on 2026-10-25 (approval: “I aapaprove”), following the earlier whole-stage approval. Planning only; no execution, commit or grant activation.
**State:** Withdrawn; every untouched row is `[X]`; no work delivered. Review 0/10 remains unchanged.
**Specification:** [milestone.md](milestone.md#stg13-c13).
**Original dependencies (withdrawn; no active edges):** [Kinet:M60](../../../kinet/tabilet/memory-bank/status-M60.md)
**Source baseline:** `84cd33a5b8db14e88b3e4dd93dcc1cb0fe83e958`; existing dirty planning/history was inspected and preserved. Reconcile live source and exact accepted producers before implementation.
**Ownership:** Withdrawn; no executable task owner. Original serial ownership and task boundaries are historical.
**Authorization:** None from this withdrawn record. Original requirements are retained only as historical text in the specification; no old or proposed publication grant transfers.

**Lineage:** Fresh successor to withdrawn UWS:C11, under Kinet STG-13; no old inventory/deployment triggers or publication grants are reactivated.

**Withdrawal disposition:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row.

**Review provenance.** Kinet Stage 13 plan review (`kinet-stage13.md`), 2026-10-25; source commit baseline not supplied. H2/H3 (High/local P1): current Ramen/Udon readers and writers depend on UWS HCL. The owner selected preservation. Original task wording and review counters are retained below.

## Tasks

| Item | State | Notes |
|---|---|---|
| C13.1 — Freeze versioned retirement and migration contract | `[X]` | Use accepted Kinet:M60 readiness, identify all public HCL/wrapper entry points and choose the next compatible version boundary from current source. Add fresh migration/refusal fixtures while preserving published baselines. **Withdrawal:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row. |
| C13.2 — Add bounded JSON/YAML canonical projection | `[X]` | Provide the replacement complete-value comparison API with cancellation/resource bounds and exact numeric/extension preservation. Establish separate presentation-module adaptation before removing its old import dependency. **Withdrawal:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row. |
| C13.3 — Remove core HCL input and legacy wrappers | `[X]` | Retire UWS HCL parsing and deprecated expression wrappers in the new version, with clear errors and migration guidance. Preserve JSON/YAML, old files and independent Render/Verify behavior. **Withdrawal:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row. |
| C13.4 — Qualify separate closures and publish handoff | `[X]` | Prove root and presentation dependency/compiled closures separately; run focused/full required checks and review. Prepare authorized exact publication for APItools/Udon/OpenUdon/Browsertools/Ramen/Kinet; local tests do not prove consumer adoption. **Withdrawal:** Withdrawn by the owner-approved 2026-10-25 rescope: retain supported HCL in Ramen/Udon and the supporting UWS codecs. This obligation is undelivered, not accepted. HCL removal and legacy-wrapper retirement move to separate unnumbered Candidate Directions; no replacement is scheduled and no active milestone depends on this row. |

## Original acceptance and verification — withdrawn, undelivered

A separately versioned UWS compatibility change removes UWS HCL input and deprecated expression wrappers, retaining JSON/YAML and an independent HCL presentation module.

- Compatibility/version choice follows the current published implementation and UWS retirement rules; frozen versions, schemas and old fixtures remain immutable.
- UWS core no longer accepts HCL or legacy wrappers and has no HCL module/compiled dependency. Bounded context-aware JSON/YAML projection replaces hcl.Import for complete-value comparison.
- Preserve large integers, decimal/exponent spelling, signed zero and extension fields. HCL display remains separate; no global HCL-free claim covers native executors or Terraform.

Original planned verification (not scheduled or performed): Focused JSON/YAML ingestion/refusal and complete-value projection tests, numeric lexemes/extensions, separate HCL Render/Verify tests, standalone root and hcl-module tests/vet plus module/compiled-dependency inspection. All evidence is pending. Preserve existing package limits. Use disposable fixtures/fake providers; identify executed, mocked, compile-only and deferred proof separately.

## Review

Review gate iterations: 0

Task reviews: not started; withdrawal is not acceptance. Whole-milestone review: 0/10, not started. Follow the package's existing persisted review and closure procedure; planning, terminal rows and successful builds do not establish acceptance.
