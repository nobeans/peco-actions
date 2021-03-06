package main

import (
	"fmt"
	act "github.com/nobeans/peco-actions/action"
	cmn "github.com/nobeans/peco-actions/common"
	"log"
	"os"
)

const (
	USAGE = `usage: peco-actions [option]
options:
  -h,--help     show this usage
  -v,--version  display the version
  -D,--debug    display the version
  --file        actions for file path(s)
  --process     actions for a process id
  --server      actions for a host/IP-address
  --git         actions for a commit id
  --generic     actions for generic only using addhoc menu`
)

var (
	Version = "X.X.X-SNAPSHOT" // replaced by ldflags
)

type (
	Options struct {
		help    bool
		version bool
		debug   bool
		file    bool
		process bool
		server  bool
		git     bool
		generic bool
	}
)

func newOptions() *Options {
	return &Options{
		help:    false,
		version: false,
		debug:   false,
		file:    false,
		process: false,
		server:  false,
		git:     false,
		generic: false,
	}
}

func parseOptions(args cmn.Args) *Options {
	opts := newOptions()
	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			opts.help = true
		case "-v", "--version":
			opts.version = true
		case "-D", "--debug":
			opts.debug = true
		case "--file":
			opts.file = true
		case "--process":
			opts.process = true
		case "--server":
			opts.server = true
		case "--git":
			opts.git = true
		case "--generic":
			opts.generic = true
		default:
			panic(fmt.Sprintf("unrecognized option %s", arg))
		}
	}
	return opts
}

func actionType(opts *Options) act.ActionType {
	switch {
	case opts.file:
		return new(act.FileActionType)
	case opts.process:
		return new(act.ProcessActionType)
	case opts.server:
		return new(act.ServerActionType)
	case opts.git:
		return new(act.GitActionType)
	case opts.generic:
		return new(act.GenericActionType)
	default:
		panic("no action type\n" + USAGE)
	}
}

func main() {
	defer cmn.ExitIfPanic()

	// Parsing options
	opts := parseOptions(os.Args)

	// Setting log level (all logging must be after this line)
	cmn.SetupGlobalLogger(!opts.debug)

	log.Printf("Original arguments: %#v", os.Args)
	log.Printf("Parsed options: %#v", opts) // must be after affecting -debug option

	switch {
	case opts.help:
		fmt.Println(USAGE)
		os.Exit(0)
	case opts.version:
		fmt.Println("peco-actions version " + Version)
		os.Exit(0)
	}

	action, err := act.ResolveAction(actionType(opts), os.Stdin)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	fmt.Printf("%s", action)

	os.Exit(0)
}
