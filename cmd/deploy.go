package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/events"
	"github.com/giantswarm/architect/workflow"
)

var (
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy the project",
		Run:   runDeploy,
	}

	dockerEmail    string
	dockerUsername string
	dockerPassword string

	kubernetesApiServer string

	kubernetesCaPath  string
	kubernetesCrtPath string
	kubernetesKeyPath string

	kubectlVersion string

	helmDirectoryPath      string
	resourcesDirectoryPath string

	deploymentEventsToken string
)

func init() {
	RootCmd.AddCommand(deployCmd)

	var defaultDockerEmail string
	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" {
		if os.Getenv("QUAY_USERNAME") != "" {
			defaultDockerEmail = ""
			defaultDockerUsername = os.Getenv("QUAY_USERNAME")
			defaultDockerPassword = os.Getenv("QUAY_PASSWORD")
		} else {
			defaultDockerEmail = os.Getenv("DOCKER_EMAIL")
			defaultDockerUsername = os.Getenv("DOCKER_USERNAME")
			defaultDockerPassword = os.Getenv("DOCKER_PASSWORD")
		}

		deploymentEventsToken = os.Getenv("DEPLOYMENT_EVENTS_TOKEN")
	}

	deployCmd.Flags().StringVar(&dockerEmail, "docker-email", defaultDockerEmail, "email to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	deployCmd.Flags().StringVar(&helmDirectoryPath, "helm-directory-path", "./helm", "directory holding helm chart")
	deployCmd.Flags().StringVar(&resourcesDirectoryPath, "resources-directory-path", "./kubernetes", "directory holding kubernetes resources")
}

func runDeploy(cmd *cobra.Command, args []string) {
	fs := afero.NewOsFs()

	clusters, err := workflow.ClustersFromEnv(fs, workingDirectory)
	if err != nil {
		log.Fatalf("could not get clusters from env: %v\n", err)
	}

	resourcesDirectoryAbsolutePath, err := filepath.Abs(resourcesDirectoryPath)
	if err != nil {
		log.Fatalf("could not get absolute path for resources directory: %v\n", err)
	}

	if os.Getenv("QUAY_USERNAME") != "" {
		registry = "quay.io"
	}

	projectInfo := workflow.ProjectInfo{
		WorkingDirectory: workingDirectory,
		Organisation:     organisation,
		Project:          project,
		Sha:              sha,

		Registry:       registry,
		DockerEmail:    dockerEmail,
		DockerUsername: dockerUsername,
		DockerPassword: dockerPassword,

		HelmDirectoryPath:                helmDirectoryPath,
		KubernetesResourcesDirectoryPath: resourcesDirectoryAbsolutePath,
		KubernetesClusters:               clusters,
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

	commands.RunCommands(workflow)

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
	environments := events.GetEnvironments(project)

	log.Printf("creating for environments: %v", environments)
	for _, environment := range environments {
		if err := events.CreateDeploymentEvent(githubClient, environment, organisation, project, sha); err != nil {
			log.Fatalf("could not create deployment event: %v", err)
		}
	}
}
