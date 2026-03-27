# Support

> How to get help with go-composable-business-types

---

## Quick Links

| Resource          | Link                                                                                          |
| ----------------- | --------------------------------------------------------------------------------------------- |
| **Documentation** | [README.md](./README.md)                                                                      |
| **Changelog**     | [CHANGELOG.md](./CHANGELOG.md)                                                                |
| **Policy**        | [POLICY.md](./POLICY.md)                                                                      |
| **Issues**        | [GitHub Issues](https://github.com/larsartmann/go-composable-business-types/issues)           |
| **Discussions**   | [GitHub Discussions](https://github.com/larsartmann/go-composable-business-types/discussions) |

---

## Getting Help

### 1. Check the Documentation

Most questions are answered in:

- **[README.md](./README.md)** — Usage examples and API overview
- **[POLICY.md](./POLICY.md)** — Versioning, breaking changes, contribution guidelines
- **Package READMEs** — `id/README.md` for detailed ID type documentation
- **Examples** — See `examples/` directory for working code

### 2. Search Existing Issues

Before opening a new issue:

1. Search [GitHub Issues](https://github.com/larsartmann/go-composable-business-types/issues) for similar problems
2. Check [closed issues](https://github.com/larsartmann/go-composable-business-types/issues?q=is%3Aissue+is%3Aclosed) — your question may have been answered

### 3. Ask in Discussions

For questions, help, or general discussion:

- **[GitHub Discussions](https://github.com/larsartmann/go-composable-business-types/discussions)** — Community Q&A
- Use categories: `Q&A`, `Show and tell`, `Ideas`

### 4. Open an Issue

For bug reports or feature requests:

- **[Bug Report](https://github.com/larsartmann/go-composable-business-types/issues/new?template=bug_report.md)** — Something not working?
- **[Feature Request](https://github.com/larsartmann/go-composable-business-types/issues/new?template=feature_request.md)** — Have an idea?

---

## Support Channels

| Channel                | Best For              | Response Time |
| ---------------------- | --------------------- | ------------- |
| **Documentation**      | Self-service help     | Immediate     |
| **GitHub Discussions** | Questions, usage help | 1-2 weeks     |
| **GitHub Issues**      | Bug reports, features | 1 week        |
| **Security Email**     | Vulnerabilities       | 48 hours      |
| **Enterprise Support** | Commercial assistance | 24 hours      |

---

## Reporting Issues

### Bug Reports

Include in your bug report:

````markdown
**Package:** Which package(s) are affected?
**Version:** What version are you using? (`go list -m`)
**Go Version:** `go version`
**OS/Arch:** e.g., `darwin/arm64`, `linux/amd64`

**What happened?**
Describe the bug clearly.

**What did you expect?**
Describe what should have happened.

**Minimal reproduction:**

```go
// Provide a minimal code example that reproduces the issue
package main

import (
    "github.com/larsartmann/go-composable-business-types/id"
)

func main() {
    // Your code here
}
```
````

**Additional context:**
Any error messages, stack traces, or relevant logs.

````

### Feature Requests

Include in your feature request:

```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Any alternative solutions or features you've considered.

**Additional context**
Any other context, use cases, or examples.
````

---

## Version Support

| Version       | Status             | Support Level              |
| ------------- | ------------------ | -------------------------- |
| v0.x (latest) | Active development | Full support               |
| v0.x (older)  | Maintenance        | Bug fixes only             |
| v1.x          | Planned            | Full support when released |

---

## Security Issues

**Do NOT open public issues for security vulnerabilities.**

Email: **security@lars.software**

Subject: `[SECURITY] go-composable-business-types: <brief description>`

See [POLICY.md#security-policy](./POLICY.md#5-security-policy) for details.

---

## Enterprise Support

For organizations requiring:

- Priority bug fixes
- Custom type development
- Migration assistance
- Architecture review
- Training and workshops

Contact: **enterprise@lars.software**

---

## Community

### Contributing

Want to contribute? See [POLICY.md#contribution-guidelines](./POLICY.md#6-contribution-guidelines):

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `just check`
5. Submit a pull request

### Code of Conduct

This project follows a code of conduct:

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Respect differing viewpoints
- Prioritize community benefit

---

## FAQ

### Q: What Go versions are supported?

**A:** Go 1.26+ is required. See [POLICY.md#go-version-support](./POLICY.md#14-go-version-support).

### Q: How do I upgrade between versions?

**A:** See [CHANGELOG.md](./CHANGELOG.md) for migration notes. Breaking changes require a major version bump.

### Q: Is there a changelog?

**A:** Yes, see [CHANGELOG.md](./CHANGELOG.md).

### Q: How do I report a breaking change?

**A:** Open an issue with the `breaking-change` label. See [POLICY.md#breaking-change-process](./POLICY.md#23-breaking-change-process).

### Q: Can I use this in production?

**A:** The library is pre-v1.0.0. While stable, API changes may occur. Pin to specific versions in production.

### Q: How do I get notified of new releases?

**A:** Watch the repository on GitHub for release notifications.

### Q: Where can I find examples?

**A:** See the `examples/` directory and [README.md](./README.md).

---

## Maintenance

This library is maintained by **Lars Artmann**.

- **Active development**: New features, bug fixes
- **Response time**: Issues within 1 week, security within 48 hours
- **Stability**: Conservative API design; minimal breaking changes

---

_Last updated: 2026-03-27_
