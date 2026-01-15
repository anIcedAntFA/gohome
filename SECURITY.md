# Security Policy

## Supported Versions

We actively support the following versions of **gohome** with security updates:

| Version | Supported          | Status      |
| ------- | ------------------ | ----------- |
| 1.x.x   | :white_check_mark: | Current     |
| < 1.0   | :x:                | End of Life |

**Recommendation:** Always use the latest stable release for the most secure experience.

## Reporting a Vulnerability

We take security issues seriously. If you discover a security vulnerability, please follow responsible disclosure:

### 🔒 Private Disclosure (Preferred)

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, report privately via:

1. **GitHub Security Advisories**
   - Go to [Security tab](https://github.com/anIcedAntFA/gohome/security)
   - Click "Report a vulnerability"
   - Fill out the form with details

2. **Email** (if GitHub is not available)
   - Send to: **security@ngockhoi96.dev**
   - Subject: `[SECURITY] gohome - [Brief Description]`
   - Include details below

### 📋 What to Include

Please provide the following information:

- **Description** - Clear explanation of the vulnerability
- **Impact** - What can an attacker do? (RCE, data leak, DoS, etc.)
- **Steps to Reproduce** - Detailed instructions to trigger the issue
- **Affected Versions** - Which versions are vulnerable?
- **Proof of Concept** - Code, screenshots, or logs demonstrating the issue
- **Suggested Fix** (optional) - If you have ideas for mitigation

### Example Report

```
Subject: [SECURITY] gohome - Command Injection in Git Operations

Description:
The git.Client.GetLogs() function does not properly sanitize the --author flag,
allowing command injection when processing malicious git author names.

Impact:
Remote Code Execution (RCE) - An attacker could execute arbitrary commands on
the victim's system by manipulating git configuration.

Steps to Reproduce:
1. Set malicious git author: git config user.name "$(whoami)"
2. Run: gohome --author '$(whoami)'
3. Observe command execution

Affected Versions: v1.0.0 - v1.1.5

Proof of Concept:
[Attached screenshot showing command execution]

Suggested Fix:
Implement strict regex validation on author input before passing to git commands.
See existing sanitization in git/client.go:sanitizeInput()
```

### ⏱️ Response Timeline

We aim to respond to security reports within:

- **Initial Response:** 48 hours
- **Triage & Assessment:** 7 days
- **Fix Development:** 14-30 days (depending on severity)
- **Public Disclosure:** After patch release

### 🎯 Severity Levels

We categorize vulnerabilities using CVSS scoring:

| Severity | CVSS Score | Response Time | Example                          |
| -------- | ---------- | ------------- | -------------------------------- |
| Critical | 9.0-10.0   | 48 hours      | Remote Code Execution            |
| High     | 7.0-8.9    | 7 days        | Privilege Escalation             |
| Medium   | 4.0-6.9    | 14 days       | Information Disclosure           |
| Low      | 0.1-3.9    | 30 days       | Minor config exposure            |

### 🔐 Disclosure Process

1. **Private Notification** - We notify you when we confirm the issue
2. **Fix Development** - We develop and test the patch
3. **Security Advisory** - We create a GitHub Security Advisory (draft)
4. **Coordinated Release**
   - Patch released in new version
   - Security advisory published
   - CVE assigned (if applicable)
   - Credits given to reporter (unless anonymous preferred)

## Security Best Practices

### For Users

- **Keep Updated** - Always use the latest version
- **Verify Downloads** - Check SHA256 checksums from GitHub releases
- **Use Official Sources** - Download only from:
  - GitHub Releases
  - Official install scripts (`get.ngockhoi96.dev/gohome`)
  - NPM (`@ngockhoi96/gohome`)
  - AUR (Arch Linux)

### For Contributors

- **Input Validation** - Always sanitize user input
  - See `internal/git/client.go:sanitizeInput()` for reference
  - Use regex to allow only safe characters
- **Path Traversal** - Validate file paths
  - Use `filepath.Clean()` to normalize paths
  - Check `internal/config/config.go:validateConfigPath()`
- **Command Injection** - Never pass unsanitized input to shell
  - Use `exec.Command()` with separate arguments
  - Avoid `bash -c` or similar constructs
- **Dependency Security** - Keep dependencies updated
  - Run `go mod tidy` regularly
  - Monitor Dependabot alerts

## Known Security Considerations

### Current Mitigations

1. **Command Injection Prevention**
   - All git command arguments sanitized via regex
   - Location: `internal/git/client.go:sanitizeInput()`
   - Pattern: `[^a-zA-Z0-9\s._@-]+` removed

2. **Path Traversal Protection**
   - Config file paths validated
   - `filepath.Clean()` used to normalize paths
   - Location: `internal/config/config.go:validateConfigPath()`

3. **No Network Operations**
   - gohome is fully offline, no outbound connections
   - No telemetry or analytics

### Out of Scope

The following are **NOT** considered security issues:

- Issues requiring physical access to the machine
- Denial of Service via local resource exhaustion
- Issues in dependencies (report to upstream)
- Social engineering attacks
- Output formatting issues that don't expose sensitive data

## Security Updates

Security patches are released as:

- **Patch versions** (1.0.x) - For minor/low severity
- **Minor versions** (1.x.0) - For medium/high severity
- **Emergency releases** - For critical vulnerabilities (within 48h)

Subscribe to:
- [GitHub Releases](https://github.com/anIcedAntFA/gohome/releases)
- [GitHub Watch](https://github.com/anIcedAntFA/gohome) → Security alerts only

## Attribution

We appreciate responsible disclosure and will credit researchers:

- **GitHub Security Advisories** - Listed as collaborators
- **Release Notes** - Acknowledged in CHANGELOG
- **Hall of Fame** - Recognized in README (optional)

You may choose to remain anonymous.

## Contact

- **Security Issues:** security@ngockhoi96.dev or [GitHub Security Advisories](https://github.com/anIcedAntFA/gohome/security)
- **General Questions:** [GitHub Discussions](https://github.com/anIcedAntFA/gohome/discussions)
- **Non-Security Bugs:** [GitHub Issues](https://github.com/anIcedAntFA/gohome/issues)

Thank you for helping keep **gohome** secure! 🔒
