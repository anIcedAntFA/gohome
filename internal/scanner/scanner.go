// Package scanner provides utilities for discovering git repositories in the filesystem.
package scanner

import (
	"os"
	"path/filepath"
)

// ScanGitRepos finds directory paths that contain a .git folder.
// It scans up to maxDepth levels deep to support structures like github.com/{org}/{repo}.
// If maxDepth is 0 or negative, it defaults to 2 levels.
func ScanGitRepos(rootPath string, maxDepth int) ([]string, error) {
	var repos []string

	// Default to 2 levels if invalid value
	if maxDepth <= 0 {
		maxDepth = 2
	}

	// 1. Check root
	if isGitRepo(rootPath) {
		repos = append(repos, rootPath)
		return repos, nil // If root is a git repo, don't scan subdirectories
	}

	// 2. Scan subdirectories recursively (up to maxDepth levels deep)
	repos, err := scanRecursive(rootPath, 0, maxDepth)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

// scanRecursive recursively scans directories up to maxDepth levels.
func scanRecursive(path string, currentDepth, maxDepth int) ([]string, error) {
	var repos []string

	// Stop if we've reached max depth
	if currentDepth > maxDepth {
		return repos, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || shouldSkip(entry.Name()) {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		// Check if this directory is a git repo
		if isGitRepo(fullPath) {
			repos = append(repos, fullPath)
			// Don't scan inside git repos (skip nested repos)
			continue
		}

		// Recursively scan subdirectories
		subRepos, err := scanRecursive(fullPath, currentDepth+1, maxDepth)
		if err != nil {
			// Log error but continue scanning other directories
			continue
		}
		repos = append(repos, subRepos...)
	}

	return repos, nil
}

func isGitRepo(path string) bool {
	gitPath := filepath.Join(path, ".git")
	_, err := os.Stat(gitPath)
	return err == nil
}

func shouldSkip(name string) bool {
	return name == ".git" || name == ".vscode" || name == ".idea"
}
