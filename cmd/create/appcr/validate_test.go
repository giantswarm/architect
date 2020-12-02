package appcr

import (
	"strconv"
	"testing"
)

func Test_validateConfigVersion(t *testing.T) {
	testCases := []struct {
		name               string
		inputConfigVersion string
		errorMatcher       func(err error) bool
	}{
		{
			name:               "case 0: valid major",
			inputConfigVersion: "3.x.x",
			errorMatcher:       nil,
		},
		{
			name:               "case 1: valid branch",
			inputConfigVersion: "my-branch",
			errorMatcher:       nil,
		},
		{
			name:               "case 2: invalid - provided minor",
			inputConfigVersion: "1.2.x",
			errorMatcher:       func(err error) bool { return err != nil },
		},
		{
			name:               "case 3: invalid - provided minor and patch",
			inputConfigVersion: "1.2.3",
			errorMatcher:       func(err error) bool { return err != nil },
		},
		{
			name:               "case 3: invalid - starts with number and dot",
			inputConfigVersion: "100.not-x.x",
			errorMatcher:       func(err error) bool { return err != nil },
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			err := validateConfigVersion(tc.inputConfigVersion)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if tc.errorMatcher != nil {
				return
			}
		})
	}
}
