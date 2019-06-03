package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/giantswarm/architect/cmd/helm"
	"github.com/giantswarm/architect/cmd/release"
)

var (
	RootCmd = &cobra.Command{
		Use:   "architect",
		Short: "Architect is a tool for managing builds within Giant Swarm release engineering.",
	}

	workingDirectory string

	registry     string
	organisation string
	project      string

	dryRun bool
)

func init() {
	// Get the current working directory, to use as a default.
	defaultWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v\n", err)
	}

	var defaultOrganisation string
	var defaultProject string
	if os.Getenv("CIRCLECI") == "true" {
		// If we're running on CircleCI, we can be smart with the organisation and project values.
		defaultOrganisation = os.Getenv("CIRCLE_PROJECT_USERNAME")
		defaultProject = os.Getenv("CIRCLE_PROJECT_REPONAME")
	} else {
		// If running elsewhere, we can attempt to infer the organisation and project value from the path.
		path := strings.Split(defaultWorkingDirectory, string(os.PathSeparator))
		defaultOrganisation = path[len(path)-2]
		defaultProject = path[len(path)-1]
	}

	var githubToken string
	{
		githubToken = os.Getenv("DEPLOYMENT_EVENTS_TOKEN")
	}

	RootCmd.PersistentFlags().StringVar(&workingDirectory, "working-directory", defaultWorkingDirectory, "working directory for architect")
	RootCmd.PersistentFlags().String("github-token", githubToken, "github OAuth access token")

	RootCmd.PersistentFlags().StringVar(&registry, "registry", "quay.io", "docker registry")
	RootCmd.PersistentFlags().StringVar(&organisation, "organisation", defaultOrganisation, "organisation who owns the project")
	RootCmd.PersistentFlags().StringVar(&project, "project", defaultProject, "name of the project")

	RootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", dryRun, "show what would be executed, but take no action")

	RootCmd.AddCommand(release.Cmd)
	RootCmd.AddCommand(helm.Cmd)
}
