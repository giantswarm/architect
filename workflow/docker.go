package workflow

import (
	"fmt"

	"github.com/giantswarm/architect/commands"
	"github.com/spf13/afero"
)

var (
	DockerBuildCommandName      = "docker-build"
	DockerRunVersionCommandName = "docker-run-version"
	DockerRunHelpCommandName    = "docker-run-help"

	DockerLoginCommandName = "docker-login"
	DockerPushCommandName  = "docker-push"
)

func checkDockerRequirements(projectInfo ProjectInfo) error {
	if projectInfo.WorkingDirectory == "" {
		return emptyWorkingDirectoryError
	}
	if projectInfo.Organisation == "" {
		return emptyOrganisationError
	}
	if projectInfo.Project == "" {
		return emptyProjectError
	}

	if projectInfo.Sha == "" {
		return emptyShaError
	}
	if projectInfo.Registry == "" {
		return emptyRegistryError
	}

	return nil
}

func NewDockerBuildCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	dockerBuild := commands.Command{
		Name: DockerBuildCommandName,
		Args: []string{
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
	}

	return dockerBuild, nil
}

func NewDockerRunVersionCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	dockerRunVersion := commands.NewDockerCommand(
		DockerRunVersionCommandName,
		commands.DockerCommandConfig{
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

func NewDockerRunHelpCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	dockerRunHelp := commands.NewDockerCommand(
		DockerRunHelpCommandName,
		commands.DockerCommandConfig{
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

func NewDockerLoginCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	if projectInfo.DockerUsername == "" {
		return commands.Command{}, emptyDockerUsernameError
	}
	if projectInfo.DockerPassword == "" {
		return commands.Command{}, emptyDockerPasswordError
	}

	// CircleCI's Docker version still expects to be given an email,
	// even though it is not used by quay.
	// If we don't specify one, use the empty string.
	if projectInfo.DockerEmail == "" {
		projectInfo.DockerEmail = `" "`
	}

	dockerLogin := commands.Command{
		Name: DockerLoginCommandName,
		Args: []string{
			"docker",
			"login",
			fmt.Sprintf("--email=%v", projectInfo.DockerEmail),
			fmt.Sprintf("--username=%v", projectInfo.DockerUsername),
			fmt.Sprintf("--password=%v", projectInfo.DockerPassword),
			projectInfo.Registry,
		},
	}

	return dockerLogin, nil
}

func NewDockerPushCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkDockerRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	dockerPush := commands.Command{
		Name: DockerPushCommandName,
		Args: []string{
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
	}

	return dockerPush, nil
}
