// Package version provides build information for the gohome CLI.
package version

import "fmt"

// Build information, injected at build time via ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// String returns the formatted version string.
func String() string {
	return fmt.Sprintf("gohome %s (commit: %s, built: %s)", Version, Commit, Date)
}

// Short returns a short version string (version only).
func Short() string {
	if Version == "dev" {
		return "gohome dev"
	}
	return fmt.Sprintf("gohome %s", Version)
}
