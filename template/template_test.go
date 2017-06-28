package template

import (
	"fmt"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/configuration/giantswarm"
	"github.com/giantswarm/architect/configuration/giantswarm/api"
	"github.com/giantswarm/architect/configuration/monitoring"
	"github.com/giantswarm/architect/configuration/monitoring/prometheus"
	"github.com/giantswarm/architect/configuration/monitoring/testbot"
	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/afero"
)

// TestTemplateHelmChart tests the TestTemplateHelmChart function.
func TestTemplateHelmChart(t *testing.T) {
	tests := []struct {
		helmPath  string
		buildInfo BuildInfo
		setUp     func(afero.Fs, string) error
		check     func(afero.Fs, string) error
	}{
		// Test that an empty helm directory does nothing.
		{
			helmPath: "/helm",
			buildInfo: BuildInfo{
				SHA: "jabberwocky",
			},
			setUp: func(fs afero.Fs, helmPath string) error {
				if err := fs.Mkdir(helmPath, permission); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			check: func(fs afero.Fs, helmPath string) error {
				fileInfos, err := afero.ReadDir(fs, helmPath)
				if err != nil {
					return microerror.MaskAny(err)
				}
				if len(fileInfos) != 0 {
					return microerror.MaskAny(multipleHelmChartsError)
				}

				return nil
			},
		},

		// Test that a chart is templated correctly.
		{
			helmPath: "/helm",
			buildInfo: BuildInfo{
				SHA: "jabberwocky",
			},
			setUp: func(fs afero.Fs, helmPath string) error {
				directories := []string{
					helmPath,
					filepath.Join(helmPath, "test-chart"),
					filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName),
				}

				for _, directory := range directories {
					if err := fs.Mkdir(directory, permission); err != nil {
						return microerror.MaskAny(err)
					}
				}

				files := []struct {
					path string
					data string
				}{
					{
						path: filepath.Join(helmPath, "test-chart", HelmChartYamlName),
						data: "version: 1.0.0-{{ .SHA }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, HelmDeploymentYamlName),
						data: "image: {{ .SHA }}",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "ingress.yaml"),
						data: "host: {{ .Values.Installation.etc }}",
					},
				}

				for _, file := range files {
					if err := afero.WriteFile(fs, file.path, []byte(file.data), permission); err != nil {
						return microerror.MaskAny(err)
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
						data: "image: jabberwocky",
					},
					{
						path: filepath.Join(helmPath, "test-chart", HelmTemplateDirectoryName, "ingress.yaml"),
						data: "host: {{ .Values.Installation.etc }}",
					},
				}

				for _, file := range files {
					bytes, err := afero.ReadFile(fs, file.path)
					if err != nil {
						return microerror.MaskAny(err)
					}
					if string(bytes) != file.data {
						return microerror.MaskAnyf(incorrectValueError, fmt.Sprintf("%v, found: %v, expected: %v", file.path, string(bytes), file.data))
					}
				}

				return nil
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()

		if err := test.setUp(fs, test.helmPath); err != nil {
			t.Fatalf("%v: unexpected error during setup: %v\n", index, err)
		}

		if err := TemplateHelmChart(fs, test.helmPath, test.buildInfo, nil, ""); err != nil {
			t.Fatalf("%v: unexpected error during templating: %v\n", index, err)
		}

		if err := test.check(fs, test.helmPath); err != nil {
			t.Fatalf("%v: unexpected error during check: %v\n", index, err)
		}
	}
}

// TestTemplateKubernetesResources tests the TemplateKubernetesResources function.
func TestTemplateKubernetesResources(t *testing.T) {
	tests := []struct {
		resourcesPath string
		config        TemplateConfiguration
		setUp         func(afero.Fs, string) error
		check         func(afero.Fs, string) error
	}{
		// Test that an empty resources dir and config
		// produces an empty directory.
		{
			resourcesPath: "/kubernetes",
			config:        TemplateConfiguration{},
			setUp: func(fs afero.Fs, resourcesPath string) error {
				if err := fs.Mkdir(resourcesPath, permission); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return microerror.MaskAny(err)
				}
				if len(fileInfos) != 0 {
					return multipleHelmChartsError
				}

				return nil
			},
		},

		// Test that a resources directory with an api ingress is templated correctly.
		{
			resourcesPath: "/kubernetes",
			config: TemplateConfiguration{
				Installation: configuration.Installation{
					V1: configuration.V1{
						GiantSwarm: giantswarm.GiantSwarm{
							API: api.API{
								Address: url.URL{
									Scheme: "https",
									Host:   "api-g8s.giantswarm.io",
								},
							},
						},
					},
				},
			},
			setUp: func(fs afero.Fs, resourcesPath string) error {
				if err := fs.Mkdir(resourcesPath, permission); err != nil {
					return microerror.MaskAny(err)
				}

				path := filepath.Join(resourcesPath, "ingress.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte("{{ .Installation.V1.GiantSwarm.API.Address.Host }}"),
					permission,
				); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if len(fileInfos) != 1 {
					return multipleFilesFoundInResourcesError
				}

				if fileInfos[0].Name() != "ingress.yml" {
					return microerror.MaskAnyf(resourceNotFoundError, "ingress.yml")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "ingress.yml"))
				if err != nil {
					return microerror.MaskAny(err)
				}

				if string(bytes) != "api-g8s.giantswarm.io" {
					return microerror.MaskAnyf(incorrectTemplatingError, "ingress address not found, found: '%v'", string(bytes))
				}

				return nil
			},
		},

		// Test that a resources directory with a prometheus address is templated correctly.
		{
			resourcesPath: "/kubernetes",
			config: TemplateConfiguration{
				Installation: configuration.Installation{
					V1: configuration.V1{
						Monitoring: monitoring.Monitoring{
							Prometheus: prometheus.Prometheus{
								Address: url.URL{
									Scheme: "https",
									Host:   "prometheus-g8s.giantswarm.io",
								},
							},
						},
					},
				},
			},
			setUp: func(fs afero.Fs, resourcesPath string) error {
				if err := fs.Mkdir(resourcesPath, permission); err != nil {
					return microerror.MaskAny(err)
				}

				path := filepath.Join(resourcesPath, "ingress.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte(`{{ .Installation.V1.Monitoring.Prometheus.Address | urlString }}`),
					permission,
				); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if len(fileInfos) != 1 {
					return multipleFilesFoundInResourcesError
				}

				if fileInfos[0].Name() != "ingress.yml" {
					return microerror.MaskAnyf(resourceNotFoundError, "ingress.yml")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "ingress.yml"))
				if err != nil {
					return microerror.MaskAny(err)
				}

				if string(bytes) != "https://prometheus-g8s.giantswarm.io" {
					return microerror.MaskAnyf(incorrectTemplatingError, "ingress address not found, found: '%v'", string(bytes))
				}

				return nil
			},
		},

		// Test that a resources directory with a testbot interval is templated correctly.
		{
			resourcesPath: "/kubernetes",
			config: TemplateConfiguration{
				Installation: configuration.Installation{
					V1: configuration.V1{
						Monitoring: monitoring.Monitoring{
							Testbot: testbot.Testbot{
								Interval: 5 * time.Minute,
							},
						},
					},
				},
			},
			setUp: func(fs afero.Fs, resourcesPath string) error {
				if err := fs.Mkdir(resourcesPath, permission); err != nil {
					return microerror.MaskAny(err)
				}

				path := filepath.Join(resourcesPath, "configmap.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte("interval: '@every {{ .Installation.V1.Monitoring.Testbot.Interval | shortDuration }}'"),
					permission,
				); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if len(fileInfos) != 1 {
					return multipleFilesFoundInResourcesError
				}

				if fileInfos[0].Name() != "configmap.yml" {
					return microerror.MaskAnyf(resourceNotFoundError, "configmap.yml")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "configmap.yml"))
				if err != nil {
					return microerror.MaskAny(err)
				}

				if string(bytes) != "interval: '@every 5m'" {
					return microerror.MaskAnyf(incorrectTemplatingError, "correct interval string not found, found: '%v'", string(bytes))
				}

				return nil
			},
		},

		// Test that a resources directory with a docker tag is templated correctly.
		{
			resourcesPath: "/kubernetes",
			config: TemplateConfiguration{
				BuildInfo: BuildInfo{
					SHA: "12345",
				},
				Installation: configuration.Installation{},
			},
			setUp: func(fs afero.Fs, resourcesPath string) error {
				if err := fs.Mkdir(resourcesPath, permission); err != nil {
					return microerror.MaskAny(err)
				}

				path := filepath.Join(resourcesPath, "deployment.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte("image: foo/bar:{{ .BuildInfo.SHA }}"),
					permission,
				); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if len(fileInfos) != 1 {
					return multipleFilesFoundInResourcesError
				}

				if fileInfos[0].Name() != "deployment.yml" {
					return microerror.MaskAnyf(resourceNotFoundError, "deployment.yml")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "deployment.yml"))
				if err != nil {
					return microerror.MaskAny(err)
				}

				if string(bytes) != "image: foo/bar:12345" {
					return microerror.MaskAnyf(incorrectShaError, "in deployment, found: '%v'", string(bytes))
				}

				return nil
			},
		},

		// Test that the older docker tag format is templated correctly.
		{
			resourcesPath: "/kubernetes",
			config: TemplateConfiguration{
				BuildInfo: BuildInfo{
					SHA: "12345",
				},
				Installation: configuration.Installation{},
			},
			setUp: func(fs afero.Fs, resourcesPath string) error {
				if err := fs.Mkdir(resourcesPath, permission); err != nil {
					return microerror.MaskAny(err)
				}

				path := filepath.Join(resourcesPath, "deployment.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte("image: registry.giantswarm.io/giantswarm/api:%%DOCKER_TAG%%"),
					permission,
				); err != nil {
					return microerror.MaskAny(err)
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if len(fileInfos) != 1 {
					return multipleFilesFoundInResourcesError
				}

				if fileInfos[0].Name() != "deployment.yml" {
					return microerror.MaskAnyf(resourceNotFoundError, "deployment.yml")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "deployment.yml"))
				if err != nil {
					return microerror.MaskAny(err)
				}

				if string(bytes) != "image: registry.giantswarm.io/giantswarm/api:12345" {
					return microerror.MaskAnyf(incorrectShaError, "in deployment, found: '%v'", string(bytes))
				}

				return nil
			},
		},
	}

	for index, test := range tests {
		fs := afero.NewMemMapFs()

		if err := test.setUp(fs, test.resourcesPath); err != nil {
			t.Fatalf("%v: unexpected error during setup: %v\n", index, err)
		}

		if err := TemplateKubernetesResources(fs, test.resourcesPath, test.config, nil, ""); err != nil {
			t.Fatalf("%v: unexpected error during templating: %v\n", index, err)
		}

		if err := test.check(fs, test.resourcesPath); err != nil {
			t.Fatalf("%v: unexpected error during check: %v\n", index, err)
		}
	}
}
