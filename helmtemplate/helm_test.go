package helmtemplate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

const (
	BranchTag     = "[[ .Branch ]]"
	SHATag        = "[[ .SHA ]]"
	VersionTag    = "[[ .Version ]]"
	AppVersionTag = "[[ .AppVersion ]]"
)

// TestTemplateHelmChartTask tests the TemplateHelmChartTask.
func TestTemplateHelmChartTask(t *testing.T) {
	testCases := []struct {
		name             string
		config           Config
		files            map[string]string
		validateFlag     bool
		taggedBuildFlag  bool
		expectedChartDir string
		errorMatcher     func(err error) bool
	}{
		{
			name: "case 0: chart templating",
			config: Config{
				Branch:     "master",
				Sha:        "ea82e754178bb2b8065aca0a0760e77ce3733649",
				Version:    "1.2.3",
				AppVersion: "1.0.0",
			},
			files: map[string]string{
				HelmChartYamlName:  "version: [[ .Version ]]\nappVersion: [[ .AppVersion ]]\n",
				HelmValuesYamlName: "branch: [[ .Branch ]]\ncommit: [[ .SHA ]]\n",
			},
			expectedChartDir: "test0",
		},
		{
			name: "case 1: chart templating + validation",
			config: Config{
				Branch:     "master",
				Sha:        "ea82e754178bb2b8065aca0a0760e77ce3733649",
				Version:    "1.2.3",
				AppVersion: "1.2.3",
			},
			files: map[string]string{
				HelmChartYamlName:  "version: [[ .Version ]]\nappVersion: [[ .AppVersion ]]\n",
				HelmValuesYamlName: "branch: [[ .Branch ]]\ncommit: [[ .SHA ]]\n",
			},
			validateFlag:     true,
			taggedBuildFlag:  true,
			expectedChartDir: "test1",
		},
		{
			name: "case 2: validate version",
			config: Config{
				Branch:  "master",
				Sha:     "ea82e754178bb2b8065aca0a0760e77ce3733649",
				Version: "1.2.3",
			},
			files: map[string]string{
				HelmChartYamlName: "version: 2.0.0\n",
			},
			validateFlag: true,
			errorMatcher: IsValidationFailedError,
		},
		{
			name: "case 3: validate appVersion",
			config: Config{
				Branch:     "master",
				Sha:        "ea82e754178bb2b8065aca0a0760e77ce3733649",
				Version:    "1.2.3",
				AppVersion: "1.0.0",
			},
			files: map[string]string{
				HelmChartYamlName: "version: [[ .Version ]]\nappVersion: 2.0.0\n",
			},
			validateFlag: true,
			errorMatcher: IsValidationFailedError,
		},
		{
			name: "case 4: validate tagged build",
			config: Config{
				Branch:     "master",
				Sha:        "ea82e754178bb2b8065aca0a0760e77ce3733649",
				Version:    "1.2.3",
				AppVersion: "1.0.0",
			},
			validateFlag:    true,
			taggedBuildFlag: true,
			errorMatcher:    IsValidationFailedError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.config.ChartDir = "/tmp/architect-testtemplatehelmcharttask"
			tc.config.Fs = afero.NewMemMapFs()

			err := setup(tc.config, tc.files)
			if err != nil {
				t.Fatalf("unexpected error during setup: %v\n", err)
			}

			task, err := NewTemplateHelmChartTask(tc.config)
			if err != nil {
				t.Fatalf("unexpected error when creating NewTemplateHelmChartTask: %v\n", err)
			}

			err = task.Run(tc.validateFlag, tc.taggedBuildFlag)
			switch {
			case err == nil && tc.errorMatcher == nil:
				err = check(tc.config, tc.expectedChartDir)
				if err != nil {
					t.Fatalf("unexpected error during check: %v\n", err)
				}
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

		})
	}
}

// setup create files inside the config.Fs filesystem at config.ChartDir.
func setup(config Config, files map[string]string) error {
	err := config.Fs.MkdirAll(config.ChartDir, permission)
	if err != nil {
		return microerror.Mask(err)
	}

	for fpath, data := range files {
		path := filepath.Join(config.ChartDir, fpath)
		dir := filepath.Base(path)
		err := config.Fs.MkdirAll(dir, permission)
		if err != nil {
			return microerror.Mask(err)
		}

		err = afero.WriteFile(config.Fs, path, []byte(data), permission)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

// check compare all files inside config.Fs filesystem at config.ChartDir directory with files inside expectedChartDir directory.
func check(config Config, expectedChartDir string) error {
	dir := filepath.Join("testdata", expectedChartDir)

	// create a file system using the test golden folder as root directory.
	expectedFs := afero.NewBasePathFs(afero.NewOsFs(), dir)
	// create a file system using the test result folder as root directory.
	resultFs := afero.NewBasePathFs(config.Fs, config.ChartDir)

	err := afero.Walk(expectedFs, "", func(path string, info os.FileInfo, err error) error {
		// stop in case an error happens while walking files.
		if err != nil {
			return microerror.Mask(err)
		}

		isDir, err := afero.IsDir(expectedFs, path)
		if err != nil {
			return microerror.Mask(err)
		}
		// continue walking down directories.
		if isDir {
			return nil
		}

		expectedData, err := afero.ReadFile(expectedFs, path)
		if err != nil {
			return microerror.Mask(err)
		}

		data, err := afero.ReadFile(resultFs, path)
		if err != nil {
			return microerror.Mask(err)
		}

		if !bytes.Equal(expectedData, data) {
			return microerror.Maskf(invalidConfigError, fmt.Sprintf("%v, diff\n%s\n", path, cmp.Diff(string(expectedData), string(data))))
		}

		return nil
	})
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
