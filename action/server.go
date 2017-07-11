package action

import (
	"errors"
	"log"
	"strings"
)

type (
	ServerActionType struct{}
)

func (ServerActionType) prompt() string {
	return "server-actions>"
}

func (ServerActionType) menuItems(hosts []string) ([]menuItem, error) {
	if len(hosts) > 1 {
		return nil, errors.New("target of server-action must be a single line")
	}

	host := strings.TrimSpace(strings.Join(hosts, " "))
	log.Printf("Host: %s", host)

	return []menuItem{
		{Label: "Ssh", Action: "ssh " + host},
		{Label: "Ping", Action: "ping " + host},
		{Label: "Traceroute", Action: "traceroute " + host},
		{Label: "Copy IP address to Clipboard", Action: "ping -c 1 " + host + " | grep PING | sed -E 's/.*\\((.*)\\):.*/\\\\1/' | pbcopy"},
	}, nil
}
