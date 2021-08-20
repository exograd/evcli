package main

import (
	"errors"
	"fmt"
	"os"
)

type App struct {
	Config *Config
	Client *Client

	projectPathOption *string
	projectIdOption   *string
	projectNameOption *string
}

func (a *App) LoadAPIKey() {
	if a.Config.API.Key == "" {
		err("missing or empty API key")
		info("\nYou need to provide an API key to interact with Eventline. " +
			"You can either edit the evcli configuration file or use the " +
			"following command:")
		info("\n\tevcli config set api.key <key>")
		info("\nAlternatively, you can set the EVCLI_API_KEY environment " +
			"variable.")
		os.Exit(1)
	}

	a.Client.APIKey = a.Config.API.Key
}

func (a *App) IdentifyCurrentProject() {
	id, err := a.identifyCurrentProject()
	if err != nil {
		die("%v", err)
	}

	trace("using project %s as current project", id)

	a.Client.ProjectId = id
}

func (a *App) identifyCurrentProject() (string, error) {
	id, err := a.loadProjectDirectory(".")
	if err != nil {
		return "", err
	} else if id != "" {
		return id, nil
	}

	if a.projectIdOption != nil {
		return *a.projectIdOption, nil
	}

	if a.projectPathOption != nil {
		id, err := a.loadProjectDirectory(*a.projectPathOption)
		if err != nil {
			return "", err
		} else if id != "" {
			return id, nil
		}
	}

	if a.projectNameOption != nil {
		name := *a.projectNameOption

		project, err := a.Client.FetchProjectByName(name)
		if err != nil {
			die("cannot fetch project %q: %v", name, err)
		}

		return project.Id, nil
	}

	return "", fmt.Errorf("cannot identify the current project")
}

func (a *App) loadProjectDirectory(dirPath string) (string, error) {
	var projectFile ProjectFile
	if err := projectFile.Read(dirPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	return projectFile.Id, nil
}
