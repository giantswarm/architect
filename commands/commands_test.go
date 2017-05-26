package commands

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestCommandString(t *testing.T) {
	tests := []struct {
		command        Command
		expectedString string
	}{
		{
			command: Command{
				Name: "docker-run",
				Args: []string{"docker", "run"},
			},
			expectedString: "docker-run:\t'docker run'",
		},
		{
			command: Command{
				Name: "docker-login",
				Args: []string{"docker", "login", "--email=foo", "--password=bar"},
			},
			expectedString: "docker-login:\t'docker login --email=foo --password=[REDACTED]'",
		},
		{
			command: Command{
				Name: "many-pass",
				Args: []string{"foo", "--first-password=bar", "--second-password=baz"},
			},
			expectedString: "many-pass:\t'foo --first-password=[REDACTED] --second-password=[REDACTED]'",
		},
		{
			command: Command{
				Name: "boolean-flag",
				Args: []string{"foo", "-password"},
			},
			expectedString: "boolean-flag:\t'foo -password'",
		},
	}

	for index, test := range tests {
		returnedString := fmt.Sprintf("%s", test.command)

		if returnedString != test.expectedString {
			t.Fatalf(
				"%v: expected string did not match returned\nexpected:\n%s\nreturned: \n%s\n",
				index,
				test.expectedString,
				returnedString,
			)
		}
	}
}

func TestNewDockerCommand(t *testing.T) {
	tests := []struct {
		config       DockerCommandConfig
		inCircle     bool
		expectedArgs []string
	}{
		// Test a simple docker command is constructed correctly
		{
			config: DockerCommandConfig{
				Volumes: []string{
					"/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				},
				Env: []string{
					"GOOS=linux",
				},
				WorkingDirectory: "/go/src/github.com/giantswarm/architect",
				Image:            "golang:1.7.5",
				Args:             []string{"go", "test", "-v"},
			},
			inCircle: false,
			expectedArgs: []string{
				"docker", "run", "--rm",
				"-v", "/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				"-e", "GOOS=linux",
				"-w", "/go/src/github.com/giantswarm/architect",
				"golang:1.7.5",
				"go", "test", "-v",
			},
		},

		// Test a similar config, but running in circle
		{
			config: DockerCommandConfig{
				Volumes: []string{
					"/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				},
				Env: []string{
					"GOOS=linux",
				},
				WorkingDirectory: "/go/src/github.com/giantswarm/architect",
				Image:            "golang:1.7.5",
				Args:             []string{"go", "test", "-v"},
			},
			inCircle: true,
			expectedArgs: []string{
				"docker", "run", "--rm=false",
				"-v", "/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				"-e", "GOOS=linux",
				"-w", "/go/src/github.com/giantswarm/architect",
				"golang:1.7.5",
				"go", "test", "-v",
			},
		},

		// Test a similar config, but running in circle
		{
			config: DockerCommandConfig{
				Volumes: []string{
					"/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				},
				Env: []string{
					"GOOS=linux",
				},
				Network:          "host",
				WorkingDirectory: "/go/src/github.com/giantswarm/architect",
				Image:            "golang:1.7.5",
				Args:             []string{"go", "test", "-v"},
			},
			inCircle: true,
			expectedArgs: []string{
				"docker", "run", "--rm=false",
				"-v", "/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				"-e", "GOOS=linux",
				"-w", "/go/src/github.com/giantswarm/architect",
				"--network=host",
				"golang:1.7.5",
				"go", "test", "-v",
			},
		},
	}

	for index, test := range tests {
		// Configure circle env vars if needed
		if test.inCircle {
			if err := os.Setenv("CIRCLECI", "true"); err != nil {
				t.Fatalf("could not set circle env var: %v", err)
			}
		} else {
			if err := os.Setenv("CIRCLECI", ""); err != nil {
				t.Fatalf("could not unset circle env var: %v", err)
			}
		}

		testCommand := NewDockerCommand("test-command", test.config)

		if !reflect.DeepEqual(test.expectedArgs, testCommand.Args) {
			t.Fatalf(
				"%v: expected args did not match returned\n expected:\n%#v\nreturned: \n%#v\n",
				index,
				test.expectedArgs,
				testCommand.Args,
			)
		}
	}
}
