package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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
)

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&dockerEmail, "docker-email", "", "email to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerUsername, "docker-username", "", "username to use to login to docker registry")
	deployCmd.Flags().StringVar(&dockerPassword, "docker-password", "", "password to use to login to docker registry")

	deployCmd.Flags().StringVar(&kubernetesApiServer, "kubernetes-api-server", "https://api.g8s.fra-1.giantswarm.io", "kubernetes api to deploy to")
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
}
