# M08 source publication evidence

Accepted root/codec implementation is
`c0b19385a3b034cd45de16726668b9150f0633f2`, whole review 3/10.
The confirmed STG11_SOURCE_PUBLICATION grant published that source by normal
fast-forward to the unchanged exact UWS origin/main; independent ls-remote
observed that full source after push. Both configured-registry module queries
resolve `v0.0.0-20261006224744-c0b19385a3b0` to that exact Git origin hash;
the codec identifies subdirectory `hcl`.

| Module | Archive sum | GoMod sum |
|---|---|---|
| github.com/OpenUdon/uws | h1:CXaT3naUUi4hKM4UfE41l5VxwDQ96Y8uVv5EOpheaU8= | h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw= |
| github.com/OpenUdon/uws/hcl | h1:v42UQhy/CCLAe2+BKUDQzuofjziBXcDqMJgPIXCvIRM= | h1:0cR/xLzEP8vJ9FAUhsLbaKVkU7UarXPc51nJjEaMP5Q= |

Downloaded root/codec archives pass their complete standalone suites in
disposable copies. An independent public consumer selects both exact modules
with workspaces disabled and no replacements, passes 1.13 schema lookup and
Render/Verify/Import, and resolves the complete
[51-module closure](m08-module-closure.json), SHA-256
`36a8a7a1152f9f49af5064294eca756bf9b00f4632dfa1d215112157991d1e96`.
The codec's standalone root C09 require supplies unchanged public tags/helpers;
ordinary MVS selects the final root in this consumer. All legacy/Horizon/HCL
compatibility remains explicit; no ambient sibling path or HCL-free claim.

Retirement closure `0870b0169f2dba5570b9b2646d0abaf6f86b3564` was independently observed
on the unchanged authorized origin/main after a normal fast-forward push.
Accepted source c0b19385a3b034cd45de16726668b9150f0633f2 is an ancestor of
that observed closure. The [full status/specification](../tabilet/docs/history/status-M08.md)
is frozen under the normal UWS retirement procedure. Reconciled Udon:M47/M48
and Kinet:M46/W17 statuses record the exact accepted source, observed closure
and limits; the prerequisite publication gate is satisfied.

Push heads use `[skip ci]` under
[GitHub's documented mechanism](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/skip-workflow-runs)
to prevent the automatic force docs deployment to another ref outside this
main-only grant. Local checks provide the qualification evidence; hosted CI or
deployed-docs success is not claimed. No tag/force/branch deletion, installation,
live ledger/provider/API/model/mail/registration or deployment action is included.
