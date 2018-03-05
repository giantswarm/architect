// Package pipeline provides functionality for moving artifacts between
// different stages in the deployment process.
package pipeline

import (
	"path/filepath"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// StartChannel determines the initial channel to push a chart when the
// deployment process starts. It has inforamtion about the chart version
// and the initial stability defined for charts.
func StartChannel(fs afero.Fs, workingDirectory, project string) (string, error) {
	c := &struct {
		Version string `yaml:"version"`
	}{}

	chartFile := filepath.Join(workingDirectory, "helm", project+"-chart", "Chart.yaml")

	y, err := afero.ReadFile(fs, chartFile)
	if err != nil {
		return "", microerror.Mask(err)
	}

	err = yaml.Unmarshal(y, c)
	if err != nil {
		return "", microerror.Mask(err)
	}

	elements := strings.Split(c.Version, ".")
	if len(elements) < 2 {
		return "", microerror.Mask(incorrectChartVersionError)
	}
	// remove initial v
	elements[0] = strings.TrimLeft(elements[0], "v")

	items := append(elements[0:2], initialStability)

	startChannel := strings.Join(items, "-")

	return startChannel, nil
}
