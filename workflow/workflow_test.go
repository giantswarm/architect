package workflow

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/architect/tasks"
)

func TestWorkflowString(t *testing.T) {
	tests := []struct {
		workflow       Workflow
		expectedString string
	}{
		// Test the empty workflow
		{
			workflow:       Workflow{},
			expectedString: "{}",
		},

		// Test one task
		{
			workflow: Workflow{
				tasks.NewExecTask(
					"foo",
					[]string{"apple", "banana"},
				),
			},
			expectedString: "{\n\tfoo:\t'apple banana'\n}",
		},

		// Test multiple tasks
		{
			workflow: Workflow{
				tasks.NewExecTask(
					"foo",
					[]string{"apple", "banana"},
				),
				tasks.NewExecTask(
					"bar",
					[]string{"cherry", "durian"},
				),
			},
			expectedString: "{\n\tfoo:\t'apple banana'\n\tbar:\t'cherry durian'\n}",
		},
	}

	for index, test := range tests {
		returnedString := fmt.Sprintf("%s", test.workflow)
		if returnedString != test.expectedString {
			t.Fatalf(
				"%v: returned string did not match expected string.\nexpected: %v\nreturned: %v\n",
				index,
				test.expectedString,
				returnedString,
			)
		}
	}
}

// TestGetBuildWorkflow tests that build workflows are correctly created for builds
func TestGetBuildWorkflow(t *testing.T) {
	projectInfo := ProjectInfo{
		WorkingDirectory: "/test/",
		Organisation:     "giantswarm",
		Project:          "test",
		Sha:              "jfkejhfkejfkejfef",
		Registry:         "quay.io",
		DockerUsername:   "test",
		DockerPassword:   "ekfnkfne",
		Goos:             "linux",
		Goarch:           "amd64",
		GolangImage:      "golang",
		GolangVersion:    "1.7.5",
	}

	tests := []struct {
		setUp             func(afero.Fs) error
		expectedTaskNames map[int]string
	}{
		// 0: Test that a project with no files produces an empty workflow.
		{
			setUp: func(fs afero.Fs) error {
				return nil
			},
			expectedTaskNames: map[int]string{},
		},

		// 1: Test that a project with only golang files produces a correct workflow.
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "main.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoPullTaskName,
				1: GoFmtTaskName,
				2: GoVetTaskName,
				3: GoTestTaskName,
				4: GoBuildTaskName,
			},
		},

		// 2: Test that a library project creates a correct workflow.
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "client"), 0644); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "client", "client.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoPullTaskName,
				1: GoFmtTaskName,
				2: GoVetTaskName,
				3: GoTestTaskName,
			},
		},

		// 3: Test that a particularly nested library project creates a correct workflow.
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "server"), 0644); err != nil {
					return microerror.Mask(err)
				}

				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "server", "client"), 0644); err != nil {
					return microerror.Mask(err)
				}

				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "server", "client", "client.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoPullTaskName,
				1: GoFmtTaskName,
				2: GoVetTaskName,
				3: GoTestTaskName,
			},
		},

		// 4: Test that a project with a golang file not named `main.go` produces a
		// golang build workflow.
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "other.go")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoPullTaskName,
				1: GoFmtTaskName,
				2: GoVetTaskName,
				3: GoTestTaskName,
				4: GoBuildTaskName,
			},
		},

		// 5: Test that a project with only a dockerfile produces a correct workflow.
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: DockerBuildTaskName,
				1: DockerLoginTaskName,
				2: DockerPushShaTaskName,
			},
		},

		// 6: Test that a project with golang files, and a dockerfile produces a correct
		// workflow.
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "main.go")); err != nil {
					return microerror.Mask(err)
				}
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoPullTaskName,
				1: GoFmtTaskName,
				2: GoVetTaskName,
				3: GoTestTaskName,
				4: GoBuildTaskName,
				5: DockerBuildTaskName,
				6: DockerRunVersionTaskName,
				7: DockerRunHelpTaskName,
				8: DockerLoginTaskName,
				9: DockerPushShaTaskName,
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		if err := test.setUp(fs); err != nil {
			t.Fatalf("test %d received unexpected error during setup: %#v", index, err)
		}

		workflow, err := NewBuild(projectInfo, fs)
		if err != nil {
			t.Fatalf("test %d received unexpected error getting build workflow: %#v", index, err)
		}

		if len(workflow) != len(test.expectedTaskNames) {
			t.Fatalf("test %d expected %d tasks, received %#v", index, len(test.expectedTaskNames), len(workflow))
		}

		for testIndex, expectedTaskName := range test.expectedTaskNames {
			if !strings.Contains(workflow[testIndex].Name(), expectedTaskName) {
				t.Fatalf(
					"test %d task %d expected name '%s' received name '%s'",
					index,
					testIndex,
					expectedTaskName,
					workflow[testIndex].Name(),
				)
			}
		}
	}
}

// TestGetDeployWorkflow tests that deploy workflows are correctly created
func TestGetDeployWorkflow(t *testing.T) {
	workingDirectory := "/test/"

	tests := []struct {
		projectInfo       ProjectInfo
		setUp             func(afero.Fs) error
		expectedTaskNames map[int]string
	}{
		// Test 1 a project with no files produces an empty workflow.
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
			},
			setUp: func(fs afero.Fs) error {
				return nil
			},
			expectedTaskNames: map[int]string{},
		},

		// Test 2 a project with only a Dockerfile productes a workflow containg
		// docker push.
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				Sha:              "1cd72a25e16e93da14f08d95bd98662f8827028e",
				Registry:         "quay.io",
				DockerUsername:   "test",
				DockerPassword:   "test",
			},
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(workingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: DockerLoginTaskName,
				1: DockerTagLatestTaskName,
				2: DockerPushLatestTaskName,
			},
		},

		// Test 3 a project with a helm directory produces a workflow that does
		// nothing.
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
			},
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(workingDirectory, "helm/"), 0644); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		if err := test.setUp(fs); err != nil {
			t.Fatalf("test %d received unexpected error during setup: %#v", index+1, err)
		}

		workflow, err := NewDeploy(test.projectInfo, fs)
		if err != nil {
			t.Fatalf("test %d received unexpected error getting build workflow: %#v", index+1, err)
		}

		if len(workflow) != len(test.expectedTaskNames) {
			t.Fatalf("test %d expected %d tasks received %d", index+1, len(test.expectedTaskNames), len(workflow))
		}

		for testIndex, expectedTaskName := range test.expectedTaskNames {
			if !strings.Contains(workflow[testIndex].Name(), expectedTaskName) {
				t.Fatalf(
					"test %d task %d expected name '%s' received name '%s'",
					index+1,
					testIndex,
					expectedTaskName,
					workflow[testIndex].Name(),
				)
			}
		}
	}
}
