package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/afero"

	"github.com/giantswarm/microerror"
)

const (
	TemplateHelmChartTaskName = "template-helm-chart"

	// TemplateHelmChartTaskString is the format for printing the
	// helm chart templating task.
	// Name of the task, the helm directory path, and the sha.
	TemplateHelmChartTaskString = "%s:\t%s sha:%s ref:%s"

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
	ref      string
	sha      string
}

// Run templates the chart's Chart.yaml and templates/deployment.yaml.
func (t TemplateHelmChartTask) Run() error {
	err := afero.Walk(t.fs, t.chartDir, func(path string, info os.FileInfo, err error) error {
		contents, err := afero.ReadFile(t.fs, path)
		if err != nil {
			microerror.Mask(err)
		}

		buildInfo := BuildInfo{
			Ref: t.ref,
			SHA: t.sha,
		}

		newTemplate := template.Must(template.New(path).Delims("[[", "]]").Parse(string(contents)))
		if err != nil {
			microerror.Mask(err)
		}

		var buf bytes.Buffer
		if err := newTemplate.Execute(&buf, buildInfo); err != nil {
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
	return fmt.Sprintf(TemplateHelmChartTaskString, t.Name(), t.chartDir, t.sha, t.ref)
}

func NewTemplateHelmChartTask(fs afero.Fs, chartDir, ref, sha string) TemplateHelmChartTask {
	return TemplateHelmChartTask{
		fs:       fs,
		chartDir: chartDir,
		ref:      ref,
		sha:      sha,
	}
}
