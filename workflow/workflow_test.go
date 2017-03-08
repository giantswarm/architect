package workflow

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

// TestGetBuildWorkflow tests that build workflows are correctly created for builds
func TestGetBuildWorkflow(t *testing.T) {
	projectInfo := ProjectInfo{
		WorkingDirectory: "/test/",
		Organisation:     "giantswarm",
		Project:          "test",
		Sha:              "jfkejhfkejfkejfef",
		Registry:         "registry.giantswarm.io",
		DockerEmail:      "test@giantswarm.io",
		DockerUsername:   "test",
		DockerPassword:   "ekfnkfne",
		Goos:             "linux",
		Goarch:           "amd64",
		GolangImage:      "golang",
		GolangVersion:    "1.7.5",
	}

	tests := []struct {
		setUp                func(afero.Fs) error
		expectedCommandNames map[int]string
	}{
		// Test a project with no files produces an empty workflow
		{
			setUp: func(fs afero.Fs) error {
				return nil
			},
			expectedCommandNames: map[int]string{},
		},

		// Test a project with only golang files produces a correct workflow
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "main.go")); err != nil {
					return err
				}

				return nil
			},
			expectedCommandNames: map[int]string{
				0: GoTestCommandName,
				1: GoBuildCommandName,
			},
		},

		// Test a project with only a dockerfile produces a correct workflow
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return err
				}

				return nil
			},
			expectedCommandNames: map[int]string{
				0: DockerBuildCommandName,
				1: DockerRunVersionCommandName,
				2: DockerRunHelpCommandName,
			},
		},

		// Test a project with golang files, and a dockerfile produces a correct workflow
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "main.go")); err != nil {
					return err
				}
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return err
				}

				return nil
			},
			expectedCommandNames: map[int]string{
				0: GoTestCommandName,
				1: GoBuildCommandName,
				2: DockerBuildCommandName,
				3: DockerRunVersionCommandName,
				4: DockerRunHelpCommandName,
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		if err := test.setUp(fs); err != nil {
			t.Fatalf("received unexpected error during setup: %v", err)
		}

		workflow, err := NewBuild(projectInfo, fs)
		if err != nil {
			t.Fatalf("received unexpected error getting build workflow: %v", err)
		}

		if len(workflow) != len(test.expectedCommandNames) {
			t.Fatalf(
				"expected %v commands, received %v",
				len(test.expectedCommandNames),
				len(workflow),
			)
		}

		for testIndex, expectedCommandName := range test.expectedCommandNames {
			if workflow[testIndex].Name != expectedCommandName {
				t.Fatalf(
					"command: %v, expected name: %v, received name: %v",
					index,
					expectedCommandName,
					workflow[index].Name,
				)
			}
		}
	}
}

// TestGetDeployWorkflow tests that deploy workflows are correctly created
func TestGetDeployWorkflow(t *testing.T) {
	projectInfo := ProjectInfo{
		WorkingDirectory:                          "/test/",
		Organisation:                              "giantswarm",
		Project:                                   "test",
		Sha:                                       "jfkejhfkejfkejfef",
		Registry:                                  "registry.giantswarm.io",
		DockerEmail:                               "test@giantswarm.io",
		DockerUsername:                            "test",
		DockerPassword:                            "ekfnkfne",
		KubernetesApiServer:                       "kubernetes.giantswarm.io",
		KubernetesCaPath:                          "/ca.pem",
		KubernetesCrtPath:                         "/crt.pem",
		KubernetesKeyPath:                         "/key.pem",
		KubectlVersion:                            "1.5.2",
		KubernetesTemplatedResourcesDirectoryPath: "/kubernetes-templated/",
	}

	tests := []struct {
		setUp                func(afero.Fs) error
		expectedCommandNames map[int]string
	}{
		// Test a project with no files produces an empty workflow
		{
			setUp: func(fs afero.Fs) error {
				return nil
			},
			expectedCommandNames: map[int]string{},
		},

		// Test a project with only a Dockerfile productes a workflow containg docker push
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return err
				}

				return nil
			},
			expectedCommandNames: map[int]string{
				0: DockerLoginCommandName,
				1: DockerPushCommandName,
			},
		},

		// Test a project with only a kubernetes directory produces a workflow containg kubernetes apply
		{
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "kubernetes/"), 0644); err != nil {
					return err
				}

				return nil
			},
			expectedCommandNames: map[int]string{
				0: KubectlClusterInfoCommandName,
				1: KubectlApplyCommandName,
			},
		},

		// Test a project with a Dockerfile and a kubernetes directory contains both docker and kubernetes commands
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return err
				}
				if err := fs.Mkdir(filepath.Join(projectInfo.WorkingDirectory, "kubernetes/"), 0644); err != nil {
					return err
				}

				return nil
			},
			expectedCommandNames: map[int]string{
				0: DockerLoginCommandName,
				1: DockerPushCommandName,
				2: KubectlClusterInfoCommandName,
				3: KubectlApplyCommandName,
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		if err := test.setUp(fs); err != nil {
			t.Fatalf("received unexpected error during setup: %v", err)
		}

		workflow, err := NewDeploy(projectInfo, fs)
		if err != nil {
			t.Fatalf("received unexpected error getting build workflow: %v", err)
		}

		if len(workflow) != len(test.expectedCommandNames) {
			t.Fatalf(
				"expected %v commands, received %v",
				len(test.expectedCommandNames),
				len(workflow),
			)
		}

		for testIndex, expectedCommandName := range test.expectedCommandNames {
			if workflow[testIndex].Name != expectedCommandName {
				t.Fatalf(
					"command: %v, expected name: %v, received name: %v",
					index,
					expectedCommandName,
					workflow[index].Name,
				)
			}
		}
	}
}
