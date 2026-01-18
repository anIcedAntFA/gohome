package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

// TestRootCommandProperties tests the root command metadata.
func TestRootCommandProperties(t *testing.T) {
	if rootCmd.Use != "gohome" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "gohome")
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}

	// Verify Long contains ASCII art
	if !strings.Contains(rootCmd.Long, "____") {
		t.Error("rootCmd.Long should contain ASCII art")
	}

	// Verify Long contains emoji indicators
	emojis := []string{"🏠", "📊", "✨", "💡"}
	for _, emoji := range emojis {
		if !strings.Contains(rootCmd.Long, emoji) {
			t.Errorf("rootCmd.Long should contain emoji %q", emoji)
		}
	}

	if rootCmd.Version == "" {
		t.Error("rootCmd.Version should not be empty")
	}

	if !rootCmd.SilenceErrors {
		t.Error("rootCmd.SilenceErrors should be true")
	}

	if rootCmd.SilenceUsage {
		t.Error("rootCmd.SilenceUsage should be false")
	}
}

// TestRootCommandHasSubcommands tests that root has expected subcommands.
func TestRootCommandHasSubcommands(t *testing.T) {
	expectedCommands := []string{
		"config",
		"report",
		"version",
		"completion",
	}

	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("rootCmd should have subcommand %q", expected)
		}
	}
}

// TestExecuteWithInvalidCommand tests error handling.
func TestExecuteWithInvalidCommand(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set invalid command
	os.Args = []string{"gohome", "invalidcommand"}

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Reset root command for clean test
	rootCmd.SetArgs([]string{"invalidcommand"})

	// Capture output
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()

	// Restore stderr and read pipe
	_ = w.Close()
	var stderrBuf bytes.Buffer
	_, _ = stderrBuf.ReadFrom(r)
	os.Stderr = oldStderr

	if err == nil {
		t.Error("Execute() with invalid command should return error")
	}

	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("Error message should contain 'unknown command', got: %v", err)
	}
}

// TestInitConfig tests the configuration initialization.
func TestInitConfig(t *testing.T) {
	// Save original viper state
	origConfigFile := viper.ConfigFileUsed()
	defer func() {
		// Cleanup after test
		viper.Reset()
		if origConfigFile != "" {
			viper.SetConfigFile(origConfigFile)
			_ = viper.ReadInConfig()
		}
	}()

	// Test 1: Default config initialization (no config file specified)
	t.Run("default_initialization", func(t *testing.T) {
		viper.Reset()
		cfgFile = ""

		initConfig()

		// Verify environment variable setup
		if viper.GetEnvPrefix() != "GOHOME" {
			t.Errorf("Env prefix = %q, want %q", viper.GetEnvPrefix(), "GOHOME")
		}
	})

	// Test 2: Custom config file specified
	t.Run("custom_config_file", func(t *testing.T) {
		viper.Reset()

		// Create a temporary config file
		tmpFile, err := os.CreateTemp("", "gohome-test-*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write minimal config
		_, err = tmpFile.WriteString(`{"days": 5}`)
		if err != nil {
			t.Fatalf("Failed to write temp file: %v", err)
		}
		_ = tmpFile.Close()

		// Set config file explicitly
		viper.SetConfigFile(tmpFile.Name())
		err = viper.ReadInConfig()
		if err != nil {
			t.Fatalf("Failed to read config: %v", err)
		}

		// Verify config file is used
		if viper.ConfigFileUsed() != tmpFile.Name() {
			t.Errorf("ConfigFileUsed() = %q, want %q", viper.ConfigFileUsed(), tmpFile.Name())
		}

		// Verify config value is loaded
		if viper.GetInt("days") != 5 {
			t.Errorf("days = %d, want 5", viper.GetInt("days"))
		}
	})

	// Test 3: Environment variable support
	t.Run("environment_variables", func(t *testing.T) {
		viper.Reset()
		cfgFile = ""

		// Set test environment variable
		_ = os.Setenv("GOHOME_DAYS", "10")
		defer func() { _ = os.Unsetenv("GOHOME_DAYS") }()

		initConfig()

		// Verify env var is read
		if viper.GetInt("days") != 10 {
			t.Errorf("days from env = %d, want 10", viper.GetInt("days"))
		}
	})
}

// TestRootCommandHelp tests that help output is generated.
func TestRootCommandHelp(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	output := buf.String()

	// Verify help contains key sections
	requiredSections := []string{
		"Usage:",
		"Available Commands:",
		"Flags:",
		"gohome",
	}

	for _, section := range requiredSections {
		if !strings.Contains(output, section) {
			t.Errorf("Help output should contain %q", section)
		}
	}
}

// TestRootCommandVersion tests version flag.
func TestRootCommandVersion(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"--version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Version flag failed: %v", err)
	}

	output := buf.String()

	// Verify version output contains "gohome" (command name is always shown)
	if !strings.Contains(output, "gohome") {
		t.Error("Version output should contain 'gohome'")
	}

	// Version format can vary (sometimes just "version X", sometimes includes version keyword)
	// So we just check that output is not empty
	if output == "" {
		t.Error("Version output should not be empty")
	}
}

// TestEnvVarKeyReplacer tests that environment variable names are correctly mapped.
func TestEnvVarKeyReplacer(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		key    string
		value  string
	}{
		{
			name:   "underscore_to_hyphen",
			envVar: "GOHOME_MAX_DEPTH",
			key:    "max-depth",
			value:  "5",
		},
		{
			name:   "simple_key",
			envVar: "GOHOME_DAYS",
			key:    "days",
			value:  "7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper state
			viper.Reset()
			defer viper.Reset()

			// Set environment variable
			_ = os.Setenv(tt.envVar, tt.value)
			defer func() { _ = os.Unsetenv(tt.envVar) }()

			// Initialize config
			cfgFile = ""
			initConfig()

			// Check if value is accessible via the key
			got := viper.GetString(tt.key)
			if got != tt.value {
				t.Errorf("viper.GetString(%q) = %q, want %q", tt.key, got, tt.value)
			}
		})
	}
}
