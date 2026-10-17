package template

import (
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/giantswarm/gitsemver/v3/pkg/gitsemver"
	"github.com/giantswarm/microerror"
)

// TestGetProjectVersions tests getProjectversion method which retrieves
// the value of the version variable in pkg/project/project.go file.
func TestGetProjectVersion(t *testing.T) {
	t.Parallel()

	architectGitTopLevelDir, err := gitsemver.TopLevel(".")
	if err != nil {
		t.Fatalf("err = %#q, want %#v", microerror.JSON(err), nil)
	}

	// project.Version() is overridden by the -X ldflags of a release build, so
	// the expected value is read from the source the function under test reads.
	source, err := os.ReadFile(filepath.Clean(path.Join(architectGitTopLevelDir, "pkg/project/project.go")))
	if err != nil {
		t.Fatalf("err = %#q, want %#v", microerror.JSON(err), nil)
	}
	match := regexp.MustCompile(`(?m)^\s*version\s*=\s*"([^"]*)"`).FindSubmatch(source)
	if match == nil {
		t.Fatalf("no version literal in pkg/project/project.go")
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
			expectedVersion: string(match[1]),
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

// TestParseString covers the declarations of version a project.go can carry;
// only a string literal yields a value, every other form an empty string.
func TestParseString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "string literal",
			source:   `var version = "1.2.3"`,
			expected: "1.2.3",
		},
		{
			name:     "string literal in a var block with several names",
			source:   "var (\n\tname, version = \"app\", \"1.2.3\"\n)",
			expected: "1.2.3",
		},
		{
			name:     "selector expression",
			source:   "import \"example.com/buildinfo\"\nvar version = buildinfo.DevVersion",
			expected: "",
		},
		{
			name:     "declaration without a value",
			source:   `var version string`,
			expected: "",
		},
		{
			name:     "no version variable",
			source:   `var name = "app"`,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			node, err := parser.ParseFile(token.NewFileSet(), "", "package project\n"+tc.source+"\n", 0)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}

			version, err := parseString(node, "version")
			if err != nil {
				t.Fatalf("error == %#v, want nil", err)
			}
			if version != tc.expected {
				t.Errorf("got %#q, expected %#q", version, tc.expected)
			}
		})
	}
}
