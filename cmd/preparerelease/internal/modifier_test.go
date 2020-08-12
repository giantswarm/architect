package internal

import (
	"fmt"
	"strconv"
	"testing"
)

const (
	fileProjectGO = `version        = "%s"`
)

func Test_modifier_addReleaseToChangelogMd(t *testing.T) {
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
