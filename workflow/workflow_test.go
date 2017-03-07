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
				0: "go-test",
				1: "go-build",
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
				0: "docker-build",
				1: "docker-run-version",
				2: "docker-run-help",
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
				0: "go-test",
				1: "go-build",
				2: "docker-build",
				3: "docker-run-version",
				4: "docker-run-help",
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		if err := test.setUp(fs); err != nil {
			t.Fatalf("received unexpected error during setup: %v", err)
		}

		workflow, err := GetBuildWorkflow(projectInfo, fs)
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
				0: "docker-login",
				1: "docker-push",
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
				0: "kubectl-cluster-info",
				1: "kubectl-apply",
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
				0: "docker-login",
				1: "docker-push",
				2: "kubectl-cluster-info",
				3: "kubectl-apply",
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		if err := test.setUp(fs); err != nil {
			t.Fatalf("received unexpected error during setup: %v", err)
		}

		workflow, err := GetDeployWorkflow(projectInfo, fs)
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
