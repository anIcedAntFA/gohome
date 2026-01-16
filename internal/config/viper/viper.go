package viper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Time period
	Hours  int  `mapstructure:"hours" json:"hours"`
	Days   int  `mapstructure:"days" json:"days"`
	Weeks  int  `mapstructure:"weeks" json:"weeks"`
	Months int  `mapstructure:"months" json:"months"`
	Years  int  `mapstructure:"years" json:"years"`
	Today  bool `mapstructure:"today" json:"today"`

	// Path and scanning
	Path     string `mapstructure:"path" json:"path"`
	MaxDepth int    `mapstructure:"max_depth" json:"max_depth"`
	Author   string `mapstructure:"author" json:"author"`

	// Output
	Format    string `mapstructure:"format" json:"format"`
	Style     string `mapstructure:"style" json:"style"`
	ShowIcon  bool   `mapstructure:"icon" json:"icon"`
	ShowScope bool   `mapstructure:"scope" json:"scope"`

	// Branch filtering
	AllBranches bool   `mapstructure:"all_branches" json:"all_branches"`
	Branch      string `mapstructure:"branch" json:"branch"`

	// Clipboard
	CopyToClipboard bool `mapstructure:"copy" json:"copy"`

	// Tasks
	Tasks        []entity.Task `mapstructure:"tasks" json:"tasks"`
	DynamicTasks []string      `mapstructure:"-" json:"-"` // Not stored in config file
}

// LoadFromViper loads config with automatic hierarchy (flags > env > config > defaults)
func LoadFromViper() *Config {
	var cfg Config

	// Set defaults
	viper.SetDefault("days", 1)
	viper.SetDefault("path", ".")
	viper.SetDefault("max_depth", 2)
	viper.SetDefault("format", "text")
	viper.SetDefault("style", "normal")
	viper.SetDefault("icon", false)
	viper.SetDefault("scope", false)
	viper.SetDefault("copy", false)

	// Unmarshal into struct (Viper handles precedence automatically)
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("⚠️ Warning: Failed to unmarshal config: %v\n", err)
		return &Config{
			Days:     1,
			Path:     ".",
			MaxDepth: 2,
			Format:   "text",
			Style:    "normal",
		}
	}

	// Auto-detect author if not set
	if cfg.Author == "" {
		cfg.Author = detectGitAuthor()
	}

	return &cfg
}

// SaveToFile saves current config to file (clean JSON without duplicates)
func (c *Config) SaveToFile() error {
	// Normalize period fields: only keep the highest priority non-zero period
	// This prevents confusion when multiple period fields are set
	c.NormalizePeriod()

	// If user doesn't have any tasks yet, add sample tasks
	if len(c.Tasks) == 0 {
		c.Tasks = []entity.Task{
			// Group 1: Communication
			{Type: "meeting", Message: "Daily Standup & Team Sync", Icon: "📅", Enabled: false},
			{Type: "collab", Message: "Pair Programming / Mentoring", Icon: "👥", Enabled: false},

			// Group 2: Quality Assurance
			{Type: "review", Message: "Code Review & PR Feedback", Icon: "👀", Enabled: true},
			{Type: "testing", Message: "Write Unit/Integration Tests", Icon: "🧪", Enabled: false},

			// Group 3: Operations
			{Type: "ops", Message: "Monitor CI/CD Pipelines & Deploy", Icon: "🚀", Enabled: false},
			{Type: "admin", Message: "Check Emails, Jira & Sentry Logs", Icon: "📮", Enabled: false},

			// Group 4: Maintenance & Knowledge
			{Type: "docs", Message: "Update Documentation / Wiki", Icon: "📝", Enabled: false},
			{Type: "learning", Message: "Tech Research & Knowledge Sharing", Icon: "📚", Enabled: true},
		}
	}

	// Determine config file path
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configPath := filepath.Join(home, ".gohome.json")

	// Write clean JSON directly to avoid Viper's duplicate keys
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

// detectGitAuthor auto-detects git author from git config
func detectGitAuthor() string {
	gitClient := git.NewClient()
	author := gitClient.GetUser(context.Background())
	return author
}

// GetPeriod returns a human-readable time period string based on the configuration.
func (c *Config) GetPeriod() string {
	// Table-driven approach: check periods in priority order
	periods := []struct {
		value int
		unit  string
	}{
		{c.Years, "year"},
		{c.Months, "month"},
		{c.Weeks, "week"},
		{c.Days, "day"},
		{c.Hours, "hour"},
	}

	// Special case: today flag
	if c.Today {
		return "midnight"
	}

	// Check each period in order
	for _, p := range periods {
		if p.value > 0 {
			return fmt.Sprintf("%d %s%s ago", p.value, p.unit, pluralize(p.value))
		}
	}

	// Default fallback
	return "24 hours ago"
}

// pluralize adds "s" suffix for plural values
func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// normalizePeriod ensures only one period field is non-zero when saving to config.
// Priority: today > years > months > weeks > days > hours
// This prevents confusion when multiple period fields are set in config file.
func (c *Config) NormalizePeriod() {
	// Special case: today flag takes highest priority
	if c.Today {
		c.Years = 0
		c.Months = 0
		c.Weeks = 0
		c.Days = 0
		c.Hours = 0
		return
	}

	// Check periods in priority order, keep only the first non-zero
	if c.Years > 0 {
		c.Months = 0
		c.Weeks = 0
		c.Days = 0
		c.Hours = 0
	} else if c.Months > 0 {
		c.Weeks = 0
		c.Days = 0
		c.Hours = 0
	} else if c.Weeks > 0 {
		c.Days = 0
		c.Hours = 0
	} else if c.Days > 0 {
		c.Hours = 0
	}
	// If only Hours is set, keep it as is
}
