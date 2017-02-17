package cmd

import (
	"fmt"
	"log"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"
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

	buildCmd.Flags().StringVar(&golangImage, "golang-image", "golang", "golang image")
	buildCmd.Flags().StringVar(&golangVersion, "golang-version", "1.7.5", "golang version")
}

func runBuild(cmd *cobra.Command, args []string) {
	testPackageArguments, err := utils.NoVendor(workingDirectory)
	if err != nil {
		log.Fatalf("could not determine test packages: %v", err)
	}

	goTest := commands.NewDockerCommand(
		"go-test",
		commands.DockerCommandConfig{
			Volumes: []string{
				fmt.Sprintf("%v:/go/src/github.com/%v/%v", workingDirectory, organisation, project),
			},
			Env: []string{
				fmt.Sprintf("GOOS=%v", goos),
				fmt.Sprintf("GOARCH=%v", goarch),
				"GOPATH=/go",
				"CGOENABLED=0",
			},
			WorkingDirectory: fmt.Sprintf("/go/src/github.com/%v/%v", organisation, project),
			Image:            fmt.Sprintf("%v:%v", golangImage, golangVersion),
			Args:             []string{"go", "test", "-v"},
		},
	)
	goTest.Args = append(goTest.Args, testPackageArguments...)

	goBuild := commands.NewDockerCommand(
		"go-build",
		commands.DockerCommandConfig{
			Volumes: []string{
				fmt.Sprintf("%v:/go/src/github.com/%v/%v", workingDirectory, organisation, project),
			},
			Env: []string{
				fmt.Sprintf("GOOS=%v", goos),
				fmt.Sprintf("GOARCH=%v", goarch),
				"GOPATH=/go",
				"CGOENABLED=0",
			},
			WorkingDirectory: fmt.Sprintf("/go/src/github.com/%v/%v", organisation, project),
			Image:            fmt.Sprintf("%v:%v", golangImage, golangVersion),
			Args:             []string{"go", "build", "-v", "-a", "-tags", "netgo"},
		},
	)

	dockerBuild := commands.Command{
		Name: "docker-build",
		Args: []string{
			"docker",
			"build",
			"-t",
			fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
			workingDirectory,
		},
	}

	dockerRunVersion := commands.NewDockerCommand(
		"docker-run-version",
		commands.DockerCommandConfig{
			Image: fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
			Args:  []string{"version"},
		},
	)

	dockerRunHelp := commands.NewDockerCommand(
		"docker-run-help",
		commands.DockerCommandConfig{
			Image: fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
			Args:  []string{"--help"},
		},
	)

	commands.RunCommands([]commands.Command{
		goTest,
		goBuild,
		dockerBuild,
		dockerRunVersion,
		dockerRunHelp,
	})
}
