package main

import (
	"github.com/exograd/go-program"
)

func addPipelineCommands() {
	var c *program.Command

	// list-pipelines
	c = p.AddCommand("list-pipelines", "list pipelines",
		cmdListPipelines)

	// abort-pipeline
	c = p.AddCommand("abort-pipeline", "abort a pipeline",
		cmdAbortPipeline)

	c.AddArgument("pipeline-id", "the pipeline to abort")

	// restart-pipeline
	c = p.AddCommand("restart-pipeline", "restart a pipeline",
		cmdRestartPipeline)

	c.AddArgument("pipeline-id", "the pipeline to restart")

	// restart-pipeline-from-failure
	c = p.AddCommand("restart-pipeline-from-failure",
		"restart a pipeline from failed or aborted tasks",
		cmdRestartPipelineFromFailure)

	c.AddArgument("pipeline-id", "the pipeline to restart")
}

func cmdListPipelines(p *program.Program) {
	app.IdentifyCurrentProject()

	pipelines, err := app.Client.FetchPipelines()
	if err != nil {
		p.Fatal("cannot fetch pipelines: %v", err)
	}

	header := []string{
		"id",
		"name",
		"event time",
		"start time",
		"duration",
		"status",
	}

	table := NewTable(header)
	for _, pipeline := range pipelines {
		row := []interface{}{
			pipeline.Id,
			pipeline.Name,
			pipeline.EventTime,
			pipeline.StartTime,
			pipeline.Duration(),
			pipeline.Status,
		}

		table.AddRow(row)
	}

	table.Write()
}

func cmdAbortPipeline(p *program.Program) {
	app.IdentifyCurrentProject()

	Id := p.ArgumentValue("pipeline-id")

	if err := app.Client.AbortPipeline(Id); err != nil {
		p.Fatal("cannot abort pipeline: %v", err)
	}

	p.Info("pipeline aborted")
}

func cmdRestartPipeline(p *program.Program) {
	app.IdentifyCurrentProject()

	Id := p.ArgumentValue("pipeline-id")

	if err := app.Client.RestartPipeline(Id); err != nil {
		p.Fatal("cannot restart pipeline: %v", err)
	}

	p.Info("pipeline restarted")
}

func cmdRestartPipelineFromFailure(p *program.Program) {
	app.IdentifyCurrentProject()

	Id := p.ArgumentValue("pipeline-id")

	if err := app.Client.RestartPipelineFromFailure(Id); err != nil {
		p.Fatal("cannot restart pipeline from failure: %v", err)
	}

	p.Info("pipeline restarted")
}
