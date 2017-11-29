package workflow

import (
	"fmt"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

const (
	DockerBuildTaskName      = "docker-build"
	DockerRunVersionTaskName = "docker-run-version"
	DockerRunHelpTaskName    = "docker-run-help"

	DockerLoginTaskName      = "docker-login"
	DockerTagLatestTaskName  = "docker-tag-latest"
	DockerPushShaTaskName    = "docker-push-sha"
	DockerPushLatestTaskName = "docker-push-latest"

	// DockerImageRefFmt is the format string used to compute the reference of the
	// Docker image used to build and push. It may look something like this.
	//
	//     quay.io/giantswarm/architect:e8363ac222255e991c126abe6673cd0f33934ac8
	//
	DockerImageRefFmt    = "%s/%s/%s:%s"
	LatestDockerImageTag = "latest"
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
			newDockerImageRef(projectInfo, projectInfo.Sha),
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
			Args:             []string{"version"},
			Image:            newDockerImageRef(projectInfo, projectInfo.Sha),
			WorkingDirectory: projectInfo.WorkingDirectory,
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
			Args:             []string{"--help"},
			Image:            newDockerImageRef(projectInfo, projectInfo.Sha),
			WorkingDirectory: projectInfo.WorkingDirectory,
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

	dockerLogin := tasks.NewExecTask(
		DockerLoginTaskName,
		[]string{
			"docker",
			"login",
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
			newDockerImageRef(projectInfo, projectInfo.Sha),
			newDockerImageRef(projectInfo, LatestDockerImageTag),
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
			newDockerImageRef(projectInfo, projectInfo.Sha),
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
			newDockerImageRef(projectInfo, LatestDockerImageTag),
		},
	)

	return dockerPush, nil
}

func newDockerImageRef(projectInfo ProjectInfo, version string) string {
	return fmt.Sprintf(
		DockerImageRefFmt,
		projectInfo.Registry,
		projectInfo.Organisation,
		projectInfo.Project,
		version,
	)
}
