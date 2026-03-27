## Summary

Brief description of the changes in this PR.

## Related Issues

Fixes # (issue number)
Updates # (issue number)

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Code refactoring
- [ ] Dependency update
- [ ] Security fix

## Affected Package(s)

<!-- Check all that apply -->

- [ ] `actor/`
- [ ] `bounded/`
- [ ] `datapoint/`
- [ ] `enums/`
- [ ] `id/`
- [ ] `locale/`
- [ ] `money/`
- [ ] `nanoid/`
- [ ] `pkg/errors/`
- [ ] `temporal/`
- [ ] `types/`
- [ ] Documentation
- [ ] Build/CI

## Changes Made

### Detailed Description

Describe your changes in detail:

### API Changes (if any)

| Before  | After   |
| ------- | ------- |
| Old API | New API |

### Migration Guide (if breaking)

```go
// Before
oldCode()

// After
newCode()
```

## Testing

- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Benchmarks added (if performance-related)
- [ ] Manual testing performed
- [ ] `just check` passes locally

### Test Coverage

- [ ] Coverage maintained or improved
- [ ] All new code paths covered

## Documentation

- [ ] README updated (if needed)
- [ ] CHANGELOG updated
- [ ] POLICY.md updated (if needed)
- [ ] Package documentation updated
- [ ] Code comments added
- [ ] Examples updated (if needed)

## Quality Gates

### Code Quality

- [ ] Code follows project style guidelines
- [ ] No linting errors (`just lint`)
- [ ] No vet errors (`go vet ./...`)
- [ ] Code is properly formatted (`go fmt`)

### Commit Messages

- [ ] Commits follow conventional commit format
- [ ] Commit messages are clear and descriptive

Example:

```
feat(id): add CompareOrZero convenience method

Adds a new CompareOrZero method that returns 0 instead of an error
when comparing IDs with non-ordered types.

Fixes #123
```

## Checklist

- [ ] I have read the [CONTRIBUTING guidelines](../POLICY.md#6-contribution-guidelines)
- [ ] I have read the [Breaking Changes Policy](../POLICY.md#2-breaking-changes-policy) (if applicable)
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Breaking Changes

- [ ] This PR contains no breaking changes
- [ ] This PR contains breaking changes (explain below)

### Breaking Changes Explanation

If this PR contains breaking changes, describe them here:

## Additional Notes

Any additional information or context for reviewers:

---

## Reviewer Checklist

For reviewers:

- [ ] Code is clear and maintainable
- [ ] Tests are comprehensive
- [ ] Documentation is accurate
- [ ] Breaking changes are justified (if any)
- [ ] Performance impact is acceptable (if applicable)
