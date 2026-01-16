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
- [Flags Reference](#flags-reference)
- [Configuration](#configuration)
- [Environment Variables](#environment-variables)
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

Location: `~/.gohome.json`

The config file stores your default settings. It's created automatically when you use `--save` or `gohome config set`.

**Example config file:**
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

### Configuration Priority

Settings are loaded in this order (highest to lowest priority):

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Configuration file** (~/.gohome.json)
4. **Default values** (lowest priority)

**Example:**
```bash
# Config file: days=5
# This overrides config file with days=3
gohome --days=3

# This uses config file value (days=5)
gohome
```

---

## 🌍 Environment Variables

All configuration can be set via environment variables with the `GOHOME_` prefix.

| Config Key | Environment Variable | Example |
|------------|---------------------|---------|
| `days` | `GOHOME_DAYS` | `GOHOME_DAYS=3` |
| `format` | `GOHOME_FORMAT` | `GOHOME_FORMAT=table` |
| `path` | `GOHOME_PATH` | `GOHOME_PATH=~/workspace` |
| `max_depth` | `GOHOME_MAX_DEPTH` | `GOHOME_MAX_DEPTH=3` |
| `author` | `GOHOME_AUTHOR` | `GOHOME_AUTHOR="John Doe"` |
| `copy` | `GOHOME_COPY` | `GOHOME_COPY=true` |

**Examples:**
```bash
# Set environment variable
export GOHOME_DAYS=7
gohome  # Uses 7 days

# One-time override
GOHOME_FORMAT=table gohome

# Multiple variables
GOHOME_DAYS=3 GOHOME_FORMAT=table GOHOME_COPY=true gohome
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
