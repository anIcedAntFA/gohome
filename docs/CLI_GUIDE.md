# 📚 gohome CLI Usage Guide

Complete guide to using gohome CLI tool for git activity reporting.

---

## 📖 Table of Contents

- [Quick Start](#quick-start)
- [Command Structure](#command-structure)
- [Commands](#commands)
  - [Default Command (report)](#default-command-report)
  - [config](#config)
  - [version](#version)
  - [completion](#completion)
- [Flags Reference](#flags-reference)
- [Configuration](#configuration)
- [Environment Variables](#environment-variables)
- [Shell Completions](#shell-completions)
- [Common Use Cases](#common-use-cases)
- [Examples](#examples)

---

## 🚀 Quick Start

```bash
# Generate today's activity report
gohome

# Last 3 days with table format
gohome --days=3 --format=table

# Save your preferences
gohome --days=7 --format=table --style=markdown --save

# Now just run with saved settings
gohome
```

---

## 🏗️ Command Structure

```
gohome [command] [flags]
```

### Available Commands

```
gohome           # Default: generate activity report
gohome report    # Same as default
gohome config    # Manage configuration
gohome version   # Show version info
```

---

## 📋 Commands

### Default Command (report)

Generate git activity report from your repositories.

**Syntax:**
```bash
gohome [flags]
gohome report [flags]
```

Both forms are identical - `gohome` automatically runs the `report` command.

**Examples:**
```bash
# Basic usage
gohome
gohome report

# With custom time period
gohome --days=7
gohome report --days=7

# All flags work on both forms
gohome -d 3 -f table
gohome report -d 3 -f table
```

---

### config

Manage gohome configuration settings.

#### config list

Display all current configuration settings.

```bash
gohome config list
```

**Output:**
```
Current Configuration:
=====================
Config file: /home/user/.gohome.json

days                 = 5
format               = table
path                 = .
max_depth            = 2
...
```

#### config get

Get a specific configuration value.

```bash
gohome config get <key>
```

**Examples:**
```bash
gohome config get days
# Output: days = 5

gohome config get format
# Output: format = table
```

#### config set

Set a configuration value and save to file.

```bash
gohome config set <key> <value>
```

**Examples:**
```bash
gohome config set days 7
gohome config set format table
gohome config set max_depth 3
```

#### config reset

Reset configuration to defaults (deletes config file).

```bash
gohome config reset
```

**Prompts for confirmation:**
```
Are you sure you want to delete /home/user/.gohome.json? (y/N):
```

---

### version

Display version information.

```bash
gohome version
```

**Output:**
```
gohome version 2.0.0
commit: abc123def
built: 2026-01-16T10:00:00Z
```

---

### completion

Generate shell completion scripts for tab-completion of commands, subcommands, and flag values.

```bash
gohome completion [bash|zsh|fish|powershell]
```

**Supported Shells:**
- Bash
- Zsh
- Fish
- PowerShell

**Installation Examples:**

**Fish (Permanent):**
```bash
gohome completion fish > ~/.config/fish/completions/gohome.fish
```

**Fish (Temporary):**
```bash
gohome completion fish | source
```

**Bash (System-wide - Linux):**
```bash
sudo gohome completion bash > /etc/bash_completion.d/gohome
```

**Bash (User-only):**
```bash
echo 'source <(gohome completion bash)' >> ~/.bashrc
source ~/.bashrc
```

**Zsh:**
```bash
# Enable completions first (if not already enabled)
echo "autoload -U compinit; compinit" >> ~/.zshrc

# Install completion
gohome completion zsh > "${fpath[1]}/_gohome"

# Restart shell
exec zsh
```

**PowerShell:**
```powershell
# Temporary (this session only)
gohome completion powershell | Out-String | Invoke-Expression

# Permanent (add to profile)
Add-Content $PROFILE "gohome completion powershell | Out-String | Invoke-Expression"
```

**What Completions Provide:**

- **Command completion**: `gohome <TAB>` → shows `config`, `report`, `version`, `completion`
- **Subcommand completion**: `gohome config <TAB>` → shows `get`, `set`, `list`, `reset`
- **Config key completion**: `gohome config get <TAB>` → shows all 17 config keys
- **Flag value completion**: `gohome --format <TAB>` → shows `text`, `table`
- **Style completion**: `gohome --style <TAB>` → shows `normal`, `markdown`, `nature`, `tech`
- **Shell completion**: `gohome completion <TAB>` → shows `bash`, `zsh`, `fish`, `powershell`

**Verifying Installation:**

After installing completions, test by typing:
```bash
gohome <TAB><TAB>        # Should show commands
gohome config <TAB><TAB>  # Should show subcommands
gohome --format <TAB>     # Should show format values
```

---

## 🚩 Flags Reference

### Time Period Flags

Control how far back to look for commits.

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--hours` | `-H` | 0 | Number of hours to look back |
| `--days` | `-d` | 1 | Number of days to look back |
| `--weeks` | `-w` | 0 | Number of weeks to look back |
| `--months` | `-m` | 0 | Number of months to look back |
| `--years` | `-y` | 0 | Number of years to look back |
| `--today` | - | false | From midnight to now |

**Note:** Only one time period flag should be used at a time. Priority: years > months > weeks > days > hours.

**Examples:**
```bash
gohome --days=3        # Last 3 days
gohome -d 7            # Last 7 days  
gohome --weeks=2       # Last 2 weeks
gohome --today         # Since midnight
```

---

### Path & Scanning Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--path` | `-p` | `.` | Root path to scan for repositories |
| `--max-depth` | - | 2 | Maximum depth to scan (0-10) |
| `--author` | `-a` | auto | Git author name (auto-detected) |

**Examples:**
```bash
gohome --path=/workspace        # Scan /workspace
gohome -p ~/projects            # Scan ~/projects
gohome --max-depth=3            # Scan 3 levels deep
gohome --author="John Doe"      # Filter by author
```

---

### Output Format Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--format` | `-f` | text | Output format: `text`, `table` |
| `--style` | `-s` | normal | Table style: `normal`, `markdown` |
| `--icon` | `-i` | false | Show commit type icons |
| `--scope` | `-c` | false | Show commit scope |

**Examples:**
```bash
gohome --format=text            # Simple text output
gohome -f table                 # Table format
gohome -f table -s markdown     # Markdown table
gohome -i -c                    # Show icons and scopes
```

---

### Branch Filtering Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--all-branches` | `-b` | false | Include all local branches |
| `--branch` | - | - | Filter by specific branch |

**Examples:**
```bash
gohome --all-branches           # All branches
gohome -b                       # All branches (short)
gohome --branch=main            # Only main branch
gohome --branch=develop         # Only develop branch
```

---

### Utility Flags

| Flag | Shorthand | Default | Description |
|------|-----------|---------|-------------|
| `--copy` | `-C` | false | Copy output to clipboard |
| `--task` | `-t` | - | Add custom task (repeatable) |
| `--save` | - | false | Save settings to config file |

**Examples:**
```bash
gohome --copy                      # Copy to clipboard
gohome -C                          # Copy (short)
gohome -t "Fix bug" -t "Review PR" # Add tasks
gohome --save                      # Save current flags
```

---

## ⚙️ Configuration

### Configuration File

**Supported Formats:** JSON, YAML, TOML

**Locations:**
- `~/.gohome.json` (JSON format - created by `--save`)
- `~/.gohome.yaml` or `~/.gohome.yml` (YAML format)
- `~/.gohome.toml` (TOML format)

The config file stores your default settings. It's created automatically when you use `--save` (generates JSON) or `gohome config set`.

**Format Precedence:** If multiple config files exist, gohome uses the first one found: `.json` > `.toml` > `.yaml` > `.yml`

**Example config files:**

<details>
<summary><strong>JSON Format</strong> (~/.gohome.json)</summary>

```json
{
  "hours": 0,
  "days": 5,
  "weeks": 0,
  "months": 0,
  "years": 0,
  "today": false,
  "path": ".",
  "max_depth": 2,
  "author": "your-name",
  "format": "table",
  "preset": "markdown",
  "show_icon": true,
  "show_scope": false,
  "all_branches": false,
  "branch": "",
  "copy": false,
  "tasks": [
    {
      "type": "review",
      "message": "Code Review & PR Feedback",
      "icon": "👀",
      "enabled": true
    }
  ]
}
```
</details>

<details>
<summary><strong>YAML Format</strong> (~/.gohome.yaml)</summary>

```yaml
# Time period
hours: 0
days: 5
weeks: 0
months: 0
years: 0
today: false

# Path and scanning
path: .
max_depth: 2
author: your-name

# Output
format: table
style: markdown
icon: true
scope: false

# Branch filtering
all_branches: false
branch: ""

# Clipboard
copy: false

# Tasks
tasks:
  - type: review
    message: "Code Review & PR Feedback"
    icon: "👀"
    enabled: true
```
</details>

<details>
<summary><strong>TOML Format</strong> (~/.gohome.toml)</summary>

```toml
# Time period
hours = 0
days = 5
weeks = 0
months = 0
years = 0
today = false

# Path and scanning
path = "."
max_depth = 2
author = "your-name"

# Output
format = "table"
style = "markdown"
icon = true
scope = false

# Branch filtering
all_branches = false
branch = ""

# Clipboard
copy = false

# Tasks
[[tasks]]
type = "review"
message = "Code Review & PR Feedback"
icon = "👀"
enabled = true
```
</details>

### Configuration Priority

Settings are loaded in this order (highest to lowest priority):

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Configuration file** (~/.gohome.json/.yaml/.toml)
4. **Default values** (lowest priority)

**Example:**
```bash
# Config file: days=5
# This overrides config file with days=3
gohome --days=3

# This uses config file value (days=5)
gohome
```

**Choosing a Format:**

- **JSON**: Auto-generated by `--save` or `config set`, best for programmatic editing
- **YAML**: Human-friendly, great for manual editing with comments. Create `~/.gohome.yaml` manually (see examples above)
- **TOML**: Clean syntax, good balance between JSON and YAML. Create `~/.gohome.toml` manually (see examples above)

**Tip:** Copy examples from the collapsible sections above and paste into your config file.

---

## 🌍 Environment Variables

All configuration can be set via environment variables with the `GOHOME_` prefix.

### Supported Environment Variables

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

### Basic Usage

```bash
# Set environment variable (persistent in current shell session)
export GOHOME_DAYS=7
gohome  # Uses 7 days

# One-time use (just for this command)
GOHOME_FORMAT=table gohome

# Multiple variables at once
GOHOME_DAYS=3 GOHOME_FORMAT=table GOHOME_COPY=true gohome
```

### Practical Examples

#### Example 1: CI/CD Pipeline
```bash
# .github/workflows/daily-report.yml
- name: Generate Daily Report
  env:
    GOHOME_DAYS: 1
    GOHOME_FORMAT: table
    GOHOME_STYLE: markdown
    GOHOME_AUTHOR: ci-bot
  run: gohome
```

#### Example 2: Docker Container
```bash
# Pass config through environment variables
docker run -e GOHOME_DAYS=7 \
           -e GOHOME_FORMAT=table \
           -e GOHOME_PATH=/workspace \
           -v $(pwd):/workspace \
           gohome:latest
```

#### Example 3: Fish Shell Aliases
```fish
# ~/.config/fish/config.fish

# Quick daily standup
alias standup="set -x GOHOME_DAYS 1; set -x GOHOME_FORMAT table; gohome"

# Weekly summary
alias weekly="set -x GOHOME_DAYS 7; set -x GOHOME_FORMAT table; set -x GOHOME_STYLE markdown; gohome"
```

#### Example 4: Bash/Zsh Profile
```bash
# ~/.bashrc or ~/.zshrc

# Set defaults via environment
export GOHOME_FORMAT=table
export GOHOME_STYLE=markdown
export GOHOME_MAX_DEPTH=3

# Aliases for common tasks
alias standup="GOHOME_DAYS=1 gohome"
alias weekly="GOHOME_DAYS=7 gohome"
```

### Precedence Testing

Environment variables have **lower priority** than command-line flags:

```bash
# Set env var
export GOHOME_DAYS=10

# This uses 10 days (from env var)
gohome

# This uses 3 days (flag overrides env var)
gohome --days 3

# Verify current precedence
gohome config list  # Shows which value is active
```

---

## 🎯 Shell Completions

Shell completions provide intelligent tab-completion for gohome commands, making the CLI faster and more discoverable.

### Features

- **Command completion**: Suggests available commands
- **Subcommand completion**: Context-aware subcommand suggestions
- **Flag completion**: Shows all available flags
- **Dynamic value completion**: Suggests valid values for flags like `--format` and `--style`
- **Config key completion**: Lists all 17 config keys for `config get/set`

### Installation by Shell

#### Fish

**Permanent Installation:**
```bash
gohome completion fish > ~/.config/fish/completions/gohome.fish
```

**Temporary (Current Session Only):**
```bash
gohome completion fish | source
```

**Testing:**
```bash
gohome <TAB>              # Shows: completion, config, help, report, version
gohome config <TAB>        # Shows: get, list, reset, set
gohome config get <TAB>    # Shows all 17 config keys
gohome --format <TAB>      # Shows: text, table
gohome --style <TAB>       # Shows: normal, markdown, nature, tech
```

#### Bash

**System-wide Installation (Linux):**
```bash
sudo gohome completion bash > /etc/bash_completion.d/gohome
```

**System-wide Installation (macOS with Homebrew):**
```bash
gohome completion bash > $(brew --prefix)/etc/bash_completion.d/gohome
```

**User-only Installation:**
```bash
echo 'source <(gohome completion bash)' >> ~/.bashrc
source ~/.bashrc
```

**Testing:**
```bash
gohome <TAB><TAB>          # Shows available commands
gohome --format <TAB><TAB> # Shows: text table
```

#### Zsh

**Setup (First time only):**
```bash
# Enable completions if not already enabled
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

**Installation:**
```bash
gohome completion zsh > "${fpath[1]}/_gohome"
# Restart shell
exec zsh
```

**Alternative (User-only):**
```bash
# Add to ~/.zshrc
echo 'source <(gohome completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

**Testing:**
```bash
gohome <TAB>            # Shows commands
gohome config get <TAB>  # Shows config keys
```

#### PowerShell

**Temporary (Current Session Only):**
```powershell
gohome completion powershell | Out-String | Invoke-Expression
```

**Permanent Installation:**
```powershell
# Add to PowerShell profile
$completionScript = "gohome completion powershell | Out-String | Invoke-Expression"
Add-Content $PROFILE $completionScript

# Reload profile
. $PROFILE
```

**Testing:**
```powershell
gohome <TAB>           # Shows commands
gohome --format <TAB>  # Shows format values
```

### Completion Examples

#### Example 1: Command Discovery
```bash
$ gohome <TAB>
completion  config  help  report  version

$ gohome completion <TAB>
bash  fish  powershell  zsh
```

#### Example 2: Config Management
```bash
$ gohome config <TAB>
get  list  reset  set

$ gohome config get <TAB>
hours  days  weeks  months  years  today  path  max_depth  author  format  style  icon  scope  all_branches  branch  copy

$ gohome config get for<TAB>
format  # Auto-completes
```

#### Example 3: Flag Values
```bash
$ gohome --format <TAB>
text  table

$ gohome --style <TAB>
normal  markdown  nature  tech

$ gohome -f table --style mark<TAB>
markdown  # Auto-completes
```

#### Example 4: Workflow Speedup
```bash
# Instead of typing the full command:
gohome --days=7 --format=table --style=markdown

# With completions, just type and press TAB:
gohome -d 7 -f ta<TAB> -s ma<TAB>
# Result: gohome -d 7 -f table -s markdown
```

### Debugging Completions

If completions aren't working, verify installation:

**Fish:**
```bash
ls ~/.config/fish/completions/gohome.fish  # Should exist
complete -C gohome  # Shows registered completions
```

**Bash:**
```bash
complete -p gohome  # Should show: complete -o default -F __start_gohome gohome
```

**Zsh:**
```bash
echo $fpath  # Check completion search paths
ls ${fpath[1]}/_gohome  # Should exist
```

**PowerShell:**
```powershell
$PROFILE  # Shows profile location
Get-Content $PROFILE | Select-String "gohome"  # Verify entry exists
```

### Uninstalling Completions

**Fish:**
```bash
rm ~/.config/fish/completions/gohome.fish
```

**Bash:**
```bash
# System-wide
sudo rm /etc/bash_completion.d/gohome

# User-only (remove from ~/.bashrc)
sed -i '/gohome completion/d' ~/.bashrc
```

**Zsh:**
```bash
rm "${fpath[1]}/_gohome"
# Or remove from ~/.zshrc if using source method
sed -i '/gohome completion/d' ~/.zshrc
```

**PowerShell:**
```powershell
# Edit profile and remove gohome completion line
notepad $PROFILE
```

---

## 💡 Common Use Cases

### 1. Daily Standup Report

Generate a quick report for your standup meeting.

```bash
# Yesterday's work
gohome --days=1

# Since last Friday (Monday standup)
gohome --days=3

# Copy to clipboard for Slack
gohome --days=1 --copy
```

### 2. Weekly Summary

Review your week's activity.

```bash
# Last 7 days, table format
gohome --days=7 --format=table --style=markdown

# Save as default for weekly reviews
gohome -d 7 -f table -s markdown --save

# Now just run weekly
gohome
```

### 3. Sprint Report

Generate reports for your sprint period.

```bash
# 2-week sprint
gohome --weeks=2 --format=table

# With all branches
gohome --weeks=2 --all-branches
```

### 4. Specific Project Report

Scan only a specific directory.

```bash
# Specific project
gohome --path=/path/to/project

# Workspace with multiple projects
gohome --path=/workspace --max-depth=3
```

### 5. Custom Tasks

Add non-commit tasks to your report.

```bash
gohome -t "Attended team meeting" -t "Code review sessions"

# Combine with commits
gohome --days=1 -t "Updated documentation" -t "Mentoring session"
```

### 6. Branch-Specific Reports

Report only from specific branches.

```bash
# Only main branch
gohome --branch=main

# Only develop branch
gohome --branch=develop

# All branches
gohome --all-branches
```

---

## 📝 Examples

### Basic Examples

```bash
# Default report (1 day)
gohome

# Last 3 days
gohome --days=3

# Last week
gohome --weeks=1

# Today only
gohome --today
```

### Formatted Output

```bash
# Simple text
gohome --format=text

# Table format
gohome --format=table

# Markdown table
gohome --format=table --style=markdown

# With icons and scopes
gohome --format=table --icon --scope
```

### Advanced Examples

```bash
# Full-featured weekly report
gohome --days=7 \
  --format=table \
  --style=markdown \
  --icon \
  --all-branches \
  --copy

# Specific project with tasks
gohome --path=/workspace/myproject \
  --days=3 \
  --format=table \
  -t "Fixed production bug" \
  -t "Updated documentation"

# Cross-team report
gohome --path=/workspace \
  --max-depth=4 \
  --weeks=2 \
  --format=table \
  --style=markdown \
  --copy
```

### Configuration Management

```bash
# Save your preferences
gohome --days=5 --format=table --style=markdown --save

# Check current settings
gohome config list

# Update specific setting
gohome config set days 7

# Get specific value
gohome config get format

# Reset to defaults
gohome config reset
```

### Using Environment Variables

```bash
# Set for session
export GOHOME_DAYS=7
export GOHOME_FORMAT=table
gohome

# One-time override
GOHOME_DAYS=3 gohome

# Override config file temporarily
GOHOME_FORMAT=text gohome  # Even if config has format=table
```

---

## 🎯 Best Practices

### 1. Save Your Defaults

Set up your preferred configuration once:

```bash
gohome --days=5 --format=table --style=markdown --icon --save
```

Then just run `gohome` daily.

### 2. Use Aliases

Add to your shell config:

```bash
# ~/.bashrc or ~/.zshrc
alias standup='gohome --days=1 --copy'
alias weekly='gohome --days=7 --format=table'
alias sprint='gohome --weeks=2 --format=table --style=markdown'
```

### 3. Project-Specific Reports

For different projects:

```bash
alias proj1='gohome --path=~/work/project1'
alias proj2='gohome --path=~/work/project2'
```

### 4. Integration with Other Tools

```bash
# Copy and paste to Slack
gohome --days=1 --copy

# Save to file
gohome --days=7 > weekly-report.md

# Email report
gohome --days=7 --format=table | mail -s "Weekly Report" team@company.com
```

---

## 🐛 Troubleshooting

### No commits found

```bash
# Check if repositories are detected
gohome --path=/your/path

# Increase scan depth
gohome --max-depth=3

# Check author name
git config user.name
gohome --author="Your Name"
```

### Config file issues

```bash
# Check config location
gohome config list

# Reset if corrupted
gohome config reset

# Recreate
gohome --days=5 --save
```

### Clipboard not working

On Linux, install clipboard tools:

```bash
# Wayland
sudo apt install wl-clipboard

# X11
sudo apt install xclip
```

---

## 📚 See Also

- [CLI Design Philosophy](./CLI_DESIGN.md) - Design patterns and best practices
- [README.md](../README.md) - Main documentation
- [CONTRIBUTING.md](../CONTRIBUTING.md) - Development guide

---

**Need help?** Open an issue on [GitHub](https://github.com/anIcedAntFA/gohome/issues)
