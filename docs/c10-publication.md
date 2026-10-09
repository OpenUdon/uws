# C10 ordinary source publication and module proof

Reviewed implementation `f01a2542410c583d0ea909dadd8f17527cc0d27d` is publicly
reachable through checkpoint `e0ee25b72e6d0a6ddfdcc753ef9ad8a1ff7c009a` on
`git@github.com-tabilet:OpenUdon/uws.git`, `refs/heads/main`. The authorized
coordinator performed a normal fast-forward push, resolved the remote again,
and independently fetched it into a fresh bare repository. Follow-up checkpoint
changes contain qualification/status evidence only. Local integration remains
at captured `4429c07bab4616ab46d9a151d50fb26b1370b1fe` until the closing gate.

Fresh ordinary Go retrieval from `https://proxy.golang.org`, with direct native
`sum.golang.org` signature/inclusion verification, proves:

- Module `github.com/OpenUdon/uws`.
- Version `v0.0.0-20261009220914-f01a2542410c`.
- Source hash `f01a2542410c583d0ea909dadd8f17527cc0d27d`.
- Archive sum `h1:TTjAn0++TdGHY0vahbXw3XsICQaavBytoqTSHu0K6o0=`.
- GoMod sum `h1:DlqFOnO9lbmYWLLIh5WicNX6NTWIuytU6mIHmxj9BVw=`.
- Public ZIP SHA256 `5f0478a4a2a87fa54e921a685d70364ade125d01cb5dacc0bcd09ec11d7a2c47`.
- All396 included filenames and contents match the reviewed source manifest.

The public and provisional ZIP container hashes differ; the native archive sum
and every source-file hash match. The nested HCL module is excluded by ordinary
Go packaging and was not downloaded or changed. Its retained unchanged-source
verification remains in [local qualification](c10-qualification.md).

A task-owned HTTPS forwarding filter restricts retrieval to the exact approved
public artifact and checksum requests, verifies upstream TLS, and rejects other
endpoints or redirects. It supplies no generated/private module artifact.
`GOPROXY` remains the public URL, no direct VCS fallback or credentials are used,
and Go independently verifies signed checksum proof. The first local attempt
refused Go's equivalent percent-encoded exclamation marks before upstream
contact; the filter was corrected and that attempt retained. The successful
retrieval used a fresh isolated cache.

After acquisition, module-only tests, races and vet and all seven frozen consumer
builds/selected closures pass with `GOWORK=off`, `GOTOOLCHAIN=local`,
`GOPROXY=off`, `GOSUMDB=off`, Go1.26.6 and bounded compiler concurrency.
APItools, Browsertools, OpenUdon, Udon and Kinet author/exec use disposable
requirement adaptations; the unchanged Kinet host builds independently. No
workspace or directory replacement substitutes for the published UWS module.
Other dependencies use retained acquired cache inputs. Actual consumer pins and
source files are unchanged; these fixtures prove compatibility, not adoption.

A first public-module race link hit the user quota despite free shared memory.
Only two inactive task-owned derived Kinet Go caches were cleared; source,
module archives, artifacts and qualification logs were preserved. The affected
race/vet and consumer checks then passed, without relaxing assertions.

[Complete publication manifest](c10-publication.json) records exact Git/module
proof, source identities, all seven selected closures and verification hashes.
The human approval is the conversation approval of Kinet suggested.txt SHA256
`f7f453d1ca3fe94068475b3f440aea89c5178f054640ac0b4aaa67c18836f75a`;
these records are evidence, not authority. Source publication supplies no
browser, deployment, live/provider action or automatic sibling adoption.
Audit remains disabled. C10.4 qualification and whole review1 pass, with no
remaining findings. The P3 retained-log locator was corrected without changing
evidence bytes. Serialized integration and accepted retirement remain pending.
