package main

import (
	"github.com/galdor/go-cmdline"
)

func cmdTask(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("list", "list tasks")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "list":
		cmd = cmdTaskList
	}

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdTaskList(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	tasks, err := app.Client.FetchTasks()
	if err != nil {
		die("cannot fetch tasks: %v", err)
	}

	header := []string{
		"id",
		"project",
		"pipeline",
		"instance",
		"status",
		"start time",
		"duration",
	}

	table := NewTable(header)
	for _, task := range tasks {
		row := []interface{}{
			task.Id,
			task.ProjectId,
			task.PipelineId,
			task.InstanceId,
			task.Status,
			task.StartTime,
			task.Duration(),
		}

		table.AddRow(row)
	}

	table.Write()
}
