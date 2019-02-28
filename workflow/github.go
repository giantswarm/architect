package workflow

import (
	"github.com/google/go-github/github"

	"github.com/giantswarm/architect/release"
	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/microerror"
)

func NewReleaseGithubTask(client *github.Client, dir string, projectInfo ProjectInfo) (tasks.Task, error) {
	err := checkReleaseRequirements(projectInfo)
	if err != nil {
		microerror.Mask(err)
	}

	githubRelease := release.ReleaseGithubTask{
		Client:       client,
		Dir:          dir,
		Organisation: projectInfo.Organisation,
		Project:      projectInfo.Project,
		Sha:          projectInfo.Sha,
		Tag:          projectInfo.Tag,
	}

	return githubRelease, nil
}

func checkReleaseRequirements(projectInfo ProjectInfo) error {
	if projectInfo.Organisation == "" {
		return microerror.Mask(emptyOrganisationError)
	}
	if projectInfo.Project == "" {
		return microerror.Mask(emptyProjectError)
	}
	if projectInfo.Sha == "" {
		return microerror.Mask(emptyShaError)
	}
	if projectInfo.Tag == "" {
		return microerror.Mask(emptyRefError)
	}

	return nil
}
