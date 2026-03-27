---
name: Breaking Change Proposal
about: Propose a breaking API change (requires major version bump)
title: "[BREAKING] "
labels: breaking-change
assignees: ""
---

## Summary

Brief summary of the proposed breaking change.

## Motivation

Why is this breaking change necessary?
What problem does it solve that cannot be addressed with backward-compatible changes?

## Current State

Describe the current API or behavior that needs to change.

```go
// Current API
package example

// Current implementation
```

## Proposed Change

Describe the new API or behavior.

```go
// Proposed new API
package example

// New implementation
```

## Migration Path

How will users migrate from the old API to the new one?

### Before

```go
// Current usage that will break
```

### After

```go
// New usage pattern
```

## Deprecation Strategy

- [ ] Deprecation period of at least 2 minor versions
- [ ] Migration guide provided
- [ ] Deprecation warnings added
- [ ] Documentation updated

## Impact Assessment

| Aspect               | Impact                              |
| -------------------- | ----------------------------------- |
| **Breaking Changes** | List specific breaking changes      |
| **Affected Users**   | Estimate how many users affected    |
| **Migration Effort** | Low / Medium / High                 |
| **Benefit**          | Describe the benefit of this change |

## Compatibility

- [ ] This change is backward-compatible (use feature request template instead)
- [ ] This change requires a major version bump (v0 → v1, or v1 → v2)
- [ ] This change can be mitigated with a deprecation period

## Timeline

| Phase                 | Target Date |
| --------------------- | ----------- |
| RFC Discussion        |             |
| Deprecation Notice    |             |
| Major Version Release |             |

## References

- Related issues: #
- Related PRs: #
- External references:

## Checklist

- [ ] I have read the [Breaking Changes Policy](../POLICY.md#2-breaking-changes-policy)
- [ ] I have considered backward-compatible alternatives
- [ ] I have provided a migration guide
- [ ] I have estimated the impact on users
- [ ] I understand this requires a major version bump
- [ ] I am willing to help implement this change

---

**Note**: Breaking changes require significant community discussion and a major version release. Please be prepared for a thorough review process.
