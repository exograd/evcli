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

	lastBuildId, err := lastBuildId()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve last build id: %w", err)
	} else if lastBuildId == nil {
		return nil, nil
	}

	p.Info("current build id: %v", currentBuildId)
	p.Info("last build id: %v", lastBuildId)

	if lastBuildId.LowerThanOrEqualTo(currentBuildId) {
		return nil, nil
	}

	return lastBuildId, nil
}

func lastBuildId() (*program.BuildId, error) {
	client := github.NewClient(nil)

	ctx := context.Background()

	org := "exograd"
	repo := "evcli"

	opts := &github.ListOptions{
		PerPage: 1,
	}

	releases, _, err := client.Repositories.ListReleases(ctx, org, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot list releases: %w", err)
	}

	if len(releases) == 0 {
		return nil, nil
	}

	release := releases[0]
	tagName := release.GetTagName()

	var buildId program.BuildId
	if err := buildId.Parse(tagName); err != nil {
		return nil, fmt.Errorf("invalid build id %q: %w", tagName, err)
	}

	return &buildId, nil
}
