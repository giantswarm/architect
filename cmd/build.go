package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

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

	golangVersion string
)

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVar(&goos, "goos", "linux", "value for $GOOS")
	buildCmd.Flags().StringVar(&goarch, "goarch", "amd64", "value for $GOARCH")

	buildCmd.Flags().StringVar(&golangVersion, "golang-version", "1.8rc2", "golang version")
}

func runBuild(cmd *cobra.Command, args []string) {
	// Replicate glide novendor, so we don't have to bother with glide at all
	testPackages := []string{}
	directories, err := ioutil.ReadDir(workingDirectory)
	if err != nil {
		log.Fatalf("could not read directories: %v\n", err)
	}
	for _, directory := range directories {
		if !directory.IsDir() {
			continue
		}

		if directory.Name() == "vendor" {
			continue
		}

		files, err := ioutil.ReadDir(directory.Name())
		if err != nil {
			log.Fatalf("could not read files: %v %v\n", directory.Name(), err)
		}
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".go") {
				testPackages = append(testPackages, directory.Name())
				break
			}
		}
	}
	testPackageArguments := []string{"."}
	for _, testPackage := range testPackages {
		testPackageArguments = append(testPackageArguments, fmt.Sprintf("./%v/...", testPackage))
	}

	// Run go test
	goTestCommandArgs := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%v:/go/src/github.com/%v/%v", workingDirectory, organisation, project),
		"-e", fmt.Sprintf("GOOS=%v", goos),
		"-e", fmt.Sprintf("GOARCH=%v", goarch),
		"-e", "GOPATH=/go",
		"-e", "CGOENABLED=0",
		"-w", fmt.Sprintf("/go/src/github.com/%v/%v", organisation, project),
		fmt.Sprintf("golang:%v", golangVersion),
		"go", "test", "-v",
	}
	goTestCommandArgs = append(goTestCommandArgs, testPackageArguments...)

	goTestCommand := exec.Command("docker", goTestCommandArgs...)

	goTestCommand.Stdout = os.Stdout
	goTestCommand.Stderr = os.Stderr

	log.Printf("running %v\n", goTestCommand.Args)
	if err := goTestCommand.Run(); err != nil {
		log.Fatalf("could not run go test command: %v\n", err)
	}

	// Run go build
	goBuildCommandArgs := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%v:/go/src/github.com/%v/%v", workingDirectory, organisation, project),
		"-e", fmt.Sprintf("GOOS=%v", goos),
		"-e", fmt.Sprintf("GOARCH=%v", goarch),
		"-e", "GOPATH=/go",
		"-e", "CGOENABLED=0",
		"-w", fmt.Sprintf("/go/src/github.com/%v/%v", organisation, project),
		fmt.Sprintf("golang:%v", golangVersion),
		"go", "build", "-v", "-a", "-tags", "netgo",
	}

	goBuildCommand := exec.Command("docker", goBuildCommandArgs...)

	goBuildCommand.Stdout = os.Stdout
	goBuildCommand.Stderr = os.Stderr

	log.Printf("running %v\n", goBuildCommand.Args)
	if err := goBuildCommand.Run(); err != nil {
		log.Fatalf("could not run go build command: %v\n", err)
	}

	// Run docker build
	dockerBuildCommandArgs := []string{
		"build",
		"-t",
		fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
		workingDirectory,
	}

	dockerBuildCommand := exec.Command("docker", dockerBuildCommandArgs...)

	dockerBuildCommand.Stdout = os.Stdout
	dockerBuildCommand.Stderr = os.Stderr

	log.Printf("running %v\n", dockerBuildCommand.Args)
	if err := dockerBuildCommand.Run(); err != nil {
		log.Fatalf("could not run docker build command: %v\n", err)
	}

	// Run docker run
	dockerRunCommandArgs := []string{
		"run",
		fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
		"version",
	}

	dockerRunCommand := exec.Command("docker", dockerRunCommandArgs...)

	dockerRunCommand.Stdout = os.Stdout
	dockerRunCommand.Stderr = os.Stderr

	log.Printf("running %v\n", dockerRunCommand.Args)
	if err := dockerRunCommand.Run(); err != nil {
		log.Fatalf("could not run docker run command: %v\n", err)
	}
}
