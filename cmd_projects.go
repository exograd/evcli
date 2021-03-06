package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/exograd/go-program"
	"github.com/qri-io/jsonpointer"
)

func addProjectCommands() {
	var c *program.Command

	// list-projects
	c = p.AddCommand("list-projects", "list projects",
		cmdListProjects)

	// initialize-project
	c = p.AddCommand("initialize-project",
		"initialize a directory for an existing project",
		cmdInitializeProject)

	c.AddArgument("name", "the name of the project")
	c.AddArgument("path", "the directory which will contain project data")

	// create-project
	c = p.AddCommand("create-project", "create a new project",
		cmdCreateProject)

	c.AddArgument("name", "the name of the project")
	c.AddArgument("path", "the directory which will contain project data")

	// delete-project
	c = p.AddCommand("delete-project", "delete a project",
		cmdDeleteProject)

	c.AddArgument("name", "the name of the project")

	// deploy-project
	c = p.AddCommand("deploy-project", "deploy resources for a project",
		cmdDeployProject)

	c.AddOption("d", "directory", "path", ".",
		"the directory containing project data")
	c.AddFlag("n", "dry-run", "validate resources but do not deploy them")

	// list-project-files
	c = p.AddCommand("list-project-files",
		"list resource files in a project directory", cmdListProjectFiles)

	c.AddOption("d", "directory", "path", ".",
		"the directory containing project data")
}

func cmdListProjects(p *program.Program) {
	projects, err := app.Client.FetchProjects()
	if err != nil {
		p.Fatal("cannot fetch projects: %v", err)
	}

	header := []string{"id", "name"}
	table := NewTable(header)
	for _, p := range projects {
		row := []interface{}{p.Id, p.Name}
		table.AddRow(row)
	}

	table.Write()
}

func cmdInitializeProject(p *program.Program) {
	name := p.ArgumentValue("name")
	dirPath := p.ArgumentValue("path")

	var projectFile ProjectFile
	if err := projectFile.Read(dirPath); err == nil {
		p.Fatal("directory %s already contains a project file for "+
			"project %q", dirPath, projectFile.Name)
	}

	project, err := app.Client.FetchProjectByName(name)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == "unknown_project" {
			p.Error("unknown project")
			p.Info("\nYou can use the create-project command to create a " +
				"new project and initialize its directory.")
			os.Exit(1)
		} else {
			p.Fatal("cannot fetch project %q: %v", name, err)
		}
	}

	projectFile.Name = name
	projectFile.Id = project.Id

	if err := projectFile.Write(dirPath); err != nil {
		p.Fatal("cannot write project file in %s: %v", dirPath, err)
	}

	p.Info("project %s initialized", project.Name)
}

func cmdCreateProject(p *program.Program) {
	name := p.ArgumentValue("name")
	dirPath := p.ArgumentValue("path")

	var projectFile ProjectFile
	if err := projectFile.Read(dirPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			projectFile.Name = name

			if err := projectFile.Write(dirPath); err != nil {
				p.Fatal("cannot write project file in %s: %v", dirPath, err)
			}
		} else {
			p.Fatal("cannot read project file in %s: %v", dirPath, err)
		}
	}

	if projectFile.Name != name {
		p.Fatal("directory %s already contains project %s",
			dirPath, projectFile.Name)
	}

	project := &Project{
		Name: name,
	}

	if err := app.Client.CreateProject(project); err != nil {
		p.Fatal("cannot create project: %v", err)
	}

	projectFile.Id = project.Id
	if err := projectFile.Write(dirPath); err != nil {
		p.Fatal("cannot write project file in %s: %v", dirPath, err)
	}

	p.Info("project %q created", project.Name)
}

func cmdDeleteProject(p *program.Program) {
	name := p.ArgumentValue("name")

	prompt := fmt.Sprintf("Do you want to delete project %q??? All resources "+
		"associated with it will be deleted as well.", name)
	if Confirm(prompt) == false {
		p.Info("deletion aborted")
		return
	}

	project, err := app.Client.FetchProjectByName(name)
	if err != nil {
		p.Fatal("cannot fetch project: %v", err)
	}

	if err := app.Client.DeleteProject(project.Id); err != nil {
		p.Fatal("cannot delete project: %v", err)
	}
}

func cmdDeployProject(p *program.Program) {
	dirPath := p.OptionValue("directory")
	dryRun := p.IsOptionSet("dry-run")

	var projectFile ProjectFile
	if err := projectFile.Read(dirPath); err != nil {
		p.Fatal("cannot read project file in %s: %v", dirPath, err)
	}

	app.Client.ProjectId = projectFile.Id

	var ignoreSet IgnoreSet
	if err := ignoreSet.LoadDirectoryIfExists(dirPath); err != nil {
		p.Fatal("cannot load ignore file: %v", err)
	}

	var resourceSet ResourceSet
	if err := resourceSet.Load(dirPath, &ignoreSet); err != nil {
		p.Fatal("cannot load resources: %v", err)
	}

	if len(resourceSet.Resources) == 0 {
		p.Fatal("no resource available")
	}

	err := app.Client.DeployProject(projectFile.Id, &resourceSet, dryRun)
	if err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.Code == "invalid_request_body" {
			invalidRequestBodyErr := apiErr.Data.(InvalidRequestBodyError)
			p.Fatal("invalid resources:\n%s",
				formatInvalidRequestBodyError(invalidRequestBodyErr,
					&resourceSet))
		}

		if dryRun {
			p.Fatal("invalid project: %v", err)
		} else {
			p.Fatal("cannot deploy project: %v", err)
		}
	}

	if dryRun {
		p.Info("project validated successfully")
	} else {
		p.Info("project deployed successfully")
	}
}

func formatInvalidRequestBodyError(err InvalidRequestBodyError, resourceSet *ResourceSet) string {
	var buf bytes.Buffer

	for i, jsvError := range err.JSVErrors {
		if i > 0 {
			buf.WriteByte('\n')
		}

		ptr, err := jsonpointer.Parse(jsvError.Pointer)
		if err != nil {
			p.Fatal("invalid json pointer %q in error response: %v", ptr, err)
		}

		if len(ptr) < 2 || ptr[0] != "specs" {
			p.Fatal("invalid json pointer %q in error response", ptr)
		}

		document, err := strconv.Atoi(ptr[1])
		if err != nil {
			p.Fatal("invalid document index %q in json pointer %q", ptr[1], ptr)
		}

		if document < 0 || document >= len(resourceSet.Resources) {
			p.Fatal("invalid document index %d in json pointer %q", document, ptr)
		}

		resource := resourceSet.Resources[document]
		resourcePtr := ptr[2:]

		documentId := fmt.Sprintf("document %d", document)

		var message string

		if len(resourcePtr) == 0 {
			message = jsvError.Reason
		} else {
			message = resourcePtr.String() + ": " + jsvError.Reason

			typeName, name := resource.TypeAndName()
			if typeName != "" && name != "" {
				documentId = fmt.Sprintf("%s %q", typeName, name)
			}
		}

		fmt.Fprintf(&buf, "%s: invalid %s: %s",
			resource.Path, documentId, message)
	}

	return buf.String()
}

func cmdListProjectFiles(p *program.Program) {
	dirPath := p.OptionValue("directory")

	var ignoreSet IgnoreSet
	if err := ignoreSet.LoadDirectoryIfExists(dirPath); err != nil {
		p.Fatal("cannot load ignore file: %v", err)
	}

	filePaths, err := FindResourceFiles(dirPath, &ignoreSet)
	if err != nil {
		p.Fatal("cannot find files: %v", err)
	}

	sort.Strings(filePaths)

	dirPathLen := len(dirPath)

	for _, filePath := range filePaths {
		relFilePath := filePath[dirPathLen+1:]
		fmt.Printf("%s\n", relFilePath)
	}
}
