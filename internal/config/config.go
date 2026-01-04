// Package config handles application configuration loading, parsing, and persistence.
// It supports both command-line flags and JSON file storage.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/anIcedAntFA/gohome/internal/entity"
)

// StringSlice is a helper type for capturing multiple -t flag values.
// It implements flag.Value interface to support repeated flags like: -t "task1" -t "task2"
type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set appends a new value to the slice when the flag is encountered.
func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// AppConfig maps directly to the JSON file
type AppConfig struct {
	Hours  int  `json:"hours"`
	Days   int  `json:"days"`
	Weeks  int  `json:"weeks"`
	Months int  `json:"months"`
	Years  int  `json:"years"`
	Today  bool `json:"today"`

	Path      string `json:"path"`
	Author    string `json:"author"`
	OutputFmt string `json:"format"`
	Preset    string `json:"preset"`

	ShowIcon        bool `json:"show_icon"`
	ShowScope       bool `json:"show_scope"`
	CopyToClipboard bool `json:"copy_to_clipboard"`

	// Static Tasks loaded from JSON file (Rich objects)
	Tasks []entity.Task `json:"tasks"`

	// Dynamic Tasks from CLI flags (Simple strings) - This field is not loaded from JSON
	DynamicTasks StringSlice `json:"-"`

	// Special flag to save config, not saved to file
	SaveConfig bool `json:"-"`
}

// getConfigFilePath returns the config file path in user's home directory
func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.gohome.json" // Fallback to current directory if error
	}
	return filepath.Join(home, ".gohome.json")
}

// validateConfigPath ensures the config file path is safe from path traversal attacks
func validateConfigPath(filePath string) error {
	// Clean the path to remove any '..' or other unsafe elements
	cleanPath := filepath.Clean(filePath)

	// Get absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return err
	}

	// Ensure the file is named .gohome.json
	if filepath.Base(absPath) != ".gohome.json" {
		return errors.New("invalid config file name")
	}

	return nil
}

// loadConfigFromFile reads and parses the config file if it exists
func loadConfigFromFile() AppConfig {
	filePath := getConfigFilePath()
	var cfg AppConfig

	// Validate path to prevent path traversal attacks
	if err := validateConfigPath(filePath); err != nil {
		fmt.Printf("⚠️ Warning: Invalid config path: %v\\n", err)
		return cfg
	}

	// Open file
	// #nosec G304 -- filePath is validated above
	file, err := os.Open(filePath)
	if err != nil {
		return cfg // Return default config if file doesn't exist
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("⚠️ Warning: Failed to close file: %v\n", err)
		}
	}()

	// Parse JSON into struct
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		// If file format is invalid, skip and use default
		fmt.Printf("⚠️ Warning: Cannot parse config file at %s: %v\n", filePath, err)
	}

	return cfg
}

// SaveToFile writes the current config to the config file
func (c *AppConfig) SaveToFile() error {
	filePath := getConfigFilePath()

	// Validate path to prevent path traversal attacks
	if err := validateConfigPath(filePath); err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	// If user doesn't have any tasks yet (first time saving), add sample tasks as reference
	if len(c.Tasks) == 0 {
		c.Tasks = []entity.Task{
			// Group 1: Communication
			{Type: "meeting", Message: "Daily Standup & Team Sync", Icon: "📅", Enabled: true},
			{Type: "collab", Message: "Pair Programming / Mentoring", Icon: "👥", Enabled: true},

			// Group 2: Quality Assurance
			{Type: "review", Message: "Code Review & PR Feedback", Icon: "👀", Enabled: true},
			{Type: "testing", Message: "Write Unit/Integration Tests", Icon: "🧪", Enabled: true},

			// Group 3: Operations
			{Type: "ops", Message: "Monitor CI/CD Pipelines & Deploy", Icon: "🚀", Enabled: true},
			{Type: "admin", Message: "Check Emails, Jira & Sentry Logs", Icon: "📮", Enabled: true},
			// Group 4: Maintenance & Knowledge
			{Type: "docs", Message: "Update Documentation / Wiki", Icon: "📝", Enabled: true},
			{Type: "learning", Message: "Tech Research & Knowledge Sharing", Icon: "📚", Enabled: true},
		}
	}

	// Create or overwrite file
	// #nosec G304 -- filePath is validated above
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("⚠️ Warning: Failed to close file: %v\n", err)
		}
	}()

	// Write JSON from struct to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Format JSON nicely (Pretty print)
	return encoder.Encode(c)
}

// Load parses command-line flags and returns application configuration.
// It merges defaults from config file with CLI arguments.
func Load() *AppConfig {
	// A. Load defaults from file first (if exists)
	fileCfg := loadConfigFromFile()

	cfg := &AppConfig{}

	// B. Define flags
	flag.IntVar(&cfg.Hours, "hours", 0, "")
	flag.IntVar(&cfg.Hours, "H", 0, "")

	flag.IntVar(&cfg.Days, "days", 0, "")
	flag.IntVar(&cfg.Days, "d", 0, "")

	flag.IntVar(&cfg.Weeks, "weeks", 0, "")
	flag.IntVar(&cfg.Weeks, "w", 0, "")

	flag.IntVar(&cfg.Months, "months", 0, "")
	flag.IntVar(&cfg.Months, "m", 0, "")

	flag.IntVar(&cfg.Years, "years", 0, "")
	flag.IntVar(&cfg.Years, "y", 0, "")

	flag.BoolVar(&cfg.Today, "today", false, "")

	flag.StringVar(&cfg.Path, "path", ".", "")
	flag.StringVar(&cfg.Path, "p", ".", "")

	flag.StringVar(&cfg.Author, "author", "", "")
	flag.StringVar(&cfg.Author, "a", "", "")

	flag.StringVar(&cfg.OutputFmt, "format", "text", "")
	flag.StringVar(&cfg.OutputFmt, "f", "text", "")

	flag.StringVar(&cfg.Preset, "style", "normal", "")
	flag.StringVar(&cfg.Preset, "s", "", "normal")

	flag.BoolVar(&cfg.ShowIcon, "icon", false, "")
	flag.BoolVar(&cfg.ShowIcon, "i", false, "")

	flag.BoolVar(&cfg.ShowScope, "scope", false, "")
	flag.BoolVar(&cfg.ShowScope, "c", false, "")

	flag.BoolVar(&cfg.CopyToClipboard, "copy", false, "")
	flag.BoolVar(&cfg.CopyToClipboard, "cp", false, "")

	flag.Var(&cfg.DynamicTasks, "task", "")
	flag.Var(&cfg.DynamicTasks, "t", "")

	// Add flag for user to save config
	flag.BoolVar(&cfg.SaveConfig, "save", false, "Save current arguments as default configuration")

	flag.Usage = printUsage
	flag.Parse()

	// Track which flags were explicitly set by user
	userSetFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		userSetFlags[f.Name] = true
	})

	// --- A. Handle time period group (mutual exclusion) ---
	// Logic: If user sets ANY time flag -> ignore all time values from file.
	// If user sets NO time flags -> use all time values from file.

	isTimeSetByUser := checkTimeFlags(userSetFlags)

	if !isTimeSetByUser {
		// User didn't specify any time flags, load from file
		cfg.Hours = fileCfg.Hours
		cfg.Days = fileCfg.Days
		cfg.Weeks = fileCfg.Weeks
		cfg.Months = fileCfg.Months
		cfg.Years = fileCfg.Years
		cfg.Today = fileCfg.Today
	}
	// Otherwise: User has specified time flags (e.g., -d 5).
	// Since flag defaults are 0/false, cfg.Today is false and cfg.Weeks is 0.
	// Only cfg.Days is 5. Logic is correct!

	// --- B. Handle other independent flags ---
	// Logic: If user didn't set a flag, use value from file (if available)

	if !isSet(userSetFlags, "path", "p") && fileCfg.Path != "" {
		cfg.Path = fileCfg.Path
	}
	if !isSet(userSetFlags, "author", "a") && fileCfg.Author != "" {
		cfg.Author = fileCfg.Author
	}
	// For output format, need careful check since CLI default is "text"
	// If user didn't specify format flag, prioritize file value
	if !isSet(userSetFlags, "format", "f") && fileCfg.OutputFmt != "" {
		cfg.OutputFmt = fileCfg.OutputFmt
	}

	// Boolean flags (Icon/Scope)
	// Note: with booleans, if file is true and user wants to disable, user must override (but bool flags rarely support --flag=false easily).
	// Simplest approach for Phase 1: If user didn't touch the flag, use value from file.
	if !isSet(userSetFlags, "icon", "i") {
		cfg.ShowIcon = fileCfg.ShowIcon
	}
	if !isSet(userSetFlags, "scope", "c") {
		cfg.ShowScope = fileCfg.ShowScope
	}
	if !isSet(userSetFlags, "copy", "cp") {
		cfg.CopyToClipboard = fileCfg.CopyToClipboard
	}

	if len(fileCfg.Tasks) > 0 {
		cfg.Tasks = fileCfg.Tasks
	}

	return cfg
}

// checkTimeFlags checks if user has set any time-related flag
func checkTimeFlags(setFlags map[string]bool) bool {
	keys := []string{
		"hours", "H",
		"days", "d",
		"weeks", "w",
		"months", "m",
		"years", "y",
		"today",
	}
	for _, k := range keys {
		if setFlags[k] {
			return true
		}
	}
	return false
}

// isSet checks if user has set a flag (including its alias)
func isSet(setFlags map[string]bool, name, alias string) bool {
	return setFlags[name] || setFlags[alias]
}

// printUsage displays a custom professional help screen
func printUsage() {
	// Header
	fmt.Fprintf(os.Stderr, "\n🚀 GO HOME TOOL (Go CLI)\n")
	fmt.Fprintf(os.Stderr, "A simple tool to aggregate git commit reports.\n\n")
	fmt.Fprintf(os.Stderr, "Config file: %s\n\n", getConfigFilePath()) // Print config file location for user

	fmt.Fprintf(os.Stderr, "USAGE:\n")
	fmt.Fprintf(os.Stderr, "  gohome [flags]\n\n")

	fmt.Fprintf(os.Stderr, "EXAMPLES:\n")
	fmt.Fprintf(os.Stderr, "  gohome -d 3\n")
	fmt.Fprintf(os.Stderr, "  gohome -f table -s markdown -i -w 1\n\n")

	fmt.Fprintf(os.Stderr, "FLAGS:\n")

	// Use tabwriter to align columns
	// minwidth=0, tabwidth=8, padding=2, padchar=' ', flags=0
	w := tabwriter.NewWriter(os.Stderr, 0, 8, 2, ' ', 0)

	// Format: Flags \t Description
	fmt.Fprintln(w, "   -H, --hours <int>\tNumber of hours to look back")
	fmt.Fprintln(w, "   -d, --days <int>\tNumber of days to look back")
	fmt.Fprintln(w, "   -w, --weeks <int>\tNumber of weeks to look back")
	fmt.Fprintln(w, "   -m, --months <int>\tNumber of months to look back")
	fmt.Fprintln(w, "   -y, --years <int>\tNumber of years to look back")
	fmt.Fprintln(w, "      --today\tLook back since midnight today")
	fmt.Fprintln(w, "\t")
	fmt.Fprintln(w, "   -p, --path <string>\tRepo path to scan (default \".\")")
	fmt.Fprintln(w, "   -a, --author <string>\tGit author (auto-detect if empty)")
	fmt.Fprintln(w, "\t")
	fmt.Fprintln(w, "   -f, --format <string>\tOutput format: text, table (default \"text\")")
	fmt.Fprintln(w, "   -s, --style <string>\tPreset style: normal, markdown (default \"normal\")")
	fmt.Fprintln(w, "   -c, --scope\tShow commit scope")
	fmt.Fprintln(w, "   -i, --icon\tShow commit type icons")
	fmt.Fprintln(w, "\t")
	fmt.Fprintln(w, "  -cp, --copy\tCopy output to system clipboard")
	fmt.Fprintln(w, "       --save\tSave current arguments as default configuration")

	_ = w.Flush() // Flush buffer to screen
	fmt.Fprintf(os.Stderr, "\n")
}

// GetPeriod returns a human-readable time period string based on the configuration.
// It returns the largest non-zero period value (years > months > weeks > days > hours).
func (c *AppConfig) GetPeriod() string {
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
