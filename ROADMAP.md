# 🗺️ Product Roadmap

This document outlines the development status and future plans for **gohome** (Git Standup Tool).

## 🚀 Phase 1: The MVP & Professional Release (v1.0.0)

**Goal:** Deliver a stable, performant CLI tool with a complete distribution pipeline, suitable for both daily use and scripting.

### Feature Checklist

**Core & Logic:**

- [x] **Git Integration:** Auto-scan directories (`.git`) to detect repositories.
- [x] **Log Parsing:** Parse `git log` output using Conventional Commits regex.
- [x] **Smart Configuration:**
  - [x] Load config from JSON file (`~/.gohome.json`).
  - [x] Support command-line flags with shorthand aliases.
  - [x] Persist settings via `--save` flag.
  - [x] Auto-detect git author from system config.
- [ ] **Concurrency:** Implement Fan-out/Fan-in pattern using Goroutines for fast scanning.
- [x] **Custom Tasks:** Allow users to append manual tasks (e.g., meetings, code reviews) via `-t` flags.

**User Interface (UI/UI):**

- [x] **Output Formats:** Support both `text` list and rich `table` format.
- [x] **Styling:** Custom table styles (markdown, nature, tech, etc.).
- [x] **Clipboard:** Cross-platform clipboard support (`--copy`).
- [x] **Feedback:** Add a Spinner/Loading indicator during the scanning process (UX).

**System & Refinements:**

- [x] **Versioning:** Implement `--version` (`-v`) flag (injected via build time).
- [ ] **Debugging:** Implement `--verbose` flag to print debug logs (scanned paths, git errors).
- [ ] **Scripting:** Implement `--quiet` (`-q`) flag to suppress banners and meta-info (output only raw data).
- [ ] **Filtering:**
  - [ ] Filter by commit types (e.g., `--types feat,fix`).
  - [ ] Exclude specific patterns/directories (e.g., `--exclude vendor,node_modules`).
- [ ] **Validation:** Better error handling for invalid paths or git errors.
- [x] **Help:** Refine help messages and examples using `tabwriter`.

**Quality Assurance:**

- [ ] **Unit Tests:** Add test coverage for `parser` and `config` packages.
- [ ] **Integration Tests:** Test the full flow with a dummy git repo.

**CI/CD & Distribution:**

- [x] **GitHub Actions:** Setup workflow for linting (`golangci-lint`) and testing on every push.
- [x] **GoReleaser Integration:** Automate release process (see [RELEASE_GUIDE.md](RELEASE_GUIDE.md)).
- [x] **Cross-Platform Builds:** Binaries for Linux (amd64/arm64), Windows, macOS (Intel/Apple Silicon).
- [x] **Installation Support:** Multiple installation methods for all platforms.
  - [x] Direct binary downloads (GitHub Releases)
  - [x] Universal installation scripts:
    - [x] `install.sh` (Linux/macOS via curl)
    - [x] `install.ps1` (Windows via PowerShell)
  - [ ] Package managers:
    - [ ] Homebrew (macOS/Linux)
    - [ ] Scoop/Chocolatey (Windows)
    - [ ] Snap (Linux universal)
    - [ ] APT/RPM packages (Debian/RedHat)
    - [x] npm (Node.js users)

**Documentation & Support:**

- [x] **README.md:**
  - [x] Comprehensive installation guide (Go install, Binary download).
  - [x] Usage examples with terminal recordings (asciinema + GIF).
  - [x] Configuration guide (flags explanation).
- [ ] **Contribution Guide:** Instructions for developers (Running tests, Linting).

---

## 🔮 Phase 2: Advanced Features & Ecosystem (v1.x.x)

**Goal:** Enhance usability with AI, interactive UI, and robust architecture.

### Architecture Refactoring

- [ ] **Migrate to Cobra:** Restructure the application to support sub-commands (e.g., `gohome config`, `gohome summary`).
- [ ] **Adopt Viper:**
  - [ ] Support Environment Variables (essential for API Keys).
  - [ ] Support YAML/TOML config formats.
  - [ ] Hierarchy management: Flag > Env > Config > Default.

### New Features

- [ ] **Static/Recurring Tasks:** Support defining daily recurring tasks (e.g., "Daily Standup") in config file.
- [ ] **Export Options:**
  - [ ] Export to JSON (for programmatic use/integration).
  - [ ] Export to Markdown (`.md`).
  - [ ] Export to HTML (Report style).
  - [ ] Export to PDF (requires external libs).
- [ ] **AI-Powered Summary:**
  - [ ] Integrate with LLMs (OpenAI/Anthropic/Gemini) to generate a concise daily summary from raw commit logs.
  - [ ] Prompt engineering for "Standup style" or "Changelog style".
- [ ] **Interactive Mode (TUI):**
  - [ ] Implement `charmbracelet/bubbletea` interface.
  - [ ] Allow users to interactively select/deselect repositories to include in the report.
- [ ] **Advanced Filtering:**
  - [ ] Exclude specific repositories or folders.
  - [ ] Filter by commit message pattern (Regex).

---

## 🧪 Phase 3: Analytics & Integrations (Future)

**Goal:** Data insights and workflow integrations.

- [ ] **Integrations:** Slack/Discord webhook support.
- [ ] **Analytics:** Commit heatmaps, contributor stats.
- [ ] **Dashboard:** A simple web-view for local history.

---

_Note: This roadmap is subject to change based on user feedback and priorities._
