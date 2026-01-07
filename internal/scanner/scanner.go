// Package scanner provides utilities for discovering git repositories in the filesystem.
package scanner

import (
	"os"
	"path/filepath"
)

// ScanGitRepos finds directory paths that contain a .git folder.
func ScanGitRepos(rootPath string) ([]string, error) {
	var repos []string

	// 1. Check root
	if isGitRepo(rootPath) {
		repos = append(repos, rootPath)
	}

	// 2. Check sub-directories
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if shouldSkip(entry.Name()) {
				continue
			}
			fullPath := filepath.Join(rootPath, entry.Name())
			if isGitRepo(fullPath) {
				repos = append(repos, fullPath)
			}
		}
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
