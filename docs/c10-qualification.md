# C10 local browser-shape candidate qualification

C10.1–3 are implemented and locally verified. C10.4 remains incomplete because
required source publication and its independent verification are pending.
Whole-milestone closing review remains 0/10. This candidate is unaccepted and
unpublished; it supplies no ordinary consumer adoption, browser, runtime,
containment, deployment or live-operation claim.

The source candidate is `f01a2542410c583d0ea909dadd8f17527cc0d27d`, on isolated
`goal/C10`, reviewed against captured integration ref `refs/heads/main` at
`4429c07bab4616ab46d9a151d50fb26b1370b1fe`. Local rebase confirmed that baseline
was unchanged. The source includes three verified task commits, a verified
browser decoder correction and four verified candidate pre-review fixes.
Qualification/handoff documents and the final blocked status are additional
uncommitted C10.4 work.

## Behavior qualified

Optional browser shape metadata preserves native profile/call identity, exact
native `id` selector, effects, confirmation, origins and symbolic slots. All
families use protocol `browser` and carry no HTTP fields or aliases. Raw schemas
and native conditional declarations retain numeric lexemes, extensions and
false/zero/empty values. Required fields, explicit empty optional inventories and
immutable resolver snapshots are checked. Complete registration1.1/1.2 slot
declarations use the exact embedded native input-slot fragment: type, label and
exactly one required/requiredWhen branch with closed fields. Required/Condition
metadata aligns with the complete raw declaration; minimal preservation
projections remain explicitly partial and indeterminate. Browser-only decoding
refuses erased empty/null HTTP fields and metadata case aliases. Source production and exact
source/subtree reproduction remain outside UWS.

Binding, flow and strict portability inspect source-bound browser body core
references while leaving native profile placeholders and private declarations
opaque. Required credential symbols map to existing `SecurityBinding.Scheme`;
host `CredentialSlot` names remain symbolic. Partial evidence stays indeterminate.
Native action integer literals/defaults use signed64 limits for Browser1.8 and
±9007199254740991 for Browser1.9/1.10 without rewriting schema bytes. Dynamic
integer references require independent const/enum/range proof. Whole-object
proofs additionally require an explicitly closed finite source inventory, no
unsupported patterns and native target coverage for every declared property.
Open or unproved extras remain indeterminate. No diagnostic establishes
credential readiness or authority.

Native numeric proof preflights complete target/source schema trees before
reading defaults, const/enum/range or descending. Its supported subset uses the
pinned2020 default and canonical Draft06/07/2019-09/2020-12 declarations.
Unsupported dialect declarations or references anywhere in schema-valued
positions make the whole tree indeterminate, retaining inherited context.
No new reference/dialect engine is introduced; generic equal-schema compatibility
remains unchanged. Literal schema-like data and property names remain data.
Symbolic array item proof remains indeterminate without complete tuple/prefix
semantics; finite constant arrays remain supported.

The [fixture-only corpus](examples/browser-shapes/v1/manifest.json) includes
eleven native action/authentication/registration profile versions and all six
exact frozen Kinet M51 artifacts. Kinet input is
`ec760d83b6e344e9e9cd7c034f9a92eba99f38b1`, manifest SHA256
`0c445a5c90d2c09be561e713c364747f4ab9a3698ea8ab46e7b7b774ccf16bad`.
Tests reproduce its complete canonical golden and positive/negative shape cases.
They run no browser or native driver.

## Local verification

Retained Go1.26.6, `GOWORK=off`, `GOTOOLCHAIN=local`, disabled checksum networking,
private GOCACHE/GOTMPDIR/TMPDIR and disposable fixture roots were used. Root and
separate unchanged HCL modules pass full tests, races and vet. Root checks were
repeated after both pre-review fixes; the unchanged standalone HCL source and
its unchanged pinned dependency evidence are retained from the prior pass.
Published-version immutability, strict MkDocs and diff checks pass. Root/codec module declarations,
dependency sums, published versions and embedded schemas are unchanged.

The exact committed root source was packaged with Go's native module ZIP helper
into a private file proxy as provisional
`v0.0.0-20261009220914-f01a2542410c`, archive sum
`h1:TTjAn0++TdGHY0vahbXw3XsICQaavBytoqTSHu0K6o0=` and GoMod sum
`h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`. All 396 included root
module filenames/content match the committed snapshot; the nested HCL module is
excluded by normal Go module packaging. ZIP-only root tests, races and vet pass.
This local artifact is packaging evidence, not published resolution.

Frozen APItools, Browsertools, OpenUdon and Udon copies and Kinet author/exec
workers build against that candidate without workspaces or directory
replacements. The unchanged Kinet host also builds; it consumes no UWS module.
Only disposable consumer go.mod/go.sum files are adapted. Udon's existing Docker
module-version replacement is retained. Full source identities, selected module
closures, artifact hashes and verification hashes are in the
[qualification manifest](c10-local-candidate-qualification.json).

Initial compilation hit `/tmp`'s user quota; completed runs used private shared
memory. A final audit found Go updating the retained cache's version list through
a borrowed symlink. Only the two introduced provisional version entries were
removed, preserving every other entry. Private version lists are now regular
copies; all frozen repository inputs and retained module source/artifact bytes
are unchanged. No provisional candidate remains in the retained cache.

## Remaining handoff

The [publication handoff](c10-publication-handoff.md) names the exact candidate
and destination approved by the human through the Kinet launch amendment at
SHA256 `f7f453d1ca3fe94068475b3f440aea89c5178f054640ac0b4aaa67c18836f75a`.
The coordinator's explicit local checkpoint exception keeps C10.4 incomplete
until publication proof exists. At this checkpoint no push, tag,
publication, download from a network, browser action, installation, host/provider
operation, user-ledger write or audit activation occurred. Required publication
and independent resolution prevent C10.4 completion and accepted closure.

## Candidate pre-review corrections

Independent read-only inspection of checkpoint
`23d445227ac09eddbed6c92bc2ba4b27f589ca99` found 0P1/2P2/0P3 without running
tests. C10-P2-1 (incomplete native registration declarations) and C10-P2-2
(open/extra-property whole-object integer proof) were persisted before fixes.
The corrected implementation and regressions pass focused tests/races/vet,
complete root checks, module-only packaging checks and all seven consumer builds.
The prior checkpoint/evidence remain retained. This is candidate qualification;
the formal closing review remains0 while required publication is unperformed.

Read-only recheck of `7b1e936e6cf7fa78c63276e585ed98a796feb8f7` found the additional
C10-P2-3 numeric reference/dialect issue (0P1/1P2/0P3); prior two findings were
fixed and the reviewer ran no dynamic tests. This finding was persisted before
correction. Effective unsafe direct/inherited/ancestor Draft07 witnesses,
isolated source proof, unproved dynamic/recursive contexts, generic equality,
reference-shaped data and safe controls pass. All affected root/module/consumer
checks were repeated for the final source. The coordinator cleared the inactive
derived C10 Go build cache between handoffs for quota recovery, preserving all
module/source/frozen/proof inputs; requalification used -p2/GOMAXPROCS2.

The next independent inspection of
`ca11bf9a02bac4c36ac5bd8b1dad2137870580dc` confirmed earlier fixes but found
C10-P2-4 (0P1/1P2/0P3): ignored Draft04 const could still falsely prove numeric
safety. The persisted correction replaces that assumption with whole-tree
supported-subset preflight. Compiled unsafe direct/inherited Draft04 and modern
array-prefix witnesses, every supported/default const/enum/range/object control,
unknown/nested dialects, reference contexts and literal schema-like data pass.
All affected root/module/consumer checks were repeated for the final source.
Required publication and formal closing review remain unperformed.
