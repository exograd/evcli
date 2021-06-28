package main

import (
	"encoding/json"
	"os"

	"github.com/galdor/go-cmdline"
)

func cmdConfig(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("show", "print the configuration and exit")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "show":
		cmd = cmdConfigShow
	}

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdConfigShow(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(app.Config); err != nil {
		die("cannot encode configuration: %v", err)
	}
}
