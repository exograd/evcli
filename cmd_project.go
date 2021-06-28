package main

import (
	"net/url"

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

	var page ProjectPage

	query := url.Values{}
	query.Add("size", "100")
	uri := &url.URL{Path: "/v0/projects", RawQuery: query.Encode()}

	err := app.Client.SendRequest("GET", uri, nil, &page)
	if err != nil {
		die("cannot fetch projects: %v", err)
	}

	header := []string{"id", "name", "description"}
	table := NewTable(header)
	for _, project := range page.Elements {
		row := []interface{}{project.Id, project.Name, project.Description}
		table.AddRow(row)
	}

	table.Write()
}

func cmdProjectCreate(args []string, app *App) {
	cl := cmdline.New()
	cl.Parse(args)

	// TODO
	die("unimplemented")
}

func cmdProjectDelete(args []string, app *App) {
	cl := cmdline.New()
	cl.AddArgument("name", "the name of the project to delete")
	cl.Parse(args)

	name := cl.ArgumentValue("name")

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
