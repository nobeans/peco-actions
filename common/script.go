package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Args []string

func Env(name string, defaultValue string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func CommandExists(name string) bool {
	// If there is the command in PATH, Command.Path() returns its full path.
	return exec.Command(name).Path != name
}

func ExitIfPanic() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", r)
		os.Exit(1)
	}
}

func SetupGlobalLogger(discard bool) {
	// You can also specify DEBUG mode by environment variable DEBUG. It's handy to debug in runtime.
	if Env("DEBUG", "") == "" && discard {
		log.SetOutput(ioutil.Discard)
		return
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
}
