package action

import (
	"errors"
	"log"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type ProcessActionType struct{}

func (ProcessActionType) prompt() string {
	return "process-actions>"
}

func (ProcessActionType) menuItems(pids []string) ([]menuItem, error) {
	if len(pids) > 1 {
		return nil, errors.New("target of process-action must be a single line")
	}

	pid := strings.TrimSpace(strings.Join(pids, " "))
	log.Printf("PID: %s", pid)

	items := []menuItem{
		{Label: "Signal HUP  (1:hang up)", Action: "kill -HUP " + pid},
		{Label: "Signal INT  (2:interrupt)", Action: "kill -INT " + pid},
		{Label: "Signal QUIT (3:quit)", Action: "kill -QUIT " + pid},
		{Label: "Signal KILL (9:kill)", Action: "kill -KILL " + pid},
		{Label: "Signal TERM (15:termination)", Action: "kill -TERM " + pid},
	}
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard", Action: "echo -n '" + pid + "' | pbcopy"})
	}
	return items, nil
}
