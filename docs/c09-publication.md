# C09 source publication evidence

The confirmed STG11_SOURCE_PUBLICATION envelope published qualified source,
review fixes and retirement closure by normal fast-forward pushes to the
exact UWS origin/main ref. No tag, force push or branch deletion occurred.

- Accepted implementation: `6a267306032edc687a298cefc8bba7019d3ad059` (whole review 3/10).
- Independently observed retirement closure: `8e5be730aa68aa4cb4f6c591a3a9425c6b8bd55c`.
- git ls-remote origin refs/heads/main returned that exact closure after push.
- The accepted implementation is an ancestor of the observed closure.
- Supplement: docs/examples/binding/v1/manifest.json, including exact upstream hashes.

Consumers pin accepted source and record observed published closure; metadata
completeness is not trust, credential readiness or execution permission. The
[full C09 record](https://github.com/OpenUdon/uws/blob/main/tabilet/docs/history/status-C09.md)
is frozen. Later metadata commits can descend from this closure without
changing accepted implementation. No installation/live-ledger/provider/API/model/mail
operation is included.
