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
	Client *github.Client

	AssetsDir    string
	Draft        bool
	Organisation string
	Project      string
	Sha          string
	Tag          string
}

// Run ensures a github release.
func (r ReleaseGithubTask) Run() error {
	return r.ensureWithDir()
}

func (r ReleaseGithubTask) Name() string {
	return ReleaseGithubTaskName
}

func (r ReleaseGithubTask) String() string {
	return fmt.Sprintf(ReleaseGithubTaskString, r.Name(), r.Sha, r.Tag)
}
