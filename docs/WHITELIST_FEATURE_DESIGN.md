# 📋 Repository Whitelist/Favorites Feature Design

**Status:** Planning Phase  
**Target Release:** v1.4.0 (Phase 2)  
**Priority:** High - Daily workflow enhancement

## 📋 Overview

The **Repository Whitelist** (also called "Favorites") feature allows users to define a curated list of repositories they're actively working on. This solves the common problem of scanning hundreds of repos when you only care about 5-10 active projects for daily standup reports.

## 🎯 Problem Statement

**Current Behavior:**
- gohome scans ALL repositories in workspace directories
- Users with many repos (50+) experience:
  - Slow scan times
  - Cluttered output with irrelevant repos
  - Difficulty finding active projects
  - Noisy reports with archived/inactive repos

**User Scenarios:**
1. **Daily Standup:** "I only want to report on my 3 active projects, not all 47 repos"
2. **Context Switching:** "I work on backend team repos this week, frontend next week"
3. **Client Projects:** "Show me only client X repos, hide internal tools"
4. **Learning:** "Track my progress on tutorial repos separately from work"

## ✨ Proposed Solution

### Feature Overview

A **whitelist system** that allows users to:
- Mark specific repositories as "favorites"
- Scan ONLY whitelisted repos with a flag
- Manage whitelist via CLI sub-commands
- Persist whitelist in config file
- Override scan scope without changing workspace paths

### Key Benefits

- ⚡ **Faster Scans:** Skip irrelevant directories entirely
- 🎯 **Focused Reports:** See only what matters today
- 🔄 **Flexible Workflows:** Switch contexts easily
- 💾 **Persistent State:** Whitelist saved across sessions
- 🚀 **Zero Breaking Changes:** Optional feature, doesn't affect default behavior

## 🏗️ Architecture Design

### Data Model

**Config File Structure:**

```json
{
  "whitelist": {
    "enabled": false,
    "repos": [
      {
        "path": "/home/user/workspace/github.com/org/project-a",
        "alias": "Project A",
        "added": "2026-01-20T10:30:00Z",
        "tags": ["backend", "active", "team-alpha"]
      },
      {
        "path": "/home/user/workspace/github.com/org/project-b",
        "alias": "Project B",
        "added": "2026-01-18T14:22:00Z",
        "tags": ["frontend", "client-x"]
      }
    ],
    "scanMode": "whitelist-only"
  }
}
```

**Entity Definition:**

```go
// internal/entity/whitelist.go
package entity

import "time"

type Whitelist struct {
    Enabled  bool             `json:"enabled" yaml:"enabled"`
    Repos    []WhitelistRepo  `json:"repos" yaml:"repos"`
    ScanMode string           `json:"scanMode" yaml:"scanMode"` // "all", "whitelist-only", "blacklist"
}

type WhitelistRepo struct {
    Path    string    `json:"path" yaml:"path"`         // Absolute path to repo
    Alias   string    `json:"alias,omitempty" yaml:"alias,omitempty"` // User-friendly name
    Added   time.Time `json:"added" yaml:"added"`       // When added to whitelist
    Tags    []string  `json:"tags,omitempty" yaml:"tags,omitempty"` // Categorization tags
}
```

### CLI Interface

#### Sub-commands

```bash
# Whitelist Management Sub-commands
gohome whitelist add <path> [--alias <name>] [--tags <tag1,tag2>]
gohome whitelist remove <path>
gohome whitelist list [--format table|json] [--tags <filter>]
gohome whitelist clear
gohome whitelist enable
gohome whitelist disable
gohome whitelist status

# Scanning with Whitelist
gohome report --whitelist-only          # Scan ONLY whitelisted repos
gohome report --whitelist-mode          # Same as above (alias)
gohome report --tags backend,active     # Scan repos matching tags
gohome report --exclude-whitelist       # Scan everything EXCEPT whitelist (blacklist mode)
```

#### Command Examples

**Add repositories:**
```bash
# Add repo with auto-detected name
gohome whitelist add ~/workspace/github.com/myorg/project-a

# Add with custom alias
gohome whitelist add ~/workspace/project-b --alias "Cool API"

# Add with tags for categorization
gohome whitelist add ~/workspace/frontend-app --tags "frontend,client-x,urgent"

# Bulk add (interactive or from file)
gohome whitelist add --from-file repos.txt
```

**List favorites:**
```bash
# Show all favorited repos
gohome whitelist list

# Output:
# ┏━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━┓
# ┃ ALIAS        ┃ PATH                           ┃ TAGS              ┃
# ┡━━━━━━━━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━━━━━━━━━┩
# │ Project A    │ .../github.com/org/project-a   │ backend, active   │
# │ Project B    │ .../github.com/org/project-b   │ frontend, client  │
# └──────────────┴────────────────────────────────┴───────────────────┘

# Filter by tags
gohome whitelist list --tags backend

# JSON export for scripting
gohome whitelist list --format json
```

**Remove repositories:**
```bash
# Remove by path
gohome whitelist remove ~/workspace/github.com/myorg/project-a

# Remove by alias
gohome whitelist remove "Project A"

# Remove by tag (remove all repos with tag)
gohome whitelist remove --tag archived
```

**Enable/Disable whitelist:**
```bash
# Enable whitelist globally (all scans will use whitelist by default)
gohome whitelist enable

# Disable whitelist (revert to scanning all repos)
gohome whitelist disable

# Check current status
gohome whitelist status
# Output: Whitelist: enabled, 5 repositories
```

**Scan with whitelist:**
```bash
# Scan ONLY whitelisted repos (even if whitelist disabled globally)
gohome report --whitelist-only

# Scan repos matching specific tags
gohome report --tags frontend,urgent

# Combine with other flags
gohome report --whitelist-only --format table --days 2
```

### Scanner Integration

**Enhanced Scanner Logic:**

```go
// internal/scanner/scanner.go
type Scanner struct {
    whitelist    *entity.Whitelist
    whitelistEnabled bool
}

func (s *Scanner) ScanGitRepos(dirs []string, maxDepth int) ([]string, error) {
    // If whitelist-only mode, skip directory scanning entirely
    if s.whitelistEnabled && s.whitelist.ScanMode == "whitelist-only" {
        return s.getWhitelistedRepoPaths(), nil
    }
    
    // Normal scan logic
    repos := s.scanRecursive(dirs, maxDepth)
    
    // Apply whitelist filter if enabled
    if s.whitelistEnabled {
        repos = s.filterByWhitelist(repos)
    }
    
    return repos, nil
}

func (s *Scanner) getWhitelistedRepoPaths() []string {
    paths := make([]string, 0, len(s.whitelist.Repos))
    for _, repo := range s.whitelist.Repos {
        if isValidGitRepo(repo.Path) {
            paths = append(paths, repo.Path)
        }
    }
    return paths
}

func (s *Scanner) filterByWhitelist(repos []string) []string {
    whitelistMap := make(map[string]bool)
    for _, repo := range s.whitelist.Repos {
        whitelistMap[repo.Path] = true
    }
    
    filtered := make([]string, 0)
    for _, repo := range repos {
        if s.whitelist.ScanMode == "whitelist-only" {
            if whitelistMap[repo] {
                filtered = append(filtered, repo)
            }
        } else if s.whitelist.ScanMode == "blacklist" {
            if !whitelistMap[repo] {
                filtered = append(filtered, repo)
            }
        }
    }
    return filtered
}
```

### Tag-Based Filtering

**Use Cases:**
```bash
# Daily standup: show only active projects
gohome report --tags active

# Client work: show only client repos
gohome report --tags client-acme

# Weekly review: show backend and frontend work
gohome report --tags backend,frontend --days 7

# Emergency fix tracking
gohome report --tags hotfix,production
```

**Tag Management:**
```bash
# Add tags to existing whitelist entry
gohome whitelist tag add ~/workspace/project-a --tags urgent,production

# Remove tags
gohome whitelist tag remove ~/workspace/project-a --tags urgent

# List all available tags
gohome whitelist tags
# Output: active, archived, backend, client-x, frontend, urgent
```

## 🎨 UI/UX Design

### Visual Indicators

**Whitelist Status in Output:**

```bash
# When whitelist is enabled
╭─────────────────────────────────────────────╮
│          DAILY STANDUP REPORT               │
│       🔖 Whitelist Mode: ON (3 repos)       │
│           2026-01-20 • @user                │
╰─────────────────────────────────────────────╯

# When showing only tagged repos
╭─────────────────────────────────────────────╮
│          DAILY STANDUP REPORT               │
│    🏷️ Tags: backend, active (2 repos)       │
│           2026-01-20 • @user                │
╰─────────────────────────────────────────────╯
```

**Repository Cards with Whitelist Indicator:**

```bash
┌─ project-a ────────────────────────────── ⭐ ┐
│ 📍 main • 5 commits • last 24h              │
│ 🏷️ backend, active, team-alpha              │
└─────────────────────────────────────────────┘
```

### Interactive Prompts (Optional Enhancement)

```bash
# Interactive add with validation
$ gohome whitelist add ~/workspace/new-project

? Enter alias (optional): Cool New API
? Add tags (comma-separated): backend,experimental
✅ Added "Cool New API" to whitelist
```

### Status Command Output

```bash
$ gohome whitelist status

📊 Whitelist Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Status:          ✅ Enabled
Mode:            Whitelist-only
Total Repos:     5
Last Modified:   2026-01-20 10:30:00

Recent Activity:
  • Project A - 3 commits today
  • Project B - 1 commit today
  • Project C - No commits

💡 Tip: Use 'gohome whitelist list' to see all favorites
```

## 🔧 Developer Experience (DX)

### Configuration Management

**Viper Integration:**

```go
// internal/config/viper/whitelist.go
func LoadWhitelist() (*entity.Whitelist, error) {
    var whitelist entity.Whitelist
    if err := viper.UnmarshalKey("whitelist", &whitelist); err != nil {
        return nil, err
    }
    return &whitelist, nil
}

func SaveWhitelist(whitelist *entity.Whitelist) error {
    viper.Set("whitelist", whitelist)
    return viper.WriteConfig()
}
```

**Environment Variable Support:**

```bash
GOHOME_WHITELIST_ENABLED=true
GOHOME_WHITELIST_SCAN_MODE=whitelist-only
```

### Error Handling

**Common Scenarios:**

1. **Invalid Path:**
   ```bash
   $ gohome whitelist add /invalid/path
   ❌ Error: Path does not exist or is not a git repository
   ```

2. **Duplicate Entry:**
   ```bash
   $ gohome whitelist add ~/workspace/project-a
   ⚠️ Warning: Repository already in whitelist as "Project A"
   ```

3. **Empty Whitelist:**
   ```bash
   $ gohome report --whitelist-only
   ⚠️ Warning: Whitelist is empty. No repositories to scan.
   💡 Tip: Add repositories with 'gohome whitelist add <path>'
   ```

4. **Whitelist Disabled:**
   ```bash
   $ gohome report --whitelist-only
   ℹ️ Note: Whitelist is disabled. Enable with 'gohome whitelist enable'
   ```

### Validation Rules

- Paths must be absolute or expandable (`~` supported)
- Paths must exist and contain `.git` directory
- Aliases must be unique
- Tags must be lowercase, alphanumeric, hyphens allowed
- Maximum 100 repos in whitelist (configurable)

## 🚀 Implementation Plan

### Phase 1: Core Infrastructure (Week 1-2)

1. **Data Models**
   - Define `Whitelist` and `WhitelistRepo` entities
   - Add Viper config integration
   - Create migration for existing configs

2. **Sub-commands**
   - Implement `gohome whitelist add`
   - Implement `gohome whitelist remove`
   - Implement `gohome whitelist list`
   - Add Cobra command structure

3. **Config Persistence**
   - Save/load whitelist from config file
   - Handle file locking for concurrent access
   - Config validation and error handling

### Phase 2: Scanner Integration (Week 3-4)

4. **Scanner Enhancement**
   - Add whitelist-aware scanning logic
   - Implement `--whitelist-only` flag
   - Skip unnecessary directory traversal when using whitelist

5. **Performance Optimization**
   - Cache whitelist lookups
   - Early exit on whitelist-only mode
   - Benchmark scan performance improvements

6. **Testing**
   - Unit tests for whitelist management
   - Integration tests for scanner with whitelist
   - Edge case testing (symlinks, invalid paths, etc.)

### Phase 3: Tag System (Week 5)

7. **Tag Management**
   - Add tags to whitelist entries
   - Implement `--tags` filtering flag
   - Tag autocomplete for shell completion

8. **Tag Commands**
   - `gohome whitelist tag add`
   - `gohome whitelist tag remove`
   - `gohome whitelist tags` (list all tags)

### Phase 4: UI/UX Polish (Week 6)

9. **Visual Enhancements**
   - Whitelist indicator in output
   - Styled whitelist list command
   - Status command with rich output

10. **Interactive Features** (Optional)
    - Interactive selection when adding repos
    - Fuzzy search for repo paths
    - Bulk operations with confirmation

### Phase 5: Documentation (Week 7)

11. **User Documentation**
    - Update CLI_GUIDE.md with whitelist examples
    - Add use case scenarios
    - Create workflow recommendations

12. **Developer Documentation**
    - API documentation for whitelist package
    - Migration guide for config updates
    - Architecture decision records (ADR)

## 📊 Success Metrics

- **Performance:** Scan time reduced by 70%+ when using whitelist-only mode
- **Adoption:** 40%+ of users enable whitelist within first week
- **Usability:** Zero confusion in user feedback
- **Reliability:** No data loss or corruption in whitelist management

## 🔮 Future Enhancements (Phase 3)

### Smart Whitelist

- **Auto-suggest:** Recommend repos based on recent commit activity
- **AI-powered:** "Show me repos I worked on this sprint"
- **Dynamic whitelist:** Auto-add repos when you commit to them

### Advanced Filtering

- **Time-based:** "Show repos active in last 7 days"
- **Author-based:** "Show repos where I'm the main contributor"
- **Commit-based:** "Show repos with >10 commits this week"

### Whitelist Presets

```bash
# Save current whitelist as preset
gohome whitelist preset save "client-work"

# Load preset
gohome whitelist preset load "client-work"

# List presets
gohome whitelist preset list
```

### Team Sharing

```bash
# Export whitelist for team sharing
gohome whitelist export --output team-repos.json

# Import team whitelist
gohome whitelist import team-repos.json --merge
```

## 🔗 Related Features

- Integrates with **UI/UX Enhancement** for styled output
- Foundation for **Interactive Mode** (TUI) in Phase 3
- Complements **Performance & Concurrency** optimizations
- Enables **Smart Summaries** with focused context

## 📝 Open Questions

1. Should whitelist be git-ignored by default or shareable?
2. Default behavior: whitelist-only or normal scan?
3. Maximum recommended whitelist size?
4. Should we support regex patterns for paths?
5. Conflict resolution when same repo has multiple paths?

## 🎯 Acceptance Criteria

- [ ] Users can add/remove repos to whitelist via CLI
- [ ] `--whitelist-only` flag scans only favorited repos
- [ ] Whitelist persists across sessions in config file
- [ ] List command shows all whitelisted repos with formatting
- [ ] Tags support for categorization and filtering
- [ ] Performance: 50%+ faster scan for 10+ repos in whitelist mode
- [ ] Zero breaking changes to existing workflows
- [ ] Comprehensive documentation and examples
- [ ] 80%+ test coverage for whitelist features
- [ ] Shell completion for whitelist sub-commands
