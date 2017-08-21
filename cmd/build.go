package cmd

import (
	"log"

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

	goos   string
	goarch string

	golangImage   string
	golangVersion string
)

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVar(&goos, "goos", "linux", "value for $GOOS")
	buildCmd.Flags().StringVar(&goarch, "goarch", "amd64", "value for $GOARCH")

	buildCmd.Flags().StringVar(&golangImage, "golang-image", "quay.io/giantswarm/golang", "golang image")
	buildCmd.Flags().StringVar(&golangVersion, "golang-version", "1.8.3", "golang version")
}

func runBuild(cmd *cobra.Command, args []string) {
	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Branch: branch,
		Sha:    sha,

		Registry: registry,

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
