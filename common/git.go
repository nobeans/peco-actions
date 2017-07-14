package common

import (
	"fmt"
	"os/exec"
	"strings"
)

func InGitRepository() bool {
	if !CommandExists("git") {
		return false
	}

	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(fmt.Sprintf("%s", out)) == "true"
}
