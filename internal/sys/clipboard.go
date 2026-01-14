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
		// Check for WSL first, then Wayland, then X11
		switch {
		case isWSL():
			// WSL uses Windows' clip.exe to access system clipboard
			cmd = exec.CommandContext(ctx, "clip.exe")
		case isWayland():
			// Prefer Wayland
			cmd = exec.CommandContext(ctx, "wl-copy")
		case isCommandAvailable("xclip"):
			// Fallback to X11, prefer xclip
			cmd = exec.CommandContext(ctx, "xclip", "-selection", "clipboard")
		case isCommandAvailable("xsel"):
			// Then xsel
			cmd = exec.CommandContext(ctx, "xsel", "--clipboard", "--input")
		default:
			// Last resort: try wl-copy (in case environment variables are incorrect)
			cmd = exec.CommandContext(ctx, "wl-copy")
		}
	default:
		return nil // Unsupported OS
	}

	// Attach text to command's Stdin
	cmd.Stdin = strings.NewReader(text)

	// Execute command
	return cmd.Run()
}

// isWSL checks if the binary is running inside Windows Subsystem for Linux.
func isWSL() bool {
	// Method 1: Check for WSL_DISTRO_NAME environment variable (common in newer WSL versions)
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}

	// Method 2: Check /proc/version for "microsoft" string
	data, err := os.ReadFile("/proc/version")
	if err == nil {
		content := strings.ToLower(string(data))
		if strings.Contains(content, "microsoft") || strings.Contains(content, "wsl") {
			return true
		}
	}

	return false
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
