package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/hook"
	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
)

var (
	unpublishCmd = &cobra.Command{
		Use:     "unpublish",
		Short:   "unpublish charts from the specified channels",
		Run:     runUnpublish,
		PreRunE: hook.PreRunE,
	}
)

func init() {
	RootCmd.AddCommand(unpublishCmd)

	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" { // nolint:goconst
		defaultDockerUsername = os.Getenv("QUAY_USERNAME")
		defaultDockerPassword = os.Getenv("QUAY_PASSWORD")
	}

	unpublishCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	unpublishCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	unpublishCmd.Flags().StringVar(&channels, "channels", "beta,testing", "channels to unpublish the charts from, separated by comma")
}

func runUnpublish(cmd *cobra.Command, args []string) {
	fs := afero.NewOsFs()

	chs := strings.Split(channels, ",")

	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Branch: cmd.Flag("branch").Value.String(),
		Sha:    cmd.Flag("sha").Value.String(),
		Tag:    cmd.Flag("tag").Value.String(),

		Version: cmd.Flag("version").Value.String(),

		Registry:       registry,
		DockerUsername: dockerUsername,
		DockerPassword: dockerPassword,

		Channels: chs,
	}

	workflow, err := workflow.NewUnpublish(projectInfo, fs)
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
