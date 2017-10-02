package cmd

import (
	"log"
	"os"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "build the project",
		Run:   runBuild,
	}

	dockerEmail    string
	dockerUsername string
	dockerPassword string

	goos          string
	goarch        string
	golangImage   string
	golangVersion string

	helmDirectoryPath string
)

func init() {
	RootCmd.AddCommand(buildCmd)

	var defaultDockerEmail string
	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" {
		defaultDockerEmail = ""
		defaultDockerUsername = os.Getenv("QUAY_USERNAME")
		defaultDockerPassword = os.Getenv("QUAY_PASSWORD")

		deploymentEventsToken = os.Getenv("DEPLOYMENT_EVENTS_TOKEN")
	}

	buildCmd.Flags().StringVar(&dockerEmail, "docker-email", defaultDockerEmail, "email to use to login to docker registry")
	buildCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	buildCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	buildCmd.Flags().StringVar(&goos, "goos", "linux", "value for $GOOS")
	buildCmd.Flags().StringVar(&goarch, "goarch", "amd64", "value for $GOARCH")

	buildCmd.Flags().StringVar(&helmDirectoryPath, "helm-directory-path", "./helm", "directory holding helm chart")

	buildCmd.Flags().StringVar(&golangImage, "golang-image", "quay.io/giantswarm/golang", "golang image")
	buildCmd.Flags().StringVar(&golangVersion, "golang-version", "1.9.0", "golang version")
}

func runBuild(cmd *cobra.Command, args []string) {
	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Branch: branch,
		Sha:    sha,

		Registry:       registry,
		DockerEmail:    dockerEmail,
		DockerUsername: dockerUsername,
		DockerPassword: dockerPassword,

		HelmDirectoryPath: helmDirectoryPath,

		Goos:          goos,
		Goarch:        goarch,
		GolangImage:   golangImage,
		GolangVersion: golangVersion,
	}

	fs := afero.NewOsFs()

	workflow, err := workflow.NewBuild(projectInfo, fs)
	if err != nil {
		log.Fatalf("could not fetch workflow: %v", err)
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
