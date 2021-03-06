package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

const (
	TemplateHelmChartTaskName = "template-helm-chart"

	// TemplateHelmChartTaskString is the format for printing the
	// helm chart templating task.
	// Name of the task, the helm directory path, the sha, and the version.
	TemplateHelmChartTaskString = "%s:\t%s sha:%s version:%s"

	// HelmChartYamlName is the name of Helm's chart yaml.
	HelmChartYamlName = "Chart.yaml"
	// HelmTemplateDirectoryName is the name of the directory that stores
	// Kubernetes resources inside a chart.
	HelmTemplateDirectoryName = "templates"
	// HelmDeploymentYamlName is the name of the file we template inside the
	// Helm template directory.
	HelmDeploymentYamlName = "deployment.yaml"
)

type TemplateHelmChartTask struct {
	fs afero.Fs

	chartDir string
	sha      string
	version  string
}

// Run templates the chart's Chart.yaml and templates/deployment.yaml.
func (t TemplateHelmChartTask) Run() error {
	err := afero.Walk(t.fs, t.chartDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return microerror.Mask(err)
		}

		if info != nil && info.IsDir() {
			return nil
		}
		if info != nil && strings.HasSuffix(info.Name(), ".tgz") {
			return nil
		}

		contents, err := afero.ReadFile(t.fs, path)
		if err != nil {
			return microerror.Mask(err)
		}

		buildInfo := BuildInfo{
			SHA:     t.sha,
			Version: t.version,
		}

		tmpl, err := template.New(path).Delims("[[", "]]").Parse(string(contents))
		if err != nil {
			return microerror.Mask(err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, buildInfo); err != nil {
			return microerror.Mask(err)
		}

		if err := afero.WriteFile(t.fs, path, buf.Bytes(), permission); err != nil {
			return microerror.Mask(err)
		}
		return nil
	})

	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (t TemplateHelmChartTask) Name() string {
	return TemplateHelmChartTaskName
}

func (t TemplateHelmChartTask) String() string {
	return fmt.Sprintf(TemplateHelmChartTaskString, t.Name(), t.chartDir, t.sha, t.version)
}

func NewTemplateHelmChartTask(fs afero.Fs, chartDir, sha, version string) TemplateHelmChartTask {
	return TemplateHelmChartTask{
		fs:       fs,
		chartDir: chartDir,
		sha:      sha,
		version:  version,
	}
}
