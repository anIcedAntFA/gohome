package viper

import (
	"context"
	"fmt"

	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
	"github.com/spf13/viper"
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
	DynamicTasks []string      `mapstructure:"-"` // Not stored in config file
}

// LoadFromViper loads config with automatic hierarchy (flags > env > config > defaults)
func LoadFromViper() *Config {
	var cfg Config

	// Set defaults
	viper.SetDefault("days", 1)
	viper.SetDefault("path", ".")
	viper.SetDefault("max_depth", 2)
	viper.SetDefault("format", "text")
	viper.SetDefault("preset", "normal")
	viper.SetDefault("show_icon", false)
	viper.SetDefault("show_scope", false)

	// Unmarshal into struct (Viper handles precedence automatically)
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("⚠️ Warning: Failed to unmarshal config: %v\n", err)
		return &Config{
			Days:     1,
			Path:     ".",
			MaxDepth: 2,
			Format:   "text",
			Preset:   "normal",
		}
	}

	// Auto-detect author if not set
	if cfg.Author == "" {
		cfg.Author = detectGitAuthor()
	}

	return &cfg
}

// SaveToFile saves current config to file
func (c *Config) SaveToFile() error {
	// Set all fields
	viper.Set("hours", c.Hours)
	viper.Set("days", c.Days)
	viper.Set("weeks", c.Weeks)
	viper.Set("months", c.Months)
	viper.Set("years", c.Years)
	viper.Set("today", c.Today)

	viper.Set("path", c.Path)
	viper.Set("max_depth", c.MaxDepth)
	viper.Set("author", c.Author)

	viper.Set("format", c.Format)
	viper.Set("preset", c.Preset)
	viper.Set("show_icon", c.ShowIcon)
	viper.Set("show_scope", c.ShowScope)

	viper.Set("all_branches", c.AllBranches)
	viper.Set("branch", c.Branch)

	viper.Set("copy_to_clipboard", c.CopyToClipboard)

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

	viper.Set("tasks", c.Tasks)

	// Write to config file
	return viper.WriteConfig()
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
