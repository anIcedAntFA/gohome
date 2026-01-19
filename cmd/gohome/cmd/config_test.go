package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TestConfigCmdInit verifies the config command and all subcommands are properly initialized.
func TestConfigCmdInit(t *testing.T) {
	if configCmd == nil {
		t.Fatal("configCmd should not be nil")
	}

	if configCmd.Use != "config" {
		t.Errorf("expected Use to be 'config', got %q", configCmd.Use)
	}

	if !configCmd.SilenceErrors {
		t.Error("expected SilenceErrors to be true")
	}

	if configCmd.SilenceUsage {
		t.Error("expected SilenceUsage to be false")
	}

	// Verify minimum args requirement
	if err := configCmd.Args(configCmd, []string{}); err == nil {
		t.Error("expected error when no args provided")
	}

	// Verify valid args
	expectedValidArgs := []string{"list", "get", "set", "reset"}
	if len(configCmd.ValidArgs) != len(expectedValidArgs) {
		t.Errorf("expected %d valid args, got %d", len(expectedValidArgs), len(configCmd.ValidArgs))
	}

	// Verify all subcommands exist
	subcommands := []struct {
		name string
		cmd  *cobra.Command
		use  string
	}{
		{"list", configListCmd, "list"},
		{"get", configGetCmd, "get <key>"},
		{"set", configSetCmd, "set <key> <value>"},
		{"reset", configResetCmd, "reset"},
	}

	for _, sc := range subcommands {
		if sc.cmd == nil {
			t.Errorf("%s subcommand should not be nil", sc.name)
		}
		if sc.cmd.Use != sc.use {
			t.Errorf("expected Use to be %q, got %q", sc.use, sc.cmd.Use)
		}
	}
}

// TestConfigCmdUnknownSubcommand verifies error handling for invalid subcommands.
func TestConfigCmdUnknownSubcommand(t *testing.T) {
	// Reset viper for clean state
	viper.Reset()

	err := configCmd.RunE(configCmd, []string{"invalid"})
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}

	expectedMsg := `unknown command "invalid" for "gohome config"`
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("expected error message to contain %q, got %q", expectedMsg, err.Error())
	}
}

// TestRunConfigList verifies the config list command output.
func TestRunConfigList(t *testing.T) {
	// Reset viper for clean state
	viper.Reset()

	// Create temp config file for testing
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".gohome.json")

	// Set some test values
	viper.Set("hours", 24)
	viper.Set("path", "/test/path")
	viper.Set("author", "Test Author")
	viper.Set("format", "table")
	viper.SetConfigFile(configPath)

	// Write config to file so ConfigFileUsed() returns the path
	if err := viper.WriteConfig(); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute command
	err := runConfigList(configListCmd, []string{})

	// Restore stdout
	_ = w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("runConfigList failed: %v", err)
	}

	// Read captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify key outputs
	expectedStrings := []string{
		"Current Configuration",
		"Config file:",
		"hours",
		"path",
		"author",
		"format",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, output)
		}
	}

	// Verify config file path in output
	if !strings.Contains(output, configPath) {
		t.Errorf("expected output to contain config file path %q", configPath)
	}
}

// TestRunConfigListNoConfigFile verifies list command with no config file.
func TestRunConfigListNoConfigFile(t *testing.T) {
	// Reset viper for clean state
	viper.Reset()

	// Ensure no config file is used
	viper.SetConfigFile("")

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runConfigList(configListCmd, []string{})

	_ = w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("runConfigList failed: %v", err)
	}

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Should show defaults message
	expectedMsg := "No config file found (using defaults)"
	if !strings.Contains(output, expectedMsg) {
		t.Errorf("expected output to contain %q, got:\n%s", expectedMsg, output)
	}
}

// TestRunConfigGet verifies getting a config value.
func TestRunConfigGet(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         interface{}
		expectedOut   string
		expectedError bool
	}{
		{
			name:        "get existing string value",
			key:         "path",
			value:       "/test/path",
			expectedOut: "path = /test/path",
		},
		{
			name:        "get existing int value",
			key:         "hours",
			value:       24,
			expectedOut: "hours = 24",
		},
		{
			name:        "get existing bool value",
			key:         "today",
			value:       true,
			expectedOut: "today = true",
		},
		{
			name:          "get non-existent key",
			key:           "nonexistent",
			value:         nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for clean state
			viper.Reset()

			// Set test value if provided
			if tt.value != nil {
				viper.Set(tt.key, tt.value)
			}

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := runConfigGet(configGetCmd, []string{tt.key})

			_ = w.Close()
			os.Stdout = old

			// Check error expectation
			if tt.expectedError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				expectedMsg := fmt.Sprintf("key %q not found", tt.key)
				if !strings.Contains(err.Error(), expectedMsg) {
					t.Errorf("expected error to contain %q, got %q", expectedMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Read captured output
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := strings.TrimSpace(buf.String())

			if output != tt.expectedOut {
				t.Errorf("expected output %q, got %q", tt.expectedOut, output)
			}
		})
	}
}

// TestRunConfigSet verifies setting config values.
func TestRunConfigSet(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         string
		expectedError bool
		errorContains string
	}{
		{
			name:  "set valid path",
			key:   "path",
			value: "/new/path",
		},
		{
			name:  "set valid hours",
			key:   "hours",
			value: "48",
		},
		{
			name:  "set valid author",
			key:   "author",
			value: "John Doe",
		},
		{
			name:  "set valid format",
			key:   "format",
			value: "text",
		},
		{
			name:  "set valid style",
			key:   "style",
			value: "markdown",
		},
		{
			name:  "set valid boolean",
			key:   "today",
			value: "true",
		},
		{
			name:          "set invalid format value",
			key:           "format",
			value:         "invalid",
			expectedError: true,
			errorContains: "format must be",
		},
		{
			name:          "set invalid style value",
			key:           "style",
			value:         "invalid",
			expectedError: true,
			errorContains: "style must be",
		},
		{
			name:          "set negative max_depth",
			key:           "max_depth",
			value:         "-1",
			expectedError: true,
			errorContains: "max-depth must be at least 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for clean state
			viper.Reset()

			// Create temp config file for testing
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, ".gohome.json")

			// Set config file path in viper (this is critical!)
			viper.SetConfigFile(configPath)

			// Set HOME env to temp dir so SaveToFile writes there
			oldHome := os.Getenv("HOME")
			_ = os.Setenv("HOME", tmpDir)
			defer func() { _ = os.Setenv("HOME", oldHome) }()

			// Initialize config with some defaults in viper
			viper.Set("path", "~/workspace")
			viper.Set("format", "table")
			viper.Set("style", "normal")
			viper.Set("max_depth", 2)

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := runConfigSet(configSetCmd, []string{tt.key, tt.value})

			_ = w.Close()
			os.Stdout = old

			// Check error expectation
			if tt.expectedError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Read captured output
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := strings.TrimSpace(buf.String())

			// Verify success message
			expectedMsg := fmt.Sprintf("✅ Set %s = %s", tt.key, tt.value)
			if !strings.Contains(output, expectedMsg) {
				t.Errorf("expected output to contain %q, got %q", expectedMsg, output)
			}

			// Verify config file was created in HOME
			homeConfigPath := filepath.Join(tmpDir, ".gohome.json")
			if _, err := os.Stat(homeConfigPath); os.IsNotExist(err) {
				t.Errorf("config file should have been created at %s", homeConfigPath)
			}
		})
	}
}

// TestRunConfigReset verifies the reset command.
func TestRunConfigReset(t *testing.T) {
	t.Run("reset with confirmation yes", func(t *testing.T) {
		// Reset viper for clean state
		viper.Reset()

		// Create temp config file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".gohome.json")
		viper.SetConfigFile(configPath)

		// Create the config file
		if err := os.WriteFile(configPath, []byte("{}"), 0o600); err != nil {
			t.Fatalf("failed to create test config file: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatal("test config file should exist before reset")
		}

		// Simulate user input "y"
		oldStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r
		_, _ = w.WriteString("y\n")
		_ = w.Close()

		// Capture stdout
		oldStdout := os.Stdout
		rOut, wOut, _ := os.Pipe()
		os.Stdout = wOut

		err := runConfigReset(configResetCmd, []string{})

		_ = wOut.Close()
		os.Stdout = oldStdout
		os.Stdin = oldStdin

		if err != nil {
			t.Fatalf("runConfigReset failed: %v", err)
		}

		// Read output
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(rOut)
		output := buf.String()

		// Verify success message
		if !strings.Contains(output, "Configuration reset to defaults") {
			t.Errorf("expected success message in output, got: %s", output)
		}

		// Verify file was deleted
		if _, err := os.Stat(configPath); !os.IsNotExist(err) {
			t.Error("config file should have been deleted")
		}
	})

	t.Run("reset with confirmation no", func(t *testing.T) {
		// Reset viper for clean state
		viper.Reset()

		// Create temp config file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, ".gohome.json")
		viper.SetConfigFile(configPath)

		// Create the config file
		if err := os.WriteFile(configPath, []byte("{}"), 0o600); err != nil {
			t.Fatalf("failed to create test config file: %v", err)
		}

		// Simulate user input "n"
		oldStdin := os.Stdin
		r, w, _ := os.Pipe()
		os.Stdin = r
		_, _ = w.WriteString("n\n")
		_ = w.Close()

		// Capture stdout
		oldStdout := os.Stdout
		rOut, wOut, _ := os.Pipe()
		os.Stdout = wOut

		err := runConfigReset(configResetCmd, []string{})

		_ = wOut.Close()
		os.Stdout = oldStdout
		os.Stdin = oldStdin

		if err != nil {
			t.Fatalf("runConfigReset failed: %v", err)
		}

		// Read output
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(rOut)
		output := buf.String()

		// Verify canceled message
		if !strings.Contains(output, "Canceled") {
			t.Errorf("expected canceled message in output, got: %s", output)
		}

		// Verify file still exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("config file should not have been deleted")
		}
	})

	t.Run("reset with no config file", func(t *testing.T) {
		// Reset viper for clean state
		viper.Reset()

		// Ensure no config file
		viper.SetConfigFile("")

		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := runConfigReset(configResetCmd, []string{})

		_ = w.Close()
		os.Stdout = old

		if err != nil {
			t.Fatalf("runConfigReset failed: %v", err)
		}

		// Read output
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		output := buf.String()

		// Verify message
		expectedMsg := "No config file found. Nothing to reset."
		if !strings.Contains(output, expectedMsg) {
			t.Errorf("expected output to contain %q, got: %s", expectedMsg, output)
		}
	})
}

// TestConfigGetCmdValidation verifies the get command argument validation.
func TestConfigGetCmdValidation(t *testing.T) {
	if configGetCmd.Args == nil {
		t.Fatal("configGetCmd.Args should not be nil")
	}

	// Test with correct number of args
	err := configGetCmd.Args(configGetCmd, []string{"path"})
	if err != nil {
		t.Errorf("expected no error with 1 arg, got: %v", err)
	}

	// Test with no args
	err = configGetCmd.Args(configGetCmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	// Test with too many args
	err = configGetCmd.Args(configGetCmd, []string{"path", "extra"})
	if err == nil {
		t.Error("expected error with 2 args")
	}
}

// TestConfigSetCmdValidation verifies the set command argument validation.
func TestConfigSetCmdValidation(t *testing.T) {
	if configSetCmd.Args == nil {
		t.Fatal("configSetCmd.Args should not be nil")
	}

	// Test with correct number of args
	err := configSetCmd.Args(configSetCmd, []string{"path", "/test/path"})
	if err != nil {
		t.Errorf("expected no error with 2 args, got: %v", err)
	}

	// Test with no args
	err = configSetCmd.Args(configSetCmd, []string{})
	if err == nil {
		t.Error("expected error with 0 args")
	}

	// Test with one arg
	err = configSetCmd.Args(configSetCmd, []string{"path"})
	if err == nil {
		t.Error("expected error with 1 arg")
	}

	// Test with too many args
	err = configSetCmd.Args(configSetCmd, []string{"path", "value", "extra"})
	if err == nil {
		t.Error("expected error with 3 args")
	}
}

// TestConfigListCmdProperties verifies properties of the list command.
func TestConfigListCmdProperties(t *testing.T) {
	if configListCmd.Use != "list" {
		t.Errorf("expected Use to be 'list', got %q", configListCmd.Use)
	}

	if configListCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if configListCmd.Long == "" {
		t.Error("Long description should not be empty")
	}

	if configListCmd.RunE == nil {
		t.Error("RunE function should not be nil")
	}
}

// TestConfigResetCmdProperties verifies properties of the reset command.
func TestConfigResetCmdProperties(t *testing.T) {
	if configResetCmd.Use != "reset" {
		t.Errorf("expected Use to be 'reset', got %q", configResetCmd.Use)
	}

	if configResetCmd.Short == "" {
		t.Error("Short description should not be empty")
	}

	if configResetCmd.Long == "" {
		t.Error("Long description should not be empty")
	}

	if configResetCmd.RunE == nil {
		t.Error("RunE function should not be nil")
	}
}
