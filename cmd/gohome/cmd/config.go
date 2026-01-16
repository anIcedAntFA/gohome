package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gohome configuration",
	Long:  `View, edit, and manage gohome configuration settings.`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Show current configuration",
	Long:  `Display all current configuration settings with their values.`,
	RunE:  runConfigList,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long:  `Retrieve the value of a specific configuration key.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Set a specific configuration key to the provided value and save to config file.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runConfigSet,
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

func runConfigList(cmd *cobra.Command, args []string) error {
	fmt.Println("Current Configuration:")
	fmt.Println("=====================")

	// Get config file location
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("Config file: %s\n\n", configFile)
	} else {
		fmt.Println("No config file found (using defaults)")
		fmt.Println()
	}

	// Display all settings
	allSettings := viper.AllSettings()
	if len(allSettings) == 0 {
		fmt.Println("No configuration settings found.")
		return nil
	}

	for key, value := range allSettings {
		fmt.Printf("%-20s = %v\n", key, value)
	}

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := viper.Get(key)

	if value == nil {
		return fmt.Errorf("key %q not found in configuration", key)
	}

	fmt.Printf("%s = %v\n", key, value)
	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key, value := args[0], args[1]

	viper.Set(key, value)

	// Try to write to existing config file, or create new one
	if err := viper.WriteConfig(); err != nil {
		// If config file doesn't exist, create it
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			viper.SetConfigFile(home + "/.gohome.json")
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	fmt.Printf("✅ Set %s = %s\n", key, value)
	return nil
}

func runConfigReset(cmd *cobra.Command, args []string) error {
	// Get config file location
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		fmt.Println("No config file found. Nothing to reset.")
		return nil
	}

	// Prompt for confirmation
	fmt.Printf("Are you sure you want to delete %s? (y/N): ", configFile)
	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" {
		fmt.Println("Cancelled.")
		return nil
	}

	// Delete config file
	if err := os.Remove(configFile); err != nil {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	fmt.Println("✅ Configuration reset to defaults.")
	return nil
}
