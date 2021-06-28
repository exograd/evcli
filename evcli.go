package main

import (
	"fmt"
	"os"

	"github.com/galdor/go-cmdline"
)

var (
	verbose bool
	quiet   bool
)

type App struct {
	Config *Config
	Client *Client
}

func main() {
	// Command line
	cl := cmdline.New()

	cl.AddFlag("v", "verbose", "print debug messages")
	cl.AddFlag("q", "quiet", "do not print status and information messages")

	cl.AddCommand("config", "interact with the evcli configuration")
	cl.AddCommand("project", "manipulate projects")

	cl.Parse(os.Args)

	// Config
	verbose = cl.IsOptionSet("verbose")
	quiet = cl.IsOptionSet("quiet")

	config, err := LoadConfig()
	if err != nil {
		die("cannot load configuration: %v", err)
	}

	// Application
	app := &App{
		Config: config,
		Client: NewClient(config),
	}

	// Commands
	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "config":
		cmd = cmdConfig
	case "project":
		cmd = cmdProject
	}

	// Main
	cmd(cl.CommandNameAndArguments(), app)
}

func trace(format string, args ...interface{}) {
	if verbose == false {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func info(format string, args ...interface{}) {
	if quiet == true {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
