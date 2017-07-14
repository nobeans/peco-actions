package common

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func InGitRepository(path string) bool {
	if !CommandExists("git") {
		return false
	}

	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if IsDirectory(path) {
		cmd.Dir = path
	} else {
		cmd.Dir = filepath.Dir(path)
	}
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(fmt.Sprintf("%s", out)) == "true"
}

func GitRepositoryRoot(path string) bool {
	if !CommandExists("git") {
		return false
	}

	if !IsDirectory(path) {
		return false
	}

	cmd := exec.Command("git", "rev-parse", "--show-cdup")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(fmt.Sprintf("%s", out)) == ""
}
