package action

import (
	"errors"
	"fmt"
	shellwords "github.com/mattn/go-shellwords"
	cmn "github.com/nobeans/peco-actions/common"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type (
	ActionType interface {
		menuItems(lines []string) ([]menuItem, error)
		prompt() string
	}

	menuItem struct {
		Label  string
		Action string
	}
)

func ResolveAction(actionType ActionType, r io.Reader) (string, error) {
	lines, err := cmn.ReadLines(r)
	if err != nil {
		return "", err
	}
	if len(lines) <= 0 {
		return "", errors.New("no input")
	}
	log.Printf("Read lines: %#v", lines)

	menuItems, err := renderMenuItems(actionType, lines)
	if err != nil {
		return "", err
	}
	log.Printf("Menu items: %#v", menuItems)

	action, err := selectSingleActionByPeco(menuItems, actionType.prompt())
	if err != nil {
		return "", err
	}
	log.Printf("Resolved action: %#v", action)

	return action, nil
}

func renderMenuItems(actionType ActionType, lines []string) ([]menuItem, error) {
	menuItems, err := actionType.menuItems(lines)
	if err != nil {
		return nil, err
	}
	log.Printf("Menu items: %#v", menuItems)

	// You can specify adhoc menu items via environment variable PECO_ACTIONS__ADHOC_MENU
	// e.g. export PECO_ACTIONS__ADHOC_MENU="A=B;C=D"
	adhocMenu := cmn.Env("PECO_ACTIONS__ADHOC_MENU", "")
	if len(adhocMenu) > 0 {
		adhocItems := []menuItem{}
		for _, adhocLine := range strings.Split(adhocMenu, ";") {
			tokens := cmn.Map(strings.SplitN(adhocLine, "=", 2), strings.TrimSpace)
			adhocItems = append(adhocItems, menuItem{
				Label:  tokens[0],
				Action: tokens[1],
			})
		}
		menuItems = append(adhocItems, menuItems...) // append to top
		log.Printf("Menu items (applied adhoc): %#v", menuItems)
	}

	return menuItems, nil
}

func selectSingleActionByPeco(menuItems []menuItem, pecoPrompt string) (string, error) {
	cmd := exec.Command("peco", "--prompt", pecoPrompt)

	// You can specify peco options via environment variable PECO_ACTIONS__PECO_OPTS
	// e.g. export PECO_ACTIONS__PECO_OPTS="--layout bottom-up"
	pecoOpts := cmn.Env("PECO_ACTIONS__PECO_OPTS", "")
	if len(pecoOpts) > 0 {
		parsedOpts, err := shellwords.Parse(pecoOpts)
		if err != nil {
			return "", err
		}
		for _, token := range parsedOpts {
			cmd.Args = append(cmd.Args, token)
		}
	}

	stdin, _ := cmd.StdinPipe()
	io.WriteString(stdin, formatMenu(menuItems))
	stdin.Close()
	out, _ := cmd.Output()

	// Parse an action part from a menu line
	action := strings.TrimSpace(regexp.MustCompile("(?m)^.*> ").ReplaceAllLiteralString(fmt.Sprintf("%s", out), ""))
	log.Printf("Selected action: %s", strconv.Quote(action))

	// Check if it's a single line
	if strings.Contains(action, "\n") {
		return "", errors.New("could not select multiple actions")
	}

	return action, nil
}

func formatMenu(menuItems []menuItem) string {
	maxLabelLen := 0
	for _, item := range menuItems {
		if len(item.Label) > maxLabelLen {
			maxLabelLen = len(item.Label)
		}
	}

	menuLines := []string{}
	for _, item := range menuItems {
		menuLines = append(menuLines, cmn.PadLeft(item.Label+" ", maxLabelLen+10, ".")+" > "+item.Action)
	}

	return strings.Join(menuLines, "\n")
}
