package viper

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/anIcedAntFA/gohome/internal/entity"
	"github.com/anIcedAntFA/gohome/internal/git"
)

// init runs once at package initialization (Viper best practice).
// Sets up aliases, defaults, and config file paths.
func init() {
	// Register key aliases ONCE (kebab-case flags → snake_case config)
	// Allows Viper to map flag "max-depth" to config field "max_depth"
	viper.RegisterAlias("max_depth", "max-depth")
	viper.RegisterAlias("all_branches", "all-branches")

	// Set config file search paths
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")
	viper.SetConfigName(".gohome")
	viper.SetConfigType("json")

	// Set defaults ONCE at initialization
	viper.SetDefault("days", 1)
	viper.SetDefault("path", ".")
	viper.SetDefault("max_depth", 2)
	viper.SetDefault("format", "text")
	viper.SetDefault("style", "normal")
	viper.SetDefault("icon", false)
	viper.SetDefault("scope", false)
	viper.SetDefault("copy", false)
	viper.SetDefault("all_branches", false)
	viper.SetDefault("today", false)

	// Auto-read config file if exists (silent fail is OK)
	_ = viper.ReadInConfig()
}

// Config represents the application configuration.
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

// LoadFromViper loads config with automatic hierarchy (flags > env > config > defaults).
// Viper precedence: Set > flags > env > config > defaults (all handled automatically).
func LoadFromViper() *Config {
	var cfg Config

	// Unmarshal all Viper values into struct (precedence handled by Viper)
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("⚠️  Warning: Failed to unmarshal config: %v\n", err)
		// Fallback to defaults
		return &Config{
			Days:     1,
			Path:     ".",
			MaxDepth: 2,
			Format:   "text",
			Style:    "normal",
		}
	}

	// Post-processing: Auto-detect author if not set
	if cfg.Author == "" {
		cfg.Author = detectGitAuthor()
	}

	return &cfg
}

// SaveToFile saves current config to file (clean JSON without duplicates).
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
	// #nosec G304 -- configPath is constructed from os.UserHomeDir() which is safe
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}

// detectGitAuthor auto-detects git author from git config.
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

// pluralize adds "s" suffix for plural values.
func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// SetValue sets a config field by key with automatic type conversion.
// Uses reflection to avoid giant switch statements (DRY principle).
// Returns error if key is unknown or value type is incompatible.
func (c *Config) SetValue(key, value string) error {
	// Normalize key: kebab-case → snake_case for consistency
	key = strings.ReplaceAll(key, "-", "_")

	// Get field name from mapping
	fieldName, err := getFieldName(key)
	if err != nil {
		return err
	}

	// Use reflection to set field value
	v := reflect.ValueOf(c).Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return fmt.Errorf("cannot set field %q", key)
	}

	// Type conversion based on field type
	return setFieldValue(field, key, value)
}

// getFieldName maps config key to struct field name.
func getFieldName(key string) (string, error) {
	fieldMap := map[string]string{
		"hours":        "Hours",
		"days":         "Days",
		"weeks":        "Weeks",
		"months":       "Months",
		"years":        "Years",
		"today":        "Today",
		"path":         "Path",
		"max_depth":    "MaxDepth",
		"author":       "Author",
		"format":       "Format",
		"style":        "Style",
		"icon":         "ShowIcon",
		"scope":        "ShowScope",
		"all_branches": "AllBranches",
		"branch":       "Branch",
		"copy":         "CopyToClipboard",
	}

	fieldName, ok := fieldMap[key]
	if !ok {
		if key == "tasks" {
			return "", fmt.Errorf("tasks cannot be set via CLI; edit ~/.gohome.json directly")
		}
		return "", fmt.Errorf("unknown config key %q", key)
	}

	return fieldName, nil
}

// setFieldValue sets a reflect.Value based on its type.
func setFieldValue(field reflect.Value, key, value string) error {
	//nolint:exhaustive // Only Int, Bool, String are supported; other types explicitly return error
	switch field.Kind() {
	case reflect.Int:
		return setIntField(field, value)

	case reflect.Bool:
		return setBoolField(field, value)

	case reflect.String:
		return setStringField(field, key, value)

	default:
		// Explicitly handle all other reflect.Kind cases to satisfy exhaustive linter
		return fmt.Errorf("unsupported field type %v for key %q", field.Kind(), key)
	}
}

// setIntField parses and sets an integer field.
func setIntField(field reflect.Value, value string) error {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("expected integer value, got %q", value)
	}
	field.SetInt(int64(intVal))
	return nil
}

// setBoolField parses and sets a boolean field.
func setBoolField(field reflect.Value, value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("expected boolean value (true/false), got %q", value)
	}
	field.SetBool(boolVal)
	return nil
}

// setStringField validates and sets a string field.
func setStringField(field reflect.Value, key, value string) error {
	// Validate enum values for specific fields
	if err := validateEnumValue(key, value); err != nil {
		return err
	}
	field.SetString(value)
	return nil
}

// validateEnumValue validates enum fields (format, style).
func validateEnumValue(key, value string) error {
	switch key {
	case "format":
		if value != "text" && value != "table" {
			return fmt.Errorf("format must be 'text' or 'table', got %q", value)
		}
	case "style":
		if value != "normal" && value != "markdown" {
			return fmt.Errorf("style must be 'normal' or 'markdown', got %q", value)
		}
	}
	return nil
}

// NormalizePeriod ensures only one period field is non-zero when saving to config.
// Priority: today > years > months > weeks > days > hours.
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
	switch {
	case c.Years > 0:
		c.Months = 0
		c.Weeks = 0
		c.Days = 0
		c.Hours = 0
	case c.Months > 0:
		c.Weeks = 0
		c.Days = 0
		c.Hours = 0
	case c.Weeks > 0:
		c.Days = 0
		c.Hours = 0
	case c.Days > 0:
		c.Hours = 0
		// If only Hours is set, keep it as is
	}
}
