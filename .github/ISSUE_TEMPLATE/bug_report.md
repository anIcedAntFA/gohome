---
name: 🪲 Bug Report
about: If something isn't working as expected in gohome CLI
title: '[BUG] '
labels: bug
assignees: ''
---

## 🪲 Describe the Bug

Provide a clear explanation of the bug within the framework of gohome CLI.

**Example:** "gohome crashes when attempting to scan repositories with the `gohome -d 7` command."

## ⚡️ Type of Bug

Please select the type of bug you are reporting:

- [ ] Command Error
- [ ] Crash/Error Message
- [ ] Installation Issue
- [ ] Configuration Issue
- [ ] Git Integration Issue
- [ ] Clipboard/Output Issue
- [ ] Documentation Issue
- [ ] Other (Please describe)

## 🔬 Steps to Reproduce

Detail out the steps to reproduce the bug.

**Example:**
1. Run the `gohome -d 7 -f table` command
2. Navigate to a directory with multiple git repositories
3. Observe the crash/error message

## 🔑 Expected Behavior

Explain what you anticipated happening.

**Example:** "The CLI should've scanned all repositories and outputted a formatted table with commits from the last 7 days."

## 🌚 Actual Behavior

Describe what actually occurred given the steps.

**Example:** "gohome crashed and displayed an error message: 'panic: runtime error: invalid memory address'."

## 📸 Screenshots or Logs

Include screenshots to better illustrate the bug, if necessary.

```bash
# Paste error logs or terminal output here
```

## 🖥️ Environment

Please provide details about your environment:

- **OS:** [e.g., macOS 14.0, Ubuntu 22.04, Windows 11]
- **Go Version:** [e.g., 1.20.5] (run `go version`)
- **gohome Version:** [e.g., v1.0.0 or commit hash] (run `gohome --version` if available)
- **Installation Method:** [e.g., `go install`, built from source, binary download]
- **Git Version:** [e.g., 2.39.0] (run `git --version`)
- **Shell:** [e.g., bash, zsh, fish, PowerShell]

## 🧰 Possible Solution

Impart any insight on a potential bug fix, if possible.

## 📝 Additional Context

Add any other context about the problem here. For example:

- Config file contents (`~/.gohome.json`)
- Any special workspace setup
- Repository structure or characteristics
- Related issues or discussions
- Were you using any specific flags or configurations?

## 🌪️ Impact

Provide details on the scale and severity of the bug:

- [ ] 🔴 Critical - Blocks main functionality
- [ ] 🟠 High - Significant impact on user experience
- [ ] 🟡 Medium - Workaround available
- [ ] 🟢 Low - Minor inconvenience

## 📚 Related Documentation

Include any relevant documentation/resources that correlate with the bug:

- Links to README sections
- Related GitHub issues
- Stack Overflow discussions
- External resources
