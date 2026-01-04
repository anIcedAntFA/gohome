package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
)

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
	ShowIcon  bool   `json:"show_icon"`
	ShowScope bool   `json:"show_scope"`

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

// loadConfigFromFile reads and parses the config file if it exists
func loadConfigFromFile() AppConfig {
	filePath := getConfigFilePath()
	var cfg AppConfig

	// Open file
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

	// Create or overwrite file
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

	flag.BoolVar(&cfg.ShowIcon, "icon", false, "")
	flag.BoolVar(&cfg.ShowIcon, "i", false, "")

	flag.BoolVar(&cfg.ShowScope, "scope", false, "")
	flag.BoolVar(&cfg.ShowScope, "c", false, "")

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
	_, _ = fmt.Fprintln(w, "  -H, --hours <int>\tNumber of hours to look back")
	_, _ = fmt.Fprintln(w, "  -d, --days <int>\tNumber of days to look back")
	_, _ = fmt.Fprintln(w, "  -w, --weeks <int>\tNumber of weeks to look back")
	_, _ = fmt.Fprintln(w, "  -m, --months <int>\tNumber of months to look back")
	_, _ = fmt.Fprintln(w, "  -y, --years <int>\tNumber of years to look back")
	_, _ = fmt.Fprintln(w, "      --today\tLook back since midnight today")
	_, _ = fmt.Fprintln(w, "\t")
	_, _ = fmt.Fprintln(w, "  -p, --path <string>\tRepo path to scan (default \".\")")
	_, _ = fmt.Fprintln(w, "  -a, --author <string>\tGit author (auto-detect if empty)")
	_, _ = fmt.Fprintln(w, "\t")
	_, _ = fmt.Fprintln(w, "  -f, --format <string>\tOutput format: text, table (default \"text\")")
	_, _ = fmt.Fprintln(w, "  -c, --scope\tShow commit scope")
	_, _ = fmt.Fprintln(w, "  -i, --icon\tShow commit type icons")
	_, _ = fmt.Fprintln(w, "\t")
	_, _ = fmt.Fprintln(w, "      --save\tSave current arguments as default configuration")

	_ = w.Flush() // Flush buffer to screen
	fmt.Fprintf(os.Stderr, "\n")
}

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
