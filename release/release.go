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

// createWithDir creates a draft github release with files contained in dir as assets.
func createWithDir(client *github.Client, dir, organisation, project, sha, tag string) error {
	release, err := create(client, organisation, project, sha, tag, false)
	if err != nil {
		return microerror.Mask(err)
	}

	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return microerror.Mask(err)
	}

	for _, f := range filesInfo {
		fd, err := os.Open(filepath.Join(dir, f.Name()))
		if err != nil {
			return microerror.Mask(err)
		}
		defer fd.Close()

		err = uploadAsset(client, organisation, project, release.GetID(), fd)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func create(client *github.Client, organisation, project, sha, tag string, draft bool) (*github.RepositoryRelease, error) {
	release := &github.RepositoryRelease{
		Draft:           &draft,
		Name:            &tag,
		TagName:         &tag,
		TargetCommitish: &sha,
	}

	createdRelease, _, err := client.Repositories.CreateRelease(
		context.TODO(),
		organisation,
		project,
		release,
	)
	if err != nil {
		return nil, microerror.Maskf(err, "could not create release for %s:%s", project, tag)
	}

	log.Printf("created release for %s sha:%s tag:%s", project, sha, tag)

	return createdRelease, nil
}

func uploadAsset(client *github.Client, organisation, project string, id int64, file *os.File) error {
	options := &github.UploadOptions{
		Name: filepath.Base(file.Name()),
	}

	_, _, err := client.Repositories.UploadReleaseAsset(
		context.TODO(),
		organisation,
		project,
		id,
		options,
		file,
	)
	if err != nil {
		return microerror.Maskf(err, "could not upload release asset id:%d file:%s", id, file.Name())
	}

	log.Printf("uploaded release asset id:%d file:%s", id, file.Name())

	return nil
}
