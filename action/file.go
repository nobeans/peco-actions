package action

import (
	cmn "github.com/nobeans/peco-actions/common"
	"log"
	"path/filepath"
	"strings"
	"regexp"
	"strconv"
)

type (
	FileActionType struct{}
)

func (FileActionType) prompt() string {
	return "file-actions>"
}

func (FileActionType) menuItems(lines []string) ([]menuItem, error) {
	paths, lineNumOfFirstFile := linesToPaths(lines)

	quotedPaths := strings.TrimSpace(strings.Join(quoteIfRequired(paths), " "))
	log.Printf("quotedPaths: %s", quotedPaths)

	expandedQuotedPaths := strings.TrimSpace(strings.Join(quoteIfRequired(expandAllPaths(paths)), " "))
	log.Printf("expandedQuotedPaths: %s", expandedQuotedPaths)

	items := []menuItem{}

	// Without directories
	// TODO only text!!!!
	if allFiles(paths) {
		log.Printf("All file exist and aren't directory: %#v", paths)
		items = append(items, menuItem{Label: "Edit", Action: editorCommandLine(quotedPaths, lineNumOfFirstFile)})
		items = append(items, menuItem{Label: "Show text", Action: "cat " + quotedPaths})
	}

	// Only if cwd in git repository and tig exists
	if cmn.CommandExists("tig") && cmn.InGitRepository() {
		items = append(items, menuItem{Label: "Tig", Action: "tig " + quotedPaths})
	}

	// Only for single target
	if len(paths) == 1 {
		if cmn.ExistFile(paths[0]) {
			if cmn.IsDirectory(paths[0]) {
				items = append(items, menuItem{Label: "Go into", Action: "cd " + quotedPaths})
			} else {
				parentDir := filepath.Dir(paths[0])
				items = append(items, menuItem{Label: "Go to parent", Action: "cd " + parentDir})
			}
		}
	}

	// Common
	items = append(items, []menuItem{
		{Label: "Open", Action: "open " + quotedPaths},
		{Label: "Show list", Action: "ls -al " + quotedPaths},
		{Label: "Show file type", Action: "file " + quotedPaths},
		{Label: "Copy to Clipboard", Action: "echo -n '" + expandedQuotedPaths + "' | pbcopy"},
	}...)

	return items, nil
}

func linesToPaths(lines []string) ([]string, int) {
	// Support "path:lineNum:lineString" as grep result (lineString is ignored)
	paths := make([]string, len(lines))
	lineNumOfFirstFile := -1
	if len(paths) > 0 && isGrepFormat(lines[0]) {
		for i, line := range lines {
			tokens := strings.SplitN(line, ":", 3)
			paths[i] = tokens[0]
			if lineNumOfFirstFile < 0 {
				lineNumOfFirstFile, _ = strconv.Atoi(tokens[1])
			}
		}
	} else {
		paths = append(paths, lines...)
	}

	// TODO 重複ファイルパスを除去する

	return paths, lineNumOfFirstFile
}

func isGrepFormat(line string) bool {
	return regexp.MustCompile("^[^:]+:[0-9]+:.*$").MatchString(line)
}

func editorCommandLine(path string, lineNum int) string {
	cl := make([]string, 0)

	cmd := cmn.Env("EDITOR", "vi")
	cl = append(cl, cmd)

	// only for vim
	// If EDITOR is "vi(m)" and line is a grep format, use editor options. Yes, I love vim.
	if regexp.MustCompile("(vi|vim)$").MatchString(cmd) {
		// Default line for first file
		if lineNum > 0 {
			cl = append(cl, "+" + strconv.Itoa(lineNum))
		}

		// Highlight in vim
		pattern := cmn.Env("PECO_ACTIONS_EDITOR_PATTERN", "")
		if len(pattern) > 0 {
			cl = append(cl, "+/\"\\c" + pattern + "\"")
		}
	}

	cl = append(cl, path)

	return strings.Join(cl, " ")
}

func expandAllPaths(paths []string) []string {
	return cmn.Map(paths, func(path string) string {
		return cmn.ExpandPath(path)
	})
}

func quoteIfRequired(paths []string) []string {
	return cmn.Map(paths, func(path string) string {
		if strings.Contains(path, " ") {
			// only if the path has spaces
			return "\"" + path + "\""
		}
		return path
	})
}

func allFiles(paths []string) bool {
	return cmn.All(paths, func(path string) bool {
		// TODO want to detect it's a TEXT file or not
		return cmn.ExistFile(path) && !cmn.IsDirectory(path)
	})
}
