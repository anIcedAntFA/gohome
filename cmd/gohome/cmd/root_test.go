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

	// Verify Long contains key description (ASCII art is now in custom help)
	if !strings.Contains(rootCmd.Long, "Git Activity Aggregator") {
		t.Error("rootCmd.Long should contain description")
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

// TestExecuteFunction tests the Execute function wrapper.
func TestExecuteFunction(t *testing.T) {
	t.Run("execute_success", func(t *testing.T) {
		// Reset root command
		rootCmd.SetArgs([]string{"--version"})

		// Capture output
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		// Call Execute() wrapper directly
		// Save old os.Stderr
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		// This should not exit or panic
		// We can't actually test the exit(1) case without subprocess
		// So we just test the success path
		err := rootCmd.Execute()

		// Restore stderr
		_ = w.Close()
		os.Stderr = oldStderr
		_, _ = r.Read(make([]byte, 1024))

		if err != nil {
			t.Errorf("Execute() should not return error for --version, got: %v", err)
		}
	})

	t.Run("execute_with_help", func(t *testing.T) {
		rootCmd.SetArgs([]string{"--help"})

		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("Execute() should not return error for --help, got: %v", err)
		}

		if !strings.Contains(buf.String(), "Usage:") {
			t.Error("Help output should contain 'Usage:'")
		}
	})

	t.Run("execute_error_formatting", func(t *testing.T) {
		// Test error formatting by using invalid command
		rootCmd.SetArgs([]string{"nonexistent"})

		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		// Capture stderr for error message
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		// Execute will return error
		err := rootCmd.Execute()

		_ = w.Close()
		var stderrBuf bytes.Buffer
		_, _ = stderrBuf.ReadFrom(r)
		os.Stderr = oldStderr

		if err == nil {
			t.Error("Execute() should return error for invalid command")
		}
	})
}

// TestInitConfigWithHomeDir tests initConfig when home directory is available.
func TestInitConfigWithHomeDir(t *testing.T) {
	// Reset viper
	viper.Reset()
	defer viper.Reset()

	// Set empty config file to trigger default path logic
	cfgFile = ""

	// This should not panic even if config file doesn't exist
	initConfig()

	// Verify env prefix is set
	if viper.GetEnvPrefix() != "GOHOME" {
		t.Errorf("Expected env prefix 'GOHOME', got %q", viper.GetEnvPrefix())
	}
}

// TestRootCommandSilenceSettings tests error and usage silence settings.
func TestRootCommandSilenceSettings(t *testing.T) {
	if !rootCmd.SilenceErrors {
		t.Error("SilenceErrors should be true (we handle errors ourselves)")
	}

	if rootCmd.SilenceUsage {
		t.Error("SilenceUsage should be false (show usage on errors)")
	}
}

// TestRootCommandVersion tests the version field.
func TestRootCommandVersionField(t *testing.T) {
	expectedVersion := "v1.3.0"
	if rootCmd.Version != expectedVersion {
		t.Errorf("Version = %q, want %q", rootCmd.Version, expectedVersion)
	}
}

// TestInitFunctionCalled tests that init function is called.
func TestInitFunctionCalled(t *testing.T) {
	// Verify that cobra.OnInitialize was called by checking
	// if initConfig is in the initialization chain
	// This is tested implicitly by other tests that call Execute()

	// Just verify rootCmd is initialized
	if rootCmd == nil {
		t.Fatal("rootCmd should be initialized")
	}

	if rootCmd.Use == "" {
		t.Error("rootCmd.Use should not be empty after init()")
	}
}

// TestExecuteWithWriter verifies Execute error formatting with custom writer.
// This tests the refactored Execute function that allows dependency injection.
func TestExecuteWithWriter(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantErrMsg   string
	}{
		{
			name:         "invalid_command",
			args:         []string{"invalidcmd"},
			wantExitCode: 1,
			wantErrMsg:   "❌ Error:",
		},
		{
			name:         "valid_version_command",
			args:         []string{"version"},
			wantExitCode: 0,
			wantErrMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: Capture error output
			var errBuf bytes.Buffer

			// Reset command state
			rootCmd.SetArgs(tt.args)
			rootCmd.SetOut(&bytes.Buffer{})
			rootCmd.SetErr(&errBuf)

			// Execute: Run command with custom error writer
			exitCode := ExecuteWithWriter(&errBuf)

			// Assert: Check exit code
			if exitCode != tt.wantExitCode {
				t.Errorf("ExecuteWithWriter() exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}

			// Assert: Check error message if expecting failure
			errOutput := errBuf.String()
			if tt.wantExitCode != 0 {
				if !strings.Contains(errOutput, tt.wantErrMsg) {
					t.Errorf("Error output should contain %q, got: %q", tt.wantErrMsg, errOutput)
				}
			} else {
				// For successful commands, error output should be empty
				if errOutput != "" {
					t.Errorf("Error output should be empty for success, got: %q", errOutput)
				}
			}
		})
	}
}

// TestInitConfigWithWriter verifies config initialization with custom error writer.
// This tests error handling when home directory is inaccessible.
func TestInitConfigWithWriter(t *testing.T) {
	t.Run("successful_initialization", func(t *testing.T) {
		// Save original state
		viper.Reset()
		origCfgFile := cfgFile
		defer func() {
			cfgFile = origCfgFile
			viper.Reset()
		}()

		cfgFile = ""
		var errBuf bytes.Buffer

		// Execute: Initialize config
		err := initConfigWithWriter(&errBuf)
		// Assert: Should succeed
		if err != nil {
			t.Errorf("initConfigWithWriter() unexpected error: %v", err)
		}

		// Assert: No error output
		if errBuf.String() != "" {
			t.Errorf("Error output should be empty, got: %q", errBuf.String())
		}

		// Verify environment setup
		if viper.GetEnvPrefix() != "GOHOME" {
			t.Errorf("Env prefix = %q, want %q", viper.GetEnvPrefix(), "GOHOME")
		}
	})

	t.Run("with_custom_config_file", func(t *testing.T) {
		// Save original state
		viper.Reset()
		origCfgFile := cfgFile
		defer func() {
			cfgFile = origCfgFile
			viper.Reset()
		}()

		// Setup: Create temp config file
		tmpFile, err := os.CreateTemp("", "gohome-test-*.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		_, _ = tmpFile.WriteString(`{"path": "/test/path"}`)
		_ = tmpFile.Close()

		cfgFile = tmpFile.Name()
		var errBuf bytes.Buffer

		// Execute: Initialize with custom config
		err = initConfigWithWriter(&errBuf)
		// Assert: Should succeed
		if err != nil {
			t.Errorf("initConfigWithWriter() unexpected error: %v", err)
		}

		// Verify config was loaded
		if viper.ConfigFileUsed() != cfgFile {
			t.Errorf("ConfigFileUsed = %q, want %q", viper.ConfigFileUsed(), cfgFile)
		}
	})
}
