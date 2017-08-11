package workflow

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"

	microerror "github.com/giantswarm/microkit/error"

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
		WorkingDirectory: "/test/",
		Organisation:     "giantswarm",
		Project:          "test",
		Sha:              "jfkejhfkejfkejfef",
		Registry:         "quay.io",
		DockerEmail:      "test@giantswarm.io",
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
		// Test a project with no files produces an empty workflow
		{
			setUp: func(fs afero.Fs) error {
				return nil
			},
			expectedTaskNames: map[int]string{},
		},

		// Test a project with only golang files produces a correct workflow
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "main.go")); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoFmtTaskName,
				1: GoVetTaskName,
				2: GoTestTaskName,
				3: GoBuildTaskName,
			},
		},

		// Test that a project with a golang file not named `main.go` produces a golang build workflow
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "other.go")); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoFmtTaskName,
				1: GoVetTaskName,
				2: GoTestTaskName,
				3: GoBuildTaskName,
			},
		},

		// Test a project with only a dockerfile produces a correct workflow
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: DockerBuildTaskName,
			},
		},

		// Test a project with golang files, and a dockerfile produces a correct workflow
		{
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "main.go")); err != nil {
					return microerror.MaskAny(err)
				}
				if _, err := fs.Create(filepath.Join(projectInfo.WorkingDirectory, "Dockerfile")); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: GoFmtTaskName,
				1: GoVetTaskName,
				2: GoTestTaskName,
				3: GoBuildTaskName,
				4: DockerBuildTaskName,
				5: DockerRunVersionTaskName,
				6: DockerRunHelpTaskName,
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

		if len(workflow) != len(test.expectedTaskNames) {
			t.Fatalf(
				"expected %v taskss, received %v",
				len(test.expectedTaskNames),
				len(workflow),
			)
		}

		for testIndex, expectedTaskName := range test.expectedTaskNames {
			if workflow[testIndex].Name() != expectedTaskName {
				t.Fatalf(
					"Task: %v, expected name: %v, received name: %v",
					index,
					expectedTaskName,
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
		projectInfo       ProjectInfo
		setUp             func(afero.Fs) error
		expectedTaskNames map[int]string
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
			expectedTaskNames: map[int]string{},
		},

		// Test a project with only a Dockerfile productes a workflow containg docker push
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				Sha:              "1cd72a25e16e93da14f08d95bd98662f8827028e",
				Registry:         "quay.io",
				DockerEmail:      "test@giantswarm.io",
				DockerUsername:   "test",
				DockerPassword:   "test",
			},
			setUp: func(fs afero.Fs) error {
				if _, err := fs.Create(filepath.Join(workingDirectory, "Dockerfile")); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: DockerLoginTaskName,
				1: DockerTagLatestTaskName,
				2: DockerPushShaTaskName,
				3: DockerPushLatestTaskName,
			},
		},

		// Test that a project with a helm directory produces a workflow containing helm push
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
			},
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(workingDirectory, "helm/"), 0644); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: template.TemplateHelmChartTaskName,
				1: HelmLoginTaskName,
				2: HelmPushTaskName,
			},
		},

		// Test a project with only a kubernetes directory produces a workflow containg kubernetes apply
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				KubernetesResourcesDirectoryPath: filepath.Join(workingDirectory, "/kubernetes/"),
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
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: template.TemplateKubernetesResourcesTaskName,
				1: KubectlClusterInfoTaskName,
				2: KubectlApplyTaskName,
			},
		},

		// Test a project with a Dockerfile and a kubernetes directory contains both docker and kubernetes tasks/
		{
			projectInfo: ProjectInfo{
				WorkingDirectory:                 workingDirectory,
				Organisation:                     "giantswarm",
				Project:                          "test",
				Sha:                              "1cd72a25e16e93da14f08d95bd98662f8827028e",
				Registry:                         "quay.io",
				DockerEmail:                      "test@giantswarm.io",
				DockerUsername:                   "test",
				DockerPassword:                   "test",
				KubernetesResourcesDirectoryPath: filepath.Join(workingDirectory, "/kubernetes/"),
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
					return microerror.MaskAny(err)
				}
				if err := fs.Mkdir(filepath.Join(workingDirectory, "kubernetes/"), 0644); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: DockerLoginTaskName,
				1: DockerTagLatestTaskName,
				2: DockerPushShaTaskName,
				3: DockerPushLatestTaskName,
				4: template.TemplateKubernetesResourcesTaskName,
				5: KubectlClusterInfoTaskName,
				6: KubectlApplyTaskName,
			},
		},

		// Test that a project with two clusters configured returns two sets of kubectl tasks.
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				KubernetesResourcesDirectoryPath: filepath.Join(workingDirectory, "/kubernetes/"),
				KubernetesClusters: []KubernetesCluster{
					KubernetesCluster{
						ApiServer:      "kubernetes-1.giantswarm.io",
						Prefix:         "1",
						CaPath:         "/1-ca.pem",
						CrtPath:        "/1-crt.pem",
						KeyPath:        "/1-key.pem",
						KubectlVersion: "1.5.2",
					},
					KubernetesCluster{
						ApiServer:      "kubernetes-2.giantswarm.io",
						Prefix:         "2",
						CaPath:         "/2-ca.pem",
						CrtPath:        "/2-crt.pem",
						KeyPath:        "/2-key.pem",
						KubectlVersion: "1.5.2",
					},
				},
			},
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(workingDirectory, "kubernetes/"), 0644); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: template.TemplateKubernetesResourcesTaskName,
				1: KubectlClusterInfoTaskName,
				2: KubectlApplyTaskName,
				3: template.TemplateKubernetesResourcesTaskName,
				4: KubectlClusterInfoTaskName,
				5: KubectlApplyTaskName,
			},
		},

		// Test a project with a non-standard kubernetes resources directory
		{
			projectInfo: ProjectInfo{
				WorkingDirectory: workingDirectory,
				Organisation:     "giantswarm",
				Project:          "test",
				KubernetesResourcesDirectoryPath: filepath.Join(workingDirectory, "/something-different"),
				KubernetesClusters: []KubernetesCluster{
					KubernetesCluster{
						ApiServer:      "kubernetes-1.giantswarm.io",
						Prefix:         "1",
						CaPath:         "/1-ca.pem",
						CrtPath:        "/1-crt.pem",
						KeyPath:        "/1-key.pem",
						KubectlVersion: "1.5.2",
					},
				},
			},
			setUp: func(fs afero.Fs) error {
				if err := fs.Mkdir(filepath.Join(workingDirectory, "something-different"), 0644); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			expectedTaskNames: map[int]string{
				0: template.TemplateKubernetesResourcesTaskName,
				1: KubectlClusterInfoTaskName,
				2: KubectlApplyTaskName,
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

		if len(workflow) != len(test.expectedTaskNames) {
			t.Fatalf(
				"%v: expected %v tasks, received %v",
				index,
				len(test.expectedTaskNames),
				len(workflow),
			)
		}

		for testIndex, expectedTaskName := range test.expectedTaskNames {
			if !strings.Contains(workflow[testIndex].Name(), expectedTaskName) {
				t.Fatalf(
					"%s: task: %s, expected name: %s, received name: %s",
					index,
					testIndex,
					expectedTaskName,
					workflow[testIndex].Name(),
				)
			}
		}
	}
}
