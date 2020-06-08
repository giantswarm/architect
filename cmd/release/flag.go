package release

import (
	"os"
)

func init() {
	defaultDockerUsername := os.Getenv("QUAY_USERNAME")
	defaultDockerPassword := os.Getenv("QUAY_PASSWORD")

	Cmd.Flags().String("docker-username", defaultDockerUsername, "username to use to login to docker registry")
	Cmd.Flags().String("docker-password", defaultDockerPassword, "password to use to login to docker registry")
}
