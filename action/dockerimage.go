package action

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type DockerImageActionType struct{}

func (DockerImageActionType) prompt() string {
	return "docker-image-actions>"
}

func (DockerImageActionType) menuItems(lines []string) ([]menuItem, error) {
	keys := make([]string, 0, len(lines))
	for i, line := range lines {
		log.Printf("Input line [%d]: %s", i, line)

		// Expected the table format of `docker images` command
		tokens := regexp.MustCompile(" +").Split(line, -1)
		if len(tokens) < 2 {
			return nil, fmt.Errorf("invalid format [%d]: %s", i, line)
		}

		imageName := strings.TrimSpace(tokens[0])
		log.Printf("Image name [%d]: %s", i, imageName)

		tag := strings.TrimSpace(tokens[1])
		log.Printf("Tag [%d]: %s", i, tag)

		keys = append(keys, fmt.Sprintf("%s:%s", imageName, tag))
	}

	items := []menuItem{
		{Label: "Delete", Action: "docker rmi " + strings.Join(keys, " ")},
	}
	if len(keys) == 1 {
		imageName := keys[0]
		items = append(items, menuItem{Label: "Exec (docker degbug)", Action: "docker debug " + imageName})
	}
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard", Action: "echo -n " + strings.Join(keys, " ") + " | pbcopy"})
	}
	return items, nil
}
