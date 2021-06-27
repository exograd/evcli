package main

import (
	"os"

	"github.com/galdor/go-cmdline"
)

func main() {
	// Command line
	cl := cmdline.New()

	cl.AddFlag("v", "verbose", "print debug messages")
	cl.AddFlag("q", "quiet", "do not print status and information messages")

	cl.AddCommand("config", "interact with the evcli configuration")

	cl.Parse(os.Args)

	// Commands
	var cmd func([]string)

	switch cl.CommandName() {
	case "config":
		cmd = cmdConfig
	}

	// Main
	cmd(cl.CommandNameAndArguments())
}
