# M09 local release qualification

Reviewed implementation source is
`b099f6803277ae94c7e9f1da0904a0140b278f20`, whole review **6/10 passed** on
2026-10-08. Both read-only whole-milestone reviewers found no remaining blocker.
The final whole review **7/10 passed** after independent ordinary proof; all five
M09 task rows are qualified. This record preserves local qualification; the separately granted
[publication and independent ordinary module proof](m09-publication.md) now
records actual source availability. Milestone acceptance/retirement remain with
the coordinator after review and downstream reconciliation.

Both module archives have proposed version
`v0.0.0-20261008043726-b099f6803277` and are generated from that exact clean Git
source using `golang.org/x/mod/zip.CreateFromVCS`, excluding workspace and
directory replacements. Source files added after that commit are evidence
envelope material, not part of these implementation archives.

| Module archive | Bytes | SHA-256 | Go archive sum |
|---|---:|---|---|
| github.com/OpenUdon/uws | 2371390 | b44da9795df9778e499d7b68949abd40d9f4804ab8ca4ff1c5bf8e1d6a677221 | h1:4xy+/HBNh1CSJDO+qOzWdV/0zC/yFCKAz2kOBWufA7g= |
| github.com/OpenUdon/uws/hcl | 46925 | 9c70fea971c0d91e7c5bd97cf79cb5c264bea7616e86163831f36f30b9cfd6b8 | h1:OJsmDK/RcFpyMAMGy84DjcX6kNalYH/E0Q/C3XwD5Uo= |

GoMod sums remain `h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`
for root and `h1:0cR/xLzEP8vJ9FAUhsLbaKVkU7UarXPc51nJjEaMP5Q=` for codec.
Archive go.mod/sum bytes remain identical to the committed sources after tests.

The complete selected-version closures retain every selected module, with sums
where present in ordinary module metadata:

| Proof | Modules | SHA-256 |
|---|---:|---|
| [Root archive closure](m09-qualified-root-module-closure.json) | 52 | f55c376f4caa2a20c6ccb5f2282c03eff184f13d6af18de4d095e4ac859233fb |
| [Standalone codec closure](m09-qualified-codec-module-closure.json) | 50 | b5842dd26dc22927ff290787600c518c633f0b0c3b801ee4dbc72bd2435fa2b9 |
| [Owner bootstrap consumer closure](m09-bootstrap-consumer-module-closure.json) | 52 | f57fe84abfcd0acc67694fb9f4b055eb5366981285745e26b33b30756fbd6b66 |

Standalone codec still selects accepted C09 root
`v0.0.0-20261006181058-6a267306032e`. Its own bounded source decoder protects the
retained helper. The combined consumer explicitly requires both proposed modules;
ordinary MVS selects new root there. No codec declared-pin upgrade is implied.
Horizon/HCL and legacy/browser compatibility remain in the graph.

Root and codec exact module-only archives pass full race suites, vet and builds
independently with Go 1.26.6, GOWORK=off, GOPROXY=off and GOTOOLCHAIN=local.
Fresh complete repository root/codec tests and affected binding/expressions/
strictjson races pass. Strict MkDocs, gofmt and staged/new-file whitespace checks
pass. Published versions 1.0–1.13, schema archive, grammar, conformance/testdata,
browser paths, module manifests and all existing frozen pins have no diff.
Ordinary Parse, model and semantic/executable validation code are unchanged.

Initial concurrent archive linking/unpacking hit `/tmp`'s separate user quota.
Successful complete serial checks use owned GOTMPDIR/TMPDIR under
`/home/peter/.cache/uws-m09-proof-b099f680`, with the isolated bootstrap module
cache on the same root filesystem. The coordinator cleared only old regenerative
compilation-cache entries. Source, module-download caches and retained proofs
were preserved; no waiver or skipped check applies.

The [public consumer fixture](m09-consumer-proof.go.txt), SHA-256
`85972965803df395afbd539bf8ae4726e8474e2ccdbe6b94706a10abdd4f770a`, passes
full races using exact archives through an isolated local file proxy and module
cache, with no workspace/replacement. It proves public Render/Verify/Import
complete-value and numeric-lexeme behavior, ordinary/executable validation,
global goto flow resolution and fresh root portability contexts. This is owner
bootstrap evidence, not public registry availability or independently published
module provenance. Later ordinary-source proof must use the configured upstream
source, not this local proxy/cache.

Retained candidates are preliminary at their original contexts. They are not
accepted releases or substitutes for the final source:

| Source | Context |
|---|---|
| 7074a1fec59bd9a9f9f653f882b0ae5a68a229a1 | Implementation rows complete; review 1 found defects. |
| edd4db40f55ca065b1d8fe6064487df7631f9a74 | Review 1 corrections; review 2 found defects. |
| 8b2f8dfb5da419284ee98e9404b5279bec463bae | Review 2 corrections; review 3 found defects; [root](m09-root-module-closure.json)/[codec](m09-codec-module-closure.json) closure bytes preserved. |
| adb0a5aa74179f335ca5c5dddab4bbd84089e39d | Review 3 corrections; review 4 found defects; [root](m09-review4-root-module-closure.json)/[codec](m09-review4-codec-module-closure.json) closure bytes preserved. |
| e89060b456ebcedb6cf3e1427fe60af9d734c04e | Review 4 corrections; review 5 found defects. |

All candidate archives remain under `/tmp/uws-m09-qualification.iO5IeV` in
their recorded qualification directories; the selected archives/consumer live
in `review5-qualified/`. Persisted status retains all 23 P2 findings, fixes,
counter, bounded reproductions and infrastructure outcomes. No source authority,
consumer approval, installation or live-run authority follows from local tests.
Stage 11 history and installed/frozen consumers stay preserved. Direction and
evolution version are unchanged.
