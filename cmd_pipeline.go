package main

import (
	"github.com/galdor/go-cmdline"
)

func cmdPipeline(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("list", "list pipelines")
	cl.AddCommand("abort", "abort a running pipeline")
	cl.AddCommand("restart", "restart a finished pipeline")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "list":
		cmd = cmdPipelineList
	case "abort":
		cmd = cmdPipelineAbort
	case "restart":
		cmd = cmdPipelineRestart
	}

	app.IdentifyCurrentProject()

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdPipelineList(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	pipelines, err := app.Client.FetchPipelines()
	if err != nil {
		die("cannot fetch pipelines: %v", err)
	}

	header := []string{
		"id",
		"name",
		"creation time",
		"start time",
		"duration",
		"status",
	}

	table := NewTable(header)
	for _, pipeline := range pipelines {
		row := []interface{}{
			pipeline.Id,
			pipeline.Name,
			pipeline.CreationTime,
			pipeline.StartTime,
			pipeline.Duration(),
			pipeline.Status,
		}

		table.AddRow(row)
	}

	table.Write()
}

func cmdPipelineAbort(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("pipeline-id", "the pipeline to abort")
	cl.Parse(args)

	Id := cl.ArgumentValue("pipeline-id")

	if err := app.Client.AbortPipeline(Id); err != nil {
		die("cannot abort pipeline: %v", err)
	}

	info("pipeline aborted")
}

func cmdPipelineRestart(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("pipeline-id", "the pipeline to restart")
	cl.Parse(args)

	Id := cl.ArgumentValue("pipeline-id")

	if err := app.Client.RestartPipeline(Id); err != nil {
		die("cannot restart pipeline: %v", err)
	}

	info("pipeline restarted")
}
