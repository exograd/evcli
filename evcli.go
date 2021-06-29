package main

import (
	"fmt"
	"os"

	"github.com/galdor/go-cmdline"
)

var (
	verbose           bool
	quiet             bool
	skipConfirmations bool

	colorOutput bool
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
	cl.AddFlag("y", "yes", "skip all confirmations")

	cl.AddCommand("api", "interact with the eventline api")
	cl.AddCommand("config", "interact with the evcli configuration")
	cl.AddCommand("project", "manipulate projects")

	cl.Parse(os.Args)

	// Config
	verbose = cl.IsOptionSet("verbose")
	quiet = cl.IsOptionSet("quiet")
	skipConfirmations = cl.IsOptionSet("yes")

	config, err := LoadConfig()
	if err != nil {
		die("cannot load configuration: %v", err)
	}

	colorOutput = config.Interface.Color

	// Application
	client, err := NewClient(config)
	if err != nil {
		die("cannot create api client: %v", err)
	}

	app := &App{
		Config: config,
		Client: client,
	}

	// Commands
	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "api":
		cmd = cmdAPI
	case "config":
		cmd = cmdConfig
	case "project":
		cmd = cmdProject
	}

	// Main
	cmd(cl.CommandNameAndArguments(), app)
}

func trace(format string, args ...interface{}) {
	if !verbose {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func info(format string, args ...interface{}) {
	if !quiet {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func die(format string, args ...interface{}) {
	msg := fmt.Sprintf("error: "+format, args...)
	fmt.Fprintf(os.Stderr, "%s\n", Colorize(ColorRed, msg))
	os.Exit(1)
}
