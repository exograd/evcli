package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/exograd/go-program"
)

func addCommandCommands() {
	var c *program.Command

	// list-commands
	c = p.AddCommand("list-commands", "list available commands",
		cmdListCommands)

	// describe-command
	c = p.AddCommand("describe-command", "print information about a command",
		cmdDescribeCommand)

	c.AddArgument("name", "the name of the command")

	// execute-command
	c = p.AddCommand("execute-command", "execute a command",
		cmdExecuteCommand)

	c.AddArgument("name", "the name of the command")
	c.AddTrailingArgument("parameter", "a parameter passed to the command")
}

func cmdListCommands(p *program.Program) {
	app.IdentifyCurrentProject()

	projects, err := app.Client.FetchCommands()
	if err != nil {
		p.Fatal("cannot fetch commands: %v", err)
	}

	header := []string{"name", "description"}
	table := NewTable(header)
	for _, c := range projects {
		row := []interface{}{c.Spec.Name, c.Spec.Description}
		table.AddRow(row)
	}

	table.Write()
}

func cmdDescribeCommand(p *program.Program) {
	app.IdentifyCurrentProject()

	name := p.ArgumentValue("name")

	command, err := app.Client.FetchCommandByName(name)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == "unknown_resource" {
			p.Fatal("unknown command %q", name)
		}

		p.Fatal("cannot fetch command: %v", err)
	}

	commandData := command.Spec.Data.(*CommandData)

	fmt.Printf("%-20s %s\n",
		Colorize(ColorYellow, "Name:"), command.Spec.Name)
	fmt.Printf("%-20s %s\n",
		Colorize(ColorYellow, "Description:"), command.Spec.Description)
	fmt.Printf("%-20s\n",
		Colorize(ColorYellow, "Parameters:"))
	for _, p := range commandData.Parameters {
		fmt.Printf("  - %s: %s\n",
			Colorize(ColorYellow, p.Name), Colorize(ColorGreen, p.Type))
		if p.Description != "" {
			fmt.Printf("    %s\n", p.Description)
		}
		if p.Default != nil {
			defaultString := fmt.Sprintf("%v", p.Default)
			fmt.Printf("    Default: %s\n", Colorize(ColorRed, defaultString))
		}
	}
}

func cmdExecuteCommand(p *program.Program) {
	app.IdentifyCurrentProject()

	name := p.ArgumentValue("name")
	parameterStrings := p.TrailingArgumentValues("parameter")

	command, err := app.Client.FetchCommandByName(name)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == "unknown_resource" {
			p.Fatal("unknown command %q", name)
		}

		p.Fatal("cannot fetch command: %v", err)
	}

	parameters, err := parseParameters(parameterStrings, command)
	if err != nil {
		p.Fatal("%v", err)
	}

	input := CommandExecutionInput{
		Parameters: parameters,
	}

	result, err := app.Client.ExecuteCommand(command.Id, &input)
	if err != nil {
		p.Fatal("cannot execute command: %v", err)
	}

	p.Info("command executed")

	nbPipelines := len(result.PipelineIds)
	if nbPipelines == 1 {
		p.Info("1 pipeline created")
	} else {
		p.Info("%d pipelines created", nbPipelines)
	}

	fmt.Printf("%s\n", result.Id)
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
