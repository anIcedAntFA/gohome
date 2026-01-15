# Contributing to gohome

Thank you for your interest in contributing to **gohome**! 🎉

We welcome contributions of all kinds: bug reports, feature requests, documentation improvements, and code contributions.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Convention](#commit-convention)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## 🤝 Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please be respectful and constructive in all interactions.

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/dl/)
- **Git** - For version control
- **golangci-lint** - For code linting
- **Make** - For build automation

### Setup Development Environment

```bash
# 1. Fork the repository on GitHub

# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/gohome.git
cd gohome

# 3. Add upstream remote
git remote add upstream https://github.com/anIcedAntFA/gohome.git

# 4. Install dependencies
go mod tidy

# 5. Verify setup
make build
make test
make lint
```

## 🔄 Development Workflow

### Working on Issues

1. **Find or create an issue** - Check existing issues or create a new one
2. **Comment on the issue** - Let others know you're working on it
3. **Create a branch** - Use descriptive names:
   ```bash
   git checkout -b feat/add-json-export
   git checkout -b fix/clipboard-wsl2
   git checkout -b docs/update-readme
   ```

### Running Locally

```bash
# Run the CLI directly
go run cmd/gohome/main.go

# Build and test
make build
./bin/gohome --help

# Run tests
go test -v ./...
make test

# Run linter
make lint
```

### Testing

- **Unit Tests** - Test individual functions/packages
- **Integration Tests** - Test end-to-end workflows
- All tests must pass before merging

Example test command:
```bash
# Run all tests
go test -v ./...

# Run specific package tests
go test -v ./internal/scanner/

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 Coding Standards

### Code Style

- Follow **[Effective Go](https://go.dev/doc/effective_go)** guidelines
- Use **gofmt** / **gofumpt** for formatting
- Run **golangci-lint** before committing
- Keep functions small and focused (< 50 lines when possible)
- Add comments for exported functions and complex logic

### Project Structure

```
gohome/
├── cmd/gohome/        # Main application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration loading/parsing
│   ├── entity/        # Data models (Commit, Task)
│   ├── git/           # Git client wrapper
│   ├── parser/        # Conventional Commits parser
│   ├── renderer/      # Output formatters (text/table)
│   ├── scanner/       # Repository discovery
│   ├── spinner/       # Terminal spinner UI
│   ├── sys/           # System utilities (clipboard)
│   └── version/       # Version info (injected at build)
├── docs/              # Documentation and guides
├── scripts/           # Installation scripts
└── npm-package/       # NPM distribution wrapper
```

### Key Patterns

- **Constructor functions**: Use `New<Type>()` pattern (e.g., `git.NewClient()`)
- **Error handling**: Return errors, use `log.Fatal()` only in main
- **Input sanitization**: Always sanitize git command arguments (see `git/client.go`)
- **Config precedence**: CLI flags override JSON file values

### Security

- **Never** execute unsanitized user input
- Use regex validation for git arguments
- Sanitize file paths to prevent traversal attacks
- See `internal/git/client.go` for sanitization examples

## 📜 Commit Convention

We follow **[Conventional Commits](https://www.conventionalcommits.org/)** with emojis:

### Format

```
<emoji> <type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

- `✨ feat` - New feature
- `🐛 fix` - Bug fix
- `📝 docs` - Documentation only
- `🎨 style` - Code style (formatting, no logic change)
- `♻️ refactor` - Code restructure (no feature/bug change)
- `⚡ perf` - Performance improvement
- `✅ test` - Add/update tests
- `🔨 chore` - Build process, dependencies, tools
- `🔧 ci` - CI/CD configuration
- `🚀 release` - Version bump, release prep

### Scope (optional)

- `scanner` - Repository discovery
- `parser` - Commit parsing
- `config` - Configuration
- `renderer` - Output formatting
- `git` - Git client
- `cli` - CLI flags/commands
- `docs` - Documentation
- `deps` - Dependencies

### Examples

```bash
✨ feat(scanner): add recursive scanning up to 2 levels deep
🐛 fix(clipboard): add WSL2 support for clipboard operations
📝 docs(readme): update installation instructions for AUR
🔨 chore(deps): upgrade golangci-lint to v1.55
♻️ refactor(config): simplify time period merging logic
⚡ perf(scanner): implement concurrent repo scanning
```

### Writing Good Commits

- Use imperative mood: "add" not "added" or "adds"
- Keep subject line under 72 characters
- Capitalize subject line
- No period at the end of subject
- Separate subject from body with blank line
- Explain **what** and **why**, not **how**

## 🔁 Pull Request Process

### Before Submitting

1. **Sync with upstream**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run checks**
   ```bash
   make lint    # Must pass
   make test    # Must pass
   make build   # Must succeed
   ```

3. **Update documentation** - If you changed:
   - CLI flags → Update `README.md` flags table
   - Config options → Update config example
   - Architecture → Update `.github/copilot-instructions.md`

4. **Update CHANGELOG.md** - Add entry under `[Unreleased]`

### Submitting PR

1. **Push to your fork**
   ```bash
   git push origin feat/your-feature
   ```

2. **Create Pull Request** on GitHub
   - Use PR template (auto-populated)
   - Link related issues with "Closes #123"
   - Add screenshots/demos for UI changes
   - Fill out all checklist items

3. **Address review feedback**
   - Make requested changes
   - Push new commits (don't force-push during review)
   - Re-request review when ready

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Self-reviewed code and added comments
- [ ] Updated documentation (README, inline docs)
- [ ] No new warnings or errors
- [ ] Added/updated tests for changes
- [ ] All tests pass locally
- [ ] Ran `gofumpt -l -w .` for formatting
- [ ] Ran `golangci-lint run` with no issues
- [ ] Updated CHANGELOG.md

### After Merge

- Delete your feature branch (GitHub will prompt)
- Sync your fork:
  ```bash
  git checkout main
  git fetch upstream
  git merge upstream/main
  git push origin main
  ```

## 🐛 Reporting Bugs

### Before Reporting

- Check [existing issues](https://github.com/anIcedAntFA/gohome/issues)
- Try the latest version
- Collect debug information

### Bug Report Template

Use GitHub issue templates - they guide you through providing:

- **Description** - Clear summary of the bug
- **Steps to Reproduce** - Exact steps to trigger the bug
- **Expected Behavior** - What should happen
- **Actual Behavior** - What actually happens
- **Environment** - OS, Go version, gohome version
- **Logs/Screenshots** - Terminal output, error messages

## 💡 Suggesting Features

### Feature Request Template

- **Problem Statement** - What problem does this solve?
- **Proposed Solution** - How should it work?
- **Alternatives Considered** - Other approaches you thought about
- **Additional Context** - Screenshots, examples from other tools

### Feature Discussion

- Open an issue **before** implementing large features
- Discuss approach with maintainers
- Check [ROADMAP.md](ROADMAP.md) for planned features

## 📚 Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)

## 🙏 Recognition

Contributors are recognized in:
- GitHub Contributors list
- Release notes (for significant contributions)
- CHANGELOG.md

## ❓ Questions?

- Open a [GitHub Discussion](https://github.com/anIcedAntFA/gohome/discussions)
- Comment on relevant issues
- Check [README.md](README.md) for documentation

Thank you for contributing to **gohome**! 🚀
