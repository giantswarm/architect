package release

import (
	"fmt"

	"github.com/google/go-github/github"
)

const (
	ReleaseGithubTaskName   = "release-github"
	ReleaseGithubTaskString = "%s: sha:%s tag:%s"
)

type ReleaseGithubTask struct {
	Client       *github.Client
	Dir          string
	Organisation string
	Project      string
	Sha          string
	Tag          string
}

func NewReleaseGithubTaskclient(client *github.Client) ReleaseGithubTask {
	task := ReleaseGithubTask{
		Client: client,
	}

	return task
}

// Run creates a draft github release.
func (r ReleaseGithubTask) Run() error {
	info := releaseInfo{
		AssetsDir:    r.Dir,
		Draft:        true,
		Organisation: r.Organisation,
		Project:      r.Project,
		Sha:          r.Sha,
		Tag:          r.Tag,
	}
	return ensureWithDir(r.Client, info)
}

func (r ReleaseGithubTask) Name() string {
	return ReleaseGithubTaskName
}

func (r ReleaseGithubTask) String() string {
	return fmt.Sprintf(ReleaseGithubTaskString, r.Name(), r.Sha, r.Tag)
}
