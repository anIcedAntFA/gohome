# 🗺️ Product Roadmap

This document outlines the development status and future plans for **gohome** (Git Standup Tool).

## ✅ Phase 1: Core MVP & Foundation (v1.0.0 - v1.2.0) — **COMPLETED**

**Goal:** Deliver a stable, production-ready CLI tool with essential features and professional distribution pipeline.

**Status:** ✅ **COMPLETED** - Phase 1 is complete. The tool is stable, tested, and ready for daily use with multiple installation methods.

### ✅ Completed Features

**Core Functionality:**
- [x] **Git Integration:** Auto-scan directories for `.git` folders with recursive scanning
- [x] **Configurable Depth:** `--max-depth` flag for nested structures like `github.com/{org}/{repo}` (default: 2 levels)
- [x] **Smart Scanner:** Skip nested repos, ignore common directories (`.git`, `.vscode`, `node_modules`)
- [x] **Log Parsing:** Conventional Commits regex parser with emoji detection
- [x] **Multi-Branch Support:** 
  - [x] `--all-branches` flag to include all local branches
  - [x] `--branch <name>` flag to filter by specific branch
- [x] **Configuration System:**
  - [x] JSON config file (`~/.gohome.json`)
  - [x] CLI flags with shorthand aliases
  - [x] `--save` flag to persist settings
  - [x] Auto-detect git author from system
- [x] **Custom Tasks:** Append manual tasks via `-t` flags (meetings, reviews, etc.)

**User Experience:**
- [x] **Output Formats:** Text list and rich table formats
- [x] **Table Styles:** Multiple presets (normal, markdown)
- [x] **Clipboard Integration:** Cross-platform `--copy` support (Linux, macOS, Windows, WSL2)
- [x] **Visual Feedback:** Loading spinner during scan operations
- [x] **Version Info:** `--version` flag with build metadata injection
- [x] **Help System:** Comprehensive help messages with `tabwriter` formatting

**Quality & Testing:**
- [x] **Unit Tests:** Comprehensive coverage for `scanner` package with table-driven tests
- [x] **CI/CD Pipeline:** GitHub Actions for linting, testing, and releases
- [x] **Code Quality:** golangci-lint integration with strict rules
- [x] **Security:** Input sanitization, path validation, security policy

**Distribution & Installation:**
- [x] **GoReleaser:** Automated cross-platform builds (Linux/macOS/Windows, amd64/arm64)
- [x] **Installation Scripts:**
  - [x] `install.sh` for Linux/macOS (curl piping)
  - [x] `install.ps1` for Windows PowerShell
- [x] **Package Distribution:**
  - [x] GitHub Releases (direct binary downloads)
  - [x] AUR (Arch Linux) - Build from source
  - [x] npm (@ngockhoi96/gohome)
  - [x] Go Install (`go install`)

**Documentation:**
- [x] **README.md:** Installation, usage, examples, flags reference
- [x] **CONTRIBUTING.md:** Development setup, coding standards, commit conventions
- [x] **SECURITY.md:** Vulnerability reporting guidelines
- [x] **Git LFS Guide:** Complete guide for LFS-tracked media files
- [x] **Release Documentation:** Versioning, release checklist, AUR setup guides
- [x] **Demo Media:** VHS-generated GIFs for visual documentation

---

## ✅ Phase 2: Architecture Modernization (v1.3.0) — **COMPLETED**

**Goal:** Modernize CLI architecture with industry-standard frameworks (Cobra/Viper), improve code quality, expand test coverage, and enhance developer experience.

**Status:** ✅ **COMPLETED** - Successfully migrated to modern architecture with comprehensive testing and documentation.

### ✅ Completed in v1.3.0

**Architecture Refactoring:**
- [x] **Cobra/Viper Integration:**
  - [x] Sub-command support: `gohome report`, `gohome config`, `gohome version`, `gohome completion`
  - [x] Auto-generated help text and documentation
  - [x] Shell completion (bash, zsh, fish, PowerShell)
  - [x] Better flag inheritance and organization
  - [x] Industry-standard CLI patterns

- [x] **Configuration Management:**
  - [x] Multi-format config: JSON, YAML, TOML support
  - [x] Environment variable binding: `GOHOME_*` prefix
  - [x] Automatic config hierarchy: Flags > Env > Config > Defaults
  - [x] Type-safe configuration access
  - [x] Config subcommands: `list`, `get`, `set`, `reset`

**Testing & Quality Improvements:**
- [x] **Comprehensive Test Suite:**
  - [x] Config commands: 90%+ coverage (list: 100%, get: 100%, set: 90%, reset: 92.9%)
  - [x] Root command: 85.7% coverage
  - [x] Completion command: 90%+ coverage
  - [x] Overall project coverage: **49.6%** (up from 14.6%)

- [x] **Code Quality:**
  - [x] Dependency injection pattern for testability
  - [x] All linting errors resolved (47 issues fixed)
  - [x] Security improvements (command injection prevention)
  - [x] Consistent error handling with emoji prefixes

- [x] **CI/CD Enhancements:**
  - [x] Codecov integration with coverage tracking
  - [x] Automated test execution on PRs
  - [x] NPM prerelease tag handling (beta/alpha/rc)
  - [x] Parallel test conflict resolution

**Documentation:**
- [x] Migration guide (v1.2 → v1.3)
- [x] Comprehensive CLI guide with examples
- [x] Codecov usage guide
- [x] Go testing best practices documentation
- [x] Viper configuration management guide

---

## 🚧 Phase 3: Enhanced Features & Ecosystem (v1.4.x - v2.0.0) — **IN PLANNING**
  - [ ] `gohome config list` - Show current configuration
  - [ ] `gohome config reset` - Reset to defaults
  - [ ] Multiple config profiles (work, personal, etc.)

### 📤 Export & Integration

- [ ] **Export Formats:**
  - [ ] JSON export for programmatic use
  - [ ] Markdown export for documentation
  - [ ] HTML report with styling
  - [ ] CSV for spreadsheet analysis
  - [ ] PDF generation (via external libs)
- [ ] **Template System:**
  - [ ] Custom output templates
  - [ ] Variable interpolation
  - [ ] Conditional rendering
- [ ] **Webhook Integration:**
  - [ ] POST results to custom endpoints
  - [ ] Slack/Discord webhook support
  - [ ] Webhook retry and error handling

---

## 🔮 Phase 3: AI & Interactive Features (v2.x.x) — **FUTURE**

**Goal:** Transform gohome into an intelligent, interactive productivity tool with AI-powered insights.

**Note:** Phase 3 depends on Phase 2 Cobra/Viper architecture for sub-commands and advanced config management.

### 🤖 AI-Powered Features

- [ ] **Smart Summaries:**
  - [ ] Integration with LLMs (OpenAI, Anthropic, Gemini, local models)
  - [ ] Generate concise daily summaries from raw commits
  - [ ] Multiple prompt styles: "Standup", "Changelog", "Executive Summary"
  - [ ] Context-aware suggestions for missing information
- [ ] **Commit Message Enhancement:**
  - [ ] AI suggestions for better commit messages
  - [ ] Automatic categorization of work
  - [ ] Sentiment analysis and productivity insights
- [ ] **Natural Language Interface:**
  - [ ] Query commits using natural language
  - [ ] "Show me all bug fixes from last week"
  - [ ] "What did I work on related to authentication?"

### 🎨 Interactive Mode (TUI)

- [ ] **Terminal UI Framework:**
  - [ ] Implement `charmbracelet/bubbletea` interface
  - [ ] Beautiful, responsive terminal UI
  - [ ] Keyboard navigation and shortcuts
- [ ] **Interactive Features:**
  - [ ] Select/deselect repositories to include
  - [ ] Live filtering and search
  - [ ] Multi-select commits for export
  - [ ] Preview reports before copying
  - [ ] Configuration editor in TUI
- [ ] **Visual Enhancements:**
  - [ ] Syntax highlighting for code diffs
  - [ ] Rich formatting with `lipgloss`
  - [ ] Custom themes and color schemes

### 📊 Advanced Task Management

- [ ] **Recurring Tasks:**
  - [ ] Define daily/weekly recurring tasks in config
  - [ ] Task templates with variables
  - [ ] Task completion tracking
- [ ] **Time Tracking:**
  - [ ] Estimate time spent per commit/task
  - [ ] Daily/weekly time summaries
  - [ ] Integration with time-tracking tools
- [ ] **Task Prioritization:**
  - [ ] Priority levels for tasks
  - [ ] Sort reports by priority
  - [ ] Highlight overdue or urgent items

---

## 📚 Phase 3 Extended: Distribution, Documentation & Marketing (v2.x.x) — **EXPANSION**

**Goal:** Complete distribution channels, establish official documentation, and grow community.

### 📦 Additional Package Managers

**Goal:** Complete package manager coverage for remaining platforms.

- [ ] **Linux Package Managers:**
  - [ ] **APT** (Debian/Ubuntu) - `.deb` packages for Debian-based distros
  - [ ] **RPM** (Fedora/RHEL/openSUSE) - `.rpm` packages for RedHat-based distros
  - [ ] **Snap** (Universal Linux) - Snap package for all distros
- [ ] **macOS Package Managers:**
  - [ ] **MacPorts** - Alternative macOS package manager
- [ ] **Windows Package Managers:**
  - [ ] **Scoop** - `scoop install gohome` (developer-focused)
  - [ ] **Chocolatey** - `choco install gohome` (enterprise-friendly)

### 📖 Official Documentation Site

**Goal:** Create a professional, searchable, version-controlled documentation website.

- [ ] **Documentation Platform:**
  - [ ] Choose platform: VitePress, Docusaurus, MkDocs, or custom
  - [ ] Domain: `docs.gohome.dev` or similar
  - [ ] Hosting: GitHub Pages, Vercel, Netlify, or Cloudflare Pages
- [ ] **Content Structure:**
  - [ ] Getting Started guide
  - [ ] Installation instructions (all platforms)
  - [ ] Configuration reference (all flags, env vars, config file)
  - [ ] Usage examples and recipes
  - [ ] API documentation (for plugin developers)
  - [ ] Migration guides (v1 → v2, etc.)
  - [ ] Troubleshooting and FAQ
  - [ ] Contributing guide
- [ ] **Features:**
  - [ ] Full-text search
  - [ ] Version selector (docs for each major version)
  - [ ] Dark/light mode
  - [ ] Code syntax highlighting
  - [ ] Copy-to-clipboard for code blocks
  - [ ] Mobile-responsive design
- [ ] **CI/CD:**
  - [ ] Auto-deploy on main branch updates
  - [ ] Preview deploys for PRs
  - [ ] Broken link checker

### 🌐 Marketing & Community Growth

**Goal:** Establish online presence and grow user base.

- [ ] **Landing Page:**
  - [ ] Professional landing page: `gohome.dev`
  - [ ] Hero section with demo GIF
  - [ ] Feature highlights with icons
  - [ ] Installation quick-start
  - [ ] Testimonials and use cases
  - [ ] Link to docs, GitHub, community
  - [ ] Analytics (privacy-focused)
- [ ] **Content Marketing:**
  - [ ] Blog/changelog section
  - [ ] Technical blog posts
  - [ ] Tutorial videos (YouTube)
  - [ ] Community showcases
- [ ] **Community Channels:**
  - [ ] GitHub Discussions
  - [ ] Discord/Slack community (if demand exists)
  - [ ] Twitter/X presence
  - [ ] Dev.to/Medium articles
- [ ] **SEO & Discovery:**
  - [ ] Optimize for search engines
  - [ ] Submit to awesome lists
  - [ ] Product Hunt launch
  - [ ] HackerNews Show HN post

---

## 🗓️ Release Strategy

### Version Numbering

- **v1.x.x:** Stable releases with incremental improvements (Phase 1 - Complete)
- **v2.x.x:** Major architecture changes (Cobra/Viper), enhanced features, AI/TUI (Phases 2-3)
- **v3.x.x:** (Reserved for future major breaking changes)

### Release Cadence

- **Patch releases (v1.x.Y):** Bug fixes, security updates (as needed)
- **Minor releases (v1.X.0):** New features, improvements (monthly)
- **Major releases (vX.0.0):** Breaking changes, architecture shifts (yearly)

---

## 📝 Notes & Priorities

### Current Focus (v1.3.x - v2.0.0)

**Priority Order:**

1. **🏗️ Architecture Migration (TOP PRIORITY):** Cobra/Viper migration is the foundation for all Phase 2 features
   - See [docs/COBRA_VIPER_MIGRATION.md](docs/COBRA_VIPER_MIGRATION.md) for complete implementation guide
   - 12-week timeline with 4 phases
   - Enables sub-commands, plugins, advanced config management
2. **🧪 Testing:** Expand unit and integration test coverage during architecture migration
3. **⭐ Repository Whitelist/Favorites:** Daily workflow feature to scan only active projects
4. **📦 Package Managers:** AUR-bin, Homebrew, Snap, winget (priority after v2.0.0 release)
5. **⚡ Performance:** Implement concurrent scanning (easier with Cobra architecture)
6. **🔧 DX:** Add `--verbose` and `--quiet` flags (native Cobra support)

### Community Feedback Wanted

- Which package managers are most important to you?
- What export formats would you use most?
- Should we prioritize AI features or interactive mode?
- What integrations would make gohome more useful?

### Design Principles

- **Simplicity First:** Core features should "just work" without configuration
- **Privacy Focused:** No telemetry, no external services required
- **Unix Philosophy:** Do one thing well, play nicely with other tools
- **Performance Matters:** Fast enough for hundreds of repositories
- **User Agency:** Users control their data and workflow

---

_Note: This roadmap is a living document and subject to change based on user feedback, priorities, and community contributions. Want to influence the roadmap? Open an issue or discussion on GitHub!_
