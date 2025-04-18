package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nobeans/peco-actions/action"
	"github.com/nobeans/peco-actions/common"
)

const (
	USAGE = `usage: peco-actions [option]
options:
  -h,--help             show this usage
  -v,--version          display the version
  -D,--debug            show debug logs
  --file                actions for file path(s)
  --process             actions for a process id
  --env                 actions for a environment variable
  --server              actions for a host/IP-address
  --git                 actions for a commit id
  --docker-container    actions for a docker container
  --docker-image        actions for a docker image
  --ollama              actions for a ollama model
  --generic             actions for generic only using adhoc menu`
)

var (
	Version = "1.9.0"
)

type (
	Options struct {
		help            bool
		version         bool
		debug           bool
		file            bool
		process         bool
		env             bool
		server          bool
		git             bool
		dockerContainer bool
		dockerImage     bool
		ollama          bool
		generic         bool
	}
)

func newOptions() *Options {
	return &Options{
		help:            false,
		version:         false,
		debug:           false,
		file:            false,
		process:         false,
		env:             false,
		server:          false,
		git:             false,
		dockerContainer: false,
		dockerImage:     false,
		ollama:          false,
		generic:         false,
	}
}

func parseOptions(args common.Args) *Options {
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
		case "--env":
			opts.env = true
		case "--server":
			opts.server = true
		case "--git":
			opts.git = true
		case "--docker-container":
			opts.dockerContainer = true
		case "--docker-image":
			opts.dockerImage = true
		case "--ollama":
			opts.ollama = true
		case "--generic":
			opts.generic = true
		default:
			panic(fmt.Sprintf("unrecognized option %s", arg))
		}
	}
	return opts
}

func actionType(opts *Options) action.Type {
	switch {
	case opts.file:
		return new(action.FileActionType)
	case opts.process:
		return new(action.ProcessActionType)
	case opts.env:
		return new(action.EnvActionType)
	case opts.server:
		return new(action.ServerActionType)
	case opts.git:
		return new(action.GitActionType)
	case opts.dockerContainer:
		return new(action.DockerContainerActionType)
	case opts.dockerImage:
		return new(action.DockerImageActionType)
	case opts.ollama:
		return new(action.OllamaActionType)
	case opts.generic:
		return new(action.GenericActionType)
	default:
		panic("no action type\n" + USAGE)
	}
}

func main() {
	defer common.ExitIfPanic()

	// Parsing options
	opts := parseOptions(os.Args)

	// Setting log level (all logging must be after this line)
	common.SetupGlobalLogger(!opts.debug)

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

	act, err := action.Resolve(actionType(opts), os.Stdin)
	if err != nil {
		panic(fmt.Sprintf("%s", err.Error()))
	}
	fmt.Printf("%s", act)

	os.Exit(0)
}
