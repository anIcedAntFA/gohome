// Package editor handles opening external text editors for interactive content editing.
package editor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Client handles opening external editors.
type Client struct {
	editor string
}

// NewClient creates a new editor client with auto-detected editor.
func NewClient() *Client {
	return &Client{
		editor: detectEditor(),
	}
}

// NewClientWithEditor creates a new editor client with a specific editor command.
func NewClientWithEditor(editorCmd string) *Client {
	return &Client{
		editor: editorCmd,
	}
}

// detectEditor returns the preferred editor from environment variables or system defaults.
// Priority: VISUAL > EDITOR > Detected Editors > Fallback.
func detectEditor() string {
	// 1. Check VISUAL env (preferred for full-screen editors)
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}

	// 2. Check EDITOR env (traditional Unix standard)
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	// 3. Check GOHOME_EDITOR env (gohome-specific)
	if editor := os.Getenv("GOHOME_EDITOR"); editor != "" {
		return editor
	}

	// 4. Try to detect available editors
	editors := getEditorCandidates()
	for _, editor := range editors {
		parts := strings.Fields(editor)
		if _, err := exec.LookPath(parts[0]); err == nil {
			return editor
		}
	}

	// 5. Platform-specific fallback
	return getPlatformFallback()
}

// getEditorCandidates returns a list of editor commands to try.
func getEditorCandidates() []string {
	return []string{
		"code --wait",          // VS Code
		"code-insiders --wait", // VS Code Insiders
		"subl -w",              // Sublime Text
		"atom --wait",          // Atom
		"nano",                 // Nano (user-friendly)
		"vim",                  // Vim
		"vi",                   // Vi (always available on Unix)
		"emacs",                // Emacs
	}
}

// getPlatformFallback returns the best fallback editor for the current platform.
func getPlatformFallback() string {
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	case "darwin":
		return "nano"
	default: // linux, freebsd, etc.
		return "nano"
	}
}

// Open opens the given content in the configured editor.
// Returns the modified content after the editor closes.
func (c *Client) Open(content string) (string, error) {
	// Add helpful instructions at the top
	instructions := c.generateInstructions()
	fullContent := instructions + content

	// Create temporary file with .txt extension for syntax highlighting
	tmpFile, err := os.CreateTemp("", "gohome-report-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write content to temp file
	if _, writeErr := tmpFile.WriteString(fullContent); writeErr != nil {
		_ = tmpFile.Close() // #nosec G104 -- already handling write error
		return "", fmt.Errorf("failed to write to temp file: %w", writeErr)
	}

	if closeErr := tmpFile.Close(); closeErr != nil {
		return "", fmt.Errorf("failed to close temp file: %w", closeErr)
	}

	// Parse editor command (handle commands with arguments like "code --wait")
	parts := strings.Fields(c.editor)
	if len(parts) == 0 {
		return "", fmt.Errorf("editor command is empty")
	}

	editorBinary := parts[0]
	editorArgs := parts[1:]
	editorArgs = append(editorArgs, tmpFile.Name())

	// Create editor command
	//nolint:noctx // Interactive command requires direct terminal control
	cmd := exec.Command(editorBinary, editorArgs...) // #nosec G204 -- editorBinary is validated via detectEditor()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run editor and wait for it to close
	if runErr := cmd.Run(); runErr != nil {
		return "", fmt.Errorf("editor exited with error: %w", runErr)
	}

	// Read modified content
	modifiedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read modified content: %w", err)
	}

	// Remove instructions and validate content
	finalContent := c.removeInstructions(string(modifiedContent))
	return c.validateAndClean(finalContent)
}

// generateInstructions creates helpful instructions for the editor.
func (c *Client) generateInstructions() string {
	return `# gohome Report Editor
#
# Instructions:
#   - Delete lines you don't want to include
#   - Edit commit messages (keep emoji format for better readability)
#   - Add custom lines (will be preserved in output)
#   - Lines starting with '#' are comments (will be removed from output)
#
# Tips:
#   - Save and close the editor to continue
#   - To cancel: close without saving or delete all content
#
# --------------------------------------------------
`
}

// removeInstructions removes the instruction header from the content.
func (c *Client) removeInstructions(content string) string {
	lines := strings.Split(content, "\n")
	contentStart := 0

	// Find the end of instructions (line with dashes)
	for i, line := range lines {
		if strings.Contains(line, "--------------------------------------------------") {
			contentStart = i + 1
			break
		}
	}

	if contentStart > 0 && contentStart < len(lines) {
		return strings.Join(lines[contentStart:], "\n")
	}

	return content
}

// validateAndClean validates and cleans the edited content.
func (c *Client) validateAndClean(content string) (string, error) {
	var cleanLines []string
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()

		// Skip comment lines (starting with #)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			// But preserve empty line in place of comment for formatting
			cleanLines = append(cleanLines, "")
			continue
		}

		// Preserve all other lines (including empty ones for formatting)
		cleanLines = append(cleanLines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan content: %w", err)
	}

	// Join lines and trim excessive whitespace at start/end
	result := strings.Join(cleanLines, "\n")
	result = strings.TrimSpace(result)

	// Remove multiple consecutive empty lines (keep max 2)
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	if result == "" {
		return "", fmt.Errorf("editing canceled: no content provided (all lines were deleted or file was empty)")
	}

	return result, nil
}

// GetEditor returns the editor command that will be used.
func (c *Client) GetEditor() string {
	return c.editor
}
