package release

import "os"

var (
	push bool
)

func init() {
	defaultDockerUsername := os.Getenv("QUAY_USERNAME")
	defaultDockerPassword := os.Getenv("QUAY_PASSWORD")

	Cmd.Flags().String("docker-username", defaultDockerUsername, "username to use to login to docker registry")
	Cmd.Flags().String("docker-password", defaultDockerPassword, "password to use to login to docker registry")
	Cmd.Flags().BoolVar(&push, "push", true, "publish helm chart as a github release")
}
