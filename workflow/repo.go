package workflow

import (
	"fmt"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

const (
	RepoCheckTaskName = "repo-check"

	// licenseTextFormat is used to generate the license text so the format can
	// be checked.
	licenseTextFormat = "Copyright 2016 - %d Giant Swarm GmbH"
)

type RepoCheckTask struct {
	fs afero.Fs
}

func (t RepoCheckTask) Name() string {
	return RepoCheckTaskName
}

func (t RepoCheckTask) Run() error {
	requiredFiles := []string{"DCO", "LICENSE"}

	for _, requiredFile := range requiredFiles {
		if _, err := t.fs.Stat(requiredFile); err != nil {
			fmt.Printf("repo does not have required file '%s', see https://github.com/giantswarm/example-opensource-repo\n", requiredFile)
			return microerror.Maskf(missingFileError, requiredFile)
		}

		fmt.Printf("repo has required file '%s'\n", requiredFile)
	}

	err := t.checkLicenseText()
	if err != nil {
		return microerror.Mask(err)
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

func (t RepoCheckTask) checkLicenseText() error {
	licenseText := fmt.Sprintf(licenseTextFormat, time.Now().Year())

	l, err := afero.ReadFile(t.fs, "LICENSE")
	if err != nil {
		return microerror.Mask(err)
	}

	license := string(l)
	if !strings.Contains(license, licenseText) {
		return microerror.Maskf(missingLicenseTextError, "LICENSE file does not contain text %#q", licenseText)
	}

	return nil
}
