package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

var (
	// get env vars
	accessToken   = os.Getenv("BF3_TOKEN")
	currentOrg    = "bf3fc6c"
	currentRepo   = "cli"
	cloneFromOrg  = os.Getenv("CLONE_FROM_ORG")
	cloneFromRepo = os.Getenv("CLONE_FROM_REPO")

	httpClient *http.Client
	gh         *github.Client
)

func init() {
	if accessToken == "" {
		log.Fatal("Please set your GitHub Personal Access Token as 'BF3_TOKEN' environment variable")
	}
	if cloneFromOrg == "" {
		log.Fatal("Please set the organization name from which you want to clone the release as 'CLONE_FROM_ORG' environment variable")
	}
	if cloneFromRepo == "" {
		log.Fatal("Please set the repo name from which you want to clone the release as 'CLONE_FROM_REPO' environment variable")
	}

	// configure http client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	httpClient = oauth2.NewClient(context.Background(), ts)

	// create github client
	gh = github.NewClient(httpClient)
}

func main() {
	currentRelease, response, err := gh.Repositories.GetLatestRelease(context.Background(), currentOrg, currentRepo)
	if err != nil && response.StatusCode != 404 {
		log.Fatal(err)
	}

	latestRelease, _, err := gh.Repositories.GetLatestRelease(context.Background(), cloneFromOrg, cloneFromRepo)
	if err != nil {
		log.Fatal(err)
	}

	if (len(latestRelease.Assets) < 6) {
		log.Fatal("Release assets have not finished uploading. Try again shortly.")
	}

	releaseExists := currentRelease.GetTagName() == latestRelease.GetTagName()
	if releaseExists {
		fmt.Fprintln(os.Stderr, currentRelease.GetTagName(), " tag already exists in this project")
	}

	if currentRelease == nil || !releaseExists {
		newRelease, err := createRelease(latestRelease)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stderr, "Release", newRelease.GetTagName(), "created")
		err = downloadAssets(latestRelease.Assets)
		if err != nil {
			log.Fatal(err)
		}

		err = uploadAssets(latestRelease.Assets, newRelease.GetID())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createRelease(release *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	newRelease, _, err := gh.Repositories.CreateRelease(context.Background(), currentOrg, currentRepo, release)
	if err != nil {
		log.Fatal(err)
	}

	return newRelease, err
}

func downloadAssets(assets []*github.ReleaseAsset) error {
	for _, asset := range assets {
		rc, _, err := gh.Repositories.DownloadReleaseAsset(context.Background(), cloneFromOrg, cloneFromRepo, asset.GetID(), http.DefaultClient)
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create the file
		out, err := os.Create(filepath.Join("/tmp", asset.GetName()))
		fmt.Fprintln(os.Stderr, "Downloaded", asset.GetName(), "to /tmp")

		if err != nil {
			return err
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, rc)
	}

	return nil
}

func uploadAssets(assets []*github.ReleaseAsset, releaseID int64) error {
	for _, asset := range assets {
		rc, _, err := gh.Repositories.DownloadReleaseAsset(context.Background(), cloneFromOrg, cloneFromRepo, asset.GetID(), http.DefaultClient)
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create the file
		out, err := os.Open(filepath.Join("/tmp", asset.GetName()))

		if err != nil {
			return err
		}
		defer out.Close()

		uploadOpts := &github.UploadOptions{
			Name:      asset.GetName(),
			Label:     asset.GetLabel(),
			MediaType: asset.GetContentType(),
		}
		fmt.Fprintln(os.Stderr, "Uploading "+asset.GetName())
		_, _, err = gh.Repositories.UploadReleaseAsset(context.Background(), currentOrg, currentRepo, releaseID, uploadOpts, out)
		if err != nil {
			return err
		}
	}

	return nil
}
