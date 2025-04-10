package action

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/nobeans/peco-actions/common"
)

type FileActionType struct{}

func (FileActionType) prompt() string {
	return "file-actions>"
}

func (FileActionType) menuItems(lines []string) ([]menuItem, error) {
	paths, lineNumOfFirstFile := linesToPaths(lines)

	quotedPaths := strings.TrimSpace(strings.Join(quoteIfRequired(paths), " "))
	log.Printf("quotedPaths: %s", quotedPaths)

	expandedQuotedPaths := strings.TrimSpace(strings.Join(quoteIfRequired(expandAllPaths(paths)), " "))
	log.Printf("expandedQuotedPaths: %s", expandedQuotedPaths)

	var items []menuItem

	// If all are text files
	if allTextFiles(paths) {
		log.Printf("All file exist and aren't directory: %#v", paths)
		items = append(items, menuItem{Label: "Edit", Action: editorCommandLine(quotedPaths, lineNumOfFirstFile)})
		items = append(items, menuItem{Label: "Show text", Action: "cat " + quotedPaths})
	}

	// Only for single target
	if len(paths) == 1 {
		if common.ExistFile(paths[0]) {
			if common.IsDirectory(paths[0]) {
				items = append(items, menuItem{Label: "Go into", Action: "cd " + quotedPaths})
			} else {
				parentDir := filepath.Dir(paths[0])
				items = append(items, menuItem{Label: "Go to parent", Action: "cd " + parentDir})
			}
		}
	}

	// If all are in git repository and tig exists
	if common.CommandExists("tig") && common.CwdInGitRepository() && allInCwd(paths) {
		items = append(items, menuItem{Label: "Tig", Action: "tig " + quotedPaths})
	}

	if common.CommandExists("idea") {
		items = append(items, menuItem{Label: "Open (IntelliJ IDEA)", Action: "idea " + quotedPaths})
	}
	if common.CommandExists("code") {
		items = append(items, menuItem{Label: "Open (Visual Studio Code)", Action: "code " + quotedPaths})
	}
	if common.CommandExists("open") {
		items = append(items, menuItem{Label: "Open (default)", Action: "open " + quotedPaths})
	}
	items = append(items, []menuItem{
		{Label: "Show list", Action: "ls -al " + quotedPaths},
		{Label: "Show file type", Action: "file " + quotedPaths},
	}...)
	if common.CommandExists("pbcopy") {
		items = append(items, menuItem{Label: "Copy to Clipboard", Action: "echo -n '" + expandedQuotedPaths + "' | pbcopy"})
	}
	items = append(items, RenderAdhocMenuItems(quotedPaths)...)
	return items, nil
}

func linesToPaths(lines []string) ([]string, int) {
	// Support "path:lineNum:lineString" as grep result (lineString is ignored)

	var paths []string
	lineNumOfFirstFile := -1

	if len(lines) > 0 && isGrepFormat(lines[0]) {
		for _, line := range lines {
			tokens := strings.SplitN(line, ":", 3)
			path := tokens[0]
			lineNum, _ := strconv.Atoi(tokens[1])

			// Remove duplication
			if !common.Include(paths, path) {
				paths = append(paths, path)

				if lineNumOfFirstFile < 0 {
					lineNumOfFirstFile = lineNum
				}
			}
		}
	} else {
		for _, path := range lines {
			// Remove duplication
			if !common.Include(paths, path) {
				paths = append(paths, path)
			}
		}
	}

	return paths, lineNumOfFirstFile
}

func isGrepFormat(line string) bool {
	return regexp.MustCompile("^[^:]+:[0-9]+:.*$").MatchString(line)
}

func editorCommandLine(path string, lineNum int) string {
	var cl []string

	cmd := common.Env("EDITOR", "vi")
	cl = append(cl, cmd)

	// only for vim
	// If EDITOR is "vi(m)" and line is a grep format, use editor options. Yes, I love vim.
	if regexp.MustCompile("(vi|vim)$").MatchString(cmd) {
		// Default line for first file
		if lineNum > 0 {
			cl = append(cl, "+"+strconv.Itoa(lineNum))
		}

		// Highlight in vim
		pattern := common.Env("PECO_ACTIONS__EDITOR_PATTERN", "")
		if len(pattern) > 0 {
			cl = append(cl, "+/\"\\c"+pattern+"\"")
		}
	}

	cl = append(cl, path)

	return strings.Join(cl, " ")
}

func expandAllPaths(paths []string) []string {
	return common.Map(paths, func(path string) string {
		return common.ExpandPath(path)
	})
}

func quoteIfRequired(paths []string) []string {
	return common.Map(paths, func(path string) string {
		if strings.Contains(path, " ") {
			// only if the path has spaces
			return "\"" + path + "\""
		}
		return path
	})
}

func allTextFiles(paths []string) bool {
	return common.All(paths, func(path string) bool {
		if common.CommandExists("file") {
			out, err := exec.Command("file", path).Output()
			if err != nil {
				return false
			}
			return regexp.MustCompile("\\b(text|JSON data)\\b").MatchString(strings.TrimSpace(fmt.Sprintf("%s", out)))
		} else {
			return common.ExistFile(path) && !common.IsDirectory(path)
		}
	})
}

func allInCwd(paths []string) bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	return common.All(paths, func(path string) bool {
		inDir := common.InDir(cwd, path)
		log.Printf("allInCwd: cwd=%s, path=%s, inDir=%v", cwd, path, inDir)
		return inDir
	})
}
