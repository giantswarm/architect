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
	client       *github.Client
	dir          string
	organisation string
	project      string
	sha          string
	tag          string
}

// Run creates a draft github release.
func (r ReleaseGithubTask) Run() error {
	return createWithDir(r.client, r.dir, r.organisation, r.project, r.sha, r.tag)
}

func (r ReleaseGithubTask) Name() string {
	return ReleaseGithubTaskName
}

func (r ReleaseGithubTask) String() string {
	return fmt.Sprintf(ReleaseGithubTaskString, r.Name(), r.sha, r.tag)
}

func NewReleaseGithubTask(client *github.Client, dir, organisation, project, sha, tag string) ReleaseGithubTask {
	return ReleaseGithubTask{
		client:       client,
		dir:          dir,
		organisation: organisation,
		project:      project,
		sha:          sha,
		tag:          tag,
	}
}
