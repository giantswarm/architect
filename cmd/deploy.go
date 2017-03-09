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

	deployCmd.Flags().StringVar(&kubernetesApiServer, "kubernetes-api-server", "https://api.g8s.fra-1.giantswarm.io", "kubernetes api to deploy to")

	deployCmd.Flags().StringVar(&kubernetesCaPath, "kubernetes-ca-path", "", "path to kubernetes ca file")
	deployCmd.Flags().StringVar(&kubernetesCrtPath, "kubernetes-crt-path", "", "path to kubernetes certificate file")
	deployCmd.Flags().StringVar(&kubernetesKeyPath, "kubernetes-key-path", "", "path to kubernetes key file")

	deployCmd.Flags().StringVar(&kubectlVersion, "kubectl-version", "1.4.7", "kubectl version")

	deployCmd.Flags().StringVar(&resourcesDirectoryPath, "resources-directory-path", "./kubernetes", "directory holding kubernetes resources")
}

func runDeploy(cmd *cobra.Command, args []string) {
	fs := afero.NewOsFs()

	// When running in CircleCI, we can attempt to grab kubernetes certificates from the environment
	if kubernetesCaPath == "" && kubernetesCrtPath == "" && kubernetesKeyPath == "" && os.Getenv("CIRCLECI") == "true" {
		var err error
		kubernetesCaPath, kubernetesCrtPath, kubernetesKeyPath, err = utils.K8SCertsFromEnv(fs, workingDirectory)
		if err != nil {
			log.Printf("could not load kubernetes certificates from env: %v\n", err)
		}
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

		KubernetesApiServer:                       kubernetesApiServer,
		KubernetesCaPath:                          kubernetesCaPath,
		KubernetesCrtPath:                         kubernetesCrtPath,
		KubernetesKeyPath:                         kubernetesKeyPath,
		KubectlVersion:                            kubectlVersion,
		KubernetesTemplatedResourcesDirectoryPath: templatedResourcesDirectoryAbsolutePath,
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
