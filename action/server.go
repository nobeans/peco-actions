package action

import (
	"errors"
	"log"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type ServerActionType struct{}

func (ServerActionType) prompt() string {
	return "server-actions>"
}

func (ServerActionType) menuItems(lines []string) ([]menuItem, error) {
	if len(lines) == 0 || len(lines) > 1 {
		return nil, errors.New("must be a single line")
	}

	host := strings.TrimSpace(lines[0])
	log.Printf("Host: %s", host)

	var items = []menuItem{}
	if common.CommandExists("ssh") {
		items = append(items, menuItem{Label: "Ssh", Action: "ssh " + host})
	}
	if common.CommandExists("ping") {
		items = append(items, menuItem{Label: "Ping", Action: "ping " + host})
	}
	if common.CommandExists("traceroute") {
		items = append(items, menuItem{Label: "Traceroute", Action: "traceroute " + host})
	}
	if common.CommandExists("tracert") {
		items = append(items, menuItem{Label: "Traceroute", Action: "tracert " + host})
	}
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy IP address to Clipboard", Action: "ping -c 1 " + host + " | grep PING | sed -E 's/.*\\((.*)\\):.*/\\\\1/' | pbcopy"})
	}
	return items, nil
}
