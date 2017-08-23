package cmd

import (
	"log"
	"runtime"

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

	goosTest  string
	goosBuild string
	goarch    string

	golangImage   string
	golangVersion string
)

func init() {
	RootCmd.AddCommand(buildCmd)

	// defaultGoosTest is the default GOOS to be used while testing.
	// As we test inside a Docker container, this should always be Linux.
	defaultGoosTest := "linux"
	// defaultGoosBuild is the default GOOS to be used while building.
	// As we build inside a Docker container, we detect the current platform.
	defaultGoosBuild := runtime.GOOS

	buildCmd.Flags().StringVar(&goosTest, "goos-test", defaultGoosTest, "value for $GOOS while testing")
	buildCmd.Flags().StringVar(&goosBuild, "goos-build", defaultGoosBuild, "value for $GOOS while building")

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

		GoosTest:      goosTest,
		GoosBuild:     goosBuild,
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
