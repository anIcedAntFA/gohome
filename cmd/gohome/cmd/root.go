package cmd

import (
	"fmt"
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
	Version:       "v1.3.0-beta.1",
	SilenceErrors: true,  // We handle error formatting ourselves
	SilenceUsage:  false, // Show usage on errors (but only once)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Format error with consistent emoji prefix
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
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
}
