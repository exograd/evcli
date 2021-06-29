package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/galdor/go-cmdline"
	"github.com/qri-io/jsonpointer"
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
	dirPath := cl.OptionValue("path")

	var projectFile ProjectFile
	if err := projectFile.Read(dirPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			projectFile.Name = name

			if err := projectFile.Write(dirPath); err != nil {
				die("cannot write project file in %s: %v", dirPath, err)
			}
		} else {
			die("cannot read project file in %s: %v", dirPath, err)
		}
	}

	if projectFile.Name != name {
		die("directory %s already contains project %s",
			dirPath, projectFile.Name)
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
	cl.AddArgument("name", "the name of the project")
	cl.AddOption("d", "directory", "path",
		"the directory containing project data")
	cl.SetOptionDefault("directory", ".")
	cl.Parse(args)

	name := cl.ArgumentValue("name")
	dirPath := cl.OptionValue("directory")

	var projectFile ProjectFile
	if err := projectFile.Read(dirPath); err != nil {
		die("cannot read project file in %s: %v", dirPath, err)
	}

	if projectFile.Name != name {
		die("directory %s contains project %s", dirPath, projectFile.Name)
	}

	var resourceSet ResourceSet
	if err := resourceSet.Load(dirPath); err != nil {
		die("cannot load resources from %s: %v", dirPath, err)
	}

	project, err := app.Client.FetchProjectByName(name)
	if err != nil {
		die("cannot fetch project: %v", err)
	}

	if err := app.Client.DeployProject(project.Id, &resourceSet); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == "invalid_request_body" {
			invalidRequestBodyErr := apiErr.Data.(InvalidRequestBodyError)
			die("invalid resources:\n%s",
				FormatInvalidRequestBodyError(invalidRequestBodyErr,
					&resourceSet))
		}

		die("cannot deploy project: %v", err)
	}
}

func FormatInvalidRequestBodyError(err InvalidRequestBodyError, resourceSet *ResourceSet) string {
	var buf bytes.Buffer

	for i, jsvError := range err.JSVErrors {
		if i > 0 {
			buf.WriteByte('\n')
		}

		ptr, err := jsonpointer.Parse(jsvError.Pointer)
		if err != nil {
			die("invalid json pointer %q in error response: %v", ptr, err)
		}

		if len(ptr) < 2 || ptr[0] != "specs" {
			die("invalid json pointer %q in error response", ptr)
		}

		document, err := strconv.Atoi(ptr[1])
		if err != nil {
			die("invalid document index %q in json pointer %q", ptr[1], ptr)
		}

		if document < 0 || document >= len(resourceSet.Resources) {
			die("invalid document index %d in json pointer %q", document, ptr)
		}

		resource := resourceSet.Resources[document]
		resourcePtr := ptr[2:]

		var message string
		if len(resourcePtr) == 0 {
			message = jsvError.Reason
		} else {
			message = resourcePtr.String() + ": " + jsvError.Reason
		}

		fmt.Fprintf(&buf, "%s: invalid document %d: %s",
			resource.Path, document, message)
	}

	return buf.String()
}
