package action

import (
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
	containerNames := make([]string, 0, len(lines))
	for i, line := range lines {
		log.Printf("Input line [%d]: %s", i, line)

		// Expected the table format of `docker ps` command
		tokens := regexp.MustCompile(" ([^ ]+)$").FindAllString(line, 1)
		if len(tokens) != 1 {
			return nil, fmt.Errorf("invalid format: %s", line)
		}

		containerName := strings.TrimSpace(tokens[0])
		log.Printf("Container name [%d]: %s", i, containerName)

		containerNames = append(containerNames, containerName)
	}

	var items []menuItem
	if len(containerNames) == 1 {
		items = append(items, menuItem{Label: "Show logs", Action: "docker logs -f " + containerNames[0]})
		items = append(items, menuItem{Label: "Exec (sh)", Action: "docker exec -it " + containerNames[0] + " sh"})
		items = append(items, menuItem{Label: "Exec (bash)", Action: "docker exec -it " + containerNames[0] + " bash"})
		items = append(items, menuItem{Label: "Exec (docker debug)", Action: "docker debug " + containerNames[0]})
	}
	items = append(items, menuItem{Label: "Kill", Action: "docker kill " + strings.Join(containerNames, " ")})
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard", Action: "echo -n " + strings.Join(containerNames, " ") + " | pbcopy"})
	}
	items = append(items, RenderAdhocMenuItems(strings.Join(containerNames, " "))...)
	return items, nil
}
