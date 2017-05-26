package commands

import (
	"fmt"
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
	Network          string
	Args             []string
}

type Command struct {
	Name string
	Args []string
}

func (c Command) String() string {
	redactionPhrases := []string{
		"password",
	}

	redactedArgs := []string{}
	for _, arg := range c.Args {
		requiresRedaction := false

		for _, redactionPhrase := range redactionPhrases {
			if strings.Contains(arg, redactionPhrase) {
				requiresRedaction = true
			}
		}

		if requiresRedaction {
			parts := strings.Split(arg, "=")
			if len(parts) == 2 {
				parts[1] = "[REDACTED]"
				arg = parts[0] + "=" + parts[1]
			}
		}

		redactedArgs = append(redactedArgs, arg)
	}

	return fmt.Sprintf("%s:\t'%s'", c.Name, strings.Join(redactedArgs, " "))
}

func NewDockerCommand(name string, config DockerCommandConfig) Command {
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

	if config.Network != "" {
		args = append(args, "--network="+config.Network)
	}

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

	log.Printf("running command: %s\n", command)

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
