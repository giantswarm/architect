package tasks

import (
	"reflect"
	"testing"
)

func TestNewDockerTask(t *testing.T) {
	tests := []struct {
		config       DockerTaskConfig
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
				Image:            "golang:1.13.0",
				Args:             []string{"go", "test"},
			},
			expectedArgs: []string{
				"docker", "run", "--rm",
				"-v", "/home/ubuntu/architect:/go/src/github.com/giantswarm/architect",
				"-e", "GOOS=linux",
				"-w", "/go/src/github.com/giantswarm/architect",
				"golang:1.13.0",
				"go", "test",
			},
		},
	}

	for index, test := range tests {
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
