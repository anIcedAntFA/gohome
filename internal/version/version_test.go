package version

import (
	"strings"
	"testing"
)

// TestString tests the formatted version output
func TestString(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		commit         string
		date           string
		wantContains   []string // Strings that should be in output
		wantNotContain []string // Strings that should NOT be in output
	}{
		{
			name:           "production_release",
			version:        "v1.2.0",
			commit:         "abc123",
			date:           "2026-01-18",
			wantContains:   []string{"gohome", "v1.2.0"},
			wantNotContain: []string{"commit:", "built:"}, // Clean format for releases
		},
		{
			name:           "production_release_no_v_prefix",
			version:        "1.3.0",
			commit:         "def456",
			date:           "2026-01-18",
			wantContains:   []string{"gohome", "1.3.0"},
			wantNotContain: []string{"commit:", "built:"},
		},
		{
			name:         "dev_build_with_commit",
			version:      "dev",
			commit:       "abc1234",
			date:         "2026-01-18T10:00:00Z",
			wantContains: []string{"gohome", "dev", "commit:", "abc1234", "built:", "2026-01-18"},
		},
		{
			name:         "commit_hash_version",
			version:      "25bd8dd",
			commit:       "25bd8dd",
			date:         "2026-01-18",
			wantContains: []string{"gohome", "25bd8dd", "commit:", "built:"},
		},
		{
			name:         "dirty_build",
			version:      "abc123-dirty",
			commit:       "abc123",
			date:         "unknown",
			wantContains: []string{"gohome", "abc123-dirty"},
		},
		{
			name:           "fallback_dev_only",
			version:        "dev",
			commit:         "none",
			date:           "unknown",
			wantContains:   []string{"gohome", "dev"},
			wantNotContain: []string{"commit:", "built:", "none", "unknown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion, origCommit, origDate := Version, Commit, Date
			defer func() {
				Version, Commit, Date = origVersion, origCommit, origDate
			}()

			// Set test values
			Version = tt.version
			Commit = tt.commit
			Date = tt.date

			got := String()

			// Check required strings
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("String() = %q, want to contain %q", got, want)
				}
			}

			// Check excluded strings
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("String() = %q, should NOT contain %q", got, notWant)
				}
			}
		})
	}
}

// TestShort tests the short version output
func TestShort(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "release_version",
			version: "v1.2.0",
			want:    "gohome v1.2.0",
		},
		{
			name:    "dev_version",
			version: "dev",
			want:    "gohome dev",
		},
		{
			name:    "commit_hash",
			version: "abc1234",
			want:    "gohome abc1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore
			origVersion := Version
			defer func() { Version = origVersion }()

			Version = tt.version
			got := Short()

			if got != tt.want {
				t.Errorf("Short() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestIsSemanticVersion tests semantic version detection
func TestIsSemanticVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		// Valid semantic versions
		{name: "v_prefix_major_minor_patch", version: "v1.2.3", want: true},
		{name: "no_v_prefix", version: "1.2.3", want: true},
		{name: "major_minor_only", version: "1.2", want: true},
		{name: "with_build_metadata", version: "1.0.0+20230118", want: true},

		// NOT semantic versions (based on implementation)
		// Note: Current implementation doesn't fully support semver with pre-release/build
		// This is intentional to distinguish from commit hashes with dirty suffix
		{name: "with_pre_release", version: "v1.2.0-beta", want: false}, // Has dash like dirty commits
		{name: "empty_string", version: "", want: false},
		{name: "dev", version: "dev", want: false},
		{name: "commit_hash_7_chars", version: "abc1234", want: false},
		{name: "commit_hash_with_dirty", version: "abc123-dirty", want: false},
		{name: "commit_hash_longer_with_dirty", version: "25bd8dd-dirty", want: false},
		{name: "no_dots", version: "1", want: false},
		{name: "no_dots_v_prefix", version: "v1", want: false},
		{name: "starts_with_letter", version: "abc.def", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSemanticVersion(tt.version)
			if got != tt.want {
				t.Errorf("isSemanticVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

// TestGetVersion tests version retrieval
func TestGetVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "ldflags_set",
			version: "v1.2.0",
			want:    "v1.2.0",
		},
		{
			name:    "dev_fallback",
			version: "dev",
			want:    "dev", // Falls back to "dev" if no VCS info
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origVersion := Version
			defer func() { Version = origVersion }()

			Version = tt.version
			got := getVersion()

			if got != tt.want {
				t.Errorf("getVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestGetCommit tests commit hash retrieval
func TestGetCommit(t *testing.T) {
	tests := []struct {
		name   string
		commit string
		want   string
	}{
		{
			name:   "ldflags_set",
			commit: "abc1234567",
			want:   "abc1234567",
		},
		{
			name:   "none_fallback",
			commit: "none",
			want:   "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origCommit := Commit
			defer func() { Commit = origCommit }()

			Commit = tt.commit
			got := getCommit()

			if got != tt.want {
				t.Errorf("getCommit() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestGetDate tests build date retrieval
func TestGetDate(t *testing.T) {
	tests := []struct {
		name string
		date string
		want string
	}{
		{
			name: "ldflags_set",
			date: "2026-01-18T10:00:00Z",
			want: "2026-01-18T10:00:00Z",
		},
		{
			name: "unknown_fallback",
			date: "unknown",
			want: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origDate := Date
			defer func() { Date = origDate }()

			Date = tt.date
			got := getDate()

			if got != tt.want {
				t.Errorf("getDate() = %q, want %q", got, tt.want)
			}
		})
	}
}
