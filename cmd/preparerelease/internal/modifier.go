package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
)

type ModifierConfig struct {
	NewVersion string
	Repo       string
	WorkingDir string
}

type Modifier struct {
	newVersion string
	repo       string
	workingDir string
}

func NewModifier(config ModifierConfig) (*Modifier, error) {
	if config.NewVersion == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.NewVersion must not be empty", config)
	}
	if config.Repo == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Repo must not be empty", config)
	}
	if config.WorkingDir == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.WorkingDir must not be empty", config)
	}

	m := &Modifier{
		newVersion: config.NewVersion,
		repo:       config.Repo,
		workingDir: config.WorkingDir,
	}

	return m, nil
}

func (m *Modifier) AddReleaseToChangelogMd() error {
	file := FileChangelogMd
	modifyFunc := m.addReleaseToChangelogMd

	err := modifyFile(filepath.Join(m.workingDir, file), modifyFunc)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (m *Modifier) addReleaseToChangelogMd(content []byte) ([]byte, error) {
	var err error

	date := time.Now().Format("2006-01-02")

	// Define replacements.

	unreleasedHeader := regexp.MustCompile(regexp.QuoteMeta("## [Unreleased]"))
	unreleasedHeaderReplacement := strings.Join([]string{
		"## [Unreleased]",
		"",
		fmt.Sprintf("## [%s] - %s", m.newVersion, date),
	}, "\n")

	// To match strings like:
	//
	//	[Unreleased]: https://github.com/giantswarm/REPOSITORY_NAME/compare/v1.2.3...HEAD
	//	[Unreleased]: https://github.com/giantswarm/REPOSITORY_NAME/compare/v1.2.3-gsalpha1...HEAD
	//	[Unreleased]: https://github.com/giantswarm/REPOSITORY_NAME/compare/v1.2.3-gs.alpha.1...HEAD
	//
	bottomLinks := regexp.MustCompile(`\[Unreleased\]:\s+https://github.com/\S+/compare/v(\d+\.\d+\.\d+(-(\w|\.)+)?)\.\.\.HEAD\s*`)
	bottomLinksReplacement := strings.Join([]string{
		fmt.Sprintf("[Unreleased]: https://github.com/%s/compare/v%s...HEAD", m.repo, m.newVersion),
		fmt.Sprintf("[%s]: https://github.com/%s/compare/v${1}...v%s", m.newVersion, m.repo, m.newVersion),
		"",
	}, "\n")

	// To match strings like:
	//
	//	[Unreleased]: https://github.com//REPOSITORY_NAME/tree/master
	//
	bottomLinksFirstRelease := regexp.MustCompile(`\[Unreleased\]:\s+https://github.com/\S+/tree/\S+\s*`)
	bottomLinksFirstReleaseReplacement := strings.Join([]string{
		fmt.Sprintf("[Unreleased]: https://github.com/%s/compare/v%s...HEAD", m.repo, m.newVersion),
		fmt.Sprintf("[%s]: https://github.com/%s/releases/tag/v%s", m.newVersion, m.repo, m.newVersion),
		"",
	}, "\n")

	// Validate.

	err = validateSingleOccurrence(content, unreleasedHeader)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	err = validateSingleOccurrence(content, bottomLinks, bottomLinksFirstRelease)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Execute replacements.
	content = unreleasedHeader.ReplaceAll(content, []byte(unreleasedHeaderReplacement))
	content = bottomLinks.ReplaceAll(content, []byte(bottomLinksReplacement))
	content = bottomLinksFirstRelease.ReplaceAll(content, []byte(bottomLinksFirstReleaseReplacement))

	return content, nil
}

func (m *Modifier) UpdateVersionInProjectGo() error {
	file := FileProjectGo
	modifyFunc := m.updateVersionInProjectGo

	err := modifyFile(filepath.Join(m.workingDir, file), modifyFunc)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (m *Modifier) updateVersionInProjectGo(content []byte) ([]byte, error) {
	var err error

	// Define replacements.

	// To match strings like:
	//
	//	version = "1.2.3"
	//	version = "1.2.3-any-suffix"
	//
	version := regexp.MustCompile(`(version\s*=\s*)"[0-9]+\.[0-9]+\.[0-9]+\S*"`)
	var finalVersion string
	{
		// Match strings like:
		//
		// version = "1.2.3-1"
		//
		// Drop version suffix as we keep previous verion

		refVersion := regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+-([0-9])+`)
		if refVersion.MatchString(m.newVersion) {
			finalVersion = strings.Split(m.newVersion, "-")[0]
		} else {
			finalVersion = m.newVersion
		}

	}

	versionReplacement := fmt.Sprintf(`$1"%s"`, finalVersion)

	// Validate.

	err = validateSingleOccurrence(content, version)
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
		return microerror.Maskf(fileNotFoundError, "file %#q not found", path)
	} else if err != nil {
		return microerror.Mask(err)
	} else if info.IsDir() {
		return microerror.Maskf(executionFailedError, "file %#q is a directory, expected regular file", path)
	}

	path = filepath.Clean(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return microerror.Mask(err)
	}

	content, err = modifyFunc(content)
	if err != nil {
		return microerror.Mask(err)
	}

	err = os.WriteFile(path, content, 0)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
