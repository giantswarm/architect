package internal

import (
	"fmt"
	"time"
	"strconv"
	"testing"
)

const (
	fileProjectGO = `version        = "%s"`
)

func Test_modifier_addReleaseToChangelogMd(t *testing.T) {
	testCases := []struct {
		name              string
		lastReleaseMinus1 string
		lastRelease       string
		newVersion        string
	}{
		{
			name:              "case 0: Production release from Production release",
			lastReleaseMinus1: "1.2.1",
			lastRelease:       "1.2.2",
			newVersion:        "1.2.3",
		},
		{
			name:              "case 1: Dev release with no dots from Production release",
			lastReleaseMinus1: "1.2.1",
			lastRelease:       "1.2.2",
			newVersion:        "1.2.3-gsalpha1",
		},
		{
			name:              "case 2: Dev release with dots from Production release",
			lastReleaseMinus1: "1.2.1",
			lastRelease:       "1.2.2",
			newVersion:        "1.2.3-gs.alpha.1",
		},
		{
			name:              "case 3: Dev release with no dots from dev release with no dots",
			lastReleaseMinus1: "1.2.1",
			lastRelease:       "1.2.2-gsalpha1",
			newVersion:        "1.2.2-gsalpha2",
		},
		{
			name:              "case 4: Dev release with dots from dev release with dots",
			lastReleaseMinus1: "1.2.1",
			lastRelease:       "1.2.2-gs.alpha.1",
			newVersion:        "1.2.2-gs.alpha.2",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			m := Modifier{
				newVersion: tc.newVersion,
				repo: "REPOSITORY_NAME",
			}

			// Current MD for Changelog before a new release is published
			changelogMD := fmt.Sprintf(`## [Unreleased]
- New Changes
## [%s] - 2021-07-14
- Old Changes

[Unreleased]: https://github.com/REPOSITORY_NAME/compare/v%s...HEAD
[%s]: https://github.com/giantswarm/REPOSITORY_NAME/compare/v%s...v%s`,
				tc.lastRelease, tc.lastRelease, tc.lastRelease, tc.lastReleaseMinus1, tc.lastRelease)

			// Changelog expected after new release is published
			expectedChangelogMD := fmt.Sprintf(`## [Unreleased]

## [%s] - %s
- New Changes
## [%s] - 2021-07-14
- Old Changes

[Unreleased]: https://github.com/REPOSITORY_NAME/compare/v%s...HEAD
[%s]: https://github.com/REPOSITORY_NAME/compare/v%s...v%s
[%s]: https://github.com/giantswarm/REPOSITORY_NAME/compare/v%s...v%s`,
				tc.newVersion, time.Now().Format("2006-01-02"), tc.lastRelease, tc.newVersion, tc.newVersion, tc.lastRelease, tc.newVersion, tc.lastRelease, tc.lastReleaseMinus1, tc.lastRelease)

			content, err := m.addReleaseToChangelogMd([]byte(changelogMD))
			if err != nil {
				t.Fatalf("actual = %s, expected nil", err)
			}

			if string(content) != expectedChangelogMD {
				t.Fatalf("expected %#q, got %#q", expectedChangelogMD, string(content))
			}
		})
	}

}

func Test_modifier_UpdateVersionInProjectGo(t *testing.T) {
	testCases := []struct {
		name              string
		currentProjectGO  string
		releaseVersion    string
		expectedProjectGO string
	}{
		{
			name:              "case 0: non-reference version release",
			currentProjectGO:  fmt.Sprintf(fileProjectGO, "0.1.0"),
			releaseVersion:    "0.2.0",
			expectedProjectGO: fmt.Sprintf(fileProjectGO, "0.2.0"),
		},
		{
			name:              "case 1: non-reference version release from dev",
			currentProjectGO:  fmt.Sprintf(fileProjectGO, "0.1.0-dev"),
			releaseVersion:    "0.2.0",
			expectedProjectGO: fmt.Sprintf(fileProjectGO, "0.2.0"),
		},
		{
			name:              "case 2: reference version release",
			currentProjectGO:  fmt.Sprintf(fileProjectGO, "0.1.0"),
			releaseVersion:    "0.1.0-1",
			expectedProjectGO: fmt.Sprintf(fileProjectGO, "0.1.0"),
		},
		{
			name:              "case 3: reference version release from dev",
			currentProjectGO:  fmt.Sprintf(fileProjectGO, "0.1.0-dev"),
			releaseVersion:    "0.1.0-1",
			expectedProjectGO: fmt.Sprintf(fileProjectGO, "0.1.0"),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			m := Modifier{
				newVersion: tc.releaseVersion,
			}

			content, err := m.updateVersionInProjectGo([]byte(tc.currentProjectGO))
			if err != nil {
				t.Fatalf("actual = %s, expected nil", err)
			}

			if string(content) != tc.expectedProjectGO {
				t.Fatalf("expected %#q, got %#q", tc.expectedProjectGO, string(content))
			}
		})
	}

}
