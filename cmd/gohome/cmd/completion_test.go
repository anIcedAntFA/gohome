package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestCompletionCommand verifies completion script generation for all supported shells.
// Tests that the generated scripts contain expected shell-specific content.
func TestCompletionCommand(t *testing.T) {
	tests := []struct {
		name           string
		shell          string
		wantContains   []string // Expected strings that must appear in output
		wantNotContain []string // Strings that should NOT appear in output
	}{
		{
			name:  "bash_completion",
			shell: "bash",
			wantContains: []string{
				"# bash completion",
				"_gohome()",
			},
		},
		{
			name:  "zsh_completion",
			shell: "zsh",
			wantContains: []string{
				"#compdef gohome",
				"_gohome()",
			},
		},
		{
			name:  "fish_completion",
			shell: "fish",
			wantContains: []string{
				"# fish completion for gohome",
				"complete -c gohome",
			},
		},
		{
			name:  "powershell_completion",
			shell: "powershell",
			wantContains: []string{
				"Register-ArgumentCompleter",
				"-CommandName",
				"gohome",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup: Create isolated test root command
			testRootCmd := &cobra.Command{
				Use:   "gohome",
				Short: "Test root command",
			}

			// Create completion command with actual implementation
			testCompletionCmd := &cobra.Command{
				Use:       "completion [bash|zsh|fish|powershell]",
				Short:     "Generate shell completion script",
				ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
				Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				Run: func(cmd *cobra.Command, args []string) {
					switch args[0] {
					case "bash":
						_ = cmd.Root().GenBashCompletion(cmd.OutOrStdout())
					case "zsh":
						_ = cmd.Root().GenZshCompletion(cmd.OutOrStdout())
					case "fish":
						_ = cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
					case "powershell":
						_ = cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
					}
				},
			}
			testRootCmd.AddCommand(testCompletionCmd)

			// Capture output using buffer
			var buf bytes.Buffer
			testRootCmd.SetOut(&buf)
			testCompletionCmd.SetOut(&buf)

			// Execute: Run the completion command for the specified shell
			testCompletionCmd.Run(testCompletionCmd, []string{tt.shell})

			got := buf.String()

			// Assert: Verify expected content is present
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Output does not contain %q\nGot first 200 chars: %q",
						want, truncate(got, 200))
				}
			}

			// Assert: Verify excluded content is absent
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("Output should NOT contain %q", notWant)
				}
			}

			// Assert: Verify output is not empty
			if got == "" {
				t.Error("completion output is empty")
			}
		})
	}
}

// TestCompletionCommandProperties verifies the completion command metadata and configuration.
func TestCompletionCommandProperties(t *testing.T) {
	if completionCmd.Use != "completion [bash|zsh|fish|powershell]" {
		t.Errorf("completionCmd.Use = %q, want %q",
			completionCmd.Use, "completion [bash|zsh|fish|powershell]")
	}

	if completionCmd.Short == "" {
		t.Error("completionCmd.Short should not be empty")
	}

	if completionCmd.Long == "" {
		t.Error("completionCmd.Long should not be empty")
	}

	if completionCmd.Run == nil {
		t.Error("completionCmd.Run should not be nil")
	}

	expectedValidArgs := []string{"bash", "zsh", "fish", "powershell"}
	if len(completionCmd.ValidArgs) != len(expectedValidArgs) {
		t.Errorf("completionCmd.ValidArgs length = %d, want %d",
			len(completionCmd.ValidArgs), len(expectedValidArgs))
	}

	for i, arg := range expectedValidArgs {
		if i >= len(completionCmd.ValidArgs) || completionCmd.ValidArgs[i] != arg {
			t.Errorf("completionCmd.ValidArgs[%d] = %q, want %q",
				i, completionCmd.ValidArgs[i], arg)
		}
	}

	if !completionCmd.DisableFlagsInUseLine {
		t.Error("completionCmd.DisableFlagsInUseLine should be true")
	}
}

// TestCompletionCommandValidation tests argument validation.
func TestCompletionCommandValidation(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "valid_bash",
			args:      []string{"bash"},
			wantError: false,
		},
		{
			name:      "valid_zsh",
			args:      []string{"zsh"},
			wantError: false,
		},
		{
			name:      "valid_fish",
			args:      []string{"fish"},
			wantError: false,
		},
		{
			name:      "valid_powershell",
			args:      []string{"powershell"},
			wantError: false,
		},
		{
			name:      "invalid_shell",
			args:      []string{"cmd"},
			wantError: true,
		},
		{
			name:      "no_args",
			args:      []string{},
			wantError: true,
		},
		{
			name:      "too_many_args",
			args:      []string{"bash", "extra"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new command instance for each test
			testCmd := &cobra.Command{
				Use:       "completion [bash|zsh|fish|powershell]",
				ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
				Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
				Run:       func(_ *cobra.Command, _ []string) {},
			}

			// Suppress output
			var buf bytes.Buffer
			testCmd.SetOut(&buf)
			testCmd.SetErr(&buf)

			// Set args and execute
			testCmd.SetArgs(tt.args)
			err := testCmd.Execute()

			gotError := err != nil
			if gotError != tt.wantError {
				t.Errorf("Execute() error = %v, wantError = %v", err, tt.wantError)
			}
		})
	}
}

// TestCompletionLongHelp tests that long help contains installation instructions.
func TestCompletionLongHelp(t *testing.T) {
	requiredSections := []string{
		"Bash:",
		"Zsh:",
		"Fish:",
		"PowerShell:",
		"source <(gohome completion",
		"gohome completion bash >",
		"gohome completion zsh >",
		"gohome completion fish >",
		"gohome completion powershell >",
	}

	longHelp := completionCmd.Long

	for _, section := range requiredSections {
		if !strings.Contains(longHelp, section) {
			t.Errorf("Long help missing section: %q", section)
		}
	}
}

// truncate truncates a string to maxLen characters for better error messages.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// TestCompletionInit tests that the completion command is properly initialized.
func TestCompletionInit(t *testing.T) {
	// Verify completionCmd is added to rootCmd
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "completion" {
			found = true
			break
		}
	}

	if !found {
		t.Error("completion command not found in root command")
	}
}

// TestCompletionRunFunction tests the Run function with actual execution.
func TestCompletionRunFunction(t *testing.T) {
	tests := []struct {
		name         string
		shell        string
		wantInOutput string
	}{
		{"run_bash", "bash", "_gohome"},
		{"run_zsh", "zsh", "#compdef"},
		{"run_fish", "fish", "complete"},
		{"run_powershell", "powershell", "Register-ArgumentCompleter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Setup Buffer để hứng output
			buf := new(bytes.Buffer)

			// 2. CẤU HÌNH LẠI CÁI completionCmd THẬT
			// Reset output của lệnh thật về buffer của mình
			completionCmd.SetOut(buf)
			completionCmd.SetErr(buf)

			// 3. Setup root command giả để gắn completionCmd vào (nếu cần context root)
			// Lưu ý: completionCmd là biến global, cẩn thận khi add vào root mới nhiều lần
			// Cách an toàn nhất là dùng trực tiếp RunE/Run hoặc Mock Root

			// Mẹo: Do completionCmd dùng cmd.Root(), ta cần đảm bảo nó có cha.
			// Nhưng completionCmd đã được add vào rootCmd thật ở hàm init().
			// Ta có thể dùng trực tiếp rootCmd để execute.

			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs([]string{"completion", tt.shell})

			// 4. CHẠY LỆNH THẬT
			err := rootCmd.Execute()
			// Hoặc nếu muốn test unit nhỏ hơn chỉ cho completionCmd:
			// completionCmd.Run(completionCmd, []string{tt.shell})
			// Nhưng cách trên (rootCmd.Execute) giống thực tế hơn.
			// 5. Assert
			if err != nil {
				t.Fatalf("Execute() failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.wantInOutput) {
				t.Errorf("Expected output to contain %q, got: %q", tt.wantInOutput, truncate(output, 50))
			}
		})
	}
}

// TestCompletionDisableFlagsInUseLine verifies the DisableFlagsInUseLine setting.
func TestCompletionDisableFlagsInUseLine(t *testing.T) {
	if !completionCmd.DisableFlagsInUseLine {
		t.Error("DisableFlagsInUseLine should be true to hide flags in usage line")
	}
}
