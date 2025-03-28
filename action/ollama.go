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

func (OllamaActionType) menuItems(lines []string) ([]menuItem, error) {
	if len(lines) == 0 || len(lines) > 1 {
		return nil, errors.New("must be a single line")
	}

	modelName := strings.TrimSpace(lines[0])
	log.Printf("Model name: %s", modelName)

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
