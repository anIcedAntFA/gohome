// Package version provides build information for the gohome CLI.
package version

import (
	"fmt"
	"runtime/debug"
)

// Build information, injected at build time via ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// getVersion returns version from ldflags or VCS info (for go install).
func getVersion() string {
	if Version != "dev" {
		return Version
	}

	// Fallback to VCS info when installed via go install
	if info, ok := debug.ReadBuildInfo(); ok {
		// For go install github.com/user/repo@v1.0.0
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}

		// Try to get from VCS
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && len(setting.Value) >= 7 {
				return setting.Value[:7] // Short commit hash
			}
		}
	}

	return Version
}

// getCommit returns commit hash from ldflags or VCS info.
func getCommit() string {
	if Commit != "none" {
		return Commit
	}

	// Try to get from VCS info
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && len(setting.Value) >= 7 {
				return setting.Value[:7]
			}
		}
	}

	return Commit
}

// getDate returns build date from ldflags or VCS info.
func getDate() string {
	if Date != "unknown" {
		return Date
	}

	// Try to get from VCS info
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.time" {
				return setting.Value
			}
		}
	}

	return Date
}

// String returns the formatted version string.
func String() string {
	v := getVersion()
	c := getCommit()
	d := getDate()

	// For go install builds, show cleaner format if no build info
	if c == "none" && d == "unknown" {
		return fmt.Sprintf("gohome %s", v)
	}

	return fmt.Sprintf("gohome %s (commit: %s, built: %s)", v, c, d)
}

// Short returns a short version string (version only).
func Short() string {
	v := getVersion()
	if v == "dev" {
		return "gohome dev"
	}
	return fmt.Sprintf("gohome %s", v)
}
