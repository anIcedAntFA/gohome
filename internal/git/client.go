package git

import (
	"os/exec"
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

// GetLogs returns raw commit messages as a slice of strings
func (c *Client) GetLogs(repoPath, author, period string) ([]string, error) {
	cmd := exec.Command("git", "log",
		"--author="+author,
		"--since="+period,
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
