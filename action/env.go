package action

import (
	"errors"
	"log"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type EnvActionType struct{}

func (EnvActionType) prompt() string {
	return "env-actions>"
}

func (EnvActionType) menuItems(lines []string) ([]menuItem, error) {
	if len(lines) == 0 || len(lines) > 1 {
		return nil, errors.New("must be a single line")
	}

	entry := strings.TrimSpace(lines[0])
	tokens := strings.Split(entry, "=")
	name := tokens[0]
	value := tokens[1]
	log.Printf("Env: %s = %s", name, value)

	items := []menuItem{
		{Label: "Unset", Action: "unset " + name},
	}
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard (full)", Action: "echo -n '" + entry + "' | pbcopy"})
		items = append(items, menuItem{Label: "Copy to Clipboard (only name)", Action: "echo -n '" + name + "' | pbcopy"})
		items = append(items, menuItem{Label: "Copy to Clipboard (only value)", Action: "echo -n '" + value + "' | pbcopy"})
	}
	items = append(items, RenderAdhocMenuItems(entry)...)
	return items, nil
}
