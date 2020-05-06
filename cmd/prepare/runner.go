package prepare

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/giantswarm/gitrepo/pkg/gitrepo"
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	ChangelogFile = "CHANGELOG.md"
	VersionFile   = "pkg/project/project.go"
)

func runPrepareRelease(cmd *cobra.Command, args []string) error {
	var err error

	absoluteCurrentFolder, err := gitrepo.TopLevel(context.Background(), ".")
	if err != nil {
		return microerror.Mask(err)
	}
	repositoryName := fmt.Sprintf("%s/%s", filepath.Base(filepath.Dir(absoluteCurrentFolder)), filepath.Base(absoluteCurrentFolder))

	version := cmd.Flag("version").Value.String()
	if version == "" {
		return microerror.Maskf(executionFailedError, "'version' can't be empty")
	}

	_, err = os.Stat(VersionFile)
	if err != nil && !os.IsNotExist(err) {
		return microerror.Mask(err)
	}

	if err == nil {
		err = replaceVersionInFile(VersionFile, version)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	err = addReleaseToChangelog(time.Now().Format("2006-01-02"), version, repositoryName)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func replaceVersionInFile(file, version string) error {
	// Read contents of the file containing current version
	filecontents, err := ioutil.ReadFile(file)
	if err != nil {
		return microerror.Mask(err)
	}
	versionFileContents := string(filecontents)

	versionRegex := regexp.MustCompile(`(version\s*=\s*)("[0-9]+\.[0-9]+\.[0-9]+-dev")`)
	currentVersion := versionRegex.FindSubmatch(filecontents)
	if len(currentVersion) < 1 {
		return microerror.Maskf(executionFailedError, "No version was found")
	}
	updatedFileContents := versionRegex.ReplaceAllString(versionFileContents, fmt.Sprintf("${1}\"%s\"", version))
	err = ioutil.WriteFile(file, []byte(updatedFileContents), 0)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func addReleaseToChangelog(date, currentVersion, repository string) error {
	// Read Changelog contents
	filecontents, err := ioutil.ReadFile(ChangelogFile)
	if err != nil {
		return microerror.Mask(err)
	}
	changelogContents := string(filecontents)

	// Check if there is 'Unreleased' work, otherwise there is no point in releasing
	tagname := fmt.Sprintf("v%s", currentVersion)
	search := "## [Unreleased]"
	if !strings.Contains(changelogContents, search) {
		return microerror.Maskf(executionFailedError, "No '[Unreleased]' work was found")
	}

	// Add new entry to the top of the changelog
	replaceWith := fmt.Sprintf("## [Unreleased]\n\n## [%s] %s", currentVersion, date)
	updatedFileContents := strings.Replace(changelogContents, search, replaceWith, 1)

	// Update links at the bottom of the changelog
	bottomLinks := regexp.MustCompile(`(\[Unreleased]:)(.*)(v[0-9]+\.[0-9]+\.[0-9]+)(...HEAD)\n`)
	updatedFileContents = bottomLinks.ReplaceAllString(updatedFileContents, fmt.Sprintf("$1${2}%s...HEAD\n\n[%s]: https://github.com/%s/compare/${3}...%s", tagname, currentVersion, repository, tagname))

	err = ioutil.WriteFile(ChangelogFile, []byte(updatedFileContents), 0)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
