package workflow

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"

	"github.com/spf13/afero"
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
	IngressTag     string
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

	KubernetesResourcesDirectoryPath string
	KubernetesClusters               []KubernetesCluster

	Goos          string
	Goarch        string
	GolangImage   string
	GolangVersion string
}

func NewBuild(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	mainGoExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "main.go"))
	if err != nil {
		return nil, err
	}
	if mainGoExists {
		goTest, err := NewGoTestCommand(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, goTest)

		goBuild, err := NewGoBuildCommand(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, goBuild)
	}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, err
	}
	if dockerFileExists {
		dockerBuild, err := NewDockerBuildCommand(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerBuild)

		dockerRunVersion, err := NewDockerRunVersionCommand(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerRunVersion)

		dockerRunHelp, err := NewDockerRunHelpCommand(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerRunHelp)
	}

	return w, nil
}

func NewDeploy(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, err
	}
	if dockerFileExists {
		dockerLogin, err := NewDockerLoginCommand(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerLogin)

		dockerPush, err := NewDockerPushCommand(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerPush)
	}

	kubernetesDirectoryExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "kubernetes/"))
	if err != nil {
		return nil, err
	}
	if kubernetesDirectoryExists {
		for _, cluster := range projectInfo.KubernetesClusters {
			if projectInfo.KubernetesResourcesDirectoryPath == "" {
				return nil, fmt.Errorf("kubernetes templated resources directory path cannot be empty")
			}

			if cluster.IngressTag == "" {
				return nil, fmt.Errorf("ingress tag cannot be empty")
			}

			// Copy /kubernetes to a per-cluster directory, and template it
			dir, subdir := filepath.Split(projectInfo.KubernetesResourcesDirectoryPath)
			templatedResourcesDirectory := filepath.Join(dir, subdir+"-"+cluster.IngressTag)

			if err := utils.CopyDir(
				fs,
				projectInfo.KubernetesResourcesDirectoryPath,
				templatedResourcesDirectory,
			); err != nil {
				return nil, err
			}

			if err := utils.TemplateKubernetesResources(fs, templatedResourcesDirectory, projectInfo.Sha, cluster.IngressTag); err != nil {
				return nil, err
			}

			kubectlClusterInfo, err := NewKubectlClusterInfoCommand(fs, cluster)
			if err != nil {
				return nil, err
			}
			w = append(w, kubectlClusterInfo)

			kubectlApply, err := NewKubectlApplyCommand(fs, cluster, templatedResourcesDirectory)
			if err != nil {
				return nil, err
			}
			w = append(w, kubectlApply)
		}
	}

	return w, nil
}
