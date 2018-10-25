package workflow

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

var (
	GoPullTaskName       = "go-pull"
	GoFmtTaskName        = "go-fmt"
	GoVetTaskName        = "go-vet"
	GoTestTaskName       = "go-test"
	GoBuildTaskName      = "go-build"
	GoConcurrentTaskName = "go-concurrent"
)

func golangFilesExist(fs afero.Fs, directory string) (bool, error) {
	err := afero.Walk(fs, directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Return io.EOF to break walk early.
		if filepath.Ext(path) == ".go" {
			return io.EOF
		}

		return nil
	})

	if err == io.EOF {
		return true, nil
	}

	if err != nil {
		return false, microerror.Mask(err)
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

func goTestable(directory string) (bool, error) {
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = directory

	out, err := cmd.Output()
	if err != nil {
		return false, microerror.Mask(err)
	}

	numPackages := strings.Count(string(out), "\n")

	return numPackages > 0, nil
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
				"/tmp/go/cache:/go/cache",
			},
			Env: []string{
				"GOPATH=/go",
				"GOCACHE=/go/cache",
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
				"/tmp/go/cache:/go/cache",
			},
			Env: []string{
				fmt.Sprintf("GOOS=%v", projectInfo.Goos),
				fmt.Sprintf("GOARCH=%v", projectInfo.Goarch),
				"GOCACHE=/go/cache",
				"GOPATH=/go",
				"CGOENABLED=0",
			},
			WorkingDirectory: fmt.Sprintf(
				"/go/src/github.com/%v/%v",
				projectInfo.Organisation,
				projectInfo.Project,
			),
			Image: fmt.Sprintf("%v:%v", projectInfo.GolangImage, projectInfo.GolangVersion),
			Args:  []string{"go", "test", "-race", "./..."},
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
				"/tmp/go/cache:/go/cache",
			},
			Env: []string{
				fmt.Sprintf("GOOS=%v", projectInfo.Goos),
				fmt.Sprintf("GOARCH=%v", projectInfo.Goarch),
				"GOCACHE=/go/cache",
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
				"-a",
				"-v",
				"-tags", "netgo",
				"-ldflags", fmt.Sprintf(
					"-w -X main.gitCommit=%s -linkmode 'external' -extldflags '-static'",
					projectInfo.Sha,
				),
			},
		},
	)

	return goBuild, nil
}
