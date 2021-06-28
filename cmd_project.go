package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/galdor/go-cmdline"
)

func cmdProject(args []string, app *App) {
	cl := cmdline.New()
	cl.AddCommand("list", "list projects")
	cl.AddCommand("create", "create a project")
	cl.AddCommand("delete", "delete a project")
	cl.AddCommand("deploy", "deploy a project")
	cl.Parse(args)

	var cmd func([]string, *App)

	switch cl.CommandName() {
	case "list":
		cmd = cmdProjectList
	case "create":
		cmd = cmdProjectCreate
	case "delete":
		cmd = cmdProjectDelete
	case "deploy":
		cmd = cmdProjectDeploy
	}

	cmd(cl.CommandNameAndArguments(), app)
}

func cmdProjectList(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	projects, err := app.Client.FetchProjects()
	if err != nil {
		die("cannot fetch projects: %v", err)
	}

	header := []string{"id", "name", "description"}
	table := NewTable(header)
	for _, project := range projects {
		row := []interface{}{project.Id, project.Name, project.Description}
		table.AddRow(row)
	}

	table.Write()
}

func cmdProjectCreate(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("name", "the name of the project")
	cl.AddArgument("path", "the directory which will contain project data")
	cl.AddOption("d", "description", "description",
		"a description of the project")
	cl.Parse(args)

	name := cl.ArgumentValue("name")
	path := cl.ArgumentValue("path")

	projectFile, err := LoadProjectFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			trace("creating project file at %s", path)

			projectFile = &ProjectFile{
				Name: name,
			}

			if err := projectFile.WriteFile(path); err != nil {
				die("cannot create project file at %s: %v", path, err)
			}
		} else {
			die("cannot load project file from %s: %v", path, err)
		}
	}

	if projectFile.Name != name {
		die("directory %s already contains project %s",
			path, projectFile.Name)
	}

	project := &Project{
		Name:        name,
		Description: cl.OptionValue("description"),
	}

	if err := app.Client.CreateProject(project); err != nil {
		die("cannot create project: %v", err)
	}

	info("project %s created", project.Id)
}

func cmdProjectDelete(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("name", "the name of the project")
	cl.Parse(args)

	name := cl.ArgumentValue("name")

	prompt := fmt.Sprintf("Do you want to delete project %q ? All resources "+
		"associated with it will be deleted as well.", name)
	if Confirm(prompt) == false {
		info("deletion aborted")
		return
	}

	project, err := app.Client.FetchProjectByName(name)
	if err != nil {
		die("cannot fetch project: %v", err)
	}

	if err := app.Client.DeleteProject(project.Id); err != nil {
		die("cannot delete project: %v", err)
	}
}

func cmdProjectDeploy(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	// TODO
	die("unimplemented")
}
