package workflow

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"

	"github.com/spf13/afero"
)

var (
	GoTestCommandName  = "go-test"
	GoBuildCommandName = "go-build"

	DockerBuildCommandName      = "docker-build"
	DockerRunVersionCommandName = "docker-run-version"
	DockerRunHelpCommandName    = "docker-run-help"

	DockerLoginCommandName = "docker-login"
	DockerPushCommandName  = "docker-push"

	KubectlClusterInfoCommandName = "kubectl-cluster-info"
	KubectlApplyCommandName       = "kubectl-apply"
)

type Workflow []commands.Command

func (w Workflow) String() string {
	if len(w) == 0 {
		return "{}"
	}

	cmdStrings := []string{}
	for _, cmd := range w {
		cmdStrings = append(cmdStrings, "\t"+cmd.String()+"\n")
	}

	return fmt.Sprintf("{\n%v}", strings.Join(cmdStrings, ""))
}

type KubernetesCluster struct {
	ApiServer      string
	CaPath         string
	CrtPath        string
	KeyPath        string
	KubectlVersion string
}

type ProjectInfo struct {
	WorkingDirectory string
	Organisation     string
	Project          string
	Sha              string

	Registry       string
	DockerEmail    string
	DockerUsername string
	DockerPassword string

	KubernetesTemplatedResourcesDirectoryPath string
	KubernetesClusters                        []KubernetesCluster

	Goos          string
	Goarch        string
	GolangImage   string
	GolangVersion string
}

func NewBuild(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	if projectInfo.WorkingDirectory == "" {
		return nil, fmt.Errorf("working directory cannot be empty")
	}
	if projectInfo.Organisation == "" {
		return nil, fmt.Errorf("organisation cannot be empty")
	}
	if projectInfo.Project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	mainGoExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "main.go"))
	if err != nil {
		return nil, err
	}
	if mainGoExists {
		if projectInfo.Goos == "" {
			return nil, fmt.Errorf("goos cannot be empty")
		}
		if projectInfo.Goarch == "" {
			return nil, fmt.Errorf("goarch cannot be empty")
		}
		if projectInfo.GolangImage == "" {
			return nil, fmt.Errorf("golang image cannot be empty")
		}
		if projectInfo.GolangVersion == "" {
			return nil, fmt.Errorf("golang version cannot be empty")
		}

		testPackageArguments, err := utils.NoVendor(projectInfo.WorkingDirectory, fs)
		if err != nil {
			return nil, err
		}

		goTest := commands.NewDockerCommand(
			GoTestCommandName,
			commands.DockerCommandConfig{
				Volumes: []string{
					fmt.Sprintf(
						"%v:/go/src/github.com/%v/%v",
						projectInfo.WorkingDirectory,
						projectInfo.Organisation,
						projectInfo.Project,
					),
				},
				Env: []string{
					fmt.Sprintf("GOOS=%v", projectInfo.Goos),
					fmt.Sprintf("GOARCH=%v", projectInfo.Goarch),
					"GOPATH=/go",
					"CGOENABLED=0",
				},
				WorkingDirectory: fmt.Sprintf(
					"/go/src/github.com/%v/%v",
					projectInfo.Organisation,
					projectInfo.Project,
				),
				Image: fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
				Args:  []string{"go", "test", "-v"},
			},
		)
		goTest.Args = append(goTest.Args, testPackageArguments...)
		w = append(w, goTest)

		goBuild := commands.NewDockerCommand(
			GoBuildCommandName,
			commands.DockerCommandConfig{
				Volumes: []string{
					fmt.Sprintf(
						"%v:/go/src/github.com/%v/%v",
						projectInfo.WorkingDirectory,
						projectInfo.Organisation,
						projectInfo.Project,
					),
				},
				Env: []string{
					fmt.Sprintf("GOOS=%v", projectInfo.Goos),
					fmt.Sprintf("GOARCH=%v", projectInfo.Goarch),
					"GOPATH=/go",
					"CGOENABLED=0",
				},
				WorkingDirectory: fmt.Sprintf(
					"/go/src/github.com/%v/%v",
					projectInfo.Organisation,
					projectInfo.Project,
				),
				Image: fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
				Args:  []string{"go", "build", "-v", "-a", "-tags", "netgo", "-ldflags", "-linkmode 'external' -extldflags '-static'"},
			},
		)
		w = append(w, goBuild)
	}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, err
	}
	if dockerFileExists {
		if projectInfo.Registry == "" {
			return nil, fmt.Errorf("registry cannot be empty")
		}
		if projectInfo.Sha == "" {
			return nil, fmt.Errorf("sha cannot be empty")
		}

		dockerBuild := commands.Command{
			Name: DockerBuildCommandName,
			Args: []string{
				"docker",
				"build",
				"-t",
				fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
				projectInfo.WorkingDirectory,
			},
		}
		w = append(w, dockerBuild)

		dockerRunVersion := commands.NewDockerCommand(
			DockerRunVersionCommandName,
			commands.DockerCommandConfig{
				Image: fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
				Args:  []string{"version"},
			},
		)
		w = append(w, dockerRunVersion)

		dockerRunHelp := commands.NewDockerCommand(
			DockerRunHelpCommandName,
			commands.DockerCommandConfig{
				Image: fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
				Args:  []string{"--help"},
			},
		)
		w = append(w, dockerRunHelp)
	}

	return w, nil
}

func NewDeploy(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	if projectInfo.WorkingDirectory == "" {
		return nil, fmt.Errorf("working directory cannot be empty")
	}
	if projectInfo.Organisation == "" {
		return nil, fmt.Errorf("organisation cannot be empty")
	}
	if projectInfo.Project == "" {
		return nil, fmt.Errorf("project cannot be empty")
	}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, err
	}
	if dockerFileExists {
		if projectInfo.Sha == "" {
			return nil, fmt.Errorf("sha cannot be empty")
		}
		if projectInfo.Registry == "" {
			return nil, fmt.Errorf("registry cannot be empty")
		}
		if projectInfo.DockerEmail == "" {
			return nil, fmt.Errorf("docker email cannot be empty")
		}
		if projectInfo.DockerUsername == "" {
			return nil, fmt.Errorf("docker username cannot be empty")
		}
		if projectInfo.DockerPassword == "" {
			return nil, fmt.Errorf("docker password cannot be empty")
		}

		dockerLogin := commands.Command{
			Name: DockerLoginCommandName,
			Args: []string{
				"docker",
				"login",
				fmt.Sprintf("--email=%v", projectInfo.DockerEmail),
				fmt.Sprintf("--username=%v", projectInfo.DockerUsername),
				fmt.Sprintf("--password=%v", projectInfo.DockerPassword),
				projectInfo.Registry,
			},
		}
		w = append(w, dockerLogin)

		dockerPush := commands.Command{
			Name: DockerPushCommandName,
			Args: []string{
				"docker",
				"push",
				fmt.Sprintf("%v/%v/%v:%v", projectInfo.Registry, projectInfo.Organisation, projectInfo.Project, projectInfo.Sha),
			},
		}
		w = append(w, dockerPush)
	}

	kubernetesDirectoryExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "kubernetes/"))
	if err != nil {
		return nil, err
	}
	if kubernetesDirectoryExists {
		if projectInfo.KubernetesTemplatedResourcesDirectoryPath == "" {
			return nil, fmt.Errorf("kubernetes templated resources directory path cannot be empty")
		}

		for _, cluster := range projectInfo.KubernetesClusters {
			if cluster.ApiServer == "" {
				return nil, fmt.Errorf("kubernetes api server cannot be empty")
			}
			if cluster.CaPath == "" {
				return nil, fmt.Errorf("kubernetes ca path cannot be empty")
			}
			if cluster.CrtPath == "" {
				return nil, fmt.Errorf("kubernetes crt path cannot be empty")
			}
			if cluster.KeyPath == "" {
				return nil, fmt.Errorf("kubernetes key path cannot be empty")
			}
			if cluster.KubectlVersion == "" {
				return nil, fmt.Errorf("kubectl version cannot be empty")
			}

			kubectlClusterInfo := commands.NewDockerCommand(
				KubectlClusterInfoCommandName,
				commands.DockerCommandConfig{
					Volumes: []string{
						fmt.Sprintf("%v:/ca.pem", cluster.CaPath),
						fmt.Sprintf("%v:/crt.pem", cluster.CrtPath),
						fmt.Sprintf("%v:/key.pem", cluster.KeyPath),
					},
					Image: fmt.Sprintf("giantswarm/kubectl:%v", cluster.KubectlVersion),
					Args: []string{
						fmt.Sprintf("--server=%v", cluster.ApiServer),
						"--certificate-authority=/ca.pem",
						"--client-certificate=/crt.pem",
						"--client-key=/key.pem",
						"cluster-info",
					},
				},
			)
			w = append(w, kubectlClusterInfo)

			kubectlApply := commands.NewDockerCommand(
				KubectlApplyCommandName,
				commands.DockerCommandConfig{
					Volumes: []string{
						fmt.Sprintf("%v:/ca.pem", cluster.CaPath),
						fmt.Sprintf("%v:/crt.pem", cluster.CrtPath),
						fmt.Sprintf("%v:/key.pem", cluster.KeyPath),
						fmt.Sprintf("%v:/kubernetes", projectInfo.KubernetesTemplatedResourcesDirectoryPath),
					},
					Image: fmt.Sprintf("giantswarm/kubectl:%v", cluster.KubectlVersion),
					Args: []string{
						fmt.Sprintf("--server=%v", cluster.ApiServer),
						"--certificate-authority=/ca.pem",
						"--client-certificate=/crt.pem",
						"--client-key=/key.pem",
						"apply", "-R", "-f", "/kubernetes",
					},
				},
			)
			w = append(w, kubectlApply)
		}
	}

	return w, nil
}
