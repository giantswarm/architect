package workflow

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/giantswarm/architect/tasks"
)

const (
	HelmImage = "quay.io/giantswarm/docker-helm:006b0db51ec484be8b1bd49990804784a9737ece"
)

var (
	HelmLoginTaskName = "helm-login"
	HelmPushTaskName  = "helm-push"
)

func cnrDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(user.HomeDir, ".cnr"), nil
}

func NewHelmLoginTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	cndDir, err := cnrDirectory()
	if err != nil {
		return nil, err
	}

	helmLogin := tasks.NewDockerTask(
		HelmLoginTaskName,
		tasks.DockerTaskConfig{
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

func NewHelmPushTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	helmDirExists, err := afero.DirExists(fs, filepath.Join(projectInfo.WorkingDirectory, "helm"))
	if err != nil {
		return nil, err
	}
	if !helmDirExists {
		return nil, noHelmDirectoryError
	}

	cnrDir, err := cnrDirectory()
	if err != nil {
		return nil, err
	}

	chartDir := filepath.Join(projectInfo.WorkingDirectory, "helm", fmt.Sprintf("%v-chart", projectInfo.Project))

	helmPush := tasks.NewDockerTask(
		HelmPushTaskName,
		tasks.DockerTaskConfig{
			WorkingDirectory: chartDir,
			Image:            HelmImage,
			Volumes: []string{
				fmt.Sprintf("%v:/root/.cnr/", cnrDir),
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
