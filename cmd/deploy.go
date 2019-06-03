package cmd

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/giantswarm/architect/cmd/hook"
	"github.com/giantswarm/architect/events"
	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/workflow"
)

var (
	deployCmd = &cobra.Command{
		Use:     "deploy",
		Short:   "deploy the project",
		Run:     runDeploy,
		PreRunE: hook.PreRunE,
	}

	deploymentEventsToken string
	group                 string
)

func init() {
	RootCmd.AddCommand(deployCmd)

	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" {
		defaultDockerUsername = os.Getenv("QUAY_USERNAME")
		defaultDockerPassword = os.Getenv("QUAY_PASSWORD")

		deploymentEventsToken = os.Getenv("DEPLOYMENT_EVENTS_TOKEN")
	}

	deployCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	deployCmd.Flags().StringVar(&helmDirectoryPath, "helm-directory-path", "./helm", "directory holding helm chart")
	deployCmd.Flags().StringVar(&group, "group", "all", "the group you want to create deployment events for. Can be 'all' or 'testing'")
}

func runDeploy(cmd *cobra.Command, args []string) {
	fs := afero.NewOsFs()

	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,

		Branch: cmd.Flag("branch").Value.String(),
		Sha:    cmd.Flag("sha").Value.String(),
		Tag:    cmd.Flag("tag").Value.String(),

		Version: cmd.Flag("version").Value.String(),

		Registry:       registry,
		DockerUsername: dockerUsername,
		DockerPassword: dockerPassword,

		HelmDirectoryPath: helmDirectoryPath,
	}

	workflow, err := workflow.NewDeploy(projectInfo, fs)
	if err != nil {
		log.Fatalf("could not get workflow: %v", err)
	}

	log.Printf("running workflow: %s\n", workflow)

	if dryRun {
		log.Printf("dry run, not actually running workflow or creating events")
		return
	}

	if err := tasks.Run(workflow); err != nil {
		log.Fatalf("could not execute workflow: %v", err)
	}

	if deploymentEventsToken == "" {
		log.Printf("no deployment events token, not creating deployments event")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: deploymentEventsToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient := github.NewClient(tc)

	log.Printf("creating deployment events")
	environments := events.GetEnvironments(project, group)

	log.Printf("creating for environments: %v", environments)
	for _, environment := range environments {
		if err := events.CreateDeploymentEvent(githubClient, environment, organisation, project, cmd.Flag("sha").Value.String()); err != nil {
			log.Fatalf("could not create deployment event: %v", err)
		}
	}
}
