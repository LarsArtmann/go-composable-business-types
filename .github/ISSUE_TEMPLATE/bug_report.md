---
name: Bug Report
about: Create a report to help us improve
title: "[BUG] "
labels: bug
assignees: ""
---

## Bug Description

A clear and concise description of what the bug is.

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

## Environment

- **Library Version**: <!-- e.g., v0.1.0 -->
- **Go Version**: <!-- Run `go version` -->
- **OS/Architecture**: <!-- e.g., darwin/arm64, linux/amd64 -->

## Steps to Reproduce

1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

## Expected Behavior

A clear and concise description of what you expected to happen.

## Actual Behavior

What actually happened instead.

## Minimal Reproduction

```go
package main

import (
    "fmt"
    "github.com/larsartmann/go-composable-business-types/id"
)

func main() {
    // Minimal code that reproduces the issue
    type UserBrand struct{}
    userID := id.NewID[UserBrand, string]("test")
    fmt.Println(userID)
}
```

## Additional Context

Add any other context about the problem here:

- Error messages
- Stack traces
- Screenshots
- Related issues

## Checklist

- [ ] I have searched existing issues and this is not a duplicate
- [ ] I have provided a minimal reproduction case
- [ ] I have tested with the latest version
- [ ] I have read the [POLICY.md](../POLICY.md)
