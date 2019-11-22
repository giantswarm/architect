package helmtemplate

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
)

const (
	// TemplateHelmChartTaskName is the architect task name for templating helm chart
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
	// HelmValuesYamlName is hte name fo the file that holds default Helm chart
	// values inside the template directory.
	HelmValuesYamlName = "values.yaml"
)

// TemplateHelmChartTask is used to run a template-helm-chart command
type TemplateHelmChartTask struct {
	fs afero.Fs

	chartDir string
	sha      string
	version  string
}

// Config holds configuration for building a new TemplateHelmChartTask
type Config struct {
	Fs afero.Fs

	ChartDir string
	Sha      string
	Version  string
}

// Run templates the chart's Chart.yaml and templates/deployment.yaml.
func (t TemplateHelmChartTask) Run() error {

	for _, file := range []string{HelmChartYamlName, HelmValuesYamlName} {
		path := path.Join(t.chartDir, file)
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
	}
	return nil
}

// Name return the name of this task
func (t TemplateHelmChartTask) Name() string {
	return TemplateHelmChartTaskName
}

func (t TemplateHelmChartTask) String() string {
	return fmt.Sprintf(TemplateHelmChartTaskString, t.Name(), t.chartDir, t.sha, t.version)
}

// NewTemplateHelmChartTask creates a new TemplateHelmChartTask
func NewTemplateHelmChartTask(config Config) (*TemplateHelmChartTask, error) {
	if config.Fs == nil {
		return nil, microerror.Maskf(incorrectValueError, "%T.Fs must not be empty", config)
	}

	if config.ChartDir == "" {
		return nil, microerror.Maskf(incorrectValueError, "%T.ChartDir must not be empty", config)
	}

	if config.Sha == "" {
		return nil, microerror.Maskf(incorrectValueError, "%T.Sha must not be empty", config)
	}

	if config.Version == "" {
		return nil, microerror.Maskf(incorrectValueError, "%T.Version must not be empty", config)
	}

	return &TemplateHelmChartTask{
		fs:       config.Fs,
		chartDir: config.ChartDir,
		sha:      config.Sha,
		version:  config.Version,
	}, nil
}
