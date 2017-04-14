package workflow

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/giantswarm/architect/commands"
	"github.com/giantswarm/architect/utils"
)

var (
	GoFmtCommandName   = "go-fmt"
	GoVetCommandName   = "go-vet"
	GoTestCommandName  = "go-test"
	GoBuildCommandName = "go-build"
)

func checkGolangRequirements(projectInfo ProjectInfo) error {
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

func NewGoFmtCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	goFmt := commands.NewDockerCommand(
		GoFmtCommandName,
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

func NewGoVetCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	packageArguments, err := utils.NoVendor(fs, projectInfo.WorkingDirectory)
	if err != nil {
		return commands.Command{}, err
	}

	goVet := commands.NewDockerCommand(
		GoVetCommandName,
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

func NewGoTestCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
		return commands.Command{}, err
	}

	packageArguments, err := utils.NoVendor(fs, projectInfo.WorkingDirectory)
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
	goTest.Args = append(goTest.Args, packageArguments...)

	return goTest, nil
}

func NewGoBuildCommand(fs afero.Fs, projectInfo ProjectInfo) (commands.Command, error) {
	if err := checkGolangRequirements(projectInfo); err != nil {
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
