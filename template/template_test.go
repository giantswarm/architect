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
	"github.com/spf13/afero"
)

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
					return err
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return err
				}
				if len(fileInfos) != 0 {
					return fmt.Errorf("multiple files found in resources directory")
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
					return err
				}

				path := filepath.Join(resourcesPath, "ingress.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte("{{ .Installation.V1.GiantSwarm.API.Address.Host }}"),
					permission,
				); err != nil {
					return err
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return err
				}

				if len(fileInfos) != 1 {
					return fmt.Errorf("did not find only one file in resources path")
				}

				if fileInfos[0].Name() != "ingress.yml" {
					return fmt.Errorf("ingress not found in resources path")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "ingress.yml"))
				if err != nil {
					return err
				}

				if string(bytes) != "api-g8s.giantswarm.io" {
					return fmt.Errorf("ingress address not found, found: '%v'", string(bytes))
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
					return err
				}

				path := filepath.Join(resourcesPath, "ingress.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte(`{{ .Installation.V1.Monitoring.Prometheus.Address | urlString }}`),
					permission,
				); err != nil {
					return err
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return err
				}

				if len(fileInfos) != 1 {
					return fmt.Errorf("did not find only one file in resources path")
				}

				if fileInfos[0].Name() != "ingress.yml" {
					return fmt.Errorf("ingress not found in resources path")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "ingress.yml"))
				if err != nil {
					return err
				}

				if string(bytes) != "https://prometheus-g8s.giantswarm.io" {
					return fmt.Errorf("ingress address not found, found: '%v'", string(bytes))
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
					return err
				}

				path := filepath.Join(resourcesPath, "configmap.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte("interval: '@every {{ .Installation.V1.Monitoring.Testbot.Interval | ShortDuration }}'"),
					permission,
				); err != nil {
					return err
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return err
				}

				if len(fileInfos) != 1 {
					return fmt.Errorf("did not find only one file in resources path")
				}

				if fileInfos[0].Name() != "configmap.yml" {
					return fmt.Errorf("configmap not found in resources path")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "configmap.yml"))
				if err != nil {
					return err
				}

				if string(bytes) != "interval: '@every 5m'" {
					return fmt.Errorf("correct interval string not found, found: '%v'", string(bytes))
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
					return err
				}

				path := filepath.Join(resourcesPath, "deployment.yml")
				if err := afero.WriteFile(
					fs,
					path,
					[]byte("image: foo/bar:{{ .BuildInfo.SHA }}"),
					permission,
				); err != nil {
					return err
				}

				return nil
			},
			check: func(fs afero.Fs, resourcesPath string) error {
				fileInfos, err := afero.ReadDir(fs, resourcesPath)
				if err != nil {
					return err
				}

				if len(fileInfos) != 1 {
					return fmt.Errorf("did not find only one file in resources path")
				}

				if fileInfos[0].Name() != "deployment.yml" {
					return fmt.Errorf("deployment not found in resources path")
				}

				bytes, err := afero.ReadFile(fs, filepath.Join(resourcesPath, "deployment.yml"))
				if err != nil {
					return err
				}

				if string(bytes) != "image: foo/bar:12345" {
					return fmt.Errorf("correct sha not found, found: '%v'", string(bytes))
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

		if err := TemplateKubernetesResources(fs, test.resourcesPath, test.config); err != nil {
			t.Fatalf("%v: unexpected error during templating: %v\n", index, err)
		}

		if err := test.check(fs, test.resourcesPath); err != nil {
			t.Fatalf("%v: unexpected error during check: %v\n", index, err)
		}
	}
}
