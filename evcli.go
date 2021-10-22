package main

import (
	"github.com/galdor/go-program"
)

var (
	p   *program.Program
	app *App

	buildId string

	skipConfirmations bool
	colorOutput       bool
)

func main() {
	// Command line
	p = program.NewProgram("evcli", "client for the eventline service")

	p.AddFlag("y", "yes", "skip all confirmations")
	p.AddFlag("", "no-color", "do not use colors")

	p.AddOption("", "project-path", "path", "",
		"the path of the current project")
	p.AddOption("", "project-id", "id", "",
		"the identifier of the current project")
	p.AddOption("p", "project-name", "name", "",
		"the name of the current project")

	addConfigCommands()
	addProjectCommands()
	addCommandCommands()
	addPipelineCommands()
	addEventCommands()

	p.AddCommand("version", "print the version of evcli and exit", cmdVersion)

	p.ParseCommandLine()

	// Config
	skipConfirmations = p.IsOptionSet("yes")

	config, err := LoadConfig()
	if err != nil {
		p.Fatal("cannot load configuration: %v", err)
	}

	colorOutput = config.Interface.Color && !p.IsOptionSet("no-color")

	// Application
	client, err := NewClient(config)
	if err != nil {
		p.Fatal("cannot create api client: %v", err)
	}

	optionValue := func(name string) *string {
		if !p.IsOptionSet(name) {
			return nil
		}

		value := p.OptionValue(name)
		return &value
	}

	app = &App{
		Config: config,
		Client: client,

		projectPathOption: optionValue("project-path"),
		projectIdOption:   optionValue("project-id"),
		projectNameOption: optionValue("project-name"),
	}

	if name := p.CommandName(); name != "config" {
		app.LoadAPIKey()
	}

	p.Run()
}
