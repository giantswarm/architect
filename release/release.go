package release

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"

	"github.com/giantswarm/microerror"
)

type releaseInfo struct {
	AssetsDir    string
	Draft        bool
	Organisation string
	Project      string
	Sha          string
	Tag          string
}

// createWithDir creates a draft github release with files contained in dir as assets.
func createWithDir(client *github.Client, info releaseInfo) error {
	ctx := context.Background()

	release, err := create(ctx, client, info)
	if err != nil {
		return microerror.Mask(err)
	}

	filesInfo, err := ioutil.ReadDir(info.AssetsDir)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, f := range filesInfo {
		fd, err := os.Open(filepath.Join(info.AssetsDir, f.Name()))
		if err != nil {
			return microerror.Mask(err)
		}
		defer fd.Close()

		err = uploadAsset(ctx, client, info, release.GetID(), fd)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func create(ctx context.Context, client *github.Client, info releaseInfo) (*github.RepositoryRelease, error) {
	release := &github.RepositoryRelease{
		Draft:           &info.Draft,
		Name:            &info.Tag,
		TagName:         &info.Tag,
		TargetCommitish: &info.Sha,
	}

	createdRelease, _, err := client.Repositories.CreateRelease(
		ctx,
		info.Organisation,
		info.Project,
		release,
	)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	log.Printf("created release for %s sha:%s tag:%s", info.Project, info.Sha, info.Tag)

	return createdRelease, nil
}

func uploadAsset(ctx context.Context, client *github.Client, info releaseInfo, releaseID int64, file *os.File) error {
	options := &github.UploadOptions{
		Name: filepath.Base(file.Name()),
	}

	_, _, err := client.Repositories.UploadReleaseAsset(
		ctx,
		info.Organisation,
		info.Project,
		releaseID,
		options,
		file,
	)
	if err != nil {
		return microerror.Mask(err)
	}

	log.Printf("uploaded release asset id:%d file:%s", releaseID, file.Name())

	return nil
}
