package workflow

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"
)

var (
	GoTestCommandName  = "go-test"
	GoBuildCommandName = "go-build"
)

func checkGolangRequirmeents(projectInfo ProjectInfo) error {
	if projectInfo.WorkingDirectory == "" {
		return fmt.Errorf("working directory cannot be empty")
	}
	if projectInfo.Organisation == "" {
		return fmt.Errorf("organisation cannot be empty")
	}
	if projectInfo.Project == "" {
		return fmt.Errorf("project cannot be empty")
	}

	if projectInfo.Goos == "" {
		return fmt.Errorf("goos cannot be empty")
	}
	if projectInfo.Goarch == "" {
		return fmt.Errorf("goarch cannot be empty")
	}
	if projectInfo.GolangImage == "" {
		return fmt.Errorf("golang image cannot be empty")
	}
	if projectInfo.GolangVersion == "" {
		return fmt.Errorf("golang version cannot be empty")
	}

	return nil
}

func NewGoTestCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkGolangRequirmeents(projectInfo); err != nil {
		return commands.Command{}, err
	}

	testPackageArguments, err := utils.NoVendor(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return commands.Command{}, err
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

	return goTest, nil
}

func NewGoBuildCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkGolangRequirmeents(projectInfo); err != nil {
		return commands.Command{}, err
	}

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

	return goBuild, nil
}
