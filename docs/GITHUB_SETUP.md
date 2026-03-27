# GitHub Repository Setup Guide

> Manual configuration required for repository settings

This guide covers manual setup tasks that cannot be automated via code.

---

## 1. Required Labels

Create these labels in GitHub (Settings → General → Labels):

| Label              | Color     | Description                                       |
| ------------------ | --------- | ------------------------------------------------- |
| `breaking-change`  | `#b60205` | API breaking changes requiring major version bump |
| `security`         | `#d93f0b` | Security-related issues and fixes                 |
| `deprecation`      | `#f9d0c4` | Deprecation notices and removal tracking          |
| `documentation`    | `#0052cc` | Documentation updates                             |
| `good-first-issue` | `#7057ff` | Good for newcomers                                |
| `help-wanted`      | `#008672` | Extra attention needed                            |
| `dependencies`     | `#0366d6` | Dependency updates (Dependabot)                   |
| `automated`        | `#ededed` | Automated PRs (Dependabot)                        |
| `github-actions`   | `#000000` | GitHub Actions related                            |

### Creating Labels via GitHub CLI

```bash
gh label create "breaking-change" --color "b60205" --description "API breaking changes requiring major version bump"
gh label create "security" --color "d93f0b" --description "Security-related issues and fixes"
gh label create "deprecation" --color "f9d0c4" --description "Deprecation notices and removal tracking"
gh label create "documentation" --color "0052cc" --description "Documentation updates"
gh label create "good-first-issue" --color "7057ff" --description "Good for newcomers"
gh label create "help-wanted" --color "008672" --description "Extra attention needed"
gh label create "dependencies" --color "0366d6" --description "Dependency updates"
gh label create "automated" --color "ededed" --description "Automated PRs"
gh label create "github-actions" --color "000000" --description "GitHub Actions related"
```

---

## 2. Branch Protection Rules

Configure branch protection for `main` and `master` (Settings → Branches → Add rule):

### Required Settings

| Setting                                                          | Value                                  |
| ---------------------------------------------------------------- | -------------------------------------- |
| **Branch name pattern**                                          | `main` (repeat for `master`)           |
| Require a pull request before merging                            | ✓                                      |
| Required approvals                                               | 1                                      |
| Dismiss stale pull request approvals when new commits are pushed | ✓                                      |
| Require review from Code Owners                                  | ✓ (if CODEOWNERS exists)               |
| Require status checks to pass before merging                     | ✓                                      |
| Require branches to be up to date before merging                 | ✓                                      |
| Status checks required                                           | `test`, `lint`, `security`, `generate` |
| Require conversation resolution before merging                   | ✓                                      |
| Require linear history                                           | ✓                                      |
| Include administrators                                           | ✓                                      |
| Restrict who can push to matching branches                       | ✓ (add maintainers)                    |

### GitHub CLI Setup

```bash
# Note: Branch protection via CLI requires GitHub App or PAT with admin permissions
# Use the GitHub UI for initial setup

# View current rules
gh api repos/:owner/:repo/branches/main/protection
```

---

## 3. Repository Settings

### General Settings (Settings → General)

| Setting                                       | Value          |
| --------------------------------------------- | -------------- |
| **Features**                                  |                |
| Wiki                                          | ✗ (use docs/)  |
| Issues                                        | ✓              |
| Sponsorships                                  | ✗              |
| Preserve this repository                      | ✗              |
| Discussions                                   | ✓              |
| Projects                                      | ✗ (use Issues) |
| **Pull Requests**                             |                |
| Allow merge commits                           | ✗              |
| Allow squash merging                          | ✓              |
| Allow rebase merging                          | ✓              |
| Always suggest updating pull request branches | ✓              |
| Allow auto-merge                              | ✓              |
| Automatically delete head branches            | ✓              |

### Security Settings

| Setting                         | Value          |
| ------------------------------- | -------------- |
| **Security**                    |                |
| Private vulnerability reporting | ✓              |
| Dependency graph                | ✓              |
| Dependabot alerts               | ✓              |
| Dependabot security updates     | ✓              |
| Dependabot version updates      | ✓ (configured) |
| Code scanning (CodeQL)          | ✓ (optional)   |
| Secret scanning                 | ✓              |
| Push protection                 | ✓              |

---

## 4. GitHub Actions Permissions

Settings → Actions → General:

| Setting                                                  | Value |
| -------------------------------------------------------- | ----- |
| **Actions permissions**                                  |       |
| Allow all actions and reusable workflows                 | ✓     |
| **Workflow permissions**                                 |       |
| Read repository contents and packages permissions        | ✗     |
| Read and write permissions                               | ✓     |
| Allow GitHub Actions to create and approve pull requests | ✗     |

---

## 5. Email Configuration

### Security Email

Set up `security@lars.software` to forward to your preferred email:

1. Configure DNS MX records for `lars.software`
2. Set up email forwarding to your personal/work email
3. Test by sending an email to verify delivery

### GitHub Security Contact

Update SECURITY.md if email changes.

---

## 6. CODEOWNERS (Optional)

Create `.github/CODEOWNERS` if you want automatic reviewer assignment:

```
# Default owners
* @larsartmann

# Package maintainers
/id/ @larsartmann
/actor/ @larsartmann
/datapoint/ @larsartmann
```

---

## 7. Verification Checklist

After setup, verify:

- [ ] All labels created and visible
- [ ] Branch protection rules active (try to push directly to main)
- [ ] PR requires approval (open test PR)
- [ ] Status checks required (verify CI runs)
- [ ] Dependabot creates first PR (check Insights → Dependency graph)
- [ ] Security email forwards correctly (send test email)
- [ ] Discussions enabled (check Discussions tab)
- [ ] Vulnerability reporting enabled (check Security tab)

---

## 8. Automation Commands Summary

```bash
# Create all labels
gh label create "breaking-change" --color "b60205" --description "API breaking changes"
gh label create "security" --color "d93f0b" --description "Security-related issues"
gh label create "deprecation" --color "f9d0c4" --description "Deprecation notices"
gh label create "documentation" --color "0052cc" --description "Documentation updates"
gh label create "good-first-issue" --color "7057ff" --description "Good for newcomers"
gh label create "help-wanted" --color "008672" --description "Extra attention needed"
gh label create "dependencies" --color "0366d6" --description "Dependency updates"
gh label create "automated" --color "ededed" --description "Automated PRs"
gh label create "github-actions" --color "000000" --description "GitHub Actions"

# Verify labels
gh label list

# Check branch protection (requires admin token)
gh api repos/larsartmann/go-composable-business-types/branches/main/protection
```

---

_Last updated: 2026-03-27_
