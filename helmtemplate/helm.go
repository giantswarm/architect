package helmtemplate

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

const (
	// HelmChartYamlName is the name of Helm's chart yaml.
	HelmChartYamlName = "Chart.yaml"
	// HelmValuesYamlName is hte name fo the file that holds default Helm chart
	// values inside the template directory.
	HelmValuesYamlName = "values.yaml"
	// HelmTemplateDirectoryName is the name of the directory that stores
	// Kubernetes resources inside a chart.
	HelmTemplateDirectoryName = "templates"
)

// TemplateHelmChartTask is used to run a template-helm-chart command
type TemplateHelmChartTask struct {
	fs afero.Fs

	chartDir     string
	branch       string
	sha          string
	chartVersion string
	appVersion   string
}

// Config holds configuration for building a new TemplateHelmChartTask
type Config struct {
	Fs afero.Fs

	ChartDir   string
	Branch     string
	Sha        string
	Version    string
	AppVersion string
}

// renderedChart is used for chart validation after it has been filled with
// values
type renderedChart struct {
	Version    string `json:"version,omitempty"`
	AppVersion string `json:"appVersion,omitempty"`
}

// Run templates the chart's Chart.yaml and templates/deployment.yaml.
func (t TemplateHelmChartTask) Run(validate, taggedBuild bool) error {
	for _, file := range []string{HelmChartYamlName, HelmValuesYamlName} {
		path := path.Join(t.chartDir, file)
		contents, err := afero.ReadFile(t.fs, path)
		if err != nil {
			return microerror.Mask(err)
		}

		// We expect versions to match for a tagged build if project.go has been found.
		if validate && taggedBuild && t.appVersion != "" && t.chartVersion != t.appVersion {
			return microerror.Maskf(
				versionMismatchError,
				"version in git tag must be equal to version in pkg/project/project.go: %q != %q",
				t.chartVersion, t.appVersion,
			)
		}

		buildInfo := BuildInfo{
			Branch:     t.branch,
			SHA:        t.sha,
			Version:    t.chartVersion,
			AppVersion: t.appVersion,
		}

		tmpl, err := template.New(path).Delims("[[", "]]").Parse(string(contents))
		if err != nil {
			return microerror.Mask(err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, buildInfo); err != nil {
			return microerror.Mask(err)
		}

		if validate && ValidateChart(buildInfo.Version, buildInfo.AppVersion, buf) != nil {
			return microerror.Mask(err)
		}

		if err := afero.WriteFile(t.fs, path, buf.Bytes(), permission); err != nil {
			return microerror.Mask(err)
		}
	}
	return nil
}

// ValidateChart makes sure version fields and values made it into the
// chart when executing the template.
func ValidateChart(version, appVersion string, chartBuf bytes.Buffer) error {
	chart := renderedChart{}
	err := yaml.Unmarshal(chartBuf.Bytes(), &chart)
	if err != nil {
		return microerror.Mask(err)
	}

	if chart.Version != version {
		return microerror.Maskf(
			versionValidationError,
			"wrong value for \"version\" in chart: got %#q, expected %#q",
			chart.Version, version,
		)
	}

	// We want to validate appVersion only when the project.go has been found,
	// i.e. appVersion is non-empty.
	if appVersion != "" && chart.AppVersion != appVersion {
		return microerror.Maskf(
			versionValidationError,
			"wrong value for \"appVersion\" in chart: got %#q, expected %#q",
			chart.AppVersion, appVersion,
		)
	}

	return nil
}

func (t TemplateHelmChartTask) String() string {
	return fmt.Sprintf("%s:\t%s sha:%s chartVersion:%s appVersion:%s", "template-helm-chart", t.chartDir, t.sha, t.chartVersion, t.appVersion)
}

// NewTemplateHelmChartTask creates a new TemplateHelmChartTask
func NewTemplateHelmChartTask(config Config) (*TemplateHelmChartTask, error) {
	if config.Fs == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Fs must not be empty", config)
	}

	if config.ChartDir == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.ChartDir must not be empty", config)
	}

	if config.Branch == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Branch must not be empty", config)
	}

	if config.Sha == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Sha must not be empty", config)
	}

	if config.Version == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Version must not be empty", config)
	}

	t := &TemplateHelmChartTask{
		fs:           config.Fs,
		chartDir:     config.ChartDir,
		branch:       config.Branch,
		sha:          config.Sha,
		chartVersion: config.Version,
		appVersion:   config.AppVersion,
	}

	return t, nil
}
