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
	fs       afero.Fs
	chartDir string
	sha      string
	version  string
}

// Run templates the chart's Chart.yaml and templates/deployment.yaml.
func (t TemplateHelmChartTask) Run() error {
	err := afero.Walk(t.fs, t.chartDir, func(path string, info os.FileInfo, err error) error {
		if info != nil && strings.HasSuffix(info.Name(), ".tgz") {
			return nil
		}

		contents, err := afero.ReadFile(t.fs, path)
		if err != nil {
			microerror.Mask(err)
		}

		buildInfo := BuildInfo{
			SHA:     t.sha,
			Version: t.version,
		}

		var tmpl *template.Template
		{
			// Taken from https://github.com/Masterminds/sprig. If
			// we need more common functions we should use that
			// library.
			replaceFunc := func(old, new, src string) string {
				return strings.Replace(src, old, new, -1)
			}

			tmpl = template.New(path)
			tmpl = tmpl.Delims("[[", "]]")
			tmpl = tmpl.Funcs(template.FuncMap{
				"replace": replaceFunc,
			})
			tmpl, err = tmpl.Parse(string(contents))
			if err != nil {
				microerror.Mask(err)
			}
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, buildInfo); err != nil {
			microerror.Mask(err)
		}

		if err := afero.WriteFile(t.fs, path, buf.Bytes(), permission); err != nil {
			microerror.Mask(err)
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
