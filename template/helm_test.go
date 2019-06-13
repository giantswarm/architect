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
		chartDir  string
		dockerTag string
		sha       string
		version   string
		setUp     func(afero.Fs, string) error
		check     func(afero.Fs, string) error
	}{
		// Test that a chart is templated correctly.
		{
			chartDir:  "/helm/test-chart",
			dockerTag: "white-rabbit",
			sha:       "jabberwocky",
			version:   "mad-hatter",
			setUp: func(fs afero.Fs, chartDir string) error {
				files := []struct {
					path string
					data string
				}{
					{
						path: filepath.Join(chartDir, HelmChartYamlName),
						data: "version: 1.0.0-[[ .SHA ]]",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, HelmDeploymentYamlName),
						data: "image: [[ .SHA ]] foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "daemonset.yaml"),
						data: "image: [[ .SHA ]] foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "otherfile.yaml"),
						data: "image: [[ .SHA ]] foo: < abc",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "subdirectory", "replicaset.yaml"),
						data: "image: [[ .SHA ]] foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "ingress.yaml"),
						data: "host: {{ .Values.Installation.etc }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "with-dockerTag.yaml"),
						data: "image: [[ .DockerTag ]]",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "with-version.yaml"),
						data: "version: [[ .Version ]]",
					},
					{
						path: filepath.Join(chartDir, "charts", "foo.tgz"),
						data: "version: [[ .SHA ]] - not templated as .tgz files are ignored",
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
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, HelmDeploymentYamlName),
						data: "image: jabberwocky foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "daemonset.yaml"),
						data: "image: jabberwocky foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "otherfile.yaml"),
						data: "image: jabberwocky foo: < abc",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "subdirectory", "replicaset.yaml"),
						data: "image: jabberwocky foo: {{ .Values.Foo }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "ingress.yaml"),
						data: "host: {{ .Values.Installation.etc }}",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "with-dockerTag.yaml"),
						data: "image: white-rabbit",
					},
					{
						path: filepath.Join(chartDir, HelmTemplateDirectoryName, "with-version.yaml"),
						data: "version: mad-hatter",
					},
					{
						path: filepath.Join(chartDir, "charts", "foo.tgz"),
						data: "version: [[ .SHA ]]",
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
		task := NewTemplateHelmChartTask(fs, test.chartDir, test.dockerTag, test.sha, test.version)

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
