package action

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	shellwords "github.com/mattn/go-shellwords"
	cmn "github.com/nobeans/peco-actions/common"
)

type Type interface {
	menuItems(lines []string) ([]menuItem, error)
	prompt() string
}

type menuItem struct {
	Label  string
	Action string
}

func Resolve(actionType Type, r io.Reader) (string, error) {
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

	act, err := selectSingleActionByPeco(menuItems, actionType.prompt())
	if err != nil {
		return "", err
	}
	log.Printf("Resolved action: %#v", act)

	return act, nil
}

func RenderAdhocMenuItems(target string) []menuItem {
	// You can specify adhoc menu items via environment variable PECO_ACTIONS__ADHOC_MENU
	// "TARGET" can be replaced with the selected target.
	// e.g. export PECO_ACTIONS__ADHOC_MENU="Show file name=echo TARGET;Show file info=ls -al TARGET"
	adhocMenu := cmn.Env("PECO_ACTIONS__ADHOC_MENU", "")
	if len(adhocMenu) == 0 {
		return nil
	}

	var items []menuItem
	for _, adhocLine := range strings.Split(adhocMenu, ";") {
		tokens := cmn.Map(strings.SplitN(adhocLine, "=", 2), strings.TrimSpace)
		items = append(items, menuItem{
			Label:  tokens[0],
			Action: strings.ReplaceAll(tokens[1], "TARGET", target),
		})
	}
	return items
}

func renderMenuItems(actionType Type, lines []string) ([]menuItem, error) {
	items, err := actionType.menuItems(lines)
	if err != nil {
		return nil, err
	}
	log.Printf("Menu items: %#v", items)
	return items, nil
}

func selectSingleActionByPeco(menuItems []menuItem, pecoPrompt string) (string, error) {
	cmdName := cmn.Env("PECO_ACTIONS__COMMAND", "peco")
	if cmdName == "fzf" {
		pecoPrompt += " "
	}
	cmd := exec.Command(cmdName, "--prompt", pecoPrompt)

	// This required supporting for fzf.
	cmd.Stderr = os.Stderr

	// You can specify command options via environment variable PECO_ACTIONS__COMMAND_OPTS
	// e.g. export PECO_ACTIONS__COMMAND_OPTS="--layout bottom-up"
	opts := cmn.Env("PECO_ACTIONS__COMMAND_OPTS", "")
	if len(opts) > 0 {
		parsedOpts, err := shellwords.Parse(opts)
		if err != nil {
			return "", err
		}
		for _, token := range parsedOpts {
			cmd.Args = append(cmd.Args, token)
		}
	}

	stdin, _ := cmd.StdinPipe()
	_, _ = io.WriteString(stdin, formatMenu(menuItems))
	_ = stdin.Close()
	out, _ := cmd.Output()
	log.Printf("Command output: %s", out)

	// Parse an action part from a menu line
	act := strings.TrimSpace(regexp.MustCompile("(?m)^.*> ").ReplaceAllLiteralString(string(out), ""))
	log.Printf("Selected action: %s", strconv.Quote(act))

	// Check if it's a single line
	if strings.Contains(act, "\n") {
		return "", errors.New("could not select multiple actions")
	}

	return act, nil
}

func formatMenu(menuItems []menuItem) string {
	maxLabelLen := 0
	for _, item := range menuItems {
		if len(item.Label) > maxLabelLen {
			maxLabelLen = len(item.Label)
		}
	}

	var menuLines []string
	for _, item := range menuItems {
		menuLines = append(menuLines, cmn.PadLeft(item.Label+" ", maxLabelLen+10, ".")+" > "+item.Action)
	}

	return strings.Join(menuLines, "\n")
}
