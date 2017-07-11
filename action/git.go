package action

import (
	"errors"
	"log"
	"strings"
)

type (
	GitActionType struct{}
)

func (GitActionType) prompt() string {
	return "git-actions>"
}

func (GitActionType) menuItems(commitIds []string) ([]menuItem, error) {
	if len(commitIds) > 1 {
		return nil, errors.New("target of git-action must be a single line")
	}

	commitId := strings.TrimSpace(strings.Join(commitIds, " "))
	log.Printf("Commit ID: %s", commitId)

	return []menuItem{
		{Label: "Checkout", Action: "git checkout " + commitId},
		{Label: "Tig", Action: "tig " + commitId},
		{Label: "Delete (safely)", Action: "git branch -d " + commitId},
		{Label: "Copy to Clipboard", Action: "echo -n " + commitId + " | pbcopy"},
	}, nil
}
