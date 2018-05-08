package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/pipeline"
	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
)

var (
	publishCmd = &cobra.Command{
		Use:   "publish",
		Short: "publish charts to the specified channels",
		Run:   runPublish,
	}

	channels      string
	pipelineStart bool
	pipelineEnd   bool
)

func init() {
	RootCmd.AddCommand(publishCmd)

	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" {
		defaultDockerUsername = os.Getenv("QUAY_USERNAME")
		defaultDockerPassword = os.Getenv("QUAY_PASSWORD")
	}

	publishCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	publishCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	publishCmd.Flags().StringVar(&channels, "channels", "beta,testing", "channels to publish the charts to, separated by comma")
	publishCmd.Flags().BoolVar(&pipelineStart, "pipeline", true, "specifies if charts should be promoted to the first channel in the pipeline")
	publishCmd.Flags().BoolVar(&pipelineEnd, "stable", false, "specifies if charts should be promoted to the last channel in the pipeline")
}

func runPublish(cmd *cobra.Command, args []string) {
	fs := afero.NewOsFs()

	var chs []string
	if pipelineEnd {
		endChannel, err := pipeline.EndChannel(fs, workingDirectory, project)
		if err != nil {
			log.Fatalf("could not get pipeline end channel: %v", err)
		}
		chs = []string{endChannel}
	} else if pipelineStart {
		startChannel, err := pipeline.StartChannel(fs, workingDirectory, project)
		if err != nil {
			log.Fatalf("could not get pipeline start channel: %v", err)
		}
		chs = []string{startChannel}
	} else {
		chs = strings.Split(channels, ",")
	}

	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Branch: branch,
		Sha:    sha,

		Registry:       registry,
		DockerUsername: dockerUsername,
		DockerPassword: dockerPassword,

		Channels: chs,
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
