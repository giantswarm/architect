package workflow

import (
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

func NewHelmLoginCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	helmLogin := commands.NewDockerCommand(
		HelmLoginCommandName,
		commands.DockerCommandConfig{
			Image: HelmImage,
		},
	)

	return helmLogin, nil
}

func NewHelmPushCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	helmPush := commands.NewDockerCommand(
		HelmPushCommandName,
		commands.DockerCommandConfig{
			Image: HelmImage,
		},
	)

	return helmPush, nil
}
