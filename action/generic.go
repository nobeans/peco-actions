package action

import (
	cmn "github.com/nobeans/peco-actions/common"
)

type (
	GenericActionType struct{}
)

func (GenericActionType) prompt() string {
	return cmn.Env("PECO_ACTIONS__PROMPT", "generic-actions>")
}

func (GenericActionType) menuItems(_ []string) ([]menuItem, error) {
	return []menuItem{}, nil
}
