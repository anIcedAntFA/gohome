// Package sys provides system-level utilities like clipboard operations.
package sys

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CopyToClipboard automatically detects OS and copies text to clipboard.
func CopyToClipboard(ctx context.Context, text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.CommandContext(ctx, "pbcopy")
	case "windows":
		// Windows uses clip command
		cmd = exec.CommandContext(ctx, "clip")
	case "linux", "freebsd", "openbsd", "netbsd":
		// Smart logic for Linux clipboard
		if isWayland() {
			// Prefer Wayland
			cmd = exec.CommandContext(ctx, "wl-copy")
		} else {
			// Fallback to X11, prefer xclip then xsel
			//nolint:gocritic // if-else is clearer than switch for command availability checks
			if isCommandAvailable("xclip") {
				cmd = exec.CommandContext(ctx, "xclip", "-selection", "clipboard")
			} else if isCommandAvailable("xsel") {
				cmd = exec.CommandContext(ctx, "xsel", "--clipboard", "--input")
			} else {
				// Last resort: try wl-copy (in case environment variables are incorrect)
				cmd = exec.CommandContext(ctx, "wl-copy")
			}
		}
	default:
		return nil // Unsupported OS
	}

	// Attach text to command's Stdin
	cmd.Stdin = strings.NewReader(text)

	// Execute command
	return cmd.Run()
}

// isWayland checks if user is running Wayland.
func isWayland() bool {
	// Check common Wayland environment variables
	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	xdgSessionType := os.Getenv("XDG_SESSION_TYPE")

	return waylandDisplay != "" || strings.EqualFold(xdgSessionType, "wayland")
}

// isCommandAvailable checks if tool exists in PATH.
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
