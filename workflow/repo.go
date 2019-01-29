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

	// dcoFileName is the name of the Developer Certificate of Origin file.
	dcoFileName = "DCO"

	// licenseFileName is the name of the License file.
	licenseFileName = "LICENSE"

	// licenseTextFormat generates the license text so the format can be
	// checked.
	licenseTextFormat = "Copyright 2016 - %d Giant Swarm GmbH"
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

	l, err := afero.ReadFile(t.fs, licenseFileName)
	if err != nil {
		return microerror.Mask(err)
	}

	license := string(l)
	if !strings.Contains(license, licenseText) {
		fmt.Printf("repo '%s' file does not contain required text '%s'\n", licenseFileName, licenseText)
		return microerror.Maskf(failedExecutionError, "license does not contain text '%s'", licenseText)
	}

	fmt.Printf("repo '%s' has required text '%s'\n", licenseFileName, licenseText)

	return nil
}
