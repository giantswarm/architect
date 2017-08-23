package workflow

import (
	"fmt"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

var (
	DockerBuildTaskName      = "docker-build"
	DockerRunVersionTaskName = "docker-run-version"
	DockerRunHelpTaskName    = "docker-run-help"

	DockerLoginTaskName      = "docker-login"
	DockerTagLatestTaskName  = "docker-tag-latest"
	DockerPushShaTaskName    = "docker-push-sha"
	DockerPushLatestTaskName = "docker-push-latest"
)

func checkDockerRequirements(projectInfo ProjectInfo) error {
	if projectInfo.WorkingDirectory == "" {
		return microerror.Mask(emptyWorkingDirectoryError)
	}
	if projectInfo.Organisation == "" {
		return microerror.Mask(emptyOrganisationError)
	}
	if projectInfo.Project == "" {
		return microerror.Mask(emptyProjectError)
	}

	if projectInfo.Sha == "" {
		return microerror.Mask(emptyShaError)
	}
	if projectInfo.Registry == "" {
		return microerror.Mask(emptyRegistryError)
	}

	return nil
}

func NewDockerBuildTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	dockerBuild := tasks.NewExecTask(
		DockerBuildTaskName,
		[]string{
			"docker",
			"build",
			"-t",
			fmt.Sprintf(
				"%v/%v/%v:%v",
				projectInfo.Registry,
				projectInfo.Organisation,
				projectInfo.Project,
				projectInfo.Sha,
			),
			projectInfo.WorkingDirectory,
		},
	)

	return dockerBuild, nil
}

func NewDockerRunVersionTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	dockerRunVersion := tasks.NewDockerTask(
		DockerRunVersionTaskName,
		tasks.DockerTaskConfig{
			Image: fmt.Sprintf(
				"%v/%v/%v:%v",
				projectInfo.Registry,
				projectInfo.Organisation,
				projectInfo.Project,
				projectInfo.Sha,
			),
			Args: []string{"version"},
		},
	)

	return dockerRunVersion, nil
}

func NewDockerRunHelpTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	dockerRunHelp := tasks.NewDockerTask(
		DockerRunHelpTaskName,
		tasks.DockerTaskConfig{
			Image: fmt.Sprintf(
				"%v/%v/%v:%v",
				projectInfo.Registry,
				projectInfo.Organisation,
				projectInfo.Project,
				projectInfo.Sha,
			),
			Args: []string{"--help"},
		},
	)

	return dockerRunHelp, nil
}

func NewDockerLoginTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	if projectInfo.DockerUsername == "" {
		return nil, microerror.Mask(emptyDockerUsernameError)
	}
	if projectInfo.DockerPassword == "" {
		return nil, microerror.Mask(emptyDockerPasswordError)
	}

	// CircleCI's Docker version still expects to be given an email,
	// even though it is not used by quay.
	// If we don't specify one, use the empty string.
	if projectInfo.DockerEmail == "" {
		projectInfo.DockerEmail = `" "`
	}

	dockerLogin := tasks.NewExecTask(
		DockerLoginTaskName,
		[]string{
			"docker",
			"login",
			fmt.Sprintf("--email=%v", projectInfo.DockerEmail),
			fmt.Sprintf("--username=%v", projectInfo.DockerUsername),
			fmt.Sprintf("--password=%v", projectInfo.DockerPassword),
			projectInfo.Registry,
		},
	)

	return dockerLogin, nil
}

func NewDockerTagLatestTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	dockerPush := tasks.NewExecTask(
		DockerTagLatestTaskName,
		[]string{
			"docker",
			"tag",
			fmt.Sprintf(
				"%v/%v/%v:%v",
				projectInfo.Registry,
				projectInfo.Organisation,
				projectInfo.Project,
				projectInfo.Sha,
			),
			fmt.Sprintf(
				"%v/%v/%v:latest",
				projectInfo.Registry,
				projectInfo.Organisation,
				projectInfo.Project,
			),
		},
	)

	return dockerPush, nil
}

func NewDockerPushShaTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	dockerPush := tasks.NewExecTask(
		DockerPushShaTaskName,
		[]string{
			"docker",
			"push",
			fmt.Sprintf(
				"%v/%v/%v:%v",
				projectInfo.Registry,
				projectInfo.Organisation,
				projectInfo.Project,
				projectInfo.Sha,
			),
		},
	)

	return dockerPush, nil
}

func NewDockerPushLatestTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	dockerPush := tasks.NewExecTask(
		DockerPushLatestTaskName,
		[]string{
			"docker",
			"push",
			fmt.Sprintf(
				"%v/%v/%v:latest",
				projectInfo.Registry,
				projectInfo.Organisation,
				projectInfo.Project,
			),
		},
	)

	return dockerPush, nil
}
