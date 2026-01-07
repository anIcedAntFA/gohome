---
name: ⭐ Feature Request
about: Propose a great feature for gohome CLI
title: '[FEATURE] '
labels: enhancement
assignees: ''
---

## 🌟 Feature Description

Provide a clear and concise description of the feature you want to suggest.

**Example:** "I'd love for the `gohome` command to support filtering commits by specific file types (e.g., only show commits that modified .go files)."

## 🚀 What's the Benefit?

Describe the benefit of implementing the feature from the user's perspective.

**Example:** "This feature would help developers focus on specific areas of their work. For instance, frontend developers could filter for .css/.html changes, while backend developers could focus on .go/.sql files."

## 🎯 Problem Statement

Is your feature request related to a problem? Please describe.

**Example:** "I'm always frustrated when I have to manually scan through all commits to find ones related to specific file types. This is time-consuming when preparing reports for specialized team meetings."

## 🛠️ Proposed Solution

Propose a solution if you have one in mind. If not, just leave this section blank.

**Example:** "Add a new flag `--file-ext` or `-e` that accepts file extensions:"

```bash
gohome -d 7 --file-ext .go --file-ext .mod
gohome -d 7 -e .ts -e .tsx
```

## 💡 Possible Alternatives

If applicable, describe any alternative solutions or features you've considered.

**Example:**
- Use regex patterns: `gohome -d 7 --pattern '.*\.go$'`
- Add include/exclude filters
- Support gitignore-style patterns

## 🗺️ Example Scenario

Describe an example where your proposed feature would be beneficial.

**Example:**

**Current workflow:**
1. Run `gohome -d 7`
2. Manually scan through 100+ commits
3. Copy only Go-related commits
4. Paste into daily standup notes
5. Takes ~15 minutes

**With this feature:**
1. Run `gohome -d 7 --file-ext .go --copy`
2. Paste into daily standup notes
3. Takes ~30 seconds ⚡

## 🖼️ Mockups or Examples

If applicable, add mockups, screenshots, or examples from other tools that demonstrate what you're looking for.

**Example output:**
```
🗓️ Period: 7 days ago
✓ Found 5 repositories

📦 backend-api (Go files only)
├─ feat(auth): implement JWT validation
├─ fix(db): resolve connection pool leak
└─ refactor(api): clean up handler structure

✓ Filtered 23/45 commits
```

## 📚 Related Documentation

Add any relevant documents or links that will help us understand the feature request better:

- Similar features in other tools (git-standup, github CLI, etc.)
- Stack Overflow discussions
- Community requests or discussions

## 🌞 Anything else you would like to add?

Add any other context or screenshots about the feature request here!

- Would you be willing to contribute a PR for this feature?
- Any performance concerns?
- Backward compatibility considerations?

## ✅ Acceptance Criteria

If you have specific requirements, list them here:

- [ ] Support multiple file extensions via repeated flags
- [ ] Work with all output formats (text, table, markdown)
- [ ] Update documentation and help text
- [ ] Add tests for filtering logic
- [ ] Maintain backward compatibility
