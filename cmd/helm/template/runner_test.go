package template

import (
	"context"
	"os"
	"strconv"
	"testing"

	repo "github.com/giantswarm/gitrepo/pkg/gitrepo"
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
			name:            "case 0: aws-operator@v8.1.1",
			repoURL:         "git@github.com:giantswarm/aws-operator.git",
			tag:             "v8.1.1",
			expectedVersion: "n/a",
		},
		{
			name:            "case 1: azure-operator@v3.0.0",
			repoURL:         "git@github.com:giantswarm/azure-operator.git",
			tag:             "v3.0.0",
			expectedVersion: "3.0.0",
		},
		{
			name:         "case 2: file not found",
			repoURL:      "git@github.com:giantswarm/gitrepo-test.git",
			errorMatcher: repo.IsFileNotFound,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Log(tc.name)

			dir := "/tmp/architect-test-getprojectversion"
			defer os.RemoveAll(dir)

			c := repo.Config{
				Dir: dir,
				URL: tc.repoURL,
			}
			repo, err := repo.New(c)
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()

			err = repo.EnsureUpToDate(ctx)
			if err != nil {
				t.Fatal(err)
			}

			version, err := getProjectVersion(repo, tc.tag)

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
