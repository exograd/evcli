package main

import (
	"fmt"
	"os"

	"github.com/galdor/go-cmdline"
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

	cl.Parse(os.Args)

	// Config
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
	}

	// Main
	cmd(cl.CommandNameAndArguments(), client)
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
