package workflow

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/template"
	"github.com/giantswarm/architect/utils"

	"github.com/spf13/afero"
)

type Workflow []tasks.Task

func (w Workflow) String() string {
	if len(w) == 0 {
		return "{}"
	}

	taskStrings := []string{}
	for _, task := range w {
		taskStrings = append(taskStrings, "\t"+task.String()+"\n")
	}

	return fmt.Sprintf("{\n%v}", strings.Join(taskStrings, ""))
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

	HelmDirectoryPath                string
	KubernetesResourcesDirectoryPath string
	KubernetesClusters               []KubernetesCluster

	Goos          string
	Goarch        string
	GolangImage   string
	GolangVersion string
}

func NewBuild(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	workingDirectoryExists, err := afero.Exists(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return nil, err
	}

	goLangFilesExist := false
	if workingDirectoryExists {
		fileInfos, err := afero.ReadDir(fs, projectInfo.WorkingDirectory)
		if err != nil {
			return nil, err
		}

		for _, fileInfo := range fileInfos {
			if filepath.Ext(fileInfo.Name()) == ".go" {
				goLangFilesExist = true
				break
			}
		}
	}
	if goLangFilesExist {
		goFmt, err := NewGoFmtTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, goFmt)

		goVet, err := NewGoVetTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, goVet)

		goTest, err := NewGoTestTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, goTest)

		goBuild, err := NewGoBuildTask(fs, projectInfo)
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
		dockerBuild, err := NewDockerBuildTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerBuild)
	}

	if goLangFilesExist && dockerFileExists {
		dockerRunVersion, err := NewDockerRunVersionTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerRunVersion)

		dockerRunHelp, err := NewDockerRunHelpTask(fs, projectInfo)
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
		dockerLogin, err := NewDockerLoginTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerLogin)

		dockerTagLatest, err := NewDockerTagLatestTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerTagLatest)

		dockerPushSha, err := NewDockerPushShaTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerPushSha)

		dockerPushLatest, err := NewDockerPushLatestTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, dockerPushLatest)
	}

	helmDirectoryExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "helm"))
	if err != nil {
		return nil, err
	}
	if helmDirectoryExists {
		helmChartTemplateTask, err := NewTemplateHelmChartTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, helmChartTemplateTask)

		helmLogin, err := NewHelmLoginTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, helmLogin)

		helmPush, err := NewHelmPushTask(fs, projectInfo)
		if err != nil {
			return nil, err
		}
		w = append(w, helmPush)
	}

	kubernetesDirectoryExists, err := afero.Exists(fs, projectInfo.KubernetesResourcesDirectoryPath)
	if err != nil {
		return nil, err
	}

	if kubernetesDirectoryExists {
		for _, cluster := range projectInfo.KubernetesClusters {
			if projectInfo.KubernetesResourcesDirectoryPath == "" {
				return nil, emptyKubernetesResourcesDirectoryPath
			}

			// Copy /kubernetes to a per-cluster directory, and template it
			dir, subdir := filepath.Split(projectInfo.KubernetesResourcesDirectoryPath)
			templatedResourcesDirectory := filepath.Join(dir, subdir+"-"+cluster.Prefix)

			if err := utils.CopyDir(
				fs,
				projectInfo.KubernetesResourcesDirectoryPath,
				templatedResourcesDirectory,
			); err != nil {
				return nil, err
			}

			config := template.TemplateConfiguration{
				BuildInfo: template.BuildInfo{
					SHA: projectInfo.Sha,
				},
				Installation: cluster.Installation,
			}

			if err := template.TemplateKubernetesResources(fs, templatedResourcesDirectory, config); err != nil {
				return nil, err
			}

			kubectlClusterInfo, err := NewKubectlClusterInfoTask(fs, cluster)
			if err != nil {
				return nil, err
			}
			w = append(w, kubectlClusterInfo)

			kubectlApply, err := NewKubectlApplyTask(fs, cluster, templatedResourcesDirectory)
			if err != nil {
				return nil, err
			}
			w = append(w, kubectlApply)
		}
	}

	return w, nil
}
