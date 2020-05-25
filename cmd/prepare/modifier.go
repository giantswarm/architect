package prepare

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/giantswarm/microerror"
)

type modifierConfig struct {
	NewVersion string
	Repo       string
	WorkingDir string
}

type modifier struct {
	newVersion string
	repo       string
	workingDir string
}

func newModifier(config modifierConfig) (*modifier, error) {
	if config.NewVersion == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.NewVersion must not be empty", config)
	}
	if config.Repo == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Repo must not be empty", config)
	}
	if config.WorkingDir == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.WorkingDir must not be empty", config)
	}

	m := &modifier{
		newVersion: config.NewVersion,
		repo:       config.Repo,
		workingDir: config.WorkingDir,
	}

	return m, nil
}

func (m *modifier) AddReleaseToChangelogMd() error {
	file := "CHANGELOG.md"
	modifyFunc := m.addReleaseToChangelogMd

	err := modifyFile(filepath.Join(m.workingDir, file), modifyFunc)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (m *modifier) addReleaseToChangelogMd(content []byte) ([]byte, error) {
	// Define replacements.

	unreleasedHeader := regexp.MustCompile(regexp.QuoteMeta("## [Unreleased]"))
	unreleasedHeaderReplacement := strings.Join([]string{
		"## [Unreleased]",
		"",
		fmt.Sprintf("## [%s] - %s", currentVersion, date),
	}, "\n")

	// To match strings like:
	//
	//	[Unreleased]: https://github.com/giantswarm/REPOSITORY_NAME/compare/v1.2.3...HEAD
	//
	bottomLinks := regexp.MustCompile(`^\[Unreleased]:\s+https://github.com/\S+/compare/v(\d+\.\d+\.\d+)\.\.\.HEAD\s*$`)
	bottomLinksReplacement := strings.Join([]string{
		fmt.Sprintf("[Unreleased]: https://github.com/%s/compare/v%s...HEAD", repository, currentVersion),
		fmt.Sprintf("[%s]: https://github.com/%s/compare/v${1}...v%s", currentVersion, repository, currentVersion),
	}, "\n")

	// To match strings like:
	//
	//	[Unreleased]: https://github.com//REPOSITORY_NAME/tree/master
	//
	bottomLinksFirstRelease := regexp.MustCompile(`^\[Unreleased]:\s+https://github.com/\S+/tree/master\s*$`)
	bottomLinksFirstReleaseReplacement := strings.Join([]string{
		fmt.Sprintf("[Unreleased]: https://github.com/%s/compare/v%s...HEAD", repository, currentVersion),
		fmt.Sprintf("[%s]: https://github.com/%s/releases/tag/v%s", currentVersion, repository, currentVersion),
	}, "\n")

	// Validate.

	err = validateSingleOccurence(content, unreleasedHeader)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = validateSingleOccurence(content, bottomLinks, bottomLinksFirstRelease)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Execute replacements.
	content = unreleasedHeader.ReplaceAll(content, []byte(unreleasedHeaderReplacement))
	content = bottomLinks.ReplaceAll(content, []byte(bottomLinksReplacement))
	content = bottomLinksFirstRelease.ReplaceAll(content, []byte(bottomLinksFirstReleaseReplacement))

	return content, nil
}

func (m *modifier) UpdateVersionInProjectGo() error {
	file := "pkg/project/project.go"
	modifyFunc := m.updateVersionInProjectGo

	err := modifyFile(filepath.Join(m.workingDir, file), modifyFunc)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (m *modifier) updateVersionInProjectGo(content []byte) ([]byte, error) {
	// Define replacements.

	// To match strings like:
	//
	//	version = "1.2.3"
	//	version = "1.2.3-any-suffix"
	//
	version := regexp.MustCompile(`version\s*=\s*"[0-9]+\.[0-9]+\.[0-9]+\S*"`)
	versionReplacement := fmt.Sprintf(`version = "%s"`, version)

	// Validate.

	err = validateSingleOccurence(content, version)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Execute replacements.
	content = version.ReplaceAll(content, []byte(versionReplacement))

	return content, nil
}

func modifyFile(path string, modifyFunc func([]byte) ([]byte, error)) error {
	// Make sure file exists and it's not a directory.
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return microerror.Maskf("file %#q not found", relPath)
	} else if err != nil {
		return microerror.Mask(err)
	} else if info.IsDir() {
		return microerror.Maskf("file %#q is a directory, expected regular file", relPath)
	}

	content, err := ioutil.ReadFile(relPath)
	if err != nil {
		return microerror.Mask(err)
	}

	content, err = modifyFunc(content)
	if err != nil {
		return microerror.Mask(err)
	}

	err = ioutil.WriteFile(changelogFile, content, 0)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func validateSingleOccurence(data []byte, regexps ...*regexp.Regexp) error {
	matches := 0
	for _, re := range regexps {
		matches += len(re.FindAllIndex(data, -1))
	}

	var combined *regexp.Regexp
	if len(regexps) == 1 {
		combined = regexps[0]
	} else {
		var patterns []string

		for _, re := range regexps {
			patterns = append(patterns, re.String())
		}

		pattern = fmt.Sprintf("(?:%s)", strings.Join(patterns, ") | (?:"))

		combined = regexp.MustCompile(pattern)
	}

	if matches == 0 {
		return nil, microerror.Maskf(executionFailedError, "no match for pattern %#q match found in project.go", combined)
	}
	if matches > 1 {
		return nil, microerror.Maskf(executionFailedError, "%d pattern %#q matches found in project.go, expected 1", matches, combined)
	}

	return nil
}
