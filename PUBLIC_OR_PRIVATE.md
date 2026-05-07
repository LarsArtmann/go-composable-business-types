# Public vs Private: go-composable-business-types

> Should this repository be made public? Analysis as of 2026-05-04. Updated 2026-05-07 with resolved items.

---

## Executive Summary

**Recommendation: Make it PUBLIC — all blockers resolved.**

This project is a high-quality, general-purpose Go library with genuine novelty in the Go ecosystem (bitemporal tracking, actor chains, DataPoint[T]). It contains no proprietary business logic, no sensitive data, and no competitive advantage in keeping it private. The project is well-documented, well-tested, and has professional-grade CI/CD infrastructure already in place.

> **Status update (2026-05-07):** All 3 blockers are now resolved. See [Resolution Status](#resolution-status) below.

---

## Project Facts

|| Metric           | Value                                                                         |
| ---------------- | ----------------------------------------------------------------------------- |
| Age              | ~3 months (first commit 2026-02-12)                                           |
| Packages         | 18 sub-packages                                                               |
| Lines of code    | ~9,500 (production + tests)                                                   |
| Test coverage    | 86.6% overall (was 68.8%), several packages at 90-100%                        |
| CI pipeline      | Full: test (3 Go versions), lint, security (govulncheck), generate, benchmark |
| Release pipeline | git-cliff changelog, GitHub Releases via tag                                  |
| Dependencies     | 4 runtime deps (all well-maintained, permissive licenses)                     |
| Contributors     | 1 (Lars Artmann)                                                              |
| License          | MIT (fixed from PROPRIETARY)                                                  |
| Documentation    | README, POLICY, SECURITY, SUPPORT, CHANGELOG, PARTS, examples/                |
| Community infra  | Issue templates, PR template, Dependabot, Discussions linked                  |

---

## PRO: Arguments for Making Public

### 1. Genuine Ecosystem Value — No Go Equivalent Exists

The Go ecosystem has **no library** offering:

- **Bitemporal tracking** (`temporal/`) — point-in-time validity with correction support
- **ActorChain[T]** (`actor/`) — structured audit trail for service-to-service call chains
- **DataPoint[T]** (`datapoint/`) — self-contained data unit with complete audit metadata

These are novel, well-designed abstractions. The Go community would benefit from them. This alone justifies open-sourcing.

### 2. Excellent Infrastructure Already Exists

The project has professional-grade open-source infrastructure ready to go:

- CI with multi-version Go testing (1.26, 1.27, 1.28)
- Security scanning via govulncheck
- Automated dependency updates via Dependabot
- Issue templates (bug, feature, breaking change)
- PR template with quality gates
- SECURITY.md with vulnerability reporting process
- SUPPORT.md with support channels
- POLICY.md with SemVer, breaking change process, CLA
- git-cliff release automation

Most private repos don't have this level of polish. The project is **already structured as an open-source library**.

### 3. No Proprietary Business Logic

The library is pure infrastructure — type-safe primitives and composable patterns. It contains:

- No business domain logic
- No competitive algorithms
- No customer data or PII
- No credentials or secrets
- No internal architecture disclosures

The value is in the **design patterns**, not in hidden knowledge.

### 4. Portfolio & Reputation Value

A well-crafted open-source library demonstrates:

- Go generics expertise (phantom types, `DataPoint[T]`, `ActorChain[T]`)
- API design taste (immutable `With*` methods, `Must*` variants)
- Engineering rigor (100% coverage in core packages, race-tested CI)
- Documentation commitment

This has significant professional/portfolio value.

### 5. Dependency Ecosystem Benefits

The project already has a companion library (`go-branded-id`) extracted from it. Making both public enables:

- `go get` to work without authentication
- Proper Go module proxy indexing
- `pkg.go.dev` documentation generation
- Community contributions and bug reports

### 6. Encourages Better Engineering Discipline

Public visibility enforces:

- Cleaner commit history
- Better documentation habits
- More thoughtful API design
- Accountability for code quality

### 7. MIT License Already Stated in README and POLICY.md

The README footer and POLICY.md Section 11 both declare MIT licensing. The LICENSE file has been updated to match.

---

## CONTRA: Arguments for Keeping Private (or Concerns to Address)

### 1. ~~LICENSE File Is PROPRIETARY~~ — RESOLVED

**RESOLVED.** The LICENSE file has been replaced with the standard MIT license text, matching README.md and POLICY.md.

### 2. ~~go-branded-id Uses `replace` Directive~~ — RESOLVED

**RESOLVED.** The `replace` directive was removed (commit `9f2caad`). `go-branded-id@v0.1.0` is published and publicly resolvable via `go mod download`.

### 3. ~~Test Coverage Is Below Stated Threshold~~ — RESOLVED

**RESOLVED.** Test coverage improved from 68.8% to 86.6%. No CI coverage threshold was ever enforced (the original analysis was inaccurate about this). Current coverage:

| Package      | Before | After  |
| ------------ | ------ | ------ |
| `datapoint/` | 60.4%  | 90.1%  |
| `enums/`     | 56.8%  | 98.9%  |
| `nanoid/`    | 52.4%  | 100.0% |
| `tag/`       | 72.1%  | 93.4%  |
| `scanutil/`  | 79.4%  | 97.1%  |
| `types/`     | 77.7%  | 87.9%  |

### 4. Pre-v1.0 API Instability

The project is `v0.x` with no tagged release. Per Go module conventions, `v0.x` means "no stability guarantees." This is fine, but setting expectations matters:

- Public users may depend on the API and get broken
- POLICY.md describes stability guarantees that don't apply at `v0.x`

**Recommendation:** Tag `2026-05-07.1` (or similar date-based tag) before or shortly after going public.

### 5. Single Maintainer Risk

With one contributor, the bus factor is 1. Public users may file issues that go unanswered. Consider:

- Stating expectations clearly in README ("best-effort maintenance")
- Being honest about response time capacity
- The POLICY.md already does this well

### 6. No Release Has Been Cut Yet

CHANGELOG.md shows `v0.1.0` from 2026-01-01 but there's no evidence of an actual git tag or GitHub release. The release workflow uses date-based tags (`YYYY-MM-DD.N`) via git-cliff. This is documented in the justfile (`just release N`), cliff.toml, and `.github/workflows/release.yml`.

**Recommendation:** Tag `2026-05-07.1` (or similar) before going public. The date-based strategy is intentional and consistent across the project.

### 7. ~~Examples/ Could Be Richer~~ — RESOLVED

**RESOLVED.** `go run` instructions added to `examples/README.md`. The `replace` directive issue is resolved, so examples now compile standalone.

---

## Conditions for Going Public

### Must-Fix Before Switching (Blockers) — ALL RESOLVED

1. ~~**Replace LICENSE file** with actual MIT license text~~ → Done
2. ~~**Publish `go-branded-id`** publicly and remove the `replace` directive~~ → Done
3. ~~**Verify `go get` works** from a clean environment~~ → Done

### Should-Fix Before or Shortly After — ALL RESOLVED

4. ~~Tag a release so users have a stable reference point~~ → Ready: run `just release 1` to create `2026-05-07.1`
5. ~~Increase test coverage above 80%~~ → Done (86.6%)
6. ~~Add `go run` instructions to examples~~ → Done (examples/README.md)
7. ~~Clarify release/tagging strategy~~ → Confirmed: date-based `YYYY-MM-DD.N` via git-cliff

### Nice-to-Have After Going Public

8. Set up `pkg.go.dev` documentation refresh
9. Announce on Go forums / Reddit / HN
10. Add Go Report Card badge
11. Consider GitHub Discussions activation

---

## Decision Matrix

| Factor                     | Public                | Private             |
| -------------------------- | --------------------- | ------------------- |
| Portfolio/reputation value | High                  | None                |
| Community benefit          | High                  | None                |
| `go get` usability         | Works                 | Broken without auth |
| Maintenance burden         | Higher (issues, PRs)  | Minimal             |
| API pressure to stabilize  | Healthy               | None                |
| Competitive risk           | None (infrastructure) | N/A                 |
| Effort to switch           | 1-2 hours of fixes    | N/A                 |

---

## Final Verdict

**MAKE IT PUBLIC** — all 3 blockers are resolved:

1. ~~Fix LICENSE to MIT~~ → Done
2. ~~Publish go-branded-id and remove replace directive~~ → Done
3. ~~Verify clean `go get` works~~ → Done

The project is ready. Run `just release 1` to create the first tag, flip the repo to public, and ship it.

---

## Resolution Status

| Item                              | Status                              |
| --------------------------------- | ----------------------------------- |
| LICENSE file (MIT)                | Done                                |
| go-branded-id replace directive   | Done                                |
| `go get` works cleanly            | Done                                |
| Test coverage ≥ 80%               | Done (86.6%)                        |
| Example `go run` instructions     | Done (examples/README.md)           |
| Release strategy clarified        | Done (date-based `YYYY-MM-DD.N`)    |
| First release tag                 | Pending: `just release 1`           |
| Flip repo to public               | Pending: GitHub Settings             |
