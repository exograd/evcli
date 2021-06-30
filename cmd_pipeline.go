package main

import (
	"github.com/galdor/go-cmdline"
)

func cmdPipeline(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("list", "list pipelines")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "list":
		cmd = cmdPipelineList
	}

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdPipelineList(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	pipelines, err := app.Client.FetchPipelines()
	if err != nil {
		die("cannot fetch pipelines: %v", err)
	}

	header := []string{"id", "project", "name", "creation time",
		"status", "start time", "end time"}
	table := NewTable(header)
	for _, p := range pipelines {
		row := []interface{}{p.Id, p.ProjectId, p.Name, p.CreationTime,
			p.Status, p.StartTime, p.EndTime}
		table.AddRow(row)
	}

	table.Write()
}
