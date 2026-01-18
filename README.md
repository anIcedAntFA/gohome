<h1 align="center">gohome</h1>

<p align="center">
  A fast, configurable Git standup & activity reporting CLI written in Go.
</p>

<p align="center">
  <sub>
    Turn your local Git commits across multiple repositories into clean, daily developer reports.
  </sub>
</p>

<p align="center">
  <a href="https://github.com/anIcedAntFA/gohome/actions">
    <img
      src="https://img.shields.io/github/actions/workflow/status/anIcedAntFA/gohome/release.yml?label=build&logo=githubactions&logoColor=white"
      alt="Build status"
    />
  </a>
  <a href="https://codecov.io/gh/anIcedAntFA/gohome">
    <img
      src="https://codecov.io/gh/anIcedAntFA/gohome/branch/main/graph/badge.svg"
      alt="Code coverage"
    />
  </a>
  <a href="https://goreportcard.com/report/github.com/anIcedAntFA/gohome">
    <img
      src="https://img.shields.io/badge/go%20report-A+-brightgreen?logo=go"
      alt="Go Report Card grade"
    />
  </a>
  <a href="https://github.com/anIcedAntFA/gohome/releases">
    <img
      src="https://img.shields.io/github/v/release/anIcedAntFA/gohome?logo=github"
      alt="Latest release version"
    />
  </a>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/anIcedAntFA/gohome">
    <img
      src="https://pkg.go.dev/badge/github.com/anIcedAntFA/gohome.svg"
      alt="Go package documentation on pkg.go.dev"
      style="margin-right:6px;"
    />
  </a>
  <img
    src="https://img.shields.io/github/downloads/anIcedAntFA/gohome/total?logo=github"
    alt="Total GitHub downloads"
  />
  <img
    src="https://img.shields.io/badge/go-%3E%3D%201.21-00ADD8?logo=go"
    alt="Go version >= 1.21"
  />
  <img
    src="https://img.shields.io/github/license/anIcedAntFA/gohome?logo=opensourceinitiative"
    alt="Project license"
  />
</p>

**Forgot what you worked on yesterday?**

**gohome** automates your daily status reporting by recursively scanning your workspace to find git repositories. It aggregates commit logs from multiple projects instantly and formats them into beautiful, ready-to-share reports.

Perfect for **Daily Standups**, **Weekly Summaries**, or tracking your **Personal Coding Habits**.

## 🎬 Quick Demo

![gohome quickstart demo](docs/demos/quickstart.gif)

*See [docs/demos/](docs/demos/) for more examples and recording guide.*

## ✨ Features

- **🚀 Auto-Discovery:** Recursively finds git repositories in your workspace (configurable depth, default 2 levels).
- **🎯 Smart Scanning:** Skip nested repos and ignore common directories (`.git`, `.vscode`, `node_modules`).
- **⚡️ Fast Performance:** Lightweight scanner optimized for large workspace structures like `github.com/{org}/{repo}`.
- **🌱 Branch Support:** Include commits from all local branches or filter by specific branch.
- **🎨 Rich Output:** Supports multiple formats (text, table) and styles (normal, markdown, nature, tech).
- **📋 Clipboard Ready:** Copy reports directly to your system clipboard with `--copy`.
- **📝 Custom Tasks:** Add manual tasks alongside git commits for complete daily reports.
- **⚙️ Smart Config:** Persist your preferences via `~/.gohome.json` or use command-line flags.

## 📦 Installation

### Quick Install (Recommended)

**Linux/macOS:**

```bash
curl -sS https://get.ngockhoi96.dev/gohome | sh
```

The install script will:

- Auto-detect your platform (Linux/macOS, x86_64/arm64)
- Download the latest release from GitHub
- Install to `/usr/local/bin` (may require sudo)
- Clean up any conflicting dev builds in `$GOPATH/bin`
- Automatically upgrade existing installations when run again

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/anIcedAntFA/gohome/main/scripts/install.ps1 | iex
```

The PowerShell script will:

- Auto-detect your architecture (x64/arm64)
- Download and extract the latest release
- Install to `%LOCALAPPDATA%\Programs\gohome`
- Automatically add to PATH
- Clean up conflicting dev builds

### NPM

If you have Node.js and npm installed:

```bash
npm install -g @ngockhoi96/gohome
```

Or using npx (no installation required):

```bash
npx @ngockhoi96/gohome --help
```

### Arch Linux (AUR)

For Arch Linux users, install from the AUR:

```bash
# Using yay
yay -S gohome

# Or using pacman
pacman -S gohome

# Or manually
git clone https://aur.archlinux.org/gohome.git
cd gohome
makepkg -si
```

### Go Install

If you have Go 1.21+ installed:

```bash
go install github.com/anIcedAntFA/gohome/cmd/gohome@latest
```

> ⚠️ **Path Configuration:** When using production releases (installed via curl or binary download), ensure `/usr/local/bin` comes **before** `$GOPATH/bin` in your `$PATH` to avoid conflicts with dev builds.
>
> **Shell Configuration Examples:**
>
> <details>
> <summary><strong>Bash</strong> (~/.bashrc or ~/.bash_profile)</summary>
>
> ```bash
> # Go environment
> export GOPATH=$HOME/go
> export PATH=$PATH:$GOPATH/bin  # Append GOPATH/bin (lower priority)
> ```
> </details>
>
> <details>
> <summary><strong>Zsh</strong> (~/.zshrc)</summary>
>
> ```zsh
> # Go environment
> export GOPATH=$HOME/go
> export PATH=$PATH:$GOPATH/bin  # Append GOPATH/bin (lower priority)
> ```
> </details>
>
> <details>
> <summary><strong>Fish</strong> (~/.config/fish/config.fish)</summary>
>
> ```fish
> # Go environment
> set -gx GOPATH $HOME/go
> set -gx PATH $PATH $GOPATH/bin  # Append GOPATH/bin (lower priority)
> # Or use fish_add_path for better management:
> # fish_add_path -aP $GOPATH/bin
> ```
> </details>
>
> <details>
> <summary><strong>PowerShell</strong> (Windows - run as Administrator)</summary>
>
> ```powershell
> # Check current PATH
> $env:Path
>
> # Add Go bin to PATH (User level - persists across sessions)
> $goPath = "$env:USERPROFILE\go\bin"
> [Environment]::SetEnvironmentVariable(
>     "Path",
>     [Environment]::GetEnvironmentVariable("Path", "User") + ";$goPath",
>     "User"
> )
>
> # Reload PATH in current session
> $env:Path = [System.Environment]::GetEnvironmentVariable("Path","User")
> ```
>
> Note: The install.ps1 script automatically adds gohome to PATH.
> </details>
>
> After updating your shell config, reload it:
> ```bash
> # Bash
> source ~/.bashrc
> 
> # Zsh
> source ~/.zshrc
> 
> # Fish
> source ~/.config/fish/config.fish
>
> # PowerShell
> $env:Path = [System.Environment]::GetEnvironmentVariable("Path","User")
> ```

### Download Binary

Download pre-built binaries from [GitHub Releases](https://github.com/anIcedAntFA/gohome/releases/latest):

1. Download the appropriate archive for your OS/architecture
2. Extract the binary
3. Move to a directory in your `$PATH`:

**Linux/macOS:**

```bash
# Extract
tar -xzf gohome_*_linux_x86_64.tar.gz
# Move to PATH
sudo mv gohome /usr/local/bin/
# Make executable
chmod +x /usr/local/bin/gohome
```

**Windows:**

```powershell
# Extract the .zip file
# Move gohome.exe to a directory in your PATH
```

### Verify Installation

```bash
gohome --version
# Production release: gohome v1.0.1
# Dev build: gohome abc1234 (commit: abc1234, built: 2026-01-10)
```

The version format differs based on how it was built:

- **Production releases** show clean version only
- **Development builds** include commit hash and build date for debugging

## 🚀 Usage

Simply run the tool in your workspace directory:

```bash
gohome
```

### 🧪 Common Examples

**1️⃣ Basic Usage (Last 1 day)**

```bash
gohome
```

**2️⃣ Look back 3 days**

```bash
gohome -d 3
```

**3️⃣ Generate a Table Report**

```bash
gohome -f table -s markdown
```

**4️⃣ Copy to Clipboard**

This is useful for pasting directly into Slack/Teams/Discord:

```bash
gohome -d 1 --copy
```

**5️⃣ Add Custom Tasks**

Add tasks that aren't tracked in git:

```bash
gohome -t "Meeting: Sprint Planning" -t "Review: PR #123"
```

**6️⃣ Include All Local Branches**

By default, gohome only shows commits from your current branch. Use `-b` to include commits from all local branches:

```bash
gohome -d 2 -b
```

This is useful when:

- You have commits on feature branches not yet merged
- You want to see all your work across multiple branches
- You're generating a standup report before merging PRs

**7️⃣ Filter by Specific Branch**

Filter commits from a specific branch instead of the current one:

```bash
gohome -d 3 --branch feature/new-ui
```

Useful for reviewing work on a specific feature branch without checking it out.

**8️⃣ Customize Scan Depth**

Control how deep gohome scans for repositories (default: 2 levels):

```bash
# Scan only 1 level deep (faster for flat structures)
gohome --max-depth 1

# Scan 3 levels deep for deeper nested repos
gohome --max-depth 3 -d 2
```

Useful when:
- You have a flat workspace structure (use `--max-depth 1`)
- You have deeply nested projects (increase to 3 or 4)
- You want faster scans by limiting depth

**9️⃣ Save Settings**

Save your favorite flags as default (so you don't have to type them next time):

```bash
gohome -p /Users/ngockhoi96/workspace -d 1 -f table --max-depth 2 --save
```

## 🔧 Configuration

**gohome** looks for a config file at `~/.gohome.json`. You can create it manually or use the `--save` flag to auto-generate it.

### Example Config

```json
{
  "hours": 0,
  "days": 1,
  "weeks": 0,
  "months": 0,
  "years": 0,
  "today": false,
  "path": "/Users/ngockhoi96/workspace",
  "author": "ngockhoi96",
  "format": "table",
  "preset": "normal",
  "max_depth": 2,
  "show_icon": true,
  "show_scope": false,
  "copy_to_clipboard": false,
  "all_branches": false,
  "branch": "",
  "tasks": [
    {
      "type": "meeting",
      "message": "Daily Standup & Team Sync",
      "icon": "📅",
      "enabled": true
    },
    {
      "type": "review",
      "message": "Code Review & PR Feedback",
      "icon": "👀",
      "enabled": true
    }
  ]
}
```

### 🧾 Flags Reference

| Flag             | Alias | Description                                           | Default     |
| ---------------- | ----- | ----------------------------------------------------- | ----------- |
| `--hours`        | `-H`  | Number of hours to look back                          | 0           |
| `--today`        |       | Report from midnight to now                           | false       |
| `--days`         | `-d`  | Number of days to look back                           | 1           |
| `--weeks`        | `-w`  | Number of weeks to look back                          | 0           |
| `--months`       | `-m`  | Number of months to look back                         | 0           |
| `--years`        | `-y`  | Number of years to look back                          | 0           |
| `--path`         | `-p`  | Root path to scan for repositories                    | `.`         |
| `--max-depth`    |       | Maximum depth to scan for repositories                | 2           |
| `--author`       | `-a`  | Git author name (auto-detected)                       | System User |
| `--format`       | `-f`  | Output format: `text`, `table`                        | `text`      |
| `--preset`       | `-s`  | Table style: `normal`, `markdown`                     | `normal`    |
| `--all-branches` | `-b`  | Include commits from all local branches               | false       |
| `--branch`       |       | Filter commits by specific branch                     | (current)   |
| `--copy`         | `-cp` | Copy output to clipboard                              | false       |
| `--icon`         | `-i`  | Show icon column (table format only)                  | false       |
| `--scope`        | `-c`  | Show scope column (table format only)                 | false       |
| `--task`         | `-t`  | Add custom task (repeatable)                          | []          |
| `--save`         |       | Save current flags as default config                  | false       |
| `--version`      | `-v`  | Show version information                              |             |
| `--help`         | `-h`  | Show help message                                     |             |

> **Note:** By default (without `-b` or `--branch`), gohome shows commits from your **current branch only**. Use `-b` to include all local branches, or `--branch <name>` to filter by a specific branch.

### 🌍 Environment Variables

All configuration options can be set via environment variables with the `GOHOME_` prefix. This is especially useful for:
- **CI/CD pipelines** (avoid creating config files)
- **Docker containers** (pass config through env vars)
- **Sensitive data** (API keys for future AI features)

**Configuration Precedence:** `CLI Flags > Environment Variables > Config File > Defaults`

#### Supported Environment Variables

| Config Key       | Environment Variable        | Example Value                      |
|------------------|-----------------------------|----------------------------------- |
| `hours`          | `GOHOME_HOURS`             | `24`                               |
| `days`           | `GOHOME_DAYS`              | `7`                                |
| `weeks`          | `GOHOME_WEEKS`             | `2`                                |
| `months`         | `GOHOME_MONTHS`            | `1`                                |
| `years`          | `GOHOME_YEARS`             | `1`                                |
| `today`          | `GOHOME_TODAY`             | `true`                             |
| `path`           | `GOHOME_PATH`              | `/home/user/workspace`             |
| `max_depth`      | `GOHOME_MAX_DEPTH`         | `3`                                |
| `author`         | `GOHOME_AUTHOR`            | `johndoe`                          |
| `format`         | `GOHOME_FORMAT`            | `table`                            |
| `style`          | `GOHOME_STYLE`             | `markdown`                         |
| `icon`           | `GOHOME_ICON`              | `true`                             |
| `scope`          | `GOHOME_SCOPE`             | `true`                             |
| `all_branches`   | `GOHOME_ALL_BRANCHES`      | `true`                             |
| `branch`         | `GOHOME_BRANCH`            | `main`                             |
| `copy`           | `GOHOME_COPY`              | `true`                             |

#### Examples

```bash
# Set defaults via environment
export GOHOME_DAYS=7
export GOHOME_FORMAT=table
export GOHOME_STYLE=markdown
gohome  # Uses env var values

# Override env var with flag
export GOHOME_DAYS=7
gohome --days 3  # Flag wins: uses 3 days

# Useful for CI/CD
export GOHOME_PATH=/workspace
export GOHOME_AUTHOR=ci-bot
export GOHOME_FORMAT=table
gohome --today

# Docker usage
docker run -e GOHOME_DAYS=7 -e GOHOME_FORMAT=table gohome:latest
```

## 🗺️ Roadmap

See [ROADMAP.md](ROADMAP.md) for the full development plan and upcoming features.

## 🤝 Contributing

Contributions are welcome! We appreciate bug reports, feature requests, documentation improvements, and code contributions.

For detailed guidelines on:
- Development setup and workflow
- Coding standards and conventions
- Commit message format (Conventional Commits with emojis)
- Pull request process
- Testing and quality assurance

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for the complete guide.

### Quick Start for Contributors

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/gohome.git
cd gohome

# 2. Install dependencies
go mod tidy

# 3. Create feature branch
git checkout -b feat/amazing-feature

# 4. Make changes and test
make build
make test
make lint

# 5. Commit using Conventional Commits
git commit -m '✨ feat(scanner): add amazing feature'

# 6. Push and open PR
git push origin feat/amazing-feature
```

## ❤️ Credits & Motivation

**gohome** is heavily inspired by the awesome [git-standup](https://github.com/kamranahmedse/git-standup) utility by [Kamran Ahmed](https://github.com/kamranahmedse).

While `git-standup` is great, **gohome** was built to address specific personal needs for daily reporting, such as:

- **Rich formatting:** Tables, icons, and custom styles.
- **Workflow integration:** Direct clipboard support.
- **Smart config:** Persisted settings for zero-setup runs.

This project also serves as a practical journey to master **Go (Golang)**, implementing concepts like Concurrency, CLI architecture, and Cross-platform distribution.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
