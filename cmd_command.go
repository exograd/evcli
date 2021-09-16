package main

import (
	"github.com/galdor/go-cmdline"
)

func cmdCommand(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("list", "list available commands")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "list":
		cmd = cmdCommandList
	}

	app.IdentifyCurrentProject()

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdCommandList(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	projects, err := app.Client.FetchCommands()
	if err != nil {
		die("cannot fetch commands: %v", err)
	}

	header := []string{"name"}
	table := NewTable(header)
	for _, c := range projects {
		row := []interface{}{c.Spec.Name}
		table.AddRow(row)
	}

	table.Write()
}
