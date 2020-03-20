package template

import (
	"context"
	"path"
	"strconv"
	"testing"

	"github.com/giantswarm/gitrepo/pkg/gitrepo"
	"github.com/giantswarm/microerror"

	"github.com/giantswarm/architect/pkg/project"
)

// TestGetProjectVersions tests getProjectversion method which retrieves
// the value of the version variable in pkg/project/project.go file.
func TestGetProjectVersion(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	architectGitTopLevelDir, err := gitrepo.TopLevel(ctx, ".")
	if err != nil {
		t.Fatalf("err = %#q, want %#v", microerror.JSON(err), nil)
	}

	testCases := []struct {
		name            string
		inputDir        string
		expectedVersion string
		errorMatcher    func(err error) bool
	}{
		{
			name:            "case 0",
			inputDir:        architectGitTopLevelDir,
			expectedVersion: project.Version(),
		},
		{
			name:            "case 1",
			inputDir:        path.Join(architectGitTopLevelDir, "non-exitent"),
			expectedVersion: "",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			version, err := getProjectVersion(tc.inputDir)

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

			if err == nil {
				if version != tc.expectedVersion {
					t.Errorf("got %#q, expected %#q\n", version, tc.expectedVersion)
				}
			}
		})
	}
}
