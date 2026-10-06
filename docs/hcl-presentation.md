# Verified HCL presentation

UWS 1.13 adds the [presentation/deprecation contract](https://github.com/OpenUdon/uws/blob/main/versions/1.13.0.md#12-verified-hcl-presentation-and-deprecation).
Its wire fields and execution rules remain the UWS 1.12 rules. The separate
`github.com/OpenUdon/uws/hcl` module owns deterministic Render, independent
lossless Verify and deprecated Import. M08.1 publishes this contract; module
implementation and qualification belong to M08.2–M08.4 and remain pending.

Supply explicit JSON or YAML bytes and the exact codec revision. A successful
view must carry source format/raw SHA-256, codec contract/revision and HCL SHA-256.
Hash agreement alone is not a semantic proof. Verify must compare the complete
source projection and required numeric lexemes without float64 or JCS equality.
Large integers, precise decimals, exponents, signed zero, strings, dynamic keys,
extensions, null/absence and empty containers must retain their meaning.
An unsupported source or failed proof suppresses the entire view.

The HCL view is derived: it does not replace source bytes, become a package input,
grant approval, or authorize execution. Label any historical digest-bound
packaged workflow.hcl separately. Changed bytes or codec revision require fresh
verification; conversion never transfers grants or reactivates a schedule.

Import is an explicit compatibility API for the inert supported HCL subset.
Its JSON output is a proposal requiring the consuming product's validation and
confirmation. It cannot load files/variables, invoke HCL functions, fetch URLs,
resolve credentials or execute workflows. Existing root `convert`/`uws1` HCL
APIs and older inputs remain available throughout Stage 11; browser migration,
input removal and transitive dependency extraction require later qualification.

Consuming workers own hard process CPU/RSS/deadline limits, mounts/network,
cancellation, teardown and private-value filtering. This specification/module
does not supply those controls or assert that UWS core is HCL-free. Source parsing,
package trust, credentials and runtime authorization remain with their owners.
