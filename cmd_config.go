package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/galdor/go-cmdline"
)

func cmdConfig(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("show", "print the configuration and exit")
	cl.AddCommand("get", "extract a value from the configuration and print it")
	cl.AddCommand("set", "set a value in the configuration")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "show":
		cmd = cmdConfigShow
	case "get":
		cmd = cmdConfigGet
	case "set":
		cmd = cmdConfigSet
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

func cmdConfigGet(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("name", "the name of the entry")
	cl.Parse(args)

	name := cl.ArgumentValue("name")

	value, err := app.Config.GetEntry(name)
	if err != nil {
		die("%v", err)
	}

	fmt.Printf("%s\n", value)
}

func cmdConfigSet(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("name", "the name of the entry")
	cl.AddArgument("value", "the value of the entry")
	cl.Parse(args)

	name := cl.ArgumentValue("name")
	value := cl.ArgumentValue("value")

	if err := app.Config.SetEntry(name, value); err != nil {
		die("%v", err)
	}

	if err := app.Config.Write(); err != nil {
		die("%v", err)
	}
}
