package cmd

import (
	"log"
	"os"
	"path"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
	"github.com/jstemmer/go-junit-report/formatter"
	"github.com/jstemmer/go-junit-report/parser"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "build the project",
		Run:   runBuild,
	}

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

	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" {
		defaultDockerUsername = os.Getenv("QUAY_USERNAME")
		defaultDockerPassword = os.Getenv("QUAY_PASSWORD")

		deploymentEventsToken = os.Getenv("DEPLOYMENT_EVENTS_TOKEN")
	}

	buildCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	buildCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	buildCmd.Flags().StringVar(&goos, "goos", "linux", "value for $GOOS")
	buildCmd.Flags().StringVar(&goarch, "goarch", "amd64", "value for $GOARCH")

	buildCmd.Flags().StringVar(&helmDirectoryPath, "helm-directory-path", "./helm", "directory holding helm chart")

	buildCmd.Flags().StringVar(&golangImage, "golang-image", "quay.io/giantswarm/golang", "golang image")
	buildCmd.Flags().StringVar(&golangVersion, "golang-version", "1.10.3", "golang version")
}

func runBuild(cmd *cobra.Command, args []string) {
	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Branch: branch,
		Sha:    sha,

		Registry:       registry,
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

	workflowErr := tasks.Run(workflow)

	// Try to upload JUnit reports, regardless of the workflow state.
	goTestOutputFilePath := path.Join(os.TempDir(), "architect-go-test")
	junitOutputDirectoryPath := path.Join(os.TempDir(), "results", "golang")
	junitOutputFilePath := path.Join(junitOutputDirectoryPath, "results.xml")

	if _, err := os.Stat(goTestOutputFilePath); err == nil {
		log.Printf("generating Junit report at: %v", junitOutputFilePath)

		goTestOutputFile, err := os.Open(goTestOutputFilePath)
		if err != nil {
			log.Fatalf("could not open go test output file: %v", err)
		}
		defer goTestOutputFile.Close()

		report, err := parser.Parse(goTestOutputFile, "")
		if err != nil {
			log.Fatalf("could not parse go test output file: %v", err)
		}

		if _, err := os.Stat(junitOutputDirectoryPath); os.IsNotExist(err) {
			if err := os.MkdirAll(junitOutputDirectoryPath, 0700); err != nil {
				log.Fatalf("could not create junit output directory: %v", err)
			}
		}

		junitOutputFile, err := os.Create(junitOutputFilePath)
		if err != nil {
			log.Fatalf("could not create junit output file: %v", err)
		}
		defer junitOutputFile.Close()

		if err := formatter.JUnitReportXML(report, false, golangVersion, junitOutputFile); err != nil {
			log.Fatalf("could not write junit report: %v", err)
		}
	}

	if workflowErr != nil {
		log.Fatalf("could not execute workflow: %v", err)
	}
}
