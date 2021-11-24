package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	// Determinate the identifier of the build to download and install
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

	// Find the URI of the evcli binary for the current platform
	osName := runtime.GOOS
	archName := runtime.GOARCH

	buildURI, err := findBuildURI(buildId, osName, archName)
	if err != nil {
		p.Fatal("cannot find build uri: %v", err)
	}

	p.Debug(1, "build uri: %s", buildURI)

	// Locate the full path of the current program
	programPath, err := locateProgramPath()
	if err != nil {
		p.Fatal("cannot locate program path: %w", err)
	}

	// Download the new evcli binary to a temporary location
	tmpPath := programPath + ".tmp"

	if err := download(buildURI, tmpPath); err != nil {
		p.Fatal("cannot download build: %w", err)
	}

	// Rename the temporary binary to the installation directory
	p.Info("installing evcli to %s", programPath)

	if err := os.Rename(tmpPath, programPath); err != nil {
		tryDeleteFile(tmpPath)
		p.Fatal("cannot rename %s to %s: %v", tmpPath, programPath, err)
	}

	if err := os.Chmod(programPath, 0755); err != nil {
		p.Fatal("cannot make %s executable: %v", programPath, err)
	}

	p.Info("evcli updated")
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
	httpClient := app.HTTPClient
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

func findBuildURI(id *program.BuildId, osName, archName string) (string, error) {
	client := github.NewClient(app.HTTPClient)

	ctx := context.Background()

	org := "exograd"
	repo := "evcli"
	tagName := id.String()

	p.Debug(1, "fetching release for build %v on os %s and arch %s",
		id, osName, archName)

	release, _, err := client.Repositories.GetReleaseByTag(ctx, org, repo,
		tagName)
	if err != nil {
		var githubErr *github.ErrorResponse
		if errors.As(err, &githubErr) && githubErr.Response.StatusCode == 404 {
			return "", fmt.Errorf("release not found")
		}

		return "", fmt.Errorf("cannot fetch release: %w", err)
	}

	assetName := "evcli-" + osName + "-" + archName

	var asset *github.ReleaseAsset
	for _, asset = range release.Assets {
		if asset.GetName() == assetName {
			break
		}
	}

	return asset.GetBrowserDownloadURL(), nil
}

func download(uri, filePath string) error {
	p.Debug(2, "downloading %s to %s", uri, filePath)

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(filePath, flags, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		tryDeleteFile(filePath)
		return fmt.Errorf("cannot create http request: %w", err)
	}

	res, err := app.HTTPClient.Do(req)
	if err != nil {
		tryDeleteFile(filePath)
		return fmt.Errorf("cannot send http request: %w", err)
	}
	defer res.Body.Close()

	n, err := io.Copy(file, res.Body)
	if err != nil {
		tryDeleteFile(filePath)
		return fmt.Errorf("cannot copy response body to %s: %w",
			filePath, err)
	}

	p.Debug(2, "%d bytes written to %s", n, filePath)

	if err := file.Sync(); err != nil {
		tryDeleteFile(filePath)
		return fmt.Errorf("cannot sync %s: %w", filePath, err)
	}

	return nil
}

func locateProgramPath() (string, error) {
	filePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("cannot find current program path: %w", err)
	}

	resolvedFilePath, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot resolve symlinks: %w", err)
	}

	return resolvedFilePath, nil
}

func tryDeleteFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			p.Error("cannot delete %s: %w", filePath, err)
		}
	}
}
