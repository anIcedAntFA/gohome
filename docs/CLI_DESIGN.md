# 🎨 CLI Design Philosophy & Best Practices

Comprehensive guide to designing effective command-line interfaces, based on research from industry-leading tools (kubectl, docker, git, gh, cargo, npm) and the gohome implementation.

---

## 📖 Table of Contents

- [Overview](#overview)
- [Command vs Subcommand Architecture](#command-vs-subcommand-architecture)
- [Real-World Examples](#real-world-examples)
- [Design Principles](#design-principles)
- [Command Structure Patterns](#command-structure-patterns)
- [Configuration Management](#configuration-management)
- [UX/DX Guidelines](#uxdx-guidelines)
- [Implementation Decisions](#implementation-decisions)
- [gohome Design Rationale](#gohome-design-rationale)

---

## 🎯 Overview

A well-designed CLI provides:
- **Intuitive** command structure
- **Consistent** behavior and naming
- **Discoverable** functionality
- **Flexible** configuration
- **Fast** execution and feedback

---

## 🏗️ Command vs Subcommand Architecture

### What are Commands and Subcommands?

```
<tool> <command> <subcommand> [flags] [arguments]
```

- **Tool**: The CLI binary (`gohome`, `kubectl`, `docker`)
- **Command**: Primary action or resource type (`report`, `config`, `version`)
- **Subcommand**: Secondary action on the command (`list`, `get`, `set`)
- **Flags**: Modifiers for behavior (`--days=3`, `--format=table`)
- **Arguments**: Positional inputs (`gohome config get days`)

### When to Use Commands vs Subcommands

#### ✅ Use Commands When:
1. **Distinct Actions**: Each represents a fundamentally different operation
   - `git commit`, `git push`, `git pull`
   - `docker build`, `docker run`, `docker ps`
   
2. **Different Resource Types**: Operating on different entities
   - `kubectl get pods`, `kubectl get services`
   - `gh pr list`, `gh issue list`

3. **Independent Workflows**: Commands that rarely interact
   - `cargo build`, `cargo test`, `cargo doc`

#### ✅ Use Subcommands When:
1. **Related Operations**: Multiple actions on the same concept
   - `gohome config list|get|set|reset` - all config operations
   - `docker compose up|down|logs` - all compose operations
   - `gh repo create|clone|fork` - all repo operations

2. **CRUD Operations**: Create, Read, Update, Delete pattern
   - `kubectl create|get|edit|delete`
   - `npm install|uninstall|update|list`

3. **Hierarchical Organization**: Grouping related functionality
   - `docker network create|ls|rm`
   - `gh secret set|list|delete`

---

## 🔍 Real-World Examples

### 1. kubectl (Kubernetes CLI)

**Structure:**
```bash
kubectl <command> <resource> [flags]
kubectl <command> <subcommand> [flags]
```

**Examples:**
```bash
# Command + Resource
kubectl get pods
kubectl describe service nginx

# Command with subcommands
kubectl config view
kubectl config set-context
kubectl config use-context
```

**Pattern:** 
- Resource-oriented (nouns): `get`, `create`, `delete` + resource type
- Config management via subcommands: `kubectl config <action>`
- Flags for filtering/formatting: `--namespace`, `--output=json`

**gohome adopts:**
- Subcommands for config management: `gohome config list|get|set`
- Flags for filtering: `--days`, `--format`

---

### 2. docker

**Structure:**
```bash
docker <command> [options] [arguments]
docker <command> <subcommand> [options]
```

**Examples:**
```bash
# Direct commands
docker build -t myapp .
docker run -p 8080:80 nginx
docker ps

# Commands with subcommands
docker compose up -d
docker network create mynetwork
docker volume ls
```

**Pattern:**
- Core operations as commands: `build`, `run`, `ps`
- Resource management via subcommands: `docker network <action>`
- Plugin-style extensions: `docker compose`, `docker buildx`

**gohome adopts:**
- Default command pattern: `gohome` = `gohome report`
- Simple core command, complex management via subcommands

---

### 3. git

**Structure:**
```bash
git <command> [flags] [arguments]
```

**Examples:**
```bash
# Primary commands
git commit -m "message"
git push origin main
git pull --rebase

# Config management
git config --global user.name "Name"
git config --list
```

**Pattern:**
- Flat command structure (minimal subcommands)
- Heavy use of flags for variants: `git log --graph --oneline`
- Config as a command: `git config`
- Global vs local scope: `--global`, `--local`

**gohome adopts:**
- Config as command: `gohome config`
- Flags for customization: `--days`, `--format`, `--style`
- Persistent configuration file

---

### 4. gh (GitHub CLI)

**Structure:**
```bash
gh <command> <subcommand> [flags]
```

**Examples:**
```bash
# Resource + action pattern
gh pr create
gh pr list
gh pr merge 123

gh issue create
gh issue list --assignee=@me
gh issue view 456

# Config
gh config set editor vim
gh config get editor
```

**Pattern:**
- Resource-first: `pr`, `issue`, `repo`, `gist`
- Actions as subcommands: `create`, `list`, `view`, `edit`
- Consistent CRUD verbs across resources

**gohome adopts:**
- Config subcommands: `list`, `get`, `set`, `reset`
- Consistent command structure

---

### 5. cargo (Rust package manager)

**Structure:**
```bash
cargo <command> [options]
```

**Examples:**
```bash
cargo build --release
cargo test --verbose
cargo run --bin myapp
cargo install ripgrep
```

**Pattern:**
- Action-oriented commands: `build`, `test`, `run`, `publish`
- Config in `Cargo.toml` (file-first)
- Minimal CLI config commands
- Heavy use of flags for variants

**gohome adopts:**
- Action-oriented: `report` (default), `config`, `version`
- Config file + CLI flags hybrid
- Flags for common variants: `--days`, `--format`

---

## 🎨 Design Principles

### 1. Principle of Least Surprise

**Bad:**
```bash
mytool --get data       # --get as flag
mytool data get         # get as subcommand
```

**Good:**
```bash
mytool get data         # Consistent: verb + noun
mytool data get         # Consistent: noun + verb
```

**gohome:**
```bash
gohome config get key   # Consistent: config + action + target
gohome config set key value
```

### 2. Progressive Disclosure

Start simple, reveal complexity when needed.

**Level 1: Beginner**
```bash
gohome                  # Simple, works immediately
```

**Level 2: Intermediate**
```bash
gohome --days=7         # Add basic flags
gohome --format=table
```

**Level 3: Advanced**
```bash
gohome --days=7 --format=table --style=markdown --icon --all-branches
```

**Level 4: Expert**
```bash
export GOHOME_FORMAT=table
gohome config set days 7
gohome                  # Customized defaults
```

### 3. Convention over Configuration

Provide sensible defaults, allow customization.

**gohome defaults:**
```bash
gohome                  # Auto: days=1, format=text, path=current
```

**Explicit:**
```bash
gohome --days=1 --format=text --path=.
```

### 4. Composability

Commands should work well together.

```bash
# Pipe to other tools
gohome | grep "feat"
gohome --format=table | wc -l

# Save to files
gohome > report.txt
gohome --format=table --style=markdown > weekly.md

# Combine with clipboard
gohome --copy && pbpaste | mail -s "Report" team@company.com
```

---

## 🏛️ Command Structure Patterns

### Pattern 1: Default Command

**When:** Tool has one primary use case.

**Example:**
```bash
gohome              # Same as: gohome report
rg "pattern"        # ripgrep (no need for 'search' command)
htop                # No subcommands needed
```

**Pros:**
- Fast to type
- Clear primary purpose
- Beginner-friendly

**Cons:**
- Less room for expansion
- May need subcommands later

**gohome uses this:**
```bash
gohome              # Primary: generate report
gohome report       # Explicit form (identical)
```

### Pattern 2: Resource + Action

**When:** Managing multiple resource types.

**Example:**
```bash
kubectl get pods
kubectl create deployment
kubectl delete service

gh pr create
gh issue list
gh repo fork
```

**Pros:**
- Clear hierarchy
- Scalable
- Consistent CRUD operations

**Cons:**
- More typing
- Steeper learning curve

### Pattern 3: Action + Resource

**When:** Actions are more important than resources.

**Example:**
```bash
git commit file.txt
git push origin main
git pull --rebase

terraform apply
terraform destroy
```

**Pros:**
- Action-first mindset
- Natural language flow
- Clear intent

**Cons:**
- Can be ambiguous with many resources

### Pattern 4: Flat Commands

**When:** Limited scope, clear separation of concerns.

**Example:**
```bash
ls -la
cd directory
pwd
```

**Pros:**
- Simple
- Fast
- Easy to learn

**Cons:**
- Hard to scale
- Namespace collisions

---

## ⚙️ Configuration Management

### Approaches from Popular Tools

#### 1. File-First (Config File is Primary)

**Examples:** git, npm, cargo

```bash
# git
git config --global user.name "Name"  # Writes to ~/.gitconfig
git config --local user.email "email"  # Writes to .git/config

# npm
npm config set registry "url"         # Writes to ~/.npmrc
```

**Pattern:**
- Config file is source of truth
- CLI commands modify the file
- File format: INI, JSON, YAML, TOML

**Pros:**
- Human-editable
- Version controllable
- Share via dotfiles

**Cons:**
- Manual editing can break things
- Schema validation needed

#### 2. CLI-First (Flags Override Everything)

**Examples:** docker, kubectl with kubeconfig

```bash
# Flags always win
kubectl get pods --context=prod       # Override config
docker run -e VAR=value nginx         # Inline config
```

**Pattern:**
- Flags are primary interface
- Config file is cache/convenience
- CLI flags > Config file

**Pros:**
- Explicit
- No hidden state
- Easy to script

**Cons:**
- Verbose
- Hard to persist changes
- Repetitive

#### 3. Hybrid (gohome approach)

```bash
# 1. Set defaults in file
gohome config set days 7              # Save to file
gohome config set format table

# 2. Use defaults
gohome                                # Uses config file

# 3. Override temporarily
gohome --days=3                       # Flag overrides config

# 4. Override permanently
gohome --days=3 --save                # Update config file

# 5. Environment variables
GOHOME_DAYS=5 gohome                  # Env override
```

**Priority Chain:**
```
CLI Flags > Environment Variables > Config File > Defaults
```

**Pattern:**
- Config file stores preferences
- Flags override temporarily
- `--save` flag persists to file
- Environment variables for CI/CD

**Pros:**
- Flexible
- Best of both worlds
- Easy migration (flags → config)

**Cons:**
- More complex implementation
- Need to document priority

---

## 💡 UX/DX Guidelines

### 1. Helpful Error Messages

**Bad:**
```
Error: Invalid input
```

**Good:**
```
❌ Error: Author not found.
  
  Tip: Use -a flag or check git config:
    $ git config user.name
    $ gohome --author="Your Name"
```

**gohome:**
```go
return fmt.Errorf("❌ Author not found. Please use -a flag or check git config")
```

### 2. Progress Feedback

**For long operations:**
```bash
🔍 Scanning repositories...  ⠋
📥 Fetching commits from project1... ⠙
📥 Fetching commits from project2... ⠹
✓ Found 15 repositories
```

**gohome uses spinners:**
```go
sp := spinner.New("🔍 Scanning repositories...")
sp.Start()
// ... operation ...
sp.Stop()
```

### 3. Success Confirmations

**Bad:**
```bash
$ gohome config set days 7
$
```

**Good:**
```bash
$ gohome config set days 7
✅ Set days = 7
```

### 4. Help Text

**Auto-generated (Cobra):**
```bash
$ gohome --help
gohome scans your workspace for git repositories...

Usage:
  gohome [flags]
  gohome [command]

Available Commands:
  config      Manage gohome configuration
  report      Generate activity report
  version     Print version information

Flags:
  -d, --days int        Number of days to look back (default 1)
  -f, --format string   Output format: text, table (default "text")
  ...
```

### 5. Completion

Shell completion for commands and flags:

```bash
# Generate completion
gohome completion bash > /etc/bash_completion.d/gohome
gohome completion zsh > "${fpath[1]}/_gohome"
gohome completion fish > ~/.config/fish/completions/gohome.fish

# Usage
$ gohome co<TAB>
config  copy

$ gohome --days=<TAB>
1  3  7  14  30
```

---

## 🛠️ Implementation Decisions

### Decision 1: Config File Location

**Options:**
- `~/.gohome.json` ✅ (chosen)
- `~/.config/gohome/config.json`
- `./.gohome.json` (project-local)

**Rationale:**
- Simple path
- Single file
- Standard for simple tools (similar to `~/.npmrc`, `~/.gitconfig`)

**Alternatives considered:**
- XDG Base Directory: `~/.config/gohome/` (more complex, overkill for one file)
- Project-local: `.gohome.json` in repo (rejected - tool is for user, not per-project)

---

### Decision 2: Default Command

**Options:**
- `gohome` runs nothing, requires subcommand ❌
- `gohome` shows help ❌
- `gohome` = `gohome report` ✅ (chosen)

**Rationale:**
- Primary use case (90% of usage)
- Fast workflow: just type `gohome`
- Similar to: `rg`, `fd`, `bat` (single-purpose tools)

**Trade-offs:**
- Explicit `gohome report` still works
- Easy to discover: `gohome --help` shows all commands

---

### Decision 3: Flag Naming

**Choices:**
- Long form: `--days`, `--format`, `--style`
- Short form: `-d`, `-f`, `-s`
- Both: `-d, --days` ✅ (chosen)

**Rationale:**
- Long form: readable in scripts
- Short form: fast for interactive use
- Both: best of both worlds

**Special cases:**
- `--copy` / `-C` (not `-c`, which is used by `--scope`)
- `--icon` / `-i` (clear abbreviation)

---

### Decision 4: Configuration Priority

**Priority Chain:**
```
1. CLI Flags (highest)
2. Environment Variables
3. Config File (~/.gohome.json)
4. Defaults (lowest)
```

**Example scenario:**
```bash
# Config file: days=7
# Command: GOHOME_DAYS=3 gohome --days=1

Result: days=1  # CLI flag wins
```

**Why this order?**
- **CLI Flags**: Explicit user intent right now
- **Env Variables**: Session/CI configuration
- **Config File**: User's saved preferences
- **Defaults**: Fallback when nothing specified

---

### Decision 5: --save Flag Behavior

**Options:**
1. Automatic save on every run ❌
2. Explicit `--save` flag ✅ (chosen)
3. `gohome config set` only ❌

**Rationale:**
- Explicit `--save`: Clear when config changes
- No automatic save: Prevent accidental overrides
- Still provide `gohome config set` for direct editing

**User flow:**
```bash
# Try settings temporarily
gohome --days=7 --format=table

# Like it? Save it
gohome --days=7 --format=table --save

# Now it's the default
gohome
```

---

### Decision 6: Subcommand Structure

**For config management:**
```bash
gohome config list       # Show all settings
gohome config get key    # Get specific value
gohome config set key value  # Set value
gohome config reset      # Reset to defaults
```

**Why subcommands for config?**
- Clear grouping of related operations
- Follows kubectl, gh, docker patterns
- Scalable (easy to add more config operations)

**Alternative rejected:**
```bash
gohome list-config       # Flat structure - harder to discover
gohome get-config days
gohome set-config days 7
```

---

## 🎯 gohome Design Rationale

### Design Goals

1. **Simplicity First**: `gohome` should just work
2. **Progressive Complexity**: Easy to start, powerful when needed
3. **Flexible Configuration**: Save preferences, override temporarily
4. **Standard Patterns**: Follow industry conventions

### Architecture Summary

```
gohome                      # Default: report with config file defaults
gohome report              # Explicit form (identical)
gohome report --days=7     # Override specific flag
gohome --save              # Persist current flags to config
gohome config list         # View configuration
gohome config set key val  # Edit configuration
gohome version             # Version info
```

### Why This Works

1. **Intuitive**: New users can type `gohome` and it works
2. **Discoverable**: `--help` reveals all features
3. **Flexible**: Power users can customize extensively
4. **Consistent**: Follows patterns from git, kubectl, gh
5. **Maintainable**: Clear command structure, easy to extend

### Future Extensibility

Easy to add new commands:

```bash
gohome export --format=json     # Export to other formats
gohome analyze                  # Statistical analysis
gohome hook install             # Git hook integration
gohome plugin list              # Plugin system
```

Easy to add config subcommands:

```bash
gohome config export            # Export config
gohome config import file.json  # Import config
gohome config validate          # Validate config
```

---

## 📊 Comparison Matrix

| Tool | Primary Pattern | Config Approach | Default Cmd | Subcommands |
|------|----------------|-----------------|-------------|-------------|
| **git** | Flat commands | File-first | ❌ | Minimal |
| **kubectl** | Resource + Action | File + Flags | ❌ | Heavy |
| **docker** | Mixed | Flags-first | ❌ | Moderate |
| **gh** | Resource + Action | File + Flags | ❌ | Heavy |
| **cargo** | Action commands | File-first | ❌ | Minimal |
| **gohome** | Default + Subcommands | Hybrid | ✅ | Moderate |

---

## ✅ Best Practices Summary

### ✅ DO:

1. **Provide a sensible default command** if there's one primary use case
2. **Use subcommands for related operations** (config, plugin management)
3. **Support both short and long flag forms** (`-d`, `--days`)
4. **Implement configuration hierarchy** (flags > env > file > defaults)
5. **Show progress for long operations** (spinners, progress bars)
6. **Provide helpful error messages** with suggestions
7. **Auto-generate shell completions**
8. **Make help text comprehensive**
9. **Follow verb + noun or noun + verb consistently**
10. **Use emoji/icons sparingly** for visual clarity

### ❌ DON'T:

1. **Don't require subcommand for primary operation** (if there's a clear default)
2. **Don't mix command patterns** (pick resource-first OR action-first)
3. **Don't auto-save config silently** (require explicit `--save`)
4. **Don't use ambiguous flag names** (`-c` for multiple things)
5. **Don't hide important operations** in obscure flags
6. **Don't make the config file required** (defaults should work)
7. **Don't break existing configs** in updates (migration paths)
8. **Don't use inconsistent naming** (camelCase vs kebab-case)

---

## 📚 Further Reading

- [Command Line Interface Guidelines](https://clig.dev/)
- [12 Factor CLI Apps](https://medium.com/@jdxcode/12-factor-cli-apps-dd3c227a0e46)
- [Heroku CLI Style Guide](https://devcenter.heroku.com/articles/cli-style-guide)
- [Google's Command Line Flags Best Practices](https://google.github.io/styleguide/shellguide.html#flags)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)

---

**This document is a living guide** - as CLI design patterns evolve and gohome grows, this will be updated to reflect new learnings and best practices.
