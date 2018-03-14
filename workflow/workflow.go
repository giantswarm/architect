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

	isGolangFilesExist, err := golangFilesExist(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if isGolangFilesExist {
		var goTasks []tasks.Task

		golangPull, err := NewGoPullTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		wrappedGolangPull := tasks.NewRetryTask(backoff.NewExponentialBackOff(), golangPull)
		goTasks = append(goTasks, wrappedGolangPull)

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

		goTest, err := NewGoTestTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		goTasks = append(goTasks, goTest)

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
			dockerPushSha, err := NewDockerPushShaTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedDockerPushSha := tasks.NewRetryTask(backoff.NewExponentialBackOff(), dockerPushSha)
			w = append(w, wrappedDockerPushSha)
		}
	}

	helmDirectory := filepath.Join(projectInfo.WorkingDirectory, "helm")
	helmDirectoryExists, err := afero.Exists(fs, helmDirectory)
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

		fileInfos, err := afero.ReadDir(fs, helmDirectory)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		if len(fileInfos) > 0 {
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

			helmPush, err := NewHelmPushTask(fs, chartDir, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			w = append(w, helmPush)
		}
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
	}

	chartDirectory := filepath.Join(projectInfo.WorkingDirectory, "helm", projectInfo.Project+"-chart")
	chartDirectoryExists, err := afero.Exists(fs, chartDirectory)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	if chartDirectoryExists {
		{
			helmPull, err := NewHelmPullTask(fs, projectInfo)
			if err != nil {
				return nil, microerror.Mask(err)
			}
			wrappedHelmPull := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmPull)

			fmt.Printf("adding wrappedHelmPull %v\n", chartDirectory)
			w = append(w, wrappedHelmPull)
		}

		helmChartTemplate, err := NewTemplateHelmChartTask(fs, chartDirectory, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmChartTemplate)

		helmLogin, err := NewHelmLoginTask(fs, projectInfo)
		if err != nil {
			return nil, microerror.Mask(err)
		}
		w = append(w, helmLogin)

		helmPromoteToChannel, err := NewHelmPromoteToChannelTask(fs, chartDirectory, projectInfo, "stable")
		if err != nil {
			return nil, microerror.Mask(err)
		}

		wrappedHelmPromoteToChannel := tasks.NewRetryTask(backoff.NewExponentialBackOff(), helmPromoteToChannel)

		w = append(w, wrappedHelmPromoteToChannel)
	}

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
