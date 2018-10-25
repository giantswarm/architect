package tasks

import (
	"os"
	"reflect"
	"testing"
)

func TestNewDockerTask(t *testing.T) {
	tests := []struct {
		config       DockerTaskConfig
		inCircle     bool
		expectedArgs []string
	}{
		// Test a simple docker task is constructed correctly
		{
			config: DockerTaskConfig{
				Volumes: []string{
					"/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				},
				Env: []string{
					"GOOS=linux",
				},
				WorkingDirectory: "/go/src/github.com/giantswarm/architect",
				Image:            "golang:1.7.5",
				Args:             []string{"go", "test"},
			},
			inCircle: false,
			expectedArgs: []string{
				"docker", "run", "--rm",
				"-v", "/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				"-e", "GOOS=linux",
				"-w", "/go/src/github.com/giantswarm/architect",
				"golang:1.7.5",
				"go", "test",
			},
		},

		// Test a similar config, but running in circle
		{
			config: DockerTaskConfig{
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

		testTask := NewDockerTask("test-task", test.config)

		if !reflect.DeepEqual(test.expectedArgs, testTask.Args) {
			t.Fatalf(
				"%v: expected args did not match returned\n expected:\n%#v\nreturned: \n%#v\n",
				index,
				test.expectedArgs,
				testTask.Args,
			)
		}
	}
}
