package action

import (
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type (
	GenericActionType struct{}
)

func (GenericActionType) prompt() string {
	return common.Env("PECO_ACTIONS__PROMPT", "generic-actions>")
}

func (GenericActionType) menuItems(lines []string) ([]menuItem, error) {
	items := []menuItem{}
	for _, line := range lines {
		items = append(items, RenderAdhocMenuItems(strings.TrimSpace(line))...)
	}
	return items, nil
}
