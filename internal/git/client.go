// Package git provides wrapper functions for executing git commands.
package git

import (
	"os/exec"
	"regexp"
	"strings"
)

// Client handles git command executions
type Client struct{}

// NewClient creates a new git client
func NewClient() *Client {
	return &Client{}
}

// GetUser retrieves the user.name from git config
func (c *Client) GetUser() string {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// sanitizeInput removes potentially dangerous characters from git arguments
func sanitizeInput(input string) string {
	// Allow only alphanumeric, spaces, dots, hyphens, underscores, and @ symbol
	re := regexp.MustCompile(`[^a-zA-Z0-9\s._@-]+`)
	return re.ReplaceAllString(input, "")
}

// GetLogs returns raw commit messages as a slice of strings
func (c *Client) GetLogs(repoPath, author, period string) ([]string, error) {
	// Sanitize inputs to prevent command injection
	safeAuthor := sanitizeInput(author)
	safePeriod := sanitizeInput(period)

	// #nosec G204 -- inputs are sanitized above
	cmd := exec.Command("git", "log",
		"--author="+safeAuthor,
		"--since="+safePeriod,
		"--pretty=format:%s",
		"--no-merges", // Exclude merge commits
	)
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
