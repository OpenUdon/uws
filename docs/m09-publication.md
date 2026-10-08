# M09 source publication and ordinary module proof

The user separately granted UWS:M09 source publication on 2026-10-08, in
response to the [exact proposal](m09-publication-proposal.md). Normal
fast-forward publication moved the unchanged configured origin
`git@github.com-tabilet:OpenUdon/uws.git`, `refs/heads/main`, from
`0a4597122a7baa7e79e46e79e4ec60dbfffc3720` to reviewed initial evidence head
`af225f8b0b5cfbd5848e28e5ce49eaebc0f228bb`. A fresh independent ls-remote
observed that exact head. Runtime implementation source is its ancestor
`b099f6803277ae94c7e9f1da0904a0140b278f20`.

Fresh ordinary queries of both public modules by the full runtime source resolve
`v0.0.0-20261008043726-b099f6803277`, timestamp `2026-10-08T04:37:26Z`,
GoVersion `1.25.4`, and Origin VCS git / URL https://github.com/OpenUdon/uws /
Hash b099f6803277ae94c7e9f1da0904a0140b278f20. Codec Origin.Subdir is hcl.
These queries use the configured proxy.golang.org/direct and sum.golang.org
source in a fresh independent module cache. No owner bootstrap file proxy,
workspace or directory replacement supplies this proof.

The [normalized ordinary download metadata](m09-ordinary-module-proof.json)
has SHA-256 `9933826f846830354f7ed15990c35a5c77cc7cd0582108afa06005d58e5c0b3b`.
Canonical Go archive and GoMod sums match the reviewed local candidate:

| Module | Go archive sum | GoMod sum |
|---|---|---|
| github.com/OpenUdon/uws | h1:4xy+/HBNh1CSJDO+qOzWdV/0zC/yFCKAz2kOBWufA7g= | h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw= |
| github.com/OpenUdon/uws/hcl | h1:OJsmDK/RcFpyMAMGy84DjcX6kNalYH/E0Q/C3XwD5Uo= | h1:0cR/xLzEP8vJ9FAUhsLbaKVkU7UarXPc51nJjEaMP5Q= |

Ordinary ZIP containers use different compression from the locally generated
VCS ZIPs. Every filename and content matches exactly: 360 root files and 21
codec files, including the codec's inherited root LICENSE. The canonical module
sums remain equal. Raw ordinary containers retain their own identities:

| Ordinary module ZIP | Bytes | SHA-256 |
|---|---:|---|
| Root | 2417310 | 328e5bec557b05536879bb7b3d159b3a5dbff7027cc2b7ea9a8471093a04013a |
| Codec | 48206 | 0c6c7c9ece74d888a76213b5b0db233ef31fd297c78145ab8d0513d1d994b21f |

Full selected artifact downloads succeed. The independently observed complete
root and codec closures match the [52-module root](m09-qualified-root-module-closure.json)
and [50-module codec](m09-qualified-codec-module-closure.json) closure hashes
from local qualification. Standalone codec selects accepted C09 root
`v0.0.0-20261006181058-6a267306032e`; its declared pin remains unchanged.

The independent [52-module ordinary consumer closure](m09-ordinary-consumer-module-closure.json)
has SHA-256 `42ee5d857dfc6b0c5875e95b630ef53e7019b057ced638a01b4f96ad2e1ded26`.
Its ordinary module is example.com/uws-m09-ordinary-proof, explicitly requiring
both exact proposed versions. MVS selects b099 root there. The committed
[public fixture](m09-consumer-proof.go.txt) has unchanged SHA-256
`85972965803df395afbd539bf8ae4726e8474e2ccdbe6b94706a10abdd4f770a` and passes
fresh full races without replacements. Its checks cover exact numeric/value
Render/Verify/Import and public ordinary/executable validation, flow and
portability behavior.

Both ordinary downloaded module-only copies pass fresh full suites, full race
suites, vet and builds independently with Go 1.26.6, GOWORK=off,
GOTOOLCHAIN=local, GOPROXY=off after ordinary artifact acquisition. Their
go.mod/sum bytes remain identical to reviewed source. Verification uses owned
root-filesystem GOTMPDIR/TMPDIR and the separate ordinary module cache under
`/home/peter/.cache/uws-m09-proof-b099f680`; proof files and logs are retained in
ordinary-download-proof/. Earlier local/bootstrap/preliminary evidence remains
at its recorded context and is not relabeled as independent publication.

Whole review 7 passed after this proof with no remaining blocker from either
read-only whole-milestone reviewer. M09.5 and all five task rows are qualified;
milestone acceptance still requires exact downstream reconciliation and closure.
This evidence does not retire the milestone or activate a consumer.
Parent owns cross-package reconciliation and final retirement. All publication
heads use [skip ci] under the retained convention to suppress the main-push
workflow's force docs deployment to another ref. Local checks are the evidence;
no hosted CI or docs-deployment result is claimed. No force push, tag, gh-pages,
host deployment, installed migration, registration or live run occurred.
