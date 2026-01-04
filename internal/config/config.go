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

// 1. Function to get config file path in user's Home directory
func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./gohome.json" // Fallback to current directory if error
	}
	return filepath.Join(home, "gohome.json")
}

// 2. Function to read config file (if exists)
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

// 3. Function to save current config to file
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
	// A. Load Default from File first (if exists)
	fileCfg := loadConfigFromFile()

	cfg := &AppConfig{}

	// B. Define Flags
	// Trick: Use values from fileCfg as default values for Flag.
	// If user doesn't specify flag, it will use value from file.
	// If file doesn't exist yet, value will be 0 or "" (zero value of struct).

	flag.IntVar(&cfg.Hours, "hours", fileCfg.Hours, "")
	flag.IntVar(&cfg.Hours, "H", fileCfg.Hours, "")

	flag.IntVar(&cfg.Days, "days", fileCfg.Days, "")
	flag.IntVar(&cfg.Days, "d", fileCfg.Days, "")

	flag.IntVar(&cfg.Weeks, "weeks", fileCfg.Weeks, "")
	flag.IntVar(&cfg.Weeks, "w", fileCfg.Weeks, "")

	flag.IntVar(&cfg.Months, "months", fileCfg.Months, "")
	flag.IntVar(&cfg.Months, "m", fileCfg.Months, "")

	flag.IntVar(&cfg.Years, "years", fileCfg.Years, "")
	flag.IntVar(&cfg.Years, "y", fileCfg.Years, "")

	flag.BoolVar(&cfg.Today, "today", fileCfg.Today, "")

	// Handle strings, note hard-coded default fallback if file is empty
	defaultPath := fileCfg.Path
	if defaultPath == "" {
		defaultPath = "."
	}
	flag.StringVar(&cfg.Path, "path", defaultPath, "")
	flag.StringVar(&cfg.Path, "p", defaultPath, "")

	flag.StringVar(&cfg.Author, "author", fileCfg.Author, "")
	flag.StringVar(&cfg.Author, "a", fileCfg.Author, "")

	defaultFmt := fileCfg.OutputFmt
	if defaultFmt == "" {
		defaultFmt = "text"
	}
	flag.StringVar(&cfg.OutputFmt, "format", defaultFmt, "")
	flag.StringVar(&cfg.OutputFmt, "f", defaultFmt, "")

	flag.BoolVar(&cfg.ShowIcon, "icon", fileCfg.ShowIcon, "")
	flag.BoolVar(&cfg.ShowIcon, "i", fileCfg.ShowIcon, "")

	flag.BoolVar(&cfg.ShowScope, "scope", fileCfg.ShowScope, "")
	flag.BoolVar(&cfg.ShowScope, "c", fileCfg.ShowScope, "")

	// Add new flag for user to command save config
	flag.BoolVar(&cfg.SaveConfig, "save", false, "Save current arguments as default configuration")

	flag.Usage = printUsage
	flag.Parse()

	return cfg
}

// printUsage: Custom professional help screen
func printUsage() {
	// Header
	fmt.Fprintf(os.Stderr, "\n🚀 GO HOME TOOL (Go CLI)\n")
	fmt.Fprintf(os.Stderr, "A simple tool to aggregate git commit reports.\n\n")
	fmt.Fprintf(os.Stderr, "Config file: %s\n\n", getConfigFilePath()) // Print config file location for user

	fmt.Fprintf(os.Stderr, "USAGE:\n")
	fmt.Fprintf(os.Stderr, "  gohome [flags]\n\n")

	fmt.Fprintf(os.Stderr, "EXAMPLES:\n")
	fmt.Fprintf(os.Stderr, "  gohome -d 3\n")
	fmt.Fprintf(os.Stderr, "  gohome -f table -s nature -i -d 7\n\n")

	fmt.Fprintf(os.Stderr, "FLAGS:\n")

	// Use Tabwriter to align columns
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
	return "1 day ago"
}

// pluralize adds "s" suffix for plural values
func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
