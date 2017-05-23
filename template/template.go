// Package template provides functions for templating Kubernetes resources,
// with standard G8S configuration.
package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"

	"github.com/giantswarm/architect/configuration"
)

var (
	// permission is a default permission to use for templated files
	permission os.FileMode = 0644
)

// BuildInfo holds information concerning the current build.
type BuildInfo struct {
	// SHA is the SHA-1 tag of the commit we are building for.
	SHA string
}

// TemplateConfiguration holds both build info, and configuration info.
type TemplateConfiguration struct {
	// BuildInfo is the configuration for the current build
	BuildInfo

	// Installation is the configuration for the installation
	configuration.Installation
}

// TemplateHelmChart takea a filesystem, a path to a directory containing a
// helm chart, and a BuildInfo struct.
// It templates the chart's Chart.yaml and templates/deployment.yaml
// with this information.
// It is an error if there are multiple charts in the helm directory.
func TemplateHelmChart(fs afero.Fs, helmPath string, buildInfo BuildInfo) error {
	fileInfos, err := afero.ReadDir(fs, helmPath)
	if err != nil {
		return err
	}

	if len(fileInfos) == 0 {
		return nil
	}

	if len(fileInfos) > 1 {
		return fmt.Errorf("Multiple charts found in helm path")
	}

	chartDirectory := fileInfos[0].Name()

	chartsYamlPath := filepath.Join(helmPath, chartDirectory, "Chart.yaml")
	deploymentPath := filepath.Join(helmPath, chartDirectory, "templates", "deployment.yaml")

	paths := []string{chartsYamlPath, deploymentPath}

	for _, path := range paths {
		exists, err := afero.Exists(fs, path)
		if err != nil {
			return err
		}

		if exists {
			contents, err := afero.ReadFile(fs, path)
			if err != nil {
				return err
			}

			t := template.Must(template.New(path).Funcs(filters).Parse(string(contents)))
			if err != nil {
				return err
			}

			var buf bytes.Buffer
			if err := t.Execute(&buf, buildInfo); err != nil {
				return err
			}

			templatedContents := buf.String()
			if err := afero.WriteFile(fs, path, []byte(templatedContents), permission); err != nil {
				return err
			}
		}
	}

	return nil
}

// TemplateKubernetesResources takes a filesystem,
// a path to a directory holding kubernetes resources,
// and an installation configuration.
// It templates the given resources, with data from the configuration,
// writing changes to the files.
func TemplateKubernetesResources(fs afero.Fs, resourcesPath string, config TemplateConfiguration) error {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		contents, err := afero.ReadFile(fs, path)
		if err != nil {
			return err
		}

		t := template.Must(template.New(path).Funcs(filters).Parse(string(contents)))
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, config); err != nil {
			return err
		}

		templatedContents := buf.String()

		// This add backwards compatability for `%%DOCKER_TAG%%`. Deprecated.
		templatedContents = strings.Replace(templatedContents, "%%DOCKER_TAG%%", config.BuildInfo.SHA, -1)

		if err := afero.WriteFile(fs, path, []byte(templatedContents), permission); err != nil {
			return err
		}

		return nil
	}

	if err := afero.Walk(fs, resourcesPath, walkFunc); err != nil {
		return err
	}

	return nil
}
