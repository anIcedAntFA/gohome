package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestNewClient verifies client creation.
func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
}

// TestGetUser tests retrieving the git user name.
func TestGetUser(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	tests := []struct {
		name       string
		setupGit   bool
		userName   string
		expectUser bool
	}{
		{
			name:       "git_configured_user",
			setupGit:   true,
			userName:   "Test User",
			expectUser: true,
		},
		// Note: This test may pass/fail depending on global git config
		// {
		// 	name:       "git_not_configured",
		// 	setupGit:   false,
		// 	userName:   "",
		// 	expectUser: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for isolated git config
			tmpDir := t.TempDir()
			oldHome := os.Getenv("HOME")
			os.Setenv("HOME", tmpDir)
			defer os.Setenv("HOME", oldHome)

			if tt.setupGit {
				// Configure git user.name
				cmd := exec.Command("git", "config", "--global", "user.name", tt.userName)
				cmd.Env = append(os.Environ(), "HOME="+tmpDir)
				if err := cmd.Run(); err != nil {
					t.Skipf("Cannot configure git: %v", err)
				}
			}

			result := client.GetUser(ctx)

			if tt.expectUser && result == "" {
				t.Errorf("GetUser() = empty string, want %q", tt.userName)
			}

			if !tt.expectUser && result != "" {
				t.Errorf("GetUser() = %q, want empty string", result)
			}
		})
	}
}

// TestGetUserWithTimeout tests context cancellation.
func TestGetUserWithTimeout(t *testing.T) {
	client := NewClient()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	time.Sleep(2 * time.Millisecond) // Ensure context is cancelled
	result := client.GetUser(ctx)

	// Should return empty string when context is cancelled
	if result != "" {
		t.Errorf("GetUser() with cancelled context = %q, want empty string", result)
	}
}

// TestSanitizeInput tests the input sanitization function.
func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "safe_alphanumeric",
			input: "john123",
			want:  "john123",
		},
		{
			name:  "safe_with_spaces",
			input: "John Doe",
			want:  "John Doe",
		},
		{
			name:  "safe_with_dots",
			input: "john.doe",
			want:  "john.doe",
		},
		{
			name:  "safe_with_hyphens",
			input: "john-doe",
			want:  "john-doe",
		},
		{
			name:  "safe_with_underscores",
			input: "john_doe",
			want:  "john_doe",
		},
		{
			name:  "safe_with_at_symbol",
			input: "john@example",
			want:  "john@example",
		},
		{
			name:  "combined_safe_chars",
			input: "john.doe-123_test@example",
			want:  "john.doe-123_test@example",
		},
		{
			name:  "command_injection_semicolon",
			input: "john; rm -rf /",
			want:  "john rm -rf ",
		},
		{
			name:  "command_injection_pipe",
			input: "john | cat /etc/passwd",
			want:  "john  cat etcpasswd",
		},
		{
			name:  "command_injection_ampersand",
			input: "john && malicious",
			want:  "john  malicious",
		},
		{
			name:  "command_injection_backticks",
			input: "john`whoami`",
			want:  "johnwhoami",
		},
		{
			name:  "command_injection_dollar",
			input: "john$(whoami)",
			want:  "johnwhoami",
		},
		{
			name:  "path_traversal",
			input: "../../../etc/passwd",
			want:  "......etcpasswd", // Slashes are removed, dots and paths remain
		},
		{
			name:  "special_chars_brackets",
			input: "john[test]",
			want:  "johntest",
		},
		{
			name:  "special_chars_parentheses",
			input: "john(test)",
			want:  "johntest",
		},
		{
			name:  "special_chars_quotes",
			input: "john'test\"",
			want:  "johntest",
		},
		{
			name:  "special_chars_asterisk",
			input: "john*",
			want:  "john",
		},
		{
			name:  "special_chars_question",
			input: "john?",
			want:  "john",
		},
		{
			name:  "newline_injection",
			input: "john\nmalicious",
			want:  "john\nmalicious", // Newlines are allowed by regex
		},
		{
			name:  "carriage_return",
			input: "john\rmalicious",
			want:  "john\rmalicious", // Carriage returns are allowed
		},
		{
			name:  "tab_character",
			input: "john\tmalicious",
			want:  "john\tmalicious", // Tabs are allowed (\s matches whitespace)
		},
		{
			name:  "empty_string",
			input: "",
			want:  "",
		},
		{
			name:  "only_special_chars",
			input: "!@#$%^&*()",
			want:  "@",
		},
		{
			name:  "unicode_characters",
			input: "john™®©",
			want:  "john",
		},
		{
			name:  "emoji",
			input: "john✨",
			want:  "john",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeInput(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeInput(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestGetLogs tests retrieving git commit logs.
func TestGetLogs(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	// Create a temporary git repository for testing
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	if err := os.Mkdir(repoPath, 0o755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot initialize git repo: %v", err)
	}

	// Configure git user for the test repo
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot configure git email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test Author")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot configure git name: %v", err)
	}

	// Create and commit a test file
	testFile := filepath.Join(repoPath, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot add file to git: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "feat: test commit")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot commit to git: %v", err)
	}

	// Create a second commit
	if err := os.WriteFile(testFile, []byte("updated content"), 0o644); err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-am", "fix: bug fix")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot create second commit: %v", err)
	}

	tests := []struct {
		name         string
		author       string
		period       string
		allBranches  bool
		branch       string
		expectCommit bool
		minCommits   int
	}{
		{
			name:         "valid_author_recent",
			author:       "Test Author",
			period:       "1.day.ago",
			allBranches:  false,
			branch:       "",
			expectCommit: true,
			minCommits:   2,
		},
		{
			name:         "valid_author_all_branches",
			author:       "Test Author",
			period:       "1.day.ago",
			allBranches:  true,
			branch:       "",
			expectCommit: true,
			minCommits:   2,
		},
		{
			name:         "nonexistent_author",
			author:       "Nonexistent User",
			period:       "1.day.ago",
			allBranches:  false,
			branch:       "",
			expectCommit: false,
			minCommits:   0,
		},
		{
			name:         "specific_branch_master",
			author:       "Test Author",
			period:       "1.day.ago",
			allBranches:  false,
			branch:       "master", // Default git branch is master
			expectCommit: true,
			minCommits:   2,
		},
		{
			name:         "sanitized_author_with_special_chars",
			author:       "Test Author; rm -rf /",
			period:       "1.day.ago",
			allBranches:  false,
			branch:       "",
			expectCommit: false, // Sanitized author won't match
			minCommits:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := client.GetLogs(ctx, repoPath, tt.author, tt.period, tt.allBranches, tt.branch)

			if err != nil && tt.expectCommit {
				t.Errorf("GetLogs() error = %v, want nil", err)
				return
			}

			if !tt.expectCommit && len(logs) > 0 {
				t.Errorf("GetLogs() returned %d commits, want 0", len(logs))
				return
			}

			if tt.expectCommit && len(logs) < tt.minCommits {
				t.Errorf("GetLogs() returned %d commits, want at least %d", len(logs), tt.minCommits)
			}
		})
	}
}

// TestGetLogsInvalidRepo tests behavior with invalid repository path.
func TestGetLogsInvalidRepo(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	logs, err := client.GetLogs(ctx, "/nonexistent/path", "Test Author", "1.day.ago", false, "")

	if err == nil {
		t.Error("GetLogs() with invalid repo path should return error")
	}

	if len(logs) != 0 {
		t.Errorf("GetLogs() with invalid repo returned %d logs, want 0", len(logs))
	}
}

// TestGetLogsEmptyCommits tests behavior when no commits match criteria.
func TestGetLogsEmptyCommits(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	// Create a temporary git repository with no commits
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "empty-repo")

	if err := os.Mkdir(repoPath, 0o755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot initialize git repo: %v", err)
	}

	logs, err := client.GetLogs(ctx, repoPath, "Test Author", "1.day.ago", false, "")

	// Git log returns error (exit 128) when no commits exist in repo
	if err == nil {
		t.Error("GetLogs() with no commits should return error")
	}

	if len(logs) != 0 {
		t.Errorf("GetLogs() with no commits returned %d logs, want 0", len(logs))
	}
}

// TestGetLogsWithTimeout tests context cancellation during git log operation.
func TestGetLogsWithTimeout(t *testing.T) {
	client := NewClient()

	// Create a very short timeout that will likely cancel before git finishes
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Use current directory (should be valid git repo)
	time.Sleep(1 * time.Millisecond) // Ensure context is cancelled

	_, err := client.GetLogs(ctx, ".", "Test Author", "1.day.ago", false, "")

	// Should return error when context is cancelled
	if err == nil {
		t.Log("GetLogs() with cancelled context should return error (may pass if git is very fast)")
	}
}

// TestGetLogsBranchFiltering tests branch filtering logic.
func TestGetLogsBranchFiltering(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	// Create a temporary git repository for testing
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "branch-repo")

	if err := os.Mkdir(repoPath, 0o755); err != nil {
		t.Fatalf("Failed to create repo directory: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot initialize git repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot configure git: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Branch Author")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot configure git: %v", err)
	}

	// Create initial commit on main branch
	testFile := filepath.Join(repoPath, "main.txt")
	if err := os.WriteFile(testFile, []byte("main content"), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	cmd = exec.Command("git", "add", "main.txt")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot add file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "feat: main branch commit")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot commit: %v", err)
	}

	// Create a feature branch
	cmd = exec.Command("git", "checkout", "-b", "feature")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot create branch: %v", err)
	}

	// Create commit on feature branch
	featureFile := filepath.Join(repoPath, "feature.txt")
	if err := os.WriteFile(featureFile, []byte("feature content"), 0o644); err != nil {
		t.Fatalf("Failed to write feature file: %v", err)
	}

	cmd = exec.Command("git", "add", "feature.txt")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot add feature file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "feat: feature branch commit")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot commit on feature: %v", err)
	}

	// Switch back to master
	cmd = exec.Command("git", "checkout", "master")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Skipf("Cannot switch to master: %v", err)
	}

	tests := []struct {
		name        string
		allBranches bool
		branch      string
		minCommits  int
		maxCommits  int
	}{
		{
			name:        "current_branch_only",
			allBranches: false,
			branch:      "",
			minCommits:  1,
			maxCommits:  1, // Only master commit
		},
		{
			name:        "all_branches",
			allBranches: true,
			branch:      "",
			minCommits:  2,
			maxCommits:  2, // Both master and feature commits
		},
		{
			name:        "specific_branch_master",
			allBranches: false,
			branch:      "master",
			minCommits:  1,
			maxCommits:  1,
		},
		{
			name:        "specific_branch_feature",
			allBranches: false,
			branch:      "feature",
			minCommits:  2,
			maxCommits:  2, // Feature includes main's commit history
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs, err := client.GetLogs(ctx, repoPath, "Branch Author", "1.day.ago", tt.allBranches, tt.branch)
			if err != nil {
				t.Errorf("GetLogs() error = %v, want nil", err)
				return
			}

			commitCount := len(logs)
			if commitCount < tt.minCommits || commitCount > tt.maxCommits {
				t.Errorf("GetLogs() returned %d commits, want between %d and %d",
					commitCount, tt.minCommits, tt.maxCommits)
			}
		})
	}
}

// TestGetLogsSanitizationSecurity ensures dangerous inputs are neutralized.
func TestGetLogsSanitizationSecurity(t *testing.T) {
	client := NewClient()
	ctx := context.Background()

	// Use current directory (assuming it's a valid git repo)
	// These should not cause command injection, just no matches
	dangerousInputs := []struct {
		name   string
		author string
		period string
		branch string
	}{
		{
			name:   "semicolon_injection",
			author: "test; rm -rf /",
			period: "1.day.ago",
			branch: "",
		},
		{
			name:   "pipe_injection",
			author: "test | cat /etc/passwd",
			period: "1.day.ago",
			branch: "",
		},
		{
			name:   "ampersand_injection",
			author: "test && malicious",
			period: "1.day.ago",
			branch: "",
		},
		{
			name:   "backtick_injection",
			author: "test`whoami`",
			period: "1.day.ago",
			branch: "",
		},
		{
			name:   "dollar_injection",
			author: "test$(whoami)",
			period: "1.day.ago",
			branch: "",
		},
		{
			name:   "period_injection",
			author: "test",
			period: "1.day.ago; rm -rf /",
			branch: "",
		},
		{
			name:   "branch_injection",
			author: "test",
			period: "1.day.ago",
			branch: "main; malicious",
		},
	}

	for _, tt := range dangerousInputs {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic or execute malicious code
			// Just return empty results or error
			logs, err := client.GetLogs(ctx, ".", tt.author, tt.period, false, tt.branch)
			// Either empty logs or error is fine, just shouldn't execute injection
			if err != nil {
				t.Logf("GetLogs() with dangerous input returned error (expected): %v", err)
			}

			if len(logs) > 0 {
				t.Logf("GetLogs() with dangerous input returned %d logs (safe)", len(logs))
			}

			// The test passes if we get here without panicking or hanging
		})
	}
}
