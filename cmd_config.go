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
	cl.AddFlag("e", "entries",
		"show a list of entries instead of the entire configuration")
	cl.Parse(args)

	if cl.IsOptionSet("entries") {
		var names []string
		for _, e := range ConfigEntries {
			names = append(names, e.Name)
		}

		table := NewTable([]string{"name", "value"})
		for _, name := range names {
			value, err := app.Config.GetEntry(name)
			if err != nil {
				warn("cannot read entry %q: %v", name, err)
				continue
			}

			table.AddRow([]interface{}{name, value})
		}

		table.Write()
	} else {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")

		if err := encoder.Encode(app.Config); err != nil {
			die("cannot encode configuration: %v", err)
		}
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
