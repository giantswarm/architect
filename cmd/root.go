package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "architect",
		Short: "architect is a tool for managing Giant Swarm release engineering",
	}

	workingDirectory string

	registry     string
	organisation string
	project      string
	sha          string
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

	// We also use the git HEAD commit sha as well.
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		log.Fatalf("could not get git sha: %v\n", err)
	}
	defaultSha := strings.TrimSpace(string(out))

	RootCmd.PersistentFlags().StringVar(&workingDirectory, "working-directory", defaultWorkingDirectory, "working directory for architect")

	RootCmd.PersistentFlags().StringVar(&registry, "registry", "registry.giantswarm.io", "docker registry")
	RootCmd.PersistentFlags().StringVar(&organisation, "organisation", defaultOrganisation, "organisation who owns the project")
	RootCmd.PersistentFlags().StringVar(&project, "project", defaultProject, "name of the project")
	RootCmd.PersistentFlags().StringVar(&sha, "sha", defaultSha, "git SHA1 to build")
}
