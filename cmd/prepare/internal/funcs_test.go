package internal

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func Test_validateSingleOccurence(t *testing.T) {
	testCases := []struct {
		name          string
		inputData     string
		inputRegexps  []*regexp.Regexp
		expectedError bool
	}{
		{
			name: "case 0",
			inputData: strings.Join([]string{
				"line 1",
				"line 2",
				`[Unreleased]: https://github.com/giantswarm/architect/compare/v1.0.0...HEAD`,
				`[1.0.0]: https://github.com/giantswarm/architect/releases/tag/v1.0.0`,
			}, "\n"),
			inputRegexps: []*regexp.Regexp{
				regexp.MustCompile(`non existent`),
				regexp.MustCompile(`\[Unreleased\]:\s+https://github.com/\S+/compare/v(\d+\.\d+\.\d+)\.\.\.HEAD\s*`),
				regexp.MustCompile(`non existent`),
			},
			expectedError: false,
		},
		{
			name: "case 0",
			inputData: strings.Join([]string{
				"line 1",
				"line 2",
				`[Unreleased]: https://github.com/giantswarm/architect/compare/v1.0.0...HEAD`,
				`[1.0.0]: https://github.com/giantswarm/architect/releases/tag/v1.0.0`,
			}, "\n"),
			inputRegexps: []*regexp.Regexp{
				regexp.MustCompile(`non existent`),
			},
			expectedError: true,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			err := validateSingleOccurence([]byte(tc.inputData), tc.inputRegexps...)

			if tc.expectedError && err == nil {
				t.Fatalf("actual = %s, expected non-nil", err)
			}

			if !tc.expectedError && err != nil {
				t.Fatalf("actual = %s, expected nil", err)
			}
		})
	}
}
