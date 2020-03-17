package helmtemplate

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

// TestTemplateHelmChartTask tests the TemplateHelmChartTask.
func TestTemplateHelmChartTask(t *testing.T) {
	tests := []struct {
		chartDir     string
		branch       string
		sha          string
		chartVersion string
		appVersion   string
		setUp        func(afero.Fs, string) error
		check        func(afero.Fs, string) error
	}{
		// Test that a chart is templated correctly.
		{
			chartDir:     "/helm/test-chart",
			branch:       "beamish-boy",
			sha:          "jabberwocky",
			chartVersion: "mad-hatter",
			appVersion:   "1.0.0",
			setUp: func(fs afero.Fs, chartDir string) error {
				files := []struct {
					path string
					data string
				}{
					{
						path: filepath.Join(chartDir, HelmChartYamlName),
						data: "version: [[ .AppVersion ]]-[[ .SHA ]]",
					},
					{
						path: filepath.Join(chartDir, HelmValuesYamlName),
						data: "Branch: [[ .Branch ]]\nVersion: [[ .Version ]]",
					},
				}

				for _, file := range files {
					dir := filepath.Base(file.path)
					if dir != "." && dir != "/" {
						if err := fs.MkdirAll(dir, permission); err != nil {
							return microerror.Mask(err)
						}
					}
					if err := afero.WriteFile(fs, file.path, []byte(file.data), permission); err != nil {
						return microerror.Mask(err)
					}
				}

				return nil
			},
			check: func(fs afero.Fs, chartDir string) error {
				files := []struct {
					path string
					data string
				}{
					{
						path: filepath.Join(chartDir, HelmChartYamlName),
						data: "version: 1.0.0-jabberwocky",
					},
					{
						path: filepath.Join(chartDir, HelmValuesYamlName),
						data: "Branch: beamish-boy\nVersion: mad-hatter",
					},
				}

				for _, file := range files {
					bytes, err := afero.ReadFile(fs, file.path)
					if err != nil {
						return microerror.Mask(err)
					}
					if string(bytes) != file.data {
						return microerror.Maskf(invalidConfigError, fmt.Sprintf("%v, found: %v, expected: %v", file.path, string(bytes), file.data))
					}
				}

				return nil
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		task, err := NewTemplateHelmChartTask(Config{
			Fs:         fs,
			ChartDir:   test.chartDir,
			Branch:     test.branch,
			Sha:        test.sha,
			Version:    test.chartVersion,
			AppVersion: test.appVersion,
		})

		if err != nil {
			t.Fatalf("%v: unexpected error when creating NewTemplateHelmChartTask: %v\n", index, err)
		}

		if err := test.setUp(fs, test.chartDir); err != nil {
			t.Fatalf("%v: unexpected error during setup: %v\n", index, err)
		}

		if err := task.Run(); err != nil {
			t.Fatalf("%v: unexpected error during templating: %v\n", index, err)
		}

		if err := test.check(fs, test.chartDir); err != nil {
			t.Fatalf("%v: unexpected error during check: %v\n", index, err)
		}
	}
}
