package workflow

import (
	"reflect"
	"strings"
	"testing"

	"github.com/giantswarm/architect/tasks"
	"github.com/spf13/afero"
)

func TestNewGoBuildTask(t *testing.T) {
	testCases := []struct {
		name         string
		projectInfo  ProjectInfo
		expectedArgs []string
		errorMatcher func(error) bool
	}{
		{
			name:         "case 0: empty projectInfo",
			projectInfo:  ProjectInfo{},
			expectedArgs: nil,
			errorMatcher: IsEmptyWorkingDirectory,
		},
		{
			name:        "case 1: test with values",
			projectInfo: testNewProjectInfo(),
			expectedArgs: []string{
				"docker",
				"run",
				"--rm",
				"-v",
				"/usr/code/:/go/src/github.com/giantswarm/architect",
				"-v",
				"/tmp/go/cache:/go/cache",
				"-e",
				"GOOS=test-goos",
				"-e",
				"GOARCH=test-goarch",
				"-e",
				"GOCACHE=/go/cache",
				"-e",
				"GOPATH=/go",
				"-e",
				"CGO_ENABLED=0",
				"-w",
				"/go/src/github.com/giantswarm/architect",
				"test-golang-image:test-golang-version",
				"go",
				"build",
				"-a",
				"-v",
				"-ldflags",
				strings.Join([]string{
					"-w",
					"-linkmode 'auto'",
					"-extldflags '-static'",
					"-X 'main.gitCommit=e8363ac222255e991c126abe6673cd0f33934ac8'",
					"-X 'github.com/giantswarm/architect/pkg/project.buildTimestamp=2019-06-04T12:40:05Z'",
					"-X 'github.com/giantswarm/architect/pkg/project.gitSHA=e8363ac222255e991c126abe6673cd0f33934ac8'",
				}, " "),
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFS := &afero.MemMapFs{}

			task, err := NewGoBuildTask(testFS, tc.projectInfo)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if err == nil {
				execTask, ok := task.(tasks.ExecTask)
				if !ok {
					t.Fatalf("wrong task type: got %T, want %T", task, tasks.ExecTask{})
				}

				if !reflect.DeepEqual(execTask.Args, tc.expectedArgs) {
					t.Fatalf("wrong expectedArgs: got %v, want %v", execTask.Args, tc.expectedArgs)
				}
			}
		})

	}

}
