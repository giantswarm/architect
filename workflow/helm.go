package workflow

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/template"
	"github.com/giantswarm/microerror"
)

const (
	HelmImage = "quay.io/giantswarm/docker-helm:006b0db51ec484be8b1bd49990804784a9737ece"
)

var (
	HelmPullTaskName  = "helm-pull"
	HelmLoginTaskName = "helm-login"
	HelmPushTaskName  = "helm-push"
)

func cnrDirectory() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", microerror.Mask(err)
	}

	return filepath.Join(user.HomeDir, ".cnr"), nil
}

func NewHelmPullTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	helmPull := tasks.NewExecTask(
		HelmPullTaskName,
		[]string{
			"docker", "pull", HelmImage,
		},
	)

	return helmPull, nil
}

func NewTemplateHelmChartTask(fs afero.Fs, chartDir string, projectInfo ProjectInfo) (tasks.Task, error) {
	templateHelmChart := template.NewTemplateHelmChartTask(
		fs,
		chartDir,
		projectInfo.Sha,
	)

	return templateHelmChart, nil
}

func NewHelmLoginTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	cndDir, err := cnrDirectory()
	if err != nil {
		return nil, microerror.Mask(err)
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

func NewHelmPushTask(fs afero.Fs, chartDir string, projectInfo ProjectInfo) (tasks.Task, error) {
	cnrDir, err := cnrDirectory()
	if err != nil {
		return nil, microerror.Mask(err)
	}

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

func NewHelmPromoteToChannelTask(fs afero.Fs, chartDir string, projectInfo ProjectInfo, channel string) (tasks.Task, error) {
	cnrDir, err := cnrDirectory()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	helmPromoteToChannel := tasks.NewDockerTask(
		fmt.Sprintf("%s-%s", HelmPushTaskName, channel),
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
				fmt.Sprintf("--channel=%s", channel),
				fmt.Sprintf("--namespace=%v", projectInfo.Organisation),
				projectInfo.Registry,
			},
		},
	)

	return helmPromoteToChannel, nil
}
