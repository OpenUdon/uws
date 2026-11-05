# Test policy

Every test, check and qualification run in this project follows this file:
task rows, milestone reviews, goal runs and their briefs, and launch
references. It governs how tests are written, run and recorded. It grants no
authority: it never permits a live service, credentials, real user data, a
physical browser, external network access, downloads, privileged commands or
external mutations, and it never removes an authorization grant's limits.
Commands and time budgets live in [tech-stack.md](tech-stack.md);
project-specific choices go under [Project specifics](#project-specifics).

**Precedence.** This policy governs test method. Where a pending milestone
specification, status row, `tech-stack.md` text, launch input or brief
prescribes a method it forbids, such as a supervising controller or receipt
reader around class O tests, an UNKNOWN class O outcome, class O attempt
counters or a test runner in another language without approval, this policy
wins on method. The acceptance outcome stays and is shown through self-checking
tests. It never overrides acceptance requirements, an authorization grant's
limits or frozen records.

## Two classes of test

| Class | What it is | Rule |
|---|---|---|
| **O: offline** | Fakes, scripted providers, disposable stores, isolated loopback fixtures and local sandboxed workers; no credentials, no external network access | The default. Every rule below applies. |
| **R: real effect** | Touches a live service, real user data, credentials, a physical browser or real browser profile, typed human input, privileged operations, metered external quota or the external network | Never a default test. Needs separately named human authority; see [Class R runs](#class-r-runs). |

All development verification is class O unless a milestone specification names
a class R check and a human separately authorizes it. The line is whether a run
can have a real effect, not which tool it uses: a synthetic or headless browser
test on loopback fixtures with a throwaway profile is class O.

## Class O rules

1. **Tests check themselves.** Each test owns its deadline, waits for every
   child process it starts, cleans up what it created, and fails on a leak or
   overrun. The runner's exit status is the result. The runner timeout is a
   tight backstop, about twice the expected runtime and within the hard limits
   in `tech-stack.md`; never raise it to hide a slow or hung test. Do not wrap
   tests in external supervising controllers, receipt readers or proof chains;
   the project's existing test entry points are the runners.
2. **PASS or FAIL, never UNKNOWN.** A timeout, signal, crash, leak, setup error
   or runner parse error is a FAIL. Cleanup deletes the run's own temp root and
   labelled resources; remove a leftover only when its label or temp root shows
   the run owns it, and report anything else. Fix the cause and run again; this policy
   adds no human approval or disposition for that. Approvals that the effective
   authorization requires, such as fresh approval for a protected run whose
   input changed outside an approved envelope, still apply. Never rerun
   unchanged code hoping for a pass: a flaky test is a bug.
3. **Run whole units.** Required gates run whole packages, modules or project
   targets, not per-test cohorts that each get their own plan, approval or
   review. Rule 7 covers checks made while changing code, and rule 8 covers
   skipping a fixture under an instrumented variant.
4. **Review code, not test runs.** Task and milestone reviews cover product and
   test code together. No readiness review before a run; no separate review of
   a runner or harness fix.
5. **Light bookkeeping.** This policy adds no attempt counters, failure streaks
   or numeric allowances. For each required gate the status file records the
   command, the source revision or saved diff hash (and artifact identity when
   one is qualified), PASS or FAIL, and one link to a log kept in ignored
   scratch; a reused result says its inputs are unchanged. No per-attempt
   narrative. **Authorization is separate:** when a run needs an explicit
   grant, that grant's scope, reservations, budget and stop conditions bind as
   the goal protocol states. An exhausted budget pauses the run for human
   direction even though the test is class O.
6. **Keep the safety defaults.** Offline and credential-free fixtures,
   disposable stores, loopback-only networking, required sandboxes, pinned
   artifacts and product limits stay mandatory wherever the project or a
   specification requires them. Where a specification requires actual
   qualification, use the real built artifact or binary, not a mock.
7. **Focused checks and required gates.** While changing code, run the
   narrowest check that exercises the change, without race, sanitizer, coverage
   or other slow variants. Before marking a row done or starting a task review
   or re-review, run the row check: build or type-check the whole project where
   the toolchain has such a step, which catches a broken signature in a unit
   whose tests do not run, then the project's static checks and the affected
   whole units. A project whose full build is slow names a narrower step in
   `tech-stack.md`. Add direct consumers when the row changes a public
   contract, a persisted format or behavior its consumers rely on, and the UI
   check when it changes UI. When it changes concurrency, add the race check on
   units whose concurrent code changed and on consumers of a changed
   concurrency contract. This limits the row check, never a required gate; a
   consumer broken in some other way surfaces at the required gates.
   The row check is the row's focused verification and its evidence, with any
   red-then-green regression; earlier narrow runs are not evidence. Record it
   as rule 5 records a gate. It may use the toolchain's test cache unless a test
   depends on inputs the cache cannot see, such as a pinned external artifact;
   record a result served from the cache as cached.
   Run required gates immediately before the whole review, not after each edit
   or before each task review. Every gate run, rerun included, is fresh and not
   served from a test-result cache: when a gate command can be served from one,
   use the toolchain's option to bypass it. Compiler and dependency caches stay.
   Rerun them before a later iteration or closure
   only when something a gate reads has changed: code, tests, build or tool
   inputs, dependencies, a moved base, or the documents and data the gate
   itself checks.
   If it is unclear what a gate reads, rerun it.
   Record the state they passed on as the full `HEAD` commit ID, plus a saved
   diff and its hash for uncommitted work. When comparing states, ignore status
   and record bookkeeping unless a gate reads it. Run gates cheapest first.
   After a failure, finish the other cheap gates, fix with focused checks, then
   rerun the failed and affected gates and
   run any gate that was skipped before the review.
   Where a stage needs end-to-end qualification with real pinned artifacts, it
   is a required gate of the milestone that completes the stage, which depends
   on every other milestone in that stage. Rule 6 governs those artifacts, and
   any class R part keeps its own authority.
   Grants and class R runs are unchanged, and this adds no counters, budgets or
   approvals.
8. **Match the fixture to the variant.** A large capacity or maximum-size
   fixture may skip an instrumented variant (race, sanitizer, coverage) when
   all of these hold:
   - the same fixture still runs in a required gate without the variant;
   - the skip is written in the test with its reason, uses the toolchain's own
     skip mechanism, and is reported as skipped, never a silent filter;
   - anything the variant checks on the path the fixture exercises, such as
     concurrency, including work that grows with the fixture, has its own
     smaller test under the variant. A property that appears only at full size,
     such as memory safety at a size boundary, keeps the variant: no skip;
   - the skip is not added to get past a failure the variant reported. A
     timeout is a FAIL with a cause: confirm the time is instrumentation
     overhead before skipping.

   The gate still runs whole units; only the named fixture is skipped.

## Capacity and unit size

These are recommendations, not gates. Numbers belong to each project, in
`tech-stack.md`.

- **Concurrency limits.** Run heavy instrumented targets serially on a shared
  host: isolated stores and ports do not isolate CPU. Plain targets may run in
  parallel. Record the measured reason for any limit in either direction; an
  unexplained serial setting is worth challenging.
- **Instrumentation scope.** When a profile shows that an instrumented run spends
  most of its time in third-party code, or a stock default applies, a project may
  exclude that code from instrumentation. Its own code stays fully instrumented
  and every test still runs. Record the saving, measured on first use, and what
  the exclusion can no longer detect in `tech-stack.md`. The milestone that
  completes a stage runs its instrumented gates without the exclusion. Never
  change an exclusion to get past a failure.
  In a goal run the options go in the `INSTRUMENTATION_SCOPE` launch field,
  pre-filled with the stock defaults, so the owner can edit or clear them before
  launch.
- **Test unit size.** When one unit's gate becomes slow or the unit has many
  importers, plan a split as its own milestone, and do not start it in the
  middle of a task. The project sets the threshold. Scoped checks only pay off
  when units are small. Keep test helpers near the domain they serve: a shared
  helper unit that many tests import invalidates their cached results on every
  edit.
- **Parallel tests.** Prefer tests that can run in parallel, each with its own
  disposable state. Note toolchain limits, such as tests that change the process
  environment.

## Use the project's own toolchain

Tests are written and run in the language and toolchain of the code they test,
as declared in `tech-stack.md`. Anything that runs tests, judges their outcome
or produces test evidence uses that same toolchain.

- Another language in those roles needs explicit human approval recording the
  reason, the exact scope and why the declared toolchain cannot do the job.
  List approved exceptions under Project specifics.
- Committed scripts already listed in `tech-stack.md` keep their role;
  extending one or adding another needs the same approval.
- Uncommitted one-off helpers in any language may inspect results, such as
  counting log lines or hashing a file for a status record. They never decide
  PASS or FAIL, never wrap or supervise a run, and never stand in for a test.
- This is toolchain discipline, independent of class; it applies to O and R.
  Using the project's language does not make a supervising controller
  acceptable: the controller ban in class O rule 1 applies in every language.

## Class R runs

Each class R run needs a human's named authority, which also states how the run
is supervised. Supervise with standard tooling, such as the test framework's own
lifecycle, process groups and throwaway profiles, not custom controllers.
UNKNOWN is allowed only when a real external effect may remain, and only for
that run. Class R runs use the
[verification-attempt gate](milestone.md#verification-attempt-gate).

## When to stop and ask

Only when:

- a test would need class R behavior;
- a run needs authority it lacks, or a change falls outside an approved grant or
  its envelope;
- an authorization budget or limit is exhausted;
- a required tool or input is missing (never download it); or
- passing would require changing product scope, a security boundary or this
  policy.

## Reviewing against this policy

Reviewers report as findings: an external controller or receipt reader around
class O tests; an UNKNOWN class O outcome; per-run approvals or readiness
reviews that this policy does not require and no effective authorization
requires; class O attempt counters other than grant budgets; a required-gate
result reused after something the gate reads changed; a required gate served
from a test-result cache; an instrumented-variant skip that breaks a condition of rule 8; test runners or evidence
generators in another language without approval. Approvals,
reservations and conditions required by effective authorization are not
findings. A milestone cannot close on acceptance evidence that depends on the
prohibited machinery: that acceptance line is unproven, which is P1.

## Project specifics

- **Toolchain**: Go 1.25.4 module declarations; native `go build`, `go vet` and `go test` in both the repository root and the separate `hcl/` module. See [verification commands](tech-stack.md#verification-commands).
- **Row checks**: Build and vet both modules, then test affected whole packages; add direct consumers, documentation-site checks and race checks when the change requires them, as [tech-stack.md](tech-stack.md#verification-commands) specifies.
- **Required gates**: Both modules' build, vet, full tests and full race tests; strict MkDocs and `git diff --check`. Go test gates use `-count=1`; preserve offline, standalone and exact-artifact requirements.
- **Stage qualification**: None scheduled for UWS. Future approved end-to-end qualification belongs to the final stage milestone with its dependencies and real pinned artifacts; existing qualification evidence stays at its recorded context.
- **Class R checks**: None scheduled; later named checks need separate human authority and supervision.
- **Approved other-language tools**: Existing Python-based MkDocs documentation tooling retains its documented role. No other-language Go test runner or test-evidence generator is approved.
