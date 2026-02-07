package editor

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestDetectEditor(t *testing.T) {
	tests := []struct {
		name           string
		visual         string
		editor         string
		gohomeEditor   string
		expectedSubstr string // We check if result contains this
	}{
		{
			name:           "VISUAL takes priority",
			visual:         "vim",
			editor:         "nano",
			gohomeEditor:   "emacs",
			expectedSubstr: "vim",
		},
		{
			name:           "EDITOR fallback when VISUAL empty",
			visual:         "",
			editor:         "nano",
			gohomeEditor:   "emacs",
			expectedSubstr: "nano",
		},
		{
			name:           "GOHOME_EDITOR when VISUAL and EDITOR empty",
			visual:         "",
			editor:         "",
			gohomeEditor:   "emacs",
			expectedSubstr: "emacs",
		},
		{
			name:           "System fallback when all env empty",
			visual:         "",
			editor:         "",
			gohomeEditor:   "",
			expectedSubstr: "", // Will be platform-specific
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Backup original env
			origVisual := os.Getenv("VISUAL")
			origEditor := os.Getenv("EDITOR")
			origGohomeEditor := os.Getenv("GOHOME_EDITOR")
			defer func() {
				_ = os.Setenv("VISUAL", origVisual)              // #nosec G104 -- test cleanup
				_ = os.Setenv("EDITOR", origEditor)              // #nosec G104 -- test cleanup
				_ = os.Setenv("GOHOME_EDITOR", origGohomeEditor) // #nosec G104 -- test cleanup
			}()

			// Set test env
			_ = os.Setenv("VISUAL", tt.visual)              // #nosec G104 -- test setup
			_ = os.Setenv("EDITOR", tt.editor)              // #nosec G104 -- test setup
			_ = os.Setenv("GOHOME_EDITOR", tt.gohomeEditor) // #nosec G104 -- test setup

			result := detectEditor()

			if tt.expectedSubstr == "" {
				// For system fallback, just check it's not empty
				if result == "" {
					t.Errorf("Expected non-empty editor, got empty string")
				}
			} else {
				if !strings.Contains(result, tt.expectedSubstr) {
					t.Errorf("Expected result to contain %q, got %q", tt.expectedSubstr, result)
				}
			}
		})
	}
}

func TestGetPlatformFallback(t *testing.T) {
	result := getPlatformFallback()

	switch runtime.GOOS {
	case "windows":
		if result != "notepad" {
			t.Errorf("Expected 'notepad' on Windows, got %q", result)
		}
	case "darwin", "linux":
		if result != "nano" {
			t.Errorf("Expected 'nano' on %s, got %q", runtime.GOOS, result)
		}
	default:
		if result == "" {
			t.Errorf("Expected non-empty fallback for %s", runtime.GOOS)
		}
	}
}

func TestGetEditorCandidates(t *testing.T) {
	candidates := getEditorCandidates()

	if len(candidates) == 0 {
		t.Error("Expected at least one editor candidate")
	}

	// Check for common editors
	expectedEditors := []string{"code", "nano", "vim"}
	for _, expected := range expectedEditors {
		found := false
		for _, candidate := range candidates {
			if strings.Contains(candidate, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find %q in candidates", expected)
		}
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.editor == "" {
		t.Error("Expected client to have an editor set")
	}
}

func TestNewClientWithEditor(t *testing.T) {
	customEditor := "my-custom-editor --wait"
	client := NewClientWithEditor(customEditor)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.editor != customEditor {
		t.Errorf("Expected editor %q, got %q", customEditor, client.editor)
	}
}

func TestGetEditor(t *testing.T) {
	customEditor := "vim -c 'set syntax=markdown'"
	client := NewClientWithEditor(customEditor)

	result := client.GetEditor()

	if result != customEditor {
		t.Errorf("Expected %q, got %q", customEditor, result)
	}
}

func TestGenerateInstructions(t *testing.T) {
	client := NewClient()
	instructions := client.generateInstructions()

	if instructions == "" {
		t.Error("Expected non-empty instructions")
	}

	// Check for key elements
	requiredStrings := []string{
		"gohome Report Editor",
		"Instructions:",
		"Delete lines",
		"Save and close",
		"--------------------------------------------------",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(instructions, required) {
			t.Errorf("Expected instructions to contain %q", required)
		}
	}
}

func TestRemoveInstructions(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Remove full instructions",
			input: `# gohome Report Editor
# Instructions...
# --------------------------------------------------

Actual content here
More content`,
			expected: "\nActual content here\nMore content",
		},
		{
			name:     "No instructions to remove",
			input:    "Just content\nNo instructions",
			expected: "Just content\nNo instructions",
		},
		{
			name: "Instructions without content",
			input: `# gohome Report Editor
# --------------------------------------------------
`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.removeInstructions(tt.input)
			if result != tt.expected {
				t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
			}
		})
	}
}

func TestValidateAndClean(t *testing.T) {
	client := NewClient()

	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name: "Valid content with comments",
			input: `✨ feat: add new feature
# This is a comment
🐛 fix: resolve bug
# Another comment
📝 docs: update README`,
			expected: `✨ feat: add new feature

🐛 fix: resolve bug

📝 docs: update README`,
			expectError: false,
		},
		{
			name: "Only comments",
			input: `# Comment 1
# Comment 2
# Comment 3`,
			expected:    "",
			expectError: true, // Empty after cleaning
		},
		{
			name: "Content with excessive whitespace",
			input: `

✨ feat: add feature


🐛 fix: bug


`,
			expected:    "✨ feat: add feature\n\n🐛 fix: bug",
			expectError: false,
		},
		{
			name:        "Empty input",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name: "Mixed valid content",
			input: `📁 Repository: myapp

✨ feat(auth): add JWT
# User added this comment
🐛 fix(api): CORS issue

📝 Additional Tasks
- Meeting: Sprint Planning`,
			expected: `📁 Repository: myapp

✨ feat(auth): add JWT

🐛 fix(api): CORS issue

📝 Additional Tasks
- Meeting: Sprint Planning`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := client.validateAndClean(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected:\n%q\nGot:\n%q", tt.expected, result)
				}
			}
		})
	}
}

func TestOpenWithInvalidEditor(t *testing.T) {
	// Test with non-existent editor
	client := NewClientWithEditor("this-editor-definitely-does-not-exist-12345")

	content := "Test content"
	_, err := client.Open(content)

	if err == nil {
		t.Error("Expected error when opening with invalid editor, got nil")
	}

	// Error should mention the editor problem
	if !strings.Contains(err.Error(), "editor") {
		t.Errorf("Expected error message to contain 'editor', got: %v", err)
	}
}

func TestOpenWithEmptyEditor(t *testing.T) {
	client := NewClientWithEditor("")

	content := "Test content"
	_, err := client.Open(content)

	if err == nil {
		t.Error("Expected error when editor command is empty, got nil")
	}

	expectedMsg := "editor command is empty"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message to contain %q, got: %v", expectedMsg, err)
	}
}

// TestOpenIntegration is a manual integration test (skipped in CI).
// To run: go test -v -run TestOpenIntegration.
func TestOpenIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Only run if explicitly requested
	if os.Getenv("RUN_EDITOR_INTEGRATION_TEST") != "1" {
		t.Skip("Skipping editor integration test (set RUN_EDITOR_INTEGRATION_TEST=1 to run)")
	}

	client := NewClient()
	testContent := `✨ feat: test feature
🐛 fix: test bug fix
📝 docs: update documentation`

	result, err := client.Open(testContent)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result from editor")
	}

	t.Logf("Editor result:\n%s", result)
}

// Benchmark tests.
func BenchmarkDetectEditor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = detectEditor()
	}
}

func BenchmarkValidateAndClean(b *testing.B) {
	client := NewClient()
	testContent := `✨ feat: add feature
# Comment here
🐛 fix: bug fix
# Another comment
📝 docs: update docs`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.validateAndClean(testContent)
	}
}
