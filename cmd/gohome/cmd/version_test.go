package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/anIcedAntFA/gohome/internal/version"
)

// TestVersionCommand tests the version command execution
func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name           string
		setupVersion   string
		setupCommit    string
		setupDate      string
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:         "full_version_info",
			setupVersion: "v1.2.0",
			setupCommit:  "abc1234",
			setupDate:    "2026-01-18",
			wantContains: []string{
				"gohome version v1.2.0",
				"commit: abc1234",
				"built: 2026-01-18",
			},
		},
		{
			name:         "dev_version",
			setupVersion: "dev",
			setupCommit:  "",
			setupDate:    "",
			wantContains: []string{
				"gohome version dev",
			},
			wantNotContain: []string{
				"commit:",
				"built:",
			},
		},
		{
			name:         "version_with_commit_only",
			setupVersion: "v1.3.0",
			setupCommit:  "def5678",
			setupDate:    "",
			wantContains: []string{
				"gohome version v1.3.0",
				"commit: def5678",
			},
			wantNotContain: []string{
				"built:",
			},
		},
		{
			name:         "version_without_metadata",
			setupVersion: "v2.0.0",
			setupCommit:  "",
			setupDate:    "",
			wantContains: []string{
				"gohome version v2.0.0",
			},
			wantNotContain: []string{
				"commit:",
				"built:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original version info
			origVersion := version.Version
			origCommit := version.Commit
			origDate := version.Date
			defer func() {
				version.Version = origVersion
				version.Commit = origCommit
				version.Date = origDate
			}()

			// Setup test version info
			version.Version = tt.setupVersion
			version.Commit = tt.setupCommit
			version.Date = tt.setupDate

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute command
			runVersion(versionCmd, []string{})

			// Restore stdout and read output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			got := buf.String()

			// Verify expected content
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Output = %q\nwant to contain %q", got, want)
				}
			}

			// Verify excluded content
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("Output = %q\nshould NOT contain %q", got, notWant)
				}
			}
		})
	}
}

// TestVersionCommandProperties tests the command metadata
func TestVersionCommandProperties(t *testing.T) {
	if versionCmd.Use != "version" {
		t.Errorf("versionCmd.Use = %q, want %q", versionCmd.Use, "version")
	}

	if versionCmd.Short == "" {
		t.Error("versionCmd.Short should not be empty")
	}

	if versionCmd.Long == "" {
		t.Error("versionCmd.Long should not be empty")
	}

	if versionCmd.Run == nil {
		t.Error("versionCmd.Run should not be nil")
	}
}
