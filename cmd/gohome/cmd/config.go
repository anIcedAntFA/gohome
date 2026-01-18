package cmd

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	viperconfig "github.com/anIcedAntFA/gohome/internal/config/viper"
)

var configCmd = &cobra.Command{
	Use:           "config",
	Short:         "Manage gohome configuration",
	Long:          `View, edit, and manage gohome configuration settings.`,
	SilenceErrors: true,  // We handle error formatting ourselves
	SilenceUsage:  false, // Show usage on errors
	// Ensure unknown subcommands show errors (by requiring a subcommand)
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: []string{"list", "get", "set", "reset"},
	RunE: func(_ *cobra.Command, args []string) error {
		// This will be called if Args validation passes but no subcommand matches
		return fmt.Errorf("unknown command %q for \"gohome config\"", args[0])
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show current configuration",
	Long:  `Display all current configuration settings with their values.`,
	RunE:  runConfigList,
}

var configGetCmd = &cobra.Command{
	Use:       "get <key>",
	Short:     "Get a configuration value",
	Long:      `Retrieve the value of a specific configuration key.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: viperconfig.GetAllConfigKeys(),
	RunE:      runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:       "set <key> <value>",
	Short:     "Set a configuration value",
	Long:      `Set a specific configuration key to the provided value and save to config file.`,
	Args:      cobra.ExactArgs(2),
	ValidArgs: viperconfig.GetAllConfigKeys(),
	RunE:      runConfigSet,
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long:  `Delete the configuration file and reset all settings to default values.`,
	RunE:  runConfigReset,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
}

func runConfigList(_ *cobra.Command, _ []string) error {
	// Load config using the clean Config struct
	cfg := viperconfig.LoadFromViper()

	fmt.Println("⚙️  Current Configuration:")
	fmt.Println("=====================")

	// Get config file location
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("📄 Config file: %s\n\n", configFile)
	} else {
		fmt.Println("No config file found (using defaults)")
		fmt.Println()
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{AutoFormat: tw.On},
				Alignment:  tw.CellAlignment{Global: tw.AlignCenter},
			},
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignLeft},
			},
		}),
	)

	table.Header([]string{"Key", "Value", "Description"})

	// Time period section
	_ = table.Append([]string{"hours", fmt.Sprintf("%d", cfg.Hours), "Number of hours to look back"})
	_ = table.Append([]string{"days", fmt.Sprintf("%d", cfg.Days), "Number of days to look back"})
	_ = table.Append([]string{"weeks", fmt.Sprintf("%d", cfg.Weeks), "Number of weeks to look back"})
	_ = table.Append([]string{"months", fmt.Sprintf("%d", cfg.Months), "Number of months to look back"})
	_ = table.Append([]string{"years", fmt.Sprintf("%d", cfg.Years), "Number of years to look back"})
	_ = table.Append([]string{"today", fmt.Sprintf("%t", cfg.Today), "Report from midnight to now"})
	_ = table.Append([]string{"", "", ""}) // Empty separator

	// Scanning section
	_ = table.Append([]string{"path", cfg.Path, "Root path to scan for repositories"})
	_ = table.Append([]string{"max_depth", fmt.Sprintf("%d", cfg.MaxDepth), "Maximum depth to scan repositories"})
	_ = table.Append([]string{"author", cfg.Author, "Git author name"})
	_ = table.Append([]string{"", "", ""}) // Empty separator

	// Output section
	_ = table.Append([]string{"format", cfg.Format, "Output format (text, table)"})
	_ = table.Append([]string{"style", cfg.Style, "Table style (normal, markdown)"})
	_ = table.Append([]string{"icon", fmt.Sprintf("%t", cfg.ShowIcon), "Show commit icons"})
	_ = table.Append([]string{"scope", fmt.Sprintf("%t", cfg.ShowScope), "Show commit scope"})
	_ = table.Append([]string{"", "", ""}) // Empty separator

	// Branch filtering
	_ = table.Append([]string{"all_branches", fmt.Sprintf("%t", cfg.AllBranches), "Include all branches"})
	_ = table.Append([]string{"branch", cfg.Branch, "Specific branch to filter"})
	_ = table.Append([]string{"", "", ""}) // Empty separator

	// Misc
	_ = table.Append([]string{"copy", fmt.Sprintf("%t", cfg.CopyToClipboard), "Copy output to clipboard"})
	_ = table.Append([]string{"tasks", fmt.Sprintf("%d task(s)", len(cfg.Tasks)), "Static recurring tasks"})

	_ = table.Render()
	return nil
}

func runConfigGet(_ *cobra.Command, args []string) error {
	key := args[0]
	value := viper.Get(key)

	if value == nil {
		return fmt.Errorf("key %q not found in configuration", key)
	}

	fmt.Printf("%s = %v\n", key, value)
	return nil
}

func runConfigSet(_ *cobra.Command, args []string) error {
	key, valueStr := args[0], args[1]

	// Load existing config or create new one
	cfg := viperconfig.LoadFromViper()

	// Use reflection-based SetValue (DRY, maintainable)
	if err := cfg.SetValue(key, valueStr); err != nil {
		return fmt.Errorf("failed to set %q: %w", key, err)
	}

	// Validate configuration before saving
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Save using the same method as --save flag (clean JSON, no duplicates)
	if err := cfg.SaveToFile(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("✅ Set %s = %s\n", key, valueStr)
	return nil
}

func runConfigReset(_ *cobra.Command, _ []string) error {
	// Get config file location
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		fmt.Println("No config file found. Nothing to reset.")
		return nil
	}

	// Prompt for confirmation
	fmt.Printf("Are you sure you want to delete %s? (y/N): ", configFile)
	var response string
	_, _ = fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		fmt.Println("Canceled.")
		return nil
	}

	// Delete config file
	if err := os.Remove(configFile); err != nil {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	fmt.Println("✅ Configuration reset to defaults.")
	return nil
}
