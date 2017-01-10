package command

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func configureCommand(command []string, workingDirectory string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)

	if workingDirectory != "" {
		cmd.Dir = workingDirectory
	}

	if cmd.Dir == "" {
		log.Printf("%v\n", strings.Join(cmd.Args, " "))
	} else {
		log.Printf("%v - %v\n", cmd.Dir, strings.Join(cmd.Args, " "))
	}

	return cmd
}

// RunCommand takes a string of arguments, and a working directory.
// Essentially, a fancy wrapper around os/exec.
func Run(command []string, workingDirectory string) error {
	cmd := configureCommand(command, workingDirectory)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func RunWithOutput(command []string, workingDirectory string) (string, error) {
	cmd := configureCommand(command, workingDirectory)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
