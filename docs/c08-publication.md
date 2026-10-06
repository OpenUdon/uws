# C08 source publication evidence

The owner-confirmed STG11_SOURCE_PUBLICATION envelope published the qualified
source, reviewed fixes and retirement closure by normal fast-forward pushes to
UWS origin refs/heads/main. No tag, force push or branch deletion was performed.

- Accepted implementation: `0411eea6fc84fbd6aa97cef94f53f301260f4844` (whole review 2/10).
- Independently observed retirement closure: `5c0c74f48d84588e3ff4f994f0713f199cfcc67c`.
- `git ls-remote origin refs/heads/main` returned that exact closure after push.
- The accepted implementation is an ancestor of the observed published closure.
- The separately pinned reference supplement is docs/examples/expressions/v1/manifest.json.

Consumers pin the accepted implementation and record the observed closure;
local sibling heads or successful builds are not publication evidence. Later
metadata commits may descend from this closure without changing the accepted
implementation. The full retired specification/status remains frozen in
[the C08 record](https://github.com/OpenUdon/uws/blob/main/tabilet/docs/history/status-C08.md). No provider, model,
credential, hosted permission, installation or live-ledger operation occurred.
