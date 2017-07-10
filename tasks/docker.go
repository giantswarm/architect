package tasks

import (
	"os"
)

type DockerTaskConfig struct {
	Volumes          []string
	Env              []string
	WorkingDirectory string
	Image            string
	Args             []string
}

func NewDockerTask(name string, config DockerTaskConfig) ExecTask {
	args := []string{
		"docker", "run",
	}

	// CircleCI struggles with intermediate images, this helps to deal with that.
	// See https://circleci.com/docs/docker/#deployment-to-a-docker-registry
	if os.Getenv("CIRCLECI") == "true" {
		args = append(args, "--rm=false")
	} else {
		args = append(args, "--rm")
	}

	for _, volume := range config.Volumes {
		args = append(args, "-v", volume)
	}

	for _, env := range config.Env {
		args = append(args, "-e", env)
	}

	args = append(args, "-w", config.WorkingDirectory)
	args = append(args, config.Image)

	for _, arg := range config.Args {
		args = append(args, arg)
	}

	newDockerTask := NewExecTask(name, args)

	return newDockerTask
}
