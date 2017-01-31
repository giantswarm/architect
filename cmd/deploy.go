package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	kubernetesResourcesDirectoryPath string
	templatedResourcesDirectoryPath  string
	removeResourceFilesAfterUse      bool
)

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&dockerEmail, "docker-email", "", "email to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerUsername, "docker-username", "", "username to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerPassword, "docker-password", "", "password to use to login to docker registry")

	deployCmd.Flags().StringVar(&kubernetesApiServer, "kubernetes-api-server", "https://api.g8s.fra-1.giantswarm.io", "kubernetes api to deploy to")

	deployCmd.Flags().StringVar(&kubernetesCaPath, "kubernetes-ca-path", "", "path to kubernetes ca file")
	deployCmd.Flags().StringVar(&kubernetesCrtPath, "kubernetes-crt-path", "", "path to kubernetes certificate file")
	deployCmd.Flags().StringVar(&kubernetesKeyPath, "kubernetes-key-path", "", "path to kubernetes key file")

	deployCmd.Flags().StringVar(&kubectlVersion, "kubectl-version", "1.5.2", "kubectl version")

	deployCmd.Flags().StringVar(&kubernetesResourcesDirectoryPath, "kubernetes-resources-directory-path", "./kubernetes", "directory holding kubernetes resources")
	deployCmd.Flags().StringVar(&templatedResourcesDirectoryPath, "templated-resources-directory-path", "./kubernetes-templated", "directory holding templated kubernetes resources")
	deployCmd.Flags().BoolVar(&removeResourceFilesAfterUse, "remove-resource-files-after-use", true, "whether to remove templated kubernetes resource files after use")
}

func runDeploy(cmd *cobra.Command, args []string) {
	if dockerEmail == "" {
		log.Fatalf("specify docker email\n")
	}
	if dockerUsername == "" {
		log.Fatalf("specify docker username\n")
	}
	if dockerPassword == "" {
		log.Fatalf("specify docker password\n")
	}

	dockerLoginCommandArgs := []string{
		"login",
		fmt.Sprintf("--email=%v", dockerEmail),
		fmt.Sprintf("--username=%v", dockerUsername),
		fmt.Sprintf("--password=%v", dockerPassword),
		registry,
	}

	dockerLoginCommand := exec.Command("docker", dockerLoginCommandArgs...)

	dockerLoginCommand.Stdout = os.Stdout
	dockerLoginCommand.Stderr = os.Stderr

	log.Printf("running %v\n", dockerLoginCommand.Args)
	if err := dockerLoginCommand.Run(); err != nil {
		log.Fatalf("could not run docker login command: %v\n", err)
	}

	dockerPushCommandArgs := []string{
		"push",
		fmt.Sprintf("%v/%v/%v:%v", registry, organisation, project, sha),
	}

	dockerPushCommand := exec.Command("docker", dockerPushCommandArgs...)

	dockerPushCommand.Stdout = os.Stdout
	dockerPushCommand.Stderr = os.Stderr

	log.Printf("running %v\n", dockerPushCommand.Args)
	if err := dockerPushCommand.Run(); err != nil {
		log.Fatalf("could not run docker push command: %v\n", err)
	}

	kubectlClusterInfoCommandArgs := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%v:/ca.pem", kubernetesCaPath),
		"-v", fmt.Sprintf("%v:/crt.pem", kubernetesCrtPath),
		"-v", fmt.Sprintf("%v:/key.pem", kubernetesKeyPath),
		fmt.Sprintf("giantswarm/kubectl:%v", kubectlVersion),
		fmt.Sprintf("--server=%v", kubernetesApiServer),
		"--certificate-authority=/ca.pem",
		"--client-certificate=/crt.pem",
		"--client-key=/key.pem",
		"cluster-info",
	}

	kubectlClusterInfoCommand := exec.Command("docker", kubectlClusterInfoCommandArgs...)

	kubectlClusterInfoCommand.Stdout = os.Stdout
	kubectlClusterInfoCommand.Stderr = os.Stderr

	log.Printf("running %v\n", kubectlClusterInfoCommand.Args)
	if err := kubectlClusterInfoCommand.Run(); err != nil {
		log.Fatalf("could not run kubectl cluster info command: %v\n", err)
	}

	log.Printf("templating kubernetes resources")

	if _, err := os.Stat(templatedResourcesDirectoryPath); os.IsNotExist(err) {
		if err := os.Mkdir(templatedResourcesDirectoryPath, 0755); err != nil {
			log.Fatalf("could not create templated resources directory: %v\n", err)
		}
	}

	files, err := ioutil.ReadDir(kubernetesResourcesDirectoryPath)
	if err != nil {
		log.Fatalf("could not read kubernetes resources directory: %v\n", err)
	}

	for _, file := range files {
		contents, err := ioutil.ReadFile(filepath.Join(kubernetesResourcesDirectoryPath, file.Name()))
		if err != nil {
			log.Fatalf("could not read file: %v\n", err)
		}

		templatedContents := strings.Replace(string(contents), "%%DOCKER_TAG%%", sha, -1)

		if err := ioutil.WriteFile(filepath.Join(templatedResourcesDirectoryPath, file.Name()), []byte(templatedContents), 0755); err != nil {
			log.Fatalf("could not write file: %v\n", err)
		}
	}

	if removeResourceFilesAfterUse {
		if err := os.RemoveAll(templatedResourcesDirectoryPath); err != nil {
			log.Fatalf("could not remove templated resources directory: %v\n", err)
		}
	}

	templatedResourcesDirectoryAbsolutePath, err := filepath.Abs(templatedResourcesDirectoryPath)
	if err != nil {
		log.Fatalf("could not get absolute path for templated resources directory: %v\n", err)
	}

	kubectlApplyCommandArgs := []string{
		"run",
		"--rm",
		"-v", fmt.Sprintf("%v:/ca.pem", kubernetesCaPath),
		"-v", fmt.Sprintf("%v:/crt.pem", kubernetesCrtPath),
		"-v", fmt.Sprintf("%v:/key.pem", kubernetesKeyPath),
		"-v", fmt.Sprintf("%v:/kubernetes", templatedResourcesDirectoryAbsolutePath),
		fmt.Sprintf("giantswarm/kubectl:%v", kubectlVersion),
		fmt.Sprintf("--server=%v", kubernetesApiServer),
		"--certificate-authority=/ca.pem",
		"--client-certificate=/crt.pem",
		"--client-key=/key.pem",
		"apply", "-f", "/kubernetes",
	}

	kubectlApplyCommand := exec.Command("docker", kubectlApplyCommandArgs...)

	kubectlApplyCommand.Stdout = os.Stdout
	kubectlApplyCommand.Stderr = os.Stderr

	log.Printf("running %v\n", kubectlApplyCommand.Args)
	if err := kubectlApplyCommand.Run(); err != nil {
		log.Fatalf("could not run kubectl apply command: %v\n", err)
	}
}
