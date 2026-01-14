package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsGitRepo tests the isGitRepo function.
func TestIsGitRepo(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(string) error
		expected bool
	}{
		{
			name: "valid git repo",
			setup: func(dir string) error {
				return os.Mkdir(filepath.Join(dir, ".git"), 0o755)
			},
			expected: true,
		},
		{
			name: "not a git repo",
			setup: func(_ string) error {
				return nil // No .git folder
			},
			expected: false,
		},
		{
			name: "git file instead of directory",
			setup: func(dir string) error {
				f, err := os.Create(filepath.Join(dir, ".git"))
				if err != nil {
					return err
				}
				return f.Close()
			},
			expected: true, // .git file exists (submodule case)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()

			// Setup test case
			if err := tt.setup(tempDir); err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			// Test
			result := isGitRepo(tempDir)
			if result != tt.expected {
				t.Errorf("isGitRepo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestShouldSkip tests the shouldSkip function.
func TestShouldSkip(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{"skip .git", ".git", true},
		{"skip .vscode", ".vscode", true},
		{"skip .idea", ".idea", true},
		{"don't skip normal dir", "src", false},
		{"don't skip node_modules", "node_modules", false},
		{"don't skip .github", ".github", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkip(tt.dirName)
			if result != tt.expected {
				t.Errorf("shouldSkip(%q) = %v, want %v", tt.dirName, result, tt.expected)
			}
		})
	}
}

// TestScanRecursive tests the scanRecursive function with various depth scenarios.
func TestScanRecursive(t *testing.T) {
	tests := []struct {
		name          string
		maxDepth      int
		setup         func(string) error
		expectedCount int
	}{
		{
			name:     "depth 0 - only current level",
			maxDepth: 0,
			setup: func(root string) error {
				// Create git repos at root level
				if err := createGitRepo(root, "repo1"); err != nil {
					return err
				}
				if err := createGitRepo(root, "repo2"); err != nil {
					return err
				}
				return nil
			},
			expectedCount: 2,
		},
		{
			name:     "depth 1 - one level deep",
			maxDepth: 1,
			setup: func(root string) error {
				// Level 0: repos at root
				if err := createGitRepo(root, "repo1"); err != nil {
					return err
				}
				// Level 1: repos in subdirs
				if err := os.Mkdir(filepath.Join(root, "org1"), 0o755); err != nil {
					return err
				}
				if err := createGitRepo(filepath.Join(root, "org1"), "repo2"); err != nil {
					return err
				}
				// Level 2: should be ignored
				if err := os.Mkdir(filepath.Join(root, "org1", "repo2", "subdir"), 0o755); err != nil {
					return err
				}
				if err := createGitRepo(filepath.Join(root, "org1", "repo2", "subdir"), "repo3"); err != nil {
					return err
				}
				return nil
			},
			expectedCount: 2, // repo1 and repo2, not repo3
		},
		{
			name:     "depth 2 - github.com/org/repo structure",
			maxDepth: 2,
			setup: func(root string) error {
				// Level 0: github.com
				// Level 1: org1, org2
				org1 := filepath.Join(root, "org1")
				org2 := filepath.Join(root, "org2")
				if err := os.Mkdir(org1, 0o755); err != nil {
					return err
				}
				if err := os.Mkdir(org2, 0o755); err != nil {
					return err
				}

				// Level 2: repos
				if err := createGitRepo(org1, "repo1"); err != nil {
					return err
				}
				if err := createGitRepo(org1, "repo2"); err != nil {
					return err
				}
				if err := createGitRepo(org2, "repo3"); err != nil {
					return err
				}

				// Level 3: should be ignored (nested inside git repo)
				if err := os.Mkdir(filepath.Join(org1, "repo1", "nested"), 0o755); err != nil {
					return err
				}
				if err := createGitRepo(filepath.Join(org1, "repo1", "nested"), "repo4"); err != nil {
					return err
				}

				return nil
			},
			expectedCount: 3, // repo1, repo2, repo3 (not repo4)
		},
		{
			name:     "skip special directories",
			maxDepth: 1,
			setup: func(root string) error {
				// Create repos
				if err := createGitRepo(root, "repo1"); err != nil {
					return err
				}

				// Create special dirs that should be skipped
				if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
					return err
				}
				if err := os.Mkdir(filepath.Join(root, ".vscode"), 0o755); err != nil {
					return err
				}
				if err := os.Mkdir(filepath.Join(root, ".idea"), 0o755); err != nil {
					return err
				}

				// Create repos inside skipped dirs (should not be detected)
				if err := createGitRepo(filepath.Join(root, ".vscode"), "hidden-repo"); err != nil {
					return err
				}

				return nil
			},
			expectedCount: 1, // Only repo1
		},
		{
			name:     "empty directory",
			maxDepth: 2,
			setup: func(_ string) error {
				return nil // Empty directory
			},
			expectedCount: 0,
		},
		{
			name:     "only non-git directories",
			maxDepth: 2,
			setup: func(root string) error {
				if err := os.Mkdir(filepath.Join(root, "folder1"), 0o755); err != nil {
					return err
				}
				if err := os.Mkdir(filepath.Join(root, "folder1", "subfolder"), 0o755); err != nil {
					return err
				}
				return nil
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()

			// Setup test case
			if err := tt.setup(tempDir); err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			// Test
			repos, err := scanRecursive(tempDir, 0, tt.maxDepth)
			if err != nil {
				t.Fatalf("scanRecursive() error = %v", err)
			}

			if len(repos) != tt.expectedCount {
				t.Errorf("scanRecursive() found %d repos, want %d", len(repos), tt.expectedCount)
				t.Logf("Repos found: %v", repos)
			}
		})
	}
}

// TestScanGitRepos tests the main ScanGitRepos function.
func TestScanGitRepos(t *testing.T) {
	tests := []struct {
		name          string
		maxDepth      int
		setup         func(string) error
		expectedCount int
	}{
		{
			name:     "root is git repo",
			maxDepth: 2,
			setup: func(root string) error {
				// Make root a git repo
				if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
					return err
				}
				// Create nested repos (should be ignored)
				if err := os.Mkdir(filepath.Join(root, "nested"), 0o755); err != nil {
					return err
				}
				return createGitRepo(filepath.Join(root, "nested"), "repo1")
			},
			expectedCount: 1, // Only root repo
		},
		{
			name:     "default depth handling - zero",
			maxDepth: 0, // Should default to 2
			setup: func(root string) error {
				org := filepath.Join(root, "org")
				if err := os.Mkdir(org, 0o755); err != nil {
					return err
				}
				return createGitRepo(org, "repo1")
			},
			expectedCount: 1,
		},
		{
			name:     "default depth handling - negative",
			maxDepth: -1, // Should default to 2
			setup: func(root string) error {
				org := filepath.Join(root, "org")
				if err := os.Mkdir(org, 0o755); err != nil {
					return err
				}
				return createGitRepo(org, "repo1")
			},
			expectedCount: 1,
		},
		{
			name:     "deep nesting with maxDepth=3",
			maxDepth: 3,
			setup: func(root string) error {
				// Level 1
				l1 := filepath.Join(root, "level1")
				if err := os.Mkdir(l1, 0o755); err != nil {
					return err
				}
				// Level 2
				l2 := filepath.Join(l1, "level2")
				if err := os.Mkdir(l2, 0o755); err != nil {
					return err
				}
				// Level 3
				l3 := filepath.Join(l2, "level3")
				if err := os.Mkdir(l3, 0o755); err != nil {
					return err
				}
				// Repos at each level
				if err := createGitRepo(l1, "repo1"); err != nil {
					return err
				}
				if err := createGitRepo(l2, "repo2"); err != nil {
					return err
				}
				return createGitRepo(l3, "repo3")
			},
			expectedCount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Setup test case
			if err := tt.setup(tempDir); err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			// Test
			repos, err := ScanGitRepos(tempDir, tt.maxDepth)
			if err != nil {
				t.Fatalf("ScanGitRepos() error = %v", err)
			}

			if len(repos) != tt.expectedCount {
				t.Errorf("ScanGitRepos() found %d repos, want %d", len(repos), tt.expectedCount)
				t.Logf("Repos found: %v", repos)
			}
		})
	}
}

// Helper function to create a git repo structure for testing.
func createGitRepo(parentDir, repoName string) error {
	repoPath := filepath.Join(parentDir, repoName)
	if err := os.MkdirAll(repoPath, 0o755); err != nil {
		return err
	}
	return os.Mkdir(filepath.Join(repoPath, ".git"), 0o755)
}
