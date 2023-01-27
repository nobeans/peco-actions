package action

import (
	"errors"
	"log"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type GitActionType struct{}

func (GitActionType) prompt() string {
	return "git-actions>"
}

func (GitActionType) menuItems(commitIds []string) ([]menuItem, error) {
	if len(commitIds) > 1 {
		return nil, errors.New("target of git-action must be a single line")
	}

	commitId := strings.TrimSpace(strings.Join(commitIds, " "))
	log.Printf("Commit ID: %s", commitId)

	items := []menuItem{
		{Label: "Switch", Action: "git switch " + strings.ReplaceAll(commitId, "remotes/origin/", "")},
	}
	if common.CommandExists("tig") {
		items = append(items, menuItem{Label: "Tig", Action: "tig " + commitId})
	}
	items = append(items, menuItem{Label: "Delete (safely)", Action: "git branch -d " + commitId})
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard", Action: "echo -n " + commitId + " | pbcopy"})
	}
	return items, nil
}
