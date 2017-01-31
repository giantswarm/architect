package commands

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Name string
	Args []string
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
