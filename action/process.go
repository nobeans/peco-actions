package action

import (
	"log"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type ProcessActionType struct{}

func (ProcessActionType) prompt() string {
	return "process-actions>"
}

func (ProcessActionType) menuItems(lines []string) ([]menuItem, error) {
	pids := strings.TrimSpace(strings.Join(lines, " ")) // support multiple PIDs
	log.Printf("PID: %s", pids)

	items := []menuItem{
		{Label: "Signal HUP  (1:hang up)", Action: "kill -HUP " + pids},
		{Label: "Signal INT  (2:interrupt)", Action: "kill -INT " + pids},
		{Label: "Signal QUIT (3:quit)", Action: "kill -QUIT " + pids},
		{Label: "Signal KILL (9:kill)", Action: "kill -KILL " + pids},
		{Label: "Signal TERM (15:termination)", Action: "kill -TERM " + pids},
	}
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard", Action: "echo -n '" + pids + "' | pbcopy"})
	}
	items = append(items, RenderAdhocMenuItems(pids)...)
	return items, nil
}
