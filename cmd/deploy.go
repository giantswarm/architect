package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"
	"github.com/giantswarm/architect/workflow"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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

	resourcesDirectoryPath string
)

func init() {
	RootCmd.AddCommand(deployCmd)

	var defaultDockerEmail string
	var defaultDockerUsername string
	var defaultDockerPassword string

	if os.Getenv("CIRCLECI") == "true" {
		defaultDockerEmail = os.Getenv("DOCKER_EMAIL")
		defaultDockerUsername = os.Getenv("DOCKER_USERNAME")
		defaultDockerPassword = os.Getenv("DOCKER_PASSWORD")
	}

	deployCmd.Flags().StringVar(&dockerEmail, "docker-email", defaultDockerEmail, "email to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerUsername, "docker-username", defaultDockerUsername, "username to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerPassword, "docker-password", defaultDockerPassword, "password to use to login to docker registry")

	deployCmd.Flags().StringVar(&resourcesDirectoryPath, "resources-directory-path", "./kubernetes", "directory holding kubernetes resources")
}

func runDeploy(cmd *cobra.Command, args []string) {
	fs := afero.NewOsFs()

	clusters, err := workflow.ClustersFromEnv(fs, workingDirectory)
	if err != nil {
		log.Fatalf("could not get clusters from env: %v\n", err)
	}

	// Manage kubernetes resource templating
	if err := utils.TemplateKubernetesResources(fs, resourcesDirectoryPath, sha); err != nil {
		log.Fatalf("could not template kubernetes resources: %v\n", err)
	}

	templatedResourcesDirectoryAbsolutePath, err := filepath.Abs(resourcesDirectoryPath)
	if err != nil {
		log.Fatalf("could not get absolute path for templated resources directory: %v\n", err)
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

		KubernetesTemplatedResourcesDirectoryPath: templatedResourcesDirectoryAbsolutePath,
		KubernetesClusters:                        clusters,
	}

	workflow, err := workflow.NewDeploy(projectInfo, fs)
	if err != nil {
		log.Fatalf("could not get workflow: %v", err)
	}

	log.Printf("running workflow: %s\n", workflow)

	if dryRun {
		log.Printf("dry run, not actually running workflow")
		return
	}

	commands.RunCommands(workflow)
}
