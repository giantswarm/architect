package workflow

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/architect/tasks"
)

var (
	GoPullTaskName  = "go-pull"
	GoFmtTaskName   = "go-fmt"
	GoVetTaskName   = "go-vet"
	GoTestTaskName  = "go-test"
	GoBuildTaskName = "go-build"
)

func golangFilesExist(fs afero.Fs, directory string) (bool, error) {
	fileInfos, err := afero.ReadDir(fs, directory)
	if err != nil {
		return false, microerror.Mask(err)
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			subFileInfos, err := afero.ReadDir(fs, path.Join(directory, fileInfo.Name()))
			if err != nil {
				return false, microerror.Mask(err)
			}

			for _, subFileInfo := range subFileInfos {
				if filepath.Ext(subFileInfo.Name()) == ".go" {
					return true, nil
				}
			}
		}

		if filepath.Ext(fileInfo.Name()) == ".go" {
			return true, nil
		}
	}

	return false, nil
}

func goBuildable(fs afero.Fs, directory string) (bool, error) {
	fileInfos, err := afero.ReadDir(fs, directory)
	if err != nil {
		return false, microerror.Mask(err)
	}

	for _, fileInfo := range fileInfos {
		if filepath.Ext(fileInfo.Name()) == ".go" {
			return true, nil
		}
	}

	return false, nil
}

func checkGolangRequirements(projectInfo ProjectInfo) error {
	if projectInfo.WorkingDirectory == "" {
		return microerror.Mask(emptyWorkingDirectoryError)
	}
	if projectInfo.Organisation == "" {
		return microerror.Mask(emptyOrganisationError)
	}
	if projectInfo.Project == "" {
		return microerror.Mask(emptyProjectError)
	}

	if projectInfo.Goos == "" {
		return microerror.Mask(emptyGoosError)
	}
	if projectInfo.Goarch == "" {
		return microerror.Mask(emptyGoarchError)
	}
	if projectInfo.GolangImage == "" {
		return microerror.Mask(emptyGolangImageError)
	}
	if projectInfo.GolangVersion == "" {
		return microerror.Mask(emptyGolangVersionError)
	}

	return nil
}

func NewGoPullTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	goPull := tasks.NewExecTask(
		GoPullTaskName,
		[]string{
			"docker", "pull", fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
		},
	)

	return goPull, nil
}

func NewGoFmtTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	goFmt := tasks.NewDockerTask(
		GoFmtTaskName,
		tasks.DockerTaskConfig{
			Volumes: []string{
				fmt.Sprintf(
					"%v:/go/src/github.com/%v/%v",
					projectInfo.WorkingDirectory,
					projectInfo.Organisation,
					projectInfo.Project,
				),
			},
			Env: []string{
				"GOPATH=/go",
			},
			WorkingDirectory: fmt.Sprintf(
				"/go/src/github.com/%v/%v",
				projectInfo.Organisation,
				projectInfo.Project,
			),
			Image: fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
			// gofmt always exits with code 0. Use grep matching everything to determine if any diff is outputed.
			// gofmt also requires specific files, so we use find to provide a list of all files.
			Args: []string{"bash", "-c", "! gofmt -d $(find . -type f -name '*.go' -not -path \"./vendor/*\") 2>&1 | grep -e '.'"},
		},
	)

	return goFmt, nil
}

func NewGoVetTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	goVet := tasks.NewDockerTask(
		GoVetTaskName,
		tasks.DockerTaskConfig{
			Volumes: []string{
				fmt.Sprintf(
					"%v:/go/src/github.com/%v/%v",
					projectInfo.WorkingDirectory,
					projectInfo.Organisation,
					projectInfo.Project,
				),
			},
			Env: []string{
				"GOPATH=/go",
			},
			WorkingDirectory: fmt.Sprintf(
				"/go/src/github.com/%v/%v",
				projectInfo.Organisation,
				projectInfo.Project,
			),
			Image: fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
			Args:  []string{"go", "vet", "./..."},
		},
	)

	return goVet, nil
}

func NewGoTestTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	goTest := tasks.NewDockerTask(
		GoTestTaskName,
		tasks.DockerTaskConfig{
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
			Args:  []string{"go", "test", "-v", "-race", "./..."},
		},
	)

	return goTest, nil
}

func NewGoBuildTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return nil, microerror.Mask(err)
	}

	goBuild := tasks.NewDockerTask(
		GoBuildTaskName,
		tasks.DockerTaskConfig{
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
			Args: []string{
				"go", "build",
				"-v",
				"-a",
				"-tags", "netgo",
				"-ldflags", fmt.Sprintf(
					"-X main.gitCommit=%s -linkmode 'external' -extldflags '-static'",
					projectInfo.Sha,
				),
			},
		},
	)

	return goBuild, nil
}
