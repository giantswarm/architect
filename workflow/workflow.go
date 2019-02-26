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
	Tag    string

	Registry       string
	DockerUsername string
	DockerPassword string

	HelmDirectoryPath string

	Goos          string
	Goarch        string
	GolangImage   string
	GolangVersion string

	Channels []string
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

	repoCheckTask, err := NewRepoCheckTask(fs, projectInfo)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	w = append(w, repoCheckTask)

	isGolangFilesExist, err := golangFilesExist(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if isGolangFilesExist {
		golangPull, err := NewGoPullTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		wrappedGolangPull := tasks.NewRetryTask(backoff.NewExponentialBackOff(), golangPull)
		w = append(w, wrappedGolangPull)

		var goTasks []tasks.Task
		goFmt, err := NewGoFmtTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		goTasks = append(goTasks, goFmt)

		isGoBuildable, err := goBuildable(fs, projectInfo.WorkingDirectory)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if isGoBuildable {
			goBuild, err := NewGoBuildTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			goTasks = append(goTasks, goBuild)
		}

		isGoTestable, err := goTestable(projectInfo.WorkingDirectory)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		if isGoTestable {
			goTest, err := NewGoTestTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			goTasks = append(goTasks, goTest)
		}

		goConcurrentTask := tasks.NewConcurrentTask(GoConcurrentTaskName, goTasks...)

		w = append(w, goConcurrentTask)
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
			dockerPushRef, err := NewDockerPushRefTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedDockerPushRef := tasks.NewRetryTask(backoff.NewExponentialBackOff(), dockerPushRef)
			w = append(w, wrappedDockerPushRef)
		}
	}

	helmTasks, err := processHelmDir(fs, projectInfo, NewHelmPushTask)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	w = append(w, helmTasks...)

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

		dockerPull, err := NewDockerPullTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, dockerPull)

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
	}

	helmTasks, err := processHelmDir(fs, projectInfo, NewHelmPromoteToStableChannelTask)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	w = append(w, helmTasks...)

	return w, nil
}

func NewPublish(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	chartDirectory := filepath.Join(projectInfo.WorkingDirectory, "helm", projectInfo.Project+"-chart")
	chartDirectoryExists, err := afero.Exists(fs, chartDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if chartDirectoryExists {
		helmLogin, err := NewHelmLoginTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmLogin)

		helmChartTemplate, err := NewTemplateHelmChartTask(fs, chartDirectory, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmChartTemplate)

		for _, c := range projectInfo.Channels {
			if c == "" {
				return nil, microerror.Mask(emptyChannelError)
			}
			helmPromoteToChannel, err := NewHelmPromoteToChannelTask(fs, chartDirectory, projectInfo, c)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			wrappedHelmPromoteToChannel := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmPromoteToChannel)

			w = append(w, wrappedHelmPromoteToChannel)
		}
	}
	return w, nil
}

func NewUnpublish(projectInfo ProjectInfo, fs afero.Fs) (Workflow, error) {
	w := Workflow{}

	chartDirectory := filepath.Join(projectInfo.WorkingDirectory, "helm", projectInfo.Project+"-chart")
	chartDirectoryExists, err := afero.Exists(fs, chartDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if chartDirectoryExists {
		helmLogin, err := NewHelmLoginTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmLogin)

		for _, c := range projectInfo.Channels {
			if c == "" {
				return nil, microerror.Mask(emptyChannelError)
			}

			helmDeleteFromChannel, err := NewHelmDeleteFromChannelTask(fs, chartDirectory, projectInfo, c)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			wrappedHelmDeleteFromChannel := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmDeleteFromChannel)

			w = append(w, wrappedHelmDeleteFromChannel)
		}
	}
	return w, nil
}

func processHelmDir(fs afero.Fs, projectInfo ProjectInfo, f func(afero.Fs, string, ProjectInfo) (tasks.Task, error)) (Workflow, error) {
	var w = Workflow{}

	helmDirectory := filepath.Join(projectInfo.WorkingDirectory, "helm")
	helmDirectoryExists, err := afero.Exists(fs, helmDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if helmDirectoryExists {
		fileInfos, err := afero.ReadDir(fs, helmDirectory)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		if len(fileInfos) > 0 {
			helmPull, err := NewHelmPullTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedHelmPull := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmPull)

			w = append(w, wrappedHelmPull)

			helmLogin, err := NewHelmLoginTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			w = append(w, helmLogin)
		}

		prefix := projectInfo.Project + "-"
		suffix := "-chart"
		for _, fileInfo := range fileInfos {
			if !fileInfo.IsDir() {
				return nil, microerror.Maskf(invalidHelmDirectoryError, "%q is not a directory", fileInfo.Name())
			}
			if !strings.HasPrefix(fileInfo.Name(), prefix) {
				return nil, microerror.Maskf(invalidHelmDirectoryError, "%q must start with %q", fileInfo.Name(), prefix)
			}
			if !strings.HasSuffix(fileInfo.Name(), suffix) {
				return nil, microerror.Maskf(invalidHelmDirectoryError, "%q must end with %q", fileInfo.Name(), suffix)
			}

			chartDir := filepath.Join(helmDirectory, fileInfo.Name())

			helmChartTemplate, err := NewTemplateHelmChartTask(fs, chartDir, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			w = append(w, helmChartTemplate)

			helmTask, err := f(fs, chartDir, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}

			wrappedHelmTask := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmTask)

			w = append(w, wrappedHelmTask)
		}
	}

	return w, nil
}
