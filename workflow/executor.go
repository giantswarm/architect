package workflow

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// Executor executes commands.
type Executor interface {
	// Run executes the given command, returning any errors.
	RunCommand(arg ...string) error
}

// DryRunExecutor does nothing, successfully.
type DryRunExecutor struct{}

func (e *DryRunExecutor) RunCommand(arg ...string) error {
	log.Printf("Executing: %v", arg)

	return nil
}

// ShellExecutor executes commands in a shell.
type ShellExecutor struct{}

func (e *ShellExecutor) RunCommand(arg ...string) error {
	command := strings.Join(arg, " ")
	args := []string{"-c", command}

	log.Printf("Executing: 'bash %v'", args)

	cmd := exec.Command("bash", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func GetExecutor(dryRun bool) Executor {
	if dryRun {
		return &DryRunExecutor{}
	}

	return &ShellExecutor{}
}
