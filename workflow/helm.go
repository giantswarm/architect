package workflow

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/giantswarm/architect/commands"
	"github.com/spf13/afero"
)

const (
	HelmImage = "quay.io/giantswarm/docker-helm:006b0db51ec484be8b1bd49990804784a9737ece"
)

var (
	HelmLoginCommandName = "helm-login"
	HelmPushCommandName  = "helm-push"
)

func cnrDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(user.HomeDir, ".cnr"), nil
}

func NewHelmLoginCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	cndDir, err := cnrDirectory()
	if err != nil {
		return commands.Command{}, err
	}

	helmLogin := commands.NewDockerCommand(
		HelmLoginCommandName,
		commands.DockerCommandConfig{
			Image: HelmImage,
			Volumes: []string{
				fmt.Sprintf("%v:/root/.cnr/", cndDir),
			},
			Args: []string{
				"registry",
				"login",
				fmt.Sprintf("--user=%v", projectInfo.DockerUsername),
				fmt.Sprintf("--password=%v", projectInfo.DockerPassword),
				projectInfo.Registry,
			},
		},
	)

	return helmLogin, nil
}

func NewHelmPushCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	helmDirExists, err := afero.DirExists(fs, filepath.Join(projectInfo.WorkingDirectory, "helm"))
	if err != nil {
		return commands.Command{}, err
	}
	if !helmDirExists {
		return commands.Command{}, fmt.Errorf("could not find helm directory")
	}

	cndDir, err := cnrDirectory()
	if err != nil {
		return commands.Command{}, err
	}

	chartDir := filepath.Join(projectInfo.WorkingDirectory, "helm", fmt.Sprintf("%v-chart", projectInfo.Project))

	helmPush := commands.NewDockerCommand(
		HelmPushCommandName,
		commands.DockerCommandConfig{
			WorkingDirectory: chartDir,
			Image:            HelmImage,
			Volumes: []string{
				fmt.Sprintf("%v:/root/.cnr/", cndDir),
				fmt.Sprintf("%v:%v", chartDir, chartDir),
			},
			Args: []string{
				"registry",
				"push",
				fmt.Sprintf("--namespace=%v", projectInfo.Organisation),
				"--force",
				projectInfo.Registry,
			},
		},
	)

	return helmPush, nil
}
