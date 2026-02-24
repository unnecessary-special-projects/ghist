package project

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const GhistDir = ".ghist"
const DBFile = "ghist.sqlite"
const ContextFile = "current_context.json"

// FindRoot walks up from startDir to find a directory containing .ghist/.
// Returns the project root (parent of .ghist/) or an error.
func FindRoot(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("resolving path: %w", err)
	}

	for {
		candidate := filepath.Join(dir, GhistDir)
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no %s directory found (walked up to filesystem root)", GhistDir)
		}
		dir = parent
	}
}

// DBPath returns the full path to the SQLite database given a project root.
func DBPath(root string) string {
	return filepath.Join(root, GhistDir, DBFile)
}

// ContextPath returns the full path to current_context.json given a project root.
func ContextPath(root string) string {
	return filepath.Join(root, GhistDir, ContextFile)
}

// GhistDirPath returns the full path to the .ghist/ directory given a project root.
func GhistDirPath(root string) string {
	return filepath.Join(root, GhistDir)
}

// DetectGitHubRepo attempts to detect the GitHub repository URL from the git remote.
// Returns a URL like "https://github.com/owner/repo" or empty string if not found.
func DetectGitHubRepo(root string) string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = root
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return parseGitHubURL(strings.TrimSpace(string(out)))
}

func parseGitHubURL(remote string) string {
	// SSH format: git@github.com:owner/repo.git
	if strings.HasPrefix(remote, "git@github.com:") {
		path := strings.TrimPrefix(remote, "git@github.com:")
		path = strings.TrimSuffix(path, ".git")
		return "https://github.com/" + path
	}
	// HTTPS format: https://github.com/owner/repo.git
	if strings.Contains(remote, "github.com/") {
		return strings.TrimSuffix(remote, ".git")
	}
	return ""
}
