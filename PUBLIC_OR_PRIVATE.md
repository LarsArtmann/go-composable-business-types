# Public vs Private: go-composable-business-types

> Should this repository be made public? Analysis as of 2026-05-04.

---

## Executive Summary

**Recommendation: Make it PUBLIC — conditionally.**

This project is a high-quality, general-purpose Go library with genuine novelty in the Go ecosystem (bitemporal tracking, actor chains, DataPoint[T]). It contains no proprietary business logic, no sensitive data, and no competitive advantage in keeping it private. The project is well-documented, well-tested, and has professional-grade CI/CD infrastructure already in place.

However, several blockers must be resolved **before** flipping the visibility switch.

---

## Project Facts

| Metric | Value |
| --- | --- |
| Age | ~3 months (first commit 2026-02-12) |
| Packages | 18 sub-packages |
| Lines of code | ~9,500 (production + tests) |
| Test coverage | 68.8% overall, several packages at 90-100% |
| CI pipeline | Full: test (3 Go versions), lint, security (govulncheck), generate, benchmark |
| Release pipeline | git-cliff changelog, GitHub Releases via tag |
| Dependencies | 4 runtime deps (all well-maintained, permissive licenses) |
| Contributors | 1 (Lars Artmann) |
| License | Currently PROPRIETARY — README says MIT |
| Documentation | README, POLICY, SECURITY, SUPPORT, CHANGELOG, PARTS, examples/ |
| Community infra | Issue templates, PR template, Dependabot, Discussions linked |

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

The README footer and POLICY.md Section 11 both declare MIT licensing. The only contradiction is the LICENSE file itself, which is currently PROPRIETARY. This inconsistency should be resolved regardless.

---

## CONTRA: Arguments for Keeping Private (or Concerns to Address)

### 1. LICENSE File Is PROPRIETARY — Must Be Fixed

**Blocker.** The LICENSE file says "PROPRIETARY LICENSE" with "strictly prohibited" language, but README.md and POLICY.md both say MIT. This contradiction must be resolved **before** going public. If the intent is MIT (as documented), replace the LICENSE file content.

### 2. go-branded-id Uses `replace` Directive

**Blocker.** `go.mod` contains:

```
replace github.com/larsartmann/go-branded-id => ../go-branded-id
```

This local replace directive will break for anyone `go get`-ing the module. Both `go-branded-id` must be published (or the replace directive removed) before going public.

### 3. Test Coverage Is Below Stated Threshold

The CI enforces 80% coverage threshold, but actual coverage is 68.8%. Several packages are below target:

| Package | Coverage |
| --- | --- |
| `datapoint/` | 60.4% |
| `enums/` | 56.8% |
| `nanoid/` | 52.4% |
| `tag/` | 72.1% |
| `scanutil/` | 79.4% |

**Not a blocker**, but should be improved — especially `datapoint/` as the flagship type. Either fix coverage or adjust the CI threshold to match reality.

### 4. Pre-v1.0 API Instability

The project is `v0.x` with no tagged release. Per Go module conventions, `v0.x` means "no stability guarantees." This is fine, but setting expectations matters:

- Public users may depend on the API and get broken
- POLICY.md describes stability guarantees that don't apply at `v0.x`

**Recommendation:** Tag `v0.1.0` (or similar) before or shortly after going public.

### 5. Single Maintainer Risk

With one contributor, the bus factor is 1. Public users may file issues that go unanswered. Consider:

- Stating expectations clearly in README ("best-effort maintenance")
- Being honest about response time capacity
- The POLICY.md already does this well

### 6. No Release Has Been Cut Yet

CHANGELOG.md shows `v0.1.0` from 2026-01-01 but there's no evidence of an actual git tag or GitHub release. The release workflow expects date-based tags (`YYYY-MM-DD.*`), not SemVer tags. This is a non-standard approach that may confuse users expecting `v0.1.0`.

**Recommendation:** Clarify the release/tagging strategy before going public.

### 7. Examples/ Could Be Richer

Two example programs exist (`examples/basic/`, `examples/datapoint/`) but they don't have `go run` instructions and may not compile standalone due to the `replace` directive issue.

---

## Conditions for Going Public

### Must-Fix Before Switching (Blockers)

1. **Replace LICENSE file** with actual MIT license text (matching README and POLICY.md)
2. **Publish `go-branded-id`** publicly and remove the `replace` directive from `go.mod`, OR keep `replace` but document it clearly
3. **Verify `go get` works** from a clean environment without local access

### Should-Fix Before or Shortly After

4. Tag a release (`v0.1.0` or date-based) so users have a stable reference point
5. Increase test coverage above 80% (or adjust CI threshold to match reality at 68%)
6. Add `go run` instructions to examples
7. Clarify release/tagging strategy (SemVer tags vs date-based tags in release.yml)

### Nice-to-Have After Going Public

8. Set up `pkg.go.dev` documentation refresh
9. Announce on Go forums / Reddit / HN
10. Add `go report card` badge
11. Consider GitHub Discussions activation

---

## Decision Matrix

| Factor | Public | Private |
| --- | --- | --- |
| Portfolio/reputation value | High | None |
| Community benefit | High | None |
| `go get` usability | Works | Broken without auth |
| Maintenance burden | Higher (issues, PRs) | Minimal |
| API pressure to stabilize | Healthy | None |
| Competitive risk | None (infrastructure) | N/A |
| Effort to switch | 1-2 hours of fixes | N/A |

---

## Final Verdict

**MAKE IT PUBLIC** after resolving the 3 blockers:

1. Fix LICENSE to MIT
2. Publish go-branded-id and remove replace directive
3. Verify clean `go get` works

The project is well-built, has genuine novelty in the Go ecosystem, and already has open-source-grade infrastructure. The only real cost is the ongoing maintenance commitment, which is manageable for a type library with a stable API surface.

The proprietary license file is almost certainly an oversight given that MIT is declared everywhere else — fix it and ship it.
