# Policy Implementation Checklist

> Tracking implementation of library policies defined in POLICY.md

**Created:** 2026-03-27  
**Policy Version:** 1.0.0

---

## Phase 1: Documentation (Completed ✓)

| Task                | Status | Notes                                       |
| ------------------- | ------ | ------------------------------------------- |
| Create POLICY.md    | ✅     | 629 lines covering all 11 policy areas      |
| Create SUPPORT.md   | ✅     | 223 lines with support channels and FAQ     |
| Update CHANGELOG.md | ✅     | Added v0.2.0 entry for policy documentation |
| Update README.md    | ✅     | Added policy link in documentation section  |
| Create SECURITY.md  | ✅     | Security policy with disclosure process     |

---

## Phase 2: GitHub Configuration

### 2.1 Issue Templates

| Template        | Purpose                                    | Status |
| --------------- | ------------------------------------------ | ------ |
| Bug Report      | Structured bug reports                     | ✅     |
| Feature Request | New feature proposals                      | ✅     |
| Breaking Change | Breaking change RFCs                       | ✅     |
| Security        | Vulnerability reports (redirects to email) | ✅     |
| Config          | Contact links and blank issue settings     | ✅     |

### 2.2 Pull Request Template

| Task                                 | Status |
| ------------------------------------ | ------ |
| Create PULL_REQUEST_TEMPLATE.md      | ✅     |
| Include conventional commit reminder | ✅     |
| Include quality gates checklist      | ✅     |

### 2.3 Labels

| Label              | Color   | Purpose                    | Status    |
| ------------------ | ------- | -------------------------- | --------- |
| `breaking-change`  | #b60205 | API breaking changes       | ⚠️ Manual |
| `security`         | #d93f0b | Security-related issues    | ⚠️ Manual |
| `deprecation`      | #f9d0c4 | Deprecation notices        | ⚠️ Manual |
| `documentation`    | #0052cc | Documentation updates      | ⚠️ Manual |
| `good-first-issue` | #7057ff | Beginner-friendly issues   | ⚠️ Manual |
| `help-wanted`      | #008672 | Community help needed      | ⚠️ Manual |
| `bug`              | -       | Bug reports (default)      | ✅        |
| `enhancement`      | -       | Feature requests (default) | ✅        |

### 2.4 Branch Protection

| Rule                           | Status           |
| ------------------------------ | ---------------- |
| Require PR reviews (1 minimum) | ⚠️ Manual        |
| Require status checks to pass  | ✅ CI configured |
| Require up-to-date branches    | ⚠️ Manual        |
| Require linear history         | ⚠️ Manual        |
| Restrict push to main          | ⚠️ Manual        |

---

## Phase 3: Automation

### 3.1 CI/CD Enhancements

| Workflow             | Purpose                          | Status                |
| -------------------- | -------------------------------- | --------------------- |
| Dependabot           | Automated dependency updates     | ✅                    |
| govulncheck          | Vulnerability scanning           | ✅ Added to CI        |
| Coverage reporting   | Code coverage tracking           | ✅ Codecov configured |
| Benchmark comparison | Performance regression detection | ✅ CI configured      |

### 3.2 Pre-commit Hooks

| Hook                 | Purpose                          | Status                     |
| -------------------- | -------------------------------- | -------------------------- |
| Conventional commits | Enforce commit message format    | ⚠️ Manual - see .git/hooks |
| Lint checks          | Prevent commits with lint errors | ⏳                         |
| Test execution       | Require tests to pass            | ⏳                         |

### 3.3 Release Automation

| Task                           | Status                    |
| ------------------------------ | ------------------------- |
| Automated changelog generation | ⏳                        |
| Automated GitHub releases      | ⚠️ Manual - see POLICY.md |
| Automated version tagging      | ⚠️ Manual - see POLICY.md |

---

## Phase 4: Governance

### 4.1 Security

| Task                                 | Status                     |
| ------------------------------------ | -------------------------- |
| Set up security@lars.software email  | ⚠️ External setup required |
| Create SECURITY.md (GitHub standard) | ✅                         |
| Document CVE process                 | ✅ In SECURITY.md          |

### 4.2 Legal

| Task                         | Status |
| ---------------------------- | ------ | --------------------------- |
| Set up CLA assistant         | ⏳     |
| Review license compatibility | ✅     | All deps MIT/Apache-2.0/BSD |

---

## Phase 5: Monitoring

### 5.1 Metrics

| Metric              | Tool       | Status           |
| ------------------- | ---------- | ---------------- |
| Download statistics | pkg.go.dev | ✅               |
| Issue response time | GitHub     | ✅ Built-in      |
| PR merge time       | GitHub     | ✅ Built-in      |
| Test coverage       | Codecov    | ✅ CI configured |

### 5.2 Alerts

| Alert                  | Trigger                | Status           |
| ---------------------- | ---------------------- | ---------------- |
| Security vulnerability | govulncheck failure    | ✅ CI configured |
| Breaking build         | CI failure             | ✅ Built-in      |
| Deprecation deadline   | 30 days before removal | ⚠️ Manual        |

---

## Implementation Priority Matrix

| Priority | Task                  | Effort | Impact | Timeline | Status      |
| -------- | --------------------- | ------ | ------ | -------- | ----------- |
| P0       | Set up security email | Low    | High   | Week 1   | ⚠️ External |
| P0       | Create SECURITY.md    | Low    | High   | Week 1   | ✅ Done     |
| P1       | Issue templates       | Medium | Medium | Week 1-2 | ✅ Done     |
| P1       | Branch protection     | Low    | High   | Week 1   | ⚠️ Manual   |
| P1       | Labels setup          | Low    | Medium | Week 1   | ⚠️ Manual   |
| P2       | Dependabot            | Low    | Medium | Week 2   | ✅ Done     |
| P2       | govulncheck           | Low    | High   | Week 2   | ✅ Done     |
| P2       | Coverage reporting    | Medium | Medium | Week 2-3 | ✅ Done     |
| P3       | Pre-commit hooks      | Medium | Low    | Week 3-4 | ⏳          |
| P3       | Release automation    | Medium | Low    | Week 4   | ⏳          |
| P4       | Benchmark comparison  | High   | Low    | Month 2  | ✅ CI done  |

---

## Success Criteria

### Documentation

- [x] POLICY.md is complete and accurate
- [x] SUPPORT.md is helpful for users
- [x] All policy references resolve correctly

### GitHub

- [x] All issue templates are functional
- [x] PR template is configured
- [ ] Labels are created and documented (⚠️ Manual - GitHub UI required)
- [ ] Branch protection is active (⚠️ Manual - GitHub UI required)
- [ ] CLA assistant is configured (⚠️ External service)

### Automation

- [x] CI passes on all PRs
- [x] Dependabot creates PRs for updates
- [x] govulncheck runs regularly
- [x] Coverage reports are generated

### Adoption

- [ ] First issue uses template
- [ ] First PR follows conventional commits
- [ ] First security report handled per policy
- [ ] First deprecation notice posted

---

## Maintenance Schedule

| Review                  | Frequency   | Responsible        |
| ----------------------- | ----------- | ------------------ |
| Policy updates          | Quarterly   | Lars Artmann       |
| Dependency audit        | Monthly     | Automated          |
| Security review         | Continuous  | Automated + manual |
| Performance benchmarks  | Per release | CI                 |
| Documentation freshness | Per release | Maintainers        |

---

## Notes

- This checklist should be updated as tasks are completed
- Mark items as ✅ when done, ⏳ when pending, ⚠️ when requires manual/external action
- Add new tasks as policy evolves
- Review quarterly for relevance
- **Status Legend**: ✅ Done | ⏳ Pending | ⚠️ Manual/External

---

_Last updated: 2026-03-27_
