package cmd

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/architect/cmd/create"
	"github.com/giantswarm/architect/cmd/helm"
	"github.com/giantswarm/architect/cmd/legacy"
	"github.com/giantswarm/architect/cmd/preparerelease"
	cmdProject "github.com/giantswarm/architect/cmd/project"
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
	var err error

	var logger micrologger.Logger
	{
		c := micrologger.Config{}

		logger, err = micrologger.New(c)
		if err != nil {
			panic(microerror.Pretty(microerror.Mask(err), true))
		}

	}

	var stderr, stdout io.Writer
	{
		stderr = os.Stderr
		stdout = os.Stdout
	}

	var legacyCmd *cobra.Command
	{
		c := legacy.Config{
			Logger: logger,
			Stderr: stderr,
			Stdout: stdout,
		}

		legacyCmd, err = legacy.New(c)
		if err != nil {
			panic(microerror.Pretty(microerror.Mask(err), true))
		}
	}

	// Get the current working directory, to use as a default.
	defaultWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v\n", err)
	}

	var defaultOrganisation string
	var defaultProject string
	if os.Getenv("CIRCLECI") == "true" { // nolint:goconst
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

	RootCmd.AddCommand(cmdProject.Cmd)
	RootCmd.AddCommand(create.Cmd)
	RootCmd.AddCommand(helm.Cmd)
	RootCmd.AddCommand(legacyCmd)
	RootCmd.AddCommand(preparerelease.Cmd)
}
