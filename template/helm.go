package template

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/spf13/afero"

	microerror "github.com/giantswarm/microkit/error"
)

const (
	TemplateHelmChartTaskName = "template-helm-chart"

	// TemplateHelmChartTaskString is the format for printing the
	// helm chart templating task.
	// Name of the task, the helm directory path, and the sha.
	TemplateHelmChartTaskString = "%s:\t%s %s"

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
	helmPath string
	sha      string
}

// Run templates the chart's Chart.yaml and templates/deployment.yaml.
// It is an error if there are multiple charts in the helm directory.
func (t TemplateHelmChartTask) Run() error {
	fileInfos, err := afero.ReadDir(t.fs, t.helmPath)
	if err != nil {
		return microerror.MaskAny(err)
	}

	if len(fileInfos) == 0 {
		return nil
	}

	if len(fileInfos) > 1 {
		return microerror.MaskAny(multipleHelmChartsError)
	}

	chartDirectory := fileInfos[0].Name()

	chartsYamlPath := filepath.Join(t.helmPath, chartDirectory, HelmChartYamlName)
	deploymentPath := filepath.Join(t.helmPath, chartDirectory, HelmTemplateDirectoryName, HelmDeploymentYamlName)

	paths := []string{chartsYamlPath, deploymentPath}

	for _, path := range paths {
		exists, err := afero.Exists(t.fs, path)
		if err != nil {
			microerror.MaskAny(err)
		}

		if exists {
			contents, err := afero.ReadFile(t.fs, path)
			if err != nil {
				microerror.MaskAny(err)
			}

			buildInfo := BuildInfo{SHA: t.sha}

			newTemplate := template.Must(template.New(path).Parse(string(contents)))
			if err != nil {
				microerror.MaskAny(err)
			}

			var buf bytes.Buffer
			if err := newTemplate.Execute(&buf, buildInfo); err != nil {
				microerror.MaskAny(err)
			}

			if err := afero.WriteFile(t.fs, path, buf.Bytes(), permission); err != nil {
				microerror.MaskAny(err)
			}
		}
	}

	return nil
}

func (t TemplateHelmChartTask) Name() string {
	return TemplateHelmChartTaskName
}

func (t TemplateHelmChartTask) String() string {
	return fmt.Sprintf(TemplateHelmChartTaskString, t.Name(), t.helmPath, t.sha)
}

func NewTemplateHelmChartTask(fs afero.Fs, helmPath, sha string) TemplateHelmChartTask {
	return TemplateHelmChartTask{
		fs:       fs,
		helmPath: helmPath,
		sha:      sha,
	}
}
