package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gohome",
	Short: "Git standup & activity reporting CLI",
	Long: `gohome scans your workspace for git repositories and generates
formatted daily standup reports from commit history.

Perfect for Daily Standups, Weekly Summaries,
or tracking your Personal Coding Habits.`,
	Version: "2.0.0-dev",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
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
	viper.AutomaticEnv()

	// Read config file (ignore error if not found)
	// Silently ignore if config file is not found
	viper.ReadInConfig()
}
