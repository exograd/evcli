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

type Client struct {
	Config *Config
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

	// Client
	client := &Client{
		Config: config,
	}

	// Commands
	var cmd func([]string, *Client)

	switch cl.CommandName() {
	case "config":
		cmd = cmdConfig
	case "project":
		cmd = cmdProject
	}

	// Main
	cmd(cl.CommandNameAndArguments(), client)
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
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
