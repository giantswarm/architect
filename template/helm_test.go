package template

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

// TestTemplateHelmChartTask tests the TemplateHelmChartTask.
func TestTemplateHelmChartTask(t *testing.T) {
	tests := []struct {
		helmPath string
		sha      string
		setUp    func(afero.Fs, string) error
		check    func(afero.Fs, string) error
	}{
		// Test that an empty helm directory does nothing.
		{
			helmPath: "/helm",
			sha:      "jabberwocky",
			setUp: func(fs afero.Fs, helmPath string) error {
				if err := fs.Mkdir(helmPath, permission); err != nil {
					return microerror.Mask(err)
				}

				return nil
			},
			check: func(fs afero.Fs, helmPath string) error {
				fileInfos, err := afero.ReadDir(fs, helmPath)
				if err != nil {
					return microerror.Mask(err)
				}
				if len(fileInfos) != 0 {
					return microerror.Mask(multipleHelmChartsError)
				}

				return nil
			},
		},

		// Test that a chart is templated correctly.
		{
			helmPath: "/helm",
			sha:      "jabberwocky",
			setUp: func(fs afero.Fs, helmPath string) error {
				directories := []string{
					helmPath,
					filepath.Join(helmPath, "test-chart"),
					filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName),
				}

				for _, directory := range directories {
					if err := fs.Mkdir(directory, permission); err != nil {
						return microerror.Mask(err)
					}
				}

				files := []struct {
					path string
					data string
				}{
					{
						path: filepath.Join(helmPath, "test-chart", HelmChartYamlName),
						data: "version: 1.0.0-[[ .SHA ]]",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, HelmDeploymentYamlName),
						data: "image: [[ .SHA ]] foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "daemonset.yaml"),
						data: "image: [[ .SHA ]] foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "otherfile.yaml"),
						data: "image: [[ .SHA ]] foo: < abc",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "subdirectory", "replicaset.yaml"),
						data: "image: [[ .SHA ]] foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "ingress.yaml"),
						data: "host: {{ .Values.Installation.etc }}",
					},
				}

				for _, file := range files {
					if err := afero.WriteFile(fs, file.path, []byte(file.data), permission); err != nil {
						return microerror.Mask(err)
					}
				}

				return nil
			},
			check: func(fs afero.Fs, helmPath string) error {
				files := []struct {
					path string
					data string
				}{
					{
						path: filepath.Join(helmPath, "test-chart", HelmChartYamlName),
						data: "version: 1.0.0-jabberwocky",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, HelmDeploymentYamlName),
						data: "image: jabberwocky foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "daemonset.yaml"),
						data: "image: jabberwocky foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "otherfile.yaml"),
						data: "image: jabberwocky foo: < abc",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "subdirectory", "replicaset.yaml"),
						data: "image: jabberwocky foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "ingress.yaml"),
						data: "host: {{ .Values.Installation.etc }}",
					},
				}

				for _, file := range files {
					bytes, err := afero.ReadFile(fs, file.path)
					if err != nil {
						return microerror.Mask(err)
					}
					if string(bytes) != file.data {
						return microerror.Maskf(incorrectValueError, fmt.Sprintf("%v, found: %v, expected: %v", file.path, string(bytes), file.data))
					}
				}

				return nil
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()
		task := NewTemplateHelmChartTask(fs, test.helmPath, test.sha)

		if err := test.setUp(fs, test.helmPath); err != nil {
			t.Fatalf("%v: unexpected error during setup: %v\n", index, err)
		}

		if err := task.Run(); err != nil {
			t.Fatalf("%v: unexpected error during templating: %v\n", index, err)
		}

		if err := test.check(fs, test.helmPath); err != nil {
			t.Fatalf("%v: unexpected error during check: %v\n", index, err)
		}
	}
}
