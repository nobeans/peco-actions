package action

import (
	cmn "github.com/nobeans/peco-actions/common"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type (
	FileActionType struct{}
)

func (FileActionType) prompt() string {
	return "file-actions>"
}

func (FileActionType) menuItems(rawPaths []string) ([]menuItem, error) {
	quotedPaths := strings.TrimSpace(strings.Join(quoteIfRequired(rawPaths), " "))
	log.Printf("quotedPaths: %s", quotedPaths)

	expandedQuotedPaths := strings.TrimSpace(strings.Join(quoteIfRequired(expandAllPaths(rawPaths)), " "))
	log.Printf("expandedQuotedPaths: %s", expandedQuotedPaths)

	items := []menuItem{}

	// Without directories
	if allFiles(rawPaths) {
		log.Printf("All file exist and aren't directory: %#v", rawPaths)
		items = append(items, []menuItem{
			{Label: "Edit", Action: os.Getenv("EDITOR") + " " + quotedPaths},
			{Label: "Show text", Action: "cat " + quotedPaths},
		}...)
	}

	// Only if cwd in git repository and tig exists
	if cmn.CommandExists("tig") && cmn.InGitRepository() {
		items = append(items, menuItem{Label: "Tig", Action: "tig " + quotedPaths})
	}

	// Only for single target
	if len(rawPaths) == 1 {
		if cmn.ExistFile(rawPaths[0]) {
			if cmn.IsDirectory(rawPaths[0]) {
				items = append(items, menuItem{Label: "Go into", Action: "cd " + quotedPaths})
			} else {
				parentDir := filepath.Dir(rawPaths[0])
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
