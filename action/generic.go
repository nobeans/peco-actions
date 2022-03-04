package action

import (
	"github.com/nobeans/peco-actions/common"
)

type (
	GenericActionType struct{}
)

func (GenericActionType) prompt() string {
	return common.Env("PECO_ACTIONS__PROMPT", "generic-actions>")
}

func (GenericActionType) menuItems(_ []string) ([]menuItem, error) {
	return []menuItem{}, nil
}
