package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

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

	branch string
	sha    string
	tag    string

	version string

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

	// Use git HEAD as defaultSha.
	var defaultSha string
	{
		out, err := exec.Command("git", "rev-parse", "HEAD").Output()
		if err != nil {
			log.Fatalf("could not get git sha: %#v\n", err)
		}
		defaultSha = strings.TrimSpace(string(out))
	}

	// Use git tag when available.
	var defaultTag string
	{
		out, err := exec.Command("git", "describe", "--tags", "--exact-match", "HEAD").Output()
		if err == nil {
			defaultTag = strings.TrimPrefix(strings.TrimSpace(string(out)), "v")
		}
	}

	// We also use the git HEAD branch as well.
	var defaultBranch string
	{
		out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		if err != nil {
			log.Fatalf("could not get git branch: %#v\n", err)
		}
		defaultBranch = strings.TrimSpace(string(out))
	}

	// Define the version we are building.
	var defaultVersion string
	{
		// version can be of three different formats:
		//   2.0.0: building a tagged version.
		//   2.0.0-3a955cbb126f0fe5d51aedf2eb84acca7b074374: building ahead of a tagged version.
		//   1.0.0-939f5c6949f83c0a7ea98a25bc9524fd2f751ffe: building a repo which has no tags.
		if defaultTag != "" {
			defaultVersion = defaultTag
		} else {
			out, err := exec.Command("git", "describe", "--tags", "--abbrev=0", "HEAD").Output()
			if err != nil {
				defaultVersion = fmt.Sprintf("1.0.0-%s", defaultSha)
			} else {
				defaultVersion = fmt.Sprintf("%s-%s", strings.TrimPrefix(strings.TrimSpace(string(out)), "v"), defaultSha)
			}

		}
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

	RootCmd.PersistentFlags().StringVar(&branch, "branch", defaultBranch, "git branch being built")
	RootCmd.PersistentFlags().StringVar(&sha, "sha", defaultSha, "git SHA1 being built")
	RootCmd.PersistentFlags().StringVar(&tag, "tag", defaultTag, "git tag being built")

	RootCmd.PersistentFlags().StringVar(&version, "version", defaultVersion, "project version being built")

	RootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", dryRun, "show what would be executed, but take no action")

	RootCmd.AddCommand(release.Cmd)
}
