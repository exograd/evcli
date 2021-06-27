package main

import (
	"encoding/json"
	"os"

	"github.com/galdor/go-cmdline"
)

func cmdConfig(args []string, client *Client) {
	cl := cmdline.New()
	cl.AddCommand("show", "print the configuration and exit")
	cl.Parse(args)

	var cmd func([]string, *Client)

	switch cl.CommandName() {
	case "show":
		cmd = cmdConfigShow
	}

	cmd(cl.CommandNameAndArguments(), client)
}

func cmdConfigShow(args []string, client *Client) {
	cl := cmdline.New()
	cl.Parse(args)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(client.Config); err != nil {
		die("cannot encode configuration: %v", err)
	}
}
