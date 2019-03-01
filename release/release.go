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

// ensureWithDir ensure github release created with files in dir as assets.
func (r ReleaseGithubTask) ensureWithDir() error {
	ctx := context.Background()

	release, err := r.getByTag(ctx)
	if err != nil {
		return microerror.Mask(err)
	}

	if release == nil {
		release, err = r.create(ctx)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	filesInfo, err := ioutil.ReadDir(r.AssetsDir)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, f := range filesInfo {
		fd, err := os.Open(filepath.Join(r.AssetsDir, f.Name()))
		if err != nil {
			return microerror.Mask(err)
		}
		defer fd.Close()

		uploaded, err := isFileUploaded(fd, release.Assets)
		if err != nil {
			return microerror.Mask(err)
		}

		if !uploaded {
			err = r.uploadAsset(ctx, release.GetID(), fd)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}

func (r ReleaseGithubTask) create(ctx context.Context) (*github.RepositoryRelease, error) {
	release := &github.RepositoryRelease{
		Draft:           &r.Draft,
		Name:            &r.Tag,
		TagName:         &r.Tag,
		TargetCommitish: &r.Sha,
	}

	createdRelease, _, err := r.Client.Repositories.CreateRelease(
		ctx,
		r.Organisation,
		r.Project,
		release,
	)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	log.Printf("created release for %s sha:%s tag:%s", r.Project, r.Sha, r.Tag)

	return createdRelease, nil
}

func (r ReleaseGithubTask) getByTag(ctx context.Context) (*github.RepositoryRelease, error) {
	release, _, err := r.Client.Repositories.GetReleaseByTag(
		ctx,
		r.Organisation,
		r.Project,
		r.Tag,
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

func (r ReleaseGithubTask) uploadAsset(ctx context.Context, releaseID int64, file *os.File) error {
	options := &github.UploadOptions{
		Name: filepath.Base(file.Name()),
	}

	_, _, err := r.Client.Repositories.UploadReleaseAsset(
		ctx,
		r.Organisation,
		r.Project,
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
