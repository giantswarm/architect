package commands

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type DockerCommandConfig struct {
	Volumes          []string
	Env              []string
	WorkingDirectory string
	Image            string
	Args             []string
}

type Command struct {
	Name string
	Args []string
}

func NewDockerCommand(name string, config DockerCommandConfig) Command {
	args := []string{
		"docker", "run",
	}

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

	newDockerCommand := Command{
		Name: name,
		Args: args,
	}

	return newDockerCommand
}

func Run(command Command) error {
	cmd := exec.Command(command.Args[0], command.Args[1:]...)

	log.Printf("running command: %v\n", command.Name)
	log.Printf("%v\n", strings.Join(cmd.Args, " "))

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RunCommands is a helper to run a slice of commands
func RunCommands(commands []Command) {
	for _, command := range commands {
		if err := Run(command); err != nil {
			log.Fatalf("could not execute command '%v': %v\n", command.Name, err)
		}
	}
}
