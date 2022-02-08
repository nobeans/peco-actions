package action

import (
	"errors"
	"log"
	"strings"

	cmn "github.com/nobeans/peco-actions/common"
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

	items := []menuItem{}

	items = append(items, menuItem{Label: "Checkout", Action: "git checkout " + strings.ReplaceAll(commitId, "remotes/origin/", "")})

	if cmn.CommandExists("tig") {
		items = append(items, menuItem{Label: "Tig", Action: "tig " + commitId})
	}

	items = append(items, []menuItem{
		{Label: "Delete (safely)", Action: "git branch -d " + commitId},
		{Label: "Copy to Clipboard", Action: "echo -n " + commitId + " | pbcopy"},
	}...)

	return items, nil
}
