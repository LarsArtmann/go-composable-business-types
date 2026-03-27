# Security Policy

> Security policies and vulnerability reporting for go-composable-business-types

---

## Supported Versions

Security updates are provided for the following versions:

| Version             | Status         | Security Support    |
| ------------------- | -------------- | ------------------- |
| Latest minor (v0.x) | ✅ Active      | Full support        |
| Previous minor      | ⚠️ Maintenance | Security fixes only |
| Older versions      | ❌ End-of-life | No support          |

---

## Reporting a Vulnerability

**Please DO NOT open public GitHub issues for security vulnerabilities.**

Instead, report security issues privately:

### Contact

- **Email**: security@lars.software
- **Subject**: `[SECURITY] go-composable-business-types: <brief description>`

### Required Information

Please include:

1. **Description**: Clear description of the vulnerability
2. **Affected versions**: Which versions are affected?
3. **Affected packages**: Which package(s) are impacted?
4. **Steps to reproduce**: Detailed steps to reproduce
5. **Proof of concept**: Code demonstrating the issue (if possible)
6. **Impact assessment**: What could an attacker do?
7. **Suggested fix**: If you have a proposed fix
8. **Your contact**: For follow-up questions

### Example Report

```
To: security@lars.software
Subject: [SECURITY] go-composable-business-types: ID validation bypass

Description:
The ID type doesn't properly validate input in certain edge cases,
allowing invalid identifiers to be created.

Affected Versions:
- v0.1.0 and later

Affected Packages:
- id/

Steps to Reproduce:
1. Create an ID with empty string
2. ...

Impact:
An attacker could bypass validation by...

Suggested Fix:
Add explicit validation in NewID constructor.

Contact: your-email@example.com
```

---

## Response Timeline

| Phase                 | Timeline              | Actions                               |
| --------------------- | --------------------- | ------------------------------------- |
| **Acknowledgment**    | Within 48 hours       | Confirm receipt, assign internal ID   |
| **Assessment**        | Within 7 days         | Severity classification, reproduction |
| **Fix Development**   | Within 30 days        | Develop and test fix                  |
| **Pre-announcement**  | 7 days before release | Notify affected users                 |
| **Public Disclosure** | With release          | CVE, fix, and post-mortem             |

**Note**: Complex vulnerabilities may require additional time. We will communicate timeline updates.

---

## Severity Classification

| Severity     | Definition                                                    | Response               |
| ------------ | ------------------------------------------------------------- | ---------------------- |
| **Critical** | Remote code execution; complete system compromise             | 48 hours               |
| **High**     | Authentication bypass; privilege escalation; data corruption  | 7 days                 |
| **Medium**   | Information disclosure; denial of service; partial compromise | 30 days                |
| **Low**      | Minor security improvements; defense in depth                 | Next scheduled release |

---

## Security Best Practices

### For Users

1. **Keep dependencies updated**: Use the latest version
2. **Validate inputs**: Never trust user input
3. **Use safe defaults**: Enable all validation options
4. **Monitor advisories**: Watch for security announcements

### For Contributors

1. **No hardcoded secrets**: Never commit credentials
2. **Input validation**: Validate all inputs at boundaries
3. **Safe defaults**: Prefer secure-by-default designs
4. **No panics on user input**: Return errors, don't panic
5. **Dependency scanning**: Check dependencies for vulnerabilities

---

## Security Features

This library implements:

| Feature                          | Implementation                                        |
| -------------------------------- | ----------------------------------------------------- |
| **Cryptographic IDs**            | `sixafter/nanoid` for FIPS 140-2 compliant IDs        |
| **Input validation**             | All constructors validate and return errors           |
| **No panics**                    | Errors returned instead of panicking on invalid input |
| **Bounds checking**              | Numeric types validate ranges                         |
| **SQL safety**                   | Scanner/Valuer use parameterized queries              |
| **No logging of sensitive data** | IDs and values are never logged internally            |

---

## Security Audit History

| Date | Auditor | Scope | Results       |
| ---- | ------- | ----- | ------------- |
| -    | -       | -     | No audits yet |

---

## Acknowledgments

We thank the following security researchers for responsibly disclosing vulnerabilities:

| Researcher | Issue | Date |
| ---------- | ----- | ---- |
| -          | -     | -    |

---

## Policy Updates

This security policy may be updated. Changes will be announced via:

- GitHub Security Advisories
- CHANGELOG.md
- GitHub Releases

---

## Questions?

For security-related questions that are NOT vulnerabilities:

- Open a [GitHub Discussion](https://github.com/larsartmann/go-composable-business-types/discussions)
- Email: security@lars.software (include `[QUESTION]` in subject)

---

_Last updated: 2026-03-27_
