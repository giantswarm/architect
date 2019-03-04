package cmd

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
)

var (
	releaseCmd = &cobra.Command{
		Use:   "release",
		Short: "release chart as github release",
		Run:   runRelease,
	}
)

func init() {
	RootCmd.AddCommand(releaseCmd)
}

func runRelease(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	fs := afero.NewOsFs()

	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Sha: sha,
		Tag: tag,
	}

	var githubClient *github.Client
	{
		if deploymentEventsToken == "" {
			log.Fatalf("no github token")
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: deploymentEventsToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		githubClient = github.NewClient(tc)
	}

	workflow, err := workflow.NewRelease(projectInfo, fs, githubClient)
	if err != nil {
		log.Fatalf("could not get workflow: %v", err)
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
