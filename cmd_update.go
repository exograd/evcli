package main

import (
	"context"
	"fmt"

	"github.com/galdor/go-program"
	"github.com/google/go-github/v40/github"
)

func addUpdateCommand() {
	// update
	p.AddCommand("update", "update the evcli program",
		cmdUpdate)
}

func cmdUpdate(p *program.Program) {
	newBuildId, err := findNewBuild()
	if err != nil {
		p.Fatal("cannot find new evcli build: %v", err)
	}

	if newBuildId == nil {
		p.Info("evcli is up-to-date")
		return
	}

	p.Info("updating to evcli %v", newBuildId)

	// TODO
}

func findNewBuild() (*program.BuildId, error) {
	var currentBuildId program.BuildId
	currentBuildId.Parse(buildId)

	p.Debug(1, "current build id: %v", currentBuildId)

	lastBuildId, err := lastBuildId()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve last build id: %w", err)
	} else if lastBuildId == nil {
		return nil, nil
	}

	p.Debug(1, "last build id: %v", lastBuildId)

	if lastBuildId.LowerThanOrEqualTo(currentBuildId) {
		return nil, nil
	}

	return lastBuildId, nil
}

func lastBuildId() (*program.BuildId, error) {
	httpClient := NewHTTPClient()
	client := github.NewClient(httpClient)

	ctx := context.Background()

	org := "exograd"
	repo := "evcli"

	release, _, err := client.Repositories.GetLatestRelease(ctx, org, repo)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve latest release: %w", err)
	}

	tagName := release.GetTagName()

	var buildId program.BuildId
	if err := buildId.Parse(tagName); err != nil {
		return nil, fmt.Errorf("invalid build id %q: %w", tagName, err)
	}

	return &buildId, nil
}
