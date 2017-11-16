package workflow

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cenk/backoff"
	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/microerror"

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

	Branch string
	Sha    string

	Registry       string
	DockerEmail    string
	DockerUsername string
	DockerPassword string

	HelmDirectoryPath string

	Goos          string
	Goarch        string
	GolangImage   string
	GolangVersion string
}

func NewBuild(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	workingDirectoryExists, err := afero.Exists(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if !workingDirectoryExists {
		return w, nil
	}

	isGolangFilesExist, err := golangFilesExist(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if isGolangFilesExist {
		{
			golangPull, err := NewGoPullTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedGolangPull := tasks.NewRetryTask(backoff.NewExponentialBackOff(), golangPull)

			w = append(w, wrappedGolangPull)
		}

		goFmt, err := NewGoFmtTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, goFmt)

		goVet, err := NewGoVetTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, goVet)

		goTest, err := NewGoTestTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, goTest)

		isGoBuildable, err := goBuildable(fs, projectInfo.WorkingDirectory)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if isGoBuildable {
			goBuild, err := NewGoBuildTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			w = append(w, goBuild)
		}
	}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if dockerFileExists {
		dockerBuild, err := NewDockerBuildTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, dockerBuild)
	}

	if isGolangFilesExist && dockerFileExists {
		dockerRunVersion, err := NewDockerRunVersionTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, dockerRunVersion)

		dockerRunHelp, err := NewDockerRunHelpTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, dockerRunHelp)
	}

	if dockerFileExists {
		{
			dockerLogin, err := NewDockerLoginTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedDockerLogin := tasks.NewRetryTask(backoff.NewExponentialBackOff(), dockerLogin)
			w = append(w, wrappedDockerLogin)
		}

		{
			dockerPushSha, err := NewDockerPushShaTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedDockerPushSha := tasks.NewRetryTask(backoff.NewExponentialBackOff(), dockerPushSha)
			w = append(w, wrappedDockerPushSha)
		}
	}

	helmDirectoryExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "helm"))
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if helmDirectoryExists {
		{
			helmPull, err := NewHelmPullTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedHelmPull := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmPull)

			w = append(w, wrappedHelmPull)
		}

		helmChartTemplate, err := NewTemplateHelmChartTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmChartTemplate)

		helmLogin, err := NewHelmLoginTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmLogin)

		helmPush, err := NewHelmPushTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmPush)
	}

	return w, nil
}

func NewDeploy(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	dockerFileExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "Dockerfile"))
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if dockerFileExists {
		{
			dockerLogin, err := NewDockerLoginTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedDockerLogin := tasks.NewRetryTask(backoff.NewExponentialBackOff(), dockerLogin)
			w = append(w, wrappedDockerLogin)
		}

		dockerTagLatest, err := NewDockerTagLatestTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, dockerTagLatest)

		{
			dockerPushLatest, err := NewDockerPushLatestTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedDockerPushLatest := tasks.NewRetryTask(backoff.NewExponentialBackOff(), dockerPushLatest)
			w = append(w, wrappedDockerPushLatest)
		}

		helmDirectoryExists, err := afero.Exists(fs, filepath.Join(projectInfo.WorkingDirectory, "helm"))
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if helmDirectoryExists {
			{
				helmPull, err := NewHelmPullTask(fs, projectInfo)
				if err != nil {
					return nil, microerror.Mask(err)
				}
				wrappedHelmPull := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmPull)

				w = append(w, wrappedHelmPull)
			}

			helmChartTemplate, err := NewTemplateHelmChartTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			w = append(w, helmChartTemplate)

			helmLogin, err := NewHelmLoginTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			w = append(w, helmLogin)

			helmPromoteToChannel, err := NewHelmPromoteToStableChannelTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			wrappedHelmPromoteToChannel := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmPromoteToChannel)

			w = append(w, wrappedHelmPromoteToChannel)
		}

	}

	return w, nil
}
