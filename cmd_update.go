package main

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/galdor/go-program"
	"github.com/google/go-github/v40/github"
)

func addUpdateCommand() {
	var c *program.Command

	// update
	c = p.AddCommand("update", "update the evcli program",
		cmdUpdate)

	c.AddOption("i", "build-id", "build-id", "",
		"force the version to update to")
}

func cmdUpdate(p *program.Program) {
	var buildId *program.BuildId

	if p.IsOptionSet("build-id") {
		s := p.OptionValue("build-id")
		buildId = new(program.BuildId)

		if err := buildId.Parse(s); err != nil {
			p.Fatal("invalid build id %q: %v", s, err)
		}
	}

	if buildId == nil {
		newBuildId, err := findNewBuildId()
		if err != nil {
			p.Fatal("cannot find new evcli build: %v", err)
		}

		if newBuildId == nil {
			p.Info("evcli is up-to-date")
			return
		}

		buildId = newBuildId
	}

	p.Info("updating to evcli %v", buildId)

	// Locate the URI of the evcli binary for the current platform
	os := runtime.GOOS
	arch := runtime.GOARCH

	buildURL, err := findBuildURI(buildId, os, arch)
	if err != nil {
		p.Fatal("cannot find build uri: %v", err)
	}

	p.Debug(1, "build uri: %s", buildURL)

	// TODO
}

func findNewBuildId() (*program.BuildId, error) {
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
		return nil, fmt.Errorf("cannot fetch latest release: %w", err)
	}

	tagName := release.GetTagName()

	var buildId program.BuildId
	if err := buildId.Parse(tagName); err != nil {
		return nil, fmt.Errorf("invalid build id %q: %w", tagName, err)
	}

	return &buildId, nil
}

func findBuildURI(id *program.BuildId, os, arch string) (string, error) {
	httpClient := NewHTTPClient()
	client := github.NewClient(httpClient)

	ctx := context.Background()

	org := "exograd"
	repo := "evcli"
	tagName := id.String()

	p.Debug(1, "fetching release for build %v on os %s and arch %s",
		id, os, arch)

	release, _, err := client.Repositories.GetReleaseByTag(ctx, org, repo,
		tagName)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Response.StatusCode == 404 {
			return "", fmt.Errorf("release not found")
		}

		return "", fmt.Errorf("cannot fetch release: %w", err)
	}

	assetName := "evcli-" + os + "-" + arch

	var asset *github.ReleaseAsset
	for _, asset = range release.Assets {
		if asset.GetName() == assetName {
			break
		}
	}

	return asset.GetBrowserDownloadURL(), nil
}
