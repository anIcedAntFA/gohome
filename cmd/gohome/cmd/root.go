package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/anIcedAntFA/gohome/internal/ui"
)

var cfgFile string

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:     "gohome",
	Short:   "Git standup & activity reporting CLI",
	Version: "v1.3.0",
	Long: `🏠 Git Activity Aggregator & Standup Report Generator

📊 gohome scans your workspace for git repositories and generates
   beautifully formatted daily standup reports from commit history.

✨ Perfect for:
   • Daily Standups & Team Sync
   • Weekly Progress Summaries
   • Personal Coding Habit Tracking
   • Developer Productivity Insights

💡 Quick Start:
   gohome --today          # Today's commits
   gohome -w 1             # Last week's commits
   gohome -f table         # Table format
   gohome config list      # View configuration`,
	SilenceErrors: true,  // We handle error formatting ourselves
	SilenceUsage:  false, // Show usage on errors (but only once)
}

// Execute runs the root command.
// This is the entry point called by main.go.
func Execute() {
	exitCode := ExecuteWithWriter(os.Stderr)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

// ExecuteWithWriter runs the root command with a custom error writer.
// This allows testing error output without calling os.Exit().
// Returns the exit code: 0 for success, 1 for error.
func ExecuteWithWriter(errWriter io.Writer) int {
	if err := rootCmd.Execute(); err != nil {
		// Format error with consistent emoji prefix
		fmt.Fprintf(errWriter, "❌ Error: %v\n", err)
		return 1
	}
	return 0
}

func init() {
	// Define theme flag at root level so it's available for help
	rootCmd.PersistentFlags().String("theme", "default", "Color theme: default, dracula, catppuccin-latte, catppuccin-mocha")
	_ = viper.BindPFlag("theme", rootCmd.PersistentFlags().Lookup("theme"))

	// Register theme flag completion
	_ = rootCmd.RegisterFlagCompletionFunc("theme", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"default", "dracula", "catppuccin-latte", "catppuccin-mocha"}, cobra.ShellCompDirectiveNoFileComp
	})

	cobra.OnInitialize(initConfig)

	// Apply themed help directly in init (before command execution)
	setupThemedHelp()
}

func setupThemedHelp() {
	// Create custom help func that checks for theme flag dynamically
	customHelpFunc := func(cmd *cobra.Command, _ []string) {
		// Try to get theme from flag first (for --theme dracula --help)
		theme := "default"
		if themeFlag := cmd.Flag("theme"); themeFlag != nil {
			flagValue := themeFlag.Value.String()
			if flagValue != "" && flagValue != "default" {
				theme = flagValue
			}
		}

		// For -h usage, check os.Args directly for --theme value
		// Note: -h --theme dracula won't work due to Cobra parsing order
		// Users should use: --theme dracula --help or --help --theme dracula
		if theme == "default" {
			for i, arg := range os.Args {
				if arg == "--theme" && i+1 < len(os.Args) {
					theme = os.Args[i+1]
					break
				}
			}
		}

		// Fallback to viper config if still default
		if theme == "default" {
			viperTheme := viper.GetString("theme")
			if viperTheme != "" {
				theme = viperTheme
			}
		}

		help := ui.FormatHelp(cmd, theme)
		fmt.Fprint(cmd.OutOrStdout(), help)
	}

	// Apply to root and all subcommands
	rootCmd.SetHelpFunc(customHelpFunc)
	for _, subCmd := range rootCmd.Commands() {
		subCmd.SetHelpFunc(customHelpFunc)
	}
}

func initConfig() {
	_ = initConfigWithWriter(os.Stderr)
}

// initConfigWithWriter initializes Viper configuration with a custom error writer.
func initConfigWithWriter(errWriter io.Writer) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprint(errWriter, err)
			return err
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".gohome")
	}

	// Environment variable support
	viper.SetEnvPrefix("GOHOME")
	// Replace underscores with hyphens in env var keys to match flag names
	// Example: GOHOME_MAX_DEPTH → max-depth → max_depth (via RegisterAlias)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Read config file (ignore error if not found)
	// Silently ignore if config file is not found
	_ = viper.ReadInConfig()
	return nil
}
