package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/galdor/go-program"
	"github.com/google/go-github/v40/github"
)

type App struct {
	Config *Config
	Client *Client

	HTTPClient *http.Client

	projectPathOption *string
	projectIdOption   *string
	projectNameOption *string
}

func (a *App) LoadAPIKey() {
	if key := os.Getenv("EVCLI_API_KEY"); key != "" {
		p.Debug(1, "using api key from EVCLI_API_KEY environment variable")
		a.Client.APIKey = key
		return
	}

	if key := a.Config.API.Key; key != "" {
		p.Debug(1, "using api key from configuration")
		a.Client.APIKey = key
		return
	}

	p.Error("missing or empty API key")
	p.Info("\nYou need to provide an API key to interact with Eventline. " +
		"You can either edit the evcli configuration file or use the " +
		"following command:")
	p.Info("\n\tevcli set-config api.key <key>")
	p.Info("\nAlternatively, you can set the EVCLI_API_KEY environment " +
		"variable.")
	os.Exit(1)
}

func (a *App) IdentifyCurrentProject() {
	id, err := a.identifyCurrentProject()
	if err != nil {
		p.Fatal("%v", err)
	}

	p.Debug(1, "using project %s as current project", id)

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
			p.Fatal("cannot fetch project %q: %v", name, err)
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

func (a *App) LookForLastBuild() {
	lastBuildId, err := a.lookForLastBuild()
	if err != nil {
		p.Error("cannot find last build: %v", err)
		return
	}

	if lastBuildId == nil {
		p.Debug(1, "evcli is up-to-date")
		return
	}

	p.Info("evcli %v is now available: run \"evcli update\" to install it")
}

func (a *App) lookForLastBuild() (*program.BuildId, error) {
	p.Debug(1, "looking for the last build")

	currentBuildId := a.currentBuildId()

	lastBuildId, err := a.lastBuildId()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve last build id: %w", err)
	} else if lastBuildId == nil {
		return nil, nil
	}

	if lastBuildId.LowerThanOrEqualTo(currentBuildId) {
		return nil, nil
	}

	return lastBuildId, nil
}

func (a *App) currentBuildId() program.BuildId {
	var id program.BuildId
	id.Parse(buildId)
	return id
}

func (a *App) lastBuildId() (*program.BuildId, error) {
	httpClient := a.HTTPClient
	client := github.NewClient(httpClient)

	ctx := context.Background()

	release, _, err := client.Repositories.GetLatestRelease(ctx,
		"exograd", "evcli")
	if err != nil {
		return nil, fmt.Errorf("cannot fetch latest release: %w", err)
	}

	tagName := release.GetTagName()

	var buildId program.BuildId
	if err := buildId.Parse(tagName); err != nil {
		return nil, fmt.Errorf("invalid build id %q: %w", tagName, err)
	}

	return &buildId, nil
}
