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

	goTest := commands.Command{
		Name: "go-test",
		Args: []string{
			"docker",
			"run",
			"--rm",
			"-v", fmt.Sprintf("%v:/go/src/github.com/%v/%v", workingDirectory, organisation, project),
			"-e", fmt.Sprintf("GOOS=%v", goos),
			"-e", fmt.Sprintf("GOARCH=%v", goarch),
			"-e", "GOPATH=/go",
			"-e", "CGOENABLED=0",
			"-w", fmt.Sprintf("/go/src/github.com/%v/%v", organisation, project),
			fmt.Sprintf("%v:%v", golangImage, golangVersion),
			"go", "test", "-v",
		},
	}
	goTest.Args = append(goTest.Args, testPackageArguments...)

	goBuild := commands.Command{
		Name: "go-build",
		Args: []string{
			"docker",
			"run",
			"--rm",
			"-v", fmt.Sprintf("%v:/go/src/github.com/%v/%v", workingDirectory, organisation, project),
			"-e", fmt.Sprintf("GOOS=%v", goos),
			"-e", fmt.Sprintf("GOARCH=%v", goarch),
			"-e", "GOPATH=/go",
			"-e", "CGOENABLED=0",
			"-w", fmt.Sprintf("/go/src/github.com/%v/%v", organisation, project),
			fmt.Sprintf("%v:%v", golangImage, golangVersion),
			"go", "build", "-v", "-a", "-tags", "netgo",
		},
	}

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

	dockerRun := commands.Command{
		Name: "docker-run",
		Args: []string{
			"docker",
			"run",
			fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
			"--help",
		},
	}

	commands.RunCommands([]commands.Command{
		goTest,
		goBuild,
		dockerBuild,
		dockerRun,
	})
}
