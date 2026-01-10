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

	// Production tagged release: version matches semantic versioning (with or without 'v' prefix)
	// Examples: v1.0.1, 1.0.1, v2.3.0-beta
	// Show clean format without build details
	if isSemanticVersion(v) {
		return fmt.Sprintf("gohome %s", v)
	}

	// Dev builds or go install without tag: show full details if available
	if c != "none" || d != "unknown" {
		return fmt.Sprintf("gohome %s (commit: %s, built: %s)", v, c, d)
	}

	// Fallback: just version
	return fmt.Sprintf("gohome %s", v)
}

// isSemanticVersion checks if version string matches semantic versioning pattern.
func isSemanticVersion(v string) bool {
	if v == "" || v == "dev" {
		return false
	}

	// Git commit hashes or dirty builds (e.g., "abc123-dirty", "25bd8dd-dirty")
	// are NOT semantic versions
	// Semantic versions use dots (1.0.1), commits use dashes for dirty suffix
	if len(v) >= 7 && v[6] == '-' {
		return false // Likely a commit hash with -dirty suffix
	}

	// Strip 'v' prefix if present
	test := v
	if v[0] == 'v' {
		test = v[1:]
	}

	// Must start with digit and contain a dot (e.g., 1.0.1, not just "1")
	// This distinguishes semantic versions from plain commit hashes
	if len(test) == 0 || test[0] < '0' || test[0] > '9' {
		return false
	}

	// Check for dot (semantic versions always have major.minor format)
	hasDot := false
	for i := 0; i < len(test); i++ {
		if test[i] == '.' {
			hasDot = true
			break
		}
	}

	return hasDot
}

// Short returns a short version string (version only).
func Short() string {
	v := getVersion()
	if v == "dev" {
		return "gohome dev"
	}
	return fmt.Sprintf("gohome %s", v)
}
