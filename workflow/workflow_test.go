package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/giantswarm/architect/commands"
	"github.com/spf13/afero"
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

		// Test one command
		{
			workflow: Workflow{
				commands.Command{
					Name: "foo",
					Args: []string{"apple", "banana"},
				},
			},
			expectedString: "{\n\tfoo:\t'apple banana'\n}",
		},

		// Test multiple commands
		{
			workflow: Workflow{
				commands.Command{
					Name: "foo",
					Args: []string{"apple", "banana"},
				},
				commands.Command{
					Name: "bar",
					Args: []string{"cherry", "durian"},
				},
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
	workingDirectory := "/test/"

	tests := []struct {
		projectInfo          ProjectInfo
		setUp                func(afero.Fs) error
		expectedCommandNames map[int]string
	}{
		// Test a project with no files produces an empty workflow
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
			},
			setUp: func(fs afero.Fs) error {
				return nil
			},
			expectedCommandNames: map[int]string{},
		},

		// Test a project with only a Dockerfile productes a workflow containg docker push
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				Sha:              "1cd72a25e16e93da14f08d95bd98662f8827028e",
				Registry:         "registry.giantswarm.io",
				DockerEmail:      "test@giantswarm.io",
				DockerUsername:   "test",
				DockerPassword:   "test",
			},
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(workingDirectory, "Dockerfile")); err != nil {
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
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				KubernetesTemplatedResourcesDirectoryPath: "/kubernetes/",
				KubernetesClusters: []KubernetesCluster{
					KubernetesCluster{
						ApiServer:      "kubernetes.giantswarm.io",
						CaPath:         "/ca.pem",
						CrtPath:        "/crt.pem",
						KeyPath:        "/key.pem",
						KubectlVersion: "1.5.2",
					},
				},
			},
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(workingDirectory, "kubernetes/"), 0644); err != nil {
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
			projectInfo: ProjectInfo{
				WorkingDirectory:                          workingDirectory,
				Organisation:                              "giantswarm",
				Project:                                   "test",
				Sha:                                       "1cd72a25e16e93da14f08d95bd98662f8827028e",
				Registry:                                  "registry.giantswarm.io",
				DockerEmail:                               "test@giantswarm.io",
				DockerUsername:                            "test",
				DockerPassword:                            "test",
				KubernetesTemplatedResourcesDirectoryPath: "/kubernetes/",
				KubernetesClusters: []KubernetesCluster{
					KubernetesCluster{
						ApiServer:      "kubernetes.giantswarm.io",
						CaPath:         "/ca.pem",
						CrtPath:        "/crt.pem",
						KeyPath:        "/key.pem",
						KubectlVersion: "1.5.2",
					},
				},
			},
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(workingDirectory, "Dockerfile")); err != nil {
					return err
				}
				if err := fs.Mkdir(filepath.Join(workingDirectory, "kubernetes/"), 0644); err != nil {
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

		// Test that a project with two clusters configured returns two sets of kubectl commands
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				KubernetesTemplatedResourcesDirectoryPath: "/kubernetes/",
				KubernetesClusters: []KubernetesCluster{
					KubernetesCluster{
						ApiServer:      "kubernetes-1.giantswarm.io",
						CaPath:         "/1-ca.pem",
						CrtPath:        "/1-crt.pem",
						KeyPath:        "/1-key.pem",
						KubectlVersion: "1.5.2",
					},
					KubernetesCluster{
						ApiServer:      "kubernetes-2.giantswarm.io",
						CaPath:         "/2-ca.pem",
						CrtPath:        "/2-crt.pem",
						KeyPath:        "/2-key.pem",
						KubectlVersion: "1.5.2",
					},
				},
			},
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(workingDirectory, "kubernetes/"), 0644); err != nil {
					return err
				}

				return nil
			},
			expectedCommandNames: map[int]string{
				0: KubectlClusterInfoCommandName,
				1: KubectlApplyCommandName,
				2: KubectlClusterInfoCommandName,
				3: KubectlApplyCommandName,
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		if err := test.setUp(fs); err != nil {
			t.Fatalf("%v: received unexpected error during setup: %v", index, err)
		}

		workflow, err := NewDeploy(test.projectInfo, fs)
		if err != nil {
			t.Fatalf("%v: received unexpected error getting build workflow: %v", index, err)
		}

		if len(workflow) != len(test.expectedCommandNames) {
			t.Fatalf(
				"%v: expected %v commands, received %v",
				index,
				len(test.expectedCommandNames),
				len(workflow),
			)
		}

		for testIndex, expectedCommandName := range test.expectedCommandNames {
			if workflow[testIndex].Name != expectedCommandName {
				t.Fatalf(
					"%v: command: %v, expected name: %v, received name: %v",
					index,
					testIndex,
					expectedCommandName,
					workflow[index].Name,
				)
			}
		}
	}
}

func TestClustersFromEnv(t *testing.T) {
	type envVar struct {
		key   string
		value string
	}

	certTestData := "test"
	certEncodedTestData := "dGVzdA=="

	tests := []struct {
		workingDirectory string
		envVars          []envVar
		clusters         []KubernetesCluster
		expectedFiles    []string
	}{
		// Test that an empty env var set returns no clusters.
		{
			workingDirectory: "/test/",
			envVars:          []envVar{},
			clusters:         []KubernetesCluster{},
			expectedFiles:    []string{},
		},

		// Test that a few G8S env vars returns no clusters.
		{
			workingDirectory: "/test/",
			envVars: []envVar{
				envVar{key: "G8S_CA", value: "test"},
			},
			clusters:      []KubernetesCluster{},
			expectedFiles: []string{},
		},

		// Test that G8S env vars returns the LW cluster
		{
			workingDirectory: "/test/",
			envVars: []envVar{
				envVar{key: "G8S_CA", value: certEncodedTestData},
				envVar{key: "G8S_CRT", value: certEncodedTestData},
				envVar{key: "G8S_KEY", value: certEncodedTestData},
			},
			clusters: []KubernetesCluster{
				KubernetesCluster{
					ApiServer:      "https://api.g8s.fra-1.giantswarm.io",
					CaPath:         "/test/g8s-ca.pem",
					CrtPath:        "/test/g8s-crt.pem",
					KeyPath:        "/test/g8s-key.pem",
					KubectlVersion: "1.4.7",
				},
			},
			expectedFiles: []string{
				"/test/g8s-ca.pem",
				"/test/g8s-crt.pem",
				"/test/g8s-key.pem",
			},
		},

		// Test that AWS env vars returns the AWS cluster
		{
			workingDirectory: "/test",
			envVars: []envVar{
				envVar{key: "AWS_CA", value: certEncodedTestData},
				envVar{key: "AWS_CRT", value: certEncodedTestData},
				envVar{key: "AWS_KEY", value: certEncodedTestData},
			},
			clusters: []KubernetesCluster{
				KubernetesCluster{
					ApiServer:      "https://api.g8s.eu-west-1.aws.adidas.private.giantswarm.io:6443",
					CaPath:         "/test/aws-ca.pem",
					CrtPath:        "/test/aws-crt.pem",
					KeyPath:        "/test/aws-key.pem",
					KubectlVersion: "1.4.7",
				},
			},
			expectedFiles: []string{
				"/test/aws-ca.pem",
				"/test/aws-crt.pem",
				"/test/aws-key.pem",
			},
		},

		// Test that G8S and AWS env vars return both clusters
		{
			workingDirectory: "/test",
			envVars: []envVar{
				envVar{key: "G8S_CA", value: certEncodedTestData},
				envVar{key: "G8S_CRT", value: certEncodedTestData},
				envVar{key: "G8S_KEY", value: certEncodedTestData},

				envVar{key: "AWS_CA", value: certEncodedTestData},
				envVar{key: "AWS_CRT", value: certEncodedTestData},
				envVar{key: "AWS_KEY", value: certEncodedTestData},
			},
			clusters: []KubernetesCluster{
				KubernetesCluster{
					ApiServer:      "https://api.g8s.fra-1.giantswarm.io",
					CaPath:         "/test/g8s-ca.pem",
					CrtPath:        "/test/g8s-crt.pem",
					KeyPath:        "/test/g8s-key.pem",
					KubectlVersion: "1.4.7",
				},
				KubernetesCluster{
					ApiServer:      "https://api.g8s.eu-west-1.aws.adidas.private.giantswarm.io:6443",
					CaPath:         "/test/aws-ca.pem",
					CrtPath:        "/test/aws-crt.pem",
					KeyPath:        "/test/aws-key.pem",
					KubectlVersion: "1.4.7",
				},
			},
			expectedFiles: []string{
				"/test/g8s-ca.pem",
				"/test/g8s-crt.pem",
				"/test/g8s-key.pem",

				"/test/aws-ca.pem",
				"/test/aws-crt.pem",
				"/test/aws-key.pem",
			},
		},
	}

	for index, test := range tests {
		for _, envVar := range test.envVars {
			if err := os.Setenv(envVar.key, envVar.value); err != nil {
				t.Fatalf("%v: unexpected error setting env var: %v", index, err)
			}
		}

		fs := afero.NewMemMapFs()

		clusters, err := ClustersFromEnv(fs, test.workingDirectory)

		for _, envVar := range test.envVars {
			if err := os.Setenv(envVar.key, ""); err != nil {
				t.Fatalf("%v: unexpected error unsetting env var: %v", index, err)
			}
		}

		if err != nil {
			t.Fatalf("%v: unexpected error getting clusters from env: %v", index, err)
		}

		if !reflect.DeepEqual(clusters, test.clusters) {
			t.Fatalf(
				"%v: expected clusters did not match returned clusters.\nexpected: %#v\nreturned: %#v\n",
				index,
				test.clusters,
				clusters,
			)
		}

		for _, expectedFile := range test.expectedFiles {
			_, err := fs.Stat(expectedFile)
			if err != nil {
				t.Fatalf("%v: unexpected error checking certificate: %v", index, err)
			}

			contents, err := afero.ReadFile(fs, expectedFile)
			if err != nil {
				t.Fatalf("%v: unexpected error checking certificate contents: %v", index, err)
			}
			if string(contents) != certTestData {
				t.Fatalf("%v: certificate did not match expected contents: %v", index, string(contents))
			}
		}
	}
}
