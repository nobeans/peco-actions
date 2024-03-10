package action

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type DockerContainerActionType struct{}

func (DockerContainerActionType) prompt() string {
	return "docker-container-actions>"
}

func (DockerContainerActionType) menuItems(lines []string) ([]menuItem, error) {
	if len(lines) > 1 {
		return nil, errors.New("target of docker-container-action must be a single line")
	}

	line := strings.TrimSpace(strings.Join(lines, " "))
	log.Printf("Input line: %s", line)

	// Expected the table format of `docker ps` command
	tokens := regexp.MustCompile(" ([^ ]+)$").FindAllString(line, 1)
	if len(tokens) != 1 {
		return nil, fmt.Errorf("invalid format: %s", line)
	}

	containerName := strings.TrimSpace(tokens[0])
	log.Printf("Container name: %s", containerName)

	items := []menuItem{
		{Label: "Show logs", Action: "docker logs -f " + containerName},
		{Label: "Kill", Action: "docker kill " + containerName},
		{Label: "Exec (sh)", Action: "docker exec -it " + containerName + " sh"},
		{Label: "Exec (bash)", Action: "docker exec -it " + containerName + " bash"},
		{Label: "Exec (docker debug)", Action: "docker debug " + containerName},
	}
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard", Action: "echo -n " + containerName + " | pbcopy"})
	}
	return items, nil
}
