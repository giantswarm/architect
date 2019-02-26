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

// ensureWithDir ensure github release created with files in dir as assets.
func ensureWithDir(client *github.Client, info releaseInfo) error {
	ctx := context.Background()

	release, err := getByTag(ctx, client, info)
	if err != nil {
		return microerror.Mask(err)
	}

	if release == nil {
		release, err = create(ctx, client, info)
		if err != nil {
			return microerror.Mask(err)
		}
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

		uploaded, err := isFileUploaded(fd, release.Assets)
		if err != nil {
			return microerror.Mask(err)
		}

		if !uploaded {
			err = uploadAsset(ctx, client, info, release.GetID(), fd)
			if err != nil {
				return microerror.Mask(err)
			}
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

func getByTag(ctx context.Context, client *github.Client, info releaseInfo) (*github.RepositoryRelease, error) {
	release, _, err := client.Repositories.GetReleaseByTag(
		ctx,
		info.Organisation,
		info.Project,
		info.Tag,
	)
	if IsNotFoundError(err) {
		// fallthrough
	} else if err != nil {
		return nil, microerror.Mask(err)
	}

	return release, nil
}

// isFileUploaded check if the file as been uploaded as a github release asset.
//
// It compares name and size, and check for uploaded status.
func isFileUploaded(fd *os.File, assets []github.ReleaseAsset) (bool, error) {
	for _, asset := range assets {
		name := filepath.Base(fd.Name())
		if asset.Name != nil && *asset.Name == name {
			info, err := fd.Stat()
			if err != nil {
				return false, microerror.Mask(err)
			}

			if asset.Size != nil && int64(*asset.Size) == info.Size() {

				if asset.State != nil && *asset.State == "uploaded" {
					return true, nil
				}
			}
		}
	}

	return false, nil
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
