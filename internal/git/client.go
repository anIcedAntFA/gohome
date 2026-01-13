// Package git provides wrapper functions for executing git commands.
package git

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
)

// Client handles git command executions.
type Client struct{}

// NewClient creates a new git client.
func NewClient() *Client {
	return &Client{}
}

// GetUser retrieves the user.name from git config.
func (c *Client) GetUser(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// sanitizeInput removes potentially dangerous characters from git arguments.
func sanitizeInput(input string) string {
	// Allow only alphanumeric, spaces, dots, hyphens, underscores, and @ symbol
	re := regexp.MustCompile(`[^a-zA-Z0-9\s._@-]+`)
	return re.ReplaceAllString(input, "")
}

// GetLogs returns raw commit messages as a slice of strings.
// If allBranches is true, it includes commits from all local branches.
// If branch is specified, it filters commits from that specific branch.
func (c *Client) GetLogs(ctx context.Context, repoPath, author, period string, allBranches bool, branch string) ([]string, error) {
	// Sanitize inputs to prevent command injection
	safeAuthor := sanitizeInput(author)
	safePeriod := sanitizeInput(period)
	safeBranch := sanitizeInput(branch)

	// Build git log arguments
	args := []string{
		"log",
		"--author=" + safeAuthor,
		"--since=" + safePeriod,
		"--pretty=format:%s",
		"--no-merges", // Exclude merge commits
	}

	// Add branch filtering (mutually exclusive with --branches)
	if branch != "" && safeBranch != "" {
		// Filter by specific branch
		args = append(args, safeBranch)
	} else if allBranches {
		// Include commits from all local branches
		args = append(args, "--branches")
	}

	// #nosec G204 -- inputs are sanitized above
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	strOutput := strings.TrimSpace(string(output))
	if strOutput == "" {
		return []string{}, nil
	}

	return strings.Split(strOutput, "\n"), nil
}
