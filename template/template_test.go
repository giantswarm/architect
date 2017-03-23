package template

import (
	"fmt"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/configuration/giantswarm"
	"github.com/giantswarm/architect/configuration/giantswarm/api"
	"github.com/spf13/afero"
)

// TestTemplateKubernetesResources tests the TemplateKubernetesResources function.
func TestTemplateKubernetesResources(t *testing.T) {
	tests := []struct {
		resourcesPath string
		config        configuration.Installation
		setUp         func(afero.Fs, string) error
		check         func(afero.Fs, string) error
	}{
		// Test that an empty resources dir and config
		// produces an empty directory.
		{
			resourcesPath: "/kubernetes",
			config:        configuration.Installation{},
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
			config: configuration.Installation{
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
			setUp: func(fs afero.Fs, resourcesPath string) error {
				if err := fs.Mkdir(resourcesPath, permission); err != nil {
					return err
				}

				ingressPath := filepath.Join(resourcesPath, "ingress.yml")
				if err := afero.WriteFile(
					fs,
					ingressPath,
					[]byte("{{ .V1.GiantSwarm.API.Address.Host }}"),
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
