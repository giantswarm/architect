package workflow

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/giantswarm/architect/tasks"
	"github.com/giantswarm/architect/utils"
)

var (
	GoFmtTaskName   = "go-fmt"
	GoVetTaskName   = "go-vet"
	GoTestTaskName  = "go-test"
	GoBuildTaskName = "go-build"
)

func checkGolangRequirements(projectInfo ProjectInfo) error {
	if projectInfo.WorkingDirectory == "" {
		return emptyWorkingDirectoryError
	}
	if projectInfo.Organisation == "" {
		return emptyOrganisationError
	}
	if projectInfo.Project == "" {
		return emptyProjectError
	}

	if projectInfo.Goos == "" {
		return emptyGoosError
	}
	if projectInfo.Goarch == "" {
		return emptyGoarchError
	}
	if projectInfo.GolangImage == "" {
		return emptyGolangImageError
	}
	if projectInfo.GolangVersion == "" {
		return emptyGolangVersionError
	}

	return nil
}

func NewGoFmtTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return nil, err
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
		return nil, err
	}

	packageArguments, err := utils.NoVendor(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return nil, err
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
			Args:  []string{"go", "vet"},
		},
	)
	goVet.Args = append(goVet.Args, packageArguments...)

	return goVet, nil
}

func NewGoTestTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return nil, err
	}

	packageArguments, err := utils.NoVendor(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return nil, err
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
			Args:  []string{"go", "test", "-v"},
		},
	)
	goTest.Args = append(goTest.Args, packageArguments...)

	return goTest, nil
}

func NewGoBuildTask(fs afero.Fs, projectInfo ProjectInfo) (tasks.Task, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return nil, err
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
