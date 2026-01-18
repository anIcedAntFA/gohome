package viper

import (
	"strings"
	"testing"
)

// TestValidate_ValidConfig tests that a valid config passes all validation checks.
func TestValidate_ValidConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "default config",
			cfg: Config{
				Days:     1,
				Path:     ".",
				MaxDepth: 2,
				Format:   "text",
				Style:    "normal",
			},
		},
		{
			name: "table format with markdown style",
			cfg: Config{
				Days:     7,
				Path:     "/home/user/repos",
				MaxDepth: 5,
				Format:   "table",
				Style:    "markdown",
			},
		},
		{
			name: "all branches enabled",
			cfg: Config{
				Days:        1,
				Path:        ".",
				MaxDepth:    1,
				Format:      "text",
				Style:       "normal",
				AllBranches: true,
				Branch:      "",
			},
		},
		{
			name: "specific branch filter",
			cfg: Config{
				Days:     1,
				Path:     ".",
				MaxDepth: 10,
				Format:   "text",
				Style:    "normal",
				Branch:   "main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err != nil {
				t.Errorf("Validate() failed for valid config: %v", err)
			}
		})
	}
}

// TestValidate_NegativeTimePeriods tests validation of negative time period values.
func TestValidate_NegativeTimePeriods(t *testing.T) {
	tests := []struct {
		name        string
		cfg         Config
		expectedErr string
	}{
		{
			name: "negative hours",
			cfg: Config{
				Hours:    -5,
				Path:     ".",
				MaxDepth: 2,
				Format:   "text",
				Style:    "normal",
			},
			expectedErr: "hours cannot be negative",
		},
		{
			name: "negative days",
			cfg: Config{
				Days:     -1,
				Path:     ".",
				MaxDepth: 2,
				Format:   "text",
				Style:    "normal",
			},
			expectedErr: "days cannot be negative",
		},
		{
			name: "negative weeks",
			cfg: Config{
				Weeks:    -2,
				Path:     ".",
				MaxDepth: 2,
				Format:   "text",
				Style:    "normal",
			},
			expectedErr: "weeks cannot be negative",
		},
		{
			name: "negative months",
			cfg: Config{
				Months:   -3,
				Path:     ".",
				MaxDepth: 2,
				Format:   "text",
				Style:    "normal",
			},
			expectedErr: "months cannot be negative",
		},
		{
			name: "negative years",
			cfg: Config{
				Years:    -1,
				Path:     ".",
				MaxDepth: 2,
				Format:   "text",
				Style:    "normal",
			},
			expectedErr: "years cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if err == nil {
				t.Errorf("Validate() should fail for %s", tt.name)
				return
			}
			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("Validate() error = %q, expected to contain %q", err.Error(), tt.expectedErr)
			}
		})
	}
}

// TestValidate_MaxDepth tests validation of MaxDepth boundaries.
func TestValidate_MaxDepth(t *testing.T) {
	tests := []struct {
		name        string
		maxDepth    int
		shouldFail  bool
		expectedErr string
	}{
		{
			name:       "min valid depth",
			maxDepth:   1,
			shouldFail: false,
		},
		{
			name:       "max valid depth",
			maxDepth:   10,
			shouldFail: false,
		},
		{
			name:       "normal depth",
			maxDepth:   5,
			shouldFail: false,
		},
		{
			name:        "too small depth",
			maxDepth:    0,
			shouldFail:  true,
			expectedErr: "max-depth must be at least 1",
		},
		{
			name:        "negative depth",
			maxDepth:    -1,
			shouldFail:  true,
			expectedErr: "max-depth must be at least 1",
		},
		{
			name:        "too large depth",
			maxDepth:    11,
			shouldFail:  true,
			expectedErr: "max-depth too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Days:     1,
				Path:     ".",
				MaxDepth: tt.maxDepth,
				Format:   "text",
				Style:    "normal",
			}

			err := cfg.Validate()
			if tt.shouldFail {
				if err == nil {
					t.Errorf("Validate() should fail for maxDepth=%d", tt.maxDepth)
					return
				}
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("Validate() error = %q, expected to contain %q", err.Error(), tt.expectedErr)
				}
			} else if err != nil {
				t.Errorf("Validate() should pass for maxDepth=%d, got error: %v", tt.maxDepth, err)
			}
		})
	}
}

// TestValidate_EmptyPath tests that empty path is rejected.
func TestValidate_EmptyPath(t *testing.T) {
	cfg := Config{
		Days:     1,
		Path:     "",
		MaxDepth: 2,
		Format:   "text",
		Style:    "normal",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() should fail for empty path")
		return
	}
	if !strings.Contains(err.Error(), "path cannot be empty") {
		t.Errorf("Validate() error = %q, expected to contain 'path cannot be empty'", err.Error())
	}
}

// TestValidate_InvalidFormat tests validation of format enum values.
func TestValidate_InvalidFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"invalid format json", "json"},
		{"invalid format xml", "xml"},
		{"empty format", ""},
		{"uppercase TEXT", "TEXT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Days:     1,
				Path:     ".",
				MaxDepth: 2,
				Format:   tt.format,
				Style:    "normal",
			}

			err := cfg.Validate()
			if err == nil {
				t.Errorf("Validate() should fail for format=%q", tt.format)
				return
			}
			if !strings.Contains(err.Error(), "invalid format") {
				t.Errorf("Validate() error = %q, expected to contain 'invalid format'", err.Error())
			}
		})
	}
}

// TestValidate_InvalidStyle tests validation of style enum values.
func TestValidate_InvalidStyle(t *testing.T) {
	tests := []struct {
		name  string
		style string
	}{
		{"invalid style fancy", "fancy"},
		{"empty style", ""},
		{"uppercase NORMAL", "NORMAL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Days:     1,
				Path:     ".",
				MaxDepth: 2,
				Format:   "text",
				Style:    tt.style,
			}

			err := cfg.Validate()
			if err == nil {
				t.Errorf("Validate() should fail for style=%q", tt.style)
				return
			}
			if !strings.Contains(err.Error(), "invalid style") {
				t.Errorf("Validate() error = %q, expected to contain 'invalid style'", err.Error())
			}
		})
	}
}

// TestValidate_BranchConflict tests that AllBranches and Branch cannot be used together.
func TestValidate_BranchConflict(t *testing.T) {
	tests := []struct {
		name        string
		allBranches bool
		branch      string
		shouldFail  bool
	}{
		{
			name:        "both set",
			allBranches: true,
			branch:      "main",
			shouldFail:  true,
		},
		{
			name:        "only AllBranches",
			allBranches: true,
			branch:      "",
			shouldFail:  false,
		},
		{
			name:        "only Branch",
			allBranches: false,
			branch:      "feature/test",
			shouldFail:  false,
		},
		{
			name:        "neither set",
			allBranches: false,
			branch:      "",
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				Days:        1,
				Path:        ".",
				MaxDepth:    2,
				Format:      "text",
				Style:       "normal",
				AllBranches: tt.allBranches,
				Branch:      tt.branch,
			}

			err := cfg.Validate()
			if tt.shouldFail {
				if err == nil {
					t.Errorf("Validate() should fail when AllBranches=%v and Branch=%q", tt.allBranches, tt.branch)
					return
				}
				if !strings.Contains(err.Error(), "cannot use --all-branches and --branch together") {
					t.Errorf("Validate() error = %q, expected branch conflict error", err.Error())
				}
			} else if err != nil {
				t.Errorf("Validate() should pass, got error: %v", err)
			}
		})
	}
}
