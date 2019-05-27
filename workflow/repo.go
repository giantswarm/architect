package workflow

import (
	"fmt"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

const (
	RepoCheckTaskName = "repo-check"

	// dcoFileName is the name of the Developer Certificate of Origin file.
	dcoFileName = "DCO"

	// licenseFileName is the name of the License file.
	licenseFileName = "LICENSE"
)

type RepoCheckTask struct {
	fs afero.Fs
}

func (t RepoCheckTask) Name() string {
	return RepoCheckTaskName
}

func (t RepoCheckTask) Run() error {
	requiredFiles := []string{dcoFileName, licenseFileName}

	for _, requiredFile := range requiredFiles {
		if _, err := t.fs.Stat(requiredFile); err != nil {
			fmt.Printf("repo does not have required file '%s', see https://github.com/giantswarm/example-opensource-repo\n", requiredFile)
			return microerror.Maskf(missingFileError, requiredFile)
		}

		fmt.Printf("repo has required file '%s'\n", requiredFile)
	}

	return nil
}

func (t RepoCheckTask) String() string {
	return RepoCheckTaskName
}

func NewRepoCheckTask(fs afero.Fs, projectInfo ProjectInfo) (RepoCheckTask, error) {
	t := RepoCheckTask{
		fs: fs,
	}

	return t, nil
}
