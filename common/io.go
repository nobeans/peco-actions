package common

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadLines(r io.Reader) ([]string, error) {
	var lines []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		text := s.Text()
		log.Printf("Line: %#v", text)
		lines = append(lines, text)
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return lines, nil
}

func ExistFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDirectory(path string) bool {
	fInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fInfo.IsDir()
}

func ExpandPath(path string) string {
	expanded, err := filepath.EvalSymlinks(path)
	if err != nil {
		// TODO it should propagate the error...
		log.Printf("Failed to filepath.EvalSymlinks: %s (because %s)", path, err.Error())
		return path
	}

	expanded, err = filepath.Abs(expanded)
	if err != nil {
		log.Printf("Failed to filepath.Abs: %s (because %s)", path, err.Error())
		return path
	}

	return expanded
}

func InDir(dir string, target string) bool {
	if !IsDirectory(dir) {
		log.Printf("Not directory: %s", dir)
		return false
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Printf("Failed to filepath.Abs(dir): %s (because %s)", target, err.Error())
		return false
	}

	absTarget, err := filepath.Abs(target)
	if err != nil {
		log.Printf("Failed to filepath.Abs(target): %s (because %s)", target, err.Error())
		return false
	}

	return strings.HasPrefix(absTarget, absDir+"/")
}
