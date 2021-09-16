package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/galdor/go-cmdline"
)

func cmdCommand(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("list", "list available commands")
	cl.AddCommand("execute", "execute a command")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "list":
		cmd = cmdCommandList
	case "execute":
		cmd = cmdCommandExecute
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

	header := []string{"name", "description"}
	table := NewTable(header)
	for _, c := range projects {
		row := []interface{}{c.Spec.Name, c.Spec.Description}
		table.AddRow(row)
	}

	table.Write()
}

func cmdCommandExecute(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("name", "the name of the command")
	cl.AddTrailingArguments("parameter", "a parameter passed to the command")
	cl.Parse(args)

	name := cl.ArgumentValue("name")
	parameterStrings := cl.TrailingArgumentsValues("parameter")

	command, err := app.Client.FetchCommandByName(name)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == "unknown_resource" {
			die("unknown command %q", name)
		}

		die("cannot fetch command: %v", err)
	}

	parameters, err := parseParameters(parameterStrings, command)
	if err != nil {
		die("%v", err)
	}

	execution := CommandExecution{
		Parameters: parameters,
	}

	result, err := app.Client.ExecuteCommand(command.Id, &execution)
	if err != nil {
		die("cannot execute command: %v", err)
	}

	info("command executed")

	nbPipelines := len(result.PipelineIds)
	if nbPipelines == 1 {
		info("1 pipeline created")
	} else {
		info("%d pipelines created", nbPipelines)
	}

	table := NewTable([]string{"Pipeline"})
	for _, pipelineId := range result.PipelineIds {
		table.AddRow([]interface{}{pipelineId})
	}
	table.Write()
}

func parseParameters(parameterStrings []string, command *Resource) (map[string]interface{}, error) {
	commandData := command.Spec.Data.(*CommandData)

	parameters := make(map[string]interface{})

	for _, s := range parameterStrings {
		name, value, err := parseParameter(s, command)
		if err != nil {
			return nil, err
		}

		parameters[name] = value
	}

	for _, p := range commandData.Parameters {
		if p.Default != nil {
			continue
		}

		if _, found := parameters[p.Name]; !found {
			return nil, fmt.Errorf("missing parameter %q", p.Name)
		}
	}

	return parameters, nil
}

func parseParameter(parameterString string, command *Resource) (string, interface{}, error) {
	commandData := command.Spec.Data.(*CommandData)

	parts := strings.SplitN(parameterString, "=", 2)
	if len(parts) != 2 {
		return "", nil,
			fmt.Errorf("invalid parameter format %q", parameterString)
	}

	name := parts[0]
	valueString := parts[1]

	var p *Parameter
	for _, cmdp := range commandData.Parameters {
		if cmdp.Name == name {
			p = cmdp
			break
		}
	}

	if p == nil {
		return "", nil, fmt.Errorf("unknown parameter %q", name)
	}

	var value interface{}

	switch p.Type {
	case "number":
		var i int64
		i, err := strconv.ParseInt(valueString, 10, 64)
		if err == nil {
			value = i
		} else {
			f, err := strconv.ParseFloat(valueString, 64)
			if err == nil {
				value = f
			} else {
				return "", nil,
					fmt.Errorf("invalid number value %q", valueString)
			}
		}

	case "string":
		value = valueString

	case "boolean":
		valueString = strings.ToLower(valueString)

		switch valueString {
		case "true":
			value = true
		case "false":
			value = false
		default:
			return "", nil, fmt.Errorf("invalid boolean value %q", valueString)
		}
	}

	return name, value, nil
}
