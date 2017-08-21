package workflow

import (
	"reflect"
	"testing"

	"github.com/giantswarm/architect/tasks"
)

func Test_Workflow_Docker_NewDockerBuildTask(t *testing.T) {
	testCases := []struct {
		TaskFunc     func() (tasks.Task, error)
		ExpectedArgs []string
	}{
		// Test 1, make sure NewDockerBuildTask works as expected.
		{
			TaskFunc: func() (tasks.Task, error) {
				return NewDockerBuildTask(nil, testNewProjectInfo())
			},
			ExpectedArgs: []string{
				"docker",
				"build",
				"-t",
				"quay.io/giantswarm/architect:master-e8363ac222255e991c126abe6673cd0f33934ac8",
				"/usr/code/",
			},
		},

		// Test 2, make sure NewDockerRunVersionTask works as expected.
		{
			TaskFunc: func() (tasks.Task, error) {
				return NewDockerRunVersionTask(nil, testNewProjectInfo())
			},
			ExpectedArgs: []string{
				"docker",
				"run",
				"--rm",
				"-w",
				"/usr/code/",
				"quay.io/giantswarm/architect:master-e8363ac222255e991c126abe6673cd0f33934ac8",
				"version",
			},
		},

		// Test 3, make sure NewDockerRunHelpTask works as expected.
		{
			TaskFunc: func() (tasks.Task, error) {
				return NewDockerRunHelpTask(nil, testNewProjectInfo())
			},
			ExpectedArgs: []string{
				"docker",
				"run",
				"--rm",
				"-w",
				"/usr/code/",
				"quay.io/giantswarm/architect:master-e8363ac222255e991c126abe6673cd0f33934ac8",
				"--help",
			},
		},

		// Test 4, make sure NewDockerLoginTask works as expected.
		{
			TaskFunc: func() (tasks.Task, error) {
				return NewDockerLoginTask(nil, testNewProjectInfo())
			},
			ExpectedArgs: []string{
				"docker",
				"login",
				"--email=\" \"",
				"--username=username",
				"--password=password",
				"quay.io",
			},
		},

		// Test 5, make sure NewDockerTagLatestTask works as expected.
		{
			TaskFunc: func() (tasks.Task, error) {
				return NewDockerTagLatestTask(nil, testNewProjectInfo())
			},
			ExpectedArgs: []string{
				"docker",
				"tag",
				"quay.io/giantswarm/architect:master-e8363ac222255e991c126abe6673cd0f33934ac8",
				"quay.io/giantswarm/architect:master-latest",
			},
		},

		// Test 6, make sure NewDockerPushShaTask works as expected.
		{
			TaskFunc: func() (tasks.Task, error) {
				return NewDockerPushShaTask(nil, testNewProjectInfo())
			},
			ExpectedArgs: []string{
				"docker",
				"push",
				"quay.io/giantswarm/architect:master-e8363ac222255e991c126abe6673cd0f33934ac8",
			},
		},

		// Test 7, make sure NewDockerPushLatestTask works as expected.
		{
			TaskFunc: func() (tasks.Task, error) {
				return NewDockerPushLatestTask(nil, testNewProjectInfo())
			},
			ExpectedArgs: []string{
				"docker",
				"push",
				"quay.io/giantswarm/architect:master-latest",
			},
		},
	}

	for i, tc := range testCases {
		task, err := tc.TaskFunc()
		if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}

		execTask, ok := task.(tasks.ExecTask)
		if !ok {
			t.Fatalf("test %d expected %T got %T", i, tasks.ExecTask{}, task)
		}

		if !reflect.DeepEqual(execTask.Args, tc.ExpectedArgs) {
			t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedArgs, execTask.Args)
		}
	}
}

func testNewProjectInfo() ProjectInfo {
	return ProjectInfo{
		WorkingDirectory: "/usr/code/",
		Organisation:     "giantswarm",
		Project:          "architect",

		Branch: "master",
		Sha:    "e8363ac222255e991c126abe6673cd0f33934ac8",

		Registry:       "quay.io",
		DockerUsername: "username",
		DockerPassword: "password",

		Goos:          "test-goos",
		Goarch:        "test-goarch",
		GolangImage:   "test-golang-image",
		GolangVersion: "test-golang-version",
	}
}
