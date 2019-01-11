package workflow

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/template"
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
		Organisation:   "giantswarm",
		Project:        "test-project",
		Sha:            "jfkejhfkejfkejfef",
		Registry:       "quay.io",
		DockerUsername: "test",
		DockerPassword: "ekfnkfne",
		Goos:           "linux",
		Goarch:         "amd64",
		GolangImage:    "golang",
		GolangVersion:  "1.7.5",
	}

	tests := []struct {
		setUp             func(afero.Fs, string) error
		expectedTaskNames []string
		errorMatcher      func(error) bool
	}{
		// Test 0 that a project with no files produces a workflow with just the repo check.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "empty"), 0644); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
			},
		},

		// Test 1 that a project with only golang files produces a correct workflow.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := afero.WriteFile(fs, filepath.Join(projectInfo.WorkingDirectory, "main.go"), []byte("package test"), 0644); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				GoPullTaskName,
				strings.Join([]string{
					GoFmtTaskName,
					GoBuildTaskName,
					GoTestTaskName,
				}, ";") + ";",
			},
		},

		// Test 2 that a library project creates a correct workflow.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "client"), 0744); err != nil {
					return microerror.Mask(err)
				}
				if err := afero.WriteFile(fs, filepath.Join(projectInfo.WorkingDirectory, "client", "client.go"), []byte("package test"), 0644); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				GoPullTaskName,
				strings.Join([]string{
					GoFmtTaskName,
					GoTestTaskName,
				}, ";") + ";",
			},
		},

		// Test 3 that a particularly nested library project creates a correct workflow.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "server", "client"), 0744); err != nil {
					return microerror.Mask(err)
				}
				if err := afero.WriteFile(fs, filepath.Join(projectInfo.WorkingDirectory, "server", "client", "client.go"), []byte("package test"), 0644); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				GoPullTaskName,
				strings.Join([]string{
					GoFmtTaskName,
					GoTestTaskName,
				}, ";") + ";",
			},
		},

		// Test 4 that a project with a golang file not named `main.go` produces a
		// golang build workflow.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := afero.WriteFile(fs, filepath.Join(projectInfo.WorkingDirectory, "other.go"), []byte("package test"), 0644); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				GoPullTaskName,
				strings.Join([]string{
					GoFmtTaskName,
					GoBuildTaskName,
					GoTestTaskName,
				}, ";") + ";",
			},
		},

		// Test 5 that a project with only a dockerfile produces a correct workflow.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				DockerBuildTaskName,
				DockerLoginTaskName,
				DockerPushShaTaskName,
			},
		},

		// Test 6 that a project with golang files, and a dockerfile produces a correct
		// workflow.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := afero.WriteFile(fs, filepath.Join(projectInfo.WorkingDirectory, "main.go"), []byte("package test"), 0644); err != nil {
					return microerror.Mask(err)
				}
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				GoPullTaskName,
				strings.Join([]string{
					GoFmtTaskName,
					GoBuildTaskName,
					GoTestTaskName,
				}, ";") + ";",
				DockerBuildTaskName,
				DockerRunVersionTaskName,
				DockerRunHelpTaskName,
				DockerLoginTaskName,
				DockerPushShaTaskName,
			},
		},

		// Test 7 that a project with multiple helm charts has all of them pushed.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "test-project-chart"), 0744); err != nil {
					return microerror.Mask(err)
				}
				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "test-project-something-chart"), 0744); err != nil {
					return microerror.Mask(err)
				}
				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "test-project-another-chart"), 0744); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				HelmPullTaskName,
				HelmLoginTaskName,
				template.TemplateHelmChartTaskName,
				HelmPushTaskName,
				template.TemplateHelmChartTaskName,
				HelmPushTaskName,
				template.TemplateHelmChartTaskName,
				HelmPushTaskName,
			},
		},

		// Test 8 that a docker image is pushed before helm chart.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}
				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "test-project-some-chart"), 0744); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				DockerBuildTaskName,
				DockerLoginTaskName,
				DockerPushShaTaskName,
				HelmPullTaskName,
				HelmLoginTaskName,
				template.TemplateHelmChartTaskName,
				HelmPushTaskName,
			},
		},

		// Test 9 that charts not starting with a project name causes an error.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "some-chart"), 0744); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			errorMatcher: IsInvalidHelmDirectory,
		},

		// Test 10 that charts not starting with a project name causes an error.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "test-project-chart"), 0744); err != nil {
					return microerror.Mask(err)
				}
				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "some-chart"), 0744); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			errorMatcher: IsInvalidHelmDirectory,
		},

		// Test 11 that a project with only golang files that have build contraints
		// do not trigger a test workflow.
		{
			setUp: func(fs afero.Fs, testDir string) error {
				projectInfo.WorkingDirectory = testDir

				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "integration"), 0744); err != nil {
					return microerror.Mask(err)
				}
				if err := afero.WriteFile(fs, filepath.Join(projectInfo.WorkingDirectory, "integration", "integration_test.go"), []byte("// +build k8srequired\n\npackage test"), 0644); err != nil {
					return microerror.Mask(err)
				}
				return nil
			},
			expectedTaskNames: []string{
				RepoCheckTaskName,
				GoPullTaskName,
				strings.Join([]string{
					GoFmtTaskName,
				}, ";") + ";",
			},
		},
	}

	for i, tc := range tests {
		fs := afero.NewOsFs()

		// Create temporary testing dir and clean up after.
		testDir, err := afero.TempDir(fs, "", "test-dir")
		if err != nil {
			t.Fatalf("test %d received unexpected error creating a test directory: %#v", i, err)
		}
		defer fs.RemoveAll(testDir)

		if err := tc.setUp(fs, testDir); err != nil {
			t.Fatalf("test %d received unexpected error during setup: %#v", i, err)
		}

		workflow, err := NewBuild(projectInfo, fs)
		if err != nil && tc.errorMatcher != nil && tc.errorMatcher(err) {
			continue
		}
		if err == nil && tc.errorMatcher != nil {
			t.Fatalf("test %d: expected error, got %#v", i, err)
		}
		if err != nil && tc.errorMatcher == nil {
			t.Fatalf("test %d: unexpected error = %#v", i, err)
		}

		taskNames := []string{}
		for _, task := range workflow {
			retryTask, ok := task.(tasks.RetryTask)
			if ok {
				taskNames = append(taskNames, retryTask.Task.Name())
			} else {
				taskNames = append(taskNames, task.Name())
			}
		}

		if !reflect.DeepEqual(taskNames, tc.expectedTaskNames) {
			t.Fatalf("test %d expected %v tasks, got %v", i, tc.expectedTaskNames, taskNames)
		}
	}
}

// TestGetDeployWorkflow tests that deploy workflows are correctly created
func TestGetDeployWorkflow(t *testing.T) {
	projectInfo := ProjectInfo{
		WorkingDirectory: "/test-project/",
		Organisation:     "giantswarm",
		Project:          "test-project",
		Sha:              "1cd72a25e16e93da14f08d95bd98662f8827028e",
		Registry:         "quay.io",
		DockerUsername:   "test",
		DockerPassword:   "test",
	}

	pushTaskName := fmt.Sprintf("%s-%s", HelmPushTaskName, "stable")

	tests := []struct {
		setUp             func(afero.Fs) error
		expectedTaskNames []string
	}{
		// Test 0 a project with no files produces an empty workflow.
		{
			setUp: func(fs afero.Fs) error {
				return nil
			},
			expectedTaskNames: []string{},
		},

		// Test 1 a project with only a Dockerfile productes a workflow containg
		// docker push.
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: []string{
				DockerLoginTaskName,
				DockerPullTaskName,
				DockerTagLatestTaskName,
				DockerPushLatestTaskName,
			},
		},

		// Test 2 a project with a helm directory produces a workflow that does
		// nothing.
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.Mask(err)
				}
				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "helm/"), 0644); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: []string{
				DockerLoginTaskName,
				DockerPullTaskName,
				DockerTagLatestTaskName,
				DockerPushLatestTaskName,
			},
		},

		// Test 3 a project with a helm/PROJECT-chart directory
		// produces a workflow with helm push.
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm", "test-project-chart"), 0644); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: []string{
				HelmPullTaskName,
				HelmLoginTaskName,
				template.TemplateHelmChartTaskName,
				pushTaskName,
			},
		},

		// Test 4 a project with a helm/PROJECT-chart and
		// helm/PROJECT-something-chart directories pushes both charts
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm/test-project-chart"), 0644); err != nil {
					return microerror.Mask(err)
				}
				if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm/test-project-something-chart"), 0644); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			expectedTaskNames: []string{
				HelmPullTaskName,
				HelmLoginTaskName,
				template.TemplateHelmChartTaskName,
				pushTaskName,
				template.TemplateHelmChartTaskName,
				pushTaskName,
			},
		},
	}

	for i, tc := range tests {
		fs := afero.NewMemMapFs()
		if err := tc.setUp(fs); err != nil {
			t.Fatalf("test %d received unexpected error during setup: %#v", i, err)
		}

		workflow, err := NewDeploy(projectInfo, fs)
		if err != nil {
			t.Fatalf("test %d received unexpected error getting build workflow: %#v", i, err)
		}

		taskNames := []string{}
		for _, task := range workflow {
			retryTask, ok := task.(tasks.RetryTask)
			if ok {
				taskNames = append(taskNames, retryTask.Task.Name())
			} else {
				taskNames = append(taskNames, task.Name())
			}
		}

		if !reflect.DeepEqual(taskNames, tc.expectedTaskNames) {
			t.Fatalf("test %d expected %v tasks, got %v", i, tc.expectedTaskNames, taskNames)
		}
	}
}

func TestGetPublishWorkflow(t *testing.T) {
	projectInfo := ProjectInfo{
		WorkingDirectory: "/test-project/",
		Organisation:     "giantswarm",
		Project:          "test-project",
		Sha:              "1cd72a25e16e93da14f08d95bd98662f8827028e",
		Registry:         "quay.io",
		DockerUsername:   "test",
		DockerPassword:   "test",
	}

	setUp := func(fs afero.Fs) error {
		if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm/test-project-chart"), 0644); err != nil {
			return microerror.Mask(err)
		}
		return nil
	}

	tcs := []struct {
		description       string
		channels          []string
		expectedTaskNames []string
		expectedError     error
	}{
		{
			description: "default channels",
			channels:    []string{"beta", "testing"},
			expectedTaskNames: []string{
				HelmLoginTaskName,
				template.TemplateHelmChartTaskName,
				fmt.Sprintf("%s-beta", HelmPushTaskName),
				fmt.Sprintf("%s-testing", HelmPushTaskName),
			},
		},
		{
			description: "single channel",
			channels:    []string{"alpha"},
			expectedTaskNames: []string{
				HelmLoginTaskName,
				template.TemplateHelmChartTaskName,
				fmt.Sprintf("%s-alpha", HelmPushTaskName),
			},
		},
		{
			description: "multiple channels",
			channels:    []string{"alpha", "beta", "testing", "unstable"},
			expectedTaskNames: []string{
				HelmLoginTaskName,
				template.TemplateHelmChartTaskName,
				fmt.Sprintf("%s-alpha", HelmPushTaskName),
				fmt.Sprintf("%s-beta", HelmPushTaskName),
				fmt.Sprintf("%s-testing", HelmPushTaskName),
				fmt.Sprintf("%s-unstable", HelmPushTaskName),
			},
		},
		{
			description:       "error on empty channels",
			channels:          []string{"alpha", "beta", "", "unstable", ""},
			expectedTaskNames: []string{},
			expectedError:     emptyChannelError,
		},
	}

	for _, tc := range tcs {
		fs := afero.NewMemMapFs()
		t.Run(tc.description, func(t *testing.T) {
			if err := setUp(fs); err != nil {
				t.Errorf("received unexpected error during setup: %v", err)
			}

			projectInfo.Channels = tc.channels
			workflow, err := NewPublish(projectInfo, fs)
			if tc.expectedError != nil && err == nil {
				t.Errorf("expected error didn't happen")
			}
			if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("received unexpected error getting build workflow: %v", err)
			}

			taskNames := []string{}
			for _, task := range workflow {
				retryTask, ok := task.(tasks.RetryTask)
				if ok {
					taskNames = append(taskNames, retryTask.Task.Name())
				} else {
					taskNames = append(taskNames, task.Name())
				}
			}

			if !reflect.DeepEqual(taskNames, tc.expectedTaskNames) {
				t.Errorf("expected %v tasks, got %v", tc.expectedTaskNames, taskNames)
			}
		})
	}
}

func TestGetUnpublishWorkflow(t *testing.T) {
	projectInfo := ProjectInfo{
		WorkingDirectory: "/test-project/",
		Organisation:     "giantswarm",
		Project:          "test-project",
		Sha:              "1cd72a25e16e93da14f08d95bd98662f8827028e",
		Registry:         "quay.io",
		DockerUsername:   "test",
		DockerPassword:   "test",
	}

	setUp := func(fs afero.Fs) error {
		if err := fs.MkdirAll(filepath.Join(projectInfo.WorkingDirectory, "helm/test-project-chart"), 0644); err != nil {
			return microerror.Mask(err)
		}
		return nil
	}

	tcs := []struct {
		description       string
		channels          []string
		expectedTaskNames []string
		expectedError     error
	}{
		{
			description: "default channels",
			channels:    []string{"beta", "testing"},
			expectedTaskNames: []string{
				"helm-login",
				fmt.Sprintf("%s-beta", HelmDeleteChannelTaskName),
				fmt.Sprintf("%s-testing", HelmDeleteChannelTaskName),
			},
		},
		{
			description: "single channel",
			channels:    []string{"alpha"},
			expectedTaskNames: []string{
				"helm-login",
				fmt.Sprintf("%s-alpha", HelmDeleteChannelTaskName),
			},
		},
		{
			description: "multiple channels",
			channels:    []string{"alpha", "beta", "testing", "unstable"},
			expectedTaskNames: []string{
				"helm-login",
				fmt.Sprintf("%s-alpha", HelmDeleteChannelTaskName),
				fmt.Sprintf("%s-beta", HelmDeleteChannelTaskName),
				fmt.Sprintf("%s-testing", HelmDeleteChannelTaskName),
				fmt.Sprintf("%s-unstable", HelmDeleteChannelTaskName),
			},
		},
		{
			description:       "error on empty channels",
			channels:          []string{"alpha", "beta", "", "unstable", ""},
			expectedTaskNames: []string{},
			expectedError:     emptyChannelError,
		},
	}

	for _, tc := range tcs {
		fs := afero.NewMemMapFs()
		t.Run(tc.description, func(t *testing.T) {
			if err := setUp(fs); err != nil {
				t.Errorf("received unexpected error during setup: %v", err)
			}

			projectInfo.Channels = tc.channels
			workflow, err := NewUnpublish(projectInfo, fs)
			if tc.expectedError != nil && err == nil {
				t.Errorf("expected error didn't happen")
			}
			if err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("received unexpected error getting build workflow: %v", err)
			}

			taskNames := []string{}
			for _, task := range workflow {
				retryTask, ok := task.(tasks.RetryTask)
				if ok {
					taskNames = append(taskNames, retryTask.Task.Name())
				} else {
					taskNames = append(taskNames, task.Name())
				}
			}

			if !reflect.DeepEqual(taskNames, tc.expectedTaskNames) {
				t.Errorf("expected %v tasks, got %v", tc.expectedTaskNames, taskNames)
			}
		})
	}
}
