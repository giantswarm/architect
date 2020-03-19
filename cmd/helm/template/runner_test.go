package template

import (
	"context"
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

	testCases := []struct {
		name            string
		repoURL         string
		tag             string
		expectedVersion string
		errorMatcher    func(err error) bool
	}{
		{
			name:            "case 0",
			expectedVersion: project.Version(),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			ctx := context.Background()

			dir, err := gitrepo.TopLevel(ctx, ".")
			if err != nil {
				t.Fatalf("err = %v, want %v", microerror.Stack(err), nil)
			}

			version, err := getProjectVersion(ctx, dir)

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
