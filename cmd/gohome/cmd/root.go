package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "gohome",
	Short: "Git standup & activity reporting CLI",
	Long: `
   ____       _   _
  / ___| ___ | | | | ___  _ __ ___   ___
 | |  _ / _ \| |_| |/ _ \| '_ ' _ \ / _ \
 | |_| | (_) |  _  | (_) | | | | | |  __/
  \____|\___/|_| |_|\___/|_| |_| |_|\___|

🏠 Git Activity Aggregator & Standup Report Generator

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
	Version:       "v1.3.0",
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
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	_ = initConfigWithWriter(os.Stderr)
}

// initConfigWithWriter initializes Viper configuration with a custom error writer.
// This allows testing error handling without calling os.Exit().
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
