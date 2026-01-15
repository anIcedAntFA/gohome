# 🏗️ Cobra & Viper Migration Guide

**Status:** 📋 Planning Phase  
**Target Version:** v2.0.0  
**Estimated Timeline:** 2-3 months (phased approach)

---

## 📖 Table of Contents

- [Overview](#overview)
- [Why Migrate?](#why-migrate)
- [Architecture Comparison](#architecture-comparison)
- [Migration Roadmap](#migration-roadmap)
- [Phase 1: Foundation](#phase-1-foundation)
- [Phase 2: Core Migration](#phase-2-core-migration)
- [Phase 3: Advanced Features](#phase-3-advanced-features)
- [Phase 4: Testing & Deployment](#phase-4-testing--deployment)
- [Implementation Checklist](#implementation-checklist)
- [Code Examples](#code-examples)
- [Breaking Changes](#breaking-changes)
- [Rollback Strategy](#rollback-strategy)

---

## 🎯 Overview

This guide provides a comprehensive step-by-step plan to migrate **gohome** from the standard `flag` package to **Cobra** (CLI framework) and **Viper** (configuration management).

### Current Architecture (v1.x.x)

```
cmd/gohome/main.go
├── flag package (CLI parsing)
├── internal/config (JSON file + flag merging)
├── internal/scanner (repository discovery)
├── internal/git (git operations)
├── internal/parser (commit parsing)
└── internal/renderer (output formatting)
```

### Target Architecture (v2.0.0)

```
cmd/gohome/main.go
├── cmd/root.go (Cobra root command)
├── cmd/report.go (default command)
├── cmd/config.go (config management subcommand)
├── cmd/version.go (version subcommand)
├── internal/config (Viper-based config)
│   ├── viper.go (config loading)
│   ├── defaults.go (default values)
│   └── validation.go (config validation)
├── internal/scanner (unchanged)
├── internal/git (unchanged)
├── internal/parser (unchanged)
├── internal/renderer (unchanged)
└── internal/plugin (new - plugin system)
    ├── loader.go
    ├── registry.go
    └── interface.go
```

---

## 💡 Why Migrate?

### Problems with Current Architecture

1. **Limited Extensibility:** Hard to add subcommands (e.g., `gohome config list`)
2. **Complex Flag Management:** Manual flag parsing and validation
3. **Config Merging Complexity:** Custom logic for flag precedence
4. **No Environment Variable Support:** Can't use env vars for sensitive data (API keys)
5. **Hard to Test:** Tight coupling between CLI parsing and business logic
6. **Limited Format Support:** Only JSON config, no YAML/TOML

### Benefits of Cobra + Viper

#### Cobra Benefits
- ✅ **Subcommands:** Native support for `gohome config`, `gohome report`, etc.
- ✅ **Auto-generated Help:** Professional help text and documentation
- ✅ **Flag Inheritance:** Parent flags available to child commands
- ✅ **Shell Completion:** Auto-generate bash/zsh/fish completions
- ✅ **Testing:** Easy to test commands in isolation
- ✅ **Industry Standard:** Used by kubectl, hugo, GitHub CLI, etc.

#### Viper Benefits
- ✅ **Multi-format Config:** JSON, YAML, TOML, HCL, env files
- ✅ **Environment Variables:** Automatic env var binding
- ✅ **Config Hierarchy:** Flags > Env > Config > Defaults (automatic)
- ✅ **Live Reloading:** Watch config file for changes
- ✅ **Type-safe Access:** Strongly typed getters
- ✅ **Remote Config:** Support for etcd, Consul (future)

---

## 🏛️ Architecture Comparison

### Current: Flag-based Architecture

```go
// cmd/gohome/main.go
func main() {
    // Manual flag definition
    flag.IntVar(&cfg.Days, "days", 0, "")
    flag.IntVar(&cfg.Days, "d", 0, "")
    flag.StringVar(&cfg.Path, "path", ".", "")
    flag.Parse()

    // Manual config merging
    fileCfg := config.LoadFromFile()
    config.Merge(&cfg, &fileCfg)

    // Direct execution
    scanner.Scan()
    git.GetLogs()
    renderer.Print()
}
```

**Issues:**
- ❌ No subcommands
- ❌ Duplicate flag definitions
- ❌ Complex merging logic
- ❌ Hard to extend
- ❌ No environment variables

### Target: Cobra + Viper Architecture

```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "gohome",
    Short: "Git standup & activity reporting CLI",
    RunE:  runReport, // Default command
}

func init() {
    cobra.OnInitialize(initConfig)
    
    // Define flags once
    rootCmd.PersistentFlags().IntP("days", "d", 1, "Number of days to look back")
    rootCmd.PersistentFlags().StringP("path", "p", ".", "Root path to scan")
    
    // Bind to viper (automatic env var support)
    viper.BindPFlag("days", rootCmd.PersistentFlags().Lookup("days"))
    viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path"))
}

func initConfig() {
    // Auto-detect config file
    viper.SetConfigName(".gohome")
    viper.SetConfigType("json") // or yaml, toml
    viper.AddConfigPath("$HOME")
    
    // Environment variable support
    viper.SetEnvPrefix("GOHOME")
    viper.AutomaticEnv()
    
    // Read config
    viper.ReadInConfig()
}

// cmd/config.go
var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage gohome configuration",
}

var configListCmd = &cobra.Command{
    Use:   "list",
    Short: "Show current configuration",
    RunE:  runConfigList,
}
```

**Benefits:**
- ✅ Subcommands: `gohome`, `gohome config list`, `gohome version`
- ✅ Single flag definition
- ✅ Automatic config hierarchy
- ✅ Environment variables: `GOHOME_DAYS=3 gohome`
- ✅ Easy to extend

---

## 🗺️ Migration Roadmap

### Timeline Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Migration Timeline                          │
├─────────────────────────────────────────────────────────────────────┤
│ Phase 1 (Week 1-2):  Foundation & Setup                            │
│   └─ Install deps, create skeleton, parallel development           │
│                                                                     │
│ Phase 2 (Week 3-6):  Core Migration                                │
│   ├─ Migrate root command                                          │
│   ├─ Migrate config system                                         │
│   └─ Add subcommands                                               │
│                                                                     │
│ Phase 3 (Week 7-9):  Advanced Features                             │
│   ├─ Environment variables                                         │
│   ├─ Shell completions                                             │
│   └─ Plugin system                                                 │
│                                                                     │
│ Phase 4 (Week 10-12): Testing & Deployment                         │
│   ├─ Comprehensive testing                                         │
│   ├─ Migration guide for users                                     │
│   └─ v2.0.0 release                                                │
└─────────────────────────────────────────────────────────────────────┘
```

### Branching Strategy

```
main (v1.x.x - stable)
  │
  ├── v2-cobra-migration (development branch)
  │     │
  │     ├── feature/cobra-skeleton
  │     ├── feature/viper-config
  │     ├── feature/subcommands
  │     ├── feature/plugins
  │     └── release/v2.0.0-beta.1
  │
  └── merge to main when stable
```

---

## 📋 Phase 1: Foundation (Week 1-2)

**Goal:** Set up Cobra/Viper skeleton without breaking existing functionality.

### 1.1 Create Development Branch

```bash
# Create and switch to migration branch
git checkout -b v2-cobra-migration

# Create feature branch
git checkout -b feature/cobra-skeleton
```

### 1.2 Install Dependencies

```bash
# Add Cobra
go get -u github.com/spf13/cobra@latest

# Add Viper
go get -u github.com/spf13/viper@latest

# Update go.mod
go mod tidy
```

**Verify versions:**
```bash
go list -m github.com/spf13/cobra
go list -m github.com/spf13/viper
```

### 1.3 Create Directory Structure

```bash
# Create cmd directory structure
mkdir -p cmd/gohome/cmd
mkdir -p internal/config/viper
mkdir -p internal/plugin

# Create skeleton files
touch cmd/gohome/cmd/root.go
touch cmd/gohome/cmd/report.go
touch cmd/gohome/cmd/config.go
touch cmd/gohome/cmd/version.go
touch internal/config/viper/viper.go
```

### 1.4 Initialize Cobra Application

**Create `cmd/gohome/cmd/root.go`:**

```go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
    Use:   "gohome",
    Short: "Git standup & activity reporting CLI",
    Long: `gohome scans your workspace for git repositories and generates
formatted daily standup reports from commit history.

Perfect for Daily Standups, Weekly Summaries, or tracking your Personal Coding Habits.`,
    Version: "2.0.0-dev",
}

// Execute runs the root command
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    // Persistent flags (available to all subcommands)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gohome.json)")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

        viper.AddConfigPath(home)
        viper.SetConfigType("json")
        viper.SetConfigName(".gohome")
    }

    // Environment variable support
    viper.SetEnvPrefix("GOHOME")
    viper.AutomaticEnv()

    // Read config file (ignore error if not found)
    if err := viper.ReadInConfig(); err == nil {
        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}
```

**Update `cmd/gohome/main.go`:**

```go
package main

import "github.com/anIcedAntFA/gohome/cmd/gohome/cmd"

func main() {
    cmd.Execute()
}
```

### 1.5 Parallel Development Setup

**Keep v1.x working while developing v2:**

```go
// cmd/gohome/main_v1.go (backup of current main.go)
// +build v1

package main

// ... existing v1 code ...
```

**Build scripts:**

```bash
# Build v1 (current stable)
go build -tags=v1 -o bin/gohome-v1 ./cmd/gohome

# Build v2 (development)
go build -o bin/gohome-v2 ./cmd/gohome
```

### ✅ Phase 1 Checklist

- [ ] Create `v2-cobra-migration` branch
- [ ] Install Cobra and Viper dependencies
- [ ] Create directory structure (`cmd/gohome/cmd/`)
- [ ] Create `root.go` with basic Cobra setup
- [ ] Create `initConfig()` with Viper initialization
- [ ] Update `main.go` to call `cmd.Execute()`
- [ ] Verify app compiles: `go build ./cmd/gohome`
- [ ] Test basic command: `./gohome --help`
- [ ] Commit: `git commit -m "🏗️ feat(arch): add Cobra/Viper skeleton"`

---

## 🔨 Phase 2: Core Migration (Week 3-6)

**Goal:** Migrate core functionality to Cobra/Viper while maintaining compatibility.

### 2.1 Migrate Report Command (Default)

**Create `cmd/gohome/cmd/report.go`:**

```go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    
    "github.com/anIcedAntFA/gohome/internal/config"
    "github.com/anIcedAntFA/gohome/internal/scanner"
    "github.com/anIcedAntFA/gohome/internal/git"
    "github.com/anIcedAntFA/gohome/internal/parser"
    "github.com/anIcedAntFA/gohome/internal/renderer"
)

var reportCmd = &cobra.Command{
    Use:   "report",
    Short: "Generate activity report (default command)",
    Long:  `Scan repositories and generate a formatted activity report.`,
    RunE:  runReport,
}

func init() {
    rootCmd.AddCommand(reportCmd)
    
    // Time period flags
    reportCmd.Flags().IntP("days", "d", 1, "Number of days to look back")
    reportCmd.Flags().IntP("hours", "H", 0, "Number of hours to look back")
    reportCmd.Flags().IntP("weeks", "w", 0, "Number of weeks to look back")
    reportCmd.Flags().IntP("months", "m", 0, "Number of months to look back")
    reportCmd.Flags().IntP("years", "y", 0, "Number of years to look back")
    reportCmd.Flags().Bool("today", false, "Report from midnight to now")
    
    // Path and author
    reportCmd.Flags().StringP("path", "p", ".", "Root path to scan for repositories")
    reportCmd.Flags().IntP("max-depth", "", 2, "Maximum depth to scan for repositories")
    reportCmd.Flags().StringP("author", "a", "", "Git author name (auto-detected if empty)")
    
    // Output formatting
    reportCmd.Flags().StringP("format", "f", "text", "Output format: text, table")
    reportCmd.Flags().StringP("style", "s", "normal", "Table style: normal, markdown")
    reportCmd.Flags().BoolP("icon", "i", false, "Show icon column")
    reportCmd.Flags().BoolP("scope", "c", false, "Show scope column")
    
    // Branch filtering
    reportCmd.Flags().BoolP("all-branches", "b", false, "Include commits from all local branches")
    reportCmd.Flags().String("branch", "", "Filter commits by specific branch")
    
    // Clipboard and tasks
    reportCmd.Flags().BoolP("copy", "cp", false, "Copy output to clipboard")
    reportCmd.Flags().StringSliceP("task", "t", []string{}, "Add custom task")
    
    // Bind all flags to viper
    viper.BindPFlags(reportCmd.Flags())
    
    // Make report the default command if no subcommand specified
    rootCmd.RunE = runReport
}

func runReport(cmd *cobra.Command, args []string) error {
    // Load configuration (Viper handles flag > env > config > default)
    cfg := config.LoadFromViper()
    
    // Business logic (same as v1)
    repos, err := scanner.ScanGitRepos(cfg.Path, cfg.MaxDepth)
    if err != nil {
        return err
    }
    
    // ... rest of the logic ...
    
    return nil
}
```

### 2.2 Migrate Config System to Viper

**Create `internal/config/viper/viper.go`:**

```go
package viper

import (
    "github.com/spf13/viper"
    "github.com/anIcedAntFA/gohome/internal/entity"
)

// Config represents the application configuration
type Config struct {
    // Time period
    Hours  int  `mapstructure:"hours"`
    Days   int  `mapstructure:"days"`
    Weeks  int  `mapstructure:"weeks"`
    Months int  `mapstructure:"months"`
    Years  int  `mapstructure:"years"`
    Today  bool `mapstructure:"today"`
    
    // Path and scanning
    Path     string `mapstructure:"path"`
    MaxDepth int    `mapstructure:"max_depth"`
    Author   string `mapstructure:"author"`
    
    // Output
    Format    string `mapstructure:"format"`
    Preset    string `mapstructure:"preset"`
    ShowIcon  bool   `mapstructure:"show_icon"`
    ShowScope bool   `mapstructure:"show_scope"`
    
    // Branch filtering
    AllBranches bool   `mapstructure:"all_branches"`
    Branch      string `mapstructure:"branch"`
    
    // Clipboard
    CopyToClipboard bool `mapstructure:"copy_to_clipboard"`
    
    // Tasks
    Tasks        []entity.Task `mapstructure:"tasks"`
    DynamicTasks []string      `mapstructure:"-"`
}

// LoadFromViper loads config with automatic hierarchy
func LoadFromViper() *Config {
    var cfg Config
    
    // Unmarshal into struct (Viper handles precedence automatically)
    if err := viper.Unmarshal(&cfg); err != nil {
        // Handle error
        return &Config{}
    }
    
    // Auto-detect author if not set
    if cfg.Author == "" {
        cfg.Author = detectGitAuthor()
    }
    
    return &cfg
}

// SaveToFile saves current config to file
func (c *Config) SaveToFile() error {
    viper.Set("hours", c.Hours)
    viper.Set("days", c.Days)
    // ... set all fields ...
    
    return viper.WriteConfig()
}

func detectGitAuthor() string {
    // Same logic as v1
    return ""
}
```

### 2.3 Add Config Subcommand

**Create `cmd/gohome/cmd/config.go`:**

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var configCmd = &cobra.Command{
    Use:   "config",
    Short: "Manage gohome configuration",
    Long:  `View, edit, and manage gohome configuration settings.`,
}

var configListCmd = &cobra.Command{
    Use:   "list",
    Short: "Show current configuration",
    RunE:  runConfigList,
}

var configGetCmd = &cobra.Command{
    Use:   "get <key>",
    Short: "Get a configuration value",
    Args:  cobra.ExactArgs(1),
    RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
    Use:   "set <key> <value>",
    Short: "Set a configuration value",
    Args:  cobra.ExactArgs(2),
    RunE:  runConfigSet,
}

var configResetCmd = &cobra.Command{
    Use:   "reset",
    Short: "Reset configuration to defaults",
    RunE:  runConfigReset,
}

func init() {
    rootCmd.AddCommand(configCmd)
    configCmd.AddCommand(configListCmd)
    configCmd.AddCommand(configGetCmd)
    configCmd.AddCommand(configSetCmd)
    configCmd.AddCommand(configResetCmd)
}

func runConfigList(cmd *cobra.Command, args []string) error {
    fmt.Println("Current Configuration:")
    fmt.Println("=====================")
    
    allSettings := viper.AllSettings()
    for key, value := range allSettings {
        fmt.Printf("%-20s = %v\n", key, value)
    }
    
    return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
    key := args[0]
    value := viper.Get(key)
    
    if value == nil {
        return fmt.Errorf("key %q not found", key)
    }
    
    fmt.Printf("%s = %v\n", key, value)
    return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
    key, value := args[0], args[1]
    
    viper.Set(key, value)
    
    if err := viper.WriteConfig(); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }
    
    fmt.Printf("Set %s = %s\n", key, value)
    return nil
}

func runConfigReset(cmd *cobra.Command, args []string) error {
    // Prompt for confirmation
    fmt.Print("Are you sure you want to reset configuration? (y/N): ")
    var response string
    fmt.Scanln(&response)
    
    if response != "y" && response != "Y" {
        fmt.Println("Cancelled.")
        return nil
    }
    
    // Delete config file
    configFile := viper.ConfigFileUsed()
    if err := os.Remove(configFile); err != nil {
        return fmt.Errorf("failed to delete config: %w", err)
    }
    
    fmt.Println("Configuration reset to defaults.")
    return nil
}
```

### 2.4 Update Version Command

**Create `cmd/gohome/cmd/version.go`:**

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/anIcedAntFA/gohome/internal/version"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Long:  `Display the version, commit hash, and build date of gohome.`,
    Run:   runVersion,
}

func init() {
    rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
    fmt.Printf("gohome version %s\n", version.Version)
    fmt.Printf("commit: %s\n", version.Commit)
    fmt.Printf("built: %s\n", version.Date)
}
```

### ✅ Phase 2 Checklist

- [ ] Create `report.go` with all flags migrated
- [ ] Bind all flags to Viper using `viper.BindPFlags()`
- [ ] Create `internal/config/viper/viper.go`
- [ ] Implement `LoadFromViper()` function
- [ ] Test config hierarchy: flag > env > config > default
- [ ] Create `config.go` subcommand
- [ ] Implement `config list`, `config get`, `config set`, `config reset`
- [ ] Update `version.go` command
- [ ] Test all commands work: `gohome`, `gohome report`, `gohome config list`
- [ ] Verify backward compatibility with v1 config files
- [ ] Test environment variables: `GOHOME_DAYS=5 gohome`
- [ ] Write unit tests for new cmd package
- [ ] Commit: `git commit -m "🏗️ feat(arch): migrate core commands to Cobra/Viper"`

---

## 🚀 Phase 3: Advanced Features (Week 7-9)

**Goal:** Add advanced features enabled by Cobra/Viper.

### 3.1 Environment Variable Support

**Document all supported env vars:**

```bash
# Time period
export GOHOME_DAYS=3
export GOHOME_HOURS=24

# Path and scanning
export GOHOME_PATH=/path/to/workspace
export GOHOME_MAX_DEPTH=3

# Output
export GOHOME_FORMAT=table
export GOHOME_PRESET=markdown

# Sensitive data (future AI features)
export GOHOME_OPENAI_API_KEY=sk-...
export GOHOME_ANTHROPIC_API_KEY=sk-ant-...
```

**Add to README.md:**

```markdown
## Environment Variables

All configuration can be set via environment variables with the `GOHOME_` prefix:

| Config Key   | Environment Variable    | Example                  |
|--------------|-------------------------|--------------------------|
| `days`       | `GOHOME_DAYS`          | `GOHOME_DAYS=3`         |
| `path`       | `GOHOME_PATH`          | `GOHOME_PATH=~/workspace`|
| `format`     | `GOHOME_FORMAT`        | `GOHOME_FORMAT=table`   |
```

### 3.2 Shell Completions

**Add completion command:**

```go
// cmd/gohome/cmd/completion.go
package cmd

import (
    "os"
    "github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion scripts",
    Long: `Generate shell completion scripts for gohome.

Examples:
  # Bash
  source <(gohome completion bash)
  
  # Zsh
  gohome completion zsh > "${fpath[1]}/_gohome"
  
  # Fish
  gohome completion fish > ~/.config/fish/completions/gohome.fish
  
  # PowerShell
  gohome completion powershell | Out-String | Invoke-Expression
`,
    DisableFlagsInUseLine: true,
    ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
    Args:                  cobra.ExactArgs(1),
    RunE:                  runCompletion,
}

func init() {
    rootCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, args []string) error {
    switch args[0] {
    case "bash":
        return rootCmd.GenBashCompletion(os.Stdout)
    case "zsh":
        return rootCmd.GenZshCompletion(os.Stdout)
    case "fish":
        return rootCmd.GenFishCompletion(os.Stdout, true)
    case "powershell":
        return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
    default:
        return fmt.Errorf("unsupported shell: %s", args[0])
    }
}
```

### 3.3 Plugin System Foundation

**Create plugin interface:**

```go
// internal/plugin/interface.go
package plugin

import "github.com/anIcedAntFA/gohome/internal/entity"

// Plugin represents a gohome plugin
type Plugin interface {
    // Name returns the plugin name
    Name() string
    
    // Version returns the plugin version
    Version() string
    
    // Init initializes the plugin with config
    Init(config map[string]interface{}) error
    
    // Execute runs the plugin
    Execute(commits []entity.Commit) error
}

// ExporterPlugin exports reports in custom formats
type ExporterPlugin interface {
    Plugin
    Export(commits []entity.Commit, output string) error
}

// FilterPlugin filters commits
type FilterPlugin interface {
    Plugin
    Filter(commits []entity.Commit) ([]entity.Commit, error)
}
```

**Create plugin loader:**

```go
// internal/plugin/loader.go
package plugin

import (
    "fmt"
    "path/filepath"
    "plugin"
)

// Loader loads plugins from files
type Loader struct {
    pluginDir string
    plugins   map[string]Plugin
}

// NewLoader creates a new plugin loader
func NewLoader(pluginDir string) *Loader {
    return &Loader{
        pluginDir: pluginDir,
        plugins:   make(map[string]Plugin),
    }
}

// Load loads a plugin from file
func (l *Loader) Load(name string) (Plugin, error) {
    // Check if already loaded
    if p, ok := l.plugins[name]; ok {
        return p, nil
    }
    
    // Build plugin path
    pluginPath := filepath.Join(l.pluginDir, name+".so")
    
    // Load plugin
    p, err := plugin.Open(pluginPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open plugin: %w", err)
    }
    
    // Lookup Plugin symbol
    symPlugin, err := p.Lookup("Plugin")
    if err != nil {
        return nil, fmt.Errorf("plugin missing Plugin symbol: %w", err)
    }
    
    // Assert type
    var gohomePlugin Plugin
    gohomePlugin, ok := symPlugin.(Plugin)
    if !ok {
        return nil, fmt.Errorf("invalid plugin type")
    }
    
    // Cache and return
    l.plugins[name] = gohomePlugin
    return gohomePlugin, nil
}

// LoadAll loads all plugins from directory
func (l *Loader) LoadAll() error {
    matches, err := filepath.Glob(filepath.Join(l.pluginDir, "*.so"))
    if err != nil {
        return err
    }
    
    for _, match := range matches {
        name := filepath.Base(match)
        name = name[:len(name)-3] // Remove .so
        
        if _, err := l.Load(name); err != nil {
            // Log error but continue
            fmt.Printf("Warning: failed to load plugin %s: %v\n", name, err)
        }
    }
    
    return nil
}
```

### 3.4 Multi-format Config Support

**Add YAML support:**

```bash
# Install YAML dependency (already included in Viper)
go get gopkg.in/yaml.v3
```

**Update config detection:**

```go
// cmd/gohome/cmd/root.go
func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }

        viper.AddConfigPath(home)
        
        // Support multiple formats
        viper.SetConfigName(".gohome")
        // Viper will try: .gohome.json, .gohome.yaml, .gohome.toml
    }

    viper.SetEnvPrefix("GOHOME")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
        fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
    }
}
```

**Example `.gohome.yaml`:**

```yaml
# gohome configuration (YAML format)
days: 1
path: /Users/ngockhoi96/workspace
format: table
preset: markdown
max_depth: 2

tasks:
  - type: meeting
    message: Daily Standup & Team Sync
    icon: 📅
    enabled: true
  - type: review
    message: Code Review & PR Feedback
    icon: 👀
    enabled: true
```

### ✅ Phase 3 Checklist

- [ ] Document all environment variables in README
- [ ] Test env var precedence: `GOHOME_DAYS=5 gohome -d 3` (flag should win)
- [ ] Create `completion.go` command
- [ ] Generate and test completions for bash, zsh, fish
- [ ] Create plugin interface (`internal/plugin/interface.go`)
- [ ] Create plugin loader (`internal/plugin/loader.go`)
- [ ] Create example plugin (exporter)
- [ ] Add YAML config support
- [ ] Add TOML config support (optional)
- [ ] Test multi-format config detection
- [ ] Write plugin documentation
- [ ] Commit: `git commit -m "✨ feat(arch): add shell completions and plugin system"`

---

## 🧪 Phase 4: Testing & Deployment (Week 10-12)

**Goal:** Comprehensive testing, documentation, and smooth v2.0.0 release.

### 4.1 Testing Strategy

**Unit Tests:**

```go
// cmd/gohome/cmd/report_test.go
package cmd

import (
    "bytes"
    "testing"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

func TestReportCommand(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {
            name:    "default flags",
            args:    []string{},
            wantErr: false,
        },
        {
            name:    "with days flag",
            args:    []string{"--days", "3"},
            wantErr: false,
        },
        {
            name:    "with multiple flags",
            args:    []string{"-d", "2", "-f", "table", "-s", "markdown"},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Reset viper
            viper.Reset()
            
            // Create command
            cmd := &cobra.Command{}
            cmd.SetArgs(tt.args)
            
            // Execute
            err := reportCmd.RunE(cmd, tt.args)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("reportCmd.RunE() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Integration Tests:**

```go
// test/integration/cobra_test.go
package integration

import (
    "os"
    "os/exec"
    "testing"
)

func TestCobraCommands(t *testing.T) {
    // Build binary
    cmd := exec.Command("go", "build", "-o", "/tmp/gohome-test", "./cmd/gohome")
    if err := cmd.Run(); err != nil {
        t.Fatalf("failed to build: %v", err)
    }
    defer os.Remove("/tmp/gohome-test")

    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {"help", []string{"--help"}, false},
        {"version", []string{"version"}, false},
        {"config list", []string{"config", "list"}, false},
        {"report", []string{"report", "-d", "1"}, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command("/tmp/gohome-test", tt.args...)
            output, err := cmd.CombinedOutput()
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Command %v failed: %v\nOutput: %s", tt.args, err, output)
            }
        })
    }
}
```

### 4.2 Migration Guide for Users

**Create `docs/V2_MIGRATION_GUIDE.md`:**

```markdown
# Migrating from v1.x to v2.0

## Overview

Version 2.0 introduces Cobra/Viper architecture with new features and some breaking changes.

## What's New

✨ **Subcommands:** `gohome config list`, `gohome version`
✨ **Environment Variables:** `GOHOME_DAYS=3 gohome`
✨ **Shell Completions:** Auto-complete for bash/zsh/fish
✨ **Multi-format Config:** JSON, YAML, TOML support
✨ **Plugin System:** Extend gohome with custom plugins

## Breaking Changes

### 1. Default Command Behavior

**v1.x:**
```bash
gohome         # Run report
gohome --help  # Show help
```

**v2.0:**
```bash
gohome         # Run report (same)
gohome report  # Explicit report command (new)
gohome --help  # Show help (same)
```

### 2. Config File Location

**v1.x:**
- Only `~/.gohome.json` supported

**v2.0:**
- `~/.gohome.json` (JSON)
- `~/.gohome.yaml` (YAML)
- `~/.gohome.toml` (TOML)

Viper automatically detects format.

### 3. Version Flag

**v1.x:**
```bash
gohome --version
gohome -v
```

**v2.0:**
```bash
gohome version      # Preferred
gohome --version    # Still works
```

## Migration Steps

1. **Backup your config:**
   ```bash
   cp ~/.gohome.json ~/.gohome.json.backup
   ```

2. **Install v2.0:**
   ```bash
   curl -sS https://get.ngockhoi96.dev/gohome | sh
   ```

3. **Verify config compatibility:**
   ```bash
   gohome config list
   ```

4. **(Optional) Convert to YAML:**
   ```bash
   gohome config convert yaml
   ```

5. **Test your workflow:**
   ```bash
   gohome -d 1
   gohome config list
   ```

## Compatibility

- ✅ **v1.x config files work in v2.0** (backward compatible)
- ✅ **All v1.x flags work in v2.0**
- ✅ **No data loss during migration**

## Rollback

If you need to rollback to v1.x:

```bash
# Download v1.x binary
curl -L https://github.com/anIcedAntFA/gohome/releases/download/v1.2.0/gohome_1.2.0_linux_x86_64.tar.gz -o gohome.tar.gz
tar -xzf gohome.tar.gz
sudo mv gohome /usr/local/bin/

# Restore config
cp ~/.gohome.json.backup ~/.gohome.json
```

## FAQ

**Q: Do I need to change my config file?**
A: No, v1.x JSON configs work in v2.0.

**Q: Can I use both v1.x and v2.0?**
A: Yes, install them as `gohome-v1` and `gohome-v2`.

**Q: Are there performance improvements?**
A: Yes, v2.0 is optimized with better error handling and cleaner architecture.

**Q: When should I migrate?**
A: v2.0 is stable and recommended for all users. v1.x will receive security updates only.
```

### 4.3 Update Documentation

**Update README.md:**
- Add subcommands section
- Add environment variables section
- Add shell completions section
- Update examples with new commands

**Update CONTRIBUTING.md:**
- Add Cobra/Viper development guidelines
- Add plugin development guide
- Update testing instructions

### 4.4 Release Process

**Pre-release checklist:**

```bash
# 1. Run all tests
make test
make lint

# 2. Build for all platforms
goreleaser build --snapshot --clean

# 3. Test binaries
./dist/gohome_linux_amd64/gohome version
./dist/gohome_darwin_amd64/gohome config list

# 4. Create beta release
git tag -a v2.0.0-beta.1 -m "v2.0.0 beta 1"
git push origin v2.0.0-beta.1

# 5. Gather feedback (2 weeks)
# ...

# 6. Final release
git tag -a v2.0.0 -m "v2.0.0: Cobra/Viper architecture"
git push origin v2.0.0
```

### ✅ Phase 4 Checklist

- [ ] Write unit tests for all cmd files
- [ ] Write integration tests
- [ ] Achieve >80% test coverage
- [ ] Create V2_MIGRATION_GUIDE.md
- [ ] Update README.md with v2.0 features
- [ ] Update CONTRIBUTING.md
- [ ] Test backward compatibility with v1.x configs
- [ ] Create v2.0.0-beta.1 release
- [ ] Gather community feedback (2 weeks)
- [ ] Fix reported issues
- [ ] Create v2.0.0 final release
- [ ] Announce on GitHub, social media
- [ ] Monitor for issues and respond quickly
- [ ] Commit: `git commit -m "🚀 release: v2.0.0 with Cobra/Viper architecture"`

---

## 📝 Implementation Checklist

### Complete Migration Checklist

```markdown
## Phase 1: Foundation (Week 1-2)
- [ ] Create v2-cobra-migration branch
- [ ] Install Cobra (github.com/spf13/cobra)
- [ ] Install Viper (github.com/spf13/viper)
- [ ] Create cmd/gohome/cmd/ directory structure
- [ ] Create root.go with Cobra root command
- [ ] Implement initConfig() with Viper
- [ ] Update main.go to call cmd.Execute()
- [ ] Test: `gohome --help` works
- [ ] Commit: "🏗️ feat(arch): add Cobra/Viper skeleton"

## Phase 2: Core Migration (Week 3-6)
- [ ] Create report.go command
- [ ] Migrate all flags from flag package
- [ ] Bind flags to Viper with viper.BindPFlags()
- [ ] Create internal/config/viper package
- [ ] Implement LoadFromViper() function
- [ ] Test config hierarchy (flag > env > config)
- [ ] Create config.go subcommand
- [ ] Implement config list subcommand
- [ ] Implement config get subcommand
- [ ] Implement config set subcommand
- [ ] Implement config reset subcommand
- [ ] Update version.go command
- [ ] Test: `gohome report -d 3` works
- [ ] Test: `GOHOME_DAYS=5 gohome` works
- [ ] Test: v1.x config files still work
- [ ] Write unit tests for cmd package
- [ ] Commit: "🏗️ feat(arch): migrate core commands to Cobra/Viper"

## Phase 3: Advanced Features (Week 7-9)
- [ ] Document all environment variables
- [ ] Test env var precedence
- [ ] Create completion.go command
- [ ] Generate bash completion
- [ ] Generate zsh completion
- [ ] Generate fish completion
- [ ] Generate PowerShell completion
- [ ] Test completions in each shell
- [ ] Create plugin/interface.go
- [ ] Create plugin/loader.go
- [ ] Create plugin/registry.go
- [ ] Write example exporter plugin
- [ ] Test plugin loading
- [ ] Add YAML config support
- [ ] Add TOML config support
- [ ] Test multi-format detection
- [ ] Write plugin development guide
- [ ] Commit: "✨ feat(arch): add shell completions and plugin system"

## Phase 4: Testing & Deployment (Week 10-12)
- [ ] Write unit tests for all commands
- [ ] Write integration tests
- [ ] Achieve >80% test coverage
- [ ] Create docs/V2_MIGRATION_GUIDE.md
- [ ] Update README.md
- [ ] Update CONTRIBUTING.md
- [ ] Update .github/copilot-instructions.md
- [ ] Test backward compatibility
- [ ] Build binaries for all platforms
- [ ] Create v2.0.0-beta.1 tag
- [ ] Release beta to community
- [ ] Gather feedback (2 weeks)
- [ ] Fix reported issues
- [ ] Update CHANGELOG.md
- [ ] Create v2.0.0 release
- [ ] Update package managers (AUR, npm, etc.)
- [ ] Announce release
- [ ] Monitor issues
- [ ] Commit: "🚀 release: v2.0.0 with Cobra/Viper architecture"
```

---

## 💻 Code Examples

### Example 1: Using Subcommands

```bash
# Default command (run report)
gohome

# Explicit report command
gohome report -d 3 -f table

# Config management
gohome config list
gohome config get days
gohome config set days 5

# Version info
gohome version

# Generate completions
gohome completion bash > /etc/bash_completion.d/gohome
```

### Example 2: Environment Variables

```bash
# Set via environment
export GOHOME_DAYS=3
export GOHOME_FORMAT=table
export GOHOME_PATH=~/workspace

# Run (uses env vars)
gohome

# Override with flags (flags win)
gohome -d 5
```

### Example 3: Multi-format Config

**JSON (~/.gohome.json):**
```json
{
  "days": 1,
  "format": "table",
  "path": "/home/user/workspace"
}
```

**YAML (~/.gohome.yaml):**
```yaml
days: 1
format: table
path: /home/user/workspace
```

**TOML (~/.gohome.toml):**
```toml
days = 1
format = "table"
path = "/home/user/workspace"
```

### Example 4: Plugin Development

```go
// plugins/csv-exporter/main.go
package main

import (
    "encoding/csv"
    "fmt"
    "os"
    
    "github.com/anIcedAntFA/gohome/internal/entity"
    "github.com/anIcedAntFA/gohome/internal/plugin"
)

type CSVExporter struct {
    config map[string]interface{}
}

func (e *CSVExporter) Name() string { return "csv-exporter" }
func (e *CSVExporter) Version() string { return "1.0.0" }

func (e *CSVExporter) Init(config map[string]interface{}) error {
    e.config = config
    return nil
}

func (e *CSVExporter) Execute(commits []entity.Commit) error {
    return nil // Not used for exporters
}

func (e *CSVExporter) Export(commits []entity.Commit, output string) error {
    file, err := os.Create(output)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    // Write header
    writer.Write([]string{"Type", "Scope", "Message", "Icon"})
    
    // Write commits
    for _, commit := range commits {
        writer.Write([]string{
            commit.Type,
            commit.Scope,
            commit.Message,
            commit.Icon,
        })
    }
    
    return nil
}

// Export plugin symbol
var Plugin plugin.ExporterPlugin = &CSVExporter{}
```

**Build plugin:**
```bash
go build -buildmode=plugin -o csv-exporter.so plugins/csv-exporter/main.go
```

**Use plugin:**
```bash
gohome report --plugin csv-exporter --output report.csv
```

---

## ⚠️ Breaking Changes

### Changes that May Affect Users

1. **Command Structure**
   - **Old:** `gohome --version`
   - **New:** `gohome version` (recommended, but `--version` still works)

2. **Config File Detection**
   - **Old:** Only `~/.gohome.json`
   - **New:** Tries `.gohome.json`, `.gohome.yaml`, `.gohome.toml` (first found wins)

3. **Error Messages**
   - More descriptive error messages with actionable suggestions
   - Different exit codes for different error types

4. **Log Output**
   - Structured logging with levels (ERROR, WARN, INFO, DEBUG)
   - Use `--quiet` to suppress logs

### Changes that Are Backward Compatible

1. **All v1.x flags work in v2.0**
2. **All v1.x config files work in v2.0**
3. **Default behavior is identical**

---

## 🔄 Rollback Strategy

### If Migration Fails

**Option 1: Fix Forward**
- Keep v2-cobra-migration branch
- Fix issues incrementally
- Merge when stable

**Option 2: Rollback**
```bash
# Abandon migration branch
git checkout main
git branch -D v2-cobra-migration

# Continue with v1.x architecture
# Plan retry with lessons learned
```

### If Release Has Critical Bug

```bash
# 1. Immediately publish v1.2.x patch
git checkout v1.2.0
git cherry-pick <critical-fix>
git tag v1.2.1
git push origin v1.2.1

# 2. Announce rollback recommendation
# 3. Fix v2.0 issue
# 4. Release v2.0.1 when ready
```

---

## 📚 Additional Resources

### Documentation to Create

- [ ] `docs/COBRA_VIPER_MIGRATION.md` (this file)
- [ ] `docs/V2_MIGRATION_GUIDE.md` (user-facing)
- [ ] `docs/PLUGIN_DEVELOPMENT.md` (plugin API docs)
- [ ] `docs/ARCHITECTURE_V2.md` (architecture diagrams)

### External References

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [Cobra Generator](https://github.com/spf13/cobra-cli)
- [Plugin Tutorial](https://golang.org/pkg/plugin/)

### Similar Projects

Study these projects that use Cobra/Viper:

- [kubectl](https://github.com/kubernetes/kubectl)
- [hugo](https://github.com/gohugoio/hugo)
- [gh (GitHub CLI)](https://github.com/cli/cli)
- [lazygit](https://github.com/jesseduffield/lazygit)

---

## 🎯 Success Criteria

### Technical Criteria

- ✅ All unit tests pass (>80% coverage)
- ✅ All integration tests pass
- ✅ Backward compatible with v1.x configs
- ✅ Performance equal or better than v1.x
- ✅ Zero data loss during migration
- ✅ Clean architecture with separation of concerns

### User Experience Criteria

- ✅ Migration is seamless (auto-detect old config)
- ✅ Help text is clear and comprehensive
- ✅ Error messages are actionable
- ✅ Documentation is complete
- ✅ Community feedback is positive

### Process Criteria

- ✅ Beta testing period (2 weeks minimum)
- ✅ Migration guide reviewed by 3+ users
- ✅ All breaking changes documented
- ✅ Rollback plan tested

---

## 📅 Timeline Summary

| Phase   | Duration | Deliverables                                    |
|---------|----------|-------------------------------------------------|
| Phase 1 | 2 weeks  | Cobra/Viper skeleton, parallel development      |
| Phase 2 | 4 weeks  | Core commands migrated, config subcommands      |
| Phase 3 | 3 weeks  | Env vars, completions, plugins, multi-format    |
| Phase 4 | 3 weeks  | Testing, docs, beta release, v2.0.0 final       |
| **Total** | **12 weeks** | **v2.0.0 with Cobra/Viper architecture** |

---

_Last Updated: January 15, 2026_  
_Status: Planning Phase - Ready for Implementation_  
_Maintainer: @anIcedAntFA_
