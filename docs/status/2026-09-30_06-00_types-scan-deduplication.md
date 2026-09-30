# Status Report — Deduplication of `types` Scan Clones (art-dupl)

| | |
|---|---|
| **Date** | 2026-09-30 06:00 CEST |
| **Session scope** | art-dupl report (3 clone groups at `-t 2`) → deduplicate to zero harmful clones |
| **Branch** | `master` |
| **Session commits** | `c9ca228` (types refactor, net −22 lines), `52272dc` (pkg/errors test-compile fix); `AGENTS.md` note pending daemon pickup |
| **Verification** | 14/14 root packages green (`-race`), 4/4 test-bearing sub-modules green, art-dupl **0 actionable** (was 3/8), gofmt/dprint/treefmt 0 changed, golangci-lint 0 findings in touched files |

## Executive Summary

The art-dupl report flagged 3 clone groups (8 clones, all `assignment`, threshold 2). Root cause: `Duration.Scan` and `Timestamp.Scan` hand-rolled the "textual source → parse with per-source error → assign" dance twice each (string + `[]byte` branches), and `Duration` repeated it a third time in `UnmarshalJSON` — while string-based types (`Email`, `URL`) already delegate to shared helpers. Fix: one `setFrom(value, source)` method per type owning that contract. All groups eliminated, not accepted. A pre-existing `pkg/errors` test-compile failure was fixed en passant because it blocked the full-suite gate. The project's end-to-end quality gates are red for environmental reasons (nix sandbox DNS, GitHub Actions billing) — documented, not caused by this session.

---

## a) FULLY DONE

| # | Item | Evidence | Scope |
|---|------|----------|-------|
| 1 | Clone group #1 eliminated: `d.Duration = 0; return nil` ×4 (Scan nil/empty-string/empty-[]byte, UnmarshalJSON) — the "empty means zero" rule now lives in exactly one place | art-dupl final run: 0 actionable groups; commit `c9ca228` | `types/types_time.go` |
| 2 | Clone group #2 eliminated: `t.Time = parsed; return nil` ×2 (Timestamp.Scan string/[]byte) | same | `types/types_sql.go` |
| 3 | Clone group #3 eliminated: `d.Duration = parsed; return nil` ×2 (Duration.Scan string/[]byte) | same | `types/types_time.go` |
| 4 | `Duration.setFrom(value, source)` extracted (`types/types_time.go:58`) — owns empty→zero reset, `time.ParseDuration`, source-context error; `Scan` (nil/string/[]byte) and `UnmarshalJSON` all delegate | commit `c9ca228`; all Duration tests pass unmodified | `types/types_time.go` |
| 5 | `Timestamp.setFrom(value, source)` extracted (`types/types_sql.go:121`) — owns RFC3339Nano parse + source-context error; `Scan` string/[]byte delegate | commit `c9ca228` | `types/types_sql.go` |
| 6 | Behavior preservation on existing paths — all Scan error strings byte-identical ("from string"/"from []byte"); float64/int64 conversions untouched; SQL NULL → zero preserved | existing test suite green with zero test edits | `types/` |
| 7 | Pre-existing `pkg/errors` test-compile failure fixed — `errors.AsType[E]` requires `E error`; `testAs[E any]` → `testAs[E error]` (`pkg/errors/errors_structured_test.go:25`) | commit `52272dc`; `go test ./pkg/errors` green; unblocks full-suite gate | `pkg/errors/` |
| 8 | Full test verification: root module 14/14 packages `ok` with `-race`; sub-modules nanoid, locale, money, datapoint `ok` | `GOEXPERIMENT=jsonv2 go test -race ./...` (both runs) | repo-wide |
| 9 | Lint & format clean on all touched files — golangci-lint 0 findings in `types_time.go`, `types_sql.go`, `errors_structured_test.go`; `gofmt -l` empty; treefmt 0 changed; dprint ✔ | BuildFlow `golangci-lint [root]` step; `nix build .#checks.x86_64-linux.format` log | touched files |
| 10 | Environmental gotcha documented in `AGENTS.md` — nix sandbox DNS refusal (go1.27.0 toolchain download), misleading treefmt "formatting failures" message (actually 0 changed files), go-licenses vs Go 1.27 incompatibility | `AGENTS.md` § Release & CI | docs |
| 11 | All session code changes committed | daemon commits `c9ca228`, `52272dc`; working tree clean except `AGENTS.md` note | git |

## b) PARTIALLY DONE

| # | Item | What works | What remains | Blocker | Effort |
|---|------|-----------|--------------|---------|--------|
| 1 | BuildFlow `--build-mode dev` quality gate | All code-owned steps ✔: golangci-lint, dprint/nix-fmt, go-mod-tidy/normalize, govulncheck, nix lints, jscpd | 7 steps fail environmentally: `nix-build`, `nix-build-verify`, and `license-check` per module (observed root + nanoid) | nix sandbox has no DNS (toolchain download refused); go-licenses incompatible with Go 1.27 (upstream #128) | M–L |
| 2 | Error-message contract for `setFrom` | JSON path improved to consistent `cannot parse %q from JSON` | No test pins the error-text contract; no CHANGELOG entry for the behavior-visible text change | None — just not done in this pass | S |
| 3 | Repo-wide lint health | Touched files: clean | 5 pre-existing warnings: `address.Address` missing Line2/State; `contact.Contact` missing Email/Phone/Website/Address; gocyclo 24 `ScanEnum`; gocyclo 24 `testEnumScanPointerCases`; gocyclo 21 `TestDurationSQL` | Pre-existing, out of session scope | M |
| 4 | art-dupl coverage | Root run scanned 54 Go files → 0 actionable | Whether all 6 sub-modules were included is unverified; suppressed-group population (53) not sampled for near-threshold clones | None | S |
| 5 | gopls diagnostics | Compiler + tests prove `errors_structured_test.go` compiles | Stale gopls error still displayed (LSP cache not refreshed via restart) | Cosmetic; tool cache | S |

## c) NOT STARTED

Observed during this session, deliberately not acted on (session scope = dedup + verify; per instructions, report only):

1. **docs-health HARVEST** of section (f) into `TODO_LIST.md` / `ROADMAP.md` — the status-report loop-closure step; waiting for user go-ahead.
2. **CHANGELOG entry** for the Duration JSON error-text change — `CHANGELOG.md` exists, untouched.
3. **Timestamp JSON support** (`MarshalJSON`/`UnmarshalJSON`) — Duration has JSON methods, Timestamp does not; asymmetry observed, intent unknown.
4. **nix sandbox toolchain fix** (vendor/pin go1.27 toolchain so flake checks build hermetically).
5. **GitHub Actions billing** — CI is fully red at account level.
6. **go-licenses replacement/config** for Go 1.27.
7. **`buildflow doctor`** — 9 tools unavailable (health check), never identified.
8. **BuildFlow binary rebuild** — binary `e881e96` vs repo HEAD `69dd195` (advisory warning).
9. **BuildFlow cache DB VACUUM** — 0.81 GB cache (31% free pages), 0.26 GB state DB.
10. **go.mod `go` line stabilization** — flip-flopped 20× in last 20 commits (tooling dispositions fighting).
11. **gocyclo refactors** — `scanutil.ScanEnum` (24), `testEnumScanPointerCases` (24), `TestDurationSQL` (21).
12. **`address.Address` / `contact.Contact` missing-field decisions** — add fields or justify suppression.
13. **testify → ginkgo/gomega migration** — `testify` is a root-module dep, banned per policy (per AGENTS.md).
14. **Benchmark sanity run** — CI claims benchmark jobs; never run locally this session.
15. **Release/tag** — no version cut; user decision.

## d) TOTALLY FUCKED UP

Nothing this session produced is broken — verified by tests, lint, and a clean dup scan. Radical honesty about the repo state observed:

1. **The project has NO green end-to-end quality gate anywhere right now.**
   - *What's broken:* local flake checks (`nix build .#check-*`) fail 7/7 historically (sandbox DNS refuses the go1.27.0 toolchain download); GitHub Actions is 100% red (billing). Quality enforcement has degraded to manual `GOEXPERIMENT=jsonv2 go test/lint` invocations.
   - *Severity:* High — regressions can slip through if the manual loop is skipped.
   - *Root cause:* environmental (sandbox network isolation) + account-level (billing). Pre-existing; not this session.
   - *Mitigation:* the manual command pair in AGENTS.md (used all session); gotcha documented.
2. **Misleading failure signal in treefmt-in-nix.** It prints "failed to finalise formatting: formatting failures detected" when the real cause is a toolchain download error (it reports **0 changed files**). Any future session or CI triage will chase a phantom formatting problem first. Documented in AGENTS.md; still fucked until the flake is fixed or the message path fixed upstream.
3. **Git history quality is degraded by heuristic daemon commits** ("chore: auto-commit N file(s)"). The dedup refactor's *why* exists only in this report, not in history. Severity: Medium (archaeology cost). Root cause: auto-commit daemon design. Mitigation: detailed reports like this one; commit critical artifacts manually where the harness allows.

## e) WHAT WE SHOULD IMPROVE

1. **Pin error-text contracts with table-driven tests** — `setFrom` now centralizes scan error construction across 5 sites; one table (`source ∈ string/[]byte/JSON`) locks the contract so future refactors can't silently reword it. Impact: prevents contract drift; Effort: S.
2. **Review the `setFrom("", "SQL NULL")` nil-case idiom** — the source label is never used on that path (empty input returns before the error branch); it is documentary-only. A reviewer may reasonably ask why a dead label is passed. Consider a named reset or accept with a rationale. Effort: S.
3. **Inspect BuildFlow auto-repairs before accepting them** — `github-actions-pinning:repair` ran and modified files unreviewed this session; I verified my own diffs but not the tool's. Impact: unreviewed automated changes can ride in with any run. Fix: always `git status`/diff after BuildFlow runs. Effort: S.
4. **Don't stack fixes blind on freshly-modified files** — my `pkg/errors` one-liner landed on a file another session had modified in the immediately-prior commit (`328771e`, 21 lines). It worked out, but I audited after, not before. Impact: low but real merge/conflict risk across concurrent sessions. Effort: S (process).
5. **Verify dedup tooling scope per module** — run art-dupl per sub-module so clones hiding behind root-run scoping can't accumulate. Effort: S.
6. **Close the status-report loop every time** — section (f) must be HARVESTed into `TODO_LIST.md`/`ROADMAP.md` or it dies in this timestamped file. Effort: S.
7. **Restore one green end-to-end gate** (see d1) — everything else is tolerable; a silent-repo-with-red-gates is not.

## f) Top 50 Things To Get Done Next

Ranked by impact; brainstorm list (HARVEST fuel) — most beyond the top ~15 are ROADMAP items, not commitments.

| # | Task | Impact | Effort | Category |
|---|------|--------|--------|----------|
| 1 | Add table-driven test pinning `setFrom` error text for string/[]byte/JSON sources | High | S | Quality |
| 2 | Run docs-health HARVEST: move this report's section (f) into TODO_LIST.md / ROADMAP.md | High | S | Documentation |
| 3 | Add CHANGELOG entry for the Duration JSON error-text change | High | S | Documentation |
| 4 | Inspect the `github-actions-pinning:repair` diff BuildFlow applied this session | High | S | Cleanup |
| 5 | Audit pre-session commit `328771e` (pkg/errors, enums_enum.go, go.mod/go.sum churn from another session) | High | S | Cleanup |
| 6 | Fix GitHub Actions billing so CI runs again | Critical | S | Bug |
| 7 | Make flake checks hermetic: vendor/pin the Go toolchain so nix checks build without network (per nix-private-go-repos pattern) | Critical | L | Bug |
| 8 | Replace or pin go-licenses with a Go-1.27-compatible license checker | Medium | M | Bug |
| 9 | Run `buildflow doctor` and resolve the 9 unavailable tools | Medium | S | Quality |
| 10 | Rebuild/reinstall BuildFlow binary (e881e96 → 69dd195) | Low | S | Cleanup |
| 11 | Stabilize the go.mod `go` line (20 flips in 20 commits — align go-version-auto-configure vs go-mod-update dispositions) | High | M | Bug |
| 12 | VACUUM/purge BuildFlow cache DB (0.81 GB) and state DB (0.26 GB) | Low | S | Cleanup |
| 13 | Refactor `scanutil.ScanEnum` (gocyclo 24) — extract case-groups into helpers | Medium | M | Quality |
| 14 | Refactor `testEnumScanPointerCases` (gocyclo 24) into table-driven subtests | Medium | S | Quality |
| 15 | Split `TestDurationSQL` (gocyclo 21) into focused per-behavior tests | Medium | S | Quality |
| 16 | Decide `address.Address` missing fields (Line2, State): add or suppress with rationale | Medium | M | Feature |
| 17 | Decide `contact.Contact` missing fields (Email, Phone, Website, Address): add or suppress | Medium | M | Feature |
| 18 | Decide Timestamp JSON parity: add `MarshalJSON`/`UnmarshalJSON` (RFC3339) or document the asymmetry | High | M | Feature |
| 19 | Audit JSON round-trip support across all wrapper types; fill gaps + tests | Medium | M | Feature |
| 20 | Run art-dupl per sub-module (nanoid/locale/money/datapoint) at `-t 2` | Medium | S | Quality |
| 21 | Sample-review the 53 art-dupl-suppressed groups for near-threshold harmful clones | Medium | M | Quality |
| 22 | Restart LSP / refresh stale gopls diagnostic on `errors_structured_test.go` | Low | S | Cleanup |
| 23 | Run `go generate ./...` and verify `enums_enum.go` has no drift | Medium | S | Quality |
| 24 | Run the benchmark suite locally (CI claims benchmark jobs; never verified locally) | Low | S | Quality |
| 25 | Inventory testify usage; migrate banned dep to ginkgo/gomega per policy | Medium | L | Cleanup |
| 26 | Ginkgo DescribeTable scan matrix for all `sql.Scanner` types (nil/string/[]byte/int64/float64 per type) | Medium | M | Quality |
| 27 | Add test for `Duration.Scan` float64 negative/fractional truncation (`time.Duration(int64(v))` semantics) | Low | S | Quality |
| 28 | Add Timestamp.Scan tests with timezone-offset RFC3339 strings (only `Z` tested today) | Low | S | Quality |
| 29 | Add nil-receiver Scan tests (`errDurationScanNil`, `errTimestampScanNil` paths) | Low | S | Quality |
| 30 | Add `errors.Is` tests for sentinel wrapping (`errDurationCannotScan`, `errTimestampCannotScan`) | Low | S | Quality |
| 31 | Add Scan∘Value round-trip identity test for Duration/Timestamp/Email/URL/Cents | Medium | S | Quality |
| 32 | Add `sql.Null[T]` interop tests for wrapper types | Low | S | Quality |
| 33 | Decide whether `setFrom` should be exported (`SetFrom`) for consumer parsing from arbitrary sources | Medium | S | Feature |
| 34 | Consider sentinel "cannot parse" errors (typed) instead of pure `fmt.Errorf` strings, for `errors.Is` matching | Medium | M | Feature |
| 35 | Document the Scan/Value/JSON support matrix per type in README | Medium | S | Documentation |
| 36 | Add `docs/DOMAIN_LANGUAGE.md` entries: "empty means zero" rule, `setFrom` contract, source-label convention | Low | S | Documentation |
| 37 | Add a CI job running art-dupl with a baseline threshold as a regression gate | Medium | M | Quality |
| 38 | `buildflow precommit install` if no pre-commit hook is active in this repo | Low | S | Quality |
| 39 | Run full `buildflow --build-mode full` once environment issues are resolved (tests + coverage) | High | M | Quality |
| 40 | Verify art-dupl's `-t` counts statements vs tokens assumption on this repo (document the threshold policy in AGENTS.md) | Low | S | Documentation |
| 41 | Consider `Timestamp` accepting unix-epoch int64 in Scan (API question; see g3) | Low | S | Feature |
| 42 | Check `bounded`, `importance`, `tag`, `temporal` packages for the same hand-rolled Scan pattern (clone prevention at source) | Medium | S | Quality |
| 43 | Tag a release after the dedup lands (version per AGENTS.md tag format) — user decision | Medium | S | Cleanup |
| 44 | Re-verify `GOEXPERIMENT=jsonv2` requirement notes against current go.mod (`go 1.26.4` + go1.27 toolchain request mismatch) and reconcile docs | Medium | S | Documentation |
| 45 | Add negative test: `Duration.UnmarshalJSON` with non-string JSON still errors after delegation | Low | S | Quality |
| 46 | Consider hoisting `errDurationCannotScan`-style sentinels into a shared scan-error type across types package | Low | M | Feature |
| 47 | Confirm license headers/legal files match `license-sync` expectations once license-check is fixed | Low | S | Cleanup |
| 48 | Add `examples/` snippet showing the `setFrom`-style delegation pattern for library consumers | Low | S | Documentation |
| 49 | Re-run this art-dupl command after tasks 42–46 to confirm the report stays at 0 actionable | Medium | S | Quality |
| 50 | Schedule the flake `treefmt` misleading-error-message fix upstream (BuildFlow/nix checker) once root cause (toolchain download) is resolved | Low | M | Bug |

## g) Questions I Cannot Answer Myself

1. **Timestamp JSON asymmetry** — `Duration` has `MarshalJSON`/`UnmarshalJSON`; `Timestamp` has neither (verified: the only RFC3339 sites are `Scan`). Is that deliberate (timestamps left to `time.Time`'s native JSON) or a gap I should close with RFC3339 JSON methods?
2. **Quality-gate strategy** — should I invest in making the flake checks hermetic (vendor/pin the go1.27 toolchain so `nix build .#check-*` work in the DNS-less sandbox), or are flake checks expected to run only in networked CI — in which case fixing GitHub Actions billing is the only real path to a green gate? These lead to very different work.
3. **Error-text contract policy** — do any consumers match on these error strings? The JSON path now reads `duration: cannot parse %q from JSON:` (was `... %q:`). I need to know whether error-message wording counts as a breaking change (→ minor/major bump + CHANGELOG) or as internal detail.

---

## Self-Review: What I Forgot / Could Do Better

- **Forgot:** CHANGELOG entry for the behavior-visible JSON error-text change; inspecting BuildFlow's `github-actions-pinning:repair` diff; the docs-health HARVEST loop-closure (constrained to report-only this session); verifying art-dupl's sub-module coverage; benchmark sanity run.
- **Could do better:** add the error-contract tests in the same pass as the refactor (the centralization was the perfect moment); run `buildflow doctor` when the 9-tool health warning appeared instead of deferring; audit the prior session's commit `328771e` *before* stacking my one-line fix on the same file.
- **Still improve:** everything in (d) — one green end-to-end gate; the gocyclo outliers and missing-field model decisions in (b3); contract-test coverage so "zero harmful duplication" claims are enforced by CI, not just observed once.
