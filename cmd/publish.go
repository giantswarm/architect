package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
)

var (
	publishCmd = &cobra.Command{
		Use:   "publish",
		Short: "publish charts to the specified channels",
		Run:   runPublish,
	}

	channels string
)

func init() {
	RootCmd.AddCommand(deployCmd)

	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" {
		defaultDockerUsername = os.Getenv("QUAY_USERNAME")
		defaultDockerPassword = os.Getenv("QUAY_PASSWORD")
	}

	publishCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	publishCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	publishCmd.Flags().StringVar(&channels, "channels", "beta,testing", "channels to publish the charts to, separated by comma")
}

func runPublish(cmd *cobra.Command, args []string) {
	fs := afero.NewOsFs()

	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Branch: branch,
		Sha:    sha,

		Registry:       registry,
		DockerUsername: dockerUsername,
		DockerPassword: dockerPassword,

		Channels: strings.Split(channels, ","),
	}

	workflow, err := workflow.NewPublish(projectInfo, fs)
	if err != nil {
		log.Fatalf("could not get workflow: %v", err)
	}

	log.Printf("running workflow: %s\n", workflow)

	if dryRun {
		log.Printf("dry run, not actually running workflow")
		return
	}

	if err := tasks.Run(workflow); err != nil {
		log.Fatalf("could not execute workflow: %v", err)
	}
}
