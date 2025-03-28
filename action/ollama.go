package action

import (
	"errors"
	"log"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type OllamaActionType struct{}

func (OllamaActionType) prompt() string {
	return "ollama-actions>"
}

func (OllamaActionType) menuItems(modelNames []string) ([]menuItem, error) {
	if len(modelNames) > 1 {
		return nil, errors.New("target of ollama-action must be a single line")
	}

	modelName := strings.TrimSpace(strings.Join(modelNames, " "))
	log.Printf("Model name: %s", modelNames)

	items := []menuItem{
		{Label: "Show", Action: "ollama show " + modelName},
		{Label: "Run", Action: "ollama run " + modelName},
		{Label: "Stop", Action: "ollama stop " + modelName},
		{Label: "Remove", Action: "ollama rm " + modelName},
	}
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard (full)", Action: "echo -n " + modelName + " | pbcopy"})
		items = append(items, menuItem{Label: "Copy to Clipboard (only base name)", Action: "echo -n " + strings.Split(modelName, ":")[0] + " | pbcopy"})
	}
	return items, nil
}
