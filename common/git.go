package common

import (
	"fmt"
	"os/exec"
	"strings"
)

func CwdInGitRepository() bool {
	if !CommandExists("git") {
		return false
	}

	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(fmt.Sprintf("%s", out)) == "true"
}
