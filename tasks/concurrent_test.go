package tasks_test

import (
	"fmt"
	"testing"

	"github.com/giantswarm/architect/tasks"
)

type testTask struct {
	shouldFail bool
}

func (t *testTask) Name() string {
	return "test-task"
}

func (t *testTask) Run() error {
	if t.shouldFail {
		return fmt.Errorf("error!")
	}
	return nil
}

func (t *testTask) String() string {
	return "test-task"
}

func TestRunConcurrent(t *testing.T) {
	tcs := []struct {
		description   string
		tasks         []tasks.Task
		expectedError bool
	}{
		{
			description: "all succeed",
			tasks: []tasks.Task{
				&testTask{},
				&testTask{},
				&testTask{},
			},
		},
		{
			description: "one fail",
			tasks: []tasks.Task{
				&testTask{shouldFail: true},
				&testTask{},
				&testTask{},
			},
			expectedError: true,
		},
		{
			description: "two fail",
			tasks: []tasks.Task{
				&testTask{shouldFail: true},
				&testTask{shouldFail: true},
				&testTask{},
			},
			expectedError: true,
		},
		{
			description: "three fail",
			tasks: []tasks.Task{
				&testTask{shouldFail: true},
				&testTask{shouldFail: true},
				&testTask{shouldFail: true},
			},
			expectedError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			task := tasks.NewConcurrentTask("test-concurrent", tc.tasks...)
			err := task.Run()
			switch {
			case err != nil && !tc.expectedError:
				t.Errorf("unexpected error %v", err)
			case err == nil && tc.expectedError:
				t.Errorf("expected error didn't happen")
			}
		})
	}
}
