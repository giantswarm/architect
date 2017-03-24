// Package template provides functions for templating Kubernetes resources,
// with standard G8S configuration.
package template

import (
	"bytes"
	"os"
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

// TemplateKubernetesResources takes a filesystem,
// a path to a directory holding kubernetes resources,
// and an installation configuration.
// It templates the given resources, with data from the configuration,
// writing changes to the files.
func TemplateKubernetesResources(fs afero.Fs, resourcesPath string, config TemplateConfiguration) error {
	funcMap := template.FuncMap{
		"ShortDuration": shortDuration,
		"urlString":     urlString,
	}

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

		t := template.Must(template.New(path).Funcs(funcMap).Parse(string(contents)))
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
